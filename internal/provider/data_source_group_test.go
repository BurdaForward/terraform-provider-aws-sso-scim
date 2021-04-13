package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGroup(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.aws-sso-scim_group.foo", "display_name", regexp.MustCompile("^groupname")),
				),
			},
		},
	})
}

const testAccDataSourceGroup = `
data "aws-sso-scim_group" "foo" {
  display_name = "groupname"
}
`
