package diags

import (
	"errors"
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

func TestWrapError_WithRequestLocation(t *testing.T) {
	original := errors.New("something failed")
	loc := "https://api.ionos.com/cloudapi/v6/requests/a1b2c3d4-e5f6-7890-abcd-ef1234567890/status"
	got := WrapError(original, &ErrorContext{RequestLocation: loc})
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

func TestWrapError_AllFields(t *testing.T) {
	original := errors.New("base error")
	loc := "https://api.ionos.com/cloudapi/v6/requests/deadbeef-dead-beef-dead-beefdeadbeef/status"
	ctx := &ErrorContext{
		ResourceID:      "res-id-1",
		ResourceName:    "my-server",
		Timeout:         "60m",
		StatusCode:      502,
		RequestLocation: loc,
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

// Tests for extractRequestId

func TestExtractRequestId_ValidURL(t *testing.T) {
	loc := "https://api.ionos.com/cloudapi/v6/requests/a1b2c3d4-e5f6-7890-abcd-ef1234567890/status"
	got := extractRequestId(loc)
	if got != "a1b2c3d4-e5f6-7890-abcd-ef1234567890" {
		t.Fatalf("expected UUID, got %q", got)
	}
}

func TestExtractRequestId_NoMatch(t *testing.T) {
	got := extractRequestId("https://api.ionos.com/cloudapi/v6/datacenters/123")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestExtractRequestId_EmptyString(t *testing.T) {
	got := extractRequestId("")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}
