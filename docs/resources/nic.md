---
layout: "ionoscloud"
page_title: "IonosCloud: nic"
sidebar_current: "docs-resource-nic"
description: |-
  Creates and manages Network Interface objects.
---

# ionoscloud_nic

Manages a NIC on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_nic" "example" {
  datacenter_id = "${ionoscloud_datacenter.example.id}"
  server_id     = "${ionoscloud_server.example.id}"
  lan           = 2
  dhcp          = true
  ip            = "${ionoscloud_ipblock.example.ips[0]}"
}
```

## Argument reference

- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `server_id` - (Required)[string] The ID of a server.
- `lan` - (Required)[integer] The LAN ID the NIC will sit on.
- `name` - (Optional)[string] The name of the LAN.
- `dhcp` - (Optional)[Boolean] Indicates if the NIC should get an IP address using DHCP (true) or not (false).
- `ip` - (Optional)[string] IP assigned to the NIC.
- `firewall_active` - (Optional)[Boolean] If this resource is set to true and is nested under a server resource firewall, with open SSH port, resource must be nested under the NIC.
- `nat` - (Optional)[Boolean] Boolean value indicating if the private IP address has outbound access to the public internet.
- `ips` - (Computed) The IP address or addresses assigned to the NIC.

## Import

Resource Nic can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_nic.mynic {datacenter uuid}/{server uuid}/{nic uuid}
```
