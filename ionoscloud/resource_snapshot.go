package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func resourceSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnapshotCreate,
		ReadContext:   resourceSnapshotRead,
		UpdateContext: resourceSnapshotUpdate,
		DeleteContext: resourceSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSnapshotImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "A name of that resource",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human readable description",
			},
			"licence_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS type of this Snapshot",
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location of that image/snapshot",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the image in GB",
			},
			"sec_auth_protection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Boolean value representing if the snapshot requires extra protection e.g. two factor protection",
			},
			"cpu_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cpu_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nic_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nic_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_virtio_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_virtio_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_scsi_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_scsi_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"volume_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	volumeId := d.Get("volume_id").(string)
	name := d.Get("name").(string)

	rsp, apiResponse, err := client.VolumesApi.DatacentersVolumesCreateSnapshotPost(ctx, dcId, volumeId).Name(name).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a snapshot: %w ", err))
		return diags
	}

	d.SetId(*rsp.Id)
	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if cloudapi.IsRequestFailed(errState) {
			d.SetId("")
		}
		return diag.FromErr(errState)
	}

	return resourceSnapshotRead(ctx, d, meta)
}

func resourceSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	snapshot, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a snapshot ID %s %w", d.Id(), err))
		return diags
	}

	if err = setSnapshotData(d, &snapshot); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	name := d.Get("name").(string)
	input := ionoscloud.SnapshotProperties{
		Name: &name,
	}

	_, apiResponse, err := client.SnapshotsApi.SnapshotsPatch(context.TODO(), d.Id()).Snapshot(input).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while restoring a snapshot ID %s %d", d.Id(), err))
		return diags
	}

	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return diag.FromErr(errState)
	}

	return resourceSnapshotRead(ctx, d, meta)
}

func resourceSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	rsp, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a snapshot ID %s %w", d.Id(), err))
		return diags
	}

	for !strings.EqualFold(*rsp.Metadata.State, constant.Available) {
		time.Sleep(30 * time.Second)
		_, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, d.Id()).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while fetching a snapshot ID %s %w", d.Id(), err))
			return diags
		}
	}

	dcId := d.Get("datacenter_id").(string)
	dc, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, dcId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a Datacenter ID %s %s", dcId, err))
		return diags
	}

	for !strings.EqualFold(*dc.Metadata.State, constant.Available) {
		time.Sleep(30 * time.Second)
		_, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, dcId).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while fetching a Datacenter ID %s %s", dcId, err))
			return diags
		}
	}

	apiResponse, err = client.SnapshotsApi.SnapshotsDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a snapshot ID %s %w", d.Id(), err))
		return diags
	}

	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		return diag.FromErr(errState)
	}

	d.SetId("")
	return nil
}

func resourceSnapshotImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).CloudApiClient

	snapshotId := d.Id()
	snapshot, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find snapshot %q", snapshotId)
		}
		return nil, fmt.Errorf("an error occured while retrieving the snapshot %q, %q", snapshotId, err)
	}

	log.Printf("[INFO] snapshot %s found: %+v", d.Id(), snapshot)

	if err = setSnapshotData(d, &snapshot); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func setSnapshotData(d *schema.ResourceData, snapshot *ionoscloud.Snapshot) error {
	d.SetId(*snapshot.Id)

	if snapshot.Properties != nil {
		if snapshot.Properties.Name != nil {
			if err := d.Set("name", *snapshot.Properties.Name); err != nil {
				return err
			}
		}

		if snapshot.Properties.Location != nil {
			if err := d.Set("location", *snapshot.Properties.Location); err != nil {
				return err
			}
		}

		if snapshot.Properties.Size != nil {
			if err := d.Set("size", *snapshot.Properties.Size); err != nil {
				return err
			}
		}

		if snapshot.Properties.Description != nil {
			if err := d.Set("description", *snapshot.Properties.Description); err != nil {
				return err
			}
		}

		if snapshot.Properties.LicenceType != nil {
			if err := d.Set("licence_type", *snapshot.Properties.LicenceType); err != nil {
				return err
			}
		}

		if snapshot.Properties.SecAuthProtection != nil {
			if err := d.Set("sec_auth_protection", *snapshot.Properties.SecAuthProtection); err != nil {
				return err
			}
		}

		if snapshot.Properties.CpuHotPlug != nil {
			if err := d.Set("cpu_hot_plug", *snapshot.Properties.CpuHotPlug); err != nil {
				return err
			}
		}

		if snapshot.Properties.CpuHotUnplug != nil {
			if err := d.Set("cpu_hot_unplug", *snapshot.Properties.CpuHotUnplug); err != nil {
				return err
			}
		}

		if snapshot.Properties.RamHotPlug != nil {
			if err := d.Set("ram_hot_plug", *snapshot.Properties.RamHotPlug); err != nil {
				return err
			}
		}

		if snapshot.Properties.RamHotUnplug != nil {
			if err := d.Set("ram_hot_unplug", *snapshot.Properties.RamHotUnplug); err != nil {
				return err
			}
		}

		if snapshot.Properties.NicHotPlug != nil {
			if err := d.Set("nic_hot_plug", *snapshot.Properties.NicHotPlug); err != nil {
				return err
			}
		}

		if snapshot.Properties.NicHotUnplug != nil {
			if err := d.Set("nic_hot_unplug", *snapshot.Properties.NicHotUnplug); err != nil {
				return err
			}
		}

		if snapshot.Properties.DiscVirtioHotPlug != nil {
			if err := d.Set("disc_virtio_hot_plug", *snapshot.Properties.DiscVirtioHotPlug); err != nil {
				return err
			}
		}

		if snapshot.Properties.DiscVirtioHotUnplug != nil {
			if err := d.Set("disc_virtio_hot_unplug", *snapshot.Properties.DiscVirtioHotUnplug); err != nil {
				return err
			}
		}

		if snapshot.Properties.DiscScsiHotPlug != nil {
			if err := d.Set("disc_scsi_hot_plug", *snapshot.Properties.DiscVirtioHotUnplug); err != nil {
				return err
			}
		}

		if snapshot.Properties.DiscScsiHotUnplug != nil {
			if err := d.Set("disc_scsi_hot_unplug", *snapshot.Properties.DiscVirtioHotUnplug); err != nil {
				return err
			}
		}
	}
	return nil
}
