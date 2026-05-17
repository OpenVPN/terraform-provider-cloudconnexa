package cloudconnexa

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

// TestAccCloudConnexaDevice_basic exercises the full standalone lifecycle of
// the device resource against the real CloudConnexa API: create, update,
// import, and destroy.
func TestAccCloudConnexaDevice_basic(t *testing.T) {
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("set %s=1 to run acceptance tests", resource.EnvTfAcc)
	}
	testAccPreCheck(t)

	rn := "cloudconnexa_device.test"
	userID := pickTestUserID(t)

	name1 := "tf-acc-" + acctest.RandStringFromCharSet(8, alphabet)
	desc1 := "first description"
	name2 := "tf-acc-" + acctest.RandStringFromCharSet(8, alphabet)
	desc2 := "second description"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaDeviceConfig(userID, name1, desc1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudConnexaDeviceExists(rn),
					resource.TestCheckResourceAttr(rn, "user_id", userID),
					resource.TestCheckResourceAttr(rn, "name", name1),
					resource.TestCheckResourceAttr(rn, "description", desc1),
					resource.TestCheckResourceAttrSet(rn, "device_id"),
				),
			},
			{
				Config: testAccCloudConnexaDeviceConfig(userID, name2, desc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudConnexaDeviceExists(rn),
					resource.TestCheckResourceAttr(rn, "name", name2),
					resource.TestCheckResourceAttr(rn, "description", desc2),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccCloudConnexaDeviceImportStateIdFunc(rn),
				ImportStateVerify: true,
			},
		},
	})
}

// pickTestUserID returns the ID of an existing user in the test org so the
// device resource has a parent to attach to. Skips when no users are present.
func pickTestUserID(t *testing.T) string {
	t.Helper()
	c, err := cloudconnexa.NewClient(
		testBaseURL,
		os.Getenv(ClientIDEnvVar),
		os.Getenv(ClientSecretEnvVar),
	)
	if err != nil {
		t.Fatalf("failed to create CloudConnexa client: %s", err)
	}
	users, err := c.Users.List()
	if err != nil {
		t.Fatalf("failed to list users: %s", err)
	}
	if len(users) == 0 {
		t.Skip("no users in test org; cannot attach a device")
	}
	return users[0].ID
}

// testAccCheckCloudConnexaDeviceExists asserts the resource is present in state
// and the device can be fetched from the API.
func testAccCheckCloudConnexaDeviceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}
		userID := rs.Primary.Attributes["user_id"]
		if userID == "" {
			return errors.New("no user_id is set")
		}
		c := testAccProvider.Meta().(*cloudconnexa.Client)
		if _, err := c.Devices.GetByID(userID, rs.Primary.ID); err != nil {
			return fmt.Errorf("device %s not retrievable: %w", rs.Primary.ID, err)
		}
		return nil
	}
}

// testAccCheckCloudConnexaDeviceDestroy confirms the device is gone from the
// API after Terraform destroys the resource.
func testAccCheckCloudConnexaDeviceDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*cloudconnexa.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_device" {
			continue
		}
		userID := rs.Primary.Attributes["user_id"]
		if userID == "" {
			continue
		}
		if _, err := c.Devices.GetByID(userID, rs.Primary.ID); err == nil {
			return fmt.Errorf("device %s still exists after destroy", rs.Primary.ID)
		}
	}
	return nil
}

// testAccCloudConnexaDeviceImportStateIdFunc returns the composite "user_id/device_id"
// import string that resourceDeviceImport expects.
func testAccCloudConnexaDeviceImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}
		return rs.Primary.Attributes["user_id"] + "/" + rs.Primary.ID, nil
	}
}

// testAccCloudConnexaDeviceConfig renders the minimal Terraform config the
// acceptance test applies for each step.
func testAccCloudConnexaDeviceConfig(userID, name, description string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
	base_url = %q
}

resource "cloudconnexa_device" "test" {
	user_id     = %q
	name        = %q
	description = %q
}
`, testBaseURL, userID, name, description)
}
