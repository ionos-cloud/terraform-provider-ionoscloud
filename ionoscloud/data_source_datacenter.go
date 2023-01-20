package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"

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
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDataCenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	var id, name, location string

	t, idOk := d.GetOk("id")
	if idOk {
		id = t.(string)
	}

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
		log.Printf("[INFO] Using data source for datacenter by id %s", id)
		datacenter, apiResponse, err = client.DataCentersApi.DatacentersFindById(ctx, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting datacenter with id %s %w", id.(string), err))
		}
		if nameOk {
			if datacenter.Properties != nil && datacenter.Properties.Name != nil && (!strings.EqualFold(*datacenter.Properties.Name, name) || !strings.Contains(*datacenter.Properties.Name, name)) {
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
			log.Printf("[INFO] Got dc [Name=%s, Location=%s]", *datacenter.Properties.Name, *datacenter.Properties.Location)
		}

	} else {
		var results []ionoscloud.Datacenter

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for datacenter by name with partial_match %t and name: %s", partialMatch, name)

		if nameOk && partialMatch {
			datacenters, apiResponse, err := client.DataCentersApi.DatacentersGet(ctx).Depth(1).Filter("name", name).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occured while fetching datacenters while searching by name: %s, %w ", name, err))
			}
			results = *datacenters.Items
		} else {
			datacenters, apiResponse, err := client.DataCentersApi.DatacentersGet(ctx).Depth(1).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occured while fetching datacenters while searching by name: %s, %w ", name, err))
			}

			results = *datacenters.Items

			if nameOk && results != nil {
				var resultsByName []ionoscloud.Datacenter
				for _, dc := range results {
					if dc.Properties != nil && dc.Properties.Name != nil && strings.EqualFold(*dc.Properties.Name, name) {
						resultsByName = append(resultsByName, dc)
					}
				}

				if resultsByName == nil {
					return diag.FromErr(fmt.Errorf("no datacenter found with the specified criteria: name = %s", name))
				} else {
					results = resultsByName
				}
			}
		}

		if locationOk {
			var resultsByLocation []ionoscloud.Datacenter
			if results != nil {
				for _, dc := range results {
					if dc.Properties.Location != nil && strings.EqualFold(*dc.Properties.Location, location) {
						resultsByLocation = append(resultsByLocation, dc)
					}
				}
			}
			if resultsByLocation == nil {
				return diag.FromErr(fmt.Errorf("no datacenter found with the specified criteria: location = %s", location))
			} else {
				results = resultsByLocation
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no datacenter found with the specified criteria: name = %s location = %s", name, location))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one datacenter found with the specified criteria: name = %s location = %s", name, location))
		} else {
			datacenter = results[0]
		}

	}

	if err := setDatacenterData(d, &datacenter); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
