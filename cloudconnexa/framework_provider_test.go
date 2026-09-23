package cloudconnexa

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMuxProviderServer builds the same mux server main.go serves: the SDK v2 provider plus the
// framework provider. Constructing it is what enforces provider-block schema parity.
func newMuxProviderServer(t *testing.T) tfprotov5.ProviderServer {
	t.Helper()
	mux, err := tf5muxserver.NewMuxServer(context.Background(),
		Provider().GRPCProvider,
		providerserver.NewProtocol5(NewFrameworkProvider()),
	)
	require.NoError(t, err)
	return mux.ProviderServer()
}

// TestUnitMuxServer_ProviderSchemaParity verifies that the SDK v2 and framework provider blocks
// are identical, which mux requires, and that the merged schema exposes both the SDK v2
// resources and the framework ephemeral resources.
func TestUnitMuxServer_ProviderSchemaParity(t *testing.T) {
	server := newMuxProviderServer(t)

	resp, err := server.GetProviderSchema(context.Background(), &tfprotov5.GetProviderSchemaRequest{})
	require.NoError(t, err)
	for _, d := range resp.Diagnostics {
		assert.NotEqual(t, tfprotov5.DiagnosticSeverityError, d.Severity, "%s: %s", d.Summary, d.Detail)
	}

	require.NotNil(t, resp.Provider)
	var names []string
	for _, a := range resp.Provider.Block.Attributes {
		names = append(names, a.Name)
	}
	assert.ElementsMatch(t, []string{"client_id", "client_secret", "base_url", "cloud_id"}, names)

	assert.Contains(t, resp.ResourceSchemas, "cloudconnexa_network_connector", "SDK v2 resources must be served through the mux")
	assert.Contains(t, resp.EphemeralResourceSchemas, "cloudconnexa_network_connector_token")
	assert.Contains(t, resp.EphemeralResourceSchemas, "cloudconnexa_host_connector_token")
	assert.Contains(t, resp.EphemeralResourceSchemas, "cloudconnexa_network_connector_profile")
	assert.Contains(t, resp.EphemeralResourceSchemas, "cloudconnexa_host_connector_profile")
	for name, s := range resp.EphemeralResourceSchemas {
		// The result attribute is named after the type name suffix: token or profile.
		value := name[strings.LastIndex(name, "_")+1:]
		var attrs []string
		for _, a := range s.Block.Attributes {
			attrs = append(attrs, a.Name)
			if a.Name == value {
				assert.True(t, a.Sensitive, "%s.%s must be sensitive", name, value)
			}
		}
		assert.ElementsMatch(t, []string{"connector_id", value}, attrs, name)
	}
}
