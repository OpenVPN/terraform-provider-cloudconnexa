package cloudconnexa

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// frameworkProvider is the terraform-plugin-framework half of the provider. It is served next
// to the SDK v2 provider through terraform-plugin-mux and currently owns only the ephemeral
// resources, which the SDK v2 cannot implement. Its provider block schema must stay identical
// to the SDK v2 one in Provider(): mux compares the two on every GetProviderSchema call and
// returns an error diagnostic when they differ, which fails every Terraform command.
type frameworkProvider struct{}

var _ provider.Provider = (*frameworkProvider)(nil)
var _ provider.ProviderWithEphemeralResources = (*frameworkProvider)(nil)

// NewFrameworkProvider returns the framework provider served alongside the SDK v2 provider.
func NewFrameworkProvider() provider.Provider {
	return &frameworkProvider{}
}

// frameworkProviderModel mirrors the provider block attributes.
type frameworkProviderModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	BaseURL      types.String `tfsdk:"base_url"`
	CloudID      types.String `tfsdk:"cloud_id"`
}

func (p *frameworkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudconnexa"
	resp.Version = version
}

func (p *frameworkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "The authentication client_id used to connect to CloudConnexa API. The value can be sourced from " +
					"the `CLOUDCONNEXA_CLIENT_ID` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The authentication client_secret used to connect to CloudConnexa API. The value can be sourced from " +
					"the `CLOUDCONNEXA_CLIENT_SECRET` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"base_url": schema.StringAttribute{
				Description: "The target CloudConnexa Base API URL in the format `https://[companyName].api.openvpn.com`",
				Optional:    true,
			},
			"cloud_id": schema.StringAttribute{
				Description: "Cloud ID",
				Optional:    true,
			},
		},
	}
}

// Configure builds the same API client as providerConfigure. Validation of the provider block
// (exactly one of base_url / cloud_id, credentials present) is left to the SDK v2 server, which
// receives the same configuration through mux, so errors are not reported twice.
func (p *frameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg frameworkProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	baseUrl, err := resolveBaseURL(cfg.BaseURL.ValueString(), cfg.CloudID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid provider configuration", err.Error())
		return
	}
	client, err := newAPIClient(ctx, baseUrl, stringOrEnv(cfg.ClientID, ClientIDEnvVar), stringOrEnv(cfg.ClientSecret, ClientSecretEnvVar))
	if err != nil {
		resp.Diagnostics.AddError("Unable to create CloudConnexa client", err.Error())
		return
	}
	resp.EphemeralResourceData = client
	resp.ResourceData = client
	resp.DataSourceData = client
}

// stringOrEnv returns the configured value, or the environment variable when it is unset,
// matching the SDK v2 schema.EnvDefaultFunc behaviour.
func stringOrEnv(v types.String, envVar string) string {
	if !v.IsNull() && !v.IsUnknown() {
		return v.ValueString()
	}
	return os.Getenv(envVar)
}

func (p *frameworkProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *frameworkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *frameworkProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		newNetworkConnectorTokenEphemeral,
		newHostConnectorTokenEphemeral,
		newNetworkConnectorProfileEphemeral,
		newHostConnectorProfileEphemeral,
	}
}
