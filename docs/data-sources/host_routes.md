---
page_title: "cloudconnexa_host_routes Data Source - terraform-provider-cloudconnexa"
subcategory: ""
description: |-
  Use cloudconnexa_host_routes data source to retrieve routes for a specific host.
---

# cloudconnexa_host_routes (Data Source)

Use `cloudconnexa_host_routes` data source to retrieve routes for a specific host.

## Example Usage

```terraform
data "cloudconnexa_host" "example" {
  name = "example-host"
}

data "cloudconnexa_host_routes" "example" {
  host_id = data.cloudconnexa_host.example.id
}

output "route_count" {
  value = length(data.cloudconnexa_host_routes.example.routes)
}

output "route_subnets" {
  value = [for r in data.cloudconnexa_host_routes.example.routes : r.subnet]
}
```

## Schema

### Required

- `host_id` (String) The ID of the host to get routes for.

### Read-Only

- `id` (String) The ID of this resource.
- `routes` (List of Object) The list of routes for the host. (see [below for nested schema](#nestedatt--routes))

<a id="nestedatt--routes"></a>
### Nested Schema for `routes`

Read-Only:

- `id` (String) The route ID.
- `type` (String) The route type.
- `subnet` (String) The route subnet.
- `domain` (String) The route domain (if applicable).
- `description` (String) The route description.
