# Test harness

`ionoscloud/resource_datacenter_list_test.go` (`ionoscloud_test`) drives the real `ListResource`
RPC against a stubbed API, through the mux `main.go:30-47` serves. Why, and what the assertions
catch: `references/why.md`. Yours lives in `ionoscloud/resource_<resource>_list_test.go`, same
external package, **no `//go:build` tag** (no `TF_ACC` or credentials; stays in lint scope).
`set<X>Data` = the resource's own state writer (`sdkv2-branch.md` §2d). Everything here is
branch-agnostic except the `ListResourceSchemas`/`GetResourceIdentitySchemas` assertions, which
prove `RawV6Schemas` and so are SDKv2-only.

## Fixture

`TestDatacenterListResource:36-73`, up to `resourceType`, is identical in every list test. Take
that block, change only the names, and take nothing else from that file: its assertions are
datacenter-specific and its stub is the wrong model. Keep all five of its `t.Fatalf` guards — the
two RPC-error ones (`:46`, `:61`) and the three that name the failure: list resource not registered
(`:53`), managed resource missing from the merged schema (`:56`), no identity declared (`:67`).
The `t.Setenv`s must precede the server, or `~/.ionos/config` overrides the stub URL;
`const <resource>ListType` is package-level.

## The shared helpers exist — write none of them

`scripts/test.sh` says which and where — on the SDKv2 branch, all of them in
`resource_datacenter_list_test.go`; re-declaring one is a compile error. You add the const,
the test, `stub<Resource>API` and its subtests.

| Helpers | For |
|---|---|
| `muxedProviderServer`, `configureProvider` | the server, configured from env |
| `decode`, `goValue`, `readValue` | `DynamicValue`/`tftypes.Value` → Go |
| `identityType`, `nullObject`, `dynamicValue` | building the RPCs' input values |
| `failOnErrorDiagnostics`, `hasErrorDiagnostic` | fatal on a diagnostic; bool, last subtest |
| `listServerAndConfig`, `listConfig` | the list server and its `filters` config |
| the list driver | `listDatacenters` on master, generic once parameterised — see below |

- **Rename the stub.** `stubCloudAPI` is pinned to `/datacenters` (`:212`); reusing it is the
  collision that does not announce itself: it compiles, then 404s every request.
- **Parameterise the driver, never fork it.** On master it is
  `listDatacenters(ctx, t, server, schema, filters)` (`:245`) and it hardcodes
  `datacenterListType`. Add a `typeName` parameter in place, give it a generic name, and fix
  its call sites — the next list test reuses it. If a generic driver already exists (someone
  got there first), just call it. **Run `scripts/test.sh` and trust what it prints over any
  name quoted here**: this is the line most likely to be stale.
- Keep the two explained `//nolint`s (`:259`, `:277`) if you touch them.
- `goValue` decodes numbers as `int64`; it has no `Map`/`Tuple` case (`why.md`).

## The stub

`stubCloudAPI:176-226` is the skeleton: a `strings.HasSuffix` match (the SDK prefixes
`/cloudapi/v6`), all else 404 **on purpose**. Three differences.

- **Items.** Result 1: everything the writer reads. Result 2: only what the API always returns,
  `name` still set even when Optional (the display-name and filter assertions need it) and a
  *different* optional one omitted. Result 3 **only when `name` is Optional**: `Name` unset, for
  the `displayName` fallback. Match the count assertion.
- **A bundle stub must fill every non-omitempty enum field on every item, minimal one included** —
  enums validate while unmarshalling, so a zero value fails the *fetch*, naming neither field nor
  item. Check too that `IONOS_API_URL` reaches that client (`why.md`).
- **Assert the query the fetch closure builds**; the reference asserts none, the one thing not to
  copy. One per explicit option —
  `if got := r.URL.Query().Get("depth"); got != "1" { t.Errorf(...) }` — a limit as a literal, not
  `constant.<X>Limit`. If it sets **nothing** deliberately, assert the absence: loop `depth`,
  `limit`, `offset`, `<filter-param>`, `t.Errorf`ing on any set.

## The subtests

```go
	t.Run("streams every <resource>", func(t *testing.T) {
		results := <driver>(ctx, t, server, <resource>ListType, listSchema, nil)
		// N = the number of stub items: 2, or 3 when the optional-name result 3 exists
		if len(results) != N {
			t.Fatalf("expected N results, got %d", len(results))
		}
		assert.Equal(t, "<display-name-1>", results[0].DisplayName)
		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "<id-1>", identity["id"])   // one per identity attribute

		resource := decode(t, results[0].Resource, resourceType)    // EVERY attribute the writer
		assert.Equal(t, "<id-1>", resource["id"])                   // fills: a wrong key must fail
		assert.ElementsMatch(t, []any{...}, resource["<set-attr>"]) // a TypeSet is unordered
		assert.Equal(t, []any{map[string]any{ /* … */ }}, resource["<block>"]) // form: :98-104
		assert.Nil(t, resource["timeouts"], "a listed <resource> has no timeouts")

		assert.Equal(t, "<display-name-2>", results[1].DisplayName) // result 2: pairing holds
		second := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "<id-2>", second["id"])
		assert.Nil(t, second["<optional-attribute>"])  // per omitted optional ATTRIBUTE — check
		// the WRITER: an omitted nested BLOCK reifies null→empty, and a list attribute d.Set() with
		// a nil slice materialises empty (SetZoneData does that with `nameservers`).

		// Only when `name` is Optional, and why N is 3 then — the displayName fallback:
		assert.Equal(t, "<id-3>", results[2].DisplayName)
		assert.Nil(t, decode(t, results[2].Resource, resourceType)["name"])
	})
```

Per `:128-148`: **one subtest per field in the `FilterAttribute(...)` allow-list** (an untested
field silently matches nothing, forever); **`applies every filter`**, two fields matching
*different* items, expecting 0; and **`rejects unknown filter fields`**, which also fails if
`FilterAttribute()` has no allow-list.

`go test ./ionoscloud/ -run Test<Resource>ListResource -count=1`; CI will not
(`verify-and-pr.md` §1).

## The mutation check — mandatory, not advice

**Prove the assertions are load-bearing before claiming coverage.** Mutate the state writer — **a
different file from the one you just wrote** (step 0's grep (7)): `ionoscloud/resource_<x>.go`, or
`services/<product>/<x>.go` on the bundle branch (`SetZoneData`: `services/dns/zone.go:60`).

> **Revert every mutation by hand.** `git checkout --` discards the whole identity implementation
> on the Cloud API branch; on the bundle branch that file is shared with the managed resource and
> its data source.

Swap **two adjacent fields or map keys of the same Go type** (so the compiler cannot stand in for
the assertion), one pair **per type the writer sets**, preferring one in a nested block:
datacenter's `max_cores` ↔ `max_ram` (`ionoscloud/resource_datacenter.go:416-420`, both `*int32`).
No adjacent pair? Change that field's value in the stub. Mutate `set<X>Identity` too: swap its
two values if it sets two (`setDatacenterIdentity` sets `id` and `location`,
`ionoscloud/resource_datacenter.go:340-353`). **A single-attribute identity has no pair — do not
skip it:** make the setter write something other than `d.Id()`, e.g. `identity.Set("id", d.Id()+"x")`.
If the test still passes, nothing is asserting the identity and the coverage claim is false.

Each: mutate → run → it **MUST** fail, naming both attributes → revert by hand → confirm green. A
**pass** means the assertion set is incomplete; add the missing assertion first, and the same to a
field only result 2 leaves unset.

**Four more, in the list resource file:** delete `.Depth(1)` and any `.Limit(...)` from the fetch
closure (sets none deliberately? invert it: *add* one); swap the `MatchesFilters` map's two values,
which only the per-field subtests catch; swap the `set<X>Data` and `set<X>Identity` calls; and, for
optional-`name`, delete the `if displayName == ""` fallback.

Report the mutation check as done only if you ran it and saw it fail.
