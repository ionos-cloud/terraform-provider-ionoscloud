---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : network security group"
sidebar_current: "docs-datasource-nsg"
description: |-
  Get information on a IonosCloud Network Security Group
---

# ionoscloud_nsg

The **NSG Data source** can be used to search for and return an existing security groups.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_nsg" "example" {
  datacenter_id  = ionoscloud_datacenter.example.id
  id             = nsg_id
}
```

### By Name & Location
```hcl
data "ionoscloud_nsg" "example" {
  datacenter_id  = ionoscloud_datacenter.example.id
  name     = "NSG Example"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `id` - (Optional) Id of an existing Network Security Group that you want to search for.
* `name` - (Optional) Name of an existing Network Security Group that you want to search for.

Either `name`, `location` or `id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `datacenter_id` - UUID of the Virtual Data Center
* `id` - UUID of the Network Security Group
* `name` - The name of the Network Security Group
* `description` - Description for the Network Security Group
* `rule_ids` - List of IDs for the Firewall Rules attached to this group
* `rules` - List of Firewall Rule objects attached to this group
