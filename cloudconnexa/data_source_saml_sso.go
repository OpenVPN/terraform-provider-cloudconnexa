package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// dataSourceSamlSso returns a Terraform data source for CloudConnexa SAML SSO configuration
func dataSourceSamlSso() *schema.Resource {
	return &schema.Resource{
		Description: "Use a `cloudconnexa_saml_sso` data source to read CloudConnexa SAML SSO configuration.",
		ReadContext: dataSourceSamlSsoRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of this resource.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether SAML SSO authentication is enabled.",
			},
			"entity_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SAML entity ID (SP entity ID).",
			},
			"sso_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SAML SSO URL (IdP SSO URL).",
			},
			"slo_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SAML Single Logout URL (IdP SLO URL).",
			},
			"attribute_mapping": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "SAML attribute mapping configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SAML attribute name for email.",
						},
						"first_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SAML attribute name for first name.",
						},
						"last_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SAML attribute name for last name.",
						},
						"groups": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SAML attribute name for groups.",
						},
					},
				},
			},
			"auto_provision_users": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether users are automatically provisioned from SAML assertions.",
			},
			"default_user_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default user group ID for auto-provisioned users.",
			},
			"require_signed_assertions": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether signed SAML assertions are required.",
			},
			"name_id_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SAML NameID format.",
			},
		},
	}
}

// dataSourceSamlSsoRead handles the read operation for the SAML SSO data source
func dataSourceSamlSsoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_ = m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	// Note: This is a placeholder implementation
	// In a real implementation, you would call something like:
	// samlConfig, err := c.SamlSso.Get()
	// if err != nil {
	//     return append(diags, diag.FromErr(err)...)
	// }
	//
	// d.SetId("saml-sso-config")
	// d.Set("enabled", samlConfig.Enabled)
	// d.Set("entity_id", samlConfig.EntityID)
	// d.Set("sso_url", samlConfig.SsoURL)
	// d.Set("slo_url", samlConfig.SloURL)
	// ... set other attributes

	d.SetId("saml-sso-config")
	return diags
}
