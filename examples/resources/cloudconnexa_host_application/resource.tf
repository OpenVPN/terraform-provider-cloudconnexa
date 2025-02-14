resource "cloudconnexa_host" "this" {
  name            = "test_host"
  description     = "Managed by Terraform"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "test.example.local"
  connector {
    name          = "test-connector"
    vpn_region_id = "eu-central-1"
  }
}

resource "cloudconnexa_host_application" "example1" {
  name        = "example1"
  description = "Managed by Terraform"
  host_id     = cloudconnexa_host.this.id
  routes {
    domain            = "example1.test.example.local"
    allow_embedded_ip = false
  }

  config {
    service_types = ["ANY"]
  }
}

resource "cloudconnexa_host_application" "example2" {
  name        = "example2"
  description = "Managed by Terraform"
  host_id     = cloudconnexa_host.this.id
  routes {
    domain            = "example2.test.example.local"
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
