package ionoscloud_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	ionoscloudsdk "github.com/ionos-cloud/sdk-go/v6"
	"github.com/stretchr/testify/assert"
)

const ipBlockListType = "ionoscloud_ipblock"

// TestIPBlockListResource drives the ListResource RPC end to end against a stubbed
// Cloud API, through the same muxed provider server that main.go serves.
func TestIPBlockListResource(t *testing.T) {
	ctx := context.Background()

	t.Setenv("IONOS_API_URL", stubIPBlocksAPI(t))
	t.Setenv("IONOS_TOKEN", "token-for-the-stub")

	server := muxedProviderServer(ctx, t)

	providerSchema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema: %v", err)
	}
	failOnErrorDiagnostics(t, "GetProviderSchema", providerSchema.Diagnostics)

	// The framework only registers a list resource that has no framework resource
	// behind it once RawV6Schemas has supplied both protocol schemas.
	if _, ok := providerSchema.ListResourceSchemas[ipBlockListType]; !ok {
		t.Fatalf("the IP Block list resource was not registered; check RawV6Schemas and the SDKv2 resource identity")
	}
	if _, ok := providerSchema.ResourceSchemas[ipBlockListType]; !ok {
		t.Fatalf("the IP Block managed resource is missing from the merged schema")
	}

	identitySchemas, err := server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		t.Fatalf("GetResourceIdentitySchemas: %v", err)
	}
	failOnErrorDiagnostics(t, "GetResourceIdentitySchemas", identitySchemas.Diagnostics)

	identitySchema := identitySchemas.IdentitySchemas[ipBlockListType]
	if identitySchema == nil {
		t.Fatalf("the SDKv2 IP Block resource does not declare a resource identity")
	}

	configureProvider(ctx, t, server, providerSchema.Provider)

	listSchema := providerSchema.ListResourceSchemas[ipBlockListType]
	resourceType := providerSchema.ResourceSchemas[ipBlockListType].ValueType()

	t.Run("streams every IP Block", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, nil, true)
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}

		assert.Equal(t, "web", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "7c3e9a12-4b5d-4f6e-8a9b-000000000001", identity["id"])
		assert.Equal(t, "de/fra", identity["location"])

		// Every attribute the writer fills is asserted here, so that mapping a value
		// to the wrong attribute fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "7c3e9a12-4b5d-4f6e-8a9b-000000000001", resource["id"])
		assert.Equal(t, "web", resource["name"])
		assert.Equal(t, "de/fra", resource["location"])
		assert.Equal(t, int64(2), resource["size"])
		assert.Equal(t, []any{"203.0.113.10", "203.0.113.11"}, resource["ips"])
		assert.Equal(t, []any{map[string]any{
			"ip":                "203.0.113.10",
			"mac":               "02:01:5e:00:00:01",
			"nic_id":            "5a1b2c3d-0000-4000-8000-0000000000a1",
			"server_id":         "5a1b2c3d-0000-4000-8000-0000000000b1",
			"server_name":       "web-server",
			"datacenter_id":     "5a1b2c3d-0000-4000-8000-0000000000c1",
			"datacenter_name":   "web-datacenter",
			"k8s_nodepool_uuid": "5a1b2c3d-0000-4000-8000-0000000000d1",
			"k8s_cluster_uuid":  "5a1b2c3d-0000-4000-8000-0000000000e1",
		}}, resource["ip_consumers"])
		assert.Nil(t, resource["timeouts"], "a listed IP Block has no timeouts")

		// The second IP Block carries only its id, name and the required properties,
		// which pins that the pairing holds past the first result and that what the
		// API left out stays null, or an empty block, instead of a zero value.
		assert.Equal(t, "db", results[1].DisplayName)

		secondIdentity := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "7c3e9a12-4b5d-4f6e-8a9b-000000000002", secondIdentity["id"])
		assert.Equal(t, "us/las", secondIdentity["location"])

		second := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "7c3e9a12-4b5d-4f6e-8a9b-000000000002", second["id"])
		assert.Equal(t, "db", second["name"])
		assert.Equal(t, "us/las", second["location"])
		assert.Equal(t, int64(1), second["size"])
		assert.Nil(t, second["ips"])
		assert.Equal(t, []any{}, second["ip_consumers"])

		// The third IP Block has no name, so it is labelled with its id.
		assert.Equal(t, "7c3e9a12-4b5d-4f6e-8a9b-000000000003", results[2].DisplayName)

		thirdIdentity := decode(t, results[2].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "7c3e9a12-4b5d-4f6e-8a9b-000000000003", thirdIdentity["id"])
		assert.Equal(t, "de/txl", thirdIdentity["location"])

		assert.Nil(t, decode(t, results[2].Resource, resourceType)["name"])
	})

	t.Run("filters by name", func(t *testing.T) {
		// An IP Block without a name is matched by an empty field_value.
		unnamed := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{"name": ""}, true)
		if len(unnamed) != 1 {
			t.Fatalf("expected 1 result, got %d", len(unnamed))
		}
		assert.Equal(t, "7c3e9a12-4b5d-4f6e-8a9b-000000000003", unnamed[0].DisplayName)

		named := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{"name": "db"}, true)
		if len(named) != 1 {
			t.Fatalf("expected 1 result, got %d", len(named))
		}
		assert.Equal(t, "db", named[0].DisplayName)
	})

	t.Run("filters by location", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{"location": "de/fra"}, true)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "web", results[0].DisplayName)
	})

	t.Run("applies every filter", func(t *testing.T) {
		// web is in de/fra and db in us/las, so no IP Block matches both filters.
		results := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{"name": "web", "location": "us/las"}, true)
		if len(results) != 0 {
			t.Fatalf("expected no results, got %d", len(results))
		}
	})

	t.Run("identity only", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, nil, false)
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}

		assert.Equal(t, "web", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "7c3e9a12-4b5d-4f6e-8a9b-000000000001", identity["id"])
		assert.Equal(t, "de/fra", identity["location"])

		for _, result := range results {
			assert.Nil(t, result.Resource, "IncludeResource was not set, so no result carries a resource")
		}
	})

	t.Run("rejects unknown filter fields", func(t *testing.T) {
		// size, ips and ip_consumers are kept out of the allow-list, see
		// ListResourceConfigSchema.
		for _, field := range []string{"nope", "size", "ips", "ip_consumers"} {
			listServer, config := listServerAndConfig(t, server, listSchema, map[string]string{field: "value"})

			resp, err := listServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
				TypeName: ipBlockListType,
				Config:   &config,
			})
			if err != nil {
				t.Fatalf("ValidateListResourceConfig: %v", err)
			}
			if !hasErrorDiagnostic(resp.Diagnostics) {
				t.Fatalf("expected a validation error for the filter field %q", field)
			}
		}
	})

	t.Run("accepts every allowed filter field", func(t *testing.T) {
		for _, field := range []string{"name", "location"} {
			listServer, config := listServerAndConfig(t, server, listSchema, map[string]string{field: "value"})

			resp, err := listServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
				TypeName: ipBlockListType,
				Config:   &config,
			})
			if err != nil {
				t.Fatalf("ValidateListResourceConfig: %v", err)
			}
			failOnErrorDiagnostics(t, "ValidateListResourceConfig", resp.Diagnostics)
		}
	})
}

// stubIPBlocksAPI serves the IP Block collection the list resource reads, and returns
// the URL to point IONOS_API_URL at. Every request for that collection must carry
// depth=1, offset=0 and limit=1000; any other path gets a 404.
func stubIPBlocksAPI(t *testing.T) string {
	t.Helper()

	ipBlocks := ionoscloudsdk.IpBlocks{
		Items: &[]ionoscloudsdk.IpBlock{
			{
				Id: new("7c3e9a12-4b5d-4f6e-8a9b-000000000001"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Name:     new("web"),
					Location: new("de/fra"),
					Size:     new(int32(2)),
					Ips:      &[]string{"203.0.113.10", "203.0.113.11"},
					IpConsumers: &[]ionoscloudsdk.IpConsumer{{
						Ip:              new("203.0.113.10"),
						Mac:             new("02:01:5e:00:00:01"),
						NicId:           new("5a1b2c3d-0000-4000-8000-0000000000a1"),
						ServerId:        new("5a1b2c3d-0000-4000-8000-0000000000b1"),
						ServerName:      new("web-server"),
						DatacenterId:    new("5a1b2c3d-0000-4000-8000-0000000000c1"),
						DatacenterName:  new("web-datacenter"),
						K8sNodePoolUuid: new("5a1b2c3d-0000-4000-8000-0000000000d1"),
						K8sClusterUuid:  new("5a1b2c3d-0000-4000-8000-0000000000e1"),
					}},
				},
			},
			{
				Id: new("7c3e9a12-4b5d-4f6e-8a9b-000000000002"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Name:     new("db"),
					Location: new("us/las"),
					Size:     new(int32(1)),
				},
			},
			{
				Id: new("7c3e9a12-4b5d-4f6e-8a9b-000000000003"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Location: new("de/txl"),
					Size:     new(int32(1)),
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/ipblocks") {
			http.NotFound(w, r)
			return
		}
		for param, want := range map[string]string{"depth": "1", "offset": "0", "limit": "1000"} {
			if got := r.URL.Query().Get(param); got != want {
				t.Errorf("expected %s=%s on the IP Block list request, got %q", param, want, got)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ipBlocks); err != nil {
			t.Errorf("failed to write the stubbed response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return server.URL
}
