resource "cloudconnexa_user_group" "ug01" {
  name                 = "ug01"
  all_regions_included = true
  connect_auth         = "ON_PRIOR_AUTH"
  internet_access      = "SPLIT_TUNNEL_ON"
  max_device           = "3"
}

# minimalistic user (not very useful in real world)
resource "cloudconnexa_user" "user1" {
  username = "test_user1"
  group_id = cloudconnexa_user_group.ug01.id
  role     = "MEMBER"
}

# more real life example
resource "cloudconnexa_user" "user2" {
  username   = "test_use2"
  group_id   = cloudconnexa_user_group.ug01.id
  role       = "MEMBER"
  first_name = "John"
  last_name  = "Doe"
  email      = "John.Doe@example.com" # replace with valid email address (or user creation will fail!)
}
