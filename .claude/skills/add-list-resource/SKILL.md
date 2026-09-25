---
name: add-list-resource
description: "Add a `terraform query` list resource, plus the resource identity it requires, to a terraform-provider-ionoscloud resource — code, unit test, docs, CHANGELOG, verification and PR review. Use when asked to make a resource queryable, add a list resource or a resource identity, support `terraform query` for an `ionoscloud_*` type, or repeat PR 1034's work on `ionoscloud_datacenter`."
user-invocable: true
argument-hint: <ionoscloud_type>, e.g. ionoscloud_ipblock
---

# Add a list resource + resource identity

Target: **$ARGUMENTS** (if empty, ask which `ionoscloud_*` type). Reproduces PR #1034
(reference: `ionoscloud/resource_datacenter_list.go`).

**Where the arc ends.** The work stays **uncommitted**: no staging, branch, commit, push,
`gh pr create` or review reply unless the user explicitly asks (`references/verify-and-pr.md` is
*how*, not permission). Report every file you created or edited (`git status --short` where git
is allowed) and the verification results, ask what next (commit, the PR, the tagged acceptance
test or a live `terraform query` — never yours to start — another resource), then wait. **No user
to answer** (a harnessed run): take your own recommendation at step 0.5; the decisions, their
evidence and the question go in the report.

## Machine limits — read first

**No subagents**, no parallel or background builds: a 12-core laptop in use. **Never**
`make testacc` or a build-tagged `go test` (`TF_ACC=1`, live cloud). No `go build ./...`,
`go test ./...`, `make test`; a package-scoped untagged `go test` is fine.

## Step 0 — Discovery, and six gates

`scripts/probe.sh discover <ionoscloud_type>` answers seven numbered questions.

**Gate 1 — branch.** (2) found, (3) not → SDKv2-backed. (3) from a `resource_*.go`, (2) not →
framework-native: more methods on the *same* struct. Only (2b), or only `data_source_*.go` /
`ephemeral_*.go` in (3) → nothing to list. Both → ask; neither → no such type.

**Gate 2 — a parentless collection GET?** `scripts/probe.sh listcalls <type>`; none → stop.

**Gate 3 — client (6).** `.NewCloudAPIClient(ctx, location)` / `...WithFailover(ctx)` → Cloud API
(`sdk-go/v6`), templates as written; a plain field (`.DNSClient`) → sdk-go-bundle product
(`sdkv2-branch.md` §2c deltas); a constructor taking a location (`.NewMongoClient`) → decision 4.

**Gate 4 — the state writer (7).** Found by **following Read**, never a name grep: any
`func(d *schema.ResourceData, obj <sdk>.<X>) error`, whatever its name, receiver or package
(`setDatacenterData`, `ionoscloud/resource_datacenter.go:357`; the by-value
`SetZoneData(d, dns.ZoneRead)`). **Diff its `d.Set`s against the schema's `Required: true`
entries** (`-generate-config-out` needs them all): a gap (a parent id, a write-only field) goes in
the docs, parent ids into the mapper (§2c). It must need nothing but one collection element: one
needing a `context.Context`, a client or more objects (`setBackupUnitData`, `setGroupData`) →
stop and report; never hand-map.

**Gate 5 — a parent in the schema?** (`ionoscloud_dns_record`: `zone_id`, cross-zone
`RecordsGet(ctx)`) → the child path (§1a, §2c): the identity carries it, the **mapper** sets it.

**Gate 6 — a bad candidate: stop.** A GET needing a parent ID (`DatacentersLansGet`) is an
unbounded N+1; an association-only resource (`ionoscloud_ipfailover`) has no API object.

Then open `references/sdkv2-branch.md` or `references/framework-native-branch.md`.

## Step 0.5 — Four decisions for the user, before any code

Present them together, each with your recommendation and its evidence, then wait.

**1. Pagination.** `scripts/probe.sh limit-cloud <Op>` / `limit-bundle <Op>`: what the client
sends unasked; no `Limit`/`Offset` on the request type → it does not page. `datasource-limit
<resource>`: a higher `.Limit(n)` in the sibling data source → send the same. The one offset/limit
loop is `services/dbaas/pgsqlv2/cluster.go` (stop on a short page: `Links.HasNext()` stays set).
#1034 declined one, as `ionoscloud/data_source_datacenter.go` makes the same call unpaginated:
fine, but **the docs must say so** (`docs-and-changelog.md` §2a).

**2. Identity**: the import ID's tuple (`splitImportID`, `ionoscloud/utils.go`): `location` if the
resource has one, then one string per parsed part, its own being `id`; never a mutable `name`.

**3. Filters**: `name`, plus the regional key as *the resource's own schema* spells it
(`location`/`region`), else another string attribute (never a bool or list), the two matching
different items; none narrows → none, say why; a one-value validator cannot narrow. The
`FilterAttribute(...)` allow-list is exactly the mapper's `MatchesFilters(...)` key set, or a
filter silently matches nothing. Matching is case-sensitive.

**4. Client, regional or global.** Cloud API → one global `NewCloudAPIClientWithFailover(ctx)`
call, as `datacenterListResource.List` (CRUD's per-location client picks an endpoint override,
not a partition: a fan-out returns each item N times; its `FilterGlobalOverrides` hole is known,
not yours). A bundle product → its own client: built per location (`.New<X>Client(ctx,
location)`) → partitioned, else one call covers every location. Partitioned: `scripts/probe.sh
locations <product>`; enumerable → fan out; **nothing enumerable → stop and ask** (one call =
one location's objects, silently incomplete).

## The plan

Identity, list resource + registration, tests, docs + CHANGELOG, ladder, report. Nothing new in
`internal/framework/identity/`; SDKv2 leaves `internal/framework/provider/provider.go` **untouched**.

## Verification ladder

`scripts/probe.sh rung0` … `rung5-scoped`, in order, one at a time: 0 `gofmt` after every edit;
1 build (SDKv2 `./ionoscloud/`; framework-native also `./internal/framework/provider/` for a new
package; a not-yet-vendored import fails here: rung 4, then 1); 2 the untagged unit test; 2b
`TestProvider$` (SDKv2), the **only** check of your `Identity`; 3 lint your own lines (keep
`--new-from-rev` and the path); 4 only if imports changed: vendor even when tidy is clean, report
`vendor/`; 5 `rung5-scoped`, once, last: the only compile of the tagged acceptance test. CI runs
no test: 2 and 2b are yours. `terraform query` itself: **only on explicit request**.

## Definition of done

- [ ] `Identity` declared, every attribute set on **every** read path; the import `d.SetId`s right
  after parsing; the legacy string import unchanged
- [ ] The resource file's diff against its pre-edit content (`git diff`, or a copy taken first)
  read line by line: every `-` line forced by the identity; no redundant `d.SetId` dropped, no
  re-worded comment, no tidying
- [ ] The list file's fetch and mapper read against the reference's (`sdkv2-branch.md` §2c): no
  check dropped, `apiResponse` never `_`
- [ ] Registered (`ionoscloud/list_resources.go`, or the service package's `ListResources()` plus
  a provider line **only** for a new package); **no model, no mapping of your own** (SDKv2)
- [ ] The tests `references/test-harness.md` lists; **the mutation check ran and failed**; the
  acceptance test written, **not run**
- [ ] The list page, the resource page's pointer section, the CHANGELOG (the **git tag** decides)
- [ ] Rungs 0–5 clean, **2b included**
- [ ] The pagination bound in the docs note and the fetch comment, read from the vendored client
- [ ] **Every comment and docs sentence states only what the code beside it does or checks, or
  what the vendored SDK or the resource's docs say**: no "every", "all" or "the whole contract"
  the code does not enforce, no reason it does not bear out, no unsourced product behaviour, no
  wording from this skill
- [ ] Uncommitted, no push or `gh` write, unless asked in so many words; reported, then asked

## Where to look things up

Paths are relative to this skill's directory. **Run the scripts, never read them**: bare, each
prints its subcommands; each check ends in a VERDICT line. `review.sh` beyond `test-unit` and
`queries` only once the user asks. References, in order: `sdkv2-branch.md` (**the traps live
here**) or `framework-native-branch.md`; `test-harness.md`; `docs-and-changelog.md`;
`verify-and-pr.md` only for a commit, PR or review.
