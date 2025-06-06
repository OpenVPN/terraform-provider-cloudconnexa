package cloudconnexa

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccCloudConnexaUserGroup_basic tests the basic functionality of creating and updating a user group resource.
// It verifies that the user group is created with the correct name and VPN region IDs, and can be updated.
//
// Parameters:
//   - t: The testing context
func TestAccCloudConnexaUserGroup_basic(t *testing.T) {
	rn := "cloudconnexa_user_group.test"
	userGroup := cloudconnexa.UserGroup{
		Name: acctest.RandStringFromCharSet(10, alphabet),
		VpnRegionIDs: []string{
			"us-east-1",
		},
	}
	userGroupChanged := userGroup
	userGroupChanged.Name = fmt.Sprintf("changed-%s", acctest.RandStringFromCharSet(10, alphabet))

	check := func(userGroup cloudconnexa.UserGroup) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckCloudConnexaUserGroupExists(rn),
			resource.TestCheckResourceAttr(rn, "name", userGroup.Name),
			resource.TestCheckResourceAttr(rn, "vpn_region_ids.0", userGroup.VpnRegionIDs[0]),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaUserGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaUserGroupConfig(userGroup),
				Check:  check(userGroup),
			},
			{
				Config: testAccCloudConnexaUserGroupConfig(userGroupChanged),
				Check:  check(userGroupChanged),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccCloudConnexaUserImportStateIdFunc(rn),
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCheckCloudConnexaUserGroupDestroy verifies that the user group has been properly destroyed.
// It checks that the user group no longer exists in the CloudConnexa API.
//
// Parameters:
//   - s: The Terraform state containing information about the destroyed resources
//
// Returns:
//   - error: An error if the user group still exists or if there was an error checking its existence
func testAccCheckCloudConnexaUserGroupDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*cloudconnexa.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_user_group" {
			continue
		}
		id := rs.Primary.Attributes["id"]
		u, err := c.UserGroups.Get(id)
		if err == nil {
			if u != nil {
				return errors.New("user still exists")
			}
		}
	}
	return nil
}

// testAccCheckCloudConnexaUserGroupExists verifies that a user group exists in the CloudConnexa API.
// It checks that the resource exists and has a valid ID.
//
// Parameters:
//   - rn: The resource name to check
//
// Returns:
//   - resource.TestCheckFunc: A function that performs the existence check
func testAccCheckCloudConnexaUserGroupExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		c := testAccProvider.Meta().(*cloudconnexa.Client)
		_, err := c.UserGroups.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

// testAccCloudConnexaUserGroupConfig generates a Terraform configuration for a user group resource.
// It creates a configuration string with the specified user group properties.
//
// Parameters:
//   - userGroup: The user group configuration to use
//
// Returns:
//   - string: A Terraform configuration string for the user group
func testAccCloudConnexaUserGroupConfig(userGroup cloudconnexa.UserGroup) string {
	idsStr, _ := json.Marshal(userGroup.VpnRegionIDs)

	return fmt.Sprintf(`
provider "cloudconnexa" {
	base_url = "https://%s.api.openvpn.com"
}
resource "cloudconnexa_user_group" "test" {
  name           = "%s"
  vpn_region_ids = %s
  connect_auth   = "ON_PRIOR_AUTH"
}
`, testCloudID, userGroup.Name, idsStr)
}
