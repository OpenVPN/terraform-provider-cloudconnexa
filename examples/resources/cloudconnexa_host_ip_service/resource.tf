resource "cloudconnexa_host" "this" {
  name            = "test_host"
  description     = "Managed by Terraform"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "test.example.local"
}

resource "cloudconnexa_host_ip_service" "example1" {
  name        = "example1"
  description = "Managed by Terraform"
  host_id     = cloudconnexa_host.this.id

  config {
    service_types = ["ANY"]
  }
}

resource "cloudconnexa_host_ip_service" "example2" {
  name        = "example2"
  description = "Managed by Terraform"
  host_id     = cloudconnexa_host.this.id

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
}

# Hosts for different IP service scenarios
resource "cloudconnexa_host" "database_cluster" {
  name            = "database-cluster-host"
  description     = "Database cluster host for IP services"
  internet_access = "RESTRICTED_INTERNET"
  domain          = "db-cluster.internal.example.com"
}

resource "cloudconnexa_host" "web_farm" {
  name            = "web-farm-host"
  description     = "Web farm host for load-balanced services"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "webfarm.example.com"
}

resource "cloudconnexa_host" "microservices" {
  name            = "microservices-host"
  description     = "Microservices platform host"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "microservices.platform.example.com"
}

resource "cloudconnexa_host" "monitoring_host" {
  name            = "monitoring-infrastructure"
  description     = "Monitoring and observability infrastructure"
  internet_access = "SPLIT_TUNNEL_ON"
  domain          = "monitoring.infra.example.com"
}

resource "cloudconnexa_host" "network_services" {
  name            = "network-services-host"
  description     = "Network infrastructure services"
  internet_access = "RESTRICTED_INTERNET"
  domain          = "netservices.internal.example.com"
}

resource "cloudconnexa_host" "legacy_systems" {
  name            = "legacy-systems-host"
  description     = "Legacy systems and applications"
  internet_access = "RESTRICTED_INTERNET"
  domain          = "legacy.internal.example.com"
}

# 1. Database IP services with specific ports
resource "cloudconnexa_host_ip_service" "postgresql_cluster" {
  name        = "postgresql-cluster-service"
  description = "PostgreSQL database cluster IP service"
  host_id     = cloudconnexa_host.database_cluster.id

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 5432
      to_port   = 5432
    }
  }
}

resource "cloudconnexa_host_ip_service" "mysql_cluster" {
  name        = "mysql-cluster-service"
  description = "MySQL database cluster IP service"
  host_id     = cloudconnexa_host.database_cluster.id

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 3306
      to_port   = 3306
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 3307 # Secondary MySQL instance
      to_port   = 3307
    }
  }
}

resource "cloudconnexa_host_ip_service" "redis_cluster" {
  name        = "redis-cluster-service"
  description = "Redis cluster IP service"
  host_id     = cloudconnexa_host.database_cluster.id

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 6379 # Redis primary
      to_port   = 6379
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 16379 # Redis Cluster bus
      to_port   = 16379
    }
  }
}

resource "cloudconnexa_host_ip_service" "mongodb_cluster" {
  name        = "mongodb-cluster-service"
  description = "MongoDB cluster IP service"
  host_id     = cloudconnexa_host.database_cluster.id

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 27017 # MongoDB default
      to_port   = 27017
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 27018 # MongoDB shard
      to_port   = 27019
    }
  }
}

# 2. Web services with HTTP/HTTPS and custom ports
resource "cloudconnexa_host_ip_service" "web_load_balancer" {
  name        = "web-load-balancer-service"
  description = "Web application load balancer IP service"
  host_id     = cloudconnexa_host.web_farm.id

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
}

resource "cloudconnexa_host_ip_service" "nginx_reverse_proxy" {
  name        = "nginx-reverse-proxy-service"
  description = "Nginx reverse proxy IP service"
  host_id     = cloudconnexa_host.web_farm.id

  config {
    service_types = ["HTTP", "HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 8081 # Admin interface
      to_port   = 8081
    }
  }
}

resource "cloudconnexa_host_ip_service" "apache_web_cluster" {
  name        = "apache-web-cluster-service"
  description = "Apache web server cluster IP service"
  host_id     = cloudconnexa_host.web_farm.id

  config {
    service_types = ["HTTP", "HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 8000 # Development port
      to_port   = 8010 # Range for multiple instances
    }
  }
}

# 3. Microservices IP services
resource "cloudconnexa_host_ip_service" "api_gateway" {
  name        = "api-gateway-service"
  description = "API Gateway IP service"
  host_id     = cloudconnexa_host.microservices.id

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
}

resource "cloudconnexa_host_ip_service" "message_queue" {
  name        = "message-queue-service"
  description = "Message queue system IP service"
  host_id     = cloudconnexa_host.microservices.id

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 5672 # RabbitMQ
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
  }
}

resource "cloudconnexa_host_ip_service" "container_orchestration" {
  name        = "kubernetes-cluster-service"
  description = "Kubernetes cluster IP service"
  host_id     = cloudconnexa_host.microservices.id

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
}

# 4. Monitoring and observability IP services
resource "cloudconnexa_host_ip_service" "prometheus_stack" {
  name        = "prometheus-monitoring-service"
  description = "Prometheus monitoring stack IP service"
  host_id     = cloudconnexa_host.monitoring_host.id

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
}

resource "cloudconnexa_host_ip_service" "grafana_dashboards" {
  name        = "grafana-dashboard-service"
  description = "Grafana dashboard IP service"
  host_id     = cloudconnexa_host.monitoring_host.id

  config {
    service_types = ["HTTPS"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 3000 # Grafana
      to_port   = 3000
    }
  }
}

resource "cloudconnexa_host_ip_service" "elk_logging" {
  name        = "elk-logging-service"
  description = "ELK stack logging IP service"
  host_id     = cloudconnexa_host.monitoring_host.id

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
}

# 5. Network infrastructure IP services
resource "cloudconnexa_host_ip_service" "dns_servers" {
  name        = "dns-infrastructure-service"
  description = "DNS infrastructure IP service"
  host_id     = cloudconnexa_host.network_services.id

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
}

resource "cloudconnexa_host_ip_service" "dhcp_service" {
  name        = "dhcp-infrastructure-service"
  description = "DHCP infrastructure IP service"
  host_id     = cloudconnexa_host.network_services.id

  config {
    service_types = []
    custom_service_types {
      protocol  = "UDP"
      from_port = 67 # DHCP server
      to_port   = 67
    }
    custom_service_types {
      protocol  = "UDP"
      from_port = 68 # DHCP client
      to_port   = 68
    }
  }
}

resource "cloudconnexa_host_ip_service" "ntp_service" {
  name        = "ntp-time-service"
  description = "NTP time synchronization IP service"
  host_id     = cloudconnexa_host.network_services.id

  config {
    service_types = []
    custom_service_types {
      protocol  = "UDP"
      from_port = 123 # NTP
      to_port   = 123
    }
  }
}

resource "cloudconnexa_host_ip_service" "vpn_concentrator" {
  name        = "vpn-concentrator-service"
  description = "VPN concentrator IP service"
  host_id     = cloudconnexa_host.network_services.id

  config {
    service_types = []
    custom_service_types {
      protocol  = "UDP"
      from_port = 1194 # OpenVPN
      to_port   = 1194
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 443 # SSL VPN
      to_port   = 443
    }
    custom_service_types {
      protocol  = "UDP"
      from_port = 500 # IPSec IKE
      to_port   = 500
    }
    custom_service_types {
      protocol  = "UDP"
      from_port = 4500 # IPSec NAT-T
      to_port   = 4500
    }
  }
}

# 6. Security and authentication IP services
resource "cloudconnexa_host_ip_service" "ldap_directory" {
  name        = "ldap-directory-service"
  description = "LDAP directory IP service"
  host_id     = cloudconnexa_host.network_services.id

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
      to_port   = 3268
    }
  }
}

resource "cloudconnexa_host_ip_service" "radius_auth" {
  name        = "radius-authentication-service"
  description = "RADIUS authentication IP service"
  host_id     = cloudconnexa_host.network_services.id

  config {
    service_types = []
    custom_service_types {
      protocol  = "UDP"
      from_port = 1812 # RADIUS auth
      to_port   = 1812
    }
    custom_service_types {
      protocol  = "UDP"
      from_port = 1813 # RADIUS accounting
      to_port   = 1813
    }
  }
}

# 7. Legacy systems with comprehensive access
resource "cloudconnexa_host_ip_service" "legacy_comprehensive" {
  name        = "legacy-systems-any-service"
  description = "Legacy systems with comprehensive IP service access"
  host_id     = cloudconnexa_host.legacy_systems.id

  config {
    service_types = ["ANY"]
  }
}

resource "cloudconnexa_host_ip_service" "legacy_mainframe" {
  name        = "legacy-mainframe-service"
  description = "Legacy mainframe IP service"
  host_id     = cloudconnexa_host.legacy_systems.id

  config {
    service_types = []
    custom_service_types {
      protocol  = "TCP"
      from_port = 23 # Telnet
      to_port   = 23
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 3270 # TN3270
      to_port   = 3270
    }
    custom_service_types {
      protocol  = "TCP"
      from_port = 5250 # TN5250
      to_port   = 5250
    }
  }
}

# Multiple IP services using for_each pattern
variable "application_ip_services" {
  description = "Application IP services to deploy"
  type = map(object({
    description   = string
    host_id       = string
    service_types = list(string)
    custom_ports = list(object({
      protocol  = string
      from_port = number
      to_port   = number
    }))
  }))
  default = {
    "jenkins-ci" = {
      description   = "Jenkins CI/CD IP service"
      host_id       = "microservices"
      service_types = ["HTTP", "HTTPS"]
      custom_ports = [
        {
          protocol  = "TCP"
          from_port = 8080
          to_port   = 8080
        }
      ]
    }
    "gitlab-runner" = {
      description   = "GitLab Runner IP service"
      host_id       = "microservices"
      service_types = ["HTTPS"]
      custom_ports = [
        {
          protocol  = "TCP"
          from_port = 8093
          to_port   = 8093
        }
      ]
    }
    "docker-registry" = {
      description   = "Docker registry IP service"
      host_id       = "microservices"
      service_types = ["HTTPS"]
      custom_ports = [
        {
          protocol  = "TCP"
          from_port = 5000
          to_port   = 5000
        }
      ]
    }
  }
}

resource "cloudconnexa_host_ip_service" "application_ip_services" {
  for_each = var.application_ip_services

  name        = each.key
  description = each.value.description
  host_id = each.value.host_id == "database_cluster" ? cloudconnexa_host.database_cluster.id : (
    each.value.host_id == "web_farm" ? cloudconnexa_host.web_farm.id : (
      each.value.host_id == "microservices" ? cloudconnexa_host.microservices.id : (
        each.value.host_id == "monitoring_host" ? cloudconnexa_host.monitoring_host.id : cloudconnexa_host.network_services.id
      )
    )
  )

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
}

# Game and multimedia IP services
locals {
  game_services = {
    "minecraft-server" = {
      description = "Minecraft game server"
      host_id     = "web_farm"
      ports = [
        { protocol = "TCP", from_port = 25565, to_port = 25565 }
      ]
    }
    "teamspeak-server" = {
      description = "TeamSpeak voice server"
      host_id     = "web_farm"
      ports = [
        { protocol = "UDP", from_port = 9987, to_port = 9987 },
        { protocol = "TCP", from_port = 10011, to_port = 10011 },
        { protocol = "TCP", from_port = 30033, to_port = 30033 }
      ]
    }
    "streaming-server" = {
      description = "Media streaming server"
      host_id     = "web_farm"
      ports = [
        { protocol = "TCP", from_port = 1935, to_port = 1935 }, # RTMP
        { protocol = "UDP", from_port = 1935, to_port = 1935 }  # RTMP UDP
      ]
    }
  }
}

resource "cloudconnexa_host_ip_service" "game_services" {
  for_each = local.game_services

  name        = each.key
  description = each.value.description
  host_id = each.value.host_id == "database_cluster" ? cloudconnexa_host.database_cluster.id : (
    each.value.host_id == "web_farm" ? cloudconnexa_host.web_farm.id : cloudconnexa_host.microservices.id
  )

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
}

# Outputs
output "database_ip_services" {
  description = "Database host IP services"
  value = {
    postgresql = cloudconnexa_host_ip_service.postgresql_cluster.id
    mysql      = cloudconnexa_host_ip_service.mysql_cluster.id
    redis      = cloudconnexa_host_ip_service.redis_cluster.id
    mongodb    = cloudconnexa_host_ip_service.mongodb_cluster.id
  }
}

output "web_ip_services" {
  description = "Web server host IP services"
  value = {
    load_balancer  = cloudconnexa_host_ip_service.web_load_balancer.id
    nginx_proxy    = cloudconnexa_host_ip_service.nginx_reverse_proxy.id
    apache_cluster = cloudconnexa_host_ip_service.apache_web_cluster.id
  }
}

output "microservice_ip_services" {
  description = "Microservices host IP services"
  value = {
    api_gateway   = cloudconnexa_host_ip_service.api_gateway.id
    message_queue = cloudconnexa_host_ip_service.message_queue.id
    kubernetes    = cloudconnexa_host_ip_service.container_orchestration.id
  }
}

output "monitoring_ip_services" {
  description = "Monitoring host IP services"
  value = {
    prometheus = cloudconnexa_host_ip_service.prometheus_stack.id
    grafana    = cloudconnexa_host_ip_service.grafana_dashboards.id
    elk_stack  = cloudconnexa_host_ip_service.elk_logging.id
  }
}

output "network_infrastructure_services" {
  description = "Network infrastructure IP services"
  value = {
    dns_servers      = cloudconnexa_host_ip_service.dns_servers.id
    dhcp_service     = cloudconnexa_host_ip_service.dhcp_service.id
    ntp_service      = cloudconnexa_host_ip_service.ntp_service.id
    vpn_concentrator = cloudconnexa_host_ip_service.vpn_concentrator.id
  }
}

output "security_services" {
  description = "Security and authentication IP services"
  value = {
    ldap_directory = cloudconnexa_host_ip_service.ldap_directory.id
    radius_auth    = cloudconnexa_host_ip_service.radius_auth.id
  }
}

output "legacy_services" {
  description = "Legacy system IP services"
  value = {
    comprehensive = cloudconnexa_host_ip_service.legacy_comprehensive.id
    mainframe     = cloudconnexa_host_ip_service.legacy_mainframe.id
  }
}

output "application_services" {
  description = "Application IP services created with for_each"
  value       = { for k, v in cloudconnexa_host_ip_service.application_ip_services : k => v.id }
}

output "game_multimedia_services" {
  description = "Game and multimedia IP services"
  value       = { for k, v in cloudconnexa_host_ip_service.game_services : k => v.id }
}

output "host_info" {
  description = "Host information for reference"
  value = {
    database_cluster = cloudconnexa_host.database_cluster.id
    web_farm         = cloudconnexa_host.web_farm.id
    microservices    = cloudconnexa_host.microservices.id
    monitoring_host  = cloudconnexa_host.monitoring_host.id
    network_services = cloudconnexa_host.network_services.id
    legacy_systems   = cloudconnexa_host.legacy_systems.id
  }
}

output "service_summary" {
  description = "Summary of IP services by category"
  value = {
    total_services = length([
      cloudconnexa_host_ip_service.postgresql_cluster.id,
      cloudconnexa_host_ip_service.mysql_cluster.id,
      cloudconnexa_host_ip_service.redis_cluster.id,
      cloudconnexa_host_ip_service.mongodb_cluster.id,
      cloudconnexa_host_ip_service.web_load_balancer.id,
      cloudconnexa_host_ip_service.nginx_reverse_proxy.id,
      cloudconnexa_host_ip_service.apache_web_cluster.id,
      cloudconnexa_host_ip_service.api_gateway.id,
      cloudconnexa_host_ip_service.message_queue.id,
      cloudconnexa_host_ip_service.container_orchestration.id,
      cloudconnexa_host_ip_service.prometheus_stack.id,
      cloudconnexa_host_ip_service.grafana_dashboards.id,
      cloudconnexa_host_ip_service.elk_logging.id,
      cloudconnexa_host_ip_service.dns_servers.id,
      cloudconnexa_host_ip_service.dhcp_service.id,
      cloudconnexa_host_ip_service.ntp_service.id,
      cloudconnexa_host_ip_service.vpn_concentrator.id,
      cloudconnexa_host_ip_service.ldap_directory.id,
      cloudconnexa_host_ip_service.radius_auth.id,
      cloudconnexa_host_ip_service.legacy_comprehensive.id,
      cloudconnexa_host_ip_service.legacy_mainframe.id
    ]) + length(var.application_ip_services) + length(local.game_services)

    by_category = {
      database_services        = 4
      web_services             = 3
      microservice_services    = 3
      monitoring_services      = 3
      network_infrastructure   = 4
      security_services        = 2
      legacy_services          = 2
      application_services     = length(var.application_ip_services)
      game_multimedia_services = length(local.game_services)
    }

    by_protocol = {
      tcp_services        = 18
      udp_services        = 8
      http_https_services = 6
      any_services        = 1
    }
  }
}
