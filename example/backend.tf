terraform {
  backend "local" {}
  required_providers {
    cloudconnexa = {
      source  = "OpenVPN/cloudconnexa"
      version = "0.0.12"
    }
  }
}
