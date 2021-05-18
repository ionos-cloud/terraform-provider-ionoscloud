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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceLocationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	locations, _, err := client.LocationApi.LocationsGet(ctx).Execute()

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
