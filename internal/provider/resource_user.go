package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "SCIM User resource.",

		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		UpdateContext: resourceUserUpdate,

		// "The givenName, familyName, userName, and displayName fields are required. "
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"given_name": {
				Description: "Given Name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"family_name": {
				Description: "Family Name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"display_name": {
				Description: "Display Name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user_name": {
				Description: "Username",
				Type:        schema.TypeString,
				Required:    true,
			},
			"active": {
				Description: "Active",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	new_user := User{
		UserName:    d.Get("user_name").(string),
		DisplayName: d.Get("display_name").(string),
		Name: Name{
			FamilyName: d.Get("family_name").(string),
			GivenName:  d.Get("given_name").(string),
		},
		Active: d.Get("active").(bool),
	}

	user, err := client.CreateUser(&new_user)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create User",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(user.ID)

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	user, err := client.ReadUser(d.Id())

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read User",
			Detail:   err.Error(),
		})
		return diags
	}

	d.Set("display_name", user.DisplayName)
	d.Set("user_name", user.UserName)
	d.Set("family_name", user.Name.FamilyName)
	d.Set("given_name", user.Name.GivenName)
	d.Set("active", user.Active)

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	err := client.DeleteUser(d.Get("id").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete User",
			Detail:   err.Error(),
		})
		return diags
	}

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	opmsg := OperationMessage{}

	required_map := map[string]string{
		"user_name":    "userName",
		"display_name": "displayName",
		"family_name":  "name.familyName",
		"given_name":   "name.givenName",
		"active":       "active",
	}

	for attribute, path := range required_map {
		// These attributes are required and can only be replaced
		if d.HasChange(attribute) {
			opmsg.Operations = append(opmsg.Operations, Operation{
				Operation: "replace",
				Path:      path,
				Value:     d.Get(attribute),
			})
		}
	}

	_, err := client.PatchUser(&opmsg, d.Id())

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update User",
			Detail:   err.Error(),
		})
		return diags
	}

	return resourceUserRead(ctx, d, meta)
}
