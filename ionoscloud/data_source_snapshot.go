package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSnapshotRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A name of that resource",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location of that image/snapshot",
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The size of the image in GB",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human readable description",
			},
			"licence_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS type of this Snapshot",
			},
			"sec_auth_protection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Boolean value representing if the snapshot requires extra protection e.g. two factor protection",
			},
			"cpu_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cpu_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nic_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nic_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_virtio_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_virtio_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_scsi_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_scsi_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	id, idOk := d.GetOk("id")
	name := d.Get("name").(string)
	location, locationOk := d.GetOk("location")
	size, sizeOk := d.GetOk("size")

	var snapshot ionoscloud.Snapshot
	var err error
	if idOk {
		/* search by ID */
		snapshot, _, err = client.SnapshotsApi.SnapshotsFindById(ctx, id.(string)).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the snapshot with ID %s: %s", id.(string), err))
			return diags
		}
	} else {
		var results []ionoscloud.Snapshot

		snapshots, _, err := client.SnapshotsApi.SnapshotsGet(ctx).Execute()

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while fetching IonosCloud locations %s", err))
			return diags
		}

		if snapshots.Items != nil {
			for _, snp := range *snapshots.Items {
				if snp.Properties.Name != nil && *snp.Properties.Name == name {
					results = append(results, snp)
				}
			}
		}

		if locationOk {
			var locationResults []ionoscloud.Snapshot
			for _, snp := range results {
				if *snp.Properties.Location == location.(string) {
					locationResults = append(locationResults, snp)
				}

			}
			results = locationResults
		}

		if sizeOk {
			var sizeResults []ionoscloud.Snapshot
			for _, snp := range results {
				if *snp.Properties.Size <= float32(size.(int)) {
					sizeResults = append(sizeResults, snp)
				}

			}
			results = sizeResults
		}

		if len(results) == 0 {
			diags := diag.FromErr(fmt.Errorf("There are no snapshots that match the search criteria "))
			return diags
		}
		snapshot = results[0]
	}

	if err = setSnapshotData(d, &snapshot); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
