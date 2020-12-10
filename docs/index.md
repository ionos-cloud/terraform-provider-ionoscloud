---
layout: "ionoscloud"
page_title: "Provider: IonosCloud"
sidebar_current: "docs-index"
description: |-
  A provider for IonosCloud.
---

# IonosCloud Provider

The IonosCloud provider gives the ability to deploy and configure resources using the IonosCloud Cloud API.

Use the navigation to the left to read about the available data sources and resources.

## Migrating from the ProfitBricks provider

### Provider Name in HCL files

The provider name changed from `profitbricks` to `ionoscloud`.
This reflects in the following change in your terraform hcl files:
`provider "profitbricks"` becomes `provider "ionoscloud"`

### Resources and Datasources in HCL files

The migration affects resource names and datasource names.
Every resource and datasource changed its prefix from `profitbricks_` to `ionoscloud_`.

In order to accommodate that, the terraform hcl files must be updated.

This can be done with a simple find and replace procedure.
For example, on Linux, sed can be used:
```
$ sed -i 's/profitbricks_/ionoscloud_/g' ./main.tf
```

On OSX the same command becomes:
```
$ sed -i bak 's/profitbricks_/ionoscloud_/g' ./main.tf
```

### Terraform State

Because of the name changes of resources and datasources, the terraform state must also be updated.
The local state, in json format, can be updated by replacing `profitbricks_` with `ionoscloud_` directly in the state file.
For example, on Linux, using:
```
$ sed -i 's/profitbricks_/ionoscloud_/g' ./terraform.tfstate
```

On OSX the same command becomes:
```
$ sed -i bak 's/profitbricks_/ionoscloud_/g' ./terraform.tfstate
```

The `provider` entries must also be updated. For example:
```
"provider": "provider[\"registry.terraform.io/hashicorp/profitbricks\"]"
```
becomes
```
"provider": "provider[\"registry.terraform.io/hashicorp/ionoscloud\"]"
```

If you manage your state using remote backends you need to take the appropriate steps specific to your backend.

### Environment Variables

The following env variables have changed:

| Old Variable Name     | New Variable Name |
|-----------------------|-------------------|
| PROFITBRICKS_USERNAME | IONOS_USERNAME    |
| PROFITBRICKS_PASSWORD | IONOS_PASSWORD    |
| PROFITBRICKS_TOKEN    | IONOS_TOKEN       |
| PROFITBRICKS_API_URL  | IONOS_API_URL     |

## Usage

The provider needs to be configured with proper credentials before it can be used.

```hcl
$ export IONOS_USERNAME="ionoscloud_username"
$ export IONOS_PASSWORD="ionoscloud_password"
$ export IONOS_API_URL="ionoscloud_cloud_api_url"
```

Or you can provide your credentials in a `.tf` configuration file as shown in this example.

## Example Usage

```hcl
provider "ionoscloud" {
  username = "ionoscloud_username"
  password = "ionoscloud_password"
  endpoint = "ionoscloud_cloud_api_url"
}

resource "ionoscloud_datacenter" "main" {
  # ...
}
```

**Note**: The credentials provided in a `.tf` file will override the credentials from environment variables.

## Configuration Reference

The following arguments are supported:

- `username` - (Required) If omitted, the `IONOS_USERNAME` environment variable is used. The username is generally an e-mail address in 'username@domain.tld' format.

- `password` - (Required) If omitted, the `IONOS_PASSWORD` environment variable is used.

- `endpoint` - (Optional) If omitted, the `IONOS_API_URL` environment variable is used, or it defaults to the current Cloud API release.

- `retries` - (Deprecated) Number of retries while waiting for a resource to be provisioned. Default value is 50. **Note**: This argument has been deprecated and replaced by the implementation of resource timeouts described below.

## Resource Timeout

Individual resources may provide a `timeouts` block to configure the amount of time a specific operation is allowed to take before being considered an error. Each resource may provide configurable timeouts for the `create`, `update`, and `delete` operations. Each resource that supports timeouts will have or inherit default values for that operation.
Users can overwrite the default values for a specific resource in the configuration.

The default `timeouts` values are:

- create - (Default 60 minutes) Used for creating a resource.
- update - (Default 60 minutes) Used for updating a resource .
- delete - (Default 60 minutes) Used for destroying a resource.
- default - (Default 60 minutes) Used for every other action on a resource.

An example of overwriting the `create`, `update`, and `delete` timeouts:

```hcl
resource "ionoscloud_server" "example" {
  name              = "server"
  datacenter_id     = "${ionoscloud_datacenter.example.id}"
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
  cpu_family        = "AMD_OPTERON"

  volume {
    name           = "new"
    image_name     = "${var.ubuntu}"
    size           = 5
    disk_type      = "SSD"
    ssh_key_path   = "${var.private_key_path}"
    image_password = "test1234"
  }

  nic {
    lan             = "${ionoscloud_lan.example.id}"
    dhcp            = true
    ip              = "${ionoscloud_ipblock.example.ips[0]}"
    firewall_active = true

    firewall {
      protocol         = "TCP"
      name             = "SSH"
      port_range_start = 22
      port_range_end   = 22
    }
  }

  timeouts {
    create = "30m"
    update = "300s"
    delete = "2h"
  }
}

```

Valid units of time should be expressed in "s", "m", "h" for "seconds", "minutes", and "hours" respectively.

Individual resources must opt-in to providing configurable `timeouts`, and attempting to configure values for a resource that does not support `timeouts`, or overwriting a specific action that the resource does not specify as an option, will result in an error.

~> **Note:** Terraform does not automatically rollback in the face of errors.
Instead, your Terraform state file will be partially updated with
any resources that successfully completed.

## Support

You are welcome to contact us with questions or comments at [IonosCloud DevOps Central](https://devops.ionos.com/).
