---
subcategory: "Database as a Service - PostgreSQL v2"
layout: "ionoscloud"
page_title: "Migrating from ionoscloud_pg_cluster to ionoscloud_pg_cluster_v2"
sidebar_current: "docs-guide-migrating-pgsql-v1-to-v2"
description: |-
  How to migrate existing PostgreSQL clusters from the legacy ionoscloud_pg_cluster resource to ionoscloud_pg_cluster_v2.
---

# Migrating PostgreSQL clusters from `ionoscloud_pg_cluster` to `ionoscloud_pg_cluster_v2`

`ionoscloud_pg_cluster_v2` is a re-implementation of the PostgreSQL cluster resource on the
newer plugin framework, backed by the v2 DBaaS PostgreSQL API. It is a **different resource type**
from the legacy `ionoscloud_pg_cluster`, with a different schema (single nested blocks instead of
lists, a required `location`, write-only credentials, etc.).

Because the two are different resource types served by different plugin frameworks, a plain
`moved {}` block **does not** work between them. Instead, the recommended migration adopts each
existing cluster into the new resource and drops the old one from state — **without touching the
running cluster**. No data is moved or recreated; only Terraform state changes.

> **Important:** This migration assumes `ionoscloud_pg_cluster` and `ionoscloud_pg_cluster_v2`
> manage the *same* underlying cluster (same UUID). It changes Terraform state only and never
> mutates, recreates, or deletes the live database.

## Prerequisites

- Terraform **1.14+** (for the `terraform query` list/config-generation workflow used below).
  Terraform **1.7+** is enough if you choose to write the `import` blocks by hand.
- A backup of your state file (`terraform state pull > state.backup.json`).

## Recommended approach: query-generate + `removed`

The `ionoscloud_pg_cluster_v2` **list resource** can enumerate your existing clusters and
generate complete `v2` configuration for each one, so you don't have to hand-write the new
resource blocks or look up cluster UUIDs.

### 1. Discover clusters and generate v2 configuration

Create a query file, for example `migrate.tfquery.hcl`:

```hcl
list "ionoscloud_pg_cluster_v2" "all" {
  provider         = ionoscloud
  include_resource = true
  # Optional: scope discovery to a single region.
  # config {
  #   filters = [
  #     { field_name = "location", field_value = "de/txl" },
  #   ]
  # }
}
```

Then generate configuration and import blocks for the discovered clusters:

```sh
terraform query -generate-config-out=generated_pg_v2.tf
```

`generated_pg_v2.tf` will contain a fully-populated `ionoscloud_pg_cluster_v2` block for each
cluster (name, version, instances, connections, maintenance window, backup, logs/metrics, …),
together with the `import` wiring keyed by the v2 identity (`<location>:<cluster_id>`).

### 2. Drop the legacy resources from state

For each legacy `ionoscloud_pg_cluster` you are migrating, add a `removed` block so Terraform
forgets it **without destroying** the live cluster:

```hcl
removed {
  from = ionoscloud_pg_cluster.example
  lifecycle {
    destroy = false
  }
}
```

Delete the corresponding `resource "ionoscloud_pg_cluster" "example" { … }` configuration.

### 3. Fill in the values that cannot be generated

A few attributes are not part of the cluster's readable state and must be added by hand to the
generated `ionoscloud_pg_cluster_v2` blocks:

| Attribute                      | Why it's missing                            | What to do                                                                 |
|--------------------------------|---------------------------------------------|----------------------------------------------------------------------------|
| `credentials.password`         | Write-only; the API never returns it.       | Set it (reuse the same value, e.g. via a variable or `random_password`).   |
| `credentials.password_version` | Provider-side marker with no v1 equivalent. | Set an initial value such as `"v1"`.                                       |
| `timeouts`                     | Optional; not generated.                    | Add only if you previously customized timeouts.                            |

`restore_from_backup` is create-time-only and not part of a cluster's readable state, so config
generation never emits it and you do not need it when adopting an existing cluster. Leave it out
(it is only relevant when initializing a brand-new cluster from a backup).

### 4. Apply

```sh
terraform plan    # confirm: imports the v2 resources, removes v1 from state, NO destroys/replaces
terraform apply
```

The plan should show the clusters being **imported** into `ionoscloud_pg_cluster_v2`, the
`ionoscloud_pg_cluster` resources being **removed from state**, and **no** create/destroy/replace
actions against the live clusters. If you see a forced replacement, stop and reconcile the
offending attribute before applying.

## Manual alternative (no config generation)

If you prefer not to use the query workflow, you can do the same thing with explicit blocks. For
each cluster, write the `ionoscloud_pg_cluster_v2` configuration yourself and pair it with:

```hcl
import {
  to = ionoscloud_pg_cluster_v2.example
  id = "de/txl:00000000-0000-0000-0000-000000000000" # "<location>:<cluster_id>"
}

removed {
  from = ionoscloud_pg_cluster.example
  lifecycle {
    destroy = false
  }
}
```

Then `terraform plan` / `terraform apply` as above. You can find the cluster UUID in the legacy
resource's state (`terraform state show ionoscloud_pg_cluster.example`).

## Rollback

Because the live cluster is never touched, rollback is just reverting your configuration changes
and restoring the pre-migration state file if you already applied:

```sh
terraform state push state.backup.json
```
