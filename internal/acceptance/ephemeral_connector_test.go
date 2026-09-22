// Package acceptance holds acceptance tests that need terraform-plugin-testing, such as the
// ephemeral resource tests that observe values through the echo provider. They live outside
// package cloudconnexa because terraform-plugin-testing/helper/resource and the SDK v2
// helper/resource both register a -sweep flag and cannot share one test binary.
//
// Credentials are read from CLOUDCONNEXA_CLIENT_ID, CLOUDCONNEXA_CLIENT_SECRET and
// CLOUDCONNEXA_BASE_URL; this package does not load the .env file.
package acceptance

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/OpenVPN/terraform-provider-cloudconnexa/cloudconnexa"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	api "github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// baseURLEnvVar mirrors the constant of the same name in package cloudconnexa's tests.
const baseURLEnvVar = "CLOUDCONNEXA_BASE_URL"

// protoV5ProviderFactories serves the provider exactly as main.go does: SDK v2 and framework
// implementations behind one mux.
var protoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	"cloudconnexa": func() (tfprotov5.ProviderServer, error) {
		mux, err := tf5muxserver.NewMuxServer(context.Background(),
			cloudconnexa.Provider().GRPCProvider,
			providerserver.NewProtocol5(cloudconnexa.NewFrameworkProvider()),
		)
		if err != nil {
			return nil, err
		}
		return mux.ProviderServer(), nil
	},
}

func preCheck(t *testing.T) {
	t.Helper()
	for _, v := range []string{cloudconnexa.ClientIDEnvVar, cloudconnexa.ClientSecretEnvVar, baseURLEnvVar} {
		if os.Getenv(v) == "" {
			t.Fatalf("%s must be set for acceptance tests", v)
		}
	}
}

// TestAccNetworkConnectorToken_ephemeral mints a token for a new network connector.
func TestAccNetworkConnectorToken_ephemeral(t *testing.T) {
	// key (44 chars base64) + id (40 chars sha1), see profile-extension ScriptController.encryptProfile
	testAccNetworkConnectorEphemeral(t, "token", regexp.MustCompile(`^[A-Za-z0-9+/=]{44}[0-9a-f]{40}$`))
}

// TestAccNetworkConnectorProfile_ephemeral retrieves the OpenVPN profile of a new network connector.
func TestAccNetworkConnectorProfile_ephemeral(t *testing.T) {
	testAccNetworkConnectorEphemeral(t, "profile", regexp.MustCompile(`(?s)<ca>.+</ca>`))
}

// testAccNetworkConnectorEphemeral creates a network and an OpenVPN connector through the SDK v2
// side of the mux, opens the cloudconnexa_network_connector_<attribute> ephemeral resource through
// the framework side, and observes its <attribute> value with the echo test provider, since it
// never appears in plan or state.
func testAccNetworkConnectorEphemeral(t *testing.T, attribute string, want *regexp.Regexp) {
	rName := acctest.RandomWithPrefix("test-" + attribute)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { preCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		CheckDestroy:             checkNetworkConnectorDestroy,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: networkConnectorEphemeralConfig(rName, attribute),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("cloudconnexa_network_connector.test", tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data"), knownvalue.StringRegexp(want)),
				},
			},
		},
	})
}

// checkNetworkConnectorDestroy confirms the API no longer returns the network connector and
// network the test created. A lookup of a deleted one fails with 400 "... not found", not 404.
func checkNetworkConnectorDestroy(s *terraform.State) error {
	client, err := api.NewClient(os.Getenv(baseURLEnvVar), os.Getenv(cloudconnexa.ClientIDEnvVar), os.Getenv(cloudconnexa.ClientSecretEnvVar))
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		var exists bool
		switch rs.Type {
		case "cloudconnexa_network_connector":
			c, err := client.NetworkConnectors.GetByID(rs.Primary.ID)
			if err != nil && !strings.Contains(err.Error(), "not found") {
				return err
			}
			exists = c != nil && err == nil
		case "cloudconnexa_network":
			n, err := client.Networks.Get(rs.Primary.ID)
			if err != nil && !strings.Contains(err.Error(), "not found") {
				return err
			}
			exists = n != nil && err == nil
		default:
			continue
		}
		if exists {
			return fmt.Errorf("%s with ID %q still exists", rs.Type, rs.Primary.ID)
		}
	}
	return nil
}

func networkConnectorEphemeralConfig(rName, attribute string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
  base_url = "%[1]s"
}

resource "cloudconnexa_network" "test" {
  name = "%[2]s"
}

resource "cloudconnexa_network_connector" "test" {
  name          = "%[2]s"
  network_id    = cloudconnexa_network.test.id
  vpn_region_id = "us-east-1"
}

ephemeral "cloudconnexa_network_connector_%[3]s" "test" {
  connector_id = cloudconnexa_network_connector.test.id
}

provider "echo" {
  data = ephemeral.cloudconnexa_network_connector_%[3]s.test.%[3]s
}

resource "echo" "test" {}
`, os.Getenv(baseURLEnvVar), rName, attribute)
}
