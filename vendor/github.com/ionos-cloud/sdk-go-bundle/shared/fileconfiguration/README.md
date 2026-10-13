# fileconfiguration

Loads the IONOS YAML configuration file (`~/.ionos/config` by default) and
exposes lookup helpers for per-product, per-location endpoint overrides,
credentials, and transport-level failover options.

This package only *reads* configuration. Applying an override to a live client
is the caller's job — see [Applying overrides](#applying-overrides).

## Loading

```go
import "github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

cfg, err := fileconfiguration.New("")                    // default path
cfg, err := fileconfiguration.New("/path/to/config")     // explicit path
cfg, err := fileconfiguration.NewFromEnv()               // path from IONOS_CONFIG_FILE

profiles := fileconfiguration.ReadProfilesFromFile()     // profiles only, never errors
```

| Function | Behaviour |
|---|---|
| `New(path)` | Empty `path` → `DefaultConfigFileName()`. Errors if the file is missing, unreadable, empty, or invalid YAML. |
| `NewFromEnv()` | `New(os.Getenv("IONOS_CONFIG_FILE"))` — falls back to the default path when the variable is unset. |
| `DefaultConfigFileName()` | `$HOME/.ionos/config`; errors if the home directory can't be determined. |
| `ReadProfilesFromFile()` | Best-effort read of just the `currentProfile` + `profiles` keys. Returns `nil` when the config path is unavailable or the file cannot be read; YAML unmarshal errors are ignored. |

### Environment variables

| Variable | Effect |
|---|---|
| `IONOS_CONFIG_FILE` | Path to the config file. Used by `NewFromEnv` and `ReadProfilesFromFile`. |
| `IONOS_CURRENT_PROFILE` | Overrides `currentProfile` from the file. Applied inside `New` **and** re-checked by `GetCurrentProfile`, so it wins even if you set `cfg.CurrentProfile` manually afterwards. |

## File format

```yaml
version: 1.0
currentProfile: testProfile

profiles:
  - name: testProfile
    environment: testEnvironment
    credentials:
      username: testUser
      password: testPass
      token: testToken
      s3AccessKey: ...
      s3SecretKey: ...

environments:
  - name: testEnvironment
    certificateAuthData: testCertData      # environment-wide CA
    products:
      - name: mariadb
        endpoints:
          - location: de/fra               # location-based endpoint
            name: mariadb.de-fra.ionos.com
            skipTlsVerify: false
          - location: de/txl
            name: mariadb.de-txl.ionos.com
            certificateAuthData: certauthdata
            skipTlsVerify: true
      - name: psql
        endpoints:
          - name: api.ionos.com/databases/postgresql   # global endpoint (no location)

failover:                                  # optional, see shared/failover
  strategy: roundRobin
  retryableMethods: [GET, PUT]
  retryOnTimeout: true
  failoverOnStatusCodes: [502, 503]
```

Structure: **profile** → names an **environment** → holds a list of
**products** → each with a list of **endpoints**.

`Version` is a `float64` wrapper with a custom `MarshalYAML`, so `1.0` is
emitted as `1.0` rather than `1` when the config is written back out.

### Endpoint fields

| Field | Meaning |
|---|---|
| `name` | The endpoint host (no scheme). Callers prepend `https://` — `shared.NewConfigurationFromOptions` does this via `getServerUrl`. |
| `location` | Region, e.g. `de/fra`. **Empty means global** — that distinction drives every lookup below. |
| `skipTlsVerify` | Disable TLS verification for this endpoint. |
| `certificateAuthData` | PEM CA bundle for this endpoint; falls back to the environment's `certificateAuthData`. |

### Product name constants

Use the exported constants rather than raw strings — they are the names the
lookups match (case-insensitively, whitespace-trimmed).

- **Regionless** (global endpoint override): `Autoscaling`, `APIGateway`, `CDN`,
  `Cert`, `Cloud`, `ContainerRegistry`, `DNS`, `Mongo`,
  `ObjectStorageManagement`, `PSQL`
- **Location-based** (endpoint per region): `InMemoryDB`, `InMemoryDBV2`,
  `Kafka`, `Logging`, `Mariadb`, `MariaDBV2`, `Monitoring`, `NFS`,
  `ObjectStorage`, `VPN`, `PSQLV2`

## Lookups

Every method is nil-receiver safe and returns `nil`/zero rather than panicking,
so `fileconfiguration.New(...)` failures can be ignored where a missing config
file is acceptable.

### Profiles and environments

```go
cfg.GetProfileNames()          // []string
cfg.GetEnvironmentNames()      // []string
cfg.GetCurrentProfile()        // *Profile, nil if unset or not found
cfg.GetEnvForCurrentProfile()  // environment name of the current profile, "" if none
```

`GetCurrentProfile` matches case-insensitively and logs a warning when no
current profile is set or the named profile is absent.

### Endpoints

```go
cfg.GetProductOverrides("mariadb")                        // *Product
cfg.GetProductLocationOverrides("mariadb", "de/fra")      // *Endpoint, exact location match
cfg.GetProductGlobalOverrides("psql", 0)                  // *Endpoint, n-th global endpoint
cfg.GetOverride("dns", "de/fra")                          // location, else first endpoint
cfg.GetLocationOverridesWithGlobalFallback("dns", "de/fra") // location, else first global
cfg.FilterOverrides("nfs", func(e Endpoint) bool { ... })  // []Endpoint by predicate
cfg.FilterGlobalOverrides("nfs")                          // []Endpoint, location == ""
cfg.FilterLocationOverrides("nfs")                        // []Endpoint, location != ""
```

`GetProductOverrides` scopes the search to the current profile's environment.
**If no current profile (or no environment) is resolvable, the environment
filter is skipped** and the first environment containing a product with that
name wins — worth knowing when a config defines the same product in several
environments.

#### `GetOverride` vs `GetLocationOverridesWithGlobalFallback`

Both try the exact location first, then fall back. They differ in what they fall
back *to*:

| | Fallback target | Fails when |
|---|---|---|
| `GetOverride` | `Endpoints[0]`, whatever it is | `Endpoints[0]` is location-based and its location ≠ the requested one → returns `nil` with a debug log, **even if global endpoints exist further down the list** |
| `GetLocationOverridesWithGlobalFallback` | the first endpoint with an empty `location` | no global endpoints exist for that product |

Prefer `GetLocationOverridesWithGlobalFallback` when the endpoint list mixes
global and location-based entries; it preserves declaration order among globals
and won't be tripped up by a location-based first element. Pass `location: ""`
to `GetOverride` for a regionless product.

### Failover

```go
opts := cfg.GetFailoverOptions()   // *failover.Options, nil when no failover: block
```

The `failover:` block deserializes straight into
[`failover.Options`](../failover/README.md) — it's the same struct the
`RoundTripper` consumes, so any field documented there is valid YAML here. The
package does not wire it up; propagating it to the runtime configuration is the
caller's responsibility.

## Applying overrides

`fileconfiguration` returns data; the `shared` package applies it:

```go
ep := cfg.GetLocationOverridesWithGlobalFallback(fileconfiguration.Mariadb, "de/fra")
if ep != nil {
    shared.OverrideLocationFor(client, "de/fra", "https://"+ep.Name, false)
}
```

`shared.OverrideLocationFor` rewrites the matching `ServerConfiguration` in
place (or appends one), marking it with the `shared.EndpointOverridden`
description so overridden servers can be identified later. Pass
`replaceServers: true` to discard all configured servers in favour of the single
override. TLS settings from the endpoint go through
`shared.SetSkipTLSVerify` / `shared.CreateTransport`.

## Files

- `fileconfiguration.go` — types, loaders, lookup helpers, product name constants
- `fileconfiguration_test.go` — loading, profile/environment resolution, override precedence, failover deserialization
