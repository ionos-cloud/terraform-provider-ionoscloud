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
//
// The shared helpers it calls - muxedProviderServer, configureProvider, listResults,
// decode and the rest - are declared once for the package in
// resource_datacenter_list_test.go.
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
		t.Fatalf("the ipblock list resource was not registered; check RawV6Schemas and the SDKv2 resource identity")
	}
	if _, ok := providerSchema.ResourceSchemas[ipBlockListType]; !ok {
		t.Fatalf("the ipblock managed resource is missing from the merged schema")
	}

	identitySchemas, err := server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		t.Fatalf("GetResourceIdentitySchemas: %v", err)
	}
	failOnErrorDiagnostics(t, "GetResourceIdentitySchemas", identitySchemas.Diagnostics)

	identitySchema := identitySchemas.IdentitySchemas[ipBlockListType]
	if identitySchema == nil {
		t.Fatalf("the SDKv2 ipblock resource does not declare a resource identity")
	}

	configureProvider(ctx, t, server, providerSchema.Provider)

	listSchema := providerSchema.ListResourceSchemas[ipBlockListType]
	resourceType := providerSchema.ResourceSchemas[ipBlockListType].ValueType()

	t.Run("streams every ip block", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, nil)
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}

		assert.Equal(t, "public-ips", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "b1a0e5f2-0000-4000-8000-000000000001", identity["id"])
		assert.Equal(t, "de/txl", identity["location"])
		assert.Len(t, identity, 2, "the ipblock identity is id plus location")

		// Every attribute the writer fills is asserted here, so that writing a value to
		// the wrong key fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "b1a0e5f2-0000-4000-8000-000000000001", resource["id"])
		assert.Equal(t, "public-ips", resource["name"])
		assert.Equal(t, "de/txl", resource["location"])
		assert.Equal(t, int64(2), resource["size"])
		assert.Equal(t, []any{"203.0.113.10", "203.0.113.11"}, resource["ips"])
		assert.Equal(t, []any{map[string]any{
			"ip":                "203.0.113.10",
			"mac":               "02:01:0b:1a:0e:5f",
			"nic_id":            "c2b1a0e5-0000-4000-8000-00000000000a",
			"server_id":         "d3c2b1a0-0000-4000-8000-00000000000b",
			"server_name":       "web-01",
			"datacenter_id":     "e4d3c2b1-0000-4000-8000-00000000000c",
			"datacenter_name":   "prod",
			"k8s_nodepool_uuid": "f5e4d3c2-0000-4000-8000-00000000000d",
			"k8s_cluster_uuid":  "a6f5e4d3-0000-4000-8000-00000000000e",
		}}, resource["ip_consumers"])
		assert.Nil(t, resource["timeouts"], "a listed ip block has no timeouts")

		// The second ip block reports only the properties the API always sets, which
		// pins that the pairing holds past the first result and that the properties the
		// API left out stay null instead of turning into zero values.
		assert.Equal(t, "spare-ips", results[1].DisplayName)

		secondIdentity := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "b1a0e5f2-0000-4000-8000-000000000002", secondIdentity["id"])
		assert.Equal(t, "de/fra", secondIdentity["location"])

		second := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "b1a0e5f2-0000-4000-8000-000000000002", second["id"])
		assert.Equal(t, "spare-ips", second["name"])
		assert.Equal(t, "de/fra", second["location"])
		assert.Equal(t, int64(1), second["size"])
		assert.Equal(t, []any{"198.51.100.7"}, second["ips"])
		// ip_consumers is Optional + Computed with an Elem: &schema.Resource{}, so
		// core_schema.go:108-112 renders it as a nested BLOCK rather than an attribute,
		// and toproto6.DynamicValue's ReifyNullCollectionBlocks turns the null block the
		// API omitted into an EMPTY list. An omitted attribute would decode to nil - see
		// the name assertion on the third result below.
		assert.Equal(t, []any{}, second["ip_consumers"])

		// The third ip block has no name at all, which is legal: name is Optional on
		// ionoscloud_ipblock. It covers the display-name fallback to the UUID, which the
		// second result cannot reach because it keeps its name.
		assert.Equal(t, "b1a0e5f2-0000-4000-8000-000000000003", results[2].DisplayName)

		third := decode(t, results[2].Resource, resourceType)
		assert.Equal(t, "b1a0e5f2-0000-4000-8000-000000000003", third["id"])
		assert.Nil(t, third["name"], "an omitted optional attribute decodes to null, not an empty string")
		assert.Equal(t, "de/fra", third["location"])
		assert.Nil(t, third["ips"], "ips is a Computed-only attribute, so an omitted one stays null")
		assert.Equal(t, []any{}, third["ip_consumers"])
	})

	// One subtest per field in the FilterAttribute allow-list: MatchesFilters returns
	// false for a field_name the mapper forgot to put in its map, so an untested field
	// is a filter that silently matches nothing while the other subtests stay green.
	t.Run("filters by name", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{"name": "spare-ips"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "spare-ips", results[0].DisplayName)
	})

	t.Run("filters by location", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{"location": "de/txl"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "public-ips", results[0].DisplayName)
	})

	// Two fields whose values match DIFFERENT stub items: pins that the filters are
	// ANDed, and that neither field is wired to the other's property.
	t.Run("applies every filter", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{
			"name":     "public-ips",
			"location": "de/fra",
		})
		if len(results) != 0 {
			t.Fatalf("expected 0 results, got %d", len(results))
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

// stubIPBlockAPI serves the ipblock collection the list resource reads, and returns the
// URL to point IONOS_API_URL at. Every other path 404s on purpose, so an unexpected
// extra request surfaces as an error diagnostic instead of succeeding.
func stubIPBlockAPI(t *testing.T) string {
	t.Helper()

	ipBlocks := ionoscloudsdk.IpBlocks{
		Items: &[]ionoscloudsdk.IpBlock{
			{
				// Every property the writer reads is set.
				Id: new("b1a0e5f2-0000-4000-8000-000000000001"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Name:     new("public-ips"),
					Location: new("de/txl"),
					Size:     new(int32(2)),
					Ips:      &[]string{"203.0.113.10", "203.0.113.11"},
					IpConsumers: &[]ionoscloudsdk.IpConsumer{{
						Ip:              new("203.0.113.10"),
						Mac:             new("02:01:0b:1a:0e:5f"),
						NicId:           new("c2b1a0e5-0000-4000-8000-00000000000a"),
						ServerId:        new("d3c2b1a0-0000-4000-8000-00000000000b"),
						ServerName:      new("web-01"),
						DatacenterId:    new("e4d3c2b1-0000-4000-8000-00000000000c"),
						DatacenterName:  new("prod"),
						K8sNodePoolUuid: new("f5e4d3c2-0000-4000-8000-00000000000d"),
						K8sClusterUuid:  new("a6f5e4d3-0000-4000-8000-00000000000e"),
					}},
				},
			},
			{
				// Only the properties the API always returns, so that ip_consumers can be
				// asserted as the empty block it reifies to. name stays set here: the
				// display-name and filter assertions need something to match on.
				Id: new("b1a0e5f2-0000-4000-8000-000000000002"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Name:     new("spare-ips"),
					Location: new("de/fra"),
					Size:     new(int32(1)),
					Ips:      &[]string{"198.51.100.7"},
				},
			},
			{
				// No name, which is legal because name is Optional on the resource. This
				// is what covers the displayName fallback in mapIPBlock.
				Id: new("b1a0e5f2-0000-4000-8000-000000000003"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Location: new("de/fra"),
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
		// The request the fetch closure builds is part of what is under test - the stub
		// answers any query string identically, so without these the options are
		// unpinned. The limit is asserted as a literal rather than as
		// constant.IPBlockLimit because that literal is the only thing tying the number
		// in docs/list-resources/ipblock.md to the code.
		if got := r.URL.Query().Get("depth"); got != "1" {
			t.Errorf("expected depth=1 on the /ipblocks request, got %q", got)
		}
		if got := r.URL.Query().Get("limit"); got != "1000" {
			t.Errorf("expected limit=1000 on the /ipblocks request, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ipBlocks); err != nil {
			t.Errorf("failed to write the stubbed response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return server.URL
}
