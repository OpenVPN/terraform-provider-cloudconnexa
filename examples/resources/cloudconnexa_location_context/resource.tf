resource "cloudconnexa_user_group" "ug01" {
  name                 = "ug01"
  all_regions_included = true
  connect_auth         = "ON_PRIOR_AUTH"
  internet_access      = "SPLIT_TUNNEL_ON"
  max_device           = "3"
}

resource "cloudconnexa_location_context" "this" {
  name            = "Location Context Policy"
  description     = "Description for Location Context Policy"
  user_groups_ids = [cloudconnexa_user_group.ug01.id]
  ip_check {
    allowed = true
    ips {
      ip          = "10.10.0.0/16"
      description = "Test subnet"
    }
    ips {
      ip          = "10.20.0.0/16"
      description = "Test subnet 2"
    }
  }
  country_check {
    allowed   = true
    countries = ["US", "GB"]
  }
  default_check {
    allowed = false
  }
}
