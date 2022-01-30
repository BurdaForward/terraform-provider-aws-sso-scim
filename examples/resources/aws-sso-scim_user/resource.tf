resource "aws-sso-scim_user" "example" {
    user_name = "john.doe@example.com"
    given_name = "John"
    family_name = "Doe"
    display_name = "John Doe"
}
