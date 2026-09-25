#!/usr/bin/env bash
# Every check ends in VERDICT: FOUND, VERDICT: NO MATCH - <meaning>, or BROKEN; never silence.
set -u
SCRIPTS=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)   # resolved before any cd

die() { printf 'BROKEN: %s\n' "$*" >&2; exit 3; }

# root - cd to the first directory, up from $PWD and then from this script, that holds the
# provider (utils/constant/constants.go and go.mod). Needs no git.
root() {
  local d
  for d in "$PWD" "$SCRIPTS"; do
    while [ -n "$d" ] && [ "$d" != / ]; do
      if [ -f "$d/utils/constant/constants.go" ] && [ -f "$d/go.mod" ]; then
        cd "$d" || die "cannot cd to $d"
        return
      fi
      d=$(dirname "$d")
    done
  done
  die "terraform-provider-ionoscloud not found - run from inside it"
}

# havegit WHY - only the checks that need git call this; without it they stop, never guess.
havegit() { git rev-parse --is-inside-work-tree >/dev/null 2>&1 || die "git unavailable - $1"; }

# need PATH... - a missing path is BROKEN (the probe cannot look), never an empty answer.
need() { local p; for p in "$@"; do [ -e "$p" ] || die "no such path: $p"; done; }

q() {
  local out; out=$(cat)
  if [ -n "$out" ]; then
    printf '== %s ==\n%s\nVERDICT: FOUND (%s line(s))\n\n' "$1" "$out" "$(printf '%s\n' "$out" | wc -l)"
  else
    printf '== %s ==\nVERDICT: NO MATCH - %s\n\n' "$1" "$2"
  fi
}

# sdk_client FILE - classify on the symbol after `SdkBundle).`, never on a hit count.
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

# op_doc_comment OP FILE - the comment block right above OP's request builder.
op_doc_comment() {
  local op=$1 f=$2 ln
  ln=$(grep -n "func (a \*[A-Za-z0-9]*ApiService) ${op}(" "$f" | head -1 | cut -d: -f1)
  [ -z "$ln" ] && return
  sed -n "1,$((ln - 1))p" "$f" | tac | sed -E '/^[[:space:]]*(\/\/|\*|\/\*)/!Q' | tac
}

# limit_probe OP - an emitted default beats the doc comment: DatacentersGet's says "the first
# 100 items", its client sends 1000.
limit_probe() {
  local op=$1 f body n doc
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
