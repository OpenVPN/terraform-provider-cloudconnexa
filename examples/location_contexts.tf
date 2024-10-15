resource "cloudconnexa_location_context" "this" {
  name = "Location Context Policy"
  description = "Description for Location Context Policy"
  user_groups_ids = []
  ip_policy {
    allowed = true
    ips {
      ip = "10.10.0.0/16"
      description = "Test subnet"
    }
    ips {
      ip = "10.20.0.0/16"
      description = "Test subnet 2"
    }
  }
  country_policy {
    allowed = true
    countries = ["US", "GB"]
  }
  default_policy {
    allowed = false
  }
}

#Advanced
resource "cloudconnexa_location_context" "this" {
  name = "Location Context Policy"
  description = "Description for Location Context Policy"
  user_groups_ids = values(cloudconnexa_user_group.this)[*].id
  ip_policy {
    allowed = true
    ips {
      ip = "10.10.0.0/16"
      description = "Test subnet"
    }
    ips {
      ip = "10.20.0.0/16"
      description = "Test subnet 2"
    }
  }
  country_policy {
    allowed = true
    countries = ["US", "GB"]
  }
  default_policy {
    allowed = false
  }
}