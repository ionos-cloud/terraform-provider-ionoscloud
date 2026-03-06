/*
 * IONOS Shared Libraries – Failover RoundTripper
 */

package failover

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"

	boff "github.com/cenkalti/backoff/v5"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/retry"
)

// Strategy selects the endpoint failover behaviour.
// It is a string type so it serializes nicely to JSON/YAML config files.
//
// Supported values:
//   - None ("none") or "": default behaviour (no endpoint failover)
//   - RoundRobin ("roundRobin"): on network-level errors, retry the request against the next server in Servers
//
// Note: comparisons should be case-insensitive.
type Strategy string

const (
	None       Strategy = "none"
	RoundRobin Strategy = "roundRobin"
)

// Endpoint represents a server endpoint with optional per-endpoint TLS
// configuration. When SkipTLSVerify or CertificateAuthData are set, a dedicated
// http.Transport is created for that endpoint at construction time, otherwise
// the default base transport is shared.
type Endpoint struct {
	URL                 string
	SkipTLSVerify       bool
	CertificateAuthData string
}

// endpointRT pairs a resolved server URL with its own transport.
type endpointRT struct {
	url string
	rt  http.RoundTripper
}

// RoundTripper is an http.RoundTripper wrapper that retries the request
// against multiple configured servers when the underlying transport returns a
// network-level error or an HTTP status code listed in
// Options.FailoverOnStatusCodes.
//
// Notes:
//   - Network errors trigger retries with exponential backoff, cycling to the
//     next server.
//   - Status codes in FailoverOnStatusCodes also trigger retries: the response
//     body is drained and the request is sent to the next server. Response
//     headers (e.g. Retry-After) are not inspected at this layer.
//   - If the request carries a body, the request must have GetBody set (the SDK
//     generates requests in a way that supports this).
//   - For non-idempotent requests (POST, PATCH, etc.), enabling failover on
//     timeouts can produce duplicates when RetryOnTimeout is set. Context
//     cancellation always stops retries immediately.
//
// The request URL is rewritten by swapping scheme/host with each server URL,
// preserving path and query.
//
// Endpoints, transports, and options are snapshotted at construction time;
// subsequent changes to the original slices or struct do not affect this instance.
type RoundTripper struct {
	endpoints   []endpointRT
	fo          Options
	defaultBase http.RoundTripper // used for pass-through (strategy=none, empty endpoints, non-retryable methods)
}

// NewRoundTripper creates a RoundTripper that cycles through
// endpoints according to fo when a request fails. Slices in fo are copied so
// callers can safely mutate the originals.
//
// For each endpoint, a dedicated http.Transport is created when SkipTLSVerify
// or CertificateAuthData are set; otherwise defaultBase is used.
//
// When the strategy is "none" or empty, RoundTrip passes through to defaultBase
// with no retry overhead.
func NewRoundTripper(endpoints []Endpoint, fo Options, defaultBase http.RoundTripper) http.RoundTripper {
	if defaultBase == nil {
		defaultBase = http.DefaultTransport
	}

	// Snapshot slices so external mutations do not affect this instance.
	methodsCopy := make([]string, len(fo.RetryableMethods))
	copy(methodsCopy, fo.RetryableMethods)
	fo.RetryableMethods = methodsCopy

	statusCodesCopy := make([]int, len(fo.FailoverOnStatusCodes))
	copy(statusCodesCopy, fo.FailoverOnStatusCodes)
	fo.FailoverOnStatusCodes = statusCodesCopy

	eps := make([]endpointRT, len(endpoints))
	for i, e := range endpoints {
		rt := defaultBase
		if e.SkipTLSVerify || e.CertificateAuthData != "" {
			rt = shared.CreateTransport(e.SkipTLSVerify, e.CertificateAuthData)
		}
		eps[i] = endpointRT{url: e.URL, rt: rt}
	}

	return &RoundTripper{
		endpoints:   eps,
		fo:          fo,
		defaultBase: defaultBase,
	}
}

// RoundTrip implements http.RoundTripper. It validates preconditions and
// dispatches to a strategy-specific method based on the configured
// Strategy:
//
//   - RoundRobin: delegates to orderedRoundTrip, which cycles through
//     servers sequentially.
//   - None / empty: passes through to the base
//     transport with no retry logic.
//   - Unknown strategy: returns an error immediately (configuration error).
//
// New strategies extend the switch in this method — each case builds its own
// server-selection function (or calls an entirely different method).
//
// This method is called by HTTPClient.Do() inside product-level callAPI. Each
// callAPI retry triggers a fresh RoundTrip invocation that cycles through all
// servers from the beginning.
func (t *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("nil request")
	}
	if t == nil {
		return nil, errors.New("nil RoundTripper")
	}
	if t.defaultBase == nil {
		// Be resilient if instantiated without constructor.
		t.defaultBase = http.DefaultTransport
	}

	numServers := len(t.endpoints)

	switch NormalizeStrategy(t.fo.Strategy) {
	case NormalizeStrategy(RoundRobin):
		if numServers == 0 {
			return nil, fmt.Errorf("failover strategy %q requires at least one server URL", t.fo.Strategy)
		}
		order := func(attempt int) int { return attempt % numServers }
		return t.orderedRoundTrip(req, order)
	case None, "": // Explicit "none" or empty strategy disables failover and passes through to base.
		return t.defaultBase.RoundTrip(req)
	default: // unknown strategy is a configuration error that we surface immediately.
		return nil, fmt.Errorf("unknown failover strategy: %s", t.fo.Strategy)
	}
}

// orderedRoundTrip is the shared retry loop for strategies that pick servers
// deterministically by attempt index. The order function maps each attempt
// (0, 1, 2, …) to a server index.
//
// Behaviour:
//   - Non-retryable HTTP methods (e.g. POST by default): passes through to the
//     base transport without retry.
//   - Network errors: retries only for an allowlist of transport errors
//     (connection refused/reset/aborted, unreachable host/network) and timeout
//     errors when Options.RetryOnTimeout is set. DNS errors and context
//     cancellation are never retried.
//   - Status codes in FailoverOnStatusCodes: drains the response body and
//     retries against the next server. Response headers (e.g. Retry-After) are
//     not inspected at this layer.
//   - All attempts exhausted: returns an error (not an HTTP response). This
//     means callAPI receives err != nil and returns immediately — its own
//     status-code-based retry logic (for 502/503/504/429) is never reached.
func (t *RoundTripper) orderedRoundTrip(req *http.Request, order serverOrder) (*http.Response, error) {
	fo := &t.fo

	if !isRetryableMethod(fo, req.Method) {
		return t.defaultBase.RoundTrip(req)
	}

	bo := fo.ExponentialBackoff.NewExponentialBackoff()

	maxRetries := fo.MaxRetries
	if maxRetries == 0 {
		maxRetries = shared.DefaultMaxRetries
	}
	// maxRetries in config means "number of retries", so total attempts = maxRetries + 1
	var lastErr error
	totalAttempts := maxRetries + 1
	for attempt := range totalAttempts {
		serverIndex := order(attempt)
		ep := t.endpoints[serverIndex]

		shared.LogDebug("[Failover] attempt=%d, serverIndex=%d, serverURL=%s, method=%s", attempt, serverIndex, ep.url, req.Method)

		resp, err := t.doFailoverAttempt(req, ep)
		if err != nil {
			if !isNetworkErrorRT(req.Context(), err, fo.RetryOnTimeout) {
				shared.LogDebug("[Failover] attempt=%d failed with non-retriable error on Servers[%d]: %v", attempt, serverIndex, err)
				return nil, err
			}

			shared.LogDebug("[Failover] attempt=%d failed with retriable error on Servers[%d]: %v", attempt, serverIndex, err)

			lastErr = err
			// Don't sleep if this was the last attempt
			if attempt < totalAttempts-1 {
				retry.BackOff(req.Context(), bo.NextBackOff())
			}
			continue
		}

		if !shouldFailoverOnStatus(fo, resp.StatusCode) {
			shared.LogDebug("[Failover] attempt=%d ends failover loop with status=%d on Servers[%d]=%s", attempt, resp.StatusCode, serverIndex, ep.url)
			return resp, nil
		}

		shared.LogDebug("[Failover] attempt=%d, status=%d triggers failover to next server", attempt, resp.StatusCode)
		if resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
		// Don't sleep if this was the last attempt
		if attempt < totalAttempts-1 {
			retry.BackOff(req.Context(), bo.NextBackOff())
		}
		lastErr = fmt.Errorf("failover status: %s", resp.Status)
	}

	return nil, lastErr
}

// doFailoverAttempt prepares and executes a single failover attempt against
// the given endpoint, using the endpoint's own transport.
func (t *RoundTripper) doFailoverAttempt(req *http.Request, ep endpointRT) (*http.Response, error) {
	targetURL, err := url.Parse(ep.url)
	if err != nil {
		return nil, fmt.Errorf("invalid server URL %q: %w", ep.url, err)
	}

	attemptReq, err := retry.CloneRequestForRetry(req)
	if err != nil {
		return nil, fmt.Errorf("failed to clone request for retry: %w", err)
	}

	attemptReq.URL.Scheme = targetURL.Scheme
	attemptReq.URL.Host = targetURL.Host
	attemptReq.Host = targetURL.Host

	shared.LogDebug("[Failover] method=%s url=%s", attemptReq.Method, attemptReq.URL.String())
	return ep.rt.RoundTrip(attemptReq)
}

func isRetryableMethod(fo *Options, method string) bool {
	m := strings.ToUpper(strings.TrimSpace(method))
	if fo == nil {
		return defaultRetryableMethods[m]
	}

	if len(fo.RetryableMethods) == 0 {
		return defaultRetryableMethods[m]
	}

	for _, v := range fo.RetryableMethods {
		if strings.ToUpper(strings.TrimSpace(v)) == m {
			return true
		}
	}
	return false
}

var defaultRetryableMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodPut:     true,
	http.MethodDelete:  true,
	http.MethodOptions: true,
}

func isNetworkErrorRT(ctx context.Context, err error, retryOnTimeout bool) bool {
	if err == nil {
		return false
	}

	// TLS certificate/handshake failures are deterministic configuration/protocol
	// issues, so fail fast and do not fail over to another endpoint.
	if isTLSFailFastError(err) {
		return false
	}

	// If caller context is done, stop failover immediately.
	if ctx != nil && ctx.Err() != nil {
		return false
	}

	// DNS errors are not retried by failover.
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return false
	}

	// Only allowlisted transport-level net.OpError values are eligible for failover.
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return isRetryableNetOpError(opErr, retryOnTimeout)
	}

	return false
}

func isRetryableNetOpError(opErr *net.OpError, retryOnTimeout bool) bool {
	if opErr == nil {
		return false
	}

	if opErr.Timeout() || errors.Is(opErr, syscall.ETIMEDOUT) {
		return retryOnTimeout
	}

	return errors.Is(opErr, syscall.ECONNREFUSED) ||
		errors.Is(opErr, syscall.ECONNRESET) ||
		errors.Is(opErr, syscall.ECONNABORTED) ||
		errors.Is(opErr, syscall.EHOSTUNREACH) ||
		errors.Is(opErr, syscall.ENETUNREACH)
}

func isTLSFailFastError(err error) bool {
	if err == nil {
		return false
	}

	var certVerifyErr *tls.CertificateVerificationError
	var unknownAuthorityErr x509.UnknownAuthorityError
	var hostnameErr x509.HostnameError
	var certInvalidErr x509.CertificateInvalidError
	var systemRootsErr x509.SystemRootsError
	var recordHeaderErr tls.RecordHeaderError
	var alertErr tls.AlertError

	switch {
	case errors.As(err, &certVerifyErr):
		return true
	case errors.As(err, &unknownAuthorityErr):
		return true
	case errors.As(err, &hostnameErr):
		return true
	case errors.As(err, &certInvalidErr):
		return true
	case errors.As(err, &systemRootsErr):
		return true
	case errors.As(err, &recordHeaderErr):
		return true
	case errors.As(err, &alertErr):
		return true
	default:
		return false
	}
}

func shouldFailoverOnStatus(fo *Options, statusCode int) bool {
	if fo == nil || len(fo.FailoverOnStatusCodes) == 0 {
		return false
	}
	for _, sc := range fo.FailoverOnStatusCodes {
		if sc == statusCode {
			return true
		}
	}
	return false
}

// serverOrder maps an attempt index (0, 1, 2, …) to a server index.
// Different strategies produce different orderings.
type serverOrder func(attempt int) int

// NormalizeStrategy returns the strategy in lower-case with surrounding
// whitespace removed, so comparisons are case-insensitive.
func NormalizeStrategy(s Strategy) Strategy {
	return Strategy(strings.TrimSpace(strings.ToLower(string(s))))
}

// Options controls transport-level endpoint failover behaviour.
// It is nested under Configuration so it can be grouped cleanly in JSON/YAML.
//
// This layer only applies when a multiserver strategy is active (currently
// only RoundRobin). With None/empty or a single server, the
// transport passes through directly to the base http.RoundTripper.
//
// # Interaction with product-level callAPI retry loop
//
// Product-level callAPI wraps HTTPClient.Do(), which invokes RoundTrip().
// Each callAPI retry triggers a fresh RoundTrip() call that cycles through
// all servers from the beginning. Worst-case total attempts:
//
//	callAPI.MaxRetries × Options.MaxRetries
//
// # FailoverOnStatusCodes behaviour
//
// Status codes listed in FailoverOnStatusCodes are handled at the transport
// level: the response body is drained, and the request is retried against the
// next server with exponential backoff. Response headers (e.g. Retry-After)
// are not inspected. If all servers return a listed code, RoundTrip returns an
// error (not an HTTP response), so callAPI receives err != nil and returns
// immediately — its own status-code retry logic is never reached.
//
// Status codes NOT listed here pass through as a normal HTTP response to
// callAPI, which has its own retry logic for 502/503/504 (fixed backoff, no
// POST retry) and 429 (honors Retry-After).
type Options struct {
	Strategy Strategy `json:"strategy,omitempty" yaml:"strategy,omitempty"`

	// RetryableMethods controls which HTTP methods are eligible for transport-level
	// failover retries. If empty/nil, the default is safe/idempotent methods.
	RetryableMethods []string `json:"retryableMethods,omitempty" yaml:"retryableMethods,omitempty"`

	// RetryOnTimeout controls whether failover retries should also happen when
	// the underlying transport experiences a timeout (for example, a network
	// timeout such as net.OpError.Timeout() / ETIMEDOUT). It does not apply to
	// context cancellation or deadline-exceeded errors (ctx.Err() != nil), which
	// always stop retries.
	RetryOnTimeout bool `json:"retryOnTimeout,omitempty" yaml:"retryOnTimeout,omitempty"`

	// FailoverOnStatusCodes controls whether the transport should fail over to the
	// next server when it receives one of these HTTP status codes.
	FailoverOnStatusCodes []int `json:"failoverOnStatusCodes,omitempty" yaml:"failoverOnStatusCodes,omitempty"`

	// MaxRetries controls how many times the transport will attempt a request using the set strategy before giving up
	// and returning the last error. If zero, it defaults to 3.
	MaxRetries int `json:"maxRetries,omitempty" yaml:"maxRetries,omitempty"`

	ExponentialBackoff *ExponentialBackoffOptions `json:"exponentialBackoff,omitempty" yaml:"exponentialBackoff,omitempty"`
}

// ExponentialBackoffOptions controls the backoff parameters for exponential backoff.
// By configuring Multiplier and RandomizationFactor, it is possible to achieve a constant backoff or a linear backoff as well.
type ExponentialBackoffOptions struct {
	// InitialInterval is the initial interval between retries. If zero, it defaults to 500ms.
	InitialInterval time.Duration `json:"initialInterval,omitempty" yaml:"initialInterval,omitempty"`

	// MaxInterval is the maximum interval between retries. If zero, it defaults to 60s.
	MaxInterval time.Duration `json:"maxInterval,omitempty" yaml:"maxInterval,omitempty"`

	// Multiplier is the factor by which the interval increases after each retry. If nil, it defaults to 1.5.
	Multiplier *float64 `json:"multiplier,omitempty" yaml:"multiplier,omitempty"`

	// RandomizationFactor is the factor used to randomize the backoff intervals. If nil, it defaults to 0.5.
	RandomizationFactor *float64 `json:"randomizationFactor,omitempty" yaml:"randomizationFactor,omitempty"`
}

// NewExponentialBackoff creates a new exponential backoff instance configured
// with the options in b. If b is nil, defaults from cenkalti/backoff are used.
func (b *ExponentialBackoffOptions) NewExponentialBackoff() boff.BackOff {
	bo := boff.NewExponentialBackOff()
	if b == nil {
		return bo
	}

	if b.InitialInterval != 0 {
		bo.InitialInterval = b.InitialInterval
	}

	if b.MaxInterval != 0 {
		bo.MaxInterval = b.MaxInterval
	}

	if b.Multiplier != nil {
		bo.Multiplier = *b.Multiplier
	}

	if b.RandomizationFactor != nil {
		bo.RandomizationFactor = *b.RandomizationFactor
	}

	return bo
}
