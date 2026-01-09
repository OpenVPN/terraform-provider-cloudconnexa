package cloudconnexa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccCloudConnexaDataSourceHostRoutes_basic tests the basic functionality of the host routes data source.
func TestAccCloudConnexaDataSourceHostRoutes_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaDataSourceHostRoutesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudconnexa_host_routes.test", "host_routes.#"),
				),
			},
		},
	})
}

// testAccCloudConnexaDataSourceHostRoutesConfig generates a Terraform configuration for host routes data source.
func testAccCloudConnexaDataSourceHostRoutesConfig() string {
	return `
provider "cloudconnexa" {
	base_url = "https://` + testCloudID + `.api.openvpn.com"
}

data "cloudconnexa_host_routes" "test" {
}
`
}
