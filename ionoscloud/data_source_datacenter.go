package ionoscloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
	"log"
)

func dataSourceDataCenter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDataCenterRead,
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func getDatacenter(client *profitbricks.Client, d *schema.ResourceData) (*profitbricks.Datacenter, error) {
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
	if !idOk && !nameOk && !locationOk {
		return nil, fmt.Errorf("either id, location or name must be set")
	}
	if idOk {
		dc, err := client.GetDatacenter(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error getting datacenter with id %s", id.(string))
		}
		if nameOk {
			if dc.Properties.Name != name {
				return nil, fmt.Errorf("name of dc (UUID=%s, name=%s) does not match expected name: %s",
					dc.ID, dc.Properties.Name, name)
			}
		}
		if locationOk {
			if dc.Properties.Location != location {
				return nil, fmt.Errorf("location of dc (UUID=%s, location=%s) does not match expected location: %s",
					dc.ID, dc.Properties.Location, location)

			}
		}
		log.Printf("[INFO] Got dc [Name=%s, Location=%s]", dc.Properties.Name, dc.Properties.Location)
		return dc, nil
	}

	datacenters, err := client.ListDatacenters()

	if err != nil {
		return nil, fmt.Errorf("an error occured while fetching datacenters: %s", err)
	}

	var results []profitbricks.Datacenter

	if nameOk {
		for _, dc := range datacenters.Items {
			if dc.Properties.Name == name {
				results = append(results, dc)
			}
		}

		if results == nil {
			return nil, fmt.Errorf("could not find a datacenter with name %s", name)
		}

	}

	if locationOk {

		if results != nil {
			for _, dc := range results {
				if dc.Properties.Location == location {
					return &dc, nil
				}
			}
			return nil, fmt.Errorf("no datacenter with name %s and location %s was found", name, location)
		} else {
			/* find the first datacenter matching the location */
			for _, dc := range datacenters.Items {
				if dc.Properties.Location == location {
					return &dc, nil
				}
			}
		}
	}

	if results == nil {
		return nil, fmt.Errorf("there are no datacenters that match the search criteria")
	}

	return &results[0], nil
}
func dataSourceDataCenterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	datacenter, err := getDatacenter(client, d)

	if err != nil {
		return fmt.Errorf("an error occured while fetching datacenters: %s", err)
	}
	d.SetId(datacenter.ID)
	err = d.Set("location", datacenter.Properties.Location)
	if err != nil {
		return err
	}
	err = d.Set("name", datacenter.Properties.Name)
	if err != nil {
		return err
	}

	return nil
}
