---
subcategory: "API Gateway"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_apigateway"
sidebar_current: "docs-datasource-apigateway"
description: |-
  Reads IonosCloud API Gateway objects.
---

# ionoscloud_apigateway

The **API Gateway data source** can be used to search for and return an existing API Gateway.
You can provide a string for the name parameter which will be compared with provisioned API Gateways.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID

```
data "ionoscloud_apigateway" "example" {
    id = <your_apigateway_id>
}
```

### By Name

Needs to have the resource be previously created, or a depends_on clause to ensure that the resource is created before
this data source is called.

```
data "ionoscloud_apigateway" "example" {
    name = "example-apigateway"
}
```

## Argument Reference

* `id` - (Optional) ID of an existing API Gateway that you want to search for.
* `name` - (Optional) Name of an existing API Gateway that you want to search for.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using the name filter. Defaults to `false`.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - ID of the API Gateway.
* `name` - The name of the API Gateway.
* `logs` - Shows whether the collection and reporting of logs for observability of this instance is enabled or disabled.
* `metrics` - Shows whether the collection and reporting of metrics for observability of this instance is enabled or disabled.
* `custom_domains`:
    * `name` - The domain name of the distribution.
    * `certificate_id` - The ID of the certificate to use for the distribution.
* `public_endpoint` - The public endpoint of the API Gateway.
