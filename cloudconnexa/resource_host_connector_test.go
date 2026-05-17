package cloudconnexa

import (
	"fmt"
	"strings"
	"testing"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccCloudConnexaConnector_basic tests the basic creation and configuration of a CloudConnexa connector.
// It verifies that all required attributes are set correctly after creation.
//
// Parameters:
//   - t: The testing context
func TestAccCloudConnexaConnector_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("test-connector")
	resourceName := "cloudconnexa_host_connector.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaConnectorConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudConnexaConnectorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "host_id"),
				),
			},
		},
	})
}

// testAccCheckCloudConnexaConnectorExists verifies that a CloudConnexa connector exists in the Terraform state.
//
// Parameters:
//   - n: The name of the resource to check
//
// Returns:
//   - resource.TestCheckFunc: A function that performs the existence check
func testAccCheckCloudConnexaConnectorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No connector ID is set")
		}
		return nil
	}
}

// testAccCheckCloudConnexaConnectorDestroy verifies that a CloudConnexa connector has been properly destroyed.
// It checks that the connector no longer exists in the CloudConnexa API.
//
// Parameters:
//   - s: The Terraform state to check
//
// Returns:
//   - error: An error if the connector still exists or if there was an error checking its existence
func testAccCheckCloudConnexaConnectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudconnexa.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_host_connector" {
			continue
		}

		connectorId := rs.Primary.ID
		connector, err := client.HostConnectors.GetByID(connectorId)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				continue
			}
			return err
		}

		if connector != nil {
			return fmt.Errorf("connector with ID '%s' still exists", connectorId)
		}
	}

	return nil
}

// testAccCloudConnexaConnectorConfigBasic generates a basic Terraform configuration for a CloudConnexa connector.
//
// Parameters:
//   - rName: The name to use for the connector resource
//
// Returns:
//   - string: The Terraform configuration as a string
func testAccCloudConnexaConnectorConfigBasic(rName string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
  base_url = "%[1]s"
}

resource "cloudconnexa_host" "test" {
  name        = "%[2]s-host"
  description = "test"
}

resource "cloudconnexa_host_connector" "test" {
  name          = "%[2]s"
  vpn_region_id = "us-east-1"
  host_id       = cloudconnexa_host.test.id
}
`, testBaseURL, rName)
}
