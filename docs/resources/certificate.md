---
subcategory: "Compute Engine"
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
  certificate = "${file("path_to_cert")}"
  certificate_chain = "${file("path_to_cert_chain")}"
  private_key = "${file("path_to_private_key")}"
}
```

**NOTE**: You can also provide the values as multiline strings, as seen below:

```hcl
resource "ionoscloud_certificate" "cert" {
  name = "add_name_here"
  certificate = <<EOT
-----BEGIN CERTIFICATE-----
cert_body_here
-----END CERTIFICATE-----
EOT
  certificate_chain = "${file("path_to_cert_chain")}"
  private_key = "${file("path_to_private_key")}"
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
terraform import ionoscloud_certificate.mycert {certificate uuid}
```
