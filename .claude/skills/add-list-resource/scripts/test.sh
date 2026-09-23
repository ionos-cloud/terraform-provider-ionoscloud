#!/usr/bin/env bash
# "Which shared list-test helpers already exist, and where?" Derived from the tree, not
# from a list written into the skill - the set grows every time this skill is used.
# Answer it before writing a helper: re-declaring one in the same package is a
# "redeclared in this block" compile error.
. "$(dirname "$0")/_lib.sh"
root
need ionoscloud   # else "no list test yet" would be a conclusion drawn from a typo

ls ionoscloud/*_list_test.go 2>/dev/null |
  q "files holding list-resource tests" \
    "you are the FIRST list test in this package - only now do you write the helpers yourself, and put them in a neutrally named file (list_resource_test_helpers_test.go), never in one named after your resource"

# Every lower-case top-level func in those files bar stub* is a shared helper: reuse it,
# do not re-declare it. One absent from the listing really is missing - add it to the
# file that already holds the others, never start a second helpers file.
grep -nE '^func [a-z][A-Za-z0-9_]*\(' ionoscloud/*_list_test.go 2>/dev/null |
  grep -v ':func stub' |
  q "shared helpers - already declared, do not write them again" \
    "no helpers exist yet (consistent with there being no list test)"

# The driver's name is not fixed: listDatacenters(ctx,t,server,schema,filters) on master,
# something generic once it has been parameterised. This prints whatever is ACTUALLY
# declared - trust it over any name quoted in the reference files, which is what rots.
grep -nE '^func (list[A-Za-z0-9]*|stub[A-Za-z0-9]*API)\(' ionoscloud/*_list_test.go 2>/dev/null |
  q "driver and stub names already taken (yours must not collide)" "none taken yet"

echo "Then, one at a time (probe.sh rung2 ./ionoscloud/ Test<Resource>ListResource, rung2b)."
