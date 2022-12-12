---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: firewall"
sidebar_current: "docs-resource-firewall"
description: |-
  Creates and manages Firewall Rules.
---

# ionoscloud\_firewall

Manages a set of **Firewall Rules** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
	name                = "Datacenter Example"
	location            = "us/las"
	description         = "Datacenter Description"
	sec_auth_protection = false
}

resource "ionoscloud_ipblock" "example" {
    location            = ionoscloud_datacenter.example.location
    size                = 2
    name                = "IP Block Example"
}

resource "ionoscloud_server" "example" {
    name                  = "Server Example"
    datacenter_id         = ionoscloud_datacenter.example.id
    cores                 = 1
    ram                   = 1024
    availability_zone     = "ZONE_1"
    cpu_family            = "AMD_OPTERON"
    image_name            = "Ubuntu-20.04"
    image_password        = random_password.server_image_password.result
    volume {
      name                = "system"
      size                = 14
      disk_type           = "SSD"
    }
    nic {
      lan                 = "1"
      dhcp                = true
      firewall_active     = true
    }
}

resource "ionoscloud_nic" "example" {
    datacenter_id         = ionoscloud_datacenter.example.id
    server_id             = ionoscloud_server.example.id
    lan                   = 2
    dhcp                  = true
    firewall_active       = true
    name                  = "Nic Example"
}

resource "ionoscloud_firewall" "example" {
    datacenter_id         = ionoscloud_datacenter.example.id
    server_id             = ionoscloud_server.example.id
    nic_id                = ionoscloud_nic.example.id
    protocol              = "ICMP"
    name                  = "Firewall Example"
    source_mac            = "00:0a:95:9d:68:16"
    source_ip             = ionoscloud_ipblock.example.ips[0]
    target_ip             = ionoscloud_ipblock.example.ips[1]
    icmp_type             = 1
    icmp_code             = 8
    type                  = "INGRESS"
}
resource "random_password" "server_image_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The Virtual Data Center ID.
* `server_id` - (Required)[string] The Server ID.
* `nic_id` - (Required)[string] The NIC ID.
* `protocol` - (Required)[string] The protocol for the rule: TCP, UDP, ICMP, ANY. Property cannot be modified after creation (disallowed in update requests).
* `name` - (Optional)[string] The name of the firewall rule.
* `source_mac` - (Optional)[string] Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Value null allows all source MAC address. Valid format: aa:bb:cc:dd:ee:ff.
* `source_ip` -  (Optional)[string] Only traffic originating from the respective IPv4 address is allowed. Value null allows all source IPs.
* `target_ip` - (Optional)[string] In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Value null allows all target IPs.
* `port_range_start` - (Optional)[int] Defines the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd null to allow all ports.
* `port_range_end` - (Optional)[int] Defines the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd null to allow all ports.
* `icmp_type` - (Optional)[string] Defines the allowed code (from 0 to 254) if protocol ICMP is chosen. Value null allows all codes.
* `icmp_code` - (Optional)[int] Defines the allowed code (from 0 to 254) if protocol ICMP is chosen.
* `type` - (Optional)[string] The type of firewall rule. If is not specified, it will take the default value INGRESS.

## Import

Resource Firewall can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_firewall.myfwrule {datacenter uuid}/{server uuid}/{nic uuid}/{firewall uuid}
```
