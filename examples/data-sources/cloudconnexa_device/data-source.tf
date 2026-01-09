# Get a specific device by ID
data "cloudconnexa_device" "example" {
  device_id = "device-id-here"
}

# Output device information
output "device_info" {
  value = {
    id          = data.cloudconnexa_device.example.id
    name        = data.cloudconnexa_device.example.name
    description = data.cloudconnexa_device.example.description
    platform    = data.cloudconnexa_device.example.platform
    status      = data.cloudconnexa_device.example.status
    user_id     = data.cloudconnexa_device.example.user_id
  }
}
