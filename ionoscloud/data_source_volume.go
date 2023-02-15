package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func dataSourceVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVolumeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"licence_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sshkey": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bus": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_plug": {
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
			"backup_unit_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"boot_server": {
				Type:        schema.TypeString,
				Description: "The UUID of the attached server.",
				Computed:    true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return diag.FromErr(errors.New("no datacenter_id was specified"))
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		diags := diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !nameOk {
		diags := diag.FromErr(errors.New("please provide either the volume id or name"))
		return diags
	}
	var volume ionoscloud.Volume
	var err error
	var apiResponse *shared.APIResponse

	if idOk {
		/* search by ID */

		volume, apiResponse, err = client.VolumesApi.DatacentersVolumesFindById(ctx, datacenterId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching volume with ID %s: %w", id.(string), err))
			return diags
		}
	} else {
		/* search by name */
		var volumes ionoscloud.Volumes
		volumes, apiResponse, err = client.VolumesApi.DatacentersVolumesGet(ctx, datacenterId.(string)).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching volumes: %w", err))
			return diags
		}

		var results []ionoscloud.Volume
		if volumes.Items != nil {
			for _, v := range *volumes.Items {
				if v.Properties != nil && v.Properties.Name != nil && *v.Properties.Name == name.(string) {
					/* volume found */
					volume, apiResponse, err = client.VolumesApi.DatacentersVolumesFindById(ctx, datacenterId.(string), *v.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						diags := diag.FromErr(fmt.Errorf("an error occurred while fetching volume %s: %w", *v.Id, err))
						return diags
					}
					results = append(results, volume)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no volume found with the specified criteria: name = %s", name.(string)))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one volume found with the specified criteria: name = %s", name.(string)))
		} else {
			volume = results[0]
		}
	}

	if err = setVolumeData(d, &volume); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
