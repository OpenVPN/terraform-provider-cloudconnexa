package cloudconnexa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// resourceSamlSso returns a Terraform resource for managing CloudConnexa SAML SSO configuration
func resourceSamlSso() *schema.Resource {
	return &schema.Resource{
		Description:   "Use `cloudconnexa_saml_sso` to configure SAML Single Sign-On authentication",
		CreateContext: resourceSamlSsoUpdate,
		ReadContext:   resourceSamlSsoRead,
		DeleteContext: resourceSamlSsoDelete,
		UpdateContext: resourceSamlSsoUpdate,
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable or disable SAML SSO authentication.",
			},
			"entity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SAML entity ID (SP entity ID).",
			},
			"sso_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SAML SSO URL (IdP SSO URL).",
			},
			"slo_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SAML Single Logout URL (IdP SLO URL).",
			},
			"certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The X.509 certificate for SAML signature verification.",
			},
			"attribute_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "SAML attribute mapping configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "email",
							Description: "SAML attribute name for email. Defaults to 'email'.",
						},
						"first_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "firstName",
							Description: "SAML attribute name for first name. Defaults to 'firstName'.",
						},
						"last_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "lastName",
							Description: "SAML attribute name for last name. Defaults to 'lastName'.",
						},
						"groups": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "groups",
							Description: "SAML attribute name for groups. Defaults to 'groups'.",
						},
					},
				},
			},
			"auto_provision_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Automatically provision users from SAML assertions.",
			},
			"default_user_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default user group ID for auto-provisioned users.",
			},
			"require_signed_assertions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Require signed SAML assertions.",
			},
			"name_id_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
				ValidateFunc: validation.StringInSlice([]string{
					"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
					"urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
					"urn:oasis:names:tc:SAML:2.0:nameid-format:transient",
				}, false),
				Description: "SAML NameID format.",
			},
		},
	}
}

// resourceSamlSsoUpdate handles the creation and update of CloudConnexa SAML SSO configuration
func resourceSamlSsoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	// Build SAML configuration from resource data
	samlConfig := buildSamlConfigFromResourceData(d)

	// Update SAML SSO configuration
	err := updateSamlSsoConfig(c, samlConfig)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// Set the resource ID
	d.SetId("saml-sso-config")

	// Read back the configuration to ensure state consistency
	return resourceSamlSsoRead(ctx, d, m)
}

// resourceSamlSsoRead reads the current state of CloudConnexa SAML SSO configuration
func resourceSamlSsoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	// Get SAML SSO configuration
	samlConfig, err := getSamlSsoConfig(c)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// Update resource data with current configuration
	updateSamlSsoResourceData(d, samlConfig)

	return diags
}

// resourceSamlSsoDelete handles the deletion of CloudConnexa SAML SSO configuration
func resourceSamlSsoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	// Disable SAML SSO by setting connect_auth back to LOCAL
	_, err := c.Settings.SetDefaultConnectAuth("LOCAL")
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId("")
	return diags
}

// buildSamlConfigFromResourceData converts Terraform resource data to SAML configuration
// This is a helper function for when the actual API is implemented
func buildSamlConfigFromResourceData(d *schema.ResourceData) map[string]interface{} {
	config := map[string]interface{}{
		"enabled":                   d.Get("enabled").(bool),
		"entity_id":                 d.Get("entity_id").(string),
		"sso_url":                   d.Get("sso_url").(string),
		"slo_url":                   d.Get("slo_url").(string),
		"certificate":               d.Get("certificate").(string),
		"auto_provision_users":      d.Get("auto_provision_users").(bool),
		"default_user_group":        d.Get("default_user_group").(string),
		"require_signed_assertions": d.Get("require_signed_assertions").(bool),
		"name_id_format":            d.Get("name_id_format").(string),
	}

	if attrMapping := d.Get("attribute_mapping").([]interface{}); len(attrMapping) > 0 {
		mapping := attrMapping[0].(map[string]interface{})
		config["attribute_mapping"] = map[string]interface{}{
			"email":      mapping["email"].(string),
			"first_name": mapping["first_name"].(string),
			"last_name":  mapping["last_name"].(string),
			"groups":     mapping["groups"].(string),
		}
	}

	return config
}

// updateSamlSsoConfig updates the SAML SSO configuration via API calls
func updateSamlSsoConfig(c *cloudconnexa.Client, config map[string]interface{}) error {
	// Enable/disable SAML authentication
	enabled := config["enabled"].(bool)
	if enabled {
		_, err := c.Settings.SetDefaultConnectAuth("SAML")
		if err != nil {
			return err
		}
	} else {
		_, err := c.Settings.SetDefaultConnectAuth("LOCAL")
		if err != nil {
			return err
		}
	}

	// Note: Additional SAML configuration would be set here when the API supports it
	// For now, we only manage the basic enable/disable functionality through connect_auth

	return nil
}

// updateSamlSsoResourceData updates the Terraform resource data with SAML configuration
func updateSamlSsoResourceData(d *schema.ResourceData, config *SamlSsoConfig) {
	_ = d.Set("enabled", config.Enabled)
	_ = d.Set("entity_id", config.EntityID)
	_ = d.Set("sso_url", config.SsoURL)
	_ = d.Set("slo_url", config.SloURL)
	// Note: We don't set certificate as it's sensitive and write-only

	// Set attribute mapping
	if config.AttributeMapping != nil {
		attributeMapping := []interface{}{
			map[string]interface{}{
				"email":      config.AttributeMapping.Email,
				"first_name": config.AttributeMapping.FirstName,
				"last_name":  config.AttributeMapping.LastName,
				"groups":     config.AttributeMapping.Groups,
			},
		}
		_ = d.Set("attribute_mapping", attributeMapping)
	}

	// Set auto provisioning settings
	if config.AutoProvisioning != nil {
		_ = d.Set("auto_provision_users", config.AutoProvisioning.Enabled)
		_ = d.Set("default_user_group", config.AutoProvisioning.DefaultUserGroup)
	}

	// Set security settings
	if config.SecuritySettings != nil {
		_ = d.Set("require_signed_assertions", config.SecuritySettings.RequireSignedAssertions)
		_ = d.Set("name_id_format", config.SecuritySettings.NameIDFormat)
	}
}
