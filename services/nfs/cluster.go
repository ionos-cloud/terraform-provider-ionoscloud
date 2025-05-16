package nfs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk "github.com/ionos-cloud/sdk-go-bundle/products/nfs/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// GetNFSClusterByID returns a cluster given an ID
func (c *Client) GetNFSClusterByID(ctx context.Context, id string, location string) (sdk.ClusterRead, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.NFS, location)
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// ListNFSClusters returns a list of all clusters
func (c *Client) ListNFSClusters(ctx context.Context, d *schema.ResourceData) (sdk.ClusterReadList, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.NFS, d.Get("location").(string))
	clusters, apiResponse, err := c.sdkClient.ClustersApi.ClustersGet(ctx).Execute()
	apiResponse.LogInfo()
	return clusters, apiResponse, err
}

// DeleteNFSCluster deletes a cluster given an ID
func (c *Client) DeleteNFSCluster(ctx context.Context, d *schema.ResourceData) (*shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.NFS, d.Get("location").(string))
	apiResponse, err := c.sdkClient.ClustersApi.ClustersDelete(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// UpdateNFSCluster updates a cluster given an ID or creates a new one if it doesn't exist
func (c *Client) UpdateNFSCluster(ctx context.Context, d *schema.ResourceData) (sdk.ClusterRead, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.NFS, d.Get("location").(string))
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersPut(ctx, d.Id()).
		ClusterEnsure(*setClusterPutRequest(d)).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// CreateNFSCluster creates a new cluster
func (c *Client) CreateNFSCluster(ctx context.Context, d *schema.ResourceData) (sdk.ClusterRead, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.NFS, d.Get("location").(string))
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersPost(ctx).
		ClusterCreate(*setClusterPostRequest(d)).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// SetNFSClusterData sets the data of the cluster in the terraform resource
func (c *Client) SetNFSClusterData(d *schema.ResourceData, cluster sdk.ClusterRead) error {
	d.SetId(cluster.Id)

	if err := d.Set("name", cluster.Properties.Name); err != nil {
		return err
	}

	if cluster.Properties.Size != nil {
		if err := d.Set("size", *cluster.Properties.Size); err != nil {
			return err
		}
	}

	if cluster.Properties.Nfs != nil {
		nfs := []map[string]interface{}{
			{
				"min_version": *cluster.Properties.Nfs.MinVersion,
			},
		}
		if err := d.Set("nfs", nfs); err != nil {
			return fmt.Errorf("error setting nfs data for the NFS Cluster with ID %s: %w", cluster.Id, err)
		}
	}

	var connections []map[string]interface{}
	for _, connection := range cluster.Properties.Connections {
		connectionData := map[string]interface{}{
			"datacenter_id": connection.DatacenterId,
			"lan":           connection.Lan,
			"ip_address":    connection.IpAddress,
		}
		connections = append(connections, connectionData)

		if err := d.Set("connections", connections); err != nil {
			return fmt.Errorf("error setting connections for the NFS Cluster with ID %s: %w", cluster.Id, err)
		}
	}

	return nil
}

// IsClusterReady checks if the cluster is ready
func (c *Client) IsClusterReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterID := d.Id()
	location := d.Get("location").(string)
	cluster, _, err := c.GetNFSClusterByID(ctx, clusterID, location)
	if err != nil {
		return false, fmt.Errorf("failed checking if Cluster %s from %s is ready: %w", clusterID, location, err)
	}

	log.Printf("[INFO] state of the cluster with ID %s is: %s", clusterID, cluster.Metadata.Status)
	if utils.IsStateFailed(cluster.Metadata.Status) {
		return false, fmt.Errorf("cluster %s is in a failed state", d.Id())
	}
	return strings.EqualFold(cluster.Metadata.Status, constant.Available), nil
}

// IsClusterDeleted checks if the cluster is deleted
func (c *Client) IsClusterDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterID := d.Id()
	location := d.Get("location").(string)
	_, apiResponse, err := c.GetNFSClusterByID(ctx, clusterID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("failed checking if Cluster %s from %s has been deleted: %w", clusterID, location, err)
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

	nfsRaw := d.Get("nfs").([]interface{})
	nfs := sdk.ClusterNfs{}
	var minVersion string
	if len(nfsRaw) > 0 && nfsRaw[0] != nil {
		nfsData := nfsRaw[0].(map[string]interface{})
		minVersion = nfsData["min_version"].(string)

		nfs = sdk.ClusterNfs{
			MinVersion: &minVersion,
		}
	}

	connectionsRaw := d.Get("connections").([]interface{})
	connections := make([]sdk.ClusterConnections, 1)
	if len(connectionsRaw) > 0 && connectionsRaw[0] != nil {
		connData := connectionsRaw[0].(map[string]interface{})
		datacenterID := connData["datacenter_id"].(string)
		lan := connData["lan"].(string)
		ipAddress := connData["ip_address"].(string)

		connections[0] = sdk.ClusterConnections{
			DatacenterId: datacenterID,
			Lan:          lan,
			IpAddress:    ipAddress,
		}
	}

	return sdk.Cluster{
		Name:        name,
		Size:        &size,
		Nfs:         &nfs,
		Connections: connections,
	}
}
