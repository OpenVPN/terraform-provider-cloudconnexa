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
resource "cloudconnexa_network" "this" {
  description     = "Test network"
  egress          = true
  name            = "my_test_network"
  internet_access = "SPLIT_TUNNEL_ON"
  default_route {
    description = "Managed by Terraform"
    subnet      = "192.168.144.0/24"
    type        = "IP_V4"
  }
  default_connector {
    description   = "Managed by Terraform"
    name          = "test-connector"
    vpn_region_id = "eu-central-1"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `default_connector` (Block List, Min: 1, Max: 1) The default connector of this network. (see [below for nested schema](#nestedblock--default_connector))
- `default_route` (Block List, Min: 1, Max: 1) The default route of this network. (see [below for nested schema](#nestedblock--default_route))
- `name` (String) The display name of the network.

### Optional

- `description` (String) The display description for this resource. Defaults to `Managed by Terraform`.
- `egress` (Boolean) Boolean to control whether this network provides an egress or not.
- `internet_access` (String) The type of internet access provided. Valid values are `SPLIT_TUNNEL_ON`, `SPLIT_TUNNEL_OFF`, or `RESTRICTED_INTERNET`. Defaults to `SPLIT_TUNNEL_ON`.

### Read-Only

- `id` (String) The ID of this resource.
- `system_subnets` (Set of String) The IPV4 and IPV6 subnets automatically assigned to this network.

<a id="nestedblock--default_connector"></a>
### Nested Schema for `default_connector`

Required:

- `name` (String) Name of the connector automatically created and attached to this network.
- `vpn_region_id` (String) The id of the region where the default connector will be deployed.

Optional:

- `description` (String) The default connection description.

Read-Only:

- `id` (String) The ID of this resource.
- `ip_v4_address` (String) The IPV4 address of the default connector.
- `ip_v6_address` (String) The IPV6 address of the default connector.
- `network_id` (String) The parent network id.
- `profile` (String) OpenVPN profile of the connector.


<a id="nestedblock--default_route"></a>
### Nested Schema for `default_route`

Required:

- `subnet` (String) The target value of the default route.

Optional:

- `description` (String) The default route description.
- `type` (String) The type of route. Valid values are `IP_V4`, `IP_V6`, and `DOMAIN`.

Read-Only:

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# Importing resource "cloudconnexa_network" currently is not supported properly (!!).
# You may use as workaround data source "cloudconnexa_network" to reference it in code.
terraform import cloudconnexa_network.example <id>
```
