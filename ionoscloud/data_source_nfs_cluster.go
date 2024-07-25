package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	nfs "github.com/ionos-cloud/sdk-go-nfs"
)

func dataSourceNFSCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNFSClusterRead,
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
				Description: "The ID (UUID) of the NFS Cluster.",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the NFS Cluster.",
				Optional:    true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "The size of the NFS Cluster.",
				Computed:    true,
			},
			"nfs": {
				Type:        schema.TypeList,
				Description: "The NFS properties of the NFS Cluster.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_version": {
							Type:        schema.TypeString,
							Description: "The minimum version of the NFS.",
							Computed:    true,
						},
					},
				},
			},
			"connections": {
				Type:        schema.TypeList,
				Description: "The connections of the NFS Cluster.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The datacenter ID of the connection.",
							Computed:    true,
						},
						"lan": {
							Type:        schema.TypeString,
							Description: "The LAN of the connection.",
							Computed:    true,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Description: "The IP address of the connection.",
							Computed:    true,
						},
					},
				},
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using the name filter.",
				Default:     false,
				Optional:    true,
			},
		},
	}
}

func dataSourceNFSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).NFSClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	partialMatch := d.Get("partial_match").(bool)
	id := idValue.(string)
	name := nameValue.(string)
	location := d.Get("location").(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the NFS Cluster ID or name"))
	}

	var cluster nfs.ClusterRead
	var err error
	if idOk {
		cluster, _, err = client.GetNFSClusterByID(ctx, id, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the NFS Cluster with ID: %s, error: %w", idValue, err))
		}
	} else {
		clusters, _, err := client.ListNFSClusters(ctx, d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching NFS Clusters: %w", err))
		}

		var results []nfs.ClusterRead
		for _, cl := range *clusters.Items {
			if cl.Properties != nil && cl.Properties.Name != nil && utils.NameMatches(*cl.Properties.Name, name, partialMatch) {
				results = append(results, cl)
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no NFS Cluster found with the specified name: %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one NFS Cluster found with the specified name: %s", name))
		} else {
			cluster = results[0]
		}
	}

	if err = client.SetNFSClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
