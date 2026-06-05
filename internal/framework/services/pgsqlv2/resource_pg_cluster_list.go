package pgsqlv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	pgsqlv2sdk "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	pgsqlv2Service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/pgsqlv2"
)

var (
	_ list.ListResource              = (*clusterResource)(nil)
	_ list.ListResourceWithConfigure = (*clusterResource)(nil)
)

// clusterWithLocation pairs a ClusterRead with the location it was fetched from,
// since the API does not include location in the response.
type clusterWithLocation struct {
	Cluster  pgsqlv2sdk.ClusterRead
	Location string
}

// NewClusterListResource creates a new list resource for pg_cluster_v2.
func NewClusterListResource() list.ListResource {
	return &clusterResource{}
}

// ListResourceConfigSchema returns the schema for the list resource config block.
func (r *clusterResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			identity.FiltersKey: identity.FilterAttribute("name", "location"),
		},
	}
}

// List fetches clusters from all regional endpoints and streams the results.
// If a location filter is set, only that endpoint is queried.
func (r *clusterResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	identity.StreamList(ctx, stream, req,
		func(ctx context.Context) ([]clusterWithLocation, error) {
			// Read filters early to skip unnecessary regional API calls.
			var filters []identity.Filter
			req.Config.GetAttribute(ctx, path.Root(identity.FiltersKey), &filters) //nolint:errcheck

			nameFilter := identity.FilterValue(filters, "name")
			locationFilter := identity.FilterValue(filters, "location")

			locations := pgsqlv2Service.AvailableLocations()
			if locationFilter != "" {
				locations = []string{locationFilter}
			}

			var all []clusterWithLocation
			for _, loc := range locations {
				client, err := r.bundle.NewPgSQLV2Client(ctx, loc)
				if err != nil {
					return nil, fmt.Errorf("failed to create client for location %q: %w", loc, err)
				}
				result, _, err := client.ListClusters(ctx, nameFilter)
				if err != nil {
					return nil, fmt.Errorf("failed to list clusters in location %q: %w", loc, err)
				}
				for _, c := range result.Items {
					all = append(all, clusterWithLocation{Cluster: c, Location: loc})
				}
			}
			return all, nil
		},
		r.mapCluster,
	)
}

// mapCluster maps a clusterWithLocation to an identity.MappedItem, or returns nil to skip it.
func (r *clusterResource) mapCluster(_ context.Context, includeResource bool, filters []identity.Filter, item clusterWithLocation) (*identity.MappedItem, diag.Diagnostics) {
	c := item.Cluster

	if !identity.MatchesFilters(map[string]string{
		"name":     c.Properties.Name,
		"location": item.Location,
	}, filters) {
		return nil, nil
	}

	mapped := &identity.MappedItem{
		DisplayName: c.Properties.Name,
		Identity: &clusterIdentityModel{
			ID:       types.StringValue(c.Id),
			Location: types.StringValue(item.Location),
		},
	}

	if !includeResource {
		return mapped, nil
	}

	model := &clusterResourceModel{
		ID:       types.StringValue(c.Id),
		Location: types.StringValue(item.Location),
		Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType,
			"update": types.StringType,
			"delete": types.StringType,
		})},
	}
	mapClusterResponseToModel(&c, model)
	mapped.Resource = model
	return mapped, nil
}
