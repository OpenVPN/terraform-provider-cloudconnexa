# Easy example To allow all resources in source and destination
resource "cloudconnexa_access_group" "full_mesh" {
  name        = "Access Group name"
  description = "Add your description here"
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

# More advanced example
resource "cloudconnexa_user_group" "ug01" {
  name                 = "ug01"
  all_regions_included = true
  connect_auth         = "ON_PRIOR_AUTH"
  internet_access      = "SPLIT_TUNNEL_ON"
  max_device           = "3"
}

resource "cloudconnexa_network" "this" {
  description     = "Test network"
  egress          = true
  name            = "my_test_network"
  internet_access = "SPLIT_TUNNEL_ON"
  default_route {
    description = "Managed by Terraform"
    subnet      = "192.168.144.0/24"
    type        = "IP_V4"
  }
  default_connector {
    description   = "Managed by Terraform"
    name          = "test-connector"
    vpn_region_id = "eu-central-1"
  }
}

resource "cloudconnexa_access_group" "limited_example1" {
  name        = "limited_example1"
  description = "Add your description here"
  source {
    type        = "USER_GROUP"
    all_covered = false
    children    = [cloudconnexa_user_group.ug01.id]
  }
  destination {
    type        = "NETWORK"
    all_covered = true
    parent      = cloudconnexa_network.this.id
  }
}
