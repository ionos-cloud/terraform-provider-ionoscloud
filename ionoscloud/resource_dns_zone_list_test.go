package ionoscloud_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/stretchr/testify/assert"
)

const dnsZoneListType = "ionoscloud_dns_zone"

// TestDNSZoneListResource drives the ListResource RPC end to end against a stubbed
// Cloud DNS API, through the same muxed provider server that main.go serves.
//
// It covers the parts of a list resource for an SDKv2 managed resource that can only
// fail at runtime: that the mux is happy with the list resource and the managed
// resource coming from different servers, that the framework registers a list resource
// with no framework resource behind it, that the protocol schemas handed over by
// RawV6Schemas convert cleanly, and that the state the resource's own writer produces
// survives the round trip through ResourceData.TfTypeResourceState into the list result,
// including the timeouts block identity.MappedItemFromResourceData nulls back out.
func TestDNSZoneListResource(t *testing.T) {
	ctx := context.Background()

	t.Setenv("IONOS_API_URL", stubDNSZoneAPI(t))
	t.Setenv("IONOS_TOKEN", "token-for-the-stub")

	server := muxedProviderServer(ctx, t)

	providerSchema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema: %v", err)
	}
	failOnErrorDiagnostics(t, "GetProviderSchema", providerSchema.Diagnostics)

	// The framework only registers a list resource that has no framework resource
	// behind it once RawV6Schemas has supplied both protocol schemas.
	if _, ok := providerSchema.ListResourceSchemas[dnsZoneListType]; !ok {
		t.Fatalf("the dns zone list resource was not registered; check RawV6Schemas and the SDKv2 resource identity")
	}
	if _, ok := providerSchema.ResourceSchemas[dnsZoneListType]; !ok {
		t.Fatalf("the dns zone managed resource is missing from the merged schema")
	}

	identitySchemas, err := server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		t.Fatalf("GetResourceIdentitySchemas: %v", err)
	}
	failOnErrorDiagnostics(t, "GetResourceIdentitySchemas", identitySchemas.Diagnostics)

	identitySchema := identitySchemas.IdentitySchemas[dnsZoneListType]
	if identitySchema == nil {
		t.Fatalf("the SDKv2 dns zone resource does not declare a resource identity")
	}

	configureProvider(ctx, t, server, providerSchema.Provider)

	listSchema := providerSchema.ListResourceSchemas[dnsZoneListType]
	resourceType := providerSchema.ResourceSchemas[dnsZoneListType].ValueType()

	t.Run("streams every dns zone", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, nil)
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}

		assert.Equal(t, "example.com", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000001", identity["id"])

		// Every attribute the mapper fills is asserted here, so that mapping a value
		// to the wrong attribute fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000001", resource["id"])
		assert.Equal(t, "example.com", resource["name"])
		assert.Equal(t, "the production zone", resource["description"])
		assert.Equal(t, true, resource["enabled"])
		assert.Equal(t, []any{"ns-ic.ui-dns.com", "ns-ic.ui-dns.de"}, resource["nameservers"])
		assert.Nil(t, resource["timeouts"], "a listed dns zone has no timeouts")

		// The second zone reports only the properties the API always sets, which pins
		// that the pairing holds past the first result and that the properties the API
		// left out do not turn into values of the first zone's.
		assert.Equal(t, "staging.example.com", results[1].DisplayName)

		stagingIdentity := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000002", stagingIdentity["id"])

		staging := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "d3b07384-d9a0-4d1e-8f2a-000000000002", staging["id"])
		assert.Equal(t, "staging.example.com", staging["name"])
		assert.Equal(t, "the staging zone", staging["description"])
		assert.Nil(t, staging["enabled"])
		// SetZoneData always sets nameservers, so a zone the API returns none for ends
		// up with an empty list rather than a null one.
		assert.Equal(t, []any{}, staging["nameservers"])
	})

	t.Run("filters by name", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{"name": "example.com"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "example.com", results[0].DisplayName)
	})

	t.Run("filters by description", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{"description": "the staging zone"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "staging.example.com", results[0].DisplayName)
	})

	t.Run("applies every filter", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{
			"name":        "example.com",
			"description": "the staging zone",
		})
		if len(results) != 0 {
			t.Fatalf("expected no result for a name and a description that belong to different zones, got %d", len(results))
		}
	})

	t.Run("rejects unknown filter fields", func(t *testing.T) {
		listServer, config := listServerAndConfig(t, server, listSchema, map[string]string{"nope": "value"})

		resp, err := listServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
			TypeName: dnsZoneListType,
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

// stubDNSZoneAPI serves the zone collection the list resource reads, asserts the query
// the fetch closure builds, and returns the URL to point IONOS_API_URL at.
func stubDNSZoneAPI(t *testing.T) string {
	t.Helper()

	zones := dns.ZoneReadList{
		Items: []dns.ZoneRead{
			{
				Id:   "d3b07384-d9a0-4d1e-8f2a-000000000001",
				Type: "zone",
				Href: "/zones/d3b07384-d9a0-4d1e-8f2a-000000000001",
				// State is an enum the SDK validates while unmarshalling, so every item
				// has to carry one or the fetch itself fails.
				Metadata: dns.MetadataWithStateNameservers{
					State:       dns.PROVISIONINGSTATE_AVAILABLE,
					Nameservers: []string{"ns-ic.ui-dns.com", "ns-ic.ui-dns.de"},
				},
				Properties: dns.Zone{
					ZoneName:    "example.com",
					Description: new("the production zone"),
					Enabled:     new(true),
				},
			},
			{
				Id:   "d3b07384-d9a0-4d1e-8f2a-000000000002",
				Type: "zone",
				Href: "/zones/d3b07384-d9a0-4d1e-8f2a-000000000002",
				Metadata: dns.MetadataWithStateNameservers{
					State: dns.PROVISIONINGSTATE_AVAILABLE,
				},
				Properties: dns.Zone{
					ZoneName:    "staging.example.com",
					Description: new("the staging zone"),
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/zones") {
			http.NotFound(w, r)
			return
		}
		// The fetch closure sets no option at all: the collection takes no depth, the
		// provider sends no limit or offset, and the name filter is applied in the
		// mapper rather than pushed down to filter.zoneName.
		for _, param := range []string{"depth", "limit", "offset", "filter.zoneName"} {
			if got := r.URL.Query().Get(param); got != "" {
				t.Errorf("expected the dns zone listing to send no %s, got %q", param, got)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(zones); err != nil {
			t.Errorf("failed to write the stubbed response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return server.URL
}
