package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"sec_auth_protection": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func getDatacenter(client *ionoscloud.APIClient, d *schema.ResourceData) (*ionoscloud.Datacenter, error) {
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
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		dc, _, err := client.DataCentersApi.DatacentersFindById(ctx, id.(string)).Execute()
		if err != nil {
			return nil, fmt.Errorf("Error getting datacenter with id %s", id.(string))
		}
		if nameOk {
			if *dc.Properties.Name != name {
				return nil, fmt.Errorf("[ERROR] Name of dc (UUID=%s, name=%s) does not match expected name: %s",
					*dc.Id, *dc.Properties.Name, name)
			}
		}
		if locationOk {
			if *dc.Properties.Location != location {
				return nil, fmt.Errorf("[ERROR] location of dc (UUID=%s, location=%s) does not match expected location: %s",
					*dc.Id, *dc.Properties.Location, location)

			}
		}
		log.Printf("[INFO] Got dc [Name=%s, Location=%s]", *dc.Properties.Name, *dc.Properties.Location)
		return &dc, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	datacenters, _, err := client.DataCentersApi.DatacentersGet(ctx).Execute()

	if err != nil {
		return nil, fmt.Errorf("An error occured while fetching datacenters %s", err)
	}

	results := []ionoscloud.Datacenter{}

	if datacenters.Items != nil {
		for _, dc := range *datacenters.Items {
			if *dc.Properties.Name == name || strings.Contains(*dc.Properties.Name, name) {
				results = append(results, dc)
			}
		}
	}

	if locationOk {
		log.Printf("[INFO] searching dcs by location***********")
		locationResults := []ionoscloud.Datacenter{}
		for _, dc := range results {
			if *dc.Properties.Location == location {
				locationResults = append(locationResults, dc)
			}
		}
		results = locationResults
	}
	log.Printf("[INFO] Results length %d *************", len(results))

	if len(results) > 1 {
		log.Printf("[INFO] Results length greater than 1")
		return nil, fmt.Errorf("There is more than one datacenters that match the search criteria")
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("There are no datacenters that match the search criteria")
	}
	return &results[0], nil
}

func dataSourceDataCenterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	datacenter, err := getDatacenter(client, d)

	if err != nil {
		return fmt.Errorf("An error occured while fetching datacenters %s", err)
	}
	d.SetId(*datacenter.Id)

	if datacenter.Properties.Location != nil {
		err := d.Set("location", datacenter.Properties.Location)
		if err != nil {
			return fmt.Errorf("Error while setting location property for datacenter %s: %s", d.Id(), err)
		}
	}

	if datacenter.Properties.Name != nil {
		err := d.Set("name", datacenter.Properties.Name)
		if err != nil {
			return fmt.Errorf("Error while setting name property for datacenter %s: %s", d.Id(), err)
		}
	}

	if datacenter.Properties.SecAuthProtection != nil {
		err := d.Set("sec_auth_protection", datacenter.Properties.SecAuthProtection)
		if err != nil {
			return fmt.Errorf("Error while setting sec_auth_protection property for datacenter %s: %s", d.Id(), err)
		}
	}

	return nil
}
