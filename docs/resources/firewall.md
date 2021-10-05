---
layout: "ionoscloud"
page_title: "IonosCloud: firewall"
sidebar_current: "docs-resource-firewall"
description: |-
  Creates and manages Firewall Rules.
---

# ionoscloud\_firewall

Manages a set of firewall rules on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_firewall" "example" {
  datacenter_id    = "${ionoscloud_datacenter.example.id}"
  server_id        = "${ionoscloud_server.example.id}"
  nic_id           = "${ionoscloud_server.example.primary_nic}"
  protocol         = "TCP"
  name             = "test"
  port_range_start = 1
  port_range_end   = 2
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The Virtual Data Center ID.
* `server_id` - (Required)[string] The Server ID.
* `nic_id` - (Required)[string] The NIC ID.
* `protocol` - (Required)[string] The protocol for the rule: TCP, UDP, ICMP, ANY.
* `name` - (Optional)[string] The name of the firewall rule.
* `source_mac` - (Optional)[string] Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff.
* `source_ip` - (Optional)[string] Only traffic originating from the respective IPv4 address is allowed.
* `target_ip` - (Optional)[string] Only traffic directed to the respective IP address of the NIC is allowed.
* `port_range_start` - (Optional)[string] Defines the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen.
* `port_range_end` - (Optional)[string] Defines the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen.
* `icmp_type` - (Optional)[int] Defines the allowed type (from 0 to 254) if the protocol ICMP is chosen.
* `icmp_code` - (Optional)[int] Defines the allowed code (from 0 to 254) if protocol ICMP is chosen.
* `type` - (Optional)[string] The type of firewall rule. If is not specified, it will take the default value INGRESS

## Import

Resource Firewall can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_firewall.myfwrule {datacenter uuid}/{server uuid}/{nic uuid}/{firewall uuid}
```
