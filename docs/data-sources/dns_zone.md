---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_dns_zone"
sidebar_current: "docs-dns_zone"
description: |-
  Get information on a DNS Zone.
---

# ionoscloud_dns_zone

The **DNS Zone** can be used to search for and return an existing DNS Zone.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

> ⚠️  Only tokens are accepted for authorization in the **ionoscloud_dns_zone** data source. Please ensure you are using tokens as other methods will not be valid.

## Example Usage

### By ID

```hcl
data "ionoscloud_dns_zone" "example" {
  id = <zone_id>
}
```

### By name
```hcl
data "ionoscloud_dns_zone" "example" {
  name = "example.com"
}
```

### By name with partial match
```hcl
data "ionoscloud_dns_zone" "example" {
  name = "example"
  partial_match = true
}
```

## Argument reference
* `id` - (Optional)[string] The ID of the DNS Zone you want to search for.
* `name` - (Optional)[string] The name of the DNS Zone you want to search for.
* `partial_match` - (Optional)[bool] Whether partial matching is allowed or not when using name argument. Default value is false.

Either `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The UUID of the DNS Zone.
* `name` - The name of the DNS Zone.
* `description` - The description of the DNS Zone.
* `enabled` - Indicates if the DNS Zone is activated or not.
* `nameservers` - A list of available name servers.