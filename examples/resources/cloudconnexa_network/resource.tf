# Basic network example
resource "cloudconnexa_network" "basic" {
  name            = "basic-network"
  description     = "Basic network for development environment"
  egress          = false
  internet_access = "SPLIT_TUNNEL_ON"
}

# Advanced network with egress
resource "cloudconnexa_network" "production" {
  name            = "production-network"
  description     = "Production network with egress capabilities"
  egress          = true
  internet_access = "SPLIT_TUNNEL_OFF"
}

# Restricted internet access network
resource "cloudconnexa_network" "secure" {
  name            = "secure-network"
  description     = "Secure network with restricted internet access"
  egress          = true
  internet_access = "RESTRICTED_INTERNET"
}

# Network with connector
resource "cloudconnexa_network" "with_connector" {
  name            = "network-with-connector"
  description     = "Network with dedicated connector"
  egress          = true
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_network_connector" "connector" {
  name          = "main-connector"
  description   = "Primary connector for the network"
  network_id    = cloudconnexa_network.with_connector.id
  vpn_region_id = "us-west-1"
}

# Network with custom routes
resource "cloudconnexa_route" "custom_route" {
  type            = "IP_V4"
  subnet          = "192.168.1.0/24"
  network_item_id = cloudconnexa_network.with_connector.id
  description     = "Custom route for internal services"
}

# Output examples
output "network_ids" {
  description = "IDs of created networks"
  value = {
    basic      = cloudconnexa_network.basic.id
    production = cloudconnexa_network.production.id
    secure     = cloudconnexa_network.secure.id
  }
}

output "network_system_subnets" {
  description = "System-assigned subnets for networks"
  value = {
    basic      = cloudconnexa_network.basic.system_subnets
    production = cloudconnexa_network.production.system_subnets
    secure     = cloudconnexa_network.secure.system_subnets
  }
}
