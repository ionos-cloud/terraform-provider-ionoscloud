package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceDataCenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataCenterRead,
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
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"features": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sec_auth_protection": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cpu_architecture": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vendor": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDataCenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	var name, location string
	id, idOk := d.GetOk("id")

	t, nameOk := d.GetOk("name")
	if nameOk {
		name = t.(string)
	}

	t, locationOk := d.GetOk("location")
	if locationOk {
		location = t.(string)
	}

	var datacenter ionoscloud.Datacenter
	var err error
	var apiResponse *ionoscloud.APIResponse

	if !idOk && !nameOk && !locationOk {
		return utils.ToDiags(d, "either id, location or name must be set", nil)
	}

	if idOk {
		datacenter, apiResponse, err = client.DataCentersApi.DatacentersFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("error getting datacenter with id %s %s", id.(string), err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		if nameOk {
			if *datacenter.Properties.Name != name {
				return utils.ToDiags(d, fmt.Sprintf("name of dc (UUID=%s, name=%s) does not match expected name: %s",
					*datacenter.Id, *datacenter.Properties.Name, name), nil)
			}
		}
		if locationOk {
			if *datacenter.Properties.Location != location {
				return utils.ToDiags(d, fmt.Sprintf("location of dc (UUID=%s, location=%s) does not match expected location: %s",
					*datacenter.Id, *datacenter.Properties.Location, location), nil)
			}
		}
		if datacenter.Properties != nil {
			log.Printf("[INFO] Got dc [Name=%s, Location=%s]", *datacenter.Properties.Name, *datacenter.Properties.Location)
		}

	} else {
		datacenters, apiResponse, err := client.DataCentersApi.DatacentersGet(ctx).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching datacenters: %s ", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}

		var results []ionoscloud.Datacenter

		if nameOk && datacenters.Items != nil {
			var resultsByDatacenter []ionoscloud.Datacenter
			for _, dc := range *datacenters.Items {
				if dc.Properties != nil && dc.Properties.Name != nil && *dc.Properties.Name == name {
					resultsByDatacenter = append(resultsByDatacenter, dc)
				}
			}

			if resultsByDatacenter == nil {
				return utils.ToDiags(d, fmt.Sprintf("no datacenter found with the specified criteria: name = %s", name), nil)
			} else {
				results = resultsByDatacenter
			}
		}

		if locationOk {
			var resultsByLocation []ionoscloud.Datacenter
			if results != nil {
				for _, dc := range results {
					if dc.Properties.Location != nil && *dc.Properties.Location == location {
						resultsByLocation = append(resultsByLocation, dc)
					}
				}
			} else if datacenters.Items != nil {
				/* find the first datacenter matching the location */
				for _, dc := range *datacenters.Items {
					if dc.Properties.Location != nil && *dc.Properties.Location == location {
						resultsByLocation = append(resultsByLocation, dc)
					}
				}
			}
			if resultsByLocation == nil {
				return utils.ToDiags(d, fmt.Sprintf("no datacenter found with the specified criteria: location = %s", location), nil)
			} else {
				results = resultsByLocation
			}
		}

		if results == nil || len(results) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no datacenter found with the specified criteria: name = %s location = %s", name, location), nil)
		} else if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one datacenter found with the specified criteria: name = %s location = %s", name, location), nil)
		} else {
			datacenter = results[0]
		}

	}

	if err := setDatacenterData(d, &datacenter); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
