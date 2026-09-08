---
subcategory: "Application Load Balancer"
layout: "ionoscloud"
page_title: "IONOS CLOUD: target_group"
description: |-
  Lists IONOS CLOUD Target Groups.
---

# List Resource: ionoscloud_target_group

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists **Target Groups** on IONOS CLOUD. A target group is the pool of backend targets an Application Load Balancer forwarding rule balances traffic across.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated query file, whose name ends in `.tfquery.hcl` (for example `queries.tfquery.hcl`), separate from your main Terraform configuration.

### List all target groups

```hcl
list "ionoscloud_target_group" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter target groups by algorithm

```hcl
list "ionoscloud_target_group" "round_robin" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [{
      field_name  = "algorithm"
      field_value = "ROUND_ROBIN"
    }]
  }
}
```

### Filter by name and algorithm

```hcl
list "ionoscloud_target_group" "prod" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "name",      field_value = "production" },
      { field_name = "algorithm", field_value = "ROUND_ROBIN" },
    ]
  }
}
```

### Generate resource configuration from existing target groups

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_target_group` resource blocks for the target groups the query returns:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_target_group` resource block for each discovered target group into `imported.tf`, which can then be used directly in your configuration.

> **Note:** The target groups are read with a single Cloud API request, which asks for up to
> **200** items — the same limit the `ionoscloud_target_group` data source requests, and twice
> the 100 the IONOS SDK falls back to when no limit is asked for. Whether the server honours the
> requested limit in full is not verified. That ceiling is considerably
> lower than the 1000 the datacenter and IP block listings use, so it is reached sooner: a
> contract holding more than 200 target groups is truncated silently — no error is raised and
> the missing target groups simply do not appear. Check that the number of generated resource
> blocks matches the number of target groups you expect before treating `imported.tf` as
> complete.

Terraform names each generated resource after the `list` block label plus an index — a `list "ionoscloud_target_group" "smoke"` block produces `ionoscloud_target_group.smoke_0`, `smoke_1`, and so on.

### ⚠️ Do not reuse a `list` block label across separate imports

Because the generated names are derived from the `list` block label, running a second query with the **same** label produces the **same** resource addresses. If a previous address is still in state, the generated configuration is silently applied to the target group already recorded at that address — not to the one you just queried.

This happens because an `import` block is idempotent: Terraform skips it when the target address is already in state, so the identity in the generated `import` block is never consulted. Deleting the generated `.tf` file does **not** remove the state entry.

Target groups have no force-new attributes, so the result is an in-place update rather than a destroy — but it is still an update applied to the **wrong** target group: the name, algorithm, protocol, targets and health checks of the one you queried are written over the one in state, and the ALB forwarding rules pointing at it start balancing across the new target list.

To avoid this:

- Use a distinct `list` block label for each query you intend to import from, or
- Clear the address deliberately: run `terraform state show ionoscloud_target_group.smoke_0` first
  and confirm it is the leftover import and not a target group you still manage, then
  `terraform state rm ionoscloud_target_group.smoke_0`. Removing the address does not delete the
  target group — it stops Terraform managing it, and it has to be re-imported to come back under
  management.

Always read the plan before applying. A clean import reports:

```
Plan: 1 to import, 0 to add, 0 to change, 0 to destroy.
```

Anything reporting changes means the address is bound to a different target group — stop and clear the state entry first.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `algorithm`.
  - `field_value` - (Required) The exact value to match against.

> **Note:** Filtering happens after the single API response above is read, so it narrows the results but does not reduce the number of API calls, and it cannot recover a target group left out by the 200-item limit.

> **Note:** Values are compared exactly and **case-sensitively** against what the API returns, so `algorithm` must be given in the API's upper-case form (`ROUND_ROBIN`, `LEAST_CONNECTION`, `RANDOM`, `SOURCE_IP`) even though the `ionoscloud_target_group` resource itself accepts it in any case. An unknown `field_name` is rejected at plan time; a wrongly-cased `field_value` simply matches nothing.

> **Note:** `protocol` is deliberately not filterable — the API allows only the value `HTTP`, so a protocol filter could never narrow the results.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute | Description                       |
|-----------|-----------------------------------|
| `id`      | The UUID of the target group.     |

Unlike `ionoscloud_datacenter` and `ionoscloud_ipblock`, the target group identity carries no `location`: the collection is global and the resource has no location attribute.

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_target_group` resource schema:

- `id` - The UUID of the target group.
- `name` - The name of the target group.
- `algorithm` - The balancing algorithm.
- `protocol` - The forwarding protocol. Only `HTTP` is returned.
- `protocol_version` - The forwarding protocol version. Ignored unless `protocol` is `HTTP`.
- `targets` - The list of targets balanced by this group. Empty when the group has none:
  - `ip` - The IP of the balanced target.
  - `port` - The port of the balanced target.
  - `weight` - The traffic weight of the balanced target.
  - `proxy_protocol` - The proxy protocol version in use.
  - `health_check_enabled` - Whether the target is included in health checks.
  - `maintenance_enabled` - Whether the target is in maintenance mode.
- `health_check` - The health check settings, returned as a **list of at most one object** (`health_check[0]`) because the schema declares it `MaxItems: 1`; an empty list when the API returns none:
  - `check_timeout` - The health check timeout, in milliseconds.
  - `check_interval` - The interval between health checks, in milliseconds.
  - `retries` - The number of retries before a target is marked down.
- `http_health_check` - The HTTP health check settings, likewise a **list of at most one object** (`http_health_check[0]`); an empty list when the API returns none:
  - `path` - The path requested by the health check.
  - `method` - The HTTP method used.
  - `match_type` - What the response is matched on.
  - `response` - The value matched against.
  - `regex` - Whether `response` is a regular expression.
  - `negate` - Whether the match is negated.
- `timeouts` - Always null; timeouts are configuration-only and are not returned by a query.
