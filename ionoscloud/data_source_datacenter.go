package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDataCenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataCenterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: false,
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDataCenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

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

	if !idOk && !nameOk && !locationOk {
		return diag.FromErr(fmt.Errorf("either id, location or name must be set"))
	}

	if idOk {
		datacenter, _, err = client.DataCenterApi.DatacentersFindById(ctx, id.(string)).Execute()
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
		log.Printf("[INFO] Got dc [Name=%s, Location=%s]", *datacenter.Properties.Name, *datacenter.Properties.Location)
	} else {
		datacenters, _, err := client.DataCenterApi.DatacentersGet(ctx).Execute()

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching datacenters: %s ", err))
		}

		var results []ionoscloud.Datacenter

		if nameOk && datacenters.Items != nil {
			var resultsByDatacenter []ionoscloud.Datacenter
			for _, dc := range *datacenters.Items {
				if dc.Properties.Name != nil && *dc.Properties.Name == name {
					resultsByDatacenter = append(resultsByDatacenter, dc)
				}
			}

			if resultsByDatacenter == nil {
				return diag.FromErr(fmt.Errorf("could not find a datacenter with name %s", name))
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
						break
					}
				}
			} else if datacenters.Items != nil {
				/* find the first datacenter matching the location */
				for _, dc := range *datacenters.Items {
					if dc.Properties.Location != nil && *dc.Properties.Location == location {
						resultsByLocation = append(resultsByLocation, dc)
						break
					}
				}
			}
			if resultsByLocation == nil {
				return diag.FromErr(fmt.Errorf("could not find a datacenter with location %s", location))
			} else {
				results = resultsByLocation
			}
		}

		if results != nil {
			datacenter = results[0]
		} else {
			return diag.FromErr(fmt.Errorf("there are no datacenters that match the search criteria"))
		}
	}

	if err := setDatacenterData(d, &datacenter); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
