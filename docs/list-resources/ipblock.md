---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ipblock"
description: |-
  Lists IONOS CLOUD IP Blocks.
---

# List Resource: ionoscloud_ipblock

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists **IP Blocks** on IONOS CLOUD. IP Blocks contain reserved public IP addresses that can be assigned to servers or other resources.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated query file, whose name ends in `.tfquery.hcl` (for example `queries.tfquery.hcl`), separate from your main Terraform configuration.

### List all IP Blocks

```hcl
list "ionoscloud_ipblock" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter IP Blocks by location

```hcl
list "ionoscloud_ipblock" "de_fra" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [{
      field_name  = "location"
      field_value = "de/fra"
    }]
  }
}
```

### Filter IP Blocks by name and location

```hcl
list "ionoscloud_ipblock" "web" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "name",     field_value = "web" },
      { field_name = "location", field_value = "de/fra" },
    ]
  }
}
```

### Generate resource configuration from existing IP Blocks

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_ipblock` resource blocks for the IP Blocks the query returns:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_ipblock` resource block for each discovered IP Block into `imported.tf`, which can then be used directly in your configuration.

> **Note:** The IP Blocks are read with a single Cloud API request, which asks for up to **1000**
> items — the same limit the `ionoscloud_ipblock` data source requests. A contract holding more
> IP Blocks than that limit is truncated silently — no error is raised and the missing IP Blocks
> simply do not appear. Check that the number of generated resource blocks matches the number of
> IP Blocks you expect before treating `imported.tf` as complete. `terraform query` also stops at
> the `list` block's own `limit`, which is 100 unless set: set `limit = 1000` in the `list` block
> to receive more than 100 of them.

Terraform names each generated resource after the `list` block label plus an index — a `list "ionoscloud_ipblock" "smoke"` block produces `ionoscloud_ipblock.smoke_0`, `smoke_1`, and so on.

### ⚠️ Do not reuse a `list` block label across separate imports

Because the generated names are derived from the `list` block label, running a second query with the **same** label produces the **same** resource addresses. If a previous address is still in state, the generated configuration is silently applied to the IP Block already recorded at that address — not to the one you just queried.

This happens because an `import` block is idempotent: Terraform skips it when the target address is already in state, so the identity in the generated `import` block is never consulted. Deleting the generated `.tf` file does **not** remove the state entry.

The consequences are not limited to a harmless diff. `location` and `size` are force-new attributes, so if the two IP Blocks differ in either, the plan **destroys the IP Block already in state** and creates a replacement — the IP Block you meant to import is never touched:

```hcl
# generated for a 2-address IP Block in de/fra, but smoke_0 in state points at a 1-address one in de/txl
resource "ionoscloud_ipblock" "smoke_0" {
  location = "de/fra"   # forces replacement of the de/txl IP Block
  size     = 2          # forces replacement as well
  name     = "web"
}
```

To avoid this:

- Use a distinct `list` block label for each query you intend to import from, or
- Remove the stale address before regenerating:

```shell
terraform state show ionoscloud_ipblock.smoke_0   # confirm the address holds the stale import
terraform state rm ionoscloud_ipblock.smoke_0
```

This only removes it from state; the IP Block itself is left in place.

Always read the plan before applying. A clean import reports:

```
Plan: 1 to import, 0 to add, 0 to change, 0 to destroy.
```

Anything reporting changes, and especially `must be replaced` with `location = "..." -> "..." # forces replacement` or `size = ... -> ... # forces replacement`, means the address is bound to a different IP Block — stop and clear the state entry first.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `location`.
  - `field_value` - (Required) The exact value to match against.

> **Note:** The IP Blocks of all locations are read with that one request, so filtering by `location` does not reduce the number of API calls; it only narrows the results. Filtering happens after the response is read, so it cannot recover an IP Block left out by the limit described above.

> **Note:** `field_value` has to match exactly, and the match is case-sensitive. An IP Block without a name is matched by an empty `name` value. A `field_name` other than `name` or `location` is rejected when the configuration is validated.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute  | Description                                                                                                        |
|------------|--------------------------------------------------------------------------------------------------------------------|
| `id`       | The UUID of the IP Block.                                                                                          |
| `location` | The regional location of the IP Block (e.g. `de/fra`). Only needed when the Cloud API endpoint is overridden per location. |

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_ipblock` resource schema:

- `id` - The UUID of the IP Block.
- `name` - The name of the IP Block.
- `location` - The regional location of the IP Block (e.g. `de/fra`).
- `size` - The number of IP addresses reserved in the IP Block.
- `ips` - The list of IP addresses associated with the IP Block.
- `ip_consumers` - Read-only. Lists consumption details of the individual IP addresses:
  - `ip` - The IP address.
  - `mac` - The MAC address.
  - `nic_id` - The ID of the NIC.
  - `server_id` - The ID of the server.
  - `server_name` - The name of the server.
  - `datacenter_id` - The ID of the datacenter.
  - `datacenter_name` - The name of the datacenter.
  - `k8s_nodepool_uuid` - The UUID of the Kubernetes node pool.
  - `k8s_cluster_uuid` - The UUID of the Kubernetes cluster.
- `timeouts` - Always null; timeouts are configuration-only and are not returned by a query.
