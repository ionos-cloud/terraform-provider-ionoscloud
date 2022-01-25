package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDataCenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataCenterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
				Type:     schema.TypeList,
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDataCenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

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
		return diag.FromErr(fmt.Errorf("either id, location or name must be set"))
	}

	if idOk {
		datacenter, apiResponse, err = client.DataCentersApi.DatacentersFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting datacenter with id %s %s", id.(string), err))
		}
		if nameOk {
			if *datacenter.Properties.Name != name {
				return diag.FromErr(fmt.Errorf("name of dc (UUID=%s, name=%s) does not match expected name: %s",
					*datacenter.Id, *datacenter.Properties.Name, name))
			}
		}
		if locationOk {
			if *datacenter.Properties.Location != location {
				return diag.FromErr(fmt.Errorf("location of dc (UUID=%s, location=%s) does not match expected location: %s",
					*datacenter.Id, *datacenter.Properties.Location, location))
			}
		}
		if datacenter.Properties != nil {
			log.Printf("[INFO] Got backupUnit [Name=%s], Location=%s, [Id=%s]", *datacenter.Properties.Name, *datacenter.Properties.Location, *datacenter.Id)
		}
	} else {

		var results []ionoscloud.Datacenter
		var resultsByDatacenter ionoscloud.Datacenters
		var resultsByLocation ionoscloud.Datacenters

		if nameOk {
			resultsByDatacenter, apiResponse, err = client.DataCentersApi.DatacentersGet(ctx).Depth(1).Filter("name", name).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit: %s", err.Error()))
			}

			if resultsByDatacenter.Items == nil || len(*resultsByDatacenter.Items) == 0 {
				return diag.FromErr(fmt.Errorf("no datacenter found with the specified name"))
			}
		}

		if locationOk {
			if resultsByDatacenter.Items != nil && len(*resultsByDatacenter.Items) > 0 {
				for _, dc := range *resultsByDatacenter.Items {
					if dc.Properties.Location != nil && *dc.Properties.Location == location {
						results = append(results, dc)
						break
					}
				}
				if results == nil || len(results) == 0 {
					return diag.FromErr(fmt.Errorf("no datacenter found with the specified location"))
				}
			} else {
				/* find the first datacenter matching the location */
				resultsByLocation, apiResponse, err = client.DataCentersApi.DatacentersGet(ctx).Depth(1).Filter("location", location).Execute()
				logApiRequestTime(apiResponse)

				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit: %s", err.Error()))
				}
				if resultsByLocation.Items == nil || len(*resultsByLocation.Items) == 0 {
					return diag.FromErr(fmt.Errorf("no datacenter found with the specified location"))
				}
			}
		}

		if results != nil && len(results) > 0 {
			datacenter = results[len(results)-1]
		} else if resultsByLocation.Items != nil && len(*resultsByLocation.Items) > 0 {
			datacenter = (*resultsByLocation.Items)[len(*resultsByLocation.Items)-1]
		} else if resultsByDatacenter.Items != nil && len(*resultsByDatacenter.Items) > 0 {
			datacenter = (*resultsByDatacenter.Items)[len(*resultsByDatacenter.Items)-1]
		}

	}

	if err := setDatacenterData(d, &datacenter); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
