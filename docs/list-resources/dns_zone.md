---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IONOS CLOUD: dns_zone"
description: |-
  Lists IONOS CLOUD DNS Zones.
---

# List Resource: ionoscloud_dns_zone

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists **DNS Zones** on IONOS CLOUD.

> ⚠️  As for the **ionoscloud_dns_zone** resource, only tokens are accepted for authorization. Please ensure you are using tokens as other methods will not be valid.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated query file, whose name ends in `.tfquery.hcl` (for example `queries.tfquery.hcl`), separate from your main Terraform configuration.

### List all DNS Zones

```hcl
list "ionoscloud_dns_zone" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter DNS Zones by name

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

### Filter DNS Zones by name and description

```hcl
list "ionoscloud_dns_zone" "production" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "name",        field_value = "example.com" },
      { field_name = "description", field_value = "production" },
    ]
  }
}
```

### Generate resource configuration from existing DNS Zones

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_dns_zone` resource blocks for the DNS Zones the query returns:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_dns_zone` resource block for each discovered DNS Zone into `imported.tf`, which can then be used directly in your configuration.

> **Note:** The DNS Zones are read with a single DNS API request. The provider sends no `limit`
> with it, so the DNS API applies its default page size of **100** DNS Zones. A contract holding
> more DNS Zones than that limit is truncated silently — no error is raised and the missing DNS
> Zones simply do not appear. Check that the number of generated resource blocks matches the
> number of DNS Zones you expect before treating `imported.tf` as complete.

Terraform names each generated resource after the `list` block label plus an index — a `list "ionoscloud_dns_zone" "smoke"` block produces `ionoscloud_dns_zone.smoke_0`, `smoke_1`, and so on.

### ⚠️ Do not reuse a `list` block label across separate imports

Because the generated names are derived from the `list` block label, running a second query with the **same** label produces the **same** resource addresses. If a previous address is still in state, the generated configuration is silently applied to the DNS Zone already recorded at that address — not to the one you just queried.

This happens because an `import` block is idempotent: Terraform skips it when the target address is already in state, so the identity in the generated `import` block is never consulted. Deleting the generated `.tf` file does **not** remove the state entry.

The consequences are not limited to a harmless diff. `name` is a force-new attribute, so if the two DNS Zones have different names, the plan **destroys the DNS Zone already in state** and creates a replacement — the DNS Zone you meant to import is never touched. Deleting a DNS Zone also deletes all of the records it contains:

```hcl
# generated for example.org, but smoke_0 in state points at example.com
resource "ionoscloud_dns_zone" "smoke_0" {
  name = "example.org"   # forces replacement of the example.com DNS Zone
}
```

To avoid this:

- Use a distinct `list` block label for each query you intend to import from, or
- Remove the stale address before regenerating:

```shell
terraform state show ionoscloud_dns_zone.smoke_0   # confirm the address holds the stale import
terraform state rm ionoscloud_dns_zone.smoke_0
```

This only removes it from state; the DNS Zone itself is left in place.

Always read the plan before applying. A clean import reports:

```
Plan: 1 to import, 0 to add, 0 to change, 0 to destroy.
```

Anything reporting changes, and especially `must be replaced` with `name = "..." -> "..." # forces replacement`, means the address is bound to a different DNS Zone — stop and clear the state entry first.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `description`.
  - `field_value` - (Required) The exact value to match against.

> **Note:** Filtering happens after the response is read, so it cannot recover a DNS Zone left out by the page limit described above.

> **Note:** `field_value` has to match exactly, and the match is case-sensitive. A DNS Zone without a description is matched by an empty `description` value. A `field_name` other than `name` or `description` is rejected when the configuration is validated.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute | Description               |
|-----------|---------------------------|
| `id`      | The UUID of the DNS Zone. |

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_dns_zone` resource schema:

- `id` - The UUID of the DNS Zone.
- `name` - The name of the DNS Zone.
- `description` - The description for the DNS Zone.
- `enabled` - Indicates if the DNS Zone is active or not.
- `nameservers` - A list of available name servers.
- `timeouts` - Always null; timeouts are configuration-only and are not returned by a query.
