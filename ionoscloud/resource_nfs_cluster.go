package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var locations = []string{"de/fra", "de/txl", "qa/de/txl"}

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
				Description: fmt.Sprintf("The location of the Network File Storage Cluster. "+
					"Available locations: '%s'", strings.Join(locations, ", '")),
				Required: true,
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
			"min_version": {
				Type:        schema.TypeString,
				Description: "The minimum Network File Storage version. Current options are '4.2'.",
				Required:    true,
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
	client := meta.(services.SdkBundle).NFSClient

	response, _, err := client.CreateNFSCluster(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating NFS Cluster: %w", err))
	}
	clusterID := *response.Id
	d.SetId(clusterID)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error checking status for NFS Cluster with ID %v: %w", clusterID, err))
	}
	if err := client.SetNFSClusterData(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNFSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, err := resourceNFSClusterImport(ctx, d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading NFS Cluster with ID: %v, error: %w", d.Id(), err))
	}
	return nil
}

func resourceNFSClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).NFSClient

	response, _, err := client.UpdateNFSCluster(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating NFS Cluster: %w", err))
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error checking status for NFS Cluster %w", err))
	}
	if err := client.SetNFSClusterData(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNFSClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).NFSClient
	clusterID := d.Id()
	_, err := client.DeleteNFSCluster(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting NFS Cluster with ID: %v, error: %w", clusterID, err))
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsClusterDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("deletion check failed for NFS Cluster with ID: %v, error: %w", clusterID, err))
	}
	return nil
}

func resourceNFSClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).NFSClient
	clusterID := d.Id()
	cluster, resp, err := client.GetNFSClusterById(ctx, d)
	if err != nil {
		if resp.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("NFS Cluster does not exist, error: %w", err)
		}
		return nil, fmt.Errorf("error importing NFS Cluster with ID: %v, error: %w", clusterID, err)
	}
	log.Printf("[INFO] Cluster found: %+v", cluster)

	if err := client.SetNFSClusterData(d, cluster); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
