resource "cloudconnexa_network" "this" {
  description     = "Test network"
  egress          = true
  name            = "my_test_network"
  internet_access = "SPLIT_TUNNEL_ON"
  default_route {
    description = "Managed by Terraform"
    subnet      = "192.168.144.0/24"
    type        = "IP_V4"
  }
  default_connector {
    description   = "Managed by Terraform"
    name          = "test-connector"
    vpn_region_id = "eu-central-1"
  }
}

resource "cloudconnexa_network_application" "example1" {
  name        = "example-application-1"
  description = "Managed by Terraform"
  network_id  = cloudconnexa_network.this.id
  routes {
    domain            = "example-application-1.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["ANY"]
  }
}

resource "cloudconnexa_network_application" "example2" {
  name        = "example-application-2"
  description = "Managed by Terraform"
  network_id  = cloudconnexa_network.this.id
  routes {
    domain            = "example-application-2.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["HTTP", "HTTPS"]
    custom_service_types {
      protocol = "TCP" # all TCP ports
    }
    custom_service_types {
      protocol  = "UDP"
      from_port = 1194
      to_port   = 1194
    }
    custom_service_types {
      protocol  = "UDP"
      from_port = 5000
      to_port   = 5010
    }
    custom_service_types {
      protocol  = "ICMP"
      from_port = 8
      to_port   = 8
    }
    custom_service_types {
      protocol  = "ICMP"
      from_port = 20
      to_port   = 22
    }
  }
}
