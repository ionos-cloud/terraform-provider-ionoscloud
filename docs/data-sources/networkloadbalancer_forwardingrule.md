---
subcategory: "Network Load Balancer"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_networkloadbalancer_forwardingrule"
sidebar_current: "docs-ionoscloud-datasource-networkloadbalancer_forwardingrule"
description: |-
  Get information on a Network Load Balancer Forwarding Rule
---

# ionoscloud_networkloadbalancer_forwardingrule

The **Network Load Balancer Forwarding Rule data source** can be used to search for and return existing network forwarding rules.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please make sure that your resources have unique names.

## Example Usage

## By ID
```hcl
data "ionoscloud_networkloadbalancer_forwardingrule" "example" {
    datacenter_id               = <datacenter_id>
    networkloadbalancer_id      = <networkloadbalancer_id>
    id  			            = <networkloadbalancer_forwardingrule_id>
}
```


## By Name
```hcl
data "ionoscloud_networkloadbalancer_forwardingrule" "example" {
    datacenter_id               = <datacenter_id>
    networkloadbalancer_id      = <networkloadbalancer_id>
    name                        = "Network Load Balancer Forwarding Rule Example"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `networkloadbalancer_id` - (Required) Network Load Balancer's UUID.
* `name` - (Optional) Name of an existing network load balancer forwarding rule that you want to search for. Search by name is case-insensitive, but the whole resource name is required (we do not support partial matching).
* `id` - (Optional) ID of the network load balancer forwarding rule you want to search for.

Both `datacenter_id` and `networkloadbalancer_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of that Network Load Balancer forwarding rule.
* `name` - The name of that Network Load Balancer forwarding rule.
* `algorithm` - Algorithm for the balancing.
* `protocol` - Protocol of the balancing.
* `listener_ip` - Listening IP. (inbound)
* `listener_port` - Listening port number. (inbound) (range: 1 to 65535)
* `health_check` - Health check attributes for Network Load Balancer forwarding rule.
    * `client_timeout` - ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data. If unset the default of 50 seconds will be used.
    * `connect_timeout` - It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed. If unset, the default of 5 seconds will be used.
    * `target_timeout` - TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side. If unset, the default of 50 seconds will be used.
    * `retries` - Retries specifies the number of retries to perform on a target VM after a connection failure. If unset, the default value of 3 will be used.
* `targets` - Array of items in that collection.
    * `ip` -  IP of a balanced target VM.
    * `port` - Port of the balanced target service. (range: 1 to 65535).
    * `weight` - Weight parameter is used to adjust the target VM's weight relative to other target VMs.
    * `health_check` -  Health check attributes for Network Load Balancer forwarding rule target.
        * `check` - Check specifies whether the target VM's health is checked.
        * `check_interval` - CheckInterval determines the duration (in milliseconds) between consecutive health checks. If unspecified a default of 2000 ms is used.
        * `maintenance` - Maintenance specifies if a target VM should be marked as down, even if it is not.
