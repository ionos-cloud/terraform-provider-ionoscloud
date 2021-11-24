---
layout: "ionoscloud"
page_title: "IonosCloud: target_group"
sidebar_current: "docs-resource-target_group"
description: |-
Creates and manages IonosCloud Target Group.
---

# ionoscloud_target_group

Manages a Target Group on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_target_group" "example" {
 name = "example"
 algorithm = "SOURCE_IP"
 protocol = "HTTP"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "123"
   health_check {
     check = true
     check_interval = 1000
     maintenance = true
   }
 }
 health_check {
     connect_timeout = 5000
     target_timeout = 50000
     retries = 2
 }
 http_health_check {
     path = "/."
     method = "GET"
     match_type = "STATUS_CODE"
     response = "200"
     regex = false
     negate = false"
   }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] A name of that Target Group.
- `algorithm` - (Required)[string] Algorithm for the balancing.
- `protocol` - (Required)[string] Protocol of the balancing.
- `targets` - (Required)[list] Array of items in that collection.
    - `ip` - (Required)[string] IP of a balanced target VM.
    - `port` - (Required)[int] Port of the balanced target service. (range: 1 to 65535).
    - `weight` - (Required)[int] Weight parameter is used to adjust the target VM's weight relative to other target VMs.
    - `health_check` - (Optional) Health check attributes for Network Load Balancer forwarding rule target.
        - `check` - (Optional)[boolean] Check specifies whether the target VM's health is checked.
        - `check_interval` - (Optional)[int] CheckInterval determines the duration (in milliseconds) between consecutive health checks. If unspecified a default of 2000 ms is used.
        - `maintenance` - (Optional)[boolean] Maintenance specifies if a target VM should be marked as down, even if it is not.
- `health_check` - (Optional) Health check attributes for Target Group.
    - `connect_timeout` - (Optional)[int] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed. If unset, the default of 5 seconds will be used.
    - `target_timeout` - (Optional)[int] TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side. If unset, the default of 50 seconds will be used.
    - `retries` - (Optional)[int] Retries specifies the number of retries to perform on a target VM after a connection failure. If unset, the default value of 3 will be used.
- `http_health_check` - (Optional) Http health check attributes for Target Group.
    - `path` - (Optional)[string] The path for the HTTP health check; default: /.
    - `method` - (Optional)[string] The method for the HTTP health check.
    - `match_type` - (Required)[string] 
    - `response` - (Required)[string] The response returned by the request.
    - `regex` - (Optional)[bool] 
    - `negate` - (Optional)[bool] 
