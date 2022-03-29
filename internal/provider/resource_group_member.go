package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Assigns a user to be member of a group.",

		CreateContext: resourceGroupMemberCreate,
		ReadContext:   resourceGroupMemberRead,
		DeleteContext: resourceGroupMemberDelete,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "Group identifier.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"user_id": {
				Description: "User identifier.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceGroupMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	_, err := client.AddGroupMember(d.Get("group_id").(string), d.Get("user_id").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create group member",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(fmt.Sprintf("%v,%v", d.Get("group_id"), d.Get("user_id")))

	return diags
}

func resourceGroupMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	is_member, resp, err := client.TestGroupMember(d.Get("group_id").(string), d.Get("user_id").(string))

	if err != nil {
		// if we get a 404, user might have vanished, so we remove this resource from the state.
		if resp.StatusCode == 404 {
			d.SetId("")
			return diags
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read group member",
			Detail:   err.Error(),
		})
		return diags
	}

	if is_member {
		d.SetId(fmt.Sprintf("%v,%v", d.Get("group_id"), d.Get("user_id")))
	}

	return diags
}

func resourceGroupMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	_, err := client.RemoveGroupMember(d.Get("group_id").(string), d.Get("user_id").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete group member",
			Detail:   err.Error(),
		})
		return diags
	}

	return diags
}
