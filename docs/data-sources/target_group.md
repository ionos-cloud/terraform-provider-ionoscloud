---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_target_group"
sidebar_current: "docs-ionoscloud_target_group"
description: |-
Get information on an Target Group
---

# ionoscloud_target_group

The Target Group data source can be used to search for and return existing Application Load Balancer Target Group.

## Example Usage

```hcl
data "ionoscloud_target_group" "example" {
  name			= "example_"
}
```

## Argument Reference

* `name` - (Optional) Name of an existing target group that you want to search for.
* `id` - (Optional) ID of the target group you want to search for.

Either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.


## Attributes Reference

The following attributes are returned by the datasource:

- `id` - The Id of that Target group
- `name` - The name of that Target Group.
- `algorithm` - Algorithm for the balancing.
- `protocol` - Protocol of the balancing.
- `targets` - Array of items in that collection.
    - `ip` - IP of a balanced target VM.
    - `port` - Port of the balanced target service. (range: 1 to 65535).
    - `weight` - Weight parameter is used to adjust the target VM's weight relative to other target VMs.
    - `health_check` - Health check attributes for Network Load Balancer forwarding rule target.
        - `check` - Check specifies whether the target VM's health is checked.
        - `check_interval` - CheckInterval determines the duration (in milliseconds) between consecutive health checks. If unspecified a default of 2000 ms is used.
        - `maintenance` - Maintenance specifies if a target VM should be marked as down, even if it is not.
- `health_check` - Health check attributes for Target Group.
    - `connect_timeout` - It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed. If unset, the default of 5 seconds will be used.
    - `target_timeout` - TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side. If unset, the default of 50 seconds will be used.
    - `retries` - Retries specifies the number of retries to perform on a target VM after a connection failure. If unset, the default value of 3 will be used.
- `http_health_check` - Http health check attributes for Target Group.
    - `path` - The path for the HTTP health check; default: /.
    - `method` - The method for the HTTP health check.
    - `match_type`
    - `response` - The response returned by the request.
    - `regex` 
    - `negate` 
