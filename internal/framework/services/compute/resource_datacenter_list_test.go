package compute_test

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	ionoscloudsdk "github.com/ionos-cloud/sdk-go/v6"
	"github.com/stretchr/testify/assert"

	fwprovider "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/provider"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/ionoscloud"
)

const datacenterListType = "ionoscloud_datacenter"

// TestDatacenterListResource drives the ListResource RPC end to end against a stubbed
// Cloud API, through the same muxed provider server that main.go serves.
//
// It covers the parts of a list resource for an SDKv2 managed resource that can only
// fail at runtime: that the mux is happy with the list resource and the managed
// resource coming from different servers, that the framework registers a list resource
// with no framework resource behind it, that the protocol schemas handed over by
// RawV6Schemas convert cleanly, and that the resource model fills the SDKv2 schema
// without a type mismatch.
func TestDatacenterListResource(t *testing.T) {
	ctx := context.Background()

	t.Setenv("IONOS_API_URL", stubCloudAPI(t))
	t.Setenv("IONOS_TOKEN", "token-for-the-stub")

	server := muxedProviderServer(ctx, t)

	providerSchema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema: %v", err)
	}
	failOnErrorDiagnostics(t, "GetProviderSchema", providerSchema.Diagnostics)

	// The framework only registers a list resource that has no framework resource
	// behind it once RawV6Schemas has supplied both protocol schemas.
	if _, ok := providerSchema.ListResourceSchemas[datacenterListType]; !ok {
		t.Fatalf("the datacenter list resource was not registered; check RawV6Schemas and the SDKv2 resource identity")
	}
	if _, ok := providerSchema.ResourceSchemas[datacenterListType]; !ok {
		t.Fatalf("the datacenter managed resource is missing from the merged schema")
	}

	identitySchemas, err := server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		t.Fatalf("GetResourceIdentitySchemas: %v", err)
	}
	failOnErrorDiagnostics(t, "GetResourceIdentitySchemas", identitySchemas.Diagnostics)

	identitySchema := identitySchemas.IdentitySchemas[datacenterListType]
	if identitySchema == nil {
		t.Fatalf("the SDKv2 datacenter resource does not declare a resource identity")
	}

	configureProvider(ctx, t, server, providerSchema.Provider)

	listSchema := providerSchema.ListResourceSchemas[datacenterListType]
	resourceType := providerSchema.ResourceSchemas[datacenterListType].ValueType()

	t.Run("streams every datacenter", func(t *testing.T) {
		results := listDatacenters(ctx, t, server, listSchema, nil)
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}

		assert.Equal(t, "prod", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000001", identity["id"])
		assert.Equal(t, "de/txl", identity["location"])

		// Every attribute the mapper fills is asserted here, so that mapping a value
		// to the wrong attribute fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000001", resource["id"])
		assert.Equal(t, "prod", resource["name"])
		assert.Equal(t, "de/txl", resource["location"])
		assert.Equal(t, "the production datacenter", resource["description"])
		assert.Equal(t, "2001:db8::/56", resource["ipv6_cidr_block"])
		assert.Equal(t, false, resource["sec_auth_protection"])
		assert.Equal(t, float64(7), resource["version"])
		assert.ElementsMatch(t, []any{"SSD", "MULTIPLE_CPU"}, resource["features"])
		assert.Equal(t, []any{map[string]any{
			"cpu_family": "INTEL_SKYLAKE",
			"max_cores":  float64(32),
			"max_ram":    float64(245760),
			"vendor":     "GenuineIntel",
		}}, resource["cpu_architecture"])
		assert.Nil(t, resource["timeouts"], "a listed datacenter has no timeouts")

		// The second datacenter reports only the properties the API always sets, which
		// pins that the pairing holds past the first result and that the properties the
		// API left out stay null instead of turning into zero values.
		assert.Equal(t, "staging", results[1].DisplayName)

		stagingIdentity := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000002", stagingIdentity["id"])
		assert.Equal(t, "de/fra", stagingIdentity["location"])

		staging := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000002", staging["id"])
		assert.Equal(t, "staging", staging["name"])
		assert.Equal(t, "de/fra", staging["location"])
		assert.Equal(t, float64(1), staging["version"])
		assert.Nil(t, staging["description"])
		assert.Nil(t, staging["ipv6_cidr_block"])
		assert.Nil(t, staging["sec_auth_protection"])
		assert.Nil(t, staging["features"])
		assert.Nil(t, staging["cpu_architecture"])
	})

	t.Run("applies filters", func(t *testing.T) {
		results := listDatacenters(ctx, t, server, listSchema, map[string]string{"location": "de/fra"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "staging", results[0].DisplayName)
	})

	t.Run("rejects unknown filter fields", func(t *testing.T) {
		listServer, config := listServerAndConfig(t, server, listSchema, map[string]string{"nope": "value"})

		resp, err := listServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
			TypeName: datacenterListType,
			Config:   &config,
		})
		if err != nil {
			t.Fatalf("ValidateListResourceConfig: %v", err)
		}
		if !hasErrorDiagnostic(resp.Diagnostics) {
			t.Fatalf("expected a validation error for an unknown filter field")
		}
	})
}

// muxedProviderServer builds the same mux of the framework and SDKv2 providers that
// main.go serves.
func muxedProviderServer(ctx context.Context, t *testing.T) tfprotov6.ProviderServer {
	t.Helper()

	sdkv2Provider := ionoscloud.Provider()
	upgraded, err := tf5to6server.UpgradeServer(ctx, sdkv2Provider.GRPCProvider)
	if err != nil {
		t.Fatalf("failed to upgrade the SDKv2 provider server: %v", err)
	}

	mux, err := tf6muxserver.NewMuxServer(ctx,
		providerserver.NewProtocol6(fwprovider.New(sdkv2Provider)),
		func() tfprotov6.ProviderServer { return upgraded },
	)
	if err != nil {
		t.Fatalf("failed to build the mux server: %v", err)
	}

	return mux.ProviderServer()
}

// stubCloudAPI serves the datacenter collection the list resource reads, and returns
// the URL to point IONOS_API_URL at.
func stubCloudAPI(t *testing.T) string {
	t.Helper()

	datacenters := ionoscloudsdk.Datacenters{
		Items: &[]ionoscloudsdk.Datacenter{
			{
				Id: new("d3b07384-d9a0-4d1e-8f2a-000000000001"),
				Properties: &ionoscloudsdk.DatacenterProperties{
					Name:              new("prod"),
					Description:       new("the production datacenter"),
					Location:          new("de/txl"),
					Version:           new(int32(7)),
					Features:          &[]string{"SSD", "MULTIPLE_CPU"},
					SecAuthProtection: new(false),
					Ipv6CidrBlock:     new("2001:db8::/56"),
					CpuArchitecture: &[]ionoscloudsdk.CpuArchitectureProperties{{
						CpuFamily: new("INTEL_SKYLAKE"),
						MaxCores:  new(int32(32)),
						MaxRam:    new(int32(245760)),
						Vendor:    new("GenuineIntel"),
					}},
				},
			},
			{
				Id: new("d3b07384-d9a0-4d1e-8f2a-000000000002"),
				Properties: &ionoscloudsdk.DatacenterProperties{
					Name:     new("staging"),
					Location: new("de/fra"),
					Version:  new(int32(1)),
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/datacenters") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(datacenters); err != nil {
			t.Errorf("failed to write the stubbed response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return server.URL
}

// configureProvider runs ConfigureProvider with an all-null config, which leaves the
// provider reading its credentials and endpoint from the environment.
func configureProvider(ctx context.Context, t *testing.T, server tfprotov6.ProviderServer, schema *tfprotov6.Schema) {
	t.Helper()

	config := dynamicValue(t, schema.ValueType(), nullObject(t, schema.ValueType()))

	resp, err := server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{
		TerraformVersion: "1.14.0",
		Config:           &config,
	})
	if err != nil {
		t.Fatalf("ConfigureProvider: %v", err)
	}
	failOnErrorDiagnostics(t, "ConfigureProvider", resp.Diagnostics)
}

// listDatacenters calls the ListResource RPC with the given filters and collects the
// results, failing on the first error diagnostic.
func listDatacenters(ctx context.Context, t *testing.T, server tfprotov6.ProviderServer, schema *tfprotov6.Schema, filters map[string]string) []tfprotov6.ListResourceResult {
	t.Helper()

	listServer, config := listServerAndConfig(t, server, schema, filters)

	stream, err := listServer.ListResource(ctx, &tfprotov6.ListResourceRequest{
		TypeName:        datacenterListType,
		Config:          &config,
		IncludeResource: true,
	})
	if err != nil {
		t.Fatalf("ListResource: %v", err)
	}

	var results []tfprotov6.ListResourceResult //nolint:prealloc // the number of streamed results is not known upfront
	for result := range stream.Results {
		failOnErrorDiagnostics(t, "ListResource result", result.Diagnostics)
		if result.Identity == nil {
			t.Fatalf("every list result must carry an identity")
		}
		if result.Resource == nil {
			t.Fatalf("IncludeResource was set, so every result must carry a resource")
		}
		results = append(results, result)
	}

	return results
}

// listServerAndConfig returns the list-resource half of the provider server together
// with a list block config carrying the given filters.
//
//nolint:staticcheck // ListResource still lives on this temporary interface in terraform-plugin-go v0.31.0
func listServerAndConfig(t *testing.T, server tfprotov6.ProviderServer, schema *tfprotov6.Schema, filters map[string]string) (tfprotov6.ProviderServerWithListResource, tfprotov6.DynamicValue) {
	t.Helper()

	listServer, ok := server.(tfprotov6.ProviderServerWithListResource)
	if !ok {
		t.Fatalf("the provider server does not implement ListResource")
	}

	return listServer, dynamicValue(t, schema.ValueType(), listConfig(t, schema.ValueType(), filters))
}

// listConfig builds the value of a list block config with the given filters.
func listConfig(t *testing.T, configType tftypes.Type, filters map[string]string) tftypes.Value {
	t.Helper()

	object, ok := configType.(tftypes.Object)
	if !ok {
		t.Fatalf("the list resource config schema is not an object, got %T", configType)
	}

	filtersType, ok := object.AttributeTypes["filters"].(tftypes.List)
	if !ok {
		t.Fatalf("the filters attribute is not a list, got %T", object.AttributeTypes["filters"])
	}

	entries := make([]tftypes.Value, 0, len(filters))
	for name, value := range filters {
		entries = append(entries, tftypes.NewValue(filtersType.ElementType, map[string]tftypes.Value{
			"field_name":  tftypes.NewValue(tftypes.String, name),
			"field_value": tftypes.NewValue(tftypes.String, value),
		}))
	}

	filtersValue := tftypes.NewValue(filtersType, nil)
	if len(entries) > 0 {
		filtersValue = tftypes.NewValue(filtersType, entries)
	}

	return tftypes.NewValue(configType, map[string]tftypes.Value{"filters": filtersValue})
}

// decode unmarshals a dynamic value into the plain Go values its attributes hold, so
// that the assertions can reach the collection and number attributes as well as the
// string ones.
func decode(t *testing.T, value *tfprotov6.DynamicValue, valueType tftypes.Type) map[string]any {
	t.Helper()

	decoded, err := value.Unmarshal(valueType)
	if err != nil {
		t.Fatalf("failed to unmarshal a list result: %v", err)
	}

	converted := goValue(t, decoded)
	attributes, ok := converted.(map[string]any)
	if !ok {
		t.Fatalf("expected a list result to decode into an object, got %T", converted)
	}

	return attributes
}

// goValue converts a tftypes.Value into the plain Go value it holds: an object becomes
// a map[string]any, a list or set becomes a []any, a number becomes a float64, and a
// null or unknown value becomes nil.
func goValue(t *testing.T, value tftypes.Value) any {
	t.Helper()

	if value.IsNull() || !value.IsKnown() {
		return nil
	}

	switch valueType := value.Type(); {
	case valueType.Is(tftypes.String):
		var s string
		readValue(t, value, &s)
		return s

	case valueType.Is(tftypes.Bool):
		var b bool
		readValue(t, value, &b)
		return b

	case valueType.Is(tftypes.Number):
		var number big.Float
		readValue(t, value, &number)
		f, _ := number.Float64()
		return f

	case valueType.Is(tftypes.List{}), valueType.Is(tftypes.Set{}):
		var elements []tftypes.Value
		readValue(t, value, &elements)
		converted := make([]any, 0, len(elements))
		for _, element := range elements {
			converted = append(converted, goValue(t, element))
		}
		return converted

	case valueType.Is(tftypes.Object{}):
		var attributes map[string]tftypes.Value
		readValue(t, value, &attributes)
		converted := make(map[string]any, len(attributes))
		for name, attribute := range attributes {
			converted[name] = goValue(t, attribute)
		}
		return converted

	default:
		t.Fatalf("cannot convert a value of type %s", valueType)
		return nil
	}
}

// readValue reads a tftypes.Value into target, failing the test if it does not fit.
func readValue(t *testing.T, value tftypes.Value, target any) {
	t.Helper()

	if err := value.As(target); err != nil {
		t.Fatalf("failed to read a value of type %s: %v", value.Type(), err)
	}
}

// identityType builds the object type of a resource identity schema.
func identityType(schema *tfprotov6.ResourceIdentitySchema) tftypes.Type {
	attributeTypes := make(map[string]tftypes.Type, len(schema.IdentityAttributes))
	for _, attribute := range schema.IdentityAttributes {
		attributeTypes[attribute.Name] = attribute.Type
	}

	return tftypes.Object{AttributeTypes: attributeTypes}
}

// nullObject builds an object value whose attributes are all null.
func nullObject(t *testing.T, objectType tftypes.Type) tftypes.Value {
	t.Helper()

	object, ok := objectType.(tftypes.Object)
	if !ok {
		t.Fatalf("expected an object type, got %T", objectType)
	}

	attributes := make(map[string]tftypes.Value, len(object.AttributeTypes))
	for name, attributeType := range object.AttributeTypes {
		attributes[name] = tftypes.NewValue(attributeType, nil)
	}

	return tftypes.NewValue(objectType, attributes)
}

func dynamicValue(t *testing.T, valueType tftypes.Type, value tftypes.Value) tfprotov6.DynamicValue {
	t.Helper()

	dv, err := tfprotov6.NewDynamicValue(valueType, value)
	if err != nil {
		t.Fatalf("failed to build a dynamic value: %v", err)
	}

	return dv
}

// failOnErrorDiagnostics fails the test on the first error diagnostic.
func failOnErrorDiagnostics(t *testing.T, rpc string, diags []*tfprotov6.Diagnostic) {
	t.Helper()

	for _, d := range diags {
		if d != nil && d.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatalf("%s returned an error diagnostic: %s: %s", rpc, d.Summary, d.Detail)
		}
	}
}

func hasErrorDiagnostic(diags []*tfprotov6.Diagnostic) bool {
	for _, d := range diags {
		if d != nil && d.Severity == tfprotov6.DiagnosticSeverityError {
			return true
		}
	}

	return false
}
