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
When this happens, please refine your search string so that it is specific enough to return only one result.

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
 * `name` - (Optional) Name of an existing Virtual Data Center that you want to search for.
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
* `security_group_id` - If `create_default_security_group` is set, will be receive the value of that default group. This will become the default security group for the datacenter, replacing the old one if already exists. This security group must already exist prior to this request. Provide this field only if the `create_default_security_group` field is missing. You cannot provide both of them. Can only be set for update requests.
* `default_created_security_group_id` - The ID of the default security group created for the datacenter. This field is only available if `create_default_security_group` is set to true.
* `cpu_architecture` - Array of features and CPU families available in a location
  * `cpu_family` - A valid CPU family name
  * `max_cores` - The maximum number of cores available
  * `max_ram` - The maximum number of RAM in MB
  * `vendor` - A valid CPU vendor name
