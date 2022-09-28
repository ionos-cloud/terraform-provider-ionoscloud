---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: nic"
sidebar_current: "docs-ionoscloud-datasource-nic"
description: |-
  Get information on a Ionos Cloud NIC
---

# ionoscloud_nic

The **Nic data source** can be used to search for and return existing nics.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By ID
```hcl
data "ionoscloud_nic" "example" {
  datacenter_id   = <datancenter_id>
  server_id       = <server_id>
  id			  = <nic_id>
}
```

### By Name
```hcl
data "ionoscloud_nic" "example" {
  datacenter_id   = <datancenter_id>
  server_id       = <server_id>
  name            = "Nic Example"
}
```

### By Name with Partial Match
```hcl
data "ionoscloud_nic" "example" {
  datacenter_id   = <datancenter_id>
  server_id       = <server_id>
  name            = "Example"
  partial_match   = true
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `server_id` - (Required)[string] The ID of a server.
* `id` - (Optional) ID of the nic you want to search for.
* `name` - (Optional)[string] The name of the LAN. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true..
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

`datacenter_id` and either `id` or `name` must be provided. 
If none, are provided, the datasource will return an error.

## Import

Resource Nic can be imported using the `resource id`, e.g.

## Attributes Reference

The following attributes are returned by the datasource:
* `id` - The id of the NIC.
* `datacenter_id` - The ID of a Virtual Data Center.
* `server_id` - The ID of a server.
* `lan` - The LAN ID the NIC will sit on.
* `name` - The name of the LAN.
* `dhcp` - Indicates if the NIC should get an IP address using DHCP (true) or not (false).
* `ips` - Collection of IP addresses assigned to a nic. Explicitly assigned public IPs need to come from reserved IP blocks, Passing value null or empty array will assign an IP address automatically.
* `firewall_active` - If this resource is set to true and is nested under a server resource firewall, with open SSH port, resource must be nested under the NIC.
* `firewall_type` - The type of firewall rules that will be allowed on the NIC. If it is not specified it will take the default value INGRESS
* `mac` - The MAC address of the NIC.
* `device_number`- The Logical Unit Number (LUN) of the storage volume. Null if this NIC was created from CloudAPI and no DCD changes were done on the Datacenter.
* `pci_slot`- The PCI slot number of the Nic.
