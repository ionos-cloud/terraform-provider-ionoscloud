#!/usr/bin/env bash
. "$(dirname "$0")/_lib.sh"
root
REPO=ionos-cloud/terraform-provider-ionoscloud

case "${1:-}" in

test-unit)
  go test ./ionoscloud/ -run "Test${2:?Resource, e.g. DNSZone}ListResource" -count=1
  go test ./ionoscloud/ -run 'TestProvider$' -count=1
  ;;

# queries - framework-native services fold query steps into lifecycle tests.
queries)
  need ionoscloud internal/framework/services
  grep -l 'TestAcc.*Query' ionoscloud/*_test.go 2>/dev/null |
    q "SDKv2 query acceptance tests" "none in ionoscloud/ yet"
  grep -rl 'Query: *true' internal/framework/services/ 2>/dev/null |
    q "framework-native query steps" "none under internal/framework/services/"
  ;;

# pr-create - gh does not apply the PR template: pass the body.
pr-create)
  need "${3:?path to the PR body}"
  gh pr create --base master \
    --title "feat: add ${2:?ionoscloud_<type>} list resource and resource identity" \
    --body-file "$3"
  ;;

review)
  n=${2:?PR number}
  echo "--- 1. inline review comments (the id that reply takes) ---"
  gh api "repos/$REPO/pulls/$n/comments" --paginate --jq '.[] | {id, in_reply_to_id, user: .user.login, path, line, body}'
  echo "--- 2. review bodies: summaries and suppressed findings ---"
  gh api "repos/$REPO/pulls/$n/reviews" --paginate --jq '.[] | {id, user: .user.login, state, submitted_at, body}'
  echo "--- 3. top-level conversation (Sonar gate, humans) ---"
  gh api "repos/$REPO/issues/$n/comments" --paginate --jq '.[] | {id, user: .user.login, body}'
  ;;

# threads - GraphQL: the REST comments API has no isResolved.
threads)
  n=${2:?PR number}
  gh api graphql -f query="{repository(owner: \"ionos-cloud\", name: \"terraform-provider-ionoscloud\") {
    pullRequest(number: $n) { reviewThreads(first: 50) { nodes {
      id isResolved isOutdated path
      comments(first: 10) { nodes { author { login } databaseId body } } } } } } }"
  ;;

# reply - the id is endpoint 1's numeric one, not a review id.
reply)
  need "${4:?path to the reply body}"
  gh api --method POST -H "Accept: application/vnd.github+json" \
    "repos/$REPO/pulls/${2:?PR number}/comments/${3:?inline comment id}/replies" -f body="$(cat "$4")"
  ;;

*) echo "usage: $0 {test-unit <Resource> | queries | pr-create <type> <body> | review <PR> | threads <PR> | reply <PR> <id> <body>}" >&2; exit 2 ;;
esac
