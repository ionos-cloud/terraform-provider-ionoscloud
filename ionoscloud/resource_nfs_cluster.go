package ionoscloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sdk "github.com/ionos-cloud/sdk-go-nfs"
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
				// Affects the Host of the SDK
				Type: schema.TypeString,
				Description: fmt.Sprintf("The location of the Network File Storage Cluster. "+
					"Available locations: %s", strings.Join(locations, ", ")),
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(locations, false)),
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
				Description: "The minimum Network File Storage version.",
				Required:    true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "The size of the Network File Storage Cluster.",
				Required:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceNFSClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sdk.APIClient)

	cluster := sdk.ClusterCreate{
		Properties: &sdk.Cluster{
			Name: sdk.ToPtr(d.Get("name").(string)),
			Connections: &[]sdk.ClusterConnections{
				{
					DatacenterId: sdk.ToPtr(d.Get("connections.0.datacenter_id").(string)),
					Lan:          sdk.ToPtr(d.Get("connections.0.lan").(string)),
					IpAddress:    sdk.ToPtr(d.Get("connections.0.ip_address").(string)),
				},
			},
			Nfs: &sdk.ClusterNfs{
				MinVersion: sdk.ToPtr(d.Get("min_version").(string)),
			},
			Size: sdk.ToPtr(int32(d.Get("size").(int))),
		},
	}

	resp, _, err := client.ClustersApi.ClustersPost(ctx).ClusterCreate(cluster).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.Id)

	return resourceNFSClusterRead(ctx, d, meta)
}

func resourceNFSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sdk.APIClient)

	clusterID := d.Id()
	cluster, _, err := client.ClustersApi.ClustersFindById(ctx, clusterID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", cluster.Properties.Name)
	d.Set("connections", []interface{}{
		map[string]interface{}{
			"datacenter_id": *(*cluster.Properties.Connections)[0].DatacenterId,
			"lan":           (*cluster.Properties.Connections)[0].Lan,
			"ip_address":    (*cluster.Properties.Connections)[0].IpAddress,
		},
	})
	d.Set("min_version", cluster.Properties.Nfs.MinVersion)
	d.Set("size", cluster.Properties.Size)

	return nil
}

func resourceNFSClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sdk.APIClient)

	clusterID := d.Id()

	cluster := sdk.ClusterEnsure{
		Properties: &sdk.Cluster{
			Name: sdk.ToPtr(d.Get("name").(string)),
			Nfs: &sdk.ClusterNfs{
				MinVersion: sdk.ToPtr(d.Get("min_version").(string)),
			},
			Size: sdk.ToPtr(int32(d.Get("size").(int))),
		},
	}

	_, _, err := client.ClustersApi.ClustersPut(ctx, clusterID).ClusterEnsure(cluster).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNFSClusterRead(ctx, d, meta)
}

func resourceNFSClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sdk.APIClient)

	clusterID := d.Id()
	_, err := client.ClustersApi.ClustersDelete(ctx, clusterID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNFSClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*sdk.APIClient)
	clusterID := d.Id()
	cluster, _, err := client.ClustersApi.ClustersFindById(ctx, clusterID).Execute()
	if err != nil {
		return nil, err
	}

	if err := d.Set("name", cluster.Properties.Name); err != nil {
		return nil, err
	}
	if err := d.Set("connections", []interface{}{
		map[string]interface{}{
			"datacenter_id": (*cluster.Properties.Connections)[0].DatacenterId,
			"lan":           (*cluster.Properties.Connections)[0].Lan,
			"ip_address":    (*cluster.Properties.Connections)[0].IpAddress,
		},
	}); err != nil {
		return nil, err
	}
	if err := d.Set("min_version", cluster.Properties.Nfs.MinVersion); err != nil {
		return nil, err
	}
	if err := d.Set("size", cluster.Properties.Size); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
