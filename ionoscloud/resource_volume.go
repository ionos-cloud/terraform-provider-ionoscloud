package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Optional: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ionoscloud.APIClient)

	var ssh_keypath []interface{}
	var image_alias string
	isSnapshot := false
	dcId := d.Get("datacenter_id").(string)
	serverId := d.Get("server_id").(string)
	imagePassword := d.Get("image_password").(string)
	ssh_keypath = d.Get("ssh_key_path").([]interface{})
	image_name := d.Get("image_name").(string)

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

	var image string
	if image_alias == "" && image_name != "" {
		if !IsValidUUID(image_name) {
			img, err := getImage(client, dcId, image_name, d.Get("disk_type").(string))
			if err != nil {
				return err
			}
			if img != nil {
				image = *img.Id
			}
			//if no image id was found with that name we look for a matching snapshot
			if image == "" {
				image = getSnapshotId(client, image_name)
				if image != "" {
					isSnapshot = true
				} else {
					ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

					if cancel != nil {
						defer cancel()
					}

					dc, _, err := client.DataCenterApi.DatacentersFindById(ctx, dcId).Execute()

					if err != nil {
						return fmt.Errorf("An error occured while fetching a Datacenter ID %s %s", dcId, err)
					}
					image_alias = getImageAlias(client, image_name, *dc.Properties.Location)
				}
			}

			if image == "" && image_alias == "" {
				return fmt.Errorf("Could not find an image/imagealias/snapshot that matches %s ", image_name)
			}
			if imagePassword == "" && len(ssh_keypath) == 0 && isSnapshot == false && img != nil && *img.Properties.Public {
				return fmt.Errorf("Either 'image_password' or 'sshkey' must be provided.")
			}
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

			if cancel != nil {
				defer cancel()
			}

			img, _, err := client.ImageApi.ImagesFindById(ctx, image_name).Execute()
			if err != nil {
				_, _, err := client.SnapshotApi.SnapshotsFindById(ctx, image_name).Execute()
				if err != nil {
					return fmt.Errorf("Error fetching image/snapshot: %s", err)
				}
				isSnapshot = true
			}
			if *img.Properties.Public == true && isSnapshot == false {
				if imagePassword == "" && len(ssh_keypath) == 0 {
					return fmt.Errorf("Either 'image_password' or 'sshkey' must be provided.")
				}
				image = image_name
			} else {
				image = image_name
			}
		}
	}

	if image_name == "" && licenceType == "" && isSnapshot == false {
		return fmt.Errorf("either 'image_name', or 'licenceType' must be set")
	}

	if isSnapshot == true && (imagePassword != "" || len(publicKeys) > 0) {
		return fmt.Errorf("you can't pass 'image_password' and/or 'ssh keys' when creating a volume from a snapshot")
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
			Bus:           &volumeBus,
			LicenceType:   &licenceType,
		},
	}

	if imagePassword != "" {
		volume.Properties.ImagePassword = &imagePassword
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

	if image_alias != "" {
		volume.Properties.ImageAlias = &image_alias
	} else {
		volume.Properties.ImageAlias = nil
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
		if image == "" && image_alias == "" {
			return fmt.Errorf("it is mandatory to provied either public image or imageAlias in conjunction with backup unit id property")
		} else {
			volume.Properties.BackupunitId = &backupUnitId
		}
	} else {
		volume.Properties.BackupunitId = nil
	}

	userData := d.Get("user_data").(string)
	if userData != "" {
		if image == "" && image_alias == "" {
			return fmt.Errorf("it is mandatory to provied either public image or imageAlias that has cloud-init compatibility in conjunction with backup unit id property ")
		} else {
			volume.Properties.UserData = &userData
		}
	} else {
		volume.Properties.UserData = nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}
	volume, apiResponse, err := client.VolumeApi.DatacentersVolumesPost(ctx, dcId).Volume(volume).Execute()

	if err != nil {
		return fmt.Errorf("an error occured while creating a volume: %s", err)
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
	volume, apiResponse, err = client.ServerApi.DatacentersServersVolumesPost(ctx, dcId, serverId).Volume(volumeToAttach).Execute()

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
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)
	volumeID := d.Id()

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}

	volume, apiResponse, err := client.VolumeApi.DatacentersVolumesFindById(ctx, dcId, volumeID).Execute()

	if err != nil {

		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("error occured while fetching volume with ID %s: %s", d.Id(), err)
	}

	_, _, err = client.ServerApi.DatacentersServersVolumesFindById(ctx, dcId, serverID, volumeID).Execute()
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
		image, _, err := client.ImageApi.ImagesFindById(ctx, *volume.Properties.Image).Execute()
		if err != nil {
			return fmt.Errorf("Error while getting image_name property for image %s: %s", *volume.Properties.Image, err)
		}
		err = d.Set("image_name", *image.Properties.Name)
		if err != nil {
			return fmt.Errorf("Error while setting image_name property for volume %s: %s", d.Id(), err)
		}
	}

	if volume.Properties.ImageAlias != nil {
		err := d.Set("image_alias", *volume.Properties.ImageAlias)
		if err != nil {
			return fmt.Errorf("Error while setting image_alias property for volume %s: %s", d.Id(), err)
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

	if volume.Properties.UserData != nil {
		err := d.Set("user_data", *volume.Properties.UserData)
		if err != nil {
			return fmt.Errorf("Error while setting user_data property for volume %s: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

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

	if d.HasChange("user_data") {
		return fmt.Errorf("User data property of resource volume is immutable ")
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

	if cancel != nil {
		defer cancel()
	}

	volume, apiResponse, err := client.VolumeApi.DatacentersVolumesPatch(ctx, dcId, d.Id()).Volume(properties).Execute()

	if err != nil {
		return fmt.Errorf("an error occured while updating volume with ID %s: %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	if d.HasChange("server_id") {
		_, newValue := d.GetChange("server_id")
		serverID := newValue.(string)
		volumeAttach, apiResponse, err := client.ServerApi.DatacentersServersVolumesPost(ctx, dcId, serverID).Volume(volume).Execute()
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
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.VolumeApi.DatacentersVolumesDelete(ctx, dcId, d.Id()).Execute()
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
