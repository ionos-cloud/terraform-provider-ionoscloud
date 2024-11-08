---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: network security group rule"
sidebar_current: "docs-resource-nsg-rule"
description: |-
  Creates and manages IonosCloud Network Security Group Firewall Rule.
---

# ionoscloud_nsg_firewallrule

Manages a **Network Security Group Rule** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name            = "Datacenter NSG Example"
  location        = "de/txl"
}

resource "ionoscloud_nsg" "example" {
  name              = "Example NSG"
  description       = "Example NSG Description"
  datacenter_id     = ionoscloud_datacenter.example.id
}

resource "ionoscloud_nsg_firewallrule" "example" {
  nsg_id            = ionoscloud_nsg.example.id
  datacenter_id     = ionoscloud_datacenter.example.id
  protocol          = "TCP"
  name              = "SG Rule"
  source_mac        = "00:0a:95:9d:68:15"
  source_ip         = "22.231.113.11"
  target_ip         = "22.231.113.75"
  type              = "EGRESS"
}
```

## Argument Reference

The following arguments are supported:
* `nsg_id` - (Required)[string] The ID of a Network Security Group.
* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `name` - (Optional)[string] The name of the Network Security Group.
* `protocol` - (Required)[string] The protocol for the rule: TCP, UDP, ICMP, ANY. Property cannot be modified after creation (disallowed in update requests).
* `name` - (Optional)[string] The name of the firewall rule.
* `source_mac` - (Optional)[string] Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Value null allows all source MAC address. Valid format: aa:bb:cc:dd:ee:ff.
* `source_ip` -  (Optional)[string] Only traffic originating from the respective IPv4 address is allowed. Value null allows all source IPs.
* `target_ip` - (Optional)(Computed)[string] In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Value null allows all target IPs.
* `port_range_start` - (Optional)[int] Defines the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd null to allow all ports.
* `port_range_end` - (Optional)[int] Defines the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd null to allow all ports.
* `icmp_type` - (Optional)[string] Defines the allowed code (from 0 to 254) if protocol ICMP is chosen. Value null allows all codes.
* `icmp_code` - (Optional)[int] Defines the allowed code (from 0 to 254) if protocol ICMP is chosen.
* `type` - (Optional)(Computed)[string] The type of firewall rule. If is not specified, it will take the default value INGRESS.

## Import

Resource Server can be imported using the `resource id`, `nsg id` and `datacenter id`, e.g.

```shell
terraform import ionoscloud_nsg.mynsg {datacenter uuid}/{nsg uuid}/{firewall uuid}
```

Or by using an `import` block.
```hcl
import {
  to = ionoscloud_nsg.imported
  id = "{datacenter uuid}/{nsg uuid}/{firewall uuid}" 
}
  
resource "ionoscloud_nsg_firewallrule" "imported" {
  nsg_id            = ionoscloud_nsg.example.id
  datacenter_id     = ionoscloud_datacenter.example.id
  protocol          = <protocol of the imported rule>
}
```