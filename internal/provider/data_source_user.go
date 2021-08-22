package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "SCIM User data source.",
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
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nick_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"preferred_language": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name_formatted": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"family_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"given_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"middle_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"honorific_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"honorific_suffix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"emails": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"phone_numbers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"formatted": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"street_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locality": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"postal_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"employee_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cost_center": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"division": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"department": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manager": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manager_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*APIClient)

	user, err := client.FindUserByUsername(d.Get("user_name").(string))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read User",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(user.ID)
	d.Set("display_name", user.DisplayName)
	d.Set("external_id", user.ExternalID)
	d.Set("nick_name", user.NickName)
	d.Set("profile_url", user.ProfileURL)
	d.Set("title", user.Title)
	d.Set("user_type", user.UserType)
	d.Set("preferred_language", user.PreferredLanguage)
	d.Set("locale", user.Locale)
	d.Set("timezone", user.Timezone)
	d.Set("active", user.Active)
	d.Set("name_formatted", user.Name.Formatted)
	d.Set("family_name", user.Name.FamilyName)
	d.Set("given_name", user.Name.GivenName)
	d.Set("middle_name", user.Name.MiddleName)
	d.Set("honorific_prefix", user.Name.HonorificPrefix)
	d.Set("honorific_suffix", user.Name.HonorificSuffix)

	emails := make([]map[string]interface{}, 0)
	for _, m := range user.Emails {
		emails = append(emails, map[string]interface{}{
			"primary": m.Primary,
			"type":    m.Type,
			"value":   m.Value,
		})
	}
	d.Set("emails", emails)

	phone_numbers := make([]map[string]interface{}, 0)
	for _, pn := range user.PhoneNumbers {
		phone_numbers = append(phone_numbers, map[string]interface{}{
			"type":  pn.Type,
			"value": pn.Value,
		})
	}
	d.Set("phone_numbers", phone_numbers)

	addresses := make([]map[string]interface{}, 0)
	for _, addr := range user.Addresses {
		addresses = append(addresses, map[string]interface{}{
			"formatted":      addr.Formatted,
			"street_address": addr.StreetAddress,
			"locality":       addr.Locality,
			"region":         addr.Region,
			"postal_code":    addr.PostalCode,
			"country":        addr.Country,
		})
	}
	d.Set("addresses", addresses)

	if user.EnterpriseUser != nil {
		d.Set("employee_number", user.EnterpriseUser.EmployeeNumber)
		d.Set("cost_center", user.EnterpriseUser.CostCenter)
		d.Set("organization", user.EnterpriseUser.Organization)
		d.Set("division", user.EnterpriseUser.Division)
		d.Set("department", user.EnterpriseUser.Department)
		d.Set("manager", user.EnterpriseUser.Manager.Value)
		d.Set("manager_ref", user.EnterpriseUser.Manager.Ref)
	}

	return diags
}
