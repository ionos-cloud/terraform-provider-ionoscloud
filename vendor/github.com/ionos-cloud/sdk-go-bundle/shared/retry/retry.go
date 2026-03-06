// IONOS Shared Libraries – Application-level retry

package retry

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// DoWithApplicationRetry executes an HTTP request with application-level retry
// logic. This is the shared equivalent of the per-product callAPI method.
//
// Retry behaviour:
//   - 502/503/504: retries with backoff (skips POST)
//   - 429: retries honoring Retry-After header, falls back to waitTime
//   - Other status codes or transport errors: returns immediately
//   - Respects context cancellation during backoff
//
// # Interaction with RoundTripper
//
// client.Do invokes the configured http.RoundTripper, which may be a
// RoundTripper. When a multi-server Strategy is active,
// transport-level network errors and status codes listed in
// FailoverOnStatusCodes are handled by the round tripper before this
// function sees the response: it either receives a successful response or
// err != nil after all failover attempts are exhausted.
//
// Status codes NOT in FailoverOnStatusCodes (including the 502/503/504/429
// that this function retries on) pass through normally.
//
// Worst-case total attempts: Configuration.MaxRetries × Options.MaxRetries.
func DoWithApplicationRetry(
	cfg *shared.Configuration,
	request *http.Request,
) (*http.Response, time.Duration, error) {
	if cfg == nil {
		return nil, 0, errors.New("nil configuration")
	}

	var resp *http.Response
	var httpRequestTime time.Duration
	var err error

	for attempt := range cfg.MaxRetries {
		resp, httpRequestTime, err = doRetryAttempt(cfg, request)
		if err != nil {
			return resp, httpRequestTime, err
		}
		shared.SdkLogger.Printf("\n try no: %d\n", attempt)
		logResponse(resp)

		backoffTime, retryErr := retryBackoff(resp, request.Method, cfg.WaitTime)
		if retryErr != nil {
			return resp, httpRequestTime, retryErr
		}
		if backoffTime < 0 {
			return resp, httpRequestTime, err
		}

		if attempt == cfg.MaxRetries-1 {
			shared.LogDebug(" Number of maximum retries exceeded (%d retries)\n", cfg.MaxRetries)
			return resp, httpRequestTime, err
		}

		// Drain body before retrying so the connection can be reused.
		drainBody(resp)

		if backoffTime > cfg.MaxWaitTime {
			backoffTime = cfg.MaxWaitTime
		}
		BackOff(request.Context(), backoffTime)
	}

	return resp, httpRequestTime, err
}

func CloneRequestForRetry(req *http.Request) (*http.Request, error) {
	clone := req.Clone(req.Context())
	if req.Body != nil {
		if req.GetBody == nil {
			return nil, errors.New("request body is not replayable (GetBody is nil)")
		}
		b, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		clone.Body = b
	}
	return clone, nil
}

// doRetryAttempt clones the request, executes it, and returns the response.
func doRetryAttempt(cfg *shared.Configuration, request *http.Request) (*http.Response, time.Duration, error) {
	clonedRequest, cloneErr := CloneRequestForRetry(request)
	if cloneErr != nil {
		return nil, 0, cloneErr
	}

	logRequest(clonedRequest)

	httpRequestStartTime := time.Now()
	clonedRequest.Close = true
	resp, err := cfg.HTTPClient.Do(clonedRequest)
	httpRequestTime := time.Since(httpRequestStartTime)

	return resp, httpRequestTime, err
}

// retryBackoff determines the backoff duration based on the response status
// code. It returns a negative duration if the request should not be retried,
// or a non-nil error if parsing the Retry-After header fails on 429 responses.
func retryBackoff(resp *http.Response, method string, waitTime time.Duration) (time.Duration, error) {
	switch resp.StatusCode {
	case http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusBadGateway:
		if method == http.MethodPost {
			return -1, nil
		}
		return waitTime, nil

	case http.StatusTooManyRequests:
		return parseTooManyRequestsBackoff(resp, waitTime)

	default:
		return -1, nil
	}
}

// parseTooManyRequestsBackoff reads the Retry-After header to determine
// the backoff duration, falling back to waitTime if the header is absent.
func parseTooManyRequestsBackoff(resp *http.Response, waitTime time.Duration) (time.Duration, error) {
	retryAfterSeconds := resp.Header.Get("Retry-After")
	if retryAfterSeconds == "" {
		return waitTime, nil
	}

	retryWait, parseErr := time.ParseDuration(retryAfterSeconds + "s")
	if parseErr != nil {
		return -1, parseErr
	}

	return retryWait, nil
}

// drainBody discards and closes the response body so the underlying
// connection can be returned to the pool.
func drainBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

// logRequest dumps the outgoing request at Debug level.
// The Authorization header is stripped unless Trace is enabled.
func logRequest(request *http.Request) {
	if !shared.SdkLogLevel.Satisfies(shared.Debug) {
		return
	}
	logReq := request.Clone(request.Context())
	if !shared.SdkLogLevel.Satisfies(shared.Trace) {
		logReq.Header.Del("Authorization")
	}
	dump, err := httputil.DumpRequestOut(logReq, true)
	if err == nil {
		shared.SdkLogger.Printf(" DumpRequestOut : %s\n", string(dump))
	} else {
		shared.SdkLogger.Printf(" DumpRequestOut err: %+v", err)
	}
}

// logResponse dumps the server response at Debug level.
// The response body is only included at Trace level to avoid leaking
// sensitive data.
func logResponse(resp *http.Response) {
	if !shared.SdkLogLevel.Satisfies(shared.Debug) {
		return
	}
	dumpBody := shared.SdkLogLevel.Satisfies(shared.Trace)
	dump, err := httputil.DumpResponse(resp, dumpBody)
	if err == nil {
		shared.SdkLogger.Printf("\n DumpResponse : %s\n", string(dump))
	} else {
		shared.SdkLogger.Printf(" DumpResponse err %+v", err)
	}
}

// BackOff sleeps for the given duration and respects context cancellation.
func BackOff(ctx context.Context, t time.Duration) {
	shared.LogDebug(" Sleeping %s before retrying request\n", t.String())
	if t <= 0 {
		return
	}
	timer := time.NewTimer(t)
	defer timer.Stop()

	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}
