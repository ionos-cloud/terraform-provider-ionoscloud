# SDKv2-backed branch

For a managed resource still implemented with `terraform-plugin-sdk/v2` (anything in
`ionoscloud/resource_*.go`). Reference implementation, read it before writing anything:

- `ionoscloud/resource_datacenter.go` — the identity side
- `ionoscloud/resource_datacenter_list.go` — the list resource, 170 lines end to end
- `ionoscloud/list_resources.go` — registration
- `internal/framework/identity/sdkv2.go` — the two shared helpers everything above leans on

`ionoscloud/resource_ipblock_list.go` is a second worked example: an optional `name` (so a
`DisplayName` fallback, §2c), an explicit `Limit`, and a state writer that is *exported*
(`IpBlockSetData`, §2d).

**Read the current set rather than a list written here — it grows every time this skill is used:**

```bash
grep -l 'ResourceIdentity{' ionoscloud/*.go   # which resources declare an identity
cat ionoscloud/list_resources.go              # which list resources are registered
```

**Provenance matters when you read them.** `ionoscloud_datacenter` (PR #1034) is the shipped,
independently reviewed reference. Everything else in that output was produced by an earlier run of
*this skill* — a worked instance of these rules, still worth reading, but not independent
corroboration of them. Where this file argues a rule from one of those, it is quoting itself.
Guidance below about a child resource or `MaxItems: 1` nested blocks may still be extrapolation
from the schemas; check the grep output before assuming an example exists.

## Why this is not just "a framework list resource"

List resources cannot be implemented in SDKv2, so the list resource is framework code
even when the resource it lists is not. Two things follow, and both look like mistakes
until you know why.

- **The list resource lives in package `ionoscloud`, beside the SDKv2 resource — not
  under `internal/framework/services/`.** It has to. A result is produced by calling the
  managed resource's own `resource<Resource>()`, its state writer and
  `set<Resource>Identity()` — the constructor and the identity setter are always unexported, and
  the writer usually is, though `IpBlockSetData` is exported — and the import cannot be turned
  around:
  `ionoscloud/provider_test.go` is an *in-package* test (`package ionoscloud`) that
  imports `internal/framework/provider`, so any `internal/framework/...` → `ionoscloud`
  import is a cycle in the test build. Exporting the setters does not help. Other providers
  colocate for the same reason: terraform-provider-aws puts
  `internal/service/batch/job_definition_list.go` in `package batch` beside
  `job_definition.go` and calls `resourceJobDefinitionFlatten`, the function
  `resourceJobDefinitionRead` calls; terraform-provider-scaleway puts
  `internal/services/vpc/vpc_list.go` in `package vpc` and calls `ResourceVPC()` and
  `setVPCState`. If someone proposes moving it "where the framework code lives", this is
  the answer.
- **The framework has to be handed the managed resource's protocol schemas by hand.**
  That is the `list.ListResourceWithRawV6Schemas` interface, and `identity.SetRawV6Schemas`
  fills it from the SDKv2 `*schema.Resource` via `ProtoSchema` / `ProtoIdentitySchema`
  plus `terraform-plugin-mux/tf5to6server/translate`
  (`internal/framework/identity/sdkv2.go:24-42`).

Two dead ends, both already walked:

- **Do not obtain the schemas via `tf5to6server.UpgradeServer` + `GetProviderSchema`.** It
  works, but it stands up a second gRPC server and builds protocol schemas for every
  resource and data source to keep one. The original #1034 did this and it was replaced.
  Do not copy terraform-provider-aws's `ListResourceWithSDKv2Resource` *struct* either —
  that base exists to carry registry-injected state and to mutate the schema; with none of
  that, a plain function is enough.
- **Do not write a Go struct mirroring the SDKv2 schema.** That was the original shape and
  it shipped broken; §2d has the story and the mechanism that replaced it.

---

## §1 — Resource identity on the SDKv2 resource

**Hard prerequisite.** A resource with no `Identity` cannot be listed:
`resourceSchema.ProtoIdentitySchema(ctx)` returns a **nil func** (not a func returning nil,
`core_schema.go:419-423`), `SetRawV6Schemas` gives up, and `RawV6Schemas` has no
diagnostics channel — the only trace is a `tflog.Error` line. Debugging that from the
terraform side is miserable, and the failure is not local: see the end of §3.

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

**Only declare what addresses the resource.** Both shipped identities are `id` + `location`,
but that is not a template. `ionoscloud_target_group`, for instance, has no `location`
attribute at all — the collection is global — so its identity would be a lone `id`. It
declares no identity yet, so there is no file to copy that from; treat it as the reasoning,
not as a worked example. Do not add a `location` for symmetry with datacenter and ipblock;
every identity attribute has to be set on every read (1b), and one that is always empty is one
more thing to keep stable.

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

	// ONLY if `location` is in the identity you declared in 1a. A resource with no
	// location declares a lone `id` and must not write this - setting an attribute the
	// identity schema does not declare errors through MapFieldWriter, and panics rather
	// than returning under TF_ACC (helper/schema/schema.go:657-659 makes panicOnError
	// `os.Getenv("TF_ACC") != ""`). `ionoscloud_dns_zone` is the locationless example.
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

**This function is also the list resource's identity mapper**, verbatim — the list resource
calls it on a `ResourceData` it built itself (§2c). It reads out of state (`d.Id()`,
`d.Get(...)`) and never touches an API response, which is what makes that reuse work.

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
  was on `resourceIPBlockUpdate` and has since been corrected — read
  `ionoscloud/resource_ipblock.go:216-228` for the right wording; datacenter's Update delegates
  to Read and carries no such comment.) The SDK carries the prior identity into the apply on its
  own — `grpc_provider.go:1425` → `:1461` → `resource_data.go:856-857` → `:521-544` — so an
  update that never touches the identity still returns one and "Missing Resource Identity After
  Update" does *not* fire. The call is a safety net for state written before the resource
  declared an identity (e.g. `-refresh=false` on an old state file), not the everyday path. Say
  that in the comment; do not copy the Create rationale onto Update.
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

**This template is the Cloud API variant**, like §2c's: it builds
`meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)` and calls a `FindById` on an
`sdk-go/v6` API service. On the bundle branch take the client the resource's own importer already
takes, drop the `location` from both the resolver signature and the client construction when the
resource has none, and expect the fetch to be a method on the service client
(`client.GetZoneById(ctx, id)`) rather than `<Api>.<FindById>(ctx, id).Execute()`.

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
  That snippet returns **two** values because a plain-import-ID resource is normally locationless,
  so its resolver is `func <resource>ImportParts(d) (string, error)`. The template above returns
  three (`<resource>ID, location string, err error`) — match whichever arity your resolver has.
  No resource in the repo carries that guard yet, because no plain-import-ID resource has been
  given an identity. `resourceTargetGroupImport` in `ionoscloud/resource_target_group.go` is the
  unguarded shape to recognise: it takes `groupIp := d.Id()` and calls
  `TargetgroupsFindByTargetGroupId` with it, with nothing in between. Adding an identity to a
  resource shaped like that is exactly what makes the guard load-bearing.
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

File: `ionoscloud/resource_<resource>_list.go`, package `ionoscloud`.

**There is no package decision to make.** The file goes next to the resource it lists, always,
for the reason at the top of this document. `internal/framework/services/compute` is data
sources only — nothing on this branch belongs there. A **framework-native** list resource still
goes in its product's `internal/framework/services/<product>/` package; that branch is
unaffected (`framework-native-branch.md`).

### 2a. Header, assertions, struct, constructor

```go
// A list resource for <ionoscloud_type>, whose managed resource is still
// implemented with terraform-plugin-sdk/v2 and therefore lives on the other half of
// the mux. The protocol schemas the framework needs come from that resource via
// identity.SetRawV6Schemas.
import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	fwidentity "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

//
// It lives in this package, next to the resource it lists, so it can call
// resource<Resource>, set<Resource>Data and set<Resource>Identity directly. That is
// the whole design: results are produced by the resource's own state writer, so this
// file declares no model of the <resource> schema and there is nothing here to keep
// in sync when that schema changes.

var (
	_ list.ListResource                 = (*<resource>ListResource)(nil)
	_ list.ListResourceWithConfigure    = (*<resource>ListResource)(nil)
	_ list.ListResourceWithRawV6Schemas = (*<resource>ListResource)(nil)
)

// <resource>ListResource lists <ionoscloud_type> instances.
type <resource>ListResource struct {
	bundle *bundleclient.SdkBundle

	// resourceSchema is the SDKv2 managed resource being listed: the source of the
	// protocol schemas returned by RawV6Schemas, and of the ResourceData every
	// result is written into.
	resourceSchema *schema.Resource
}

// New<Resource>ListResource creates a new list resource for <ionoscloud_type>.
func New<Resource>ListResource() list.ListResource {
	return &<resource>ListResource{resourceSchema: resource<Resource>()}
}
```

- **The constructor takes no arguments.** It calls `resource<Resource>()` itself, so it *is* a
  `func() list.ListResource` and registration needs no wrapper closure (§3). Do not add a
  `*schema.Resource` parameter back; there is nothing to inject and nothing to fail.
- `resource<Resource>()` is called **once**, at construction, and the resulting `*schema.Resource`
  is kept for the life of the list resource. It is used for two things and never mutated:
  `RawV6Schemas` reads the protocol schemas off it, and every mapped item gets a fresh
  `ResourceData` from it (§2c). Do not call it per item — one `*schema.Resource` serves any
  number of `Data()` calls.
- **Use `constant.<Resource>Resource`, do not declare a local type-name const.** The type name
  must exactly match the `ResourcesMap` key in `ionoscloud/provider.go`, or terraform has no
  managed resource to attach the results to — and `provider.go:97,98,136` keys that map with
  `constant.DatacenterResource`, `constant.IpBlockResource`, `constant.TargetGroupResource`.
  Referencing the same constant makes the match structural instead of a convention two string
  literals have to keep. This file is in `package ionoscloud`, so the constant is one import away
  (`utils/constant`); the list resources in `internal/framework/services/` predate that and use
  literals because they were written in a package that did not import it.
- **The reference implementation does not follow that rule — this is a known divergence, not a
  precedent.** `ionoscloud/resource_datacenter_list.go` declares
  `const datacenterResourceType = "ionoscloud_datacenter"` and passes it to both
  `SetRawV6Schemas` and `Metadata`, even though `constant.DatacenterResource` exists and is what
  keys `ResourcesMap`. It predates the rule and has not been migrated.
  `ionoscloud/resource_ipblock_list.go` is the one that gets this right
  (`constant.IpBlockResource` in both places). Copy ipblock here, not datacenter, and do not
  read datacenter's local const as sanctioning a new one.

### 2b. The four boilerplate methods

```go
// RawV6Schemas hands the framework the protocol schemas of the SDKv2 managed
// resource. A framework-native list resource inherits them from the resource
// itself; this is only needed because <ionoscloud_type> lives on the SDKv2 side.
func (r *<resource>ListResource) RawV6Schemas(ctx context.Context, _ list.RawV6SchemaRequest, resp *list.RawV6SchemaResponse) {
	fwidentity.SetRawV6Schemas(ctx, resp, constant.<Resource>Resource, r.resourceSchema)
}

// Metadata returns the type name of the managed resource being listed. It must match
// the SDKv2 resource exactly, otherwise terraform has no resource to attach the
// results to.
func (r *<resource>ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = constant.<Resource>Resource
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
			fwidentity.FiltersKey: fwidentity.FilterAttribute("<field1>", "<field2>"),
		},
	}
}
```

- **`internal/framework/identity` is imported as `fwidentity` here.** In package `ionoscloud`,
  `identity` is already the name of the SDKv2 `*schema.IdentityData` that `d.Identity()` returns
  and that `set<Resource>Identity` operates on (`ionoscloud/resource_datacenter.go:341`), so an
  unaliased import would make `identity.` mean two unrelated things a few lines apart. Go imports
  are file-scoped and the two never collide in one file today, so this is **readability, not a
  compiler requirement** — but both SDKv2 list resources do it, and the SDKv2 identity is the
  reading a maintainer of this package expects. Keep the alias.
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

**This template is the Cloud API (`sdk-go/v6`) variant.** For an sdk-go-bundle product take the
deltas listed after it — they are not optional adjustments, four of these lines do not compile.

```go
func (r *<resource>ListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	fwidentity.StreamList(ctx, stream, req,
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
//
// The mapping itself is set<Resource>Data, the same state writer resource<Resource>Read
// uses, run against a ResourceData built from the live schema. Nothing here knows what
// attributes a <resource> has.
func (r *<resource>ListResource) map<Resource>(_ context.Context, includeResource bool, filters []fwidentity.Filter, item <sdk>.<Resource>) (*fwidentity.MappedItem, diag.Diagnostics) {
	var diags diag.Diagnostics

	if item.Id == nil || item.Properties == nil {
		return nil, nil
	}

	// Build the map from the accessors directly. Do NOT lift each filterable field into a
	// single-use local first - the map keys already name the values, so the locals only
	// add a line each and a second place to get the pairing wrong.
	if !fwidentity.MatchesFilters(map[string]string{
		"<field1>": shared.ToValueDefault(item.Properties.<Field1>),
		"<field2>": item.Properties.<Field2>,  // a value field needs no ToValueDefault
		// … exactly the keys passed to FilterAttribute above
	}, filters) {
		return nil, nil
	}

	data := r.resourceSchema.Data(&terraform.InstanceState{})
	if err := set<Resource>Data(data, &item); err != nil {
		diags.AddError("Failed to map the <resource>", err.Error())
		return nil, diags
	}

	// CHILD RESOURCES ONLY: the parent ids. The state writer usually does not set them -
	// they are not in the API object, they come from the fetch context - and the identity
	// setter is about to read them back out. This is the one exception to "no attribute
	// mapping of your own": set only parent ids, and only the ones the writer provably
	// leaves unset. `ionoscloud_dns_record` is the live case, where SetRecordData never
	// touches zone_id.
	//   if err := data.Set("<parent>_id", <parentID>); err != nil { … }

	// set<Resource>Identity reads id and location back out of the ResourceData, so it
	// has to run after set<Resource>Data.
	if err := set<Resource>Identity(data); err != nil {
		diags.AddError("Failed to map the <resource> identity", err.Error())
		return nil, diags
	}

	mapped, err := fwidentity.MappedItemFromResourceData(name, data, includeResource)
	if err != nil {
		diags.AddError("Failed to convert the <resource> state", err.Error())
		return nil, diags
	}

	return mapped, diags
}
```

Those last three statements are the entire mapping. §2d explains why that is enough.

#### The sdk-go-bundle deltas

Step 0's grep (6) told you which branch you are on. If the resource's CRUD takes a `<Product>Client`
off the bundle rather than `NewCloudAPIClient*`, five things in the template above change.
`ionoscloud_dns_zone` is the worked example — it was written from this list, and this list from it.

1. **The client is the resource's own, and it is already on the bundle.** Not
   `NewCloudAPIClientWithFailover`, which returns a `*ionoscloud.APIClient` that cannot reach a
   bundle product at all. Prefer the service package's own list helper where there is one, so the
   list resource reads exactly what the sibling data source reads:

   ```go
   zones, apiResponse, err := r.bundle.DNSClient.ListZones(ctx, "")
   ```

2. **Check for `Depth` rather than assuming it away.** Most bundle collections have none — it is
   primarily a Cloud API parameter, and dns, nfs, vpn, cert and kafka all lack it — but
   `vmautoscaling`, which backs `ionoscloud_autoscaling_group`, **does** take one, and
   `services/autoscaling/groups.go` already passes it on the by-id path. So grep before you drop:

   ```bash
   grep -n 'func (r Api<X>GetRequest) Depth' vendor/github.com/ionos-cloud/sdk-go-bundle/products/<product>/v2/api_<x>.go
   ```

   No hit → drop it, the collection returns full properties already. A hit → pass `.Depth(1)` for
   the same reason the Cloud API branch does.

3. **`Items` is a value slice**, not `*[]T`, so the `items.Items == nil` guard and the
   `*items.Items` deref both go away: `return zones.Items, nil`.

4. **The element's `Id` and `Properties` are values, not pointers.** `if item.Id == nil ||
   item.Properties == nil` does not compile. The only unusable item is one with an empty id:

   ```go
   if zone.Id == "" {
       return nil, nil
   }
   ```

   Fields *inside* `Properties` are still a mix — `ZoneName` is a plain `string`, `Description` a
   `*string` — so `shared.ToValueDefault` still applies to the pointer ones and must not be used on
   the value ones.

5. **The writer may be a method on that client**, called through the bundle:
   `r.bundle.DNSClient.SetZoneData(data, zone)`. Note it takes its object **by value** — drop the
   `&` the Cloud API template passes (`set<Resource>Data(data, &item)`). See §2d.

6. **The import block changes too.** §2a's list is Cloud-API-only: `ionoscloud
   "github.com/ionos-cloud/sdk-go/v6"` is unused on this branch and will not compile. Import the
   product package instead — `dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"` — and
   keep `github.com/ionos-cloud/sdk-go-bundle/shared` only if you still call `ToValueDefault` on a
   pointer field.

**Push a filter down only when it removes API calls.** A bundle list helper may accept a
server-side filter — `ListZones(ctx, filterName)` sends `filter.zoneName`. `StreamList`'s fetch
closure is `func(context.Context) ([]T, error)` and is never handed the filters, but that is not a
wall: the closure can capture `req` and re-read the config itself, which is exactly what the three
regional list resources do —
`internal/framework/services/pgsqlv2/resource_pg_cluster_list.go` calls
`req.Config.GetAttribute(ctx, path.Root(identity.FiltersKey), &filters)` inside the fetch, then
`identity.FilterValue(filters, "name")` and a `location` filter that narrows the fan-out, with the
comment "Read filters early to skip unnecessary regional API calls". `mariadbv2` and
`inmemorydbv2` are identical.

That is the rule: **push down when it saves round trips** (a fan-out you can narrow, a collection
you would otherwise page through), and otherwise pass the no-filter form and let the mapper do it —
which is what the two Cloud API list resources and `ionoscloud_dns_zone` do, because a single
global call costs the same either way. **Either way the mapper still re-checks**: pgsqlv2 pushes
`name` down *and* calls `identity.MatchesFilters`, because a server-side filter's semantics are
not guaranteed to match the exact compare the docs promise.

**`DisplayName` when `name` is optional.** `DisplayName` is the label `terraform query` prints
for each row. Datacenter's `name` is `Required: true`, so `shared.ToValueDefault(item.Properties.Name)`
is always populated — but that is not general: `ionoscloud/resource_ipblock.go:51-54` has
`"name": {Type: schema.TypeString, Optional: true}`, and the SDK property is
`Name *string \`json:"name,omitempty"\``. Nothing enforces non-emptiness (the framework only
reads `DisplayName` to detect a diagnostics-only event,
`fwserver/server_listresource.go:189`), so an unnamed item silently renders as a blank row.
**A third case: some resources have no `Name` property at all** — `ionoscloud_cdn_distribution`
carries `Domain`, a user carries `Email`. Pick the closest human identifier the properties do
carry, fall back to the id when it is empty, and never leave `DisplayName` blank: the framework
reads a blank one as a diagnostics-only event. Which property to label rows with is a judgement
call, so put it in the step 0.5 report rather than deciding it silently in the mapper.

**Check the SDKv2 schema**: if `name` is not `Required: true`, fall back —

A local **is** justified here, because the value is used twice — once to test for empty, once to
pass on:

```go
	displayName := shared.ToValueDefault(item.Properties.Name)
	if displayName == "" {
		displayName = *item.Id
	}
```

and pass `displayName` to `MappedItemFromResourceData` instead of `name`. `resource_ipblock_list.go`
is the worked example; see also the matching stub warning in `test-harness.md` — the default
two-result stub template asserts a display name that an optional-name resource cannot produce.

`shared.ToValueDefault[T any](ptr *T) T` (`github.com/ionos-cloud/sdk-go-bundle/shared`,
vendored at `vendor/.../shared/utils.go:21-27`) is the repo-wide nil-safe deref. Use it from any
package and add the import; **never declare a local copy**. A hand-rolled `valueOrZero` helper
used to live in package `compute` and was deleted in favour of this — do not reintroduce it,
anywhere.

**`StreamList` semantics you must respect:**

| mapper returns | effect |
|---|---|
| `(nil, nil)` | item skipped silently, loop continues |
| `(nil, diags-with-error)` | item skipped; only `diags[0].Detail()` is logged at WARN, nothing reaches the user |
| `(non-nil, nil)` with `Identity == nil` | error result pushed, **stream aborts** |
| `(non-nil, diags-with-error)` | **stream aborts** — the diags are appended before the error check |
| `(non-nil, nil)`, identity set | result pushed |

So: `return nil, diags` for a recoverable per-item problem, `return mapped, nil` to include.
That is why every failure branch in the template above returns `nil, diags` and only the tail
returns `mapped`. `objectstorage/resource_bucket_list.go` returns `(mapped, diags)` on a per-item
error, which truncates the whole listing on one bad item. **Do not copy that.**

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

### 2d. No model — the result is the resource's own state

The three statements at the end of the mapper are the whole mapping:

```go
	data := r.resourceSchema.Data(&terraform.InstanceState{})
	set<Resource>Data(data, &item)     // the resource's OWN state writer, the one Read uses
	set<Resource>Identity(data)        // AFTER the writer — it reads values back out
	fwidentity.MappedItemFromResourceData(displayName, data, includeResource)
```

**Do not add a Go struct mirroring the SDKv2 schema.** `MappedItem.Identity` and
`MappedItem.Resource` are `any`, so one compiles; it is the shape this design replaced.

#### What the design requires of the state writer

The templates spell the writer `set<Resource>Data` because that is the majority name in
`package ionoscloud`, but **the name is not the contract, the signature is.** Grep before you
write the mapper:

**Follow Read; do not grep for a name.** A name grep misses on both sides: `func .*SetData`
requires the literal "SetData", `func set.*Data` also matches unrelated helpers, and for a bundle
product the writer is not in `ionoscloud/` at all. The read function always calls it, so read the
read function:

```bash
grep -n 'ReadContext:' ionoscloud/resource_<resource>.go   # names the read func
sed -n '/func <thatFunc>/,/^}/p' ionoscloud/resource_<resource>.go
```

The mapper can call any writer shaped `func(d *schema.ResourceData, obj <sdk>.<X>) error` — one
`ResourceData`, one object of the element type the collection endpoint returns, an `error`. Neither
the spelling, the receiver, nor the package matters:

- `setDatacenterData` — unexported, package-level, takes its object by pointer.
- `IpBlockSetData` — exported, package-level; the list resource calls it directly, no adapter.
- `SetZoneData` — a **method** on the product's service client (`services/dns/zone.go`), taking its
  object by value. The mapper already holds the bundle, so it is reachable as
  `r.bundle.DNSClient.SetZoneData(data, zone)` and living outside package `ionoscloud` costs
  nothing. It also calls `d.SetId` itself, which not every writer does — check, because the
  identity setter reads the id back out of the `ResourceData`.

**A writer that needs arguments the mapper cannot supply is telling you something about the
resource.** Several writers in this package take more than that pair: `setBackupUnitData` also
wants an `*ionoscloud.Contracts`, `setApplicationLoadBalancerData` an `*ionoscloud.FlowLog`, and
`setResourceServerData`, `setResourceVCPUServerData` and `setGroupData` want a `context.Context`
plus an `*ionoscloud.APIClient` so they can fetch more. The mapper is handed only
`(ctx, includeResource, filters, item)` — one already-fetched element of the collection
(`internal/framework/identity/list.go`, the `mapper(ctx, req.IncludeResource, filters, item)`
call). It has no `Contracts` and no `FlowLog`, and although it can reach a client, because it is
a method and `r.bundle` is right there, every use of one costs an extra request per item.

So: **a resource is listable under this design when its full state is derivable from one element
of the collection response at `Depth(1)`.** When it is not, pick deliberately rather than
discovering it when the mapper will not compile:

1. list it anyway and accept that the writer's extra inputs are absent — defensible only if the
   attributes they fill are genuinely optional, which for `server` and `group` they are not;
2. do the extra per-item fetch inside the mapper, and own the N+1 request cost a listing over
   hundreds of items then pays;
3. do not list this resource yet — usually the right answer, and a finding to report back rather
   than a problem to code around.

This is a constraint on the *resource*, and it is invisible until you read the writer's
signature. Read it in step 0, not after the mapper fails to compile.

#### Why no model is possible

Both halves of a list result are rendered from **one** `configschema.Block`:

- the *type* the framework builds the result against comes from `resp.ProtoV6Schema`, which
  `SetRawV6Schemas` fills with `translate.Schema(resourceSchema.ProtoSchema(ctx)())`, and
  `ProtoSchema` is `convert.ConfigSchemaToProto(ctx, r.CoreConfigSchema())`
  (`core_schema.go:407-416`);
- the *value* comes from `rd.TfTypeResourceState()`, which is
  `schemaMap(d.schema).CoreConfigSchema()` plus the injected `id` and `timeouts`
  (`resource_data.go:84-163`).

Two renderings of the same source, so they cannot disagree — where a hand-written struct is a
third, independent rendering that has to be kept equal to both by hand. Same for the identity:
`ProtoIdentitySchema` and `TfTypeIdentityState` both go through
`schemaMap(r.Identity.SchemaMap()).CoreConfigSchema()` (`core_schema.go:397-404`,
`resource_data.go:67-82`).

The motivation, in two sentences: PR #1034 shipped a `datacenterResourceModel` mirroring the
datacenter schema while an unrelated Confidential Computing PR was independently adding
`cpu_architecture.enabled_features` to that schema, and the merge left `master` failing every
`terraform query` on `ionoscloud_datacenter` with *"Object defines fields not found in struct:
enabled_features"*. The fix was to delete the duplicate, not to guard it.

#### What each of the three lines does

**`r.resourceSchema.Data(&terraform.InstanceState{})`** returns a `*schema.ResourceData` built
from the live schema over an empty state (`resource.go:1405-1426`). Two details matter:

- it goes through `schemaMapWithIdentity{r.SchemaMap(), r.Identity.SchemaMap()}`
  (`schema.go:650-674`), so the `ResourceData` carries the **identity** schema as well — that is
  what lets `set<Resource>Identity`'s `d.Identity()` succeed. `r.Apply` builds its `ResourceData`
  from the same pair (`resource.go:909-910`).
- it then copies `result.timeouts = r.Timeouts`. That single assignment is what makes
  `TfTypeResourceState` materialise a `timeouts` block at all — it only injects one when
  `d.timeouts != nil` (`resource_data.go:98-145`). It is also the reason `nullTimeouts` exists,
  below.

**Use `Data`, not `TestResourceData`.** The latter builds the `ResourceData` struct directly
and never assigns `timeouts` (`resource.go:1431-1436`), so `TfTypeResourceState` would omit the
`timeouts` block while the result type — built from `r.CoreConfigSchema()`, which injects it
whenever `r.Timeouts != nil` (`core_schema.go:325-330`) — still declares it, and `Set` would
reject the value on the type comparison.

Build a fresh one **per item**. It is state, and reusing it across items leaks values from the
previous one into any attribute the writer leaves untouched.

**Order: writer, then identity setter.** `set<Resource>Identity` reads back out of the
`ResourceData` — `d.Id()` and `d.Get("location")` — so it can only run once the writer has put
them there. This is the same ordering constraint as §1c, for the same reason.

**The writer must call `d.SetId`.** Both conversions go through `ResourceData.State()`, which
returns `nil` while the ID is empty (`resource_data.go:425-434`), and both then fail with
`"state is nil, call SetId() on ResourceData first"`. Every `set<Resource>Data` in this package
already begins with `d.SetId(*<resource>.Id)` because Read needs it too — check yours does before
assuming it.

**`fwidentity.MappedItemFromResourceData(displayName, data, includeResource)`**
(`internal/framework/identity/sdkv2.go:44-89`) does the rest:

- `rd.TfTypeIdentityState()` → `MappedItem.Identity`, and errors when it comes back nil, which is
  the "you forgot the identity setter" message;
- when `includeResource` is false it stops there — `terraform query` without `include_resource`
  never pays for the resource conversion;
- otherwise `rd.TfTypeResourceState()` → `nullTimeouts(...)` → `MappedItem.Resource`.

Both values are `tftypes.Value` passed **by value, not by pointer**. `Set` type-asserts a
`tftypes.Value` and takes a documented fast path — it compares
`d.Schema.Type().TerraformType(ctx)` with `v.Type()` and assigns
(`vendor/.../internal/fwschemadata/data_set.go:20-35`). A `*tftypes.Value` misses the assertion
and falls through to struct reflection, which is exactly the machinery being avoided.

`nullTimeouts` nulls the whole `timeouts` block out of the converted value. The flatmap shims
cannot tell a null single block from an empty one, so they always materialise it as an object of
null attributes; SDKv2's own `ReadResource` nulls it back out on every read of its own resources
for the same reason. A listed resource has no timeouts either — they are configuration, not
state — so doing the same keeps a query result identical to what a refresh writes. The unit test
asserts it (`ionoscloud/resource_datacenter_list_test.go`).

#### What the value looks like on the wire

Still worth knowing, because the test asserts against it and the shapes are not obvious.

`toproto6.DynamicValue` runs `ReifyNullCollectionBlocks`
(`vendor/.../internal/toproto6/dynamic_value.go:29-30`, reached from
`toproto6.ListResourceResultWithResource` → `toproto6.Resource`), which turns a null list/set
**block** into an **empty** collection. Whether a nested schema is a block or an attribute is
decided in `core_schema.go`: a nested schema that is `Computed` only becomes a protocol
**attribute** (`:102-106` — "Computed-only schemas are always handled as attributes", which is
datacenter's `cpu_architecture`), while one that is `Optional` (± `Computed`) becomes a nested
**block** (`:108-112`, `:201-205` — which is ipblock's `ip_consumers`). So when the API omits it:

- a Computed-only nested **attribute** decodes to `nil` → assert `assert.Nil`;
- a nested **block** decodes to `[]any{}` → assert `assert.Equal(t, []any{}, second["<block>"])`.

`MaxItems: 1` changes nothing here: `core_schema.go:201-205` sets `Nesting = NestingList` from
`s.Type` before it ever looks at `MaxItems`, so such a block is still a list, is reified when
absent, and decodes as a **one-element list** when present — assert `[]any{map[string]any{...}}`.
No shipped list resource has a `MaxItems: 1` block yet; `ionoscloud_target_group`'s
`health_check` and `http_health_check` are the shape to expect if you are the one to list it. The
injected `timeouts` is `NestingSingle`, so it is not reified and stays null, which is what
`nullTimeouts` leaves behind.

#### What is still coupled to the schema

Nothing about the resource's *attributes*. When someone adds an attribute to
`ionoscloud/resource_<resource>.go`, they update `set<Resource>Data` — which they have to do
anyway for Read — and the list resource follows for free. There is no model to grep for, and no
rebase-before-merge hazard for schema changes.

Three things do remain hand-written and do need review when the resource changes:

1. the filter keys — the strings in `FilterAttribute(...)` and in the `MatchesFilters` map must
   stay identical, and both must name real, meaningful fields;
2. the `Identity` block and `set<Resource>Identity` (§1), which the list resource calls directly;
3. the `DisplayName` fallback, if `name` stops being `Required`.

---

## §3 — Registration

`ionoscloud/list_resources.go`, one line per list resource:

```go
// ListResources returns the list resources for the managed resources in this package.
//
// List resources can only be implemented with terraform-plugin-framework, so a list
// resource for a resource that is still on SDKv2 is framework code living in this
// package. It has to be: the results are produced by the resource's own state writer,
// which is unexported, and no package under internal/framework can import this one -
// ionoscloud/provider_test.go is an in-package test that imports
// internal/framework/provider, so an import back would be a cycle in the test build.
func ListResources() []func() list.ListResource {
	return []func() list.ListResource{
		NewDatacenterListResource,
		NewIPBlockListResource,
	}
}
```

That is the entire registration. The constructors are named directly, not wrapped in
`func() list.ListResource { … }`, because they already have that signature (§2a).

**No lookup, no `ok` branch, no nil guard.** Adding a list resource is one line here and
nothing else. `New<Resource>ListResource` calls `resource<Resource>()` in-package: it cannot
miss, cannot return nil and has no error path, so there is nothing to guard and nothing to
degrade to. If you find yourself reaching for a `*schema.Provider` or a `ResourcesMap` lookup
to get at the schema, you are in the wrong package — the schema is one function call away (§2a),
and no provider is threaded to the framework side any more (§4).

**Never register a list resource without schemas.** The framework raises "ListResource Type
Defined without a Matching Managed Resource Type"
(`vendor/.../fwserver/server_listresources.go:99-109`) when a list resource has no matching
framework resource *and* `RawV6Schemas` supplied neither schema, and `GetProviderSchema` —
terraform's very first RPC — returns on that error
(`vendor/.../fwserver/server_getproviderschema.go:84-88`), killing the **whole provider** for
every operation, not just `terraform query`.

The one remaining way to walk into it is **removing or forgetting the `Identity` block on the
SDKv2 resource** (§1). Then `ProtoIdentitySchema` is a nil func, `SetRawV6Schemas` returns
having set nothing but a `tflog.Error`, and the provider dies at the first RPC. There is no
guard left that could quietly unregister the list resource instead, so treat the `Identity`
block as load-bearing for the entire provider, not just for this feature. The unit test is
what catches it: `TestDatacenterListResource` asserts
`providerSchema.ListResourceSchemas[<type>]` exists right after `GetProviderSchema`
(`ionoscloud/resource_datacenter_list_test.go:44-57`).

---

## §4 — Provider wiring

**For a new SDKv2-backed list resource there is no provider wiring to do.** One line in
`ionoscloud/list_resources.go` (§3) is the whole registration. `internal/framework/provider/provider.go`,
`main.go`, `xpprovider`, `internal/acctest` and `internal/framework/identity/sdkv2.go` need no
edit at all. This is the part of the design that got simpler; do not go looking for a second
place to register.

For context, the wiring that already exists:

```go
// New creates a new provider. sdkv2ListResources are the list resources declared
// alongside the SDKv2 managed resources they list, normally ionoscloud.ListResources();
// passing none is allowed, and then only the framework-native list resources are served.
func New(sdkv2ListResources ...func() list.ListResource) provider.Provider {
	return &IonosCloudProvider{sdkv2ListResources: sdkv2ListResources}
}
```

(`internal/framework/provider/provider.go:64-69`), and the slice is appended alongside the
framework-native packages in `ListResources` (`:336-353`):

```go
	listResources := [][]func() list.ListResource{
		objectstorage.ListResources(),
		objectstoragemanagement.ListResources(),
		pgsqlv2.ListResources(),
		inmemorydbv2.ListResources(),
		mariadbv2.ListResources(),
		p.sdkv2ListResources,
	}
```

`New` is variadic so that `provider.New(ionoscloud.ListResources()...)` reads as a call and an
empty slice needs no special case. Every construction site already passes it:

| call site | |
|---|---|
| `main.go:41` | the served provider |
| `xpprovider/provider.go:16` | crossplane / upjet |
| `internal/acctest/acctest.go:55` | acceptance tests |
| `internal/framework/provider/provider_test.go:60` | framework provider tests |
| `ionoscloud/provider_test.go:69`, `:91` | the v5 and v6 mux factories |
| `ionoscloud/resource_datacenter_list_test.go:164` | the list-resource unit test harness |

Touch these only if you add a **new** entry point that builds the provider. A framework-native
list resource in a **new** service package is the one case that still needs a `provider.go`
edit: add that package's `ListResources()` to the slice above.

`resp.ListResourceData = client` is already set in `Configure` alongside the other three.
