---
page_title: "cloudconnexa_host_connector Resource - terraform-provider-cloudconnexa"
subcategory: ""
description: |-
  Use cloudconnexa_connector to create an CloudConnexa connector.
  ~> NOTE: This only creates the CloudConnexa connector object. Additional manual steps are required to associate a host in your infrastructure with the connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.
---

# cloudconnexa_host_connector (Resource)

Use `cloudconnexa_connector` to create an CloudConnexa connector.

~> NOTE: This only creates the CloudConnexa connector object. Additional manual steps are required to associate a host in your infrastructure with the connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.

~> **WARNING: Status Field Limitation** The `status` field is **write-only**. The CloudConnexa API does not return the connector status in GET responses. This means:
- Terraform can suspend/activate connectors via the `status` field
- Terraform **cannot detect** external status changes (e.g., via CloudConnexa UI)
- Default value is `ACTIVE`. If you need to suspend, explicitly set `status = "SUSPENDED"`

## Example Usage

```terraform
# Hosts for different connector scenarios
resource "cloudconnexa_host" "onprem_datacenter" {
  name            = "onprem-datacenter"
  description     = "On-premises data center host"
  internet_access = "RESTRICTED_INTERNET"
  domain          = "datacenter.onprem.example.com"
}

resource "cloudconnexa_host" "branch_office_ny" {
  name            = "branch-office-new-york"
  description     = "New York branch office host"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "ny.branch.example.com"
}

resource "cloudconnexa_host" "branch_office_london" {
  name            = "branch-office-london"
  description     = "London branch office host"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "london.branch.example.com"
}

resource "cloudconnexa_host" "cloud_gateway" {
  name            = "cloud-gateway"
  description     = "Cloud services gateway host"
  internet_access = "SPLIT_TUNNEL_OFF"
  domain          = "gateway.cloud.example.com"
}

resource "cloudconnexa_host" "backup_site" {
  name            = "disaster-recovery-site"
  description     = "Disaster recovery backup site"
  internet_access = "RESTRICTED_INTERNET"
  domain          = "dr.backup.example.com"
}

resource "cloudconnexa_host" "edge_location" {
  name            = "edge-computing-node"
  description     = "Edge computing location"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "edge.compute.example.com"
}

resource "cloudconnexa_host" "partner_network" {
  name            = "partner-integration"
  description     = "Partner network integration point"
  internet_access = "RESTRICTED_INTERNET"
  domain          = "partner.integration.example.com"
}

# 1. Primary data center connector (US East)
resource "cloudconnexa_host_connector" "primary_datacenter" {
  name          = "primary-datacenter-connector"
  description   = "Primary data center connector in US East"
  host_id       = cloudconnexa_host.onprem_datacenter.id
  vpn_region_id = "us-east-1"
}

# 2. Branch office connectors for different regions
resource "cloudconnexa_host_connector" "ny_branch_connector" {
  name          = "new-york-branch-connector"
  description   = "New York branch office connector"
  host_id       = cloudconnexa_host.branch_office_ny.id
  vpn_region_id = "us-east-1"
}

resource "cloudconnexa_host_connector" "london_branch_connector" {
  name          = "london-branch-connector"
  description   = "London branch office connector"
  host_id       = cloudconnexa_host.branch_office_london.id
  vpn_region_id = "eu-west-1"
}

# 3. Multi-region cloud gateway connectors
resource "cloudconnexa_host_connector" "cloud_gateway_us" {
  name          = "cloud-gateway-us-connector"
  description   = "Cloud gateway connector for US region"
  host_id       = cloudconnexa_host.cloud_gateway.id
  vpn_region_id = "us-west-1"
}

resource "cloudconnexa_host_connector" "cloud_gateway_eu" {
  name          = "cloud-gateway-eu-connector"
  description   = "Cloud gateway connector for EU region"
  host_id       = cloudconnexa_host.cloud_gateway.id
  vpn_region_id = "eu-central-1"
}

# 4. Disaster recovery site connector
resource "cloudconnexa_host_connector" "dr_site_connector" {
  name          = "disaster-recovery-connector"
  description   = "Disaster recovery site connector"
  host_id       = cloudconnexa_host.backup_site.id
  vpn_region_id = "us-west-2"
}

# 4a. Suspended connector example - temporarily disable a connector
resource "cloudconnexa_host_connector" "maintenance_connector" {
  name          = "maintenance-connector"
  description   = "Connector under maintenance - suspended"
  host_id       = cloudconnexa_host.backup_site.id
  vpn_region_id = "us-west-2"
  status        = "SUSPENDED" # Connector is suspended during maintenance
}

# 4b. Connector that can be toggled between active and suspended
variable "dr_connector_active" {
  description = "Whether the DR connector should be active"
  type        = bool
  default     = false # DR connector is suspended by default until failover
}

resource "cloudconnexa_host_connector" "dr_failover_connector" {
  name          = "dr-failover-connector"
  description   = "DR connector - activated only during failover"
  host_id       = cloudconnexa_host.backup_site.id
  vpn_region_id = "us-west-2"
  status        = var.dr_connector_active ? "ACTIVE" : "SUSPENDED"
}

# 5. Edge computing connectors
resource "cloudconnexa_host_connector" "edge_us_connector" {
  name          = "edge-us-connector"
  description   = "Edge computing connector for US region"
  host_id       = cloudconnexa_host.edge_location.id
  vpn_region_id = "us-west-1"
}

resource "cloudconnexa_host_connector" "edge_eu_connector" {
  name          = "edge-eu-connector"
  description   = "Edge computing connector for EU region"
  host_id       = cloudconnexa_host.edge_location.id
  vpn_region_id = "eu-central-1"
}

resource "cloudconnexa_host_connector" "edge_apac_connector" {
  name          = "edge-apac-connector"
  description   = "Edge computing connector for APAC region"
  host_id       = cloudconnexa_host.edge_location.id
  vpn_region_id = "ap-southeast-1"
}

# 6. Partner network integration connector
resource "cloudconnexa_host_connector" "partner_connector" {
  name          = "partner-network-connector"
  description   = "Partner network integration connector"
  host_id       = cloudconnexa_host.partner_network.id
  vpn_region_id = "us-east-1"
}

# Multiple connectors using for_each pattern
variable "regional_connectors" {
  description = "Regional connector configurations"
  type = map(object({
    description   = string
    host_id       = string
    vpn_region_id = string
  }))
  default = {
    "frankfurt-connector" = {
      description   = "Frankfurt office connector"
      host_id       = "branch_office_london"
      vpn_region_id = "eu-central-1"
    }
    "singapore-connector" = {
      description   = "Singapore office connector"
      host_id       = "edge_location"
      vpn_region_id = "ap-southeast-1"
    }
    "tokyo-connector" = {
      description   = "Tokyo office connector"
      host_id       = "edge_location"
      vpn_region_id = "ap-northeast-1"
    }
    "sydney-connector" = {
      description   = "Sydney office connector"
      host_id       = "branch_office_ny"
      vpn_region_id = "ap-southeast-2"
    }
  }
}

resource "cloudconnexa_host_connector" "regional_connectors" {
  for_each = var.regional_connectors

  name        = each.key
  description = each.value.description
  host_id = each.value.host_id == "onprem_datacenter" ? cloudconnexa_host.onprem_datacenter.id : (
    each.value.host_id == "branch_office_ny" ? cloudconnexa_host.branch_office_ny.id : (
      each.value.host_id == "branch_office_london" ? cloudconnexa_host.branch_office_london.id : (
        each.value.host_id == "cloud_gateway" ? cloudconnexa_host.cloud_gateway.id : (
          each.value.host_id == "edge_location" ? cloudconnexa_host.edge_location.id : cloudconnexa_host.partner_network.id
        )
      )
    )
  )
  vpn_region_id = each.value.vpn_region_id
}

# Redundant connector pairs for high availability
locals {
  ha_connector_pairs = {
    "primary-ha" = {
      primary = {
        name          = "primary-ha-connector-1"
        description   = "Primary HA connector (main)"
        host_id       = cloudconnexa_host.onprem_datacenter.id
        vpn_region_id = "us-east-1"
      }
      secondary = {
        name          = "primary-ha-connector-2"
        description   = "Primary HA connector (backup)"
        host_id       = cloudconnexa_host.backup_site.id
        vpn_region_id = "us-east-1"
      }
    }
    "europe-ha" = {
      primary = {
        name          = "europe-ha-connector-1"
        description   = "Europe HA connector (main)"
        host_id       = cloudconnexa_host.branch_office_london.id
        vpn_region_id = "eu-west-1"
      }
      secondary = {
        name          = "europe-ha-connector-2"
        description   = "Europe HA connector (backup)"
        host_id       = cloudconnexa_host.branch_office_london.id
        vpn_region_id = "eu-central-1"
      }
    }
  }
}

resource "cloudconnexa_host_connector" "ha_primary_connectors" {
  for_each = local.ha_connector_pairs

  name          = each.value.primary.name
  description   = each.value.primary.description
  host_id       = each.value.primary.host_id
  vpn_region_id = each.value.primary.vpn_region_id
}

resource "cloudconnexa_host_connector" "ha_secondary_connectors" {
  for_each = local.ha_connector_pairs

  name          = each.value.secondary.name
  description   = each.value.secondary.description
  host_id       = each.value.secondary.host_id
  vpn_region_id = each.value.secondary.vpn_region_id
}

# Environment-specific connectors
variable "environment_connectors" {
  description = "Environment-specific connector configurations"
  type = map(object({
    description = string
    region      = string
  }))
  default = {
    "production-primary" = {
      description = "Production primary connector"
      region      = "us-east-1"
    }
    "production-backup" = {
      description = "Production backup connector"
      region      = "us-west-1"
    }
    "staging-connector" = {
      description = "Staging environment connector"
      region      = "eu-central-1"
    }
    "development-connector" = {
      description = "Development environment connector"
      region      = "eu-west-1"
    }
  }
}

resource "cloudconnexa_host_connector" "environment_connectors" {
  for_each = var.environment_connectors

  name        = each.key
  description = each.value.description
  host_id = contains(["production-primary", "production-backup"], each.key) ? cloudconnexa_host.onprem_datacenter.id : (
    each.key == "staging-connector" ? cloudconnexa_host.branch_office_london.id : cloudconnexa_host.edge_location.id
  )
  vpn_region_id = each.value.region
}

# Outputs
output "primary_connectors" {
  description = "Primary data center and infrastructure connectors"
  value = {
    primary_datacenter = cloudconnexa_host_connector.primary_datacenter.id
    cloud_gateway_us   = cloudconnexa_host_connector.cloud_gateway_us.id
    cloud_gateway_eu   = cloudconnexa_host_connector.cloud_gateway_eu.id
    dr_site_connector  = cloudconnexa_host_connector.dr_site_connector.id
  }
}

output "branch_office_connectors" {
  description = "Branch office connectors"
  value = {
    ny_branch         = cloudconnexa_host_connector.ny_branch_connector.id
    london_branch     = cloudconnexa_host_connector.london_branch_connector.id
    partner_connector = cloudconnexa_host_connector.partner_connector.id
  }
}

output "edge_connectors" {
  description = "Edge computing connectors"
  value = {
    edge_us   = cloudconnexa_host_connector.edge_us_connector.id
    edge_eu   = cloudconnexa_host_connector.edge_eu_connector.id
    edge_apac = cloudconnexa_host_connector.edge_apac_connector.id
  }
}

output "regional_connectors" {
  description = "Regional connectors created with for_each"
  value       = { for k, v in cloudconnexa_host_connector.regional_connectors : k => v.id }
}

output "ha_connectors" {
  description = "High availability connector pairs"
  value = {
    ha_primary   = { for k, v in cloudconnexa_host_connector.ha_primary_connectors : k => v.id }
    ha_secondary = { for k, v in cloudconnexa_host_connector.ha_secondary_connectors : k => v.id }
  }
}

output "environment_connectors" {
  description = "Environment-specific connectors"
  value       = { for k, v in cloudconnexa_host_connector.environment_connectors : k => v.id }
}

output "host_info" {
  description = "Host information for reference"
  value = {
    onprem_datacenter    = cloudconnexa_host.onprem_datacenter.id
    branch_office_ny     = cloudconnexa_host.branch_office_ny.id
    branch_office_london = cloudconnexa_host.branch_office_london.id
    cloud_gateway        = cloudconnexa_host.cloud_gateway.id
    backup_site          = cloudconnexa_host.backup_site.id
    edge_location        = cloudconnexa_host.edge_location.id
    partner_network      = cloudconnexa_host.partner_network.id
  }
}

output "connector_summary" {
  description = "Summary of connectors by type and region"
  value = {
    total_connectors = length([
      cloudconnexa_host_connector.primary_datacenter.id,
      cloudconnexa_host_connector.ny_branch_connector.id,
      cloudconnexa_host_connector.london_branch_connector.id,
      cloudconnexa_host_connector.cloud_gateway_us.id,
      cloudconnexa_host_connector.cloud_gateway_eu.id,
      cloudconnexa_host_connector.dr_site_connector.id,
      cloudconnexa_host_connector.edge_us_connector.id,
      cloudconnexa_host_connector.edge_eu_connector.id,
      cloudconnexa_host_connector.edge_apac_connector.id,
      cloudconnexa_host_connector.partner_connector.id
    ]) + length(var.regional_connectors) + length(local.ha_connector_pairs) * 2 + length(var.environment_connectors)

    us_region_connectors = [
      cloudconnexa_host_connector.primary_datacenter.id,
      cloudconnexa_host_connector.ny_branch_connector.id,
      cloudconnexa_host_connector.cloud_gateway_us.id,
      cloudconnexa_host_connector.dr_site_connector.id,
      cloudconnexa_host_connector.edge_us_connector.id,
      cloudconnexa_host_connector.partner_connector.id
    ]

    eu_region_connectors = [
      cloudconnexa_host_connector.london_branch_connector.id,
      cloudconnexa_host_connector.cloud_gateway_eu.id,
      cloudconnexa_host_connector.edge_eu_connector.id
    ]

    apac_region_connectors = [
      cloudconnexa_host_connector.edge_apac_connector.id
    ]

    infrastructure_types = {
      datacenter = [cloudconnexa_host_connector.primary_datacenter.id]
      branch_offices = [
        cloudconnexa_host_connector.ny_branch_connector.id,
        cloudconnexa_host_connector.london_branch_connector.id
      ]
      cloud_gateways = [
        cloudconnexa_host_connector.cloud_gateway_us.id,
        cloudconnexa_host_connector.cloud_gateway_eu.id
      ]
      edge_computing = [
        cloudconnexa_host_connector.edge_us_connector.id,
        cloudconnexa_host_connector.edge_eu_connector.id,
        cloudconnexa_host_connector.edge_apac_connector.id
      ]
      disaster_recovery   = [cloudconnexa_host_connector.dr_site_connector.id]
      partner_integration = [cloudconnexa_host_connector.partner_connector.id]
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `host_id` (String) The id of the network with which this connector is associated.
- `name` (String) The connector display name.
- `vpn_region_id` (String) The ID of the region where the connector will be deployed. Actual list of available regions can be obtained from data_source_vpn_regions.

### Optional

- `description` (String) The description for the UI. Defaults to `Managed by Terraform`.
- `status` (String) The status of the connector. Valid values are `ACTIVE` or `SUSPENDED`. When set to `SUSPENDED`, the connector will be suspended. Note: This is a write-only field - the API does not return connector status.

### Read-Only

- `connection_status` (String) The connection status of the connector.
- `id` (String) The ID of this resource.
- `ip_v4_address` (String) The IPV4 address of the connector.
- `ip_v6_address` (String) The IPV6 address of the connector.
- `profile` (String, Sensitive) OpenVPN profile of the connector.
- `token` (String, Sensitive) Connector token.

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
terraform import cloudconnexa_host_connector.example <id>
```
