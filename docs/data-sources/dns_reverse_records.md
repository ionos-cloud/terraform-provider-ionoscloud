---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_dns_reverse_records"
sidebar_current: "docs-dns_reverse_record"
description: |-
  Get information on DNS Rverse Records.
---

# ionoscloud_dns_reverse_records

The **DNS Rverse Records** can be used to search for and return existing DNS Rverse Records.
Multiple matches will be returned.

> ⚠️  Only tokens are accepted for authorization in the **ionoscloud_dns_reverse_records** data source. Please ensure you are using tokens as other methods will not be valid.

## Example Usage

### By name
```hcl
data "ionoscloud_dns_reverse_records" "example" {
  name = "recordexample"
}
```

### By name with partial match
```hcl
data "ionoscloud_dns_reverse_records" "example" {
  name = "record"
  partial_match = true
}
```

### By IPs
```hcl
data "ionoscloud_dns_reverse_records" "example" {
  ips = ["exampleIP1", "exampleIP2"]
}
```

## Argument reference
* `name` - (Optional)[string] The name of the DNS Rverse Record you want to search for.
* `ips` - (Optional)[list of string] The IPs of the DNS Rverse Records you want to search for.
* `partial_match` - (Optional)[bool] Whether partial matching is allowed or not when using name argument. Default value is false.


## Attributes Reference

The following attributes are returned by the datasource:
* `reverse_records` list of
    * `id` - The UUID of the DNS Rverse Record.
    * `name` - The reverse DNS record name.
    * `ip` - Specifies for which IP address the reverse record should be created. The IP addresses needs to be owned by the contract.
    * `description` - Description stored along with the reverse DNS record to describe its usage.
