package ionoscloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

func dataSourceGpu() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGpuRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"server_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"model": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceGpuRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	datacenterID := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the GPU id or name"))
	}

	var gpu ionoscloud.Gpu
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		gpu, apiResponse, err = client.GraphicsProcessingUnitCardsApi.
			DatacentersServersGPUsFindById(ctx, datacenterID, serverID, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching GPU with ID %s: %w", id.(string), err))
		}
	} else {
		/* search by name */
		var gpus ionoscloud.Gpus
		gpus, apiResponse, err = client.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsGet(ctx, datacenterID, serverID).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching GPUs: %w", err))
		}

		var results []ionoscloud.Gpu
		if gpus.Items != nil {
			for _, g := range *gpus.Items {
				if g.Properties != nil && g.Properties.Name != nil && *g.Properties.Name == name.(string) {
					/* GPU found */
					gpu, apiResponse, err = client.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsFindById(ctx, datacenterID, serverID, *g.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diag.FromErr(fmt.Errorf("an error occurred while fetching GPU %s: %w", *g.Id, err))
					}
					results = append(results, gpu)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no GPU found with the specified criteria: name = %s", name.(string)))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one GPU found with the specified criteria: name = %s", name.(string)))
		} else {
			gpu = results[0]
		}
	}

	if gpu.Id == nil {
		return diag.FromErr(fmt.Errorf("GPU returned without an ID"))
	}

	d.SetId(*gpu.Id)

	if err := d.Set("datacenter_id", datacenterID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("server_id", serverID); err != nil {
		return diag.FromErr(err)
	}

	if gpu.Properties != nil {
		if gpu.Properties.Name != nil {
			if err := d.Set("name", *gpu.Properties.Name); err != nil {
				return diag.FromErr(err)
			}
		}
		if gpu.Properties.Vendor != nil {
			if err := d.Set("vendor", *gpu.Properties.Vendor); err != nil {
				return diag.FromErr(err)
			}
		}
		if gpu.Properties.Type != nil {
			if err := d.Set("type", *gpu.Properties.Type); err != nil {
				return diag.FromErr(err)
			}
		}
		if gpu.Properties.Model != nil {
			if err := d.Set("model", *gpu.Properties.Model); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}
