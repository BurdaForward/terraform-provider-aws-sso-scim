package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceGroup(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"aws-sso-scim_group.foo", "id")),
			},
			{
				Config: testAccResourceGroupExternal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"aws-sso-scim_group.external_group", "id")),
			},
		},
	})
}

const testAccResourceGroup = `
resource "aws-sso-scim_group" "foo" {
  display_name = "terraform-test-temporary-group"
}
`

const testAccResourceGroupExternal = `
resource "aws-sso-scim_group" "external_group" {
  display_name = "ExampleGroup"
  external_id = "e5a41517-bcd6-4b8b-8590-487ae996de44"
}
`
