# Failover

Transport-level endpoint failover for the IONOS Go SDKs.

`failover.RoundTripper` is an `http.RoundTripper` wrapper that retries a request
against several configured server endpoints when the underlying transport
returns a network-level error, or when the response carries a status code that
was explicitly configured for failover. It sits *below* the product-level
`callAPI` retry loop, so failover happens before the SDK ever sees a response.

## Quick start

```go
import (
    "net/http"
    "time"

    "github.com/ionos-cloud/sdk-go-bundle/shared/failover"
)

endpoints := []failover.Endpoint{
    {URL: "https://api.example.com"},
    {URL: "https://api-backup.example.com", SkipTLSVerify: true},
}

opts := failover.Options{
    Strategy:              failover.RoundRobin,
    MaxRetries:            3,
    RetryOnTimeout:        true,
    FailoverOnStatusCodes: []int{502, 503},
    ExponentialBackoff: &failover.ExponentialBackoffOptions{
        InitialInterval: 500 * time.Millisecond,
        MaxInterval:     5 * time.Second,
    },
}

client := &http.Client{
    Transport: failover.NewRoundTripper(endpoints, opts, http.DefaultTransport),
}
```

Via the config file (`shared/fileconfiguration`), the same options are read from
a `failover:` block:

```yaml
failover:
  strategy: roundRobin
  retryableMethods:
    - GET
    - PUT
  retryOnTimeout: true
  failoverOnStatusCodes:
    - 502
    - 503
```

## Request flow

```
HTTPClient.Do
   └─ RoundTripper.RoundTrip
        ├─ strategy none / empty      → pass through to base transport
        ├─ strategy unknown           → error (configuration error)
        └─ strategy roundRobin        → orderedRoundTrip
                 ├─ method not retryable → pass through to base transport
                 └─ for attempt = 0 .. MaxRetries:
                        server = endpoints[attempt % len(endpoints)]
                        rewrite scheme/host, clone request (GetBody)
                        send via that endpoint's transport
                        ├─ retryable network error  → backoff, next server
                        ├─ status in FailoverOnStatusCodes → drain body, backoff, next server
                        └─ otherwise                → return response
                    all attempts exhausted → return last error (no response)
```

Total attempts are `MaxRetries + 1` (`MaxRetries` means *retries*, not
attempts). If `MaxRetries` is `0`, `shared.DefaultMaxRetries` is used.

Only the scheme and host are swapped per endpoint; path and query are
preserved. Requests with a body must have `GetBody` set — the SDK generates
requests that way, and `retry.CloneRequestForRetry` relies on it.

## Strategies

| Strategy | Value | Behaviour |
|---|---|---|
| `failover.None` | `"none"` or `""` | Failover disabled; `RoundTrip` passes straight through to the base transport with no retry overhead. |
| `failover.RoundRobin` | `"roundRobin"` | Attempt `n` targets `endpoints[n % len(endpoints)]`, cycling through servers sequentially. |

Comparison is case-insensitive and whitespace-trimmed (`NormalizeStrategy`).
`roundRobin` with an empty endpoint list is an error; an unrecognised strategy
string is surfaced immediately as a configuration error.

Adding a strategy means adding a case to the `switch` in `RoundTrip` that builds
its own `serverOrder` function (`func(attempt int) int`) — or dispatches to an
entirely different method.

## What triggers a failover

**Network errors** — failover happens only for an allowlist of transport-level
errors. The error must be (or wrap) a `*net.OpError` carrying one of these
syscall errnos:

| Errno                                   | Typical cause                                                                    | Failover                         |
|-----------------------------------------|----------------------------------------------------------------------------------|----------------------------------|
| `syscall.ECONNREFUSED`                  | Nothing listening on the target port — server down or not yet started.           | always                           |
| `syscall.ECONNRESET`                    | Peer sent RST: server crashed, restarted, or a middlebox dropped the connection. | always                           |
| `syscall.ECONNABORTED`                  | Connection aborted locally, typically an overflowing accept queue on the server. | always                           |
| `syscall.EHOSTUNREACH`                  | No route to that host — host down or unreachable from this network.              | always                           |
| `syscall.ENETUNREACH`                   | The whole target network is unreachable.                                         | always                           |
| `opErr.Timeout()` / `syscall.ETIMEDOUT` | Connection or I/O timed out.                                                     | only when `RetryOnTimeout: true` |

These errors often mean that no reply was received from that endpoint, but
they do not prove the server did not process the request (for example, a
connection reset can happen after the request is applied). Retrying
non-idempotent methods can therefore duplicate side effects; keep the default
idempotent method allowlist unless the API is known to make retries safe.

Matching is done with `errors.Is`, so wrapped errors are detected — see
`isRetryableNetOpError`. An errno that arrives *not* wrapped in a `*net.OpError`
is not retried; the `*net.OpError` type check is the outer gate.

**HTTP status codes** — any code listed in `FailoverOnStatusCodes`. The response
body is drained and closed, then the request goes to the next server. Response
headers (e.g. `Retry-After`) are *not* inspected at this layer.

## What never triggers a failover
| Case                                                       | Why                                                                                                                                                                                                                                                                                                      |
|------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| DNS errors (`*net.DNSError`)                               | Indicate misconfiguration, not a server being down.                                                                                                                                                                                                                                                      |
| Context cancelled / deadline exceeded (`ctx.Err() != nil`) | The caller aborted; retrying is pointless and hides the cancellation.                                                                                                                                                                                                                                    |
| TLS / certificate failures                                 | Deterministic configuration or protocol problems: `tls.CertificateVerificationError`, `x509.UnknownAuthorityError`, `x509.HostnameError`, `x509.CertificateInvalidError`, `x509.SystemRootsError`, `tls.RecordHeaderError`, `tls.AlertError`. Failing over would just repeat the same handshake failure. |
| Redirect and other protocol errors                         | Not in the allowlist; there is no `url.Error` catch-all.                                                                                                                                                                                                                                                 |
| Non-retryable HTTP methods                                 | Passed straight to the base transport.                                                                                                                                                                                                                                                                   |

Default retryable methods are the safe/idempotent ones: `GET`, `HEAD`, `PUT`,
`DELETE`, `OPTIONS`. Setting `RetryableMethods` replaces that set entirely.
Adding `POST`/`PATCH` — especially together with `RetryOnTimeout` — can produce
duplicate side effects, since a timed-out request may already have been applied
server-side.

## Per-endpoint TLS

Each `Endpoint` may carry its own `SkipTLSVerify` / `CertificateAuthData`. When
either is set, a dedicated `http.Transport` is built for that endpoint at
construction time (`shared.CreateTransport`); otherwise the endpoint shares the
base transport passed to `NewRoundTripper`. Endpoints, transports, and options
are snapshotted in the constructor, so mutating the original slices or struct
afterwards has no effect on a live instance.

## Backoff

Backoff between attempts is exponential, via `cenkalti/backoff/v5`, and is
skipped after the final attempt. `ExponentialBackoffOptions` overrides only the
fields you set; the rest keep the library defaults:

| Field                 | Default |
|-----------------------|---------|
| `InitialInterval`     | 500ms   |
| `MaxInterval`         | 60s     |
| `Multiplier`          | 1.5     |
| `RandomizationFactor` | 0.5     |

Setting `Multiplier` to `1.0` yields constant backoff; combined with
`RandomizationFactor: 0` you get exact fixed delays. Sleeping honours the
request context, so cancellation interrupts a pending backoff.

## Interaction with the application-level retry loop

`retry.DoWithApplicationRetry` (the shared equivalent of each product's
`callAPI`) wraps `client.Do`, which in turn invokes `RoundTrip`. The two layers
compose as follows:

- Every application-level retry starts a **fresh** `RoundTrip`, which cycles
  through the servers from the beginning. Worst case total attempts:
  `Configuration.MaxRetries × failover.Options.MaxRetries`.
- If all failover attempts fail, `RoundTrip` returns an **error, not a
  response**. The application-level loop therefore returns immediately and its
  own status-code retry logic never runs.
- That matters for overlapping codes: putting `502`/`503`/`504` in
  `FailoverOnStatusCodes` moves their handling to the transport, bypassing the
  application-level fixed-backoff retry (and, for `429`, the `Retry-After`
  handling). Codes *not* listed pass through as a normal response and are
  retried by the application layer as usual.

## Options reference

| Option                  | Type                         | Meaning                                                                                  |
|-------------------------|------------------------------|------------------------------------------------------------------------------------------|
| `Strategy`              | `Strategy`                   | `none` (default) or `roundRobin`.                                                        |
| `RetryableMethods`      | `[]string`                   | Methods eligible for failover. Empty → `GET`, `HEAD`, `PUT`, `DELETE`, `OPTIONS`.        |
| `RetryOnTimeout`        | `bool`                       | Also fail over on transport timeouts. Never applies to context errors.                   |
| `FailoverOnStatusCodes` | `[]int`                      | Status codes that trigger failover to the next server.                                   |
| `MaxRetries`            | `int`                        | Number of retries (total attempts = `MaxRetries + 1`). `0` → `shared.DefaultMaxRetries`. |
| `ExponentialBackoff`    | `*ExponentialBackoffOptions` | Backoff tuning; `nil` → library defaults.                                                |

All fields carry `json`/`yaml` tags, so `Options` can be embedded in
configuration files as-is.

## Debugging

Every decision is logged through `shared.LogDebug` with a `[Failover]` prefix —
attempt index, server index and URL, method, target URL, and whether an error or
status code was treated as retryable. Enable SDK debug logging to trace a
failover sequence.

## Files

- `failover_roundtripper.go` — `RoundTripper`, `Options`, strategies, error classification
- `failover_roundtripper_test.go` — behavioural tests (uses typed errors such as `*net.OpError` and `*net.DNSError`, not string-matched errors)
- `helpers_test.go` — test helpers
