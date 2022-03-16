package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
)

func dataSourceLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocationRead,
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
			"image_aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	name, nameOk := d.GetOk("name")
	feature, featureOk := d.GetOk("feature")

	if !nameOk && !featureOk {
		return diag.FromErr(fmt.Errorf("either 'name' or 'feature' must be provided"))
	}

	request := client.LocationsApi.LocationsGet(ctx).Depth(1)

	if featureOk {
		request = request.Filter("features", feature.(string))
	}

	locations, apiResponse, err := request.Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching locations: %s", err.Error()))
	}
	var results []ionoscloud.Location

	if nameOk && locations.Items != nil {
		for _, loc := range *locations.Items {
			if loc.Properties != nil && loc.Properties.Name != nil && *loc.Properties.Name == name.(string) {
				results = append(results, loc)
			}
		}
	}

	log.Printf("[INFO] Results length %d *************", len(results))

	var location ionoscloud.Location

	if results == nil || len(results) == 0 {
		return diag.FromErr(fmt.Errorf("no location found with the specified criteria: name = %s, feature = %s", name.(string), feature.(string)))
	} else {
		location = results[0]
	}

	if err := setLocationData(d, &location); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setLocationData(d *schema.ResourceData, location *ionoscloud.Location) error {

	if location.Id != nil {
		d.SetId(*location.Id)
	}

	if location.Properties != nil {
		var cpuArchitectures []interface{}
		for _, cpuArchitecture := range *location.Properties.CpuArchitecture {
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

			cpuArchitectures = append(cpuArchitectures, architectureEntry)
		}

		if len(cpuArchitectures) > 0 {
			if err := d.Set("cpu_architecture", cpuArchitectures); err != nil {
				return fmt.Errorf("error while setting cpu_architecture property for datacenter %s: %s", d.Id(), err)
			}
		}

		var imageAliases []string
		for _, imageAlias := range *location.Properties.ImageAliases {
			imageAliases = append(imageAliases, imageAlias)
		}

		if len(imageAliases) > 0 {
			if err := d.Set("image_aliases", imageAliases); err != nil {
				return fmt.Errorf("error while setting image_aliases property for datacenter %s: %s", d.Id(), err)
			}
		}
	}

	return nil
}
