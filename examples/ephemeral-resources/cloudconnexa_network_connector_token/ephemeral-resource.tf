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

# Mints a setup token every time Terraform opens it (once during plan, once during apply).
# The value never reaches the plan or the state file.
ephemeral "cloudconnexa_network_connector_token" "office" {
  connector_id = cloudconnexa_network_connector.office.id
}

# Hand the token to a secret store through a write-only argument. The secret is only rewritten
# when secret_string_wo_version changes, so bump it to rotate the stored token.
resource "aws_secretsmanager_secret" "connector_token" {
  name = "cloudconnexa/office-connector-token"
}

resource "aws_secretsmanager_secret_version" "connector_token" {
  secret_id                = aws_secretsmanager_secret.connector_token.id
  secret_string_wo         = ephemeral.cloudconnexa_network_connector_token.office.token
  secret_string_wo_version = 1
}
