package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "SCIM Group data source.",
		ReadContext: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Description: "Reference by displayName attribute.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*APIClient)

	group, err := client.FindGroupByDisplayname(d.Get("display_name").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Group",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(group.ID)

	return diags
}
