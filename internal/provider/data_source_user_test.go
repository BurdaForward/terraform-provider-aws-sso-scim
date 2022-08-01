package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aws-sso-scim_user.foo", "id")),
			},
			{
				Config: testAccDataSourceUserFull,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aws-sso-scim_user.foo", "id"),
					resource.TestCheckResourceAttr("data.aws-sso-scim_user.foo", "display_name", "terraform-test-temporary-user-display_nameXX"),
					resource.TestCheckResourceAttr("data.aws-sso-scim_user.foo", "user_name", "terraform-test-temporary-user-user_nameXX"),
					resource.TestCheckResourceAttr("data.aws-sso-scim_user.foo", "family_name", "terraform-test-temporary-user-family_nameXX"),
					resource.TestCheckResourceAttr("data.aws-sso-scim_user.foo", "given_name", "terraform-test-temporary-user-given_nameXX"),
					resource.TestCheckResourceAttr("data.aws-sso-scim_user.foo", "email_address", "terraformtestXX@burda-forward.de"),
					resource.TestCheckResourceAttr("data.aws-sso-scim_user.foo", "email_type", "work"),
				),
			},
		},
	})
}

const testAccDataSourceUser = `
data "aws-sso-scim_user" "foo" {
  user_name = "terraform-test-permanent-user"
}
`

const testAccDataSourceUserFull = `
resource "aws-sso-scim_user" "foo" {
  display_name = "terraform-test-temporary-user-display_nameXX"
  user_name = "terraform-test-temporary-user-user_nameXX"
  family_name = "terraform-test-temporary-user-family_nameXX"
  given_name = "terraform-test-temporary-user-given_nameXX"
  active = false
	email_address = "terraformtestXX@burda-forward.de"
	email_type = "work"
}

data "aws-sso-scim_user" "foo" {
  user_name = "terraform-test-temporary-user-user_nameXX"
	depends_on = [
		aws-sso-scim_user.foo
	]
}
`
