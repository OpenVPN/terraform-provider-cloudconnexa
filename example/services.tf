data "cloudconnexa_network" "test-net" {
  name = "test-net"
}

resource "cloudconnexa_ip_service" "test-service" {
  name              = "test-service"
  type              = "IP_SOURCE"
  description       = "test-description"
  routes            = ["10.0.0.2/32"]
  network_item_type = "NETWORK"
  network_item_id   = data.cloudconnexa_network.test-net.network_id

  config {
    service_types = ["ANY"]
  }
}
