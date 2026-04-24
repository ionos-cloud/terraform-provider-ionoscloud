package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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
			"location": {
				Type:        schema.TypeString,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				Optional:    true,
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
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	datacenterID := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("ID and name cannot be both specified at the same time"), nil)
	}
	if !idOk && !nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("please provide either the GPU ID or name"), nil)
	}

	idStr, idIsString := id.(string)
	if idOk && (!idIsString || idStr == "") {
		return diagutil.ToDiags(d, fmt.Errorf("the provided ID is not valid"), nil)
	}

	nameStr, nameIsString := name.(string)
	if nameOk && (!nameIsString || nameStr == "") {
		return diagutil.ToDiags(d, fmt.Errorf("the provided name is not valid"), nil)
	}

	var gpu ionoscloud.Gpu
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		gpu, apiResponse, err = client.GraphicsProcessingUnitCardsApi.
			DatacentersServersGPUsFindById(ctx, datacenterID, serverID, idStr).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching GPU with ID %s: %w", idStr, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
	} else {
		/* search by name */
		var gpus ionoscloud.Gpus
		gpus, apiResponse, err = client.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsGet(ctx, datacenterID, serverID).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching GPUs: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}

		var results []ionoscloud.Gpu
		if gpus.Items != nil {
			for _, g := range *gpus.Items {
				if g.Properties != nil && g.Properties.Name != nil && *g.Properties.Name == nameStr {
					/* GPU found */
					if g.Id == nil {
						return diagutil.ToDiags(d, fmt.Errorf("GPU found with name %s but returned without an ID", nameStr), nil)
					}
					gpu, apiResponse, err = client.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsFindById(ctx, datacenterID, serverID, *g.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching GPU %s: %w", *g.Id, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
					}
					results = append(results, gpu)
				}
			}
		}

		switch {
		case len(results) == 0:
			return diagutil.ToDiags(d, fmt.Errorf("no GPU found with the specified criteria: name = %s", nameStr), nil)
		case len(results) > 1:
			return diagutil.ToDiags(d, fmt.Errorf("more than one GPU found with the specified criteria: name = %s", nameStr), nil)
		default:
			gpu = results[0]
		}
	}

	if gpu.Id == nil {
		return diagutil.ToDiags(d, fmt.Errorf("GPU returned without an ID"), nil)
	}

	d.SetId(*gpu.Id)

	if err := d.Set("datacenter_id", datacenterID); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	if err := d.Set("server_id", serverID); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	if gpu.Properties != nil {
		if gpu.Properties.Name != nil {
			if err := d.Set("name", *gpu.Properties.Name); err != nil {
				return diagutil.ToDiags(d, err, nil)
			}
		}
		if gpu.Properties.Vendor != nil {
			if err := d.Set("vendor", *gpu.Properties.Vendor); err != nil {
				return diagutil.ToDiags(d, err, nil)
			}
		}
		if gpu.Properties.Type != nil {
			if err := d.Set("type", *gpu.Properties.Type); err != nil {
				return diagutil.ToDiags(d, err, nil)
			}
		}
		if gpu.Properties.Model != nil {
			if err := d.Set("model", *gpu.Properties.Model); err != nil {
				return diagutil.ToDiags(d, err, nil)
			}
		}
	}

	return nil
}
