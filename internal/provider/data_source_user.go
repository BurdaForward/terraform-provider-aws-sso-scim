package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "Allows you to reference an existing user by user name and get the internal ID.",
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Description: "Reference by userName attribute.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"given_name": {
				Description: "Given name for the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"family_name": {
				Description: "Family name for the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email_address": {
				Description: "Primary email address.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email_type": {
				Description: "Usage type of the email adress, e.g. 'work'.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"active": {
				Description: "Set user to be active. Defaults to `false`.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*APIClient)

	user, _, err := client.FindUserByUsername(d.Get("user_name").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read User",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(user.ID)
	d.Set("user_name", user.UserName)
	d.Set("display_name", user.DisplayName)
	// check if name is set?
	d.Set("given_name", user.Name.GivenName)
	d.Set("family_name", user.Name.FamilyName)

	// AWS SSO SCIM only allows a single value for multi-valued properties like emails, so executes only once
	for _, v := range user.Emails {
		d.Set("email_address", v.Value)
		d.Set("email_type", v.Type)
	}

	return diags
}
