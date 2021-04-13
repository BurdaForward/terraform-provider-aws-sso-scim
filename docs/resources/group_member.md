---
page_title: "aws-sso-scim_group_member Resource - terraform-provider-aws-sso-scim"
subcategory: ""
description: |-
  SCIM Group resource member.
---

# Resource `aws-sso-scim_group_member`

SCIM Group resource member.

## Example Usage

```terraform
data "aws-sso-scim_user" "example" {
  display_name = "foo"
}

data "aws-sso-scim_group" "example" {
  display_name = "bar"
}

resource "aws-sso-scim_group_member" "example" {
  user_id  = data.aws-sso-scim_user.example.id
  group_id = data.aws-sso-scim_group.example.id
}
```

## Schema

### Required

- **group_id** (String) Group identifier.
- **user_id** (String) User identifier.

### Optional

- **id** (String) The ID of this resource.


