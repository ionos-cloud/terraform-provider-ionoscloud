---
layout: "ionoscloud"
page_title: "ProfitBricks : ionoscloud_private_crossconnect"
sidebar_current: "docs-ionoscloud-datasource-private-crossconnect"
description: |-
Get information on a Ionos Cloud Private Crossconnects
---

# ionoscloud\_private_crossconnect

The private crossconnect data source can be used to search for and return existing private crossconnects.

## Example Usage

```hcl
data "ionoscloud_private_crossconnect" "pcc_example" {
  name     = "My PCC"
}
```

## Argument Reference

* `name` - (Optional) Name of an existing private crossconnect that you want to search for.
* `id` - (Optional) ID of the private crossconnect you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the private cross-connection.
* `name` - The name of the private cross-connection.
* `description` - The description for the private cross-connection.
* `peers` - Lists LAN's joined to this private cross connect
    * `lan_id`
    * `lan_name`
    * `datacenter_id`
    * `datacenter_name`
    * `location`
* `connectable_datacenters` - Lists datacenters that can be joined to this private cross connect
    * `id`
    * `name`
    * `location`
