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
