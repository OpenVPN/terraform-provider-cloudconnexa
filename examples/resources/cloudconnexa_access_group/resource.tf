# Prerequisites: User groups for different access levels
resource "cloudconnexa_user_group" "developers" {
  name                 = "developers"
  description          = "Development team"
  all_regions_included = true
  connect_auth         = "ON_PRIOR_AUTH"
  internet_access      = "SPLIT_TUNNEL_ON"
  max_device           = "5"
}

resource "cloudconnexa_user_group" "admins" {
  name                 = "admins"
  description          = "System administrators"
  all_regions_included = true
  connect_auth         = "EVERY_TIME"
  internet_access      = "SPLIT_TUNNEL_OFF"
  max_device           = "3"
}

resource "cloudconnexa_user_group" "contractors" {
  name                 = "contractors"
  description          = "External contractors"
  all_regions_included = false
  connect_auth         = "EVERY_TIME"
  internet_access      = "RESTRICTED_INTERNET"
  max_device           = "2"
}

# Networks for different environments
resource "cloudconnexa_network" "production" {
  name            = "production-network"
  description     = "Production environment network"
  egress          = true
  internet_access = "SPLIT_TUNNEL_OFF"
}

resource "cloudconnexa_network" "staging" {
  name            = "staging-network"
  description     = "Staging environment network"
  egress          = true
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_network" "development" {
  name            = "development-network"
  description     = "Development environment network"
  egress          = false
  internet_access = "SPLIT_TUNNEL_ON"
}

# Host for database access
resource "cloudconnexa_host" "database_server" {
  name        = "database-server"
  description = "Production database server"
}

# Applications on networks
resource "cloudconnexa_network_application" "prod_api" {
  name        = "production-api"
  description = "Production API application"
  network_id  = cloudconnexa_network.production.id

  routes {
    domain            = "api.production.example.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["HTTPS", "SSH"]
  }
}

resource "cloudconnexa_network_application" "staging_api" {
  name        = "staging-api"
  description = "Staging API application"
  network_id  = cloudconnexa_network.staging.id

  routes {
    domain            = "api.staging.example.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["HTTPS", "SSH"]
  }
}

resource "cloudconnexa_host_application" "database_app" {
  name        = "database-application"
  description = "Database access application"
  host_id     = cloudconnexa_host.database_server.id

  routes {
    domain            = "db.production.example.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["CUSTOM"]
    custom_service_types {
      protocol  = "TCP"
      from_port = 5432
      to_port   = 5432
    }
  }
}

# 1. Full mesh access (emergency/admin use case)
resource "cloudconnexa_access_group" "emergency_full_access" {
  name        = "Emergency Full Access"
  description = "Emergency access group for critical situations - grants full access to admins"

  source {
    type        = "USER_GROUP"
    all_covered = false
    children    = [cloudconnexa_user_group.admins.id]
  }

  destination {
    type        = "NETWORK"
    all_covered = true
  }

  destination {
    type        = "HOST"
    all_covered = true
  }
}

# 2. Production access (admin-only)
resource "cloudconnexa_access_group" "production_access" {
  name        = "Production Access"
  description = "Production environment access for administrators only"

  source {
    type        = "USER_GROUP"
    all_covered = false
    children    = [cloudconnexa_user_group.admins.id]
  }

  destination {
    type        = "NETWORK"
    all_covered = true
    parent      = cloudconnexa_network.production.id
  }
}

# 3. Staging access (admins and developers)
resource "cloudconnexa_access_group" "staging_access" {
  name        = "Staging Access"
  description = "Staging environment access for admins and developers"

  source {
    type        = "USER_GROUP"
    all_covered = false
    children = [
      cloudconnexa_user_group.admins.id,
      cloudconnexa_user_group.developers.id
    ]
  }

  destination {
    type        = "NETWORK"
    all_covered = true
    parent      = cloudconnexa_network.staging.id
  }
}

# 4. Development access (all internal users)
resource "cloudconnexa_access_group" "development_access" {
  name        = "Development Access"
  description = "Development environment access for all internal users"

  source {
    type        = "USER_GROUP"
    all_covered = false
    children = [
      cloudconnexa_user_group.admins.id,
      cloudconnexa_user_group.developers.id
    ]
  }

  destination {
    type        = "NETWORK"
    all_covered = true
    parent      = cloudconnexa_network.development.id
  }
}

# 5. Database access (admin-only, specific application)
resource "cloudconnexa_access_group" "database_access" {
  name        = "Database Access"
  description = "Direct database access for administrators only"

  source {
    type        = "USER_GROUP"
    all_covered = false
    children    = [cloudconnexa_user_group.admins.id]
  }

  destination {
    type        = "HOST"
    all_covered = false
    children    = [cloudconnexa_host_application.database_app.id]
  }
}

# 6. Contractor limited access (staging only, specific applications)
resource "cloudconnexa_access_group" "contractor_access" {
  name        = "Contractor Limited Access"
  description = "Limited access for external contractors to staging environment"

  source {
    type        = "USER_GROUP"
    all_covered = false
    children    = [cloudconnexa_user_group.contractors.id]
  }

  destination {
    type        = "NETWORK"
    all_covered = false
    children    = [cloudconnexa_network_application.staging_api.id]
  }
}

# 7. Cross-environment user group access (developers can access other developers)
resource "cloudconnexa_access_group" "dev_to_dev_communication" {
  name        = "Developer Communication"
  description = "Allow developers to communicate with each other across environments"

  source {
    type        = "USER_GROUP"
    all_covered = false
    children    = [cloudconnexa_user_group.developers.id]
  }

  destination {
    type        = "USER_GROUP"
    all_covered = false
    children    = [cloudconnexa_user_group.developers.id]
  }
}

# Outputs for access group management
output "access_groups" {
  description = "Created access groups and their IDs"
  value = {
    emergency_full_access    = cloudconnexa_access_group.emergency_full_access.id
    production_access        = cloudconnexa_access_group.production_access.id
    staging_access           = cloudconnexa_access_group.staging_access.id
    development_access       = cloudconnexa_access_group.development_access.id
    database_access          = cloudconnexa_access_group.database_access.id
    contractor_access        = cloudconnexa_access_group.contractor_access.id
    dev_to_dev_communication = cloudconnexa_access_group.dev_to_dev_communication.id
  }
}

output "network_info" {
  description = "Network information for reference"
  value = {
    production  = cloudconnexa_network.production.id
    staging     = cloudconnexa_network.staging.id
    development = cloudconnexa_network.development.id
  }
}

output "user_group_info" {
  description = "User group information for reference"
  value = {
    developers  = cloudconnexa_user_group.developers.id
    admins      = cloudconnexa_user_group.admins.id
    contractors = cloudconnexa_user_group.contractors.id
  }
}
