---
page_title: "cloudconnexa_device Data Source - terraform-provider-cloudconnexa"
subcategory: ""
description: |-
  Use cloudconnexa_device data source to retrieve a specific device by ID.
---

# cloudconnexa_device (Data Source)

Use `cloudconnexa_device` data source to retrieve a specific device by ID.

## Example Usage

```terraform
data "cloudconnexa_device" "example" {
  device_id = "device-uuid-here"
}

output "device_name" {
  value = data.cloudconnexa_device.example.name
}

output "device_status" {
  value = data.cloudconnexa_device.example.status
}
```

## Schema

### Required

- `device_id` (String) The device ID.

### Read-Only

- `id` (String) The device ID.
- `name` (String) The device name.
- `description` (String) The device description.
- `platform` (String) The device platform.
- `status` (String) The device status.
- `user_id` (String) The user ID associated with this device.
