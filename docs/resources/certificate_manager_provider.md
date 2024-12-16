---
subcategory: "Certificate Manager Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_auto_certificate_provider"
sidebar_current: "docs-resource-auto_certificate_provider"
description: |-
  Creates and manages Certificate Manager provider objects.
---

# ionoscloud_auto_certificate_provider

Manages a **CM provider**. 

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
```

## Argument reference

* `name` - (Required)[string] The name of the certificate provider.
* `email` - (Required)[string] The email address of the certificate requester.
* `location` - (Optional)[string] The location of the provider.
* `server` - (Required)[string] The URL of the certificate provider.
* `external_account_binding` - (Optional)[list] External account binding details.
  * `key_id` - (Required)[string] The key ID of the external account binding.
  * `key_secret` - (Required)[string] The key secret of the external account binding

## Import

The resource can be imported using the `provider_id` and the `location`, separated by `:`, e.g.

```shell
terraform import ionoscloud_auto_certificate_provider.example location:provider_id
```
