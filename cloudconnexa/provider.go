package cloudconnexa

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ClientIDEnvVar is the environment variable name for the CloudConnexa client ID
const ClientIDEnvVar = "CLOUDCONNEXA_CLIENT_ID"

// ClientSecretEnvVar is the environment variable name for the CloudConnexa client secret
const ClientSecretEnvVar = "CLOUDCONNEXA_CLIENT_SECRET"

// cloudIDPattern validates that cloud_id contains only alphanumeric characters and hyphens
var cloudIDPattern = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

// version represents the current version of the Terraform provider
var version = "v1.3.0"

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
			"cloudconnexa_device":              resourceDevice(),
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
			"cloudconnexa_sessions":            dataSourceSessions(),
			"cloudconnexa_devices":             dataSourceDevices(),
			"cloudconnexa_device":              dataSourceDevice(),
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
	baseUrl, err := resolveBaseURL(d.Get("base_url").(string), d.Get("cloud_id").(string))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	cloudConnexaClient, err := newAPIClient(ctx, baseUrl, d.Get("client_id").(string), d.Get("client_secret").(string))
	if err != nil {
		return nil, diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Unable to create CloudConnexa client",
			Detail:   err.Error(),
		}}
	}
	return cloudConnexaClient, nil
}

// resolveBaseURL returns the API base URL from the provider's base_url, or derives it from
// cloud_id when base_url is empty. Shared by the SDK v2 and framework provider servers.
func resolveBaseURL(baseUrl, cloudId string) (string, error) {
	if baseUrl != "" || cloudId == "" {
		return baseUrl, nil
	}
	if !cloudIDPattern.MatchString(cloudId) {
		return "", fmt.Errorf("invalid cloud_id format: must contain only alphanumeric characters and hyphens")
	}
	return "https://" + cloudId + ".api.openvpn.com", nil
}

// apiClients caches clients by credentials so that the SDK v2 and framework servers behind the
// mux, which are both configured on every Terraform command, authenticate once and share one
// rate-limit back-off state.
var apiClients = struct {
	sync.Mutex
	m map[[3]string]*cloudconnexa.Client
}{m: map[[3]string]*cloudconnexa.Client{}}

// newAPIClient builds the CloudConnexa API client used as provider meta by both the SDK v2
// and the framework provider servers, so retry logging and the user agent stay identical.
func newAPIClient(ctx context.Context, baseUrl, clientId, clientSecret string) (*cloudconnexa.Client, error) {
	key := [3]string{baseUrl, clientId, clientSecret}
	apiClients.Lock()
	defer apiClients.Unlock()
	if c, ok := apiClients.m[key]; ok {
		return c, nil
	}
	cloudConnexaClient, err := cloudconnexa.NewClientWithOptions(baseUrl, clientId, clientSecret, &cloudconnexa.ClientOptions{
		OnRetry: newRetryLogger(ctx),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CloudConnexa client with base URL '%s': %v", baseUrl, err)
	}
	cloudConnexaClient.UserAgent = fmt.Sprintf("terraform-provider-cloudconnexa/%v", version)
	apiClients.m[key] = cloudConnexaClient
	return cloudConnexaClient, nil
}

// newRetryLogger returns a hook for cloudconnexa.ClientOptions.OnRetry that reports every
// rate-limit retry through the provider logger, so it is visible with TF_LOG=WARN or lower.
// ctx is the configure context: it carries the provider logger, and the SDK client builds its
// requests without a context, so the hook cannot take the logger from the request itself.
func newRetryLogger(ctx context.Context) func(*http.Request, int, time.Duration) {
	return func(req *http.Request, attempt int, wait time.Duration) {
		tflog.Warn(ctx, "CloudConnexa API rate limit reached, retrying request", map[string]interface{}{
			"method":  req.Method,
			"path":    req.URL.Path,
			"attempt": attempt,
			"wait":    wait.String(),
		})
	}
}
