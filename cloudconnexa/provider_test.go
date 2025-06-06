package cloudconnexa

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// alphabet is a constant string containing all lowercase letters of the English alphabet
const alphabet = "abcdefghigklmnopqrstuvwxyz"

// testCloudID stores the CloudConnexa test organization ID from environment variables
var testCloudID = os.Getenv("CLOUDCONNEXA_TEST_ORGANIZATION")

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
}
