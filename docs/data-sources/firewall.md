---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: firewall"
sidebar_current: "docs-datasource-firewall"
description: |-
  Get Information on a IonosCloud Firewall
---

# ionoscloud\_firewall

The **Firewall data source** can be used to search for and return an existing FirewallRules. 
You can provide a string for either id or name parameters which will be compared with provisioned Firewall Rules.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By ID
```hcl
data "ionoscloud_firewall" "example" {
  datacenter_id = <datacenter_id>
  server_id     = <server_id>
  nic_id        = <nic_id>
  id            = <firewall_id>
}
```

### By Name
```hcl
data "ionoscloud_firewall" "example" {
  datacenter_id   = <datacenter_id>
  server_id       = <server_id>
  nic_id          = <nic_id>
  name            = "Firewall Rule Example"
}
```

### By Name with Partial Match
```hcl
data "ionoscloud_firewall" "example" {
  datacenter_id   = <datacenter_id>
  server_id       = <server_id>
  nic_id          = <nic_id>
  name            = "Firewall"
  partial_match   = true
}
```

### By Type
```hcl
data "ionoscloud_firewall" "example" {
  datacenter_id   = <datacenter_id>
  server_id       = <server_id>
  nic_id          = <nic_id>
  type            = "INGRESS"
}
```

### By Protocol
```hcl
data "ionoscloud_firewall" "example" {
  datacenter_id   = <datacenter_id>
  server_id       = <server_id>
  nic_id          = <nic_id>
  protocol            = "ICMP"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) ID of the firewall rule you want to search for.
* `name` - (Optional) Name of an existing firewall rule that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true..
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.
* `datacenter_id` - (Required) The Virtual Data Center ID.
* `server_id` - (Required) The Server ID.
* `nic_id` - (Required) The NIC ID.
* `type` - (Optional) Type of the firewall rule you want to search for.
* `protocol` - (Optional) Protocol of the firewall rule you want to search for.

Either `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the firewall rule.
* `name` - The name of the firewall rule.
* `protocol` - The protocol for the rule: TCP, UDP, ICMP, ANY. This property is immutable.
* `source_mac` - Only traffic originating from the respective MAC address is allowed.
* `source_ip` - Only traffic originating from the respective IPv4 address is allowed.
* `target_ip` - Only traffic directed to the respective IP address of the NIC is allowed.
* `port_range_start` - Defines the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen.
* `port_range_end` - Defines the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen.
* `icmp_type` - Defines the allowed type (from 0 to 254) if the protocol ICMP is chosen.
* `icmp_code` - Defines the allowed code (from 0 to 254) if protocol ICMP is chosen.