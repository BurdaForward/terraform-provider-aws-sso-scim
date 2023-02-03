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
				Config: testAccResourceGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("aws-sso-scim_group.foo", "id"),
					resource.TestCheckResourceAttr("aws-sso-scim_group.foo", "display_name", "terraform-test-temporary-group2"),
				),
			},
			{
				Config: testAccResourceGroupExternal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"aws-sso-scim_group.external_group", "id")),
			},
			{
				Config: testAccResourceGroupExternalUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("aws-sso-scim_group.external_group", "id"),
					resource.TestCheckResourceAttr("aws-sso-scim_group.external_group", "display_name", "terraform-test-temporary-external-group2"),
				),
			},
		},
	})
}

const testAccResourceGroup = `
resource "aws-sso-scim_group" "foo" {
  display_name = "terraform-test-temporary-group"
}
`

const testAccResourceGroupUpdate = `
resource "aws-sso-scim_group" "foo" {
  display_name = "terraform-test-temporary-group2"
}
`

const testAccResourceGroupExternal = `
resource "aws-sso-scim_group" "external_group" {
  display_name = "terraform-test-temporary-external-group"
  external_id  = "e5a41517-bcd6-4b8b-8590-487ae996de44"
}
`

const testAccResourceGroupExternalUpdate = `
resource "aws-sso-scim_group" "external_group" {
  display_name = "terraform-test-temporary-external-group2"
  external_id  = "e5a41517-bcd6-4b8b-8590-487ae996de44"
}
`
