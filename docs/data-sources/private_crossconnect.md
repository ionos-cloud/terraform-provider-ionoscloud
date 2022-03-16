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
When this happens, please refine your search string so that it is specific enough to return only one result.

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

## Argument Reference

* `name` - (Optional) Name of an existing private crossconnect that you want to search for.
* `id` - (Optional) ID of the private crossconnect you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

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