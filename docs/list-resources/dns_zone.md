---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IONOS CLOUD: dns_zone"
description: |-
  Lists IONOS CLOUD DNS Zones.
---

# List Resource: ionoscloud_dns_zone

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists [DNS Zones](https://docs.ionos.com/cloud/network-services/cloud-dns/overview) on IONOS CLOUD.

> ⚠️  Only tokens are accepted for authorization in the **ionoscloud_dns_zone** list resource. Please ensure you are using tokens as other methods will not be valid.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated query file, whose name ends in `.tfquery.hcl` (for example `queries.tfquery.hcl`), separate from your main Terraform configuration.

### List all zones

```hcl
list "ionoscloud_dns_zone" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter zones by name

```hcl
list "ionoscloud_dns_zone" "example_com" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [{
      field_name  = "name"
      field_value = "example.com"
    }]
  }
}
```

### Filter zones by name and description

```hcl
list "ionoscloud_dns_zone" "prod" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "name",        field_value = "example.com" },
      { field_name = "description", field_value = "production zone" },
    ]
  }
}
```

### Generate resource configuration from existing zones

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_dns_zone` resource blocks for the zones the query returns:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_dns_zone` resource block for each discovered zone into `imported.tf`, which can then be used directly in your configuration.

> **Note:** The zones are read with a single Cloud DNS API request. The provider sends no
> `limit` with it, so the Cloud DNS API applies its default page size of **100** zones. A
> contract holding more zones than that limit is truncated silently — no error is raised and
> the missing zones simply do not appear. Check that the number of generated resource blocks
> matches the number of zones you expect before treating `imported.tf` as complete.

Terraform names each generated resource after the `list` block label plus an index — a `list "ionoscloud_dns_zone" "smoke"` block produces `ionoscloud_dns_zone.smoke_0`, `smoke_1`, and so on.

### ⚠️ Do not reuse a `list` block label across separate imports

Because the generated names are derived from the `list` block label, running a second query with the **same** label produces the **same** resource addresses. If a previous address is still in state, the generated configuration is silently applied to the zone already recorded at that address — not to the one you just queried.

This happens because an `import` block is idempotent: Terraform skips it when the target address is already in state, so the identity in the generated `import` block is never consulted. Deleting the generated `.tf` file does **not** remove the state entry.

The consequences are not limited to a harmless diff. `name` is a force-new attribute, so if the two zones have different names, the plan **destroys the zone already in state** and creates a replacement — the zone you meant to import is never touched, and every DNS record that belonged to the destroyed zone goes with it:

```hcl
# generated for example.com, but smoke_0 in state points at staging.example.com
resource "ionoscloud_dns_zone" "smoke_0" {
  name        = "example.com"   # forces replacement of staging.example.com
  description = "production zone"
}
```

To avoid this:

- Use a distinct `list` block label for each query you intend to import from, or
- Confirm what the stale address is bound to and remove it before regenerating:

```shell
terraform state show ionoscloud_dns_zone.smoke_0
terraform state rm ionoscloud_dns_zone.smoke_0
```

`terraform state rm` stops Terraform managing that zone until it is imported again; the zone and its records are left in place.

Always read the plan before applying. A clean import reports:

```
Plan: 1 to import, 0 to add, 0 to change, 0 to destroy.
```

Anything reporting changes, and especially `must be replaced` with `name = "..." -> "..." # forces replacement`, means the address is bound to a different zone — stop and clear the state entry first.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `description`.
  - `field_value` - (Required) The exact value to match against.

> **Note:** Filtering happens after the response described above is read, so it cannot recover a zone left out by the page limit. Filter values are matched exactly and are case-sensitive; an unknown `field_name` is rejected at plan time, but a mis-cased `field_value` simply matches nothing.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute | Description             |
|-----------|-------------------------|
| `id`      | The UUID of the DNS zone. |

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_dns_zone` resource schema:

- `id` - The UUID of the DNS zone.
- `name` - The name of the DNS zone.
- `description` - The description of the DNS zone.
- `enabled` - Indicates if the DNS zone is active or not.
- `nameservers` - The list of name servers assigned to the DNS zone.
- `timeouts` - Always null; timeouts are configuration-only and are not returned by a query.
