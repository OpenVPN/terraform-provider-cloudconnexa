terraform {
  required_providers {
    cloudconnexa = {
      source  = "OpenVPN/cloudconnexa"
      version = "~> 1"
    }
  }
}

provider "cloudconnexa" {
  client_id     = var.cloudconnexa_client_id
  client_secret = var.cloudconnexa_client_secret
}

# Create a User Group
resource "cloudconnexa_user_group" "this" {
  # ...
}
