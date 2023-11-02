---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: certificate"
sidebar_current: "docs-datasource-certificate"
description: |-
  Get Information on a certificate
---

# ionoscloud_certificate

The **Certificate data source** can be used to search for and return an existing certificate.
You can provide a string for either id or name parameters which will be compared with provisioned certificates.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_certificate" "example" {
  id			= <certificate_id>
}
```

### By Name
```hcl
data "ionoscloud_certificate" "example" {
  name			= "Certificate Name Example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of an existing certificate that you want to search for.
* `id` - (Optional) ID of the certificate you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the certificate.
* `name` - The name of the certificate.
* `certificate` - Certificate body. 
* `certificate_chain` - Certificate chain.