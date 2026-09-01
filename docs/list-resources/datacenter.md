---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IONOS CLOUD: datacenter"
description: |-
  Lists IONOS CLOUD Virtual Data Centers.
---

# List Resource: ionoscloud_datacenter

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists Virtual [Data Centers](https://docs.ionos.com/cloud/set-up-ionos-cloud/get-started/configure-data-center) on IONOS CLOUD.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated query file, whose name ends in `.tfquery.hcl` (for example `queries.tfquery.hcl`), separate from your main Terraform configuration.

### List all datacenters

```hcl
list "ionoscloud_datacenter" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter datacenters by location

```hcl
list "ionoscloud_datacenter" "de_txl" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [{
      field_name  = "location"
      field_value = "de/txl"
    }]
  }
}
```

### Filter datacenters by name and location

```hcl
list "ionoscloud_datacenter" "prod" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "name",     field_value = "production" },
      { field_name = "location", field_value = "de/txl" },
    ]
  }
}
```

### Generate resource configuration from existing datacenters

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_datacenter` resource blocks for all existing datacenters:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_datacenter` resource block for each discovered datacenter into `imported.tf`, which can then be used directly in your configuration.

Terraform names each generated resource after the `list` block label plus an index — a `list "ionoscloud_datacenter" "smoke"` block produces `ionoscloud_datacenter.smoke_0`, `smoke_1`, and so on.

### ⚠️ Do not reuse a `list` block label across separate imports

Because the generated names are derived from the `list` block label, running a second query with the **same** label produces the **same** resource addresses. If a previous address is still in state, the generated configuration is silently applied to the datacenter already recorded at that address — not to the one you just queried.

This happens because an `import` block is idempotent: Terraform skips it when the target address is already in state, so the identity in the generated `import` block is never consulted. Deleting the generated `.tf` file does **not** remove the state entry.

The consequences are not limited to a harmless diff. `location` is a force-new attribute, so if the two datacenters are in different locations, the plan **destroys the datacenter already in state** and creates a replacement — the datacenter you meant to import is never touched:

```hcl
# generated for a datacenter in de/fra, but smoke_0 in state points at one in de/txl
resource "ionoscloud_datacenter" "smoke_0" {
  location = "de/fra"          # forces replacement of the de/txl datacenter
  name     = "docker-machine-data-center"
}
```

To avoid this:

- Use a distinct `list` block label for each query you intend to import from, or
- Remove the stale address before regenerating: `terraform state rm ionoscloud_datacenter.smoke_0`

Always read the plan before applying. A clean import reports:

```
Plan: 1 to import, 0 to add, 0 to change, 0 to destroy.
```

Anything reporting changes, and especially `must be replaced` with `location = "..." -> "..." # forces replacement`, means the address is bound to a different datacenter — stop and clear the state entry first.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `location`.
  - `field_value` - (Required) The exact value to match against.

> **Note:** The Cloud API returns datacenters from every location in a single response, so filtering by `location` does not reduce the number of API calls; it only narrows the results.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute  | Description                                                                                              |
|------------|----------------------------------------------------------------------------------------------------------|
| `id`       | The UUID of the datacenter.                                                                              |
| `location` | The location the datacenter lives in (e.g. `de/txl`). Only needed when the Cloud API endpoint is overridden per location. |

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_datacenter` resource schema:

- `id` - The UUID of the datacenter.
- `name` - The name of the datacenter.
- `description` - A description for the datacenter.
- `location` - The physical location of the datacenter (e.g. `de/txl`).
- `version` - The version of the datacenter, incremented with every change.
- `features` - List of features supported by the location the datacenter is in.
- `sec_auth_protection` - Whether a two-factor protection is enabled for the datacenter.
- `cpu_architecture` - Array of features and CPU families available in the location:
  - `cpu_family` - A valid CPU family name (e.g. `INTEL_SKYLAKE`).
  - `max_cores` - The maximum number of cores available.
  - `max_ram` - The maximum RAM size in MB.
  - `vendor` - A valid CPU vendor name.
- `ipv6_cidr_block` - The auto-assigned /56 IPv6 CIDR block, if IPv6 is enabled for the datacenter.
- `timeouts` - Always null; timeouts are configuration-only and are not returned by a query.
