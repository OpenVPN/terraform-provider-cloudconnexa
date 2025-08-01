---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudconnexa_settings Resource - terraform-provider-cloudconnexa"
subcategory: ""
description: |-
  Use cloudconnexa_settings to define settings
---

# cloudconnexa_settings (Resource)

Use `cloudconnexa_settings` to define settings

## Example Usage

```terraform
# 1. Production environment settings with high security
resource "cloudconnexa_settings" "production" {
  # Authentication and security settings for production
  allow_trusted_devices = false # Strict security for production
  two_factor_auth       = true
  connect_auth          = "ALWAYS_REQUIRED" # Always require authentication

  # DNS configuration for production
  dns_proxy_enabled  = true
  default_dns_suffix = "prod.example.com"

  dns_servers {
    primary_ip_v4   = "8.8.8.8"
    secondary_ip_v4 = "8.8.4.4"
    primary_ip_v6   = "2001:4860:4860::8888"
    secondary_ip_v6 = "2001:4860:4860::8844"
  }

  # Internal DNS zones for production services
  dns_zones {
    name      = "internal.prod.example.com"
    addresses = ["10.100.1.10", "10.100.1.11", "10.100.1.12"]
  }

  dns_zones {
    name      = "db.prod.example.com"
    addresses = ["10.100.2.10", "10.100.2.11"]
  }

  dns_zones {
    name      = "api.prod.example.com"
    addresses = ["10.100.3.10", "10.100.3.11"]
  }

  dns_zones {
    name      = "monitoring.prod.example.com"
    addresses = ["10.100.4.10"]
  }

  # Restrictive device management for production
  device_allowance_per_user     = 2
  device_allowance_force_update = true
  device_enforcement            = "ENFORCE"

  # Profile and connection settings
  profile_distribution = "MANUAL" # Manual distribution for production control
  connection_timeout   = 600      # 10 minutes for production

  # Production client options
  client_options = [
    "block-outside-dns",
    "dhcp-option DNS 8.8.8.8",
    "dhcp-option DNS 8.8.4.4",
    "persist-tun",
    "persist-key",
    "comp-lzo",
    "verb 3"
  ]

  # Full mesh topology for production reliability
  topology       = "FULL_MESH"
  default_region = "us-east-1"
  snat           = true

  # Production domain routing subnet
  domain_routing_subnet {
    ip_v4_address = "192.168.200.0/24"
    ip_v6_address = "2001:db8:200::/64"
  }

  # Production subnet configuration
  subnet {
    ip_v4_address = ["10.100.0.0/16", "172.20.0.0/16"]
    ip_v6_address = ["2001:db8:100::/48"]
  }
}

# 2. Staging environment settings with moderate security
resource "cloudconnexa_settings" "staging" {
  # Authentication and security settings for staging
  allow_trusted_devices = true
  two_factor_auth       = true
  connect_auth          = "ON_PRIOR_AUTH"

  # DNS configuration for staging
  dns_proxy_enabled  = true
  default_dns_suffix = "staging.example.com"

  dns_servers {
    primary_ip_v4   = "1.1.1.1"
    secondary_ip_v4 = "1.0.0.1"
    primary_ip_v6   = "2606:4700:4700::1111"
    secondary_ip_v6 = "2606:4700:4700::1001"
  }

  # Staging DNS zones
  dns_zones {
    name      = "internal.staging.example.com"
    addresses = ["10.110.1.10", "10.110.1.11"]
  }

  dns_zones {
    name      = "db.staging.example.com"
    addresses = ["10.110.2.10"]
  }

  dns_zones {
    name      = "api.staging.example.com"
    addresses = ["10.110.3.10"]
  }

  # Moderate device management for staging
  device_allowance_per_user     = 3
  device_allowance_force_update = true
  device_enforcement            = "LEARN_AND_ENFORCE"

  # Profile and connection settings
  profile_distribution = "AUTOMATIC"
  connection_timeout   = 450

  # Staging client options
  client_options = [
    "block-outside-dns",
    "dhcp-option DNS 1.1.1.1",
    "dhcp-option DNS 1.0.0.1",
    "persist-tun",
    "comp-lzo"
  ]

  # Hub and spoke topology for staging efficiency
  topology       = "HUB_AND_SPOKE"
  default_region = "us-west-1"
  snat           = true

  # Staging domain routing subnet
  domain_routing_subnet {
    ip_v4_address = "192.168.110.0/24"
    ip_v6_address = "2001:db8:110::/64"
  }

  # Staging subnet configuration
  subnet {
    ip_v4_address = ["10.110.0.0/16"]
    ip_v6_address = ["2001:db8:110::/48"]
  }
}

# 3. Development environment settings with relaxed security
resource "cloudconnexa_settings" "development" {
  # Authentication and security settings for development
  allow_trusted_devices = true
  two_factor_auth       = false # Relaxed for development
  connect_auth          = "AUTO_LOGIN"

  # DNS configuration for development
  dns_proxy_enabled  = true
  default_dns_suffix = "dev.example.com"

  dns_servers {
    primary_ip_v4   = "8.8.8.8"
    secondary_ip_v4 = "1.1.1.1"
  }

  # Development DNS zones
  dns_zones {
    name      = "internal.dev.example.com"
    addresses = ["10.120.1.10"]
  }

  dns_zones {
    name      = "db.dev.example.com"
    addresses = ["10.120.2.10"]
  }

  dns_zones {
    name      = "localhost.dev.example.com"
    addresses = ["127.0.0.1"]
  }

  # Relaxed device management for development
  device_allowance_per_user     = 10
  device_allowance_force_update = false
  device_enforcement            = "LEARN_ONLY"

  # Profile and connection settings
  profile_distribution = "AUTOMATIC"
  connection_timeout   = 180 # Shorter timeout for development

  # Development client options
  client_options = [
    "dhcp-option DNS 8.8.8.8",
    "dhcp-option DNS 1.1.1.1",
    "persist-tun"
  ]

  # Hub and spoke topology for development simplicity
  topology       = "HUB_AND_SPOKE"
  default_region = "us-west-2"
  snat           = false # Direct routing for development

  # Development domain routing subnet
  domain_routing_subnet {
    ip_v4_address = "192.168.120.0/24"
    ip_v6_address = "2001:db8:120::/64"
  }

  # Development subnet configuration
  subnet {
    ip_v4_address = ["10.120.0.0/16", "172.30.0.0/16"]
    ip_v6_address = ["2001:db8:120::/48"]
  }
}

# 4. Enterprise settings with comprehensive security and multiple regions
resource "cloudconnexa_settings" "enterprise" {
  # Maximum security settings for enterprise
  allow_trusted_devices = false
  two_factor_auth       = true
  connect_auth          = "ALWAYS_REQUIRED"

  # Enterprise DNS configuration with custom servers
  dns_proxy_enabled  = true
  default_dns_suffix = "corp.enterprise.com"

  dns_servers {
    primary_ip_v4   = "10.0.1.53" # Internal DNS server
    secondary_ip_v4 = "10.0.2.53" # Secondary internal DNS
    primary_ip_v6   = "2001:db8:1::53"
    secondary_ip_v6 = "2001:db8:2::53"
  }

  # Comprehensive DNS zones for enterprise
  dns_zones {
    name      = "internal.corp.enterprise.com"
    addresses = ["10.0.10.10", "10.0.10.11", "10.0.10.12"]
  }

  dns_zones {
    name      = "database.corp.enterprise.com"
    addresses = ["10.0.20.10", "10.0.20.11", "10.0.20.12"]
  }

  dns_zones {
    name      = "applications.corp.enterprise.com"
    addresses = ["10.0.30.10", "10.0.30.11", "10.0.30.12"]
  }

  dns_zones {
    name      = "monitoring.corp.enterprise.com"
    addresses = ["10.0.40.10", "10.0.40.11"]
  }

  dns_zones {
    name      = "backup.corp.enterprise.com"
    addresses = ["10.0.50.10", "10.0.50.11"]
  }

  dns_zones {
    name      = "partner.corp.enterprise.com"
    addresses = ["10.0.60.10"]
  }

  # Strict device management for enterprise
  device_allowance_per_user     = 1
  device_allowance_force_update = true
  device_enforcement            = "ENFORCE"

  # Enterprise profile and connection settings
  profile_distribution = "MANUAL"
  connection_timeout   = 900 # 15 minutes for enterprise

  # Enterprise client options with security focus
  client_options = [
    "block-outside-dns",
    "dhcp-option DNS 10.0.1.53",
    "dhcp-option DNS 10.0.2.53",
    "persist-tun",
    "persist-key",
    "comp-lzo",
    "cipher AES-256-GCM",
    "auth SHA256",
    "tls-version-min 1.2",
    "verb 1"
  ]

  # Full mesh topology for enterprise reliability
  topology       = "FULL_MESH"
  default_region = "us-east-1"
  snat           = true

  # Enterprise domain routing subnet
  domain_routing_subnet {
    ip_v4_address = "192.168.0.0/24"
    ip_v6_address = "2001:db8::/64"
  }

  # Enterprise subnet configuration
  subnet {
    ip_v4_address = ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"]
    ip_v6_address = ["2001:db8::/32"]
  }
}

# 5. Remote work settings optimized for distributed teams
resource "cloudconnexa_settings" "remote_work" {
  # Balanced security for remote work
  allow_trusted_devices = true
  two_factor_auth       = true
  connect_auth          = "ON_PRIOR_AUTH"

  # DNS configuration optimized for remote workers
  dns_proxy_enabled  = true
  default_dns_suffix = "remote.company.com"

  dns_servers {
    primary_ip_v4   = "1.1.1.1" # Fast public DNS
    secondary_ip_v4 = "8.8.8.8" # Reliable backup
    primary_ip_v6   = "2606:4700:4700::1111"
    secondary_ip_v6 = "2001:4860:4860::8888"
  }

  # Remote work DNS zones
  dns_zones {
    name      = "office.company.com"
    addresses = ["10.200.1.10", "10.200.1.11"]
  }

  dns_zones {
    name      = "resources.company.com"
    addresses = ["10.200.2.10"]
  }

  dns_zones {
    name      = "collaboration.company.com"
    addresses = ["10.200.3.10", "10.200.3.11"]
  }

  dns_zones {
    name      = "vpn.company.com"
    addresses = ["10.200.4.10"]
  }

  # Flexible device management for remote work
  device_allowance_per_user     = 5 # Multiple devices for remote workers
  device_allowance_force_update = true
  device_enforcement            = "LEARN_AND_ENFORCE"

  # Remote work profile and connection settings
  profile_distribution = "AUTOMATIC"
  connection_timeout   = 3600 # 1 hour for long remote sessions

  # Remote work client options
  client_options = [
    "dhcp-option DNS 1.1.1.1",
    "dhcp-option DNS 8.8.8.8",
    "persist-tun",
    "persist-key",
    "comp-lzo",
    "fast-io",
    "sndbuf 0",
    "rcvbuf 0"
  ]

  # Hub and spoke for remote work efficiency
  topology       = "HUB_AND_SPOKE"
  default_region = "us-central-1"
  snat           = true

  # Remote work domain routing subnet
  domain_routing_subnet {
    ip_v4_address = "192.168.200.0/24"
    ip_v6_address = "2001:db8:200::/64"
  }

  # Remote work subnet configuration
  subnet {
    ip_v4_address = ["10.200.0.0/16"]
    ip_v6_address = ["2001:db8:200::/48"]
  }
}

# 6. Guest network settings with restricted access
resource "cloudconnexa_settings" "guest_network" {
  # Minimal security for guest access
  allow_trusted_devices = false
  two_factor_auth       = false
  connect_auth          = "AUTO_LOGIN"

  # Basic DNS configuration for guests
  dns_proxy_enabled  = false # Direct DNS for guests
  default_dns_suffix = "guest.example.com"

  dns_servers {
    primary_ip_v4   = "8.8.8.8"
    secondary_ip_v4 = "8.8.4.4"
  }

  # Limited DNS zones for guests
  dns_zones {
    name      = "portal.guest.example.com"
    addresses = ["10.250.1.10"]
  }

  dns_zones {
    name      = "wifi.guest.example.com"
    addresses = ["10.250.2.10"]
  }

  # Restrictive device management for guests
  device_allowance_per_user     = 1
  device_allowance_force_update = false
  device_enforcement            = "LEARN_ONLY"

  # Guest profile and connection settings
  profile_distribution = "AUTOMATIC"
  connection_timeout   = 120 # Short timeout for guests

  # Minimal client options for guests
  client_options = [
    "dhcp-option DNS 8.8.8.8",
    "dhcp-option DNS 8.8.4.4"
  ]

  # Simple topology for guest network
  topology       = "HUB_AND_SPOKE"
  default_region = "us-west-1"
  snat           = true

  # Guest domain routing subnet
  domain_routing_subnet {
    ip_v4_address = "192.168.250.0/24"
  }

  # Guest subnet configuration
  subnet {
    ip_v4_address = ["10.250.0.0/16"]
  }
}

# 7. High availability settings with multiple DNS zones and regions
resource "cloudconnexa_settings" "high_availability" {
  # High security for critical infrastructure
  allow_trusted_devices = false
  two_factor_auth       = true
  connect_auth          = "ALWAYS_REQUIRED"

  # Redundant DNS configuration
  dns_proxy_enabled  = true
  default_dns_suffix = "ha.critical.com"

  dns_servers {
    primary_ip_v4   = "8.8.8.8"
    secondary_ip_v4 = "1.1.1.1"
    primary_ip_v6   = "2001:4860:4860::8888"
    secondary_ip_v6 = "2606:4700:4700::1111"
  }

  # Multiple DNS zones for high availability
  dns_zones {
    name      = "primary.ha.critical.com"
    addresses = ["10.10.1.10", "10.10.1.11", "10.10.1.12"]
  }

  dns_zones {
    name      = "secondary.ha.critical.com"
    addresses = ["10.20.1.10", "10.20.1.11", "10.20.1.12"]
  }

  dns_zones {
    name      = "database.ha.critical.com"
    addresses = ["10.10.2.10", "10.10.2.11", "10.20.2.10", "10.20.2.11"]
  }

  dns_zones {
    name      = "loadbalancer.ha.critical.com"
    addresses = ["10.10.3.10", "10.20.3.10"]
  }

  dns_zones {
    name      = "monitor.ha.critical.com"
    addresses = ["10.10.4.10", "10.20.4.10"]
  }

  dns_zones {
    name      = "backup.ha.critical.com"
    addresses = ["10.30.1.10", "10.30.1.11"]
  }

  # Strict device management for HA
  device_allowance_per_user     = 2
  device_allowance_force_update = true
  device_enforcement            = "ENFORCE"

  # HA profile and connection settings
  profile_distribution = "MANUAL"
  connection_timeout   = 1800 # 30 minutes for critical systems

  # HA client options with maximum reliability
  client_options = [
    "block-outside-dns",
    "dhcp-option DNS 8.8.8.8",
    "dhcp-option DNS 1.1.1.1",
    "persist-tun",
    "persist-key",
    "keepalive 10 60",
    "ping-timer-rem",
    "comp-lzo",
    "cipher AES-256-GCM",
    "auth SHA256",
    "tls-version-min 1.3",
    "verb 3"
  ]

  # Full mesh topology for maximum availability
  topology       = "FULL_MESH"
  default_region = "us-east-1"
  snat           = true

  # HA domain routing subnet with redundancy
  domain_routing_subnet {
    ip_v4_address = "192.168.10.0/24"
    ip_v6_address = "2001:db8:10::/64"
  }

  # HA subnet configuration with multiple ranges
  subnet {
    ip_v4_address = ["10.10.0.0/16", "10.20.0.0/16", "10.30.0.0/16"]
    ip_v6_address = ["2001:db8:10::/48", "2001:db8:20::/48"]
  }
}

# Outputs for different environment settings
output "production_settings" {
  description = "Production environment settings ID"
  value       = cloudconnexa_settings.production.id
}

output "staging_settings" {
  description = "Staging environment settings ID"
  value       = cloudconnexa_settings.staging.id
}

output "development_settings" {
  description = "Development environment settings ID"
  value       = cloudconnexa_settings.development.id
}

output "enterprise_settings" {
  description = "Enterprise settings ID"
  value       = cloudconnexa_settings.enterprise.id
}

output "remote_work_settings" {
  description = "Remote work settings ID"
  value       = cloudconnexa_settings.remote_work.id
}

output "guest_network_settings" {
  description = "Guest network settings ID"
  value       = cloudconnexa_settings.guest_network.id
}

output "high_availability_settings" {
  description = "High availability settings ID"
  value       = cloudconnexa_settings.high_availability.id
}

output "settings_summary" {
  description = "Summary of all CloudConnexa settings configurations"
  value = {
    total_configurations = 7

    by_security_level = {
      high_security = [
        cloudconnexa_settings.production.id,
        cloudconnexa_settings.enterprise.id,
        cloudconnexa_settings.high_availability.id
      ]
      medium_security = [
        cloudconnexa_settings.staging.id,
        cloudconnexa_settings.remote_work.id
      ]
      low_security = [
        cloudconnexa_settings.development.id,
        cloudconnexa_settings.guest_network.id
      ]
    }

    by_environment_type = {
      production  = cloudconnexa_settings.production.id
      staging     = cloudconnexa_settings.staging.id
      development = cloudconnexa_settings.development.id
      enterprise  = cloudconnexa_settings.enterprise.id
    }

    by_use_case = {
      corporate_office        = cloudconnexa_settings.enterprise.id
      remote_workforce        = cloudconnexa_settings.remote_work.id
      guest_access            = cloudconnexa_settings.guest_network.id
      critical_infrastructure = cloudconnexa_settings.high_availability.id
    }

    by_topology = {
      full_mesh = [
        cloudconnexa_settings.production.id,
        cloudconnexa_settings.enterprise.id,
        cloudconnexa_settings.high_availability.id
      ]
      hub_and_spoke = [
        cloudconnexa_settings.staging.id,
        cloudconnexa_settings.development.id,
        cloudconnexa_settings.remote_work.id,
        cloudconnexa_settings.guest_network.id
      ]
    }

    by_authentication = {
      always_required = [
        cloudconnexa_settings.production.id,
        cloudconnexa_settings.enterprise.id,
        cloudconnexa_settings.high_availability.id
      ]
      on_prior_auth = [
        cloudconnexa_settings.staging.id,
        cloudconnexa_settings.remote_work.id
      ]
      auto_login = [
        cloudconnexa_settings.development.id,
        cloudconnexa_settings.guest_network.id
      ]
    }

    dns_configuration = {
      internal_dns_servers = [
        cloudconnexa_settings.enterprise.id
      ]
      public_dns_primary = [
        cloudconnexa_settings.production.id,
        cloudconnexa_settings.development.id,
        cloudconnexa_settings.guest_network.id,
        cloudconnexa_settings.high_availability.id
      ]
      cloudflare_dns_primary = [
        cloudconnexa_settings.staging.id,
        cloudconnexa_settings.remote_work.id
      ]
    }

    device_management = {
      strict_single_device = [
        cloudconnexa_settings.enterprise.id,
        cloudconnexa_settings.guest_network.id
      ]
      limited_devices = [
        cloudconnexa_settings.production.id,
        cloudconnexa_settings.high_availability.id
      ]
      moderate_devices = [
        cloudconnexa_settings.staging.id
      ]
      flexible_devices = [
        cloudconnexa_settings.remote_work.id,
        cloudconnexa_settings.development.id
      ]
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `access_visibility_enabled` (Boolean)
- `allow_trusted_devices` (Boolean)
- `client_options` (List of String)
- `connect_auth` (String)
- `connection_timeout` (Number)
- `default_dns_suffix` (String)
- `default_region` (String)
- `device_allowance_force_update` (Boolean)
- `device_allowance_per_user` (Number)
- `device_enforcement` (String)
- `dns_log_enabled` (Boolean)
- `dns_proxy_enabled` (Boolean)
- `dns_servers` (Block List, Max: 1) (see [below for nested schema](#nestedblock--dns_servers))
- `dns_zones` (Block List) (see [below for nested schema](#nestedblock--dns_zones))
- `domain_routing_subnet` (Block List, Max: 1) (see [below for nested schema](#nestedblock--domain_routing_subnet))
- `profile_distribution` (String)
- `snat` (Boolean)
- `subnet` (Block List, Max: 1) (see [below for nested schema](#nestedblock--subnet))
- `topology` (String)
- `two_factor_auth` (Boolean)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--dns_servers"></a>
### Nested Schema for `dns_servers`

Required:

- `primary_ip_v4` (String)

Optional:

- `secondary_ip_v4` (String)


<a id="nestedblock--dns_zones"></a>
### Nested Schema for `dns_zones`

Required:

- `name` (String)

Optional:

- `addresses` (List of String)


<a id="nestedblock--domain_routing_subnet"></a>
### Nested Schema for `domain_routing_subnet`

Required:

- `ip_v4_address` (String)

Optional:

- `ip_v6_address` (String)


<a id="nestedblock--subnet"></a>
### Nested Schema for `subnet`

Optional:

- `ip_v4_address` (List of String)
- `ip_v6_address` (List of String)

## Import

Import is supported using the following syntax:

```shell
# Settings is a singleton resource, so the import ID can be any value
# The actual ID will always be set to "settings" during import
terraform import cloudconnexa_settings.example settings
```
