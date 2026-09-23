#!/usr/bin/env bash
# Docs + CHANGELOG probes. Run it, do not read it.
. "$(dirname "$0")/_lib.sh"
root

case "${1:-}" in

# changelog-heading - append under the topmost "## X.Y.Z", or open a new one? THE TAG IS
# THE AUTHORITY: an un-bumped heading is not evidence of being unreleased (#1034 shipped
# bullets under a version tagged while the PR was open). Do what this prints, either
# way, and re-run it right before merge.
changelog-heading)
  need CHANGELOG.md
  top=$(sed -n 's/^## //p' CHANGELOG.md | head -1 | awk '{print $1}')
  [ -n "$top" ] || die "no '## X.Y.Z' heading in CHANGELOG.md"
  if git rev-parse -q --verify "refs/tags/v$top" >/dev/null; then
    printf 'VERDICT: v%s is ALREADY TAGGED - open a new "## <next>" heading above it.\n' "$top"
  else
    printf 'VERDICT: v%s is untagged, i.e. unreleased - append your bullets under it.\n' "$top"
  fi
  ;;

# limit-cloudapi <Op> | limit-bundle <Op> - what limit does the client send with no
# .Limit(n)? One probe for both SDKs: it reads the emission, so it cannot answer
# "nothing" merely because a product generates params differently. Op: DatacentersGet...
limit-cloudapi|limit-bundle) limit_probe "${2:?operation, e.g. DatacentersGet}" ;;

# limit-datasource <resource> - does the sibling data source pass an explicit .Limit(n)?
# No match = it sends none, so do not claim a number on the list page.
limit-datasource) exec "$SCRIPTS/probe.sh" datasource-limit "${2:?resource}" ;;

# forcenew <resource> - a force-new attribute at all? NO MATCH = none, so the
# label-collision warning's destroy paragraph and HCL are false: replace them with the
# in-place damage to the wrong object, do not just delete them. A CustomizeDiff hit is
# not automatically a substitute - an immutable-field check errors out instead of
# planning a replacement. Read the hit before believing it.
forcenew)
  f="ionoscloud/resource_${2:?resource}.go"
  need "$f"
  grep -nE 'ForceNew|CustomizeDiff' "$f" |
    q "force-new / CustomizeDiff in $f" "no force-new attribute: the destroy paragraph does not apply"
  ;;

*) echo "usage: $0 {changelog-heading | limit-cloudapi <Op> | limit-bundle <Op> | limit-datasource <resource> | forcenew <resource>}" >&2; exit 2 ;;
esac
