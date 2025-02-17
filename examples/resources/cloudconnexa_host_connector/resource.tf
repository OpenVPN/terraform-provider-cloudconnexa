resource "cloudconnexa_host" "this" {
  name            = "test_host"
  description     = "Managed by Terraform"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "test.example.local"
}

resource "cloudconnexa_host_connector" "this" {
  name          = "test_connector"
  description   = "Managed by Terraform"
  host_id       = cloudconnexa_host.this.id
  vpn_region_id = "eu-central-2"
}
