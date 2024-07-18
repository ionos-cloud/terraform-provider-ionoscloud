---
subcategory: "Cdn"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_cdn_distribution"
sidebar_current: "docs-ionoscloud_cdn_distribution"
description: |-
  Get information on an CDN Distribution
---

# ionoscloud_cdn_distribution

The Distribution data source can be used to search for and return an existing Distributions.
You can provide a string for the domain parameter which will be compared with provisioned Distributions.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique domains.

## Example Usage

### By Id
```hcl
data "ionoscloud_cdn_distribution" "example" {
  id = <distr_id>
}
```

### By Domain
```hcl
data "ionoscloud_cdn_distribution" "example" {
  domain = "example.com"
}
```

### By Domain with Partial Match
```hcl
data "ionoscloud_cdn_distribution" "example" {
  domain    		= "example"
  partial_match = true
}
```

## Argument Reference

* `id` - (Optional) ID of the distribution you want to search for.
* `domain` - (Optional) Domain of an existing distribution that you want to search for. Search by domain is case-insensitive. The whole resource domain is required if `partial_match` parameter is not set to true.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using domain argument. Default value is false.

Either `domain` or `id` must be provided. If none, or both of `domain` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

- `domain` - The domain of the distribution.
- `certificate_id` - The ID of the certificate to use for the distribution. You can create certificates with the [certificate](certificate.md) resource.
- `routing_rules` - The routing rules for the distribution.
    - `scheme` - The scheme of the routing rule.
    - `prefix` - The prefix of the routing rule.
    - `upstream` - A map of properties for the rule
        * `host` - The upstream host that handles the requests if not already cached. This host will be protected by the WAF if the option is enabled.
        * `caching` - Enable or disable caching. If enabled, the CDN will cache the responses from the upstream host. Subsequent requests for the same resource will be served from the cache.
        * `waf` - Enable or disable WAF to protect the upstream host.
        * `rate_limit_class` - Rate limit class that will be applied to limit the number of incoming requests per IP.
        * `geo_restrictions` - A map of geo_restrictions
            * `allow_list` - List of allowed countries
            * `block_list` - List of blocked countries
