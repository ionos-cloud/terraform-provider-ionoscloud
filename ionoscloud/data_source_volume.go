package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func dataSourceVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVolumeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image": {
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
			"expose_serial": {
				Type: schema.TypeBool,
				Description: "If set to `true` will expose the serial id of the disk attached to the server. " +
					"If set to `false` will not expose the serial id. Some operating systems or software solutions require the serial id to be exposed to work properly. " +
					"Exposing the serial can influence licensed software (e.g. Windows) behavior",
				Computed: true,
			},
			"require_legacy_bios": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the image requires the legacy BIOS for compatibility or specific needs.",
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				Optional:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceVolumeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return diagutil.ToDiags(d, fmt.Errorf("no datacenter_id was specified"), nil)
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("ID and name cannot be both specified in the same time"), nil)
	}
	if !idOk && !nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("please provide either the volume ID or name"), nil)
	}
	var volume ionoscloud.Volume
	var err error
	var apiResponse *ionoscloud.APIResponse

	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	if idOk {
		/* search by ID */

		volume, apiResponse, err = client.VolumesApi.DatacentersVolumesFindById(ctx, datacenterId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching volume with ID %s: %w", id.(string), err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
	} else {
		/* search by name */
		var volumes ionoscloud.Volumes
		volumes, apiResponse, err = client.VolumesApi.DatacentersVolumesGet(ctx, datacenterId.(string)).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching volumes: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}

		var results []ionoscloud.Volume
		if volumes.Items != nil {
			for _, v := range *volumes.Items {
				if v.Properties != nil && v.Properties.Name != nil && *v.Properties.Name == name.(string) {
					/* volume found */
					volume, apiResponse, err = client.VolumesApi.DatacentersVolumesFindById(ctx, datacenterId.(string), *v.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching volume %s: %w", *v.Id, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
					}
					results = append(results, volume)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diagutil.ToDiags(d, fmt.Errorf("no volume found with the specified criteria: name = %s", name.(string)), nil)
		} else if len(results) > 1 {
			return diagutil.ToDiags(d, fmt.Errorf("more than one volume found with the specified criteria: name = %s", name.(string)), nil)
		} else {
			volume = results[0]
		}
	}

	if err = setVolumeData(d, &volume); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}
