data "cloudconnexa_settings" "current" {
  # This data source reads the current CloudConnexa settings
}

# Example of using the data source output
output "current_dns_settings" {
  value = {
    dns_proxy_enabled  = data.cloudconnexa_settings.current.dns_proxy_enabled
    default_dns_suffix = data.cloudconnexa_settings.current.default_dns_suffix
    dns_servers        = data.cloudconnexa_settings.current.dns_servers
  }
}

output "security_settings" {
  value = {
    allow_trusted_devices = data.cloudconnexa_settings.current.allow_trusted_devices
    two_factor_auth       = data.cloudconnexa_settings.current.two_factor_auth
    connect_auth          = data.cloudconnexa_settings.current.connect_auth
    device_enforcement    = data.cloudconnexa_settings.current.device_enforcement
  }
}

output "network_topology" {
  value = {
    topology       = data.cloudconnexa_settings.current.topology
    default_region = data.cloudconnexa_settings.current.default_region
    snat           = data.cloudconnexa_settings.current.snat
  }
}
