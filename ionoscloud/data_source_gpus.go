package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

func dataSourceGpus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGpusRead,
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
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gpus": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceGpusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	datacenterID := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)

	gpus, apiResponse, err := client.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsGet(ctx, datacenterID, serverID).Depth(1).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching GPUs for server %s in datacenter %s: %w", serverID, datacenterID, err))
	}

	d.SetId(fmt.Sprintf("%s/gpus", serverID))

	var gpuList []map[string]interface{}
	if gpus.Items != nil {
		for _, gpu := range *gpus.Items {
			gpuMap := make(map[string]interface{})
			if gpu.Id != nil {
				gpuMap["id"] = *gpu.Id
			}
			if gpu.Properties != nil {
				if gpu.Properties.Name != nil {
					gpuMap["name"] = *gpu.Properties.Name
				}
				if gpu.Properties.Vendor != nil {
					gpuMap["vendor"] = *gpu.Properties.Vendor
				}
				if gpu.Properties.Type != nil {
					gpuMap["type"] = *gpu.Properties.Type
				}
				if gpu.Properties.Model != nil {
					gpuMap["model"] = *gpu.Properties.Model
				}
			}
			gpuList = append(gpuList, gpuMap)
		}
	}

	if err := d.Set("gpus", gpuList); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
