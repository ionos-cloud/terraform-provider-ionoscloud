package compute_test

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

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

const ipblockListType = "ionoscloud_ipblock"

// TestIPBlockListResource drives the ListResource RPC end to end against a stubbed
// Cloud API, through the same muxed provider server that main.go serves.
//
// It covers the parts of a list resource for an SDKv2 managed resource that can only
// fail at runtime: that the mux is happy with the list resource and the managed
// resource coming from different servers, that the framework registers a list resource
// with no framework resource behind it, that the protocol schemas handed over by
// RawV6Schemas convert cleanly, and that the resource model fills the SDKv2 schema
// without a type mismatch.
//
// The shared helpers it calls (muxedProviderServer, decode, listServerAndConfig, ...)
// are declared once for the package in resource_datacenter_list_test.go.
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
	if _, ok := providerSchema.ListResourceSchemas[ipblockListType]; !ok {
		t.Fatalf("the ip block list resource was not registered; check RawV6Schemas and the SDKv2 resource identity")
	}
	if _, ok := providerSchema.ResourceSchemas[ipblockListType]; !ok {
		t.Fatalf("the ip block managed resource is missing from the merged schema")
	}

	identitySchemas, err := server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		t.Fatalf("GetResourceIdentitySchemas: %v", err)
	}
	failOnErrorDiagnostics(t, "GetResourceIdentitySchemas", identitySchemas.Diagnostics)

	identitySchema := identitySchemas.IdentitySchemas[ipblockListType]
	if identitySchema == nil {
		t.Fatalf("the SDKv2 ip block resource does not declare a resource identity")
	}

	configureProvider(ctx, t, server, providerSchema.Provider)

	listSchema := providerSchema.ListResourceSchemas[ipblockListType]
	resourceType := providerSchema.ResourceSchemas[ipblockListType].ValueType()

	t.Run("streams every ip block", func(t *testing.T) {
		results := listResults(ctx, t, server, ipblockListType, listSchema, nil)
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}

		assert.Equal(t, "prod-block", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Len(t, identity, 2, "the ip block identity is id + location")
		assert.Equal(t, "a1b2c3d4-0000-0000-0000-000000000001", identity["id"])
		assert.Equal(t, "de/txl", identity["location"])

		// Every attribute the mapper fills is asserted here, so that mapping a value
		// to the wrong attribute fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "a1b2c3d4-0000-0000-0000-000000000001", resource["id"])
		assert.Equal(t, "prod-block", resource["name"])
		assert.Equal(t, "de/txl", resource["location"])
		assert.Equal(t, int64(2), resource["size"])
		// ips is a TypeList, so unlike a TypeSet its order is part of the value.
		assert.Equal(t, []any{"192.0.2.10", "192.0.2.11"}, resource["ips"])
		assert.Equal(t, []any{map[string]any{
			"ip":                "192.0.2.10",
			"mac":               "02:01:0b:9d:4d:ce",
			"nic_id":            "11111111-0000-0000-0000-00000000000a",
			"server_id":         "22222222-0000-0000-0000-00000000000b",
			"server_name":       "web-1",
			"datacenter_id":     "33333333-0000-0000-0000-00000000000c",
			"datacenter_name":   "prod",
			"k8s_nodepool_uuid": "44444444-0000-0000-0000-00000000000d",
			"k8s_cluster_uuid":  "55555555-0000-0000-0000-00000000000e",
		}}, resource["ip_consumers"])
		assert.Nil(t, resource["timeouts"], "a listed ip block has no timeouts")

		// The second ip block reports only the properties the API always sets, which
		// pins that the pairing holds past the first result and that the properties the
		// API left out stay null instead of turning into zero values.
		assert.Equal(t, "staging-block", results[1].DisplayName)

		stagingIdentity := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "a1b2c3d4-0000-0000-0000-000000000002", stagingIdentity["id"])
		assert.Equal(t, "de/fra", stagingIdentity["location"])

		staging := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "a1b2c3d4-0000-0000-0000-000000000002", staging["id"])
		assert.Equal(t, "staging-block", staging["name"])
		assert.Equal(t, "de/fra", staging["location"])
		assert.Equal(t, int64(1), staging["size"])
		assert.Nil(t, staging["ips"])
		// ip_consumers is the one attribute that does not come back null when the API
		// omits it. It is Optional as well as Computed, so core_schema.go:108-112 routes
		// its *Resource Elem to BlockTypes rather than Attributes - and on the way out
		// to the protocol, toproto6/dynamic_value.go:30 runs ReifyNullCollectionBlocks,
		// which turns a null list/set *block* into an empty one. (datacenter's
		// cpu_architecture is Computed-only, so it stays an attribute, skips the reify
		// step and does decode to nil - as does `ips` here, whose Elem is a *Schema.)
		assert.Equal(t, []any{}, staging["ip_consumers"])

		// name is Optional on ionoscloud_ipblock, so an unnamed block is legitimate.
		// The display name has to fall back to the ID rather than render as a blank
		// row, and the name attribute itself must still come back null.
		assert.Equal(t, "a1b2c3d4-0000-0000-0000-000000000003", results[2].DisplayName)

		unnamed := decode(t, results[2].Resource, resourceType)
		assert.Equal(t, "a1b2c3d4-0000-0000-0000-000000000003", unnamed["id"])
		assert.Nil(t, unnamed["name"])
		assert.Equal(t, "es/vit", unnamed["location"])
	})

	t.Run("filters by location", func(t *testing.T) {
		results := listResults(ctx, t, server, ipblockListType, listSchema, map[string]string{"location": "de/fra"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "staging-block", results[0].DisplayName)
	})

	// Both advertised filter fields need their own case: MatchesFilters returns false
	// for a field_name that is missing from the mapper's map, so dropping "name" there
	// would leave the filter permanently matching nothing while the location case above
	// stays green.
	t.Run("filters by name", func(t *testing.T) {
		results := listResults(ctx, t, server, ipblockListType, listSchema, map[string]string{"name": "prod-block"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "prod-block", results[0].DisplayName)
	})

	// A name and a location that each match a *different* ip block, which pins that the
	// filters are ANDed and that neither is wired to the other's property.
	t.Run("applies every filter", func(t *testing.T) {
		results := listResults(ctx, t, server, ipblockListType, listSchema, map[string]string{
			"name":     "prod-block",
			"location": "de/fra",
		})
		if len(results) != 0 {
			t.Fatalf("expected 0 results, got %d", len(results))
		}
	})

	t.Run("rejects unknown filter fields", func(t *testing.T) {
		listServer, config := listServerAndConfig(t, server, listSchema, map[string]string{"nope": "value"})

		resp, err := listServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
			TypeName: ipblockListType,
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

// stubIPBlockAPI serves the ip block collection the list resource reads, and returns
// the URL to point IONOS_API_URL at. Every other path 404s on purpose, so an
// unexpected extra request surfaces as an error diagnostic instead of succeeding.
func stubIPBlockAPI(t *testing.T) string {
	t.Helper()

	ipBlocks := ionoscloudsdk.IpBlocks{
		Items: &[]ionoscloudsdk.IpBlock{
			{
				Id: new("a1b2c3d4-0000-0000-0000-000000000001"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Name:     new("prod-block"),
					Location: new("de/txl"),
					Size:     new(int32(2)),
					Ips:      &[]string{"192.0.2.10", "192.0.2.11"},
					IpConsumers: &[]ionoscloudsdk.IpConsumer{{
						Ip:              new("192.0.2.10"),
						Mac:             new("02:01:0b:9d:4d:ce"),
						NicId:           new("11111111-0000-0000-0000-00000000000a"),
						ServerId:        new("22222222-0000-0000-0000-00000000000b"),
						ServerName:      new("web-1"),
						DatacenterId:    new("33333333-0000-0000-0000-00000000000c"),
						DatacenterName:  new("prod"),
						K8sNodePoolUuid: new("44444444-0000-0000-0000-00000000000d"),
						K8sClusterUuid:  new("55555555-0000-0000-0000-00000000000e"),
					}},
				},
			},
			{
				// Only the properties the API always sets: location and size are the
				// two IpBlockProperties fields without `omitempty`.
				Id: new("a1b2c3d4-0000-0000-0000-000000000002"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Name:     new("staging-block"),
					Location: new("de/fra"),
					Size:     new(int32(1)),
				},
			},
			{
				// An unnamed block, which the schema allows, for the display-name
				// fallback.
				Id: new("a1b2c3d4-0000-0000-0000-000000000003"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Location: new("es/vit"),
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
		// The request the list resource builds is part of what is under test: drop the
		// explicit Limit and the SDK client falls back to limit=100, silently capping the
		// listing below the data source reading the same collection; drop Depth(1) and
		// the API returns links instead of properties.
		// Asserted as a literal, not as constant.IPBlockLimit, so that changing the
		// constant fails here too - docs/list-resources/ipblock.md documents 1000 as
		// the ceiling, and nothing else ties the doc to the code.
		if got := r.URL.Query().Get("limit"); got != "1000" {
			t.Errorf("expected limit=1000 on the ipblocks request, got %q; constant.IPBlockLimit is %d", got, constant.IPBlockLimit)
		}
		if got := r.URL.Query().Get("depth"); got != "1" {
			t.Errorf("expected depth=1 on the ipblocks request, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ipBlocks); err != nil {
			t.Errorf("failed to write the stubbed response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return server.URL
}
