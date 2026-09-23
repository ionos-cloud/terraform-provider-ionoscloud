#!/usr/bin/env bash
# add-list-resource: step 0 and the verification ladder. Run it, do not read it.
# Read-only except rung*, which compile/test/lint: one at a time, never in parallel.
. "$(dirname "$0")/_lib.sh"
root

usage() { cat >&2 <<'U'
usage: scripts/probe.sh <subcommand> [args]
  discover <ionoscloud_type> [Const] [suffix] [resource]  step 0, all seven questions
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
  # 1. ResourcesMap is keyed by the CONSTANT, never the literal: an unresolvable type
  #    name is itself the answer - no such resource.
  [ -n "$c" ] || c=$(grep -E "^[[:space:]]*[A-Za-z0-9]+[[:space:]]+=[[:space:]]+\"$t\"" \
      utils/constant/constants.go | awk '{print $1}' | head -1)
  [ -n "$c" ] || die "no constant holds \"$t\" - the type does not exist"
  printf '== 1. type constant ==\nconstant.%s = "%s"\n\n' "$c" "$t"
  # 2/3. Which branch. ResourcesMap and DataSourcesMap key off the SAME constant, so
  #      match the FACTORY. Only a DataSource hit, or only a data_source_*.go in (3):
  #      no managed resource, nothing to list.
  grep -nE "constant\.$c: *[Rr]esource" ionoscloud/provider.go |
    q "2. SDKv2 ResourcesMap (listable, the harder branch)" "not an SDKv2 managed resource"
  grep -nE "constant\.$c: *[Dd]ataSource" ionoscloud/provider.go |
    q "2b. SDKv2 DataSourcesMap (this alone: nothing to list)" "no SDKv2 data source"
  grep -rn "ProviderTypeName + \"_$sfx\"" internal/framework/services/ --include=*.go |
    grep -v _test.go |
    q "3. framework-native (only a resource_*.go hit counts)" "not framework-native"
  # 4. The list resource requires an identity; a second one is a bug.
  { grep -n 'ResourceIdentity' "ionoscloud/resource_$r.go" 2>/dev/null
    grep -rn 'IdentitySchema' internal/framework/services/ --include=*.go | grep -i "$sfx"
  } | q "4. identity already declared" "none yet - you add it"
  # 5. Registration is the authority, not the docs page.
  { grep -oE 'New[A-Za-z0-9]+ListResource' ionoscloud/list_resources.go
    ls docs/list-resources/ 2>/dev/null
  } | grep -Ei "$(printf '%s' "$sfx" | sed 's/_/_?/g')" |
    q "5. list resource already registered / documented" "none yet - you add it"
  # 6/7. SDKv2 gates only: a framework-native resource has no ionoscloud/resource_*.go.
  if [ -f "ionoscloud/resource_$r.go" ]; then
    sdk_client "ionoscloud/resource_$r.go"
    "$SCRIPTS/sdkv2.sh" writer "$r"
  else
    printf '== 6/7. SDK client and state writer ==\nVERDICT: NO MATCH - no ionoscloud/resource_%s.go; gates 3 and 4 do not apply. If (3)\n         matched you are on the framework-native branch.\n\n' "$r"
  fi
  ;;

# listcalls [filter] - is there a PARENTLESS collection GET at all? None for yours =>
# it cannot be listed; stop. [A-Za-z0-9]+ is load-bearing: without the digits this
# silently misses K8sGet.
# listcalls [resource_or_product] - every parentless collection GET in the vendored SDKs.
# The filter matches the SDK FILE PATH as well as the signature, because the product name
# lives only in the path: mariadb's collection GET is ClustersGet in
# products/dbaas/mariadb/v2/api_clusters.go, and a signature-only filter says "not
# vendored" and wrongly stops the run at gate 2. Underscore-separated parts must ALL
# appear somewhere in path+signature, so `listcalls mariadb_cluster` narrows to mariadb
# while a bare `cluster` does not drown you. An empty result here is a REAL stop.
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
        dropped="$dropped $p"          # narrowing on this part would empty the set
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

# datasource-limit - a hit means a bare fetch would cap the list resource BELOW its own
# data source. No match = the data source sends none either; claim no number in docs.
datasource-limit)
  need "ionoscloud/data_source_${2:?resource}.go"
  grep -n 'Limit(' "ionoscloud/data_source_$2.go" |
    q "explicit .Limit(n) in the sibling data source" "it sends none either"
  ;;

# locations <product> - only when the resource is partitioned by location. No match =>
# partitioned with nothing to enumerate: stop and ask, do not ship a partial listing.
locations)
  need "services/${2:?product}/"
  grep -rn 'AvailableLocations\|Valid[A-Za-z]*Locations *=' --include=*.go "services/$2/" |
    q "enumerable locations for $2" "not enumerable in-tree"
  ;;

# ---- verification ladder: cheapest first, strictly sequential, one at a time. ----
# rung0 - never with an empty argument list: gofmt then reads stdin and reports green
# having inspected nothing.
rung0) gofmt -l -d ./ionoscloud ./internal ./utils ;;
# rung1 - only the packages you touched. framework-native: the service package, plus
# ./internal/framework/provider/ ONLY if the list resource went in a new one.
rung1) shift; go build "${@:-./ionoscloud/}" ;;
rung2) go test "${2:?pkg}" -run "${3:?TestName}" -count=1 ;;
# rung2b - the only thing that runs InternalIdentityValidate() on your Identity. Rung 2
# passes with a malformed one and nothing in CI checks.
rung2b) go test ./ionoscloud/ -run 'TestProvider$' -count=1 ;;
# rung3 - your own lines only. Never drop --new-from-rev (hundreds of pre-existing
# issues) nor the path argument (`make lint` has none: all 60 packages, ~45 linters).
rung3) shift; golangci-lint run --new-from-rev "$(git merge-base origin/master HEAD)" "$@" ;;
# rung4 - only if imports changed. Run the vendor half EVEN WHEN tidy is clean: a new
# subpackage of a module already in go.mod moves vendor/modules.txt alone (#1034).
rung4) go mod tidy && git diff --exit-code -- go.mod go.sum
       go mod vendor && git status --porcelain vendor/ ;;
rung5-scoped) shift; go vet -tags=all ./ionoscloud/ ./internal/framework/provider/ ./internal/acctest/ "$@" ;;
# rung5-full - EXPENSIVE, heavier than `go build ./...`: ASK THE USER FIRST. The only
# thing that compiles the tagged acceptance file holding your Query: true step.
rung5-full) go vet -tags=all ./... ;;

*) usage ;;
esac
