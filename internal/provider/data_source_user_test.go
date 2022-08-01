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
					resource.TestCheckResourceAttrSet("data.aws-sso-scim_user.foo", "id"),
					resource.TestCheckResourceAttr("data.aws-sso-scim_user.foo", "display_name", "terraform-test-permanent-user"),
					resource.TestCheckResourceAttr("data.aws-sso-scim_user.foo", "user_name", "terraform-test-permanent-user"),
					resource.TestCheckResourceAttr("data.aws-sso-scim_user.foo", "family_name", "permanent-user"),
					resource.TestCheckResourceAttr("data.aws-sso-scim_user.foo", "given_name", "terraform-test"),
					resource.TestCheckResourceAttr("data.aws-sso-scim_user.foo", "email_address", "terraformtest@burda-forward.de"),
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
