---
subcategory: "Network Load Balancer"
layout: "ionoscloud"
page_title: "IonosCloud: networkloadbalancer_forwardingrule"
sidebar_current: "docs-resource-networkloadbalancer_forwardingrule"
description: |-
  Creates and manages Network Load Balancer Forwarding Rule objects.
---

# ionoscloud_networkloadbalancer_forwardingrule

Manages a Network Load Balancer Forwarding Rule on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_networkloadbalancer_forwardingrule" "example" {
 datacenter_id = ionoscloud_datacenter.example.id
 networkloadbalancer_id = ionoscloud_networkloadbalancer.example.id
 name = "example"
 algorithm = "SOURCE_IP"
 protocol = "TCP"
 listener_ip = "10.12.118.224"
 listener_port = "8081"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "123"
   health_check {
     check = true
     check_interval = 1000
   }
 }
}
```

## Argument reference

- `name` - (Required)[string] A name of that Network Load Balancer forwarding rule.
- `algorithm` - (Required)[string] Algorithm for the balancing.
- `protocol` - (Required)[string] Protocol of the balancing.
- `listener_ip` - (Required)[string] Listening IP. (inbound)
- `listener_port` - (Required)[int] Listening port number. (inbound) (range: 1 to 65535)
- `health_check` - (Optional) Health check attributes for Network Load Balancer forwarding rule.
    - `client_timeout` - (Optional)[int] ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data. If unset the default of 50 seconds will be used.
    - `connect_timeout` - (Optional)[int] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed. If unset, the default of 5 seconds will be used.
    - `target_timeout` - (Optional)[int] TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side. If unset, the default of 50 seconds will be used.
    - `retries` - (Optional)[int] Retries specifies the number of retries to perform on a target VM after a connection failure. If unset, the default value of 3 will be used.
- `targets` - (Required)[list] Array of items in that collection.
    - `ip` - (Required)[string] IP of a balanced target VM.
    - `port` - (Required)[int] Port of the balanced target service. (range: 1 to 65535).
    - `weight` - (Required)[int] Weight parameter is used to adjust the target VM's weight relative to other target VMs.
    - `health_check` - (Optional) Health check attributes for Network Load Balancer forwarding rule target.
         - `check` - (Optional)[boolean] Check specifies whether the target VM's health is checked.
         - `check_interval` - (Optional)[int] CheckInterval determines the duration (in milliseconds) between consecutive health checks. If unspecified a default of 2000 ms is used.
         - `maintenance` - (Optional)[boolean] Maintenance specifies if a target VM should be marked as down, even if it is not.
- `datacenter_id` - (Required)[string] A Datacenter's UUID.
- `natgateway_id` - (Required)[string] Network Load Balancer's UUID.

## Import

A Network Load Balancer Forwarding Rule resource can be imported using its `resource id`, the `datacenter id` and the `networkloadbalancer id` e.g.

```shell
terraform import ionoscloud_networkloadbalancer_forwardingrule.my_networkloadbalancer_forwardingrule {datacenter uuid}/{networkloadbalancer uuid}/{networkloadbalancer_forwardingrule uuid}
```