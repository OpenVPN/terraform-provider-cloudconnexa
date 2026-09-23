# Requires Terraform >= 1.10 (ephemeral resources).

variable "app_server_ip" {
  description = "Public IP of the server the connector is installed on."
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

ephemeral "cloudconnexa_host_connector_token" "app" {
  connector_id = cloudconnexa_host_connector.app.id
}

# Provisioner and connection blocks are ephemeral contexts, so the token can be used to
# bootstrap the host directly without ever being stored.
resource "null_resource" "bootstrap" {
  connection {
    type = "ssh"
    host = var.app_server_ip
    user = "ubuntu"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo openvpn-connector-setup --token '${ephemeral.cloudconnexa_host_connector_token.app.token}'",
    ]
  }
}
