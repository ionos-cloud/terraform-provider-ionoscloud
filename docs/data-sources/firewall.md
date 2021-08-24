---
layout: "ionoscloud"
page_title: "IonosCloud: firewall"
sidebar_current: "docs-datasource-firewall"
description: |-
Get Information on a IonosCloud Firewall
---

# ionoscloud\_firewall

The firewall data source can be used to search for and return an existing FirewallRules. You can provide a string for either id or name parameters which will be compared with provisioned Firewall Rules. If a single match is found, it will be returned. If your search results in multiple matches, an error will be generated. When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_firewall" "test_firewall" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  name	= "test_ds"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of an existing firewall rule that you want to search for.
* `id` - (Optional) ID of the firewall rule you want to search for.
* `datacenter_id` - (Required) The Virtual Data Center ID.
* `server_id` - (Required) The Server ID.
* `nic_id` - (Required) The NIC ID.

Either `name` or   `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the firewall rule.
* `name` - The name of the firewall rule.
* `protocol` - The protocol for the rule: TCP, UDP, ICMP, ANY.
* `source_mac` - Only traffic originating from the respective MAC address is allowed. 
* `source_ip` - Only traffic originating from the respective IPv4 address is allowed.
* `target_ip` - Only traffic directed to the respective IP address of the NIC is allowed.
* `port_range_start` - Defines the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen.
* `port_range_end` - Defines the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen.
* `icmp_type` - Defines the allowed type (from 0 to 254) if the protocol ICMP is chosen.
* `icmp_code` - Defines the allowed code (from 0 to 254) if protocol ICMP is chosen.
