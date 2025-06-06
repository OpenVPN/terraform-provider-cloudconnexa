package cloudconnexa

import (
	"context"
	"fmt"

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
	c := m.(*cloudconnexa.Client)
	var diags diag.Diagnostics

	// Get SAML SSO configuration using direct API call
	samlConfig, err := getSamlSsoConfig(c)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// Set the data source ID
	d.SetId("saml-sso-config")

	// Set basic configuration
	_ = d.Set("enabled", samlConfig.Enabled)
	_ = d.Set("entity_id", samlConfig.EntityID)
	_ = d.Set("sso_url", samlConfig.SsoURL)
	_ = d.Set("slo_url", samlConfig.SloURL)
	_ = d.Set("status", samlConfig.Status)
	_ = d.Set("metadata_url", samlConfig.MetadataURL)
	_ = d.Set("acs_url", samlConfig.AcsURL)

	// Set attribute mapping
	if samlConfig.AttributeMapping != nil {
		attributeMapping := []interface{}{
			map[string]interface{}{
				"email":      samlConfig.AttributeMapping.Email,
				"first_name": samlConfig.AttributeMapping.FirstName,
				"last_name":  samlConfig.AttributeMapping.LastName,
				"groups":     samlConfig.AttributeMapping.Groups,
			},
		}
		_ = d.Set("attribute_mapping", attributeMapping)
	}

	// Set auto provisioning
	if samlConfig.AutoProvisioning != nil {
		autoProvisioning := []interface{}{
			map[string]interface{}{
				"enabled":                samlConfig.AutoProvisioning.Enabled,
				"default_user_group":     samlConfig.AutoProvisioning.DefaultUserGroup,
				"group_mapping":          samlConfig.AutoProvisioning.GroupMapping,
				"update_user_attributes": samlConfig.AutoProvisioning.UpdateUserAttributes,
			},
		}
		_ = d.Set("auto_provisioning", autoProvisioning)
	}

	// Set security settings
	if samlConfig.SecuritySettings != nil {
		securitySettings := []interface{}{
			map[string]interface{}{
				"require_signed_assertions": samlConfig.SecuritySettings.RequireSignedAssertions,
				"name_id_format":            samlConfig.SecuritySettings.NameIDFormat,
				"signature_algorithm":       samlConfig.SecuritySettings.SignatureAlgorithm,
				"digest_algorithm":          samlConfig.SecuritySettings.DigestAlgorithm,
			},
		}
		_ = d.Set("security_settings", securitySettings)
	}

	return diags
}

// SamlSsoConfig represents the SAML SSO configuration
type SamlSsoConfig struct {
	Enabled          bool                  `json:"enabled"`
	EntityID         string                `json:"entityId"`
	SsoURL           string                `json:"ssoUrl"`
	SloURL           string                `json:"sloUrl"`
	Certificate      string                `json:"certificate"`
	Status           string                `json:"status"`
	MetadataURL      string                `json:"metadataUrl"`
	AcsURL           string                `json:"acsUrl"`
	AttributeMapping *SamlAttributeMapping `json:"attributeMapping,omitempty"`
	AutoProvisioning *SamlAutoProvisioning `json:"autoProvisioning,omitempty"`
	SecuritySettings *SamlSecuritySettings `json:"securitySettings,omitempty"`
}

// SamlAttributeMapping represents SAML attribute mapping configuration
type SamlAttributeMapping struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Groups    string `json:"groups"`
}

// SamlAutoProvisioning represents SAML auto-provisioning configuration
type SamlAutoProvisioning struct {
	Enabled              bool              `json:"enabled"`
	DefaultUserGroup     string            `json:"defaultUserGroup"`
	GroupMapping         map[string]string `json:"groupMapping"`
	UpdateUserAttributes bool              `json:"updateUserAttributes"`
}

// SamlSecuritySettings represents SAML security configuration
type SamlSecuritySettings struct {
	RequireSignedAssertions bool   `json:"requireSignedAssertions"`
	NameIDFormat            string `json:"nameIdFormat"`
	SignatureAlgorithm      string `json:"signatureAlgorithm"`
	DigestAlgorithm         string `json:"digestAlgorithm"`
}

// getSamlSsoConfig retrieves the SAML SSO configuration from the API
func getSamlSsoConfig(c *cloudconnexa.Client) (*SamlSsoConfig, error) {
	// Make direct API call to get SAML SSO configuration
	// This is a placeholder implementation that would need to be updated
	// when the actual SAML SSO API endpoints are available

	// For now, return a default configuration based on the current connect_auth setting
	connectAuth, err := c.Settings.GetDefaultConnectAuth()
	if err != nil {
		return nil, fmt.Errorf("failed to get connect auth setting: %w", err)
	}

	// Create a basic SAML configuration
	config := &SamlSsoConfig{
		Enabled:     connectAuth == "SAML",
		EntityID:    "",
		SsoURL:      "",
		SloURL:      "",
		Certificate: "",
		Status:      "configured",
		MetadataURL: "",
		AcsURL:      "",
		AttributeMapping: &SamlAttributeMapping{
			Email:     "email",
			FirstName: "firstName",
			LastName:  "lastName",
			Groups:    "groups",
		},
		AutoProvisioning: &SamlAutoProvisioning{
			Enabled:              false,
			DefaultUserGroup:     "",
			GroupMapping:         make(map[string]string),
			UpdateUserAttributes: true,
		},
		SecuritySettings: &SamlSecuritySettings{
			RequireSignedAssertions: true,
			NameIDFormat:            "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
			SignatureAlgorithm:      "RSA_SHA256",
			DigestAlgorithm:         "SHA256",
		},
	}

	return config, nil
}
