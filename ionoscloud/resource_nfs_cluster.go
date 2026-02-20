package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/nfs/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/nfs"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceNFSCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNFSClusterCreate,
		ReadContext:   resourceNFSClusterRead,
		UpdateContext: resourceNFSClusterUpdate,
		DeleteContext: resourceNFSClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNFSClusterImport,
		},
		Schema: map[string]*schema.Schema{
			"location": {
				Type: schema.TypeString,
				Description: fmt.Sprintf(
					"The location of the Network File Storage Cluster. "+
						"Available locations: '%s'", strings.Join(nfs.ValidNFSLocations, ", '"),
				),
				Optional: true,
				ForceNew: true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the Network File Storage Cluster.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Network File Storage Cluster.",
				Required:    true,
			},
			"connections": {
				Type:        schema.TypeList,
				Description: "The network connections for the Network File Storage Cluster.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The datacenter to connect your instance to.",
							Required:    true,
						},
						"lan": {
							Type:        schema.TypeString,
							Description: "The numeric LAN ID to connect your instance to.",
							Required:    true,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Description: "The IP address and subnet for your instance.",
							Required:    true,
						},
					},
				},
			},
			"nfs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_version": {
							Description: "The minimum Network File Storage version",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "The size of the Network File Storage Cluster. Minimum size is 2.",
				Required:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNFSClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).NFSClient

	response, apiResponse, err := client.CreateNFSCluster(ctx, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error creating NFS Cluster: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	clusterID := response.Id
	d.SetId(clusterID)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error checking status for NFS Cluster with ID %v: %s", clusterID, err), nil)
	}
	if err := client.SetNFSClusterData(d, response); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}

func resourceNFSClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).NFSClient

	response, apiResponse, err := client.UpdateNFSCluster(ctx, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error updating NFS Cluster: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error checking status for NFS Cluster %s", err), nil)
	}
	if err := client.SetNFSClusterData(d, response); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}

func resourceNFSClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).NFSClient
	clusterID := d.Id()
	apiResponse, err := client.DeleteNFSCluster(ctx, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error deleting NFS Cluster with ID: %v, error: %s", clusterID, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsClusterDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("deletion check failed for NFS Cluster with ID: %v, error: %s", clusterID, err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}
	return nil
}

func resourceNFSClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).NFSClient
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, utils.ToError(d, "invalid import, expected ID in the format '<location>:<replica_set_id>'", nil)
	}
	location := parts[0]
	id := parts[1]

	err := d.Set("location", location)
	if err != nil {
		return nil, utils.ToError(d, fmt.Sprintf("failed setting location %s: %s", location, err), nil)
	}
	err = d.Set("id", id)
	if err != nil {
		return nil, utils.ToError(d, fmt.Sprintf("failed setting id %s: %s", id, err), nil)
	}

	cluster, err := findCluster(ctx, d, id, location, client)
	if err != nil {
		return nil, utils.ToError(d, fmt.Sprintf("error finding NFS Cluster: %s", err), nil)
	}
	if err := client.SetNFSClusterData(d, cluster); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceNFSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).NFSClient
	cluster, err := findCluster(ctx, d, d.Id(), d.Get("location").(string), client)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error finding NFS Cluster: %s", err), nil)
	}
	if errSetData := client.SetNFSClusterData(d, cluster); errSetData != nil {
		return utils.ToDiags(d, fmt.Sprintf("failed to set NFS Cluster data: %s", errSetData), nil)
	}
	return nil
}

func findCluster(ctx context.Context, d *schema.ResourceData, id, location string, client *nfs.Client) (ionoscloud.ClusterRead, error) {
	cluster, resp, err := client.GetNFSClusterByID(ctx, id, location)
	if err != nil {
		if resp.HttpNotFound() {
			d.SetId("")
			return ionoscloud.ClusterRead{},
				fmt.Errorf("NFS Cluster %s does not exist in %s: %w", id, location, err)
		}
		return ionoscloud.ClusterRead{},
			fmt.Errorf("couldn't find NFS Cluster %s in %s: %w", id, location, err)
	}
	log.Printf("[INFO] Cluster found: %+v", cluster)
	return cluster, nil
}
