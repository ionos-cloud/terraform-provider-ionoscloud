package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceIpBlock() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceIpBlockRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"location": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ips": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ip_consumers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"k8s_nodepool_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"k8s_cluster_uuid": {
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

//nolint:gocyclo
func datasourceIpBlockRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id, idOk := d.GetOk("id")

	var name, location string

	t, nameOk := d.GetOk("name")
	if nameOk {
		name = t.(string)
	}

	t, locationOk := d.GetOk("location")
	if locationOk {
		location = t.(string)
	}
	var ipBlock ionoscloud.IpBlock
	var err error
	var apiResponse *ionoscloud.APIResponse

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	if !idOk && !nameOk && !locationOk {
		return diagutil.ToDiags(d, fmt.Errorf("either id, location or name must be set"), nil)
	}
	if idOk {
		ipBlock, apiResponse, err = client.IPBlocksApi.IpblocksFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error getting ip block with id %s %w", id.(string), err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		if nameOk {
			if ipBlock.Properties != nil && *ipBlock.Properties.Name != name {
				return diagutil.ToDiags(d, fmt.Errorf("name of ip block (UUID=%s, name=%s) does not match expected name: %s",
					*ipBlock.Id, *ipBlock.Properties.Name, name), nil)
			}
		}
		if locationOk {
			if ipBlock.Properties != nil && *ipBlock.Properties.Location != location {
				return diagutil.ToDiags(d, fmt.Errorf("location of ip block (UUID=%s, location=%s) does not match expected location: %s",
					*ipBlock.Id, *ipBlock.Properties.Location, location), nil)
			}
		}
		log.Printf("[INFO] Got ip block [Name=%s, Location=%s]", *ipBlock.Properties.Name, *ipBlock.Properties.Location)
	} else {

		ipBlocks, apiResponse, err := client.IPBlocksApi.IpblocksGet(ctx).Depth(1).Limit(constant.IPBlockLimit).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching ipBlocks: %w ", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}

		var results []ionoscloud.IpBlock

		if nameOk && ipBlocks.Items != nil {
			for _, block := range *ipBlocks.Items {
				if block.Properties != nil && block.Properties.Name != nil && *block.Properties.Name == name {
					results = append(results, block)
				}
			}

			if results == nil {
				return diagutil.ToDiags(d, fmt.Errorf("no ip block found with the specified criteria name %s", name), nil)
			}
		}

		if locationOk {
			if results != nil {
				var locationResults []ionoscloud.IpBlock
				for _, block := range results {
					if block.Properties != nil && block.Properties.Location != nil && *block.Properties.Location == location {
						locationResults = append(locationResults, block)
					}
				}
				results = locationResults
			} else if ipBlocks.Items != nil {
				/* find the first ipblock matching the location */
				for _, block := range *ipBlocks.Items {
					if block.Properties != nil && block.Properties.Location != nil && *block.Properties.Location == location {
						results = append(results, block)
					}
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diagutil.ToDiags(d, fmt.Errorf("no ip block found with the specified criteria name = %s, location = %s", name, location), nil)
		} else if len(results) > 1 {
			return diagutil.ToDiags(d, fmt.Errorf("more than one ip block found with the specified criteria name = %s, location = %s", name, location), nil)
		} else {
			ipBlock = results[0]
		}

	}

	if err := IpBlockSetData(d, &ipBlock); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}
