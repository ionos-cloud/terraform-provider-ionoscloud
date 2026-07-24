package diags

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestWrapError_NilError(t *testing.T) {
	ctx := &ErrorContext{ResourceID: "some-id"}
	got := (*Enricher)(nil).WrapError(nil, ctx)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestWrapError_NilContext(t *testing.T) {
	original := errors.New("something failed")
	got := (*Enricher)(nil).WrapError(original, nil)
	if got.Error() != original.Error() {
		t.Fatalf("expected %q, got %q", original.Error(), got.Error())
	}
}

func TestWrapError_EmptyContext(t *testing.T) {
	original := errors.New("something failed")
	got := (*Enricher)(nil).WrapError(original, &ErrorContext{})
	if got.Error() != original.Error() {
		t.Fatalf("expected %q, got %q", original.Error(), got.Error())
	}
}

func TestWrapError_WithResourceID(t *testing.T) {
	original := errors.New("something failed")
	got := (*Enricher)(nil).WrapError(original, &ErrorContext{ResourceID: "abc-123"})
	if !strings.Contains(got.Error(), "Resource ID: abc-123") {
		t.Fatalf("expected output to contain 'Resource ID: abc-123', got %q", got.Error())
	}
}

func TestWrapError_WithResourceName(t *testing.T) {
	original := errors.New("something failed")
	got := (*Enricher)(nil).WrapError(original, &ErrorContext{ResourceName: "my-resource"})
	if !strings.Contains(got.Error(), "Resource Name: my-resource") {
		t.Fatalf("expected output to contain 'Resource Name: my-resource', got %q", got.Error())
	}
}

func TestWrapError_WithTimeout(t *testing.T) {
	original := errors.New("something failed")
	got := (*Enricher)(nil).WrapError(original, &ErrorContext{Timeout: "30m"})
	if !strings.Contains(got.Error(), "Configured timeout: 30m") {
		t.Fatalf("expected output to contain 'Configured timeout: 30m', got %q", got.Error())
	}
}

func TestWrapError_WithRequestID(t *testing.T) {
	original := errors.New("something failed")
	got := (*Enricher)(nil).WrapError(original, &ErrorContext{RequestID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"})
	if !strings.Contains(got.Error(), "Request ID: a1b2c3d4-e5f6-7890-abcd-ef1234567890") {
		t.Fatalf("expected output to contain request ID, got %q", got.Error())
	}
}

func TestWrapError_WithStatusCode500(t *testing.T) {
	original := errors.New("internal server error")
	got := (*Enricher)(nil).WrapError(original, &ErrorContext{StatusCode: 500})
	if !strings.Contains(got.Error(), "This is an API Error") {
		t.Fatalf("expected output to contain API error message, got %q", got.Error())
	}
}

func TestWrapError_WithAdditionalInfo(t *testing.T) {
	original := errors.New("something failed")
	got := (*Enricher)(nil).WrapError(original, &ErrorContext{
		AdditionalInfo: map[string]string{
			"Cluster ID":    "cluster-abc",
			"Datacenter ID": "dc-123",
		},
	})
	msg := got.Error()
	if !strings.Contains(msg, "Cluster ID: cluster-abc") {
		t.Fatalf("expected output to contain 'Cluster ID: cluster-abc', got %q", msg)
	}
	if !strings.Contains(msg, "Datacenter ID: dc-123") {
		t.Fatalf("expected output to contain 'Datacenter ID: dc-123', got %q", msg)
	}
	// Verify sorted order: Cluster ID before Datacenter ID
	clusterIdx := strings.Index(msg, "Cluster ID")
	dcIdx := strings.Index(msg, "Datacenter ID")
	if clusterIdx > dcIdx {
		t.Fatalf("expected AdditionalInfo keys in sorted order, got Cluster ID at %d, Datacenter ID at %d", clusterIdx, dcIdx)
	}
}

func TestWrapError_AllFields(t *testing.T) {
	original := errors.New("base error")
	ctx := &ErrorContext{
		ResourceID:     "res-id-1",
		ResourceName:   "my-server",
		Timeout:        "60m",
		StatusCode:     502,
		RequestID:      "deadbeef-dead-beef-dead-beefdeadbeef",
		AdditionalInfo: map[string]string{"Cluster ID": "cl-1"},
	}
	got := (*Enricher)(nil).WrapError(original, ctx)
	msg := got.Error()

	for _, want := range []string{
		"base error",
		"Resource ID: res-id-1",
		"Resource Name: my-server",
		"Configured timeout: 60m",
		"Request ID: deadbeef-dead-beef-dead-beefdeadbeef",
		"This is an API Error",
		"Cluster ID: cl-1",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("expected output to contain %q, got %q", want, msg)
		}
	}
}

func TestWrapError_PreservesOriginalError(t *testing.T) {
	original := errors.New("root cause")
	got := (*Enricher)(nil).WrapError(original, &ErrorContext{ResourceID: "id-1"})
	if !errors.Is(got, original) {
		t.Fatalf("expected errors.Is to match the original error")
	}
}

func TestBuildContextString_NilContext(t *testing.T) {
	got := buildContextString(nil, "")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

// Tests for the per-configuration Enricher (replacement for the package-level
// contract number resolver).

// TestEnricher_Independence is the core anti-contamination guarantee: two
// configurations no longer share state, so configuring b cannot change what a
// reports. This is exactly the case the package-level resolver got wrong.
func TestEnricher_Independence(t *testing.T) {
	a := NewEnricher("111", "", nil)
	b := NewEnricher("222", "", nil)

	gotA := a.WrapError(errors.New("x"), &ErrorContext{ResourceID: "a"}).Error()
	gotB := b.WrapError(errors.New("x"), &ErrorContext{ResourceID: "b"}).Error()

	if !strings.Contains(gotA, "Contract number: 111") {
		t.Fatalf("enricher a should report 111, got %q", gotA)
	}
	if !strings.Contains(gotB, "Contract number: 222") {
		t.Fatalf("enricher b should report 222, got %q", gotB)
	}
}

// TestEnricher_SharedAcrossValueCopy mirrors how the bundle is delivered:
// SDKv2 returns the bundle by value (`return *client`) and the per-product
// constructors have value receivers that copy c.Diags. Because Diags is a
// pointer, every copy must observe the same instance and the same sync.Once.
func TestEnricher_SharedAcrossValueCopy(t *testing.T) {
	type bundle struct{ Diags *Enricher } // stands in for SdkBundle

	orig := bundle{Diags: NewEnricher("", "", func() string { return "55555" })}
	cp := orig // value copy

	if orig.Diags != cp.Diags {
		t.Fatalf("value copy must share the same *Enricher instance")
	}
	// Resolve through the copy; it must be visible through the original.
	_ = cp.Diags.ContractNumber()
	if got := orig.Diags.ContractNumber(); got != "55555" {
		t.Fatalf("expected shared resolved value 55555, got %q", got)
	}
}

// TestEnricher_LazyResolvedOnceConcurrent proves the lazy API fallback is
// resolved exactly once per instance, race-free under concurrency (run with
// -race). This is the data race the shared `called` bool in the old resolver
// could not avoid.
func TestEnricher_LazyResolvedOnceConcurrent(t *testing.T) {
	var calls int32
	e := NewEnricher("", "", func() string {
		atomic.AddInt32(&calls, 1)
		return "98765"
	})

	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			if got := e.ContractNumber(); got != "98765" {
				t.Errorf("expected 98765, got %q", got)
			}
		}()
	}
	wg.Wait()

	if c := atomic.LoadInt32(&calls); c != 1 {
		t.Fatalf("expected fallback resolved exactly once, got %d", c)
	}
}

// TestEnricher_NilSafe covers an unconfigured bundle (Diags still nil): the
// methods must degrade to plain wrapping without a contract number.
func TestEnricher_NilSafe(t *testing.T) {
	var e *Enricher

	if got := e.ContractNumber(); got != "" {
		t.Fatalf("nil enricher should return empty contract number, got %q", got)
	}

	original := errors.New("boom")
	got := e.WrapError(original, &ErrorContext{ResourceID: "x"})
	if !errors.Is(got, original) {
		t.Fatalf("expected wrapped error to match original")
	}
	if strings.Contains(got.Error(), "Contract number") {
		t.Fatalf("nil enricher should not add a contract number, got %q", got.Error())
	}
}

// Tests for ExtractRequestID

func TestExtractRequestID_ValidURL(t *testing.T) {
	u, _ := url.Parse("https://api.ionos.com/cloudapi/v6/requests/a1b2c3d4-e5f6-7890-abcd-ef1234567890/status")
	got := ExtractRequestID(u)
	if got != "a1b2c3d4-e5f6-7890-abcd-ef1234567890" {
		t.Fatalf("expected UUID, got %q", got)
	}
}

func TestExtractRequestID_NilURL(t *testing.T) {
	got := ExtractRequestID(nil)
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestExtractRequestID_NoMatch(t *testing.T) {
	u, _ := url.Parse("https://api.ionos.com/cloudapi/v6/datacenters/123")
	got := ExtractRequestID(u)
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

// Tests for ContractNumberFromToken

func buildTestJWT(payload string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256"}`))
	body := base64.RawURLEncoding.EncodeToString([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString([]byte("fakesignature"))
	return fmt.Sprintf("%s.%s.%s", header, body, sig)
}

func TestContractNumberFromToken_ValidToken(t *testing.T) {
	got := ContractNumberFromToken(buildTestJWT(`{"identity":{"contractNumber":31884391}}`))
	if got != "31884391" {
		t.Fatalf("expected '31884391', got %q", got)
	}
}

func TestContractNumberFromToken_StringContractNumber(t *testing.T) {
	got := ContractNumberFromToken(buildTestJWT(`{"identity":{"contractNumber":"12345"}}`))
	if got != "12345" {
		t.Fatalf("expected '12345', got %q", got)
	}
}

func TestContractNumberFromToken_MissingIdentity(t *testing.T) {
	got := ContractNumberFromToken(buildTestJWT(`{"iss":"ionoscloud"}`))
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestContractNumberFromToken_MissingContractNumber(t *testing.T) {
	got := ContractNumberFromToken(buildTestJWT(`{"identity":{"role":"owner"}}`))
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestContractNumberFromToken_InvalidJWT(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{"empty string", ""},
		{"no dots", "notajwt"},
		{"one part", "header.payload"},
		{"invalid base64", "header.!!!invalid!!!.signature"},
		{"invalid json", fmt.Sprintf("header.%s.signature", base64.RawURLEncoding.EncodeToString([]byte("not json")))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContractNumberFromToken(tt.token)
			if got != "" {
				t.Fatalf("expected empty string, got %q", got)
			}
		})
	}
}
