package cloudconnexa

import (
	"errors"
	"fmt"
	"testing"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccCloudConnexaHostRoute_basic tests the basic functionality of creating a host route resource.
func TestAccCloudConnexaHostRoute_basic(t *testing.T) {
	rn := "cloudconnexa_host_route.test"
	hostName := acctest.RandStringFromCharSet(10, alphabet)
	subnet := "10.0.0.0/24"
	var routeID string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaHostRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaHostRouteConfig(hostName, subnet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudConnexaHostRouteExists(rn, &routeID),
					resource.TestCheckResourceAttr(rn, "subnet", subnet),
					resource.TestCheckResourceAttr(rn, "description", "Managed by Terraform"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCheckCloudConnexaHostRouteDestroy verifies that the host route has been properly destroyed.
func testAccCheckCloudConnexaHostRouteDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudconnexa.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_host_route" {
			continue
		}
		route, err := client.HostRoutes.GetByID(rs.Primary.ID)
		if err == nil && route != nil {
			return errors.New("host route still exists")
		}
	}
	return nil
}

// testAccCheckCloudConnexaHostRouteExists verifies that a host route exists.
func testAccCheckCloudConnexaHostRouteExists(n string, routeID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		client := testAccProvider.Meta().(*cloudconnexa.Client)
		_, err := client.HostRoutes.GetByID(rs.Primary.ID)
		if err != nil {
			return err
		}
		*routeID = rs.Primary.ID
		return nil
	}
}

// testAccCloudConnexaHostRouteConfig generates a Terraform configuration for a host route resource.
func testAccCloudConnexaHostRouteConfig(hostName, subnet string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
	base_url = "https://%s.api.openvpn.com"
}

resource "cloudconnexa_host" "test" {
	name           = "%s"
	description    = "Test host for route"
	internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_host_route" "test" {
	host_id     = cloudconnexa_host.test.id
	subnet      = "%s"
	description = "Managed by Terraform"
}
`, testCloudID, hostName, subnet)
}
