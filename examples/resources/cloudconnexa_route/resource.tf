resource "cloudconnexa_network" "this" {
  description     = "Test network"
  egress          = true
  name            = "my_test_network"
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_route" "this" {
  description     = "Managed by Terraform"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.this.id
  subnet          = "192.168.1.0/24"
}
