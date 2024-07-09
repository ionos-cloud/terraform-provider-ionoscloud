---
subcategory: "Cdn"
layout: "ionoscloud"
page_title: "IonosCloud: distribution"
sidebar_current: "docs-resource-distribution"
description: |-
  Creates and manages IonosCloud CDN Distributions.
---

# ionoscloud_distribution

Manages a **Distribution** on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_distribution" "example" {
  domain         = "example.com"
  certificate_id = ionoscloud_certificate.cert.id
  routing_rules {
    scheme = "https"
    prefix = "/api"
    upstream {
      host                = "server.example.com"
      caching             = true
      waf                 = true
      rate_limit_class    = "none"
      geo_restrictions {
        allow_list = [ "CN", "RU"]
      }
    }
  }
  routing_rules {
    scheme = "http/https"
    prefix = "/api2"
    upstream {
      host                = "server2.example.com"
      caching             = false
      waf                 = false
      rate_limit_class    = "R10000"
      geo_restrictions {
        block_list = [ "CN", "RU"]
      }
    }
  }
}

#optionally you can add a certificate to the distribution
resource "ionoscloud_certificate" "cert" {
  name = "add_name_here"
  certificate = "${file("path_to_cert")}"
  certificate_chain = "${file("path_to_cert_chain")}"
  private_key = "${file("path_to_private_key")}"
}
```

## Argument Reference

The following arguments are supported:

- `domain` - (Required)[string] The domain of the distribution.
- `certificate_id` - (Required)[string] The ID of the certificate to use for the distribution. You can create certificates with the [certificate](certificate.md) resource.
- `routing_rules` - (Required)[list] The routing rules for the distribution.
    - `scheme` - (Required)[string] The scheme of the routing rule.
    - `prefix` - (Required)[string] The prefix of the routing rule.
    - `upstream` - (Required)[map] - A map of properties for the rule
        * `host` - (Required)[string] The upstream host that handles the requests if not already cached. This host will be protected by the WAF if the option is enabled.
        * `caching` - (Required)[bool] Enable or disable caching. If enabled, the CDN will cache the responses from the upstream host. Subsequent requests for the same resource will be served from the cache.
        * `waf` - (Required)[bool] Enable or disable WAF to protect the upstream host.
        * `rate_limit_class` - (Required)[string] Rate limit class that will be applied to limit the number of incoming requests per IP.
        * `geo_restrictions` - (Optional)[map] - A map of geo_restrictions
            * `allow_list` - (Optional)[string] List of allowed countries
            * `block_list` - (Optional)[string] List of blocked countries

## Import

Resource Distribution can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_distribution.myDistribution {distribution uuid}
```
