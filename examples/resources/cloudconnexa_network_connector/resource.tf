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

resource "cloudconnexa_network_connector" "this" {
  name          = "test_connector"
  description   = "Managed by Terraform"
  vpn_region_id = "eu-central-2"
  network_id    = cloudconnexa_network.this.id
}
