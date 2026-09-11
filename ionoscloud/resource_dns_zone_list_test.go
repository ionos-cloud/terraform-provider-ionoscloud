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

// TestDNSZoneListResource drives the ListResource RPC end to end against a stubbed DNS
// API, through the same muxed provider server that main.go serves.
//
// It covers the parts of a list resource for an SDKv2 managed resource that can only
// fail at runtime: that the mux is happy with the list resource and the managed
// resource coming from different servers, that the framework registers a list resource
// with no framework resource behind it, that the protocol schemas handed over by
// RawV6Schemas convert cleanly, and that the state the resource's own writer produces
// survives the round trip through ResourceData.TfTypeResourceState into the list result,
// including the timeouts block identity.MappedItemFromResourceData nulls back out.
//
// Unlike the datacenter and ipblock tests, the stub here answers an sdk-go-bundle
// product: the DNS client is built from the same IONOS_API_URL, but the collection is
// /zones and the response model has value fields rather than pointers.
//
// The shared helpers it calls - muxedProviderServer, configureProvider, listResults,
// listServerAndConfig, decode, identityType and the rest - are declared once for the
// package in resource_datacenter_list_test.go.
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
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}

		assert.Equal(t, "example.com", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Len(t, identity, 1, "a dns zone has no location, so the identity is a lone id")
		assert.Equal(t, "a1a07384-d9a0-4d1e-8f2a-000000000001", identity["id"])

		// Every attribute the writer fills is asserted here, so that writing a value to
		// the wrong key fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "a1a07384-d9a0-4d1e-8f2a-000000000001", resource["id"])
		assert.Equal(t, "example.com", resource["name"])
		assert.Equal(t, "Production zone", resource["description"])
		assert.Equal(t, true, resource["enabled"])
		assert.Equal(t, []any{"ns1.example.net", "ns2.example.net"}, resource["nameservers"])
		assert.Nil(t, resource["timeouts"], "a listed dns zone has no timeouts")

		// The second zone reports only the properties the API always sets, which pins
		// that the pairing holds past the first result and that the properties the API
		// left out stay null instead of turning into zero values. `description` is
		// Optional and `enabled` Optional+Computed - both plain protocol attributes, so
		// an omitted one decodes to nil rather than to "" or false.
		assert.Equal(t, "staging.example.com", results[1].DisplayName)

		stagingIdentity := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "a1a07384-d9a0-4d1e-8f2a-000000000002", stagingIdentity["id"])

		staging := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "a1a07384-d9a0-4d1e-8f2a-000000000002", staging["id"])
		assert.Equal(t, "staging.example.com", staging["name"])
		assert.Nil(t, staging["description"])
		assert.Nil(t, staging["enabled"])
		// `nameservers` is a Computed list attribute, but SetZoneData calls d.Set on it
		// unconditionally, and setting a nil []string on a TypeList materialises an
		// EMPTY list rather than leaving the attribute null. So this one decodes to
		// []any{} where description and enabled decode to nil - the null-vs-zero
		// distinction turns on what the writer does, not only on the schema type.
		assert.Equal(t, []any{}, staging["nameservers"])
	})

	t.Run("filters by name", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{"name": "staging.example.com"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "staging.example.com", results[0].DisplayName)
	})

	t.Run("filters by description", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{"description": "Shared zone"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "dev.example.com", results[0].DisplayName)
	})

	// The two values match DIFFERENT stub items, so this pins that the filters are
	// ANDed and that neither field is wired to the other's property.
	t.Run("applies every filter", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{
			"name":        "example.com",
			"description": "Shared zone",
		})
		if len(results) != 0 {
			t.Fatalf("expected 0 results, got %d", len(results))
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

// stubDNSZoneAPI serves the zone collection the list resource reads, and returns the
// URL to point IONOS_API_URL at. Every other path 404s on purpose, so an unexpected
// extra request surfaces as an error diagnostic instead of succeeding.
func stubDNSZoneAPI(t *testing.T) string {
	t.Helper()

	zones := dns.ZoneReadList{
		Items: []dns.ZoneRead{
			{
				// Result 1: every property the writer reads is set.
				Id: "a1a07384-d9a0-4d1e-8f2a-000000000001",
				Properties: dns.Zone{
					ZoneName:    "example.com",
					Description: new("Production zone"),
					Enabled:     new(true),
				},
				Metadata: dns.MetadataWithStateNameservers{
					State:       dns.PROVISIONINGSTATE_AVAILABLE,
					Nameservers: []string{"ns1.example.net", "ns2.example.net"},
				},
			},
			{
				// Result 2: only the properties the API always returns, so that the
				// optional ones can be asserted null. `name` stays set here - it is
				// Required on the resource, and the display-name and filter assertions
				// need a name to match on.
				Id: "a1a07384-d9a0-4d1e-8f2a-000000000002",
				Properties: dns.Zone{
					ZoneName: "staging.example.com",
				},
				// State is not omitempty and the bundle SDK rejects an empty enum on
				// unmarshal, so even the minimal stub item has to carry one.
				Metadata: dns.MetadataWithStateNameservers{State: dns.PROVISIONINGSTATE_AVAILABLE},
			},
			{
				// Result 3: a second described zone, so the description filter matches
				// something other than result 1 and the AND case can pair a name and a
				// description that belong to different zones.
				Id: "a1a07384-d9a0-4d1e-8f2a-000000000003",
				Properties: dns.Zone{
					ZoneName:    "dev.example.com",
					Description: new("Shared zone"),
					Enabled:     new(false),
				},
				Metadata: dns.MetadataWithStateNameservers{
					State:       dns.PROVISIONINGSTATE_AVAILABLE,
					Nameservers: []string{"ns1.example.net"},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/zones") {
			http.NotFound(w, r)
			return
		}
		// The request the fetch closure builds is part of what is under test, and here
		// what is under test is that it asks for NOTHING: /zones has no depth parameter,
		// no limit is set (the listing is deliberately unpaginated, like the data
		// source), and the name filter is applied in the mapper rather than pushed down
		// with filter.zoneName. The stub answers any query string identically, so
		// without these the fetch options are unpinned.
		for _, param := range []string{"depth", "limit", "offset", "filter.zoneName"} {
			if got := r.URL.Query().Get(param); got != "" {
				t.Errorf("expected no %s on the /zones request, got %q", param, got)
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
