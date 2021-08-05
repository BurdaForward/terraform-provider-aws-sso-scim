data "aws_sso_scim_user" "example" {
  display_name = "foo"
}

data "aws_sso_scim_group" "example" {
  display_name = "bar"
}

resource "aws_sso_scim_group_member" "example" {
  user_id  = data.aws_sso_scim_user.example.id
  group_id = data.aws_sso_scim_group.example.id
}
