---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IONOS CLOUD : ionoscloud_dns_zone_dnssec"
sidebar_current: "docs-dns_zone_dnssec"
description: |-
  Get the DNSSEC key information of a DNS Zone.
---

# ionoscloud_dns_zone_dnssec

Reads the DNSSEC signing key of a DNS Zone, including the values needed for the DS record. Fails if DNSSEC is not enabled for the zone.

> ⚠️  Only tokens are accepted for authorization in the **ionoscloud_dns_zone_dnssec** data source. Please ensure you are using tokens as other methods will not be valid.

## Example Usage

```hcl
data "ionoscloud_dns_zone_dnssec" "example" {
  zone_id = "zone_id"
}

output "ds_record" {
  value = data.ionoscloud_dns_zone_dnssec.example.ds_record
}
```

## Argument reference

* `zone_id` - (Required)[string] The ID (UUID) of the DNS zone.

## Attributes reference

* `id` - The ID of the data source, equal to `zone_id`.
* `algorithm` - The signing algorithm.
* `nsec_mode` - The authenticated denial of existence mode.
* `ds_record` - The DS record to hand over to the registrar, formatted as `<key_tag> <algorithm_number> <digest_type> <digest>`.
* `key_tag` - The key tag of the zone signing key.
* `algorithm_number` - The IANA DNS Security Algorithm Number.
* `digest_type` - The IANA DS digest type number.
* `digest_algorithm_mnemonic` - The mnemonic of the digest algorithm.
* `digest` - The digest of the DNSKEY.
* `flags` - The flags of the DNSKEY record.
* `public_key` - The Base64 encoded public key of the DNSKEY record.
* `composed_key_data` - The DNSKEY RDATA: flags, protocol, algorithm and public key.
