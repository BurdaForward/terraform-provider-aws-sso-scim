package provider

import (
	"time"
)

type Meta struct {
	ResourceType string    `json:"resourceType"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	Location     string    `json:"location,omitempty"`
	Version      string    `json:"version,omitempty"`
}

type Name struct {
	Formatted       *string `json:"formatted,omitempty"`
	FamilyName      string  `json:"familyName,omitempty"`
	GivenName       string  `json:"givenName,omitempty"`
	MiddleName      *string `json:"middleName,omitempty"`
	HonorificPrefix *string `json:"honorificPrefix,omitempty"`
	HonorificSuffix *string `json:"honorificSuffix,omitempty"`
}

type Email struct {
	Value   string `json:"value"`
	Type    string `json:"type"`
	Primary bool   `json:"primary"`
}

type PhoneNumber struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Address struct {
	Formatted     *string `json:"formatted,omitempty"`
	StreetAddress *string `json:"streetAddress,omitempty"`
	Locality      *string `json:"locality,omitempty"`
	Region        *string `json:"region,omitempty"`
	PostalCode    *string `json:"postalCode,omitempty"`
	Country       *string `json:"country,omitempty"`
}

type User struct {
	Meta              Meta            `json:"meta,omitempty"`
	ID                string          `json:"id"`
	ExternalID        string          `json:"externalId,omitempty"`
	UserName          string          `json:"userName"`
	Name              Name            `json:"name,omitempty"`
	DisplayName       string          `json:"displayName,omitempty"`
	NickName          *string         `json:"nickName,omitempty"`
	ProfileURL        *string         `json:"profileUrl,omitempty"`
	Title             *string         `json:"title,omitempty"`
	UserType          *string         `json:"userType,omitempty"`
	PreferredLanguage *string         `json:"preferredLanguage,omitempty"`
	Locale            *string         `json:"locale,omitempty"`
	Timezone          *string         `json:"timezone,omitempty"`
	Active            bool            `json:"active,omitempty"`
	Emails            []Email         `json:"emails,omitempty"`
	PhoneNumbers      []PhoneNumber   `json:"phoneNumbers,omitempty"`
	Addresses         []Address       `json:"addresses,omitempty"`
	EnterpriseUser    *EnterpriseUser `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`
	Schemas           []string        `json:"schemas"`
	Roles             []string        `json:"roles,omitempty"`
	Groups            []string        `json:"groups,omitempty"`
}

type Manager struct {
	Value string `json:"value"`
	Ref   string `json:"$ref,omitempty"`
	// displayName not supported by AWS
}

type EnterpriseUser struct {
	EmployeeNumber *string  `json:"employeeNumber"`
	CostCenter     *string  `json:"costCenter"`
	Organization   *string  `json:"organization"`
	Division       *string  `json:"division"`
	Department     *string  `json:"department"`
	Manager        *Manager `json:"manager"`
}

type Member struct {
	Value string `json:"value"`
	Ref   string `json:"$ref,omitempty"`
}

type Group struct {
	Meta        Meta     `json:"meta,omitempty"`
	ID          string   `json:"id"`
	ExternalID  string   `json:"externalId,omitempty"`
	DisplayName string   `json:"displayName"`
	Members     []Member `json:"members,omitempty"`
}

type UserListResponse struct {
	TotalResults int      `json:"totalResults"`
	StartIndex   int      `json:"startIndex,omitempty"`
	ItemsPerPage int      `json:"itemsPerPage,omitempty"`
	Resources    []User   `json:"Resources,omitempty"`
	Schemas      []string `json:"schemas"`
}

type GroupListResponse struct {
	TotalResults int      `json:"totalResults"`
	StartIndex   int      `json:"startIndex,omitempty"`
	ItemsPerPage int      `json:"itemsPerPage,omitempty"`
	Resources    []Group  `json:"Resources,omitempty"`
	Schemas      []string `json:"schemas"`
}

type Operation struct {
	Operation string      `json:"op"`
	Value     interface{} `json:"value"`
	Path      string      `json:"path,omitempty"`
}

type OperationMessage struct {
	Schemas    []string    `json:"schemas"`
	Operations []Operation `json:"Operations"`
}
