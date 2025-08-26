package cloudconnexa

import (
	"context"
	"fmt"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ClientIDEnvVar is the environment variable name for the CloudConnexa client ID
const ClientIDEnvVar = "CLOUDCONNEXA_CLIENT_ID"

// ClientSecretEnvVar is the environment variable name for the CloudConnexa client secret
const ClientSecretEnvVar = "CLOUDCONNEXA_CLIENT_SECRET"

// version represents the current version of the Terraform provider
var version = "v1.1.2"

// Token represents the authentication token structure returned by the CloudConnexa API
type Token struct {
	AccessToken string `json:"access_token"`
}

// Provider returns a Terraform provider for CloudConnexa.
// It configures the provider schema, available resources, and data sources.
//
// Returns:
//   - *schema.Provider: A configured Terraform provider instance
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Description: "The authentication client_id used to connect to CloudConnexa API. The value can be sourced from " +
					"the `CLOUDCONNEXA_CLIENT_ID` environment variable.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(ClientIDEnvVar, nil),
			},
			"client_secret": {
				Description: "The authentication client_secret used to connect to CloudConnexa API. The value can be sourced from " +
					"the `CLOUDCONNEXA_CLIENT_SECRET` environment variable.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(ClientSecretEnvVar, nil),
			},
			"base_url": {
				Description:  "The target CloudConnexa Base API URL in the format `https://[companyName].api.openvpn.com`",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"base_url", "cloud_id"},
			},
			"cloud_id": {
				Description: "Cloud ID",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudconnexa_network":             resourceNetwork(),
			"cloudconnexa_network_connector":   resourceNetworkConnector(),
			"cloudconnexa_host_connector":      resourceHostConnector(),
			"cloudconnexa_route":               resourceRoute(),
			"cloudconnexa_dns_record":          resourceDnsRecord(),
			"cloudconnexa_user":                resourceUser(),
			"cloudconnexa_host":                resourceHost(),
			"cloudconnexa_user_group":          resourceUserGroup(),
			"cloudconnexa_network_ip_service":  resourceNetworkIPService(),
			"cloudconnexa_host_ip_service":     resourceHostIPService(),
			"cloudconnexa_host_application":    resourceHostApplication(),
			"cloudconnexa_network_application": resourceNetworkApplication(),
			"cloudconnexa_location_context":    resourceLocationContext(),
			"cloudconnexa_access_group":        resourceAccessGroup(),
			"cloudconnexa_settings":            resourceSettings(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cloudconnexa_network":             dataSourceNetwork(),
			"cloudconnexa_network_connector":   dataSourceNetworkConnector(),
			"cloudconnexa_host_connector":      dataSourceHostConnector(),
			"cloudconnexa_user":                dataSourceUser(),
			"cloudconnexa_user_group":          dataSourceUserGroup(),
			"cloudconnexa_vpn_region":          dataSourceVpnRegion(),
			"cloudconnexa_vpn_regions":         dataSourceVpnRegions(),
			"cloudconnexa_network_routes":      dataSourceNetworkRoutes(),
			"cloudconnexa_host":                dataSourceHost(),
			"cloudconnexa_network_ip_service":  dataSourceNetworkIPService(),
			"cloudconnexa_host_ip_service":     dataSourceHostIPService(),
			"cloudconnexa_host_application":    dataSourceHostApplication(),
			"cloudconnexa_network_application": dataSourceNetworkApplication(),
			"cloudconnexa_location_context":    dataSourceLocationContext(),
			"cloudconnexa_access_group":        dataSourceAccessGroup(),
			"cloudconnexa_settings":            dataSourceSettings(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure configures the CloudConnexa client with the provided credentials and base URL.
// It sets up the client with proper authentication and user agent information.
//
// Parameters:
//   - ctx: The context for the operation
//   - d: The Terraform resource data containing provider configuration
//
// Returns:
//   - interface{}: The configured CloudConnexa client
//   - diag.Diagnostics: Diagnostics containing any errors that occurred during configuration
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	baseUrl := d.Get("base_url").(string)
	if baseUrl == "" {
		cloudId := d.Get("cloud_id").(string)
		baseUrl = "https://" + cloudId + ".api.openvpn.com"
	}
	cloudConnexaClient, err := cloudconnexa.NewClient(baseUrl, clientId, clientSecret)
	var diags diag.Diagnostics
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create CloudConnexa client",
			Detail:   fmt.Sprintf("Failed to create CloudConnexa client with base URL '%s': %v", baseUrl, err),
		})
		return nil, diags
	}
	cloudConnexaClient.UserAgent = fmt.Sprintf("terraform-provider-cloudconnexa/%v", version)
	return cloudConnexaClient, nil
}
