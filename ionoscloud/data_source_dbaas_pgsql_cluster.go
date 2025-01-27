package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dbaas "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
)

func dataSourceDbaasPgSqlCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasPgSqlReadCluster,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Description:      "The id of your cluster.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The name of your cluster.",
				Optional:    true,
				Computed:    true,
			},
			"postgres_version": {
				Type:        schema.TypeString,
				Description: "The PostgreSQL version of your cluster.",
				Computed:    true,
			},
			"instances": {
				Type:        schema.TypeInt,
				Description: "The total number of instances in the cluster (one master and n-1 standbys)",
				Computed:    true,
			},
			"cores": {
				Type:        schema.TypeInt,
				Description: "The number of CPU cores per replica.",
				Computed:    true,
			},
			"ram": {
				Type:        schema.TypeInt,
				Description: "The amount of memory per instance in megabytes. Has to be a multiple of 1024.",
				Computed:    true,
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Description: "The amount of storage per instance in megabytes.",
				Computed:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "The storage type used in your cluster.",
				Computed:    true,
			},
			"connection_pooler": {
				Type:        schema.TypeList,
				Description: "Configuration options for the connection pooler",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"pool_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Represents different modes of connection pooling for the connection pooler",
						},
					},
				},
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
						"cidr": {
							Type:        schema.TypeString,
							Description: "The IP and subnet for the database. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24",
							Computed:    true,
						},
					},
				},
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation (disallowed in update requests)",
				Computed:    true,
			},
			"backup_location": {
				Type:        schema.TypeString,
				Description: "The Object Storage location where the backups will be stored.",
				Computed:    true,
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
			"synchronization_mode": {
				Type:        schema.TypeString,
				Description: "Represents different modes of replication.",
				Computed:    true,
			},
			"from_backup": {
				Type:        schema.TypeList,
				Description: "The PostgreSQL version of your cluster.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_id": {
							Type:        schema.TypeString,
							Description: "The unique ID of the backup you want to restore.",
							Computed:    true,
						},
						"recovery_target_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely.",
						},
					},
				},
			},
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The DNS name pointing to your cluster",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDbaasPgSqlReadCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).PsqlClient

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

	var cluster dbaas.ClusterResponse
	var err error

	if idOk {
		/* search by ID */
		cluster, _, err = client.GetCluster(ctx, id.(string))
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the dbaas cluster with ID %s: %w", id.(string), err))
			return diags
		}
	} else {
		clusters, _, err := client.ListClusters(ctx, "")

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching dbaas clusters: %w", err))
			return diags
		}

		var results []dbaas.ClusterResponse

		if clusters.Items != nil && len(clusters.Items) > 0 {
			for _, clusterItem := range clusters.Items {
				if clusterItem.Properties != nil && clusterItem.Properties.DisplayName != nil && strings.EqualFold(*clusterItem.Properties.DisplayName, name.(string)) {
					results = append(results, clusterItem)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no DBaaS cluster found with the specified name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one DBaaS cluster found with the specified criteria name = %s", name))
		} else {
			cluster = results[0]
		}

	}

	if err := dbaasService.SetPgSqlClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil

}
