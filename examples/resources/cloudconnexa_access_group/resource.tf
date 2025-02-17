# Easy example to allow all resources in source and destination
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
}

# "cloudconnexa_network_application" added here to create at least 1 entity which will be "child" to resource "cloudconnexa_network"
# otherwise creation of resource "cloudconnexa_access_group" "limited_example1" will fail.
resource "cloudconnexa_network_application" "example1" {
  name        = "example-application-1"
  description = "Managed by Terraform"
  network_id  = cloudconnexa_network.this.id
  routes {
    domain            = "example-application-1.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["ANY"]
  }

  depends_on = [cloudconnexa_network.this]
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
