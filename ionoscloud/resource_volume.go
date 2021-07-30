package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeCreate,
		ReadContext:   resourceVolumeRead,
		UpdateContext: resourceVolumeUpdate,
		DeleteContext: resourceVolumeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVolumeImporter,
		},
		Schema: map[string]*schema.Schema{
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"disk_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"image_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"licence_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssh_key_path": {
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"datacenter_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"cpu_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ram_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"backup_unit_id": {
				Type:     schema.TypeString,
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	var imageAlias string
	var sshKeyPath []interface{}
	isSnapshot := false
	dcId := d.Get("datacenter_id").(string)
	serverId := d.Get("server_id").(string)
	imagePassword := d.Get("image_password").(string)
	sshKeyPath = d.Get("ssh_key_path").([]interface{})
	imageInput := d.Get("image_name").(string)
	licenceType := d.Get("licence_type").(string)
	diskType := d.Get("disk_type").(string)

	var publicKeys []string
	if len(sshKeyPath) != 0 {
		for _, path := range sshKeyPath {
			log.Printf("[DEBUG] Reading file %s", path)
			publicKey, err := readPublicKey(path.(string))
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error fetching sshkey from file (%s) (%s)", path, err.Error()))
				return diags
			}
			publicKeys = append(publicKeys, publicKey)
		}
	}

	var image string
	if imageInput != "" {
		if !IsValidUUID(imageInput) {
			img, err := getImage(client, dcId, imageInput, diskType)
			if err != nil {
				diags := diag.FromErr(err)
				return diags
			}
			if img != nil {
				image = *img.Id
			}
			// if no image id was found with that name we look for a matching snapshot
			if image == "" {
				image = getSnapshotId(client, imageInput)
				if image != "" {
					isSnapshot = true
				} else {
					diags := diag.FromErr(fmt.Errorf("no image or snapshot with id %s found", imageInput))
					return diags
				}
			}

			if imagePassword == "" && len(sshKeyPath) == 0 && isSnapshot == false && img.Properties.Public != nil && *img.Properties.Public {
				diags := diag.FromErr(fmt.Errorf("either 'image_password' or 'ssh_key_path' must be provided"))
				return diags
			}

		} else {

			img, apiResponse, err := client.ImagesApi.ImagesFindById(ctx, imageInput).Execute()

			if err != nil {

				if apiResponse != nil && apiResponse.StatusCode == 404 {
					snapshot, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, imageInput).Execute()

					if err != nil {
						if apiResponse != nil && apiResponse.StatusCode == 404 {
							diags := diag.FromErr(fmt.Errorf("image/snapshot %s not found: %s", imageInput, err))
							return diags
						} else {
							diags := diag.FromErr(fmt.Errorf("an error occured while fetching snapshot %s: %s", imageInput, err))
							return diags
						}
					}
					image = *snapshot.Id
					isSnapshot = true
				} else {
					diags := diag.FromErr(fmt.Errorf("error fetching image %s: %s", imageInput, err))
					return diags
				}
			} else {
				if isSnapshot == false && *img.Properties.Public == true {
					if imagePassword == "" && len(sshKeyPath) == 0 {
						diags := diag.FromErr(fmt.Errorf("either 'image_password' or 'sshkey' must be provided"))
						return diags
					}
				}
				image = *img.Id
			}
		}
	}

	if isSnapshot == true && (imagePassword != "" || len(publicKeys) > 0) {
		diags := diag.FromErr(fmt.Errorf("you can't pass 'image_password' and/or 'ssh keys' when creating a volume from a snapshot"))
		return diags
	}

	volumeName := d.Get("name").(string)
	volumeSize := float32(d.Get("size").(int))
	volumeType := d.Get("disk_type").(string)
	volumeBus := d.Get("bus").(string)

	volume := ionoscloud.Volume{
		Properties: &ionoscloud.VolumeProperties{
			Name:        &volumeName,
			Size:        &volumeSize,
			Type:        &volumeType,
			LicenceType: &licenceType,
		},
	}

	if imagePassword != "" {
		volume.Properties.ImagePassword = &imagePassword
	} else {
		volume.Properties.ImagePassword = nil
	}

	if volumeBus != "" {
		volume.Properties.Bus = &volumeBus
	} else {
		volume.Properties.Bus = nil
	}

	if licenceType != "" {
		volume.Properties.LicenceType = &licenceType
	} else {
		volume.Properties.LicenceType = nil
	}

	if v, ok := d.GetOk("cpu_hot_plug"); ok {
		vBool := v.(bool)
		volume.Properties.CpuHotPlug = &vBool
	}

	if v, ok := d.GetOk("ram_hot_plug"); ok {
		vBool := v.(bool)
		volume.Properties.RamHotPlug = &vBool
	}

	if v, ok := d.GetOk("nic_hot_unplug"); ok {
		vBool := v.(bool)
		volume.Properties.NicHotPlug = &vBool
	}

	if v, ok := d.GetOk("nic_hot_unplug"); ok {
		vBool := v.(bool)
		volume.Properties.NicHotUnplug = &vBool
	}

	if v, ok := d.GetOk("disc_virtio_hot_plug"); ok {
		vBool := v.(bool)
		volume.Properties.DiscVirtioHotPlug = &vBool
	}

	if v, ok := d.GetOk("disc_virtio_hot_unplug"); ok {
		vBool := v.(bool)
		volume.Properties.DiscVirtioHotUnplug = &vBool
	}

	if image != "" {
		volume.Properties.Image = &image
	} else {
		volume.Properties.Image = nil
	}

	if len(publicKeys) != 0 {
		volume.Properties.SshKeys = &publicKeys

	} else {
		volume.Properties.SshKeys = nil
	}

	if _, ok := d.GetOk("availability_zone"); ok {
		raw := d.Get("availability_zone").(string)
		volume.Properties.AvailabilityZone = &raw
	}

	userData := d.Get("user_data").(string)
	if userData != "" {
		if image == "" && imageAlias == "" {
			diags := diag.FromErr(fmt.Errorf("it is mandatory to provied either public image that has cloud-init compatibility in conjunction with backup unit id property "))
			return diags
		} else {
			volume.Properties.UserData = &userData
		}
	} else {
		volume.Properties.UserData = nil
	}

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesPost(ctx, dcId).Volume(volume).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a volume: %s", err))
		return diags
	}

	d.SetId(*volume.Id)

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

	volumeToAttach := ionoscloud.Volume{Id: volume.Id}
	volume, apiResponse, err = client.ServersApi.DatacentersServersVolumesPost(ctx, dcId, serverId).Volume(volumeToAttach).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s", dcId, serverId, *volumeToAttach.Id, err))
		return diags
	}

	sErr := d.Set("server_id", serverId)

	if sErr != nil {
		diags := diag.FromErr(fmt.Errorf("error while setting serverId %s: %s", serverId, sErr))
		return diags
	}

	// Wait, catching any errors
	_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			sErr := d.Set("server_id", "")
			if sErr != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting serverId: %s", sErr))
				return diags
			}
		}
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceVolumeRead(ctx, d, meta)
}

func resourceVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)
	volumeID := d.Id()

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, dcId, volumeID).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching volume with ID %s: %s", d.Id(), err))
		return diags
	}

	_, _, err = client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, serverID, volumeID).Execute()
	if err != nil {
		if err2 := d.Set("server_id", ""); err2 != nil {
			diags := diag.FromErr(err2)
			return diags
		}
	}

	if volume.Properties.Name != nil {
		err := d.Set("name", *volume.Properties.Name)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting name property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	if volume.Properties.Type != nil {
		err := d.Set("disk_type", *volume.Properties.Type)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting type property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	if volume.Properties.Size != nil {
		err := d.Set("size", *volume.Properties.Size)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting size property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	if volume.Properties.Bus != nil {
		err := d.Set("bus", *volume.Properties.Bus)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting bus property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	/*
		if volume.Properties.Image != nil {
			err = d.Set("image_name", *volume.Properties.Image)
			if err != nil {
				return fmt.Errorf("error while setting image property for volume %s: %s", d.Id(), err)
			}
		}
	*/

	if volume.Properties.CpuHotPlug != nil {
		err := d.Set("cpu_hot_plug", *volume.Properties.CpuHotPlug)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting cpu_hot_plug property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	if volume.Properties.RamHotPlug != nil {
		err := d.Set("ram_hot_plug", *volume.Properties.RamHotPlug)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting ram_hot_plug property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	if volume.Properties.NicHotPlug != nil {
		err := d.Set("nic_hot_plug", *volume.Properties.NicHotPlug)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting nic_hot_plug property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	if volume.Properties.NicHotUnplug != nil {
		err := d.Set("nic_hot_unplug", *volume.Properties.NicHotUnplug)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting nic_hot_unplug property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	if volume.Properties.DiscVirtioHotPlug != nil {
		err := d.Set("disc_virtio_hot_plug", *volume.Properties.DiscVirtioHotPlug)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting disc_virtio_hot_plug property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	if volume.Properties.DiscVirtioHotUnplug != nil {
		err := d.Set("disc_virtio_hot_unplug", *volume.Properties.DiscVirtioHotUnplug)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting disc_virtio_hot_unplug property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	if volume.Properties.BackupunitId != nil {
		err := d.Set("backup_unit_id", *volume.Properties.BackupunitId)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting backup_unit_id property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	if volume.Properties.UserData != nil {
		err := d.Set("user_data", *volume.Properties.UserData)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting user_data property for volume %s: %s", d.Id(), err))
			return diags
		}
	}

	return nil
}

func resourceVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	properties := ionoscloud.VolumeProperties{}
	dcId := d.Get("datacenter_id").(string)

	if d.HasChange("name") {
		_, newValue := d.GetChange("name")
		newValueStr := newValue.(string)
		properties.Name = &newValueStr
	}
	if d.HasChange("disk_type") {
		_, newValue := d.GetChange("disk_type")
		newValueStr := newValue.(string)
		properties.Type = &newValueStr
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
	if d.HasChange("availability_zone") {
		_, newValue := d.GetChange("availability_zone")
		newValueStr := newValue.(string)
		properties.AvailabilityZone = &newValueStr
	}
	if d.HasChange("cpu_hot_plug") {
		_, newValue := d.GetChange("cpu_hot_plug")
		newValueBool := newValue.(bool)
		properties.CpuHotPlug = &newValueBool
	}
	if d.HasChange("ram_hot_plug") {
		_, newValue := d.GetChange("ram_hot_plug")
		newValueBool := newValue.(bool)
		properties.RamHotPlug = &newValueBool
	}
	if d.HasChange("nic_hot_plug") {
		_, newValue := d.GetChange("nic_hot_plug")
		newValueBool := newValue.(bool)
		properties.NicHotPlug = &newValueBool
	}

	if d.HasChange("nic_hot_unplug") {
		_, newValue := d.GetChange("nic_hot_unplug")
		newValueBool := newValue.(bool)
		properties.NicHotUnplug = &newValueBool
	}

	if d.HasChange("disc_virtio_hot_plug") {
		_, newValue := d.GetChange("disc_virtio_hot_plug")
		newValueBool := newValue.(bool)
		properties.DiscVirtioHotPlug = &newValueBool
	}

	if d.HasChange("disc_virtio_hot_unplug") {
		_, newValue := d.GetChange("disc_virtio_hot_unplug")
		newValueBool := newValue.(bool)
		properties.DiscVirtioHotUnplug = &newValueBool
	}

	if d.HasChange("user_data") {
		diags := diag.FromErr(fmt.Errorf("User data property of resource volume is immutable "))
		return diags
	}

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesPatch(ctx, dcId, d.Id()).Volume(properties).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating volume with ID %s: %s", d.Id(), err))
		return diags

	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	if apiResponse != nil && apiResponse.StatusCode > 299 {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a volume ID %s %s", d.Id(), err))
		return diags
	}

	if d.HasChange("server_id") {
		_, newValue := d.GetChange("server_id")
		serverID := newValue.(string)
		_, apiResponse, err := client.ServersApi.DatacentersServersVolumesPost(ctx, dcId, serverID).Volume(volume).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s",
				dcId, serverID, *volume.Id, err))
			return diags
		}

		// Wait, catching any errors
		_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(errState)
			return diags
		}
	}

	return resourceVolumeRead(ctx, d, meta)
}

func resourceVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)

	apiResponse, err := client.VolumesApi.DatacentersVolumesDelete(ctx, dcId, d.Id()).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a volume ID %s %s", d.Id(), err))
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
