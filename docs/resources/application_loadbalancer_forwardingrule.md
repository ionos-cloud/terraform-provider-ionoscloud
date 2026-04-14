---
subcategory: "Application Load Balancer"
layout: "ionoscloud"
page_title: "IonosCloud: application_loadbalancer_forwardingrule"
sidebar_current: "docs-resource-application_loadbalancer_forwardingrule"
description: |-
  Creates and manages IonosCloud Application Load Balancer Forwarding Rule.
---

# ionoscloud_application_loadbalancer_forwardingrule

Manages an **Application Load Balancer Forwarding Rule** on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_datacenter" "example" {
  name                = "Datacenter Example"
  location            = "us/las"
  description         = "datacenter description"
  sec_auth_protection = false
}

resource "ionoscloud_lan" "example_1" {
  datacenter_id         = ionoscloud_datacenter.example.id
  public                = true
  name                  = "Lan Example"
}

resource "ionoscloud_lan" "example_2" {
  datacenter_id         = ionoscloud_datacenter.example.id
  public                = true
  name                  = "Lan Example"
}

resource "ionoscloud_application_loadbalancer" "example" {
  datacenter_id               = ionoscloud_datacenter.example.id
  name                        = "ALB Example"
  listener_lan                = ionoscloud_lan.example_1.id
  ips                         = [ "10.12.118.224"]
  target_lan                  = ionoscloud_lan.example_2.id
  lb_private_ips              = [ "10.13.72.225/24"]
}

resource "ionoscloud_application_loadbalancer_forwardingrule" "example" {
  datacenter_id               = ionoscloud_datacenter.example.id
  application_loadbalancer_id = ionoscloud_application_loadbalancer.example.id
  name                        = "ALB FR Example"
  protocol                    = "HTTP"
  listener_ip                 = "10.12.118.224"
  listener_port               = 8080
  client_timeout              = 1000
  http_rules {
    name                    = "http_rule"
    type                    = "REDIRECT"
    drop_query              = true
    location                =  "www.ionos.com"
    status_code             =  301
    conditions {
      type                = "HEADER"
      condition           = "EQUALS"
      negate              = true
      key                 = "key"
      value               = "10.12.120.224/24"
    }
  }
  http_rules {
      name                    = "http_rule_2"
      type                    = "STATIC"
      drop_query              = false
      status_code             = 303
      response_message        = "Response"
      content_type            = "text/plain"
      conditions {
        type                = "QUERY"
        condition           = "MATCHES"
        negate              = false
        key                 = "key"
        value               = "10.12.120.224/24"
      }
  }
  server_certificates = [ ionoscloud_certificate.cert.id ]
}
#optionally you can add a certificate to the application load balancer
resource "ionoscloud_certificate" "cert" {
  name = "add_name_here"
  certificate = "your_certificate"
  certificate_chain = "your_certificate_chain"
  private_key = "your_private_key"
}
```

## Argument Reference

The following arguments are supported:

- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `application_loadbalancer_id` - (Required)[string] The ID of Application Load Balancer.
- `name` - (Required)[string] The name of the Application Load Balancer forwarding rule.
- `protocol` - (Required)[string] Balancing protocol.
- `listener_ip` - (Required)[string] Listening (inbound) IP.
- `listener_port` - (Required)[int] Listening (inbound) port number; valid range is 1 to 65535.
- `client_timeout` - (Optional)[int] The maximum time in milliseconds to wait for the client to acknowledge or send data; default is 50,000 (50 seconds).
- `server_certificates` - (Optional)[list] Array of certificate ids. You can create certificates with the [certificate](certificate_manager_certificate.md) resource.
- `http_rules` - (Optional)[list] Array of items in that collection
    - `name` - (Required)[string] The unique name of the Application Load Balancer HTTP rule.
    - `type` - (Required)[string] Type of the Http Rule.
    - `target_group` - (Optional)[string] The UUID of the target group; mandatory for FORWARD action.
    - `drop_query` - (Optional)[bool] Default is false; valid only for REDIRECT actions.
    - `location` - (Optional)[string] The location for redirecting; mandatory and valid only for REDIRECT actions.
    - `status_code` - (Optional)[int] Valid only for REDIRECT and STATIC actions. For REDIRECT actions, default is 301 and possible values are 301, 302, 303, 307, and 308. For STATIC actions, default is 503 and valid range is 200 to 599.
    - `response_message` - (Optional)[string] The response message of the request; mandatory for STATIC action.
    - `content_type` - (Optional)[string] Valid only for STATIC actions.
    - `conditions` - (Optional)[list] - An array of items in the collection.The action is only performed if each and every condition is met; if no conditions are set, the rule will always be performed.
        * `type` - (Required)[string] Type of the Http Rule condition.
        * `condition` - (Required)[string] Matching rule for the HTTP rule condition attribute; mandatory for HEADER, PATH, QUERY, METHOD, HOST, and COOKIE types; must be null when type is SOURCE_IP.
        * `negate` - (Optional)[bool] Specifies whether the condition is negated or not; the default is False.
        * `key` - (Optional)[string] Must be null when type is PATH, METHOD, HOST, or SOURCE_IP. Key can only be set when type is COOKIES, HEADER, or QUERY.
        * `value` - (Optional)[string] Mandatory for conditions CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH; must be null when condition is EXISTS; should be a valid CIDR if provided and if type is SOURCE_IP.

## Import

Resource Application Load Balancer Forwarding Rule can be imported using the `resource id`, `alb id` and `datacenter id`, e.g.

```shell
terraform import ionoscloud_application_loadbalancer_forwardingrule.my_application_loadbalancer_forwardingrule datacenter uuid/application_loadbalancer uuid/application_loadbalancer_forwardingrule uuid
```