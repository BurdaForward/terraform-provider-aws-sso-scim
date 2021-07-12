package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"aws-sso-scim_user":  dataSourceUser(),
				"aws-sso-scim_group": dataSourceGroup(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"aws-sso-scim_user":         resourceUser(),
				"aws-sso-scim_group":        resourceGroup(),
				"aws-sso-scim_group_member": resourceGroupMember(),
			},
			Schema: map[string]*schema.Schema{
				"endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("AWS_SSO_SCIM_ENDPOINT", nil),
				},
				"token": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("AWS_SSO_SCIM_TOKEN", nil),
				},
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		endpoint := d.Get("endpoint").(string)
		token := d.Get("token").(string)
		userAgent := p.UserAgent("terraform-provider-aws-sso-scim", version)

		apiClient, err := NewClient(endpoint, token, userAgent)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return apiClient, diags
	}
}
