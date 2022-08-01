package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceUser(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"aws-sso-scim_user.foo", "id")),
			},
			{
				Config: testAccResourceUserUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("aws-sso-scim_user.foo", "id"),
					resource.TestCheckResourceAttr("aws-sso-scim_user.foo", "display_name", "terraform-test-temporary-user-display_name2"),
				),
			},
			{
				Config: testAccResourceUserWithEmail,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("aws-sso-scim_user.foo", "id"),
					resource.TestCheckResourceAttr("aws-sso-scim_user.foo", "display_name", "terraform-test-temporary-user-display_name3"),
					resource.TestCheckResourceAttr("aws-sso-scim_user.foo", "email_address", "terraformtest3@burda-forward.de"),
				),
			},
			{
				Config: testAccResourceUserWithEmailUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("aws-sso-scim_user.foo", "id"),
					resource.TestCheckResourceAttr("aws-sso-scim_user.foo", "display_name", "terraform-test-temporary-user-display_name4"),
					resource.TestCheckResourceAttr("aws-sso-scim_user.foo", "email_address", "terraformtest4@burda-forward.de"),
				),
			},
			{
				Config: testAccResourceUserWithoutEmailUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("aws-sso-scim_user.foo", "id"),
					resource.TestCheckResourceAttr("aws-sso-scim_user.foo", "display_name", "terraform-test-temporary-user-display_name5"),
				),
			},
		},
	})
}

const testAccResourceUser = `
resource "aws-sso-scim_user" "foo" {
  display_name = "terraform-test-temporary-user-display_name"
  user_name = "terraform-test-temporary-user-user_name"
  family_name = "terraform-test-temporary-user-family_name"
  given_name = "terraform-test-temporary-user-given_name"
  active = true
}
`

const testAccResourceUserUpdate = `
resource "aws-sso-scim_user" "foo" {
  display_name = "terraform-test-temporary-user-display_name2"
  user_name = "terraform-test-temporary-user-user_name2"
  family_name = "terraform-test-temporary-user-family_name2"
  given_name = "terraform-test-temporary-user-given_name2"
  active = false
}
`

const testAccResourceUserWithEmail = `
resource "aws-sso-scim_user" "foo" {
  display_name = "terraform-test-temporary-user-display_name3"
  user_name = "terraform-test-temporary-user-user_name3"
  family_name = "terraform-test-temporary-user-family_name3"
  given_name = "terraform-test-temporary-user-given_name3"
  active = true
	email_address = "terraformtest3@burda-forward.de"
}
`

const testAccResourceUserWithEmailUpdate = `
resource "aws-sso-scim_user" "foo" {
  display_name = "terraform-test-temporary-user-display_name4"
  user_name = "terraform-test-temporary-user-user_name4"
  family_name = "terraform-test-temporary-user-family_name4"
  given_name = "terraform-test-temporary-user-given_name4"
  active = false
	email_address = "terraformtest4@burda-forward.de"
	email_type = "work"
}
`

const testAccResourceUserWithoutEmailUpdate = `
resource "aws-sso-scim_user" "foo" {
  display_name = "terraform-test-temporary-user-display_name5"
  user_name = "terraform-test-temporary-user-user_name5"
  family_name = "terraform-test-temporary-user-family_name5"
  given_name = "terraform-test-temporary-user-given_name5"
  active = false
}
`
