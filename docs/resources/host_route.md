---
page_title: "cloudconnexa_host_route Resource - terraform-provider-cloudconnexa"
subcategory: ""
description: |-
  Use cloudconnexa_host_route to create a route on a CloudConnexa host.
---

# cloudconnexa_host_route (Resource)

Use `cloudconnexa_host_route` to create a route on a CloudConnexa host.

## Example Usage

```terraform
resource "cloudconnexa_host" "example" {
  name           = "example-host"
  description    = "Example host"
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_host_route" "example" {
  host_id     = cloudconnexa_host.example.id
  subnet      = "10.0.0.0/24"
  description = "Example route"
}
```

## Schema

### Required

- `host_id` (String) The ID of the host on which to create the route.
- `subnet` (String) The subnet CIDR for the route.

### Optional

- `description` (String) The description for the route. Defaults to `Managed by Terraform`.

### Read-Only

- `id` (String) The ID of this resource.
- `type` (String) The type of route.

## Import

Import is supported using the following syntax:

```shell
terraform import cloudconnexa_host_route.example <route_id>
```
