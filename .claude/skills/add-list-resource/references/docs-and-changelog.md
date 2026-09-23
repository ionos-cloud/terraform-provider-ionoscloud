# Docs and CHANGELOG

Docs are hand-written; nothing validates them and no index needs updating.
`broken-link-checker.yml` checks changed `docs/` files and `CHANGELOG.md`: every URL and the
`../list-resources/<resource>.md` back-link must resolve. Why: `references/why.md`.

## §1 — New page `docs/list-resources/<resource>.md`

Frontmatter as `docs/list-resources/datacenter.md:1-7`, with `subcategory` verbatim from
`docs/resources/<resource>.md` and `page_title: "IONOS CLOUD: <resource>"` — bare stem, **no
`ionoscloud_` prefix**, even when the resource page prefixes it.

Section order, wording and HCL: that same page, `:9-142`, ending in the SDKv2-only `timeouts`
bullet; insert §2a then §2b after its `-generate-config-out` paragraph. Read it, do not copy it
wholesale: `:68-72` uses the phrasing §2a forbids and `:95` omits the confirmation §2b requires.

## §2a — Pagination caveat

Name the limit the client really sends, per endpoint: `scripts/docs.sh
limit-cloudapi <Operation>` (or `limit-bundle <Operation>`, the same probe) and
`limit-datasource <resource>` — SKILL.md decision 1, evidence in
`references/decisions-evidence.md`. **Never write "the API's default page limit"**: the fallback is the SDK client's, not the endpoint's, and `.Limit(n)`
overrides it. Never invent a number.
No `Limit`/`Offset` → no paging: drop the truncation warning.

**Never hedge one either.** The page states product behaviour; it never reports on what we did or
did not check. No "not verified", no "a figure the provider does not verify", no "reportedly".
Name the figure and what applies it, or name the bound with no figure at all — never print one and
then disown it. Silent truncation is behaviour, not a confession: it stays, stated as fact.

```markdown
> **Note:** The <things> are read with a single <API name> request[, which asks for up to **<n>**
> items — the same limit the `<type>` data source requests | , which asks for up to **<n>** items,
> the limit the IONOS SDK client sends when the caller sets none | . The provider sends no `limit`
> with it, so the <API name> applies its default page size of **<n>** <things> | . The provider
> sends no `limit` with it, so the results are bounded by the <API name>'s own page size].
```

Branch on the probe's verdict, in that order: an explicit `.Limit(n)` matching the data source; `the
client sends limit=N when the caller sets none`; nothing sent unless you call `.Limit(n)`, and that
operation's own doc comment names the default (`ZonesGet`: "Default limit is the first 100 items");
the same, with no documented default.

Close with `datacenter.md:70-72` verbatim, <thing>-substituted (silent truncation; check the count
before trusting `imported.tf`). Companion note under `## Argument Reference`: `datacenter.md:113` —
no `location` → drop its first clause. A third note, **regional products only**, in that same
section: `docs/list-resources/pg_cluster_v2.md:76-78` (supported `location` values, then the
"adding a `location` filter reduces the query to a single endpoint call" performance note).

## §2b — Label-collision warning — every page

`datacenter.md:74-103`, plus the step it lacks: `terraform state show <type>.smoke_0` to confirm the
address is the leftover import before `terraform state rm`, which stops Terraform managing it until
re-import.

Run `scripts/docs.sh forcenew <resource>` first. No force-new attribute → that page's
destroy paragraph, HCL and `must be replaced` advice are false; do not just delete them (a hazard
with no cost) but replace them with one paragraph on the in-place damage to the **wrong** object —
e.g. target groups: name, algorithm, protocol, targets, health checks overwritten, ALB rules
rebalancing.

## §3 — Append to `docs/resources/<resource>.md`

After `## Import`: `docs/resources/datacenter.md:84-120` — the `identity` import block,
`### Identity Schema`, then `## Query (List Resource)` with the list-all block and back-link.
`RequiredForImport` → `#### Required`; no others → no `#### Optional`.

## §4 — CHANGELOG.md

Two bullets under `### Features` (list resource, identity); copy the two `ionoscloud_datacenter`
bullets under `## 6.7.37`.

A `## X.Y.Z` heading is opened once per release cycle, by the first PR after a release; later PRs
append under it (#1025 opened `## 6.7.36`; #1029-#1031 only added bullets). Duplicating it is the
mistake to avoid.

**Check the git tag, not the heading: `scripts/docs.sh changelog-heading` is the
authority — do what it prints, either way.** #1034 is why: its bullets went under `## 6.7.36`, right
when the PR opened and wrong when it merged, because `v6.7.36` was tagged in between, outside the
merge's ancestry. It advertised an unreleased feature under a released version; a later PR moved
both bullets to `## 6.7.37`. An un-bumped heading is not evidence of being unreleased; re-check
right before merge.

Format: `## <semver>`, no date or link — the topmost heading may carry an ` -- upcoming release`
suffix, as `## 6.7.38` does; resource-scoped bullets lead with the backticked type and a colon. No
`Unreleased`, no changelog tooling.
