---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_dns_zone"
sidebar_current: "docs-resource-dns_zone"
description: |-
  Creates and manages DNS Zone objects.
---

# ionoscloud_dns_zone

Manages a **DNS Zone**.

## Example Usage

```hcl
resource "ionoscloud_dns_zone" "example" {
  name = "example.com"
  description = "description"
  enabled = false
}
```

## Argument reference

* `name` - (Required)[string] The name of the DNS Zone.
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
terraform import ionoscloud_dns_zone.example {zone_id}
```