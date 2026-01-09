# Get all routes for a specific host
data "cloudconnexa_host_routes" "example" {
  host_id = cloudconnexa_host.example.id
}

# Output routes information
output "host_routes" {
  value = data.cloudconnexa_host_routes.example.routes
}

output "ipv4_routes" {
  value = [
    for route in data.cloudconnexa_host_routes.example.routes : route
    if route.type == "IP_V4"
  ]
}
