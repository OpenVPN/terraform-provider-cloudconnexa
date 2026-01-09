# Create a host for the route example
resource "cloudconnexa_host" "example" {
  name            = "example-host"
  description     = "Example host for route"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "example.internal.com"
}

# Create an IPv4 route on a host
resource "cloudconnexa_host_route" "internal_network" {
  host_id     = cloudconnexa_host.example.id
  subnet      = "10.0.0.0/24"
  description = "Internal network route"
}

# Create multiple routes for different subnets
resource "cloudconnexa_host_route" "database_subnet" {
  host_id     = cloudconnexa_host.example.id
  subnet      = "10.0.100.0/24"
  description = "Database subnet route"
}

resource "cloudconnexa_host_route" "application_subnet" {
  host_id     = cloudconnexa_host.example.id
  subnet      = "10.0.200.0/24"
  description = "Application subnet route"
}

# Using for_each for multiple routes
variable "host_subnets" {
  description = "Map of subnets to create routes for"
  type = map(object({
    subnet      = string
    description = string
  }))
  default = {
    "web" = {
      subnet      = "10.1.0.0/24"
      description = "Web tier subnet"
    }
    "api" = {
      subnet      = "10.2.0.0/24"
      description = "API tier subnet"
    }
    "backend" = {
      subnet      = "10.3.0.0/24"
      description = "Backend tier subnet"
    }
  }
}

resource "cloudconnexa_host_route" "dynamic" {
  for_each = var.host_subnets

  host_id     = cloudconnexa_host.example.id
  subnet      = each.value.subnet
  description = each.value.description
}

# Outputs
output "route_ids" {
  value = {
    internal_network   = cloudconnexa_host_route.internal_network.id
    database_subnet    = cloudconnexa_host_route.database_subnet.id
    application_subnet = cloudconnexa_host_route.application_subnet.id
  }
}

output "dynamic_routes" {
  value = { for k, v in cloudconnexa_host_route.dynamic : k => v.id }
}
