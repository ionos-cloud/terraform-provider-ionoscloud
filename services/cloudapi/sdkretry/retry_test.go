// Package sdkretry contains black-box tests for the Cloud API SDK client's
// HTTP retry behaviour, which the provider relies on. The tests drive the SDK
// through a public GET/POST operation against an httptest server so they do not
// depend on any generated internals and survive SDK regeneration.
package sdkretry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// newMockClient returns an APIClient whose requests are served by handler, with
// a tiny WaitTime so retry backoff does not slow the tests down. The SDK keeps
// the /cloudapi/v6 base path from its default server and only swaps scheme+host.
func newMockClient(t *testing.T, maxRetries int, handler http.HandlerFunc) *ionoscloud.APIClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	cfg := ionoscloud.NewConfiguration("", "", "", srv.URL)
	cfg.MaxRetries = maxRetries
	cfg.WaitTime = time.Millisecond
	cfg.MaxWaitTime = 10 * time.Millisecond
	return ionoscloud.NewAPIClient(cfg)
}

// TestRetryOn500ForGET covers the two GET outcomes: the client keeps retrying a
// 500 and succeeds if a later attempt does, and it gives up after MaxRetries
// attempts, returning the last 500 rather than retrying forever.
func TestRetryOn500ForGET(t *testing.T) {
	t.Run("retries then succeeds", func(t *testing.T) {
		var calls int32
		client := newMockClient(t, 3, func(w http.ResponseWriter, _ *http.Request) {
			if atomic.AddInt32(&calls, 1) < 3 { // fail twice, then succeed
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"dc-1","type":"collection","items":[]}`))
		})

		_, apiResp, err := client.DataCentersApi.DatacentersGet(context.Background()).Depth(0).Execute()
		if err != nil {
			t.Fatalf("expected success after retries, got error: %v", err)
		}
		if apiResp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", apiResp.StatusCode, http.StatusOK)
		}
		if got := atomic.LoadInt32(&calls); got != 3 {
			t.Fatalf("attempts = %d, want 3 (1 initial + 2 retries)", got)
		}
	})

	t.Run("exhausts retries and returns the last 500", func(t *testing.T) {
		const maxRetries = 3
		var calls int32
		client := newMockClient(t, maxRetries, func(w http.ResponseWriter, _ *http.Request) {
			atomic.AddInt32(&calls, 1)
			w.WriteHeader(http.StatusInternalServerError)
		})

		_, apiResp, err := client.DataCentersApi.DatacentersGet(context.Background()).Depth(0).Execute()
		if err == nil {
			t.Fatal("expected an error after retries are exhausted, got nil")
		}
		if apiResp == nil || apiResp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %v, want %d", apiResp, http.StatusInternalServerError)
		}
		if got := atomic.LoadInt32(&calls); got != maxRetries {
			t.Fatalf("attempts = %d, want %d (MaxRetries)", got, maxRetries)
		}
	})
}

// TestNoRetryOn500ForPOST guards the idempotency rule: a 500 on a non-GET
// (potentially non-idempotent) request must not be retried.
func TestNoRetryOn500ForPOST(t *testing.T) {
	var calls int32
	client := newMockClient(t, 3, func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusInternalServerError)
	})

	// The mock returns 500 regardless of the body, so an empty request is enough
	// to exercise the POST path.
	dc := ionoscloud.DatacenterPost{Properties: &ionoscloud.DatacenterPropertiesPost{}}
	_, _, err := client.DataCentersApi.DatacentersPost(context.Background()).Datacenter(dc).Execute()
	if err == nil {
		t.Fatal("expected an error from the 500 response, got nil")
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("attempts = %d, want 1 (POST 500 must not be retried)", got)
	}
}
