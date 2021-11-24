---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_application_loadbalancer_forwardingrule"
sidebar_current: "docs-ionoscloud_application_loadbalancer_forwardingrule"
description: |-
Get information on an Application Load Balancer Forwarding Rule
---

# ionoscloud_application_loadbalancer_forwardingrule

The Application Load Balancer Forwarding Rule data source can be used to search for and return existing Application Load Balancer Forwarding Rules.

## Example Usage

### By Id
```hcl
data "ionoscloud_application_loadbalancer_forwardingrule" "alb_fwr_example" {
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
  application_loadbalancer_id = ionoscloud_application_loadbalancer.alb_example.id
  id    		= <alb_fwr_id>
}
```

### By Name
```hcl
data "ionoscloud_application_loadbalancer_forwardingrule" "alb_fwr_example" {
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
  application_loadbalancer_id = ionoscloud_application_loadbalancer.alb_example.id
  name    		= "alb_fwr_example"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `application_loadbalancer_id` - (Required) Application Load Balancer's UUID.
* `name` - (Optional) Name of an existing application load balancer that you want to search for.
* `id` - (Optional) ID of the application load balancer you want to search for.

Both `datacenter_id` and `application_loadbalancer_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.


## Attributes Reference

The following attributes are returned by the datasource:

- `id` - Id of Application Load Balancer Forwarding Rule
- `name` - A name of that Application Load Balancer Forwarding Rule.
- `protocol` - Protocol of the balancing.
- `listener_ip` - Listening IP. (inbound).
- `listener_port` - Listening port number. (inbound) (range: 1 to 65535).
- `health_check`
    - `client_timeout` - ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data. If unset the default of 50 seconds will be used.
- `server certificates` - Array of items in that collection.
- `http_rules` 
    - `name` - A name of that Application Load Balancer http rule.
    - `type` - Type of the Http Rule.
    - `target_group` - The UUID of the target group; mandatory for FORWARD action.
    - `drop_query` - Default is false; true for REDIRECT action.
    - `location` - The location for redirecting; mandatory for REDIRECT action.
    - `status_code` - On REDIRECT action it can take the value 301, 302, 303, 307, 308; on STATIC action it is between 200 and 599.
    - `response_message` - The response message of the request; mandatory for STATIC action.
    - `content_type` -  Will be provided by the PAAS Team; default application/json.
    - `conditions` 
        - `type` -  Type of the Http Rule condition.
        - `condition` - Condition of the Http Rule condition.
        - `negate` - Specifies whether the condition is negated or not; default: false.
        - `key` 
        - `value` 