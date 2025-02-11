---
subcategory: "Certificate Manager Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_auto_certificate"
sidebar_current: "docs-resource-auto_certificate"
description: |-
  Creates and manages Certificate Manager AutoCertificate objects.
---

# ionoscloud_auto_certificate

Manages a **CM AutoCertificate**. 

## Example Usage

```hcl
resource "ionoscloud_auto_certificate_provider" "example" {
  name = "Let's Encrypt"
  email = "user@example.com"
  location = "de/fra"
  server = "https://acme-v02.api.letsencrypt.org/directory"
  external_account_binding {
    key_id = "some-key-id"
    key_secret = "secret"
  }
}

resource "ionoscloud_auto_certificate" "example" {
  provider_id = ionoscloud_auto_certificate_provider.example.id
  common_name = "www.example.com"
  location = ionoscloud_auto_certificate_provider.example.location
  key_algorithm = "rsa4096"
  name = "My Auto renewed certificate"
  subject_alternative_names = ["app.example.com"]
}
```

## Argument reference

* `provider_id` - (Required)[string] The certificate provider used to issue the certificates.
* `location` - (Required)[string] The location of the auto-certificate.
* `common_name` - (Required)[string] The common name (DNS) of the certificate to issue. The common name needs to be part of a zone in IONOS Cloud DNS.
* `key_algorithm` - (Required)[string] The key algorithm used to generate the certificate.
* `name` - (Required)[string] A certificate name used for management purposes.
* `subject_alternative_names` - (Optional)[list][string] Optional additional names to be added to the issued certificate. The additional names needs to be part of a zone in IONOS Cloud DNS.
* `last_issued_certificate_id` - (Computed)[string] The ID of the last certificate that was issued.

## Import

The resource can be imported using the `auto_certificate_id` and the `location`, separated by `:`, e.g.

```shell
terraform import ionoscloud_auto_certificate.example location:auto_certificate_id
```
