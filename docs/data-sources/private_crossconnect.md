---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_private_crossconnect"
sidebar_current: "docs-ionoscloud-datasource-private-crossconnect"
description: |-
  Get information on a Ionos Cloud Private Crossconnects
---

# ionoscloud\_private_crossconnect

The **Private Crossconnect data source** can be used to search for and return existing private crossconnects.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By ID
```hcl
data "ionoscloud_private_crossconnect" "example" {
  id       = <private_crossconnect_id>
}
```

### By Name
```hcl
data "ionoscloud_private_crossconnect" "example" {
  name     = "PCC Example"
}
```

### By Name with Partial Match
```hcl
data "ionoscloud_private_crossconnect" "example" {
  name          = "Example"
  partial_match = true
}
```

### By Connectable Datasources
```hcl
data "ionoscloud_private_crossconnect" "example" {
  name     = "PCC Example"
  connectable_datacenters {
    id = <datacenter_id>
    name = "Datacenter Example"
    location = "us/las"
  }
  connectable_datacenters {
    id = <datacenter_id>
    name = "Datacenter Example 2"
    location = "us/las"
  }
}
```

## Argument Reference

* `id` - (Optional) ID of the private crossconnect you want to search for.
* `name` - (Optional) Name of an existing private crossconnect that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true..
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.
* `connectable_datacenters` - (Optional) A list of Connectable Datacenters of the private crossconnect you want to search for.

Either `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of the found private cross connect
* `name` - Name of private cross connect 
* `description` - Description of private cross connect
* `peers` - Lists LAN's joined to this private cross connect
  * `lan_id` - The id of the cross-connected LAN
  * `lan_name` - The name of the cross-connected LAN
  * `datacenter_id` - The id of the cross-connected datacenter
  * `datacenter_name` - The name of the cross-connected datacenter
  * `location` - The location of the cross-connected datacenter
* `connectable_datacenters` - Lists datacenters that can be joined to this private cross connect
  * `id` - The UUID of the connectable datacenter
  * `name` - The name of the connectable datacenter
  * `location` - The physical location of the connectable datacenter