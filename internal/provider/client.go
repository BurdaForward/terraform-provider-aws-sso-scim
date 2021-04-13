package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type APIClient struct {
	BaseURL    *url.URL
	Token      string
	httpClient *http.Client
	UserAgent  string
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

	if resp.StatusCode == 401 {
		return nil, errors.New("401 unauthorized")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

func (c *APIClient) ListUsers() ([]User, error) {
	req, err := c.newRequest("GET", "Users", "", nil)
	if err != nil {
		return nil, err
	}
	var userLR UserListResponse
	_, err = c.do(req, &userLR)
	return userLR.Resources, err
}

func (c *APIClient) FindUserByUsername(username string) (User, error) {
	filter := fmt.Sprintf("userName eq \"%v\"", username)

	req, err := c.newRequest("GET", "Users", filter, nil)

	if err != nil {
		return User{}, err
	}

	var userLR UserListResponse
	_, err = c.do(req, &userLR)

	if userLR.TotalResults != 1 || len(userLR.Resources) != 1 {
		return User{}, fmt.Errorf("user \"%v\" not found", username)
	}

	return userLR.Resources[0], err
}

func (c *APIClient) FindGroupByDisplayname(displayname string) (Group, error) {
	filter := fmt.Sprintf("displayName eq \"%v\"", displayname)

	req, err := c.newRequest("GET", "Groups", filter, nil)

	if err != nil {
		return Group{}, err
	}

	var groupLR GroupListResponse
	_, err = c.do(req, &groupLR)

	if groupLR.TotalResults != 1 || len(groupLR.Resources) != 1 {
		return Group{}, fmt.Errorf("group \"%v\" not found", displayname)
	}

	return groupLR.Resources[0], err
}

func (c *APIClient) CreateGroup(displayname string) (Group, error) {

	body := map[string]interface{}{"displayName": displayname, "members": []string{}}

	req, err := c.newRequest("POST", "Groups", "", body)

	if err != nil {
		return Group{}, err
	}

	var groupResponse Group
	_, err = c.do(req, &groupResponse)

	return groupResponse, err
}

func (c *APIClient) ReadGroup(id string) (Group, error) {

	req, err := c.newRequest("GET", fmt.Sprintf("Groups/%v", id), "", nil)

	if err != nil {
		return Group{}, err
	}

	var groupResponse Group
	_, err = c.do(req, &groupResponse)

	return groupResponse, err
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
