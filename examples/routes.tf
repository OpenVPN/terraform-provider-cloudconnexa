resource "cloudconnexa_route" "this" {
  for_each = {
    for key, route in var.routes : route.subnet => route
  }
  network_item_id = var.networks["example-network"]
  type            = "IP_V4"
  subnet          = each.value.subnet
  description     = each.value.description
}
