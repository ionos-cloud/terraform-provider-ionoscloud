package ionoscloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func resourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceSnapshotCreate,
		Read:   resourceSnapshotRead,
		Update: resourceSnapshotUpdate,
		Delete: resourceSnapshotDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client
	dcId := d.Get("datacenter_id").(string)
	volumeId := d.Get("volume_id").(string)
	//name := d.Get("name").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)
	if cancel != nil {
		defer cancel()
	}
	rsp, apiResponse, err := client.VolumesApi.DatacentersVolumesCreateSnapshotPost(ctx, dcId, volumeId).Execute()
	if err != nil {
		return fmt.Errorf("An error occured while creating a snapshot: %s ", err)
	}

	d.SetId(*rsp.Id)
	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceSnapshotRead(d, meta)
}

func resourceSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	rsp, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error occured while fetching a snapshot ID %s %s", d.Id(), err)
	}

	d.Set("name", rsp.Properties.Name)
	return nil
}

func resourceSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	name := d.Get("name").(string)
	input := ionoscloud.SnapshotProperties{
		Name: &name,
	}

	_, apiResponse, err := client.SnapshotsApi.SnapshotsPatch(context.TODO(), d.Id()).Snapshot(input).Execute()
	if err != nil {
		return fmt.Errorf("An error occured while restoring a snapshot ID %s %d", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceSnapshotRead(d, meta)
}

func resourceSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}
	rsp, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, d.Id()).Execute()
	if err != nil {
		return fmt.Errorf("An error occured while fetching a snapshot ID %s %s", d.Id(), err)
	}

	for *rsp.Metadata.State != "AVAILABLE" {
		time.Sleep(30 * time.Second)
		_, _, err := client.SnapshotsApi.SnapshotsFindById(ctx, d.Id()).Execute()

		if err != nil {
			return fmt.Errorf("An error occured while fetching a snapshot ID %s %s", d.Id(), err)
		}
	}

	dcId := d.Get("datacenter_id").(string)
	dc, _, err := client.DataCentersApi.DatacentersFindById(ctx, dcId).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while fetching a Datacenter ID %s %s", dcId, err)
	}

	for *dc.Metadata.State != "AVAILABLE" {
		time.Sleep(30 * time.Second)
		_, _, err := client.DataCentersApi.DatacentersFindById(ctx, dcId).Execute()

		if err != nil {
			return fmt.Errorf("An error occured while fetching a Datacenter ID %s %s", dcId, err)
		}
	}

	_, apiResponse, err = client.SnapshotsApi.SnapshotsDelete(ctx, d.Id()).Execute()
	if err != nil {
		return fmt.Errorf("An error occured while deleting a snapshot ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
