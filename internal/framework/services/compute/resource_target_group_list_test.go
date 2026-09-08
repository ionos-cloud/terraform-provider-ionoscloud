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

const targetGroupListType = "ionoscloud_target_group"

// TestTargetGroupListResource drives the ListResource RPC end to end against a stubbed
// Cloud API, through the same muxed provider server that main.go serves.
//
// It covers the parts of a list resource for an SDKv2 managed resource that can only
// fail at runtime: that the mux is happy with the list resource and the managed
// resource coming from different servers, that the framework registers a list resource
// with no framework resource behind it, that the protocol schemas handed over by
// RawV6Schemas convert cleanly, and that the resource model fills the SDKv2 schema
// without a type mismatch.
//
// The shared helpers it calls (muxedProviderServer, decode, listResults, ...) are
// declared once for the package in resource_datacenter_list_test.go.
func TestTargetGroupListResource(t *testing.T) {
	ctx := context.Background()

	t.Setenv("IONOS_API_URL", stubTargetGroupAPI(t))
	t.Setenv("IONOS_TOKEN", "token-for-the-stub")

	server := muxedProviderServer(ctx, t)

	providerSchema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema: %v", err)
	}
	failOnErrorDiagnostics(t, "GetProviderSchema", providerSchema.Diagnostics)

	// The framework only registers a list resource that has no framework resource
	// behind it once RawV6Schemas has supplied both protocol schemas.
	if _, ok := providerSchema.ListResourceSchemas[targetGroupListType]; !ok {
		t.Fatalf("the target group list resource was not registered; check RawV6Schemas and the SDKv2 resource identity")
	}
	if _, ok := providerSchema.ResourceSchemas[targetGroupListType]; !ok {
		t.Fatalf("the target group managed resource is missing from the merged schema")
	}

	identitySchemas, err := server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		t.Fatalf("GetResourceIdentitySchemas: %v", err)
	}
	failOnErrorDiagnostics(t, "GetResourceIdentitySchemas", identitySchemas.Diagnostics)

	identitySchema := identitySchemas.IdentitySchemas[targetGroupListType]
	if identitySchema == nil {
		t.Fatalf("the SDKv2 target group resource does not declare a resource identity")
	}

	configureProvider(ctx, t, server, providerSchema.Provider)

	listSchema := providerSchema.ListResourceSchemas[targetGroupListType]
	resourceType := providerSchema.ResourceSchemas[targetGroupListType].ValueType()

	t.Run("streams every target group", func(t *testing.T) {
		results := listResults(ctx, t, server, targetGroupListType, listSchema, nil)
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}

		assert.Equal(t, "prod-tg", results[0].DisplayName)

		// The identity is the shared identity.Model, so it carries `id` and nothing else.
		identity := decode(t, results[0].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "aa000000-0000-0000-0000-000000000001", identity["id"])
		assert.Len(t, identity, 1, "the target group identity is a lone id")

		// Every attribute the mapper fills is asserted here, so that mapping a value
		// to the wrong attribute fails the test.
		resource := decode(t, results[0].Resource, resourceType)
		assert.Equal(t, "aa000000-0000-0000-0000-000000000001", resource["id"])
		assert.Equal(t, "prod-tg", resource["name"])
		assert.Equal(t, "ROUND_ROBIN", resource["algorithm"])
		assert.Equal(t, "HTTP", resource["protocol"])
		assert.Equal(t, "HTTP1", resource["protocol_version"])
		assert.Equal(t, []any{map[string]any{
			"ip":                   "192.0.2.10",
			"port":                 int64(8080),
			"weight":               int64(5),
			"proxy_protocol":       "v2",
			"health_check_enabled": true,
			"maintenance_enabled":  false,
		}}, resource["targets"])
		// health_check and http_health_check are MaxItems: 1 in the SDKv2 schema, which
		// is still NestingList in the protocol - so they decode as one-element lists,
		// not bare objects.
		assert.Equal(t, []any{map[string]any{
			"check_timeout":  int64(2000),
			"check_interval": int64(1000),
			"retries":        int64(3),
		}}, resource["health_check"])
		assert.Equal(t, []any{map[string]any{
			"path":       "/healthz",
			"method":     "GET",
			"match_type": "STATUS_CODE",
			"response":   "200",
			"regex":      false,
			"negate":     true,
		}}, resource["http_health_check"])
		assert.Nil(t, resource["timeouts"], "a listed target group has no timeouts")

		// The second target group reports only the properties the API always sets, which
		// pins that the pairing holds past the first result and that the properties the
		// API left out stay null instead of turning into zero values.
		assert.Equal(t, "staging-tg", results[1].DisplayName)

		stagingIdentity := decode(t, results[1].Identity.IdentityData, identityType(identitySchema))
		assert.Equal(t, "aa000000-0000-0000-0000-000000000002", stagingIdentity["id"])

		staging := decode(t, results[1].Resource, resourceType)
		assert.Equal(t, "aa000000-0000-0000-0000-000000000002", staging["id"])
		assert.Equal(t, "staging-tg", staging["name"])
		assert.Equal(t, "LEAST_CONNECTION", staging["algorithm"])
		assert.Equal(t, "HTTP", staging["protocol"])
		// protocol_version is a plain ATTRIBUTE, so an omitted one decodes to nil...
		assert.Nil(t, staging["protocol_version"])
		// ...while targets, health_check and http_health_check are nesting BLOCKS
		// (Optional + Computed with an Elem: &schema.Resource{}), and toproto6 runs
		// ReifyNullCollectionBlocks on the way out, which turns a null list block into an
		// empty one. Asserting nil here is the mistake this comment exists to prevent.
		assert.Equal(t, []any{}, staging["targets"])
		assert.Equal(t, []any{}, staging["health_check"])
		assert.Equal(t, []any{}, staging["http_health_check"])
	})

	// One subtest per field in the FilterAttribute allow-list: MatchesFilters returns
	// false for a field_name the mapper forgot to put in its map, so an untested field
	// is a filter that silently matches nothing while the others stay green.
	t.Run("filters by name", func(t *testing.T) {
		results := listResults(ctx, t, server, targetGroupListType, listSchema, map[string]string{"name": "staging-tg"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "staging-tg", results[0].DisplayName)
	})

	t.Run("filters by algorithm", func(t *testing.T) {
		results := listResults(ctx, t, server, targetGroupListType, listSchema, map[string]string{"algorithm": "ROUND_ROBIN"})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		assert.Equal(t, "prod-tg", results[0].DisplayName)
	})

	// Values that each match a different target group, which pins that the filters are
	// ANDed and that no field is wired to another field's property.
	t.Run("applies every filter", func(t *testing.T) {
		results := listResults(ctx, t, server, targetGroupListType, listSchema, map[string]string{
			"name":      "prod-tg",
			"algorithm": "LEAST_CONNECTION",
		})
		if len(results) != 0 {
			t.Fatalf("expected 0 results, got %d", len(results))
		}
	})

	// "protocol" is here alongside a plainly bogus field because it is a resource
	// attribute that is deliberately NOT filterable - it only ever has one value. That
	// decision lives in the filter allow-list and in the docs, so without this case
	// nothing would fail if someone widened the allow-list.
	for _, field := range []string{"nope", "protocol"} {
		t.Run("rejects the "+field+" filter field", func(t *testing.T) {
			listServer, config := listServerAndConfig(t, server, listSchema, map[string]string{field: "value"})

			resp, err := listServer.ValidateListResourceConfig(ctx, &tfprotov6.ValidateListResourceConfigRequest{
				TypeName: targetGroupListType,
				Config:   &config,
			})
			if err != nil {
				t.Fatalf("ValidateListResourceConfig: %v", err)
			}
			if !hasErrorDiagnostic(resp.Diagnostics) {
				t.Fatalf("expected a validation error for the %q filter field", field)
			}
		})
	}
}

// stubTargetGroupAPI serves the target group collection the list resource reads, and
// returns the URL to point IONOS_API_URL at. Every other path 404s on purpose, so an
// unexpected extra request surfaces as an error diagnostic instead of succeeding.
func stubTargetGroupAPI(t *testing.T) string {
	t.Helper()

	targetGroups := ionoscloudsdk.TargetGroups{
		Items: &[]ionoscloudsdk.TargetGroup{
			{
				Id: new("aa000000-0000-0000-0000-000000000001"),
				Properties: &ionoscloudsdk.TargetGroupProperties{
					Name:            new("prod-tg"),
					Algorithm:       new("ROUND_ROBIN"),
					Protocol:        new("HTTP"),
					ProtocolVersion: new("HTTP1"),
					Targets: &[]ionoscloudsdk.TargetGroupTarget{{
						Ip:                 new("192.0.2.10"),
						Port:               new(int32(8080)),
						Weight:             new(int32(5)),
						ProxyProtocol:      new("v2"),
						HealthCheckEnabled: new(true),
						MaintenanceEnabled: new(false),
					}},
					HealthCheck: &ionoscloudsdk.TargetGroupHealthCheck{
						CheckTimeout:  new(int32(2000)),
						CheckInterval: new(int32(1000)),
						Retries:       new(int32(3)),
					},
					HttpHealthCheck: &ionoscloudsdk.TargetGroupHttpHealthCheck{
						Path:      new("/healthz"),
						Method:    new("GET"),
						MatchType: new("STATUS_CODE"),
						Response:  new("200"),
						Regex:     new(false),
						Negate:    new(true),
					},
				},
			},
			{
				// Only the properties the API always sets: name, algorithm and protocol
				// are the three TargetGroupProperties fields without `omitempty`.
				Id: new("aa000000-0000-0000-0000-000000000002"),
				Properties: &ionoscloudsdk.TargetGroupProperties{
					Name:      new("staging-tg"),
					Algorithm: new("LEAST_CONNECTION"),
					Protocol:  new("HTTP"),
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/targetgroups") {
			http.NotFound(w, r)
			return
		}
		// The request the list resource builds is part of what is under test: drop the
		// explicit Limit and the SDK client falls back to limit=100, silently capping the
		// listing below the data source reading the same collection; drop Depth(1) and
		// the API returns links instead of properties.
		// Asserted as a literal, not as constant.TargetGroupLimit, so that changing the
		// constant fails here too - docs/list-resources/target_group.md documents 200 as
		// the ceiling, and nothing else ties the doc to the code.
		if got := r.URL.Query().Get("limit"); got != "200" {
			t.Errorf("expected limit=200 on the targetgroups request, got %q; constant.TargetGroupLimit is %d", got, constant.TargetGroupLimit)
		}
		if got := r.URL.Query().Get("depth"); got != "1" {
			t.Errorf("expected depth=1 on the targetgroups request, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(targetGroups); err != nil {
			t.Errorf("failed to write the stubbed response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return server.URL
}
