# Framework-native branch

For a resource in `internal/framework/services/<service>/resource_*.go` (paths below relative to
`internal/framework/services/`): **no `RawV6Schemas`, no second package**; the list resource is
more methods on the *same* struct. Read, don't copy, `pgsqlv2/resource_pg_cluster_list.go` (pg
list: regional, two attributes) and `objectstoragemanagement/resource_object_storage_accesskey_list.go`
(a lone `id`, not regional).

## §1 — Identity on the managed resource

`pgsqlv2/resource_pg_cluster_v2.go`: `_ resource.ResourceWithIdentity` + `<x>IdentityModel`
`:26-35`, `IdentitySchema` `:107-113`, `resp.Identity.Set` on **every** state-producing path
(Create `:377`, Read `:410`, Update `:478`; not Delete), `ImportState`'s identity branch `:523-535`
before the string import. A lone `id` → the shared `identity.Model` (accesskey list `:62`).

## §2 — The list resource

`_ list.ListResource`, `_ list.ListResourceWithConfigure`,
`New<X>ListResource() list.ListResource { return &<x>Resource{} }` (pg list `:19-22`, `:32-35`); no
`Metadata`, `Configure` or constructor argument: inherited. `ListResourceConfigSchema`, `List` and
the mapper as `sdkv2-branch.md` §2b/§2c, except: **reuse the resource's model and
`map<X>ResponseToModel`**, never a second (pg list `:105-115`); **null timeouts explicitly**
(`:107-111`), keying the `attr.Type` map exactly as the resource's `timeouts.Opts` (pg's lacks
`read`): a wrong map fails only at runtime. Register in the package's `ListResources()`
(`pgsqlv2/resources.go:15-19`); add the package to `internal/framework/provider/provider.go:339-346`
if absent. Tests and the acceptance step: `test-harness.md`.

## Defects in the five — do not propagate

1. `//nolint:errcheck` on the `req.Config.GetAttribute` discard (pg list `:53`, inmemorydbv2 `:49`):
   it returns `diag.Diagnostics`, and a bare directive fails `nolintlint`. Omit it (mariadbv2 does).
2. `objectstorage/resource_bucket_list.go`'s `(mapped, diags-with-error)` returns
   (`sdkv2-branch.md` §2c); its N+1 per-item `GetBucketLocation` is S3-forced.
3. Filter names drift: use the resource's own attribute names; `id` duplicates the identity.
4. None of the five has a test — a gap, not a convention.
