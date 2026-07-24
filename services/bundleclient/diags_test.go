package bundleclient_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

// bundleMeta builds an SdkBundle carrying its own contract-number Enricher and
// returns it BY VALUE — mirroring the SDKv2 provider, whose ConfigureContextFunc
// ends in `return *client, nil` (ionoscloud/provider.go:342). The dynamic type
// of `meta` handed to every CRUD func is therefore an SdkBundle value, which is
// what bundleclient.ToDiags/ToError type-assert against.
func bundleMeta(t *testing.T, contractNumber string) bundleclient.SdkBundle {
	t.Helper()
	b := bundleclient.New(context.Background(), clientoptions.TerraformClientOptions{}, nil)
	b.Diags = diagutil.NewEnricher(contractNumber, "", nil)
	return *b
}

// emptyResourceData returns a throwaway *schema.ResourceData. A fresh one is
// used per goroutine in the concurrency test since ResourceData is not safe for
// concurrent use.
func emptyResourceData() *schema.ResourceData {
	res := &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
	return res.TestResourceData()
}

// firstSummary returns the Summary of the first diagnostic, or "" when there
// are none. diag.FromErr places the full (enriched) error string in Summary.
func firstSummary(diags diag.Diagnostics) string {
	if len(diags) == 0 {
		return ""
	}
	return diags[0].Summary
}

func mustSummary(t *testing.T, diags diag.Diagnostics) string {
	t.Helper()
	if len(diags) == 0 {
		t.Fatalf("expected at least one diagnostic, got none")
	}
	return diags[0].Summary
}

// TestToDiags_PerConfigContractNumber is the anti-contamination guarantee at the
// SDKv2 entry point that the upjet no-fork path actually uses (ts.Meta →
// bundleclient.ToDiags(meta, ...)). Two provider configurations enrich errors
// with their OWN contract number; neither can observe the other's. This is the
// exact case the old package-level resolver got wrong.
func TestToDiags_PerConfigContractNumber(t *testing.T) {
	metaA := bundleMeta(t, "111")
	metaB := bundleMeta(t, "222")

	msgA := mustSummary(t, bundleclient.ToDiags(metaA, emptyResourceData(), errors.New("boom"), &diagutil.ErrorContext{ResourceID: "dc-a"}))
	msgB := mustSummary(t, bundleclient.ToDiags(metaB, emptyResourceData(), errors.New("boom"), &diagutil.ErrorContext{ResourceID: "dc-b"}))

	if !strings.Contains(msgA, "Contract number: 111") {
		t.Errorf("meta A should report contract 111, got %q", msgA)
	}
	if strings.Contains(msgA, "Contract number: 222") {
		t.Errorf("meta A leaked meta B's contract number (contamination): %q", msgA)
	}
	if !strings.Contains(msgB, "Contract number: 222") {
		t.Errorf("meta B should report contract 222, got %q", msgB)
	}
	if strings.Contains(msgB, "Contract number: 111") {
		t.Errorf("meta B leaked meta A's contract number (contamination): %q", msgB)
	}
}

// TestToDiags_ConcurrentNoContamination reproduces the live upjet scenario:
// many resources from two ProviderConfigs enriching errors at the same time.
// With the per-config Enricher there is no shared mutable state, so every
// goroutine must see only its own contract number. Run with -race to also prove
// there is no data race on the resolution path.
func TestToDiags_ConcurrentNoContamination(t *testing.T) {
	metaA := bundleMeta(t, "111")
	metaB := bundleMeta(t, "222")

	const n = 100
	var wg sync.WaitGroup
	failures := make(chan string, 2*n)

	for range n {
		wg.Add(2)
		go func() {
			defer wg.Done()
			msg := firstSummary(bundleclient.ToDiags(metaA, emptyResourceData(), errors.New("x"), &diagutil.ErrorContext{ResourceID: "a"}))
			if !strings.Contains(msg, "Contract number: 111") || strings.Contains(msg, "Contract number: 222") {
				failures <- fmt.Sprintf("config A saw wrong contract number: %q", msg)
			}
		}()
		go func() {
			defer wg.Done()
			msg := firstSummary(bundleclient.ToDiags(metaB, emptyResourceData(), errors.New("x"), &diagutil.ErrorContext{ResourceID: "b"}))
			if !strings.Contains(msg, "Contract number: 222") || strings.Contains(msg, "Contract number: 111") {
				failures <- fmt.Sprintf("config B saw wrong contract number: %q", msg)
			}
		}()
	}

	wg.Wait()
	close(failures)
	for msg := range failures {
		t.Error(msg)
	}
}

// TestToDiags_NilError: a nil error must yield nil diagnostics, regardless of
// the configured contract number.
func TestToDiags_NilError(t *testing.T) {
	meta := bundleMeta(t, "111")
	if got := bundleclient.ToDiags(meta, emptyResourceData(), nil, nil); got != nil {
		t.Fatalf("expected nil diagnostics for nil error, got %v", got)
	}
}

// TestToDiags_WithoutEnricher_Safe covers the metas that carry no Enricher: a
// nil meta, a non-SdkBundle meta, and an SdkBundle whose Diags was never set
// (e.g. a bundle built outside a Configure that wired the contract number). All
// must degrade to plain error wrapping — no panic, no contract number.
func TestToDiags_WithoutEnricher_Safe(t *testing.T) {
	bundleNoDiags := bundleclient.New(context.Background(), clientoptions.TerraformClientOptions{}, nil)

	cases := []struct {
		name string
		meta any
	}{
		{"nil meta", nil},
		{"non-bundle meta", "not a bundle"},
		{"bundle with nil Diags", *bundleNoDiags},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := mustSummary(t, bundleclient.ToDiags(tc.meta, emptyResourceData(), errors.New("boom"), &diagutil.ErrorContext{ResourceID: "x"}))
			if !strings.Contains(msg, "boom") {
				t.Errorf("original error should be preserved, got %q", msg)
			}
			if !strings.Contains(msg, "Resource ID: x") {
				t.Errorf("resource context should still be applied, got %q", msg)
			}
			if strings.Contains(msg, "Contract number") {
				t.Errorf("no contract number expected without an enricher, got %q", msg)
			}
		})
	}
}

// TestToError_PerConfigContractNumber is the ToError counterpart of the
// anti-contamination test: the same scoping guarantee for the error-returning
// helper, plus nil-error passthrough.
func TestToError_PerConfigContractNumber(t *testing.T) {
	errA := bundleclient.ToError(bundleMeta(t, "111"), emptyResourceData(), errors.New("boom"), &diagutil.ErrorContext{ResourceID: "a"})
	if errA == nil || !strings.Contains(errA.Error(), "Contract number: 111") {
		t.Fatalf("expected contract 111 in wrapped error, got %v", errA)
	}
	if strings.Contains(errA.Error(), "Contract number: 222") {
		t.Fatalf("wrapped error leaked another config's contract number: %q", errA.Error())
	}

	if got := bundleclient.ToError(bundleMeta(t, "111"), emptyResourceData(), nil, nil); got != nil {
		t.Fatalf("expected nil error for nil input, got %v", got)
	}
}
