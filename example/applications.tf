data "cloudconnexa_network" "test-net" {
  name = "test-net"
}

resource "cloudconnexa_application" "application_full_access" {
  name              = "application-1"
  network_item_type = "NETWORK"
  network_item_id   = data.cloudconnexa_network.test-net.network_id
  routes {
    domain = "application.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["ANY"]
  }
}

resource "cloudconnexa_application" "application_custom_access" {
  name              = "application-1"
  network_item_type = "NETWORK"
  network_item_id   = data.cloudconnexa_network.test-net.network_id

  routes {
    domain = "application.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["HTTP", "HTTPS", "CUSTOM"]
    custom_service_types {
      protocol = "TCP" //all tcp ports
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
