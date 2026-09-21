package mariadbv2

import (
	"context"
	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	mariadbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetCluster retrieves a cluster by ID.
func (c *Client) GetCluster(ctx context.Context, clusterID string) (mariadbv3.ClusterRead, *shared.APIResponse, error) {
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// ListClusters retrieves a list of clusters with optional name filter.
func (c *Client) ListClusters(ctx context.Context, filterName string) (mariadbv3.ClusterReadList, *shared.APIResponse, error) {
	request := c.sdkClient.ClustersApi.ClustersGet(ctx)
	if filterName != "" {
		request = request.FilterName(filterName)
	}
	clusters, apiResponse, err := request.Execute()
	apiResponse.LogInfo()
	return clusters, apiResponse, err
}

// CreateCluster creates a new MariaDB v2 cluster.
func (c *Client) CreateCluster(ctx context.Context, clusterCreate mariadbv3.ClusterCreate) (mariadbv3.ClusterRead, *shared.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.ClustersApi.ClustersPost(ctx).ClusterCreate(clusterCreate).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

// UpdateCluster updates a cluster using PUT semantics (full replacement).
func (c *Client) UpdateCluster(ctx context.Context, clusterEnsure mariadbv3.ClusterEnsure, clusterID string) (mariadbv3.ClusterRead, *shared.APIResponse, error) {
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
		tflog.Debug(ctx, "MariaDB v2 cluster state", map[string]any{"state": *cluster.Metadata.State})
		if *cluster.Metadata.State == mariadbv3.MARIADBCLUSTERSTATES_AVAILABLE {
			return nil
		}
		if *cluster.Metadata.State == mariadbv3.MARIADBCLUSTERSTATES_FAILED {
			msg := "cluster entered FAILED state"
			if cluster.Metadata.StatusMessage != nil {
				msg = fmt.Sprintf("cluster entered FAILED state: %s", *cluster.Metadata.StatusMessage)
			}
			return backoff.Permanent(errors.New(msg))
		}
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
		tflog.Debug(ctx, "MariaDB v2 cluster state", map[string]any{"state": *cluster.Metadata.State})
	}
	return fmt.Errorf("cluster with ID %s is not deleted yet", clusterID)
}
