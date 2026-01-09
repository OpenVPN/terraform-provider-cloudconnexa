---
page_title: "cloudconnexa_devices Data Source - terraform-provider-cloudconnexa"
subcategory: ""
description: |-
  Use cloudconnexa_devices data source to retrieve device information.
---

# cloudconnexa_devices (Data Source)

Use `cloudconnexa_devices` data source to retrieve device information.

## Example Usage

```terraform
# Get all devices
data "cloudconnexa_devices" "all" {
}

# Get devices for a specific user
data "cloudconnexa_devices" "user_devices" {
  user_id = "user-uuid-here"
}

output "device_count" {
  value = length(data.cloudconnexa_devices.all.devices)
}
```

## Schema

### Optional

- `user_id` (String) Filter devices by user ID.

### Read-Only

- `id` (String) The ID of this resource.
- `devices` (List of Object) The list of devices. (see [below for nested schema](#nestedatt--devices))

<a id="nestedatt--devices"></a>
### Nested Schema for `devices`

Read-Only:

- `id` (String) The device ID.
- `name` (String) The device name.
- `description` (String) The device description.
- `platform` (String) The device platform (e.g., Windows, macOS, iOS, Android).
- `status` (String) The device status (ACTIVE, INACTIVE, BLOCKED, PENDING).
- `user_id` (String) The user ID associated with this device.
