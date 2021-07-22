---
layout: "ionoscloud"
page_title: "IonosCloud: application_loadbalancer_forwarding_rule"
sidebar_current: "docs-resource-application_loadbalancer_forwarding_rule"
description: |-
Creates and manages IonosCloud Application Load Balancer Forwarding Rule.
---

# ionoscloud_application_loadbalancer_forwarding_rule

Manages an Application Load Balancer Forwarding Rule on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_application_loadbalancer_forwardingrule" "example" {
 datacenter_id = ionoscloud_datacenter.example.id
 application_loadbalancer_id = ionoscloud_application_loadbalancer.example.id
 name = "example"
 protocol = "HTTP"
 listener_ip = "10.12.118.224"
 listener_port = 8080
 health_check {
     client_timeout = 1000
 }
 server_certificates = ["fb007eed-f3a8-4cbd-b529-2dba508c7599"]
 http_rules {
   name = "http_rule"
   type = "REDIRECT"
   drop_query = true
   location =  "www.ionos.com"
   conditions {
     type = "HEADER"
     condition = "EQUALS"
     value = "something"
   }
 }
}
```

## Argument Reference

The following arguments are supported:

- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `application_loadbalancer_id` - (Required)[string] The ID of Application Load Balancer.
- `name` - (Required)[string] A name of that Application Load Balancer forwarding rule.
- `protocol` - (Required)[string] Protocol of the balancing.
- `listener_ip` - (Required)[string] Listening IP. (inbound).
- `listener_port` - (Required)[int] Listening port number. (inbound) (range: 1 to 65535).
- `health_check` - (Optional)
    - `client_timeout` - (Optional)[int] ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data. If unset the default of 50 seconds will be used.
- `server certificates` - (Optional)[list] Array of items in that collection.
- `http_rules` - (Optional)[list]
  - `name` - (Required)[string] A name of that Application Load Balancer http rule.
  - `type` - (Required)[string] Type of the Http Rule.
  - `target_group` - (Optional)[string] The UUID of the target group; mandatory for FORWARD action.
  - `drop_query` - (Optional)[bool] Default is false; true for REDIRECT action.
  - `location` - (Optional)[string] The location for redirecting; mandatory for REDIRECT action.
  - `status_code` - (Optional)[int] On REDIRECT action it can take the value 301, 302, 303, 307, 308; on STATIC action it is between 200 and 599.
  - `response_message` - (Optional)[string] The response message of the request; mandatory for STATIC action.
  - `content_type` - (Optional)[string] Will be provided by the PAAS Team; default application/json.
  - `conditions` - (Optional)[list] 
    - `type` - (Required)[string] Type of the Http Rule condition.
    - `condition` - (Required)[string] Condition of the Http Rule condition.
    - `negate` - (Optional)[bool] Specifies whether the condition is negated or not; default: false.
    - `key` - (Optional)[string] 
    - `value` - (Optional)[string]