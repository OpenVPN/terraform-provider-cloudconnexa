resource "cloudconnexa_host" "this" {
  name            = "test_host"
  description     = "Managed by Terraform"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "test.example.local"
}
