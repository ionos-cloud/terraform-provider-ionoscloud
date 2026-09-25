#!/usr/bin/env bash
. "$(dirname "$0")/_lib.sh"
root

case "${1:-}" in

identities)
  grep -ln 'ResourceIdentity{' ionoscloud/*.go |
    q "SDKv2 resources declaring an Identity" "none yet - you are the first"
  grep -oE 'New[A-Za-z0-9]+ListResource' ionoscloud/list_resources.go | sort -u |
    q "registered list resources" "none registered yet"
  ;;

# A match => pass .Depth(1); none => it returns full properties already.
depth)
  s=${3:?api stem, e.g. zones}; s=${s#api_}; s=${s%.go}
  f="vendor/github.com/ionos-cloud/sdk-go-bundle/products/${2:?product}/v2/api_$s.go"
  need "$f"
  grep -nE "func \(r ${4:?request type}(Request)?\) Depth" "$f" |
    q "Depth on $4" "no Depth on this collection - drop .Depth(1)"
  ;;

writer)
  f="ionoscloud/resource_${2:?resource}.go"
  need "$f"
  fn=$(grep -o 'ReadContext: *[A-Za-z0-9_]*' "$f" | head -1 | awk '{print $2}')
  if [ -z "$fn" ]; then
    printf '== 7. state writer ==\nVERDICT: NO MATCH - no ReadContext in %s; follow ResourcesMap in ionoscloud/provider.go.\n\n' "$f"
    exit 0
  fi
  # Candidates: everything Read calls with d as its first argument.
  cands=$(sed -n "/^func ${fn}(/,/^}/p" "$f" | grep -oE '[A-Za-z0-9_.]+\(d, ' |
          sed 's/(d, //; s/.*\.//' | sort -u)
  if [ -z "$cands" ]; then
    printf '== 7. state writer (%s calls) ==\nVERDICT: NO MATCH - %s calls nothing with (d, ...) first.\n         It writes state inline; a mapper cannot reuse it. Read it before designing one.\n\n' "$fn" "$fn"
    exit 0
  fi
  grep -rnE "^func (\([^)]*\) )?($(printf '%s' "$cands" | paste -sd'|'))\(" ionoscloud/ services/ internal/ |
    q "7. state writer signatures ($fn calls: $(printf '%s' "$cands" | paste -sd' '))" \
      "the candidates are declared outside ionoscloud/, services/ and internal/"
  printf 'Contract: func(d *schema.ResourceData, obj <sdk>.<X>) error - a receiver is fine,\nthe signature is what matters. Then gate 4: does it d.Set every Required attribute?\n\n'
  ;;

*) echo "usage: $0 {identities | depth <product> <api-stem> <RequestType> | writer <resource>}
  e.g. depth dns zones ApiZonesGet (the stem of api_zones.go)" >&2; exit 2 ;;
esac
