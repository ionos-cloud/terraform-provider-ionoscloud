package mongov2

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
	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	mongov2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mongov2"
)

var (
	_ list.ListResource              = (*clusterResource)(nil)
	_ list.ListResourceWithConfigure = (*clusterResource)(nil)
)

type clusterWithLocation struct {
	Cluster  mongov3.ClusterRead
	Location string
}

// NewClusterListResource creates a new list resource for mongo_cluster_v2.
func NewClusterListResource() list.ListResource {
	return &clusterResource{}
}

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
			var filters []identity.Filter
			req.Config.GetAttribute(ctx, path.Root(identity.FiltersKey), &filters)

			nameFilter := identity.FilterValue(filters, "name")

			locations := mongov2service.AvailableLocations()
			if loc := identity.FilterValue(filters, "location"); loc != "" {
				locations = []string{loc}
			}

			var all []clusterWithLocation
			for _, loc := range locations {
				client, err := r.bundle.NewMongoDBV2Client(ctx, loc)
				if err != nil {
					return nil, fmt.Errorf("failed to create MongoDB V2 client, location: %s, error: %w", loc, err)
				}
				result, _, err := client.ListClusters(ctx, nameFilter, "")
				if err != nil {
					return nil, fmt.Errorf("failed to list MongoDB V2 clusters, location: %s, error: %w", loc, err)
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

func (r *clusterResource) mapCluster(ctx context.Context, includeResource bool, filters []identity.Filter, item clusterWithLocation) (*identity.MappedItem, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	c := item.Cluster
	name := clusterName(c.Properties)

	if !identity.MatchesFilters(map[string]string{
		"name":     name,
		"location": item.Location,
	}, filters) {
		return nil, nil
	}

	mapped := &identity.MappedItem{
		DisplayName: name,
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
	diagnostics.Append(mapClusterResponseToModel(ctx, &c, model)...)
	mapped.Resource = model
	return mapped, diagnostics
}

// clusterName returns the cluster name from whichever edition variant is set.
func clusterName(props mongov3.Cluster) string {
	if props.EnterpriseCluster != nil {
		return props.EnterpriseCluster.Name
	}
	if props.BusinessCluster != nil {
		return props.BusinessCluster.Name
	}
	return ""
}
