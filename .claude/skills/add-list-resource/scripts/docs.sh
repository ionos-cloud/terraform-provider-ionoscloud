#!/usr/bin/env bash
. "$(dirname "$0")/_lib.sh"
root

case "${1:-}" in

changelog-heading)
  need CHANGELOG.md
  top=$(sed -n 's/^## //p' CHANGELOG.md | head -1 | awk '{print $1}')
  [ -n "$top" ] || die "no '## X.Y.Z' heading in CHANGELOG.md"
  havegit "the tag check is NOT RUN: report it, never infer it from the heading"
  if git rev-parse -q --verify "refs/tags/v$top" >/dev/null; then
    printf 'VERDICT: v%s is ALREADY TAGGED - open a new "## <next>" heading above it.\n' "$top"
  else
    printf 'VERDICT: v%s is untagged, i.e. unreleased - append your bullets under it.\n' "$top"
  fi
  ;;

limit-cloudapi|limit-bundle) limit_probe "${2:?operation, e.g. DatacentersGet}" ;;

limit-datasource) exec "$SCRIPTS/probe.sh" datasource-limit "${2:?resource}" ;;

# forcenew - read a CustomizeDiff hit: an immutable-field check errors, it does not replace.
forcenew)
  f="ionoscloud/resource_${2:?resource}.go"
  need "$f"
  grep -nE 'ForceNew|CustomizeDiff' "$f" |
    q "force-new / CustomizeDiff in $f" "no force-new attribute: the destroy paragraph does not apply"
  ;;

*) echo "usage: $0 {changelog-heading | limit-cloudapi <Op> | limit-bundle <Op> | limit-datasource <resource> | forcenew <resource>}" >&2; exit 2 ;;
esac
