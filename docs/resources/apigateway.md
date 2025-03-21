---
subcategory: "API Gateway"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_apigateway"
sidebar_current: "docs-resource-apigateway"
description: |-
  API Gateway is an application that acts as a "front door" for backend services and APIs, handling client requests and routing them to the appropriate backend.
---

# ionoscloud_apigateway

An [API gateway](https://api.ionos.com/docs/apigateway/v1/#tag/APIGateways) consists of the generic rules and configurations.

## Usage example

```hcl
resource "ionoscloud_apigateway" "example" {
    name              = "example-gateway"
    metrics           = true
}
```

## Argument reference

* `id` - (Computed)[string] The ID of the API Gateway.
* `name` - (Required)[string] The name of the API Gateway.
* `logs` - (Optional)[bool] Enable or disable logging. Defaults to `false`. **NOTE**: Central Logging must be enabled through the Logging API to enable this feature.
* `metrics` - (Optional)[bool] Enable or disable metrics. Defaults to `false`.
* `custom_domains` - (Optional)[list] Custom domains for the API Gateway, a list that contains elements with the following structure:
    * `name` - (Required)[string] The domain name. Externally reachable.
    * `certificate_id` - (Optional)[string] The certificate ID for the domain. Must be a valid certificate in UUID form.
* `public_endpoint` - (Computed)[string] The public endpoint of the API Gateway.

## Import

In order to import an API Gateway, you can define an empty API Gateway resource in the plan:

```
resource "ionoscloud_apigateway" "example" {

}
```


The resource can be imported using the `gateway_id`, for example:

```
terraform import ionoscloud_apigateway.example gateway_id
```
