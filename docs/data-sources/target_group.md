---
subcategory: "Application Load Balancer"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_target_group"
sidebar_current: "docs-ionoscloud_target_group"
description: |-
  Get information on an Target Group
---

# ionoscloud_target_group

The **Target Group** data source can be used to search for and return an existing Application Load Balancer Target Group.
You can provide a string for the name parameter which will be compared with provisioned Application Load Balancer Target Groups.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By Id
```hcl
data "ionoscloud_target_group" "example" {
  id  = <target_group_id>
}
```

### By Name
```hcl
data "ionoscloud_target_group" "example" {
  name  = "Target Group Example"
}
```

### By Name with Partial Match
```hcl
data "ionoscloud_target_group" "example" {
  name          = "Example"
  partial_match = true
}
```

## Argument Reference

* `id` - (Optional) ID of the target group you want to search for.
* `name` - (Optional) Name of an existing target group that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

Either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.


## Attributes Reference

The following attributes are returned by the datasource:

- `id` - The Id of that Target group
- `name` - The name of that Target Group.
- `algorithm` - Balancing algorithm.
- `protocol` - Balancing protocol.
- `protocol_version` - The forwarding protocol version. Value is ignored when protocol is not 'HTTP'.
- `targets` - Array of items in the collection
  - `ip` - The IP of the balanced target VM.
  - `port` - The port of the balanced target service; valid range is 1 to 65535.
  - `weight` - Traffic is distributed in proportion to target weight, relative to the combined weight of all targets. A target with higher weight receives a greater share of traffic. Valid range is 0 to 256 and default is 1; targets with weight of 0 do not participate in load balancing but still accept persistent connections. It is best use values in the middle of the range to leave room for later adjustments.
  - `proxy_protocol` - The proxy protocol version.
  - `health_check_enabled` - Makes the target available only if it accepts periodic health check TCP connection attempts; when turned off, the target is considered always available. The health check only consists of a connection attempt to the address and port of the target. Default is True.
  - `maintenance_enabled` - Maintenance mode prevents the target from receiving balanced traffic.
- `health_check` - Health check attributes for Target Group.
  - `check_timeout` - The maximum time in milliseconds to wait for a target to respond to a check. For target VMs with 'Check Interval' set, the lesser of the two  values is used once the TCP connection is established.
  - `check_interval` - The interval in milliseconds between consecutive health checks; default is 2000.
  - `retries` - The maximum number of attempts to reconnect to a target after a connection failure. Valid range is 0 to 65535, and default is three reconnection.
- `http_health_check` - Http health check attributes for Target Group
  - `path` - The path (destination URL) for the HTTP health check request; the default is /.
  - `method` - The method for the HTTP health check.
  - `match_type` 
  - `response` - The response returned by the request, depending on the match type.
  - `regex`
  - `negate` 
