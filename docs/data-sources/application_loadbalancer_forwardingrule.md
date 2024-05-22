---
subcategory: "Application Load Balancer"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_application_loadbalancer_forwardingrule"
sidebar_current: "docs-ionoscloud_application_loadbalancer_forwardingrule"
description: |-
  Get information on an Application Load Balancer Forwarding Rule
---

# ionoscloud_application_loadbalancer_forwardingrule

The Application Load Balancer Forwarding Rule data source can be used to search for and return an existing Application Load Balancer Forwarding Rules.
You can provide a string for the name parameter which will be compared with provisioned Application Load Balancers Forwarding Rules.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By Id
```hcl
data "ionoscloud_application_loadbalancer_forwardingrule" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  application_loadbalancer_id = ionoscloud_application_loadbalancer.example.id
  id    		= <alb_fwr_id>
}
```

### By Name
```hcl
data "ionoscloud_application_loadbalancer_forwardingrule" "example" {
  datacenter_id               = ionoscloud_datacenter.example.id
  application_loadbalancer_id = ionoscloud_application_loadbalancer.example.id
  name    		              = "ALB FR Example"
}
```

### By Name with Partial Match
```hcl
data "ionoscloud_application_loadbalancer_forwardingrule" "example" {
  datacenter_id               = ionoscloud_datacenter.example.id
  application_loadbalancer_id = ionoscloud_application_loadbalancer.example.id
  name    		              = "Example"
  partial_match               = true
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `application_loadbalancer_id` - (Required) Application Load Balancer's UUID.
* `id` - (Optional) ID of the application load balancer you want to search for.
* `name` - (Optional) Name of an existing application load balancer that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

Both `datacenter_id` and `application_loadbalancer_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

- `id` - Id of Application Load Balancer Forwarding Rule
- `name` - The name of the Application Load Balancer forwarding rule.
- `protocol` - Balancing protocol.
- `listener_ip` - Listening (inbound) IP.
- `listener_port` - Listening (inbound) port number; valid range is 1 to 65535.
- `client_timeout` - The maximum time in milliseconds to wait for the client to acknowledge or send data; default is 50,000 (50 seconds).
- `server certificates` - Array of items in that collection.
- `http_rules` - Array of items in that collection
    - `name` - The unique name of the Application Load Balancer HTTP rule.
    - `type` - Type of the Http Rule.
    - `target_group` - The UUID of the target group; mandatory for FORWARD action.
    - `drop_query` - Default is false; valid only for REDIRECT actions.
    - `location` - The location for redirecting; mandatory and valid only for REDIRECT actions.
    - `status_code` - Valid only for REDIRECT and STATIC actions. For REDIRECT actions, default is 301 and possible values are 301, 302, 303, 307, and 308. For STATIC actions, default is 503 and valid range is 200 to 599.
    - `response_message` - The response message of the request; mandatory for STATIC action.
    - `content_type` - Valid only for STATIC actions.
    - `conditions` - An array of items in the collection.The action is only performed if each and every condition is met; if no conditions are set, the rule will always be performed.
        * `type` - Type of the Http Rule condition.
        * `condition` - Matching rule for the HTTP rule condition attribute; mandatory for HEADER, PATH, QUERY, METHOD, HOST, and COOKIE types; must be null when type is SOURCE_IP.
        * `negate` - Specifies whether the condition is negated or not; the default is False.
        *  `key` - Must be null when type is PATH, METHOD, HOST, or SOURCE_IP. Key can only be set when type is COOKIES, HEADER, or QUERY.
        *  `value` - Mandatory for conditions CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH; must be null when condition is EXISTS; should be a valid CIDR if provided and if type is SOURCE_IP.
