package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	nfs2 "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/nfs"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	"github.com/ionos-cloud/sdk-go-bundle/products/nfs/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func dataSourceNFSCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNFSClusterRead,
		Schema: map[string]*schema.Schema{
			"location": {
				Type: schema.TypeString,
				Description: fmt.Sprintf("The location of the Network File Storage Cluster. "+
					"Available locations: '%s'", strings.Join(nfs2.ValidNFSLocations, ", '")),
				Required: true,
				ForceNew: true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The ID (UUID) of the NFS Cluster.",
				Computed:    true,
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the NFS Cluster.",
				Computed:    true,
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
	client := meta.(bundleclient.SdkBundle).NFSClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	partialMatch := d.Get("partial_match").(bool)
	id := idValue.(string)
	name := nameValue.(string)
	location := d.Get("location").(string)

	if idOk && nameOk {
		return utils.ToDiags(d, "ID and name cannot be both specified at the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the NFS Cluster ID or name", nil)
	}

	var cluster nfs.ClusterRead
	var apiResponse *shared.APIResponse
	var err error
	if idOk {
		cluster, apiResponse, err = client.GetNFSClusterByID(ctx, id, location)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the NFS Cluster with ID: %s, error: %s", idValue, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		clusters, apiResponse, err := client.ListNFSClusters(ctx, d)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching NFS Clusters: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}

		var results []nfs.ClusterRead
		for _, cl := range clusters.Items {
			if utils.NameMatches(cl.Properties.Name, name, partialMatch) {
				results = append(results, cl)
			}
		}

		switch {
		case len(results) == 0:
			return utils.ToDiags(d, fmt.Sprintf("no NFS Cluster found with the specified name: %s", name), nil)
		case len(results) > 1:
			return utils.ToDiags(d, fmt.Sprintf("more than one NFS Cluster found with the specified name: %s", name), nil)
		default:
			cluster = results[0]
		}
	}

	if err = client.SetNFSClusterData(d, cluster); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
