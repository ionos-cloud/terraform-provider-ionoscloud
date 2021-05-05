package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func dataSourceLocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"feature": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourceLocationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	locations, _, err := client.LocationsApi.LocationsGet(ctx).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while fetching IonosCloud locations %s", err)
	}

	name, nameOk := d.GetOk("name")
	feature, featureOk := d.GetOk("features")

	if !nameOk && !featureOk {
		return fmt.Errorf("Either 'name' or 'feature' must be provided.")
	}
	results := []ionoscloud.Location{}

	for _, loc := range *locations.Items {
		if *loc.Properties.Name == name.(string) || strings.Contains(*loc.Properties.Name, name.(string)) {
			results = append(results, loc)
		}
	}

	if featureOk {
		locationResults := []ionoscloud.Location{}
		for _, loc := range results {
			for _, f := range *loc.Properties.Features {
				if f == feature.(string) {
					locationResults = append(locationResults, loc)
				}
			}
		}
		results = locationResults
	}
	log.Printf("[INFO] Results length %d *************", len(results))

	cpuArchitectures := make([]interface{}, 0)
	for _, loc := range results {
		cpuArchitectures = make([]interface{}, len(*loc.Properties.CpuArchitecture))
		for index, cpuArchitecture := range *loc.Properties.CpuArchitecture {
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

	if len(results) > 1 {
		log.Printf("[INFO] Results length greater than 1")
		return fmt.Errorf("there is more than one location that match the search criteria")
	}

	if len(results) == 0 {
		return fmt.Errorf("There are no locations that match the search criteria")
	}

	d.SetId(*results[0].Id)

	return nil
}
