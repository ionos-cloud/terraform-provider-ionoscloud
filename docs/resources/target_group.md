---
subcategory: "Application Load Balancer"
layout: "ionoscloud"
page_title: "IonosCloud: target_group"
sidebar_current: "docs-resource-target_group"
description: |-
Creates and manages IonosCloud Target Group.
---

# ionoscloud_target_group

Manages a **Target Group** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_target_group" "example" {
    name                      = "Target Group Example" 
    algorithm                 = "ROUND_ROBIN"
    protocol                  = "HTTP"
    targets {
        ip                    = "22.231.2.2"
        port                  = "8080"
        weight                = "1"
        health_check_enabled  = true
        maintenance_enabled   = false
    }
    health_check {
        check_timeout         = 5000
        check_interval        = 50000
        retries               = 2
    }
    http_health_check {
        path                  = "/."
        method                = "GET"
        match_type            = "STATUS_CODE"
        response              = "200"
        regex                 = true
        negate                = true
    }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The name of the target group.
- `algorithm` - (Required)[string] Balancing algorithm.
- `protocol` - (Required)[string] Balancing protocol.
- `targets` - (Required)[list] Array of items in the collection
    - `ip` - (Required)[string] The IP of the balanced target VM.
    - `port` - (Required)[int] The port of the balanced target service; valid range is 1 to 65535.
    - `weight` - (Required)[int] Traffic is distributed in proportion to target weight, relative to the combined weight of all targets. A target with higher weight receives a greater share of traffic. Valid range is 0 to 256 and default is 1; targets with weight of 0 do not participate in load balancing but still accept persistent connections. It is best use values in the middle of the range to leave room for later adjustments.
    - `health_check_enabled` - (Optional)[bool] Makes the target available only if it accepts periodic health check TCP connection attempts; when turned off, the target is considered always available. The health check only consists of a connection attempt to the address and port of the target. Default is True.
    - `maintenance_enabled` - (Optional)[bool] Maintenance mode prevents the target from receiving balanced traffic.
- `health_check` - (Optional) Health check attributes for Target Group.
    - `check_timeout` - (Optional)[int] The maximum time in milliseconds to wait for a target to respond to a check. For target VMs with 'Check Interval' set, the lesser of the two  values is used once the TCP connection is established.
    - `check_interval` - (Optional)[int] The interval in milliseconds between consecutive health checks; default is 2000.
    - `retries` - (Optional)[int] The maximum number of attempts to reconnect to a target after a connection failure. Valid range is 0 to 65535, and default is three reconnection.
- `http_health_check` - (Optional) Http health check attributes for Target Group
    - `path` - (Optional)[string] The path (destination URL) for the HTTP health check request; the default is /.
    - `method` - (Optional)[string] The method for the HTTP health check.
    - `match_type` - (Required)[string] 
    - `response` - (Required)[string] The response returned by the request, depending on the match type.
    - `regex` - (Optional)[bool] 
    - `negate` - (Optional)[bool] 

## Import

Resource Target Group can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_target_group.myTargetGroup {target group uuid}
```