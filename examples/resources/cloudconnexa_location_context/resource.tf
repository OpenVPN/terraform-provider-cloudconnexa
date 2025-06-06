# User groups for different security levels
resource "cloudconnexa_user_group" "high_security_users" {
  name                 = "high-security-users"
  description          = "Users requiring high security location checks"
  all_regions_included = false
  connect_auth         = "EVERY_TIME"
  internet_access      = "RESTRICTED_INTERNET"
  max_device           = "2"
}

resource "cloudconnexa_user_group" "office_workers" {
  name                 = "office-workers"
  description          = "Regular office workers"
  all_regions_included = true
  connect_auth         = "ON_PRIOR_AUTH"
  internet_access      = "SPLIT_TUNNEL_ON"
  max_device           = "3"
}

resource "cloudconnexa_user_group" "remote_contractors" {
  name                 = "remote-contractors"
  description          = "Remote contractors with location restrictions"
  all_regions_included = false
  connect_auth         = "EVERY_TIME"
  internet_access      = "RESTRICTED_INTERNET"
  max_device           = "2"
}

resource "cloudconnexa_user_group" "executives" {
  name                 = "executives"
  description          = "Executive team with flexible access"
  all_regions_included = true
  connect_auth         = "ON_PRIOR_AUTH"
  internet_access      = "SPLIT_TUNNEL_OFF"
  max_device           = "5"
}

resource "cloudconnexa_user_group" "developers" {
  name                 = "developers"
  description          = "Development team"
  all_regions_included = true
  connect_auth         = "ON_PRIOR_AUTH"
  internet_access      = "SPLIT_TUNNEL_ON"
  max_device           = "4"
}

# 1. Strict office-only access policy
resource "cloudconnexa_location_context" "office_only_policy" {
  name            = "Office Only Access Policy"
  description     = "Strict policy allowing access only from office IP ranges"
  user_groups_ids = [cloudconnexa_user_group.high_security_users.id]

  ip_check {
    allowed = true
    ips {
      ip          = "203.0.113.0/24"
      description = "Main office network"
    }
    ips {
      ip          = "198.51.100.0/24"
      description = "Secondary office network"
    }
    ips {
      ip          = "192.0.2.0/24"
      description = "Branch office network"
    }
  }

  country_check {
    allowed   = true
    countries = ["US"]
  }

  default_check {
    allowed = false
  }
}

# 2. Multi-country office access
resource "cloudconnexa_location_context" "global_office_policy" {
  name        = "Global Office Access Policy"
  description = "Access policy for global offices with country and IP restrictions"
  user_groups_ids = [
    cloudconnexa_user_group.office_workers.id,
    cloudconnexa_user_group.developers.id
  ]

  ip_check {
    allowed = true
    ips {
      ip          = "203.0.113.0/24"
      description = "US Headquarters"
    }
    ips {
      ip          = "198.51.100.0/24"
      description = "UK Office"
    }
    ips {
      ip          = "192.0.2.0/24"
      description = "Germany Office"
    }
    ips {
      ip          = "172.16.0.0/24"
      description = "Singapore Office"
    }
    ips {
      ip          = "10.0.0.0/24"
      description = "Canada Office"
    }
  }

  country_check {
    allowed   = true
    countries = ["US", "GB", "DE", "SG", "CA"]
  }

  default_check {
    allowed = false
  }
}

# 3. Executive flexible access policy
resource "cloudconnexa_location_context" "executive_flexible_policy" {
  name            = "Executive Flexible Access Policy"
  description     = "Flexible access policy for executives with broader location allowances"
  user_groups_ids = [cloudconnexa_user_group.executives.id]

  ip_check {
    allowed = true
    ips {
      ip          = "203.0.113.0/24"
      description = "Corporate headquarters"
    }
    ips {
      ip          = "198.51.100.0/24"
      description = "Regional offices"
    }
    ips {
      ip          = "0.0.0.0/0"
      description = "Allow from any IP for executives"
    }
  }

  country_check {
    allowed = true
    countries = [
      "US", "GB", "DE", "FR", "IT", "ES", "NL", "SE", "NO", "DK",
      "CA", "AU", "NZ", "JP", "SG", "HK", "CH", "AT", "BE", "LU"
    ]
  }

  default_check {
    allowed = false
  }
}

# 4. Contractor restricted access
resource "cloudconnexa_location_context" "contractor_restricted_policy" {
  name            = "Contractor Restricted Access Policy"
  description     = "Restricted access policy for external contractors"
  user_groups_ids = [cloudconnexa_user_group.remote_contractors.id]

  ip_check {
    allowed = false
    ips {
      ip          = "203.0.113.0/24"
      description = "Office network - allowed for contractors"
    }
  }

  country_check {
    allowed   = true
    countries = ["US", "CA", "GB"]
  }

  default_check {
    allowed = false
  }
}

# 5. Development team remote work policy
resource "cloudconnexa_location_context" "dev_remote_policy" {
  name            = "Development Team Remote Work Policy"
  description     = "Remote work policy for development team with specific country restrictions"
  user_groups_ids = [cloudconnexa_user_group.developers.id]

  ip_check {
    allowed = true
    ips {
      ip          = "203.0.113.0/24"
      description = "Main office"
    }
    ips {
      ip          = "10.0.0.0/8"
      description = "Private networks allowed"
    }
    ips {
      ip          = "172.16.0.0/12"
      description = "Private networks allowed"
    }
    ips {
      ip          = "192.168.0.0/16"
      description = "Private networks allowed"
    }
  }

  country_check {
    allowed   = true
    countries = ["US", "CA", "GB", "DE", "NL", "AU", "IN"]
  }

  default_check {
    allowed = false
  }
}

# 6. Emergency access policy (very permissive)
resource "cloudconnexa_location_context" "emergency_access_policy" {
  name        = "Emergency Access Policy"
  description = "Emergency access policy with minimal restrictions"
  user_groups_ids = [
    cloudconnexa_user_group.executives.id,
    cloudconnexa_user_group.high_security_users.id
  ]

  ip_check {
    allowed = true
    ips {
      ip          = "0.0.0.0/0"
      description = "Allow from any IP during emergency"
    }
  }

  country_check {
    allowed = true
    countries = [
      "US", "GB", "DE", "FR", "IT", "ES", "NL", "SE", "NO", "DK",
      "CA", "AU", "NZ", "JP", "SG", "HK", "CH", "AT", "BE", "LU",
      "FI", "IE", "PT", "GR", "CZ", "PL", "HU", "SK", "SI", "EE",
      "LV", "LT", "CY", "MT", "BG", "RO", "HR"
    ]
  }

  default_check {
    allowed = false
  }
}

# 7. Geo-blocking policy (block specific countries)
resource "cloudconnexa_location_context" "geo_blocking_policy" {
  name        = "Geo-blocking Security Policy"
  description = "Security policy blocking access from high-risk countries"
  user_groups_ids = [
    cloudconnexa_user_group.office_workers.id,
    cloudconnexa_user_group.developers.id,
    cloudconnexa_user_group.remote_contractors.id
  ]

  ip_check {
    allowed = true
    ips {
      ip          = "203.0.113.0/24"
      description = "Trusted office network"
    }
    ips {
      ip          = "198.51.100.0/24"
      description = "Secondary office network"
    }
  }

  country_check {
    allowed = true
    countries = [
      "US", "CA", "GB", "DE", "FR", "IT", "ES", "NL", "SE", "NO",
      "DK", "FI", "IE", "PT", "GR", "CH", "AT", "BE", "LU", "AU",
      "NZ", "JP", "SG", "HK", "KR", "TW"
    ]
  }

  default_check {
    allowed = false
  }
}

# Multiple location contexts using for_each pattern
variable "regional_policies" {
  description = "Regional location context policies"
  type = map(object({
    description    = string
    user_group_ids = list(string)
    allowed_ips = list(object({
      ip          = string
      description = string
    }))
    allowed_countries = list(string)
    default_allowed   = bool
  }))
  default = {
    "asia-pacific" = {
      description    = "Asia Pacific regional access policy"
      user_group_ids = ["office_workers", "developers"]
      allowed_ips = [
        {
          ip          = "172.16.1.0/24"
          description = "Singapore office"
        },
        {
          ip          = "172.16.2.0/24"
          description = "Hong Kong office"
        }
      ]
      allowed_countries = ["SG", "HK", "JP", "AU", "NZ", "KR"]
      default_allowed   = false
    }
    "europe" = {
      description    = "European regional access policy"
      user_group_ids = ["office_workers", "developers"]
      allowed_ips = [
        {
          ip          = "172.16.10.0/24"
          description = "London office"
        },
        {
          ip          = "172.16.11.0/24"
          description = "Frankfurt office"
        }
      ]
      allowed_countries = ["GB", "DE", "FR", "IT", "ES", "NL", "CH"]
      default_allowed   = false
    }
  }
}

resource "cloudconnexa_location_context" "regional_policies" {
  for_each = var.regional_policies

  name        = "Regional Policy - ${title(replace(each.key, "-", " "))}"
  description = each.value.description
  user_groups_ids = [
    for group_name in each.value.user_group_ids :
    group_name == "office_workers" ? cloudconnexa_user_group.office_workers.id :
    group_name == "developers" ? cloudconnexa_user_group.developers.id :
    group_name == "contractors" ? cloudconnexa_user_group.remote_contractors.id :
    cloudconnexa_user_group.high_security_users.id
  ]

  ip_check {
    allowed = true
    dynamic "ips" {
      for_each = each.value.allowed_ips
      content {
        ip          = ips.value.ip
        description = ips.value.description
      }
    }
  }

  country_check {
    allowed   = true
    countries = each.value.allowed_countries
  }

  default_check {
    allowed = each.value.default_allowed
  }
}

# Time-based location policy (for demonstration)
resource "cloudconnexa_location_context" "business_hours_policy" {
  name        = "Business Hours Location Policy"
  description = "Location policy with additional business hours considerations"
  user_groups_ids = [
    cloudconnexa_user_group.office_workers.id,
    cloudconnexa_user_group.developers.id
  ]

  ip_check {
    allowed = true
    ips {
      ip          = "203.0.113.0/24"
      description = "Office network - always allowed"
    }
    ips {
      ip          = "10.0.0.0/8"
      description = "Home office networks"
    }
    ips {
      ip          = "172.16.0.0/12"
      description = "VPN client networks"
    }
  }

  country_check {
    allowed   = true
    countries = ["US", "CA", "GB", "DE", "AU"]
  }

  default_check {
    allowed = false
  }
}

# Outputs
output "location_context_policies" {
  description = "Created location context policies"
  value = {
    office_only_policy           = cloudconnexa_location_context.office_only_policy.id
    global_office_policy         = cloudconnexa_location_context.global_office_policy.id
    executive_flexible_policy    = cloudconnexa_location_context.executive_flexible_policy.id
    contractor_restricted_policy = cloudconnexa_location_context.contractor_restricted_policy.id
    dev_remote_policy            = cloudconnexa_location_context.dev_remote_policy.id
    emergency_access_policy      = cloudconnexa_location_context.emergency_access_policy.id
    geo_blocking_policy          = cloudconnexa_location_context.geo_blocking_policy.id
    business_hours_policy        = cloudconnexa_location_context.business_hours_policy.id
  }
}

output "regional_policies" {
  description = "Regional location context policies"
  value       = { for k, v in cloudconnexa_location_context.regional_policies : k => v.id }
}

output "user_groups_info" {
  description = "User groups used in location contexts"
  value = {
    high_security_users = cloudconnexa_user_group.high_security_users.id
    office_workers      = cloudconnexa_user_group.office_workers.id
    remote_contractors  = cloudconnexa_user_group.remote_contractors.id
    executives          = cloudconnexa_user_group.executives.id
    developers          = cloudconnexa_user_group.developers.id
  }
}

output "policy_summary" {
  description = "Summary of location context policies and their scope"
  value = {
    office_only = {
      id          = cloudconnexa_location_context.office_only_policy.id
      user_groups = 1
      countries   = 1
      ip_ranges   = 3
    }
    global_office = {
      id          = cloudconnexa_location_context.global_office_policy.id
      user_groups = 2
      countries   = 5
      ip_ranges   = 5
    }
    executive_flexible = {
      id          = cloudconnexa_location_context.executive_flexible_policy.id
      user_groups = 1
      countries   = 20
      ip_ranges   = 3
    }
    geo_blocking = {
      id          = cloudconnexa_location_context.geo_blocking_policy.id
      user_groups = 3
      countries   = 26
      ip_ranges   = 2
    }
  }
}
