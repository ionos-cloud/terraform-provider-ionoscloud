package ionoscloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/ssh"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServerImport,
		},
		Schema: map[string]*schema.Schema{
			// Server parameters
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

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	var imageAlias string
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
	dcId := d.Get("datacenter_id").(string)

	isSnapshot := false

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
			diags := diag.FromErr(err)
			return diags
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

	var sshkeyPath []interface{}

	if v, ok := d.GetOk("volume.0.ssh_key_path"); ok {
		sshkeyPath = v.([]interface{})
		if err := d.Set("ssh_key_path", v.([]interface{})); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	} else if v, ok := d.GetOk("ssh_key_path"); ok {
		sshkeyPath = v.([]interface{})
		if err := d.Set("ssh_key_path", v.([]interface{})); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	} else {
		if err := d.Set("ssh_key_path", [][]string{}); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	var image, imageName string
	if v, ok := d.GetOk("volume.0.image_name"); ok {
		imageName = v.(string)
		if err := d.Set("image_name", v.(string)); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	} else if v, ok := d.GetOk("image_name"); ok {
		imageName = v.(string)
	} else {
		diags := diag.FromErr(fmt.Errorf("either 'image_name' or 'volume.0.image_name' must be provided"))
		return diags
	}
	if !IsValidUUID(imageName) {
		img, err := getImage(ctx, client, dcId, imageName, *volume.Type)
		if err != nil {
			diags := diag.FromErr(err)
			return diags
		}
		if img != nil {
			image = *img.Id
		}
		// if no image id was found with that name we look for a matching snapshot
		if image == "" {
			image = getSnapshotId(ctx, client, imageName)
			if image != "" {
				isSnapshot = true
			} else {
				dc, _, err := client.DataCenterApi.DatacentersFindById(ctx, dcId).Execute()
				if err != nil {
					diags := diag.FromErr(fmt.Errorf("error fetching datacenter %s: (%s)", dcId, err))
					return diags
				}
				imageAlias = getImageAlias(ctx, client, imageName, *dc.Properties.Location)
			}
		}
		if image == "" && imageAlias == "" {
			diags := diag.FromErr(fmt.Errorf("could not find an image/imagealias/snapshot that matches %s ", imageName))
			return diags
		}
		if volume.ImagePassword == nil && len(sshkeyPath) == 0 && isSnapshot == false && img.Properties.Public != nil && *img.Properties.Public {
			diags := diag.FromErr(fmt.Errorf("either 'image_password' or 'ssh_key_path' must be provided"))
			return diags
		}
	} else {
		img, apiResponse, err := client.ImageApi.ImagesFindById(ctx, imageName).Execute()

		if apiResponse != nil && apiResponse.StatusCode == 404 {

			_, apiResponse, err = client.SnapshotApi.SnapshotsFindById(ctx, imageName).Execute()

			if err != nil {
				diags := diag.FromErr(fmt.Errorf("could not fetch image/snapshot: %s", err))
				return diags
			}

			isSnapshot = true

		} else if err != nil {
			payload := ""
			if apiResponse != nil {
				payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
			}
			diags := diag.FromErr(fmt.Errorf("error fetching image/snapshot: %s %s", err, payload))
			return diags
		}

		if *img.Properties.Public == true && isSnapshot == false {

			if volume.ImagePassword == nil && len(sshkeyPath) == 0 {
				diags := diag.FromErr(fmt.Errorf("either 'image_password' or 'ssh_key_path' must be provided"))
				return diags
			}

			img, err := getImage(ctx, client, d.Get("datacenter_id").(string), imageName, *volume.Type)

			if err != nil {
				diags := diag.FromErr(err)
				return diags
			}

			if img != nil {
				image = *img.Id
			}
		} else {
			img, _, err := client.ImageApi.ImagesFindById(ctx, imageName).Execute()
			if err != nil {
				_, _, err := client.SnapshotApi.SnapshotsFindById(ctx, imageName).Execute()
				if err != nil {
					diags := diag.FromErr(fmt.Errorf("error fetching image/snapshot: %s", err))
					return diags
				}
				isSnapshot = true
			}

			if img.Properties.Public != nil && *img.Properties.Public == true && isSnapshot == false {
				if volume.ImagePassword == nil && len(sshkeyPath) == 0 {
					diags := diag.FromErr(fmt.Errorf("either 'image_password' or 'ssh_key_path' must be provided"))
					return diags
				}
				image = imageName
			} else {
				image = imageName
			}
		}
	}

	if len(sshkeyPath) != 0 {
		var publicKeys []string
		for _, path := range sshkeyPath {
			log.Printf("[DEBUG] Reading file %s", path)
			publicKey, err := readPublicKey(path.(string))
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error fetching sshkey from file (%s) %s", path, err.Error()))
				return diags
			}
			publicKeys = append(publicKeys, publicKey)
		}
		if len(publicKeys) > 0 {
			volume.SshKeys = &publicKeys
		}
	}

	if image == "" && volume.LicenceType == nil && imageAlias == "" && !isSnapshot {
		diags := diag.FromErr(fmt.Errorf("either 'image', 'licenceType', or 'imageAlias' must be set"))
		return diags
	}

	if isSnapshot == true && (volume.ImagePassword != nil || len(sshkeyPath) > 0) {
		diags := diag.FromErr(fmt.Errorf("passwords/SSH keys are not supported for snapshots"))
		return diags
	}

	if imageAlias != "" {
		volume.ImageAlias = &imageAlias
	} else {
		volume.ImageAlias = nil
	}

	if image != "" {
		volume.Image = &image
	} else {
		volume.Image = nil
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
		nic.Properties.Nat = boolAddr(d.Get("nic.0.nat").(bool))

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

	server, apiResponse, err := client.ServerApi.DatacentersServersPost(ctx, d.Get("datacenter_id").(string)).Server(request).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error creating server: (%s)", err))
		return diags
	}

	if server.Id != nil {
		d.SetId(*server.Id)
	}

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
	server, _, err = client.ServerApi.DatacentersServersFindById(ctx, d.Get("datacenter_id").(string), *server.Id).Execute()

	if err != nil {
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("error fetching server: (%s) %s", err, payload))
		return diags
	}

	firewallRules, _, err := client.NicApi.DatacentersServersNicsFirewallrulesGet(ctx, d.Get("datacenter_id").(string),
		*server.Id, *(*server.Entities.Nics.Items)[0].Id).Execute()

	if err != nil {
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching firewall rules: %s %s", err, payload))
		return diags
	}

	if firewallRules.Items != nil {
		if len(*firewallRules.Items) > 0 {
			if err := d.Set("firewallrule_id", *(*firewallRules.Items)[0].Id); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
	}

	if (*server.Entities.Nics.Items)[0].Id != nil {
		err := d.Set("primary_nic", *(*server.Entities.Nics.Items)[0].Id)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting primary nic %s: %s", d.Id(), err))
			return diags
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
	return resourceServerRead(ctx, d, meta)
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

func resourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)
	serverId := d.Id()

	server, apiResponse, err := client.ServerApi.DatacentersServersFindById(ctx, dcId, serverId).Execute()
	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a server ID %s %s %s", d.Id(), err, payload))
		return diags
	}

	if server.Properties.Name != nil {
		if err := d.Set("name", *server.Properties.Name); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if server.Properties.Cores != nil {
		if err := d.Set("cores", *server.Properties.Cores); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if server.Properties.Ram != nil {
		if err := d.Set("ram", *server.Properties.Ram); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if server.Properties.AvailabilityZone != nil {
		if err := d.Set("availability_zone", *server.Properties.AvailabilityZone); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if server.Properties.CpuFamily != nil {
		if err := d.Set("cpu_family", *server.Properties.CpuFamily); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if server.Entities.Volumes != nil &&
		len(*server.Entities.Volumes.Items) > 0 &&
		(*server.Entities.Volumes.Items)[0].Properties.Image != nil {
		if err := d.Set("boot_image", *(*server.Entities.Volumes.Items)[0].Properties.Image); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if primarynic, ok := d.GetOk("primary_nic"); ok {
		if err := d.Set("primary_nic", primarynic.(string)); err != nil {
			diags := diag.FromErr(err)
			return diags
		}

		nic, apiResponse, err := client.NicApi.DatacentersServersNicsFindById(ctx, dcId, serverId, primarynic.(string)).Execute()
		if err != nil {
			payload := ""
			if apiResponse != nil {
				payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
			}
			diags := diag.FromErr(fmt.Errorf("error occured while fetching nic %s for server ID %s %s %s", primarynic.(string), d.Id(), err, payload))
			return diags
		}

		if len(*nic.Properties.Ips) > 0 {
			if err := d.Set("primary_ip", (*nic.Properties.Ips)[0]); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}

		network := map[string]interface{}{
			"dhcp":            *nic.Properties.Dhcp,
			"nat":             *nic.Properties.Nat,
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
			firewall, _, err := client.NicApi.DatacentersServersNicsFirewallrulesFindById(ctx, dcId, serverId, primarynic.(string), firewallId.(string)).Execute()
			if err != nil {
				payload := ""
				if apiResponse != nil {
					payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
				}
				diags := diag.FromErr(fmt.Errorf("error occured while fetching firewallrule %s for server ID %s %s %s", firewallId.(string), serverId, err, payload))
				return diags
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
			diags := diag.FromErr(fmt.Errorf("[ERROR] unable saving nic to state IonosCloud Server (%s): %s", serverId, err))
			return diags
		}
	}

	if server.Properties.BootVolume != nil {
		if server.Properties.BootVolume.Id != nil {
			if err := d.Set("boot_volume", *server.Properties.BootVolume.Id); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
		volumeObj, _, err := client.ServerApi.DatacentersServersVolumesFindById(ctx, dcId, serverId, *server.Properties.BootVolume.Id).Execute()
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
				diags := diag.FromErr(fmt.Errorf("[DEBUG] Error saving volume to state for IonosCloud server (%s): %s", d.Id(), err))
				return diags
			}
		}
	}

	bootVolume, ok := d.GetOk("boot_volume")
	if ok && len(bootVolume.(string)) > 0 {
		_, _, err = client.ServerApi.DatacentersServersVolumesFindById(ctx, dcId, d.Id(), bootVolume.(string)).Execute()
		if err != nil {
			if err := d.Set("volume", nil); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
	}

	if server.Properties.BootCdrom != nil {
		if err := d.Set("boot_cdrom", *server.Properties.BootCdrom.Id); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}
	return nil
}

func boolAddr(b bool) *bool {
	return &b
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)
	request := ionoscloud.ServerProperties{}

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

	server, apiResponse, err := client.ServerApi.DatacentersServersPatch(ctx, dcId, d.Id()).Server(request).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while updating server ID %s: %s", d.Id(), err))
		return diags
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}
	// Volume stuff
	if d.HasChange("volume") {
		bootVolume := d.Get("boot_volume").(string)
		_, _, err := client.ServerApi.DatacentersServersVolumesFindById(ctx, dcId, d.Id(), bootVolume).Execute()

		if err != nil {
			volume := ionoscloud.Volume{
				Id: &bootVolume,
			}
			_, apiResponse, err := client.ServerApi.DatacentersServersVolumesPost(ctx, dcId, d.Id()).Volume(volume).Execute()
			if err != nil {
				payload := ""
				if apiResponse != nil {
					payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
				}
				diags := diag.FromErr(fmt.Errorf("an error occured while attaching a volume dcId: %s server_id: %s ID: %s %s %s", dcId, d.Id(), bootVolume, err, payload))
				return diags
			}

			// Wait, catching any errors
			_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
			if errState != nil {
				diags := diag.FromErr(errState)
				return diags
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

		_, apiResponse, err := client.VolumeApi.DatacentersVolumesPatch(ctx, d.Get("datacenter_id").(string), bootVolume).Volume(properties).Execute()

		if err != nil {
			payload := ""
			if apiResponse != nil {
				payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
			}
			diags := diag.FromErr(fmt.Errorf("error patching volume (%s) (%s) %s", d.Id(), err, payload))
			return diags
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(errState)
			return diags
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

		properties.Nat = boolAddr(d.Get("nic.0.nat").(bool))
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
		_, apiResponse, err := client.NicApi.DatacentersServersNicsPatch(ctx, d.Get("datacenter_id").(string), *server.Id, *nic.Id).Nic(properties).Execute()
		if err != nil {
			payload := ""
			if apiResponse != nil {
				payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
			}
			diags := diag.FromErr(fmt.Errorf("error updating nic (%s) %s", err, payload))
			return diags
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(errState)
			return diags
		}

	}

	return resourceServerRead(ctx, d, meta)
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)
	dcId := d.Get("datacenter_id").(string)

	server, apiResponse, err := client.ServerApi.DatacentersServersFindById(ctx, dcId, d.Id()).Execute()

	if err != nil {
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a server ID %s %s %s", d.Id(), err, payload))
		return diags
	}

	if server.Properties.BootVolume != nil {
		_, apiResponse, err := client.VolumeApi.DatacentersVolumesDelete(ctx, dcId, *server.Properties.BootVolume.Id).Execute()

		if err != nil {
			payload := ""
			if apiResponse != nil {
				payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
			}
			diags := diag.FromErr(fmt.Errorf("error occured while delete volume %s of server ID %s %s %s", *server.Properties.BootVolume.Id, d.Id(), err, payload))
			return diags
		}
		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(errState)
			return diags
		}
	}

	_, apiResponse, err = client.ServerApi.DatacentersServersDelete(ctx, dcId, d.Id()).Execute()
	if err != nil {
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a server ID %s %s %s", d.Id(), err, payload))
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
