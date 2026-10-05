# Verification, PR and review

Only when the user asks for a commit, a PR or a review.

## §1 — CI gates compilation, not behaviour

CI's `go vet -tags=all ./...` compiles test files, but **no `pull_request` workflow runs
`go test`**: rungs 2 and 2b (`scripts/review.sh test-unit <Resource>`) are yours. Lint
(`--new-from-rev=origin/<base>`) sees only the diff; traps: `testifylint` (no float
`assert.Equal`), `revive` import alias `^[a-z][a-z0-9]*$`, `goconst` (5 uses), `gocyclo` (25),
`prealloc`, `tagalign`, `nolintlint` (no bare `//nolint`), `errorlint`, `nilerr`, `goimports`
groups (`golangci-lint fmt ./ionoscloud/...`); `lll` and `dupl` are off.

## §2 — Commit and PR

Branch `feat/<resource>-list-resource-identity`; a one-line conventional subject, no body, the
trailer `Co-Authored-By: Claude <Model> <noreply@anthropic.com>` naming **your actual model**. PR
title `feat: add <ionoscloud_type> list resource and resource identity` (`scripts/review.sh
pr-create`). `master`: 1 human approval, **every** thread resolved; a push dismisses approvals, so
batch review fixes. Resolving threads is on request too. **No `/test` literal in any PR comment**,
even mid-sentence (it starts a failing job): say "the test-trigger comment".

## §3 — Reading the review

`scripts/review.sh review <N>` (`gh pr view` omits inline comments); **read endpoint 2, the review
bodies**: the bot's suppressed comments are only there. Bot login:
`copilot-pull-request-reviewer[bot]` in `/reviews`, `Copilot` inline. Reply: a one-line verdict,
evidence (`path:line`, or a named green check), the outcome (`Leaving as is.` or what changed),
never a bare "won't fix". Copilot raises pagination: cite `ionoscloud/data_source_datacenter.go:147`
(the same unpaginated call) and the one offset/limit loop, `services/dbaas/pgsqlv2/cluster.go`. Its
false positives: Go-semantics claims wrong for this module's Go version, helpers that do not exist.
