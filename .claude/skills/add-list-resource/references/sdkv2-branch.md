# SDKv2-backed branch

For a managed resource on `terraform-plugin-sdk/v2` (`ionoscloud/resource_*.go`). Reference:
`ionoscloud_datacenter`, **not to copy wholesale** (§2a, §2c). The list resource calls
`internal/framework/identity/` (`sdkv2.go`, `filter.go`, `list.go`).

## §1 — Resource identity on the SDKv2 resource

### 1a. The `Identity` block

Beside `Importer` in `resource<Resource>()`, shaped as `resource_datacenter.go:27-47`, comment
included. Each attribute: exactly one of `RequiredForImport`/`OptionalForImport`, a scalar or a
`TypeList` of scalars, only a `Description` (else `TestProvider$` fails). Required only if the
import cannot proceed without it (datacenter's `location` is Optional: `NewCloudAPIClient(ctx, "")`
falls back to the global endpoint). **Only what addresses the resource**: no `location` → a lone
`id`; a child adds a `RequiredForImport` string per parent ID, as its schema spells it, before
`id` (no shipped precedent: say so).

### 1b. The identity setter

`setDatacenterIdentity` (`resource_datacenter.go:337-355`): `d.Identity()`, then one
`identity.Set(...)` per attribute, each error wrapped with that name and `d.Id()`. **Exactly the
declared attributes** (an undeclared one panics under `TF_ACC`), **all of them on every read** (a
partial set fails "Unexpected Identity Change"). It reads state, not the API: it is also the list
resource's identity mapper (§2d).

### 1c. Call sites

Read (`resource_datacenter.go:167-195`): the 404 branch returns before it; the tail is writer,
*then* setter. Create too, unless it ends in Read; Delete needs nothing. **Update** not ending in
`return resource<Resource>Read(...)` calls the setter once the update succeeded, its comment giving
the real reasons: a safety net for state written before the identity existed (the SDK carries the
prior identity into the apply, so "Missing Resource Identity After Update" does not fire), equal to
the prior identity because no identity attribute can change in an Update (the id is
server-assigned, the rest ForceNew: check) — never "derived from state, so it cannot differ". It
also needs the acceptance update step (`test-harness.md`).

### 1d. The importer — dual mode, and the `d.SetId` that is easy to miss

An identity-based import sends an **empty `req.ID`**; the identifier is only in `d.Identity()`.
`resourceDatacenterImport` (`resource_datacenter.go:272-309`) and `datacenterImportParts`
(`:311-335`) are the shape: resolver → `d.SetId(id)` → client → fetch → writer → setter; the
identity branch (`:316-322`) is **a pure prepend** to the legacy string path (delimiter, part count,
error message kept), which a non-nil `identityErr` falls back to. A child
`d.Set("<parent>_id", …)`s before the setter. The `d.SetId` comment names only what reads `d.Id()`
below it (not datacenter's "and the log line below" unless a log line does).

**Additive only.** Prepend the resolver and its `d.SetId(id)`, append the setter after the writer,
leave every other line as found, **including a post-fetch `d.SetId(*obj.Id)` now redundant**; no
re-worded comment or tidying anywhere in the file.

**The empty-ID guard.** `identity.GetOk("id")` reports `exists=false` for `""`, so an empty identity
`id` falls through to the legacy path, where `splitImportID` (`ionoscloud/utils.go`) yields `[""]`
and `validateImportIDParts` rejects it (Copilot's "missing validation" there is a false positive).
**A plain import ID never calls `splitImportID`**: `""` reaches the API. Such a resolver adds the
guard, commented with both inputs it catches (an empty import string, an identity `id` of `""`):
`if d.Id() == "" { return "", fmt.Errorf("invalid import identifier: expected a <resource> UUID, got an empty string") }`

## §2 — The list resource

`ionoscloud/resource_<resource>_list.go`, package `ionoscloud`, beside the resource.

### 2a. Header, assertions, struct, constructor

`resource_datacenter_list.go:3-53`, its three `var _ list.X = (*T)(nil)` assertions (`:34-38`)
load-bearing (§2b). Changes:

- **`constant.<Resource>Resource`**, the `ResourcesMap` key (`ionoscloud/provider.go`), not a local
  type-name const (datacenter's `:32` is a divergence, not a precedent).
- **The constructor takes no arguments**: `New<Resource>ListResource() list.ListResource` returns
  `&<resource>ListResource{resourceSchema: resource<Resource>()}`, never rebuilt per item.
- **Import `internal/framework/identity` as `fwidentity`**: `identity` is the SDKv2
  `*schema.IdentityData` here.

### 2b. The four boilerplate methods

`resource_datacenter_list.go:55-94`, the constant in `RawV6Schemas` and `Metadata` (no
`req.ProviderTypeName` prefix). `Configure`'s `req.ProviderData == nil` early return is mandatory.
Filters: `fwidentity.FiltersKey: fwidentity.FilterAttribute("<f1>", "<f2>")` (no arguments: nothing
validated). `Metadata`/`Configure` take `resource.MetadataRequest`/`resource.ConfigureRequest`,
**not** the vendored `list.*` decoys, a compile error only thanks to the `var _` assertions: drop
`ListResourceWithConfigure`'s and `Configure` never runs, so `List` nil-derefs `r.bundle`.

### 2c. `List` and the mapper

`resource_datacenter_list.go:96-170` is the Cloud API (`sdk-go/v6`) shape: `List` wraps
`fwidentity.StreamList(ctx, stream, req, fetch, r.map<Resource>)`; the mapper nil-guards, filters,
then does §2d. **Every check in it stays, whichever SDK**: `apiResponse` named, never `_`, and
logged by `:112-114`'s `tflog.Debug`; each `err` returned (`:104-117`) or added as a diagnostic
(`:150-166`). Only delta 3 retires a nil guard, on a value that cannot be nil. Deviations:

- **Not `:96`'s "fetches every datacenter on the contract"**: the `List` and fetch comments name the
  one request and its bound (up to `n`, or the API's default page of `n`), as the docs note does;
  never "every", "the whole contract", "covers the whole contract".
- **`MatchesFilters` values inline, from the accessors**: no local lifted out (`:140-141`: not the
  pattern); pointers through `shared.ToValueDefault`, **never a local copy of it**.
- **`.Limit(n)` when decision 1 says so**, the fetch comment naming `n` and its source.
- **Filters stay in the mapper**, whatever server-side filter a bundle helper offers.
- **A child** sets the parent ids, from the fetch context, *before* the setter: the one mapping of
  your own, only for ids the writer provably leaves unset.
- **One spelling of the thing** in every message you add, as the resource's docs write it.

#### The sdk-go-bundle deltas

If CRUD takes a `<Product>Client` off the bundle (step 0's (6)):

1. **The resource's own client**; prefer the list helper the sibling data source uses
   (`r.bundle.DNSClient.ListZones(ctx, "")`), its `apiResponse` still logged as `:112-114` even
   when the helper already calls `apiResponse.LogInfo()`.
2. **`Depth`**: `scripts/sdkv2.sh depth <product> <api-stem> <RequestType>`; most have none.
3. **Values, not pointers**: `Items` and each element's `Id`/`Properties` are values, so the nil
   guards and the `*items.Items` deref go; an empty id is the only unusable item. Fields *inside*
   `Properties` still mix. The product package replaces `sdk-go/v6`.
4. **The writer may be a method on that client, taking its object by value**
   (`r.bundle.DNSClient.SetZoneData(data, zone)`): no `&`.

#### `DisplayName`

It labels each `terraform query` row; blank, the row is unlabelled. A `Required` `name`: the
accessor straight in, no local. An `Optional` one: a `displayName` local from the accessor
(`shared.ToValueDefault(item.Properties.Name)`), set to `*item.Id` when empty. **No `Name` at all**
(`ionoscloud_cdn_distribution`: `Domain`; a user: `Email`): the closest human identifier, reported
at step 0.5.

**`StreamList` semantics** (`list.go`): `(nil, nil)` skips silently; `(nil, diags-with-error)`
skips, logging `diags[0].Detail()` at WARN only; a non-nil item with an error diagnostic, or with
`Identity == nil`, **aborts the whole stream**. So `return nil, diags` for a recoverable per-item
problem, and `return mapped, diags` to include, as the reference ends (no error in them there; a
warning reaches the user). **Never `(mapped, diags-with-error)`**
(`objectstorage/resource_bucket_list.go:76,:80`). `fetch` runs once, before streaming: pagination
lives inside it; `req.Limit` is ignored.

### 2d. No model — the result is the resource's own state

```go
	data := r.resourceSchema.Data(&terraform.InstanceState{})
	set<Resource>Data(data, &item)     // the resource's OWN state writer, the one Read uses
	set<Resource>Identity(data)        // AFTER the writer — it reads values back out
	fwidentity.MappedItemFromResourceData(displayName, data, includeResource)
```

**No Go struct mirroring the schema**: it compiles, and it is the design that shipped `master`
broken.

- **`Data`, not `TestResourceData`**, which never assigns `timeouts`, so `Set` rejects it.
- **A fresh `ResourceData` per item**, or one item's values leak into what the writer skips.
- **The writer calls `d.SetId`** (check; else `data.SetId(<id>)` first), or both conversions fail:
  `"state is nil, call SetId() on ResourceData first"`.

`timeouts` comes back null. Omitted values decode as: a scalar or primitive-list attribute → nil;
a nested `*schema.Resource` that is `Optional` (± `Computed`) is a block → `[]any{}`,
`Computed`-only is an attribute → nil; a list `d.Set()` from a nil slice → `[]any{}`. `MaxItems: 1`
changes nothing.

## §3 — Registration

`New<Resource>ListResource` in `ionoscloud/list_resources.go`'s slice, named directly: **no
closure, lookup, `ok` branch or nil guard**. Nothing else to wire (`provider.New` takes the whole
slice). **A list resource without schemas** — no `Identity` block (§1) — fails
`GetProviderSchema`, terraform's first RPC ("ListResource Type Defined without a Matching Managed
Resource Type"), killing the whole provider.
