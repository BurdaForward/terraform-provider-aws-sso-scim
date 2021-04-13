---
page_title: "aws-sso-scim_group Resource - terraform-provider-aws-sso-scim"
subcategory: ""
description: |-
  SCIM Group resource.
---

# Resource `aws-sso-scim_group`

SCIM Group resource.

## Example Usage

```terraform
resource "aws-sso-scim_group" "example" {
  display_name = "bar"
}
```

## Schema

### Required

- **display_name** (String) displayName attribute.

### Read-only

- **id** (String) The ID of this resource.


