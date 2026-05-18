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

// TestAccCloudConnexaAccessGroup_basic exercises the create/update/import lifecycle
// of a minimal access group (one source, one destination, both all_covered=true).
// The update step also flips the destination type from NETWORK to HOST so the
// destination conversion map is exercised for more than one input.
func TestAccCloudConnexaAccessGroup_basic(t *testing.T) {
	rn := "cloudconnexa_access_group.test"
	ag := cloudconnexa.AccessGroup{
		Name:        acctest.RandStringFromCharSet(10, alphabet),
		Description: "tf-acc " + acctest.RandStringFromCharSet(10, alphabet),
		Source: []cloudconnexa.AccessItem{
			{Type: "USER_GROUP", AllCovered: true},
		},
		Destination: []cloudconnexa.AccessItem{
			{Type: "NETWORK", AllCovered: true},
		},
	}
	agChanged := cloudconnexa.AccessGroup{
		Name:        fmt.Sprintf("changed-%s", acctest.RandStringFromCharSet(10, alphabet)),
		Description: "tf-acc changed " + acctest.RandStringFromCharSet(10, alphabet),
		Source: []cloudconnexa.AccessItem{
			{Type: "USER_GROUP", AllCovered: true},
		},
		Destination: []cloudconnexa.AccessItem{
			{Type: "HOST", AllCovered: true},
		},
	}
	var agID string

	check := func(ag cloudconnexa.AccessGroup) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckCloudConnexaAccessGroupExists(rn, &agID),
			resource.TestCheckResourceAttr(rn, "name", ag.Name),
			resource.TestCheckResourceAttr(rn, "description", ag.Description),
			resource.TestCheckResourceAttr(rn, "source.#", "1"),
			resource.TestCheckResourceAttr(rn, "destination.#", "1"),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "source.*", map[string]string{
				"type":        ag.Source[0].Type,
				"all_covered": "true",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "destination.*", map[string]string{
				"type":        ag.Destination[0].Type,
				"all_covered": "true",
			}),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaAccessGroupConfigBasic(ag),
				Check:  check(ag),
			},
			{
				Config: testAccCloudConnexaAccessGroupConfigBasic(agChanged),
				Check:  check(agChanged),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccCloudConnexaAccessGroupImportStateIdFunc(rn),
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccCloudConnexaAccessGroup_withChildren covers the explicit-children
// code path of setAccessGroupData (the branch that surfaces children only for
// USER_GROUP or when all_covered=false) and the NETWORK <-> NETWORK_SERVICE
// conversion maps. The access group references a user_group and a network
// created in the same Terraform run, so the children/parent IDs are real.
func TestAccCloudConnexaAccessGroup_withChildren(t *testing.T) {
	rn := "cloudconnexa_access_group.test"
	agName := acctest.RandStringFromCharSet(10, alphabet)
	ugName := acctest.RandStringFromCharSet(10, alphabet)
	netName := "tf-acc-" + acctest.RandString(10)
	var agID string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaAccessGroupConfigWithChildren(agName, ugName, netName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudConnexaAccessGroupExists(rn, &agID),
					resource.TestCheckResourceAttr(rn, "name", agName),
					resource.TestCheckResourceAttr(rn, "source.#", "1"),
					resource.TestCheckResourceAttr(rn, "destination.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "source.*", map[string]string{
						"type":        "USER_GROUP",
						"all_covered": "false",
						"children.#":  "1",
					}),
					resource.TestCheckTypeSetElemAttrPair(rn, "source.*.children.*", "cloudconnexa_user_group.ug", "id"),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "destination.*", map[string]string{
						"type":        "NETWORK",
						"all_covered": "true",
					}),
					resource.TestCheckTypeSetElemAttrPair(rn, "destination.*.parent", "cloudconnexa_network.net", "id"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccCloudConnexaAccessGroupImportStateIdFunc(rn),
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccCloudConnexaAccessGroup_fullMesh verifies that multiple source and
// destination blocks of different types round-trip correctly through create
// and read. This guards against regressions where sets get collapsed or where
// a type is silently dropped during conversion.
func TestAccCloudConnexaAccessGroup_fullMesh(t *testing.T) {
	rn := "cloudconnexa_access_group.test"
	agName := acctest.RandStringFromCharSet(10, alphabet)
	var agID string

	typeAttrs := func(t string) map[string]string {
		return map[string]string{"type": t, "all_covered": "true"}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCloudConnexaAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudConnexaAccessGroupConfigFullMesh(agName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudConnexaAccessGroupExists(rn, &agID),
					resource.TestCheckResourceAttr(rn, "name", agName),
					resource.TestCheckResourceAttr(rn, "source.#", "3"),
					resource.TestCheckResourceAttr(rn, "destination.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "source.*", typeAttrs("NETWORK")),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "source.*", typeAttrs("HOST")),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "source.*", typeAttrs("USER_GROUP")),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "destination.*", typeAttrs("NETWORK")),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "destination.*", typeAttrs("HOST")),
					resource.TestCheckTypeSetElemNestedAttrs(rn, "destination.*", typeAttrs("USER_GROUP")),
				),
			},
		},
	})
}

// testAccCheckCloudConnexaAccessGroupExists verifies that the access group at
// the given resource address is reachable via the CloudConnexa API, and writes
// its ID through the provided pointer.
func testAccCheckCloudConnexaAccessGroupExists(rn string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		c := testAccProvider.Meta().(*providerMeta).Client
		ag, err := c.AccessGroups.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if ag == nil {
			return fmt.Errorf("access group with ID %q not found", rs.Primary.ID)
		}
		*id = ag.ID
		return nil
	}
}

// testAccCheckCloudConnexaAccessGroupDestroy verifies that every access group
// previously tracked in state has been removed from the CloudConnexa API.
func testAccCheckCloudConnexaAccessGroupDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*providerMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudconnexa_access_group" {
			continue
		}
		ag, err := c.AccessGroups.Get(rs.Primary.ID)
		if err == nil && ag != nil {
			return errors.New("access group still exists")
		}
	}
	return nil
}

// testAccCloudConnexaAccessGroupImportStateIdFunc returns the primary ID of
// the resource for use in ImportState steps.
func testAccCloudConnexaAccessGroupImportStateIdFunc(rn string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return "", fmt.Errorf("not found: %s", rn)
		}
		return rs.Primary.ID, nil
	}
}

// testAccCloudConnexaAccessGroupConfigBasic renders a Terraform config with a
// single access group whose name, description, source type, and destination
// type come from the supplied AccessGroup. all_covered is hard-coded to true
// because this template is used by the basic test only.
func testAccCloudConnexaAccessGroupConfigBasic(ag cloudconnexa.AccessGroup) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
  base_url = "%[1]s"
}

resource "cloudconnexa_access_group" "test" {
  name        = "%[2]s"
  description = "%[3]s"

  source {
    type        = "%[4]s"
    all_covered = true
  }

  destination {
    type        = "%[5]s"
    all_covered = true
  }
}
`, testBaseURL, ag.Name, ag.Description, ag.Source[0].Type, ag.Destination[0].Type)
}

// testAccCloudConnexaAccessGroupConfigWithChildren renders a config that
// creates a user_group and a network alongside the access group, and wires
// their IDs into the access group's source.children and destination.parent.
func testAccCloudConnexaAccessGroupConfigWithChildren(agName, ugName, netName string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
  base_url = "%[1]s"
}

resource "cloudconnexa_user_group" "ug" {
  name           = "%[3]s"
  vpn_region_ids = ["us-east-1"]
  connect_auth   = "ON_PRIOR_AUTH"
}

resource "cloudconnexa_network" "net" {
  name = "%[4]s"
}

resource "cloudconnexa_route" "net_route" {
  network_item_id = cloudconnexa_network.net.id
  type            = "IP_V4"
  subnet          = "10.42.0.0/16"
  description     = "tf-acc route for access group test"
}

resource "cloudconnexa_access_group" "test" {
  name        = "%[2]s"
  description = "tf-acc with children"

  source {
    type        = "USER_GROUP"
    all_covered = false
    children    = [cloudconnexa_user_group.ug.id]
  }

  destination {
    type        = "NETWORK"
    all_covered = true
    parent      = cloudconnexa_network.net.id
  }

  depends_on = [cloudconnexa_route.net_route]
}
`, testBaseURL, agName, ugName, netName)
}

// testAccCloudConnexaAccessGroupConfigFullMesh renders a config with three
// source blocks and three destination blocks (NETWORK / HOST / USER_GROUP),
// each with all_covered=true.
func testAccCloudConnexaAccessGroupConfigFullMesh(agName string) string {
	return fmt.Sprintf(`
provider "cloudconnexa" {
  base_url = "%[1]s"
}

resource "cloudconnexa_access_group" "test" {
  name        = "%[2]s"
  description = "tf-acc full mesh"

  source {
    type        = "NETWORK"
    all_covered = true
  }
  source {
    type        = "HOST"
    all_covered = true
  }
  source {
    type        = "USER_GROUP"
    all_covered = true
  }

  destination {
    type        = "NETWORK"
    all_covered = true
  }
  destination {
    type        = "HOST"
    all_covered = true
  }
  destination {
    type        = "USER_GROUP"
    all_covered = true
  }
}
`, testBaseURL, agName)
}
