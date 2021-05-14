package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceVolumeCreate,
		Read:   resourceVolumeRead,
		Update: resourceVolumeUpdate,
		Delete: resourceVolumeDelete,
		Schema: map[string]*schema.Schema{
			"image": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Required: true,
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
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Optional: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(SdkBundle).Client

	var ssh_keypath []interface{}
	isSnapshot := false
	dcId := d.Get("datacenter_id").(string)
	serverId := d.Get("server_id").(string)
	imagePassword := d.Get("image_password").(string)
	ssh_keypath = d.Get("ssh_key_path").([]interface{})
	image := d.Get("image").(string)
	licenceType := d.Get("licence_type").(string)

	var publicKeys []string
	if len(ssh_keypath) != 0 {
		for _, path := range ssh_keypath {
			log.Printf("[DEBUG] Reading file %s", path)
			publicKey, err := readPublicKey(path.(string))
			if err != nil {
				return fmt.Errorf("Error fetching sshkey from file (%s) (%s)", path, err.Error())
			}
			publicKeys = append(publicKeys, publicKey)
		}
	}

	if image != "" {
		if !IsValidUUID(image) {
			return fmt.Errorf("Image is not a valid UUID")
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

			if cancel != nil {
				defer cancel()
			}

			img, apiResponse, err := client.ImagesApi.ImagesFindById(ctx, image).Execute()

			_, rsp := err.(ionoscloud.GenericOpenAPIError)

			if err != nil {
				return fmt.Errorf("Error fetching image %s: (%s) - %+v", image, err, rsp)
			}

			if apiResponse.Response.StatusCode == 404 {
				_, _, err := client.SnapshotsApi.SnapshotsFindById(ctx, image).Execute()

				if _, ok := err.(ionoscloud.GenericOpenAPIError); !ok {
					if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
						return fmt.Errorf("image/snapshot: %s Not Found", string(apiResponse.Payload))
					}
				}

				isSnapshot = true
			}
			if *img.Properties.Public == true && isSnapshot == false {
				if imagePassword == "" && len(ssh_keypath) == 0 {
					return fmt.Errorf("Either 'image_password' or 'sshkey' must be provided.")
				}
			}
		}
	}

	if image == "" && licenceType == "" && isSnapshot == false {
		return fmt.Errorf("Either 'image', or 'licenceType' must be set.")
	}

	if isSnapshot == true && (imagePassword != "" || len(publicKeys) > 0) {
		return fmt.Errorf("You can't pass 'image_password' and/or 'ssh keys' when creating a volume from a snapshot")
	}

	volumeName := d.Get("name").(string)
	volumeSize := float32(d.Get("size").(int))
	volumeType := d.Get("disk_type").(string)
	volumeBus := d.Get("bus").(string)

	volume := ionoscloud.Volume{
		Properties: &ionoscloud.VolumeProperties{
			Name:          &volumeName,
			Size:          &volumeSize,
			Type:          &volumeType,
			ImagePassword: &imagePassword,
			Bus:           &volumeBus,
		},
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

	backupUnitId := d.Get("backup_unit_id").(string)
	if IsValidUUID(backupUnitId) {
		if image == "" {
			return fmt.Errorf("It is mandatory to provied public image in conjunction with backup unit id property")
		} else {
			volume.Properties.BackupunitId = &backupUnitId
		}
	} else {
		volume.Properties.BackupunitId = nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}
	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesPost(ctx, dcId).Volume(volume).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while creating a volume: %s", err)
	}

	d.SetId(*volume.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	volumeToAttach := ionoscloud.Volume{Id: volume.Id}
	volume, apiResponse, err = client.ServersApi.DatacentersServersVolumesPost(ctx, dcId, serverId).Volume(volumeToAttach).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s", dcId, serverId, *volume.Id, err)
	}

	sErr := d.Set("server_id", serverId)

	if sErr != nil {
		return fmt.Errorf("Error while setting serverId %s: %s", serverId, sErr)
	}

	// Wait, catching any errors
	_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			sErr := d.Set("server_id", "")
			if sErr != nil {
				return fmt.Errorf("Error while setting serverId: %s", sErr)
			}
		}
		return errState
	}

	return resourceVolumeRead(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(SdkBundle).Client
	dcId := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)
	volumeID := d.Id()

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, dcId, volumeID).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error occured while fetching a volume ID %s %s", d.Id(), err)
	}

	if apiResponse.Response.StatusCode > 299 {
		return fmt.Errorf("An error occured while fetching a volume ID %s", d.Id())

	}

	_, _, err = client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, serverID, volumeID).Execute()
	if err != nil {
		d.Set("server_id", "")
	}

	if volume.Properties.Name != nil {
		err := d.Set("name", *volume.Properties.Name)
		if err != nil {
			return fmt.Errorf("Error while setting name property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.Type != nil {
		err := d.Set("disk_type", *volume.Properties.Type)
		if err != nil {
			return fmt.Errorf("Error while setting type property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.Size != nil {
		err := d.Set("size", *volume.Properties.Size)
		if err != nil {
			return fmt.Errorf("Error while setting size property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.Bus != nil {
		err := d.Set("bus", *volume.Properties.Bus)
		if err != nil {
			return fmt.Errorf("Error while setting bus property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.Image != nil {
		err = d.Set("image", *volume.Properties.Image)
		if err != nil {
			return fmt.Errorf("Error while setting image property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.CpuHotPlug != nil {
		err := d.Set("cpu_hot_plug", *volume.Properties.CpuHotPlug)
		if err != nil {
			return fmt.Errorf("Error while setting cpu_hot_plug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.RamHotPlug != nil {
		err := d.Set("ram_hot_plug", *volume.Properties.RamHotPlug)
		if err != nil {
			return fmt.Errorf("Error while setting ram_hot_plug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.NicHotPlug != nil {
		err := d.Set("nic_hot_plug", *volume.Properties.NicHotPlug)
		if err != nil {
			return fmt.Errorf("Error while setting nic_hot_plug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.NicHotUnplug != nil {
		err := d.Set("nic_hot_unplug", *volume.Properties.NicHotUnplug)
		if err != nil {
			return fmt.Errorf("Error while setting nic_hot_unplug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.DiscVirtioHotPlug != nil {
		err := d.Set("disc_virtio_hot_plug", *volume.Properties.DiscVirtioHotPlug)
		if err != nil {
			return fmt.Errorf("Error while setting disc_virtio_hot_plug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.DiscVirtioHotUnplug != nil {
		err := d.Set("disc_virtio_hot_unplug", *volume.Properties.DiscVirtioHotUnplug)
		if err != nil {
			return fmt.Errorf("Error while setting disc_virtio_hot_unplug property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.BackupunitId != nil {
		err := d.Set("backup_unit_id", *volume.Properties.BackupunitId)
		if err != nil {
			return fmt.Errorf("Error while setting backup_unit_id property for volume %s: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client
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

	if d.HasChange("backup_unit_id") {
		return fmt.Errorf("Backup unit id property is immutable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

	if cancel != nil {
		defer cancel()
	}

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesPatch(ctx, dcId, d.Id()).Volume(properties).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while updating a volume ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	if apiResponse.Response.StatusCode > 299 {
		return fmt.Errorf("An error occured while updating a volume ID %s", d.Id())

	}

	if d.HasChange("server_id") {
		_, newValue := d.GetChange("server_id")
		serverID := newValue.(string)
		volumeAttach, apiResponse, err := client.ServersApi.DatacentersServersVolumesPost(ctx, dcId, serverID).Volume(volume).Execute()
		if err != nil {
			return fmt.Errorf("An error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s", dcId, serverID, *volumeAttach.Id, err)
		}

		// Wait, catching any errors
		_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
		if errState != nil {
			return errState
		}
	}

	return resourceVolumeRead(d, meta)
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client
	dcId := d.Get("datacenter_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.VolumesApi.DatacentersVolumesDelete(ctx, dcId, d.Id()).Execute()
	if err != nil {
		return fmt.Errorf("An error occured while deleting a volume ID %s %s", d.Id(), err)

	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
