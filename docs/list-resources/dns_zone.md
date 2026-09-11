---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IONOS CLOUD: dns_zone"
description: |-
  Lists IONOS CLOUD DNS Zones.
---

# List Resource: ionoscloud_dns_zone

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists **DNS Zones** on IONOS CLOUD — see the [`ionoscloud_dns_zone` resource](../resources/dns_zone.md) for how to manage one.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated query file, whose name ends in `.tfquery.hcl` (for example `queries.tfquery.hcl`), separate from your main Terraform configuration.

### List all DNS zones

```hcl
list "ionoscloud_dns_zone" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter DNS zones by name

```hcl
list "ionoscloud_dns_zone" "by_name" {
  provider         = ionoscloud
  include_resource = true

  config {
    filters = [
      {
        field_name  = "name"
        field_value = "example.com"
      }
    ]
  }
}
```

### Filter DNS zones by name and description

```hcl
list "ionoscloud_dns_zone" "production" {
  provider         = ionoscloud
  include_resource = true

  config {
    filters = [
      { field_name = "name",        field_value = "example.com" },
      { field_name = "description", field_value = "Production zone" },
    ]
  }
}
```

### Generate resource configuration from existing DNS zones

```shell
terraform query -generate-config-out=imported.tf
```

> **Note:** The DNS zones are read with a single `GET /zones` request. Unlike the Cloud API
> endpoints, the IONOS DNS SDK sends no `limit` at all when none is asked for, so the page size is
> whatever the server applies by default — this provider does not set one, and the value the server
> uses is not verified here. A contract holding more DNS zones than one response returns is
> truncated silently — no error is raised and the missing zones simply do not appear. Check that the
> number of generated resource blocks matches the number of zones you expect before treating
> `imported.tf` as complete.

### ⚠️ Do not reuse a `list` block label across separate imports

`terraform query -generate-config-out` names each generated resource `<list block label>_<index>`,
so two separate query runs that reuse the same label generate the same addresses. If the first run's
`import` blocks are already in state, a second run silently maps those addresses onto **different**
zones: an `import` block is skipped when its address is already in state, so the identity in the
generated file is never consulted.

```hcl
# first run, imported.tf
import {
  to = ionoscloud_dns_zone.zones_0
  identity = { id = "..." }   # zone A
}

# second run, same label
import {
  to = ionoscloud_dns_zone.zones_0   # already in state, points at zone A
  identity = { id = "..." }          # zone B - never read
}
```

`name` is ForceNew on `ionoscloud_dns_zone`, so a mismatch of this kind is not a harmless rename: it
plans a destroy and recreate of a zone you did not intend to touch. Use a distinct label per query
run, and read the plan before applying.

## Argument Reference

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `description`.
  - `field_value` - (Required) The exact value to match against. Matching is case-sensitive and compares against the value the API returned.

> **Note:** Filtering happens after the response is read, so it cannot recover a DNS zone left out by the page size described above.

> **Note:** `description` is optional on `ionoscloud_dns_zone`. A zone created without one is still returned, and matches only a `description` filter whose `field_value` is the empty string.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute | Description              |
|-----------|--------------------------|
| `id`      | The UUID of the DNS zone. |

A DNS zone has no location, so the identity is a lone `id` — the same plain UUID the resource's
legacy import string accepts.

## Attributes Reference

When `include_resource = true`, each result carries the full `ionoscloud_dns_zone` resource state:

- `id` - The UUID of the DNS zone.
- `name` - The name of the DNS zone (a domain name).
- `description` - The description of the DNS zone. Optional on the resource, so it can be null.
- `enabled` - Whether the DNS zone is active.
- `nameservers` - The list of name servers assigned to the zone.
