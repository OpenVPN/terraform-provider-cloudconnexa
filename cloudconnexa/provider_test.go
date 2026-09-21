package cloudconnexa

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// alphabet is a constant string containing all lowercase letters of the English alphabet
const alphabet = "abcdefghigklmnopqrstuvwxyz"

// BaseURLEnvVar is the environment variable holding the full CloudConnexa API base URL
// (e.g. https://example.api.openvpn.com) used by the acceptance test provider config.
const BaseURLEnvVar = "CLOUDCONNEXA_BASE_URL"

// testBaseURL stores the full CloudConnexa base URL used by the acceptance tests.
var testBaseURL = os.Getenv(BaseURLEnvVar)

// testAccProvider holds the Terraform provider instance for testing
var testAccProvider *schema.Provider

// testAccProviderFactories maps provider names to factory functions that create provider instances
var testAccProviderFactories map[string]func() (*schema.Provider, error)

// init initializes the test provider and its factory function
func init() {
	testAccProvider = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"cloudconnexa": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

// TestProvider validates the provider configuration and tests error handling for missing credentials
//
// Parameters:
//   - t: The testing context
func TestProvider(t *testing.T) {
	err := Provider().InternalValidate()
	require.NoError(t, err)

	// must have the required error when the credentials are not set
	t.Setenv(ClientIDEnvVar, "")
	t.Setenv(ClientSecretEnvVar, "")
	rc := terraform.ResourceConfig{}
	diags := Provider().Configure(context.Background(), &rc)
	assert.True(t, diags.HasError())

	for _, d := range diags {
		assert.Truef(t, strings.Contains(d.Detail, cloudconnexa.ErrCredentialsRequired.Error()),
			"error message does not contain the expected error: %s", d.Detail)
	}
}

// testAccPreCheck verifies that required environment variables are set before running acceptance tests
//
// Parameters:
//   - t: The testing context
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(ClientIDEnvVar); v == "" {
		t.Fatalf("%s must be set for acceptance tests", ClientIDEnvVar)
	}
	if v := os.Getenv(ClientSecretEnvVar); v == "" {
		t.Fatalf("%s must be set for acceptance tests", ClientSecretEnvVar)
	}
	if testBaseURL == "" {
		t.Fatalf("%s must be set for acceptance tests (full URL, e.g. https://example.api.openvpn.com)", BaseURLEnvVar)
	}
}

// TestUnitRetryLogger_WritesWarning verifies that the OnRetry hook emits one warn-level log
// entry carrying the request method and path, the retry number and the wait.
func TestUnitRetryLogger_WritesWarning(t *testing.T) {
	var buf bytes.Buffer
	ctx := tflogtest.RootLogger(context.Background(), &buf)

	req, err := http.NewRequest(http.MethodPost, "https://example.api.openvpn.com/api/v1/networks", nil)
	if err != nil {
		t.Fatal(err)
	}
	newRetryLogger(ctx)(req, 2, 4*time.Second)

	entries, err := tflogtest.MultilineJSONDecode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d: %v", len(entries), entries)
	}
	got := entries[0]
	want := map[string]interface{}{
		"@level":   "warn",
		"@message": "CloudConnexa API rate limit reached, retrying request",
		"method":   "POST",
		"path":     "/api/v1/networks",
		"attempt":  float64(2),
		"wait":     "4s",
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %v, want %v", k, got[k], v)
		}
	}
}
