package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Creates a new user.",

		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		UpdateContext: resourceUserUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		// "The givenName, familyName, userName, and displayName fields are required. "
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"given_name": {
				Description: "Given name for the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"family_name": {
				Description: "Family name for the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"display_name": {
				Description: "Display name for the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user_name": {
				Description: "Username for the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email_address": {
				Description: "Primary email address.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"email_type": {
				Description: "Usage type of the email adress, e.g. 'work'.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"active": {
				Description: "Set user to be active. Defaults to `false`.",
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

	if d.Get("email_address") != "" {
		new_user.Emails = []Email{
			{
				Value:   d.Get("email_address").(string),
				Primary: true,
			},
		}

		if d.Get("email_type") != "" {
			new_user.Emails[0].Type = d.Get("email_type").(string)
		}
	}

	user, _, err := client.CreateUser(&new_user)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create User",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(user.ID)

	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	user, resp, err := client.ReadUser(d.Id())

	if err != nil {
		// if we get a 404, user maybe has vanished, so we remove this resource from the state.
		if resp.StatusCode == 404 {
			d.SetId("")
			return diags
		}

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

	// AWS SSO SCIM only allows a single value for multi-valued properties like emails, so executes only once
	for _, v := range user.Emails {
		d.Set("email_address", v.Value)
		d.Set("email_type", v.Type)
	}

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	_, err := client.DeleteUser(d.Get("id").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete User",
			Detail:   err.Error(),
		})
		return diags
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	diags := diag.Diagnostics{}

	user, resp, err := client.ReadUser(d.Id())

	user.Meta = Meta{}

	if err != nil {
		// if we get a 404, user maybe has vanished, so we remove this resource from the state.
		if resp.StatusCode == 404 {
			d.SetId("")
			return diags
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read User",
			Detail:   err.Error(),
		})
		return diags
	}

	user.UserName = d.Get("user_name").(string)
	user.DisplayName = d.Get("display_name").(string)
	user.Name.FamilyName = d.Get("family_name").(string)
	user.Name.GivenName = d.Get("given_name").(string)
	user.Active = d.Get("active").(bool)

	if d.Get("email_address") != "" {
		user.Emails = []Email{
			{
				Value:   d.Get("email_address").(string),
				Primary: true,
			},
		}
		if d.Get("email_type") != "" {
			user.Emails[0].Type = d.Get("email_type").(string)
		}
	} else {
		user.Emails = []Email{}
	}

	_, resp, err = client.PutUser(user, d.Id())

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
