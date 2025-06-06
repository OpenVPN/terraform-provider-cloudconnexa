# Networks for different environments
resource "cloudconnexa_network" "production" {
  name            = "production-network"
  description     = "Production environment network"
  egress          = true
  internet_access = "SPLIT_TUNNEL_OFF"
}

resource "cloudconnexa_network" "staging" {
  name            = "staging-network"
  description     = "Staging environment network"
  egress          = true
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_network" "development" {
  name            = "development-network"
  description     = "Development environment network"
  egress          = false
  internet_access = "SPLIT_TUNNEL_ON"
}

# Hosts for routing examples
resource "cloudconnexa_host" "database_server" {
  name            = "database-server"
  description     = "Production database server"
  internet_access = "RESTRICTED_INTERNET"
  domain          = "db.production.example.com"
}

resource "cloudconnexa_host" "api_server" {
  name            = "api-server"
  description     = "REST API server"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "api.example.com"
}

# IPv4 routes for different subnet ranges
resource "cloudconnexa_route" "prod_internal_network" {
  description     = "Production internal network range"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.production.id
  subnet          = "10.0.0.0/16"
}

resource "cloudconnexa_route" "prod_database_subnet" {
  description     = "Production database subnet"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.production.id
  subnet          = "10.0.100.0/24"
}

resource "cloudconnexa_route" "staging_network" {
  description     = "Staging environment network"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.staging.id
  subnet          = "10.1.0.0/16"
}

resource "cloudconnexa_route" "dev_network" {
  description     = "Development network range"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.development.id
  subnet          = "10.2.0.0/16"
}

# IPv6 routes
resource "cloudconnexa_route" "prod_ipv6_network" {
  description     = "Production IPv6 network"
  type            = "IP_V6"
  network_item_id = cloudconnexa_network.production.id
  subnet          = "2001:db8::/32"
}

resource "cloudconnexa_route" "staging_ipv6_network" {
  description     = "Staging IPv6 network"
  type            = "IP_V6"
  network_item_id = cloudconnexa_network.staging.id
  subnet          = "2001:db8:1::/48"
}

# Domain-based routes
resource "cloudconnexa_route" "prod_domain_route" {
  description     = "Production domain routing"
  type            = "DOMAIN"
  network_item_id = cloudconnexa_network.production.id
  domain          = "*.production.example.com"
}

resource "cloudconnexa_route" "api_domain_route" {
  description     = "API domain routing"
  type            = "DOMAIN"
  network_item_id = cloudconnexa_network.production.id
  domain          = "api.example.com"
}

resource "cloudconnexa_route" "staging_domain_route" {
  description     = "Staging domain routing"
  type            = "DOMAIN"
  network_item_id = cloudconnexa_network.staging.id
  domain          = "*.staging.example.com"
}

resource "cloudconnexa_route" "dev_domain_route" {
  description     = "Development domain routing"
  type            = "DOMAIN"
  network_item_id = cloudconnexa_network.development.id
  domain          = "*.dev.example.com"
}

# Host-based routes
resource "cloudconnexa_route" "database_route" {
  description     = "Database server routing"
  type            = "IP_V4"
  network_item_id = cloudconnexa_host.database_server.id
  subnet          = "10.0.200.0/24"
}

resource "cloudconnexa_route" "api_server_route" {
  description     = "API server internal routing"
  type            = "IP_V4"
  network_item_id = cloudconnexa_host.api_server.id
  subnet          = "10.0.150.0/24"
}

# Multiple routes using for_each pattern
variable "network_subnets" {
  description = "Map of network subnets to create routes for"
  type = map(object({
    description = string
    type        = string
    subnet      = optional(string)
    domain      = optional(string)
    network_id  = string
  }))
  default = {
    "prod-web-tier" = {
      description = "Production web tier subnet"
      type        = "IP_V4"
      subnet      = "10.0.10.0/24"
      network_id  = "production"
    }
    "prod-app-tier" = {
      description = "Production application tier subnet"
      type        = "IP_V4"
      subnet      = "10.0.20.0/24"
      network_id  = "production"
    }
    "prod-data-tier" = {
      description = "Production data tier subnet"
      type        = "IP_V4"
      subnet      = "10.0.30.0/24"
      network_id  = "production"
    }
    "staging-services" = {
      description = "Staging services subnet"
      type        = "IP_V4"
      subnet      = "10.1.10.0/24"
      network_id  = "staging"
    }
    "microservices-domain" = {
      description = "Microservices domain routing"
      type        = "DOMAIN"
      domain      = "*.microservices.example.com"
      network_id  = "production"
    }
  }
}

resource "cloudconnexa_route" "network_subnets" {
  for_each = var.network_subnets

  description = each.value.description
  type        = each.value.type
  subnet      = each.value.subnet
  domain      = each.value.domain
  network_item_id = each.value.network_id == "production" ? cloudconnexa_network.production.id : (
    each.value.network_id == "staging" ? cloudconnexa_network.staging.id : cloudconnexa_network.development.id
  )
}

# Complex routing scenario with conditional logic
locals {
  # Define environment-specific CIDR blocks
  env_cidrs = {
    production  = "10.0.0.0/16"
    staging     = "10.1.0.0/16"
    development = "10.2.0.0/16"
  }

  # Define service-specific subnets
  service_subnets = {
    web      = "/24"
    app      = "/24"
    database = "/24"
    cache    = "/28"
  }
}

# Service-specific routes for production
resource "cloudconnexa_route" "prod_service_routes" {
  for_each = local.service_subnets

  description     = "Production ${each.key} service subnet"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.production.id
  subnet          = "${cidrsubnet(local.env_cidrs.production, tonumber(substr(each.value, 1, -1)) - 16, index(keys(local.service_subnets), each.key))}${each.value}"
}

# Route with specific use cases
resource "cloudconnexa_route" "vpn_client_range" {
  description     = "VPN client IP range for production access"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.production.id
  subnet          = "172.16.0.0/12"
}

resource "cloudconnexa_route" "kubernetes_pods" {
  description     = "Kubernetes pod network range"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.production.id
  subnet          = "192.168.0.0/16"
}

resource "cloudconnexa_route" "docker_bridge" {
  description     = "Docker bridge network"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.development.id
  subnet          = "172.17.0.0/16"
}

# Content Delivery Network domains
resource "cloudconnexa_route" "cdn_domains" {
  description     = "CDN domain routing"
  type            = "DOMAIN"
  network_item_id = cloudconnexa_network.production.id
  domain          = "*.cdn.example.com"
}

# Third-party integration domains
resource "cloudconnexa_route" "external_apis" {
  description     = "External API integrations"
  type            = "DOMAIN"
  network_item_id = cloudconnexa_network.production.id
  domain          = "*.partners.example.com"
}

# Outputs
output "network_routes" {
  description = "Network-based routes created"
  value = {
    prod_internal_network = cloudconnexa_route.prod_internal_network.id
    prod_database_subnet  = cloudconnexa_route.prod_database_subnet.id
    staging_network       = cloudconnexa_route.staging_network.id
    dev_network           = cloudconnexa_route.dev_network.id
    prod_ipv6_network     = cloudconnexa_route.prod_ipv6_network.id
    staging_ipv6_network  = cloudconnexa_route.staging_ipv6_network.id
  }
}

output "domain_routes" {
  description = "Domain-based routes created"
  value = {
    prod_domain_route    = cloudconnexa_route.prod_domain_route.id
    api_domain_route     = cloudconnexa_route.api_domain_route.id
    staging_domain_route = cloudconnexa_route.staging_domain_route.id
    dev_domain_route     = cloudconnexa_route.dev_domain_route.id
    cdn_domains          = cloudconnexa_route.cdn_domains.id
    external_apis        = cloudconnexa_route.external_apis.id
  }
}

output "host_routes" {
  description = "Host-based routes created"
  value = {
    database_route   = cloudconnexa_route.database_route.id
    api_server_route = cloudconnexa_route.api_server_route.id
  }
}

output "dynamic_routes" {
  description = "Dynamically created routes"
  value = {
    network_subnets     = { for k, v in cloudconnexa_route.network_subnets : k => v.id }
    prod_service_routes = { for k, v in cloudconnexa_route.prod_service_routes : k => v.id }
  }
}

output "special_purpose_routes" {
  description = "Special purpose routes for infrastructure"
  value = {
    vpn_client_range = cloudconnexa_route.vpn_client_range.id
    kubernetes_pods  = cloudconnexa_route.kubernetes_pods.id
    docker_bridge    = cloudconnexa_route.docker_bridge.id
  }
}
