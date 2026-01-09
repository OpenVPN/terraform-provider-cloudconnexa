# Get all sessions
data "cloudconnexa_sessions" "all" {
}

# Get active sessions only
data "cloudconnexa_sessions" "active" {
  status = "ACTIVE"
}

# Get sessions within a date range
data "cloudconnexa_sessions" "date_range" {
  status     = "COMPLETED"
  start_date = "2024-01-01T00:00:00Z"
  end_date   = "2024-12-31T23:59:59Z"
}

# Output session information
output "all_sessions" {
  value = data.cloudconnexa_sessions.all.sessions
}

output "active_session_count" {
  value = length(data.cloudconnexa_sessions.active.sessions)
}

output "session_users" {
  value = distinct([
    for session in data.cloudconnexa_sessions.active.sessions : session.user_name
  ])
}
