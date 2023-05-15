---
subcategory: "DNS as a Service"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_dns_record"
sidebar_current: "docs-dns_record"
description: |-
  Get information on a DNS Record.
---

# ionoscloud_dns_record

⚠️ **Note:** DNSaaS is currently in the Early Access (EA) phase.
We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.

The **DNS Record** can be used to search for and return an existing DNS Record.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By ID

```hcl
data "ionoscloud_dns_record" "example" {
  id = <record_id>
  zone_id = <zone_id>
}
```

### By name
```hcl
data "ionoscloud_dns_record" "example" {
  name = "recordexample"
  zone_id = <zone_id>
}
```

### By name with partial match
```hcl
data "ionoscloud_dns_record" "example" {
  name = "record"
  partial_match = true
  zone_id = <zone_id>
}
```

## Argument reference
* `zone_id` - (Required)[string] The ID of the DNS Zone in which the DNS Record can be found.
* `id` - (Optional)[string] The ID of the DNS Record you want to search for.
* `name` - (Optional)[string] The name of the DNS Record you want to search for.
* `partial_match` - (Optional)[bool] Whether partial matching is allowed or not when using name argument. Default value is false.

Either `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.


## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The UUID of the DNS Record.
* `name` - The name of the DNS Record.
* `type` - The type of the DNS Record.
* `content` - The content of the DNS Record.
* `ttl` - The time to live of the DNS Record.
* `enabled` - Indicates if the DNS Record is active or not.