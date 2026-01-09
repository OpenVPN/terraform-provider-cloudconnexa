# First, get the list of devices to find the device ID
data "cloudconnexa_devices" "all" {
}

# Manage an existing device - update its name and description
resource "cloudconnexa_device" "example" {
  device_id   = data.cloudconnexa_devices.all.devices[0].id
  name        = "my-work-laptop"
  description = "Work laptop managed by Terraform"
}

# Manage a device for a specific user
data "cloudconnexa_devices" "user_devices" {
  user_id = "user-uuid-here"
}

resource "cloudconnexa_device" "user_device" {
  device_id   = data.cloudconnexa_devices.user_devices.devices[0].id
  name        = "user-mobile-device"
  description = "Mobile device for remote access"
}

# Output managed device information
output "managed_device" {
  value = {
    id       = cloudconnexa_device.example.id
    name     = cloudconnexa_device.example.name
    platform = cloudconnexa_device.example.platform
    status   = cloudconnexa_device.example.status
    user_id  = cloudconnexa_device.example.user_id
  }
}
