---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: datacenter"
sidebar_current: "docs-resource-datacenter"
description: |-
  Creates and manages IonosCloud Virtual Data Center.
---

# ionoscloud_datacenter

Manages a Virtual [Data Center](https://docs.ionos.com/cloud/set-up-ionos-cloud/get-started/configure-data-center) on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name                = "Datacenter Example"
  location            = "us/las"
  description         = "datacenter description"
  sec_auth_protection = false
}
```

## Attaching a NSG to a Datacenter

#### A single Network Security Group can be attached at any time to a Datacenter. To do this, use the `ionoscloud_datacenter_nsg_selection` and provide the IDs of the NSG and Datacenter to link them. 
#### Deleting the resource or setting the empty string for the `nsg_id` field will de-attach any previously linked NSG from the Datacenter.

```hcl
resource "ionoscloud_datacenter" "example" {
  name            = "Datacenter NSG Example"
  location        = "de/txl"
}
resource "ionoscloud_nsg" "example" {
  name              = "Example NSG"
  description       = "Example NSG Description"
  datacenter_id     = ionoscloud_datacenter.example.id
}
resource "ionoscloud_datacenter_nsg_selection" "example"{
  datacenter_id     = ionoscloud_datacenter.example.id
  nsg_id            = ionoscloud_nsg.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required)[string] The name of the Virtual Data Center.
* `location` - (Required)[string] The regional location where the Virtual Data Center will be created. This argument is immutable.
* `description` - (Optional)[string] Description for the Virtual Data Center.
* `sec_auth_protection` - (Optional) [bool] Boolean value representing if the data center requires extra protection e.g. two factor protection
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
terraform import ionoscloud_datacenter.mydc datacenter uuid
```
