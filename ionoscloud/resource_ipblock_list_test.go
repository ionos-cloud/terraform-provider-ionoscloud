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
// listServerAndConfig, decode, identityType and the rest - are declared once for the
// package in resource_datacenter_list_test.go.
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
		assert.Len(t, identity, 2, "the ip block identity is id plus location")
		assert.Equal(t, "b1b07384-d9a0-4d1e-8f2a-000000000001", identity["id"])
		assert.Equal(t, "de/txl", identity["location"])

		// Every attribute the writer fills is asserted here, so that writing a value to
		// the wrong key fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "b1b07384-d9a0-4d1e-8f2a-000000000001", resource["id"])
		assert.Equal(t, "prod-ips", resource["name"])
		assert.Equal(t, "de/txl", resource["location"])
		assert.Equal(t, int64(2), resource["size"])
		assert.Equal(t, []any{"203.0.113.1", "203.0.113.2"}, resource["ips"])
		assert.Equal(t, []any{map[string]any{
			"ip":                "203.0.113.1",
			"mac":               "02:01:0b:0d:8f:b1",
			"nic_id":            "c1c07384-d9a0-4d1e-8f2a-00000000000a",
			"server_id":         "c2c07384-d9a0-4d1e-8f2a-00000000000b",
			"server_name":       "web-01",
			"datacenter_id":     "c3c07384-d9a0-4d1e-8f2a-00000000000c",
			"datacenter_name":   "prod",
			"k8s_nodepool_uuid": "c4c07384-d9a0-4d1e-8f2a-00000000000d",
			"k8s_cluster_uuid":  "c5c07384-d9a0-4d1e-8f2a-00000000000e",
		}}, resource["ip_consumers"])
		assert.Nil(t, resource["timeouts"], "a listed ip block has no timeouts")

		// The second ip block reports only the properties the API always sets, which
		// pins that the pairing holds past the first result and that the properties the
		// API left out stay null instead of turning into zero values.
		assert.Equal(t, "staging-ips", results[1].DisplayName)

		stagingIdentity := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "b1b07384-d9a0-4d1e-8f2a-000000000002", stagingIdentity["id"])
		assert.Equal(t, "de/fra", stagingIdentity["location"])

		staging := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "b1b07384-d9a0-4d1e-8f2a-000000000002", staging["id"])
		assert.Equal(t, "staging-ips", staging["name"])
		assert.Equal(t, "de/fra", staging["location"])
		assert.Equal(t, int64(1), staging["size"])
		// `ips` is a Computed-only list of strings, so it is a protocol ATTRIBUTE and an
		// omitted one decodes to nil (core_schema.go:102-106).
		assert.Nil(t, staging["ips"])
		// `ip_consumers` is Optional+Computed with an Elem: &schema.Resource, so it is a
		// nested BLOCK (core_schema.go:108-112, :201-205). ReifyNullCollectionBlocks
		// turns a null list block into an empty one, so an omitted block decodes to
		// []any{} rather than nil (toproto6/dynamic_value.go:29-30).
		assert.Equal(t, []any{}, staging["ip_consumers"])

		// The third ip block has no name, which is legal - `name` is Optional on
		// ionoscloud_ipblock - and covers the display-name fallback to the UUID.
		assert.Equal(t, "b1b07384-d9a0-4d1e-8f2a-000000000003", results[2].DisplayName)

		unnamed := decode(t, results[2].Resource, resourceType)
		assert.Equal(t, "b1b07384-d9a0-4d1e-8f2a-000000000003", unnamed["id"])
		assert.Nil(t, unnamed["name"])
	})

	t.Run("filters by name", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{"name": "staging-ips"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "staging-ips", results[0].DisplayName)
	})

	t.Run("filters by location", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{"location": "de/fra"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "staging-ips", results[0].DisplayName)
	})

	// The two values match DIFFERENT stub items, so this pins that the filters are
	// ANDed and that neither field is wired to the other's property.
	t.Run("applies every filter", func(t *testing.T) {
		results := listResults(ctx, t, server, ipBlockListType, listSchema, map[string]string{
			"name":     "prod-ips",
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

// stubIPBlockAPI serves the ip block collection the list resource reads, and returns
// the URL to point IONOS_API_URL at. Every other path 404s on purpose, so an
// unexpected extra request surfaces as an error diagnostic instead of succeeding.
func stubIPBlockAPI(t *testing.T) string {
	t.Helper()

	ipBlocks := ionoscloudsdk.IpBlocks{
		Items: &[]ionoscloudsdk.IpBlock{
			{
				// Result 1: every property the writer reads is set.
				Id: new("b1b07384-d9a0-4d1e-8f2a-000000000001"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Name:     new("prod-ips"),
					Location: new("de/txl"),
					Size:     new(int32(2)),
					Ips:      &[]string{"203.0.113.1", "203.0.113.2"},
					IpConsumers: &[]ionoscloudsdk.IpConsumer{{
						Ip:              new("203.0.113.1"),
						Mac:             new("02:01:0b:0d:8f:b1"),
						NicId:           new("c1c07384-d9a0-4d1e-8f2a-00000000000a"),
						ServerId:        new("c2c07384-d9a0-4d1e-8f2a-00000000000b"),
						ServerName:      new("web-01"),
						DatacenterId:    new("c3c07384-d9a0-4d1e-8f2a-00000000000c"),
						DatacenterName:  new("prod"),
						K8sNodePoolUuid: new("c4c07384-d9a0-4d1e-8f2a-00000000000d"),
						K8sClusterUuid:  new("c5c07384-d9a0-4d1e-8f2a-00000000000e"),
					}},
				},
			},
			{
				// Result 2: only the properties the API always returns, so that the
				// optional ones can be asserted null. `name` stays set here - the
				// display-name and filter assertions need a name to match on.
				Id: new("b1b07384-d9a0-4d1e-8f2a-000000000002"),
				Properties: &ionoscloudsdk.IpBlockProperties{
					Name:     new("staging-ips"),
					Location: new("de/fra"),
					Size:     new(int32(1)),
				},
			},
			{
				// Result 3: no name at all, which `name` being Optional allows. It is
				// what reaches the displayName fallback in mapIPBlock; result 2 cannot,
				// because it keeps `name` set.
				Id: new("b1b07384-d9a0-4d1e-8f2a-000000000003"),
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
		// The request the fetch closure builds is part of what is under test - the stub
		// answers any query string identically, so without these the options are unpinned.
		if got := r.URL.Query().Get("depth"); got != "1" {
			t.Errorf("expected depth=1 on the /ipblocks request, got %q", got)
		}
		// Asserted as a literal rather than as constant.IPBlockLimit: this number is the
		// only thing tying docs/list-resources/ipblock.md to the code, and the point of
		// the assertion is that the fetch asks for more than the 100 the SDK client
		// falls back to for /ipblocks.
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
