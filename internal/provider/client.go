package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

const (
	// Time out requests after 10 seconds
	ClientTimeout int = 10
)

type RLHttpClient struct {
	client      *http.Client
	RateLimiter *rate.Limiter
}

type APIClient struct {
	BaseURL    *url.URL
	Token      string
	httpClient *RLHttpClient
	UserAgent  string
}

func (c *RLHttpClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.RateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewClient(endpoint string, token string, UserAgent string) (*APIClient, error) {

	if endpoint == "" || token == "" {
		return nil, fmt.Errorf("token and endpoint are required")
	}

	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	rl := rate.NewLimiter(rate.Every(time.Second/10), 10)

	h := &http.Client{
		Timeout: time.Duration(ClientTimeout) * time.Second,
	}

	rlClient := &RLHttpClient{
		client:      h,
		RateLimiter: rl,
	}

	c := &APIClient{
		httpClient: rlClient,
		BaseURL:    baseURL,
		Token:      token,
		UserAgent:  UserAgent,
	}

	return c, nil
}

func (c *APIClient) newRequest(method, path string, filter string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	if filter != "" {
		q := u.Query()
		q.Set("filter", filter)
		u.RawQuery = q.Encode()
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.Token))
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *APIClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == 401:
		return nil, errors.New("401 unauthorized")
	case resp.StatusCode == 404:
		return nil, errors.New("404 not found")
	case resp.StatusCode == 409:
		return nil, errors.New("409 conflict, resource already exists")
	case resp.StatusCode == 429:
		return nil, errors.New("429 ThrottlingException")
	case resp.StatusCode == 200 || resp.StatusCode == 201:
		err = json.NewDecoder(resp.Body).Decode(v)
		return resp, err
	case resp.StatusCode == 204:
		return resp, err
	case resp.StatusCode <= 299 && resp.StatusCode >= 200:
		return resp, err
	default:
		return nil, fmt.Errorf("unexpected HTTP status code: %v", resp.StatusCode)
	}
}

func (c *APIClient) ListUsers() (*[]User, error) {
	req, err := c.newRequest("GET", "Users", "", nil)
	if err != nil {
		return nil, err
	}
	var userLR UserListResponse
	_, err = c.do(req, &userLR)
	return &userLR.Resources, err
}

func (c *APIClient) FindUserByUsername(username string) (*User, error) {
	filter := fmt.Sprintf("userName eq \"%v\"", username)

	req, err := c.newRequest("GET", "Users", filter, nil)

	if err != nil {
		return nil, err
	}

	var userLR UserListResponse
	_, err = c.do(req, &userLR)
	if err != nil {
		return nil, err
	}

	if userLR.TotalResults != 1 || len(userLR.Resources) != 1 {
		return nil, fmt.Errorf("user \"%v\" not found", username)
	}

	return &userLR.Resources[0], err
}

func (c *APIClient) FindGroupByDisplayname(displayname string) (*Group, error) {
	filter := fmt.Sprintf("displayName eq \"%v\"", displayname)

	req, err := c.newRequest("GET", "Groups", filter, nil)

	if err != nil {
		return nil, err
	}

	var groupLR GroupListResponse
	_, err = c.do(req, &groupLR)
	if err != nil {
		return nil, err
	}

	if groupLR.TotalResults != 1 || len(groupLR.Resources) != 1 {
		return nil, fmt.Errorf("group \"%v\" not found", displayname)
	}

	return &groupLR.Resources[0], err
}

func (c *APIClient) CreateGroup(displayname string) (*Group, error) {

	body := map[string]interface{}{"displayName": displayname, "members": []string{}}

	req, err := c.newRequest("POST", "Groups", "", body)

	if err != nil {
		return nil, err
	}

	var groupResponse Group
	_, err = c.do(req, &groupResponse)

	return &groupResponse, err
}

func (c *APIClient) ReadGroup(id string) (*Group, error) {

	req, err := c.newRequest("GET", fmt.Sprintf("Groups/%v", id), "", nil)

	if err != nil {
		return nil, err
	}

	var groupResponse Group
	_, err = c.do(req, &groupResponse)

	return &groupResponse, err
}

func (c *APIClient) DeleteGroup(id string) error {

	req, err := c.newRequest("DELETE", fmt.Sprintf("Groups/%v", id), "", nil)

	if err != nil {
		return err
	}

	_, err = c.do(req, nil)

	return err
}

func (c *APIClient) TestGroupMember(group_id string, user_id string) (bool, error) {
	filter := fmt.Sprintf("id eq \"%v\" and members eq \"%v\"", group_id, user_id)

	req, err := c.newRequest("GET", "Groups", filter, nil)
	if err != nil {
		return false, err
	}

	var groupLR GroupListResponse
	_, err = c.do(req, &groupLR)
	if err != nil {
		return false, err
	}

	if groupLR.TotalResults != 1 || len(groupLR.Resources) != 1 {
		return false, err
	}

	return true, err
}

func (c *APIClient) AddGroupMember(group_id string, user_id string) error {

	opmsg := OperationMessage{
		Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
		Operations: []Operation{
			{
				Operation: "add",
				Path:      "members",
				Value: []map[string]string{
					{"value": user_id},
				},
			},
		},
	}

	req, err := c.newRequest("PATCH", fmt.Sprintf("Groups/%v", group_id), "", opmsg)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)

	return err
}

func (c *APIClient) RemoveGroupMember(group_id string, user_id string) error {

	opmsg := OperationMessage{
		Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
		Operations: []Operation{
			{
				Operation: "remove",
				Path:      "members",
				Value: []map[string]string{
					{"value": user_id},
				},
			},
		},
	}

	req, err := c.newRequest("PATCH", fmt.Sprintf("Groups/%v", group_id), "", opmsg)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)

	return err
}
