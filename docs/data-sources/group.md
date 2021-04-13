---
page_title: "aws-sso-scim_group Data Source - terraform-provider-aws-sso-scim"
subcategory: ""
description: |-
  SCIM Group data source.
---

# Data Source `aws-sso-scim_group`

SCIM Group data source.

## Example Usage

```terraform
data "aws-sso-scim_group" "example" {
  display_name = "bar"
}
```

## Schema

### Required

- **display_name** (String) Reference by displayName attribute.

### Read-only

- **id** (String) The ID of this resource.


