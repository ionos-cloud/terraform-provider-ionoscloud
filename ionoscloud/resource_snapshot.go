package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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
		CustomizeDiff: checkSnapshotUpdateOnlyAttrs,
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
				Optional:    true,
				Description: "Human readable description",
			},
			"licence_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "OS type of this Snapshot",
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the image in GB",
			},
			"sec_auth_protection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Boolean value representing if the snapshot requires extra protection e.g. two factor protection",
			},
			"cpu_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cpu_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ram_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nic_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"nic_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disc_virtio_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disc_virtio_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"require_legacy_bios": {
				Type:        schema.TypeBool,
				Description: "Indicates if the image requires the legacy BIOS for compatibility or specific needs.",
				Optional:    true,
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	dcId := d.Get("datacenter_id").(string)
	volumeId := d.Get("volume_id").(string)
	name := d.Get("name").(string)
	snapshot := ionoscloud.NewCreateSnapshot()
	snapshot.Properties = ionoscloud.NewCreateSnapshotProperties()
	props := snapshot.Properties
	props.Name = new(name)
	if v, ok := d.GetOk("description"); ok {
		props.Description = new(v.(string))
	}
	if v, ok := d.GetOk("licence_type"); ok {
		props.LicenceType = new(v.(string))
	}
	if v, ok := d.GetOk("sec_auth_protection"); ok {
		props.SecAuthProtection = new(v.(bool))
	}
	rsp, apiResponse, err := client.VolumesApi.DatacentersVolumesCreateSnapshotPost(ctx, dcId, volumeId).Snapshot(*snapshot).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating a snapshot: %w ", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	d.SetId(*rsp.Id)
	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			d.SetId("")
		}
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceSnapshotRead(ctx, d, meta)
}

func resourceSnapshotRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	snapshot, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error occurred while fetching a snapshot: %w", err), nil)
	}

	if err = setSnapshotData(d, &snapshot); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourceSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	// use ionoscloud.SnapshotProperties struct to initialize update input instead of ionoscloud.NewSnapshotProperties(),
	// because `require_legacy_bios` is set to `true` by default in the latter, which can cause unintended updates
	input := ionoscloud.SnapshotProperties{}
	input.Name = &name

	if d.HasChange("description") {
		input.Description = new(d.Get("description").(string))
	}
	if d.HasChange("licence_type") {
		input.LicenceType = new(d.Get("licence_type").(string))
	}
	if d.HasChange("sec_auth_protection") {
		input.SecAuthProtection = new(d.Get("sec_auth_protection").(bool))
	}
	if d.HasChange("cpu_hot_plug") {
		input.CpuHotPlug = new(d.Get("cpu_hot_plug").(bool))
	}
	if d.HasChange("nic_hot_plug") {
		input.NicHotPlug = new(d.Get("nic_hot_plug").(bool))
	}
	if d.HasChange("nic_hot_unplug") {
		input.NicHotUnplug = new(d.Get("nic_hot_unplug").(bool))
	}
	if d.HasChange("ram_hot_plug") {
		input.RamHotPlug = new(d.Get("ram_hot_plug").(bool))
	}
	if d.HasChange("disc_virtio_hot_plug") {
		input.DiscVirtioHotPlug = new(d.Get("disc_virtio_hot_plug").(bool))
	}
	if d.HasChange("disc_virtio_hot_unplug") {
		input.DiscVirtioHotUnplug = new(d.Get("disc_virtio_hot_unplug").(bool))
	}
	if d.HasChange("require_legacy_bios") {
		input.RequireLegacyBios = new(d.Get("require_legacy_bios").(bool))
	}

	_, apiResponse, err := client.SnapshotsApi.SnapshotsPatch(ctx, d.Id()).Snapshot(input).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while restoring a snapshot: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceSnapshotRead(ctx, d, meta)
}

func resourceSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	apiResponse, err := client.SnapshotsApi.SnapshotsDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while deleting a snapshot: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	d.SetId("")
	return nil
}

func resourceSnapshotImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	importID := d.Id()

	location, parts := splitImportID(importID, ":")
	if len(parts) != 1 {
		return nil, fmt.Errorf("invalid import identifier: expected one of <location>:<snapshot-id> or <snapshot-id>, got: %s", importID)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, fmt.Errorf("failed validating import identifier %q: %w", importID, err)
	}

	snapshotID := parts[0]
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return nil, err
	}

	snapshot, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, snapshotID).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find snapshot %q", snapshotID), nil)
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while retrieving the snapshot %q, %w", snapshotID, err), nil)
	}

	log.Printf("[INFO] snapshot %s found: %+v", importID, snapshot)

	if err = setSnapshotData(d, &snapshot); err != nil {
		return nil, diagutil.ToError(d, err, nil)
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

		if snapshot.Properties.RequireLegacyBios != nil {
			if err := d.Set("require_legacy_bios", *snapshot.Properties.RequireLegacyBios); err != nil {
				return err
			}
		}
	}
	return nil
}

// check that update-only attributes are not explicitly
// set during snapshot creation and return an error if any are found.
func checkSnapshotUpdateOnlyAttrs(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	// Only validate during creation (no ID yet)
	if diff.Id() != "" {
		return nil
	}

	updateOnlyAttrs := []string{
		"cpu_hot_plug", "ram_hot_plug", "nic_hot_plug", "nic_hot_unplug",
		"disc_virtio_hot_plug", "disc_virtio_hot_unplug",
	}

	rawConfig := diff.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil
	}
	configMap := rawConfig.AsValueMap()

	var invalidAttrs []string
	for _, attr := range updateOnlyAttrs {
		if v, ok := configMap[attr]; ok && !v.IsNull() {
			invalidAttrs = append(invalidAttrs, attr)
		}
	}

	if len(invalidAttrs) > 0 {
		return fmt.Errorf(
			"the following attributes cannot be set during snapshot creation and can only be "+
				"updated after the snapshot exists: %s",
			strings.Join(invalidAttrs, ", "),
		)
	}

	return nil
}
