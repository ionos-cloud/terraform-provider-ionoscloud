---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: datacenter"
sidebar_current: "docs-resource-datacenter"
description: |-
  Creates and manages IonosCloud Virtual Data Center.
---

# ionoscloud\_datacenter

Manages a Virtual **Data Center** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name                = "Datacenter Example"
  location            = "us/las"
  description         = "datacenter description"
  sec_auth_protection = false
}
```

## Creating or setting default NSG for a datacenter
#### If `create_default_security_group` is set at Datacenter creation, a default NSG is created together with the datacenter, it can also be set at update to create it later. Its value is then saved to `default_created_security_group_id`.
#### To set a custom NSG as default for the datacenter, set an ID value for `security_group_id`. Can only be set for update requests.
###### Note: must specify `security_group_id` ID as string, referencing a NSG is not possible due to resource reference cycle between datacenter and nsg.
#### Unsetting `security_group_id` will unset the default security group from the datacenter.
```hcl
resource "ionoscloud_datacenter" "example" {
  name            = "Datacenter NSG Example"
  location        = "de/txl"
  create_default_security_group = true
# security_group_id = "your custom security group ID. only on update"
}

resource "ionoscloud_nsg" "example" {
  name              = "Example NSG"
  description       = "Example NSG Description"
  datacenter_id     = ionoscloud_datacenter.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required)[string] The name of the Virtual Data Center.
* `location` - (Required)[string] The regional location where the Virtual Data Center will be created. This argument is immutable.
* `description` - (Optional)[string] Description for the Virtual Data Center.
* `sec_auth_protection` - (Optional) [bool] Boolean value representing if the data center requires extra protection e.g. two factor protection
* `create_default_security_group` - (Optional) [bool] If true, a default security group, with predefined rules, will be created for the datacenter. Default value is false. Provide this field only if the `security_group_id` field is missing. You cannot provide both of them.
* `security_group_id` - (Optional) [string] If `create_default_security_group` is set, it will receive the value of that default group. This will become the default security group for the datacenter, replacing the old one if already exists. This security group must already exist prior to this request. Provide this field only if the `create_default_security_group` field is missing. You cannot provide both of them. Can only be set for update requests.
* `default_created_security_group_id` - (Computed)[string] The ID of the default security group created for the datacenter. This field is only available if `create_default_security_group` is set to true.
* `version` - (Computed) The version of that Data Center. Gets incremented with every change
* `features` - (Computed) List of features supported by the location this data center is part of
* `ipv6_cidr_block` - (Computed) The automatically-assigned /56 IPv6 CIDR block if IPv6 is enabled on this virtual data center
* `cpu_architecture` - (Computed) Array of features and CPU families available in a location
  * `cpu_family` - A valid CPU family name
  * `max_cores` - The maximum number of cores available
  * `max_ram` - The maximum number of RAM in MB
  * `vendor` - A valid CPU vendor name

## Import

Resource Datacenter can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_datacenter.mydc {datacenter uuid}
```
