# Docs and CHANGELOG

No docs index to update. Every URL and the `../list-resources/<resource>.md` back-link must resolve
(`broken-link-checker.yml`); a docs.ionos.com link only if a page already has it, never
constructed. Every claim sourced (SKILL.md's rule).

## §1 — New page `docs/list-resources/<resource>.md`

Frontmatter as `docs/list-resources/datacenter.md:1-7`: `subcategory` verbatim from
`docs/resources/<resource>.md`; `page_title: "IONOS CLOUD: <resource>"`, **no `ionoscloud_`
prefix**. Body: that page's `:9-142`, **adapted, not copied** (your filters, your attributes,
ending in the SDKv2-only `timeouts` bullet), `:68-72` replaced by §2a, `:74-103` adapted per §2b.

## §2a — Pagination note

From `scripts/docs.sh limit-cloudapi <Op>` (or `limit-bundle <Op>`) and `limit-datasource
<resource>`. **Never "the API's default page limit"**, no number the probe did not print, **no
hedge** ("not verified", "reportedly"): the figure and what applies it, or the bound without one.
No `Limit`/`Offset` on the request: it does not page, no note.

```markdown
> **Note:** The <things> are read with a single <API name> request[, which asks for up to **<n>**
> items — the same limit the `<type>` data source requests | , which asks for up to **<n>** items,
> the limit the IONOS SDK client sends when the caller sets none | . The provider sends no `limit`
> with it, so the <API name> applies its default page size of **<n>** <things> | . The provider
> sends no `limit` with it, so the results are bounded by the <API name>'s own page size]. A
> contract holding more <things> than that limit is truncated silently — no error is raised and the
> missing <things> simply do not appear. Check that the number of generated resource blocks matches
> the number of <things> you expect before treating `imported.tf` as complete.
```

The bracket, by the probe's verdict: `.Limit(n)` matching the data source; the client sends
`limit=N`; nothing sent, a documented default (`ZonesGet`: "first 100 items"); nothing sent, none
documented. `n` above 100: add that `terraform query` stops at the `list` block's `limit`, 100
unless set. Under `## Argument Reference`: `datacenter.md:113`'s note (no `location` → drop its
first clause); a second, that `field_value` matches exactly and case-sensitively and an unknown
`field_name` is rejected at validation; **regional products only**, a third:
`docs/list-resources/pg_cluster_v2.md:76-78`.

## §2b — Label-collision warning — every page

`datacenter.md:74-103`, adapted: the destroy paragraph names every force-new attribute
(`scripts/docs.sh forcenew <resource>`), the HCL shows each differing, and the destroy's cost is
named only where the resource's docs or vendored API doc comments state it. The remedy is a shell
block, then "This only removes it from state; the <thing> itself is left in place.":

```shell
terraform state show <type>.smoke_0   # confirm the address holds the stale import
terraform state rm <type>.smoke_0
```

No force-new attribute → the destroy paragraph, HCL and `must be replaced` advice are false:
replace them with one paragraph on the in-place damage to the **wrong** object.

## §3 — Append to `docs/resources/<resource>.md`

After `## Import`, as `docs/resources/datacenter.md:84-120`: the `identity` import block,
`### Identity Schema` (`#### Required`; `#### Optional` only for `OptionalForImport`), and
`## Query (List Resource)` with the list-all block and back-link.

## §4 — CHANGELOG.md

Two bullets under `### Features`, shaped like `ionoscloud_datacenter`'s under `## 6.7.37`. **The
git tag decides, not the heading**: do what `scripts/docs.sh changelog-heading` prints (untagged:
under the topmost; tagged: a new `## X.Y.Z+1` above it); re-check before merge. Without git:
under the topmost, the check reported **not run**, never inferred from the heading or its
` -- upcoming release`. No duplicate heading, date, link or `Unreleased`.
