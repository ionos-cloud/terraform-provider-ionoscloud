# Implementation Notes — InMemoryDB v2

This file tracks changes made inside the Terraform provider that depend on work to be done in **external repositories**, as well as open issues that need to be resolved before this implementation is production-ready.

---

## 1. SDK must be published to `sdk-go-bundle`

The generated SDK lives locally at `sdks/generated_locally` and is wired into the provider via a `replace` directive in `go.mod`:

```
replace github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3 => /home/adeatcu/work/claude-workspace/sdks/generated_locally
```

**Required action:** The SDK needs to be reviewed, published, and released as a proper sub-module of `sdk-go-bundle` at path `products/dbaas/inmemorydb/v3` (module major version 3, following the same pattern as `products/dbaas/psql/v3`). Once published, the `replace` directive in the provider's `go.mod` must be removed and replaced with a real version pin.

---

## 2. `InMemoryDBV2` constant needed in `sdk-go-bundle/shared`

The provider uses `fileconfiguration.InMemoryDBV2 = "inmemorydbv2"` (from the `shared/fileconfiguration` package) to look up endpoint overrides from the user's config file. This constant does not exist in the upstream `sdk-go-bundle/shared` repository — it was added manually to `vendor/` as a workaround.

**Required action:** Add the following constant to `sdk-go-bundle/shared/fileconfiguration/fileconfiguration.go`, alongside the existing `PSQLV2`:

```go
InMemoryDBV2 = "inmemorydbv2"
```

---

## 3. SDK generation: `Version` constant/type collision

**Problem:** The openapi-generator always emits a `Version` constant in `client.go`:

```go
const Version = "1.0.0"
```

The InMemoryDB v2 spec defines a schema named `Version` (the entity describing a supported cluster version), which the generator turns into:

```go
type Version struct { ... }
```

Go does not allow the same identifier to be both a constant and a type in the same package, so the generated SDK does not compile.

**Local workaround (applied in vendor/):** Renamed `const Version` → `const SDKVersion` in `vendor/github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3/client.go`. The provider references `inmemorydbv3sdk.SDKVersion` accordingly.

**Proper fix options:**
1. Rename the `Version` constant in the openapi-generator Go template to `SDKVersion` (fixes all future SDKs with a `Version` model).
2. Rename the `Version` schema in the InMemoryDB v2 spec to something more specific (e.g. `ClusterVersion`), following the pgsql convention of `PostgresVersion`.

---

## 4. Production endpoints

The `LocationToURL` map in `services/dbaas/inmemorydbv2/client.go` currently uses `dev.` endpoints:

```go
"de/txl": "https://dev.in-memory-db.de-txl.ionos.com/v2",
```

**Required action:** Once the v2 API is promoted to production, replace all `dev.` prefixes with the production URLs (matching the swagger spec servers):

```go
"de/txl": "https://in-memory-db.de-txl.ionos.com/v2",
```

---

## 5. v2 API not yet fully deployed (production)

The production endpoints (`https://in-memory-db.*.ionos.com/v2`) return `400 paas-validation-1: no matching operation was found` for all v2 routes. The dev endpoints work for cluster creation (confirmed: cluster reaches AVAILABLE state) but provisioning is slow (~6 minutes) and the dev environment appears unstable for long-running operations.

**Required action:** Backend team to promote the v2 API to production before acceptance tests can be run against production endpoints.

---

## 6. `snapshotHours` not returned in GET response

When creating a cluster with `snapshot.snapshotHours = [0, 12]`, the API accepts the value but does not return it in subsequent GET responses (`snapshotHours` is omitempty in the SDK model). This causes a Terraform state inconsistency error after apply:

```
.snapshot.snapshot_hours: was [0, 12], but now null
```

**Required action:** Confirm with the API team whether `snapshotHours` should be returned in GET responses. If yes, fix the API. If no (intentionally write-only), the provider needs to preserve the planned value from state rather than reading it from the API response.
