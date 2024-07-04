---
subcategory: "API Gateway"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_apigateway"
sidebar_current: "docs-resource-apigateway"
description: |-
  API Gateway is an application that acts as a "front door" for backend services and APIs, handling client requests and routing them to the appropriate backend.
---

# ionoscloud_apigateway

An API gateway consists of the generic rules and configurations.

## Usage example

```
resource "ionoscloud_apigateway" "example" {
    name              = "example-gateway"
    logs              = true
    metrics           = true
    
    custom_domains {
        name           = "example.com"
        certificate_id = "00000000-0000-0000-0000-000000000000"
    }
    
    custom_domains {
        name           = "example.org"
        certificate_id = "00000000-0000-0000-0000-000000000000"
    }
}
```

## Argument reference

* `name` - (Required)[string] The name of the API Gateway.
* `logs` - (Optional)[bool] Enable or disable logging. Defaults to `false`.
* `metrics` - (Optional)[bool] Enable or disable metrics. Defaults to `false`.
* `custom_domains` - (Optional)[list] Custom domains for the API Gateway, a list that contains elements with the following structure:
    * `name` - (Required)[string] The domain name.
    * `certificate_id` - (Required)[string] The certificate ID for the domain.
* `public_endpoint` - (Computed)[string] The public endpoint of the API Gateway.

## Import

In order to import an API Gateway, you can define an empty API Gateway resource in the plan:

```
resource "ionoscloud_apigateway" "example" {

}
```


The resource can be imported using the `gateway_id`, for example:

```
terraform import ionoscloud_apigateway.example {gateway_id}
```

