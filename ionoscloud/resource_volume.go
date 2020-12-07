package ionoscloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

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
				image = img.ID
			}
			//if no image id was found with that name we look for a matching snapshot
			if image == "" {
				image = getSnapshotId(client, image_name)
				if image != "" {
					isSnapshot = true
				} else {
					dc, err := client.GetDatacenter(dcId)

					if err != nil {
						return fmt.Errorf("An error occured while fetching a Datacenter ID %s %s", dcId, err)
					}

					image_alias = getImageAlias(client, image_name, dc.Properties.Location)
				}
			}

			if image == "" && image_alias == "" {
				return fmt.Errorf("Could not find an image/imagealias/snapshot that matches %s ", image_name)
			}
			if imagePassword == "" && len(ssh_keypath) == 0 && isSnapshot == false && img != nil && img.Properties.Public {
				return fmt.Errorf("Either 'image_password' or 'sshkey' must be provided.")
			}
		} else {
			img, err := client.GetImage(image_name)
			if err != nil {
				_, err := client.GetSnapshot(image_name)
				if err != nil {
					return fmt.Errorf("Error fetching image/snapshot: %s", err)
				}
				isSnapshot = true
			}
			if img.Properties.Public == true && isSnapshot == false {
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
		return fmt.Errorf("Either 'image_name', or 'licenceType' must be set.")
	}

	if isSnapshot == true && (imagePassword != "" || len(publicKeys) > 0) {
		return fmt.Errorf("You can't pass 'image_password' and/or 'ssh keys' when creating a volume from a snapshot")
	}

	volume := &profitbricks.Volume{
		Properties: profitbricks.VolumeProperties{
			Name:          d.Get("name").(string),
			Size:          d.Get("size").(int),
			Type:          d.Get("disk_type").(string),
			ImagePassword: imagePassword,
			Image:         image,
			ImageAlias:    image_alias,
			Bus:           d.Get("bus").(string),
			LicenceType:   licenceType,
		},
	}

	if len(publicKeys) != 0 {
		volume.Properties.SSHKeys = publicKeys

	} else {
		volume.Properties.SSHKeys = nil
	}

	if _, ok := d.GetOk("availability_zone"); ok {
		raw := d.Get("availability_zone").(string)
		volume.Properties.AvailabilityZone = raw
	}

	volume, err := client.CreateVolume(dcId, *volume)

	if err != nil {
		return fmt.Errorf("An error occured while creating a volume: %s", err)
	}

	d.SetId(volume.ID)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, volume.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	volume, err = client.AttachVolume(dcId, serverId, volume.ID)
	if err != nil {
		return fmt.Errorf("An error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s", dcId, serverId, volume.ID, err)
	}

	d.Set("server_id", serverId)
	// Wait, catching any errors
	_, errState = getStateChangeConf(meta, d, volume.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.Set("server_id", "")
		}
		return errState
	}

	return resourceVolumeRead(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcId := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)
	volumeID := d.Id()

	volume, err := client.GetVolume(dcId, d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error occured while fetching a volume ID %s %s", d.Id(), err)
	}

	if volume.StatusCode > 299 {
		return fmt.Errorf("An error occured while fetching a volume ID %s %s", d.Id(), volume.Response)

	}

	_, err = client.GetAttachedVolume(dcId, serverID, volumeID)
	if err != nil {
		d.Set("server_id", "")
	}

	d.Set("name", volume.Properties.Name)
	d.Set("disk_type", volume.Properties.Type)
	d.Set("size", volume.Properties.Size)
	d.Set("bus", volume.Properties.Bus)
	d.Set("image_name", volume.Properties.Image)
	d.Set("image_alias", volume.Properties.ImageAlias)

	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	properties := profitbricks.VolumeProperties{}
	dcId := d.Get("datacenter_id").(string)

	if d.HasChange("name") {
		_, newValue := d.GetChange("name")
		properties.Name = newValue.(string)
	}
	if d.HasChange("disk_type") {
		_, newValue := d.GetChange("disk_type")
		properties.Type = newValue.(string)
	}
	if d.HasChange("size") {
		_, newValue := d.GetChange("size")
		properties.Size = newValue.(int)
	}
	if d.HasChange("bus") {
		_, newValue := d.GetChange("bus")
		properties.Bus = newValue.(string)
	}
	if d.HasChange("availability_zone") {
		_, newValue := d.GetChange("availability_zone")
		properties.AvailabilityZone = newValue.(string)
	}

	volume, err := client.UpdateVolume(dcId, d.Id(), properties)

	if err != nil {
		return fmt.Errorf("An error occured while updating a volume ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, volume.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	if volume.StatusCode > 299 {
		return fmt.Errorf("An error occured while updating a volume ID %s %s", d.Id(), volume.Response)

	}

	if d.HasChange("server_id") {
		_, newValue := d.GetChange("server_id")
		serverID := newValue.(string)
		volumeAttach, err := client.AttachVolume(dcId, serverID, volume.ID)
		if err != nil {
			return fmt.Errorf("An error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s", dcId, serverID, volumeAttach.ID, err)
		}

		// Wait, catching any errors
		_, errState = getStateChangeConf(meta, d, volumeAttach.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
		if errState != nil {
			return errState
		}
	}

	return resourceVolumeRead(d, meta)
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcId := d.Get("datacenter_id").(string)

	resp, err := client.DeleteVolume(dcId, d.Id())
	if err != nil {
		return fmt.Errorf("An error occured while deleting a volume ID %s %s", d.Id(), err)

	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
