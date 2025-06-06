# Basic IPsec network configuration
resource "cloudconnexa_ipsec_network" "basic" {
  name            = "basic-ipsec-network"
  description     = "Basic IPsec network for site-to-site connectivity"
  egress          = true
  internet_access = "RESTRICTED_INTERNET"

  ipsec_config {
    remote_gateway       = "203.0.113.10"
    remote_networks      = ["192.168.1.0/24", "192.168.2.0/24"]
    pre_shared_key       = var.ipsec_psk
    ike_version          = "IKEv2"
    encryption_algorithm = "AES256"
    hash_algorithm       = "SHA256"
    dh_group             = "14"
    pfs_group            = "14"
    ike_lifetime         = 28800
    ipsec_lifetime       = 3600
    dpd_timeout          = 30
    nat_traversal        = true
  }

  connector {
    name          = "basic-ipsec-connector"
    description   = "IPsec connector for basic network"
    vpn_region_id = "us-east-1"
  }
}

# Advanced IPsec network with high security settings
resource "cloudconnexa_ipsec_network" "high_security" {
  name            = "high-security-ipsec"
  description     = "High security IPsec network for sensitive data"
  egress          = true
  internet_access = "RESTRICTED_INTERNET"

  ipsec_config {
    remote_gateway       = "198.51.100.20"
    remote_networks      = ["10.10.0.0/16", "10.20.0.0/16"]
    pre_shared_key       = var.high_security_psk
    ike_version          = "IKEv2"
    encryption_algorithm = "AES256"
    hash_algorithm       = "SHA384"
    dh_group             = "19"
    pfs_group            = "19"
    ike_lifetime         = 14400 # 4 hours
    ipsec_lifetime       = 1800  # 30 minutes
    dpd_timeout          = 15
    nat_traversal        = false
  }

  connector {
    name          = "high-security-connector"
    description   = "High security IPsec connector"
    vpn_region_id = "us-east-1"
  }
}

# IPsec network for AWS VPC connectivity
resource "cloudconnexa_ipsec_network" "aws_vpc" {
  name            = "aws-vpc-ipsec"
  description     = "IPsec tunnel to AWS VPC"
  egress          = true
  internet_access = "SPLIT_TUNNEL_ON"

  ipsec_config {
    remote_gateway       = aws_vpn_connection.main.tunnel1_address
    remote_networks      = [aws_vpc.main.cidr_block]
    pre_shared_key       = aws_vpn_connection.main.tunnel1_preshared_key
    ike_version          = "IKEv1" # AWS VPN Gateway compatibility
    encryption_algorithm = "AES128"
    hash_algorithm       = "SHA1"
    dh_group             = "2"
    pfs_group            = "2"
    ike_lifetime         = 28800
    ipsec_lifetime       = 3600
    dpd_timeout          = 30
    nat_traversal        = true
  }

  connector {
    name          = "aws-vpc-connector"
    description   = "Connector for AWS VPC tunnel"
    vpn_region_id = "us-east-1"
  }
}

# IPsec network for Azure VNet connectivity
resource "cloudconnexa_ipsec_network" "azure_vnet" {
  name            = "azure-vnet-ipsec"
  description     = "IPsec tunnel to Azure VNet"
  egress          = true
  internet_access = "SPLIT_TUNNEL_ON"

  ipsec_config {
    remote_gateway       = azurerm_virtual_network_gateway.main.bgp_settings[0].peering_address
    remote_networks      = [azurerm_virtual_network.main.address_space[0]]
    pre_shared_key       = var.azure_vpn_psk
    ike_version          = "IKEv2"
    encryption_algorithm = "AES256"
    hash_algorithm       = "SHA256"
    dh_group             = "14"
    pfs_group            = "14"
    ike_lifetime         = 28800
    ipsec_lifetime       = 3600
    dpd_timeout          = 30
    nat_traversal        = true
  }

  connector {
    name          = "azure-vnet-connector"
    description   = "Connector for Azure VNet tunnel"
    vpn_region_id = "eu-central-1"
  }
}

# IPsec network for Google Cloud VPC connectivity
resource "cloudconnexa_ipsec_network" "gcp_vpc" {
  name            = "gcp-vpc-ipsec"
  description     = "IPsec tunnel to Google Cloud VPC"
  egress          = true
  internet_access = "SPLIT_TUNNEL_ON"

  ipsec_config {
    remote_gateway       = google_compute_vpn_gateway.main.ip_address
    remote_networks      = [google_compute_network.main.ipv4_range]
    pre_shared_key       = var.gcp_vpn_psk
    ike_version          = "IKEv2"
    encryption_algorithm = "AES256"
    hash_algorithm       = "SHA256"
    dh_group             = "14"
    pfs_group            = "14"
    ike_lifetime         = 28800
    ipsec_lifetime       = 3600
    dpd_timeout          = 30
    nat_traversal        = true
  }

  connector {
    name          = "gcp-vpc-connector"
    description   = "Connector for GCP VPC tunnel"
    vpn_region_id = "us-west-1"
  }
}

# Multi-region IPsec network with redundancy
resource "cloudconnexa_ipsec_network" "multi_region_primary" {
  name            = "multi-region-ipsec-primary"
  description     = "Primary IPsec tunnel for multi-region setup"
  egress          = true
  internet_access = "RESTRICTED_INTERNET"

  ipsec_config {
    remote_gateway       = "203.0.113.30"
    remote_networks      = ["172.16.0.0/12"]
    pre_shared_key       = var.primary_ipsec_psk
    ike_version          = "IKEv2"
    encryption_algorithm = "AES256"
    hash_algorithm       = "SHA256"
    dh_group             = "16"
    pfs_group            = "16"
    ike_lifetime         = 28800
    ipsec_lifetime       = 3600
    dpd_timeout          = 30
    nat_traversal        = true
  }

  connector {
    name          = "multi-region-primary-connector"
    description   = "Primary connector for multi-region setup"
    vpn_region_id = "us-east-1"
  }
}

resource "cloudconnexa_ipsec_network" "multi_region_secondary" {
  name            = "multi-region-ipsec-secondary"
  description     = "Secondary IPsec tunnel for multi-region setup"
  egress          = true
  internet_access = "RESTRICTED_INTERNET"

  ipsec_config {
    remote_gateway       = "203.0.113.31"
    remote_networks      = ["172.16.0.0/12"]
    pre_shared_key       = var.secondary_ipsec_psk
    ike_version          = "IKEv2"
    encryption_algorithm = "AES256"
    hash_algorithm       = "SHA256"
    dh_group             = "16"
    pfs_group            = "16"
    ike_lifetime         = 28800
    ipsec_lifetime       = 3600
    dpd_timeout          = 30
    nat_traversal        = true
  }

  connector {
    name          = "multi-region-secondary-connector"
    description   = "Secondary connector for multi-region setup"
    vpn_region_id = "us-west-1"
  }
}

# IPsec network for legacy systems with older encryption
resource "cloudconnexa_ipsec_network" "legacy_systems" {
  name            = "legacy-systems-ipsec"
  description     = "IPsec tunnel for legacy systems compatibility"
  egress          = false
  internet_access = "RESTRICTED_INTERNET"

  ipsec_config {
    remote_gateway       = "198.51.100.40"
    remote_networks      = ["192.168.100.0/24"]
    pre_shared_key       = var.legacy_ipsec_psk
    ike_version          = "IKEv1"
    encryption_algorithm = "3DES"
    hash_algorithm       = "MD5"
    dh_group             = "2"
    pfs_group            = "2"
    ike_lifetime         = 86400 # 24 hours for legacy compatibility
    ipsec_lifetime       = 28800 # 8 hours
    dpd_timeout          = 60
    nat_traversal        = true
  }

  connector {
    name          = "legacy-systems-connector"
    description   = "Connector for legacy systems"
    vpn_region_id = "us-east-1"
  }
}

# IPsec networks using for_each for multiple sites
variable "branch_offices" {
  description = "Branch office IPsec configurations"
  type = map(object({
    remote_gateway  = string
    remote_networks = list(string)
    vpn_region_id   = string
    description     = string
  }))
  default = {
    "new-york" = {
      remote_gateway  = "203.0.113.50"
      remote_networks = ["192.168.10.0/24"]
      vpn_region_id   = "us-east-1"
      description     = "New York branch office"
    }
    "london" = {
      remote_gateway  = "203.0.113.51"
      remote_networks = ["192.168.20.0/24"]
      vpn_region_id   = "eu-west-1"
      description     = "London branch office"
    }
    "tokyo" = {
      remote_gateway  = "203.0.113.52"
      remote_networks = ["192.168.30.0/24"]
      vpn_region_id   = "ap-northeast-1"
      description     = "Tokyo branch office"
    }
  }
}

resource "cloudconnexa_ipsec_network" "branch_offices" {
  for_each = var.branch_offices

  name            = "${each.key}-ipsec-network"
  description     = each.value.description
  egress          = true
  internet_access = "SPLIT_TUNNEL_ON"

  ipsec_config {
    remote_gateway       = each.value.remote_gateway
    remote_networks      = each.value.remote_networks
    pre_shared_key       = var.branch_office_psks[each.key]
    ike_version          = "IKEv2"
    encryption_algorithm = "AES256"
    hash_algorithm       = "SHA256"
    dh_group             = "14"
    pfs_group            = "14"
    ike_lifetime         = 28800
    ipsec_lifetime       = 3600
    dpd_timeout          = 30
    nat_traversal        = true
  }

  connector {
    name          = "${each.key}-ipsec-connector"
    description   = "IPsec connector for ${each.value.description}"
    vpn_region_id = each.value.vpn_region_id
  }
}

# Routes for IPsec networks
resource "cloudconnexa_route" "ipsec_routes" {
  for_each = var.branch_offices

  description     = "Route to ${each.key} branch office"
  type            = "IP_V4"
  network_item_id = cloudconnexa_ipsec_network.branch_offices[each.key].id
  subnet          = each.value.remote_networks[0]
}

# Variables for IPsec configuration
variable "ipsec_psk" {
  description = "Pre-shared key for basic IPsec tunnel"
  type        = string
  sensitive   = true
}

variable "high_security_psk" {
  description = "Pre-shared key for high security IPsec tunnel"
  type        = string
  sensitive   = true
}

variable "azure_vpn_psk" {
  description = "Pre-shared key for Azure VPN connection"
  type        = string
  sensitive   = true
}

variable "gcp_vpn_psk" {
  description = "Pre-shared key for GCP VPN connection"
  type        = string
  sensitive   = true
}

variable "primary_ipsec_psk" {
  description = "Pre-shared key for primary IPsec tunnel"
  type        = string
  sensitive   = true
}

variable "secondary_ipsec_psk" {
  description = "Pre-shared key for secondary IPsec tunnel"
  type        = string
  sensitive   = true
}

variable "legacy_ipsec_psk" {
  description = "Pre-shared key for legacy systems IPsec tunnel"
  type        = string
  sensitive   = true
}

variable "branch_office_psks" {
  description = "Pre-shared keys for branch office IPsec tunnels"
  type        = map(string)
  sensitive   = true
}

# Outputs
output "ipsec_networks" {
  description = "IPsec network configurations"
  value = {
    basic         = cloudconnexa_ipsec_network.basic.id
    high_security = cloudconnexa_ipsec_network.high_security.id
    aws_vpc       = cloudconnexa_ipsec_network.aws_vpc.id
    azure_vnet    = cloudconnexa_ipsec_network.azure_vnet.id
    gcp_vpc       = cloudconnexa_ipsec_network.gcp_vpc.id
  }
}

output "branch_office_networks" {
  description = "Branch office IPsec network IDs"
  value = {
    for k, v in cloudconnexa_ipsec_network.branch_offices : k => v.id
  }
}

output "ipsec_connectors" {
  description = "IPsec connector information"
  value = {
    for k, v in cloudconnexa_ipsec_network.branch_offices : k => {
      connector_id  = v.connector[0].id
      ip_v4_address = v.connector[0].ip_v4_address
      vpn_region_id = v.connector[0].vpn_region_id
    }
  }
}
