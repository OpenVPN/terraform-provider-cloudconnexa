# Terraform v1.0.0 migration guide

Important: this migration guide is created to outline migration from older versions of Terraform provider (which where using "/beta" endpoint) to Terraform provider v1.0.0 (which is using "/v1" endpoint).

In v1.0.0 we intruduced BREAKING CHANGES, which we will cover in this migration guide.

## 1/ Data sources will use only "id" field

Prior to v0.5.0 data sources used "name" field as only supported option.
In v0.5.0 we added ability to use either "name" or "id" fields.
In v1.0.0 we removed ability to use "name" fields.

Example:

before:

```hcl
data "cloudconnexa_network" {
  name = "my_network"
}
```

after:

```hcl
data "cloudconnexa_network" {
  id = "5cbe8e9c-c2c5-4a15-9bcf-491cce213adf"
}
```

## 2/ Changes in argument's values for certain resources

This was done to better reflect the name of the options in the UI with name of the values in API and Terraform.
Table below contains references on what changed, and new value:

| Resource | Argument | Beta endpoint (old value) | v1 endpoint (new value) |
| :--------: | :-------: | :-------: | :-------: |
| cloudconnexa_user_group | connect_auth | AUTH | ON_PRIOR_AUTH |
| cloudconnexa_user_group | connect_auth | AUTO | NO_AUTH |
| cloudconnexa_user_group | connect_auth | STRICT_AUTH | EVERY_TIME |
| cloudconnexa_user_group | internet_access | BLOCKED | RESTRICTED_INTERNET |
| cloudconnexa_user_group | internet_access | GLOBAL_INTERNET | SPLIT_TUNNEL_OFF |
| cloudconnexa_user_group | internet_access | LOCAL | SPLIT_TUNNEL_ON |
| cloudconnexa_network | internet_access | BLOCKED | RESTRICTED_INTERNET |
| cloudconnexa_network | internet_access | GLOBAL_INTERNET | SPLIT_TUNNEL_OFF |
| cloudconnexa_network | internet_access | LOCAL | SPLIT_TUNNEL_ON |
| cloudconnexa_host | internet_access | BLOCKED | RESTRICTED_INTERNET |
| cloudconnexa_host | internet_access | GLOBAL_INTERNET | SPLIT_TUNNEL_OFF |
| cloudconnexa_host | internet_access | LOCAL | SPLIT_TUNNEL_ON |
| cloudconnexa_location_context | n/a | default_policy | default_check |
| cloudconnexa_location_context | n/a | country_policy | country_check |
| cloudconnexa_location_context | n/a | ip_policy | ip_check |


Code example for "cloudconnexa_user_group":

before:

```hcl
resource "cloudconnexa_user_group" "ug01" {
  name                 = "ug01"
  all_regions_included = true
  connect_auth         = "AUTH"
  internet_access      = "LOCAL"
  max_device           = "15"
}
```

after:

```hcl
resource "cloudconnexa_user_group" "ug01" {
  name                 = "ug01"
  all_regions_included = true
  connect_auth         = "ON_PRIOR_AUTH"
  internet_access      = "SPLIT_TUNNEL_ON"
  max_device           = "15"
}
```

Code example for "cloudconnexa_location_context":

before:

```hcl
resource "cloudconnexa_location_context" "this" {
  name            = "Location Context Policy"
  description     = "Description for Location Context Policy"
  user_groups_ids = []
  ip_policy {
    allowed = true
    ips {
      ip          = "10.10.0.0/16"
      description = "Test subnet"
    }
    ips {
      ip          = "10.20.0.0/16"
      description = "Test subnet 2"
    }
  }
  country_policy {
    allowed   = true
    countries = ["US", "GB"]
  }
  default_policy {
    allowed = false
  }
}
```

after:

```hcl
resource "cloudconnexa_location_context" "this" {
  name            = "Location Context Policy"
  description     = "Description for Location Context Policy"
  user_groups_ids = []
  ip_check {
    allowed = true
    ips {
      ip          = "10.10.0.0/16"
      description = "Test subnet"
    }
    ips {
      ip          = "10.20.0.0/16"
      description = "Test subnet 2"
    }
  }
  country_check {
    allowed   = true
    countries = ["US", "GB"]
  }
  default_check {
    allowed = false
  }
}
```

Code example for "cloudconnexa_network":

before:

```hcl
resource "cloudconnexa_network" "this" {
  description     = "This is test network"
  name            = "my_network"
  internet_access = "LOCAL"
  default_route {
    subnet = "192.168.144.0/24"
  }
  default_connector {
    name          = "test-connector"
    vpn_region_id = "eu-central-1"
  }
}
```

after:

```hcl
resource "cloudconnexa_network" "this" {
  description     = "This is test network"
  name            = "my_network"
  internet_access = "SPLIT_TUNNEL_ON"
  default_route {
    subnet = "192.168.144.0/24"
  }
  default_connector {
    name          = "test-connector"
    vpn_region_id = "eu-central-1"
  }
}
```

Code example for "cloudconnexa_host":

before:

```hcl
resource "cloudconnexa_host" "this" {
  name        = "my_host"
  description = "This is test host"
  connector {
    name          = "test"
    vpn_region_id = "eu-central-2"
  }
  internet_access = "LOCAL"
}
```

after:

```hcl
resource "cloudconnexa_host" "this" {
  name        = "my_host"
  description = "This is test host"
  connector {
    name          = "test"
    vpn_region_id = "eu-central-2"
  }
  internet_access = "SPLIT_TUNNEL_ON"
}
```

## 3/ Some resources were splitted into separate ones

Before:

- cloudconnexa_application
- cloudconnexa_connector
- cloudconnexa_ip_service

After:

- cloudconnexa_host_application
- cloudconnexa_host_connector
- cloudconnexa_host_ip_service
- cloudconnexa_network_application
- cloudconnexa_network_connector
- cloudconnexa_network_ip_service

If you used previously "cloudconnexa_application", "cloudconnexa_connector" or "cloudconnexa_ip_service" you will have one option:

- Remove from state, rename resource and re-import them into Terraform.

We tried to use "terraform state mv" command as well as "moved" block - but Terraform didn't liked this:

output when running "terraform state mv":

```shell
│ Error: Invalid state move request
│ 
│ Cannot move cloudconnexa_application.test1 to cloudconnexa_network_application.test1: resource types don't match.
```

output when using "moved" block:

```shell
 Error: Resource type mismatch
│ 
│ This statement declares a move from cloudconnexa_application.test1 to cloudconnexa_network_application.test1, which is a resource of a different type.
```

### Remove from state and then import
Let's imagine you have two resources and you use Terraform provider v0.5.1:

```hcl
data "cloudconnexa_network" "test-net" {
  id = "e0a62eed-d034-4cec-8f59-062d96b9f2ab"
}

resource "cloudconnexa_application" "test1" {
  name              = "example-application-1"
  network_item_type = "NETWORK"
  network_item_id   = data.cloudconnexa_network.test-net.id
  routes {
    domain            = "example-application-1.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["ANY"]
  }
}

resource "cloudconnexa_application" "test2" {
  name              = "example-application-2"
  network_item_type = "NETWORK"
  network_item_id   = data.cloudconnexa_network.test-net.id
  routes {
    domain            = "example-application-2.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["ANY"]
  }
}
```
To perform migration follow next procedure:

- run "terraform plan" to get IDs of resources:

```shell
data.cloudconnexa_network.test-net: Reading...
data.cloudconnexa_network.test-net: Read complete after 0s [id=e0a62eed-d034-4cec-8f59-062d96b9f2ab]
cloudconnexa_network_application.test2: Refreshing state... [id=48be819c-f5f7-4c67-9720-30fd908cbda4]
cloudconnexa_network_application.test1: Refreshing state... [id=b1ed3722-0da2-49d5-88f4-515b4ce52690]
```

- Remove from state

```shell
terraform state rm cloudconnexa_application.test1
terraform state rm cloudconnexa_application.test2
```

- Update reference to new provider version

Specify new version of the provider:

```hcl
terraform {
  required_providers {
    cloudconnexa = {
      source = "OpenVPN/cloudconnexa"
      version = "1.0.0"
    }
  }
}
```
 and initialize it:

```shell
terraform init -upgrade
```

- Edit the resource names and code

Updated code (note that we removed "network_item_type" and renamed "network_item_id"):

```hcl
resource "cloudconnexa_network_application" "test1" {
  name = "example-application-1"
  network_id = data.cloudconnexa_network.test-net.id
  routes {
    domain            = "example-application-1.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["ANY"]
  }
}

resource "cloudconnexa_network_application" "test2" {
  name = "example-application-2"
  network_id = data.cloudconnexa_network.test-net.id
  routes {
    domain            = "example-application-2.com"
    allow_embedded_ip = false
  }

  config {
    service_types = ["ANY"]
  }
}
```

- Import resources back into Terraform

```shell
terraform import cloudconnexa_network_application.test1 b1ed3722-0da2-49d5-88f4-515b4ce52690
terraform import cloudconnexa_network_application.test2 48be819c-f5f7-4c67-9720-30fd908cbda4
```

After import Terraform will want to make minor (expected) change, apply it:

```shell
$ terraform apply
.......................
(skipped not essential output)

Terraform will perform the following actions:

  # cloudconnexa_network_application.test1 will be updated in-place
  ~ resource "cloudconnexa_network_application" "test1" {
        id          = "b1ed3722-0da2-49d5-88f4-515b4ce52690"
        name        = "example-application-1"
      + network_id  = "e0a62eed-d034-4cec-8f59-062d96b9f2ab"
        # (1 unchanged attribute hidden)

        # (2 unchanged blocks hidden)
    }

  # cloudconnexa_network_application.test2 will be updated in-place
  ~ resource "cloudconnexa_network_application" "test2" {
        id          = "48be819c-f5f7-4c67-9720-30fd908cbda4"
        name        = "example-application-2"
      + network_id  = "e0a62eed-d034-4cec-8f59-062d96b9f2ab"
        # (1 unchanged attribute hidden)

        # (2 unchanged blocks hidden)
    }

Plan: 0 to add, 2 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes
```

After you applied it, now when you run "terraform plan" - it should return that all is ok:

```shell
No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

PS. This is simple example, for use cases when you have multiple resources and you create them via for_each you may follow this approach [https://developer.hashicorp.com/terraform/language/import#import-multiple-instances-with-for_each]

## 4/ Starting with v1.0.0 when creating resource "cloudconnexa_network" route and connector are to be created separatelly

Previously before Terraform provider v1.0.0 to create "cloudconnexa_network" code looked like this:

```hcl
resource "cloudconnexa_network" "this" {
  name            = "my_test_network"
  description     = "Test network"
  egress          = true
  internet_access = "LOCAL"
  default_route {
    description = "Managed by Terraform"
    subnet      = "192.168.144.0/24"
    type        = "IP_V4"
  }
  default_connector {
    description   = "Managed by Terraform"
    name          = "test-connector"
    vpn_region_id = "eu-central-1"
  }
}
```

Due to the way how API and Terraform provider did things it lead to issue when it was not possible to import already existing networks into Terraform.

Now starting with v1.0.0 code will be like this:

```hcl
resource "cloudconnexa_network" "this" {
  name            = "my_test_network"
  description     = "Test network"
  egress          = true
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_network_connector" "this" {
  name          = "test-connector"
  description   = "Managed by Terraform"
  vpn_region_id = "eu-central-1"
  network_id    = cloudconnexa_network.this.id
}

resource "cloudconnexa_route" "this" {
  description     = "Managed by Terraform"
  type            = "IP_V4"
  subnet          = "192.168.144.0/24"
  network_item_id = cloudconnexa_network.this.id
}
```

When migrating from older versions of provider to v1.0.0 you will have to remove old Network, Connector and Route from Terraform state.

If you have Terraform provider for example v0.5.1 with next code:

```hcl
resource "cloudconnexa_network" "this" {
  name            = "my_test_network"
  description     = "Test network"
  egress          = true
  internet_access = "LOCAL"
  default_route {
    description = "Managed by Terraform"
    subnet      = "192.168.144.0/24"
    type        = "IP_V4"
  }
  default_connector {
    description   = "Managed by Terraform"
    name          = "test-connector"
    vpn_region_id = "eu-central-1"
  }
}
```

and then switch to v1.0.0 version of the provider

```hcl
terraform {
  required_providers {
    cloudconnexa = {
      source = "OpenVPN/cloudconnexa"
      version = "1.0.0"
    }
  }
}
```

you will need to update value of field "internet_access" from "LOCAL" to "SPLIT_TUNNEL_ON".
After you try to run "terraform plan" you will get next output:

```shell
│ Error: Unsupported block type
│ 
│   on main.tf line 6, in resource "cloudconnexa_network" "this":
│    6:   default_route {
│ 
│ Blocks of type "default_route" are not expected here.
╵
╷
│ Error: Unsupported block type
│ 
│   on main.tf line 11, in resource "cloudconnexa_network" "this":
│   11:   default_connector {
│ 
│ Blocks of type "default_connector" are not expected here.
╵
```

Now if you were to remove/comment those blocks ("default_route" and "default_connector") and try again, you will see next output:

```shell
$ terraform plan 
cloudconnexa_network.this: Refreshing state... [id=4099c335-94b0-44c0-83a6-73c4d1417c1c]
╷
│ Error: Failed to load plugin schemas
│ 
│ Error while loading schemas for plugin components: Failed to obtain provider schema: Could not load the schema for provider registry.terraform.io/openvpn/cloudconnexa: failed to instantiate
│ provider "registry.terraform.io/openvpn/cloudconnexa" to obtain schema: unavailable provider "registry.terraform.io/openvpn/cloudconnexa"..
╵
```

Now you need to remove from state:

```shell
terraform state rm cloudconnexa_network.this
```
after that make sure that code was updated:

```hcl
resource "cloudconnexa_network" "this" {
  name            = "my_test_network"
  description     = "Test network"
  egress          = true
  internet_access = "SPLIT_TUNNEL_ON"
}
```

Now import it (ID can be found in CloudConnexa Admin Portal):

```shell
terraform import cloudconnexa_network.this <id>
```

after import will finish try run "terraform plan", it should return "No changes. Your infrastructure matches the configuration."

Next step is to import Network Connector. At this point your code will look sililar to this

```hcl
resource "cloudconnexa_network" "this" {
  name            = "my_test_network"
  description     = "Test network"
  egress          = true
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_network_connector" "this" {
  name          = "test-connector"
  description   = "Managed by Terraform"
  vpn_region_id = "eu-central-1"
  network_id    = cloudconnexa_network.this.id
}
```
Now import it (ID can be found in CloudConnexa Admin Portal):

```shell
terraform import cloudconnexa_network_connector.this <id>
```

after import will finish most likelly you will see similar output from "terraform plan"

```shell
Terraform will perform the following actions:

  # cloudconnexa_network_connector.this will be updated in-place
  ~ resource "cloudconnexa_network_connector" "this" {
      + description   = "Managed by Terraform"
        id            = "d5c750a6-cd4e-4b6e-84d7-97a0c298484e"
        name          = "test-connector"
        # (6 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

Just run "terraform apply" and re-run "terraform plan", it should return "No changes. Your infrastructure matches the configuration."

Importing "route" will require to use "Swagger". Go to https://[replace_with_your_cloud_id].openvpn.com/api and click "Swagger" button, click "Authorize" and provide API credentials.

After that on Swagger page find "Network" section and under it find "GET /api/v1/networks/{id} Get existing network", click "Try it out" and paste ID of the Network (can be found in CloudConnexa Admin Portal) and click "Execute".

It wil give you information about your Network, there you will find block similar to (take a look at "id" field, it will be needed later):

```shell
  "routes": [
    {
      "id": "0dccf62f-7083-43a7-a8a3-38406a21b842",
      "type": "IP_V4",
      "subnet": "192.168.144.0/24",
      "description": "Managed by Terraform"
    }
  ],
```

Now we need to add configuration for Terraform, your final Terraform code would look similar to this:

```hcl
resource "cloudconnexa_network" "this" {
  name            = "my_test_network"
  description     = "Test network"
  egress          = true
  internet_access = "SPLIT_TUNNEL_ON"
}

resource "cloudconnexa_network_connector" "this" {
  name          = "test-connector"
  description   = "Managed by Terraform"
  vpn_region_id = "eu-central-1"
  network_id    = cloudconnexa_network.this.id
}

resource "cloudconnexa_route" "this" {
  description     = "Managed by Terraform"
  type            = "IP_V4"
  network_item_id = cloudconnexa_network.this.id
  subnet          = "192.168.144.0/24"
}
```
Now we can import remaining resource (use value of "id" field from Swagger):

```shell
terraform import cloudconnexa_route.this <id>
```

After import will be done run "terraform plan", it should retun "No changes. Your infrastructure matches the configuration".

Migration complete.
