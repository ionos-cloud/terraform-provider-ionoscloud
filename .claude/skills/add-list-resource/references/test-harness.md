# Test harness

Model: `ionoscloud/resource_datacenter_list_test.go`. Yours:
`ionoscloud/resource_<resource>_list_test.go`, package `ionoscloud_test`, **no `//go:build` tag**;
only its `ListResourceSchemas`/`GetResourceIdentitySchemas` assertions are SDKv2-specific.

## Fixture and helpers

`TestDatacenterListResource:36-73`, up to `resourceType`, is the fixture: copy it, renaming only,
its five `t.Fatalf` guards included; nothing else from that test. The `t.Setenv`s precede the
server, or `~/.ionos/config` wins. `const <resource>ListType` is package-level.

**Shared helpers exist — write none**: `scripts/test.sh` prints them, their file and the driver's
signature. `decode`/`goValue` give numbers as `int64`, never `float64`; `hasErrorDiagnostic` is the
bool, `failOnErrorDiagnostics` the fatal.

- **Your own stub**: `stubCloudAPI` is pinned to `/datacenters`; reused, it 404s every request.
- **The driver: call it, never fork it.** Already `listResults(ctx, t, server, typeName, schema,
  filters, includeResource)`? Call it. Else parameterise master's `listDatacenters` in place:
  rename it `listResults`, add `typeName` and `includeResource bool` (the nil-Resource `t.Fatalf`
  only when true), fix the datacenter call sites, keep its explained `//nolint`s.

## The stub

`stubCloudAPI` (`:176-226`) is the skeleton: a `strings.HasSuffix` match, all else 404 on purpose.

- **Items.** Result 1: every field the writer reads, lists and blocks included. Result 2:
  **minimal** — the id, `name`, required attributes and SDK-validated enums, nothing else. Result
  3 **only for an Optional `name`**: `Name` unset.
- **Bundle stubs fill every non-omitempty enum** on every item, commented: the SDK validates enums
  while unmarshalling, and a missing one fails the *fetch*, naming neither field nor item. Check
  that `IONOS_API_URL` reaches the client (`utils/loadedconfig` defers `ChangeConfigURL` over it).
- **Assert the query the fetch sends**, each value a literal (`"1000"`, not `constant.<X>Limit`):
  the Cloud API client always sends `depth`, `offset` and `limit` (`0`, `0`, `100` unless set), so
  assert all three; a bundle client sends only what is set: assert every other parameter it knows
  is absent (`r.URL.Query().Has`). The stub's comment names exactly what it checks.

## The subtests

```go
	t.Run("streams every <resource>", func(t *testing.T) {
		results := listResults(ctx, t, server, <resource>ListType, listSchema, nil, true)
		if len(results) != N { // N = stub items
			t.Fatalf("expected N results, got %d", len(results))
		}
		assert.Equal(t, "<display-name-1>", results[0].DisplayName)
		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "<id-1>", identity["id"]) // every identity attribute, on every result

		resource := decode(t, results[0].Resource, resourceType)    // EVERY attribute the writer
		assert.Equal(t, "<id-1>", resource["id"])                   // fills: a wrong key must fail
		assert.ElementsMatch(t, []any{...}, resource["<set-attr>"]) // a TypeSet is unordered
		assert.Equal(t, []any{map[string]any{ /* … */ }}, resource["<block>"]) // form: :98-104
		assert.Nil(t, resource["timeouts"], "a listed <resource> has no timeouts")

		second := decode(t, results[1].Resource, resourceType) // result 2: DisplayName, identity,
		assert.Nil(t, second["<optional-attribute>"])           // required values, and each
		assert.Equal(t, []any{}, second["<block-or-list>"])     // omitted one as sdkv2 §2d says

		assert.Equal(t, "<id-3>", results[2].DisplayName) // Optional name: the fallback
		assert.Nil(t, decode(t, results[2].Resource, resourceType)["name"])
	})
```

Then, as `:128-148`: a subtest per `FilterAttribute(...)` field (Optional `name`: one matching the
unnamed item and another, in order); `applies every filter`, two fields matching *different*
items, expecting 0; `identity only` (`includeResource` false: identity and DisplayName set,
`Resource` nil); `rejects unknown filter fields`, over `"nope"` and each attribute kept out of the
allow-list on purpose (its reason a comment on `FilterAttribute`); `accepts every allowed filter
field`, the same validation, error-free for each allowed one.

## The mutation check — mandatory

Mutate the state writer, **not your file** (`ionoscloud/resource_<x>.go`, or
`services/<product>/<x>.go`); **revert by hand**, never `git checkout --` (it may hold your
identity).

- Swap the **key strings** of two adjacent `d.Set` calls or map entries of one Go type, values
  and nil guards left in place (datacenter: `max_cores` ↔ `max_ram`), one pair per type, one in a
  nested block if any; no pair: change that field's value in the stub.
- The identity setter: swap its two values, or write `d.Id()+"x"` for a lone `id`.
- The list file: drop `.Depth(1)` and any `.Limit(...)` (none sent? add one, or a server-side
  filter to a bundle helper); swap the two `MatchesFilters` values; swap the writer and setter
  calls; delete the `displayName` fallback; drop a field from `FilterAttribute(...)`.

Each: mutate → run → it **must fail** → revert → green. A pass means an assertion is missing: add
it. Claim the check done only if you saw each one fail.

## The acceptance query step

`TestAcc<Resource>Query`, in the resource's existing acceptance-test file, keeping its `//go:build`
line. **Write it, never run it**; rung 5 compiles it. Model: `TestAccDataCenterQuery`
(`ionoscloud/resource_datacenter_test.go:133-212`): its `TestCase` fields, `ExpectIdentity` check,
`filters` blocks and identity-import step (and why no `ImportStateVerify`), but your own fixture:

```go
	const (
		<resource>Name = "tf-test-<resource>-query"
		<resource>Addr = constant.<Const> + ".test_<resource>_query"
		otherLocation  = "de/txl" // the zero-result filter value: set by NO step
	)

	// Own name, own label, own create step — NOT testAccCheck<Resource>ConfigBasic:
	// ExpectLength asserts a CONTRACT-WIDE total, so a same-named resource made by
	// another test in the package makes ExpectLength(1) flap.  (copy this comment in)

	// Steps: create; Query: true + no filters → ExpectIdentity; Query: true + filters matching
	// exactly one → ExpectLength 1; Query: true + same name, otherLocation → ExpectLength 0;
	// only if Update does not end in Read (sdkv2-branch.md 1c), a non-ForceNew change with
	//   ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
	//     plancheck.ExpectResourceAction(<resource>Addr, plancheck.ResourceActionUpdate)}},
	//   Check: resource.TestCheckResourceAttr(<resource>Addr, "<attribute>", <updated value>),
	// then the ImportBlockWithResourceIdentity import.
```

- **`ProtoV6ProviderFactories`, not `ProviderFactories`**, which cannot see the list resource.
- **`Query: true` on every list step**, or its `QueryResultChecks` silently never run.
- **Keep the zero-result step**: it alone proves the filter is not a no-op. Locationless: another
  filtered field, its value a const of its own.
