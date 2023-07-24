package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"log"
)

func dataSourceIpBlock() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceIpBlockRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"ips": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ip_consumers": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
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

func datasourceIpBlockRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id, idOk := data.GetOk("id")

	var name, location string

	t, nameOk := data.GetOk("name")
	if nameOk {
		name = t.(string)
	}

	t, locationOk := data.GetOk("location")
	if locationOk {
		location = t.(string)
	}
	var ipBlock ionoscloud.IpBlock
	var err error
	client := meta.(services.SdkBundle).CloudApiClient
	var apiResponse *ionoscloud.APIResponse

	if !idOk && !nameOk && !locationOk {
		return diag.FromErr(fmt.Errorf("either id, location or name must be set"))
	}
	if idOk {
		ipBlock, apiResponse, err = client.IPBlocksApi.IpblocksFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting ip block with id %s %w", id.(string), err))
		}
		if nameOk {
			if ipBlock.Properties != nil && *ipBlock.Properties.Name != name {
				return diag.FromErr(fmt.Errorf("name of ip block (UUID=%s, name=%s) does not match expected name: %s",
					*ipBlock.Id, *ipBlock.Properties.Name, name))
			}
		}
		if locationOk {
			if ipBlock.Properties != nil && *ipBlock.Properties.Location != location {
				return diag.FromErr(fmt.Errorf("location of ip block (UUID=%s, location=%s) does not match expected location: %s",
					*ipBlock.Id, *ipBlock.Properties.Location, location))
			}
		}
		log.Printf("[INFO] Got ip block [Name=%s, Location=%s]", *ipBlock.Properties.Name, *ipBlock.Properties.Location)
	} else {

		ipBlocks, apiResponse, err := client.IPBlocksApi.IpblocksGet(ctx).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching ipBlocks: %w ", err))
		}

		var results []ionoscloud.IpBlock

		if nameOk && ipBlocks.Items != nil {
			for _, block := range *ipBlocks.Items {
				if block.Properties != nil && block.Properties.Name != nil && *block.Properties.Name == name {
					results = append(results, block)
				}
			}

			if results == nil {
				return diag.FromErr(fmt.Errorf("no ip block found with the specified criteria name %s", name))
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
			return diag.FromErr(fmt.Errorf("no ip block found with the specified criteria name = %s, location = %s", name, location))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one ip block found with the specified criteria name = %s, location = %s", name, location))
		} else {
			ipBlock = results[0]
		}

	}

	if err := IpBlockSetData(data, &ipBlock); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
