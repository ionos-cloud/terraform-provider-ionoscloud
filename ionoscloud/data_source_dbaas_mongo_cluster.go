package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
)

func dataSourceDbaasMongoCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasMongoReadCluster,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Description:      "The id of your cluster.",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "a weekly 4 hour-long window, during which maintenance might occur",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"day_of_the_week": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"mongodb_version": {
				Type:        schema.TypeString,
				Description: "The MongoDB version of your cluster.",
				Computed:    true,
			},
			"instances": {
				Type:        schema.TypeInt,
				Description: "The total number of instances in the cluster (one master and n-1 standbys)",
				Computed:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The friendly name of your cluster.",
				Optional:    true,
			},
			"location": {
				Type: schema.TypeString,
				Description: "The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation (disallowed in update requests)." +
					"Available locations: de/txl, gb/lhr, es/vit",
				Computed: true,
			},
			"connections": {
				Type:        schema.TypeList,
				Description: "Details about the network connection for your cluster.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The datacenter to connect your cluster to.",
							Computed:    true,
						},
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The LAN to connect your cluster to.",
							Computed:    true,
						},
						"cidr_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The list of IPs and subnet for your cluster.\n          Note the following unavailable IP ranges:\n          10.233.64.0/18\n          10.233.0.0/18\n          10.233.114.0/24 		\n example: [192.168.1.100/24, 192.168.1.101/24]",
						},
					},
				},
			},
			"template_id": {
				Type:        schema.TypeString,
				Description: "The unique ID of the template, which specifies the number of cores, storage size, and memory.",
				Computed:    true,
			},
			"connection_string": {
				Type:        schema.TypeString,
				Description: "The connection string for your cluster.",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDbaasMongoReadCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MongoClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("display_name")

	if idOk && nameOk {
		diags := diag.FromErr(errors.New("id and display_name cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !nameOk {
		diags := diag.FromErr(errors.New("please provide either the dbaas cluster id or display_name"))
		return diags
	}

	var cluster mongo.ClusterResponse
	var err error

	if idOk {
		/* search by ID */
		cluster, _, err = client.GetCluster(ctx, id.(string))
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the dbaas mongo cluster with ID %s: %w", id.(string), err))
			return diags
		}
	} else {
		clusters, _, err := client.ListClusters(ctx, "")

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching dbaas mongo clusters: %w", err))
			return diags
		}

		var results []mongo.ClusterResponse

		if clusters.Items != nil && len(*clusters.Items) > 0 {
			for _, clusterItem := range *clusters.Items {
				if clusterItem.Properties != nil && clusterItem.Properties.DisplayName != nil && strings.EqualFold(*clusterItem.Properties.DisplayName, name.(string)) {
					results = append(results, clusterItem)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no DBaaS mongo cluster found with the specified name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one DBaaS mongo cluster found with the specified criteria name = %s", name))
		} else {
			cluster = results[0]
		}

	}

	if err := dbaasService.SetDbaasMongoDBClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil

}
