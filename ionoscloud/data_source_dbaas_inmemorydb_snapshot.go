package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

func dataSourceDBaaSInMemoryDBSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInMemoryDBSnapshotRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the snapshot.",
				Required:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The location of the snapshot.",
				Required:    true,
			},
			"metadata": {
				Type:        schema.TypeList,
				Description: "The metadata of the snapshot.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_date": {
							Type:        schema.TypeString,
							Description: "The ISO 8601 creation timestamp.",
							Computed:    true,
						},
						"last_modified_date": {
							Type:        schema.TypeString,
							Description: "The ISO 8601 modified timestamp.",
							Computed:    true,
						},
						"replica_set_id": {
							Type:        schema.TypeString,
							Description: "The ID of the InMemoryDB replica set the snapshot is taken from.",
							Computed:    true,
						},
						"snapshot_time": {
							Type:        schema.TypeString,
							Description: "The time the snapshot was dumped from the replica set.",
							Computed:    true,
						},
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The ID of the datacenter the snapshot was created in. Please mind, that the snapshot is not available in other datacenters.",
							Computed:    true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceInMemoryDBSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).InMemoryDBClient
	id := d.Get("id").(string)
	location := d.Get("location").(string)

	snapshot, _, err := client.GetSnapshot(ctx, id, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching the InMemoryDB snapshot with ID: %v, error: %w", id, err))
	}
	if err := client.SetSnapshotData(d, snapshot); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
