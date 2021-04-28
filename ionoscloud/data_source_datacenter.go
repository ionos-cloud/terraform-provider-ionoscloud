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
		err := d.Set("location", *datacenter.Properties.Location)
		if err != nil {
			return fmt.Errorf("Error while setting location property for datacenter %s: %s", d.Id(), err)
		}
	}

	if datacenter.Properties.Name != nil {
		err := d.Set("name", *datacenter.Properties.Name)
		if err != nil {
			return fmt.Errorf("Error while setting name property for datacenter %s: %s", d.Id(), err)
		}
	}

	if datacenter.Properties.Version != nil {
		err := d.Set("version", *datacenter.Properties.Version)
		if err != nil {
			return fmt.Errorf("Error while setting version property for datacenter %s: %s", d.Id(), err)
		}
	}

	if datacenter.Properties.Features != nil && len(*datacenter.Properties.Features) > 0 {
		err := d.Set("features", *datacenter.Properties.Features)
		if err != nil {
			return fmt.Errorf("Error while setting features property for datacenter %s: %s", d.Id(), err)
		}
	}

	if datacenter.Properties.SecAuthProtection != nil {
		err := d.Set("sec_auth_protection", *datacenter.Properties.SecAuthProtection)
		if err != nil {
			return fmt.Errorf("Error while setting sec_auth_protection property for datacenter %s: %s", d.Id(), err)
		}
	}

	cpuArchitectures := make([]interface{}, 0)
	if datacenter.Properties.CpuArchitecture != nil && len(*datacenter.Properties.CpuArchitecture) > 0 {
		cpuArchitectures = make([]interface{}, len(*datacenter.Properties.CpuArchitecture))
		for index, cpuArchitecture := range *datacenter.Properties.CpuArchitecture {
			architectureEntry := make(map[string]interface{})

			if cpuArchitecture.CpuFamily != nil {
				architectureEntry["cpu_family"] = *cpuArchitecture.CpuFamily
			}

			if cpuArchitecture.MaxCores != nil {
				architectureEntry["max_cores"] = *cpuArchitecture.MaxCores
			}

			if cpuArchitecture.MaxRam != nil {
				architectureEntry["max_ram"] = *cpuArchitecture.MaxRam
			}

			if cpuArchitecture.Vendor != nil {
				architectureEntry["vendor"] = *cpuArchitecture.Vendor
			}

			cpuArchitectures[index] = architectureEntry
		}
	}
	if len(cpuArchitectures) > 0 {
		if err := d.Set("cpu_architecture", cpuArchitectures); err != nil {
			return fmt.Errorf("Error while setting cpu_architecture property for datacenter %s: %s", d.Id(), err)
		}
	}

	return nil
}
