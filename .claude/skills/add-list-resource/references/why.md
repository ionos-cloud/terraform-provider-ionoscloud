# Why the skill says what it says

Never opened on a normal run; it answers "why not X" for a reviewer, or anyone about to delete
something redundant-looking. Citations are `master`. Unqualified `core_schema.go`,
`resource_data.go` and `resource.go` are under
`vendor/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema/`, `sdkv2.go` and `list.go` under
`internal/framework/identity/`, `dc_list_test.go` is `ionoscloud/resource_datacenter_list_test.go`.

1. No Go model · 2. `package ionoscloud` · 3. The state-writer gate · 4. `Data()`, `SetId`,
`nullTimeouts` · 5. `[]any{}` vs nil · 6. `//nolint`, the fan-out · 7. The unit test ·
8. Verification · 9. CI, lint, review · 10. Docs

## 1. No Go model of the schema

Type and value are two renderings of one `configschema.Block` (`core_schema.go:407-416`,
`resource_data.go:84-163`); a struct is a third, kept equal by hand. #1034 shipped
`datacenterResourceModel`, a concurrent PR added `cpu_architecture.enabled_features`, and master then
failed every `terraform query` on `ionoscloud_datacenter`: *"Object defines fields not found in
struct: enabled_features"*.

## 2. Why it lives in `package ionoscloud`

It needs the resource's unexported `resource<X>()`, writer and `set<X>Identity()`, and the import
will not reverse: `ionoscloud/provider_test.go` is an *in-package* test importing
`internal/framework/provider`, so `internal/framework/...` → `ionoscloud` is a cycle in the test
build. The provider still needs no edit: `New(sdkv2ListResources ...func() list.ListResource)` takes
the whole slice (`internal/framework/provider/provider.go:64-68`), and every call site already passes
`ionoscloud.ListResources()...`.

## 3. Why the state-writer check is a gate

`-generate-config-out` yields usable blocks only if the writer fills every `Required: true` attribute
from the API object, and the mapper hands it one already-fetched element (`list.go`). A writer
needing a `ctx`, a client for further calls, or an object the collection GET omits ends the job:
`setResourceServerData` takes a `ctx` and an `APIClient`, `setBackupUnitData` a `Contracts`.
Step 1 comes first for the same reason: with no `Identity`, `SetRawV6Schemas`
gives up, and `RawV6Schemas` has no diagnostics channel — the only trace is a `tflog.Error` line.

## 4. `Data()`, not `TestResourceData()`; `d.SetId`; `nullTimeouts`

`Data(&terraform.InstanceState{})` alone builds the `ResourceData` with an identity
(`resource.go:1405-1426`) and copies `result.timeouts = r.Timeouts` (`resource.go:1414-1418`), the
only thing that materialises a `timeouts` block; `TestResourceData` omits a block the type declares,
and `Set` rejects on type. Both conversions go through `ResourceData.State()`, nil while
the ID is empty (`resource_data.go:425-434`) — hence `d.SetId`. `nullTimeouts` (`sdkv2.go:99-121`)
then nulls the block back out, as SDKv2's own `ReadResource` does.

## 5. Why an absent nested block decodes as `[]any{}`

`toproto6.DynamicValue` runs `ReifyNullCollectionBlocks`
(`toproto6/dynamic_value.go:29-30`): a null list/set **block** becomes an
**empty** collection, and `core_schema.go`, not the API, decides which a nested schema is. `Computed`
only → *attribute* (`:102-106`), absent → `assert.Nil`; `Optional` (± `Computed`) → *block*
(`:108-112`, `:201-205`), absent → `[]any{}`. `MaxItems: 1` changes nothing — `:201-205` sets
`Nesting` from `s.Type` before reading `MaxItems` — though that is extrapolation, none having
shipped.

## 6. No bare `//nolint:errcheck`; what the fan-out is for

`Config.GetAttribute` returns `diag.Diagnostics`, not an `error` (`vendor/.../tfsdk/config.go:34`), so
`errcheck` has nothing to report, and a new bare directive fails `nolintlint` `require-explanation`
where the pgsqlv2 one predates the rule. Its fan-out narrows `AvailableLocations()` by the
`location` filter only to save round trips: the
mapper re-checks `MatchesFilters` regardless, a server-side filter being no guarantee of the compare
the docs promise. `AvailableLocations()` exists nowhere outside `services/dbaas/`, so an
SDKv2-backed resource has no fan-out.

## 7. The unit test

Untagged, so it needs no `TF_ACC` or credentials (105 of `ionoscloud/`'s 113 test files do) and stays
in golangci-lint's scope; external, so it runs through the surface `main.go` uses.
What it pins, with no model, is *the state the writer produces surviving the round trip through
`TfTypeResourceState` into the list-result schema* (`dc_list_test.go:26-35`), so its per-attribute
asserts **replaced the deleted model** — a value under the wrong key has no
compiler and no `tfsdk` tag to catch it. The rest catch silent failures: `ListResourceSchemas` an
unregistered type, `GetResourceIdentitySchemas` a missing `Identity`,
`assert.Nil(t, resource["timeouts"])` (`dc_list_test.go:105`) `nullTimeouts`.

Traps:

- **`int64`, never `float64`**: `testifylint` `float-compare` rejects `assert.Equal` on a float and is
  not excluded for test files; it cost a red CI round on #1034. The `IsInt()` guard
  (`dc_list_test.go:366`) is a tripwire, not scaffolding; `goValue` (`dc_list_test.go:342`) decodes
  numbers as `int64` and has no `Map`/`Tuple` case.
- **The mutation check is mandatory**: #1034's original test decoded `map[string]string`, skipping
  every non-string attribute, so a writer swapping two ints passed.
- **Bundle stubs must fill every non-omitempty enum field**: bundle models validate enums while
  unmarshalling, so a zero value fails the *fetch* naming neither field nor item — a `dns_zone` stub
  failed on an unset `ZoneRead.Metadata.State`: `is not a valid ProvisioningState`.
- **`IONOS_API_URL` is not enough for a bundle product**: clients read `IONOS_API_URL_<PRODUCT>`,
  and `utils/loadedconfig/loadedconfig.go:50` *defers* `ChangeConfigURL` over it, so a regional
  client with a location set takes a **production** URL. DNS is the counter-example
  (`services/dns/client.go:26`).
- **Assert the query string**: no `.Depth(1)` → every `Properties` nil, writer never reached; no
  `.Limit(...)` → truncation at the SDK client's fallback (100 for `/ipblocks`,
  `sdk-go/v6/api_ip_blocks.go:486-490`). Assert it as a **literal**, the only thing tying the docs
  page to the code.

## 8. Verification

- **One agent; the ban is on tags, not on `go test`.** A previous fan-out drove the maintainer's
  12-core laptop to load 55, and `make test` runs all 60 non-vendor packages at `-parallel=4`
  (`GNUmakefile:1`).
- **Rung 0 hard-codes its paths**: `gofmt` with no arguments reads stdin and exits clean, and
  `master...HEAD` is a commit range while this work is never committed, so the list expands to
  nothing and the rung reports green having inspected nothing. CI (`-l -d .`) fails instead.
- **Rung 1 can fail for a rung-4 reason**: the repo vendors (`vendor/modules.txt` → `-mod=vendor`),
  so an import of a not-yet-vendored subpackage breaks the *build*, as a cannot-find-module error.
- **Rung 5 is worth its cost**: 136 `_test.go` files sit behind build tags, invisible to a plain
  `go vet ./...`, including the `Query: true` step just written — which is the only thing evaluating
  `QueryResultChecks` (`terraform-plugin-testing/helper/resource/testing_new.go:370-373`), so on an
  apply step they never run and it passes. Nothing in the ladder reaches `terraform query`.

## 9. CI, lint and review

- **CI gates compilation, not behaviour.** `build.yml` runs `go vet -tags=all ./...`, loading each
  package *with its test files*, so a non-compiling test turns CI red. **No `pull_request` workflow
  runs `go test`** — the test workflows are
  `workflow_dispatch`/`schedule`, so green CI does not mean the tests pass — and
  `golangci-lint --new-from-rev=origin/<base>` means lint sees the diff, never the whole file.
- **Linters**: `_test.go` is exempt from `dupl`/`errcheck`/`gocyclo`/`gosec`/`unparam`/`unused`, not
  from `testifylint`, `revive`, `errorlint`, `goconst`, `prealloc`; `testifylint` alone has forced a
  commit here (§7). `revive` `import-alias-naming` wants `^[a-z][a-z0-9]*$`, hence
  `fwidentity` over `fw_provider`, and `goconst` fires at three occurrences, hence the type-name
  const.
- **The bot**: pagination is the finding a list-resource PR gets from Copilot, twice on #1034, and
  its false positives are Go-semantics claims wrong for this module's Go version, or helpers that do
  not exist. On scope, cite the precedent: `ionoscloud/data_source_datacenter.go:147`
  issues the identical unpaginated call, and the one offset/limit loop provider-wide is
  `services/dbaas/pgsqlv2/cluster.go:20-49`.
- **`/test`'s two halves disagree**: `pr-comment-trigger.yaml` gates the step on a plain
  `contains(comment.body, '/test')`, its script on `/^\/test\s+(\S+)/` — so a mid-sentence mention
  dispatches nothing but still runs a job that fails.

## 10. Docs

- **"The API's default page limit" is banned**: the fallback is sent by the *SDK client*, not by the
  endpoint, and an explicit `.Limit(n)` means it bounds nothing anyway. A doc comment's default
  describes the endpoint only where nothing sends a limit: the one naming "the first 100 items" sits
  on `DatacentersGet`, whose client sends 1000 (`sdk-go/v6/api_data_centers.go:489`), so there it
  describes nothing. Say the limit the code sends; where none is sent, the default that operation's
  own doc comment gives, or the bound with no figure. Inventing one is worse — and hedging one is
  not on the table (`docs-and-changelog.md` §2a).
  `docs/list-resources/datacenter.md:68-72` carries the banned
  phrasing and is why the rule exists.
- **Known divergence, not rules**: `docs/list-resources/s3_bucket.md` lacks the two-filter example
  and `## Identity Attributes`; `docs/resources/pg_cluster_v2.md:191` links to a `psql_cluster_v2.md`
  that does not exist — hence checking the back-link, not assuming it.
- **No version string exists in the source**: `main.go`, `GNUmakefile` and `.goreleaser.yml` carry
  none, the `CHANGELOG.md` heading *is* the version, and releases are cut by pushing a `v*` tag —
  hence the git-tag check, not a file to bump.