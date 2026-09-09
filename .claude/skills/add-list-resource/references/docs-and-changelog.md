# Docs and CHANGELOG

Docs here are **hand-written**. There is no `tfplugindocs`, no `templates/`, no `.web-docs/`,
no `go:generate`. Nothing verifies that docs match the schema.

**No index file needs updating.** `docs/list-resources/` is discovered by folder name.
`gitbook_docs/summary.md` indexes only `docs/resources/` and `docs/data-sources/` and has never
contained a list resource; `docs/index.md` is a landing page with no per-page index; `README.md`
does not mention list resources. `website/profitbricks.erb` is dead legacy — do not touch it.

**Links are checked on the PR.** `.github/workflows/broken-link-checker.yml` runs
`markdown-link-check` on modified files under `docs/` plus `CHANGELOG.md` and `README.md`. So
every `docs.ionos.com` URL and the `../list-resources/<resource>.md` back-link must resolve.
There is a pre-existing broken one to not copy: `docs/resources/pg_cluster_v2.md` links to
`../list-resources/psql_cluster_v2.md`, but the file is `pg_cluster_v2.md`.

---

## 1. New page: `docs/list-resources/<resource>.md`

Frontmatter conventions:

- `subcategory` — copied **verbatim** from `docs/resources/<resource>.md`.
- `layout` — always `"ionoscloud"`.
- `page_title` — `"IONOS CLOUD: <resource>"`, the **bare type stem**, no `ionoscloud_` prefix.
  All eight existing list pages use this form. Do **not** copy the resource page's title: about a
  third of `docs/resources/` prefixes the type (`docs/resources/pg_cluster_v2.md` has
  `IONOS CLOUD: ionoscloud_pg_cluster_v2` while `docs/list-resources/pg_cluster_v2.md` has
  `IONOS CLOUD: pg_cluster_v2`). The list page does not disambiguate itself in the title; only
  the H1 says "List Resource:".
- `description` — `Lists IONOS CLOUD <Plural Thing>.`

````markdown
---
subcategory: "<Subcategory>"
layout: "ionoscloud"
page_title: "IONOS CLOUD: <resource>"
description: |-
  Lists IONOS CLOUD <Things>.
---

# List Resource: <ionoscloud_type>

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists [<Things>](<docs_url>) on IONOS CLOUD.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated query file, whose name ends in `.tfquery.hcl` (for example `queries.tfquery.hcl`), separate from your main Terraform configuration.

### List all <things>

```hcl
list "<ionoscloud_type>" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter <things> by <field>

```hcl
list "<ionoscloud_type>" "<label>" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [{
      field_name  = "<field>"
      field_value = "<value>"
    }]
  }
}
```

### Filter <things> by <field_a> and <field_b>

```hcl
list "<ionoscloud_type>" "prod" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "<field_a>", field_value = "<value_a>" },
      { field_name = "<field_b>", field_value = "<value_b>" },
    ]
  }
}
```

### Generate resource configuration from existing <things>

Use `terraform query` with `-generate-config-out` to produce ready-to-use `<ionoscloud_type>` resource blocks for the <things> the query returns:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `<ionoscloud_type>` resource block for each discovered <thing> into `imported.tf`, which can then be used directly in your configuration.

<!-- pagination caveat here, if the implementation makes a single unpaginated call — §2a -->

<!-- label-collision warning here — §2b -->

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: <filter_fields>.
  - `field_value` - (Required) The exact value to match against.

<!-- regional services only: the supported-locations line and the performance note — §2c -->

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute  | Description                        |
|------------|------------------------------------|
| `id`       | The UUID of the <thing>.           |
| `<extra>`  | <description>.                     |

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `<ionoscloud_type>` resource schema:

- `id` - ...
- `<attr>` - ...
- `<block>` - <Block description>:
  - `<sub_attr>` - ...
- `timeouts` - Always null; timeouts are configuration-only and are not returned by a query.

> **Note:** `<attr>` is not available via the list resource as <reason>.
````

The `timeouts` bullet applies to the **SDKv2 branch** — the SDKv2 schema carries a `timeouts`
block that surfaces in the result as null. The framework-native pages do not have it.

Two things the two oldest pages (`s3_bucket.md`, `object_storage_accesskey.md`) lack: the
two-filter example, and the `## Identity Attributes` section. Follow the newer form instead —
`ipblock.md` and `target_group.md` are the closest models for an SDKv2-backed resource
(`target_group.md` for one with no `location`), then `datacenter.md` and `pg_cluster_v2.md`.

The supported-locations line and the performance note appear only on the three regional cluster
pages — `datacenter.md` has neither, because the Cloud API is not regional. Omit them for a
non-regional resource; the template marks them "regional services only" for that reason.

---

## 2. Reusable prose — copy near-verbatim rather than rewriting

### 2a. Pagination caveat — for any single-unpaginated-call implementation

Two variants — pick by whether the fetch passes an explicit `.Limit(n)`.

*No explicit limit (the SDK's own fallback applies):*

```markdown
> **Note:** The <things> are read with a single <API name> request, so the results are
> bounded by the page limit the IONOS SDK requests by default. A contract holding more
> <things> than that limit is truncated silently — no error is raised and the missing
> <things> simply do not appear. Check that the number of generated resource blocks matches
> the number of <things> you expect before treating `imported.tf` as complete.
```

*Explicit `.Limit(n)` (decision 1 said to match the data source):*

```markdown
> **Note:** The <things> are read with a single <API name> request, which asks for up to
> **<n>** items — the same limit the `<ionoscloud_type>` data source requests, and <k> times
> the <sdk-fallback> the IONOS SDK falls back to when no limit is asked for. Whether the
> server honours the requested limit in full is not verified. A contract holding more
> <things> than one response returns is truncated silently — no error is raised and the
> missing <things> simply do not appear. Check that the number of generated resource blocks
> matches the number of <things> you expect before treating `imported.tf` as complete.
```

**Never write "the API's default page limit".** The fallback is sent by the *SDK client*, not
applied by the endpoint (`sdk-go/v6/api_<x>.go`, the `DefaultQueryParams.Get("limit")` branch),
and when the fetch passes an explicit `.Limit(n)` the endpoint default is not what bounds the
result at all. Name the real number and say where it comes from. Keep the "not verified" hedge:
nothing in the repo confirms the server honours a requested limit in full.

Its companion under `## Argument Reference`, when the API returns everything in one response
so filtering saves no calls:

```markdown
> **Note:** The Cloud API returns <things> from every location in a single response, so filtering by `location` does not reduce the number of API calls; it only narrows the results. Filtering happens after that response is read, so it cannot recover a <thing> left out by the page limit described above.
```

If the resource has **no** location — `ionoscloud_target_group` does not — drop the first
clause rather than inventing a location: "Filtering happens after the single API response above
is read, so it narrows the results but does not reduce the number of API calls, and it cannot
recover a <thing> left out by the <n>-item limit."'

**Look the limit up per endpoint, and name the number you found.** There is no module-wide
default: the generated code sends `limit=1000` for most Cloud API collections but **100** for
`/ipblocks`, `/targetgroups` and user management, and the vendored doc comments say "Default
limit is the first 100 items" throughout regardless.

```bash
grep -n 'parameterToString(1\?0*, "")' vendor/github.com/ionos-cloud/sdk-go/v6/api_<x>.go
grep -n 'Limit(' ionoscloud/data_source_<resource>.go
```

If the fetch passes an explicit `.Limit(n)` — match whatever the sibling data source already
does — say `n`. If it does not, say the endpoint's default. Note which of the two it is, and
that whether the server honours the request is unverified. What you must not do is copy
datacenter's "1000" onto an endpoint that defaults to 100.

### 2b. Label-collision warning — recommended on every page

This documents a real near-miss: a smoke test renamed the maintainer's live datacenter and
nearly destroyed it. Only `<ionoscloud_type>` and `<thing>` change — but the section has two
halves and the second is conditional.

**First check whether the resource has a force-new attribute at all:**

```bash
grep -n 'ForceNew' ionoscloud/resource_<resource>.go
grep -n 'CustomizeDiff' ionoscloud/resource_<resource>.go
```

Many good list candidates have none — backup units, target groups, users, groups, CDN
distributions, certificates, DNS reverse records and autoscaling groups all return zero — and a
`CustomizeDiff` is usually not a substitute: an immutable-field check like
`resource_certificate_manager_certificate.go:64` *errors out* rather than planning a
replacement. If there is no force-new attribute, drop the destroy paragraph and its HCL block — but do not
just stop, or the section trails off having described a hazard without saying what it costs.
Replace them with one paragraph naming the in-place-update damage for *this* resource, e.g. for
`ionoscloud_target_group`: "Target groups have no force-new attributes, so the result is an
in-place update rather than a destroy — but it is still an update applied to the **wrong**
target group: the name, algorithm, protocol, targets and health checks of the one you queried
are written over the one in state, and the ALB forwarding rules pointing at it start balancing
across the new target list." Then trim the closing plan-reading advice to match — with nothing
force-new there is no `must be replaced` line to look for. The silent-wrong-target risk is real for every resource; the destroy-and-recreate
risk is not, and describing a destroy that cannot happen costs the whole warning its
credibility.

````markdown
Terraform names each generated resource after the `list` block label plus an index — a `list "<ionoscloud_type>" "smoke"` block produces `<ionoscloud_type>.smoke_0`, `smoke_1`, and so on.

### ⚠️ Do not reuse a `list` block label across separate imports

Because the generated names are derived from the `list` block label, running a second query with the **same** label produces the **same** resource addresses. If a previous address is still in state, the generated configuration is silently applied to the <thing> already recorded at that address — not to the one you just queried.

This happens because an `import` block is idempotent: Terraform skips it when the target address is already in state, so the identity in the generated `import` block is never consulted. Deleting the generated `.tf` file does **not** remove the state entry.

The consequences are not limited to a harmless diff. `<force_new_attr>` is a force-new attribute, so if the two <things> differ in it, the plan **destroys the <thing> already in state** and creates a replacement — the <thing> you meant to import is never touched:

```hcl
# generated for a <thing> with <force_new_attr> = "<b>", but smoke_0 in state points at one with "<a>"
resource "<ionoscloud_type>" "smoke_0" {
  <force_new_attr> = "<b>"     # forces replacement of the "<a>" <thing>
  name             = "<example-name>"
}
```

To avoid this:

- Use a distinct `list` block label for each query you intend to import from, or
- Clear the address deliberately: run `terraform state show <ionoscloud_type>.smoke_0` first
  and confirm it is the leftover import and not a <thing> you still manage, then
  `terraform state rm <ionoscloud_type>.smoke_0`. Removing the address does not delete the
  <thing> — it stops Terraform managing it, and it has to be re-imported to come back under
  management.

Always read the plan before applying. A clean import reports:

```
Plan: 1 to import, 0 to add, 0 to change, 0 to destroy.
```

Anything reporting changes, and especially `must be replaced` with `<force_new_attr> = "..." -> "..." # forces replacement`, means the address is bound to a different <thing> — stop and clear the state entry first.
````

### 2c. Regional services

```markdown
Supported `location` values: `de/fra`, `de/fra/1`, ...

> **Performance note:** When no `location` filter is set, the provider queries every regional endpoint in sequence. Adding a `location` filter reduces the query to a single endpoint call.
```

---

## 3. Append to `docs/resources/<resource>.md`

At the end of the file, after the existing `## Import` section:

````markdown
In Terraform v1.12.0 and later, the [`import` block](https://developer.hashicorp.com/terraform/language/import) can also be used with the `identity` attribute:

```hcl
import {
  to = <ionoscloud_type>.example
  identity = {
    id = "<thing> uuid"
  }
}

resource "<ionoscloud_type>" "example" {
  ### Configuration omitted for brevity ###
}
```

### Identity Schema

#### Required

* `id` (String) The UUID of the <thing>.

#### Optional

* `<extra>` (String) <description>.

## Query (List Resource)

<Things> can be listed using `terraform query` (requires Terraform 1.14+). List blocks must be placed in a dedicated query file, whose name ends in `.tfquery.hcl` (for example `queries.tfquery.hcl`).

```hcl
list "<ionoscloud_type>" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

See the [`<ionoscloud_type>` list resource documentation](../list-resources/<resource>.md) for filters and the full attribute reference.
````

Identity attributes that are `RequiredForImport` go under `#### Required`; drop the
`#### Optional` heading entirely if there are none.

`mariadb_cluster_v2.md` and `inmemorydb_cluster_v2.md` have neither section despite having
list resources and identity imports. That is drift, not a rule — follow the
datacenter/pg/s3/accesskey form.

---

## 4. CHANGELOG.md

Two entries, under `### Features`:

```markdown
- `<ionoscloud_type>`: New list resource, queryable with `terraform query` (requires Terraform 1.14+).
- `<ionoscloud_type>`: Add a resource identity (`id`, `<extra>`), which also enables `import` blocks with an `identity` attribute.
```

### The version-heading rule — the common mistake

The PR checklist says "Changelog updated **and version incremented**", but the observed rule is
subtler:

- There is **no version string anywhere in the source**. `main.go`, `GNUmakefile` and
  `.goreleaser.yml` carry none. The `CHANGELOG.md` heading *is* the version; releases are cut
  by pushing a `v*` tag.
- **The `## X.Y.Z` heading is created once per release cycle, by the first PR after a release.
  Every later PR appends under it.** `## 6.7.36` was created by #1025; #1029, #1030 and #1031
  all added bullets under it without bumping.

So: **read the top of `CHANGELOG.md` first.** If the topmost `## X.Y.Z` heading has no
corresponding git tag yet, append under the right `###` section (creating the section if
missing). Only create a new `## X.Y.Z+1` heading if the topmost one is already released.
Adding a duplicate heading is the mistake to avoid.

**#1034 is the counter-example, so check the tag and not just the heading.** It appended its two
`ionoscloud_datacenter` bullets under `## 6.7.36`, which was correct when the PR was opened and
wrong by the time it merged: `v6.7.36` had been tagged in between, at a commit that is not an
ancestor of the merge. The result advertised an unreleased feature under a released version, and a
later PR had to move both bullets into `## 6.7.37`. The heading being un-bumped is not evidence
that it is unreleased — the tag is:

```bash
top=$(sed -n 's/^## //p' CHANGELOG.md | head -1)
if git rev-parse -q --verify "refs/tags/v$top" >/dev/null; then
  echo "v$top is ALREADY TAGGED - add a new '## <next>' heading above it"
else
  echo "v$top is unreleased - append under it"
fi
```

Re-check this immediately before merge, not only when you write the entry.

Format: `## <semver>` (no date, no link), then `### Features` / `### Fixes` / `### Testing` /
`### Refactor` / `### Documentation` / `### Chore`. Resource-scoped entries lead with the
backticked type name and a colon; provider-wide entries have no prefix. There is no
`Unreleased` section, no `.changelog/` directory, and no changelog tooling.

Issue links use the full form: `(fixes [#1021](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/1021))`.
Any URL in a new bullet is link-checked.
