# User groups for different departments
resource "cloudconnexa_user_group" "developers" {
  name                 = "developers"
  description          = "Development team with full access"
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
  description          = "External contractors with limited access"
  all_regions_included = false
  connect_auth         = "EVERY_TIME"
  internet_access      = "RESTRICTED_INTERNET"
  max_device           = "2"
}

# Admin user example
resource "cloudconnexa_user" "admin_user" {
  username   = "admin.user"
  group_id   = cloudconnexa_user_group.admins.id
  role       = "ADMIN"
  first_name = "Admin"
  last_name  = "User"
  email      = "admin@example.com"
}

# Developer with device management
resource "cloudconnexa_user" "developer" {
  username   = "john.developer"
  group_id   = cloudconnexa_user_group.developers.id
  role       = "MEMBER"
  first_name = "John"
  last_name  = "Developer"
  email      = "john.developer@example.com"

  devices {
    name         = "Johns-Laptop"
    description  = "MacBook Pro for development"
    ipv4_address = "10.0.1.100"
  }
}

# Contractor user with limited permissions
resource "cloudconnexa_user" "contractor" {
  username   = "jane.contractor"
  group_id   = cloudconnexa_user_group.contractors.id
  role       = "MEMBER"
  first_name = "Jane"
  last_name  = "Contractor"
  email      = "jane.contractor@external.com"

  devices {
    name        = "Contractor-Device"
    description = "Temporary device for project work"
  }
}

# Basic user without email (for testing)
resource "cloudconnexa_user" "basic_user" {
  username = "test.user"
  group_id = cloudconnexa_user_group.developers.id
  role     = "MEMBER"
}

# Multiple users example using for_each
variable "team_members" {
  description = "List of team members to create"
  type = map(object({
    email      = string
    first_name = string
    last_name  = string
    group_id   = string
    role       = string
  }))
  default = {
    "alice.smith" = {
      email      = "alice.smith@example.com"
      first_name = "Alice"
      last_name  = "Smith"
      group_id   = "developers"
      role       = "MEMBER"
    }
    "bob.johnson" = {
      email      = "bob.johnson@example.com"
      first_name = "Bob"
      last_name  = "Johnson"
      group_id   = "developers"
      role       = "MEMBER"
    }
  }
}

resource "cloudconnexa_user" "team_members" {
  for_each = var.team_members

  username   = each.key
  first_name = each.value.first_name
  last_name  = each.value.last_name
  email      = each.value.email
  group_id   = each.value.group_id == "developers" ? cloudconnexa_user_group.developers.id : cloudconnexa_user_group.contractors.id
  role       = each.value.role
}

# Outputs
output "user_groups" {
  description = "Created user groups"
  value = {
    developers  = cloudconnexa_user_group.developers.id
    admins      = cloudconnexa_user_group.admins.id
    contractors = cloudconnexa_user_group.contractors.id
  }
}

output "users" {
  description = "Created users"
  value = {
    admin      = cloudconnexa_user.admin_user.id
    developer  = cloudconnexa_user.developer.id
    contractor = cloudconnexa_user.contractor.id
    basic      = cloudconnexa_user.basic_user.id
  }
}

output "team_member_ids" {
  description = "Team member user IDs"
  value       = { for k, v in cloudconnexa_user.team_members : k => v.id }
}
