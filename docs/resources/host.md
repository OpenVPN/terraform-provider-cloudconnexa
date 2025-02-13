---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudconnexa_host Resource - terraform-provider-cloudconnexa"
subcategory: ""
description: |-
  Use cloudconnexa_host to create an CloudConnexa host.
---

# cloudconnexa_host (Resource)

Use `cloudconnexa_host` to create an CloudConnexa host.

## Example Usage

```terraform
resource "cloudconnexa_host" "this" {
  name            = "test_host"
  description     = "Managed by Terraform"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "test.example.local"
  connector {
    name          = "test-connector"
    vpn_region_id = "eu-central-1"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `connector` (Block Set, Min: 1) The set of connectors to be associated with this host. Can be defined more than once. (see [below for nested schema](#nestedblock--connector))
- `name` (String) The display name of the host.

### Optional

- `description` (String) The description for the UI. Defaults to `Managed by Terraform`.
- `domain` (String) The domain of the host.
- `internet_access` (String) The type of internet access provided. Valid values are `SPLIT_TUNNEL_ON`, `SPLIT_TUNNEL_OFF`, or `RESTRICTED_INTERNET`. Defaults to `SPLIT_TUNNEL_ON`.

### Read-Only

- `id` (String) The ID of this resource.
- `system_subnets` (Set of String) The IPV4 and IPV6 subnets automatically assigned to this host.

<a id="nestedblock--connector"></a>
### Nested Schema for `connector`

Required:

- `name` (String) Name of the connector associated with this host.
- `vpn_region_id` (String) The id of the region where the connector will be deployed.

Optional:

- `description` (String) The description for the UI. Defaults to `Managed by Terraform`.

Read-Only:

- `host_id` (String) The host id.
- `id` (String)
- `ip_v4_address` (String) The IPV4 address of the connector.
- `ip_v6_address` (String) The IPV6 address of the connector.
- `profile` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import cloudconnexa_host.example <id>
```
