data "cloudconnexa_network" "test-net" {
  name = "test-net"
}

resource "cloudconnexa_network_ip_service" "service_full_access" {
  name        = "example-service-1"
  type        = "IP_SOURCE"
  description = "test-description"
  routes = ["10.0.0.2/32"]
  network_id  = data.cloudconnexa_network.test-net.id
  config {
    service_types = ["ANY"]
  }
}

resource "cloudconnexa_network_ip_service" "service_custom_access" {
  name       = "example-service-2"
  type       = "IP_SOURCE"
  network_id = data.cloudconnexa_network.test-net.id
  routes = ["10.0.0.2/32"]
  config {
    service_types = ["HTTP", "HTTPS"]
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

locals {
  created_by = "Managed by terraform"
}

variable "service_custom_access_advanced" {
  description = "xxx"
  type        = any
  default = {
    "example-service-3" = { route = ["10.0.0.3/32", "10.0.1.3/32"] }
    "example-service-4" = { route = ["10.0.0.4/32"] }
  }
}

resource "cloudconnexa_netork_ip_service" "service_custom_access_advanced" {
  for_each   = var.service_custom_access_advanced
  name       = each.key
  type       = "SERVICE_DESTINATION"
  description = try(each.value.description, local.created_by)
  network_id = data.cloudconnexa_network.test-net.id
  routes     = each.value.route
  config {
    service_types = ["ANY"]
  }
}