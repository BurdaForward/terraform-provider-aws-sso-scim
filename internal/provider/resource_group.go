package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Creates a new group.",

		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Description: "Display name for the group",
				Type:        schema.TypeString,
				Required:    true,
			},
			"external_id": {
				Description: "External ID for the group. This cannot be changed after creation.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	new_group := Group{
		DisplayName: d.Get("display_name").(string),
		ExternalID:  d.Get("external_id").(string),
	}

	group, _, err := client.CreateGroup(&new_group)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Group",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(group.ID)

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	group, resp, err := client.ReadGroup(d.Id())

	if err != nil {
		// if we get a 404, group maybe has vanished, so we remove this resource from the state.
		if resp.StatusCode == 404 {
			d.SetId("")
			return diags
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Group",
			Detail:   err.Error(),
		})
		return diags
	}

	d.Set("display_name", group.DisplayName)
	d.Set("external_id", group.ExternalID)

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	group, _, err := client.ReadGroup(d.Id())

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Group",
			Detail:   err.Error(),
		})
		return diags
	}

	group.DisplayName = d.Get("display_name").(string)
	group.ExternalID = d.Get("external_id").(string)

	opmsg := OperationMessage{
		Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
		Operations: []Operation{
			{
				Operation: "replace",
				Path:      "displayName",
				Value:     group.DisplayName,
			},
		},
	}

	if group.ExternalID != "" {
		opmsg.Operations = append(opmsg.Operations, Operation{
			Operation: "replace",
			Path:      "externalId",
			Value:     group.ExternalID,
		})
	}

	_, _, err = client.PatchGroup(&opmsg, d.Id())

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update User",
			Detail:   err.Error(),
		})
		return diags
	}

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	_, err := client.DeleteGroup(d.Get("id").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Group",
			Detail:   err.Error(),
		})
		return diags
	}

	return resourceGroupRead(ctx, d, meta)
}
