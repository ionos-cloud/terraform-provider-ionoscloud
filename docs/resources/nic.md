---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: nic"
sidebar_current: "docs-resource-nic"
description: |-
  Creates and manages Network Interface objects.
---

# ionoscloud_nic

Manages a [NIC](https://docs.ionos.com/cloud/set-up-ionos-cloud/get-started/configure-data-center#connect-to-the-internet) on IonosCloud.

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

resource "ionoscloud_lan" "example"{
  datacenter_id     = ionoscloud_datacenter.example.id
  public            = true
  name              = "Lan"
}

resource "ionoscloud_server" "example" {
  name                  = "Server Example"
  datacenter_id         = ionoscloud_datacenter.example.id
  cores                 = 1
  ram                   = 1024
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
  lan                   = ionoscloud_lan.example.id
  name                  = "NIC"
  dhcp                  = true
  firewall_active       = true
  firewall_type         = "INGRESS"
  ips                   = [ ionoscloud_ipblock.example.ips[0], ionoscloud_ipblock.example.ips[1] ]
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
```

## Example Usage with IPv6

```hcl
resource "ionoscloud_datacenter" "example" {
  name                = "Datacenter Example"
  location            = "us/las"
  description         = "Datacenter Description"
  sec_auth_protection = false
}

resource "ionoscloud_lan" "example"{
  datacenter_id     = ionoscloud_datacenter.example.id
  public            = true
  name              = "IPv6 Enabled LAN"
  ipv6_cidr_block   = cidrsubnet(ionoscloud_datacenter.example.ipv6_cidr_block,8,2)
}

resource "ionoscloud_server" "example" {
  name                  = "Server Example"
  datacenter_id         = ionoscloud_datacenter.example.id
  cores                 = 1
  ram                   = 1024
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
  lan                   = ionoscloud_lan.example.id
  name                  = "IPv6 Enabled NIC"
  dhcp                  = true
  firewall_active       = true
  firewall_type         = "INGRESS"
  dhcpv6                = false
  ipv6_cidr_block       = cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14)
  ipv6_ips              = [ 
                              cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14),10),
                              cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14),20),
                              cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14),30)
                          ]
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
```
## Example configuring Flowlog

```hcl
resource "ionoscloud_nic" "example" {
  datacenter_id         = ionoscloud_datacenter.example.id
  server_id             = ionoscloud_server.example.id
  lan                   = ionoscloud_lan.example.id
  name                  = "IPV6 and Flowlog Enabled NIC"
  dhcp                  = true
  firewall_active       = true
  firewall_type         = "INGRESS"
  dhcpv6                = false
  ipv6_cidr_block       = cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14)
  ipv6_ips              = [
    cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14),10),
    cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14),20),
    cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14),30)
  ]
  flowlog {
    action    = "ACCEPTED"
    bucket    = "flowlog-bucket"
    direction = "INGRESS"
    name      = "flowlog"
  }  
}
```

This will configure flowlog for accepted ingress traffic and will log it into an existing IONOS Object Storage bucket named `flowlog-bucket`. Any s3 compatible client can be used to create it. Adding a flowlog does not force re-creation of the NIC, but changing any other field than 
`name` will. Deleting a flowlog will also force NIC re-creation.

## Argument reference

- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `server_id` - (Required)[string] The ID of a server.
- `lan` - (Required)[integer] The LAN ID the NIC will sit on.
- `name` - (Optional)[string] The name of the LAN.
- `dhcp` - (Optional)[Boolean] Indicates if the NIC should get an IP address using DHCP (true) or not (false).
- `dhcpv6` - (Optional)[Boolean] Indicates if the NIC should get an IPv6 address using DHCP (true) or not (false).
- `ipv6_cidr_block` - (Computed, Optional) Automatically assigned /80 IPv6 CIDR block if the NIC is connected to an IPv6 enabled LAN. You can also specify an /80 IPv6 CIDR block for the NIC on your own, which must be inside the /64 IPv6 CIDR block of the LAN and unique.
- `ips` - (Optional)[list] Collection of IP addresses assigned to a NIC. Explicitly assigned public IPs need to come from reserved IP blocks, Passing value null or empty array will assign an IP address automatically.
- `ipv6_ips` - (Optional)[list] Collection of IPv6 addresses assigned to a NIC. Explicitly assigned public IPs need to come from the NIC's Ipv6 CIDR block, Passing value null or empty array will assign an IPv6 address automatically from the NIC's CIDR block.
- `firewall_active` - (Optional)[Boolean] If this resource is set to true and is nested under a server resource firewall, with open SSH port, resource must be nested under the NIC.
- `firewall_type` - (Optional) [String] The type of firewall rules that will be allowed on the NIC. If it is not specified it will take the default value INGRESS
- `id` - (Computed) The ID of the NIC.
- `mac` - (Optional) The MAC address of the NIC. Can be set on creation only. If not set, one will be assigned automatically by the API. Immutable, update forces re-creation. 
* `device_number`- (Computed) The Logical Unit Number (LUN) of the storage volume. Null if this NIC was created from CloudAPI and no DCD changes were done on the Datacenter.
* `pci_slot`- (Computed) The PCI slot number of the Nic.
* `flowlog` - (Optional) Only 1 flow log can be configured. Only the name field can change as part of an update. Flow logs holistically capture network information such as source and destination IP addresses, source and destination ports, number of packets, amount of bytes, the start and end time of the recording, and the type of protocol – and log the extent to which your instances are being accessed.
  - `action` - (Required) Specifies the action to be taken when the rule is matched. Possible values: ACCEPTED, REJECTED, ALL. Immutable, update forces re-creation.
  - `bucket` - (Required) Specifies the IONOS Object Storage bucket where the flow log data will be stored. The bucket must exist. Immutable, update forces re-creation.
  - `direction` - (Required) Specifies the traffic direction pattern. Valid values: INGRESS, EGRESS, BIDIRECTIONAL. Immutable, update forces re-creation.
  - `name` - (Required) Specifies the name of the flow log.
- `security_groups_ids` - (Optional) The list of Security Group IDs for the resource. 
    
⚠️ **Note:**: Removing the `flowlog` forces re-creation of the NIC resource.  

## Import

Resource **Nic** can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_nic.mynic datacenter uuid/server uuid/nic uuid
```
## Working with load balancers
Please be aware that when using a NIC in a load balancer, the load balancer will
change the NIC's ID behind the scenes, therefore the plan will always report this change
trying to revert the state to the one specified by your terraform file.
In order to prevent this, use the "lifecycle meta-argument" when declaring your NIC,
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