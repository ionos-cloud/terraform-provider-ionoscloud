---
subcategory: "Certificate Manager Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_auto_certificate_provider"
sidebar_current: "docs-datasource-auto_certificate_provider"
description: |-
  Get Information on Certificate Manager Provider
---

# ionoscloud_auto_certificate_provider

The **CM Provider data source** can be used to search for and return an existing certificate manager provider.
You can provide a string for either id or name parameters which will be compared with provisioned providers.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_auto_certificate_provider" "example" {
  id			= "provider_id"
  location      = "provider_location"
}
```

### By Name
```hcl
data "ionoscloud_auto_certificate_provider" "example" {
  name			= "Provider Name Example"
  location      = "provider_location"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required)[string] The location of the provider. Available locations: `de/fra`, `de/fra/2`
* `name` - (Optional)[string] Name of an existing provider that you want to search for.
* `id` - (Optional)[string] ID of the provider you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `email` - [string] The email address of the certificate requester.
* `server` - [string] The URL of the certificate provider.
* `external_account_binding` - [list]
  * `key_id` - [string] The key ID of the external account binding.
  * `key_secret` - [string] The key secret of the external account binding