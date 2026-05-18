package cloudconnexa

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccCloudConnexaService_basic tests the basic creation and configuration of a CloudConnexa IP service.
// It verifies that the service can be created and updated with different names.
//
// Parameters:
//   - t: The testing context
func TestAccCloudConnexaService_basic(t *testing.T) {
	rn := "cloudconnexa_host_ip_service.test"
	hostName := acctest.RandStringFromCharSet(10, alphabet)
	service := cloudconnexa.IPService{
		Name: acctest.RandStringFromCharSet(10, alphabet),
	}
	serviceChanged := service
	serviceChanged.Name = fmt.Sprintf("changed-%s", acctest.RandStringFromCharSet(10, alphabet))

	check := func(service cloudconnexa.IPService) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckCloudConnexaServiceExists(rn, hostName),
			resource.TestCheckResourceAttr(rn, "name", service.Name),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaServiceConfig(service, hostName),
				Check:  check(service),
			},
			{
				Config: testAccCloudConnexaServiceConfig(serviceChanged, hostName),
				Check:  check(serviceChanged),
			},
		},
	})
}

// testAccCheckCloudConnexaServiceExists verifies that a CloudConnexa IP service exists in the Terraform state.
//
// Parameters:
//   - rn: The name of the resource to check
//   - networkId: The ID of the network associated with the service
//
// Returns:
//   - resource.TestCheckFunc: A function that performs the existence check
func testAccCheckCloudConnexaServiceExists(rn, _ string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		c := testAccProvider.Meta().(*providerMeta).Client
		_, err := c.HostIPServices.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

// testAccCheckCloudConnexaServiceDestroy verifies that a CloudConnexa IP service has been properly destroyed.
// It checks that the service no longer exists in the CloudConnexa API.
//
// Parameters:
//   - state: The Terraform state to check
//
// Returns:
//   - error: An error if the service still exists or if there was an error checking its existence
func testAccCheckCloudConnexaServiceDestroy(state *terraform.State) error {
	c := testAccProvider.Meta().(*providerMeta).Client
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "cloudconnexa_host_ip_service" {
			continue
		}
		id := rs.Primary.Attributes["id"]
		s, err := c.HostIPServices.Get(id)
		if err == nil || s != nil {
			return fmt.Errorf("service still exists")
		}
	}
	return nil
}

// testAccCloudConnexaServiceConfig generates a Terraform configuration for testing CloudConnexa IP services.
//
// Parameters:
//   - service: The IP service configuration to use
//   - networkName: The name of the network to associate with the service
//
// Returns:
//   - string: The generated Terraform configuration
func testAccCloudConnexaServiceConfig(service cloudconnexa.IPService, hostName string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
	base_url = "%s"
}

resource "cloudconnexa_host" "test" {
	name = "%s"
	description = "test"
}

resource "cloudconnexa_host_ip_service" "test" {
	name = "%s"
	description = "test"
	host_id = cloudconnexa_host.test.id
	config {
		service_types = ["ANY"]
	}
}
`, testBaseURL, hostName, service.Name)
}
