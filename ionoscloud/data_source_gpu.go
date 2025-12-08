package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

func dataSourceGpu() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGpuRead,
		Schema: map[string]*schema.Schema{
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
			"gpu_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"name": {
				Type:     schema.TypeString,
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

/* example response:
{
  "id": "f9cba8aa-5847-4681-a488-342880f710ef",
  "type": "gpu",
  "href": "https://api.ionos.com/cloudapi/v6/datacenters/e0045d38-db36-48ed-9775-905968bead77/servers/cf6c8a6f-c652-4db2-8422-472726e6da8c/gpus/f9cba8aa-5847-4681-a488-342880f710ef",
  "metadata": {
    "etag": "5c9dd5532478dde25f2d7349c56f62ff",
    "createdDate": "2025-12-08T11:12:35Z",
    "createdBy": "terraform-v6@cloud.ionos.com",
    "createdByUserId": "4df59ddd-94d3-4a86-99f0-411536964cbf",
    "lastModifiedDate": "2025-12-08T11:12:35Z",
    "lastModifiedBy": "terraform-v6@cloud.ionos.com",
    "lastModifiedByUserId": "4df59ddd-94d3-4a86-99f0-411536964cbf",
    "state": "AVAILABLE"
  },
  "properties": {
    "name": "GPU NVIDIA Corporation GH100 [H200 NVL] 1",
    "vendor": "NVIDIA Corporation",
    "type": "passthrough",
    "model": "GH100 [H200 NVL]"
  }
}
*/

func dataSourceGpuRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	datacenterID := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)
	gpuID := d.Get("gpu_id").(string)

	gpu, apiResponse, err := client.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsFindById(
		ctx, datacenterID, serverID, gpuID,
	).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch GPU %s: %w", gpuID, err))
	}

	if gpu.Id == nil {
		return diag.FromErr(fmt.Errorf("GPU %s returned without an ID", gpuID))
	}

	d.SetId(*gpu.Id)

	if err := d.Set("datacenter_id", datacenterID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("server_id", serverID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("gpu_id", *gpu.Id); err != nil {
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
