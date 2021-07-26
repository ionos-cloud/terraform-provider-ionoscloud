package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func resourceSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnapshotCreate,
		ReadContext:   resourceSnapshotRead,
		UpdateContext: resourceSnapshotUpdate,
		DeleteContext: resourceSnapshotDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"volume_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"datacenter_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	volumeId := d.Get("volume_id").(string)
	name := d.Get("name").(string)

	rsp, apiResponse, err := client.VolumeApi.DatacentersVolumesCreateSnapshotPost(ctx, dcId, volumeId).Name(name).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a snapshot: %s", err))
		return diags
	}

	d.SetId(*rsp.Id)
	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceSnapshotRead(ctx, d, meta)
}

func resourceSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	rsp, apiResponse, err := client.SnapshotApi.SnapshotsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a snapshot ID %s %s", d.Id(), err))
		return diags
	}

	if err := d.Set("name", rsp.Properties.Name); err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	return nil
}

func resourceSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	name := d.Get("name").(string)
	input := ionoscloud.SnapshotProperties{
		Name: &name,
	}

	_, apiResponse, err := client.SnapshotApi.SnapshotsPatch(context.TODO(), d.Id()).Snapshot(input).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while restoring a snapshot ID %s %d", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceSnapshotRead(ctx, d, meta)
}

func resourceSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	rsp, apiResponse, err := client.SnapshotApi.SnapshotsFindById(ctx, d.Id()).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a snapshot ID %s %s", d.Id(), err))
		return diags
	}

	for *rsp.Metadata.State != "AVAILABLE" {
		time.Sleep(30 * time.Second)
		_, _, err := client.SnapshotApi.SnapshotsFindById(ctx, d.Id()).Execute()

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while fetching a snapshot ID %s %s", d.Id(), err))
			return diags
		}
	}

	dcId := d.Get("datacenter_id").(string)
	dc, _, err := client.DataCenterApi.DatacentersFindById(ctx, dcId).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a Datacenter ID %s %s", dcId, err))
		return diags
	}

	for *dc.Metadata.State != "AVAILABLE" {
		time.Sleep(30 * time.Second)
		_, _, err := client.DataCenterApi.DatacentersFindById(ctx, dcId).Execute()

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while fetching a Datacenter ID %s %s", dcId, err))
			return diags
		}
	}

	_, apiResponse, err = client.SnapshotApi.SnapshotsDelete(ctx, d.Id()).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a snapshot ID %s %s", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	d.SetId("")
	return nil
}
