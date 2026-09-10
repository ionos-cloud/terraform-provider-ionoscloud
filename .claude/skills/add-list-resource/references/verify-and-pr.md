# Verification, PR, and answering the review

## 1. What CI actually runs on a PR

Eight checks, from three workflows that fire on `pull_request`. A fourth (`PR Comment Trigger`
→ `E2E Tests`) exists but needs a human and real credentials — §5.

**Compile-level breakage is gated. Behaviour is not.** `.github/workflows/build.yml` runs
`go vet -tags=all ./...` on every pull request and every push to master, and `go vet` loads each
package *together with its test files* — so a list-resource unit test that does not compile
turns CI red, and so does a tagged acceptance test, since almost every tagged file in the repo
opts into `all` (the three exceptions are `//go:build alb`, `natgateway` and `nlb`, which that
job never compiles).

But **no `pull_request`-triggered workflow runs `go test`.** `test.yml` and every
`*-test-run.yml` are `workflow_dispatch`/`schedule` only, and `e2e.yaml` is `workflow_dispatch`
as well. Nothing in CI *executes* the list-resource test. Green CI means your test compiles, not
that it passes — running it is on you, and neither of these costs credentials, network or more
than a second:

```bash
go test ./ionoscloud/ -run 'Test<Resource>ListResource' -count=1   # ladder rung 2
go test ./ionoscloud/ -run 'TestProvider$' -count=1                # ladder rung 2b
```

| Check | Workflow | Command |
|---|---|---|
| `Run gofmt check` | Sonarcloud and gofmt | `gofmt -l -d` over `.` |
| `Run go vet with all build tags` | Sonarcloud and gofmt | `go vet -tags=all ./...` — compiles test files too |
| `Check go mod tidy and vendor` | Sonarcloud and gofmt | `go mod tidy` then `go mod vendor`, both must leave the tree clean |
| `SonarCloud` / `SonarCloud Code Analysis` | Sonarcloud and gofmt | no local equivalent |
| `detect-noop` | CI golangci-lint | skips lint entirely for a docs-only PR |
| `lint` | CI golangci-lint | `golangci-lint run --timeout 10m0s --verbose --new-from-rev=origin/<base branch>` (v2.10.1) |
| `build` | Broken Link Checker | `markdown-link-check` on modified `docs/**`, `CHANGELOG.md`, `README.md` |

Local ladder: see the SKILL.md **Verification ladder** section. Do not re-derive it.

### Linters most likely to bite

`.golangci.yml` enables ~45 linters under `linters.enable` plus golangci-lint v2's `standard` group
(`errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`). `_test.go` files are excluded
from `dupl`/`errcheck`/`gocyclo`/`gosec`/`unparam`/`unused` — but **not** from `testifylint`,
`revive`, `errorlint`, `goconst`, `prealloc`, `modernize`, `bodyclose` or `noctx`.

1. **`testifylint`** — the only one that has already forced a commit on this work.
   `float-compare` rejects `assert.Equal` on a float. See `test-harness.md`.
   Its other relevant checks: `expected-actual` (expected value first), `require-error`
   (`require`, not `assert`, for error assertions), `len`, `empty`, `nil-compare`.
2. **`revive`**'s configured `import-alias-naming` with `allowRegex: "^[a-z][a-z0-9]*$"` — an
   import alias may contain **no underscore and no uppercase letter**. `ionoscloudsdk`,
   `fwidentity` and `fwprovider` pass; `fw_provider` and `sdkV2` fail.
3. **`goconst`** (min-len 3, min-occurrences 5) — hence the package-level type-name const in
   the test.
4. **`prealloc`** — `make([]T, 0, len(src))`, not `var xs []T` before a range loop.
5. **`gocyclo`** (min-complexity 25) — the list resource has nothing left that grows: `List` is
   a fetch closure plus a call, and the mapper is a filter check, two setter calls and a
   conversion. What can trip it is the resource's own state writer (`setDatacenterData`,
   `IpBlockSetData` — the name is per-resource, and not always unexported) if you extend it while
   you are in there; split into a helper rather than growing one function.
6. **`tagalign`** (`sort: true`, order `json, tfsdk, mapstructure`) — struct tags must be sorted
   and aligned. `` `tfsdk:"id" json:"id"` `` fails. Only the framework-native branch can trip
   it, on the managed resource's own model; the SDKv2 branch declares no tagged struct anywhere.
7. **`nolintlint`** (`require-explanation`, `require-specific`) — a bare `//nolint` is itself
   an error. Two live in the shared test helpers; keep their explanations when you move them.
8. **`errorlint` / `nilerr`** — `%w` and `errors.Is`, never `==` on a sentinel; `nilerr`
   catches `if err != nil { return nil }`, easy to write when swallowing a 404.
9. **Formatters** `gofmt` + `goimports` with
   `local-prefixes: github.com/ionos-cloud/terraform-provider-ionoscloud` — three import
   groups: stdlib, third-party, then this module. The CI **gofmt** job runs plain `gofmt` and
   does *not* check grouping; only golangci-lint's `goimports` does. Fix with
   `golangci-lint fmt ./ionoscloud/...` (add
   `./internal/framework/services/<service>/...` on the framework-native branch).

**`lll` is NOT enabled** (it has a settings block but never appears under `enable`), nor is
`dupl`. Do not wrap lines for a linter that is off.

---

## 2. Branch, commit, PR

Branch: `feat/<resource>-list-resource-identity`. The repo's dominant style is
`<type>/<kebab-slug>` with the same vocabulary as the commit prefix.

Commits: conventional-commit prefixes (`feat:`, `fix:`, `refactor:`, `doc:`, `test:`, `chore:`),
one-line subject, no body. **Every merge to master is a squash**, so the *PR title* becomes the
master commit subject — that is why the checklist's first item is about the PR name.

> **Everything in this section is a write the maintainer owns.** `git commit`, `git push`,
> `gh pr create` and posting a review reply all happen **only on an explicit request** — never
> as the last step of the checklist. This section documents *how*, for when you are asked. It
> is not permission. Opening a PR also fires CI (`.github/workflows/build.yml`,
> `linter.yaml`), so it is not a private action either.

**Do not run `git commit` or `git push` unless explicitly asked.** The maintainer manages their
own commits. When they do ask, use a single-line subject and a Claude co-author trailer naming
**whichever model you actually are** — check the session rather than copying a name from here:

```
<one-line subject>

Co-Authored-By: Claude <Model> <noreply@anthropic.com>
```

PR title, matching #1034:

```
feat: add <ionoscloud_type> list resource and resource identity
```

The template at `.github/pull_request_template.md` is not applied by `gh pr create` — pass it
yourself. Filling in "What does this fix or implement?" is better than shipping the raw
template, though most merged PRs leave it unfilled. The checklist's `label: upcoming release`
is stale — no such label exists.

**Public write on the org repo — only on an explicit request from the user.**

```bash
gh pr create --base master \
  --title "feat: add <ionoscloud_type> list resource and resource identity" \
  --body-file /path/to/pr-body.md
```

### Branch protection consequences

`master` requires **1 human approval** with `dismiss_stale_reviews: true` and
`required_conversation_resolution: true`; no status check is *required*.

- Pushing a commit after an approval **throws the approval away** — batch review fixes into
  one push.
- **Every** review thread, including each Copilot inline thread, must be marked **resolved**
  before merge. Replying is not enough.
- Resolving is a write operation on someone else's review — only with explicit approval.

---

## 3. Reading the review — three endpoints, all needed

`gh pr view` shows **none** of the inline review comments.

```bash
R=ionos-cloud/terraform-provider-ionoscloud

# 1. INLINE review comments — the one gh pr view misses
gh api repos/$R/pulls/<N>/comments --paginate \
  --jq '.[] | {id, in_reply_to_id, user: .user.login, path, line, body}'

# 2. REVIEW BODIES — the summaries AND the <details>Suppressed comments</details>
#    findings that never became inline comments
gh api repos/$R/pulls/<N>/reviews --paginate \
  --jq '.[] | {id, user: .user.login, state, submitted_at, body}'

# 3. TOP-LEVEL conversation comments (Sonar quality gate, human notes, test triggers)
gh api repos/$R/issues/<N>/comments --paginate --jq '.[] | {id, user: .user.login, body}'
```

**Read endpoint 2.** The bot buries findings it decided not to post inline inside a
`<details>Suppressed comments</details>` block in the *review body* — they appear in neither
`gh pr view` nor the comments API. On #1034 the pagination finding first appeared only there.

The bot's login differs per endpoint: `copilot-pull-request-reviewer[bot]` in `/reviews`,
`Copilot` in `/pulls/{n}/comments`. Filter on both.

Unresolved threads need GraphQL (the REST comments API has no `isResolved`):

```bash
gh api graphql -f query='
{
  repository(owner: "ionos-cloud", name: "terraform-provider-ionoscloud") {
    pullRequest(number: <N>) {
      reviewThreads(first: 50) {
        nodes { id isResolved isOutdated path comments(first: 10) { nodes { author { login } databaseId body } } }
      }
    }
  }
}'
```

Reply to an inline comment — no `gh pr` subcommand for this; use the inline comment's numeric
`id` from endpoint 1, not the review id. Write the reply to a scratchpad file to avoid
shell-quoting damage:

```bash
gh api --method POST -H "Accept: application/vnd.github+json" \
  repos/$R/pulls/<N>/comments/<COMMENT_ID>/replies \
  -f body="$(cat /path/to/reply.md)"
```

---

## 4. What the bot review looks like on this repo

- **Copilot reviews automatically on every push**, re-reviewing all files and re-raising
  findings it previously suppressed. Expect 2–4 reviews on a multi-push PR, 1–4 inline
  comments each.
- **Recurring true positives:** docs drifting from code; the same fix not applied to a sibling
  resource; missing regression coverage for the linked issue; vague CHANGELOG wording; a
  leftover `replace` directive in `go.mod`.
- **Recurring false positives:** confident claims about Go semantics that are wrong for this
  module's Go version, and prescriptions to use a helper that does not exist in the package.
  **Verify before complying.**
- **Pagination is the finding a list-resource PR will get.** It was raised twice on #1034,
  once suppressed and once inline. Decide the answer at step 0.5, not in review.
- Human reviewers often reply *inside* the Copilot thread rather than opening a new one.

### The reply shape that worked on #1034

Three parts, every time:

1. **A one-line verdict** — agree / disagree / agree-but-out-of-scope.
2. **Evidence** — `path:line` citations into this repo or `vendor/`, or a named green CI check.
3. **An explicit outcome** — `Leaving as is.` / `No change made.` / what changed instead.

Never a bare "won't fix". Two worked examples:

> *(false positive, refuted with a language fact)* `new(value)` does compile here. As of Go
> 1.26 `new` accepts an expression and returns a pointer to it; this module declares
> `go 1.26.3` and the package builds and tests green. No change made.

Note what that reply deliberately leaves out: #1034's version also argued "there is no `ptr`
helper in this package". True of the repo, but the vendored SDK exports `ionoscloudsdk.ToPtr`
and the test already imports it — so the claim invites a correction and adds nothing. Refute on
the language fact alone.

> *(real finding, declined with a scope argument)* Correct that this reads a single page.
> We have decided not to add pagination in this PR, and to document the limitation instead.
> No datacenter read path in the provider paginates today — `ionoscloud/data_source_datacenter.go:147`
> issues the identical unpaginated call, and provider-wide there is exactly one offset/limit
> loop (`services/dbaas/pgsqlv2/cluster.go:20-49`). So this is a pre-existing provider-wide gap
> rather than something this PR introduces, better fixed uniformly in its own change. What did
> change: `docs/list-resources/datacenter.md` now states … One correction to the detail in the
> comment: the vendored SDK's doc comment says "Default limit is the first 100 items", but the
> generated code for *this* endpoint sends `limit=1000`
> (`vendor/github.com/ionos-cloud/sdk-go/v6/api_data_centers.go:489`); which the server honours
> is unverified.

**Do not reuse that number.** It is per-endpoint — `/ipblocks`, `/targetgroups` and user
management default to 100 — so re-run the grep in `docs-and-changelog.md` §2a for your own
endpoint before writing a number into a public reply. Re-check every `path:line` in a canned
reply before posting it, too; they drift.

Correcting a factual error in the bot's comment, where there is one, is worth doing — it stops
the same wrong number propagating into the docs.

---

## 5. Acceptance tests (optional, costs real cloud resources)

The SDKv2 acceptance test for the resource is tagged
(`//go:build compute || all || <resource>`, e.g. `ionoscloud/resource_datacenter_test.go:1`) and
does not run on a PR. **Never fire it unprompted** — it runs against real IONOS credentials and
creates live resources.

What fires it: `.github/workflows/pr-comment-trigger.yaml` — an `issue_comment` on a PR from an
OWNER/MEMBER/COLLABORATOR/CONTRIBUTOR. Read its two halves separately, because they do not
agree:

- the step's `if` is a plain `contains(github.event.comment.body, '/test')`, so **any**
  top-level PR comment containing that literal anywhere runs the job;
- the script then matches `/^\/test\s+(\S+)/` against the *trimmed* body and `throw`s
  `Comment must be in the format: /test tagname` when it does not match.

So a mid-sentence mention does not dispatch anything, but it does run a job that fails. Only a
comment that *begins* with `/test <tag>` dispatches `e2e.yaml`, which checks out the PR head and
runs `go test ./... -v -timeout 6h -tags "<tag>"` with `TF_ACC: true`, `TF_LOG: DEBUG` and the
real `IONOS_*` secrets. Inline review replies use a different event
(`pull_request_review_comment`) and cannot trigger it either way. **Do not put that literal in a
PR comment at all** — refer to it as "the test-trigger comment".

That dispatch is also the only automated path that would ever execute your untagged unit test —
`-tags` adds to the default build, it does not replace it — but it needs a human, a member-level
author and live cloud resources, so it is not a way to run a unit test. Run rung 2 yourself.

### Writing the query step (deliverable 5)

Writing this is part of the arc; **running it is not**. The reference is
`ionoscloud/resource_datacenter_test.go:133-213` (`TestAccDataCenterQuery`, 92 lines added by
`79ad9715`) — copy its shape, not just the fragments below.
`TestAccIPBlockQuery` (`ionoscloud/resource_ipblock_test.go:144-224`) is the same shape for an
optional-`name` resource, and the only other standalone `TestAcc<Resource>Query` in `ionoscloud/`.
The five framework-native list resources fold the same steps into their lifecycle tests instead —
`internal/framework/services/pgsqlv2/resource_pg_cluster_v2_test.go` is the closest match. An
earlier run of this skill also produced a `target_group` list resource, but that branch (PR #1041)
closed unmerged, so `target_group` has no identity, no list resource and no query test today:
nothing in the tree demonstrates a resource with no `location`.

```go
func TestAcc<Resource>Query(t *testing.T) {
	const (
		<resource>Name = "tf-test-<resource>-query"
		<resource>Addr = constant.<Const> + ".test_<resource>_query"
		otherLocation  = "de/txl" // only for a resource that HAS a location
	)

	// The name and the resource label must be unique to THIS test, not the suite's shared
	// fixture name/label. querycheck.ExpectLength asserts a CONTRACT-WIDE total, so if any
	// other test in the package creates a resource with the same name, ExpectLength(1)
	// flaps. Give the query test its own create step rather than reusing
	// testAccCheck<Resource>ConfigBasic, which is keyed on constant.<Const>TestResource.
	// The zero-result step needs a filter field that can actually discriminate: if the
	// resource has no location, use another allow-listed field whose value the fixture does
	// not use, never a field with one legal value.

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		// `terraform query` and list blocks were introduced in Terraform 1.14.
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheck<Resource>DestroyCheck,
		Steps: []resource.TestStep{
			{Config: /* a plain create step — no Query field */},
			{Query: true, Config: /* list block, no filters */, QueryResultChecks: /* ExpectIdentity */},
			{Query: true, Config: /* list block, filters matching exactly one */, QueryResultChecks: /* ExpectLength 1 */},
			{Query: true, Config: /* same name, a value of a DISCRIMINATING filter field that the fixture does not use */, QueryResultChecks: /* ExpectLength 0 */},
			{ResourceName: <resource>Addr, ImportState: true,
			 ImportStateKind: resource.ImportBlockWithResourceIdentity},
		},
	})
}
```

Four things that are easy to get wrong:

- **`ProtoV6ProviderFactories`, not `ProviderFactories`.** The list resource lives on the
  framework half of the mux; the `ionoscloud/` tests using `testAccProviderFactories`
  cannot see it. Use `testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider)`
  (`ionoscloud/provider_test.go:103`), which builds the mux through
  `ProtoV6ProviderServerFactory` — the same
  `provider.New(ListResources()...)` wiring `main.go` uses (`ionoscloud/provider_test.go:91`).
- **`Query: true` is load-bearing on every list step.** It is what routes the step to
  `testStepNewQuery` (`vendor/.../terraform-plugin-testing/helper/resource/testing_new.go:369-373`
  — note `terraform-plugin-testing`, not the SDK, which has a `helper/resource` of its own with
  no such file), what writes the config
  as a query file instead of a `.tf`, and the *only* thing that evaluates `QueryResultChecks`.
  Grafting a `QueryResultChecks` block onto an ordinary apply step is the silent failure:
  the checks are simply never run and the step passes.
- **Include the zero-result step.** `ExpectLength(addr, 0)` with a filter value that matches
  nothing is what proves the filter is evaluated at all; without it a no-op filter passes.
- `SkipBelow(tfversion.Version1_14_0)` — without it the test fails on any older CLI in CI.

The two fragments worth having in full:

```go
					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectIdentity(<resource>Addr, map[string]knownvalue.Check{
							"id":       knownvalue.NotNull(),
							"location": knownvalue.StringExact("us/las"),
						}),
					},
```

```go
			// Import through the resource identity that the list results carry. This kind
			// already checks that the import succeeds, that the plan it leaves behind is a
			// no-op and that the planned identity matches the one in state; ImportStateVerify
			// cannot be combined with it, only ImportCommandWithID reads that field.
			{
				ResourceName:    <resource>Addr,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
			},
```

A malformed identity schema is caught much earlier and cheaply — but only if you run **ladder
rung 2b**. `ionoscloud/provider_test.go`'s `TestProvider` calls `Provider().InternalValidate()`,
which runs `InternalIdentityValidate()` on every resource in `ResourcesMap` **that declares an
identity** — the call sits behind an `if r.Identity != nil` guard
(`vendor/.../helper/schema/provider.go:220-225`), so this rung catches a *malformed* `Identity`
block, never a *forgotten* one. Nothing else reaches it: CI compiles the test
files (`go vet -tags=all ./...`) but runs none of them, rung 1 compiles no test file at all, and
rung 2's `GetResourceIdentitySchemas` assertion passes for a malformed identity.
`go test ./ionoscloud/ -run 'TestProvider$' -count=1`.
