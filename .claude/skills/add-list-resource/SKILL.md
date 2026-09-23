---
name: add-list-resource
description: "Add a `terraform query` list resource, plus the resource identity it requires, to a terraform-provider-ionoscloud resource — code, unit test, docs, CHANGELOG, verification and PR review. Use when asked to make a resource queryable, add a list resource or a resource identity, support `terraform query` for an `ionoscloud_*` type, or repeat PR 1034's work on `ionoscloud_datacenter`."
user-invocable: true
argument-hint: <ionoscloud_type>, e.g. ionoscloud_ipblock
---

# Add a list resource + resource identity

Target: **$ARGUMENTS** (if empty, ask which `ionoscloud_*` type).

Reproduces PR #1034's work for `ionoscloud_datacenter` (reference:
`ionoscloud/resource_datacenter_list.go`). A list resource with no identity, or a test that passes
when the mapper is wrong, is not done — real #1034 findings.

**Where the arc ends.** The run ends with the work **uncommitted in the working tree** — nothing
staged, no branch, no commit, no push, no `gh` write. Report what changed (`git status --short` plus
the verification results), then **ask the user what else, if anything, to do**. Name the
candidates: commit, and with what subject; push the branch and open the PR; run rung 5
(`go vet -tags=all ./...`) if you left it to CI; run the tagged acceptance test or a live
`terraform query` check — both spend real IONOS credentials and create real cloud resources, never
yours to start; carry on to another resource.

Then wait. `git commit`, `git push`, `gh pr create` and replying to a review comment are writes the
maintainer makes, or explicitly asks you to make — never a step you take because the checklist has
one line left; `references/verify-and-pr.md` documents *how*, for when you are asked, and is not
permission.

---

## Machine limits — read first

12-core laptop in use while you work; a fan-out once hit load 55 (`references/why.md`).

- **Do not fan out to subagents**, and never build in parallel or in the background.
- **Never** `make testacc` or any `go test` with a build tag: `GNUmakefile:22-23` sets `TF_ACC=1` and
  `-tags`, so those suites use real credentials and create live cloud resources.
- **Avoid** `go build ./...`, `go test ./...`, `make test`: `TEST` is all 60 non-vendor packages at
  `-parallel=4` — load, not cloud. An **untagged, package-scoped `go test` is always safe**; if it
  drags, stop.

---

## Step 0 — Discovery, and six gates

Everything downstream depends on one branch: **SDKv2 or framework-native?** Settle it and the gates
below first; `scripts/probe.sh discover <ionoscloud_type> <Const> <suffix> <resource>` answers their
seven numbered questions.

**Gate 1 — which branch, or none.** (2) not (3) → **SDKv2-backed**, the harder path. (3) from a
`resource_*.go`, not (2) → **framework-native**: far less work, more methods on the *same* struct. Only `[Dd]ataSource` in (2), or `data_source_*.go`/`ephemeral_*.go` in (3) → **no
managed resource**, nothing to list. Both → ask; neither → no such type.

**Gate 2 — a parentless collection GET?** `scripts/probe.sh listcalls` prints every one in the
vendored SDKs. None for yours → **stop**.

**Gate 3 — which of three client shapes (6)?** `.NewCloudAPIClient(ctx, location)` or
`...WithFailover(ctx)` → the **Cloud API** (`sdk-go/v6`), templates as written. A plain **field**
(`.DNSClient`, `.NFSClient`, …) → an **sdk-go-bundle product** whose client is already on the bundle
(deltas: `sdkv2-branch.md` §2c). A **constructor taking a location** (`.NewMongoClient`) → a product
built per location — decision 4's answer.

**Gate 4 — the state writer. (7) is a gate, not background reading.** Two questions, the second the
one everybody skips: *can the mapper call this writer*, and *does it fill every **Required**
attribute of the schema from the API object* — `-generate-config-out` yields usable blocks only when
it does. Diff its `d.Set` calls against the schema's `Required: true` entries; a gap (a parent id, a
write-only field) goes in the docs, and parent ids are set in the mapper (§2c).

Find it by **following Read**, never by grepping for a name. **The name is not the contract, the
signature is:** the mapper calls any writer shaped `func(d *schema.ResourceData, obj <sdk>.<X>)
error` — spelling, receiver and package do not matter.

- `setDatacenterData` (`ionoscloud/resource_datacenter.go:357`) — unexported, package-level, by
  pointer.
- `IpBlockSetData` (`ionoscloud/resource_ipblock.go:259`) — exported, package-level, called directly.
- `SetZoneData` (`services/dns/zone.go:60`) — a **method** on the product's service client, by value:
  `r.bundle.DNSClient.SetZoneData(data, zone)`. It also calls `d.SetId`, which not every writer does
  — check: the identity setter reads the id back out.

`set<X>Data` is a placeholder, not a name to grep for. A writer needing a `context.Context`, an API
client for further calls, or an object the collection GET does not return ends the job: stop and say
so rather than hand-map attributes.

**Gate 5 — does a schema attribute name a parent?** A parentless GET does not mean no parent:
`ionoscloud_dns_record` passes gate 2 on a cross-zone `RecordsGet(ctx)`, yet `zone_id` belongs in its
identity — the child-resource path (§1a, §2c), where the **mapper** sets the parent id because
`SetRecordData` (`services/dns/record.go`) does not.

**Gate 6 — a bad candidate?** A collection GET needing a parent ID (`DatacentersLansGet`,
`K8sNodepoolsGet`) is an unbounded N+1 over every parent, with no precedent here; association-only
resources (`ionoscloud_ipfailover`, `ionoscloud_datacenter_nsg_selection`) have no API object. Stop.

**Only then open the branch reference** (`sdkv2-branch.md` or `framework-native-branch.md`) — long
files, and the gates stop a doomed run first.

---

## Step 0.5 — Four decisions for the user, before any code

Agreed up front, not discovered in review. Present them together, briefly, with your
recommendation, then wait. Evidence: `references/decisions-evidence.md`.

**1. Pagination.** Decide explicitly. Provider-wide there is exactly **one** offset/limit loop
(`services/dbaas/pgsqlv2/cluster.go`; it stops on a short page because `Links.HasNext()` stays
populated past the last page and spins forever), S3 uses continuation tokens, everything else is
unpaginated. On #1034 the user **declined** pagination and documented the limitation instead: fixing
only the list resource would leave it inconsistent with the identical unpaginated call in
`ionoscloud/data_source_datacenter.go`. That precedent is reusable — but then **the docs must say
so** (wording in `references/docs-and-changelog.md`), and Copilot will raise pagination on the PR.

**2. Identity attributes** — the tuple the import ID parses into (`splitImportID`,
`ionoscloud/utils.go`), never a mutable one such as `name`: Terraform fails refresh when an identity
attribute changes. **3. Filter fields** — the `identity.FilterAttribute(...)` allow-list must be
exactly the key set of the mapper's `identity.MatchesFilters(...)` map, or a filter silently matches
nothing; two string attributes that really narrow.

**4. Which client, regional or global.** Cloud API → always **global**: one
`NewCloudAPIClientWithFailover(ctx)` call, as in `datacenterListResource.List` — those collections
are not partitioned, and fanning out returns every item N times. A bundle product → the resource's
own client; settle whether *this resource* is partitioned and its locations enumerable, else ask.

*A known hole, not yours to fix.* `NewCloudAPIClientWithFailover` resolves its endpoint from
**global** cloud overrides only (`FilterGlobalOverrides` keeps entries whose `location` is empty —
`services/bundleclient/bundleclient.go:443`), so a file config overriding the cloud product per
location with no global entry hard-errors `no global failover endpoints configured for "cloud"`,
where the resource's own `NewCloudAPIClient(ctx, location)` is fine. `ionoscloud_datacenter` has the
same hole: follow it anyway.

---

## The plan

1 identity, 2 list resource + registration, 3 unit test + mutation check, 4 docs + CHANGELOG,
5 ladder, 6 stop and ask. Step 1 gates step 2: without an `Identity` nothing can be listed, and the
only trace is a `tflog.Error` line.

The branch reference walks the files. On the SDKv2 branch they all sit in package `ionoscloud` beside
the resource — forced by an import cycle (`references/why.md`), so framework code there is not a
mistake — and the new `_list_test.go` (package `ionoscloud_test`) reuses the helpers in
`ionoscloud/resource_datacenter_list_test.go`. `internal/framework/provider/provider.go` is **not**
touched, nothing new belongs in `internal/framework/identity/`.

---

## Verification ladder

Cheapest first, one at a time; `scripts/probe.sh rung0` … `rung5-full` runs them,
`references/why.md` says why.

0. **`gofmt`** after every edit — never with an empty argument list (it reads stdin and passes).
1. **Build the packages you touched** — main loop. SDKv2: `./ionoscloud/` covers registration too;
   framework-native: also `./internal/framework/provider/` if the service package is new.

**Rung 1 can fail for a reason only rung 4 fixes — the one permitted jump in the ladder.** The repo
vendors its dependencies, so an import of a not-yet-vendored subpackage breaks *this* rung — a
cannot-find-module / missing-package error naming that import path. Run rung 4, then re-run rung 1.
2. **The untagged unit test** — main loop, no credentials, no network.
2b. **`TestProvider$` (SDKv2)** — the **only** thing that validates your new `Identity`; rung 2
   passes for a malformed one.
3. **Lint your own lines** — keep `--new-from-rev` and the path argument.
4. **Only if imports changed** — the vendor half even when tidy is clean; **report `vendor/`** (#1034).
5. **Expensive, once, at the end — ask first**: the only rung that compiles your tagged
   `Query: true` step. Offer `rung5-scoped`.

**No test job runs on PRs** — rungs 2 and 2b are on you. The ladder tests the ListResource *RPC*,
never `terraform query`: that needs `dev_overrides` and real credentials, **only on explicit
request** (`docs-and-changelog.md` §2b).

---

## Definition of done

- [ ] `Identity` on the managed resource, every attribute set on **every** read path
- [ ] `d.SetId(...)` right after parsing the import parts (SDKv2); the legacy string import still
      works unchanged
- [ ] `git diff ionoscloud/resource_<resource>.go` read line by line: every `-` line is one the
      identity forced. No `d.SetId` dropped as redundant, no re-worded comment, no drive-by tidy
- [ ] List resource registered in a `ListResources()` — `ionoscloud/list_resources.go` (SDKv2), or
      the `<service>` package's plus a `<service>.ListResources(),` provider line **only** when that
      package is new
- [ ] **No model, no attribute mapping of your own** (SDKv2): `map<X>` does the filter check and
      nothing else — every result is `r.resourceSchema.Data(&terraform.InstanceState{})` →
      `set<X>Data` → `set<X>Identity` → `identity.MappedItemFromResourceData`, in that order. No
      `tfsdk`-tagged struct
- [ ] Unit test asserts every mapped attribute and stub result, and null vs zero (an omitted nested
      **block** decodes to `[]any{}`, not nil); **the mutation check ran and failed**
- [ ] One filter subtest **per field** in the allow-list, plus an AND case; the stub asserts the
      fetch's query params, absences included
- [ ] A `Query: true` step on the tagged acceptance test — a **separate** `TestAcc<Resource>Query`
      with its own fixture, never extra steps on the existing one (`verify-and-pr.md` §4); written,
      **not run**; plus a non-ForceNew update step if Update writes the identity
- [ ] `docs/list-resources/<resource>.md` + the pointer section on `docs/resources/<resource>.md`
- [ ] CHANGELOG entry under the correct `## X.Y.Z` heading — decided by the **git tag, not the
      heading**: append under the topmost heading while no `v<that version>` tag exists; once it is
      tagged, open a new `## X.Y.Z+1` heading above it. The check and the #1034
      counter-example: `references/docs-and-changelog.md` §4 — re-run it just before merge
- [ ] Ladder rungs 0–5 clean, **2b included**
- [ ] The pagination decision written where a reviewer will find it, with this endpoint's real
      default limit — or, when the client sends none (the bundle norm) or it does not page, say so.
      Never a number you did not read out of the vendored client
- [ ] Left uncommitted: nothing staged, nothing committed, no `git commit`, `git push` or `gh`
      write — unless the user asked for it in so many words
- [ ] Reported what changed, then **asked what to do next** — the run ends on that question, not on
      a silent stop or a step you chose yourself

## Where to look things up

Paths are relative to this skill's directory. **Run the scripts, do not read them** — run one bare
for its subcommands and arguments; every check ends in a verdict line saying what "no match" means.
(`scripts/_lib.sh` is shared code, sourced by the rest, never run on its own.)

- `scripts/probe.sh` — step 0 (`discover`, `listcalls`), step 0.5 limits and locations, the ladder
  (`rung0` … `rung5-full`).
- `scripts/sdkv2.sh` — steps 1–2 on the SDKv2 branch: `identities`, `depth`, `writer`.
- `scripts/docs.sh` — step 4; `changelog-heading` is the tag authority.
- `scripts/test.sh` — step 3: which shared test helpers exist and where.
- `scripts/review.sh` — `test-unit`, `queries` during the run; `pr-create`, `review`,
  `threads`, `reply` only once the user asks.

- `references/decisions-evidence.md` — step 0.5, while the four decisions are open.
- `references/sdkv2-branch.md` — SDKv2 branch, after the gates. **The traps live here.**
- `references/framework-native-branch.md` — the other branch.
- `references/test-harness.md` — step 3: fixture, stub, subtests, mutation check.
- `references/docs-and-changelog.md` — step 4: the two doc pages, the CHANGELOG heading rule.
- `references/verify-and-pr.md` — PR and review replies; on request only.
- `references/why.md` — the rationale and the dead ends; not opened on a normal run.
