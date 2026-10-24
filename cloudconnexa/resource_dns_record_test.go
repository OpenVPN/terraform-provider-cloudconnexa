package cloudconnexa

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			{
				Config:             testAccCloudConnexaDnsRecordConfig(domainName),
				Check:              testAccCheckCloudConnexaDnsRecordDisappears(resourceName),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// testAccCheckCloudConnexaDnsRecordDisappears deletes the DNS record behind
// Terraform's back, so the follow-up refresh must drop it from state and plan
// a re-create instead of failing on the API's 404.
func testAccCheckCloudConnexaDnsRecordDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		client := testAccProvider.Meta().(*cloudconnexa.Client)
		return client.DNSRecords.Delete(rs.Primary.ID)
	}
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
		r, err := client.DNSRecords.GetByID(recordId)
		if err != nil {
			if isNotFoundErr(err) {
				continue
			}
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
  base_url = "%[1]s"
}

resource "cloudconnexa_dns_record" "test" {
  domain          = "%[2]s"
  description     = "test description"
  ip_v4_addresses = ["192.168.1.1", "192.168.1.2"]
  ip_v6_addresses = ["2001:db8:85a3:0:0:8a2e:370:7334", "2001:db8:85a3:0:0:8a2e:370:7335"]
}
`, testBaseURL, domainName)
}

// TestUnitResourceDnsRecordRead_NotFound covers a record deleted out of band:
// the API's 404 removes the record from state instead of failing the refresh.
func TestUnitResourceDnsRecordRead_NotFound(t *testing.T) {
	c := newHostUnitTestClient(t, hostsHandlerError(http.StatusNotFound))
	d := schema.TestResourceDataRaw(t, resourceDnsRecord().Schema, map[string]interface{}{
		"domain": "x.example.com",
	})
	d.SetId("record-id")
	diags := resourceDnsRecordRead(context.Background(), d, c)
	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Empty(t, d.Id())
}

// TestUnitResourceDnsRecordRead_Error covers non-404 failures, including the
// API's former 400 "not found" response, which must no longer be swallowed.
func TestUnitResourceDnsRecordRead_Error(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusInternalServerError} {
		c := newHostUnitTestClient(t, hostsHandlerError(status))
		d := schema.TestResourceDataRaw(t, resourceDnsRecord().Schema, map[string]interface{}{
			"domain": "x.example.com",
		})
		d.SetId("record-id")
		diags := resourceDnsRecordRead(context.Background(), d, c)
		assert.True(t, diags.HasError(), "status %d", status)
		assert.Equal(t, "record-id", d.Id(), "status %d", status)
	}
}

// TestUnitResourceDnsRecordDelete_NotFound covers deleting a record that is
// already gone: the API's 404 counts as a successful delete.
func TestUnitResourceDnsRecordDelete_NotFound(t *testing.T) {
	c := newHostUnitTestClient(t, hostsHandlerError(http.StatusNotFound))
	d := schema.TestResourceDataRaw(t, resourceDnsRecord().Schema, map[string]interface{}{
		"domain": "x.example.com",
	})
	d.SetId("record-id")
	diags := resourceDnsRecordDelete(context.Background(), d, c)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
}

// TestUnitResourceDnsRecordDelete_Error covers non-404 delete failures.
func TestUnitResourceDnsRecordDelete_Error(t *testing.T) {
	c := newHostUnitTestClient(t, hostsHandlerError(http.StatusInternalServerError))
	d := schema.TestResourceDataRaw(t, resourceDnsRecord().Schema, map[string]interface{}{
		"domain": "x.example.com",
	})
	d.SetId("record-id")
	diags := resourceDnsRecordDelete(context.Background(), d, c)
	assert.True(t, diags.HasError())
}
