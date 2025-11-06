---
subcategory: "CDN"
layout: "ionoscloud"
page_title: "IonosCloud: cdn_distribution"
sidebar_current: "docs-resource-cdn-distribution"
description: |-
  Creates and manages IonosCloud CDN Distributions.
---

# ionoscloud_cdn_distribution

Manages a [CDN Distribution](https://docs.ionos.com/cloud/network-services/cdn/overview#how-does-cdn-work) on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_cdn_distribution" "example" {
  domain         = "example.com"
  certificate_id = ionoscloud_certificate.cert.id
  routing_rules {
    scheme = "https"
    prefix = "/api"
    upstream {
      host                = "server.example.com"
      caching             = true
      waf                 = true
      sni_mode            = "distribution"
      rate_limit_class    = "R500"
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
      sni_mode            = "origin"
      rate_limit_class    = "R10"
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
- `certificate_id` - (Required)[string] The ID of the certificate to use for the distribution. You can create certificates with the [certificate](certificate_manager_certificate.md) resource.
- `routing_rules` - (Required)[list] The routing rules for the distribution.
    - `scheme` - (Required)[string] The scheme of the routing rule.
    - `prefix` - (Required)[string] The prefix of the routing rule.
    - `upstream` - (Required)[map] - A map of properties for the rule
        * `host` - (Required)[string] The upstream host that handles the requests if not already cached. This host will be protected by the WAF if the option is enabled.
        * `caching` - (Required)[bool] Enable or disable caching. If enabled, the CDN will cache the responses from the upstream host. Subsequent requests for the same resource will be served from the cache.
        * `waf` - (Required)[bool] Enable or disable WAF to protect the upstream host.
        * `sni_mode` - (Required)[string] The SNI (Server Name Indication) mode of the upstream. It supports two modes: 1) `distribution`: for outgoing connections to the upstream host, the CDN requires the upstream host to present a valid certificate that matches the configured domain of the CDN distribution; 2) `origin`: for outgoing connections to the upstream host, the CDN requires the upstream host to present a valid certificate that matches the configured upstream/origin hostname.
        * `rate_limit_class` - (Required)[string] Rate limit class that will be applied to limit the number of incoming requests per IP.
        * `geo_restrictions` - (Optional)[map] - A map of geo_restrictions
            * `allow_list` - (Optional)[string] List of allowed countries
            * `block_list` - (Optional)[string] List of blocked countries
## Attributes Reference

- `public_endpoint_v4` - IP of the distribution, it has to be included on the domain DNS Zone as A record.
- `public_endpoint_v6` - IP of the distribution, it has to be included on the domain DNS Zone as AAAA record.
- `resource_urn` - Unique resource indentifier.

## Import

Resource Distribution can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_cdn_distribution.myDistribution distribution uuid
```
