# Basic SAML SSO configuration
resource "cloudconnexa_saml_sso" "basic" {
  enabled     = true
  entity_id   = "https://mycompany.cloudconnexa.com"
  sso_url     = "https://idp.mycompany.com/sso/saml"
  slo_url     = "https://idp.mycompany.com/slo/saml"
  certificate = file("${path.module}/saml-cert.pem")

  attribute_mapping {
    email      = "email"
    first_name = "firstName"
    last_name  = "lastName"
    groups     = "groups"
  }

  auto_provision_users      = true
  default_user_group        = cloudconnexa_user_group.saml_users.id
  require_signed_assertions = true
  name_id_format            = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
}

# User group for SAML-authenticated users
resource "cloudconnexa_user_group" "saml_users" {
  name                 = "saml-users"
  description          = "Users authenticated via SAML SSO"
  connect_auth         = "SAML"
  internet_access      = "SPLIT_TUNNEL_ON"
  max_device           = 3
  all_regions_included = true
}

# Advanced SAML SSO configuration with Azure AD
resource "cloudconnexa_saml_sso" "azure_ad" {
  enabled     = true
  entity_id   = "https://mycompany.cloudconnexa.com"
  sso_url     = "https://login.microsoftonline.com/${var.azure_tenant_id}/saml2"
  slo_url     = "https://login.microsoftonline.com/${var.azure_tenant_id}/saml2"
  certificate = var.azure_ad_certificate

  attribute_mapping {
    email      = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"
    first_name = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname"
    last_name  = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname"
    groups     = "http://schemas.microsoft.com/ws/2008/06/identity/claims/groups"
  }

  auto_provision_users      = true
  default_user_group        = cloudconnexa_user_group.azure_ad_users.id
  require_signed_assertions = true
  name_id_format            = "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent"
}

# User group for Azure AD users
resource "cloudconnexa_user_group" "azure_ad_users" {
  name            = "azure-ad-users"
  description     = "Users authenticated via Azure AD SAML"
  connect_auth    = "SAML"
  internet_access = "RESTRICTED_INTERNET"
  max_device      = 2
  vpn_region_ids  = ["us-east-1", "eu-central-1"]
}

# SAML SSO configuration with Okta
resource "cloudconnexa_saml_sso" "okta" {
  enabled     = true
  entity_id   = "https://mycompany.cloudconnexa.com"
  sso_url     = "https://mycompany.okta.com/app/cloudconnexa/exk1234567890/sso/saml"
  slo_url     = "https://mycompany.okta.com/app/cloudconnexa/exk1234567890/slo/saml"
  certificate = var.okta_certificate

  attribute_mapping {
    email      = "email"
    first_name = "firstName"
    last_name  = "lastName"
    groups     = "groups"
  }

  auto_provision_users      = true
  default_user_group        = cloudconnexa_user_group.okta_users.id
  require_signed_assertions = true
  name_id_format            = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
}

# User group for Okta users
resource "cloudconnexa_user_group" "okta_users" {
  name            = "okta-users"
  description     = "Users authenticated via Okta SAML"
  connect_auth    = "SAML"
  internet_access = "SPLIT_TUNNEL_ON"
  max_device      = 5
  vpn_region_ids  = ["us-east-1", "us-west-1", "eu-central-1"]
}

# SAML SSO configuration with Google Workspace
resource "cloudconnexa_saml_sso" "google_workspace" {
  enabled     = true
  entity_id   = "https://mycompany.cloudconnexa.com"
  sso_url     = "https://accounts.google.com/o/saml2/idp?idpid=${var.google_idp_id}"
  slo_url     = "https://accounts.google.com/o/saml2/logout?idpid=${var.google_idp_id}"
  certificate = var.google_workspace_certificate

  attribute_mapping {
    email      = "email"
    first_name = "firstName"
    last_name  = "lastName"
    groups     = "groups"
  }

  auto_provision_users      = false # Manual provisioning for Google Workspace
  require_signed_assertions = true
  name_id_format            = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
}

# Multiple user groups for different departments with SAML auth
resource "cloudconnexa_user_group" "saml_engineering" {
  name            = "saml-engineering"
  description     = "Engineering team with SAML authentication"
  connect_auth    = "SAML"
  internet_access = "SPLIT_TUNNEL_ON"
  max_device      = 5
  vpn_region_ids  = ["us-east-1", "us-west-1", "eu-central-1", "ap-southeast-1"]
}

resource "cloudconnexa_user_group" "saml_sales" {
  name            = "saml-sales"
  description     = "Sales team with SAML authentication"
  connect_auth    = "SAML"
  internet_access = "SPLIT_TUNNEL_ON"
  max_device      = 3
  vpn_region_ids  = ["us-east-1", "eu-central-1"]
}

resource "cloudconnexa_user_group" "saml_support" {
  name            = "saml-support"
  description     = "Support team with SAML authentication"
  connect_auth    = "SAML"
  internet_access = "RESTRICTED_INTERNET"
  max_device      = 2
  vpn_region_ids  = ["us-east-1"]
}

# Variables for SAML configuration
variable "azure_tenant_id" {
  description = "Azure AD tenant ID"
  type        = string
  sensitive   = true
}

variable "azure_ad_certificate" {
  description = "Azure AD SAML certificate"
  type        = string
  sensitive   = true
}

variable "okta_certificate" {
  description = "Okta SAML certificate"
  type        = string
  sensitive   = true
}

variable "google_idp_id" {
  description = "Google Workspace IdP ID"
  type        = string
  sensitive   = true
}

variable "google_workspace_certificate" {
  description = "Google Workspace SAML certificate"
  type        = string
  sensitive   = true
}

# Data source example
data "cloudconnexa_saml_sso" "current" {}

# Outputs
output "saml_sso_status" {
  description = "Current SAML SSO configuration status"
  value = {
    enabled   = data.cloudconnexa_saml_sso.current.enabled
    entity_id = data.cloudconnexa_saml_sso.current.entity_id
    sso_url   = data.cloudconnexa_saml_sso.current.sso_url
  }
}

output "saml_user_groups" {
  description = "SAML-enabled user groups"
  value = {
    engineering = cloudconnexa_user_group.saml_engineering.id
    sales       = cloudconnexa_user_group.saml_sales.id
    support     = cloudconnexa_user_group.saml_support.id
  }
}
