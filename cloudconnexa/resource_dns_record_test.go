package cloudconnexa

import (
	"fmt"
	"testing"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccCloudConnexaDnsRecord_basic tests the basic functionality of creating a DNS record resource.
// It verifies that the DNS record is created with the correct domain name, description, and IP addresses.
//
// Parameters:
//   - t: The testing context
func TestAccCloudConnexaDnsRecord_basic(t *testing.T) {
	resourceName := "cloudconnexa_dns_record.test"
	domainName := "test.cloudconnexa.com"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaDnsRecordConfig(domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "ip_v4_addresses.0", "192.168.1.1"),
					resource.TestCheckResourceAttr(resourceName, "ip_v4_addresses.1", "192.168.1.2"),
					resource.TestCheckResourceAttr(resourceName, "ip_v6_addresses.0", "2001:db8:85a3:0:0:8a2e:370:7334"),
					resource.TestCheckResourceAttr(resourceName, "ip_v6_addresses.1", "2001:db8:85a3:0:0:8a2e:370:7335"),
				),
			},
		},
	})
}

// testAccCheckCloudConnexaDnsRecordDestroy verifies that the DNS record has been properly destroyed.
// It checks that the DNS record no longer exists in the CloudConnexa API.
//
// Parameters:
//   - s: The Terraform state containing information about the destroyed resources
//
// Returns:
//   - error: An error if the DNS record still exists or if there was an error checking its existence
func testAccCheckCloudConnexaDnsRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudconnexa.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_dns_record" {
			continue
		}

		recordId := rs.Primary.ID
		r, err := client.DNSRecords.GetDNSRecord(recordId)

		if err != nil {
			return err
		}

		if r != nil {
			return fmt.Errorf("DNS record with ID '%s' still exists", recordId)
		}
	}

	return nil
}

// testAccCloudConnexaDnsRecordConfig generates the Terraform configuration for testing DNS record creation.
// It creates a provider configuration and a DNS record resource with the specified domain name.
//
// Parameters:
//   - domainName: The domain name to use for the DNS record
//
// Returns:
//   - string: The Terraform configuration as a string
func testAccCloudConnexaDnsRecordConfig(domainName string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
  base_url = "https://%[1]s.api.openvpn.com"
}

resource "cloudconnexa_dns_record" "test" {
  domain          = "%[2]s"
  description     = "test description"
  ip_v4_addresses = ["192.168.1.1", "192.168.1.2"]
  ip_v6_addresses = ["2001:db8:85a3:0:0:8a2e:370:7334", "2001:db8:85a3:0:0:8a2e:370:7335"]
}
`, testCloudID, domainName)
}
