---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_dns_zone_dnssec"
sidebar_current: "docs-resource-dns_zone_dnssec"
description: |-
  Enables DNSSEC for a DNS Zone and exposes the DS record.
---

# ionoscloud_dns_zone_dnssec

Enables **DNSSEC** for a DNS Zone by managing its signing key. The values needed to publish the DS record at your registrar are exported.

> ⚠️  Only tokens are accepted for authorization in the **ionoscloud_dns_zone_dnssec** resource. Please ensure you are using tokens as other methods will not be valid.

A zone can have only one DNSSEC key. The API does not allow modifying a key, so changing any argument other than `timeouts` replaces it (see [Changing arguments](#changing-arguments)).

## Example Usage

```hcl
resource "ionoscloud_dns_zone" "example" {
  name = "example.com"
}

resource "ionoscloud_dns_zone_dnssec" "example" {
  zone_id          = ionoscloud_dns_zone.example.id
  ksk_bits         = 2048
  zsk_bits         = 1024
  nsec_mode        = "NSEC3"
  nsec3_iterations = 10
  nsec3_salt_bits  = 64
  validity         = 120
}

# Hand this over to your registrar.
output "ds_record" {
  value = ionoscloud_dns_zone_dnssec.example.ds_record
}
```

## Argument reference

* `zone_id` - (Required)[string] The ID (UUID) of the DNS zone. Changing it recreates the resource.
* `ksk_bits` - (Required)[int] The size in bits of the key signing key. One of `1024`, `2048`, `4096`.
* `zsk_bits` - (Required)[int] The size in bits of the zone signing key. One of `1024`, `2048`, `4096`.
* `nsec_mode` - (Required)[string] The authenticated denial of existence mode. One of `NSEC`, `NSEC3`.
* `nsec3_iterations` - (Optional)[int] The number of NSEC3 iterations. Defaults to `0`.
* `nsec3_salt_bits` - (Optional)[int] The NSEC3 salt length in bits. Defaults to `64`.
* `validity` - (Required)[int] The signature validity in days.

## Attributes reference

* `id` - The ID of the resource, equal to `zone_id`.
* `algorithm` - The signing algorithm of the key. The API only supports `RSASHA256`.
* `ds_record` - The DS record to hand over to the registrar, formatted as `<key_tag> <algorithm_number> <digest_type> <digest>`.
* `key_tag` - The key tag of the zone signing key.
* `algorithm_number` - The IANA DNS Security Algorithm Number (`8` for `RSASHA256`).
* `digest_type` - The IANA DS digest type number (e.g. `2` for SHA-256).
* `digest_algorithm_mnemonic` - The mnemonic of the digest algorithm, e.g. `SHA-256`.
* `digest` - The digest of the DNSKEY.
* `flags` - The flags of the DNSKEY record.
* `public_key` - The Base64 encoded public key of the DNSKEY record.
* `composed_key_data` - The DNSKEY RDATA: flags, protocol, algorithm and public key.

## Value limits

The provider only restricts the values the API defines as enumerations (`ksk_bits`, `zsk_bits`, `nsec_mode`). Numeric limits are enforced by the API and reported when the key is created. At the time of writing they are: `validity` between 90 and 365 days, `nsec3_iterations` between 0 and 50, `nsec3_salt_bits` between 64 and 128 in multiples of 8, and `ksk_bits` greater than or equal to `zsk_bits`.

## Timeouts

* `create` - (Optional) How long to wait for the key to become available. Defaults to 60 minutes.

## Changing arguments

The DNSSEC API only supports creating, reading and deleting a key; there is no update operation. Changing `zone_id`, `ksk_bits`, `zsk_bits`, `nsec_mode`, `nsec3_iterations`, `nsec3_salt_bits` or `validity` therefore plans a replacement: the key is deleted and a new one is created. Only `timeouts` can be changed in place.

Enabling and disabling DNSSEC are acknowledged by the API and carried out asynchronously, one operation after another per zone. The provider does not wait for the removal of a deleted key. If the zone still has operations in progress when the replacement key is requested, the API rejects the request with `409 ... too many operations in progress`; the provider repeats the request until the `create` timeout is reached. Replacing a key is therefore a slow operation that should be rare, and shortly after a replacement the key attributes (`key_tag`, `digest`, `ds_record`, ...) may still describe the previous key until the API has processed the operations. Run `terraform apply -refresh-only` to pick up the current key before handing the DS record to your registrar.

If the key does not become available within the `create` timeout, a warning is shown and the resource is kept in the state with empty key attributes (`ds_record`, `key_tag`, ...). They are populated by a refresh (e.g. `terraform apply -refresh-only`) once the API returns the key; until then each refresh shows a warning. If the key never becomes available, e.g. because the API failed to set it up, recreate it with `terraform apply -replace=ionoscloud_dns_zone_dnssec.<name>`.

A new key has a new `key_tag` and `digest`, so the DS record at your registrar must be updated as well.

## Import

The DNSSEC key of a zone can be imported using the `zone_id`:

```shell
terraform import ionoscloud_dns_zone_dnssec.example zone_id
```

The API does not return `ksk_bits`, `zsk_bits`, `nsec3_iterations`, `nsec3_salt_bits` and `validity`. After an import they are taken from the configuration, so make sure they match the existing key.
