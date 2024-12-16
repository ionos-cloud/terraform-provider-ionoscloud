---
subcategory: "API Gateway"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_apigateway_route"
sidebar_current: "docs-datasource-apigateway-route"
description: |-
  Reads IonosCloud API Gateway Route objects.
---

# ionoscloud_apigateway_route

The **API Gateway Route data source** can be used to search for and return an existing API Gateway route.
You can provide a string for the name parameter which will be compared with provisioned API Gateway routes.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID

```hcl
data "ionoscloud_apigateway_route" "example" {
  id = "your_apigateway_route_id"
  gateway_id = "your_gateway_id"
}
```

### By Name

Needs to have the resource be previously created, or a depends_on clause to ensure that the resource is created before
this data source is called.

```hcl
data "ionoscloud_apigateway_route" "example" {
  name       = "apigateway-route"
  gateway_id = "your_gateway_id"
}
```

## Argument Reference

* `id` - (Optional) ID of an existing API Gateway Route that you want to search for.
* `name` - (Optional) Name of an existing API Gateway Route that you want to search for.
* `gateway_id` - (Required) The ID of the API Gateway that the route belongs to.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - ID of the API Gateway Route.
* `name` - The name of the API Gateway Route.
* `websocket` - Shows whether websocket support is enabled or disabled.
* `type` - This field specifies the protocol used by the ingress to route traffic to the backend service.
* `paths` - The paths that the route should match.
* `methods` - The HTTP methods that the route should match.
* `upstreams`:
    * `scheme` - The target URL of the upstream.
    * `loadbalancer` - The load balancer algorithm.
    * `host` - The host of the upstream.
    * `port` - The port of the upstream.
    * `weight` - Weight with which to split traffic to the upstream.
