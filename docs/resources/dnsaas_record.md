---
subcategory: "DNS as a Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_dns_record"
sidebar_current: "docs-resource-dns_record"
description: |-
  Creates and manages DNS Record objects.
---

# ionoscloud_dns_record

⚠️ **Note:** DNSaaS is currently in the Early Access (EA) phase.
We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.

Manages a **DNS Record**.

## Example Usage

```hcl
resource "ionoscloud_dns_zone" "example" {
  name = "example.com"
  description = "description"
  enabled = false
}

resource "ionoscloud_dns_record" "recordexample" {
  zone_id = ionoscloud_dns_zone.example.id
  name = "recordexample"
  type = "CNAME"
  content = "1.2.3.4"
  ttl = 2000
  priority = 1024
  enabled = false
}
```

## Argument reference

* `name` - (Required)[string] The name of the DNS Record.
* `type` - (Required)[string] The type of the DNS Record, can have one of these values: `A, AAAA, CNAME, ALIAS, MX, NS, SRV, TXT, CAA, SSHFP, TLSA, SMIMEA, DS, HTTPS, SVCB, OPENPGPKEY, CERT, URI, RP, LOC`
* `content` - (Required)[string] The content of the DNS Record.
* `ttl` - (Optional)[int] Time to live for the DNS Record.
* `priority` - (Optional)[int] The priority for the DNS Record.
* `enabled` - (Optional)[bool] Indicates if the DNS Record is active or not.
* `zone_id` - (Required)[string] The DNS Zone ID in which the DNS Record will be created.

## Import

In order to import a DNS Record, you can define an empty DNS Record resource in the plan:
```hcl
resource "ionoscloud_dns_record" "example" {
  
}
```
The resource can be imported using the `zone_id` and the `record_id`, for example:

```shell
terraform import ionoscloud_dns_record.example {zone_id}/{record_id}
```