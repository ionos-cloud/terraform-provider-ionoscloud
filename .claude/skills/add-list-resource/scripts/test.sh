#!/usr/bin/env bash
. "$(dirname "$0")/_lib.sh"
root
need ionoscloud

ls ionoscloud/*_list_test.go 2>/dev/null |
  q "files holding list-resource tests" \
    "you are the FIRST list test in this package - only now do you write the helpers yourself, and put them in a neutrally named file (list_resource_test_helpers_test.go), never in one named after your resource"

# A helper absent here is missing: add it to the file holding the others, never a second one.
grep -nE '^func [a-z][A-Za-z0-9_]*\(' ionoscloud/*_list_test.go 2>/dev/null |
  grep -v ':func stub' |
  q "shared helpers - already declared, do not write them again" \
    "no helpers exist yet (consistent with there being no list test)"

grep -nE '^func (list[A-Za-z0-9]*|stub[A-Za-z0-9]*API)\(' ionoscloud/*_list_test.go 2>/dev/null |
  q "driver (its signature: typeName, includeResource yet?) and stub names already taken" "none taken yet"

echo "Then, one at a time (probe.sh rung2 ./ionoscloud/ Test<Resource>ListResource, rung2b)."
