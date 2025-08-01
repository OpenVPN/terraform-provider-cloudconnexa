---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudconnexa_network Resource - terraform-provider-cloudconnexa"
subcategory: ""
description: |-
  Use cloudconnexa_network to create an CloudConnexa Network.
---

# cloudconnexa_network (Resource)

Use `cloudconnexa_network` to create an CloudConnexa Network.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The display name of the network.

### Optional

- `description` (String) The display description for this resource. Defaults to `Managed by Terraform`.
- `egress` (Boolean) Boolean to control whether this network provides an egress or not.
- `internet_access` (String) The type of internet access provided. Valid values are `SPLIT_TUNNEL_ON`, `SPLIT_TUNNEL_OFF`, or `RESTRICTED_INTERNET`. Defaults to `SPLIT_TUNNEL_ON`.
- `tunneling_protocol` (String) The tunneling protocol used for this network.

### Read-Only

- `id` (String) The ID of this resource.
- `system_subnets` (Set of String) The IPV4 and IPV6 subnets automatically assigned to this network.

## Import

Import is supported using the following syntax:

```shell
terraform import cloudconnexa_network.example <id>
```
