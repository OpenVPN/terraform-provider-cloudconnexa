package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/OpenVPN/terraform-provider-cloudconnexa/cloudconnexa"
	"github.com/gruntwork-io/terratest/modules/terraform"
	api "github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CloudConnexaHostKey is the environment variable key for the CloudConnexa host URL
const (
	CloudConnexaHostKey = "OVPN_HOST"
)

// TestCreationDeletion tests the creation and deletion of CloudConnexa resources.
// It verifies that a connector can be created and becomes online within the expected timeframe.
//
// Parameters:
//   - t: The testing context
func TestCreationDeletion(t *testing.T) {
	validateEnvVars(t)

	terraformOptions := &terraform.Options{
		NoColor: os.Getenv("NO_COLOR") == "1",

		// The path to where our Terraform code is located
		TerraformDir: "./setup",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	t.Cleanup(func() {
		terraform.Destroy(t, terraformOptions)
	})

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	hostID := terraform.Output(t, terraformOptions, "host_id")
	connectorID := terraform.Output(t, terraformOptions, "connector_id")

	assert.NotEmpty(t, hostID)
	assert.NotEmpty(t, connectorID)

	client, err := api.NewClient(
		os.Getenv(CloudConnexaHostKey),
		os.Getenv(cloudconnexa.ClientIDEnvVar),
		os.Getenv(cloudconnexa.ClientSecretEnvVar),
	)
	require.NoError(t, err)

	// Total waiting time: 1min
	totalAttempts := 10
	attemptWaitingTime := 6 * time.Second

	connectorWasOnline := false
	for i := 0; i < totalAttempts; i++ {
		t.Logf("Waiting for connector to be online (%d/%d)", i+1, totalAttempts)
		connector, err := client.NetworkConnectors.GetByID(connectorID)
		require.NoError(t, err, "Invalid connector ID in output")
		if connector.ConnectionStatus == "online" {
			connectorWasOnline = true
			break
		}
		time.Sleep(attemptWaitingTime)
	}
	assert.True(t, connectorWasOnline)
}

// validateEnvVars checks that all required environment variables are set.
// It validates the CloudConnexa host, client ID, and client secret environment variables.
//
// Parameters:
//   - t: The testing context
func validateEnvVars(t *testing.T) {
	validateEnvVar(t, CloudConnexaHostKey)
	validateEnvVar(t, cloudconnexa.ClientIDEnvVar)
	validateEnvVar(t, cloudconnexa.ClientSecretEnvVar)
}

// validateEnvVar checks if a specific environment variable is set.
// It fails the test if the environment variable is empty.
//
// Parameters:
//   - t: The testing context
//   - envVar: The name of the environment variable to validate
func validateEnvVar(t *testing.T, envVar string) {
	fmt.Println(os.Getenv(envVar))
	require.NotEmptyf(t, os.Getenv(envVar), "%s must be set for acceptance tests", envVar)
}
