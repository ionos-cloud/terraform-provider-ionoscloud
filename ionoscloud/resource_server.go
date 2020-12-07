package ionoscloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
	"golang.org/x/crypto/ssh"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServerImport,
		},
		Schema: map[string]*schema.Schema{
			// Server parameters
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cores": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"licence_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"boot_volume": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"boot_cdrom": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_family": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"boot_image": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"primary_nic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"firewallrule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Computed:      true,
				ConflictsWith: []string{"volume.0.image_password"},
			},
			"image_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"volume.0.image_name"},
			},
			"ssh_key_path": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"volume.0.ssh_key_path"},
				Optional:      true,
				Computed:      true,
			},
			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_name": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"image_name"},
							Deprecated:    "Please use image_name under server level",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if d.Get("image_name").(string) == new {
									return true
								}
								return false
							},
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
							Type:          schema.TypeString,
							Optional:      true,
							Deprecated:    "Please use image_password under server level",
							ConflictsWith: []string{"image_password"},
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if d.Get("image_password").(string) == new {
									return true
								}
								return false
							},
						},
						"licence_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ssh_key_path": {
							Type:       schema.TypeList,
							Elem:       &schema.Schema{Type: schema.TypeString},
							Optional:   true,
							Deprecated: "Please use ssh_key_path under server level",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if k == "volume.0.ssh_key_path.#" {
									if d.Get("ssh_key_path.#") == new {
										return true
									}
								}

								ssh_key_path := d.Get("volume.0.ssh_key_path").([]interface{})
								old_ssh_key_path := d.Get("ssh_key_path").([]interface{})

								if len(diffSlice(convertSlice(ssh_key_path), convertSlice(old_ssh_key_path))) == 0 {
									return true
								}

								return false
							},
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
							Computed: true,
						},
						"image_aliases": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
			"nic": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lan": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dhcp": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if new == "" {
									return true
								}
								return false
							},
						},
						"ips": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"nat": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"firewall_active": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"firewall": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"protocol": {
										Type:     schema.TypeString,
										Required: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											if strings.ToLower(old) == strings.ToLower(new) {
												return true
											}
											return false
										},
									},
									"source_mac": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source_ip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"target_ip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ips": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
									},
									"port_range_start": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											if v.(int) < 1 && v.(int) > 65534 {
												errors = append(errors, fmt.Errorf("Port start range must be between 1 and 65534"))
											}
											return
										},
									},

									"port_range_end": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											if v.(int) < 1 && v.(int) > 65534 {
												errors = append(errors, fmt.Errorf("Port end range must be between 1 and 65534"))
											}
											return
										},
									},
									"icmp_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"icmp_code": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	var image_alias string
	request := profitbricks.Server{
		Properties: profitbricks.ServerProperties{
			Name:  d.Get("name").(string),
			Cores: d.Get("cores").(int),
			RAM:   d.Get("ram").(int),
		},
	}
	dcId := d.Get("datacenter_id").(string)

	isSnapshot := false
	if v, ok := d.GetOk("availability_zone"); ok {
		request.Properties.AvailabilityZone = v.(string)
	}

	if v, ok := d.GetOk("cpu_family"); ok {
		if v.(string) != "" {
			request.Properties.CPUFamily = v.(string)
		}
	}

	volume := profitbricks.VolumeProperties{
		Size: d.Get("volume.0.size").(int),
		Type: d.Get("volume.0.disk_type").(string),
	}

	if v, ok := d.GetOk("volume.0.image_password"); ok {
		volume.ImagePassword = v.(string)
		d.Set("image_password", v.(string))
	}

	if v, ok := d.GetOk("image_password"); ok {
		volume.ImagePassword = v.(string)
	}

	if v, ok := d.GetOk("volume.0.licence_type"); ok {
		volume.LicenceType = v.(string)
	}

	if v, ok := d.GetOk("volume.0.availability_zone"); ok {
		volume.AvailabilityZone = v.(string)
	}

	if v, ok := d.GetOk("volume.0.name"); ok {
		volume.Name = v.(string)
	}

	if v, ok := d.GetOk("volume.0.bus"); ok {
		volume.Bus = v.(string)
	}

	var sshkey_path []interface{}

	if v, ok := d.GetOk("volume.0.ssh_key_path"); ok {
		sshkey_path = v.([]interface{})
		d.Set("ssh_key_path", v.([]interface{}))
	} else if v, ok := d.GetOk("ssh_key_path"); ok {
		sshkey_path = v.([]interface{})
		d.Set("ssh_key_path", v.([]interface{}))
	} else {
		d.Set("ssh_key_path", [][]string{})
	}

	var image, image_name string
	if v, ok := d.GetOk("volume.0.image_name"); ok {
		image_name = v.(string)
		d.Set("image_name", v.(string))
	} else if v, ok := d.GetOk("image_name"); ok {
		image_name = v.(string)
	} else {
		return fmt.Errorf("Either 'image_name' or 'volume.0.image_name' must be provided.")
	}

	if !IsValidUUID(image_name) {
		img, err := getImage(client, dcId, image_name, volume.Type)
		if err != nil {
			return err
		}
		if img != nil {
			image = img.ID
		}
		// if no image id was found with that name we look for a matching snapshot
		if image == "" {
			image = getSnapshotId(client, image_name)
			if image != "" {
				isSnapshot = true
			} else {
				dc, err := client.GetDatacenter(dcId)
				if err != nil {
					return fmt.Errorf("Error fetching datacenter %s: (%s)", dcId, err)
				}
				image_alias = getImageAlias(client, image_name, dc.Properties.Location)
			}
		}
		if image == "" && image_alias == "" {
			return fmt.Errorf("Could not find an image/imagealias/snapshot that matches %s ", image_name)
		}
		if volume.ImagePassword == "" && len(sshkey_path) == 0 && isSnapshot == false && img.Properties.Public {
			return fmt.Errorf("Either 'image_password' or 'ssh_key_path' must be provided.")
		}
	} else {
		img, err := client.GetImage(image_name)

		apiError, rsp := err.(profitbricks.ApiError)

		if err != nil {
			return fmt.Errorf("Error fetching image %s: (%s) - %+v", image_name, err, rsp)
		}

		if apiError.HttpStatusCode() == 404 {

			img, err := client.GetSnapshot(image_name)

			if apiError, ok := err.(profitbricks.ApiError); !ok {
				if apiError.HttpStatusCode() == 404 {
					return fmt.Errorf("image/snapshot: %s Not Found", img.Response)
				}
			}

			isSnapshot = true

		} else {
			if err != nil {
				return fmt.Errorf("Error fetching image/snapshot: %s", err)
			}
		}

		if img.Properties.Public == true && isSnapshot == false {

			if volume.ImagePassword == "" && len(sshkey_path) == 0 {
				return fmt.Errorf("Either 'image_password' or 'ssh_key_path' must be provided.")
			}

			img, err := getImage(client, d.Get("datacenter_id").(string), image_name, volume.Type)

			if err != nil {
				return err
			}

			if img != nil {
				image = img.ID
			}
		} else {
			img, err := client.GetImage(image_name)
			if err != nil {
				img, err := client.GetSnapshot(image_name)
				if err != nil {
					return fmt.Errorf("Error fetching image/snapshot: %s", img.Response)
				}
				isSnapshot = true
			}
			if img.Properties.Public == true && isSnapshot == false {
				if volume.ImagePassword == "" && len(sshkey_path) == 0 {
					return fmt.Errorf("Either 'image_password' or 'ssh_key_path' must be provided.")
				}
				image = image_name
			} else {
				image = image_name
			}
		}
	}

	if len(sshkey_path) != 0 {
		var publicKeys []string
		for _, path := range sshkey_path {
			log.Printf("[DEBUG] Reading file %s", path)
			publicKey, err := readPublicKey(path.(string))
			if err != nil {
				return fmt.Errorf("Error fetching sshkey from file (%s) %s", path, err.Error())
			}
			publicKeys = append(publicKeys, publicKey)
		}
		if len(publicKeys) > 0 {
			volume.SSHKeys = publicKeys
		}
	}

	if image == "" && volume.LicenceType == "" && image_alias == "" && !isSnapshot {
		return fmt.Errorf("Either 'image', 'licenceType', or 'imageAlias' must be set.")
	}

	if isSnapshot == true && (volume.ImagePassword != "" || len(sshkey_path) > 0) {
		return fmt.Errorf("Passwords/SSH keys are not supported for snapshots.")
	}

	volume.ImageAlias = image_alias
	volume.Image = image

	request.Entities = &profitbricks.ServerEntities{
		Volumes: &profitbricks.Volumes{
			Items: []profitbricks.Volume{
				{
					Properties: volume,
				},
			},
		},
	}

	if _, ok := d.GetOk("nic"); ok {
		nic := profitbricks.Nic{Properties: &profitbricks.NicProperties{
			Lan: d.Get("nic.0.lan").(int),
		}}

		if v, ok := d.GetOk("nic.0.name"); ok {
			nic.Properties.Name = v.(string)
		}

		nic.Properties.Dhcp = boolAddr(d.Get("nic.0.dhcp").(bool))
		nic.Properties.FirewallActive = boolAddr(d.Get("nic.0.firewall_active").(bool))
		nic.Properties.Nat = boolAddr(d.Get("nic.0.nat").(bool))

		if v, ok := d.GetOk("nic.0.ip"); ok {
			ips := strings.Split(v.(string), ",")
			if len(ips) > 0 {
				nic.Properties.Ips = ips
			}
		}

		log.Printf("[DEBUG] dhcp nic before%t", *nic.Properties.Dhcp)
		request.Entities.Nics = &profitbricks.Nics{
			Items: []profitbricks.Nic{
				nic,
			},
		}
		log.Printf("[DEBUG] dhcp nic after %t", *nic.Properties.Dhcp)
		log.Printf("[DEBUG] dhcp %t", *request.Entities.Nics.Items[0].Properties.Dhcp)

		if _, ok := d.GetOk("nic.0.firewall"); ok {
			firewall := profitbricks.FirewallRule{
				Properties: profitbricks.FirewallruleProperties{
					Protocol: d.Get("nic.0.firewall.0.protocol").(string),
				},
			}

			if v, ok := d.GetOk("nic.0.firewall.0.name"); ok {
				firewall.Properties.Name = v.(string)
			}

			if v, ok := d.GetOk("nic.0.firewall.0.source_mac"); ok {
				val := v.(string)
				firewall.Properties.SourceMac = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.source_ip"); ok {
				val := v.(string)
				firewall.Properties.SourceIP = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.target_ip"); ok {
				val := v.(string)
				firewall.Properties.TargetIP = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.port_range_start"); ok {
				val := v.(int)
				firewall.Properties.PortRangeStart = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.port_range_end"); ok {
				val := v.(int)
				firewall.Properties.PortRangeEnd = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.icmp_type"); ok {
				tempIcmpType := v.(string)
				if tempIcmpType != "" {
					i, _ := strconv.Atoi(tempIcmpType)
					firewall.Properties.IcmpType = &i
				}
			}
			if v, ok := d.GetOk("nic.0.firewall.0.icmp_code"); ok {
				tempIcmpCode := v.(string)
				if tempIcmpCode != "" {
					i, _ := strconv.Atoi(tempIcmpCode)
					firewall.Properties.IcmpCode = &i
				}
			}
			request.Entities.Nics.Items[0].Entities = &profitbricks.NicEntities{
				FirewallRules: &profitbricks.FirewallRules{
					Items: []profitbricks.FirewallRule{
						firewall,
					},
				},
			}
		}
	}

	if len(request.Entities.Nics.Items[0].Properties.Ips) == 0 {
		request.Entities.Nics.Items[0].Properties.Ips = nil
	}

	server, err := client.CreateServer(d.Get("datacenter_id").(string), request)

	if err != nil {
		return fmt.Errorf(
			"Error creating server: (%s)", err)
	}
	d.SetId(server.ID)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, server.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}
	server, err = client.GetServer(d.Get("datacenter_id").(string), server.ID)
	if err != nil {
		return fmt.Errorf("Error fetching server: (%s)", err)
	}

	firewallRules, err := client.ListFirewallRules(d.Get("datacenter_id").(string), server.ID, server.Entities.Nics.Items[0].ID)

	if len(firewallRules.Items) > 0 {
		d.Set("firewallrule_id", firewallRules.Items[0].ID)
	}

	d.Set("primary_nic", server.Entities.Nics.Items[0].ID)
	if len(server.Entities.Nics.Items[0].Properties.Ips) > 0 {
		d.SetConnInfo(map[string]string{
			"type":     "ssh",
			"host":     server.Entities.Nics.Items[0].Properties.Ips[0],
			"password": request.Entities.Volumes.Items[0].Properties.ImagePassword,
		})
	}
	return resourceServerRead(d, meta)
}

func GetFirewallResource(d *schema.ResourceData, path string) profitbricks.FirewallRule {

	firewall := profitbricks.FirewallRule{
		Properties: profitbricks.FirewallruleProperties{},
	}
	if v, ok := d.GetOk(path + ".protocol"); ok {
		firewall.Properties.Protocol = v.(string)
	}

	if v, ok := d.GetOk(path + ".name"); ok {
		firewall.Properties.Name = v.(string)
	}

	if v, ok := d.GetOk(path + ".source_mac"); ok {
		val := v.(string)
		firewall.Properties.SourceMac = &val
	}

	if v, ok := d.GetOk(path + ".source_ip"); ok {
		val := v.(string)
		firewall.Properties.SourceIP = &val
	}

	if v, ok := d.GetOk(path + ".target_ip"); ok {
		val := v.(string)
		firewall.Properties.TargetIP = &val
	}

	if v, ok := d.GetOk(path + ".port_range_start"); ok {
		val := v.(int)
		firewall.Properties.PortRangeStart = &val
	}

	if v, ok := d.GetOk(path + ".port_range_end"); ok {
		val := v.(int)
		firewall.Properties.PortRangeEnd = &val
	}

	if v, ok := d.GetOk(path + ".icmp_type"); ok {
		tempIcmpType := v.(string)
		if tempIcmpType != "" {
			i, _ := strconv.Atoi(tempIcmpType)
			firewall.Properties.IcmpType = &i
		}
	}
	if v, ok := d.GetOk(path + ".icmp_code"); ok {
		tempIcmpCode := v.(string)
		if tempIcmpCode != "" {
			i, _ := strconv.Atoi(tempIcmpCode)
			firewall.Properties.IcmpCode = &i
		}
	}
	return firewall
}

func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcId := d.Get("datacenter_id").(string)
	serverId := d.Id()

	server, err := client.GetServer(dcId, serverId)
	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error occured while fetching a server ID %s %s", d.Id(), err)
	}
	d.Set("name", server.Properties.Name)
	d.Set("cores", server.Properties.Cores)
	d.Set("ram", server.Properties.RAM)
	d.Set("availability_zone", server.Properties.AvailabilityZone)
	d.Set("cpu_family", server.Properties.CPUFamily)
	d.Set("boot_image", server.Entities.Volumes.Items[0].Properties.Image)

	if primarynic, ok := d.GetOk("primary_nic"); ok {
		d.Set("primary_nic", primarynic.(string))

		nic, err := client.GetNic(dcId, serverId, primarynic.(string))
		if err != nil {
			return fmt.Errorf("Error occured while fetching nic %s for server ID %s %s", primarynic.(string), d.Id(), err)
		}

		if len(nic.Properties.Ips) > 0 {
			d.Set("primary_ip", nic.Properties.Ips[0])
		}

		network := map[string]interface{}{
			"lan":             nic.Properties.Lan,
			"name":            nic.Properties.Name,
			"dhcp":            *nic.Properties.Dhcp,
			"nat":             *nic.Properties.Nat,
			"firewall_active": *nic.Properties.FirewallActive,
			"ips":             nic.Properties.Ips,
		}

		if len(nic.Properties.Ips) > 0 {
			network["ip"] = nic.Properties.Ips[0]
		}

		if firewall_id, ok := d.GetOk("firewallrule_id"); ok {
			firewall, err := client.GetFirewallRule(dcId, serverId, primarynic.(string), firewall_id.(string))
			if err != nil {
				return fmt.Errorf("Error occured while fetching firewallrule %s for server ID %s %s", firewall_id.(string), serverId, err)
			}

			fw := map[string]interface{}{
				"protocol": firewall.Properties.Protocol,
				"name":     firewall.Properties.Name,
			}

			if firewall.Properties.SourceMac != nil {
				fw["source_mac"] = *firewall.Properties.SourceMac
			}

			if firewall.Properties.SourceIP != nil {
				fw["source_ip"] = *firewall.Properties.SourceIP
			}

			if firewall.Properties.TargetIP != nil {
				fw["target_ip"] = *firewall.Properties.TargetIP
			}

			if firewall.Properties.PortRangeStart != nil {
				fw["port_range_start"] = *firewall.Properties.PortRangeStart
			}

			if firewall.Properties.PortRangeEnd != nil {
				fw["port_range_end"] = *firewall.Properties.PortRangeEnd
			}

			if firewall.Properties.IcmpType != nil {
				fw["icmp_type"] = *firewall.Properties.IcmpType
			}

			if firewall.Properties.IcmpCode != nil {
				fw["icmp_code"] = *firewall.Properties.IcmpCode
			}

			network["firewall"] = []map[string]interface{}{fw}
		}

		networks := []map[string]interface{}{network}
		if err := d.Set("nic", networks); err != nil {
			return fmt.Errorf("[ERROR] unable saving nic to state IonosCloud Server (%s): %s", serverId, err)
		}
	}

	if server.Properties.BootVolume != nil {
		d.Set("boot_volume", server.Properties.BootVolume.ID)
		volumeObj, err := client.GetAttachedVolume(dcId, serverId, server.Properties.BootVolume.ID)
		if err == nil {
			volumeItem := map[string]interface{}{
				"name":              volumeObj.Properties.Name,
				"disk_type":         volumeObj.Properties.Type,
				"size":              volumeObj.Properties.Size,
				"licence_type":      volumeObj.Properties.LicenceType,
				"bus":               volumeObj.Properties.Bus,
				"availability_zone": volumeObj.Properties.AvailabilityZone,
			}

			volumesList := []map[string]interface{}{volumeItem}
			if err := d.Set("volume", volumesList); err != nil {
				return fmt.Errorf("[DEBUG] Error saving volume to state for IonosCloud server (%s): %s", d.Id(), err)
			}
		}
	}

	_, err = client.GetAttachedVolume(dcId, d.Id(), d.Get("boot_volume").(string))
	if err != nil {
		d.Set("volume", nil)
	}

	if server.Properties.BootCdrom != nil {
		d.Set("boot_cdrom", server.Properties.BootCdrom.ID)
	}
	return nil
}

func boolAddr(b bool) *bool {
	return &b
}

func resourceServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcId := d.Get("datacenter_id").(string)

	request := profitbricks.ServerProperties{}

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		request.Name = n.(string)
	}
	if d.HasChange("cores") {
		_, n := d.GetChange("cores")
		request.Cores = n.(int)
	}
	if d.HasChange("ram") {
		_, n := d.GetChange("ram")
		request.RAM = n.(int)
	}
	if d.HasChange("availability_zone") {
		_, n := d.GetChange("availability_zone")
		request.AvailabilityZone = n.(string)
	}
	if d.HasChange("cpu_family") {
		_, n := d.GetChange("cpu_family")
		request.CPUFamily = n.(string)
	}
	server, err := client.UpdateServer(dcId, d.Id(), request)

	if err != nil {
		return fmt.Errorf("Error occured while updating server ID %s %s", d.Id(), err)
	}

	_, errState := getStateChangeConf(meta, d, server.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}
	// Volume stuff
	if d.HasChange("volume") {
		boot_volume := d.Get("boot_volume").(string)
		_, err = client.GetAttachedVolume(dcId, d.Id(), boot_volume)

		if err != nil {

			volumeAttach, err := client.AttachVolume(dcId, d.Id(), boot_volume)
			if err != nil {
				return fmt.Errorf("An error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s", dcId, d.Id(), boot_volume, err)
			}

			// Wait, catching any errors
			_, errState = getStateChangeConf(meta, d, volumeAttach.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
			if errState != nil {
				return errState
			}
		}

		properties := profitbricks.VolumeProperties{}

		if v, ok := d.GetOk("volume.0.name"); ok {
			properties.Name = v.(string)
		}

		if v, ok := d.GetOk("volume.0.size"); ok {
			properties.Size = v.(int)
		}

		if v, ok := d.GetOk("volume.0.bus"); ok {
			properties.Bus = v.(string)
		}

		volume, err := client.UpdateVolume(d.Get("datacenter_id").(string), boot_volume, properties)

		if err != nil {
			return fmt.Errorf("Error patching volume (%s) (%s)", d.Id(), err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, volume.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
		if errState != nil {
			return errState
		}
	}

	// Nic stuff
	if d.HasChange("nic") {
		nic := &profitbricks.Nic{}
		for _, n := range server.Entities.Nics.Items {
			if n.ID == d.Get("primary_nic").(string) {
				nic = &n
				break
			}
		}

		properties := profitbricks.NicProperties{
			Lan: d.Get("nic.0.lan").(int),
		}

		if v, ok := d.GetOk("nic.0.name"); ok {
			properties.Name = v.(string)
		}

		if v, ok := d.GetOk("nic.0.ip"); ok {
			ips := strings.Split(v.(string), ",")
			if len(ips) > 0 {
				nic.Properties.Ips = ips
			}
		}

		properties.Dhcp = boolAddr(d.Get("nic.0.dhcp").(bool))

		properties.Nat = boolAddr(d.Get("nic.0.nat").(bool))
		properties.FirewallActive = boolAddr(d.Get("nic.0.firewall_active").(bool))

		if d.HasChange("nic.0.firewall") {

			firewall := GetFirewallResource(d, "nic.0.firewall")
			nic.Entities = &profitbricks.NicEntities{
				FirewallRules: &profitbricks.FirewallRules{
					Items: []profitbricks.FirewallRule{
						firewall,
					},
				},
			}
		}
		mProp, _ := json.Marshal(properties)
		log.Printf("[DEBUG] Updating props: %s", string(mProp))
		nic, err := client.UpdateNic(d.Get("datacenter_id").(string), server.ID, nic.ID, properties)
		if err != nil {
			return fmt.Errorf(
				"Error updating nic (%s)", err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, nic.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
		if errState != nil {
			return errState
		}

	}

	return resourceServerRead(d, meta)
}

func resourceServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcId := d.Get("datacenter_id").(string)

	server, err := client.GetServer(dcId, d.Id())

	if err != nil {
		return fmt.Errorf("Error occured while fetching a server ID %s %s", d.Id(), err)
	}

	if server.Properties.BootVolume != nil {
		resp, err := client.DeleteVolume(dcId, server.Properties.BootVolume.ID)
		if err != nil {
			return fmt.Errorf("Error occured while delete volume %s of server ID %s %s", server.Properties.BootVolume.ID, d.Id(), err)
		}
		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
		if errState != nil {
			return errState
		}
	}

	resp, err := client.DeleteServer(dcId, d.Id())
	if err != nil {
		return fmt.Errorf("An error occured while deleting a server ID %s %s", d.Id(), err)

	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}

// Reads public key from file and returns key string iff valid
func readPublicKey(path string) (key string, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	if err != nil {
		return "", err
	}
	return string(ssh.MarshalAuthorizedKey(pubKey)[:]), nil
}
