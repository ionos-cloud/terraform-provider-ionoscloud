package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
	client := meta.(*ionoscloud.APIClient)
	var apiResponse *ionoscloud.APIResponse

	if !idOk && !nameOk && !locationOk {
		return diag.FromErr(fmt.Errorf("either id, location or name must be set"))
	}
	if idOk {
		ipBlock, apiResponse, err = client.IPBlocksApi.IpblocksFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting ip block with id %s %s", id.(string), err))
		}
		if nameOk {
			if *ipBlock.Properties.Name != name {
				return diag.FromErr(fmt.Errorf("name of ip block (UUID=%s, name=%s) does not match expected name: %s",
					*ipBlock.Id, *ipBlock.Properties.Name, name))
			}
		}
		if locationOk {
			if *ipBlock.Properties.Location != location {
				return diag.FromErr(fmt.Errorf("location of ip block (UUID=%s, location=%s) does not match expected location: %s",
					*ipBlock.Id, *ipBlock.Properties.Location, location))
			}
		}
		log.Printf("[INFO] Got ip block [Name=%s, Location=%s]", *ipBlock.Properties.Name, *ipBlock.Properties.Location)
	} else {

		ipBlocks, apiResponse, err := client.IPBlocksApi.IpblocksGet(ctx).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching ipBlocks: %s ", err))
		}

		var results []ionoscloud.IpBlock

		if nameOk && ipBlocks.Items != nil {
			for _, block := range *ipBlocks.Items {
				if block.Properties.Name != nil && *block.Properties.Name == name {
					results = append(results, block)
					//found based on name only, save this in case we don't find based on location
					if !locationOk {
						ipBlock = results[0]
					}
				}
			}

			if results == nil {
				return diag.FromErr(fmt.Errorf("could not find an ip block with name %s", name))
			}
		}

		if locationOk {
			if results != nil {
				for _, block := range results {
					if block.Properties.Location != nil && *block.Properties.Location == location {
						ipBlock = block
						break
					}
				}
			} else if ipBlocks.Items != nil {
				/* find the first ipblock matching the location */
				for _, block := range *ipBlocks.Items {
					if block.Properties.Location != nil && *block.Properties.Location == location {
						ipBlock = block
						break
					}
				}
			}
		}

	}

	if ipBlock.Id == nil {
		return diag.FromErr(fmt.Errorf("there are no ip blocks that match the search criteria id = %s, name = %s, location = %s", id, name, location))
	}

	if err := IpBlockSetData(data, &ipBlock); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
