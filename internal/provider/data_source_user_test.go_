package provider

import (
	"regexp"
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
					resource.TestMatchResourceAttr(
						"data.aws-sso-scim_user.foo", "user_name", regexp.MustCompile("^username")),
				),
			},
		},
	})
}

const testAccDataSourceUser = `
data "aws-sso-scim_user" "foo" {
  user_name = "username"
}
`
