data "cloudconnexa_network" "test-net" {
  name = "test-net"
}

resource "cloudconnexa_application" "application_full_access" {
  name              = "example-application-1"
  network_item_type = "NETWORK"
  network_item_id   = data.cloudconnexa_network.test-net.id
  routes {
    domain = "example-application-1.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["ANY"]
  }
}

resource "cloudconnexa_application" "application_custom_access" {
  name              = "example-application-2"
  network_item_type = "NETWORK"
  network_item_id   = data.cloudconnexa_network.test-net.id

  routes {
    domain = "example-application-2.com"
    allow_embedded_ip = false
  }

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

variable "application_custom_access_advanced" {
  description = "xxx"
  type        = any
  default = {
    "example-application-3" = { route = [{ domain = "example-application-3.com", allow_embedded_ip = true }, { domain = "example-application-33.com", allow_embedded_ip = false }] }
    "example-application-4" = { route = [{ domain = "example-application-4.com", allow_embedded_ip = false }] }
  }
}

resource "cloudconnexa_application" "application_custom_access_advanced" {
  for_each          = var.application_custom_access_advanced
  name              = each.key
  description       = try(each.value.description, local.created_by)
  network_item_type = "NETWORK"
  network_item_id   = data.cloudconnexa_network.test-net.id
  config {
    service_types = ["ANY"]
  }

  dynamic "routes" {
    for_each = each.value.route

    content {
      domain            = routes.value.domain
      allow_embedded_ip = routes.value.allow_embedded_ip
    }
  }
}
