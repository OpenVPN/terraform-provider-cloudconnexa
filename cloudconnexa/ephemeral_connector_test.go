package cloudconnexa

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newConnectorEphemeralTestClient serves the OAuth token endpoint and the network connector
// lookups, plus the given handler on endpoint, and returns a client pointed at it.
func newConnectorEphemeralTestClient(t *testing.T, endpoint string, handler http.HandlerFunc) *cloudconnexa.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/oauth/token":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"tok"}`))
		case endpoint:
			handler(w, r)
		case "/api/v1/networks/connectors/conn-1":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"conn-1","name":"conn","tunnelingProtocol":"OPENVPN"}`))
		case "/api/v1/networks/connectors/ipsec-1":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"ipsec-1","name":"conn","tunnelingProtocol":"IPSEC"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	c, err := cloudconnexa.NewClientWithOptions(server.URL, "id", "secret", &cloudconnexa.ClientOptions{AllowInsecureHTTP: true, MaxRetries: -1})
	require.NoError(t, err)
	return c
}

// openConnectorEphemeral runs Configure and Open on the ephemeral resource with the given
// connector_id and returns the response and the value of its token or profile attribute, so tests
// exercise the real request/response types.
func openConnectorEphemeral(t *testing.T, e ephemeral.EphemeralResource, client *cloudconnexa.Client, connectorID string) (*ephemeral.OpenResponse, string) {
	t.Helper()
	ctx := context.Background()

	var schemaResp ephemeral.SchemaResponse
	e.Schema(ctx, ephemeral.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())

	if client != nil {
		var cfgResp ephemeral.ConfigureResponse
		e.(ephemeral.EphemeralResourceWithConfigure).Configure(ctx, ephemeral.ConfigureRequest{ProviderData: client}, &cfgResp)
		require.False(t, cfgResp.Diagnostics.HasError())
	}

	attribute := e.(*connectorEphemeral).attribute
	objType := schemaResp.Schema.Type().TerraformType(ctx)
	raw := tftypes.NewValue(objType, map[string]tftypes.Value{
		"connector_id": tftypes.NewValue(tftypes.String, connectorID),
		attribute:      tftypes.NewValue(tftypes.String, nil),
	})
	req := ephemeral.OpenRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: raw}}
	resp := &ephemeral.OpenResponse{Result: tfsdk.EphemeralResultData{Schema: schemaResp.Schema, Raw: raw}}

	e.Open(ctx, req, resp)

	var value types.String
	if !resp.Diagnostics.HasError() {
		require.False(t, resp.Result.GetAttribute(ctx, path.Root(attribute), &value).HasError())
	}
	return resp, value.ValueString()
}

// TestUnitConnectorEphemeral_Open verifies that each variant requests its own endpoint exactly
// once per Open and returns the response body as its sensitive result, next to the connector_id.
func TestUnitConnectorEphemeral_Open(t *testing.T) {
	tests := []struct {
		name        string
		newResource func() ephemeral.EphemeralResource
		connectorID string
		endpoint    string
	}{
		{"network token", newNetworkConnectorTokenEphemeral, "conn-1", "/api/v1/networks/connectors/conn-1/profile/encrypt"},
		{"host token", newHostConnectorTokenEphemeral, "conn-2", "/api/v1/hosts/connectors/conn-2/profile/encrypt"},
		{"network profile", newNetworkConnectorProfileEphemeral, "conn-1", "/api/v1/networks/connectors/conn-1/profile"},
		{"host profile", newHostConnectorProfileEphemeral, "conn-2", "/api/v1/hosts/connectors/conn-2/profile"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls atomic.Int32
			client := newConnectorEphemeralTestClient(t, tt.endpoint, func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				assert.Equal(t, http.MethodPost, r.Method)
				_, _ = w.Write([]byte("credential"))
			})

			resp, value := openConnectorEphemeral(t, tt.newResource(), client, tt.connectorID)

			require.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)
			assert.Equal(t, "credential", value)
			var connectorID types.String
			require.False(t, resp.Result.GetAttribute(context.Background(), path.Root("connector_id"), &connectorID).HasError())
			assert.Equal(t, tt.connectorID, connectorID.ValueString())
			assert.Equal(t, int32(1), calls.Load(), "exactly one request per Open")
		})
	}
}

// TestUnitConnectorEphemeral_ProfileNotReady verifies that the 202-with-empty-body answer both
// endpoints give while the profile is still being generated is surfaced as an error, not an
// empty value.
func TestUnitConnectorEphemeral_ProfileNotReady(t *testing.T) {
	tests := map[string]struct {
		newResource func() ephemeral.EphemeralResource
		endpoint    string
	}{
		"token":   {newNetworkConnectorTokenEphemeral, "/api/v1/networks/connectors/conn-1/profile/encrypt"},
		"profile": {newNetworkConnectorProfileEphemeral, "/api/v1/networks/connectors/conn-1/profile"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			client := newConnectorEphemeralTestClient(t, tt.endpoint, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			})

			resp, _ := openConnectorEphemeral(t, tt.newResource(), client, "conn-1")

			require.True(t, resp.Diagnostics.HasError())
			assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "not ready")
		})
	}
}

// TestUnitConnectorEphemeral_APIError verifies that an API error is returned as a diagnostic.
func TestUnitConnectorEphemeral_APIError(t *testing.T) {
	client := newConnectorEphemeralTestClient(t, "/api/v1/networks/connectors/conn-1/profile", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"errorMessage":"not found"}`, http.StatusNotFound)
	})

	resp, _ := openConnectorEphemeral(t, newNetworkConnectorProfileEphemeral(), client, "conn-1")

	require.True(t, resp.Diagnostics.HasError())
	assert.Equal(t, "Failed to retrieve network connector profile", resp.Diagnostics.Errors()[0].Summary())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "404")
}

// TestUnitNetworkConnectorEphemeral_IPsec verifies that neither a token nor a profile is
// requested for an IPsec connector, which has neither.
func TestUnitNetworkConnectorEphemeral_IPsec(t *testing.T) {
	tests := map[string]struct {
		newResource func() ephemeral.EphemeralResource
		endpoint    string
	}{
		"token":   {newNetworkConnectorTokenEphemeral, "/api/v1/networks/connectors/ipsec-1/profile/encrypt"},
		"profile": {newNetworkConnectorProfileEphemeral, "/api/v1/networks/connectors/ipsec-1/profile"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var calls atomic.Int32
			client := newConnectorEphemeralTestClient(t, tt.endpoint, func(w http.ResponseWriter, _ *http.Request) {
				calls.Add(1)
				_, _ = w.Write([]byte("unexpected"))
			})

			resp, _ := openConnectorEphemeral(t, tt.newResource(), client, "ipsec-1")

			require.True(t, resp.Diagnostics.HasError())
			assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "OPENVPN")
			assert.Equal(t, int32(0), calls.Load(), "no token or profile must be requested for an IPsec connector")
		})
	}
}

// TestUnitConnectorEphemeral_Unconfigured verifies that Open fails cleanly when the provider has
// not been configured, instead of dereferencing a nil client.
func TestUnitConnectorEphemeral_Unconfigured(t *testing.T) {
	resp, _ := openConnectorEphemeral(t, newNetworkConnectorProfileEphemeral(), nil, "conn-1")

	require.True(t, resp.Diagnostics.HasError())
	assert.Equal(t, "Provider not configured", resp.Diagnostics.Errors()[0].Summary())
}
