# Simple
resource "cloudconnexa_user_group" "this" {
  name           = "test-group"
  vpn_region_ids = ["eu-central-1"]
  connect_auth   = "AUTH"
}

# Advanced
resource "cloudconnexa_user_group" "this" {
  for_each             = var.user_groups
  name                 = each.key
  connect_auth         = try(each.value.connect_auth, "AUTO")
  internet_access      = try(each.value.internet_access, "LOCAL")
  max_device           = try(each.value.max_device, 3)
  all_regions_included = true
}
