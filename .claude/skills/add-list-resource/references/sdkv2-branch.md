# SDKv2-backed branch

For a managed resource still implemented with `terraform-plugin-sdk/v2` (anything in
`ionoscloud/resource_*.go`). Reference implementation, read it before writing anything:

- `ionoscloud/resource_datacenter.go` — the identity side
- `internal/framework/services/compute/resource_datacenter_list.go` — the list resource
- `internal/framework/services/compute/resources.go` — registration

## Why this is not just "a framework list resource"

List resources cannot be implemented in SDKv2. So the list resource lives on the framework
half of the mux while the managed resource stays on the SDKv2 half, and the framework has to
be handed the managed resource's protocol schemas by hand. That is the
`list.ListResourceWithRawV6Schemas` interface, and `identity.SetRawV6Schemas` fills it from
the SDKv2 `*schema.Resource` via `ProtoSchema` / `ProtoIdentitySchema` plus
`terraform-plugin-mux/tf5to6server/translate`.

Two consequences that look like mistakes and are not:

- **The `*schema.Provider` parameter on `ListResources` cannot be removed.** Package
  `<service>` cannot import package `ionoscloud`: `ionoscloud/provider_test.go` is an
  in-package test that imports the framework provider, so the reverse import is a cycle.
  Exporting the resource factory does not help. The `*schema.Resource` has to arrive through
  `sdkv2Provider.ResourcesMap`. If someone proposes "cleaning this up", this is the answer.
- **Do not obtain the schemas via `tf5to6server.UpgradeServer` + `GetProviderSchema`.** It
  works, but it stands up a second gRPC server and builds protocol schemas for every resource
  and data source to keep one. The original #1034 did this and it was replaced. Do not copy
  terraform-provider-aws's `ListResourceWithSDKv2Resource` *struct* either — that base exists
  to carry registry-injected state and to mutate the schema; with none of that, a plain
  function is enough.

---

## §1 — Resource identity on the SDKv2 resource

**Hard prerequisite.** A resource with no `Identity` cannot be listed:
`resourceSchema.ProtoIdentitySchema(ctx)` returns a **nil func** (not a func returning nil),
`SetRawV6Schemas` gives up, and `RawV6Schemas` has no diagnostics channel — the only trace is
a `tflog.Error` line. Debugging that from the terraform side is miserable.

### 1a. The `Identity` block

```go
func resource<Resource>() *schema.Resource {
	return &schema.Resource{
		CreateContext: resource<Resource>Create,
		ReadContext:   resource<Resource>Read,
		UpdateContext: resource<Resource>Update,
		DeleteContext: resource<Resource>Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resource<Resource>Import,
		},
		// The identity is what a `list "<ionoscloud_type>"` block streams back for
		// each <resource> it finds, and what an import block can be written against.
		// Terraform requires every read of a resource that declares an identity to
		// return one, see set<Resource>Identity.
		Identity: &schema.ResourceIdentity{
			Version: 0,
			SchemaFunc: func() map[string]*schema.Schema {
				return map[string]*schema.Schema{
					"id": {
						Type:              schema.TypeString,
						RequiredForImport: true,
						Description:       "The UUID of the <resource>.",
					},
					"location": {
						Type:              schema.TypeString,
						OptionalForImport: true,
						Description:       "The location the <resource> lives in. Only needed when the Cloud API endpoint is overridden per location.",
					},
				}
			},
		},
		Schema:   map[string]*schema.Schema{ /* unchanged */ },
		Timeouts: &resourceDefaultTimeouts,
	}
}
```

Constraints enforced by `(*ResourceIdentity).InternalIdentityValidate()` — and already
checked by `ionoscloud/provider_test.go`'s `Provider().InternalValidate()`, so a malformed
identity fails a *unit* test, not an acceptance run:

- exactly one of `RequiredForImport` / `OptionalForImport` per attribute — never both, never
  neither;
- only `TypeBool`, `TypeFloat`, `TypeInt`, `TypeString`, and `TypeList` of those.
  `TypeMap` and `TypeSet` are rejected;
- no `ForceNew`, `Required`, `Optional`, `Computed`, `Sensitive`, validators or defaults.
  `Description` is allowed.

`RequiredForImport` vs `OptionalForImport`: mark it Required only if the import genuinely
cannot proceed without it. Datacenter's `location` is Optional because
`bundleclient.SdkBundle.NewCloudAPIClient(ctx, "")` tolerates an empty location and falls back
to the global endpoint. A product whose client cannot be built without a location makes it
Required — that is what the framework-native clusters do.

**Child resource** — one attribute per parent ID, named as the resource's own schema names
them, all Required:

```go
			Identity: &schema.ResourceIdentity{
				Version: 0,
				SchemaFunc: func() map[string]*schema.Schema {
					return map[string]*schema.Schema{
						"datacenter_id": {Type: schema.TypeString, RequiredForImport: true, Description: "..."},
						"server_id":     {Type: schema.TypeString, RequiredForImport: true, Description: "..."},
						"id":            {Type: schema.TypeString, RequiredForImport: true, Description: "..."},
						"location":      {Type: schema.TypeString, OptionalForImport: true, Description: "..."},
					}
				},
			},
```

No child resource in the repo declares an identity yet — this is the extrapolation from the
existing composite import IDs, not shipped code. Say so if you build one.

### 1b. The identity setter

```go
// set<Resource>Identity writes the resource identity from the <resource> already in
// state. Terraform errors out with "Missing Resource Identity After Read" if a
// resource that declares an identity finishes a read without returning one.
func set<Resource>Identity(d *schema.ResourceData) error {
	identity, err := d.Identity()
	if err != nil {
		return err
	}

	if err := identity.Set("id", d.Id()); err != nil {
		return fmt.Errorf("error while setting id identity attribute for <resource> %s: %w", d.Id(), err)
	}

	if err := identity.Set("location", d.Get("location")); err != nil {
		return fmt.Errorf("error while setting location identity attribute for <resource> %s: %w", d.Id(), err)
	}

	return nil
}
```

**Set every attribute, on every read, forever.** SDKv2's `ReadResource` compares the whole
identity object with `RawEquals` against the previous one and fails with "Unexpected Identity
Change" on any difference. Setting only some attributes avoids the "missing identity" error
(the check is "null or *every* attribute null") but walks into the stability check instead.

### 1c. Call sites

```go
func resource<Resource>Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// ... build client, GET the object

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")   // drift: no identity needed, terraform short-circuits on an empty ID
			return nil
		}
		return diagutil.ToDiags(d, err, nil)
	}

	if err := set<Resource>Data(d, &<resource>); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	// must run AFTER set<Resource>Data: the identity reads attributes that only the
	// data setter fills in (this matters most on an identity-based import).
	if err := set<Resource>Identity(d); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}
```

- The 404 branch correctly returns **before** setting an identity: `grpc_provider.go` returns
  a null state early when the ID is empty, before any identity handling.
- **Create** always needs the same call unless it ends with
  `return resource<Resource>Read(ctx, d, meta)`. There is no prior identity to inherit
  (`create == true`, `grpc_provider.go:1571-1573`), so without it `ApplyResourceChange` fails
  with "Missing Resource Identity After Create".
- **Update is different, and the reason is easy to get wrong.** (The comment that got it wrong
  was on `resourceIPBlockUpdate` and has since been corrected — read it there for the right
  wording; datacenter's Update delegates to Read and carries no such comment.) The SDK carries the prior
  identity into the apply on its own — `grpc_provider.go:1425` → `:1461` →
  `resource_data.go:856-857` → `:521-544` — so an update that never touches the identity still
  returns one and "Missing Resource Identity After Update" does *not* fire. The call is a
  safety net for state written before the resource declared an identity (e.g. `-refresh=false`
  on an old state file), not the everyday path. Say that in the comment; do not copy the
  Create rationale onto Update.
- **Whatever Update writes must equal the prior identity byte for byte.**
  `grpc_provider.go:1585` compares the returned identity against the planned one with
  `RawEquals` and fails with "Unexpected Identity Change" unless `ResourceBehavior`
  sets `MutableIdentity`. Derive from state (`d.Get(...)`), never from an API response.
- **Prefer ending Update with `return resource<Resource>Read(ctx, d, meta)`** where the
  resource has computed attributes to refresh — it solves both points at once and is what
  datacenter does. Only add the explicit write when the update genuinely cannot delegate,
  as `resourceIPBlockUpdate` cannot.
- If Update does *not* delegate to read, the acceptance test needs an update step that changes
  a **non-ForceNew** attribute only, or `resource<Resource>Update` is never entered at all.
  Check the resource's existing `*ConfigUpdate` fixture first: `ionoscloud/resource_ipblock_test.go`
  flips `size`, which is `ForceNew`, so that step is a destroy-and-create and covers nothing.
- Delete needs nothing.

### 1d. The importer — dual mode, and the `d.SetId` that is easy to miss

Terraform sends an **empty `req.ID`** for an identity-based import; the real identifier is
only in `d.Identity()`. So the importer must look at the identity first, and must set the ID
itself before anything that reports on the resource.

```go
func resource<Resource>Import(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	<resource>ID, location, err := <resource>ImportParts(d)
	if err != nil {
		return nil, err
	}

	// Terraform sends an empty ID for an identity-based import, so the ID has to be
	// set here for the error diagnostics and the log line below to name the resource.
	d.SetId(<resource>ID)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return nil, err
	}

	<resource>, apiResponse, err := client.<Api>.<FindById>(ctx, <resource>ID).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find <resource> %q", <resource>ID), nil)
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while retrieving the <resource>: %w", err), nil)
	}

	tflog.Info(ctx, "<resource> imported", map[string]any{"resource_id": d.Id()})

	if err := set<Resource>Data(d, &<resource>); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	if err := set<Resource>Identity(d); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}

// <resource>ImportParts resolves the <resource> to import, either from the resource
// identity - which is how an import block with an `identity` argument, and the
// import config that `terraform query` generates, address a <resource> - or from the
// "<location>:<<resource>-id>" import string.
func <resource>ImportParts(d *schema.ResourceData) (<resource>ID, location string, err error) {
	if identity, identityErr := d.Identity(); identityErr == nil {
		if id, ok := identity.GetOk("id"); ok {
			loc, _ := identity.Get("location").(string)
			<resource>ID, _ = id.(string)
			return <resource>ID, loc, nil
		}
	}

	// unchanged legacy string path: keep the delimiter, the part count and the error
	// message the resource already used.
	importID := d.Id()
	location, parts := splitImportID(importID, ":")
	if len(parts) != 1 {
		return "", "", fmt.Errorf("invalid import identifier: expected one of <location>:<<resource>-id> or <<resource>-id>, got: %s", importID)
	}

	if err := validateImportIDParts(parts); err != nil {
		return "", "", fmt.Errorf("failed validating import identifier %q: %w", importID, err)
	}

	return parts[0], location, nil
}
```

Three things worth knowing:

- `identity.GetOk("id")` returns `exists=false` for the schema type's **zero value**, so an
  empty identity `id` does not take the identity branch — it falls through to the legacy
  path, where `splitImportID` yields `[""]` and `validateImportIDParts` rejects it. No API
  call with an empty ID is reachable. **Copilot will flag this as a missing validation. It is
  a false positive** — this exact exchange happened on #1034.
  **But that invariant is `splitImportID`'s, not the resolver's.** A resource with a plain
  (non-composite) import ID has no reason to call `splitImportID` at all, and then nothing
  rejects `""` — the fall-through returns `d.Id()` verbatim and the API is called with an
  empty ID. If your resolver does not go through `splitImportID`, add the guard yourself:
  `if d.Id() == "" { return "", fmt.Errorf("invalid import identifier: expected a <resource> UUID, got an empty string") }`.
  See `targetGroupImportID` in `ionoscloud/resource_target_group.go`.
- `d.Identity()` returns an *error*, not nil, when the resource declares no identity schema —
  which is why the resolver treats a non-nil `identityErr` as "fall back to the string import"
  rather than a failure.
- The legacy string path must be preserved unchanged. This is a pure prepend.
- Child resources: `d.SetId(parts[n])` **plus** `d.Set("datacenter_id", …)` etc. before
  `set<Resource>Identity` runs, since the setter reads the parent IDs back out of state.

`splitImportID` and `validateImportIDParts` live in `ionoscloud/utils.go` and are used by most
SDKv2 resources (`grep -l 'splitImportID(' ionoscloud/*.go`). Note the unrelated `splitImportID` in
`internal/framework/services/objectstorage/resource_object.go` — different package, different
contract.

---

## §2 — The list resource

File: `internal/framework/services/<service>/resource_<resource>_list.go`, package `<service>`.

Which package: Cloud API resources (`github.com/ionos-cloud/sdk-go/v6`) go in `compute`;
sdk-go-bundle products go in a package named after the product (`dns`, `vpn`, `cert`, …),
matching the existing `services/<product>/` call layer. `compute` holds only data sources and
list resources — an SDKv2 list resource needs no framework managed resource in its package.

### 2a. Header, assertions, struct, constructor

```go
// A list resource for <ionoscloud_type>, whose managed resource is still
// implemented with terraform-plugin-sdk/v2 and therefore lives on the other half of
// the mux. The protocol schemas the framework needs come from that resource via
// identity.SetRawV6Schemas, and the Identity it requires is declared alongside it in
// ionoscloud/resource_<resource>.go.

const <resource>ResourceType = "<ionoscloud_type>"

var (
	_ list.ListResource                 = (*<resource>ListResource)(nil)
	_ list.ListResourceWithConfigure    = (*<resource>ListResource)(nil)
	_ list.ListResourceWithRawV6Schemas = (*<resource>ListResource)(nil)
)

// <resource>ListResource lists <ionoscloud_type> instances.
type <resource>ListResource struct {
	bundle *bundleclient.SdkBundle

	// resourceSchema is the SDKv2 managed resource being listed, the source of the
	// protocol schemas returned by RawV6Schemas.
	resourceSchema *schema.Resource
}

// New<Resource>ListResource creates a new list resource for <ionoscloud_type>.
// <resource>Resource is the SDKv2 managed resource being listed; it is the source of
// the protocol schemas the framework needs.
func New<Resource>ListResource(<resource>Resource *schema.Resource) list.ListResource {
	return &<resource>ListResource{resourceSchema: <resource>Resource}
}
```

The type-name constant must exactly match the `ResourcesMap` key. There is also a
`utils/constant` constant for most types. The convention here is a **local const** for the
type name (all three list resources declare their own), even though the package does now
import `utils/constant` for other reasons — `resource_ipblock_list.go` and
`resource_target_group_list.go` both pull in `constant.IPBlockLimit` / `constant.TargetGroupLimit`
for their explicit `.Limit(...)`. So "the package does not import utils/constant" is no longer
a reason for the local const; consistency with the existing three is. Do not silently change it.

### 2b. The four boilerplate methods

```go
// RawV6Schemas hands the framework the protocol schemas of the SDKv2 managed
// resource. A framework-native list resource inherits them from the resource
// itself; this is only needed because <ionoscloud_type> lives on the SDKv2 side.
func (r *<resource>ListResource) RawV6Schemas(ctx context.Context, _ list.RawV6SchemaRequest, resp *list.RawV6SchemaResponse) {
	identity.SetRawV6Schemas(ctx, resp, <resource>ResourceType, r.resourceSchema)
}

// Metadata returns the type name of the managed resource being listed. It must match
// the SDKv2 resource exactly, otherwise terraform has no resource to attach the
// results to.
func (r *<resource>ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = <resource>ResourceType
}

// Configure stores the client bundle shared by the provider.
func (r *<resource>ListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected List Resource Configure Type",
			fmt.Sprintf("Expected *bundleclient.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.bundle = clientBundle
}

// ListResourceConfigSchema returns the schema for the list resource config block.
func (r *<resource>ListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			identity.FiltersKey: identity.FilterAttribute("<field1>", "<field2>"),
		},
	}
}
```

- `Metadata` sets the **full** type name; it does not prepend `req.ProviderTypeName` the way
  framework-native resources do.
- `Metadata` and `Configure` take `resource.MetadataRequest` / `resource.ConfigureRequest`,
  **not** `list.MetadataRequest` / `list.ConfigureRequest`. Both decoys exist in the vendored
  package (`list/metadata.go`, `list/configure.go`), and picking one is a *compile* error —
  because of the three `var _` assertions above, which is exactly why they must stay. Delete
  the `ListResourceWithConfigure` assertion and it becomes silent instead: `Configure` is not
  part of `list.ListResource`, so a `list.ConfigureRequest` signature compiles, the framework's
  `listResource.(list.ListResourceWithConfigure)` type assertion just fails
  (`vendor/.../internal/fwserver/server_listresource.go:105`), `Configure` is never called,
  `r.bundle` stays nil and `List` nil-derefs at runtime.
- The `req.ProviderData == nil` early return is mandatory — `Configure` is called once before
  provider configuration completes.
- `FilterAttribute(allowed...)` attaches `stringvalidator.OneOf`, so an unknown `field_name`
  is a plan-time error. Calling it with *no* arguments leaves it unvalidated and a typo
  silently matches nothing.

### 2c. `List` and the mapper

```go
func (r *<resource>ListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	identity.StreamList(ctx, stream, req,
		func(ctx context.Context) ([]<sdk>.<Resource>, error) {
			client, err := r.bundle.NewCloudAPIClientWithFailover(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to create the Cloud API client: %w", err)
			}

			// Depth(1) is what makes the API return the properties of every <resource>
			// instead of just its links. Filtering stays client-side, in the mapper.
			// Add .Limit(...) if this endpoint's SDK default is lower than the sibling
			// data source's — /ipblocks, /targetgroups and user management default to 100,
			// not 1000 (SKILL.md step 0.5, decision 1).
			items, apiResponse, err := client.<Api>.<ListCall>(ctx).Depth(1).Execute()
			if apiResponse != nil {
				tflog.Debug(ctx, "listed <resource>s", map[string]any{"status_code": apiResponse.SafeStatusCode()})
			}
			if err != nil {
				return nil, fmt.Errorf("failed to list <resource>s: %w", err)
			}
			if items.Items == nil {
				return nil, nil
			}

			return *items.Items, nil
		},
		r.map<Resource>,
	)
}

// map<Resource> maps a <resource> to an identity.MappedItem, or returns nil to skip it.
func (r *<resource>ListResource) map<Resource>(_ context.Context, includeResource bool, filters []identity.Filter, item <sdk>.<Resource>) (*identity.MappedItem, diag.Diagnostics) {
	if item.Id == nil || item.Properties == nil {
		return nil, nil
	}

	name := shared.ToValueDefault(item.Properties.Name)
	// … one local per filterable / identity field

	if !identity.MatchesFilters(map[string]string{
		"<field1>": name,
		// … exactly the keys passed to FilterAttribute above
	}, filters) {
		return nil, nil
	}

	mapped := &identity.MappedItem{
		DisplayName: name,
		Identity: &<resource>IdentityModel{
			// one field per attribute of the SDKv2 schema.ResourceIdentity
		},
	}

	if !includeResource {
		return mapped, nil
	}

	mapped.Resource = &<resource>ResourceModel{
		// one field per attribute of the SDKv2 resource schema, direct pointer copies
	}

	return mapped, nil
}
```

**`DisplayName` when `name` is optional.** `DisplayName` is the label `terraform query` prints
for each row. Datacenter's `name` is `Required: true`, so `shared.ToValueDefault(item.Properties.Name)`
is always populated — but that is not general: `ionoscloud/resource_ipblock.go` has
`"name": {Type: schema.TypeString, Optional: true}`, and the SDK property is
`Name *string \`json:"name,omitempty"\``. Nothing enforces non-emptiness (the framework only
reads `DisplayName` to detect a diagnostics-only event,
`fwserver/server_listresource.go:189`), so an unnamed item silently renders as a blank row.
**Check the SDKv2 schema**: if `name` is not `Required: true`, fall back —

```go
	displayName := shared.ToValueDefault(item.Properties.Name)
	if displayName == "" {
		displayName = *item.Id
	}
```

and see the matching stub warning in `test-harness.md` — the default two-result stub template
asserts a display name that an optional-name resource cannot produce.

`shared.ToValueDefault[T any](ptr *T) T` (`github.com/ionos-cloud/sdk-go-bundle/shared`,
vendored at `vendor/.../shared/utils.go:21-27`) is the repo-wide nil-safe deref — 39 existing
call sites outside `vendor/`. Use it from any package and add the import; **never declare a
local copy**. A hand-rolled `valueOrZero` helper used to live in package `compute` and was
deleted in favour of this — do not reintroduce it, in `compute` or anywhere else.

**`StreamList` semantics you must respect:**

| mapper returns | effect |
|---|---|
| `(nil, nil)` | item skipped silently, loop continues |
| `(nil, diags-with-error)` | item skipped; only `diags[0].Detail()` is logged at WARN, nothing reaches the user |
| `(non-nil, nil)` with `Identity == nil` | error result pushed, **stream aborts** |
| `(non-nil, diags-with-error)` | **stream aborts** — the diags are appended before the error check |
| `(non-nil, nil)`, identity set | result pushed |

So: `return nil, diags` for a recoverable per-item problem, `return mapped, nil` to include.
`objectstorage/resource_bucket_list.go` returns `(mapped, diags)` on a per-item error, which
truncates the whole listing on one bad item. **Do not copy that.**

Also: `fetch` runs eagerly, exactly once, before streaming — the "stream" is over an
already-materialised slice, and any pagination must live inside `fetch`. `req.Limit` is
ignored by `StreamList`. `fetch` errors lose all structure except `err.Error()` under a fixed
`"Failed to list resources"` summary, so wrap with context.

**Regional products** fan out inside `fetch`, copying
`internal/framework/services/pgsqlv2/resource_pg_cluster_list.go`. Note that the filters have
to be decoded a *second* time there, because `StreamList` decodes them for the mapper but does
not pass them to `fetch`:

```go
			var filters []identity.Filter
			req.Config.GetAttribute(ctx, path.Root(identity.FiltersKey), &filters)

			locations := <product>service.AvailableLocations()
			if loc := identity.FilterValue(filters, "location"); loc != "" {
				locations = []string{loc}
			}
```

**No `//nolint` on that line — copy the mariadbv2 form.** `Config.GetAttribute` returns
`diag.Diagnostics`, not an `error` (`vendor/.../tfsdk/config.go:34`; `diag.Diagnostics` has no
`Error()` method), so `errcheck` has nothing to report and a directive is dead weight. pgsqlv2
(`resource_pg_cluster_list.go:53`) and inmemorydbv2 (`:49`) do carry `//nolint:errcheck`, but
**bare, with no explanation** — they pass only because CI lints the diff
(`--new-from-rev`), so they are grandfathered and your new copy is not. A bare `//nolint`
fails `nolintlint`'s `require-explanation` (`.golangci.yml:82-84`) on your own line. Where you
*do* need a directive, the repo's well-formed shape is
`//nolint:prealloc // the number of streamed results is not known upfront`.

Pushing a filter down to the API does not exempt you from re-checking it in the mapper. Do not
add this second decode to a non-regional resource that does not need it — and note that
**no Cloud API resource is regional**; see SKILL.md step 0.5 decision 4. There is no
`AvailableLocations()` outside `services/dbaas/{pgsqlv2,mariadbv2,inmemorydbv2}`, so for a
Cloud API resource the `<product>service.AvailableLocations()` line above cannot be filled in
at all — that is the signal you are on the wrong branch, not something to work around.

### 2d. The resource model — the single biggest trap

**Resource model → plain Go pointers. Identity model → framework `types.*`.**

**A lone-`id` identity needs no struct at all.** `internal/framework/identity` already exports
`Model{ID types.String \`tfsdk:"id"\`}`. When the SDKv2 `Identity` declares exactly one `id`
attribute — which is the case for any resource addressed by bare UUID, with no location and no
parent IDs — use it directly and declare nothing per-resource:

```go
	Identity: &identity.Model{ID: types.StringValue(*item.Id)},
```

`resource_target_group_list.go` is the worked example. Declare the struct below **only** when
the identity has more than one attribute (datacenter and ipblock both add `location`):

```go
// <resource>IdentityModel mirrors the resource identity declared by the SDKv2
// <ionoscloud_type> resource.
type <resource>IdentityModel struct {
	ID       types.String `tfsdk:"id"`
	Location types.String `tfsdk:"location"`
}

// <resource>ResourceModel mirrors the SDKv2 schema of <ionoscloud_type>. It has
// to cover every attribute and block in that schema, because the framework fills the
// whole resource object from it.
//
// The fields are plain Go pointers rather than the types.String / types.Int64 values
// a framework-native resource would use. The schema behind this model is converted
// from the SDKv2 protocol schema, where an SDKv2 TypeInt arrives as a protocol
// Number and becomes a NumberAttribute - assigning a types.Int64 to it fails the
// type check. Plain Go types reflect into whichever framework type the converted
// schema ended up with, and the pointer carries the null/absent distinction.
type <resource>ResourceModel struct {
	ID *string `tfsdk:"id"`
	// …
	Timeouts *<resource>TimeoutsModel `tfsdk:"timeouts"`
}
```

Both models are `any` at the `MappedItem` boundary, so **nothing is checked at compile time** —
a wrong type is a runtime `Set` diagnostic that aborts the stream.

#### The model is coupled to a schema you do not own

The object type is built from the **live** SDKv2 schema, through `identity.SetRawV6Schemas`. So the
model is not merely *derived from* `ionoscloud/resource_<resource>.go` once — it stays welded to it.
Anyone who later adds an attribute to that resource, without ever hearing about this list resource,
breaks it. The break is invisible to `go build`, because the pairing is `any` on both sides:

```
Mismatch between struct and object type: Object defines fields not found in struct: enabled_features.
Struct: compute.datacenterCPUArchitectureModel
Object type: types.ObjectType["cpu_family":..., "enabled_features":..., "max_cores":..., ...]
```

**This is not hypothetical — it shipped.** PR #1034 added the datacenter list resource while a
Confidential Computing PR was independently adding `cpu_architecture.enabled_features` to the
datacenter SDKv2 schema. Each PR was green against its own base; neither was rebased on the other.
The merge of the second one left `master` red, with `terraform query` on `ionoscloud_datacenter`
failing for every user at the first result. It was found only by running the unit test on `master`.

Three consequences worth internalising:

1. **The unit test is the guard.** `Test<Resource>ListResource` fills the whole resource object from the live
   schema, so it fails the moment the schemas diverge. That is precisely why the test drives the real
   `ListResource` RPC instead of unit-testing the mapper. Never weaken it to a direct mapper call.
2. **Rebase before merging, not just before opening.** A green CI on your branch says nothing about
   the schema as it will exist at merge time. If any PR in flight touches
   `ionoscloud/resource_<resource>.go`, rebase and re-run the test before merge.
3. **The fix is always the same three edits**, and all three are required: add the field to the model
   with the schema key as its `tfsdk` tag, populate it in the mapper from the sdk-go property, and
   give the test fixture a non-nil value for it so the mapping is actually exercised. Adding only the
   model field silences the error while the attribute streams back permanently null.

When you touch a resource that has a list resource, grep for it first:

```bash
grep -rln "<resource>ResourceModel" internal/framework/services/
```

#### Deriving the model mechanically

1. Open `ionoscloud/resource_<resource>.go`, find `Schema: map[string]*schema.Schema{…}`.
2. One field per key, `tfsdk:"<the map key verbatim>"`. The Go field name is free.
3. **Add `ID *string \`tfsdk:"id"\``** — SDKv2 injects the implicit `id` into every resource's
   protocol schema and it is not in the `Schema:` map.
4. If the resource sets `Timeouts:`, add a `timeouts` block model with one `*string` per
   **non-nil field** of that `schema.ResourceTimeout`. The shared `resourceDefaultTimeouts`
   (`ionoscloud/provider.go`) sets `Create/Update/Delete/Default` — so `create`, `default`,
   `delete`, `update` and **no `read`**. If the resource has no `Timeouts:` field, omit the
   block entirely: an extra `timeouts` field is as wrong as a missing one.
5. Type by SDKv2 type:

| SDKv2 declaration | protocol type | model Go type |
|---|---|---|
| `TypeString` | String | `*string` |
| `TypeBool` | Bool | `*bool` |
| `TypeInt` | **Number** | `*int32` / `*int64` — **match the sdk-go field** |
| `TypeSet`/`TypeList` + `Elem: &schema.Schema{Type: TypeString}` | Set/List of String | `*[]string` |
| `TypeSet`/`TypeList` + `Elem: &schema.Resource{…}` | List/Set of Object | `*[]<resource><Nested>Model` + a nested model + a `map…` helper |

The last row's protocol *representation* splits, though the model field does not: a nested
schema that is `Computed` only becomes a protocol **attribute**
(`vendor/.../helper/schema/core_schema.go:102-106` — "Computed-only schemas are always handled
as attributes", which is datacenter's `cpu_architecture`), while one that is `Optional`
(± `Computed`) becomes a nested **block** with `NestingList`/`NestingSet` (`:111-112`,
`:201-205` — which is ipblock's `ip_consumers`). **The model field is `*[]<Nested>Model`
either way**, because blocks and attributes fold into the same object type
(`tfprotov6/schema.go:183-195`), and the test decodes both transparently since `resourceType`
comes from `ResourceSchemas[type].ValueType()`. It matters in two places:

1. Reading a schema dump by hand — the block sits under `Block.BlockTypes`, not
   `Block.Attributes`.
2. **The test's null assertion.** `toproto6.DynamicValue` runs `ReifyNullCollectionBlocks`
   (`vendor/.../internal/toproto6/dynamic_value.go:29-30`, reached from
   `toproto6.ListResourceResultWithResource` → `toproto6.Resource`), which turns a null list/set
   **block** into an **empty** collection. So when the API omits it, a Computed-only nested
   *attribute* decodes to `nil` (datacenter's `cpu_architecture`) but a nested *block* decodes
   to `[]any{}` (ipblock's `ip_consumers`). Assert `assert.Equal(t, []any{}, second["<block>"])`,
   not `assert.Nil`. A `MaxItems: 1` block is still `NestingList` (`core_schema.go:201-205`) and
   is reified too; only the injected `timeouts`, which is `NestingSingle`, stays null.

6. Recurse into nested `&schema.Resource{}` schemas. **Nested models get no `id` and no
   `timeouts`** — those are injected only at the top level.
7. Numeric widths must match the sdk-go struct field you assign from, because the mapper does
   direct pointer copies (`Name: item.Properties.Name` — no dereference, no conversion). If
   they differ you cannot pointer-copy and must allocate a converted value.

**Mirror the terraform schema, not the API model.** `DatacenterProperties` has
`GpuArchitecture` and `DefaultSecurityGroupId` that are not in the terraform schema and
correctly not in the model. Deriving from the SDK struct is a wrong turn.

**Not exercised by any reference implementation**, so treat as unverified and confirm against
the test: `TypeFloat` (→ `*float32`/`*float64` matching sdk-go), `TypeMap` (→ `*map[string]string`),
`Elem: &schema.Schema{Type: TypeInt}`, `Sensitive`. A `TypeMap` or `TypeFloat` attribute also
needs a change in the test decoder — see `test-harness.md`.

`MaxItems: 1` **is** exercised now: `resource_target_group_list.go` maps two of them
(`health_check`, `http_health_check`). The model field is `*[]…Model` like any other block —
`core_schema.go:201-205` sets `Nesting = NestingList` before it ever looks at `MaxItems` — and
the value decodes as a **one-element list**, not a bare object. Assert
`[]any{map[string]any{...}}`, and `[]any{}` when the API omits it.

Nested-block mapper, for every `Elem: &schema.Resource{}`. **Two shapes — pick by what the
sdk-go side looks like, not by what the terraform schema looks like.**

*(a) sdk-go field is a slice* — the common case, template below.

*(b) sdk-go field is a single struct pointer* — normal for a `MaxItems: 1` block, e.g.
`TargetGroupProperties.HealthCheck *TargetGroupHealthCheck`. There is no `*[]` to range over
and no `len` to preallocate from, but **the model field is still `*[]…Model`** — wrap the one
object. Making it a bare `*<Nested>Model` to match the SDK compiles fine (`MappedItem.Resource`
is `any`) and fails only at runtime, as a `Resource.Set` type-mismatch diagnostic that aborts
the whole result stream:

```go
func map<Resource><Nested>(item *<sdk>.<NestedProperties>) *[]<resource><Nested>Model {
	if item == nil {
		return nil
	}

	return &[]<resource><Nested>Model{{ /* direct pointer copies */ }}
}
```

*(a), the slice case:*

```go
func map<Resource><Nested>(items *[]<sdk>.<NestedProperties>) *[]<resource><Nested>Model {
	if items == nil {
		return nil
	}

	mapped := make([]<resource><Nested>Model, 0, len(*items))
	for _, item := range *items {
		mapped = append(mapped, <resource><Nested>Model{
			// direct pointer copies
		})
	}

	return &mapped
}
```

`nil` in → `nil` out. That preserves the null/absent distinction *in the model*; whether it
survives to the wire depends on attribute-vs-block — see §2d, a null list/set **block** is
reified to an empty collection. Package-level func, no receiver.
Preallocate with `make(..., 0, len(...))` — `prealloc` is enabled.

---

## §3 — Registration

`internal/framework/services/<service>/resources.go`:

```go
// ListResources returns the list of list resources for the <service> package.
//
// <ionoscloud_type> is still an SDKv2 resource, so its list resource needs that
// resource's schema. sdkv2Provider is the only way to reach it: package ionoscloud
// cannot be imported here, because its own in-package tests import the framework
// provider, which would make the import a cycle.
//
// When the schema cannot be reached the list resource is left unregistered rather
// than registered without one - the framework fails GetProviderSchema, terraform's
// first RPC, for a list resource with no schemas, which would take down the whole
// provider instead of just `terraform query` on <ionoscloud_type>.
func ListResources(sdkv2Provider *schema.Provider) []func() list.ListResource {
	// provider.New documents nil as a supported argument (list resources for SDKv2
	// resources then simply cannot be served), so this guard is not defensive padding.
	if sdkv2Provider == nil {
		return nil
	}

	var listResources []func() list.ListResource

	if r, ok := sdkv2Provider.ResourcesMap[<resource>ResourceType]; ok {
		listResources = append(listResources, func() list.ListResource { return New<Resource>ListResource(r) })
	}

	return listResources
}
```

**Use the accumulating form above even for the first entry.** `compute` was converted to it
when `ionoscloud_ipblock` was added — the original `if !ok { return nil }` shape was correct
for exactly one resource, but `sdkv2-branch.md` routes the whole Cloud API family into
`compute`, so the second entry turned that `return nil` into "a miss on resource B silently
unregisters resource A". Each lookup must be independent. Keep the `sdkv2Provider == nil`
guard in every version — only the `!ok` early return is the one that has to go.

**Shape by count — decide once, do not re-litigate per resource.** Up to three entries, add
another `if r, ok := …` block. At the fourth, switch to a table and keep the guard:

```go
	for _, e := range []struct {
		typeName string
		newList  func(*schema.Resource) list.ListResource
	}{
		{datacenterResourceType, NewDatacenterListResource},
		{ipblockResourceType, NewIPBlockListResource},
		// …
	} {
		if r, ok := sdkv2Provider.ResourcesMap[e.typeName]; ok {
			listResources = append(listResources, func() list.ListResource { return e.newList(r) })
		}
	}
```

Safe as written — `go.mod` declares `go 1.26.3`, so `e` and `r` are per-iteration.

Returning a partially-populated slice is fine: `provider.ListResources` only appends it
(`internal/framework/provider/provider.go:334-349`) and the framework re-keys everything into
a map (`fwserver/server_listresources.go:113`), so slice order and length carry no meaning.

**Never register a list resource without schemas.** The framework raises "ListResource Type
Defined without a Matching Managed Resource Type" and `GetProviderSchema` — terraform's very
first RPC — returns on that error, killing the **whole provider** for every operation, not
just `terraform query`. Degrading to "unregistered" is deliberate. `SetRawV6Schemas`'s own
nil-guard only logs; the `resources.go` guard is the one that protects the provider.

---

## §4 — Provider wiring

Already generic; `provider.New`, `main.go`, `xpprovider`, `internal/acctest` and
`internal/framework/identity/sdkv2.go` all need **no change**.

- A list resource **inside `compute`**: no `provider.go` change at all.
- A list resource in a **new** service package: give that package's `ListResources` the
  `*schema.Provider` parameter and add one line to
  `internal/framework/provider/provider.go`'s `ListResources`:

```go
		<service>.ListResources(p.sdkv2Provider),
```

`resp.ListResourceData = client` is already set in `Configure` alongside the other three.
