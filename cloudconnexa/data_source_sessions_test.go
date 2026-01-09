package cloudconnexa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccCloudConnexaDataSourceSessions_basic tests the basic functionality of the sessions data source.
func TestAccCloudConnexaDataSourceSessions_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaDataSourceSessionsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudconnexa_sessions.test", "sessions.#"),
				),
			},
		},
	})
}

// TestAccCloudConnexaDataSourceSessions_withStatus tests the sessions data source with status filter.
func TestAccCloudConnexaDataSourceSessions_withStatus(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaDataSourceSessionsConfigWithStatus("ACTIVE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cloudconnexa_sessions.test", "sessions.#"),
				),
			},
		},
	})
}

// testAccCloudConnexaDataSourceSessionsConfig generates a Terraform configuration for sessions data source.
func testAccCloudConnexaDataSourceSessionsConfig() string {
	return `
provider "cloudconnexa" {
	base_url = "https://` + testCloudID + `.api.openvpn.com"
}

data "cloudconnexa_sessions" "test" {
}
`
}

// testAccCloudConnexaDataSourceSessionsConfigWithStatus generates a Terraform configuration with status filter.
func testAccCloudConnexaDataSourceSessionsConfigWithStatus(status string) string {
	return `
provider "cloudconnexa" {
	base_url = "https://` + testCloudID + `.api.openvpn.com"
}

data "cloudconnexa_sessions" "test" {
	status = "` + status + `"
}
`
}
