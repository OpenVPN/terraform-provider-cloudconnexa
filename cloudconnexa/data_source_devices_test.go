package cloudconnexa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccCloudConnexaDataSourceDevices_basic tests the basic functionality of the devices data source.
func TestAccCloudConnexaDataSourceDevices_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaDataSourceDevicesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudconnexa_devices.test", "devices.#"),
				),
			},
		},
	})
}

// testAccCloudConnexaDataSourceDevicesConfig generates a Terraform configuration for devices data source.
func testAccCloudConnexaDataSourceDevicesConfig() string {
	return `
provider "cloudconnexa" {
	base_url = "https://` + testCloudID + `.api.openvpn.com"
}

data "cloudconnexa_devices" "test" {
}
`
}
