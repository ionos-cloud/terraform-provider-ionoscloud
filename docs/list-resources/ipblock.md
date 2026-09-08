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

### List all IP blocks

```hcl
list "ionoscloud_ipblock" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter IP blocks by location

```hcl
list "ionoscloud_ipblock" "us_las" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [{
      field_name  = "location"
      field_value = "us/las"
    }]
  }
}
```

### Filter IP blocks by name and location

```hcl
list "ionoscloud_ipblock" "prod" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "name",     field_value = "production" },
      { field_name = "location", field_value = "us/las" },
    ]
  }
}
```

### Generate resource configuration from existing IP blocks

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_ipblock` resource blocks for the IP blocks the query returns:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_ipblock` resource block for each discovered IP block into `imported.tf`, which can then be used directly in your configuration.

> **Note:** The IP blocks are read with a single Cloud API request, which asks for up to
> 1000 items — the same limit the `ionoscloud_ipblock` data source requests, and ten times
> the 100 the IONOS SDK falls back to when no limit is asked for. Whether
> the server honours the requested limit in full is not verified. A contract holding more
> IP blocks than one response returns is truncated silently — no error is raised and the
> missing IP blocks simply do not appear. Check that the number of generated resource
> blocks matches the number of IP blocks you expect before treating `imported.tf` as
> complete.

Terraform names each generated resource after the `list` block label plus an index — a `list "ionoscloud_ipblock" "smoke"` block produces `ionoscloud_ipblock.smoke_0`, `smoke_1`, and so on.

### ⚠️ Do not reuse a `list` block label across separate imports

Because the generated names are derived from the `list` block label, running a second query with the **same** label produces the **same** resource addresses. If a previous address is still in state, the generated configuration is silently applied to the IP block already recorded at that address — not to the one you just queried.

This happens because an `import` block is idempotent: Terraform skips it when the target address is already in state, so the identity in the generated `import` block is never consulted. Deleting the generated `.tf` file does **not** remove the state entry.

The consequences are not limited to a harmless diff. Both `location` and `size` are force-new attributes, so if the two IP blocks differ in either, the plan **destroys the IP block already in state** and creates a replacement — the IP block you meant to import is never touched. Because the replacement reserves fresh addresses, the public IPs of the destroyed block are lost:

```hcl
# generated for an IP block in de/fra, but smoke_0 in state points at one in us/las
resource "ionoscloud_ipblock" "smoke_0" {
  location = "de/fra"          # forces replacement of the us/las IP block
  size     = 1
  name     = "IP Block Example"
}
```

To avoid this:

- Use a distinct `list` block label for each query you intend to import from, or
- Clear the address deliberately: run `terraform state show ionoscloud_ipblock.smoke_0` first
  and confirm it is the leftover import and not an IP block you still manage, then
  `terraform state rm ionoscloud_ipblock.smoke_0`. Removing the address does not release the
  IP block — it stops Terraform managing it, and it has to be re-imported to come back under
  management.

Always read the plan before applying. A clean import reports:

```
Plan: 1 to import, 0 to add, 0 to change, 0 to destroy.
```

Anything reporting changes, and especially `must be replaced` with a `# forces replacement` marker on `location` or `size`, means the address is bound to a different IP block — stop and clear the state entry first.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `location`.
  - `field_value` - (Required) The exact value to match against.

> **Note:** The Cloud API returns IP blocks from every location in a single response, so filtering by `location` does not reduce the number of API calls; it only narrows the results. Filtering happens after that response is read, so it cannot recover an IP block left out by the page limit described above.

> **Note:** `name` is optional on `ionoscloud_ipblock`. Filtering by `name` therefore never matches an unnamed IP block, and unnamed blocks are listed under their UUID instead of a name.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute  | Description                                                                                                             |
|------------|-------------------------------------------------------------------------------------------------------------------------|
| `id`       | The UUID of the IP block.                                                                                               |
| `location` | The location the IP block lives in (e.g. `us/las`). Only needed when the Cloud API endpoint is overridden per location.  |

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_ipblock` resource schema:

- `id` - The UUID of the IP block.
- `name` - The name of the IP block. Optional, so it may be null.
- `location` - The regional location the IP block is reserved in (e.g. `us/las`).
- `size` - The number of IP addresses reserved for the block.
- `ips` - The list of IP addresses associated with the block.
- `ip_consumers` - Consumption detail for each IP of the block, reporting what it is currently attached to. Empty when none of the addresses are in use:
  - `ip` - The IP address.
  - `mac` - The MAC address of the NIC the IP is attached to.
  - `nic_id` - The UUID of the NIC the IP is attached to.
  - `server_id` - The UUID of the server the IP is attached to.
  - `server_name` - The name of that server.
  - `datacenter_id` - The UUID of the datacenter the IP is used in.
  - `datacenter_name` - The name of that datacenter.
  - `k8s_nodepool_uuid` - The UUID of the Kubernetes node pool using the IP, if any.
  - `k8s_cluster_uuid` - The UUID of the Kubernetes cluster using the IP, if any.
- `timeouts` - Always null; timeouts are configuration-only and are not returned by a query.
