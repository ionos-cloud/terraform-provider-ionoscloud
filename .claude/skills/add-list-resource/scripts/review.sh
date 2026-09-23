#!/usr/bin/env bash
# Unit-test rungs, and the gh calls for opening and answering a PR.
# pr-create and reply are PUBLIC WRITES on the org repo, and opening a PR fires CI:
# run them only when the user asked for that in so many words.
. "$(dirname "$0")/_lib.sh"
root
REPO=ionos-cloud/terraform-provider-ionoscloud

case "${1:-}" in

# test-unit <Resource> - rungs 2 and 2b: no credentials, no network. No pull_request
# workflow runs `go test` at all, so a failure here is one CI will never catch, and 2b
# is the only cheap thing that runs InternalIdentityValidate on your Identity.
test-unit)
  go test ./ionoscloud/ -run "Test${2:?Resource, e.g. DNSZone}ListResource" -count=1
  go test ./ionoscloud/ -run 'TestProvider$' -count=1
  ;;

# queries - which resources already ship a `terraform query` acceptance test?
# No match in ionoscloud/ => none there; framework-native services fold the query steps
# into their lifecycle tests instead, which is the second listing.
queries)
  need ionoscloud internal/framework/services
  grep -l 'TestAcc.*Query' ionoscloud/*_test.go 2>/dev/null |
    q "SDKv2 query acceptance tests" "none in ionoscloud/ yet"
  grep -rl 'Query: *true' internal/framework/services/ 2>/dev/null |
    q "framework-native query steps" "none under internal/framework/services/"
  ;;

# pr-create <ionoscloud_type> <body-file> - PUBLIC WRITE. gh does NOT apply
# .github/pull_request_template.md; pass the body yourself. Title matches #1034 so the
# squash subject is right.
pr-create)
  need "${3:?path to the PR body}"
  gh pr create --base master \
    --title "feat: add ${2:?ionoscloud_<type>} list resource and resource identity" \
    --body-file "$3"
  ;;

# review <PR> - what the review actually said. All three endpoints are needed: `gh pr
# view` shows no inline comments, and (2) is the only place the bot's <details>Suppressed
# comments</details> findings appear. Nothing on one endpoint is not "there was no
# review". The bot's login differs per endpoint (copilot-pull-request-reviewer[bot] in
# /reviews, Copilot in /pulls/N/comments).
review)
  n=${2:?PR number}
  echo "--- 1. inline review comments (the id that reply takes) ---"
  gh api "repos/$REPO/pulls/$n/comments" --paginate --jq '.[] | {id, in_reply_to_id, user: .user.login, path, line, body}'
  echo "--- 2. review bodies: summaries and suppressed findings ---"
  gh api "repos/$REPO/pulls/$n/reviews" --paginate --jq '.[] | {id, user: .user.login, state, submitted_at, body}'
  echo "--- 3. top-level conversation (Sonar gate, humans) ---"
  gh api "repos/$REPO/issues/$n/comments" --paginate --jq '.[] | {id, user: .user.login, body}'
  ;;

# threads <PR> - which threads are still unresolved? The REST comments API has no
# isResolved, so this needs GraphQL. All must be resolved before merge, but resolving is
# a write on someone else's review: only with explicit approval.
threads)
  n=${2:?PR number}
  gh api graphql -f query="{repository(owner: \"ionos-cloud\", name: \"terraform-provider-ionoscloud\") {
    pullRequest(number: $n) { reviewThreads(first: 50) { nodes {
      id isResolved isOutdated path
      comments(first: 10) { nodes { author { login } databaseId body } } } } } } }"
  ;;

# reply <PR> <inline-comment-id> <body-file> - PUBLIC WRITE. There is no gh pr
# subcommand for replying to an inline comment. The id is the numeric one from endpoint
# 1, NOT a review id. The body comes from a file so shell quoting cannot damage it.
reply)
  need "${4:?path to the reply body}"
  gh api --method POST -H "Accept: application/vnd.github+json" \
    "repos/$REPO/pulls/${2:?PR number}/comments/${3:?inline comment id}/replies" -f body="$(cat "$4")"
  ;;

*) echo "usage: $0 {test-unit <Resource> | queries | pr-create <type> <body> | review <PR> | threads <PR> | reply <PR> <id> <body>}" >&2; exit 2 ;;
esac
