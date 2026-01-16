---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_dns_reverse_record"
sidebar_current: "docs-dns_reverse_record"
description: |-
  Get information on a DNS Reverse Record.
---

# ionoscloud_dns_reverse_record

The **DNS Reverse Record** can be used to search for and return an existing DNS Reverse Record.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

> ⚠️  Only tokens are accepted for authorization in the **ionoscloud_dns_reverse_record** data source. Please ensure you are using tokens as other methods will not be valid.

## Example Usage

### By ID

```hcl
data "ionoscloud_dns_reverse_record" "example" {
  id = "record_id"
}
```

### By name
```hcl
data "ionoscloud_dns_reverse_record" "example" {
  name = "recordexample"
}
```

### By name with partial match
```hcl
data "ionoscloud_dns_reverse_record" "example" {
  name = "record"
  partial_match = true
}
```

### By IP
```hcl
data "ionoscloud_dns_reverse_record" "example" {
  ip = "exampleIP"
}
```

## Argument reference
* `id` - (Optional)[string] The ID of the DNS Reverse Record you want to search for.
* `name` - (Optional)[string] The name of the DNS Reverse Record you want to search for.
* `ip` - (Optional)[string] The IP of the DNS Reverse Record you want to search for.
* `partial_match` - (Optional)[bool] Whether partial matching is allowed or not when using name argument. Default value is false.

Either `id`, `ip` or `name` must be provided. If none, or more are provided, the datasource will return an error.


## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The UUID of the DNS Reverse Record.
* `name` - The reverse DNS record name.
* `ip` - Specifies for which IP address the reverse record should be created. The IP addresses needs to be owned by the contract.
* `description` - Description stored along with the reverse DNS record to describe its usage.
