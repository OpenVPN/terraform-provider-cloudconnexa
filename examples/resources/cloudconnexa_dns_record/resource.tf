resource "cloudconnexa_dns_record" "this" {
  domain          = "test.example.local"
  description     = "Managed by Terraform"
  ip_v4_addresses = ["192.168.144.5"]
}
