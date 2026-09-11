# Test harness

Reference: `ionoscloud/resource_datacenter_list_test.go` (package `ionoscloud_test`). It drives
the real `ListResource` RPC end to end against a stubbed Cloud API, through the same mux
`main.go` serves — no acceptance test, no credentials, and it runs in well under a second. It is
the older of the two list-resource tests in `ionoscloud/` (datacenter, ipblock); the five
framework-native list resources under `internal/framework/services/` have none. Copy it, and
read `ionoscloud/resource_ipblock_list_test.go` too — it is the smaller worked example, and the
one that covers an `Optional` `name`, a two-attribute identity and an omitted nested block.

New file: `ionoscloud/resource_<resource>_list_test.go`, package `ionoscloud_test` — beside the
list resource, which lives beside the SDKv2 resource it lists. The **external** test package is
what both existing tests use: it keeps names as generic as `decode`, `listResults` and
`dynamicValue` out of package `ionoscloud`'s own namespace, and it forces the test through the
same exported surface `main.go` uses, `ionoscloud.Provider()` and `ionoscloud.ListResources()`.

No `//go:build` tag. That is deliberate — it keeps the test cheap to run and puts it in
golangci-lint's scope (which only analyses the default build). `ionoscloud/` has 114 test files
and only nine of them are untagged, so an untagged file costs nothing to run and, unlike the
other 105, needs no `TF_ACC` and no credentials.

---

## What this test pins — there is no model any more

The list resource declares no struct mirroring the SDKv2 schema. Each result is

```go
	data := r.resourceSchema.Data(&terraform.InstanceState{})
	set<X>Data(data, &item)
	set<X>Identity(data)
	fwidentity.MappedItemFromResourceData(displayName, data, includeResource)
```

so the code under test is the resource's **own state writer** (`set<X>Data` above is a
placeholder: `setDatacenterData` is unexported, `IpBlockSetData` is exported — use whichever
name the resource already has), the one `resource<X>Read` calls,
plus the conversion in `internal/framework/identity/sdkv2.go`. What the assertions pin is
therefore not "the resource model fills the SDKv2 schema without a type mismatch" — nothing
fills anything — but: **the state the resource's own writer produces survives the round trip
through `ResourceData.TfTypeResourceState` into the list-result schema**, including the
`timeouts` block that `MappedItemFromResourceData` nulls back out
(`internal/framework/identity/sdkv2.go:82-89`, `nullTimeouts` at `:99-121`).

Put that in the test's doc comment. Datacenter's says exactly it
(`ionoscloud/resource_datacenter_list_test.go:26-35`), and ipblock's repeats it and then adds a
pointer to where the shared helpers live. If you meet a header still saying "the resource model
fills the SDKv2 schema", that is the pre-refactor wording — do not copy it forward.

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

	// ... the subtests
}
```

Set the env vars **before** building the server: the framework provider's `Configure` reads
`IONOS_API_URL` into `clientOptions.Endpoint`, and `bundleclient.shouldApplyOverridesCustomEnv`
short-circuits to `false` whenever it is set — so the developer's real `~/.ionos/config`
overrides are ignored and the test cannot accidentally hit a live endpoint.

Declare `const <resource>ListType = "<ionoscloud_type>"` at package level (`goconst` flags a
literal repeated 5+ times).

---

## The shared helpers already exist — copy none of them

**Check first:**

```bash
ls ionoscloud/*_list_test.go
```

For the SDKv2 branch the answer is always "they are there". Every list resource for an SDKv2
managed resource lands in package `ionoscloud` — there is no per-service split on this side of
the mux — so there is exactly **one** test package for all of them, and the helpers are already
declared in it, below `TestDatacenterListResource` in `resource_datacenter_list_test.go`. They
are package-level declarations shared by every file in `ionoscloud_test`, and re-declaring them
is 14 `redeclared in this block` errors.

The only things you write per resource are:

- `const <resource>ListType = "<ionoscloud_type>"`
- `Test<Resource>ListResource`
- `stub<Resource>API` — **rename it.** The existing `stubCloudAPI` is pinned to `/datacenters`
  (`ionoscloud/resource_datacenter_list_test.go:212`); reusing the name is the one collision
  that does not announce itself, because a same-named stub compiles and then 404s every request
  your resource makes.
- the subtests

**Do not copy the driver.** `listResults(ctx, t, server, typeName, schema, filters)` already
takes the type name as a parameter for exactly this reason: the type name is the *only* thing
that differs per resource, so call the shared one
(`ionoscloud/resource_datacenter_list_test.go:248`). If your resource needs something the
shared driver does not do yet, add the parameter to it in place — that is exactly what turned
it from a datacenter-only helper into `listResults` — rather than growing a second
near-identical copy.

`resource_ipblock_list_test.go` is the worked example of this: a type-name const, one
`Test<X>ListResource`, one `stub<X>API` and its subtests, and nothing else. Its header comment
says so ("The shared helpers it calls … are declared once for the package in
resource_datacenter_list_test.go") — keep that line, it is what stops the next author
re-copying them.

Do **not** change the reference test's *assertions or stub* — those are what make
`resource_datacenter_list_test.go` the reference. Adding to, or parameterising, the shared
helpers it calls is fine and expected.

**Where the shared helpers live.** They are in `resource_datacenter_list_test.go` for
historical reasons, so every author in `ionoscloud/` has to learn from a comment that the
infrastructure lives in another resource's file. Do not reproduce that in a *new* package: for
a framework-native list resource in `internal/framework/services/<service>/` (which has no list
test at all today), put them in `list_resource_test_helpers_test.go` from the start, never
inside a file named after one resource. That copy has to import
`.../v6/ionoscloud` for `Provider()` and `ListResources()` — legal from an **external**
`<service>_test` package, since `ionoscloud`'s non-test files import nothing under
`internal/framework/services/` (only `internal/framework/identity`, `internal/serverutil` and
`internal/tf/writeonly`); a cycle from an in-package one.

The table below is for that **new** package. Take them verbatim from
`ionoscloud/resource_datacenter_list_test.go`:

| Helper | What it does |
|---|---|
| `muxedProviderServer(ctx, t)` | `ionoscloud.Provider()` → `tf5to6server.UpgradeServer` → `tf6muxserver.NewMuxServer(ctx, providerserver.NewProtocol6(fwprovider.New(ionoscloud.ListResources()...)), upgraded)` — mirrors `main.go:31-47` |
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
terraform-plugin-go, `resource_datacenter_list_test.go:280`) and the `//nolint:prealloc` on the
results slice (`:262`). `nolintlint` requires directives to be both specific and explained; a
bare `//nolint` is itself an error.

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
// stub<Resource>API serves the <resource> collection the list resource reads, and returns
// the URL to point IONOS_API_URL at. Every other path 404s on purpose, so an unexpected
// extra request surfaces as an error diagnostic instead of succeeding.
func stub<Resource>API(t *testing.T) string {
	t.Helper()

	<resource>s := ionoscloudsdk.<Resource>s{
		Items: &[]ionoscloudsdk.<Resource>{
			{
				// Result 1: every property the writer reads is set.
				Id: new("00000000-0000-0000-0000-000000000001"),
				Properties: &ionoscloudsdk.<Resource>Properties{ /* all fields */ },
			},
			{
				// Result 2: only the properties the API always returns, so that the
				// optional ones can be asserted null. Keep `name` set here even when the
				// schema makes it optional — the display-name and filter assertions below
				// are what make result 2 load-bearing, and they need a name to match on.
				// Omit some *other* optional attribute instead.
				Id: new("00000000-0000-0000-0000-000000000002"),
				Properties: &ionoscloudsdk.<Resource>Properties{ /* required fields only */ },
			},
			{
				// Result 3 — ONLY when `name` is Optional in the SDKv2 schema. Covers the
				// displayName fallback (`ionoscloud/resource_ipblock_list.go:166-171`),
				// which result 2 cannot reach because it keeps `name` set. Without this
				// the fallback branch is dead code the test never enters.
				Id: new("00000000-0000-0000-0000-000000000003"),
				Properties: &ionoscloudsdk.<Resource>Properties{ /* required fields, Name unset */ },
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
		//
		// On a fetch that sets NOTHING - an sdk-go-bundle collection has no depth, and
		// decision 1 may leave the limit unset - assert the ABSENCE instead, or the
		// request is entirely unpinned and a pushed-down filter added later goes
		// unnoticed:
		//   for _, param := range []string{"depth", "limit", "offset", "<filter-param>"} {
		//       if got := r.URL.Query().Get(param); got != "" {
		//           t.Errorf("expected no %s on the <collection-path> request, got %q", param, got)
		//       }
		//   }
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(<resource>s); err != nil {
			t.Errorf("failed to write the stubbed response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return server.URL
}
```

- **An sdk-go-bundle stub must fill every non-omitempty enum field.** Bundle models validate
  enums while unmarshalling, so a zero-valued one fails the *fetch*, not an assertion, and the
  test dies with a diagnostic that names neither the field nor the item — the `ionoscloud_dns_zone`
  stub first failed with `failed to list dns zones:  is not a valid ProvisioningState` because
  `ZoneRead.Metadata.State` was left unset. Set it on **every** stub item, including the minimal
  one whose whole point is that its optional fields are absent. The Cloud API models do not do
  this, so nothing in the datacenter or ipblock tests prepares you for it.
- `new("prod")` is the **Go 1.26 built-in** `new(value)`, not a local helper. `go.mod` declares
  `go 1.26.3`, so it compiles. Copilot has claimed it does not (this exact exchange happened on
  #1034); it is wrong. There is no repo-local `ptr`/`toPtr` helper, but note the vendored SDK
  does export one — `ionoscloudsdk.ToPtr` (`sdk-go/v6/utils.go:30`), already imported by the
  reference test — so "no helper exists" is not the argument to make. Either compiles; `new`
  needs no import.
- **Assert the query params the fetch closure sets** — at minimum `depth`, and `limit`
  whenever decision 1 in `SKILL.md` made you pass an explicit one. Dropping `.Depth(1)` in
  production means every `Properties` is nil and the writer is never reached; dropping
  `.Limit(...)` means silent truncation at the SDK client's fallback (100 for `/ipblocks`,
  `sdk-go/v6/api_ip_blocks.go:487-489` — note that fallback is the *client's*, not the
  endpoint's). Neither is visible to a stub that ignores the query string, and the mutation
  check below never reaches the fetch closure.
  The reference stub is the wrong model for this one thing: `stubCloudAPI` asserts no query
  parameters at all. Copy the assertion block from
  `ionoscloud/resource_ipblock_list_test.go:245-255` instead, including the comment
  explaining why the limit is asserted as a literal rather than as `constant.<X>Limit` — that
  literal is the only thing tying `docs/list-resources/<resource>.md` to the code.
- Match on `strings.HasSuffix`, not an exact path — the SDK prefixes `/cloudapi/v6`.
  Everything else 404s **on purpose**: an unexpected extra call (a pagination follow-up, a
  second region) surfaces as an error diagnostic rather than silently succeeding.
- A **regional** product fans out, so the handler must serve more than one path and the test
  must assert the union. A **paginated** one makes more than one request, so key off
  `r.URL.Query()`.
- **sdk-go-bundle products: `IONOS_API_URL` does not necessarily win.** Several product clients
  read `IONOS_API_URL_<PRODUCT>` as a real endpoint override (`grep -rn IONOS_API_URL_ services/`
  lists them; `services/cert/provider.go:27-36`, `services/vpn/client.go:75-86` and
  `services/objectstoragemanagement/accesskeys.go:44-54` are the shapes it takes) — and
  `utils/loadedconfig/loadedconfig.go:50`
  *defers* `ChangeConfigURL`, so it overwrites whatever `IONOS_API_URL` set. Worse, a regional
  client called with a non-empty location replaces the endpoint from that product's
  `locationToURL` map — a **production** URL — regardless of either variable. So for a bundle
  product, set `IONOS_API_URL_<PRODUCT>` as well, and read the product's client constructor
  before assuming your stub is reachable at all.

---

## The subtests

```go
	t.Run("streams every <resource>", func(t *testing.T) {
		results := listResults(ctx, t, server, <resource>ListType, listSchema, nil)
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}

		assert.Equal(t, "<display-name-1>", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "<id-1>", identity["id"])
		// one assert per identity attribute, plus an assert.Len on the identity map when
		// its shape is the point (ipblock's is id+location, so its test asserts Len 2)

		// Every attribute the writer fills is asserted here, so that writing a value to
		// the wrong key fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "<id-1>", resource["id"])
		// one assert per attribute set<X>Data sets; int attributes as int64(...)
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
		// one assert.Nil per omitted optional ATTRIBUTE - but check the WRITER before
		// assuming nil. Two things produce []any{} instead:
		//   - an omitted nested BLOCK (Optional ± Computed with an Elem: &schema.Resource{}):
		//     the framework reifies a null list/set block to an empty one;
		//   - a list ATTRIBUTE the writer d.Set()s unconditionally with a nil slice, which
		//     materialises an EMPTY list rather than leaving the attribute null. SetZoneData
		//     does exactly this with `nameservers`, so on ionoscloud_dns_zone that one
		//     attribute decodes to []any{} while description and enabled decode to nil.
		// Schema type alone does not settle it; what the writer does is half the answer.

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

The last subtest is fully generic and fails if the list resource declares
`identity.FilterAttribute()` with no allowed-field list. When a *deliberate* non-filterable
attribute exists, loop that subtest's body over both it and the bogus field
(`for _, field := range []string{"nope", "<not-filterable>"}`), so that widening the allow-list
cannot pass unnoticed. Neither existing test needs that today: datacenter and ipblock both
allow every field their mapper puts in the `MatchesFilters` map.

---

## What each assertion is actually catching

- **`ListResourceSchemas[type]` AND `ResourceSchemas[type]`, from one `GetProviderSchema`.**
  The first catches a list resource that is silently *not registered* — a typo in the type
  constant, or forgetting to add `New<Resource>ListResource` to `ListResources()` in
  `ionoscloud/list_resources.go`, produces a provider that starts fine and just has no
  list resource. The second catches a double registration: the merged schema must still serve
  the SDKv2 managed resource under the same name, pinning that the framework side registered
  *only* a list resource (`tf6muxserver` refuses duplicate type names).
- **`GetResourceIdentitySchemas`** catches a missing `Identity` on the SDKv2 resource, and
  proves it survives the protocol upgrade. The schema is then reused as the decode type for
  the identity, so a mismatch between the identity schema and what `set<X>Identity` writes
  surfaces as an unmarshal failure rather than a silent pass.
- **The per-attribute asserts are what replaced the deleted model.** Nothing in the build
  restates a `<resource>`'s shape any more, so if `set<X>Data` writes a value under the wrong
  key there is no compiler, no `tfsdk` tag and no reviewer diff to catch it — only these
  assertions and `Set`'s own type check on the way in
  (`vendor/.../terraform-plugin-framework/internal/fwschemadata/data_set.go:21-36`, the
  `tftypes.Value` fast path `MappedItemFromResourceData` is written to hit). Assert **every**
  attribute the writer fills.
- **Asserting the SECOND result** catches three distinct bugs: *pairing drift* (item N's
  identity attached to item N's resource, not shifted or hoisted), *zero-value leakage*
  (an omitted *attribute* comes back null, not `""`/`false`/`[]` — an omitted nested *block*
  is the documented exception and comes back `[]any{}`), and *filter false-positives* (a filter
  that is silently a no-op fails on the count; one matching the wrong item fails on the display
  name).
- **`assert.Nil(t, resource["timeouts"], …)`** pins `nullTimeouts`. `Resource.Data` always
  attaches the resource's `Timeouts` (`vendor/.../helper/schema/resource.go:1414-1418`), and
  the flatmap shims cannot tell a null single block from an empty one, so
  `TfTypeResourceState` hands back a `timeouts` object of null attributes;
  `MappedItemFromResourceData` nulls the whole block back out
  (`internal/framework/identity/sdkv2.go:82-89`), exactly as SDKv2's own `ReadResource` does.
  Drop that step and a query result carries a timeouts object that a refresh never writes —
  this assertion is the only thing that would notice.
- **`result.Identity == nil` / `result.Resource == nil` inside the driver** catch the
  `MappedItem.Identity must not be nil` contract violation
  (`internal/framework/identity/list.go:59-63`) and a mapper that ignores `includeResource`.

---

## Running it

```bash
go test ./ionoscloud/ -run 'Test<Resource>ListResource' -count=1
```

Subtest names: spaces become underscores —
`-run 'Test<Resource>ListResource/applies_every_filter'`.

No `TF_ACC`, no `-tags`, no credentials. Nothing in CI runs this for you — see
`verify-and-pr.md` §1.

---

## The mutation check — mandatory, not advice

The original #1034 test looked thorough and asserted almost nothing: its decoder returned
`map[string]string` and skipped every non-string attribute, so a writer that swapped two int
fields still passed. **Prove your assertions are load-bearing before claiming coverage.**

The thing to mutate is no longer a mapper of the list resource's own. It is `set<X>Data`, in
`ionoscloud/resource_<resource>.go` — **a different file from the one you just wrote**, and the
same file that holds the `Identity` block and `set<X>Identity` you added in step 1. So:

> **Revert the mutation by hand.** `git checkout -- ionoscloud/resource_<resource>.go` would
> throw the whole identity implementation away with it.

Swap **two adjacent struct fields, or two adjacent map keys, of the same Go type** — same type,
so it still compiles and the failure has to come from an assertion, not the compiler.

List the writer's fields by Go type, then mutate one same-typed adjacent pair **per distinct
type present**, preferring a pair inside a nested block (datacenter's `max_cores` ↔ `max_ram`
in the `cpu_architecture` entry, `ionoscloud/resource_datacenter.go:415-421`, is the ideal
case: both are `*int32`, it exercises the nested path and the `int64` decode at once). Most
resources will not have that luxury — ipblock's nested `IpConsumer` is nine `*string` fields
and its only number, `size`, has no same-typed neighbour at all. Where a type has no pair,
verify it the other way: **change that field's value in the stub** and confirm the assertion
fails. A `*string` pair in the top-level writer is always available and always
worth doing, as is swapping the two *values* in `set<X>Identity`
(`identity.Set("id", d.Get("location"))` and vice versa) when the identity has two attributes
— datacenter's and ipblock's are both `id` + `location`. A lone-`id` identity has no pair, so
change its value in the stub instead.

```bash
# 1. mutate one pair in ionoscloud/resource_<resource>.go, then:
go test ./ionoscloud/ -run 'Test<Resource>ListResource' -count=1
#    -> MUST fail, with the assertion diff naming the two swapped attributes.
# 2. revert BY HAND (see the warning above), then
# 3. confirm green again.
```

**Four mutations the field-swap cannot reach — they live in the list resource file, which
`set<X>Data` mutations never touch. Run them too:**

- **Delete `.Depth(1)`** (and the explicit `.Limit(...)`, if there is one) from the fetch
  closure. The test MUST fail — on the stub's query assertions, which is why they are not
  optional.
- **Swap the two values in the `MatchesFilters` map** (`"name": location, "location": name`).
  Only the per-field filter subtests catch this; with a single filter subtest it passes.
- **Swap the `set<X>Data` and `set<X>Identity` calls.** The identity setter reads its values
  back out of the `ResourceData`, so running it first leaves the identity empty. The identity
  assertions on `results[0]` are the only thing that fails.
- **Only for an optional-`name` resource:** delete the
  `if displayName == "" { displayName = *item.Id }` fallback. The result-3 display-name
  assertion MUST fail.

If step 1 **passes**, the assertion set is incomplete — an attribute the writer fills is not
asserted, or is asserted against a value equal to its neighbour's. Add the missing assertion
before moving on. Do the same on a field only stub result 2 leaves unset, to confirm the
`assert.Nil` block is real too.

Report the mutation check as done only if you actually ran it and saw it fail.
