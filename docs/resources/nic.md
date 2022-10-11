---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: nic"
sidebar_current: "docs-resource-nic"
description: |-
  Creates and manages Network Interface objects.
---

# ionoscloud_nic

Manages a **NIC** on IonosCloud.

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
    image_password        = "K3tTj8G14a3EgKyNeeiY"
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
    name                  = "Nic Example"
    lan                   = 2
    dhcp                  = true
    firewall_active       = true
    firewall_type         = "INGRESS"
    ips                   = [ ionoscloud_ipblock.example.ips[0], ionoscloud_ipblock.example.ips[1] ]
}
```

## Argument reference

- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `server_id` - (Required)[string] The ID of a server.
- `lan` - (Required)[integer] The LAN ID the NIC will sit on.
- `name` - (Optional)[string] The name of the LAN.
- `dhcp` - (Optional)[Boolean] Indicates if the NIC should get an IP address using DHCP (true) or not (false).
- `ips` - (Optional)[list] Collection of IP addresses assigned to a nic. Explicitly assigned public IPs need to come from reserved IP blocks, Passing value null or empty array will assign an IP address automatically.
- `firewall_active` - (Optional)[Boolean] If this resource is set to true and is nested under a server resource firewall, with open SSH port, resource must be nested under the NIC.
- `firewall_type` - (Optional) [String] The type of firewall rules that will be allowed on the NIC. If it is not specified it will take the default value INGRESS
- `id` - (Computed) The ID of the NIC.
- `mac` - (Computed) The MAC address of the NIC.
* `device_number`- (Computed) The Logical Unit Number (LUN) of the storage volume. Null if this NIC was created from CloudAPI and no DCD changes were done on the Datacenter.
* `pci_slot`- (Computed) The PCI slot number of the Nic.

## Import

Resource **Nic** can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_nic.mynic {datacenter uuid}/{server uuid}/{nic uuid}
```
## Working with load balancers
Please be aware that when using a nic in a load balancer, the load balancer will
change the nic's ID behind the scenes, therefore the plan will always report this change
trying to revert the state to the one specified by your terraform file.
In order to prevent this, use the "lifecycle meta-argument" when declaring your nic,
in order to ignore changes to the `lan` attribute:

Here's an example:

```
resource "ionoscloud_nic" "example" {
  datacenter_id     = ionoscloud_datacenter.foobar.id
  server_id         = ionoscloud_server.example.id
  lan               = "2"
  dhcp              = true
  firewall_active   = true
  name              = "updated"
  lifecycle {
    ignore_changes  = [ lan ]
  }
}
```