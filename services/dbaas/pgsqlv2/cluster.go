package pgsqlv2

import (
	"context"
	"fmt"
	"log"

	"github.com/cenkalti/backoff/v4"
	pgsqlv2 "github.com/ionos-cloud/pgsqlv2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetCluster retrieves a cluster by its ID.
func (c *Client) GetCluster(ctx context.Context, clusterID string) (pgsqlv2.ClusterRead, *shared.APIResponse, error) {
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// ListClusters retrieves a list of clusters. An optional name filter can be used.
func (c *Client) ListClusters(ctx context.Context, filterName string) (pgsqlv2.ClusterReadList, *shared.APIResponse, error) {
	request := c.sdkClient.ClustersApi.ClustersGet(ctx)
	if filterName != "" {
		request = request.FilterName(filterName)
	}
	clusters, apiResponse, err := request.Execute()
	apiResponse.LogInfo()
	return clusters, apiResponse, err
}

// CreateCluster creates a new cluster.
func (c *Client) CreateCluster(ctx context.Context, clusterCreate pgsqlv2.ClusterCreate) (pgsqlv2.ClusterRead, *shared.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPost(ctx).ClusterCreate(clusterCreate).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

// UpdateCluster updates a cluster by its ID using PUT semantics (full replacement).
func (c *Client) UpdateCluster(ctx context.Context, clusterEnsure pgsqlv2.ClusterEnsure, clusterID string) (pgsqlv2.ClusterRead, *shared.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPut(ctx, clusterID).ClusterEnsure(clusterEnsure).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

// DeleteCluster deletes a cluster by its ID.
func (c *Client) DeleteCluster(ctx context.Context, clusterID string) (*shared.APIResponse, error) {
	apiResponse, err := c.sdkClient.ClustersApi.ClustersDelete(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsClusterReady checks if the cluster has reached the PROVISIONED state.
func (c *Client) IsClusterReady(ctx context.Context, clusterID string) error {
	cluster, _, err := c.GetCluster(ctx, clusterID)
	if err != nil {
		return backoff.Permanent(err)
	}
	log.Printf("[DEBUG] PostgreSQL v2 cluster state: %v", cluster.Metadata.State)
	if cluster.Metadata.State != nil {
		if *cluster.Metadata.State == pgsqlv2.POSTGRESCLUSTERSTATES_PROVISIONED {
			return nil
		}
		// TODO -- Check if for 'Failed' status we keep polling or if we exit immediately.
		if *cluster.Metadata.State == pgsqlv2.POSTGRESCLUSTERSTATES_FAILED {
			return backoff.Permanent(fmt.Errorf("cluster reached FAILED state"))
		}
	}
	return fmt.Errorf("cluster is not ready, current state: %v", cluster.Metadata.State)
}

// IsClusterDeleted checks if the cluster has been deleted (returns 404).
func (c *Client) IsClusterDeleted(ctx context.Context, clusterID string) error {
	cluster, apiResponse, err := c.GetCluster(ctx, clusterID)
	if err != nil {
		if apiResponse != nil && apiResponse.HttpNotFound() {
			return nil
		}
		return backoff.Permanent(err)
	}
	log.Printf("[DEBUG] PostgreSQL v2 cluster state: %v", cluster.Metadata.State)
	return fmt.Errorf("cluster with ID: %s is not deleted yet", clusterID)
}
