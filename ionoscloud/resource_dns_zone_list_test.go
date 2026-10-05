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
// DNS API, through the same muxed provider server that main.go serves.
func TestDNSZoneListResource(t *testing.T) {
	ctx := context.Background()

	t.Setenv("IONOS_API_URL", stubDNSZonesAPI(t))
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
		t.Fatalf("the DNS Zone list resource was not registered; check RawV6Schemas and the SDKv2 resource identity")
	}
	if _, ok := providerSchema.ResourceSchemas[dnsZoneListType]; !ok {
		t.Fatalf("the DNS Zone managed resource is missing from the merged schema")
	}

	identitySchemas, err := server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		t.Fatalf("GetResourceIdentitySchemas: %v", err)
	}
	failOnErrorDiagnostics(t, "GetResourceIdentitySchemas", identitySchemas.Diagnostics)

	identitySchema := identitySchemas.IdentitySchemas[dnsZoneListType]
	if identitySchema == nil {
		t.Fatalf("the SDKv2 DNS Zone resource does not declare a resource identity")
	}

	configureProvider(ctx, t, server, providerSchema.Provider)

	listSchema := providerSchema.ListResourceSchemas[dnsZoneListType]
	resourceType := providerSchema.ResourceSchemas[dnsZoneListType].ValueType()

	t.Run("streams every DNS Zone", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, nil, true)
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}

		assert.Equal(t, "example.com", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "9b2f4c1e-6a3d-4e8b-9c7f-000000000001", identity["id"])

		// Every attribute the writer fills is asserted here, so that mapping a value
		// to the wrong attribute fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "9b2f4c1e-6a3d-4e8b-9c7f-000000000001", resource["id"])
		assert.Equal(t, "example.com", resource["name"])
		assert.Equal(t, "the production zone", resource["description"])
		assert.Equal(t, false, resource["enabled"])
		assert.Equal(t, []any{"ns-ic.ui-dns.com", "ns-ic.ui-dns.de"}, resource["nameservers"])
		assert.Nil(t, resource["timeouts"], "a listed DNS Zone has no timeouts")

		// The second DNS Zone carries only its id, name and state, which pins that the
		// pairing holds past the first result and that what the API left out stays
		// null, or an empty list, instead of a zero value.
		assert.Equal(t, "example.org", results[1].DisplayName)

		secondIdentity := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "9b2f4c1e-6a3d-4e8b-9c7f-000000000002", secondIdentity["id"])

		second := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "9b2f4c1e-6a3d-4e8b-9c7f-000000000002", second["id"])
		assert.Equal(t, "example.org", second["name"])
		assert.Nil(t, second["description"])
		assert.Nil(t, second["enabled"])
		assert.Equal(t, []any{}, second["nameservers"])
	})

	t.Run("filters by name", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{"name": "example.org"}, true)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "example.org", results[0].DisplayName)
	})

	t.Run("filters by description", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{"description": "the production zone"}, true)
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "example.com", results[0].DisplayName)

		// A DNS Zone without a description is matched by an empty field_value.
		undescribed := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{"description": ""}, true)
		if len(undescribed) != 1 {
			t.Fatalf("expected 1 result, got %d", len(undescribed))
		}
		assert.Equal(t, "example.org", undescribed[0].DisplayName)
	})

	t.Run("applies every filter", func(t *testing.T) {
		// The description belongs to example.com, so no DNS Zone matches both filters.
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{"name": "example.org", "description": "the production zone"}, true)
		if len(results) != 0 {
			t.Fatalf("expected no results, got %d", len(results))
		}
	})

	t.Run("identity only", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, nil, false)
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}

		assert.Equal(t, "example.com", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "9b2f4c1e-6a3d-4e8b-9c7f-000000000001", identity["id"])

		for _, result := range results {
			assert.Nil(t, result.Resource, "IncludeResource was not set, so no result carries a resource")
		}
	})

	t.Run("rejects unknown filter fields", func(t *testing.T) {
		// enabled and nameservers are kept out of the allow-list, see
		// ListResourceConfigSchema.
		for _, field := range []string{"nope", "enabled", "nameservers"} {
			listServer, config := listServerAndConfig(t, server, listSchema, map[string]string{field: "value"})

			resp, err := listServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
				TypeName: dnsZoneListType,
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
		for _, field := range []string{"name", "description"} {
			listServer, config := listServerAndConfig(t, server, listSchema, map[string]string{field: "value"})

			resp, err := listServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
				TypeName: dnsZoneListType,
				Config:   &config,
			})
			if err != nil {
				t.Fatalf("ValidateListResourceConfig: %v", err)
			}
			failOnErrorDiagnostics(t, "ValidateListResourceConfig", resp.Diagnostics)
		}
	})
}

// stubDNSZonesAPI serves the DNS Zone collection the list resource reads, and returns
// the URL to point IONOS_API_URL at. A request for that collection must carry none
// of the query parameters ZonesGet knows - filter.state, filter.zoneName, offset and
// limit; any other path gets a 404.
func stubDNSZonesAPI(t *testing.T) string {
	t.Helper()

	// Every item sets Metadata.State: the SDK validates that enum while unmarshalling,
	// and an empty one fails the whole fetch.
	zones := dns.ZoneReadList{
		Items: []dns.ZoneRead{
			{
				Id: "9b2f4c1e-6a3d-4e8b-9c7f-000000000001",
				Metadata: dns.MetadataWithStateNameservers{
					State:       dns.PROVISIONINGSTATE_AVAILABLE,
					Nameservers: []string{"ns-ic.ui-dns.com", "ns-ic.ui-dns.de"},
				},
				Properties: dns.Zone{
					ZoneName:    "example.com",
					Description: new("the production zone"),
					Enabled:     new(false),
				},
			},
			{
				Id: "9b2f4c1e-6a3d-4e8b-9c7f-000000000002",
				Metadata: dns.MetadataWithStateNameservers{
					State: dns.PROVISIONINGSTATE_AVAILABLE,
				},
				Properties: dns.Zone{
					ZoneName: "example.org",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/zones") {
			http.NotFound(w, r)
			return
		}
		for _, param := range []string{"filter.state", "filter.zoneName", "offset", "limit"} {
			if r.URL.Query().Has(param) {
				t.Errorf("expected no %s on the DNS Zone list request, got %q", param, r.URL.Query().Get(param))
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
