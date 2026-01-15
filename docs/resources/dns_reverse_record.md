---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_dns_reverse_record"
sidebar_current: "docs-resource-dns_reverse_record"
description: |-
  Creates and manages DNS ReverseRecord objects.
---

# ionoscloud_dns_reverse_record

Manages a [DNS Reverse Record](https://docs.ionos.com/cloud/network-services/cloud-dns/overview).

> ⚠️  Only tokens are accepted for authorization in the **ionoscloud_dns_reverse_record** resource. Please ensure you are using tokens as other methods will not be valid.

## Example Usage

```hcl
resource "ionoscloud_ipblock" "example" {
  location = "de/fra"
  size = 1
  name = "example_ipblock"
}

resource "ionoscloud_dns_reverse_record" "recordexample" {
  name = "reverse.record.example.com"
  description = "example description"
  ip = ionoscloud_ipblock.example.ips[0]
}
```

## Argument reference

* `name` - (Required)[string] The reverse DNS record name.
* `ip` - (Required)[string] Specifies for which IP address the reverse record should be created. The IP addresses needs to be owned by the contract.
* `description` - (Optional)[string] Description stored along with the reverse DNS record to describe its usage.

## Import

In order to import a DNS Reverse Record, you can define an empty DNS Reverse Record resource in the plan:
```hcl
resource "ionoscloud_dns_reverse_record" "example" {
  
}
```
The resource can be imported using the `record_id`, for example:

```shell
terraform import ionoscloud_dns_reverse_record.example record_id
```
