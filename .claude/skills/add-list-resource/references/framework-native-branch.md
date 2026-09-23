# Framework-native branch

Resource already on terraform-plugin-framework (`internal/framework/services/<service>/resource_*.go`):
**no `RawV6Schemas`, no protocol-schema bridge, no second package** — the list resource is more methods on
the *same* struct. Read (don't copy) `internal/framework/services/pgsqlv2/resource_pg_cluster_list.go` for
the regional, two-attribute shape; `.../objectstoragemanagement/resource_object_storage_accesskey_list.go`
for single-`id`/non-regional.

## §1 — Identity: on the *managed* resource, not the list resource

All in `pgsqlv2/resource_pg_cluster_v2.go`: `_ resource.ResourceWithIdentity` + `<x>IdentityModel` `:26-35`,
`IdentitySchema` `:107-113`, `resp.Identity.Set` on **every** state-producing path — Create `:377`, Read
`:410`, Update `:478`, not Delete — and the identity branch of `ImportState` `:523-535`, ahead of the
string-import path. Lone `id` → shared `identity.Model` (accesskey `:62`).

## §2 — The list resource

`_ list.ListResource` / `_ list.ListResourceWithConfigure` on the same struct, plus
`New<X>ListResource() list.ListResource { return &<x>Resource{} }` (`resource_pg_cluster_list.go:19-22,
:32-35`). No `Metadata`, `Configure`, `RawV6Schemas` or constructor arg — all inherited.
`ListResourceConfigSchema`, `List` and the mapper match `sdkv2-branch.md` §2b/§2c; only the payload differs:

- **Reuse the resource's model and `map<X>ResponseToModel`; never write a second of either** —
  `resource_pg_cluster_list.go:105-115` fills `resource_pg_cluster_v2.go:41`'s `clusterResourceModel`.
  Correctness, not convenience.
- **Null timeouts explicitly** (`:107-111`) with an `attr.Type` map matching the resource's own
  `timeouts.Opts` keys exactly — key sets differ across the five; a wrong map is a runtime mismatch, not a
  compile error.

Register it in the package's `ListResources()` (`pgsqlv2/resources.go:15-19`); add the package to
`internal/framework/provider/provider.go:339-346` if absent.

## Defects in the five — do not propagate

1. Bare `//nolint:errcheck` on the `req.Config.GetAttribute` discard (pg list `:53`, inmemorydb `:49`) is
   doubly wrong: `GetAttribute` returns `diag.Diagnostics`, not `error`, and bare directives fail
   `nolintlint: require-explanation` (`.golangci.yml:82-84`). Omit it (mariadbv2 does).
2. `objectstorage/resource_bucket_list.go` returns `(mapped, diags-with-error)` per item (`:76,:80`), so one
   bad item truncates the listing; return `nil, diags` to skip, `mapped, nil` to include. Its per-item
   `GetBucketLocation` (`:51`, before the filter check `:57`) is N+1, S3-forced.
3. Filter vocabulary drifts (`name,location` ×3; bucket `name,region`; accesskey `id,description,accesskey`)
   — use the resource's own attribute names; `id` duplicates the identity.
4. None of the five has a test — a gap, not a convention. `test-harness.md` applies, minus its
   `ListResourceSchemas`/`GetResourceIdentitySchemas` assertions (they prove `RawV6Schemas`).
