# Requires Terraform >= 1.10 (ephemeral resources).

variable "app_server_ip" {
  description = "Public IP of the server the connector runs on; OpenVPN 3 Linux must be installed on it."
  type        = string
}

resource "cloudconnexa_host" "app" {
  name = "app-server"
}

resource "cloudconnexa_host_connector" "app" {
  name          = "app-connector"
  host_id       = cloudconnexa_host.app.id
  vpn_region_id = "us-east-1"
}

ephemeral "cloudconnexa_host_connector_profile" "app" {
  connector_id = cloudconnexa_host_connector.app.id
}

# Provisioner and connection blocks are ephemeral contexts, so the profile can be copied to the
# host directly without ever being stored.
resource "null_resource" "connector" {
  connection {
    type = "ssh"
    host = var.app_server_ip
    user = "ubuntu"
  }

  provisioner "file" {
    content     = ephemeral.cloudconnexa_host_connector_profile.app.profile
    destination = "/tmp/cloudconnexa.ovpn"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo openvpn3 config-import --config /tmp/cloudconnexa.ovpn --name cloudconnexa --persistent",
      "rm /tmp/cloudconnexa.ovpn",
      "sudo systemctl enable --now openvpn3-session@cloudconnexa.service",
    ]
  }
}
