package cloudconnexa

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Fixtures and helpers
// ---------------------------------------------------------------------------

// fullIpSecConfig is every ipsec_config attribute except the dead peer detection
// ones, which are what the individual tests vary. Values are deliberately distinct
// so a field swapped for its neighbour shows up as a failure.
func fullIpSecConfig() map[string]cty.Value {
	return map[string]cty.Value{
		"platform":                        cty.StringVal("AWS"),
		"authentication_type":             cty.StringVal("SHARED_SECRET"),
		"remote_site_public_ip":           cty.StringVal("203.0.113.10"),
		"pre_shared_key":                  cty.StringVal("psk-value"),
		"ca_certificate":                  cty.StringVal("ca-cert"),
		"peer_certificate":                cty.StringVal("peer-cert"),
		"remote_gateway_certificate":      cty.StringVal("gateway-cert"),
		"peer_certificate_private_key":    cty.StringVal("peer-key"),
		"peer_certificate_key_passphrase": cty.StringVal("peer-passphrase"),
		"protocol_version":                cty.StringVal("IKE_V2"),
		"startup_action":                  cty.StringVal("START"),
		"phase_1_encryption_algorithms":   cty.ListVal([]cty.Value{cty.StringVal("AES256"), cty.StringVal("AES128")}),
		"phase_1_integrity_algorithms":    cty.ListVal([]cty.Value{cty.StringVal("SHA2_256")}),
		"phase_1_diffie_hellman_groups":   cty.ListVal([]cty.Value{cty.StringVal("G_14"), cty.StringVal("G_19")}),
		"phase_1_lifetime_sec":            cty.NumberIntVal(28800),
		"phase_2_encryption_algorithms":   cty.ListVal([]cty.Value{cty.StringVal("AES128_GCM_16")}),
		"phase_2_integrity_algorithms":    cty.ListVal([]cty.Value{cty.StringVal("SHA2_512")}),
		"phase_2_diffie_hellman_groups":   cty.ListVal([]cty.Value{cty.StringVal("G_20")}),
		"phase_2_lifetime_sec":            cty.NumberIntVal(3600),
		"margin_time_sec":                 cty.NumberIntVal(270),
		"fuzz_percent":                    cty.NumberIntVal(100),
		"replay_window_size":              cty.NumberIntVal(1024),
		"hostname":                        cty.StringVal("connector-host"),
		"domain":                          cty.StringVal("example.com"),
	}
}

// withAttrs returns a copy of base with extra merged over it.
func withAttrs(base map[string]cty.Value, extra map[string]cty.Value) map[string]cty.Value {
	out := make(map[string]cty.Value, len(base)+len(extra))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}

// ipSecRawConfig builds a resource-level raw config carrying a single ipsec_config
// block in which only the named attributes are set. Every other attribute is null,
// which is how Terraform represents one the user did not write — the distinction the
// dead peer detection rules turn on.
func ipSecRawConfig(t *testing.T, attrs map[string]cty.Value) cty.Value {
	t.Helper()
	resourceType := schema.InternalMap(resourceNetworkConnector().Schema).CoreConfigSchema().ImpliedType()
	blockType := resourceType.AttributeType("ipsec_config").ElementType()

	block := make(map[string]cty.Value)
	for name, attrType := range blockType.AttributeTypes() {
		if value, ok := attrs[name]; ok {
			block[name] = value
		} else {
			block[name] = cty.NullVal(attrType)
		}
	}

	values := make(map[string]cty.Value)
	for name, attrType := range resourceType.AttributeTypes() {
		values[name] = cty.NullVal(attrType)
	}
	values["ipsec_config"] = cty.ListVal([]cty.Value{cty.ObjectVal(block)})
	return cty.ObjectVal(values)
}

// ctyToGo converts the cty values used in these fixtures back to the plain Go
// representation TestResourceDataRaw expects, so the flatmap attributes and the raw
// config are always generated from one description of the config.
func ctyToGo(v cty.Value) interface{} {
	switch {
	case v.IsNull():
		return nil
	case v.Type() == cty.String:
		return v.AsString()
	case v.Type() == cty.Bool:
		return v.True()
	case v.Type() == cty.Number:
		i, _ := v.AsBigFloat().Int64()
		return int(i)
	case v.Type().IsListType():
		out := []interface{}{}
		for it := v.ElementIterator(); it.Next(); {
			_, element := it.Element()
			out = append(out, ctyToGo(element))
		}
		return out
	default:
		panic("unsupported fixture type: " + v.Type().GoString())
	}
}

// ipSecConfigBlock renders attrs as the plain Go map a Terraform configuration would
// produce, dropping the attributes left unset.
func ipSecConfigBlock(attrs map[string]cty.Value) map[string]interface{} {
	block := make(map[string]interface{})
	for name, value := range attrs {
		if converted := ctyToGo(value); converted != nil {
			block[name] = converted
		}
	}
	return block
}

// networkConnectorState builds the InstanceState for a connector carrying an
// ipsec_config block described by attrs. The flatmap attributes (what Get reads) and
// the raw config (what tells an unset attribute from a false one) are both derived
// from attrs, so the two can never disagree.
//
// The state is assembled by hand rather than via ResourceData.State(), which needs a
// non-empty ID and does not carry RawConfig through.
func networkConnectorState(t *testing.T, id string, attrs map[string]cty.Value) *terraform.InstanceState {
	t.Helper()
	seed := schema.TestResourceDataRaw(t, resourceNetworkConnector().Schema, map[string]interface{}{
		"name":          "test-connector",
		"description":   "managed by unit test",
		"network_id":    "net-1",
		"vpn_region_id": "us-east-1",
		"ipsec_config":  []interface{}{ipSecConfigBlock(attrs)},
	})
	seed.SetId("seed")

	state := seed.State()
	require.NotNil(t, state)
	state.ID = id
	state.Attributes["id"] = id
	state.RawConfig = ipSecRawConfig(t, attrs)
	return state
}

// networkConnectorData returns a *schema.ResourceData over networkConnectorState.
func networkConnectorData(t *testing.T, id string, attrs map[string]cty.Value) *schema.ResourceData {
	t.Helper()
	return resourceNetworkConnector().Data(networkConnectorState(t, id, attrs))
}

// networkConnectorConfig renders attrs as the ResourceConfig Terraform would hand the
// SDK for a connector carrying that ipsec_config block.
func networkConnectorConfig(attrs map[string]cty.Value) *terraform.ResourceConfig {
	return terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":          "test-connector",
		"description":   "managed by unit test",
		"network_id":    "net-1",
		"vpn_region_id": "us-east-1",
		"ipsec_config":  []interface{}{ipSecConfigBlock(attrs)},
	})
}

// networkConnectorUpdateData returns a *schema.ResourceData carrying a real diff
// between the prior state (from) and the configuration being applied (to), so
// HasChanges reports the change the way it does during a genuine apply. A ResourceData
// built from state alone sees no change and the update is skipped entirely.
func networkConnectorUpdateData(t *testing.T, id string, from, to map[string]cty.Value) *schema.ResourceData {
	t.Helper()
	resource := resourceNetworkConnector()

	state := networkConnectorState(t, id, from)
	// GetRawConfig must report the configuration being applied, not the prior one.
	state.RawConfig = ipSecRawConfig(t, to)

	diff, err := resource.Diff(context.Background(), state, networkConnectorConfig(to), nil)
	require.NoError(t, err)
	require.NotNil(t, diff, "expected a diff between the two configurations")

	d, err := schema.InternalMap(resource.Schema).Data(state, diff)
	require.NoError(t, err)
	return d
}

// capturingServer records every request path and body, and answers each route with
// the supplied responses. It stands in for the CloudConnexa API.
type capturingServer struct {
	mu       sync.Mutex
	requests []capturedRequest
}

type capturedRequest struct {
	method string
	path   string
	body   []byte
}

// newNetworkConnectorTestClient wires a client to a local server that answers the
// OAuth handshake plus the supplied handler. Mirrors newHostUnitTestClient.
func newNetworkConnectorTestClient(t *testing.T, handler http.Handler) (*cloudconnexa.Client, *capturingServer) {
	t.Helper()
	captured := &capturingServer{}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/oauth/token", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"unit-test-token"}`))
	})
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		captured.mu.Lock()
		captured.requests = append(captured.requests, capturedRequest{method: r.Method, path: r.URL.Path, body: body})
		captured.mu.Unlock()
		r.Body = io.NopCloser(bytes.NewReader(body))

		// The client defaults to one write per four seconds, which would make any test
		// issuing a second write (create then start the tunnel) sleep. Report a
		// generous budget through the rate limit headers the client reads, the way a
		// real response does, instead of waiting it out.
		w.Header().Set("X-RateLimit-Replenish-Rate", "1000")
		w.Header().Set("X-RateLimit-Replenish-Time", "1")
		w.Header().Set("X-RateLimit-Remaining", "1000")

		handler.ServeHTTP(w, r)
	}))

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	c, err := cloudconnexa.NewClientWithOptions(server.URL, "test-id", "test-secret", &cloudconnexa.ClientOptions{
		AllowInsecureHTTP: true,
	})
	require.NoError(t, err)
	return c, captured
}

// bodyFor returns the decoded JSON body of the first request matching method+path.
func (s *capturingServer) bodyFor(t *testing.T, method, path string) map[string]interface{} {
	t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, request := range s.requests {
		if request.method == method && request.path == path {
			var decoded map[string]interface{}
			require.NoError(t, json.Unmarshal(request.body, &decoded), "body: %s", request.body)
			return decoded
		}
	}
	t.Fatalf("no %s request to %s; saw %v", method, path, s.paths())
	return nil
}

// called reports whether a request with the given method and path was made.
func (s *capturingServer) called(method, path string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, request := range s.requests {
		if request.method == method && request.path == path {
			return true
		}
	}
	return false
}

func (s *capturingServer) paths() []string {
	out := []string{}
	for _, request := range s.requests {
		out = append(out, request.method+" "+request.path)
	}
	return out
}

// ---------------------------------------------------------------------------
// Wire payload: what actually reaches the API
// ---------------------------------------------------------------------------

// ipSecJSON digs the ipSecConfig object out of a captured request body.
func ipSecJSON(t *testing.T, body map[string]interface{}) map[string]interface{} {
	t.Helper()
	config, ok := body["ipSecConfig"].(map[string]interface{})
	require.True(t, ok, "no ipSecConfig in body: %v", body)
	return config
}

func ikeJSON(t *testing.T, body map[string]interface{}) map[string]interface{} {
	t.Helper()
	ike, ok := ipSecJSON(t, body)["ikeProtocol"].(map[string]interface{})
	require.True(t, ike != nil && ok, "no ikeProtocol in body")
	return ike
}

// connectorEchoHandler answers create, ipsec start and read with the supplied
// connector JSON.
func connectorEchoHandler(connectorJSON string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/networks/connectors/conn-1/ipsec/start" {
			_, _ = w.Write([]byte(`{}`))
			return
		}
		_, _ = w.Write([]byte(connectorJSON))
	})
}

const ipSecConnectorResponse = `{
  "id":"conn-1",
  "name":"test-connector",
  "networkItemId":"net-1",
  "networkItemType":"NETWORK",
  "vpnRegionId":"us-east-1",
  "tunnelingProtocol":"IPSEC",
  "ipSecConfig":{
    "platform":"AWS",
    "authenticationType":"SHARED_SECRET",
    "remoteSitePublicIp":"203.0.113.10",
    "preSharedKey":"psk-value",
    "caCertificate":"ca-cert",
    "peerCertificate":"peer-cert",
    "remoteGatewayCertificate":"gateway-cert",
    "peerCertificatePrivateKey":"peer-key",
    "peerCertificateKeyPassphrase":"peer-passphrase",
    "hostname":"connector-host",
    "domain":"example.com",
    "ikeProtocol":{
      "protocolVersion":"IKE_V2",
      "startupAction":"START",
      "autoInitiate":true,
      "deadPeerDetection":{"timeoutSec":30,"deadPeerHandling":"RESTART"},
      "phase1":{"encryptionAlgorithms":["AES256","AES128"],"integrityAlgorithms":["SHA2_256"],"diffieHellmanGroups":["G_14","G_19"],"lifetimeSec":28800},
      "phase2":{"encryptionAlgorithms":["AES128_GCM_16"],"integrityAlgorithms":["SHA2_512"],"diffieHellmanGroups":["G_20"],"lifetimeSec":3600},
      "rekey":{"marginTimeSec":270,"fuzzPercent":100,"replayWindowSize":1024}
    }
  }
}`

// TestUnitNetworkConnectorIpSecPayload_AutoInitiate asserts the create request carries
// autoInitiate and omits deadPeerDetection entirely. The API rejects a deadPeerDetection
// object whose fields are unset, so an empty one on the wire would be a hard failure.
func TestUnitNetworkConnectorIpSecPayload_AutoInitiate(t *testing.T) {
	for name, autoInitiate := range map[string]bool{"true": true, "false": false} {
		t.Run(name, func(t *testing.T) {
			c, captured := newNetworkConnectorTestClient(t, connectorEchoHandler(ipSecConnectorResponse))
			d := networkConnectorData(t, "", withAttrs(fullIpSecConfig(), map[string]cty.Value{
				"auto_initiate": cty.BoolVal(autoInitiate),
			}))

			diags := resourceNetworkConnectorCreate(context.Background(), d, c)
			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			ike := ikeJSON(t, captured.bodyFor(t, http.MethodPost, "/api/v1/networks/connectors"))
			assert.Equal(t, autoInitiate, ike["autoInitiate"])
			assert.NotContains(t, ike, "deadPeerDetection",
				"deadPeerDetection must be absent when auto_initiate is set")
		})
	}
}

// TestUnitNetworkConnectorIpSecPayload_DeadPeerDetection asserts a config still using the
// deprecated attributes sends deadPeerDetection with both of its required fields and no
// autoInitiate.
func TestUnitNetworkConnectorIpSecPayload_DeadPeerDetection(t *testing.T) {
	c, captured := newNetworkConnectorTestClient(t, connectorEchoHandler(ipSecConnectorResponse))
	d := networkConnectorData(t, "", withAttrs(fullIpSecConfig(), map[string]cty.Value{
		"timeout_sec":        cty.NumberIntVal(45),
		"dead_peer_handling": cty.StringVal("RESTART"),
	}))

	diags := resourceNetworkConnectorCreate(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	ike := ikeJSON(t, captured.bodyFor(t, http.MethodPost, "/api/v1/networks/connectors"))
	assert.NotContains(t, ike, "autoInitiate", "autoInitiate must not be invented")
	require.Contains(t, ike, "deadPeerDetection")
	deadPeerDetection := ike["deadPeerDetection"].(map[string]interface{})
	assert.Equal(t, float64(45), deadPeerDetection["timeoutSec"])
	assert.Equal(t, "RESTART", deadPeerDetection["deadPeerHandling"])
}

// TestUnitNetworkConnectorIpSecPayload_AllFields walks every ipsec_config attribute from
// configuration to the JSON key the API expects, so a mis-wired field is caught here
// rather than by a live request.
func TestUnitNetworkConnectorIpSecPayload_AllFields(t *testing.T) {
	c, captured := newNetworkConnectorTestClient(t, connectorEchoHandler(ipSecConnectorResponse))
	d := networkConnectorData(t, "", withAttrs(fullIpSecConfig(), map[string]cty.Value{
		"auto_initiate": cty.True,
	}))

	diags := resourceNetworkConnectorCreate(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	body := captured.bodyFor(t, http.MethodPost, "/api/v1/networks/connectors")
	assert.Equal(t, "test-connector", body["name"])
	assert.Equal(t, "net-1", body["networkItemId"])
	assert.Equal(t, "NETWORK", body["networkItemType"])
	assert.Equal(t, "us-east-1", body["vpnRegionId"])

	ipSec := ipSecJSON(t, body)
	assert.Equal(t, "AWS", ipSec["platform"])
	assert.Equal(t, "SHARED_SECRET", ipSec["authenticationType"])
	assert.Equal(t, "203.0.113.10", ipSec["remoteSitePublicIp"])
	assert.Equal(t, "psk-value", ipSec["preSharedKey"])
	assert.Equal(t, "ca-cert", ipSec["caCertificate"])
	assert.Equal(t, "peer-cert", ipSec["peerCertificate"])
	assert.Equal(t, "gateway-cert", ipSec["remoteGatewayCertificate"])
	assert.Equal(t, "peer-key", ipSec["peerCertificatePrivateKey"])
	assert.Equal(t, "peer-passphrase", ipSec["peerCertificateKeyPassphrase"])
	assert.Equal(t, "connector-host", ipSec["hostname"])
	assert.Equal(t, "example.com", ipSec["domain"])

	ike := ikeJSON(t, body)
	assert.Equal(t, "IKE_V2", ike["protocolVersion"])
	assert.Equal(t, "START", ike["startupAction"])

	phase1 := ike["phase1"].(map[string]interface{})
	assert.Equal(t, []interface{}{"AES256", "AES128"}, phase1["encryptionAlgorithms"])
	assert.Equal(t, []interface{}{"SHA2_256"}, phase1["integrityAlgorithms"])
	assert.Equal(t, []interface{}{"G_14", "G_19"}, phase1["diffieHellmanGroups"])
	assert.Equal(t, float64(28800), phase1["lifetimeSec"])

	phase2 := ike["phase2"].(map[string]interface{})
	assert.Equal(t, []interface{}{"AES128_GCM_16"}, phase2["encryptionAlgorithms"])
	assert.Equal(t, []interface{}{"SHA2_512"}, phase2["integrityAlgorithms"])
	assert.Equal(t, []interface{}{"G_20"}, phase2["diffieHellmanGroups"])
	assert.Equal(t, float64(3600), phase2["lifetimeSec"])

	rekey := ike["rekey"].(map[string]interface{})
	assert.Equal(t, float64(270), rekey["marginTimeSec"])
	assert.Equal(t, float64(100), rekey["fuzzPercent"])
	assert.Equal(t, float64(1024), rekey["replayWindowSize"])
}

// ---------------------------------------------------------------------------
// CustomizeDiff, through the real SDK plumbing
// ---------------------------------------------------------------------------

// TestUnitNetworkConnectorCustomizeDiff drives the mutual-exclusion rule through
// Resource.Diff so the wiring — CustomizeDiff registration and GetRawConfig — is
// covered, not just the rule function.
func TestUnitNetworkConnectorCustomizeDiff(t *testing.T) {
	tests := []struct {
		name      string
		attrs     map[string]cty.Value
		wantError string
	}{
		{
			name:  "auto_initiate alone",
			attrs: map[string]cty.Value{"auto_initiate": cty.True},
		},
		{
			name: "deprecated pair alone",
			attrs: map[string]cty.Value{
				"timeout_sec":        cty.NumberIntVal(30),
				"dead_peer_handling": cty.StringVal("RESTART"),
			},
		},
		{
			name: "both groups",
			attrs: map[string]cty.Value{
				"auto_initiate":      cty.True,
				"timeout_sec":        cty.NumberIntVal(30),
				"dead_peer_handling": cty.StringVal("RESTART"),
			},
			wantError: `ipsec_config: "auto_initiate" conflicts with the deprecated "timeout_sec" and "dead_peer_handling"; set "auto_initiate" alone`,
		},
		{
			name:      "neither group",
			attrs:     map[string]cty.Value{},
			wantError: `ipsec_config: "auto_initiate" must be set (or, deprecated, both "timeout_sec" and "dead_peer_handling")`,
		},
		{
			name:      "half the deprecated pair",
			attrs:     map[string]cty.Value{"dead_peer_handling": cty.StringVal("NONE")},
			wantError: `ipsec_config: the deprecated "timeout_sec" and "dead_peer_handling" must be set together`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := withAttrs(fullIpSecConfig(), tt.attrs)
			state := networkConnectorState(t, "conn-1", attrs)
			_, err := resourceNetworkConnector().Diff(context.Background(), state, networkConnectorConfig(attrs), nil)
			if tt.wantError == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantError)
		})
	}
}

// TestUnitValidateIpSecDeadPeerConfig_NoIpSecConfig covers the configurations the rule
// must stay out of: an OpenVPN connector with no ipsec_config block, and the null config
// Terraform sends when the resource is being destroyed.
func TestUnitValidateIpSecDeadPeerConfig_NoIpSecConfig(t *testing.T) {
	resourceType := schema.InternalMap(resourceNetworkConnector().Schema).CoreConfigSchema().ImpliedType()

	t.Run("empty ipsec_config list", func(t *testing.T) {
		values := make(map[string]cty.Value)
		for name, attrType := range resourceType.AttributeTypes() {
			values[name] = cty.NullVal(attrType)
		}
		values["ipsec_config"] = cty.ListValEmpty(resourceType.AttributeType("ipsec_config").ElementType())
		require.NoError(t, validateIpSecDeadPeerConfig(cty.ObjectVal(values)))
	})

	t.Run("null config", func(t *testing.T) {
		require.NoError(t, validateIpSecDeadPeerConfig(cty.NullVal(resourceType)))
	})

	t.Run("unset config", func(t *testing.T) {
		require.NoError(t, validateIpSecDeadPeerConfig(cty.NilVal))
	})

	t.Run("unknown config", func(t *testing.T) {
		require.NoError(t, validateIpSecDeadPeerConfig(cty.UnknownVal(resourceType)))
	})
}

// ---------------------------------------------------------------------------
// CRUD lifecycle
// ---------------------------------------------------------------------------

// TestUnitNetworkConnectorCreate_IpSec covers the create path for an IPsec connector:
// the ID is persisted, the tunnel is started, and no OpenVPN profile or token is
// fetched for a connector that is not OPENVPN.
func TestUnitNetworkConnectorCreate_IpSec(t *testing.T) {
	c, captured := newNetworkConnectorTestClient(t, connectorEchoHandler(ipSecConnectorResponse))
	d := networkConnectorData(t, "", withAttrs(fullIpSecConfig(), map[string]cty.Value{
		"auto_initiate": cty.True,
	}))

	diags := resourceNetworkConnectorCreate(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, "conn-1", d.Id())
	assert.True(t, captured.called(http.MethodPost, "/api/v1/networks/connectors/conn-1/ipsec/start"),
		"IPsec tunnel must be started; saw %v", captured.paths())
	assert.False(t, captured.called(http.MethodPost, "/api/v1/networks/connectors/conn-1/profile"),
		"no OpenVPN profile should be fetched for an IPsec connector")
}

// TestUnitNetworkConnectorCreate_Error covers the create failure branch: the resource
// must surface a diagnostic and keep an empty ID.
func TestUnitNetworkConnectorCreate_Error(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid ipsec config"}`))
	})
	c, _ := newNetworkConnectorTestClient(t, handler)
	d := networkConnectorData(t, "", withAttrs(fullIpSecConfig(), map[string]cty.Value{
		"auto_initiate": cty.True,
	}))

	diags := resourceNetworkConnectorCreate(context.Background(), d, c)
	assert.True(t, diags.HasError())
	assert.Empty(t, d.Id(), "ID must remain empty when Create fails")
}

// TestUnitNetworkConnectorCreate_StartIPsecError covers the branch where the connector
// is created but the tunnel fails to start.
func TestUnitNetworkConnectorCreate_StartIPsecError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/networks/connectors/conn-1/ipsec/start" {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"error":"cannot start"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(ipSecConnectorResponse))
	})
	c, _ := newNetworkConnectorTestClient(t, handler)
	d := networkConnectorData(t, "", withAttrs(fullIpSecConfig(), map[string]cty.Value{
		"auto_initiate": cty.True,
	}))

	diags := resourceNetworkConnectorCreate(context.Background(), d, c)
	assert.True(t, diags.HasError(), "a failed IPsec start must surface")
	assert.Equal(t, "conn-1", d.Id(), "the connector was created, so its ID must be kept")
}

// TestUnitNetworkConnectorRead_IpSec covers the read path: every ipsec_config attribute
// the API reports lands in state, including the derived dead peer detection values that
// accompany autoInitiate.
func TestUnitNetworkConnectorRead_IpSec(t *testing.T) {
	c, captured := newNetworkConnectorTestClient(t, connectorEchoHandler(ipSecConnectorResponse))
	d := networkConnectorData(t, "conn-1", withAttrs(fullIpSecConfig(), map[string]cty.Value{
		"auto_initiate": cty.True,
	}))

	diags := resourceNetworkConnectorRead(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, "AWS", d.Get("ipsec_config.0.platform"))
	assert.Equal(t, "SHARED_SECRET", d.Get("ipsec_config.0.authentication_type"))
	assert.Equal(t, "203.0.113.10", d.Get("ipsec_config.0.remote_site_public_ip"))
	assert.Equal(t, "psk-value", d.Get("ipsec_config.0.pre_shared_key"))
	assert.Equal(t, "ca-cert", d.Get("ipsec_config.0.ca_certificate"))
	assert.Equal(t, "peer-cert", d.Get("ipsec_config.0.peer_certificate"))
	assert.Equal(t, "gateway-cert", d.Get("ipsec_config.0.remote_gateway_certificate"))
	assert.Equal(t, "peer-key", d.Get("ipsec_config.0.peer_certificate_private_key"))
	assert.Equal(t, "peer-passphrase", d.Get("ipsec_config.0.peer_certificate_key_passphrase"))
	assert.Equal(t, "IKE_V2", d.Get("ipsec_config.0.protocol_version"))
	assert.Equal(t, "START", d.Get("ipsec_config.0.startup_action"))
	assert.Equal(t, "connector-host", d.Get("ipsec_config.0.hostname"))
	assert.Equal(t, "example.com", d.Get("ipsec_config.0.domain"))

	assert.Equal(t, []interface{}{"AES256", "AES128"}, d.Get("ipsec_config.0.phase_1_encryption_algorithms"))
	assert.Equal(t, []interface{}{"SHA2_256"}, d.Get("ipsec_config.0.phase_1_integrity_algorithms"))
	assert.Equal(t, []interface{}{"G_14", "G_19"}, d.Get("ipsec_config.0.phase_1_diffie_hellman_groups"))
	assert.Equal(t, 28800, d.Get("ipsec_config.0.phase_1_lifetime_sec"))
	assert.Equal(t, []interface{}{"AES128_GCM_16"}, d.Get("ipsec_config.0.phase_2_encryption_algorithms"))
	assert.Equal(t, []interface{}{"SHA2_512"}, d.Get("ipsec_config.0.phase_2_integrity_algorithms"))
	assert.Equal(t, []interface{}{"G_20"}, d.Get("ipsec_config.0.phase_2_diffie_hellman_groups"))
	assert.Equal(t, 3600, d.Get("ipsec_config.0.phase_2_lifetime_sec"))
	assert.Equal(t, 270, d.Get("ipsec_config.0.margin_time_sec"))
	assert.Equal(t, 100, d.Get("ipsec_config.0.fuzz_percent"))
	assert.Equal(t, 1024, d.Get("ipsec_config.0.replay_window_size"))

	assert.Equal(t, true, d.Get("ipsec_config.0.auto_initiate"))
	assert.Equal(t, 30, d.Get("ipsec_config.0.timeout_sec"))
	assert.Equal(t, "RESTART", d.Get("ipsec_config.0.dead_peer_handling"))

	assert.False(t, captured.called(http.MethodPost, "/api/v1/networks/connectors/conn-1/profile"),
		"no OpenVPN profile should be fetched for an IPsec connector")
}

// TestUnitNetworkConnectorRead_Error covers the read failure branch.
func TestUnitNetworkConnectorRead_Error(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	})
	c, _ := newNetworkConnectorTestClient(t, handler)
	d := networkConnectorData(t, "conn-1", withAttrs(fullIpSecConfig(), map[string]cty.Value{
		"auto_initiate": cty.True,
	}))

	diags := resourceNetworkConnectorRead(context.Background(), d, c)
	assert.True(t, diags.HasError())
}

// TestUnitNetworkConnectorUpdate_IpSec covers the update path: the changed ipsec_config
// is sent as a PUT and the tunnel is restarted.
func TestUnitNetworkConnectorUpdate_IpSec(t *testing.T) {
	c, captured := newNetworkConnectorTestClient(t, connectorEchoHandler(ipSecConnectorResponse))
	d := networkConnectorUpdateData(t, "conn-1",
		withAttrs(fullIpSecConfig(), map[string]cty.Value{
			"timeout_sec":        cty.NumberIntVal(45),
			"dead_peer_handling": cty.StringVal("RESTART"),
		}),
		withAttrs(fullIpSecConfig(), map[string]cty.Value{
			"auto_initiate": cty.True,
		}),
	)

	diags := resourceNetworkConnectorUpdate(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	ike := ikeJSON(t, captured.bodyFor(t, http.MethodPut, "/api/v1/networks/connectors/conn-1"))
	assert.Equal(t, true, ike["autoInitiate"])
	assert.NotContains(t, ike, "deadPeerDetection",
		"migrating onto auto_initiate must stop sending the deprecated values")
	assert.True(t, captured.called(http.MethodPost, "/api/v1/networks/connectors/conn-1/ipsec/start"),
		"IPsec tunnel must be restarted after an update; saw %v", captured.paths())
}

// TestUnitNetworkConnectorUpdate_Error covers the update failure branch.
func TestUnitNetworkConnectorUpdate_Error(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(ipSecConnectorResponse))
	})
	c, _ := newNetworkConnectorTestClient(t, handler)
	d := networkConnectorUpdateData(t, "conn-1",
		withAttrs(fullIpSecConfig(), map[string]cty.Value{"auto_initiate": cty.False}),
		withAttrs(fullIpSecConfig(), map[string]cty.Value{"auto_initiate": cty.True}),
	)

	diags := resourceNetworkConnectorUpdate(context.Background(), d, c)
	assert.True(t, diags.HasError())
}

// TestUnitNetworkConnectorDelete covers the delete path and its error branch.
func TestUnitNetworkConnectorDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))
		})
		c, captured := newNetworkConnectorTestClient(t, handler)
		d := networkConnectorData(t, "conn-1", withAttrs(fullIpSecConfig(), map[string]cty.Value{
			"auto_initiate": cty.True,
		}))

		diags := resourceNetworkConnectorDelete(context.Background(), d, c)
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
		assert.True(t, captured.called(http.MethodDelete, "/api/v1/networks/connectors/conn-1"))
	})

	t.Run("error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"boom"}`))
		})
		c, _ := newNetworkConnectorTestClient(t, handler)
		d := networkConnectorData(t, "conn-1", withAttrs(fullIpSecConfig(), map[string]cty.Value{
			"auto_initiate": cty.True,
		}))

		diags := resourceNetworkConnectorDelete(context.Background(), d, c)
		assert.True(t, diags.HasError())
	})
}

// TestUnitDataSourceNetworkConnectorRead_IpSec covers the data source, which shares
// ipSecConfigSchema with the resource and so must surface auto_initiate too.
func TestUnitDataSourceNetworkConnectorRead_IpSec(t *testing.T) {
	c, _ := newNetworkConnectorTestClient(t, connectorEchoHandler(ipSecConnectorResponse))
	d := schema.TestResourceDataRaw(t, dataSourceNetworkConnector().Schema, map[string]interface{}{
		"id": "conn-1",
	})

	diags := dataSourceNetworkConnectorRead(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, "AWS", d.Get("ipsec_config.0.platform"))
	assert.Equal(t, true, d.Get("ipsec_config.0.auto_initiate"))
	assert.Equal(t, "RESTART", d.Get("ipsec_config.0.dead_peer_handling"))
	assert.Equal(t, 30, d.Get("ipsec_config.0.timeout_sec"))
}

// ---------------------------------------------------------------------------
// Non-IPsec branches of the same CRUD functions
// ---------------------------------------------------------------------------

const openVPNConnectorResponse = `{
  "id":"conn-1",
  "name":"test-connector",
  "networkItemId":"net-1",
  "networkItemType":"NETWORK",
  "vpnRegionId":"us-east-1",
  "tunnelingProtocol":"OPENVPN"
}`

// openVPNHandler answers the connector routes plus the profile and token endpoints.
func openVPNHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/networks/connectors/conn-1/profile":
			_, _ = w.Write([]byte("openvpn-profile"))
		case "/api/v1/networks/connectors/conn-1/profile/encrypt":
			_, _ = w.Write([]byte("connector-token"))
		default:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(openVPNConnectorResponse))
		}
	})
}

// TestUnitNetworkConnectorCreate_OpenVPN covers the branch an IPsec connector skips:
// an OPENVPN connector fetches its profile and token and starts no tunnel.
func TestUnitNetworkConnectorCreate_OpenVPN(t *testing.T) {
	c, captured := newNetworkConnectorTestClient(t, openVPNHandler())
	d := schema.TestResourceDataRaw(t, resourceNetworkConnector().Schema, map[string]interface{}{
		"name":          "test-connector",
		"network_id":    "net-1",
		"vpn_region_id": "us-east-1",
	})

	diags := resourceNetworkConnectorCreate(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, "conn-1", d.Id())
	assert.Equal(t, "openvpn-profile", d.Get("profile"))
	assert.Equal(t, "connector-token", d.Get("token"))
	assert.False(t, captured.called(http.MethodPost, "/api/v1/networks/connectors/conn-1/ipsec/start"),
		"no IPsec tunnel should be started for an OPENVPN connector")
}

// TestUnitNetworkConnectorRead_OpenVPN covers the read path of an OPENVPN connector: the profile
// is refreshed, but no token is minted (each token request creates a new one on the API side).
func TestUnitNetworkConnectorRead_OpenVPN(t *testing.T) {
	c, captured := newNetworkConnectorTestClient(t, openVPNHandler())
	d := schema.TestResourceDataRaw(t, resourceNetworkConnector().Schema, map[string]interface{}{
		"name":  "test-connector",
		"token": "token-from-create",
	})
	d.SetId("conn-1")

	diags := resourceNetworkConnectorRead(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, "openvpn-profile", d.Get("profile"))
	assert.Equal(t, "token-from-create", d.Get("token"), "token must be left as set on create")
	assert.True(t, captured.called(http.MethodPost, "/api/v1/networks/connectors/conn-1/profile"))
	assert.False(t, captured.called(http.MethodPost, "/api/v1/networks/connectors/conn-1/profile/encrypt"),
		"read must not mint a new token")
	assert.Empty(t, d.Get("ipsec_config"), "an OPENVPN connector has no ipsec_config")
}

// hostConnectorHandler answers the host connector routes plus the profile and token endpoints.
func hostConnectorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/hosts/connectors/conn-9/profile":
			_, _ = w.Write([]byte("host-profile"))
		case "/api/v1/hosts/connectors/conn-9/profile/encrypt":
			_, _ = w.Write([]byte("host-token"))
		default:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"conn-9","name":"host-connector","networkItemId":"host-1","networkItemType":"HOST","vpnRegionId":"us-east-1"}`))
		}
	})
}

// TestUnitHostConnectorRead_DoesNotMintToken covers the same rule for host connectors.
func TestUnitHostConnectorRead_DoesNotMintToken(t *testing.T) {
	c, captured := newNetworkConnectorTestClient(t, hostConnectorHandler())
	d := schema.TestResourceDataRaw(t, resourceHostConnector().Schema, map[string]interface{}{
		"name":  "host-connector",
		"token": "token-from-create",
	})
	d.SetId("conn-9")

	diags := resourceHostConnectorRead(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, "host-profile", d.Get("profile"))
	assert.Equal(t, "token-from-create", d.Get("token"))
	assert.False(t, captured.called(http.MethodPost, "/api/v1/hosts/connectors/conn-9/profile/encrypt"),
		"read must not mint a new token")
}

// TestUnitDataSourceNetworkConnectorRead_OpenVPN covers the opposite rule for the data source:
// it keeps no state to hold a token in, so every read mints one next to the profile.
func TestUnitDataSourceNetworkConnectorRead_OpenVPN(t *testing.T) {
	c, _ := newNetworkConnectorTestClient(t, openVPNHandler())
	d := schema.TestResourceDataRaw(t, dataSourceNetworkConnector().Schema, map[string]interface{}{
		"id": "conn-1",
	})

	diags := dataSourceNetworkConnectorRead(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, "openvpn-profile", d.Get("profile"))
	assert.Equal(t, "connector-token", d.Get("token"))
}

// TestUnitDataSourceHostConnectorRead covers the same for the host connector data source.
func TestUnitDataSourceHostConnectorRead(t *testing.T) {
	c, _ := newNetworkConnectorTestClient(t, hostConnectorHandler())
	d := schema.TestResourceDataRaw(t, dataSourceHostConnector().Schema, map[string]interface{}{
		"id": "conn-9",
	})

	diags := dataSourceHostConnectorRead(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, "host-profile", d.Get("profile"))
	assert.Equal(t, "host-token", d.Get("token"))
}

// TestUnitNetworkConnectorUpdate_Status covers the suspend and activate branch, which
// runs ahead of the field comparison and is independent of ipsec_config.
func TestUnitNetworkConnectorUpdate_Status(t *testing.T) {
	tests := map[string]string{
		"SUSPENDED": "/api/v1/networks/connectors/conn-1/suspend",
		"ACTIVE":    "/api/v1/networks/connectors/conn-1/activate",
	}

	for status, wantPath := range tests {
		t.Run(status, func(t *testing.T) {
			c, captured := newNetworkConnectorTestClient(t, openVPNHandler())

			resource := resourceNetworkConnector()
			state := networkConnectorState(t, "conn-1", withAttrs(fullIpSecConfig(), map[string]cty.Value{
				"auto_initiate": cty.True,
			}))
			state.Attributes["status"] = "PLACEHOLDER"

			diff, err := resource.Diff(context.Background(), state, networkConnectorConfig(withAttrs(
				fullIpSecConfig(), map[string]cty.Value{"auto_initiate": cty.True},
			)), nil)
			require.NoError(t, err)
			require.NotNil(t, diff)
			diff.Attributes["status"] = &terraform.ResourceAttrDiff{Old: "PLACEHOLDER", New: status}

			d, err := schema.InternalMap(resource.Schema).Data(state, diff)
			require.NoError(t, err)

			diags := resourceNetworkConnectorUpdate(context.Background(), d, c)
			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
			assert.True(t, captured.called(http.MethodPut, wantPath),
				"expected a call to %s; saw %v", wantPath, captured.paths())
		})
	}
}

// TestUnitNetworkConnectorUpdate_StatusError covers a failing suspend.
func TestUnitNetworkConnectorUpdate_StatusError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/networks/connectors/conn-1/suspend" {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"error":"cannot suspend"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(ipSecConnectorResponse))
	})
	c, _ := newNetworkConnectorTestClient(t, handler)

	resource := resourceNetworkConnector()
	state := networkConnectorState(t, "conn-1", withAttrs(fullIpSecConfig(), map[string]cty.Value{
		"auto_initiate": cty.True,
	}))
	// A status differing from the schema default is what produces a diff to attach to.
	state.Attributes["status"] = "PLACEHOLDER"

	diff, err := resource.Diff(context.Background(), state, networkConnectorConfig(withAttrs(
		fullIpSecConfig(), map[string]cty.Value{"auto_initiate": cty.True},
	)), nil)
	require.NoError(t, err)
	require.NotNil(t, diff)
	diff.Attributes["status"] = &terraform.ResourceAttrDiff{Old: "PLACEHOLDER", New: "SUSPENDED"}

	d, err := schema.InternalMap(resource.Schema).Data(state, diff)
	require.NoError(t, err)

	diags := resourceNetworkConnectorUpdate(context.Background(), d, c)
	assert.True(t, diags.HasError())
}
