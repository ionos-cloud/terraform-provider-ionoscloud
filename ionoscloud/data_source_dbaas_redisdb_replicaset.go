package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	redisdb "github.com/ionos-cloud/sdk-go-dbaas-redis"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func dataSourceDBaaSRedisDBReplicaSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDBaaSRedisDBReplicaSetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the Redis Replicaset.",
				Optional:    true,
			},
			"location": {
				Type:             schema.TypeString,
				Description:      "The cluster location",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(constant.MariaDBClusterLocations, false)),
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of the Redis Replicaset.",
				Optional:    true,
			},
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The DNS name of the Redis Replicaset.",
				Computed:    true,
			},
			"connections": {
				Type:        schema.TypeList,
				Description: "The network connection for your Replicaset. Only one connection is allowed.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Description: "The IP and subnet for your Replicaset.",
							Computed:    true,
						},
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The datacenter to connect your Replicaset to.",
							Computed:    true,
						},
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The numeric LAN ID to connect your Replicaset to.",
							Computed:    true,
						},
					},
				},
			},
			"credentials": {
				Type:        schema.TypeList,
				Description: "The credentials for your Replicaset.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:        schema.TypeString,
							Description: "The password for your Replicaset.",
							Computed:    true,
						},
						"username": {
							Type:        schema.TypeString,
							Description: "The username for your Replicaset.",
							Computed:    true,
						},
					},
				},
			},
			"eviction_policy": {
				Type:        schema.TypeString,
				Description: "The eviction policy of the Redis Replicaset.",
				Computed:    true,
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
			"persistence_mode": {
				Type:        schema.TypeString,
				Description: "The persistence mode of the Redis Replicaset.",
				Computed:    true,
			},
			"redis_version": {
				Type:        schema.TypeString,
				Description: "The version of Redis used in the Replicaset.",
				Computed:    true,
			},
			"replicas": {
				Type:        schema.TypeInt,
				Description: "The number of replicas in the Replicaset.",
				Computed:    true,
			},
			"resources": {
				Type:        schema.TypeList,
				Description: "The resources allocated to the Replicaset.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"storage": {
							Type:        schema.TypeInt,
							Description: "The amount of storage per instance in gigabytes (GB).",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDBaaSRedisDBReplicaSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).RedisDBClient
	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("display_name")
	location := d.Get("location").(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and display_name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the Redis replicaset ID or display_name"))
	}

	var cluster redisdb.ReplicaSetRead
	var err error

	if idOk {
		// search by ID
		cluster, _, err = client.GetReplicaSet(ctx, id.(string), location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the Redis cluster with ID %v: %w", id.(string), err))
		}
	} else {
		// list, then filter by name
		clusters, _, err := client.ListReplicaSets(ctx, name.(string), location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching Redis clusters: %w", err))
		}

		var results []redisdb.ReplicaSetRead

		if clusters.Items != nil {
			for _, clusterItem := range *clusters.Items {
				if clusterItem.Properties != nil && clusterItem.Properties.DisplayName != nil && strings.EqualFold(*clusterItem.Properties.DisplayName, name.(string)) {
					results = append(results, clusterItem)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no Redis cluster found with the specified display name: %v", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one Redis cluster found with the specified criteria name: %v", name))
		} else {
			cluster = results[0]
		}

	}

	if err := client.SetRedisDBReplicaSetData(d, cluster); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
