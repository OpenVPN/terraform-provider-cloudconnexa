resource "cloudconnexa_user_group" "ug01" {
  name                 = "ug01"
  all_regions_included = true
  connect_auth         = "ON_PRIOR_AUTH"
  internet_access      = "SPLIT_TUNNEL_ON"
  max_device           = "3"
}

resource "cloudconnexa_user_group" "ug02" {
  name = "ug02"
  vpn_region_ids = [
    "eu-central-1",
    "eu-central-2"
  ]
  connect_auth    = "ON_PRIOR_AUTH"
  internet_access = "SPLIT_TUNNEL_ON"
  max_device      = "3"
}
