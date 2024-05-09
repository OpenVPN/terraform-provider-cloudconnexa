terraform {
  required_providers {
    cloudconnexa = {
      version = "0.0.11"
      source  = "cloudconnexa.dev/openvpn/cloudconnexa"
    }
  }
}

provider "cloudconnexa" {
  base_url = ""
}

variable "host_name" {
  default = "tf-autotest"
  type    = string
}

resource "cloudconnexa_host" "host" {
  name            = "TEST_HOST_NAME"
  description     = "Terraform test description 2"
  internet_access = "LOCAL"

  connector {
    name          = "test"
    vpn_region_id = "us-west-1"
  }

  provider = cloudconnexa
}

locals {
  connector_profile = [for connector in cloudconnexa_host.host.connector : connector.profile][0]
}


output "host_id" {
  value = cloudconnexa_host.host.id
}

output "connector_id" {
  value = [for connector in cloudconnexa_host.host.connector : connector.id][0]
}
