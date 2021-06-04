package ionoscloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"template_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
			"type": {
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"image_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Computed:      true,
				ConflictsWith: []string{"volume.0.image_password"},
			},
			"image": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"volume.0.image"},
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
						"image": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"image"},
							Deprecated:    "Please use image under server level",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if d.Get("image").(string) == new {
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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

								sshKeyPath := d.Get("volume.0.ssh_key_path").([]interface{})
								oldSshKeyPath := d.Get("ssh_key_path").([]interface{})

								if len(diffSlice(convertSlice(sshKeyPath), convertSlice(oldSshKeyPath))) == 0 {
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
					},
				},
			},
			"nic": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
										ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
												errors = append(errors, fmt.Errorf("port start range must be between 1 and 65534"))
											}
											return
										},
									},

									"port_range_end": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											if v.(int) < 1 && v.(int) > 65534 {
												errors = append(errors, fmt.Errorf("port end range must be between 1 and 65534"))
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
	client := meta.(*ionoscloud.APIClient)

	datacenterId := d.Get("datacenter_id").(string)
	serverName := d.Get("name").(string)
	serverCores := int32(d.Get("cores").(int))
	serverRam := int32(d.Get("ram").(int))
	request := ionoscloud.Server{
		Properties: &ionoscloud.ServerProperties{
			Name:  &serverName,
			Cores: &serverCores,
			Ram:   &serverRam,
		},
	}

	isSnapshot := false

	if v, ok := d.GetOk("template_uuid"); ok {
		vStr := v.(string)
		request.Properties.TemplateUuid = &vStr
	}

	if v, ok := d.GetOk("type"); ok {
		vStr := v.(string)
		request.Properties.Type = &vStr
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		vStr := v.(string)
		request.Properties.AvailabilityZone = &vStr
	}

	if v, ok := d.GetOk("cpu_family"); ok {
		if v.(string) != "" {
			vStr := v.(string)
			request.Properties.CpuFamily = &vStr
		}
	}

	volumeSize := float32(d.Get("volume.0.size").(int))
	volumeType := d.Get("volume.0.disk_type").(string)
	volume := ionoscloud.VolumeProperties{
		Size: &volumeSize,
		Type: &volumeType,
	}

	if v, ok := d.GetOk("volume.0.image_password"); ok {
		vStr := v.(string)
		volume.ImagePassword = &vStr
		if err := d.Set("image_password", vStr); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("image_password"); ok {
		vStr := v.(string)
		volume.ImagePassword = &vStr
	}

	if v, ok := d.GetOk("volume.0.licence_type"); ok {
		vStr := v.(string)
		volume.LicenceType = &vStr
	}

	if v, ok := d.GetOk("volume.0.availability_zone"); ok {
		vStr := v.(string)
		volume.AvailabilityZone = &vStr
	}

	if v, ok := d.GetOk("volume.0.name"); ok {
		vStr := v.(string)
		volume.Name = &vStr
	}

	if v, ok := d.GetOk("volume.0.bus"); ok {
		vStr := v.(string)
		volume.Bus = &vStr
	}

	var sshKeyPath []interface{}

	if v, ok := d.GetOk("volume.0.ssh_key_path"); ok {
		sshKeyPath = v.([]interface{})
		if err := d.Set("ssh_key_path", v.([]interface{})); err != nil {
			return err
		}
	} else if v, ok := d.GetOk("ssh_key_path"); ok {
		sshKeyPath = v.([]interface{})
		if err := d.Set("ssh_key_path", v.([]interface{})); err != nil {
			return err
		}
	} else {
		if err := d.Set("ssh_key_path", [][]string{}); err != nil {
			return err
		}
	}

	var image string
	var imageInput string

	if v, ok := d.GetOk("volume.0.image"); ok {
		imageInput = v.(string)
		if err := d.Set("image", v.(string)); err != nil {
			return err
		}
	} else if v, ok := d.GetOk("image"); ok {
		imageInput = v.(string)
	} else {
		return fmt.Errorf("either 'image' or 'volume.0.image' must be provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if !IsValidUUID(imageInput) {
		img, err := getImage(client, datacenterId, imageInput, *volume.Type)
		if err != nil {
			return err
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
				return fmt.Errorf("no image or snapshot with id %s found", imageInput)
			}
		}

		if volume.ImagePassword == nil && len(sshKeyPath) == 0 && isSnapshot == false && img.Properties.Public != nil && *img.Properties.Public {
			return fmt.Errorf("either 'image_password' or 'ssh_key_path' must be provided")
		}
	} else {
		img, apiResponse, err := client.ImagesApi.ImagesFindById(ctx, imageInput).Execute()

		if err != nil {
			if apiResponse != nil && apiResponse.Response.StatusCode == 404 {

				log.Printf("[DEBUG] image %s not found; trying snapshots\n", imageInput)
				snap, apiResponse, err := client.SnapshotsApi.SnapshotsFindById(ctx, imageInput).Execute()

				if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
					return fmt.Errorf("image/snapshot: %s Not Found", string(apiResponse.Payload))
				} else if err != nil {
					return fmt.Errorf("error fetching image/snapshot info: %s, %s", err, responseBody(apiResponse))
				}

				isSnapshot = true

				image = *snap.Id

			} else {
				return fmt.Errorf("error fetching image info: %s, %s", err, responseBody(apiResponse))
			}
		} else {
			if img.Properties != nil && img.Properties.Public != nil && *img.Properties.Public == true && isSnapshot == false {
				if volume.ImagePassword == nil && len(sshKeyPath) == 0 {
					return fmt.Errorf("either 'image_password' or 'ssh_key_path' must be provided")
				}
			}

			image = *img.Id

		}
	}

	if isSnapshot == true && (volume.ImagePassword != nil || len(sshKeyPath) > 0) {
		return fmt.Errorf("passwords/SSH keys are not supported for snapshots")
	}

	volume.Image = &image

	if len(sshKeyPath) != 0 {
		var publicKeys []string
		for _, path := range sshKeyPath {
			log.Printf("[DEBUG] Reading file %s", path)
			publicKey, err := readPublicKey(path.(string))
			if err != nil {
				return fmt.Errorf("error fetching sshkey from file (%s) %s", path, err.Error())
			}
			publicKeys = append(publicKeys, publicKey)
		}
		if len(publicKeys) > 0 {
			volume.SshKeys = &publicKeys
		}
	}

	request.Entities = &ionoscloud.ServerEntities{
		Volumes: &ionoscloud.AttachedVolumes{
			Items: &[]ionoscloud.Volume{
				{
					Properties: &volume,
				},
			},
		},
	}

	if _, ok := d.GetOk("nic"); ok {
		lanInt := int32(d.Get("nic.0.lan").(int))
		nic := ionoscloud.Nic{Properties: &ionoscloud.NicProperties{
			Lan: &lanInt,
		}}

		if v, ok := d.GetOk("nic.0.name"); ok {
			vStr := v.(string)
			nic.Properties.Name = &vStr
		}

		nic.Properties.Dhcp = boolAddr(d.Get("nic.0.dhcp").(bool))
		nic.Properties.FirewallActive = boolAddr(d.Get("nic.0.firewall_active").(bool))

		if v, ok := d.GetOk("nic.0.ip"); ok {
			ips := strings.Split(v.(string), ",")
			if len(ips) > 0 {
				nic.Properties.Ips = &ips
			}
		}

		log.Printf("[DEBUG] dhcp nic before%t", *nic.Properties.Dhcp)

		request.Entities.Nics = &ionoscloud.Nics{
			Items: &[]ionoscloud.Nic{
				nic,
			},
		}

		log.Printf("[DEBUG] dhcp nic after %t", *nic.Properties.Dhcp)
		log.Printf("[DEBUG] dhcp %t", *(*request.Entities.Nics.Items)[0].Properties.Dhcp)

		if _, ok := d.GetOk("nic.0.firewall"); ok {
			protocolStr := d.Get("nic.0.firewall.0.protocol").(string)
			firewall := ionoscloud.FirewallRule{
				Properties: &ionoscloud.FirewallruleProperties{
					Protocol: &protocolStr,
				},
			}

			if v, ok := d.GetOk("nic.0.firewall.0.name"); ok {
				vStr := v.(string)
				firewall.Properties.Name = &vStr
			}

			if v, ok := d.GetOk("nic.0.firewall.0.source_mac"); ok {
				val := v.(string)
				firewall.Properties.SourceMac = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.source_ip"); ok {
				val := v.(string)
				firewall.Properties.SourceIp = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.target_ip"); ok {
				val := v.(string)
				firewall.Properties.TargetIp = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.port_range_start"); ok {
				val := int32(v.(int))
				firewall.Properties.PortRangeStart = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.port_range_end"); ok {
				val := int32(v.(int))
				firewall.Properties.PortRangeEnd = &val
			}

			if v, ok := d.GetOk("nic.0.firewall.0.icmp_type"); ok {
				tempIcmpType := v.(string)
				if tempIcmpType != "" {
					i, _ := strconv.Atoi(tempIcmpType)
					iInt32 := int32(i)
					firewall.Properties.IcmpType = &iInt32
				}
			}
			if v, ok := d.GetOk("nic.0.firewall.0.icmp_code"); ok {
				tempIcmpCode := v.(string)
				if tempIcmpCode != "" {
					i, _ := strconv.Atoi(tempIcmpCode)
					iInt32 := int32(i)
					firewall.Properties.IcmpCode = &iInt32
				}
			}
			(*request.Entities.Nics.Items)[0].Entities = &ionoscloud.NicEntities{
				Firewallrules: &ionoscloud.FirewallRules{
					Items: &[]ionoscloud.FirewallRule{
						firewall,
					},
				},
			}
		}
	}

	if (*request.Entities.Nics.Items)[0].Properties.Ips != nil {
		if len(*(*request.Entities.Nics.Items)[0].Properties.Ips) == 0 {
			*(*request.Entities.Nics.Items)[0].Properties.Ips = nil
		}
	}

	server, apiResponse, err := client.ServersApi.DatacentersServersPost(ctx, d.Get("datacenter_id").(string)).Server(request).Execute()

	if err != nil {
		return fmt.Errorf(
			"Error creating server: (%s) \n apiEroor: %v", err, responseBody(apiResponse))
	}
	d.SetId(*server.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}
	server, _, err = client.ServersApi.DatacentersServersFindById(ctx, d.Get("datacenter_id").(string), *server.Id).Execute()

	if err != nil {
		return fmt.Errorf("error fetching server: (%s)", err)
	}

	firewallRules, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, d.Get("datacenter_id").(string),
		*server.Id, *(*server.Entities.Nics.Items)[0].Id).Execute()

	if err != nil {
		return fmt.Errorf("an error occurred while fetching firewall rules: %s", err)
	}

	if firewallRules.Items != nil {
		if len(*firewallRules.Items) > 0 {
			if err := d.Set("firewallrule_id", *(*firewallRules.Items)[0].Id); err != nil {
				return err
			}
		}
	}

	if (*server.Entities.Nics.Items)[0].Id != nil {
		err := d.Set("primary_nic", *(*server.Entities.Nics.Items)[0].Id)
		if err != nil {
			return fmt.Errorf("error while setting primary nic %s: %s", d.Id(), err)
		}
	}

	if (*server.Entities.Nics.Items)[0].Properties.Ips != nil &&
		len(*(*server.Entities.Nics.Items)[0].Properties.Ips) > 0 &&
		request.Entities.Volumes.Items != nil &&
		len(*request.Entities.Volumes.Items) > 0 &&
		(*request.Entities.Volumes.Items)[0].Properties != nil &&
		(*request.Entities.Volumes.Items)[0].Properties.ImagePassword != nil {

		d.SetConnInfo(map[string]string{
			"type":     "ssh",
			"host":     (*(*server.Entities.Nics.Items)[0].Properties.Ips)[0],
			"password": *(*request.Entities.Volumes.Items)[0].Properties.ImagePassword,
		})
	}

	return resourceServerRead(d, meta)
}

func GetFirewallResource(d *schema.ResourceData, path string) ionoscloud.FirewallRule {

	firewall := ionoscloud.FirewallRule{
		Properties: &ionoscloud.FirewallruleProperties{},
	}
	if v, ok := d.GetOk(path + ".protocol"); ok {
		vStr := v.(string)
		firewall.Properties.Protocol = &vStr
	}

	if v, ok := d.GetOk(path + ".name"); ok {
		vStr := v.(string)
		firewall.Properties.Name = &vStr
	}

	if v, ok := d.GetOk(path + ".source_mac"); ok {
		val := v.(string)
		firewall.Properties.SourceMac = &val
	}

	if v, ok := d.GetOk(path + ".source_ip"); ok {
		val := v.(string)
		firewall.Properties.SourceIp = &val
	}

	if v, ok := d.GetOk(path + ".target_ip"); ok {
		val := v.(string)
		firewall.Properties.TargetIp = &val
	}

	if v, ok := d.GetOk(path + ".port_range_start"); ok {
		val := int32(v.(int))
		firewall.Properties.PortRangeStart = &val
	}

	if v, ok := d.GetOk(path + ".port_range_end"); ok {
		val := int32(v.(int))
		firewall.Properties.PortRangeEnd = &val
	}

	if v, ok := d.GetOk(path + ".icmp_type"); ok {
		tempIcmpType := v.(string)
		if tempIcmpType != "" {
			i, _ := strconv.Atoi(tempIcmpType)
			iInt32 := int32(i)
			firewall.Properties.IcmpType = &iInt32
		}
	}
	if v, ok := d.GetOk(path + ".icmp_code"); ok {
		tempIcmpCode := v.(string)
		if tempIcmpCode != "" {
			i, _ := strconv.Atoi(tempIcmpCode)
			iInt32 := int32(i)
			firewall.Properties.IcmpCode = &iInt32
		}
	}
	return firewall
}

func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)
	serverId := d.Id()

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}

	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, serverId).Execute()
	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error occured while fetching a server ID %s %s", d.Id(), err)
	}

	if server.Properties.TemplateUuid != nil {
		if err := d.Set("template_uuid", *server.Properties.TemplateUuid); err != nil {
			return err
		}
	}

	if server.Properties.Name != nil {
		if err := d.Set("name", *server.Properties.Name); err != nil {
			return err
		}
	}

	if server.Properties.Cores != nil {
		if err := d.Set("cores", *server.Properties.Cores); err != nil {
			return err
		}
	}

	if server.Properties.Ram != nil {
		if err := d.Set("ram", *server.Properties.Ram); err != nil {
			return err
		}
	}

	if server.Properties.Type != nil {
		if err := d.Set("type", *server.Properties.Type); err != nil {
			return err
		}
	}

	if server.Properties.AvailabilityZone != nil {
		if err := d.Set("availability_zone", *server.Properties.AvailabilityZone); err != nil {
			return err
		}
	}

	if server.Properties.CpuFamily != nil {
		if err := d.Set("cpu_family", *server.Properties.CpuFamily); err != nil {
			return err
		}
	}

	if server.Entities.Volumes != nil && len(*server.Entities.Volumes.Items) > 0 {
		if err := d.Set("boot_image", *(*server.Entities.Volumes.Items)[0].Properties.Image); err != nil {
			return err
		}
	}

	if primarynic, ok := d.GetOk("primary_nic"); ok {
		if err := d.Set("primary_nic", primarynic.(string)); err != nil {
			return err
		}

		nic, _, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, dcId, serverId, primarynic.(string)).Execute()
		if err != nil {
			return fmt.Errorf("error occured while fetching nic %s for server ID %s %s", primarynic.(string), d.Id(), err)
		}

		if len(*nic.Properties.Ips) > 0 {
			if err := d.Set("primary_ip", (*nic.Properties.Ips)[0]); err != nil {
				return err
			}
		}

		network := map[string]interface{}{
			"dhcp":            *nic.Properties.Dhcp,
			"firewall_active": *nic.Properties.FirewallActive,
		}

		if nic.Properties.Lan != nil {
			network["lan"] = *nic.Properties.Lan
		}

		if nic.Properties.Name != nil {
			network["name"] = *nic.Properties.Name
		}

		if nic.Properties.Ips != nil {
			network["ips"] = *nic.Properties.Ips
		}

		if nic.Properties.Mac != nil {
			network["mac"] = *nic.Properties.Mac
		}

		if len(*nic.Properties.Ips) > 0 {
			network["ip"] = (*nic.Properties.Ips)[0]
		}

		if firewallId, ok := d.GetOk("firewallrule_id"); ok {
			firewall, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, dcId, serverId, primarynic.(string), firewallId.(string)).Execute()
			if err != nil {
				return fmt.Errorf("error occured while fetching firewallrule %s for server ID %s %s", firewallId.(string), serverId, err)
			}

			fw := map[string]interface{}{
				"protocol": *firewall.Properties.Protocol,
				"name":     *firewall.Properties.Name,
			}

			if firewall.Properties.SourceMac != nil {
				fw["source_mac"] = *firewall.Properties.SourceMac
			}

			if firewall.Properties.SourceIp != nil {
				fw["source_ip"] = *firewall.Properties.SourceIp
			}

			if firewall.Properties.TargetIp != nil {
				fw["target_ip"] = *firewall.Properties.TargetIp
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
		if server.Properties.BootVolume.Id != nil {
			if err := d.Set("boot_volume", *server.Properties.BootVolume.Id); err != nil {
				return err
			}
		}
		volumeObj, _, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, serverId, *server.Properties.BootVolume.Id).Execute()
		if err == nil {
			volumeItem := map[string]interface{}{}

			if volumeObj.Properties.Name != nil {
				volumeItem["name"] = *volumeObj.Properties.Name
			}

			if volumeObj.Properties.Type != nil {
				volumeItem["disk_type"] = *volumeObj.Properties.Type
			}

			if volumeObj.Properties.Size != nil {
				volumeItem["size"] = *volumeObj.Properties.Size
			}

			if volumeObj.Properties.LicenceType != nil {
				volumeItem["licence_type"] = *volumeObj.Properties.LicenceType
			}

			if volumeObj.Properties.Bus != nil {
				volumeItem["bus"] = *volumeObj.Properties.Bus
			}

			if volumeObj.Properties.AvailabilityZone != nil {
				volumeItem["availability_zone"] = *volumeObj.Properties.AvailabilityZone
			}

			volumesList := []map[string]interface{}{volumeItem}
			if err := d.Set("volume", volumesList); err != nil {
				return fmt.Errorf("[DEBUG] Error saving volume to state for IonosCloud server (%s): %s", d.Id(), err)
			}
		}
	}

	bootVolume, ok := d.GetOk("boot_volume")
	if ok && len(bootVolume.(string)) > 0 {
		_, _, err = client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, d.Id(), bootVolume.(string)).Execute()
		if err != nil {
			if err := d.Set("volume", nil); err != nil {
				return err
			}
		}
	}

	if server.Properties.BootCdrom != nil {
		if err := d.Set("boot_cdrom", *server.Properties.BootCdrom.Id); err != nil {
			return err
		}
	}
	return nil
}

func boolAddr(b bool) *bool {
	return &b
}

func resourceServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)
	request := ionoscloud.ServerProperties{}

	if d.HasChange("template_uuid") {
		_, n := d.GetChange("template_uuid")
		nStr := n.(string)
		request.TemplateUuid = &nStr
	}
	if d.HasChange("name") {
		_, n := d.GetChange("name")
		nStr := n.(string)
		request.Name = &nStr
	}
	if d.HasChange("cores") {
		_, n := d.GetChange("cores")
		nInt := int32(n.(int))
		request.Cores = &nInt
	}
	if d.HasChange("ram") {
		_, n := d.GetChange("ram")
		nInt := int32(n.(int))
		request.Ram = &nInt
	}
	if d.HasChange("type") {
		_, n := d.GetChange("type")
		nStr := n.(string)
		request.Type = &nStr
	}
	if d.HasChange("availability_zone") {
		_, n := d.GetChange("availability_zone")
		nStr := n.(string)
		request.AvailabilityZone = &nStr
	}
	if d.HasChange("cpu_family") {
		_, n := d.GetChange("cpu_family")
		nStr := n.(string)
		request.CpuFamily = &nStr
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)
	if cancel != nil {
		defer cancel()
	}

	server, apiResponse, err := client.ServersApi.DatacentersServersPatch(ctx, dcId, d.Id()).Server(request).Execute()

	if err != nil {
		return fmt.Errorf("error occured while updating server ID %s %s", d.Id(), err)
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}
	// Volume stuff
	if d.HasChange("volume") {
		bootVolume := d.Get("boot_volume").(string)
		_, _, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, d.Id(), bootVolume).Execute()

		if err != nil {
			volume := ionoscloud.Volume{
				Id: &bootVolume,
			}
			_, apiResponse, err := client.ServersApi.DatacentersServersVolumesPost(ctx, dcId, d.Id()).Volume(volume).Execute()
			if err != nil {
				return fmt.Errorf("an error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s", dcId, d.Id(), bootVolume, err)
			}

			// Wait, catching any errors
			_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
			if errState != nil {
				return errState
			}
		}

		properties := ionoscloud.VolumeProperties{}

		if v, ok := d.GetOk("volume.0.name"); ok {
			vStr := v.(string)
			properties.Name = &vStr
		}

		if v, ok := d.GetOk("volume.0.size"); ok {
			vInt := float32(v.(int))
			properties.Size = &vInt
		}

		if v, ok := d.GetOk("volume.0.bus"); ok {
			vStr := v.(string)
			properties.Bus = &vStr
		}

		_, apiResponse, err := client.VolumesApi.DatacentersVolumesPatch(ctx, d.Get("datacenter_id").(string), bootVolume).Volume(properties).Execute()

		if err != nil {
			return fmt.Errorf("error patching volume (%s) (%s)", d.Id(), err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
		if errState != nil {
			return errState
		}
	}

	// Nic stuff
	if d.HasChange("nic") {
		nic := &ionoscloud.Nic{}
		for _, n := range *server.Entities.Nics.Items {
			nicStr := d.Get("primary_nic").(string)
			if *n.Id == nicStr {
				nic = &n
				break
			}
		}

		lan := int32(d.Get("nic.0.lan").(int))
		properties := ionoscloud.NicProperties{
			Lan: &lan,
		}

		if v, ok := d.GetOk("nic.0.name"); ok {
			vStr := v.(string)
			properties.Name = &vStr
		}

		if v, ok := d.GetOk("nic.0.ip"); ok {
			ips := strings.Split(v.(string), ",")
			if len(ips) > 0 {
				properties.Ips = &ips
			}
		}

		properties.Dhcp = boolAddr(d.Get("nic.0.dhcp").(bool))

		properties.FirewallActive = boolAddr(d.Get("nic.0.firewall_active").(bool))

		if d.HasChange("nic.0.firewall") {

			firewall := GetFirewallResource(d, "nic.0.firewall")
			nic.Entities = &ionoscloud.NicEntities{
				Firewallrules: &ionoscloud.FirewallRules{
					Items: &[]ionoscloud.FirewallRule{
						firewall,
					},
				},
			}
		}
		mProp, _ := json.Marshal(properties)
		log.Printf("[DEBUG] Updating props: %s", string(mProp))
		_, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsPatch(ctx, d.Get("datacenter_id").(string), *server.Id, *nic.Id).Nic(properties).Execute()
		if err != nil {
			return fmt.Errorf(
				"error updating nic (%s)", err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
		if errState != nil {
			return errState
		}

	}

	return resourceServerRead(d, meta)
}

func resourceServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)
	dcId := d.Get("datacenter_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	server, _, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, d.Id()).Execute()

	if err != nil {
		return fmt.Errorf("error occured while fetching a server ID %s %s", d.Id(), err)
	}

	if server.Properties.BootVolume != nil {
		apiResponse, err := client.VolumesApi.DatacentersVolumesDelete(ctx, dcId, *server.Properties.BootVolume.Id).Execute()

		if err != nil {
			return fmt.Errorf("error occured while delete volume %s of server ID %s %s", *server.Properties.BootVolume.Id, d.Id(), err)
		}
		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
		if errState != nil {
			return errState
		}
	}

	apiResponse, err := client.ServersApi.DatacentersServersDelete(ctx, dcId, d.Id()).Execute()
	if err != nil {
		return fmt.Errorf("an error occured while deleting a server ID %s %s", d.Id(), err)

	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
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
