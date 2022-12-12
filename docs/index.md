---
layout: "ionoscloud"
page_title: "Provider: IonosCloud"
sidebar_current: "docs-index"
description: |-
  A provider for IonosCloud.
---

# IonosCloud Provider

The IonosCloud provider gives the ability to deploy and configure resources using the IonosCloud APIs.

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
```bash
$ sed -i 's/profitbricks_/ionoscloud_/g' ./main.tf
```

On OSX the same command becomes:
```bash
$ sed -i bak 's/profitbricks_/ionoscloud_/g' ./main.tf
```

### Terraform State

Because of the name changes of resources and datasources, the terraform state must also be updated.
The local state, in json format, can be updated by replacing `profitbricks_` with `ionoscloud_` directly in the state file.
For example, on Linux, using:
```bash
$ sed -i 's/profitbricks_/ionoscloud_/g' ./terraform.tfstate
```

On OSX the same command becomes:
```bash
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

You can set the environment variables for HTTP basic authentication:

```bash
export IONOS_USERNAME="username"
export IONOS_PASSWORD="password"
```

Or you can use token authentication:

```bash
export IONOS_TOKEN="token"
```

Also, you can overwrite the api endpoint: `api.ionos.com` via the following environment variable:

```bash
export IONOS_API_URL="api-url"
```

Note: if `IONOS_API_URL` environment variable is set, it is recommended to have the `api.ionos.com` value (not `api.ionos.com/cloudapi/v5`).


Another way of configuring it, is by providing your credentials/api_url in a `.tf` configuration file as shown in the below example.

```hcl
provider "ionoscloud" {
  username          = var.ionos_username
  password          = var.ionos_password
  token             = var.ionos_token
  endpoint          = var.ionos_api_url
}
```

You can either explicitly write them in the .tf file or use var.name as in the example above. For setting the var.name, environment variables can be used. The environment variables must be in the format TF_VAR_name and this will be checked last for a value. For example:

```bash
export TF_VAR_ionos_username="username"
export TF_VAR_ionos_password="password"
```

## Debugging

In the default mode, the Terraform provider returns only HTTP client errors. These usually consist only of the HTTP status code. There is no clear description of the problem. But if you want to see the API call error messages as well, you need to set the SDK and Terraform provider environment variables.

You can enable logging now using the `IONOS_LOG_LEVEL` env variable. Allowed values: `off`, `debug` and `trace`. Defaults to `off`.

⚠️ **Note:** We recommend you only use `trace` level for debugging purposes. Disable it in your production environments because it can log sensitive data. It logs the full request and response without encryption, even for an HTTPS call.
Verbose request and response logging can also significantly impact your application’s performance.

```bash
$ export IONOS_LOG_LEVEL=debug
```

⚠️ **Note:** `IONOS_DEBUG` is now deprecated and will be removed in a future release.

⚠️ **Note:** We recommend you only use `IONOS_DEBUG` for debugging purposes. Disable it in your production environments because it can log sensitive data. It logs the full request and response without encryption, even for an HTTPS call.
Verbose request and response logging can also significantly impact your application’s performance.

```bash
$ export TF_LOG=debug
$ export IONOS_DEBUG=true
$ terraform apply
```
now you can see the response body incl. api error message:
```json
{
  "httpStatus" : 422,
  "messages" : [ {
    "errorCode" : "200",
    "message" : "[VDC-yy-xxxx] Operation cannot be executed since this Kubernetes Nodepool is already marked for deletion. Current state of the resource is FAILED_DESTROYING."
  }]
}
```


## Example Usage

```hcl

terraform {
  required_providers {
    ionoscloud = {
      source = "ionos-cloud/ionoscloud"
      version = "= 6.2.1"
    }
  }
}

provider "ionoscloud" {
  username = "ionoscloud_username"
  password = "ionoscloud_password"
  endpoint = "ionoscloud_cloud_api_url"
}

resource "ionoscloud_datacenter" "main" {
  # ...
}
```

### Important Notes
* The `required_providers` section needs to be specified in order for terraform to be
  able to find and download the ionoscloud provider
* The credentials provided in a `.tf` file will override the credentials from environment variables.
* Note that there are two major versions `v5` for the Ionos Cloud API v5 and `v6` that works with API v6 - unless you
  specify a strict version constraint, the latest, which is `v6` will be used

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
    name                  = "Server Example"
    datacenter_id         = ionoscloud_datacenter.example.id
    cores                 = 1
    ram                   = 1024
    availability_zone     = "ZONE_1"
    cpu_family            = "AMD_OPTERON"
    image_name            = data.ionoscloud_image.example.id
    image_password        = random_password.server_image_password.result
    type                  = "ENTERPRISE"
    volume {
        name              = "system"
        size              = 5
        disk_type         = "SSD Standard"
        user_data         = "foo"
        bus               = "VIRTIO"
        availability_zone = "ZONE_1"
    }
    nic {
        lan               = ionoscloud_lan.example.id
        name              = "system"
        dhcp              = true
        firewall_active   = true
        firewall_type     = "BIDIRECTIONAL"
        ips               = [ ionoscloud_ipblock.example.ips[0], ionoscloud_ipblock.example.ips[1] ]
        firewall {
            protocol          = "TCP"
            name              = "SSH"
            port_range_start  = 22
            port_range_end    = 22
            source_mac        = "00:0a:95:9d:68:17"
            source_ip         = ionoscloud_ipblock.example.ips[2]
            target_ip         = ionoscloud_ipblock.example.ips[3]
            type              = "EGRESS"
        }
    }
    timeouts {
      create = "30m"
      update = "300s"
      delete = "2h"
    }
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
```

Valid units of time should be expressed in "s", "m", "h" for "seconds", "minutes", and "hours" respectively.

Individual resources must opt-in to providing configurable `timeouts`, and attempting to configure values for a resource that does not support `timeouts`, or overwriting a specific action that the resource does not specify as an option, will result in an error.

~> **Note:** Terraform does not automatically rollback in the face of errors.
Instead, your Terraform state file will be partially updated with
any resources that successfully completed.
