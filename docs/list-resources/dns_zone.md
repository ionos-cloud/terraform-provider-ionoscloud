---
subcategory: "Cloud DNS"
layout: "ionoscloud"
page_title: "IONOS CLOUD: dns_zone"
description: |-
  Lists IONOS CLOUD DNS Zones.
---

# List Resource: ionoscloud_dns_zone

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

⚠️ **Note:** Only tokens are accepted for authorization in the **ionoscloud_dns_zone** list resource. Please ensure you are using tokens as other methods will not be valid.

Lists [DNS Zones](https://docs.ionos.com/cloud/network-services/cloud-dns/overview) on IONOS CLOUD.

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
list "ionoscloud_dns_zone" "example" {
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

### Filter DNS zones by name and description

```hcl
list "ionoscloud_dns_zone" "prod" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "name", field_value = "example.com" },
      { field_name = "description", field_value = "production zone" },
    ]
  }
}
```

### Generate resource configuration from existing DNS zones

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_dns_zone` resource blocks for the zones the query returns:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_dns_zone` resource block for each discovered zone into `imported.tf`, which can then be used directly in your configuration.

> **Note:** The zones are read with a single `GET /zones` request. Unlike the Cloud API endpoints, the IONOS DNS SDK sends no `limit` at all when none is asked for, so the page size is whatever the server applies by default — this provider does not set one, and the value the server uses is not verified here. The vendored SDK's own documentation for that operation states "Default limit is the first 100 items. Use pagination query parameters for listing more items (up to 1000)", which is the only evidence available; it is a doc comment, not generated code, and nothing in this provider confirms it. A contract holding more zones than one response returns is truncated silently — no error is raised and the missing zones simply do not appear. Check that the number of generated resource blocks matches the number of zones you expect before treating `imported.tf` as complete.

Terraform names each generated resource after the `list` block label plus an index — a `list "ionoscloud_dns_zone" "smoke"` block produces `ionoscloud_dns_zone.smoke_0`, `smoke_1`, and so on.

### ⚠️ Do not reuse a `list` block label across separate imports

Because the generated names are derived from the `list` block label, running a second query with the **same** label produces the **same** resource addresses. If a previous address is still in state, the generated configuration is silently applied to the zone already recorded at that address — not to the one you just queried.

This happens because an `import` block is idempotent: Terraform skips it when the target address is already in state, so the identity in the generated `import` block is never consulted. Deleting the generated `.tf` file does **not** remove the state entry.

The consequences are not limited to a harmless diff. `name` is a force-new attribute on `ionoscloud_dns_zone`, so if the two zones differ in it, the plan **destroys the zone already in state** and creates a replacement — the zone you meant to import is never touched, and destroying a zone takes every DNS record in it with it:

```hcl
# generated for a zone named "b.example.com", but smoke_0 in state points at "a.example.com"
resource "ionoscloud_dns_zone" "smoke_0" {
  name        = "b.example.com"     # forces replacement of the "a.example.com" zone
  description = "production zone"
}
```

To avoid this:

- Use a distinct `list` block label for each query you intend to import from, or
- Clear the address deliberately: run `terraform state show ionoscloud_dns_zone.smoke_0` first
  and confirm it is the leftover import and not a zone you still manage, then
  `terraform state rm ionoscloud_dns_zone.smoke_0`. Removing the address does not delete the
  zone — it stops Terraform managing it, and it has to be re-imported to come back under
  management.

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

> **Note:** `enabled` and `nameservers` are attributes of the resource but are not filterable. Filter values are compared as exact strings, which a boolean and a list of names have nothing to match against.

> **Note:** Filtering happens after the single API response above is read, so it narrows the results but does not reduce the number of API calls, and it cannot recover a zone left out by the page limit described above.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute | Description             |
|-----------|-------------------------|
| `id`      | The UUID of the DNS zone. |

A DNS zone is not addressed by location, so — unlike `ionoscloud_datacenter` and `ionoscloud_ipblock` — its identity is the UUID alone.

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_dns_zone` resource schema:

- `id` - The UUID of the DNS zone.
- `name` - The name of the DNS zone.
- `description` - The description of the DNS zone. Null when the zone has none.
- `enabled` - Whether the zone is active. Null when the API does not report it.
- `nameservers` - The list of name servers assigned to the zone. Empty when the API returns none.
- `timeouts` - Always null; timeouts are configuration-only and are not returned by a query.
