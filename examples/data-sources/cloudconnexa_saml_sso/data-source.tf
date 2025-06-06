# Basic SAML SSO data source
data "cloudconnexa_saml_sso" "current" {}

# Output current SAML SSO configuration
output "saml_sso_config" {
  description = "Current SAML SSO configuration"
  value = {
    enabled                   = data.cloudconnexa_saml_sso.current.enabled
    entity_id                 = data.cloudconnexa_saml_sso.current.entity_id
    sso_url                   = data.cloudconnexa_saml_sso.current.sso_url
    slo_url                   = data.cloudconnexa_saml_sso.current.slo_url
    auto_provision_users      = data.cloudconnexa_saml_sso.current.auto_provision_users
    default_user_group        = data.cloudconnexa_saml_sso.current.default_user_group
    require_signed_assertions = data.cloudconnexa_saml_sso.current.require_signed_assertions
    name_id_format            = data.cloudconnexa_saml_sso.current.name_id_format
  }
}

# Output SAML attribute mapping
output "saml_attribute_mapping" {
  description = "SAML attribute mapping configuration"
  value = length(data.cloudconnexa_saml_sso.current.attribute_mapping) > 0 ? {
    email      = data.cloudconnexa_saml_sso.current.attribute_mapping[0].email
    first_name = data.cloudconnexa_saml_sso.current.attribute_mapping[0].first_name
    last_name  = data.cloudconnexa_saml_sso.current.attribute_mapping[0].last_name
    groups     = data.cloudconnexa_saml_sso.current.attribute_mapping[0].groups
  } : null
}

# Conditional logic based on SAML SSO status
locals {
  saml_enabled = data.cloudconnexa_saml_sso.current.enabled

  # Create different user group configurations based on SAML status
  user_group_config = local.saml_enabled ? {
    connect_auth = "SAML"
    description  = "SAML-authenticated users"
    } : {
    connect_auth = "ON_PRIOR_AUTH"
    description  = "Local authentication users"
  }
}

# Example of using SAML SSO data in other resources
resource "cloudconnexa_user_group" "conditional_auth" {
  name            = "conditional-auth-group"
  description     = local.user_group_config.description
  connect_auth    = local.user_group_config.connect_auth
  internet_access = "SPLIT_TUNNEL_ON"
  max_device      = 3
  vpn_region_ids  = ["us-east-1", "eu-central-1"]
}

# Check if SAML is properly configured
output "saml_configuration_status" {
  description = "SAML configuration validation"
  value = {
    is_enabled             = data.cloudconnexa_saml_sso.current.enabled
    has_sso_url            = data.cloudconnexa_saml_sso.current.sso_url != null && data.cloudconnexa_saml_sso.current.sso_url != ""
    has_entity_id          = data.cloudconnexa_saml_sso.current.entity_id != null && data.cloudconnexa_saml_sso.current.entity_id != ""
    auto_provision_enabled = data.cloudconnexa_saml_sso.current.auto_provision_users
    has_default_group      = data.cloudconnexa_saml_sso.current.default_user_group != null && data.cloudconnexa_saml_sso.current.default_user_group != ""

    # Overall configuration status
    is_fully_configured = (
      data.cloudconnexa_saml_sso.current.enabled &&
      data.cloudconnexa_saml_sso.current.sso_url != null &&
      data.cloudconnexa_saml_sso.current.sso_url != "" &&
      data.cloudconnexa_saml_sso.current.entity_id != null &&
      data.cloudconnexa_saml_sso.current.entity_id != ""
    )
  }
}
