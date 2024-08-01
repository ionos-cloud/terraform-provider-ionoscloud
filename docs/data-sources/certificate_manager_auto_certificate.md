---
subcategory: "Certificate Manager Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_auto_certificate"
sidebar_current: "docs-datasource-auto_certificate"
description: |-
  Get Information on Certificate Manager AutoCertificate
---

# ionoscloud_auto_certificate

The **CM AutoCertificate data source** can be used to search for and return an existing auto-certificate.
You can provide a string for either id or name parameters which will be compared with provisioned auto-certificates.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_auto_certificate" "example" {
  id			= <auto_certificate_id>
  location      = <auto_certificate_location>
}
```

### By Name
```hcl
data "ionoscloud_auto_certificate" "example" {
  name			= "AutoCertificate Name Example"
  location      = <auto_certificate_location>
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required)[string] The location of the auto-certificate.
* `name` - (Optional)[string] Name of an existing auto-certificate that you want to search for.
* `id` - (Optional)[string] ID of the auto-certificate you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `common_name` - [string] The common name (DNS) of the certificate to issue. The common name needs to be part of a zone in IONOS Cloud DNS.
* `key_algorithm` - [string] The key algorithm used to generate the certificate.
* `subject_alternative_names` - [list][string] Optional additional names to be added to the issued certificate. The additional names needs to be part of a zone in IONOS Cloud DNS.
* `last_issued_certificate_id` - [string] The ID of the last certificate that was issued.