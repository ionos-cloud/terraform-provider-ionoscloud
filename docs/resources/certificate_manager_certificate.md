---
subcategory: "Certificate Manager Service"
layout: "ionoscloud"
page_title: "IonosCloud: certificate"
sidebar_current: "docs-resource-certificate"
description: |-
  Creates and manages a certificate.
---

# ionoscloud_certificate

Manages a **Certificate** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_certificate" "cert" {
  name = "add_name_here"
  certificate = "tour_certificate"
  certificate_chain = "your_certificate_chain"
  private_key = "your_private_key"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required)[string] The certificate name
* `certificate` - (Required)[string] The certificate body. Pem encoded. Immutable.
* `private_key` - (Required)[string] The certificate private key. Immutable. Sensitive.
* `certificate_chain` - (Optional)[string] The certificate chain. Pem encoded. Immutable.

## Import

Resource certificate can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_certificate.mycert certificate uuid
```
