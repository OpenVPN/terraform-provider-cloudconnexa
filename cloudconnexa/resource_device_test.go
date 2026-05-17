package cloudconnexa

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// TestAccCloudConnexaDevice_basic tests the basic functionality of the device resource.
// Note: This test requires at least one existing device in the CloudConnexa account.
func TestAccCloudConnexaDevice_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfNoDevices(t) },
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

// skipIfNoDevices skips the test when the test org has no enrolled devices,
// since the test config indexes data.cloudconnexa_devices.all.devices[0].
func skipIfNoDevices(t *testing.T) {
	client, err := cloudconnexa.NewClient(
		testBaseURL,
		os.Getenv(ClientIDEnvVar),
		os.Getenv(ClientSecretEnvVar),
	)
	if err != nil {
		t.Fatalf("failed to create CloudConnexa client: %s", err)
	}
	devices, err := client.Devices.ListAll()
	if err != nil {
		t.Fatalf("failed to list devices: %s", err)
	}
	if len(devices) == 0 {
		t.Skip("no devices enrolled in the test org; skipping device resource test")
	}
}

// testAccCloudConnexaDeviceConfig generates a Terraform configuration for device resource.
// This configuration first fetches all devices and then manages the first one found.
func testAccCloudConnexaDeviceConfig() string {
	return `
provider "cloudconnexa" {
	base_url = "` + testBaseURL + `"
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
