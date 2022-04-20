package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
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
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
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

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	locationValue, locationOk := d.GetOk("location")
	sizeValue, sizeOk := d.GetOk("size")

	id := idValue.(string)
	name := nameValue.(string)
	location := locationValue.(string)
	size := float32(sizeValue.(int))

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the server id or name"))
	}

	var snapshot ionoscloud.Snapshot
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		snapshot, apiResponse, err = client.SnapshotsApi.SnapshotsFindById(ctx, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the snapshot with ID %s: %s", id, err))
			return diags
		}
	} else {
		var results []ionoscloud.Snapshot

		request := client.SnapshotsApi.SnapshotsGet(ctx).Depth(1)

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for snapshot by name with partial_match %t and name: %s", partialMatch, name)

		if nameOk && partialMatch {
			request = request.Filter("name", name)
		}

		snapshots, apiResponse, err := request.Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while fetching IonosCloud locations %s", err))
			return diags
		}

		if partialMatch {
			results = *snapshots.Items
		} else {
			if nameOk && snapshots.Items != nil {
				for _, snp := range *snapshots.Items {
					if snp.Properties != nil && snp.Properties.Name != nil && strings.EqualFold(*snp.Properties.Name, name) {
						results = append(results, snp)
					}
				}
			}
		}

		if locationOk {
			var locationResults []ionoscloud.Snapshot
			for _, snp := range results {
				if *snp.Properties.Location == location {
					locationResults = append(locationResults, snp)
				}

			}
			results = locationResults
		}

		if sizeOk {
			var sizeResults []ionoscloud.Snapshot
			for _, snp := range results {
				if snp.Properties != nil && snp.Properties.Size != nil && *snp.Properties.Size == size {
					sizeResults = append(sizeResults, snp)
				}

			}
			results = sizeResults
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no snapshot found with the specified criteria: name = %s, location = %s, size = %v", name, location, size))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one snapshot found with the specified criteria: name = %s, location = %s, size = %v", name, location, size))
		} else {
			snapshot = results[0]
		}
	}

	if err = setSnapshotData(d, &snapshot); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
