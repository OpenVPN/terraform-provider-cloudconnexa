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

// TestAccCloudConnexaUser_basic tests the basic functionality of creating and updating a user resource.
// It verifies that the user is created with the correct username, email, first name, and last name, and can be updated.
//
// Parameters:
//   - t: The testing context
func TestAccCloudConnexaUser_basic(t *testing.T) {
	rn := "cloudconnexa_user.test"
	user := cloudconnexa.User{
		Username:  acctest.RandStringFromCharSet(10, alphabet),
		FirstName: acctest.RandStringFromCharSet(10, alphabet),
		LastName:  acctest.RandStringFromCharSet(10, alphabet),
		Email:     fmt.Sprintf("terraform-tests+%s@devopenvpn.in", acctest.RandString(10)),
	}
	userChanged := user
	userChanged.Email = fmt.Sprintf("terraform-tests+changed%s@devopenvpn.in", acctest.RandString(10))
	var userID string

	check := func(user cloudconnexa.User) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckCloudConnexaUserExists(rn, &userID),
			resource.TestCheckResourceAttr(rn, "username", user.Username),
			resource.TestCheckResourceAttr(rn, "email", user.Email),
			resource.TestCheckResourceAttr(rn, "first_name", user.FirstName),
			resource.TestCheckResourceAttr(rn, "last_name", user.LastName),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaUserConfig(user),
				Check:  check(user),
			},
			{
				Config: testAccCloudConnexaUserConfig(userChanged),
				Check:  check(userChanged),
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

// testAccCheckCloudConnexaUserDestroy verifies that the user has been properly destroyed.
// It checks that the user no longer exists in the CloudConnexa API.
//
// Parameters:
//   - s: The Terraform state containing information about the destroyed resources
//
// Returns:
//   - error: An error if the user still exists or if there was an error checking its existence
func testAccCheckCloudConnexaUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudconnexa.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_user" {
			continue
		}
		username := rs.Primary.Attributes["username"]
		u, err := client.Users.Get(username)
		if err == nil {
			if u != nil {
				return errors.New("user still exists")
			}
		}
	}
	return nil
}

// testAccCheckCloudConnexaUserExists verifies that a user exists in the CloudConnexa API.
// It checks that the resource exists and has a valid ID.
//
// Parameters:
//   - n: The resource name to check
//   - teamID: A pointer to store the team ID if needed
//
// Returns:
//   - resource.TestCheckFunc: A function that performs the existence check
func testAccCheckCloudConnexaUserExists(n string, teamID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		client := testAccProvider.Meta().(*cloudconnexa.Client)
		_, err := client.Users.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

// testAccCloudConnexaUserImportStateIdFunc generates a function to retrieve the import state ID.
// It returns the primary ID of the resource for import operations.
//
// Parameters:
//   - n: The resource name to get the import state ID for
//
// Returns:
//   - resource.ImportStateIdFunc: A function that returns the import state ID
func testAccCloudConnexaUserImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}
		return rs.Primary.ID, nil
	}
}

// testAccCloudConnexaUserConfig generates a Terraform configuration for a user resource.
// It creates a configuration string with the specified user properties and associated user group.
//
// Parameters:
//   - user: The user configuration to use
//
// Returns:
//   - string: A Terraform configuration string for the user
func testAccCloudConnexaUserConfig(user cloudconnexa.User) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
	base_url = "https://%s.api.openvpn.com"
}
resource "cloudconnexa_user_group" "userGroup1" {
  name           = "test-group"
  vpn_region_ids = ["eu-central-1"]
  connect_auth   = "ON_PRIOR_AUTH"
}

resource "cloudconnexa_user" "test" {
	username   = "%s"
	email      = "%s"
	first_name = "%s"
	last_name  = "%s"
	group_id   = cloudconnexa_user_group.userGroup1.id
}
`, testCloudID, user.Username, user.Email, user.FirstName, user.LastName)
}
