package diags

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestWrapError_NilError(t *testing.T) {
	ctx := &ErrorContext{ResourceID: "some-id"}
	got := WrapError(nil, ctx)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestWrapError_NilContext(t *testing.T) {
	original := errors.New("something failed")
	got := WrapError(original, nil)
	if got.Error() != original.Error() {
		t.Fatalf("expected %q, got %q", original.Error(), got.Error())
	}
}

func TestWrapError_EmptyContext(t *testing.T) {
	original := errors.New("something failed")
	got := WrapError(original, &ErrorContext{})
	if got.Error() != original.Error() {
		t.Fatalf("expected %q, got %q", original.Error(), got.Error())
	}
}

func TestWrapError_WithResourceID(t *testing.T) {
	original := errors.New("something failed")
	got := WrapError(original, &ErrorContext{ResourceID: "abc-123"})
	if !strings.Contains(got.Error(), "Resource ID: abc-123") {
		t.Fatalf("expected output to contain 'Resource ID: abc-123', got %q", got.Error())
	}
}

func TestWrapError_WithResourceName(t *testing.T) {
	original := errors.New("something failed")
	got := WrapError(original, &ErrorContext{ResourceName: "my-resource"})
	if !strings.Contains(got.Error(), "Resource Name: my-resource") {
		t.Fatalf("expected output to contain 'Resource Name: my-resource', got %q", got.Error())
	}
}

func TestWrapError_WithTimeout(t *testing.T) {
	original := errors.New("something failed")
	got := WrapError(original, &ErrorContext{Timeout: "30m"})
	if !strings.Contains(got.Error(), "Configured timeout: 30m") {
		t.Fatalf("expected output to contain 'Configured timeout: 30m', got %q", got.Error())
	}
}

func TestWrapError_WithRequestID(t *testing.T) {
	original := errors.New("something failed")
	got := WrapError(original, &ErrorContext{RequestID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"})
	if !strings.Contains(got.Error(), "Request ID: a1b2c3d4-e5f6-7890-abcd-ef1234567890") {
		t.Fatalf("expected output to contain request ID, got %q", got.Error())
	}
}

func TestWrapError_WithStatusCode500(t *testing.T) {
	original := errors.New("internal server error")
	got := WrapError(original, &ErrorContext{StatusCode: 500})
	if !strings.Contains(got.Error(), "This is an API Error") {
		t.Fatalf("expected output to contain API error message, got %q", got.Error())
	}
}

func TestWrapError_WithAdditionalInfo(t *testing.T) {
	original := errors.New("something failed")
	got := WrapError(original, &ErrorContext{
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
	got := WrapError(original, ctx)
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
	got := WrapError(original, &ErrorContext{ResourceID: "id-1"})
	if !errors.Is(got, original) {
		t.Fatalf("expected errors.Is to match the original error")
	}
}

func TestBuildContextString_NilContext(t *testing.T) {
	got := buildContextString(nil)
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

// Tests for contract number resolver

func TestSetupContractNumberResolver_ExplicitContractNumber(t *testing.T) {
	SetupContractNumberResolver("12345678", "", func() string { return "should-not-be-called" })
	original := errors.New("something failed")
	got := WrapError(original, &ErrorContext{ResourceID: "abc"})
	if !strings.Contains(got.Error(), "Contract number: 12345678") {
		t.Fatalf("expected output to contain 'Contract number: 12345678', got %q", got.Error())
	}
}

func TestSetupContractNumberResolver_FromToken(t *testing.T) {
	token := buildTestJWT(`{"identity":{"contractNumber":99999}}`)
	SetupContractNumberResolver("", token, func() string { return "should-not-be-called" })
	original := errors.New("something failed")
	got := WrapError(original, &ErrorContext{ResourceID: "abc"})
	if !strings.Contains(got.Error(), "Contract number: 99999") {
		t.Fatalf("expected output to contain 'Contract number: 99999', got %q", got.Error())
	}
}

func TestSetupContractNumberResolver_APIFallback(t *testing.T) {
	calls := 0
	SetupContractNumberResolver("", "", func() string {
		calls++
		return "77777"
	})
	WrapError(errors.New("err1"), &ErrorContext{ResourceID: "a"})
	WrapError(errors.New("err2"), &ErrorContext{ResourceID: "b"})
	if calls != 1 {
		t.Fatalf("expected API fallback to be called once, got %d", calls)
	}
}

func TestSetupContractNumberResolver_NoResolver(t *testing.T) {
	contractNumberFunc = nil
	original := errors.New("something failed")
	got := WrapError(original, &ErrorContext{ResourceID: "abc"})
	if strings.Contains(got.Error(), "Contract number") {
		t.Fatalf("expected no contract number, got %q", got.Error())
	}
}

func TestSetupContractNumberResolver_EmptyFallback(t *testing.T) {
	SetupContractNumberResolver("", "", func() string { return "" })
	original := errors.New("something failed")
	got := WrapError(original, &ErrorContext{ResourceID: "abc"})
	if strings.Contains(got.Error(), "Contract number") {
		t.Fatalf("expected no contract number, got %q", got.Error())
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
