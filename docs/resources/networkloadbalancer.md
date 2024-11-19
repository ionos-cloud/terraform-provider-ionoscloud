---
subcategory: "Network Load Balancer"
layout: "ionoscloud"
page_title: "IonosCloud: networkloadbalancer"
sidebar_current: "docs-resource-networkloadbalancer"
description: |-
  Creates and manages Network Load Balancer objects.
---

# ionoscloud_networkloadbalancer

Manages a **Network Load Balancer**  on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name                    = "Datacenter Example"
  location                = "us/las"
  description             = "Datacenter Description"
  sec_auth_protection     = false
}

resource "ionoscloud_lan" "example1" {
  datacenter_id           = ionoscloud_datacenter.example.id
  public                  = false
  name                    = "Lan Example 1"
}

resource "ionoscloud_lan" "example2" {
  datacenter_id           = ionoscloud_datacenter.example.id
  public                  = false
  name                    = "Lan Example 2"
}

resource "ionoscloud_networkloadbalancer" "example" {
  datacenter_id           = ionoscloud_datacenter.example.id
  name                    = "example"
  listener_lan            = ionoscloud_lan.example1.id
  target_lan              = ionoscloud_lan.example2.id
  ips                     = ["10.12.118.224"]
  lb_private_ips          = ["10.13.72.225/24"]
  central_logging         = true
  logging_format          = "%%{+Q}o %%{-Q}ci - - [%trg] %r %ST %B \"\" \"\" %cp %ms %ft %b %s %TR %Tw %Tc %Tr %Ta %tsc %ac %fc %bc %sc %rc %sq %bq %CC %CS %hrl %hsl"
}
```

## Example configuring Flowlog

```hcl
resource "ionoscloud_networkloadbalancer" "example" {
  datacenter_id           = ionoscloud_datacenter.example.id
  name                    = "nlb with flowlog"
  listener_lan            = ionoscloud_lan.example1.id
  target_lan              = ionoscloud_lan.example2.id
  ips                     = ["10.12.118.224"]
  lb_private_ips          = ["10.13.72.225/24"]
  flowlog {
    action    = "ALL"
    bucket    = "flowlog-bucket"
    direction = "INGRESS"
    name      = "flowlog"
  }  
}

```

This will configure flowlog for ALL(rejected and accepted) ingress traffic and will log it into an existing ionos bucket named `flowlog-bucket`. Any s3 compatible client can be used to create it. Adding a flowlog does not force re-creation or the nic, but changing any other field than
`name` will. Deleting a flowlog will also force nic re-creation.

## Argument reference

- `name` - (Required)[string] A name of that Network Load Balancer.
- `listener_lan` - (Required)[int] Id of the listening LAN. (inbound)
- `ips` - (Optional)[list] Collection of IP addresses of the Network Load Balancer. (inbound and outbound) IP of the listenerLan must be a customer reserved IP for the public load balancer and private IP for the private load balancer.
- `target_lan` - (Required)[int] Id of the balanced private target LAN. (outbound)
- `lb_private_ips` - (Optional)[list] Collection of private IP addresses with subnet mask of the Network Load Balancer. IPs must contain valid subnet mask. If user will not provide any IP then the system will generate one IP with /24 subnet.
- `central_logging` - (Optional)[bool] Turn logging on and off for this product. Default value is 'false'.
- `logging_lormat` - (Optional)[string] Specifies the format of the logs.
- `datacenter_id` - (Required)[string] A Datacenter's UUID.
- `flowlog` - (Optional)[list] Only 1 flow log can be configured. Only the name field can change as part of an update. Flow logs holistically capture network information such as source and destination IP addresses, source and destination ports, number of packets, amount of bytes, the start and end time of the recording, and the type of protocol – and log the extent to which your instances are being accessed.
    - `action` - (Required)[string] Specifies the action to be taken when the rule is matched. Possible values: ACCEPTED, REJECTED, ALL. Immutable, forces re-creation.
    - `bucket` - (Required)[string] Specifies the IONOS Object Storage bucket where the flow log data will be stored. The bucket must exist. Immutable, forces re-creation.
    - `direction` - (Required)[string] Specifies the traffic direction pattern. Valid values: INGRESS, EGRESS, BIDIRECTIONAL. Immutable, forces re-creation.
    - `name` - (Required)[string] Specifies the name of the flow log.

⚠️ **Note:**: Removing the `flowlog` forces re-creation of the network load balancer resource.

## Import

A Network Load Balancer resource can be imported using its `resource id` and the `datacenter id` e.g.

```shell
terraform import ionoscloud_networkloadbalancer.my_networkloadbalancer {datacenter uuid}/{networkloadbalancer uuid}
```