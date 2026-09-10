---
name: add-list-resource
description: "Add a `terraform query` list resource, plus the resource identity it requires, to a resource in terraform-provider-ionoscloud — code, unit test, docs, CHANGELOG, verification and PR review. Use when asked to make a resource queryable, add a list resource, add a resource identity, support `terraform query` for an `ionoscloud_*` type, or repeat what PR 1034 did for `ionoscloud_datacenter`."
user-invocable: true
argument-hint: <ionoscloud_type>, e.g. ionoscloud_ipblock
---

# Add a list resource + resource identity

Target resource: **$ARGUMENTS** (if empty, ask which `ionoscloud_*` type before doing anything else).

This skill reproduces the work done for `ionoscloud_datacenter` in PR #1034 — the reference
example, `ionoscloud/resource_datacenter_list.go` — for another resource. The
deliverable is the whole arc, not just the Go file:

1. resource identity on the managed resource
2. the list resource
3. registration — one name in a `ListResources()` slice
4. a unit test that actually asserts something
5. a query step on the resource's tagged acceptance test (written, not run)
6. `docs/list-resources/<x>.md` + an identity/query section on `docs/resources/<x>.md`
7. `CHANGELOG.md`
8. verification — **then stop and report**

A list resource without an identity, or with a test that passes when the mapper is wrong,
is not done. Both were real findings on #1034.

**Where the arc ends.** Deliver the files and the verification results, then stop.
`git commit`, `git push`, `gh pr create` and replying to a review comment are writes the
maintainer makes, or explicitly asks you to make — never a step you take because the
checklist has one line left. `references/verify-and-pr.md` documents *how*, for when
you are asked; it is not permission.

---

## Machine limits — read before running anything

The maintainer works on a 12-core laptop while you work, with an IDE, gopls and often a kind
cluster running. A previous agent fan-out drove it to load average 55 and made it unusable.

- **Do not fan out to subagents for this task.** It is sequential work, and the expensive part
  is compilation — N agents × ~10 Go compile threads swamps 12 cores. One agent, one build at
  a time.
- **Never** `make testacc`, or any `go test` with a build tag. `GNUmakefile:22-23` sets
  `TF_ACC=1` and `-tags`; the tagged suites run against real IONOS credentials and create live
  cloud resources (CI budgets them at 240m–6h — `.github/workflows/compute-test-run.yml:32`,
  `e2e.yaml:65`).
- **Avoid** `go build ./...`, `go test ./...` and `make test` — not because they reach the cloud
  (`make test` sets no `TF_ACC` and no tag, `GNUmakefile:17-20`, so no acceptance suite is even
  compiled) but because `TEST?=$$(go list ./... |grep -v 'vendor')` expands to all 60 non-vendor
  packages and `-parallel=4` saturates the maintainer's cores. Getting the reason right matters:
  an **untagged, package-scoped `go test` is always safe** and the ladder below depends on it.
- Scope every build and test to a package path. Follow the verification ladder at the bottom.
- Never run commands in parallel or in the background.
- If the user reports slowness, check `uptime` and stop what you are running.

---

## Placeholders used throughout

| Placeholder | Meaning | Datacenter example |
|---|---|---|
| `<ionoscloud_type>` | terraform type name | `ionoscloud_datacenter` |
| `<Resource>` | Go-exported name | `Datacenter` |
| `<resource>` | lower-case name / doc filename stem | `datacenter` |
| `<service>` | framework service package — **framework-native branch only**; an SDKv2-backed list resource lives in package `ionoscloud`, beside the resource it lists | `pgsqlv2` |

---

## Step 0 — Discovery

Everything downstream depends on one branch: **is the managed resource SDKv2 or
framework-native?** Determine it before writing anything.

```bash
# 1. Resolve the type-name constant (ResourcesMap is keyed by constants, never by
#    string literals — a naive grep of provider.go finds nothing).
grep -n '= "<ionoscloud_type>"' utils/constant/constants.go

# 2. Is it in the SDKv2 ResourcesMap?  Match on the FACTORY, not the constant: ResourcesMap
#    (provider.go:96) and DataSourcesMap (:156) frequently key off the SAME constant — a bare
#    `grep -n 'constant.DatacenterResource:'` returns :97 AND :157 — and some types have only
#    a *DataSource constant, so the constant's own suffix does not settle it either.
#    (The [Rr]/[Dd] classes are for the three exported autoscaling factories.)
grep -nE 'constant.<Const>: *[Rr]esource'   ionoscloud/provider.go   # ResourcesMap   -> listable
grep -nE 'constant.<Const>: *[Dd]ataSource' ionoscloud/provider.go   # DataSourcesMap -> not, on its own

# 3. Is it framework-native?  Most hits here are data sources (and one ephemeral resource),
#    so check the FILE the hit is in: only a `resource_*.go` hit means framework-native.
grep -rn 'ProviderTypeName + "_<suffix>"' internal/framework/services/ --include=*.go | grep -v _test.go

# 4. Does it already declare an identity?
grep -rn 'Identity: &schema.ResourceIdentity' ionoscloud/resource_<resource>.go
grep -rn 'IdentitySchema' internal/framework/services/<service>/

# 5. Does a list resource already exist?  The SDKv2-backed ones live in ionoscloud/,
#    the framework-native ones under internal/framework/services/.
ls docs/list-resources/
find ionoscloud internal/framework/services -name '*_list.go'
cat ionoscloud/list_resources.go
```

- Hit in (2), not (3) → **SDKv2-backed**. Read `references/sdkv2-branch.md`. This is the
  harder path and the one this skill is mostly about.
- Hit in (3) from a `resource_*.go`, not (2) → **framework-native**. Read
  `references/framework-native-branch.md`. Much less work: the list resource is more methods
  on the *same* struct as the managed resource, in the same package, so there is no
  `RawV6Schemas` and no SDKv2 bridge at all.
- Only a `[Dd]ataSource` hit in (2), or only `data_source_*.go` / `ephemeral_*.go` hits in
  (3) → there is **no managed resource** of that type. It cannot be listed; say so and stop.
- Hits in both (2) and (3) → no type in the tree does this today, so treat it as a genuine
  ambiguity rather than picking: report both hits and ask.
- Hits in neither → the type does not exist. Stop and say so.

Then find the API call that lists them — **if there isn't one, stop**:

```bash
# every parentless collection GET in the vendored SDKs.
# note [A-Za-z0-9]+ — without the digits this silently misses K8sGet.
grep -rhoE "func \(a \*[A-Za-z0-9]+\) [A-Za-z0-9]+(Get|List)\(ctx _?context\.Context\) Api[A-Za-z0-9]+Request" \
  vendor/github.com/ionos-cloud/ | sed 's/^func (a \*//' | sort -u
```

A resource whose collection GET needs a parent ID (`DatacentersLansGet(ctx, datacenterId)`,
`K8sNodepoolsGet(ctx, clusterId)`, …) is a **bad candidate**: listing it means enumerating
every parent first, an unbounded N+1 with no precedent in this repo. So are
association-only resources with no API object of their own (`ionoscloud_ipfailover`,
`ionoscloud_server_boot_device_selection`, `ionoscloud_datacenter_nsg_selection`). Say so
and stop rather than inventing a fan-out.

---

## Step 0.5 — Four decisions to settle with the user, before writing code

The maintainer wants these agreed up front, not discovered in review. Present them together,
briefly, with your recommendation — then wait.

**1. Pagination.** Decide explicitly; do not default silently. The provider-wide state:
there is exactly **one** offset/limit loop in the whole provider
(`services/dbaas/pgsqlv2/cluster.go` — it stops on a short page and deliberately does *not*
use `Links.HasNext()`, which stays populated past the last page and spins forever). The S3
paths use continuation tokens. Everything else is unpaginated. On #1034 the user **declined**
pagination for the datacenter list resource and documented the limitation instead, because
fixing only the list resource would leave it inconsistent with the identical unpaginated call
in `ionoscloud/data_source_datacenter.go`. That precedent is defensible and reusable — but if
you follow it, **the docs must say so** (wording in `references/docs-and-changelog.md`), and
Copilot will raise pagination on the PR, so have the answer ready.

Before deciding, look up **this endpoint's** default limit and what the sibling data source
already does. The default is *not* uniform — it is 1000 for most Cloud API collections but
**100** for `/ipblocks`, `/targetgroups` and user management, and nothing in the provider sets
`DefaultQueryParams`:

```bash
grep -n 'parameterToString(1\?0*, "")' vendor/github.com/ionos-cloud/sdk-go/v6/api_<x>.go
grep -n 'Limit(' ionoscloud/data_source_<resource>.go
```

If the data source sets an explicit `Limit(...)` — `data_source_ipblock.go:146` passes
`constant.IPBlockLimit` (1000) — then a bare `.Depth(1)` fetch caps the *list resource* below
its own data source, which is the exact inconsistency the decline argument rests on. Match the
data source's `Limit()`, and name that endpoint's real number in the docs note.

**2. Identity attributes.** Read them off the resource's existing import ID — that tuple is
already exactly what an importer needs. `splitImportID` (in `ionoscloud/utils.go`) parses
`"<location>:<id-1><del><id-2>…"`, so the identity is `location` (optional, when the resource
has one) plus one string attribute per parsed part, the resource's own being `id`. A child
resource carries its parent IDs, named exactly as its own schema names them
(`datacenter_id`, `server_id`), all `RequiredForImport`. Never put a mutable attribute such
as `name` in an identity — Terraform fails refresh if any identity attribute changes.

**3. Filter fields.** The allow-list passed to `identity.FilterAttribute(...)` must be
exactly the key set of the map handed to `identity.MatchesFilters(...)` in the mapper.
Nothing enforces this; a mismatch means a filter silently matches nothing. Default to
`name` plus the regional key under whatever name *the resource's own schema* uses
(`location` vs `region`).

Two checks before you settle the list:

- **Drop any field that cannot narrow.** A field whose schema validator admits exactly one
  value is a filter that always matches everything — noise in the allow-list and a
  meaningless doc example. `ionoscloud_target_group`'s `protocol` is the case in point:
  `StringInSlice([]string{"HTTP"}, true)`, and the API model says "Only the value 'HTTP' is
  allowed". It was proposed as a filter field on a `target_group` list resource and dropped for
  exactly this reason; that list resource was never merged, so only the managed resource's own
  schema is in the tree to look at. Read the `ValidateDiagFunc` of every candidate.
- **`MatchesFilters` is an exact, case-sensitive compare** (`filter.go:57-60`) against the
  value the API returned, but an SDKv2 `StringInSlice(..., true)` validator lets the *resource*
  accept any case. So a filterable enum is case-sensitive in a `list` block and
  case-insensitive in the resource. A bad `field_name` is a loud plan-time error; a
  wrongly-cased `field_value` is silent. Say so in the docs whenever a filterable field is a
  case-insensitively-validated enum.

**4. Regional or global.** Decide from the **shape of the collection endpoint**, not from the
managed resource's client constructor. The constructor is a red herring: `resource_datacenter.go`
calls `NewCloudAPIClient(ctx, location)` on every CRUD path, yet the datacenter *list* resource
uses `NewCloudAPIClientWithFailover(ctx)` and makes one call — because `location` there only
selects an endpoint override (`bundleclient.go:393-417`), it does not partition the collection.

- **Cloud API (`sdk-go/v6`) — always global.** `/datacenters`, `/ipblocks` and the rest are
  single un-partitioned collections that return objects from every location in one response.
  Use `NewCloudAPIClientWithFailover(ctx)` and one call; copy the client construction inside
  `datacenterListResource.List` (`ionoscloud/resource_datacenter_list.go`). Its doc comment
  says "intended for resources that do not have a location attribute", but datacenter has one
  and uses it anyway — for a *collection read* that is correct. Fanning out here re-reads the
  same global collection N times and returns every item N times.

  *Known limitation of the pattern, not something you have to solve.*
  `NewCloudAPIClientWithFailover` resolves its endpoint from **global** cloud overrides only
  (`FilterGlobalOverrides` keeps the entries whose `location` is empty), so a file config that
  overrides the cloud product with per-location endpoints and no global one makes it hard-error
  with `no global failover endpoints configured for "cloud"` — on a configuration where the
  managed resource's own `NewCloudAPIClient(ctx, location)` is fine, because that one falls back
  the other way (location override first, global second). The shipped `ionoscloud_datacenter`
  list resource has the same hole. Follow the pattern anyway; just have the answer ready if
  review raises it.
- **Regional = an sdk-go-bundle DBaaS product that exposes `AvailableLocations()`** in
  `services/dbaas/<product>/client.go` — today exactly `pgsqlv2`, `mariadbv2`, `inmemorydbv2`.
  Only then fan out; copy `internal/framework/services/pgsqlv2/resource_pg_cluster_list.go`.
  There is **no `AvailableLocations()` for the Cloud API**, so the fan-out snippet in
  `sdkv2-branch.md` §2c cannot even be filled in for one.

```bash
grep -rn 'func AvailableLocations' --include=*.go services/   # no hit for your product -> global
```

---

## The plan

| Step | What | Reference |
|---|---|---|
| 1 | Resource identity on the managed resource | `references/sdkv2-branch.md` §1 (or `framework-native-branch.md` §1) |
| 2 | The list resource + its one registration line | `references/sdkv2-branch.md` §2–3 (or `framework-native-branch.md` §2–3) |
| 3 | The unit test, then the mutation check | `references/test-harness.md` |
| 4 | Docs + CHANGELOG | `references/docs-and-changelog.md` |
| 5 | Verification ladder | below |
| 6 | **Stop and report.** Commit / PR / review replies only when the user asks | `references/verify-and-pr.md` |

Do them in that order. Step 1 is a hard prerequisite for step 2: a resource with no
`Identity` cannot be listed at all — `identity.SetRawV6Schemas` gives up, and
`RawV6Schemas` has no diagnostics channel, so the only trace is a `tflog.Error` line.

### The files

On the SDKv2 branch every line you write is in **package `ionoscloud`, next to the resource
being listed**. It has to be: the list resource calls that resource's own package-level
`resource<X>()`, state writer and `set<X>Identity()`, and no package under
`internal/framework` can import package `ionoscloud` — `ionoscloud/provider_test.go` is an
in-package test that imports `internal/framework/provider`, so the reverse import is a cycle
in the test build. Framework code in `ionoscloud/` looks wrong and is not; a list resource can
only be written with terraform-plugin-framework, whichever half of the mux its resource is on.

The state writer's name is **not** uniform, so read the resource file instead of assuming:
`resource_datacenter.go` has the unexported `setDatacenterData`, `resource_ipblock.go` the
exported `IpBlockSetData`. `set<X>Data` everywhere in this skill is a placeholder for whatever
that resource calls its writer, not a name to grep for.

| SDKv2-backed | framework-native |
|---|---|
| `ionoscloud/resource_<resource>.go` — the `Identity` block, `set<X>Identity`, its call sites, the dual-mode importer | `internal/framework/services/<service>/resource_<resource>.go` — `IdentitySchema` + `Identity.Set` on every state-producing path |
| **new** `ionoscloud/resource_<resource>_list.go` | **new** `internal/framework/services/<service>/resource_<resource>_list.go` |
| **new** `ionoscloud/resource_<resource>_list_test.go` (package `ionoscloud_test`) — the shared helpers (`muxedProviderServer`, `listResults`, `decode`, `goValue`, …) are already in `ionoscloud/resource_datacenter_list_test.go`; reuse them, do not re-copy | **new** `internal/framework/services/<service>/resource_<resource>_list_test.go` |
| one name in `ionoscloud/list_resources.go`'s `ListResources()` | one name in the service package's `ListResources()`, and — only for a **new** service package — one `<service>.ListResources(),` line in `internal/framework/provider/provider.go` |
| `ionoscloud/resource_<resource>_test.go` — the `Query: true` acceptance step | same |
| docs + CHANGELOG | same |

`internal/framework/provider/provider.go` is **not** touched on the SDKv2 branch.
`provider.New(sdkv2ListResources ...func() list.ListResource)` takes the whole slice, and every
call site (`main.go`, `xpprovider/provider.go`, `internal/acctest/acctest.go`,
`internal/framework/provider/provider_test.go`, `ionoscloud/provider_test.go`, and the test
harness itself) passes `ionoscloud.ListResources()...` and never names an individual resource.

### Shared plumbing that already exists — do not rewrite it

`internal/framework/identity/`:

| Symbol | File | What it does |
|---|---|---|
| `FiltersKey` | `filter.go` | `"filters"`, the config attribute name |
| `Filter` | `filter.go` | `{FieldName, FieldValue types.String}` |
| `FilterAttribute(allowed...)` | `filter.go` | the `filters` list attribute; attaches `stringvalidator.OneOf` |
| `FilterValue(filters, name)` | `filter.go` | pull one filter value out, to push down to the API |
| `MatchesFilters(fields, filters)` | `filter.go` | client-side AND match, exact string compare |
| `StreamList[T](ctx, stream, req, fetch, mapper)` | `list.go` | the whole List body |
| `MappedItem` | `identity.go` | `{DisplayName string; Identity, Resource any}` |
| `Model` | `identity.go` | reusable identity model for a lone `id` (framework-native branch — the SDKv2 branch takes its identity out of the `ResourceData`) |
| `SetRawV6Schemas(...)` | `sdkv2.go` | SDKv2 → protocol-v6 schema bridge (SDKv2 branch only) |
| `MappedItemFromResourceData(name, rd, incl)` | `sdkv2.go` | turns a `ResourceData` the resource's own writers have filled into a `MappedItem` — identity and resource as `tftypes.Value`, `timeouts` nulled back out (SDKv2 branch only) |

Nothing new belongs in that package for an ordinary resource.

---

## Verification ladder

Cheapest first. Substitute `<service>` and `<Resource>`. Strictly sequential.

**Rung 0 — after every edit, free:**
```bash
gofmt -l -d ./ionoscloud ./internal ./utils
```
**Never let `gofmt` run with an empty argument list — it silently reads stdin and exits
clean.** That is what `gofmt -l -d $(git diff --name-only master...HEAD ...)` does here:
`master...HEAD` is a *commit* range, your work is uncommitted (you are not committing it,
see "Where the arc ends"), so the substitution expands to nothing and the rung reports green
having inspected none of your new files. Formatting then fails in CI instead
(`.github/workflows/build.yml`, `Run gofmt check`, `-l -d` over `.`).

**Rung 1 — type-check the packages you touched. Main iteration loop:**
```bash
# SDKv2 branch — the identity, the list resource and its registration are all in here:
go build ./ionoscloud/

# framework-native branch:
go build ./internal/framework/services/<service>/
# additionally, ONLY if the list resource went in a NEW service package (i.e. you added a
# `<service>.ListResources(),` line to internal/framework/provider/provider.go):
go build ./internal/framework/provider/
```
On the SDKv2 branch this really is one command: `ionoscloud/list_resources.go` is compiled by
the same build as the list resource, so a typo in the registration line fails here rather
than at runtime. On the framework-native branch it is not — `./internal/framework/services/<service>/`
does not reach `internal/framework/provider`, so a slip in the wiring line you were just told
to add is invisible until rung 5 unless you run the second command.

**Rung 1 can fail for a reason only rung 4 fixes — the one permitted jump in the ladder.** The
repo vendors its dependencies (`vendor/modules.txt` exists, so `-mod=vendor` is the default) and
the compiler sees only packages already under `vendor/`. An import of a not-yet-vendored
subpackage therefore breaks *this* rung, as a cannot-find-module / missing-package error naming
that import path — it does not wait for rung 4. When that is the failure, go run rung 4, then
come back and re-run rung 1.

**Rung 2 — the untagged unit test. Main iteration loop. No credentials, no network:**
```bash
# SDKv2 branch (the test is ionoscloud/resource_<resource>_list_test.go, package ionoscloud_test):
go test ./ionoscloud/ -run 'Test<Resource>ListResource' -count=1

# framework-native branch:
go test ./internal/framework/services/<service>/ -run 'Test<Resource>ListResource' -count=1
```

**Rung 2b — validate the identity schema (SDKv2 branch). No credentials, no network:**
```bash
go test ./ionoscloud/ -run 'TestProvider$' -count=1
```
This is the **only** thing that executes `Provider().InternalValidate()` →
`InternalIdentityValidate()` on your new `Identity`
(`vendor/.../helper/schema/provider.go:220-225`), which is what turns a malformed identity
into a unit-test failure instead of a runtime one. The skill elsewhere calls that check
"free" — it is free only because you run this rung. Nothing in CI runs it, rung 1 compiles no
test file, and rung 2's `GetResourceIdentitySchemas` assertion passes for a malformed identity.
On the SDKv2 branch rung 2 has already built this package's test binary, so this is one more
`-run` over it; on the framework-native branch rung 2 was scoped to another package and this
is the one small compile you pay for.

**Rung 3 — lint only your own lines, once the code compiles and the test is green:**
```bash
# SDKv2 branch:
golangci-lint run --new-from-rev $(git merge-base origin/master HEAD) ./ionoscloud/...

# framework-native branch:
golangci-lint run --new-from-rev $(git merge-base origin/master HEAD) \
  ./internal/framework/services/<service>/... ./internal/framework/provider/...
```
`--new-from-rev` narrows the *report*, not the work: `make lint` (`GNUmakefile:14-15`) is the
same flags with **no path argument**, so golangci-lint loads and analyses all 60 non-vendor
packages under ~45 linters. On this machine that is the rung most likely to reproduce the
load-55 incident, so pass paths in the loop and save `make lint` for a single run at the end.
Never drop `--new-from-rev` — the whole-package baseline in these packages runs to a couple of
hundred pre-existing issues that are not yours and not actionable, and CI lints only the diff.

**Rung 4 — only if you added or removed an import:**
```bash
go mod tidy && git diff --exit-code -- go.mod go.sum
go mod vendor && git status --porcelain vendor/
```
**The trap:** importing a *new subpackage of a module already in `go.mod`* changes
`vendor/modules.txt` and adds vendored sources while leaving `go.mod`/`go.sum` untouched.
That is exactly what happened on #1034 with `terraform-plugin-mux/tf5to6server/translate`.
Run the vendor half even when the tidy half is clean, and **report `vendor/` as part of the
change** so it goes in the same commit — do not run `git commit` yourself.

**Rung 5 — expensive, exactly once, at the very end:**
```bash
go vet -tags=all ./...
```
136 of the repo's 157 `_test.go` files sit behind build tags and are invisible to a plain
`go vet ./...`. That includes the `Query: true` step you wrote — the resource's acceptance
test file is tagged (`ionoscloud/resource_datacenter_test.go` is
`//go:build compute || all || datacenter`), so this rung is the only thing that compiles it.
It also catches a signature change reaching test files in packages you never opened.

This compiles all 60 non-vendor packages **and their test variants**, so it is strictly heavier
than the `go build ./...` the machine-limits section tells you to avoid. **Tell the user before
you start it**, and offer to leave it to CI (`.github/workflows/build.yml`, `Run go vet with all
build tags`) and run only the scoped form:
```bash
go vet -tags=all ./ionoscloud/ ./internal/framework/provider/ ./internal/acctest/
# framework-native branch: add ./internal/framework/services/<service>/
```

Note: **no test job runs on PRs.** Green CI does not mean the tests pass. Rungs 2 and 2b are
on you.

Note also what the ladder does **not** prove: it exercises the ListResource *RPC*, never the
`terraform query` command. There is no repo tooling for that — no `dev_overrides` example
anywhere in the tree, no make target. A manual check needs a `~/.terraformrc` `dev_overrides`
block pointing at a locally built binary plus a `.tfquery.hcl` file, and it runs against real
credentials. **Only on explicit request**, and only with a throwaway `list` block label — see
the hazard in `references/docs-and-changelog.md` §2b.

---

## Definition of done

- [ ] `Identity` on the managed resource, every attribute set on **every** read path
- [ ] `d.SetId(...)` right after parsing the import parts (SDKv2 branch)
- [ ] Legacy string import still works, byte-for-byte unchanged
- [ ] List resource registered — one name in `ionoscloud/list_resources.go`'s `ListResources()`
      (SDKv2 branch); or one name in the service package's `ListResources()`, plus the
      `<service>.ListResources(),` line in the framework provider **only** when that service
      package is new (framework-native branch)
- [ ] **No model, and no attribute mapping of your own** (SDKv2 branch): your `map<X>`
      does the filter check and nothing else — every result is
      `r.resourceSchema.Data(&terraform.InstanceState{})` → `set<X>Data` → `set<X>Identity`
      (in that order — the identity setter reads values back out) →
      `identity.MappedItemFromResourceData`. If you wrote a `tfsdk`-tagged struct, delete it
- [ ] Unit test asserts every mapped attribute, every stub result, and the null-vs-zero
      distinction (an omitted nested **block** can decode to `[]any{}` rather than nil —
      `test-harness.md`) — and **the mutation check was run and failed**
- [ ] One filter subtest **per field** in the `FilterAttribute(...)` allow-list, plus an AND
      case; and the stub asserts the query params the fetch closure sets (`depth`, and `limit`
      when an explicit one is passed)
- [ ] A `Query: true` step on the resource's tagged acceptance test (written; **not run**),
      plus a non-ForceNew update step whenever Update writes the identity itself
- [ ] `docs/list-resources/<resource>.md` + the pointer section on `docs/resources/<resource>.md`
- [ ] CHANGELOG entry under the correct `## X.Y.Z` heading — decided by the **git tag, not the
      heading**: append under the topmost heading only while no `v<that version>` tag exists; once
      that version is tagged, create a new `## X.Y.Z+1` heading above it. The check and the #1034
      counter-example are in `references/docs-and-changelog.md` §4 — re-run it just before merge
- [ ] Ladder rungs 0–5 clean, **2b included**
- [ ] The pagination decision is written down somewhere a reviewer will find it, with this
      endpoint's real default limit
- [ ] Reported, not committed: no `git commit`, `git push` or `gh` write unless the user asked

## Reference files

- `references/sdkv2-branch.md` — the SDKv2 path: identity on the resource, the list resource
  in package `ionoscloud`, registration. **The traps live here.**
- `references/framework-native-branch.md` — the shorter path, and what *not* to copy from
  the five existing examples.
- `references/test-harness.md` — the test file, the verbatim-copy helpers, the mutation check.
- `references/docs-and-changelog.md` — doc template, reusable prose, CHANGELOG rules.
- `references/verify-and-pr.md` — CI checks, `gh` commands, what the bot review looks like
  and how the #1034 replies were written.
