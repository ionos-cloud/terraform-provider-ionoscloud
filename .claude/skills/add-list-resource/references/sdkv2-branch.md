# SDKv2-backed branch

For a managed resource still on `terraform-plugin-sdk/v2` (`ionoscloud/resource_*.go`). Reference
implementation, all on master, cited one thing at a time — **do not copy it wholesale**, it
models two constructs this file forbids (§2a, §2c):

| file | lines |
|---|---|
| `ionoscloud/resource_datacenter.go` | identity `:27-47`, read `:167-195`, setter `:337-355`, importer `:272-309`, resolver `:311-335` |
| `ionoscloud/resource_datacenter_list.go` | struct+ctor `:34-53`, four methods `:55-94`, `List` `:96-126`, mapper `:128-170` |
| `internal/framework/identity/` | all the list resource calls: `sdkv2.go`, `filter.go` |

**Provenance.** `ionoscloud_datacenter` (PR #1034) is the only shipped, independently reviewed
instance, and on master the only resource with an `Identity` at all. Anything else you find
registered came from an earlier run of *this skill* — worth reading, not independent corroboration;
where this file argues a rule from one of those, it is quoting itself.
`scripts/sdkv2.sh identities` lists the current set; `references/why.md` carries rationale a
normal run does not need.

---

## §1 — Resource identity on the SDKv2 resource

**Hard prerequisite.** A resource with no `Identity` cannot be listed: `ProtoIdentitySchema(ctx)`
returns a **nil func** (not a func returning nil, `core_schema.go:419-423`), `SetRawV6Schemas` gives
up, and the only trace is a `tflog.Error`, `RawV6Schemas` having no diagnostics channel. That
failure is not local: see §3.

### 1a. The `Identity` block

Goes in `resource<Resource>()` beside `Importer`; `resource_datacenter.go:27-47` is it verbatim. An
identity attribute takes exactly one of `RequiredForImport`/`OptionalForImport`, a scalar type or a
`TypeList` of scalars, and nothing but a `Description` — `provider_test.go`'s `InternalValidate()`
fails a *unit* test on any violation. Mark one Required only if the import cannot proceed without it;
datacenter's `location` is Optional because `NewCloudAPIClient(ctx, "")` falls back to the global
endpoint.

**Only declare what addresses the resource.** `id` + `location` is not a template:
`ionoscloud_target_group` has no `location` attribute at all, its collection being global, so its
identity would be a lone `id`. It declares no identity yet, so there is no file to copy that from —
treat it as the reasoning, not a worked example. Never add an attribute for symmetry; every one must
be set on every read (1b). A **child resource** adds one `RequiredForImport` string per parent ID,
spelled as its own schema spells it, before `id` — but none in the repo declares an identity yet, so
that is extrapolation from the composite import IDs, not shipped code. Say so if you build one.

### 1b. The identity setter

`setDatacenterIdentity` (`resource_datacenter.go:337-355`): `d.Identity()`, then one
`identity.Set(...)` per declared attribute, each error wrapped with that name and `d.Id()`.

- **Set only the attributes you declared in 1a.** Setting an undeclared one errors through
  `MapFieldWriter` and **panics rather than returning under `TF_ACC`**
  (`helper/schema/schema.go:657-659`: `panicOnError` is `os.Getenv("TF_ACC") != ""`).
- **Set every attribute, on every read, forever.** `ReadResource` compares the whole identity with
  `RawEquals` and fails with "Unexpected Identity Change" on any difference. Setting only *some*
  escapes the missing-identity error (that check is "null or *every* attribute null") and walks into
  the stability check instead.
- **This is also the list resource's identity mapper**, verbatim (§2c), because it reads out of
  state and never from an API response.

### 1c. Call sites

Read (`resource_datacenter.go:167-195`): the 404 branch returns **before** setting an identity,
correctly; the tail is writer *then* identity setter, the identity reading attributes only the
writer fills in. Create needs the same call unless it delegates to Read; Delete needs nothing.

**Update is different, and the reason is easy to get wrong.** The SDK carries the prior identity into
the apply itself, so "Missing Resource Identity After Update" does *not* fire; an explicit call there
is only a safety net for state written before the resource declared an identity. Say that in the
comment; do not copy the Create rationale onto Update. What Update writes must also equal the prior
identity byte for byte (`grpc_provider.go:1585`, `RawEquals`), so derive it from state — or end
Update with `return resource<Resource>Read(...)`, which does both and is what datacenter does. If it
does not delegate, the acceptance test needs an update step changing a **non-ForceNew** attribute or
Update is never entered: `ionoscloud/resource_ipblock_test.go:109` feeds the update step
`testAccCheckIPBlockConfigUpdate` (`:207-212`), which flips `size` from 1 to 2 — and `size` is
`ForceNew` (`ionoscloud/resource_ipblock.go:40-44`), so that step is a destroy-and-create covering
nothing.

### 1d. The importer — dual mode, and the `d.SetId` that is easy to miss

Terraform sends an **empty `req.ID`** for an identity-based import; the real identifier is only in
`d.Identity()`. `resourceDatacenterImport` (`resource_datacenter.go:272-309`) and its resolver
`datacenterImportParts` (`:311-335`) are the shape end to end: resolver → `d.SetId(id)` → client →
fetch → writer → setter. The identity branch is `:316-322`; below it is the legacy string
path, unchanged — **a pure prepend**, so keep its delimiter, part count and error message. A non-nil
`identityErr` means "fall back to the string import", since `d.Identity()` errors rather than
returning nil when no identity schema is declared. Child resources need `d.Set("<parent>_id", …)`
before the setter reads it back.

**Additive only.** Prepend the resolver and its `d.SetId(id)`, append `set<Resource>Identity(d)`
after the writer, and leave every other line of that function as you found it — including a
post-fetch `d.SetId(*obj.Id)` the new one makes redundant. `resourceIpBlockImporter` lost exactly
that line on the first pass: a behaviour change to a reviewed function, in a diff that was supposed
to only add an import mode. The same holds file-wide, `Read` and `Update` included: this task adds
an identity, it does not get to re-word comments or tidy code it happens to be reading.

**The empty-ID guard.** `identity.GetOk("id")` returns `exists=false` for the schema type's **zero
value**, so an empty identity `id` falls through to the legacy path, where `splitImportID`
(`ionoscloud/utils.go`) yields `[""]` and `validateImportIDParts` rejects it. No API call with an
empty ID is reachable. **Copilot will flag this as a missing validation. It is a false positive** —
this exact exchange happened on PR 1034. **But that invariant is `splitImportID`'s, not the
resolver's.** A resource with a plain (non-composite) import ID never calls `splitImportID`, and then
nothing rejects `""`: the fall-through returns `d.Id()` verbatim and the API is called with an empty
ID. If your resolver skips `splitImportID`, add the guard yourself:
`if d.Id() == "" { return "", fmt.Errorf("invalid import identifier: expected a <resource> UUID, got an empty string") }`.
That snippet returns **two** values because a plain-import-ID resource is normally locationless, its
resolver being `func <resource>ImportParts(d) (string, error)` where datacenter's returns three —
match whichever arity yours has. No resource in the repo carries that guard yet, because no
plain-import-ID resource has been given an identity; `resourceTargetGroupImport`
(`ionoscloud/resource_target_group.go:339-347`) is the unguarded shape to recognise, taking
`groupIp := d.Id()` straight into `TargetgroupsFindByTargetGroupId`.

---

## §2 — The list resource

File: `ionoscloud/resource_<resource>_list.go`, package `ionoscloud`. **There is no package decision
to make** — it goes next to the resource it lists, always (`references/why.md` has the import-cycle
proof). Framework-native list resources go in their product's package.

### 2a. Header, assertions, struct, constructor

`resource_datacenter_list.go:3-53`, with the change in the first bullet. Its three `var _ list.X =
(*T)(nil)` assertions at `:34-38` — `ListResource`, `ListResourceWithConfigure`,
`ListResourceWithRawV6Schemas` — are load-bearing (§2b); copy all three.

- **Use `constant.<Resource>Resource`, not a local type-name const.** The type name must match the
  `ResourcesMap` key in `ionoscloud/provider.go` exactly or terraform has no managed resource to
  attach results to, and that map is keyed with `constant.*Resource`.
- **The reference implementation does not follow that rule — this is a known divergence, not a
  precedent.** `ionoscloud/resource_datacenter_list.go:32` declares
  `const datacenterResourceType = "ionoscloud_datacenter"` and passes it to both `SetRawV6Schemas`
  and `Metadata`, even though `constant.DatacenterResource` exists and is what keys `ResourcesMap`.
  It predates the rule and has not been migrated. Do not read it as sanctioning a new one.
- **The constructor takes no arguments** — `New<Resource>ListResource() list.ListResource` returns
  `&<resource>ListResource{resourceSchema: resource<Resource>()}`, so it *is* a
  `func() list.ListResource` and registration needs no wrapper closure (§3). It runs once; the
  `*schema.Resource` is never mutated and never rebuilt per item.
- **Import `internal/framework/identity` as `fwidentity`**: here `identity` already means the SDKv2
  `*schema.IdentityData` from `d.Identity()`.

### 2b. The four boilerplate methods

`resource_datacenter_list.go:55-94`, substituting `constant.<Resource>Resource` for the local const
in `RawV6Schemas` and `Metadata` (which sets the **full** type name, no `req.ProviderTypeName`
prefix). `Configure`'s `req.ProviderData == nil` early return is mandatory. Filter fields go in
`ListResourceConfigSchema` as `fwidentity.FiltersKey: fwidentity.FilterAttribute("<f1>", "<f2>")`;
with *no* arguments they go unvalidated and a typo matches nothing instead of failing the plan.
`Metadata` and `Configure` take `resource.MetadataRequest` / `resource.ConfigureRequest`, **not** the
`list.*` decoys also in the vendored package — a *compile* error only because of the three `var _`
assertions. Delete the `ListResourceWithConfigure` one and it goes silent: `Configure` never runs,
`r.bundle` stays nil, `List` nil-derefs.

### 2c. `List` and the mapper

`resource_datacenter_list.go:96-170` is the Cloud API (`sdk-go/v6`) shape end to end: `List` wraps
`fwidentity.StreamList(ctx, stream, req, fetch, r.map<Resource>)`, the mapper nil-guards, filters,
then does §2d. Deviations:

- **Build the filter map from the accessors, every value, inline** — no local lifted out of it, not
  even one used again later: a local is a second place to get the pairing wrong, and it reads
  unevenly next to the sibling key that did inline its accessor. `:140-141` lift `name` and
  `location` out; that part of the reference implementation is not the pattern to copy. The map's
  keys must match `FilterAttribute`'s exactly; pointer fields need `ToValueDefault`, values do not.
- **`Depth(1)` may need `.Limit(...)`** if this endpoint's SDK default is below the sibling data
  source's: `/ipblocks`, `/targetgroups` and user management default to 100, not 1000 (SKILL.md step
  0.5, decision 1).
- **Filters stay in the mapper.** `StreamList` never hands them to the fetch closure, and the regional
  trick for getting them (re-reading `req.Config` inside the fetch, to narrow a fan-out) buys nothing
  here: **no Cloud API resource is regional** (SKILL.md step 0.5 decision 4) and
  `AvailableLocations()` exists only under `services/dbaas/`, whose list resources are
  framework-native, so wanting that line means the wrong branch. Pushing a bundle helper's
  server-side filter down still never excuses the mapper's re-check (`references/why.md`).
- **Child resources only**: set the parent ids *before* the identity setter reads them back out; the
  writer usually does not, they coming from the fetch context, not the API object. The one exception
  to "no attribute mapping of your own": only parent ids, only ones the writer provably leaves unset
  (`SetRecordData` never touching `zone_id` is the live case).

#### The sdk-go-bundle deltas

If the resource's CRUD takes a `<Product>Client` off the bundle rather than `NewCloudAPIClient*`
(step 0's grep (6)), four things change:

1. **The client is the resource's own, already on the bundle** — `NewCloudAPIClientWithFailover`
   returns a `*ionoscloud.APIClient` that cannot reach a bundle product. Prefer the service package's
   list helper, so the list resource reads what the sibling data source does
   (`r.bundle.DNSClient.ListZones(ctx, "")`).
2. **Check for `Depth`, do not assume it away** —
   `scripts/sdkv2.sh depth <product> <api-stem> <RequestType>`. dns and most have none;
   `vmautoscaling` does.
3. **Values, not pointers**: `Items` is a value slice and each element's `Id`/`Properties` are
   values, so both nil guards and the `*items.Items` deref go away and the only unusable item has an
   empty id. Fields *inside* `Properties` are still a mix, so `ToValueDefault` applies to the pointer
   ones. The product package replaces `sdk-go/v6` in the imports.
4. **The writer may be a method on that client** — `r.bundle.DNSClient.SetZoneData(data, zone)` —
   taking its object **by value**, so drop the `&` the Cloud API shape passes (§2d).

#### `DisplayName` when `name` is optional

`DisplayName` labels each row of `terraform query` output, and a blank one the framework reads as a
diagnostics-only event (`fwserver/server_listresource.go:189`). Datacenter's `name` is
`Required: true`, so `shared.ToValueDefault(item.Properties.Name)` is always populated — not general:
`ionoscloud/resource_ipblock.go:30-33` has `"name"` `Optional` against an SDK `Name *string`, so an
unnamed item renders blank. **A third case: some resources have no `Name` at all**
(`ionoscloud_cdn_distribution` has `Domain`, a user `Email`) — pick the closest human identifier they
do carry, and report that judgement in step 0.5. Fall back to `*item.Id` when it comes back empty:
initialise a `displayName` local **from the accessor itself**, reassign it when empty, and pass it to
`MappedItemFromResourceData`. With no fallback to make, pass the accessor straight in — no local.
`test-harness.md` has the matching stub warning. `shared.ToValueDefault` is the repo-wide nil-safe
deref; **never declare a local copy**.

**`StreamList` semantics you must respect.** `(nil, nil)` skips the item silently;
`(nil, diags-with-error)` skips it and logs only `diags[0].Detail()` at WARN, where no user sees it;
anything **non-nil** returned with an error diagnostic or with `Identity == nil` **aborts the whole
stream**. So `return nil, diags` for a recoverable per-item problem, `return mapped, nil` to include;
`objectstorage/resource_bucket_list.go:76,:80` returns `(mapped, diags-with-error)`, truncating the
whole listing — **do not copy that**. `fetch` runs eagerly, once, before streaming: pagination lives inside
it and `req.Limit` is ignored.

### 2d. No model — the result is the resource's own state

```go
	data := r.resourceSchema.Data(&terraform.InstanceState{})
	set<Resource>Data(data, &item)     // the resource's OWN state writer, the one Read uses
	set<Resource>Identity(data)        // AFTER the writer — it reads values back out
	fwidentity.MappedItemFromResourceData(displayName, data, includeResource)
```

**Do not add a Go struct mirroring the SDKv2 schema.** `MappedItem.Identity` and
`MappedItem.Resource` are `any`, so one compiles — it is the shape this design replaced, and it
shipped `master` broken (`references/why.md`). Three rules follow, each having bitten:

- **Use `Data`, not `TestResourceData`** — the latter never assigns `timeouts`, so the value omits a
  block the result *type* still declares and `Set` rejects it on type comparison.
- **Build a fresh `ResourceData` per item**, or the previous item's values leak into whatever the
  writer leaves untouched.
- **Order: writer, then identity setter**, which reads `d.Id()` and `d.Get("location")` back out —
  and **the writer must call `d.SetId`**, or both conversions fail with `"state is nil, call SetId()
  on ResourceData first"`. Every `set<Resource>Data` here does; check yours.

`nullTimeouts` nulls the `timeouts` block out of the result, so a query equals what a refresh writes
(`ionoscloud/resource_datacenter_list_test.go:105`); absent *nested* schemas decode unevenly, and
`references/why.md` has the block-vs-attribute rule the test must match, including that `MaxItems: 1`
does not make a block single (no shipped list resource has one yet). Nothing else stays coupled to
the schema — an added attribute goes into `set<Resource>Data`, which Read needs anyway. Three things
do: the filter keys (`FilterAttribute(...)` and the `MatchesFilters` map must stay identical and name
real fields), the `Identity` block and `set<Resource>Identity` (§1), and the `DisplayName` fallback.

#### What the design requires of the state writer

The templates spell the writer `set<Resource>Data`, the majority name in `package ionoscloud`, but
**the name is not the contract, the signature is.** **Follow Read; do not grep for a name**:
`func .*SetData` needs the literal "SetData", `func set.*Data` also matches unrelated helpers, and
for a bundle product the writer is not in `ionoscloud/` at all. Read the read function —
`scripts/sdkv2.sh writer <resource>`.

The mapper can call any writer shaped `func(d *schema.ResourceData, obj <sdk>.<X>) error` — one
`ResourceData`, one object of the element type the collection returns, an `error`. All three kinds
exist on master:

- `setDatacenterData` (`ionoscloud/resource_datacenter.go:357`) — unexported, package-level, by
  pointer.
- `IpBlockSetData` (`ionoscloud/resource_ipblock.go:259`) — exported, package-level; called directly,
  no adapter.
- `SetZoneData` (`services/dns/zone.go:60`) — a **method** on the product's service client, taking
  its object by value, reachable as `r.bundle.DNSClient.SetZoneData(data, zone)`. It also calls
  `d.SetId` itself, which not every writer does — check, because the identity setter reads the id
  back out.

**A writer that needs arguments the mapper cannot supply is telling you something about the
resource.** `setBackupUnitData` also wants an `*ionoscloud.Contracts`; `setResourceServerData` and
`setGroupData` want a `context.Context` and an `*ionoscloud.APIClient` to fetch more, while the
mapper gets one already-fetched element (`internal/framework/identity/list.go:49`) and reaching for
a client costs a request per item. So: **a resource is listable under this design when its full
state is derivable from one element of the collection response at `Depth(1)`.** When it is not,
choose deliberately — list it anyway and accept the missing inputs (only if the attributes they fill
are genuinely optional; for `server` and `group` they are not), pay the per-item fetch and its N+1
cost, or do not list it yet, usually the right answer and a finding to report back. Read the writer's
signature in step 0, not after the mapper fails to compile.

---

## §3 — Registration

One line in `ionoscloud/list_resources.go`: add `New<Resource>ListResource` to the returned slice,
named directly, not wrapped in a closure, because it already has that signature (§2a). That file is itself
the pattern. **No lookup, no `ok` branch, no nil guard**: the constructor is
in-package and has no error path.

**Never register a list resource without schemas.** One whose `RawV6Schemas` supplied neither schema
makes `GetProviderSchema`, terraform's very first RPC, fail with "ListResource Type Defined without a
Matching Managed Resource Type" (`fwserver/server_listresources.go:99-109`) — killing the **whole
provider** for every operation, not just `terraform query`. The one remaining way in is **removing or
forgetting the `Identity` block on the SDKv2 resource** (§1): `SetRawV6Schemas` then sets nothing but
a `tflog.Error`, and no guard could quietly unregister the list resource instead. Treat that block as
load-bearing for the entire provider; `TestDatacenterListResource` catches it
(`ionoscloud/resource_datacenter_list_test.go:52-57`).

---

## §4 — Provider wiring

**For a new SDKv2-backed list resource there is none.** §3's one line is the whole registration and
no other file needs an edit: `provider.New` is already variadic over `sdkv2ListResources` and appends
it alongside the framework-native packages (`internal/framework/provider/provider.go:64-68`,
`:339-345`), and all six construction sites pass it the whole slice (`main.go:41`,
`xpprovider/provider.go:16`, `internal/acctest/acctest.go:55`,
`internal/framework/provider/provider_test.go:60`, `ionoscloud/resource_datacenter_list_test.go:164`
and `ionoscloud/provider_test.go:91`, which spells it unqualified because it is in-package). A
framework-native list resource in a **new** service package is the one case that still needs a
`provider.go` edit.
