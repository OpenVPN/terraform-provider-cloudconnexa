---
page_title: "cloudconnexa_sessions Data Source - terraform-provider-cloudconnexa"
subcategory: ""
description: |-
  Use cloudconnexa_sessions data source to retrieve VPN session information.
---

# cloudconnexa_sessions (Data Source)

Use `cloudconnexa_sessions` data source to retrieve VPN session information.

## Example Usage

```terraform
# Get all active sessions
data "cloudconnexa_sessions" "active" {
  status = "ACTIVE"
}

# Get sessions within a date range
data "cloudconnexa_sessions" "recent" {
  start_date = "2024-01-01T00:00:00Z"
  end_date   = "2024-01-31T23:59:59Z"
}

output "active_session_count" {
  value = length(data.cloudconnexa_sessions.active.sessions)
}
```

## Schema

### Optional

- `status` (String) Filter sessions by status. Valid values are `ACTIVE`, `COMPLETED`, or `FAILED`.
- `start_date` (String) Filter sessions starting from this date (RFC3339 format).
- `end_date` (String) Filter sessions until this date (RFC3339 format).

### Read-Only

- `id` (String) The ID of this resource.
- `sessions` (List of Object) The list of VPN sessions. (see [below for nested schema](#nestedatt--sessions))

<a id="nestedatt--sessions"></a>
### Nested Schema for `sessions`

Read-Only:

- `session_id` (String) The session ID.
- `user_id` (String) The user ID associated with this session.
- `device_id` (String) The device ID used in this session.
- `region_id` (String) The VPN region ID.
- `bytes_in` (Number) The number of bytes received.
- `bytes_out` (Number) The number of bytes sent.
- `connector_name` (String) The connector name.
- `user_name` (String) The username.
- `device_name` (String) The device name.
- `client_ip` (String) The client IP address.
- `start_date_time` (String) The session start date and time.
- `vpn_ipv4` (String) The VPN IPv4 address assigned to the client.
- `network_name` (String) The network name.
- `region_name` (String) The region name.
- `connection_status` (String) The connection status.
