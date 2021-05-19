package ionoscloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
	"log"
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
	client := meta.(*profitbricks.Client)

	locations, err := client.ListLocations()

	if err != nil {
		return fmt.Errorf("An error occured while fetching IonosCloud locations %s", err)
	}

	name, nameOk := d.GetOk("name")
	feature, featureOk := d.GetOk("features")

	if !nameOk && !featureOk {
		return fmt.Errorf("Either 'name' or 'feature' must be provided.")
	}
	results := []profitbricks.Location{}

	for _, loc := range locations.Items {
		if loc.Properties.Name == name.(string)  {
			results = append(results, loc)
		}
	}

	if featureOk {
		locationResults := []profitbricks.Location{}
		for _, loc := range results {
			for _, f := range loc.Properties.Features {
				if f == feature.(string) {
					locationResults = append(locationResults, loc)
				}
			}
		}
		results = locationResults
	}
	log.Printf("[INFO] Results length %d *************", len(results))

	if len(results) == 0 {
		return fmt.Errorf("There are no locations that match the search criteria")
	}

	d.SetId(results[0].ID)

	return nil
}
