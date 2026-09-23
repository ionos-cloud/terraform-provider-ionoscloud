# Verification, PR, and answering the review

Why: `references/why.md`. Commands: `scripts/review.sh` (`test-unit`, `queries`,
`pr-create`, `review`, `threads`, `reply`) — run it, don't read it.

## §1 — CI gates compilation, not behaviour

`go vet -tags=all ./...` (`.github/workflows/build.yml`) compiles test files, so a unit or
tagged acceptance test that does not build turns CI red — but **no `pull_request` workflow runs
`go test`**, so green CI means your test compiles, not that it passes. Rungs 2 and 2b are yours:
`scripts/review.sh test-unit <Resource>`. Lint runs with
**`--new-from-rev=origin/<base>`, so it sees only the diff**. Ladder: SKILL.md.

Linters that bite (`.golangci.yml`; why.md gives each rule): `testifylint` (no
`assert.Equal` on floats), `revive` `import-alias-naming` (alias must match `^[a-z][a-z0-9]*$`),
`goconst` (5 occurrences), `prealloc`, `gocyclo` (25), `tagalign`, `nolintlint` (no bare
`//nolint`), `errorlint`/`nilerr`, `goimports` grouping (`golangci-lint fmt ./ionoscloud/...`).
`lll` and `dupl` are **not** enabled; do not wrap lines for them.

## §2 — Where the arc ends — everything below is on request only

**The skill run ends before this section**, with the work **uncommitted in the working tree**
and a question to the user about what to do next. Nothing here happens until that comes back
"commit it" / "open the PR".

**Do not run `git commit` or `git push` unless explicitly asked** — the maintainer manages their
own commits. `gh pr create` and posting or resolving review comments are equally theirs —
public writes, and opening a PR fires CI. This is *how*, not permission.

- Branch `feat/<resource>-list-resource-identity`; conventional-prefix commit, one-line
  subject, no body, plus a `Co-Authored-By: Claude <Model> <noreply@anthropic.com>` trailer
  naming **the model you actually are** — check the session, do not copy a name here.
- PR title, matching #1034: `feat: add <ionoscloud_type> list resource and resource identity`.
  Merges are squashes, so this becomes the master subject.
- `master` needs 1 human approval, dismisses stale reviews, and requires **every** thread
  resolved. Batch review fixes into one push: a push after approval discards it.

## §3 — Reading the review

`gh pr view` shows **none** of the inline comments; `review <N>` hits all three endpoints that
matter. **Read the review bodies** (endpoint 2): the bot buries findings it chose not to post
inline in a `<details>Suppressed comments</details>` block there, invisible to both — on #1034
the pagination finding appeared only there. Its login differs per endpoint,
`copilot-pull-request-reviewer[bot]` in `/reviews` and `Copilot` in `/pulls/<N>/comments`;
filter on both.

Every reply has three parts: a one-line **verdict**, **evidence** (`path:line` into this repo or
`vendor/`, or a named green check), an explicit **outcome** (`Leaving as is.` / what changed
instead). Never a bare "won't fix". Expect pagination — decide that answer at step 0.5, not in
review. Worked replies, recurring false positives: why.md.

## §4 — The acceptance query step (deliverable 5)

**Writing it is part of the arc; running it is not.** Acceptance tests are tagged
(`//go:build compute || all || <resource>`, `ionoscloud/resource_datacenter_test.go:1`), never
run on a PR, and cost real credentials and live resources — never fire one unprompted. **Do
not write the `/test` literal into a PR comment at all** — call it "the test-trigger comment";
a mid-sentence mention starts a job that then fails (why.md).

`TestAccDataCenterQuery`, `ionoscloud/resource_datacenter_test.go:133-212`, is the only query
acceptance test on master: read it for the exact `ExpectIdentity` check, `filters = [...]`
blocks and identity-import step (whose comment says why `ImportStateVerify` cannot join that
kind), then write your own around your resource's fixture — do not copy it.

```go
func TestAcc<Resource>Query(t *testing.T) {
	const (
		<resource>Name = "tf-test-<resource>-query"
		<resource>Addr = constant.<Const> + ".test_<resource>_query"
		otherLocation  = "de/txl" // only for a resource that HAS a location
	)

	// Own name, own label, own create step — NOT testAccCheck<Resource>ConfigBasic:
	// ExpectLength asserts a CONTRACT-WIDE total, so a same-named resource made by
	// another test in the package makes ExpectLength(1) flap.

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		// `terraform query` and list blocks were introduced in Terraform 1.14.
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheck<Resource>DestroyCheck,
		Steps: []resource.TestStep{
			{Config: /* plain create step — no Query field */},
			{Query: true, Config: /* list block, no filters */, QueryResultChecks: /* ExpectIdentity */},
			{Query: true, Config: /* filters matching exactly one */, QueryResultChecks: /* ExpectLength 1 */},
			{Query: true, Config: /* same name + a DISCRIMINATING filter value the fixture lacks */, QueryResultChecks: /* ExpectLength 0 */},
			{ResourceName: <resource>Addr, ImportState: true,
			 ImportStateKind: resource.ImportBlockWithResourceIdentity},
		},
	})
}
```

- **`ProtoV6ProviderFactories`, not `ProviderFactories`** — the list resource is on the
  framework half of the mux, invisible to `testAccProviderFactories`;
  `testAccProtoV6ProviderFactoriesInternal` (`ionoscloud/provider_test.go:103`) builds it.
- **`Query: true` is load-bearing on every list step**: it alone evaluates `QueryResultChecks`.
  Graft them onto an apply step and they silently never run.
- **Keep the zero-result step.** `ExpectLength(addr, 0)` on a value matching nothing proves the
  filter runs; without it a no-op filter passes. Locationless: pick a field with more than one
  legal value.
