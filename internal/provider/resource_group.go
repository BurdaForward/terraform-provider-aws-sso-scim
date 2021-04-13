package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "SCIM Group resource.",

		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		DeleteContext: resourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Description: "displayName attribute.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(APIClient)
	diags := diag.Diagnostics{}

	group, err := client.CreateGroup(d.Get("display_name").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Group",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(group.ID)

	return diags
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(APIClient)
	diags := diag.Diagnostics{}

	group, err := client.ReadGroup(d.Get("id").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Group",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(group.ID)
	d.Set("display_name", group.DisplayName)

	return diags
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(APIClient)
	diags := diag.Diagnostics{}

	err := client.DeleteGroup(d.Get("id").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Group",
			Detail:   err.Error(),
		})
		return diags
	}

	return diags
}
