resource "cloudconnexa_user_group" "ug01" {
  name                 = "ug01"
  all_regions_included = true
  connect_auth         = "ON_PRIOR_AUTH"
  internet_access      = "SPLIT_TUNNEL_ON"
  max_device           = "3"
}

resource "cloudconnexa_user" "user1" {
  username = "test_user"
  group_id = cloudconnexa_user_group.ug01.id
  role     = "MEMBER"
}
