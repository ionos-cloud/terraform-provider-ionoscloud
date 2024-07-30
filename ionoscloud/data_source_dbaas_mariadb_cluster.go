package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dbaas "github.com/ionos-cloud/sdk-go-dbaas-mariadb"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func dataSourceDBaaSMariaDBCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMariaDBClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Description:      "The id of your cluster.",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"location": {
				Type:             schema.TypeString,
				Description:      "The cluster location",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(constant.Locations, false)),
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The friendly name of your cluster.",
				Optional:    true,
			},
			"mariadb_version": {
				Type:        schema.TypeString,
				Description: "The MariaDB version of your cluster.",
				Computed:    true,
			},
			"instances": {
				Type:        schema.TypeInt,
				Description: "The total number of instances in the cluster (one primary and n-1 secondary).",
				Computed:    true,
			},
			"cores": {
				Type:        schema.TypeInt,
				Description: "The number of CPU cores per instance.",
				Computed:    true,
			},
			"ram": {
				Type:        schema.TypeInt,
				Description: "The amount of memory per instance in gigabytes (GB).",
				Computed:    true,
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Description: "The amount of storage per instance in gigabytes (GB).",
				Computed:    true,
			},
			"connections": {
				Type:        schema.TypeList,
				Description: "The network connection for your cluster. Only one connection is allowed.",
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
							Description: "The numeric LAN ID to connect your cluster to.",
							Computed:    true,
						},
						"cidr": {
							Type:        schema.TypeString,
							Description: "The IP and subnet for your cluster.",
							Computed:    true,
						},
					},
				},
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "A weekly 4 hour-long window, during which maintenance might occur.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Description: "Start of the maintenance window in UTC time.",
							Computed:    true,
						},
						"day_of_the_week": {
							Type:        schema.TypeString,
							Description: "The name of the week day.",
							Computed:    true,
						},
					},
				},
			},
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The DNS name pointing to your cluster.",
				Computed:    true,
			},
		},
	}
}

func dataSourceMariaDBClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MariaDBClient
	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("display_name")
	location := d.Get("location").(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and display_name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the MariaDB cluster ID or display_name"))
	}

	var cluster dbaas.ClusterResponse
	var err error

	if idOk {
		/* search by ID */
		cluster, _, err = client.GetCluster(ctx, id.(string), location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the MariaDB cluster with ID %v: %w", id.(string), err))
		}
	} else {
		clusters, _, err := client.ListClusters(ctx, "", location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching MariaDB clusters: %w", err))
		}

		var results []dbaas.ClusterResponse

		if clusters.Items != nil {
			for _, clusterItem := range *clusters.Items {
				if clusterItem.Properties != nil && clusterItem.Properties.DisplayName != nil && strings.EqualFold(*clusterItem.Properties.DisplayName, name.(string)) {
					results = append(results, clusterItem)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no MariaDB cluster found with the specified display name: %v", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one MariaDB cluster found with the specified criteria name: %v", name))
		} else {
			cluster = results[0]
		}

	}

	if err := client.SetMariaDBClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
