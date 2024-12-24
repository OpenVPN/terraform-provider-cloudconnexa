#To select all resources in source and destination
resource "cloudconnexa_access_group" "full_mesh" {
  name = "Access Group name"
  description = "Add your description here"
  source {
    type = "NETWORK"
    all_covered = true
  }
  source {
    type = "HOST"
    all_covered = true
  }
  source {
    type = "USER_GROUP"
    all_covered = true
  }
  destination {
    type = "NETWORK"
    all_covered = true
  }
  destination {
    type = "HOST"
    all_covered = true
  }
  destination {
    type = "USER_GROUP"
    all_covered = true
  }
}
#To select specific network in source or destination
resource "cloudconnexa_access_group" "full_mesh" {
  name = "Access Group name"
  description = "Add your description here"
  source {
    type = "NETWORK"
    all_covered = true
    parent = "00000000-0000-0000-0000-000000000001"
  }
  destination {
    type = "NETWORK"
    all_covered = true
    parent = "00000000-0000-0000-0000-000000000002"
  }
}
#To select specific network ip service or application
#Note: only network IP services with selected type: IP_SOURCE can be used in an access group source
resource "cloudconnexa_access_group" "full_mesh" {
  name = "Access Group name"
  description = "Add your description here"
  source {
    type = "NETWORK"
    all_covered = false
    parent = "00000000-0000-0000-0000-000000000001"
    children = ["00000000-0000-0000-0000-000000000002"]
  }
  destination {
    type = "NETWORK"
    all_covered = true
    parent = "00000000-0000-0000-0000-000000000003"
    children = ["00000000-0000-0000-0000-000000000004"]
  }
}
#To select specific host in source or destination
resource "cloudconnexa_access_group" "full_mesh" {
  name = "Access Group name"
  description = "Add your description here"
  source {
    type = "HOST"
    all_covered = true
    parent = "00000000-0000-0000-0000-000000000001"
  }
  destination {
    type = "HOST"
    all_covered = true
    parent = "00000000-0000-0000-0000-000000000002"
  }
}
#To select specific host ip service or application
resource "cloudconnexa_access_group" "full_mesh" {
  name = "Access Group name"
  description = "Add your description here"
  source {
    type = "HOST"
    all_covered = true
    parent = "00000000-0000-0000-0000-000000000001"
  }
  destination {
    type = "HOST"
    all_covered = true
    parent = "00000000-0000-0000-0000-000000000002"
    children = ["00000000-0000-0000-0000-000000000003"]
  }
}
#To select specific user group in source or destination
resource "cloudconnexa_access_group" "full_mesh" {
  name = "Access Group name"
  description = "Add your description here"
  source {
    type = "USER_GROUP"
    all_covered = false
    children = ["00000000-0000-0000-0000-000000000001"]
  }
  destination {
    type = "USER_GROUP"
    all_covered = false
    children = ["00000000-0000-0000-0000-000000000002"]
  }
}