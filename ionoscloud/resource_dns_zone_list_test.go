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
// Unlike the datacenter and ipblock list resources, ionoscloud_dns_zone is on an
// sdk-go-bundle product: the writer under test here is dnsservice.Client.SetZoneData,
// a method on the DNS service client in services/dns/zone.go, not a package-level
// function in package ionoscloud.
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

		assert.Equal(t, "prod.example.com", results[0].DisplayName)

		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Len(t, identity, 1, "a dns zone has no location, so its identity is a lone id")
		assert.Equal(t, "e1e07384-d9a0-4d1e-8f2a-000000000001", identity["id"])

		// Every attribute the writer fills is asserted here, so that writing a value to
		// the wrong key fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "e1e07384-d9a0-4d1e-8f2a-000000000001", resource["id"])
		assert.Equal(t, "prod.example.com", resource["name"])
		assert.Equal(t, "the production zone", resource["description"])
		assert.Equal(t, true, resource["enabled"])
		assert.Equal(t, []any{"ns-ic.ui-dns.com", "ns-ic.ui-dns.de"}, resource["nameservers"])
		assert.Nil(t, resource["timeouts"], "a listed dns zone has no timeouts")

		// The second zone pins that the pairing holds past the first result: a different
		// description on a different name is what the AND filter subtest below relies on.
		assert.Equal(t, "staging.example.com", results[1].DisplayName)

		stagingIdentity := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "e1e07384-d9a0-4d1e-8f2a-000000000002", stagingIdentity["id"])

		staging := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "e1e07384-d9a0-4d1e-8f2a-000000000002", staging["id"])
		assert.Equal(t, "staging.example.com", staging["name"])
		assert.Equal(t, "the staging zone", staging["description"])
		assert.Equal(t, false, staging["enabled"])

		// The third zone reports only the properties the API always sets, which pins
		// that the properties the API left out stay null instead of turning into zero
		// values.
		assert.Equal(t, "minimal.example.com", results[2].DisplayName)

		minimal := decode(t, results[2].Resource, resourceType)
		assert.Equal(t, "e1e07384-d9a0-4d1e-8f2a-000000000003", minimal["id"])
		assert.Equal(t, "minimal.example.com", minimal["name"])
		// SetZoneData only writes description and enabled when the API returned them, so
		// an omitted one is left null rather than written as "" or false.
		assert.Nil(t, minimal["description"])
		assert.Nil(t, minimal["enabled"])
		// nameservers is different, and the schema type does not explain it: it is a
		// Computed-only list, so it is a protocol ATTRIBUTE and an unset one would
		// decode to nil - but SetZoneData d.Set()s it unconditionally, and setting a nil
		// slice materialises an EMPTY list. Half the answer is what the writer does.
		assert.Equal(t, []any{}, minimal["nameservers"])
	})

	t.Run("filters by name", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{"name": "staging.example.com"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "staging.example.com", results[0].DisplayName)
	})

	t.Run("filters by description", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{"description": "the production zone"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "prod.example.com", results[0].DisplayName)
	})

	// The two values match DIFFERENT stub items, so this pins that the filters are
	// ANDed and that neither field is wired to the other's property.
	t.Run("applies every filter", func(t *testing.T) {
		results := listResults(ctx, t, server, dnsZoneListType, listSchema, map[string]string{
			"name":        "prod.example.com",
			"description": "the staging zone",
		})
		if len(results) != 0 {
			t.Fatalf("expected 0 results, got %d", len(results))
		}
	})

	// `enabled` and `nameservers` are real attributes of the resource that are
	// deliberately left out of the allow-list, because MatchesFilters only compares
	// strings. Looping over them here means widening the allow-list cannot pass
	// unnoticed.
	t.Run("rejects unknown filter fields", func(t *testing.T) {
		for _, field := range []string{"nope", "enabled", "nameservers"} {
			listServer, config := listServerAndConfig(t, server, listSchema, map[string]string{field: "value"})

			resp, err := listServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
				TypeName: dnsZoneListType,
				Config:   &config,
			})
			if err != nil {
				t.Fatalf("ValidateListResourceConfig(%s): %v", field, err)
			}
			if !hasErrorDiagnostic(resp.Diagnostics) {
				t.Fatalf("expected a validation error for the filter field %q", field)
			}
		}
	})
}

// stubDNSZoneAPI serves the zone collection the list resource reads, and returns the URL
// to point IONOS_API_URL at. Every other path 404s on purpose, so an unexpected extra
// request surfaces as an error diagnostic instead of succeeding.
//
// services/dns/client.go builds its configuration straight from clientOptions.Endpoint
// and there is no IONOS_API_URL_DNS, so IONOS_API_URL alone is what redirects the DNS
// client here.
func stubDNSZoneAPI(t *testing.T) string {
	t.Helper()

	zones := dns.ZoneReadList{
		Items: []dns.ZoneRead{
			{
				// Result 1: every property the writer reads is set.
				Id: "e1e07384-d9a0-4d1e-8f2a-000000000001",
				Properties: dns.Zone{
					ZoneName:    "prod.example.com",
					Description: new("the production zone"),
					Enabled:     new(true),
				},
				Metadata: dns.MetadataWithStateNameservers{
					// State is not omitempty and ProvisioningState validates itself
					// while unmarshalling, so leaving it unset fails the fetch with
					// " is not a valid ProvisioningState" rather than an assertion.
					State:       dns.PROVISIONINGSTATE_AVAILABLE,
					Nameservers: []string{"ns-ic.ui-dns.com", "ns-ic.ui-dns.de"},
				},
			},
			{
				// Result 2: a second fully-populated zone whose name and description both
				// differ from result 1's, which is what the per-field and AND filter
				// subtests discriminate on.
				Id: "e1e07384-d9a0-4d1e-8f2a-000000000002",
				Properties: dns.Zone{
					ZoneName:    "staging.example.com",
					Description: new("the staging zone"),
					Enabled:     new(false),
				},
				Metadata: dns.MetadataWithStateNameservers{
					State:       dns.PROVISIONINGSTATE_AVAILABLE,
					Nameservers: []string{"ns-ic.ui-dns.org"},
				},
			},
			{
				// Result 3: only the properties the API always returns, so that the
				// optional ones can be asserted null. `name` stays set - it is Required
				// on the resource and non-nullable in the SDK model, and the display-name
				// assertion needs it.
				Id: "e1e07384-d9a0-4d1e-8f2a-000000000003",
				Properties: dns.Zone{
					ZoneName: "minimal.example.com",
				},
				Metadata: dns.MetadataWithStateNameservers{
					State: dns.PROVISIONINGSTATE_AVAILABLE,
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/zones") {
			http.NotFound(w, r)
			return
		}
		// The request the fetch closure builds is part of what is under test, and this
		// one deliberately builds the bare collection GET: ApiZonesGetRequest has no
		// depth parameter, the fetch asks for no limit (matching the data source) and
		// pushes no filter down. Asserting the ABSENCE of all four is the only thing
		// pinning that - the stub answers any query string identically, so a limit or a
		// pushed-down filter added later would otherwise go unnoticed.
		for _, param := range []string{"depth", "limit", "offset", "filter.zoneName", "filter.state"} {
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
