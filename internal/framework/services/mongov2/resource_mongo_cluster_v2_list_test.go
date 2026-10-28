package mongov2_test

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
	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/stretchr/testify/assert"

	fwprovider "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/provider"
)

const (
	mongoClusterListType = "ionoscloud_mongo_cluster_v2"
	// mongoListLocation is the single location the stub serves; the list config filters on it
	// so only this regional endpoint is queried and the other locations are never contacted.
	mongoListLocation = "gb/bhx"
)

// TestMongoClusterV2ListResource drives the ListResource RPC of the MongoDB v2 cluster
// resource end to end against a stubbed regional API, with include_resource = true.
//
// Going through the RPC (rather than calling the mapper directly) is the point: it exercises
// the protocol serialization of the mapped model against the list schema — the path
// configuration generation uses — including the edition-specific nested attributes and the
// hand-built null timeouts object.
//
// The regional endpoint is stubbed through a file-config override: a temporary IONOS config
// file maps the mongodbv2 product at this location to the stub server, and the provider is
// pointed at it via the IONOS_CONFIG_FILE environment variable.
func TestMongoClusterV2ListResource(t *testing.T) {
	ctx := context.Background()

	stubMongoAPI(t)
	t.Setenv("IONOS_TOKEN", "token-for-the-stub")

	server := providerserver.NewProtocol6(fwprovider.New())()

	providerSchema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema: %v", err)
	}
	failOnErrorDiagnostics(t, "GetProviderSchema", providerSchema.Diagnostics)

	if _, ok := providerSchema.ListResourceSchemas[mongoClusterListType]; !ok {
		t.Fatalf("the mongo cluster v2 list resource was not registered")
	}

	identitySchemas, err := server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		t.Fatalf("GetResourceIdentitySchemas: %v", err)
	}
	failOnErrorDiagnostics(t, "GetResourceIdentitySchemas", identitySchemas.Diagnostics)

	identitySchema := identitySchemas.IdentitySchemas[mongoClusterListType]
	if identitySchema == nil {
		t.Fatalf("the mongo cluster v2 resource does not declare a resource identity")
	}

	configureProvider(ctx, t, server, providerSchema.Provider)

	listSchema := providerSchema.ListResourceSchemas[mongoClusterListType]
	resourceType := providerSchema.ResourceSchemas[mongoClusterListType].ValueType()

	results := listResults(ctx, t, server, mongoClusterListType, listSchema, map[string]string{"location": mongoListLocation}, true)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	t.Run("enterprise cluster", func(t *testing.T) {
		assert.Equal(t, "enterprise-cluster", results[0].DisplayName)

		id := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "11111111-1111-1111-1111-111111111111", id["id"])
		assert.Equal(t, mongoListLocation, id["location"])

		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "11111111-1111-1111-1111-111111111111", resource["id"])
		assert.Equal(t, mongoListLocation, resource["location"])
		assert.Equal(t, "enterprise-cluster", resource["name"])
		assert.Equal(t, "ENTERPRISE", resource["edition"])
		assert.Equal(t, "7.0", resource["version"])
		assert.Equal(t, "the enterprise cluster", resource["description"])
		assert.Equal(t, "mongodb://enterprise.example.com", resource["connection_uri"])
		assert.Equal(t, true, resource["bi_connector_enabled"])
		assert.Equal(t, true, resource["logs_enabled"])
		assert.Equal(t, false, resource["metrics_enabled"])

		// ENTERPRISE reports a flat instances block and sharding; BUSINESS-only fields are null.
		assert.Equal(t, map[string]any{
			"cores":        int64(2),
			"count":        int64(3),
			"ram":          int64(16),
			"storage_size": int64(100),
		}, resource["instances"])
		assert.Equal(t, map[string]any{
			"enabled": true,
			"shards":  int64(3),
		}, resource["sharding"])
		assert.Nil(t, resource["instance_template_id"])
		assert.Nil(t, resource["instance_count"])

		assertCommonFields(t, resource)
	})

	t.Run("business cluster", func(t *testing.T) {
		assert.Equal(t, "business-cluster", results[1].DisplayName)

		id := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "22222222-2222-2222-2222-222222222222", id["id"])
		assert.Equal(t, mongoListLocation, id["location"])

		resource := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "22222222-2222-2222-2222-222222222222", resource["id"])
		assert.Equal(t, mongoListLocation, resource["location"])
		assert.Equal(t, "business-cluster", resource["name"])
		assert.Equal(t, "BUSINESS", resource["edition"])
		assert.Equal(t, "7.0", resource["version"])

		// BUSINESS reports the template id and instance count; ENTERPRISE-only blocks and
		// bi_connector_enabled stay null.
		assert.Equal(t, "template-xl", resource["instance_template_id"])
		assert.Equal(t, int64(3), resource["instance_count"])
		assert.Nil(t, resource["instances"])
		assert.Nil(t, resource["sharding"])
		assert.Nil(t, resource["bi_connector_enabled"])

		assertCommonFields(t, resource)
	})

	t.Run("identity only", func(t *testing.T) {
		results := listResults(ctx, t, server, mongoClusterListType, listSchema, map[string]string{"location": mongoListLocation}, false)
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}
		for _, result := range results {
			assert.Nil(t, result.Resource, "no resource must be streamed when include_resource is false")
			assert.NotNil(t, result.Identity, "every result must carry an identity")
		}
	})
}

// assertCommonFields asserts the attributes both editions share, including that the API's
// never-returned password stays null and that a listed cluster carries no timeouts block.
func assertCommonFields(t *testing.T, resource map[string]any) {
	t.Helper()

	assert.Equal(t, map[string]any{
		"datacenter_id":      "dc-id",
		"lan_id":             "1",
		"instance_addresses": []any{"10.0.0.1"},
	}, resource["connections"])

	assert.Equal(t, map[string]any{
		"day_of_the_week": "Sunday",
		"time":            "10:00:00",
	}, resource["maintenance_window"])

	assert.Equal(t, map[string]any{
		"location":       "eu-central-2",
		"retention_days": int64(7),
	}, resource["snapshot"])

	credentials, ok := resource["credentials"].(map[string]any)
	if !ok {
		t.Fatalf("expected credentials to decode into an object, got %T", resource["credentials"])
	}
	assert.Equal(t, "admin", credentials["username"])
	assert.Equal(t, "testdb", credentials["database"])
	assert.Nil(t, credentials["password"], "the API never returns the password")

	assert.Nil(t, resource["timeouts"], "a listed cluster has no timeouts")
	assert.Nil(t, resource["restore_from_snapshot"], "restore_from_snapshot is write-only and never read back")
}

// stubMongoAPI serves the cluster collection the list resource reads and writes a temporary
// file config that overrides the MongoDB v2 endpoint for mongoListLocation to point at it.
func stubMongoAPI(t *testing.T) {
	t.Helper()

	clusters := mongov3.ClusterReadList{
		Items: []mongov3.ClusterRead{
			{
				Id:       "11111111-1111-1111-1111-111111111111",
				Metadata: mongov3.ClusterMetadata{ConnectionUri: new("mongodb://enterprise.example.com")},
				Properties: mongov3.EnterpriseClusterAsCluster(&mongov3.EnterpriseCluster{
					Name:               "enterprise-cluster",
					Description:        new("the enterprise cluster"),
					Version:            "7.0",
					Edition:            mongov3.MONGOEDITION_ENTERPRISE,
					BiConnectorEnabled: new(true),
					LogsEnabled:        new(true),
					MetricsEnabled:     new(false),
					Instances:          mongov3.InstanceConfiguration{Count: 3, Ram: 16, Cores: 2, StorageSize: 100},
					Sharding:           &mongov3.MongoSharding{Enabled: new(true), Shards: new(int32(3))},
					Connection:         commonConnection(),
					MaintenanceWindow:  commonMaintenanceWindow(),
					Snapshot:           commonSnapshot(),
					Credentials:        commonCredentials(),
				}),
			},
			{
				Id:       "22222222-2222-2222-2222-222222222222",
				Metadata: mongov3.ClusterMetadata{ConnectionUri: new("mongodb://business.example.com")},
				Properties: mongov3.BusinessClusterAsCluster(&mongov3.BusinessCluster{
					Name:               "business-cluster",
					Version:            "7.0",
					Edition:            mongov3.MONGOEDITION_BUSINESS,
					InstanceTemplateId: "template-xl",
					InstanceCount:      3,
					Connection:         commonConnection(),
					MaintenanceWindow:  commonMaintenanceWindow(),
					Snapshot:           commonSnapshot(),
					Credentials:        commonCredentials(),
				}),
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/clusters") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(clusters); err != nil {
			t.Errorf("failed to write the stubbed response: %v", err)
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
`, fileconfiguration.MongoDBV2, mongoListLocation, server.URL)
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatalf("failed to write the file config: %v", err)
	}
	t.Setenv(shared.IonosFilePathEnvVar, configPath)
}

func commonConnection() mongov3.MongoClusterConnection {
	return mongov3.MongoClusterConnection{DatacenterId: "dc-id", LanId: "1", InstanceAddresses: []string{"10.0.0.1"}}
}

func commonMaintenanceWindow() mongov3.MaintenanceWindow {
	return mongov3.MaintenanceWindow{Time: "10:00:00", DayOfTheWeek: mongov3.DAYOFTHEWEEK_SUNDAY}
}

func commonSnapshot() mongov3.ClusterSnapshot {
	return mongov3.ClusterSnapshot{Location: "eu-central-2", RetentionDays: 7}
}

func commonCredentials() *mongov3.MongoUser {
	return &mongov3.MongoUser{Username: "admin", Database: "testdb"}
}

// --- protocol test helpers ---

// configureProvider runs ConfigureProvider with an all-null config, which leaves the
// provider reading its credentials from the environment.
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

// listResults calls the ListResource RPC of typeName with the given filters and collects
// the results, failing on the first error diagnostic.
//
//nolint:staticcheck // ListResource still lives on this temporary interface in terraform-plugin-go
func listResults(ctx context.Context, t *testing.T, server tfprotov6.ProviderServer, typeName string, schema *tfprotov6.Schema, filters map[string]string, includeResource bool) []tfprotov6.ListResourceResult {
	t.Helper()

	listServer, ok := server.(tfprotov6.ProviderServerWithListResource)
	if !ok {
		t.Fatalf("the provider server does not implement ListResource")
	}

	config := dynamicValue(t, schema.ValueType(), listConfig(t, schema.ValueType(), filters))

	stream, err := listServer.ListResource(ctx, &tfprotov6.ListResourceRequest{
		TypeName:        typeName,
		Config:          &config,
		IncludeResource: includeResource,
	})
	if err != nil {
		t.Fatalf("ListResource: %v", err)
	}

	var results []tfprotov6.ListResourceResult //nolint:prealloc // the number of streamed results is not known upfront
	for result := range stream.Results {
		failOnErrorDiagnostics(t, "ListResource result", result.Diagnostics)
		results = append(results, result)
	}

	return results
}

// listConfig builds the value of a list block config with the given filters.
func listConfig(t *testing.T, configType tftypes.Type, filters map[string]string) tftypes.Value {
	t.Helper()

	object := configType.(tftypes.Object)
	filtersType := object.AttributeTypes["filters"].(tftypes.List)

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

// decode unmarshals a dynamic value into the plain Go values its attributes hold.
func decode(t *testing.T, value *tfprotov6.DynamicValue, valueType tftypes.Type) map[string]any {
	t.Helper()

	decoded, err := value.Unmarshal(valueType)
	if err != nil {
		t.Fatalf("failed to unmarshal a list result: %v", err)
	}

	return goValue(t, decoded).(map[string]any)
}

// goValue converts a tftypes.Value into the plain Go value it holds: an object becomes a
// map[string]any, a list or set becomes a []any, a whole number becomes an int64, and a
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
		if !number.IsInt() {
			t.Fatalf("expected a whole number, got %s", number.String())
		}
		n, _ := number.Int64()
		return n
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

	object := objectType.(tftypes.Object)
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

func failOnErrorDiagnostics(t *testing.T, rpc string, diags []*tfprotov6.Diagnostic) {
	t.Helper()

	for _, d := range diags {
		if d != nil && d.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatalf("%s returned an error diagnostic: %s: %s", rpc, d.Summary, d.Detail)
		}
	}
}
