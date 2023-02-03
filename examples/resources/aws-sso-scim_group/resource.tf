resource "aws-sso-scim_group" "example" {
  display_name = "bar"
}

resource "aws-sso-scim_group" "external_example" {
  display_name = "bar"
  external_id  = "e5a41517-bcd6-4b8b-8590-487ae996de44"
}
