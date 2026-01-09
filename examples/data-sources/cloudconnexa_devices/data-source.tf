# Get all devices
data "cloudconnexa_devices" "all" {
}

# Get devices for a specific user
data "cloudconnexa_devices" "user_devices" {
  user_id = "user-id-here"
}

# Output device information
output "all_devices" {
  value = data.cloudconnexa_devices.all.devices
}

output "active_devices" {
  value = [
    for device in data.cloudconnexa_devices.all.devices : device
    if device.status == "ACTIVE"
  ]
}
