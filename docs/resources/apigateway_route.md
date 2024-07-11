---
subcategory: "ApiGateway"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_apigateway_route"
sidebar_current: "docs-resource-apigateway-route"
description: |-
  Creates and manages IonosCloud API Gateway Route objects.
---

# ionoscloud_apigateway_route

Manages an **API Gateway Route** on IonosCloud.

## Example Usage

This resource will create an operational API Gateway Route. After this section completes, the provisioner can be called.

```hcl
resource "ionoscloud_apigateway_route" "apigateway_route" {
  name = "apigateway-route"
  type = "http"
  paths = [
    "/foo/*",
    "/bar"
  ]
  methods = [
    "GET",
    "POST"
  ]
  websocket = false
  upstreams {
    scheme       = "http"
    loadbalancer = "round-robin"
    host         = "example.com"
    port         = 80
    weight       = 100
  }
  gateway_id = <your_apigateway_id>
}
```

## Argument reference

* `id` - (Computed)[string] The ID of the API Gateway Route.
* `name` - (Required)[string] Name of the API Gateway Route. Only alphanumeric characters are allowed.
* `gateway_id` - (Required)[string] The ID of the API Gateway that the route belongs to.
* `type` - (Optional)[string] This field specifies the protocol used by the ingress to route traffic to the backend
  service. Default value: `http`.
* `paths` - (Required)[list] The paths that the route should match. Minimum items: 1.
* `methods` - (Required)[list] The HTTP methods that the route should match. Minimum items: 1. Possible values: `GET`,
  `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS`, `HEAD`, `CONNECT`, `TRACE`.
* `websocket` - (Optional)[bool] To enable websocket support. Default value: `false`.
* `upstreams` - (Required) Upstreams information of the API Gateway Route. Minimum items: 1.
    * `scheme` - (Optional)[string] The target URL of the upstream. Default value: `http`.
    * `host` -  (Required)[string] The host of the upstream.
    * `port` -  (Optional)[int] The port of the upstream. Default value: `80`.
    * `loadbalancer` - (Optional)[string] The load balancer algorithm. Default value: `round-robin`.
    * `weight` - (Optional)[int] Weight with which to split traffic to the upstream. Default value: `100`.

## Import

ApiGateway route can be imported using the `apigateway route id`:

```shell
terraform import ionoscloud_apigateway_route.myroute {apigateway uuid}:{apigateway route uuid}
```
