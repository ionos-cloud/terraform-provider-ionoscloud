---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_dns_zone"
sidebar_current: "docs-resource-dns_zone"
description: |-
  Creates and manages DNS Zone objects.
---

# ionoscloud_dns_zone

Manages a **DNS Zone**.

> ⚠️  Only tokens are accepted for authorization in the **ionoscloud_dns_zone** resource. Please ensure you are using tokens as other methods will not be valid.

## Example Usage

```hcl
resource "ionoscloud_dns_zone" "example" {
  name = "example.com"
  description = "description"
  enabled = false
}
```

## Argument reference

* `name` - (Required)[string] The name of the DNS Zone. This property is immutable.
* `description` - (Optional)[string] The description for the DNS Zone.
* `enabled` - (Optional)[bool] Indicates if the DNS Zone is active or not. Default is `true`.

## Import

In order to import a DNS Zone, you can define an empty DNS Zone resource in the plan:

```hcl
resource "ionoscloud_dns_zone" "example" {
  
}
```

The resource can be imported using the `zone_id`, for example:

```shell
terraform import ionoscloud_dns_zone.examplezone_id
```

In Terraform v1.12.0 and later, the [`import` block](https://developer.hashicorp.com/terraform/language/import) can also be used with the `identity` attribute:

```hcl
import {
  to = ionoscloud_dns_zone.example
  identity = {
    id = "zone_id"
  }
}

resource "ionoscloud_dns_zone" "example" {
  ### Configuration omitted for brevity ###
}
```

### Identity Schema

#### Required

* `id` (String) The UUID of the DNS zone.

## Query (List Resource)

DNS zones can be listed using `terraform query` (requires Terraform 1.14+). List blocks must be placed in a dedicated query file, whose name ends in `.tfquery.hcl` (for example `queries.tfquery.hcl`).

```hcl
list "ionoscloud_dns_zone" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

See the [`ionoscloud_dns_zone` list resource documentation](../list-resources/dns_zone.md) for filters and the full attribute reference.
