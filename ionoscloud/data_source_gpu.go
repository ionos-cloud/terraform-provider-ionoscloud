package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceGpu() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGpuRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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

//nolint:gocyclo
func dataSourceGpuRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	datacenterID := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return utils.ToDiags(d, "ID and name cannot be both specified at the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the GPU ID or name", nil)
	}

	idStr, idIsString := id.(string)
	if idOk && (!idIsString || idStr == "") {
		return utils.ToDiags(d, "the provided ID is not valid", nil)
	}

	nameStr, nameIsString := name.(string)
	if nameOk && (!nameIsString || nameStr == "") {
		return utils.ToDiags(d, "the provided name is not valid", nil)
	}

	var gpu ionoscloud.Gpu
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		gpu, apiResponse, err = client.GraphicsProcessingUnitCardsApi.
			DatacentersServersGPUsFindById(ctx, datacenterID, serverID, idStr).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching GPU with ID %s: %s", idStr, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		/* search by name */
		var gpus ionoscloud.Gpus
		gpus, apiResponse, err = client.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsGet(ctx, datacenterID, serverID).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching GPUs: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}

		var results []ionoscloud.Gpu
		if gpus.Items != nil {
			for _, g := range *gpus.Items {
				if g.Properties != nil && g.Properties.Name != nil && *g.Properties.Name == nameStr {
					/* GPU found */
					if g.Id == nil {
						return utils.ToDiags(d, fmt.Sprintf("GPU found with name %s but returned without an ID", nameStr), nil)
					}
					gpu, apiResponse, err = client.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsFindById(ctx, datacenterID, serverID, *g.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching GPU %s: %s", *g.Id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
					}
					results = append(results, gpu)
				}
			}
		}

		switch {
		case len(results) == 0:
			return utils.ToDiags(d, fmt.Sprintf("no GPU found with the specified criteria: name = %s", nameStr), nil)
		case len(results) > 1:
			return utils.ToDiags(d, fmt.Sprintf("more than one GPU found with the specified criteria: name = %s", nameStr), nil)
		default:
			gpu = results[0]
		}
	}

	if gpu.Id == nil {
		return utils.ToDiags(d, "GPU returned without an ID", nil)
	}

	d.SetId(*gpu.Id)

	if err := d.Set("datacenter_id", datacenterID); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	if err := d.Set("server_id", serverID); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	if gpu.Properties != nil {
		if gpu.Properties.Name != nil {
			if err := d.Set("name", *gpu.Properties.Name); err != nil {
				return utils.ToDiags(d, err.Error(), nil)
			}
		}
		if gpu.Properties.Vendor != nil {
			if err := d.Set("vendor", *gpu.Properties.Vendor); err != nil {
				return utils.ToDiags(d, err.Error(), nil)
			}
		}
		if gpu.Properties.Type != nil {
			if err := d.Set("type", *gpu.Properties.Type); err != nil {
				return utils.ToDiags(d, err.Error(), nil)
			}
		}
		if gpu.Properties.Model != nil {
			if err := d.Set("model", *gpu.Properties.Model); err != nil {
				return utils.ToDiags(d, err.Error(), nil)
			}
		}
	}

	return nil
}
