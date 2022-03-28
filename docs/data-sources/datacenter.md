---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : datacenter"
sidebar_current: "docs-datasource-datacenter"
description: |-
  Get information on a IonosCloud Data Centers
---

# ionoscloud\_datacenter

The **Datacenter data source** can be used to search for and return an existing Virtual Data Center.
You can provide a string for the name and location parameters which will be compared with provisioned Virtual Data Centers.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.


## Example Usage

### By ID 
```hcl
data "ionoscloud_datacenter" "example" {
  id       = <datacenter_id>
}
```

### By Name & Location
```hcl
data "ionoscloud_datacenter" "example" {
  name     = "Datacenter Example"
  location = "us/las"
}
```

## Argument Reference

 * `id` - (Optional) Id of an existing Virtual Data Center that you want to search for.
 * `name` - (Optional) Name of an existing Virtual Data Center that you want to search for.Search by name is case-insensitive, but the whole resource name is required (we do not support partial matching).
 * `location` - (Optional) Id of the existing Virtual Data Center's location.

Either `name`, `location` or `id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:


* `id` - UUID of the Virtual Data Center
* `name` - The name of the Virtual Data Center
* `location` - The regional location where the Virtual Data Center will be created
* `description` - Description for the Virtual Data Center
* `version` - The version of that Data Center. Gets incremented with every change
* `features` - List of features supported by the location this data center is part of
* `sec_auth_protection` - Boolean value representing if the data center requires extra protection e.g. two factor protection
* `cpu_architecture` - Array of features and CPU families available in a location
  * `cpu_family` - A valid CPU family name
  * `max_cores` - The maximum number of cores available
  * `max_ram` - The maximum number of RAM in MB
  * `vendor` - A valid CPU vendor name
