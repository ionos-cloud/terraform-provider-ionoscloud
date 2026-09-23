#!/usr/bin/env bash
# Shared helpers. Sourced by the other scripts, never run on its own.
# Every check ends in exactly one of:
#   VERDICT: FOUND ...     the lines above are the answer
#   VERDICT: NO MATCH - X  the repo really has none; X says what that means
#   BROKEN: ... (exit 3)   the probe could not look (wrong tree, path gone)
# Silence is never a verdict: a check that cannot fail answers nothing.
set -u
SCRIPTS=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)   # resolved before any cd

die() { printf 'BROKEN: %s\n' "$*" >&2; exit 3; }

# cd to the provider repo - the caller's cwd if that is it, else this skill's own
# checkout - so every path below is repo-relative from wherever you ran this.
root() {
  local d r
  for d in . "$SCRIPTS"; do
    r=$(git -C "$d" rev-parse --show-toplevel 2>/dev/null) || continue
    [ -f "$r/utils/constant/constants.go" ] || continue
    cd "$r" || die "cannot cd to $r"; return
  done
  die "terraform-provider-ionoscloud not found - run from inside it"
}

# need PATH... - a missing path means the probe cannot look, not that the answer is
# empty. A *_list.go firing here means you are on a tree without it.
need() { local p; for p in "$@"; do [ -e "$p" ] || die "no such path: $p"; done; }

# q LABEL NOMATCH_MEANING - reads a check's output on stdin, prints it with a verdict.
q() {
  local out; out=$(cat)
  if [ -n "$out" ]; then
    printf '== %s ==\n%s\nVERDICT: FOUND (%s line(s))\n\n' "$1" "$out" "$(printf '%s\n' "$out" | wc -l)"
  else
    printf '== %s ==\nVERDICT: NO MATCH - %s\n\n' "$1" "$2"
  fi
}

# sdk_client FILE - which client does CRUD take off the bundle? Classified on the
# MATCHED TEXT, never a hit count: resource_target_group.go (Cloud API) and
# resource_dns_zone.go (bundle) both have 5 `bundleclient.SdkBundle).` hits; only the
# symbol after the dot separates them.
sdk_client() {
  need "$1"
  local hits
  hits=$(grep -oh 'bundleclient\.SdkBundle)\.[A-Za-z0-9]*' "$1" | sed 's/.*)\.//' | sort -u)
  printf '== 6. SDK client taken off the bundle in %s ==\n' "$1"
  if [ -z "$hits" ]; then
    printf 'VERDICT: NO MATCH - CRUD takes no client off the bundle here; read meta by hand.\n\n'; return
  fi
  printf '%s\n' "$hits"
  if printf '%s\n' "$hits" | grep -q '^NewCloudAPIClient'; then
    printf 'VERDICT: Cloud API (sdk-go/v6) - the SDKv2 templates apply as written.\n'
    printf '%s\n' "$hits" | grep -q WithFailover &&
      printf '         WithFailover honours the global endpoint override: a fetch that ignores\n         it reads a different server than CRUD (SKILL.md, FilterGlobalOverrides).\n'
  elif printf '%s\n' "$hits" | grep -q '^New'; then
    printf 'VERDICT: a per-location constructor - the product is built per location, so this\n         is the regional decision; do not ship a silently partial listing.\n'
  else
    printf 'VERDICT: an sdk-go-bundle product (a plain field). Deltas: sdkv2-branch.md 2c.\n'
  fi
  printf '\n'
}

# limit_probe OP - how many items does the vendored client ask for on this collection
# when the caller sets no .Limit(n)? OP is the operation: DatacentersGet, ZonesGet...
# The file is resolved here, so sdk-go/v6 vs sdk-go-bundle does not matter: both
# emissions are read. Never read the doc comment - the only sdk-go/v6 comment naming a
# default says "the first 100 items" and sits on DatacentersGet, whose client sends 1000.
# NO MATCH on the file = no such operation is vendored; get the name from `listcalls`.
# The doc comment block immediately above the request builder for <op>.
op_doc_comment() {
  local op=$1 f=$2 ln
  ln=$(grep -n "func (a \*[A-Za-z0-9]*ApiService) ${op}(" "$f" | head -1 | cut -d: -f1)
  [ -z "$ln" ] && return
  sed -n "1,$((ln - 1))p" "$f" | tac | sed -E '/^[[:space:]]*(\/\/|\*|\/\*)/!Q' | tac
}

limit_probe() {
  local op=$1 f body n
  f=$(grep -rl "func (a \*[A-Za-z0-9]*ApiService) ${op}Execute" vendor/github.com/ionos-cloud/ 2>/dev/null | head -1)
  if [ -z "$f" ]; then
    printf '== limit on %s ==\nVERDICT: NO MATCH - nothing vendored has %sExecute; check the name with listcalls.\n\n' "$op" "$op"; return
  fi
  body=$(sed -n "/) ${op}Execute(/,/^}/p" "$f")
  printf '== limit on %s (%s) ==\n' "$op" "$f"
  printf '%s\n' "$body" | grep -n '"limit"\|r\.limit'
  n=$(printf '%s\n' "$body" | grep -oE 'Add\("limit", parameterToString\([0-9]+' | grep -oE '[0-9]+$' | head -1)
  if [ -n "$n" ]; then
    printf 'VERDICT: the client sends limit=%s when the caller sets none. Write %s, never "the\n         API default": it is the SDK client, and .Limit(k) overrides it.\n\n' "$n" "$n"
  elif printf '%s\n' "$body" | grep -q 'r\.limit != nil'; then
    printf 'VERDICT: the parameter exists, but NOTHING is sent unless you call .Limit(n); the\n         page size is then the server default.\n'
    doc=$(op_doc_comment "$op" "$f" | grep -iE 'limit|pagination')
    if [ -n "$doc" ]; then
      printf '         %s doc comment names it - that figure, as the API%ss, no hedge:\n' "$op" "'"
      printf '%s\n' "$doc" | sed 's/^/         /'
      printf '\n'
    else
      printf '         Its doc comment names no default: name the bound, no number, never a hedge.\n\n'
    fi
  else
    printf 'VERDICT: NO MATCH - no limit parameter: no paging, so drop the truncation warning\n         from the docs page. Check Offset the same way.\n\n'
  fi
}
