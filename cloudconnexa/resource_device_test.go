package cloudconnexa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccCloudConnexaDevice_basic tests the basic functionality of the device resource.
// Note: This test requires at least one existing device in the CloudConnexa account.
func TestAccCloudConnexaDevice_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaDeviceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("cloudconnexa_device.test", "id"),
					resource.TestCheckResourceAttrSet("cloudconnexa_device.test", "name"),
					resource.TestCheckResourceAttrSet("cloudconnexa_device.test", "device_id"),
				),
			},
		},
	})
}

// testAccCloudConnexaDeviceConfig generates a Terraform configuration for device resource.
// This configuration first fetches all devices and then manages the first one found.
func testAccCloudConnexaDeviceConfig() string {
	return `
provider "cloudconnexa" {
	base_url = "https://` + testCloudID + `.api.openvpn.com"
}

data "cloudconnexa_devices" "all" {
}

resource "cloudconnexa_device" "test" {
	device_id   = data.cloudconnexa_devices.all.devices[0].id
	name        = "terraform-test-device"
	description = "Device managed by Terraform acceptance test"
}
`
}
