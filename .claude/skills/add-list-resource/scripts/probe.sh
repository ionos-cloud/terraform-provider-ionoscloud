#!/usr/bin/env bash
# Read-only except rung*, which build/test/lint: never in parallel.
. "$(dirname "$0")/_lib.sh"
root

usage() { cat >&2 <<'U'
usage: scripts/probe.sh <subcommand> [args]
  discover <ionoscloud_type> [Const] [suffix] [resource]  step 0, all seven questions
                            (Const is optional: the full identifier, e.g. IpBlockResource)
  listcalls [filter]        every parentless collection GET in the vendored SDKs
  limit-cloud <Op> | limit-bundle <Op>   what limit the client sends (one probe, both SDKs)
  datasource-limit <resource> | locations <product>
  rung0 | rung1 [pkg...] | rung2 <pkg> <Test> | rung2b | rung3 <pkg...> | rung4 |
  rung5-scoped | rung5-full
U
exit 2; }

case "${1:-}" in

discover)
  t=${2:-}; [ -n "$t" ] || usage
  c=${3:-}; sfx=${4:-${t#ionoscloud_}}; r=${5:-${t#ionoscloud_}}
  need utils/constant/constants.go ionoscloud/provider.go \
       ionoscloud/list_resources.go internal/framework/services
  # Every constant holding the type: a data source and a resource often share one string.
  all=$(grep -E "^[[:space:]]*[A-Za-z0-9]+[[:space:]]+=[[:space:]]+\"$t\"" \
      utils/constant/constants.go | awk '{print $1}')
  [ -n "$all" ] || die "no constant holds \"$t\" - the type does not exist"
  alt=$(printf '%s\n' "$all" | paste -sd'|')
  if [ -n "$c" ]; then
    printf '%s\n' "$all" | grep -qx "$c" ||
      die "constant.$c does not hold \"$t\" - pass the full identifier ($(printf '%s\n' "$all" | paste -sd' ')) or omit it"
  else
    c=$(grep -oE "constant\.($alt): *[Rr]esource" ionoscloud/provider.go | head -1 | sed 's/constant\.//; s/:.*//')
    [ -n "$c" ] || c=$(printf '%s\n' "$all" | head -1)
  fi
  printf '== 1. type constant ==\nconstant.%s = "%s" (every constant holding it: %s)\n\n' \
    "$c" "$t" "$(printf '%s\n' "$all" | paste -sd' ')"
  # 2/2b. Both maps key off those constants: match the FACTORY.
  grep -nE "constant\.($alt): *[Rr]esource" ionoscloud/provider.go |
    q "2. SDKv2 ResourcesMap (listable, the harder branch)" "not an SDKv2 managed resource"
  grep -nE "constant\.($alt): *[Dd]ataSource" ionoscloud/provider.go |
    q "2b. SDKv2 DataSourcesMap (this alone: nothing to list)" "no SDKv2 data source"
  grep -rn "ProviderTypeName + \"_$sfx\"" internal/framework/services/ --include=*.go |
    grep -v _test.go |
    q "3. framework-native (only a resource_*.go hit counts)" "not framework-native"
  # 4. A second identity is a bug.
  { grep -n 'ResourceIdentity' "ionoscloud/resource_$r.go" 2>/dev/null
    grep -rn 'IdentitySchema' internal/framework/services/ --include=*.go | grep -i "$sfx"
  } | q "4. identity already declared" "none yet - you add it"
  # 5. Registration is the authority, not the docs page.
  { grep -oE 'New[A-Za-z0-9]+ListResource' ionoscloud/list_resources.go
    ls docs/list-resources/ 2>/dev/null
  } | grep -Ei "$(printf '%s' "$sfx" | sed 's/_/_?/g')" |
    q "5. list resource already registered / documented" "none yet - you add it"
  if [ -f "ionoscloud/resource_$r.go" ]; then
    sdk_client "ionoscloud/resource_$r.go"
    "$SCRIPTS/sdkv2.sh" writer "$r"
  else
    printf '== 6/7. SDK client and state writer ==\nVERDICT: NO MATCH - no ionoscloud/resource_%s.go; gates 3 and 4 do not apply. If (3)\n         matched you are on the framework-native branch.\n\n' "$r"
  fi
  ;;

# The filter matches the SDK file path too, where the product name lives (mariadb's ClustersGet:
# products/dbaas/mariadb/v2/api_clusters.go); [A-Za-z0-9]+ keeps K8sGet. Every _-part must match.
listcalls)
  need vendor/github.com/ionos-cloud
  out=$(grep -rnoE "func \(a \*[A-Za-z0-9]+\) [A-Za-z0-9]+(Get|List)\(ctx _?context\.Context\) Api[A-Za-z0-9]+Request" \
        vendor/github.com/ionos-cloud/ | sed 's/func (a \*//' | sort -u)
  [ -n "$out" ] || die "the collection-GET sweep matched nothing at all - wrong tree, or vendor/ stripped"
  f=${2:-}; f=${f#ionoscloud_}
  kept=; dropped=
  if [ -n "$f" ]; then
    IFS='_' read -r -a parts <<<"$f"
    for p in "${parts[@]}"; do
      [ -n "$p" ] || continue
      if narrowed=$(printf '%s\n' "$out" | grep -i -- "$p") && [ -n "$narrowed" ]; then
        out=$narrowed; kept="$kept $p"
      else
        dropped="$dropped $p"
      fi
    done
  fi
  { [ -n "$out" ] && printf '%s\n' "$out"; } |
    q "parentless collection GET candidates${kept:+ matching$kept}" \
      "no such collection is vendored: this resource cannot be listed"
  [ -n "$dropped" ] && printf 'NOTE: no candidate mentions%s. The lines above are this product'"'"'s OTHER\n      collections - confirm the operation returns %s objects before treating this\n      as listable. A sibling collection is not your resource.\n\n' "$dropped" "${2:-the resource}"
  [ -n "$f" ] && printf 'CONFIRM: a candidate counts only if its response items are the same type the state\n         writer takes. Check the Api*Request'"'"'s Execute() return type.\n\n'
  ;;

limit-cloud|limit-bundle) limit_probe "${2:?operation, e.g. DatacentersGet}" ;;

datasource-limit)
  need "ionoscloud/data_source_${2:?resource}.go"
  grep -n 'Limit(' "ionoscloud/data_source_$2.go" |
    q "explicit .Limit(n) in the sibling data source" "it sends none either"
  ;;

locations)
  need "services/${2:?product}/"
  grep -rn 'AvailableLocations\|Valid[A-Za-z]*Locations *=' --include=*.go "services/$2/" |
    q "enumerable locations for $2" "not enumerable in-tree"
  ;;

rung0)
  out=$(gofmt -l -d ./ionoscloud ./internal ./utils) || die "gofmt failed - see the error above"
  if [ -n "$out" ]; then printf '%s\nVERDICT: NOT CLEAN - gofmt the files above\n' "$out"; exit 1; fi
  echo 'VERDICT: gofmt clean' ;;
rung1) shift; go build "${@:-./ionoscloud/}" && echo 'VERDICT: builds' ;;
rung2) go test "${2:?pkg}" -run "${3:?TestName}" -count=1 ;;
# rung2b - the only thing that runs InternalIdentityValidate() on your Identity.
rung2b) go test ./ionoscloud/ -run 'TestProvider$' -count=1 ;;
# rung3 - never drop --new-from-rev nor the path argument (`make lint` has none: all packages).
rung3) shift; havegit "rung 3 needs git merge-base: lint where git is allowed"
       golangci-lint run --new-from-rev "$(git merge-base origin/master HEAD)" "$@" ;;
# rung4 - the vendor half runs EVEN WHEN tidy is clean: a new subpackage of a module already
# in go.mod moves vendor/modules.txt alone.
rung4) havegit "rung 4 diffs go.mod and vendor/ with git"
       go mod tidy && git diff --exit-code -- go.mod go.sum; mod=$?
       go mod vendor && git status --porcelain vendor/ && exit "$mod" ;;
rung5-scoped) shift; go vet -tags=all ./ionoscloud/ ./internal/framework/provider/ ./internal/acctest/ "$@" &&
  echo 'VERDICT: vet clean - the tagged tests compile' ;;
rung5-full) go vet -tags=all ./... && echo 'VERDICT: vet clean' ;;

*) usage ;;
esac
