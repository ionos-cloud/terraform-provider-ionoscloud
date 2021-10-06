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

	found := false

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

		found = true
	} else {
		datacenters, _, err := client.DataCenterApi.DatacentersGet(ctx).Execute()

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching datacenters: %s ", err))
		}

		var results []ionoscloud.Datacenter

		if nameOk && datacenters.Items != nil {
			for _, dc := range *datacenters.Items {
				if dc.Properties.Name != nil && *dc.Properties.Name == name {
					results = append(results, dc)
				}
			}

			if results == nil {
				return diag.FromErr(fmt.Errorf("could not find a datacenter with name %s", name))
			}
		}

		if locationOk {
			if results != nil {
				for _, dc := range results {
					if dc.Properties.Location != nil && *dc.Properties.Location == location {
						datacenter = dc
						found = true
						break
					}
				}
			} else if datacenters.Items != nil {
				/* find the first datacenter matching the location */
				for _, dc := range *datacenters.Items {
					if dc.Properties.Location != nil && *dc.Properties.Location == location {
						datacenter = dc
						found = true
						break
					}
				}
			}
		}

	}

	if !found {
		return diag.FromErr(fmt.Errorf("there are no datacenters that match the search criteria"))
	}

	if datacenter.Id != nil {
		if err := d.Set("id", *datacenter.Id); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := setDatacenterData(d, &datacenter); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setDatacenterData(d *schema.ResourceData, datacenter *ionoscloud.Datacenter) error {

	if datacenter.Id != nil {
		d.SetId(*datacenter.Id)
	}

	if datacenter.Properties != nil {
		if datacenter.Properties.Location != nil {
			err := d.Set("location", *datacenter.Properties.Location)
			if err != nil {
				return fmt.Errorf("error while setting location property for datacenter %s: %s", d.Id(), err)
			}
		}

		if datacenter.Properties.Description != nil {
			err := d.Set("description", *datacenter.Properties.Description)
			if err != nil {
				return fmt.Errorf("error while setting description property for datacenter %s: %s", d.Id(), err)
			}
		}

		if datacenter.Properties.Name != nil {
			err := d.Set("name", *datacenter.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for datacenter %s: %s", d.Id(), err)
			}
		}

		if datacenter.Properties.Version != nil {
			err := d.Set("version", *datacenter.Properties.Version)
			if err != nil {
				return fmt.Errorf("error while setting version property for datacenter %s: %s", d.Id(), err)
			}
		}

		if datacenter.Properties.Features != nil && len(*datacenter.Properties.Features) > 0 {
			err := d.Set("features", *datacenter.Properties.Features)
			if err != nil {
				return fmt.Errorf("error while setting features property for datacenter %s: %s", d.Id(), err)
			}
		}

		if datacenter.Properties.SecAuthProtection != nil {
			err := d.Set("sec_auth_protection", *datacenter.Properties.SecAuthProtection)
			if err != nil {
				return fmt.Errorf("error while setting sec_auth_protection property for datacenter %s: %s", d.Id(), err)
			}
		}

	}

	return nil
}
