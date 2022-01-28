package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strconv"
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
	client := meta.(SdkBundle).CloudApiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	location, locationOk := d.GetOk("location")
	size, sizeOk := d.GetOk("size")

	var snapshot ionoscloud.Snapshot
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		snapshot, apiResponse, err = client.SnapshotsApi.SnapshotsFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the snapshot with ID %s: %s", id.(string), err))
			return diags
		}
	} else {

		request := client.SnapshotsApi.SnapshotsGet(ctx).Depth(1)

		if nameOk {
			request = request.Filter("name", name.(string))
		}

		if locationOk {
			request = request.Filter("location", location.(string))
		}

		if sizeOk {
			request = request.Filter("size", strconv.Itoa(size.(int)))
		}

		snapshots, apiResponse, err := request.Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while fetching IonosCloud locations %s", err))
			return diags
		}

		if snapshots.Items != nil && len(*snapshots.Items) > 0 {
			snapshot = (*snapshots.Items)[len(*snapshots.Items)-1]
			log.Printf("[INFO] %v snapshots found matching the search criteria. Getting the latest snapshot from the list %v", len(*snapshots.Items), *snapshot.Id)
		} else {
			return diag.FromErr(fmt.Errorf("no snapshot found with the specified criteria: location %s, size %s", location.(string), strconv.Itoa(size.(int))))
		}
	}

	if err = setSnapshotData(d, &snapshot); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
