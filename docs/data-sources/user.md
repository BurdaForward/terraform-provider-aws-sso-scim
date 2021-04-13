---
page_title: "aws-sso-scim_user Data Source - terraform-provider-aws-sso-scim"
subcategory: ""
description: |-
  SCIM User data source.
---

# Data Source `aws-sso-scim_user`

SCIM User data source.

## Example Usage

```terraform
data "aws-sso-scim_user" "example" {
  user_name = "foo"
}
```

## Schema

### Required

- **user_name** (String) Reference by userName attribute.

### Read-only

- **display_name** (String)
- **id** (String) The ID of this resource.


