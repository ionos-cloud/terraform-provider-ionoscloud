package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSnapshotRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	snapshots, _, err := client.SnapshotsApi.SnapshotsGet(ctx).Execute()

	if err != nil {
		return fmt.Errorf("an error occured while fetching IonosCloud locations %s", err)
	}

	name := d.Get("name").(string)
	location, locationOk := d.GetOk("location")
	size, sizeOk := d.GetOk("size")
	var results []ionoscloud.Snapshot

	if snapshots.Items != nil {
		for _, snp := range *snapshots.Items {
			if snp.Properties.Name != nil && *snp.Properties.Name == name {
				results = append(results, snp)
			}
		}
	}

	if locationOk {
		var locationResults []ionoscloud.Snapshot
		for _, snp := range results {
			if *snp.Properties.Location == location.(string) {
				locationResults = append(locationResults, snp)
			}

		}
		results = locationResults
	}

	if sizeOk {
		var sizeResults []ionoscloud.Snapshot
		for _, snp := range results {
			if *snp.Properties.Size <= size.(float32) {
				sizeResults = append(sizeResults, snp)
			}

		}
		results = sizeResults
	}

	if len(results) == 0 {
		return fmt.Errorf("There are no snapshots that match the search criteria ")
	}

	err = d.Set("name", results[0].Properties.Name)
	if err != nil {
		return err
	}

	d.SetId(*results[0].Id)

	return nil
}
