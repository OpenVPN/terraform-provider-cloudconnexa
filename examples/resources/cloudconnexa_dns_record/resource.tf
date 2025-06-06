# Basic A record for single service
resource "cloudconnexa_dns_record" "web_server" {
  domain          = "web.example.com"
  description     = "Production web server"
  ip_v4_addresses = ["10.0.1.100"]
}

# Multiple A records for load balancing
resource "cloudconnexa_dns_record" "api_servers" {
  domain      = "api.example.com"
  description = "API server cluster for load balancing"
  ip_v4_addresses = [
    "10.0.2.100",
    "10.0.2.101",
    "10.0.2.102"
  ]
}

# IPv6 AAAA record
resource "cloudconnexa_dns_record" "ipv6_service" {
  domain      = "ipv6.example.com"
  description = "IPv6-enabled service"
  ip_v6_addresses = [
    "2001:db8::1",
    "2001:db8::2"
  ]
}

# Dual-stack service (both IPv4 and IPv6)
resource "cloudconnexa_dns_record" "dual_stack_service" {
  domain      = "dualstack.example.com"
  description = "Dual-stack service supporting both IPv4 and IPv6"
  ip_v4_addresses = [
    "10.0.3.100",
    "10.0.3.101"
  ]
  ip_v6_addresses = [
    "2001:db8:1::1",
    "2001:db8:1::2"
  ]
}

# Database server with specific IP
resource "cloudconnexa_dns_record" "database_primary" {
  domain          = "db-primary.internal.example.com"
  description     = "Primary database server"
  ip_v4_addresses = ["10.0.100.10"]
}

resource "cloudconnexa_dns_record" "database_replica" {
  domain          = "db-replica.internal.example.com"
  description     = "Database replica server"
  ip_v4_addresses = ["10.0.100.11"]
}

# Microservices DNS records
resource "cloudconnexa_dns_record" "user_service" {
  domain      = "users.microservices.example.com"
  description = "User management microservice"
  ip_v4_addresses = [
    "10.0.10.100",
    "10.0.10.101"
  ]
}

resource "cloudconnexa_dns_record" "auth_service" {
  domain      = "auth.microservices.example.com"
  description = "Authentication microservice"
  ip_v4_addresses = [
    "10.0.11.100",
    "10.0.11.101"
  ]
}

resource "cloudconnexa_dns_record" "payment_service" {
  domain      = "payments.microservices.example.com"
  description = "Payment processing microservice"
  ip_v4_addresses = [
    "10.0.12.100",
    "10.0.12.101",
    "10.0.12.102"
  ]
}

# Environment-specific DNS records
resource "cloudconnexa_dns_record" "staging_api" {
  domain          = "api.staging.example.com"
  description     = "Staging environment API server"
  ip_v4_addresses = ["10.1.2.100"]
}

resource "cloudconnexa_dns_record" "dev_api" {
  domain          = "api.dev.example.com"
  description     = "Development environment API server"
  ip_v4_addresses = ["10.2.2.100"]
}

# Monitoring and infrastructure services
resource "cloudconnexa_dns_record" "prometheus" {
  domain          = "prometheus.monitoring.example.com"
  description     = "Prometheus monitoring server"
  ip_v4_addresses = ["10.0.50.10"]
}

resource "cloudconnexa_dns_record" "grafana" {
  domain          = "grafana.monitoring.example.com"
  description     = "Grafana dashboard server"
  ip_v4_addresses = ["10.0.50.11"]
}

resource "cloudconnexa_dns_record" "elasticsearch" {
  domain      = "elasticsearch.logging.example.com"
  description = "Elasticsearch cluster"
  ip_v4_addresses = [
    "10.0.60.10",
    "10.0.60.11",
    "10.0.60.12"
  ]
}

resource "cloudconnexa_dns_record" "kibana" {
  domain          = "kibana.logging.example.com"
  description     = "Kibana log analysis interface"
  ip_v4_addresses = ["10.0.60.20"]
}

# Cache and storage services
resource "cloudconnexa_dns_record" "redis_cluster" {
  domain      = "redis.cache.example.com"
  description = "Redis cache cluster"
  ip_v4_addresses = [
    "10.0.70.10",
    "10.0.70.11",
    "10.0.70.12"
  ]
}

resource "cloudconnexa_dns_record" "memcached" {
  domain      = "memcached.cache.example.com"
  description = "Memcached distributed cache"
  ip_v4_addresses = [
    "10.0.71.10",
    "10.0.71.11"
  ]
}

# Content delivery and static assets
resource "cloudconnexa_dns_record" "cdn_origin" {
  domain      = "origin.cdn.example.com"
  description = "CDN origin server"
  ip_v4_addresses = [
    "10.0.80.10",
    "10.0.80.11"
  ]
}

resource "cloudconnexa_dns_record" "static_assets" {
  domain          = "static.example.com"
  description     = "Static assets server"
  ip_v4_addresses = ["10.0.81.10"]
}

# Multiple DNS records using for_each pattern
variable "service_endpoints" {
  description = "Map of service endpoints to create DNS records for"
  type = map(object({
    description     = string
    ip_v4_addresses = list(string)
    ip_v6_addresses = optional(list(string), [])
  }))
  default = {
    "mail.example.com" = {
      description     = "Mail server"
      ip_v4_addresses = ["10.0.90.10"]
    }
    "ftp.example.com" = {
      description     = "FTP server"
      ip_v4_addresses = ["10.0.91.10"]
    }
    "vpn.example.com" = {
      description     = "VPN gateway"
      ip_v4_addresses = ["10.0.92.10"]
      ip_v6_addresses = ["2001:db8:2::1"]
    }
    "backup.example.com" = {
      description     = "Backup server"
      ip_v4_addresses = ["10.0.93.10", "10.0.93.11"]
    }
  }
}

resource "cloudconnexa_dns_record" "service_endpoints" {
  for_each = var.service_endpoints

  domain          = each.key
  description     = each.value.description
  ip_v4_addresses = each.value.ip_v4_addresses
  ip_v6_addresses = length(each.value.ip_v6_addresses) > 0 ? each.value.ip_v6_addresses : null
}

# Regional DNS records
locals {
  regions = {
    us-east    = "10.1.0.0/16"
    us-west    = "10.2.0.0/16"
    eu-central = "10.3.0.0/16"
    ap-south   = "10.4.0.0/16"
  }
}

resource "cloudconnexa_dns_record" "regional_api" {
  for_each = local.regions

  domain          = "api-${each.key}.example.com"
  description     = "Regional API server for ${each.key}"
  ip_v4_addresses = [cidrhost(each.value, 100)]
}

# Development team DNS records
variable "dev_environments" {
  description = "Development environment DNS records"
  type = map(object({
    team_name = string
    base_ip   = string
  }))
  default = {
    "frontend-team" = {
      team_name = "Frontend Development Team"
      base_ip   = "10.100.1.0"
    }
    "backend-team" = {
      team_name = "Backend Development Team"
      base_ip   = "10.100.2.0"
    }
    "devops-team" = {
      team_name = "DevOps Team"
      base_ip   = "10.100.3.0"
    }
  }
}

resource "cloudconnexa_dns_record" "dev_team_environments" {
  for_each = var.dev_environments

  domain      = "${each.key}.dev.example.com"
  description = "Development environment for ${each.value.team_name}"
  ip_v4_addresses = [
    cidrhost("${each.value.base_ip}/24", 10),
    cidrhost("${each.value.base_ip}/24", 11)
  ]
}

# Outputs
output "production_dns_records" {
  description = "Production DNS records"
  value = {
    web_server       = cloudconnexa_dns_record.web_server.id
    api_servers      = cloudconnexa_dns_record.api_servers.id
    database_primary = cloudconnexa_dns_record.database_primary.id
    database_replica = cloudconnexa_dns_record.database_replica.id
  }
}

output "microservice_dns_records" {
  description = "Microservice DNS records"
  value = {
    user_service    = cloudconnexa_dns_record.user_service.id
    auth_service    = cloudconnexa_dns_record.auth_service.id
    payment_service = cloudconnexa_dns_record.payment_service.id
  }
}

output "monitoring_dns_records" {
  description = "Monitoring and infrastructure DNS records"
  value = {
    prometheus    = cloudconnexa_dns_record.prometheus.id
    grafana       = cloudconnexa_dns_record.grafana.id
    elasticsearch = cloudconnexa_dns_record.elasticsearch.id
    kibana        = cloudconnexa_dns_record.kibana.id
  }
}

output "cache_dns_records" {
  description = "Cache and storage DNS records"
  value = {
    redis_cluster = cloudconnexa_dns_record.redis_cluster.id
    memcached     = cloudconnexa_dns_record.memcached.id
  }
}

output "dynamic_dns_records" {
  description = "Dynamically created DNS records"
  value = {
    service_endpoints     = { for k, v in cloudconnexa_dns_record.service_endpoints : k => v.id }
    regional_api          = { for k, v in cloudconnexa_dns_record.regional_api : k => v.id }
    dev_team_environments = { for k, v in cloudconnexa_dns_record.dev_team_environments : k => v.id }
  }
}

output "ipv6_dns_records" {
  description = "IPv6-enabled DNS records"
  value = {
    ipv6_service       = cloudconnexa_dns_record.ipv6_service.id
    dual_stack_service = cloudconnexa_dns_record.dual_stack_service.id
  }
}
