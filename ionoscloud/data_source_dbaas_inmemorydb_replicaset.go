package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	inMemoryDB "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

func dataSourceDBaaSInMemoryDBReplicaSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReplicaSetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the InMemoryDB Replicaset.",
				Optional:    true,
				Computed:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The replica set location",
				Optional:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of the InMemoryDB Replicaset.",
				Optional:    true,
				Computed:    true,
			},
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The DNS name of the InMemoryDB Replicaset.",
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
				Description: "The eviction policy of the InMemoryDB Replicaset.",
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
				Description: "The persistence mode of the InMemoryDB Replicaset.",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The version of InMemoryDB used in the Replicaset.",
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
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceReplicaSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).InMemoryDBClient
	id, idOk := d.GetOk("id")
	displayName, displayNameOk := d.GetOk("display_name")
	location := d.Get("location").(string)

	if idOk && displayNameOk {
		return diag.FromErr(fmt.Errorf("ID and display_name cannot be both specified at the same time"))
	}
	if !idOk && !displayNameOk {
		return diag.FromErr(fmt.Errorf("please provide either the InMemoryDB replicaset ID or display_name"))
	}

	var replica inMemoryDB.ReplicaSetRead
	var err error

	if idOk {
		// search by ID
		replica, _, err = client.GetReplicaSet(ctx, id.(string), location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the InMemoryDB replica set with ID %v: %w", id.(string), err))
		}
	} else {
		// list, then filter by name
		clusters, _, err := client.ListReplicaSets(ctx, displayName.(string), location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching InMemoryDB replica sets: %w", err))
		}

		var results []inMemoryDB.ReplicaSetRead

		if clusters.Items != nil {
			for _, clusterItem := range clusters.Items {
				if strings.EqualFold(clusterItem.Properties.DisplayName, displayName.(string)) {
					results = append(results, clusterItem)
				}
			}
		}

		if len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no InMemoryDB replica set found with the specified display name: %v", displayName))
		}
		if len(results) > 1 {
			var ids []string
			for _, r := range results {
				ids = append(ids, r.Id)
			}
			return diag.FromErr(fmt.Errorf("more than one InMemoryDB replica set found with the specified criteria name '%v': (%v)", displayName, strings.Join(ids, ", ")))
		}
		replica = results[0]
	}

	if err := client.SetReplicaSetData(d, replica); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
