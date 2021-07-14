---
layout: "ionoscloud"
page_title: "IonosCloud: datacenter"
sidebar_current: "docs-resource-datacenter"
description: |-
  Creates and manages IonosCloud Virtual Data Center.
---

# ionoscloud\_datacenter

Manages a Virtual Data Center on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name        = "datacenter name"
  location    = "us/las"
  description = "datacenter description"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required)[string] The name of the Virtual Data Center.
* `location` - (Required)[string] The regional location where the Virtual Data Center will be created.
* `description` - (Optional)[string] Description for the Virtual Data Center.
* `sec_auth_protection` - (Optional) [bool] Boolean value representing if the data center requires extra protection e.g. two factor protection


## Import

Resource Datacenter can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_datacenter.mydc {datacenter uuid}
```
