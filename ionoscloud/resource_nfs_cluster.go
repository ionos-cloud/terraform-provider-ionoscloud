package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go-nfs"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/nfs"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

// ValidNFSLocations is a list of valid locations for the Network File Storage Cluster.
var ValidNFSLocations = []string{"de/fra", "de/txl", "qa/de/txl"}

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
					"Available locations: '%s'", strings.Join(ValidNFSLocations, ", '")),
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(ValidNFSLocations, false)),
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
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID: %q, expected ID in the format '<location>:<replica_set_id>'", d.Id())
	}
	location := parts[0]
	if !slices.Contains(ValidNFSLocations, location) {
		return nil, fmt.Errorf("invalid import ID: %q, location must be one of '%s'", d.Id(), strings.Join(ValidNFSLocations, ", '"))
	}
	id := parts[1]

	d.Set("location", location)
	d.Set("id", id)

	cluster, err := findCluster(ctx, d, id, location, client)
	if err != nil {
		return nil, fmt.Errorf("error finding NFS Cluster: %w", err)
	}
	if err := client.SetNFSClusterData(d, cluster); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceNFSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).NFSClient
	cluster, err := findCluster(ctx, d, d.Id(), d.Get("location").(string), client)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error finding NFS Cluster: %w", err))
	}
	if errSetData := client.SetNFSClusterData(d, cluster); errSetData != nil {
		return diag.FromErr(fmt.Errorf("failed to set NFS Cluster data: %w", errSetData))
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
