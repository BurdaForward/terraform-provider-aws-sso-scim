---
page_title: "aws-sso-scim Provider"
subcategory: ""
description: |-
  
---

# aws-sso-scim Provider



## Example Usage

```terraform
provider "aws-sso-scim" {
  endpoint = "https://scim.eu-central-1.amazonaws.com/<someid>/scim/v2/"
  token    = "***"
}
```

## Schema

### Optional

- **endpoint** (String)
- **token** (String)
