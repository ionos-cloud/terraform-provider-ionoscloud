package nfs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk "github.com/ionos-cloud/sdk-go-nfs"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// GetNFSClusterById returns a cluster given an ID
func (c *Client) GetNFSClusterById(ctx context.Context, id string, location string) (sdk.ClusterRead, *sdk.APIResponse, error) {
	cluster, apiResponse, err := c.Location(location).sdkClient.ClustersApi.ClustersFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// ListNFSClusters returns a list of all clusters
func (c *Client) ListNFSClusters(ctx context.Context, d *schema.ResourceData) (sdk.ClusterReadList, *sdk.APIResponse, error) {
	clusters, apiResponse, err := c.Location(d.Get("location").(string)).sdkClient.ClustersApi.ClustersGet(ctx).Execute()
	apiResponse.LogInfo()
	return clusters, apiResponse, err
}

// DeleteNFSCluster deletes a cluster given an ID
func (c *Client) DeleteNFSCluster(ctx context.Context, d *schema.ResourceData) (*sdk.APIResponse, error) {
	apiResponse, err := c.Location(d.Get("location").(string)).sdkClient.ClustersApi.ClustersDelete(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// UpdateNFSCluster updates a cluster given an ID or creates a new one if it doesn't exist
func (c *Client) UpdateNFSCluster(ctx context.Context, d *schema.ResourceData) (sdk.ClusterRead, *sdk.APIResponse, error) {
	cluster, apiResponse, err := c.Location(d.Get("location").(string)).sdkClient.ClustersApi.ClustersPut(ctx, d.Id()).
		ClusterEnsure(*setClusterPutRequest(d)).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// CreateNFSCluster creates a new cluster
func (c *Client) CreateNFSCluster(ctx context.Context, d *schema.ResourceData) (sdk.ClusterRead, *sdk.APIResponse, error) {
	cluster, apiResponse, err := c.Location(d.Get("location").(string)).sdkClient.ClustersApi.ClustersPost(ctx).
		ClusterCreate(*setClusterPostRequest(d)).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// SetNFSClusterData sets the data of the cluster in the terraform resource
func (c *Client) SetNFSClusterData(d *schema.ResourceData, cluster sdk.ClusterRead) error {
	d.SetId(*cluster.Id)

	if cluster.Properties == nil {
		return fmt.Errorf("expected properties in the response for the NFS Cluster with ID: %s, but received 'nil' instead", *cluster.Id)
	}

	if cluster.Metadata == nil {
		return fmt.Errorf("expected metadata in the response for the NFS Cluster with ID: %s, but received 'nil' instead", *cluster.Id)
	}

	if cluster.Properties.Name != nil {
		if err := d.Set("name", *cluster.Properties.Name); err != nil {
			return err
		}
	}

	if cluster.Properties.Size != nil {
		if err := d.Set("size", *cluster.Properties.Size); err != nil {
			return err
		}
	}

	if cluster.Properties.Nfs != nil && cluster.Properties.Nfs.MinVersion != nil {
		if err := d.Set("min_version", *cluster.Properties.Nfs.MinVersion); err != nil {
			return err
		}
	}

	if cluster.Properties.Connections != nil {
		var connections []map[string]interface{}
		for _, connection := range *cluster.Properties.Connections {
			connectionData := map[string]interface{}{
				"datacenter_id": *connection.DatacenterId,
				"lan":           *connection.Lan,
				"ip_address":    *connection.IpAddress,
			}
			connections = append(connections, connectionData)
		}

		if err := d.Set("connections", connections); err != nil {
			return fmt.Errorf("error setting connections for the NFS Cluster with ID %s: %w", *cluster.Id, err)
		}
	}

	return nil
}

// IsClusterReady checks if the cluster is ready
func (c *Client) IsClusterReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterID := d.Id()
	cluster, _, err := c.GetNFSClusterById(ctx, "", "")
	if err != nil {
		return true, fmt.Errorf("status check failed for Cluster ID: %v, error: %w", clusterID, err)
	}

	if cluster.Metadata == nil || cluster.Metadata.Status == nil {
		return false, fmt.Errorf("metadata or status is empty for Cluster ID: %v", clusterID)
	}

	log.Printf("[INFO] state of the cluster with ID %s is: %s ", clusterID, *cluster.Metadata.Status)
	return strings.EqualFold(*cluster.Metadata.Status, constant.Available), nil
}

// IsClusterDeleted checks if the cluster is deleted
func (c *Client) IsClusterDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterID := d.Id()
	_, apiResponse, err := c.GetNFSClusterById(ctx, "", "")
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("check failed for Cluster deletion status, ID: %v, error: %w", clusterID, err)
	}
	return false, nil
}

func setClusterPostRequest(d *schema.ResourceData) *sdk.ClusterCreate {
	return sdk.NewClusterCreate(setClusterConfig(d))
}

func setClusterPutRequest(d *schema.ResourceData) *sdk.ClusterEnsure {
	clusterID := d.Id()
	cluster := setClusterConfig(d)

	return sdk.NewClusterEnsure(clusterID, cluster)
}

func setClusterConfig(d *schema.ResourceData) sdk.Cluster {
	name := d.Get("name").(string)
	size := int32(d.Get("size").(int))
	minVersion := d.Get("min_version").(string)
	connectionsRaw := d.Get("connections").([]interface{})

	connections := make([]sdk.ClusterConnections, len(connectionsRaw))
	for i, conn := range connectionsRaw {
		connData := conn.(map[string]interface{})
		datacenterID := connData["datacenter_id"].(string)
		lan := connData["lan"].(string)
		ipAddress := connData["ip_address"].(string)

		connectionObj := sdk.ClusterConnections{
			DatacenterId: &datacenterID,
			Lan:          &lan,
			IpAddress:    &ipAddress,
		}

		connections[i] = connectionObj
	}

	return sdk.Cluster{
		Name: &name,
		Size: &size,
		Nfs: &sdk.ClusterNfs{
			MinVersion: &minVersion,
		},
		Connections: &connections,
	}
}
