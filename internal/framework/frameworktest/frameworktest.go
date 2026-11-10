// Package frameworktest provides helpers for driving the plugin-framework provider's protocol RPCs
// against a stubbed regional API, without provisioning any real infrastructure. It is used by the
// DBaaS v2 resource tests that verify partial-state-on-failure behavior.
package frameworktest

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/stretchr/testify/require"

	fwprovider "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/provider"
)

// FailedProvisioningStub starts a fake regional DBaaS API whose cluster creation succeeds (POST
// returns an id and a PROVISIONING state) but whose cluster never becomes available (GET reports a
// FAILED state), and points the provider at it through a temporary IONOS file config. It models the
// "POST ok, WaitUntilAvailable fails" scenario shared by every DBaaS v2 cluster resource.
type FailedProvisioningStub struct {
	Product   string // fileconfiguration product name, e.g. fileconfiguration.PSQLV2
	Location  string // regional location the stub serves; must match the planned cluster's location
	ClusterID string // id the POST returns and the GET is keyed by
}

// Serve starts the stub and wires the provider to it. Call it before NewConfiguredServer.
func (s FailedProvisioningStub) Serve(t *testing.T) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/clusters"):
			w.WriteHeader(http.StatusCreated)
			writeJSON(t, w, clusterBody(s.ClusterID, "PROVISIONING"))
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/clusters/"):
			w.WriteHeader(http.StatusOK)
			writeJSON(t, w, clusterBody(s.ClusterID, "FAILED"))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	configPath := filepath.Join(t.TempDir(), "ionos-config.yaml")
	config := fmt.Sprintf(`version: 1.0
environments:
  - name: test
    products:
      - name: %s
        endpoints:
          - location: %s
            name: %s
`, s.Product, s.Location, server.URL)
	require.NoError(t, os.WriteFile(configPath, []byte(config), 0o600), "write file config")
	t.Setenv(shared.IonosFilePathEnvVar, configPath)
}

// clusterBody is the minimal ClusterRead JSON every DBaaS v2 product shares: an id and a state.
func clusterBody(id, state string) map[string]any {
	return map[string]any{"id": id, "metadata": map[string]any{"state": state}}
}

// NewConfiguredServer builds the provider server, fetches its schema and runs ConfigureProvider with
// an all-null config so credentials are read from the environment. A stub must be served first.
func NewConfiguredServer(ctx context.Context, t *testing.T) (tfprotov6.ProviderServer, *tfprotov6.GetProviderSchemaResponse) {
	t.Helper()

	t.Setenv("IONOS_TOKEN", "token-for-the-stub")

	server := providerserver.NewProtocol6(fwprovider.New())()

	schema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err, "GetProviderSchema")
	AssertNoErrorDiags(t, "GetProviderSchema", schema.Diagnostics)

	// An empty object leaves every provider attribute null (non-null object with null attributes).
	providerType := schema.Provider.ValueType()
	config, err := tfprotov6.NewDynamicValue(providerType, buildValue(t, providerType, map[string]any{}))
	require.NoError(t, err, "marshal provider config")

	resp, err := server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{
		TerraformVersion: "1.14.0",
		Config:           &config,
	})
	require.NoError(t, err, "ConfigureProvider")
	AssertNoErrorDiags(t, "ConfigureProvider", resp.Diagnostics)

	return server, schema
}

// ResourceType returns the tftypes.Type of a registered resource, failing if it is not registered.
func ResourceType(t *testing.T, schema *tfprotov6.GetProviderSchemaResponse, typeName string) tftypes.Type {
	t.Helper()

	rs, ok := schema.ResourceSchemas[typeName]
	require.Truef(t, ok, "resource %q is not registered", typeName)
	return rs.ValueType()
}

// Create drives ApplyResourceChange as a Create (null prior state) with the planned state parsed from
// plannedJSON (attributes omitted from the JSON default to null), and returns the response.
func Create(ctx context.Context, t *testing.T, server tfprotov6.ProviderServer, typeName string, resourceType tftypes.Type, plannedJSON string) *tfprotov6.ApplyResourceChangeResponse {
	t.Helper()

	var raw any
	require.NoError(t, json.Unmarshal([]byte(plannedJSON), &raw), "parse planned JSON")

	plannedValue, err := tfprotov6.NewDynamicValue(resourceType, buildValue(t, resourceType, raw))
	require.NoError(t, err, "marshal planned state")
	priorValue, err := tfprotov6.NewDynamicValue(resourceType, tftypes.NewValue(resourceType, nil))
	require.NoError(t, err, "marshal prior state")

	resp, err := server.ApplyResourceChange(ctx, &tfprotov6.ApplyResourceChangeRequest{
		TypeName:     typeName,
		PriorState:   &priorValue,
		PlannedState: &plannedValue,
		Config:       &plannedValue,
	})
	require.NoError(t, err, "ApplyResourceChange")
	return resp
}

// StateStrings reads the named string attributes from a state DynamicValue. Attributes that are null,
// unknown or non-string are omitted from the result.
func StateStrings(t *testing.T, state *tfprotov6.DynamicValue, resourceType tftypes.Type, names ...string) map[string]string {
	t.Helper()

	decoded, err := state.Unmarshal(resourceType)
	require.NoError(t, err, "unmarshal state")

	var attributes map[string]tftypes.Value
	require.NoError(t, decoded.As(&attributes), "read state object")

	out := make(map[string]string, len(names))
	for _, name := range names {
		attribute, ok := attributes[name]
		if !ok || attribute.IsNull() || !attribute.IsKnown() || !attribute.Type().Is(tftypes.String) {
			continue
		}
		var s string
		require.NoErrorf(t, attribute.As(&s), "read attribute %q", name)
		out[name] = s
	}
	return out
}

// HasErrorDiag reports whether any error diagnostic's summary or detail contains substr.
func HasErrorDiag(diags []*tfprotov6.Diagnostic, substr string) bool {
	for _, d := range diags {
		if d != nil && d.Severity == tfprotov6.DiagnosticSeverityError &&
			(strings.Contains(d.Summary, substr) || strings.Contains(d.Detail, substr)) {
			return true
		}
	}
	return false
}

// FormatDiags renders diagnostics for assertion failure messages.
func FormatDiags(diags []*tfprotov6.Diagnostic) string {
	var b strings.Builder
	for _, d := range diags {
		if d == nil {
			continue
		}
		fmt.Fprintf(&b, "\n  [%v] %s: %s", d.Severity, d.Summary, d.Detail)
	}
	return b.String()
}

// AssertNoErrorDiags fails the test if diags contains any error diagnostic.
func AssertNoErrorDiags(t *testing.T, rpc string, diags []*tfprotov6.Diagnostic) {
	t.Helper()

	for _, d := range diags {
		if d != nil && d.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatalf("%s returned an error diagnostic: %s: %s", rpc, d.Summary, d.Detail)
		}
	}
}

// buildValue converts a decoded-JSON Go value (from json.Unmarshal into any) into a tftypes.Value of
// the given type. Object attributes absent from the Go map default to null, so a planned state only
// needs to spell out the attributes that matter. It covers the types the cluster schemas use:
// objects, lists, strings, numbers and bools.
func buildValue(t *testing.T, typ tftypes.Type, raw any) tftypes.Value {
	t.Helper()

	if raw == nil {
		return tftypes.NewValue(typ, nil)
	}

	switch {
	case typ.Is(tftypes.String):
		return tftypes.NewValue(typ, raw.(string))
	case typ.Is(tftypes.Bool):
		return tftypes.NewValue(typ, raw.(bool))
	case typ.Is(tftypes.Number):
		return tftypes.NewValue(typ, big.NewFloat(raw.(float64)))
	case typ.Is(tftypes.List{}):
		elementType := typ.(tftypes.List).ElementType
		items := raw.([]any)
		elements := make([]tftypes.Value, len(items))
		for i, item := range items {
			elements[i] = buildValue(t, elementType, item)
		}
		return tftypes.NewValue(typ, elements)
	case typ.Is(tftypes.Object{}):
		object := typ.(tftypes.Object)
		provided, ok := raw.(map[string]any)
		require.Truef(t, ok, "expected an object for type %s, got %T", typ, raw)
		attributes := make(map[string]tftypes.Value, len(object.AttributeTypes))
		for name, attributeType := range object.AttributeTypes {
			attributes[name] = buildValue(t, attributeType, provided[name])
		}
		return tftypes.NewValue(typ, attributes)
	default:
		t.Fatalf("buildValue: unsupported type %s", typ)
		return tftypes.Value{}
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, body any) {
	t.Helper()

	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Errorf("failed to write the stubbed response: %v", err)
	}
}
