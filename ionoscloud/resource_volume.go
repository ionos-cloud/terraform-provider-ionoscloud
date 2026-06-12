package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapiimage"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapilocation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const HDDImage = "HDD"

func resourceVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeCreate,
		ReadContext:   resourceVolumeRead,
		UpdateContext: resourceVolumeUpdate,
		DeleteContext: resourceVolumeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVolumeImporter,
		},
		CustomizeDiff: checkVolumeImmutableFields,
		Schema: map[string]*schema.Schema{
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old != "" {
						return true
					}
					return false
				},
			},
			"image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"disk_type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"image_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"licence_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssh_key_path": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"ssh_keys": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"sshkey": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bus": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2", "ZONE_3"}, true)),
			},
			"cpu_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_plug": {
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
			"backup_unit_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pci_slot": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"boot_server": {
				Type:        schema.TypeString,
				Description: "The UUID of the attached server.",
				Computed:    true,
			},
			"server_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				ForceNew:    true,
			},
			"expose_serial": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "If set to `true` will expose the serial id of the disk attached to the server. " +
					"If set to `false` will not expose the serial id. Some operating systems or software solutions require the serial id to be exposed to work properly. " +
					"Exposing the serial can influence licensed software (e.g. Windows) behavior",
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

func checkVolumeImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ any) error {

	// we do not want to check in case of resource creation
	if diff.Id() == "" {
		return nil
	}

	if diff.HasChange("availability_zone") {
		return fmt.Errorf("availability_zone %s", ImmutableError)
	}

	if diff.HasChange("user_data") {
		return fmt.Errorf("user_data %s", ImmutableError)
	}

	if diff.HasChange("backup_unit_id") {
		return fmt.Errorf("backup_unit_id %s", ImmutableError)
	}

	if diff.HasChange("image_name") {
		return fmt.Errorf("image_name %s", ImmutableError)
	}

	if diff.HasChange("disk_type") {
		return fmt.Errorf("disk_type %s", ImmutableError)
	}

	if diff.HasChange("availability_zone") {
		return fmt.Errorf("availability_zone %s", ImmutableError)
	}
	return nil

}

func resourceVolumeCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var image, imageAlias string

	dcID := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	// create volume object with data to be used for image
	volumeProperties, err := getVolumeData(ctx, d, "", "")
	if err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	volume := ionoscloud.Volume{
		Properties: volumeProperties,
	}
	image, imageAlias, err = getImage(ctx, client, d, *volume.Properties)
	if err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	if image != "" {
		volume.Properties.Image = &image
	} else {
		volume.Properties.Image = nil
	}

	if imageAlias != "" {
		volume.Properties.ImageAlias = &imageAlias
	} else {
		volume.Properties.ImageAlias = nil
	}

	if userData, ok := d.GetOk("user_data"); ok {
		if image == "" && imageAlias == "" {
			return diagutil.ToDiags(d, fmt.Errorf("it is mandatory to provide either public image that has cloud-init compatibility in conjunction with user_data property "), nil)
		} else {
			userData := userData.(string)
			volume.Properties.UserData = &userData
		}
	}

	if backupUnitID, ok := d.GetOk("backup_unit_id"); ok {
		if utils.IsValidUUID(backupUnitID.(string)) {
			if image == "" && imageAlias == "" {
				return diagutil.ToDiags(d, fmt.Errorf("it is mandatory to provide either public image that has cloud-init compatibility in conjunction with backup_unit_id property "), nil)
			} else {
				backupUnitID := backupUnitID.(string)
				volume.Properties.BackupunitId = &backupUnitID
			}
		} else {
			return diagutil.ToDiags(d, fmt.Errorf("the backup_unit_id that you specified is not a valid UUID"), nil)
		}
	}

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesPost(ctx, dcID).Volume(volume).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating a volume: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	d.SetId(*volume.Id)

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			d.SetId("")
		}
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	volumeToAttach := ionoscloud.Volume{Id: volume.Id}
	volume, apiResponse, err = client.ServersApi.DatacentersServersVolumesPost(ctx, dcID, serverID).Volume(volumeToAttach).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while attaching a volume dcID: %s server_id: %s ID: %s Response: %w", dcID, serverID, *volumeToAttach.Id, err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	sErr := d.Set("server_id", serverID)

	if sErr != nil {
		return diagutil.ToDiags(d, fmt.Errorf("error while setting serverID %s: %w", serverID, sErr), nil)
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			if sErr := d.Set("server_id", ""); sErr != nil {
				return diagutil.ToDiags(d, fmt.Errorf("error while setting serverID: %w", sErr), nil)
			}
		}
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceVolumeRead(ctx, d, meta)
}

func resourceVolumeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	dcID := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)
	volumeID := d.Id()
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, dcID, volumeID).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error occurred while fetching volume: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	_, apiResponse, err = client.ServersApi.DatacentersServersVolumesFindById(ctx, dcID, serverID, volumeID).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if err2 := d.Set("server_id", ""); err2 != nil {
			requestLocation, _ := apiResponse.SafeLocation()
			return diagutil.ToDiags(d, err2, &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
		}
	}

	if err := setVolumeData(d, &volume); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourceVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	properties := ionoscloud.VolumeProperties{}
	dcID := d.Get("datacenter_id").(string)
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		_, newValue := d.GetChange("name")
		newValueStr := newValue.(string)
		properties.Name = &newValueStr
	}

	if d.HasChange("size") {
		_, newValue := d.GetChange("size")
		newValueFloat32 := float32(newValue.(int))
		properties.Size = &newValueFloat32
	}
	if d.HasChange("bus") {
		_, newValue := d.GetChange("bus")
		newValueStr := newValue.(string)
		properties.Bus = &newValueStr
	}

	if d.HasChange("expose_serial") {
		_, n := d.GetChange("expose_serial")
		nBool := n.(bool)
		properties.ExposeSerial = &nBool
	}

	if d.HasChange("require_legacy_bios") {
		_, newValue := d.GetChange("require_legacy_bios")
		requireLegacyBios := newValue.(bool)
		properties.RequireLegacyBios = &requireLegacyBios
	}

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesPatch(ctx, dcID, d.Id()).Volume(properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating volume: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})

	}

	// Wait, catching any errors
	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	if apiResponse.SafeStatusCode() > 299 {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating a volume, status code: %d", apiResponse.SafeStatusCode()), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	if d.HasChange("server_id") {
		_, newValue := d.GetChange("server_id")
		serverID := newValue.(string)
		volumeToAttach := ionoscloud.Volume{Id: volume.Id}
		_, apiResponse, err := client.ServersApi.DatacentersServersVolumesPost(ctx, dcID, serverID).Volume(volumeToAttach).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			requestLocation, _ := apiResponse.SafeLocation()
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while attaching a volume dcID: %s server_id: %s ID: %s Response: %w",
				dcID, serverID, *volume.Id, err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
		}

		if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
			requestLocation, _ := apiResponse.SafeLocation()
			return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
		}
	}

	return resourceVolumeRead(ctx, d, meta)
}

func resourceVolumeDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	dcID := d.Get("datacenter_id").(string)
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	apiResponse, err := client.VolumesApi.DatacentersVolumesDelete(ctx, dcID, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while deleting a volume: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})

	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	d.SetId("")
	return nil
}

func resourceVolumeImporter(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	importID := d.Id()

	location, parts := splitImportID(importID, "/")
	if len(parts) != 3 {
		return nil, diagutil.ToError(d, fmt.Errorf(
			"invalid import identifier: expected one of <location>:<datacenter-id>/<server-id>/<volume-id> "+
				"or <datacenter-id>/<server-id>/<volume-id>, got: %s", importID,
		), nil)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("failed validating import identifier %q: %w", importID, err), nil)
	}

	dcID := parts[0]
	srvID := parts[1]
	volumeID := parts[2]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return nil, err
	}

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, dcID, volumeID).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("volume does not exist %q", volumeID), nil)
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while trying to find the volume %q", volumeID), nil)
	}

	tflog.Info(ctx, "volume found", map[string]any{"volume_id": *volume.Id, "datacenter_id": dcID})

	d.SetId(*volume.Id)
	if err := d.Set("datacenter_id", dcID); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	if err := d.Set("server_id", srvID); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	if err := d.Set("location", location); err != nil {
		return nil, err
	}

	if err := setVolumeData(d, &volume); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}

func setVolumeData(d *schema.ResourceData, volume *ionoscloud.Volume) error {

	if volume.Id != nil {
		d.SetId(*volume.Id)
	}

	if volume.Properties.Name != nil {
		err := d.Set("name", *volume.Properties.Name)
		if err != nil {
			return fmt.Errorf("error while setting name property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.Type != nil {
		err := d.Set("disk_type", *volume.Properties.Type)
		if err != nil {
			return fmt.Errorf("error while setting type property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.Size != nil {
		err := d.Set("size", *volume.Properties.Size)
		if err != nil {
			return fmt.Errorf("error while setting size property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.Bus != nil {
		err := d.Set("bus", *volume.Properties.Bus)
		if err != nil {
			return fmt.Errorf("error while setting bus property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.Image != nil {
		err := d.Set("image", *volume.Properties.Image)
		if err != nil {
			return fmt.Errorf("error while setting image property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.AvailabilityZone != nil {
		err := d.Set("availability_zone", *volume.Properties.AvailabilityZone)
		if err != nil {
			return fmt.Errorf("error while setting availability_zone property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.CpuHotPlug != nil {
		err := d.Set("cpu_hot_plug", *volume.Properties.CpuHotPlug)
		if err != nil {
			return fmt.Errorf("error while setting cpu_hot_plug property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.RamHotPlug != nil {
		err := d.Set("ram_hot_plug", *volume.Properties.RamHotPlug)
		if err != nil {
			return fmt.Errorf("error while setting ram_hot_plug property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.NicHotPlug != nil {
		err := d.Set("nic_hot_plug", *volume.Properties.NicHotPlug)
		if err != nil {
			return fmt.Errorf("error while setting nic_hot_plug property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.NicHotUnplug != nil {
		err := d.Set("nic_hot_unplug", *volume.Properties.NicHotUnplug)
		if err != nil {
			return fmt.Errorf("error while setting nic_hot_unplug property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.DiscVirtioHotPlug != nil {
		err := d.Set("disc_virtio_hot_plug", *volume.Properties.DiscVirtioHotPlug)
		if err != nil {
			return fmt.Errorf("error while setting disc_virtio_hot_plug property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.DiscVirtioHotUnplug != nil {
		err := d.Set("disc_virtio_hot_unplug", *volume.Properties.DiscVirtioHotUnplug)
		if err != nil {
			return fmt.Errorf("error while setting disc_virtio_hot_unplug property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.BackupunitId != nil {
		err := d.Set("backup_unit_id", *volume.Properties.BackupunitId)
		if err != nil {
			return fmt.Errorf("error while setting backup_unit_id property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.UserData != nil {
		err := d.Set("user_data", *volume.Properties.UserData)
		if err != nil {
			return fmt.Errorf("error while setting user_data property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.DeviceNumber != nil {
		err := d.Set("device_number", *volume.Properties.DeviceNumber)
		if err != nil {
			return fmt.Errorf("error while setting device_number property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.BootServer != nil {
		err := d.Set("boot_server", *volume.Properties.BootServer)
		if err != nil {
			return fmt.Errorf("error while setting boot_server property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.ExposeSerial != nil {
		err := d.Set("expose_serial", *volume.Properties.ExposeSerial)
		if err != nil {
			return fmt.Errorf("error while setting exposeSerial property for volume %s: %w", d.Id(), err)
		}
	}

	if volume.Properties.RequireLegacyBios != nil {
		if err := d.Set("require_legacy_bios", *volume.Properties.RequireLegacyBios); err != nil {
			return fmt.Errorf("error while setting require_legacy_bios property for volume with ID: %v, error: %w", d.Id(), err)
		}
	}

	return nil
}

func getVolumeData(ctx context.Context, d *schema.ResourceData, path, serverType string) (*ionoscloud.VolumeProperties, error) {
	volume := ionoscloud.VolumeProperties{}

	if !strings.EqualFold(serverType, constant.GpuType) {
		// For GPU servers, disk_type is set to "SSD Premium" by default and cannot (yet) be changed
		volumeType := d.Get(path + "disk_type").(string)
		volume.Type = &volumeType
	}

	if v, ok := d.GetOk(path + "availability_zone"); ok {
		vStr := v.(string)
		volume.AvailabilityZone = &vStr
	}

	if v, ok := d.GetOk(path + "image_password"); ok {
		vStr := v.(string)
		volume.ImagePassword = &vStr
		if err := d.Set("image_password", vStr); err != nil {
			return nil, err
		}
	}

	if v, ok := d.GetOk("image_password"); ok {
		vStr := v.(string)
		volume.ImagePassword = &vStr
	}

	if v, ok := d.GetOk(path + "licence_type"); ok {
		vStr := v.(string)
		volume.LicenceType = &vStr
	}

	if v, ok := d.GetOk(path + "bus"); ok {
		vStr := v.(string)
		volume.Bus = &vStr
	}

	var volumeSize float32
	if !strings.EqualFold(serverType, constant.CubeType) &&
		!strings.EqualFold(serverType, constant.GpuType) {
		volumeSize = float32(d.Get(path + "size").(int))
		if volumeSize > 0 {
			volume.Size = &volumeSize
		}
	}

	if v, ok := d.GetOk(path + "name"); ok {
		vStr := v.(string)
		volume.Name = &vStr
	}

	var sshKeys []any

	if serverType != constant.VCPUType {
		if v, ok := d.GetOk(path + "ssh_key_path"); ok {
			sshKeys = v.([]any)
		} else if v, ok := d.GetOk("ssh_key_path"); ok {
			sshKeys = v.([]any)
		} else {
			if err := d.Set("ssh_key_path", [][]string{}); err != nil {
				return nil, err
			}
		}
	}

	if v, ok := d.GetOk(path + "ssh_keys"); ok {
		sshKeys = v.([]any)
	} else if v, ok := d.GetOk("ssh_keys"); ok {
		sshKeys = v.([]any)
	}

	if len(sshKeys) != 0 {
		var publicKeys []string
		for _, path := range sshKeys {
			if path == nil {
				return nil, fmt.Errorf("ssh_keys or ssh_key_path contains empty value")
			}

			tflog.Debug(ctx, "reading ssh key file", map[string]any{"path": path})
			publicKey, err := utils.ReadPublicKey(ctx, path.(string))
			if err != nil {
				return nil, err
			}
			publicKeys = append(publicKeys, publicKey)
		}
		if len(publicKeys) > 0 {
			volume.SshKeys = &publicKeys
		}
	}

	if v, ok := d.GetOk(path + "expose_serial"); ok {
		val := v.(bool)
		volume.ExposeSerial = &val
	}

	if v, ok := d.GetOk(path + "require_legacy_bios"); ok {
		requireLegacyBios := v.(bool)
		volume.RequireLegacyBios = &requireLegacyBios
	}

	return &volume, nil
}

// getImage is used for the entire logic for finding the image/snapshot provided by the user
func getImage(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData, volume ionoscloud.VolumeProperties) (image, imageAlias string, err error) {
	var imageName string
	dcID := d.Get("datacenter_id").(string)
	isSnapshot := false

	if v, ok := d.GetOk("volume.0.image_name"); ok {
		imageName = v.(string)
		if err := d.Set("image_name", v.(string)); err != nil {
			return image, imageAlias, err
		}
	} else if v, ok := d.GetOk("image_name"); ok {
		imageName = v.(string)
	}

	if imageName != "" {
		if !utils.IsValidUUID(imageName) {
			images, err := cloudapiimage.GetAllImages(ctx, client)
			if err != nil {
				return image, imageAlias, fmt.Errorf("error while fetching the list of images: %w", err)
			}

			dc, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, dcID).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				return image, imageAlias, fmt.Errorf("error fetching datacenter %s: (%w)", dcID, err)
			}

			locationIDs := cloudapilocation.ResolveParentLocation(ctx, client, *dc.Properties.Location)

			img, rejectedImg := resolveVolumeImageName(imageName, images, locationIDs)
			if img != nil {
				image = *img.Id
			}
			// if no image id was found with that name we look for a matching snapshot
			if image == "" {
				tflog.Debug(ctx, "looking for a snapshot by name", map[string]any{"image_name": imageName})
				image = getSnapshotID(ctx, client, imageName)
				if image != "" {
					isSnapshot = true
				} else {
					imageAlias = cloudapiimage.GetImageAlias(imageName, images, locationIDs)
					if imageAlias == "" {
						if rejectedImg != nil {
							return image, imageAlias, fmt.Errorf(
								"image '%s' was found (name: '%s') with type '%s' in location '%s'; "+
									"volume requires an image of type '%s' in location '%s'",
								imageName, *rejectedImg.Properties.Name, *rejectedImg.Properties.ImageType,
								*rejectedImg.Properties.Location, HDDImage, *dc.Properties.Location)
						}
						return image, imageAlias, fmt.Errorf("could not find an image/imagealias/snapshot that matches %s", imageName)
					}
				}
			}

			if volume.ImagePassword == nil && (volume.SshKeys == nil || len(*volume.SshKeys) == 0) && isSnapshot == false &&
				(img == nil || (img.Properties.Public != nil && *img.Properties.Public)) {
				return image, imageAlias, fmt.Errorf("volume, either 'image_password' or 'ssh_key_path'/'ssh_keys' must be provided")
			}
		} else {
			img, apiResponse, err := client.ImagesApi.ImagesFindById(ctx, imageName).Execute()
			logApiRequestTime(apiResponse)
			// here we search for snapshot if we do not find img based on imageName
			if apiResponse.SafeStatusCode() == 404 {

				snapshot, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, imageName).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return image, imageAlias, fmt.Errorf("could not fetch image/snapshot: %w", err)
				}

				isSnapshot = true
				if snapshot.Id != nil {
					image = *snapshot.Id
				}
			} else if err != nil {
				return image, imageAlias, fmt.Errorf("error fetching image/snapshot: %w", err)
			}

			if isSnapshot == false && img.Properties.Public != nil && *img.Properties.Public == true {

				if volume.ImagePassword == nil && (volume.SshKeys == nil || len(*volume.SshKeys) == 0) {
					return image, imageAlias, fmt.Errorf("public image, either 'image_password' or 'ssh_key_path'/'ssh_keys' must be provided")
				}

				dc, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, dcID).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return image, imageAlias, fmt.Errorf("error fetching datacenter %s: (%w)", dcID, err)
				}

				locationIDs := cloudapilocation.ResolveParentLocation(ctx, client, *dc.Properties.Location)

				images, err := cloudapiimage.GetAllImages(ctx, client)
				if err != nil {
					return image, imageAlias, fmt.Errorf("error while fetching the list of images: %w", err)
				}

				img, rejectedImg := resolveVolumeImageName(imageName, images, locationIDs)
				if rejectedImg != nil {
					tflog.Debug(ctx, "image matched by name but filtered out", map[string]any{"name": *rejectedImg.Properties.Name, "type": *rejectedImg.Properties.ImageType, "location": *rejectedImg.Properties.Location})
				}

				if img != nil {
					image = *img.Id
				}
			} else {
				img, apiResponse, err := client.ImagesApi.ImagesFindById(ctx, imageName).Execute()

				logApiRequestTime(apiResponse)
				if err != nil {
					// we want to search for snapshot again, but we check for
					// image != "" to be sure we didn't find it when we searched above for it
					if apiResponse.SafeStatusCode() == 404 && image != "" {
						snapshot, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, imageName).Execute()
						logApiRequestTime(apiResponse)
						if err != nil {
							return image, imageAlias, fmt.Errorf("error fetching image/snapshot: %w", err)
						}
						if snapshot.Id != nil {
							image = *snapshot.Id
						}
						isSnapshot = true
					} else {
						return image, imageAlias, err
					}

				} else {
					if isSnapshot == false && img.Properties.Public != nil && *img.Properties.Public == true {
						if volume.ImagePassword == nil && (volume.SshKeys == nil || len(*volume.SshKeys) == 0) {
							return image, imageAlias, fmt.Errorf("either 'image_password' or 'ssh_key_path'/'ssh_keys' must be provided for imageName %s ", imageName)
						}
						image = imageName
					} else {
						image = imageName
					}
				}
			}
		}
	}

	if image == "" && volume.LicenceType == nil && imageAlias == "" && !isSnapshot {
		return image, imageAlias, fmt.Errorf("either 'image_name', 'licence_type', or 'image_alias' must be set")
	}

	if isSnapshot == true && (volume.ImagePassword != nil || volume.SshKeys != nil && len(*volume.SshKeys) > 0) {
		return image, imageAlias, fmt.Errorf("passwords/SSH keys are not supported for snapshots")
	}

	return image, imageAlias, nil
}

func getSnapshotID(ctx context.Context, client *ionoscloud.APIClient, snapshotName string) string {

	if snapshotName == "" {
		return ""
	}

	snapshots, apiResponse, err := client.SnapshotsApi.SnapshotsGet(ctx).Depth(1).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		tflog.Error(ctx, "error while fetching the list of snapshots", map[string]any{"error": err.Error()})
	}

	if len(*snapshots.Items) > 0 {
		for _, i := range *snapshots.Items {
			imgName := ""
			if i.Properties != nil && i.Properties.Name != nil && *i.Properties.Name != "" {
				imgName = *i.Properties.Name
			}

			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(snapshotName)) {
				return *i.Id
			}
		}
	}
	return ""
}

// resolveVolumeImageName scans the given images for imageName: match is the HDD image in
// one of the given locations matched by id or name (exact wins over partial), skipped is
// a name match filtered out for wrong type/location.
func resolveVolumeImageName(imageName string, images []ionoscloud.Image, locations []string) (match, skipped *ionoscloud.Image) {

	if imageName == "" {
		return nil, nil
	}

	var partialMatch *ionoscloud.Image
	var nameMatchWrongTypeOrLocation *ionoscloud.Image
	for _, imageEntry := range images {
		if imageEntry.Properties != nil && imageEntry.Properties.Name != nil && *imageEntry.Properties.Name != "" {

			nameMatches := (imageEntry.Id != nil && strings.EqualFold(imageName, *imageEntry.Id)) ||
				strings.EqualFold(*imageEntry.Properties.Name, imageName) ||
				strings.Contains(strings.ToLower(*imageEntry.Properties.Name), strings.ToLower(imageName))

			if *imageEntry.Properties.ImageType != HDDImage || !cloudapilocation.LocationInSet(locations, *imageEntry.Properties.Location) {
				if nameMatchWrongTypeOrLocation == nil && nameMatches {
					nameMatchWrongTypeOrLocation = &imageEntry
				}
				continue
			}
			// Return the image entry if the name is an exact match
			if strings.EqualFold(imageName, *imageEntry.Id) || strings.EqualFold(*imageEntry.Properties.Name, imageName) {
				return &imageEntry, nil
			}
			// Save the first image entry which is a partial match and return it if no exact matches were found
			if partialMatch == nil && strings.Contains(strings.ToLower(*imageEntry.Properties.Name), strings.ToLower(imageName)) {
				partialMatch = &imageEntry
			}
		}
	}
	return partialMatch, nameMatchWrongTypeOrLocation
}
