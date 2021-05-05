package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func dataSourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSnapshotRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
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

func dataSourceSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	snapshots, _, err := client.SnapshotsApi.SnapshotsGet(ctx).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while fetching IonosCloud locations %s", err)
	}

	name := d.Get("name").(string)
	location, locationOk := d.GetOk("location")
	size, sizeOk := d.GetOk("size")
	results := []ionoscloud.Snapshot{}

	for _, snp := range *snapshots.Items {
		if strings.Contains(strings.ToLower(*snp.Properties.Name), strings.ToLower(name)) {
			results = append(results, snp)
		}
	}

	if locationOk {
		locationResults := []ionoscloud.Snapshot{}
		for _, snp := range results {
			if *snp.Properties.Location == location.(string) {
				locationResults = append(locationResults, snp)
			}

		}
		results = locationResults
	}

	if sizeOk {
		sizeResults := []ionoscloud.Snapshot{}
		for _, snp := range results {
			if *snp.Properties.Size <= size.(float32) {
				sizeResults = append(sizeResults, snp)
			}

		}
		results = sizeResults
	}

	if len(results) > 1 {
		return fmt.Errorf("There is more than one snapshot that match the search criteria ")
	}

	if len(results) == 0 {
		return fmt.Errorf("There are no snapshots that match the search criteria ")
	}

	d.Set("name", results[0].Properties.Name)

	d.SetId(*results[0].Id)

	return nil
}
