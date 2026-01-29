package cloudconnexa

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccCloudConnexaSettings_import tests that importing the settings resource
// properly populates all attributes from the API.
func TestAccCloudConnexaSettings_import(t *testing.T) {
	rn := "cloudconnexa_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaSettingsConfig(),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateId:     "settings",
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccCloudConnexaSettings_basic tests that the settings resource can be created
// and read back with the correct values.
func TestAccCloudConnexaSettings_basic(t *testing.T) {
	rn := "cloudconnexa_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaSettingsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "id", "settings"),
					resource.TestCheckResourceAttrSet(rn, "allow_trusted_devices"),
					resource.TestCheckResourceAttrSet(rn, "two_factor_auth"),
					resource.TestCheckResourceAttrSet(rn, "dns_proxy_enabled"),
					resource.TestCheckResourceAttrSet(rn, "device_allowance_force_update"),
					resource.TestCheckResourceAttrSet(rn, "snat"),
					resource.TestCheckResourceAttrSet(rn, "dns_log_enabled"),
					resource.TestCheckResourceAttrSet(rn, "access_visibility_enabled"),
					resource.TestCheckResourceAttrSet(rn, "connect_auth"),
					resource.TestCheckResourceAttrSet(rn, "device_allowance_per_user"),
					resource.TestCheckResourceAttrSet(rn, "device_enforcement"),
					resource.TestCheckResourceAttrSet(rn, "profile_distribution"),
					resource.TestCheckResourceAttrSet(rn, "connection_timeout"),
					resource.TestCheckResourceAttrSet(rn, "default_region"),
					resource.TestCheckResourceAttrSet(rn, "topology"),
				),
			},
		},
	})
}

// testAccCloudConnexaSettingsConfig generates a minimal Terraform configuration
// for the settings resource that relies on reading existing settings from the API.
func testAccCloudConnexaSettingsConfig() string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
	base_url = "https://%s.api.openvpn.com"
}

resource "cloudconnexa_settings" "test" {
}
`, testCloudID)
}
