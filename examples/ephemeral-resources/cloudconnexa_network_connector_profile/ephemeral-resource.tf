# Requires Terraform >= 1.10 (ephemeral resources) and >= 1.11 for write-only arguments.

resource "cloudconnexa_network" "office" {
  name            = "office"
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_network_connector" "office" {
  name          = "office-connector"
  network_id    = cloudconnexa_network.office.id
  vpn_region_id = "us-east-1"
}

# Returns the connector's current profile whenever Terraform opens it (once during plan, once
# during apply). The value never reaches the plan or the state file.
ephemeral "cloudconnexa_network_connector_profile" "office" {
  connector_id = cloudconnexa_network_connector.office.id
}

# Hand the profile to a secret store through a write-only argument, for the connector host to
# fetch at boot. The secret is only written when secret_string_wo_version changes, so bump it to
# store the profile again.
resource "aws_secretsmanager_secret" "connector_profile" {
  name = "cloudconnexa/office-connector-profile"
}

resource "aws_secretsmanager_secret_version" "connector_profile" {
  secret_id                = aws_secretsmanager_secret.connector_profile.id
  secret_string_wo         = ephemeral.cloudconnexa_network_connector_profile.office.profile
  secret_string_wo_version = 1
}
