resource "cloudconnexa_settings" "example" {
  # Authentication and security settings
  allow_trusted_devices = true
  two_factor_auth       = true
  connect_auth          = "ON_PRIOR_AUTH"

  # DNS configuration
  dns_proxy_enabled  = true
  default_dns_suffix = "example.com"

  dns_servers {
    primary_ip_v4   = "8.8.8.8"
    secondary_ip_v4 = "8.8.4.4"
  }

  dns_zones {
    name      = "internal.example.com"
    addresses = ["10.0.1.10", "10.0.1.11"]
  }

  dns_zones {
    name      = "db.example.com"
    addresses = ["10.0.2.10"]
  }

  # Device management
  device_allowance_per_user     = 5
  device_allowance_force_update = true
  device_enforcement            = "LEARN_AND_ENFORCE"

  # Profile and connection settings
  profile_distribution = "AUTOMATIC"
  connection_timeout   = 300

  # Client options
  client_options = [
    "block-outside-dns",
    "dhcp-option DNS 8.8.8.8"
  ]

  # Network topology
  topology       = "FULL_MESH"
  default_region = "us-west-1"
  snat           = true

  # Domain routing subnet
  domain_routing_subnet {
    ip_v4_address = "192.168.100.0/24"
    ip_v6_address = "2001:db8::/64"
  }

  # Subnet configuration
  subnet {
    ip_v4_address = ["10.0.0.0/8", "172.16.0.0/12"]
    ip_v6_address = ["2001:db8:1::/64"]
  }
}
