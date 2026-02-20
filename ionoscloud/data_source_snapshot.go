package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSnapshotRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A name of that resource",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Location of that image/snapshot",
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
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
			"require_legacy_bios": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the image requires the legacy BIOS for compatibility or specific needs.",
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	location, locationOk := d.GetOk("location")
	size, sizeOk := d.GetOk("size")

	if idOk && nameOk {
		return utils.ToDiags(d, "id and name cannot be both specified in the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the server id or name", nil)
	}

	var snapshot ionoscloud.Snapshot
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		snapshot, apiResponse, err = client.SnapshotsApi.SnapshotsFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the snapshot with ID %s: %s", id.(string), err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		var results []ionoscloud.Snapshot

		snapshots, apiResponse, err := client.SnapshotsApi.SnapshotsGet(ctx).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching IonosCloud locations %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}

		if snapshots.Items != nil {
			for _, snp := range *snapshots.Items {
				if snp.Properties != nil && snp.Properties.Name != nil && *snp.Properties.Name == name.(string) {
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
				if snp.Properties != nil && snp.Properties.Size != nil && *snp.Properties.Size == float32(size.(int)) {
					sizeResults = append(sizeResults, snp)
				}

			}
			results = sizeResults
		}

		if results == nil || len(results) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no snapshot found with the specified criteria: name = %s, location = %s, size = %v", name.(string), location.(string), size.(int)), nil)
		} else if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one snapshot found with the specified criteria: name = %s, location = %s, size = %v", name.(string), location.(string), size.(int)), nil)
		} else {
			snapshot = results[0]
		}
	}

	if err = setSnapshotData(d, &snapshot); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
