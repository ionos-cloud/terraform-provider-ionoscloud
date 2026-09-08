# Framework-native branch

For a managed resource already implemented with terraform-plugin-framework (anything under
`internal/framework/services/<service>/resource_*.go` with a `Metadata` method building its
name from `req.ProviderTypeName`).

This is substantially less work than the SDKv2 branch: **no `RawV6Schemas`, no
`*schema.Provider` threading, no `ResourcesMap` lookup, and no hand-written resource model.**
The list resource is just more methods on the *same* struct as the managed resource, which
already supplies `Metadata`, `Schema`, `IdentitySchema` and `Configure`.

Best template: `internal/framework/services/pgsqlv2/resource_pg_cluster_list.go`. It is the
only one of the five with doc comments on every non-obvious part, it shows the
two-attribute identity + regional fan-out + timeouts-nulling triad, and it reuses the
resource's existing mapper instead of writing a second one. Use
`objectstoragemanagement/resource_object_storage_accesskey_list.go` as the secondary template
for the single-`id`-identity, non-regional case.

---

## §1 — Resource identity on the managed resource

On the **managed** resource, never on the list resource. Add the interface assertion and the
`IdentitySchema` method:

```go
var (
	_ resource.ResourceWithImportState = (*<x>Resource)(nil)
	_ resource.ResourceWithConfigure   = (*<x>Resource)(nil)
	_ resource.ResourceWithIdentity    = (*<x>Resource)(nil)
)

type <x>IdentityModel struct {
	ID       types.String `tfsdk:"id"`
	Location types.String `tfsdk:"location"`
}

func (r *<x>Resource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id":       identityschema.StringAttribute{RequiredForImport: true},
			"location": identityschema.StringAttribute{RequiredForImport: true},
		},
	}
}
```

Import `"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"`.

For a lone `id`, reuse the shared `identity.Model` instead of declaring a package-local
struct — that is what the accesskey resource does.

**Write the identity on every state-producing path** — Create, Read, Update (Delete needs
none), and read it in ImportState:

```go
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &<x>IdentityModel{ID: plan.ID, Location: plan.Location})...)
```

```go
func (r *<x>Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.Identity != nil {
		// ... resp.Diagnostics.Append(req.Identity.Get(ctx, &id)...)
	}
	// ... existing string-import path unchanged
}
```

Where the identity model lives: next to the `IdentitySchema` that defines it, in the managed
resource file. (Only the SDKv2 branch puts it in the list file, because there is no framework
managed resource file to put it in.)

---

## §2 — The list resource

```go
var (
	_ list.ListResource              = (*<x>Resource)(nil)
	_ list.ListResourceWithConfigure = (*<x>Resource)(nil)
)

// New<X>ListResource creates a new list resource for <ionoscloud_type>.
func New<X>ListResource() list.ListResource {
	return &<x>Resource{}
}
```

No `Metadata`, no `Configure`, no `RawV6Schemas`, no constructor argument — all inherited
from the managed resource. `list.ListResource.Metadata` is declared with
`resource.MetadataRequest` and `list.ListResourceWithConfigure.Configure` with
`resource.ConfigureRequest` precisely so one method satisfies both interfaces.

`ListResourceConfigSchema`, `List` and the mapper are **identical in shape to the SDKv2
branch** — see `sdkv2-branch.md` §2b/§2c, including the `StreamList` fatal-vs-skip table and
the regional fan-out pattern. Two differences:

**Resource model — reuse the existing one.** The managed resource already has a model struct
with `types.String` / `types.Int32` / `timeouts.Value` fields and a
`map<X>ResponseToModel` function. Use both. Do not write a second model or a second mapper.

**Timeouts must be explicitly nulled**, and the attribute-type map must match the resource's
own `timeouts.Opts` **exactly**:

```go
	Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
		"create": types.StringType,
		"update": types.StringType,
		"delete": types.StringType,
	})},
```

There are three different key sets among the five existing list resources
(`create/update/delete`, `create/read/update/delete`, `create/read/delete`). Read the target
resource's `timeouts.Block(ctx, timeouts.Opts{...})` and mirror those keys. Copying the wrong
map is a runtime type mismatch, not a compile error.

---

## §3 — Registration

```go
// ListResources returns the list of list resources for the package.
func ListResources() []func() list.ListResource {
	return []func() list.ListResource{
		New<X>ListResource,
	}
}
```

Then one line in `internal/framework/provider/provider.go`'s `ListResources`:

```go
		<service>.ListResources(),
```

No `*schema.Provider`, no nil-guards, no `ResourcesMap`.

---

## Inconsistencies among the five existing examples — do not propagate these

1. **pgsqlv2 (`:53`) and inmemorydbv2 (`:49`) carry a bare `//nolint:errcheck`** on the
   deliberate `req.Config.GetAttribute` discard. Both are wrong twice over: `GetAttribute`
   returns `diag.Diagnostics`, not an `error`, so `errcheck` has nothing to report; and a bare
   directive fails `nolintlint`'s `require-explanation` (`.golangci.yml:82-84`). They survive
   only because CI lints the diff (`--new-from-rev`) and they predate it. **Use the mariadbv2
   form — no directive at all.**
2. **objectstorage bucket returns `(mapped, diags-with-error)`** on a per-item failure, which
   makes one bad bucket truncate the whole listing. The contract is `return nil, diags` to
   skip, `return mapped, nil` to include.
3. **Bucket does a per-item extra API call inside the mapper**, before the filter check —
   N+1 (or 2N) requests. Unavoidable for the S3 API, a bad default for a Cloud API resource
   where one `Depth(1)` call returns everything.
4. **mariadbv2 and inmemorydbv2 have their doc comments stripped.** Don't teach the
   stripped-down form.
5. **Filter-field vocabulary drifts** — four use `("name","location")`, bucket
   `("name","region")`, accesskey `("id","description","accesskey")`. Filtering on `id` is
   redundant with the identity. Match the *resource's own attribute names*, not a house style.
6. **None of the five has a test.** That is a gap, not a convention — write one. The harness in
   `test-harness.md` works for this branch too; drop the `ListResourceSchemas` /
   `GetResourceIdentitySchemas` assertions that exist specifically to prove `RawV6Schemas`
   worked.
