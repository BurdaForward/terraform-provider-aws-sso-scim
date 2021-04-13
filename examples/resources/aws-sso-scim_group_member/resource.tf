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
