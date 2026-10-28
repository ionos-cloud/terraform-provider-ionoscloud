package mongov2

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetCluster retrieves a cluster by ID.
func (c *Client) GetCluster(ctx context.Context, clusterID string) (mongov3.ClusterRead, *shared.APIResponse, error) {
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// ListClusters retrieves a list of clusters with optional name and edition filters.
func (c *Client) ListClusters(ctx context.Context, filterName, filterEdition string) (mongov3.ClusterReadList, *shared.APIResponse, error) {
	const pageSize int32 = 100
	var (
		all         []mongov3.ClusterRead
		offset      int32
		apiResponse *shared.APIResponse
	)
	for {
		request := c.sdkClient.ClustersApi.ClustersGet(ctx).Offset(offset).Limit(pageSize)
		if filterName != "" {
			request = request.FilterName(filterName)
		}
		if filterEdition != "" {
			request = request.FilterEdition(mongov3.MongoEdition(filterEdition))
		}
		page, resp, err := request.Execute()
		resp.LogInfo()
		apiResponse = resp
		if err != nil {
			return mongov3.ClusterReadList{}, apiResponse, err
		}
		all = append(all, page.Items...)
		// A short (or empty) page marks the end of the result set. Links.HasNext() is
		// not used: it has been observed to stay populated past the last page.
		if len(page.Items) < int(pageSize) {
			page.Items = all
			return page, apiResponse, nil
		}
		offset += pageSize
	}
}

// CreateCluster creates a new MongoDB v2 cluster.
func (c *Client) CreateCluster(ctx context.Context, clusterCreate mongov3.ClusterCreate) (mongov3.ClusterRead, *shared.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPost(ctx).ClusterCreate(clusterCreate).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

// UpdateCluster updates a cluster using PUT semantics (full replacement).
func (c *Client) UpdateCluster(ctx context.Context, clusterEnsure mongov3.ClusterEnsure, clusterID string) (mongov3.ClusterRead, *shared.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPut(ctx, clusterID).ClusterEnsure(clusterEnsure).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

// DeleteCluster deletes a cluster by ID.
func (c *Client) DeleteCluster(ctx context.Context, clusterID string) (*shared.APIResponse, error) {
	apiResponse, err := c.sdkClient.ClustersApi.ClustersDelete(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsClusterReady checks if the cluster has reached the AVAILABLE state.
func (c *Client) IsClusterReady(ctx context.Context, clusterID string) error {
	cluster, _, err := c.GetCluster(ctx, clusterID)
	if err != nil {
		return backoff.Permanent(err)
	}
	if cluster.Metadata.State != nil {
		tflog.Debug(ctx, "MongoDB v2 cluster state", map[string]any{"state": *cluster.Metadata.State})
		if *cluster.Metadata.State == mongov3.MONGOCLUSTERSTATES_AVAILABLE {
			return nil
		}
		// A cluster in FAILED state can recover to AVAILABLE, so it is treated like any
		// other not-ready state: a retryable error that keeps polling until the timeout.
		return fmt.Errorf("cluster is not ready, current state: %s", *cluster.Metadata.State)
	}
	return backoff.Permanent(fmt.Errorf("can't read cluster state, state is nil"))
}

// IsClusterDeleted checks if the cluster has been deleted (404 response).
func (c *Client) IsClusterDeleted(ctx context.Context, clusterID string) error {
	cluster, apiResponse, err := c.GetCluster(ctx, clusterID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return nil
		}
		return backoff.Permanent(err)
	}
	if cluster.Metadata.State != nil {
		tflog.Debug(ctx, "MongoDB v2 cluster state", map[string]any{"state": *cluster.Metadata.State})
	}
	return fmt.Errorf("cluster with ID %s is not deleted yet", clusterID)
}
