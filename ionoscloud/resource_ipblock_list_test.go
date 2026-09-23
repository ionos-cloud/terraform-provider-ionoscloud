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
//
// It covers the parts of a list resource for an SDKv2 managed resource that can only
// fail at runtime: that the mux is happy with the list resource and the managed
// resource coming from different servers, that the framework registers a list resource
// with no framework resource behind it, that the protocol schemas handed over by
// RawV6Schemas convert cleanly, and that the state the resource's own writer produces
// survives the round trip through ResourceData.TfTypeResourceState into the list result,
// including the timeouts block identity.MappedItemFromResourceData nulls back out.
func TestIPBlockListResource(t *testing.T) {
	ctx := context.Background()

	t.Setenv("IONOS_API_URL", stubIPBlockAPI(t))
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
		t.Fatalf("the ip block list resource was not registered; check RawV6Schemas and the SDKv2 resource identity")
	}
	if _, ok := providerSchema.ResourceSchemas[ipBlockListType]; !ok {
		t.Fatalf("the ip block managed resource is missing from the merged schema")
	}

	identitySchemas, err := server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		t.Fatalf("GetResourceIdentitySchemas: %v", err)
	}
	failOnErrorDiagnostics(t, "GetResourceIdentitySchemas", identitySchemas.Diagnostics)

	identitySchema := identitySchemas.IdentitySchemas[ipBlockListType]
	if identitySchema == nil {
		t.Fatalf("the SDKv2 ip block resource does not declare a resource identity")
	}

	configureProvider(ctx, t, server, providerSchema.Provider)

	listSchema := providerSchema.ListResourceSchemas[ipBlockListType]
	resourceType := providerSchema.ResourceSchemas[ipBlockListType].ValueType()

	t.Run("streams every ip block", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, nil)
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}

		assert.Equal(t, "prod-ips", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000001", identity["id"])
		assert.Equal(t, "de/txl", identity["location"])

		// Every attribute the mapper fills is asserted here, so that mapping a value
		// to the wrong attribute fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000001", resource["id"])
		assert.Equal(t, "prod-ips", resource["name"])
		assert.Equal(t, "de/txl", resource["location"])
		assert.Equal(t, int64(2), resource["size"])
		assert.Equal(t, []any{"192.0.2.10", "192.0.2.11"}, resource["ips"])
		assert.Equal(t, []any{map[string]any{
			"ip":                "192.0.2.10",
			"mac":               "02:01:0b:0d:8f:cd",
			"nic_id":            "c0b0b0b0-0000-4000-8000-000000000001",
			"server_id":         "c0b0b0b0-0000-4000-8000-000000000002",
			"server_name":       "webserver",
			"datacenter_id":     "c0b0b0b0-0000-4000-8000-000000000003",
			"datacenter_name":   "prod",
			"k8s_nodepool_uuid": "c0b0b0b0-0000-4000-8000-000000000004",
			"k8s_cluster_uuid":  "c0b0b0b0-0000-4000-8000-000000000005",
		}}, resource["ip_consumers"])
		assert.Nil(t, resource["timeouts"], "a listed ip block has no timeouts")

		// The second ip block reports only the properties the API always sets, which
		// pins that the pairing holds past the first result and that the properties the
		// API left out do not turn into values of the first block's.
		assert.Equal(t, "staging-ips", results[1].DisplayName)

		stagingIdentity := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000002", stagingIdentity["id"])
		assert.Equal(t, "de/fra", stagingIdentity["location"])

		staging := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000002", staging["id"])
		assert.Equal(t, "staging-ips", staging["name"])
		assert.Equal(t, "de/fra", staging["location"])
		assert.Equal(t, int64(1), staging["size"])
		assert.Equal(t, []any{"198.51.100.7"}, staging["ips"])
		// ip_consumers is a nested block, and an omitted one reifies as an empty list
		// rather than null.
		assert.Equal(t, []any{}, staging["ip_consumers"])

		// name is Optional on an ip block, so an unnamed one has to fall back to its id
		// for a display name: the framework reads a blank one as a diagnostics-only event.
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000003", results[2].DisplayName)

		unnamed := decode(t, results[2].Resource, resourceType)
		assert.Nil(t, unnamed["name"])
		assert.Equal(t, "de/fra", unnamed["location"])
	})

	t.Run("filters by name", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{"name": "prod-ips"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "prod-ips", results[0].DisplayName)
	})

	t.Run("filters by location", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{"location": "de/fra"})
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}
		assert.Equal(t, "staging-ips", results[0].DisplayName)
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000003", results[1].DisplayName)
	})

	t.Run("applies every filter", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{
			"name":     "prod-ips",
			"location": "de/fra",
		})
		if len(results) != 0 {
			t.Fatalf("expected no result for a name and a location that belong to different ip blocks, got %d", len(results))
		}
	})

	t.Run("rejects unknown filter fields", func(t *testing.T) {
		listServer, config := listServerAndConfig(t, server, listSchema, map[string]string{"nope": "value"})

		resp, err := listServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
			TypeName: ipBlockListType,
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

// stubIPBlockAPI serves the ip block collection the list resource reads, asserts the
// query the fetch closure builds, and returns the URL to point IONOS_API_URL at.
func stubIPBlockAPI(t *testing.T) string {
	t.Helper()

	ipBlocks := ionoscloudsdk.IpBlocks{
		Items: &[]ionoscloudsdk.IpBlock{
			{
				Id: new("d3b07384-d9a0-4d1e-8f2a-000000000001"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Name:     new("prod-ips"),
					Location: new("de/txl"),
					Size:     new(int32(2)),
					Ips:      &[]string{"192.0.2.10", "192.0.2.11"},
					IpConsumers: &[]ionoscloudsdk.IpConsumer{{
						Ip:              new("192.0.2.10"),
						Mac:             new("02:01:0b:0d:8f:cd"),
						NicId:           new("c0b0b0b0-0000-4000-8000-000000000001"),
						ServerId:        new("c0b0b0b0-0000-4000-8000-000000000002"),
						ServerName:      new("webserver"),
						DatacenterId:    new("c0b0b0b0-0000-4000-8000-000000000003"),
						DatacenterName:  new("prod"),
						K8sNodePoolUuid: new("c0b0b0b0-0000-4000-8000-000000000004"),
						K8sClusterUuid:  new("c0b0b0b0-0000-4000-8000-000000000005"),
					}},
				},
			},
			{
				Id: new("d3b07384-d9a0-4d1e-8f2a-000000000002"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Name:     new("staging-ips"),
					Location: new("de/fra"),
					Size:     new(int32(1)),
					Ips:      &[]string{"198.51.100.7"},
				},
			},
			{
				Id: new("d3b07384-d9a0-4d1e-8f2a-000000000003"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Location: new("de/fra"),
					Size:     new(int32(1)),
					Ips:      &[]string{"198.51.100.8"},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/ipblocks") {
			http.NotFound(w, r)
			return
		}
		// The two options the fetch closure sets explicitly. The limit is spelled out
		// rather than read from constant.IPBlockLimit, so that changing the constant
		// does not silently change what this asserts.
		if got := r.URL.Query().Get("depth"); got != "1" {
			t.Errorf("expected the ip block listing to request depth=1, got %q", got)
		}
		if got := r.URL.Query().Get("limit"); got != "1000" {
			t.Errorf("expected the ip block listing to request limit=1000, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ipBlocks); err != nil {
			t.Errorf("failed to write the stubbed response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return server.URL
}
