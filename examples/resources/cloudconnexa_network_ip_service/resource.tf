resource "cloudconnexa_network" "this" {
  description     = "Test network"
  egress          = true
  name            = "my_test_network"
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_route" "this" {
  description     = "Managed by Terraform"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.this.id
  subnet          = "192.168.144.0/24"
}

resource "cloudconnexa_network_ip_service" "example1" {
  name        = "example1"
  description = "Managed by Terraform"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.this.id
  routes      = ["192.168.144.1/32"]
  config {
    service_types = ["ANY"]
  }

  depends_on = [cloudconnexa_route.this]
}

resource "cloudconnexa_network_ip_service" "example2" {
  name        = "example2"
  description = "Managed by Terraform"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.this.id
  routes      = ["192.168.144.2/32"]
  config {
    service_types = ["HTTP", "HTTPS"]
    custom_service_types {
      protocol = "TCP" # all TCP ports
    }
    custom_service_types {
      protocol  = "UDP"
      from_port = 1194
      to_port   = 1194
    }
    custom_service_types {
      protocol  = "UDP"
      from_port = 5000
      to_port   = 5010
    }
    custom_service_types {
      protocol  = "ICMP"
      from_port = 8
      to_port   = 8
    }
    custom_service_types {
      protocol  = "ICMP"
      from_port = 20
      to_port   = 22
    }
  }

  depends_on = [cloudconnexa_route.this]
}

# Networks for different IP service scenarios
resource "cloudconnexa_network" "production_services" {
  description     = "Production services network"
  egress          = true
  name            = "production-services-network"
  internet_access = "RESTRICTED_INTERNET"
}

resource "cloudconnexa_network" "database_network" {
  description     = "Database services network"
  egress          = false
  name            = "database-network"
  internet_access = "RESTRICTED_INTERNET"
}

resource "cloudconnexa_network" "web_services" {
  description     = "Web services network"
  egress          = true
  name            = "web-services-network"
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_network" "microservices_network" {
  description     = "Microservices platform network"
  egress          = true
  name            = "microservices-network"
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_network" "monitoring_network" {
  description     = "Monitoring and observability network"
  egress          = true
  name            = "monitoring-network"
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_network" "infrastructure_network" {
  description     = "Infrastructure services network"
  egress          = true
  name            = "infrastructure-network"
  internet_access = "RESTRICTED_INTERNET"
}

# Routes for different service subnets
resource "cloudconnexa_route" "production_subnet" {
  description     = "Production services subnet"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.production_services.id
  subnet          = "10.100.0.0/24"
}

resource "cloudconnexa_route" "database_subnet" {
  description     = "Database services subnet"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.database_network.id
  subnet          = "10.200.0.0/24"
}

resource "cloudconnexa_route" "web_services_subnet" {
  description     = "Web services subnet"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.web_services.id
  subnet          = "10.110.0.0/24"
}

resource "cloudconnexa_route" "microservices_subnet" {
  description     = "Microservices subnet"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.microservices_network.id
  subnet          = "10.120.0.0/24"
}

resource "cloudconnexa_route" "monitoring_subnet" {
  description     = "Monitoring services subnet"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.monitoring_network.id
  subnet          = "10.130.0.0/24"
}

resource "cloudconnexa_route" "infrastructure_subnet" {
  description     = "Infrastructure services subnet"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.infrastructure_network.id
  subnet          = "10.140.0.0/24"
}

# Additional routes for IPv6 and specialized subnets
resource "cloudconnexa_route" "ipv6_services_subnet" {
  description     = "IPv6 services subnet"
  type            = "IP_V6"
  network_item_id = cloudconnexa_network.production_services.id
  subnet          = "2001:db8:100::/64"
}

resource "cloudconnexa_route" "kubernetes_pods_subnet" {
  description     = "Kubernetes pods subnet"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.microservices_network.id
  subnet          = "172.16.0.0/16"
}

resource "cloudconnexa_route" "docker_bridge_subnet" {
  description     = "Docker bridge network subnet"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.microservices_network.id
  subnet          = "172.17.0.0/16"
}

# 1. Database IP services with specific database ports
resource "cloudconnexa_network_ip_service" "postgresql_cluster" {
  name        = "postgresql-cluster-service"
  description = "PostgreSQL database cluster IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.database_network.id
  routes      = ["10.200.0.10/32", "10.200.0.11/32", "10.200.0.12/32"] # Primary, secondary, read replica

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 5432
      to_port   = 5432
    }
  }

  depends_on = [cloudconnexa_route.database_subnet]
}

resource "cloudconnexa_network_ip_service" "mysql_cluster" {
  name        = "mysql-cluster-service"
  description = "MySQL cluster IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.database_network.id
  routes      = ["10.200.0.20/32", "10.200.0.21/32"] # Master, slave

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 3306
      to_port   = 3306
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 33060 # MySQL X Protocol
      to_port   = 33060
    }
  }

  depends_on = [cloudconnexa_route.database_subnet]
}

resource "cloudconnexa_network_ip_service" "redis_cluster" {
  name        = "redis-cluster-service"
  description = "Redis cluster IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.database_network.id
  routes      = ["10.200.0.30/32", "10.200.0.31/32", "10.200.0.32/32"] # Redis cluster nodes

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 6379 # Redis
      to_port   = 6379
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 16379 # Redis Cluster bus
      to_port   = 16379
    }
  }

  depends_on = [cloudconnexa_route.database_subnet]
}

resource "cloudconnexa_network_ip_service" "mongodb_replica_set" {
  name        = "mongodb-replica-set-service"
  description = "MongoDB replica set IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.database_network.id
  routes      = ["10.200.0.40/32", "10.200.0.41/32", "10.200.0.42/32"] # Primary, secondary, arbiter

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 27017
      to_port   = 27017
    }
  }

  depends_on = [cloudconnexa_route.database_subnet]
}

# 2. Web services with HTTP/HTTPS and custom applications
resource "cloudconnexa_network_ip_service" "nginx_load_balancer" {
  name        = "nginx-load-balancer-service"
  description = "Nginx load balancer IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.web_services.id
  routes      = ["10.110.0.10/32", "10.110.0.11/32"] # Load balancer pair

  config {
    service_types = ["HTTP", "HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 8080 # Alternative HTTP
      to_port   = 8080
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 8443 # Alternative HTTPS
      to_port   = 8443
    }
  }

  depends_on = [cloudconnexa_route.web_services_subnet]
}

resource "cloudconnexa_network_ip_service" "apache_web_farm" {
  name        = "apache-web-farm-service"
  description = "Apache web server farm IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.web_services.id
  routes      = ["10.110.0.20/30"] # Web server farm subnet

  config {
    service_types = ["HTTP", "HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 8080
      to_port   = 8090 # Multiple web server instances
    }
  }

  depends_on = [cloudconnexa_route.web_services_subnet]
}

resource "cloudconnexa_network_ip_service" "nodejs_applications" {
  name        = "nodejs-applications-service"
  description = "Node.js applications IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.web_services.id
  routes      = ["10.110.0.100/29"] # Node.js application cluster

  config {
    service_types = ["HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 3000
      to_port   = 3010 # Multiple Node.js instances
    }
  }

  depends_on = [cloudconnexa_route.web_services_subnet]
}

# 3. Microservices IP services
resource "cloudconnexa_network_ip_service" "api_gateway_cluster" {
  name        = "api-gateway-cluster-service"
  description = "API Gateway cluster IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.microservices_network.id
  routes      = ["10.120.0.10/32", "10.120.0.11/32"] # API Gateway HA pair

  config {
    service_types = ["HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 8080 # API Gateway
      to_port   = 8080
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 9080 # Management API
      to_port   = 9080
    }
  }

  depends_on = [cloudconnexa_route.microservices_subnet]
}

resource "cloudconnexa_network_ip_service" "kubernetes_cluster" {
  name        = "kubernetes-cluster-service"
  description = "Kubernetes cluster IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.microservices_network.id
  routes      = ["10.120.0.20/30"] # Kubernetes master nodes

  config {
    service_types = ["HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 6443 # Kubernetes API server
      to_port   = 6443
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 2379 # etcd
      to_port   = 2380
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 10250 # kubelet
      to_port   = 10250
    }
  }

  depends_on = [cloudconnexa_route.microservices_subnet, cloudconnexa_route.kubernetes_pods_subnet]
}

resource "cloudconnexa_network_ip_service" "service_mesh" {
  name        = "istio-service-mesh-service"
  description = "Istio service mesh IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.microservices_network.id
  routes      = ["172.16.0.0/16"] # All Kubernetes pods

  config {
    service_types = ["HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 15010 # Pilot discovery
      to_port   = 15010
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 15011 # Pilot
      to_port   = 15011
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 8080 # Citadel
      to_port   = 8080
    }
  }

  depends_on = [cloudconnexa_route.kubernetes_pods_subnet]
}

# 4. Monitoring and observability IP services
resource "cloudconnexa_network_ip_service" "prometheus_stack" {
  name        = "prometheus-monitoring-service"
  description = "Prometheus monitoring stack IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.monitoring_network.id
  routes      = ["10.130.0.10/30"] # Prometheus cluster

  config {
    service_types = ["HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 9090 # Prometheus
      to_port   = 9090
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 9093 # Alertmanager
      to_port   = 9093
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 9100 # Node exporter
      to_port   = 9100
    }
  }

  depends_on = [cloudconnexa_route.monitoring_subnet]
}

resource "cloudconnexa_network_ip_service" "grafana_cluster" {
  name        = "grafana-dashboard-service"
  description = "Grafana dashboard cluster IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.monitoring_network.id
  routes      = ["10.130.0.20/32", "10.130.0.21/32"] # Grafana HA

  config {
    service_types = ["HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 3000 # Grafana
      to_port   = 3000
    }
  }

  depends_on = [cloudconnexa_route.monitoring_subnet]
}

resource "cloudconnexa_network_ip_service" "elk_stack" {
  name        = "elk-logging-service"
  description = "ELK stack logging IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.monitoring_network.id
  routes      = ["10.130.0.30/28"] # ELK cluster subnet

  config {
    service_types = ["HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 9200 # Elasticsearch
      to_port   = 9200
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 5601 # Kibana
      to_port   = 5601
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 5044 # Logstash
      to_port   = 5044
    }
    custom_service_types {
      protocol  = "UDP"
      from_port = 5044 # Logstash UDP
      to_port   = 5044
    }
  }

  depends_on = [cloudconnexa_route.monitoring_subnet]
}

# 5. Infrastructure services IP services
resource "cloudconnexa_network_ip_service" "dns_infrastructure" {
  name        = "dns-infrastructure-service"
  description = "DNS infrastructure IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.infrastructure_network.id
  routes      = ["10.140.0.10/32", "10.140.0.11/32"] # Primary and secondary DNS

  config {
    service_types = []
    custom_service_types {
      protocol  = "UDP"
      from_port = 53 # DNS UDP
      to_port   = 53
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 53 # DNS TCP
      to_port   = 53
    }
  }

  depends_on = [cloudconnexa_route.infrastructure_subnet]
}

resource "cloudconnexa_network_ip_service" "ntp_infrastructure" {
  name        = "ntp-time-service"
  description = "NTP time synchronization IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.infrastructure_network.id
  routes      = ["10.140.0.20/32", "10.140.0.21/32"] # NTP servers

  config {
    service_types = []
    custom_service_types {
      protocol  = "UDP"
      from_port = 123 # NTP
      to_port   = 123
    }
  }

  depends_on = [cloudconnexa_route.infrastructure_subnet]
}

resource "cloudconnexa_network_ip_service" "ldap_directory" {
  name        = "ldap-directory-service"
  description = "LDAP directory IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.infrastructure_network.id
  routes      = ["10.140.0.30/32", "10.140.0.31/32"] # LDAP primary and replica

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 389 # LDAP
      to_port   = 389
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 636 # LDAPS
      to_port   = 636
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 3268 # Global Catalog
      to_port   = 3269
    }
  }

  depends_on = [cloudconnexa_route.infrastructure_subnet]
}

# 6. IPv6 services
resource "cloudconnexa_network_ip_service" "ipv6_web_services" {
  name        = "ipv6-web-services"
  description = "IPv6 web services IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.production_services.id
  routes      = ["2001:db8:100::10/128", "2001:db8:100::11/128"] # IPv6 web servers

  config {
    service_types = ["HTTP", "HTTPS"]
  }

  depends_on = [cloudconnexa_route.ipv6_services_subnet]
}

# 7. Message queue and streaming services
resource "cloudconnexa_network_ip_service" "message_queue_cluster" {
  name        = "message-queue-cluster-service"
  description = "Message queue cluster IP service"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.microservices_network.id
  routes      = ["10.120.0.50/29"] # Message queue cluster

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 5672 # RabbitMQ AMQP
      to_port   = 5672
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 15672 # RabbitMQ Management
      to_port   = 15672
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 9092 # Kafka
      to_port   = 9092
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 2181 # Zookeeper
      to_port   = 2181
    }
  }

  depends_on = [cloudconnexa_route.microservices_subnet]
}

# 8. Comprehensive legacy services
resource "cloudconnexa_network_ip_service" "legacy_comprehensive" {
  name        = "legacy-systems-any-service"
  description = "Legacy systems with comprehensive IP service access"
  type        = "IP_SOURCE"
  network_id  = cloudconnexa_network.production_services.id
  routes      = ["10.100.0.200/29"] # Legacy systems subnet

  config {
    service_types = ["ANY"]
  }

  depends_on = [cloudconnexa_route.production_subnet]
}

# Multiple IP services using for_each pattern
variable "application_ip_services" {
  description = "Application IP services to deploy"
  type = map(object({
    description   = string
    network_name  = string
    routes        = list(string)
    service_types = list(string)
    custom_ports = list(object({
      protocol  = string
      from_port = number
      to_port   = number
    }))
  }))
  default = {
    "jenkins-ci-cluster" = {
      description   = "Jenkins CI/CD cluster IP service"
      network_name  = "microservices"
      routes        = ["10.120.0.60/30"]
      service_types = ["HTTP", "HTTPS"]
      custom_ports = [
        {
          protocol  = "TCP"
          from_port = 8080
          to_port   = 8080
        }
      ]
    }
    "artifactory-cluster" = {
      description   = "Artifactory repository cluster IP service"
      network_name  = "microservices"
      routes        = ["10.120.0.70/32", "10.120.0.71/32"]
      service_types = ["HTTPS"]
      custom_ports = [
        {
          protocol  = "TCP"
          from_port = 8081
          to_port   = 8081
        }
      ]
    }
    "sonarqube-service" = {
      description   = "SonarQube code quality IP service"
      network_name  = "microservices"
      routes        = ["10.120.0.80/32"]
      service_types = ["HTTPS"]
      custom_ports = [
        {
          protocol  = "TCP"
          from_port = 9000
          to_port   = 9000
        }
      ]
    }
  }
}

resource "cloudconnexa_network_ip_service" "application_ip_services" {
  for_each = var.application_ip_services

  name        = each.key
  description = each.value.description
  type        = "IP_SOURCE"
  network_id = each.value.network_name == "production" ? cloudconnexa_network.production_services.id : (
    each.value.network_name == "database" ? cloudconnexa_network.database_network.id : (
      each.value.network_name == "web_services" ? cloudconnexa_network.web_services.id : (
        each.value.network_name == "microservices" ? cloudconnexa_network.microservices_network.id : (
          each.value.network_name == "monitoring" ? cloudconnexa_network.monitoring_network.id : cloudconnexa_network.infrastructure_network.id
        )
      )
    )
  )
  routes = each.value.routes

  config {
    service_types = each.value.service_types

    dynamic "custom_service_types" {
      for_each = each.value.custom_ports
      content {
        protocol  = custom_service_types.value.protocol
        from_port = custom_service_types.value.from_port
        to_port   = custom_service_types.value.to_port
      }
    }
  }

  depends_on = [
    cloudconnexa_route.production_subnet,
    cloudconnexa_route.database_subnet,
    cloudconnexa_route.web_services_subnet,
    cloudconnexa_route.microservices_subnet,
    cloudconnexa_route.monitoring_subnet,
    cloudconnexa_route.infrastructure_subnet
  ]
}

# Game and media services
locals {
  game_media_services = {
    "minecraft-servers" = {
      description  = "Minecraft game servers"
      network_name = "web_services"
      routes       = ["10.110.0.200/29"]
      ports = [
        { protocol = "TCP", from_port = 25565, to_port = 25565 },
        { protocol = "UDP", from_port = 25565, to_port = 25565 }
      ]
    }
    "media-streaming" = {
      description  = "Media streaming servers"
      network_name = "web_services"
      routes       = ["10.110.0.210/30"]
      ports = [
        { protocol = "TCP", from_port = 1935, to_port = 1935 }, # RTMP
        { protocol = "UDP", from_port = 1935, to_port = 1935 }, # RTMP UDP
        { protocol = "TCP", from_port = 8080, to_port = 8080 }  # HLS
      ]
    }
    "voice-chat-servers" = {
      description  = "Voice chat servers"
      network_name = "web_services"
      routes       = ["10.110.0.220/30"]
      ports = [
        { protocol = "UDP", from_port = 9987, to_port = 9987 },   # TeamSpeak voice
        { protocol = "TCP", from_port = 10011, to_port = 10011 }, # ServerQuery
        { protocol = "TCP", from_port = 30033, to_port = 30033 }  # FileTransfer
      ]
    }
  }
}

resource "cloudconnexa_network_ip_service" "game_media_services" {
  for_each = local.game_media_services

  name        = each.key
  description = each.value.description
  type        = "IP_SOURCE"
  network_id = each.value.network_name == "production" ? cloudconnexa_network.production_services.id : (
    each.value.network_name == "web_services" ? cloudconnexa_network.web_services.id : cloudconnexa_network.microservices_network.id
  )
  routes = each.value.routes

  config {
    service_types = []

    dynamic "custom_service_types" {
      for_each = each.value.ports
      content {
        protocol  = custom_service_types.value.protocol
        from_port = custom_service_types.value.from_port
        to_port   = custom_service_types.value.to_port
      }
    }
  }

  depends_on = [
    cloudconnexa_route.production_subnet,
    cloudconnexa_route.web_services_subnet,
    cloudconnexa_route.microservices_subnet
  ]
}

# Outputs
output "database_network_services" {
  description = "Database network IP services"
  value = {
    postgresql = cloudconnexa_network_ip_service.postgresql_cluster.id
    mysql      = cloudconnexa_network_ip_service.mysql_cluster.id
    redis      = cloudconnexa_network_ip_service.redis_cluster.id
    mongodb    = cloudconnexa_network_ip_service.mongodb_replica_set.id
  }
}

output "web_network_services" {
  description = "Web services network IP services"
  value = {
    nginx_lb    = cloudconnexa_network_ip_service.nginx_load_balancer.id
    apache_farm = cloudconnexa_network_ip_service.apache_web_farm.id
    nodejs_apps = cloudconnexa_network_ip_service.nodejs_applications.id
    ipv6_web    = cloudconnexa_network_ip_service.ipv6_web_services.id
  }
}

output "microservices_network_services" {
  description = "Microservices network IP services"
  value = {
    api_gateway   = cloudconnexa_network_ip_service.api_gateway_cluster.id
    kubernetes    = cloudconnexa_network_ip_service.kubernetes_cluster.id
    service_mesh  = cloudconnexa_network_ip_service.service_mesh.id
    message_queue = cloudconnexa_network_ip_service.message_queue_cluster.id
  }
}

output "monitoring_network_services" {
  description = "Monitoring network IP services"
  value = {
    prometheus = cloudconnexa_network_ip_service.prometheus_stack.id
    grafana    = cloudconnexa_network_ip_service.grafana_cluster.id
    elk_stack  = cloudconnexa_network_ip_service.elk_stack.id
  }
}

output "infrastructure_network_services" {
  description = "Infrastructure network IP services"
  value = {
    dns_servers    = cloudconnexa_network_ip_service.dns_infrastructure.id
    ntp_servers    = cloudconnexa_network_ip_service.ntp_infrastructure.id
    ldap_directory = cloudconnexa_network_ip_service.ldap_directory.id
  }
}

output "legacy_services" {
  description = "Legacy network IP services"
  value = {
    comprehensive = cloudconnexa_network_ip_service.legacy_comprehensive.id
  }
}

output "application_services" {
  description = "Application IP services created with for_each"
  value       = { for k, v in cloudconnexa_network_ip_service.application_ip_services : k => v.id }
}

output "game_media_services" {
  description = "Game and media IP services"
  value       = { for k, v in cloudconnexa_network_ip_service.game_media_services : k => v.id }
}

output "network_info" {
  description = "Network information for reference"
  value = {
    production_services    = cloudconnexa_network.production_services.id
    database_network       = cloudconnexa_network.database_network.id
    web_services           = cloudconnexa_network.web_services.id
    microservices_network  = cloudconnexa_network.microservices_network.id
    monitoring_network     = cloudconnexa_network.monitoring_network.id
    infrastructure_network = cloudconnexa_network.infrastructure_network.id
  }
}

output "routes_info" {
  description = "Route information for reference"
  value = {
    production_subnet      = cloudconnexa_route.production_subnet.id
    database_subnet        = cloudconnexa_route.database_subnet.id
    web_services_subnet    = cloudconnexa_route.web_services_subnet.id
    microservices_subnet   = cloudconnexa_route.microservices_subnet.id
    monitoring_subnet      = cloudconnexa_route.monitoring_subnet.id
    infrastructure_subnet  = cloudconnexa_route.infrastructure_subnet.id
    ipv6_services_subnet   = cloudconnexa_route.ipv6_services_subnet.id
    kubernetes_pods_subnet = cloudconnexa_route.kubernetes_pods_subnet.id
    docker_bridge_subnet   = cloudconnexa_route.docker_bridge_subnet.id
  }
}

output "service_summary" {
  description = "Summary of network IP services"
  value = {
    total_services = length([
      cloudconnexa_network_ip_service.postgresql_cluster.id,
      cloudconnexa_network_ip_service.mysql_cluster.id,
      cloudconnexa_network_ip_service.redis_cluster.id,
      cloudconnexa_network_ip_service.mongodb_replica_set.id,
      cloudconnexa_network_ip_service.nginx_load_balancer.id,
      cloudconnexa_network_ip_service.apache_web_farm.id,
      cloudconnexa_network_ip_service.nodejs_applications.id,
      cloudconnexa_network_ip_service.api_gateway_cluster.id,
      cloudconnexa_network_ip_service.kubernetes_cluster.id,
      cloudconnexa_network_ip_service.service_mesh.id,
      cloudconnexa_network_ip_service.prometheus_stack.id,
      cloudconnexa_network_ip_service.grafana_cluster.id,
      cloudconnexa_network_ip_service.elk_stack.id,
      cloudconnexa_network_ip_service.dns_infrastructure.id,
      cloudconnexa_network_ip_service.ntp_infrastructure.id,
      cloudconnexa_network_ip_service.ldap_directory.id,
      cloudconnexa_network_ip_service.ipv6_web_services.id,
      cloudconnexa_network_ip_service.message_queue_cluster.id,
      cloudconnexa_network_ip_service.legacy_comprehensive.id
    ]) + length(var.application_ip_services) + length(local.game_media_services)

    by_network = {
      production_services    = 2
      database_network       = 4
      web_services           = 3 + length(local.game_media_services)
      microservices_network  = 4 + length([for k, v in var.application_ip_services : k if v.network_name == "microservices"])
      monitoring_network     = 3
      infrastructure_network = 3
    }

    by_protocol = {
      tcp_services        = 15
      udp_services        = 6
      http_https_services = 8
      any_services        = 1
      mixed_protocols     = 3
    }

    by_service_type = {
      database_services       = 4
      web_services            = 4
      microservices           = 4
      monitoring_services     = 3
      infrastructure_services = 3
      legacy_services         = 1
      application_services    = length(var.application_ip_services)
      game_media_services     = length(local.game_media_services)
    }

    ip_address_ranges = {
      production_range      = "10.100.0.0/24"
      database_range        = "10.200.0.0/24"
      web_services_range    = "10.110.0.0/24"
      microservices_range   = "10.120.0.0/24"
      monitoring_range      = "10.130.0.0/24"
      infrastructure_range  = "10.140.0.0/24"
      ipv6_range            = "2001:db8:100::/64"
      kubernetes_pods_range = "172.16.0.0/16"
      docker_bridge_range   = "172.17.0.0/16"
    }
  }
}
