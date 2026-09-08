# Test harness

Reference: `internal/framework/services/compute/resource_datacenter_list_test.go`
(package `compute_test`). It drives the real `ListResource` RPC end to end against a stubbed
Cloud API, through the same mux `main.go` serves — no acceptance test, no credentials, and it
runs in well under a second. It is the oldest of the three list-resource tests in `compute`
(datacenter, ipblock, target_group); the five framework-native list resources have none. Copy it,
and read `resource_target_group_list_test.go` too for the lone-`id` identity and `MaxItems: 1` cases.

New file: `internal/framework/services/<service>/resource_<resource>_list_test.go`,
package `<service>_test`. The external test package is load-bearing: it is what lets the test
import `.../v6/ionoscloud` even though the non-test package must not.

No `//go:build` tag. That is deliberate — it keeps the test cheap to run and puts it in
golangci-lint's scope (which only analyses the default build).

---

## Shape

```go
func Test<Resource>ListResource(t *testing.T) {
	ctx := context.Background()

	t.Setenv("IONOS_API_URL", stub<Resource>API(t))
	t.Setenv("IONOS_TOKEN", "token-for-the-stub")

	server := muxedProviderServer(ctx, t)

	providerSchema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema: %v", err)
	}
	failOnErrorDiagnostics(t, "GetProviderSchema", providerSchema.Diagnostics)

	// The framework only registers a list resource that has no framework resource
	// behind it once RawV6Schemas has supplied both protocol schemas.
	if _, ok := providerSchema.ListResourceSchemas[<resource>ListType]; !ok {
		t.Fatalf("the <resource> list resource was not registered; check RawV6Schemas and the SDKv2 resource identity")
	}
	if _, ok := providerSchema.ResourceSchemas[<resource>ListType]; !ok {
		t.Fatalf("the <resource> managed resource is missing from the merged schema")
	}

	identitySchemas, err := server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		t.Fatalf("GetResourceIdentitySchemas: %v", err)
	}
	failOnErrorDiagnostics(t, "GetResourceIdentitySchemas", identitySchemas.Diagnostics)

	identitySchema := identitySchemas.IdentitySchemas[<resource>ListType]
	if identitySchema == nil {
		t.Fatalf("the SDKv2 <resource> resource does not declare a resource identity")
	}

	configureProvider(ctx, t, server, providerSchema.Provider)

	listSchema := providerSchema.ListResourceSchemas[<resource>ListType]
	resourceType := providerSchema.ResourceSchemas[<resource>ListType].ValueType()

	// ... three subtests
}
```

Set the env vars **before** building the server: the framework provider's `Configure` reads
`IONOS_API_URL` into `clientOptions.Endpoint`, and `bundleclient.shouldApplyOverridesCustomEnv`
short-circuits to `false` whenever it is set — so the developer's real `~/.ionos/config`
overrides are ignored and the test cannot accidentally hit a live endpoint.

Declare `const <resource>ListType = "<ionoscloud_type>"` at package level (`goconst` flags a
literal repeated 5+ times).

---

## Copy these helpers character-exact — but only into a package that does not have them

**Check first: does the target package already contain a `*_list_test.go`?**

```bash
ls internal/framework/services/<service>/*_list_test.go
```

If it does, **copy none of the helpers below.** They are package-level declarations, shared by
every file in the test package, and re-declaring them is 14 `redeclared in this block` errors.
This is not a corner case — `sdkv2-branch.md` §2 routes the entire Cloud API family into
`compute`, and `compute` already has all of them
(everything below `TestDatacenterListResource` in `resource_datacenter_list_test.go`). When they
are already there, the only things you write per resource are:

- `const <resource>ListType = "<ionoscloud_type>"`
- `Test<Resource>ListResource`
- `stub<Resource>API` — **rename it.** The existing `stubCloudAPI` is pinned to `/datacenters`
  (`resource_datacenter_list_test.go:209`); reusing the name is the one collision that does not
  announce itself, because a same-named stub compiles and then 404s every request your resource
  makes.
- the subtests

**Do not copy the driver.** It differs between resources only by type name, so call the shared
`listResults(ctx, t, server, typeName, schema, filters)`. If you find it still hard-coded to one
resource, parameterise it in place — adding a `typeName string` parameter and updating the
existing call sites is a rename, not a refactor, and it is what stops copy N+1. Copying it was
the original instruction here and it is what produced two near-identical drivers before
`ionoscloud_ipblock` landed.

Do **not** change the reference test's *assertions or stub* — those are what make
`resource_datacenter_list_test.go` the reference. Adding to, or parameterising, the shared
helpers it calls is fine and expected.

**Where the shared helpers live.** If the package has no `*_list_test.go` yet, put them in
`list_resource_test_helpers_test.go` from the start — never inside a file named after one
resource. `compute` has them in `resource_datacenter_list_test.go` for historical reasons, so
every author there has to learn from a comment that the infrastructure lives in another
resource's file; do not reproduce that in a new package.

The table below is for a **new** service package with no list test yet. Take them verbatim from
`internal/framework/services/compute/resource_datacenter_list_test.go`:

| Helper | What it does |
|---|---|
| `muxedProviderServer(ctx, t)` | `ionoscloud.Provider()` → `tf5to6server.UpgradeServer` → `tf6muxserver.NewMuxServer(providerserver.NewProtocol6(fwprovider.New(sdkv2Provider)), upgraded)` — mirrors `main.go` |
| `configureProvider(ctx, t, server, schema)` | `ConfigureProvider` with an all-null config so credentials come from the env |
| `decode(t, value, valueType)` | `*tfprotov6.DynamicValue` → `map[string]any` |
| `goValue(t, value)` | recursive `tftypes.Value` → plain Go |
| `readValue(t, value, target)` | `value.As(target)` with a fatal on mismatch |
| `identityType(schema)` | identity schema → `tftypes.Object` type |
| `nullObject(t, objectType)` | all-null object value |
| `dynamicValue(t, valueType, value)` | `tfprotov6.NewDynamicValue` with a fatal |
| `failOnErrorDiagnostics(t, rpc, diags)` | fail on the first error diagnostic |
| `hasErrorDiagnostic(diags)` | bool, for the negative test |
| `listServerAndConfig(t, server, schema, filters)` | casts to `ProviderServerWithListResource`, builds the config |
| `listConfig(t, configType, filters)` | builds the `filters` list value |

Take the driver as `listResults(ctx, t, server, typeName, schema, filters)` — parameterised by
type name, one copy per package. `stubCloudAPI` → `stub<Resource>API`, one per resource.

Two lint notes on the copied code: keep the `//nolint:staticcheck` on `listServerAndConfig`
(with its explanation — `ListResource` still lives on a temporary interface in
terraform-plugin-go) and the `//nolint:prealloc` on the results slice. `nolintlint` requires
directives to be both specific and explained; a bare `//nolint` is itself an error.

### `goValue` — the part that matters

```go
	case valueType.Is(tftypes.Number):
		var number big.Float
		readValue(t, value, &number)
		// Every number in the <resource> schema is a TypeInt, so this keeps a single
		// integer representation and fails loudly rather than silently switching to a
		// float if that ever stops being true.
		if !number.IsInt() {
			t.Fatalf("expected a whole number, got %s", number.String())
		}
		n, _ := number.Int64()
		return n
```

**Decode numbers as `int64`, never `float64`.** `testifylint`'s `float-compare` rule rejects
`assert.Equal` on a float and demands `InDelta`/`InEpsilon`; it is **not** excluded in test
files by `.golangci.yml`. This cost a red CI round on #1034. The `IsInt()` guard is a
tripwire: if the schema ever grows a genuine `TypeFloat`, the test fails loudly instead of
degrading silently. Porting to a resource that really has a float attribute means adding a
float branch **and** switching those assertions to `assert.InDelta` — not deleting the guard.

**`goValue` has no `tftypes.Map` and no `tftypes.Tuple` case** — the default branch
`t.Fatalf`s. `schema.TypeMap` is common in this provider, so a resource with one needs:

```go
	case valueType.Is(tftypes.Map{}):
		var elements map[string]tftypes.Value
		readValue(t, value, &elements)
		converted := make(map[string]any, len(elements))
		for name, element := range elements {
			converted[name] = goValue(t, element)
		}
		return converted
```

`tftypes.Value.As` accepts only `*string`, `*big.Float`, `*bool`, `*map[string]Value`,
`*[]Value` (and `**` variants). `tftypes` `Is()` is a plain kind check, so
`valueType.Is(tftypes.List{})` with a zero-value argument works.

---

## The stub

```go
// stubCloudAPI serves the <resource> collection the list resource reads, and returns
// the URL to point IONOS_API_URL at.
func stubCloudAPI(t *testing.T) string {
	t.Helper()

	<resource>s := <sdkpkg>.<Resource>s{
		Items: &[]<sdkpkg>.<Resource>{
			{
				// Result 1: every property the mapper reads is set.
				Id: new("00000000-0000-0000-0000-000000000001"),
				Properties: &<sdkpkg>.<Resource>Properties{ /* all fields */ },
			},
			{
				// Result 2: only the properties the API always returns, so that the
				// optional ones can be asserted null. Keep `name` set here even when the
				// schema makes it optional — the display-name and filter assertions below
				// are what make result 2 load-bearing, and they need a name to match on.
				// Omit some *other* optional attribute instead.
				Id: new("00000000-0000-0000-0000-000000000002"),
				Properties: &<sdkpkg>.<Resource>Properties{ /* required fields only */ },
			},
			{
				// Result 3 — ONLY when `name` is Optional in the SDKv2 schema. Covers the
				// displayName fallback from sdkv2-branch.md §2c, which result 2 cannot
				// reach because it keeps `name` set. Without this the fallback branch is
				// dead code the test never enters.
				Id: new("00000000-0000-0000-0000-000000000003"),
				Properties: &<sdkpkg>.<Resource>Properties{ /* required fields, Name unset */ },
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/<collection-path>") {
			http.NotFound(w, r)
			return
		}
		// The request the fetch closure builds is part of what is under test — the stub
		// answers any query string identically, so without these the options are unpinned.
		if got := r.URL.Query().Get("depth"); got != "1" {
			t.Errorf("expected depth=1 on the <collection-path> request, got %q", got)
		}
		// …and one per explicit .Limit(...)/other option the closure passes.
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(<resource>s); err != nil {
			t.Errorf("failed to write the stubbed response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return server.URL
}
```

- `new("prod")` is the **Go 1.26 built-in** `new(value)`, not a local helper. `go.mod` declares
  `go 1.26.3`, so it compiles. Copilot has claimed it does not (this exact exchange happened on
  #1034); it is wrong. There is no repo-local `ptr`/`toPtr` helper, but note the vendored SDK
  does export one — `ionoscloudsdk.ToPtr` (`sdk-go/v6/utils.go:30`), already imported by the
  reference test — so "no helper exists" is not the argument to make. Either compiles; `new`
  needs no import.
- **Assert the query params the fetch closure sets** — at minimum `depth`, and `limit`
  whenever decision 1 in `SKILL.md` made you pass an explicit one. Dropping `.Depth(1)` in
  production means every `Properties` is nil and the mapper skips every item; dropping
  `.Limit(...)` means silent truncation at the SDK's fallback (100 for `/ipblocks`,
  `sdk-go/v6/api_ip_blocks.go:487-489` — note that fallback is the *client's*, not the
  endpoint's). Neither is visible to a stub that ignores the query string, and the mutation
  check below never reaches the fetch closure.
- Match on `strings.HasSuffix`, not an exact path — the SDK prefixes `/cloudapi/v6`.
  Everything else 404s **on purpose**: an unexpected extra call (a pagination follow-up, a
  second region) surfaces as an error diagnostic rather than silently succeeding.
- A **regional** product fans out, so the handler must serve more than one path and the test
  must assert the union. A **paginated** one makes more than one request, so key off
  `r.URL.Query()`.
- **sdk-go-bundle products: `IONOS_API_URL` does not necessarily win.** Six product clients
  read `IONOS_API_URL_<PRODUCT>` as a real endpoint override — `services/cert/provider.go:27-36`,
  `services/vpn/client.go:75-86`, `services/objectstoragemanagement/accesskeys.go:44-54`, plus
  `nfs`, `dbaas/mariadb`, `dbaas/inmemorydb` — and `utils/loadedconfig/loadedconfig.go:51`
  *defers* `ChangeConfigURL`, so it overwrites whatever `IONOS_API_URL` set. Worse, a regional
  client called with a non-empty location replaces the endpoint from that product's
  `locationToURL` map — a **production** URL — regardless of either variable. So for a bundle
  product, set `IONOS_API_URL_<PRODUCT>` as well, and read the product's client constructor
  before assuming your stub is reachable at all.

---

## The three subtests

```go
	t.Run("streams every <resource>", func(t *testing.T) {
		results := listResults(ctx, t, server, <resource>ListType, listSchema, nil)
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}

		assert.Equal(t, "<display-name-1>", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "<id-1>", identity["id"])
		// one assert per identity attribute

		// Every attribute the mapper fills is asserted here, so that mapping a value
		// to the wrong attribute fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "<id-1>", resource["id"])
		// one assert per attribute the mapper sets; int attributes as int64(...)
		assert.ElementsMatch(t, []any{...}, resource["<set-attribute>"])   // TypeSet: order is not stable
		assert.Equal(t, []any{map[string]any{ /* nested block */ }}, resource["<nested-block>"])
		assert.Nil(t, resource["timeouts"], "a listed <resource> has no timeouts")

		// The second <resource> reports only the properties the API always sets, which
		// pins that the pairing holds past the first result and that the properties the
		// API left out stay null instead of turning into zero values.
		assert.Equal(t, "<display-name-2>", results[1].DisplayName)

		second := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "<id-2>", second["id"])
		assert.Nil(t, second["<optional-attribute>"])
		// one assert.Nil per omitted optional ATTRIBUTE. For an omitted nested BLOCK
		// (Optional ± Computed with an Elem: &schema.Resource{}) assert []any{} instead —
		// the framework reifies a null list/set block to an empty one. See §2d of
		// sdkv2-branch.md; ipblock's ip_consumers is the worked example.

		// Only when `name` is Optional: the display-name fallback, and the name still null.
		assert.Equal(t, "<id-3>", results[2].DisplayName)
		third := decode(t, results[2].Resource, resourceType)
		assert.Nil(t, third["name"])
	})

	// ONE SUBTEST PER FIELD in the FilterAttribute(...) allow-list. MatchesFilters returns
	// false for a field_name the mapper forgot to put in its map, so an untested field is a
	// filter that silently matches nothing, forever, while the other subtests stay green.
	t.Run("filters by <field1>", func(t *testing.T) {
		results := listResults(ctx, t, server, <resource>ListType, listSchema, map[string]string{"<field1>": "<value-matching-only-result-2>"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "<display-name-2>", results[0].DisplayName)
	})

	// … one more like it per remaining allowed field …

	// Two fields whose values match DIFFERENT stub items: pins that the filters are ANDed,
	// and that neither field is wired to the other's property.
	t.Run("applies every filter", func(t *testing.T) {
		results := listResults(ctx, t, server, <resource>ListType, listSchema, map[string]string{
			"<field1>": "<value-matching-result-1>",
			"<field2>": "<value-matching-result-2>",
		})
		if len(results) != 0 {
			t.Fatalf("expected 0 results, got %d", len(results))
		}
	})

	t.Run("rejects unknown filter fields", func(t *testing.T) {
		listServer, config := listServerAndConfig(t, server, listSchema, map[string]string{"nope": "value"})

		resp, err := listServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
			TypeName: <resource>ListType,
			Config:   &config,
		})
		if err != nil {
			t.Fatalf("ValidateListResourceConfig: %v", err)
		}
		if !hasErrorDiagnostic(resp.Diagnostics) {
			t.Fatalf("expected a validation error for an unknown filter field")
		}
	})
```

The third subtest is fully generic and fails if the list resource declares
`identity.FilterAttribute()` with no allowed-field list.

---

## What each assertion is actually catching

- **`ListResourceSchemas[type]` AND `ResourceSchemas[type]`, from one `GetProviderSchema`.**
  The first catches a list resource that is silently *not registered* — a typo in the type
  constant, or forgetting to wire `<service>.ListResources(p.sdkv2Provider)` into
  `provider.go`, produces a provider that starts fine and just has no list resource. The
  second catches the mirror problem: the merged schema must still serve the SDKv2 managed
  resource under the same name, pinning that the framework side registered *only* a list
  resource (`tf6muxserver` refuses duplicate type names).
- **`GetResourceIdentitySchemas`** catches a missing `Identity` on the SDKv2 resource, and
  proves it survives the protocol upgrade. The schema is then reused as the decode type for
  the identity, so a mismatch between the SDKv2 identity schema and the framework identity
  model surfaces as an unmarshal failure rather than a silent pass.
- **Asserting the SECOND result** catches three distinct bugs: *pairing drift* (item N's
  identity attached to item N's resource, not shifted or hoisted), *zero-value leakage*
  (an omitted *attribute* comes back null, not `""`/`false`/`[]` — the whole reason the model
  uses pointers; an omitted nested *block* is the documented exception and comes back `[]any{}`,
  see sdkv2-branch.md §2d), and *filter false-positives* (a filter that is silently a no-op fails on the
  count; one matching the wrong item fails on the display name).
- **`assert.Nil(t, resource["timeouts"], …)`** pins that the injected `timeouts` block exists
  in the model and is left null. Without the field in the model, `Resource.Set` fails at
  runtime.
- **`result.Identity == nil` / `result.Resource == nil` inside the driver** catch the
  `MappedItem.Identity must not be nil` contract violation and a mapper that ignores
  `includeResource`.

---

## Running it

```bash
go test ./internal/framework/services/<service>/ -run 'Test<Resource>ListResource' -count=1
```

Subtest names: spaces become underscores —
`-run 'Test<Resource>ListResource/applies_filters'`.

No `TF_ACC`, no `-tags`, no credentials.

---

## The mutation check — mandatory, not advice

The original #1034 test looked thorough and asserted almost nothing: its decoder returned
`map[string]string` and skipped every non-string attribute, so a mapper that swapped two int
fields still passed. **Prove your assertions are load-bearing before claiming coverage.**

Swap **two adjacent struct fields of the same Go type** in the mapper — same type, so it still
compiles and the failure has to come from an assertion, not the compiler.

List your mapper's fields by Go type, then mutate one same-typed adjacent pair **per distinct
type present**, preferring a pair inside a nested-block mapper (datacenter's `MaxCores` ↔
`MaxRam` is the ideal case: it exercises the nested path and the `int64` decode at once). Most
resources will not have that luxury — ipblock's nested `IpConsumer` is nine `*string` fields
and its only number, `size`, has no same-typed neighbour at all. Where a type has no pair,
verify it the other way: **change that field's value in the stub** and confirm the assertion
fails. A `*string` pair in the top-level mapper and the identity's own fields are always
available and always worth doing.

```bash
# 1. mutate one pair, then:
go test ./internal/framework/services/<service>/ -run 'Test<Resource>ListResource' -count=1
#    -> MUST fail, with the assertion diff naming the two swapped attributes.
# 2. revert:
git checkout -- internal/framework/services/<service>/resource_<resource>_list.go
# 3. confirm green again.
```

**Two mutations the field-swap cannot reach — run them too:**

- **Delete `.Depth(1)`** (and the explicit `.Limit(...)`, if there is one) from the fetch
  closure. The test MUST fail. The field-swap check only covers the mapper, so without this
  the request options are asserted by nothing.
- **Only for an optional-`name` resource:** delete the
  `if displayName == "" { displayName = *item.Id }` fallback. The result-3 display-name
  assertion MUST fail.

If step 1 **passes**, the assertion set is incomplete — an attribute the mapper fills is not
asserted, or is asserted against a value equal to its neighbour's. Add the missing assertion
before moving on. Do the same on a field only stub result 2 leaves unset, to confirm the
`assert.Nil` block is real too.

Report the mutation check as done only if you actually ran it and saw it fail.
