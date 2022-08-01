package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"strconv"
)

func resourceCubeServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCubeServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceCubeServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServerImport,
		},
		CustomizeDiff: checkServerImmutableFields,

		Schema: map[string]*schema.Schema{
			// Cube Server parameters
			"template_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"availability_zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2"}, true)),
			},
			"boot_volume": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"boot_cdrom": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu_family": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: DiffToLower,
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
							Computed:   true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if k == "volume.0.ssh_key_path.#" {
									if d.Get("ssh_key_path.#") == new {
										return true
									}
								}

								sshKeyPath := d.Get("volume.0.ssh_key_path").([]interface{})
								oldSshKeyPath := d.Get("ssh_key_path").([]interface{})

								if len(utils.DiffSlice(convertSlice(sshKeyPath), convertSlice(oldSshKeyPath))) == 0 {
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
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.All(validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2", "ZONE_3"}, true)),
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
						"device_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backup_unit_id": {
							Type:        schema.TypeString,
							Description: "The uuid of the Backup Unit that user has access to. The property is immutable and is only allowed to be set on a new volume creation. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property.",
							Optional:    true,
							Computed:    true,
						},
						"user_data": {
							Type:        schema.TypeString,
							Description: "The cloud-init configuration for the volume as base64 encoded string. The property is immutable and is only allowed to be set on a new volume creation. It is mandatory to provide either 'public image' or 'imageAlias' that has cloud-init compatibility in conjunction with this property.",
							Optional:    true,
							Computed:    true,
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
						"ips": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
							Optional: true,
						},
						"firewall_active": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"firewall_type": {
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
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: DiffToLower,
										ValidateFunc:     validation.All(validation.StringIsNotWhiteSpace),
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
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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

func resourceCubeServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	server := ionoscloud.Server{
		Properties: &ionoscloud.ServerProperties{},
	}
	volume := ionoscloud.VolumeProperties{}

	var sshKeyPath []interface{}
	var publicKeys []string
	var image, imageAlias, imageInput string
	var isSnapshot bool
	var diags diag.Diagnostics
	var password, licenceType string

	dcId := d.Get("datacenter_id").(string)

	serverName := d.Get("name").(string)
	server.Properties.Name = &serverName

	templateUuid := d.Get("template_uuid").(string)
	server.Properties.TemplateUuid = &templateUuid

	if v, ok := d.GetOk("availability_zone"); ok {
		vStr := v.(string)
		server.Properties.AvailabilityZone = &vStr
	}

	if v, ok := d.GetOk("cpu_family"); ok {
		if v.(string) != "" {
			vStr := v.(string)
			server.Properties.CpuFamily = &vStr
		}
	}

	volumeType := d.Get("volume.0.disk_type").(string)
	volume.Type = &volumeType

	if v, ok := d.GetOk("volume.0.image_password"); ok {
		vStr := v.(string)
		volume.ImagePassword = &vStr
		if err := d.Set("image_password", password); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if v, ok := d.GetOk("image_password"); ok {
		password = v.(string)
		volume.ImagePassword = &password
	}

	if v, ok := d.GetOk("volume.0.licence_type"); ok {
		licenceType = v.(string)
		volume.LicenceType = &licenceType
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

	if v, ok := d.GetOk("volume.0.backup_unit_id"); ok {
		vStr := v.(string)
		volume.BackupunitId = &vStr
	}

	if v, ok := d.GetOk("volume.0.user_data"); ok {
		vStr := v.(string)
		volume.UserData = &vStr
	}

	if v, ok := d.GetOk("volume.0.ssh_key_path"); ok {
		sshKeyPath = v.([]interface{})
		if err := d.Set("ssh_key_path", v.([]interface{})); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	} else if v, ok := d.GetOk("ssh_key_path"); ok {
		sshKeyPath = v.([]interface{})
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

	if _, ok := d.GetOk("boot_cdrom"); ok {
		resId := d.Get("boot_cdrom").(string)
		server.Properties.BootCdrom = &ionoscloud.ResourceReference{
			Id: &resId,
		}
	}

	if _, ok := d.GetOk("boot_volume"); ok {
		diags := diag.FromErr(fmt.Errorf("boot_volume argument can be set only in update requests \n"))
		return diags
	}

	if len(sshKeyPath) != 0 {
		for _, path := range sshKeyPath {
			log.Printf("[DEBUG] Reading file %s", path)
			publicKey, err := readPublicKey(path.(string))
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error fetching sshkey from file (%s) %w", path, err))
				return diags
			}
			publicKeys = append(publicKeys, publicKey)
		}
		if len(publicKeys) > 0 {
			volume.SshKeys = &publicKeys
		}
	}

	if v, ok := d.GetOk("image_name"); ok {
		imageInput = v.(string)
	}

	if imageInput != "" {
		image, imageAlias, isSnapshot, diags = checkImage(ctx, client, imageInput, password, licenceType, dcId, sshKeyPath)
		if diags != nil {
			return diags
		}
	}

	if isSnapshot == true && (volume.ImagePassword != nil && *volume.ImagePassword != "" || len(publicKeys) > 0) {
		diags := diag.FromErr(fmt.Errorf("you can't pass 'image_password' and/or 'ssh keys' when creating a volume from a snapshot"))
		return diags
	}

	if image != "" {
		volume.Image = &image
	} else {
		volume.Image = nil
	}

	if imageAlias != "" {
		volume.ImageAlias = &imageAlias
	} else {
		volume.ImageAlias = nil
	}

	server.Entities = &ionoscloud.ServerEntities{
		Volumes: &ionoscloud.AttachedVolumes{
			Items: &[]ionoscloud.Volume{
				{
					Properties: &volume,
				},
			},
		},
	}
	var primaryNic *ionoscloud.Nic
	if server.Entities.Nics != nil && server.Entities.Nics.Items != nil && len(*server.Entities.Nics.Items) > 0 {
		primaryNic = &(*server.Entities.Nics.Items)[0]
	}
	// Nic Arguments
	if _, ok := d.GetOk("nic"); ok {
		lanInt := int32(d.Get("nic.0.lan").(int))
		nic := ionoscloud.Nic{Properties: &ionoscloud.NicProperties{
			Lan: &lanInt,
		}}

		if v, ok := d.GetOk("nic.0.name"); ok {
			vStr := v.(string)
			nic.Properties.Name = &vStr
		}

		dhcp := d.Get("nic.0.dhcp").(bool)
		fwActive := d.Get("nic.0.firewall_active").(bool)
		nic.Properties.Dhcp = &dhcp
		nic.Properties.FirewallActive = &fwActive

		if v, ok := d.GetOk("nic.0.firewall_type"); ok {
			v := v.(string)
			nic.Properties.FirewallType = &v
		}

		if v, ok := d.GetOk("nic.0.ips"); ok {
			raw := v.([]interface{})
			if raw != nil && len(raw) > 0 {
				var ips []string
				for _, rawIp := range raw {
					ip := rawIp.(string)
					ips = append(ips, ip)
				}
				if ips != nil && len(ips) > 0 {
					nic.Properties.Ips = &ips
				}
			}
		}

		log.Printf("[DEBUG] dhcp nic before %t", *nic.Properties.Dhcp)

		server.Entities.Nics = &ionoscloud.Nics{
			Items: &[]ionoscloud.Nic{
				nic,
			},
		}
		primaryNic = &(*server.Entities.Nics.Items)[0]
		log.Printf("[DEBUG] dhcp nic after %t", *nic.Properties.Dhcp)
		log.Printf("[DEBUG] dhcp %t", *primaryNic.Properties.Dhcp)

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

			if v, ok := d.GetOk("nic.0.firewall.0.type"); ok {
				val := v.(string)
				firewall.Properties.Type = &val
			}

			primaryNic.Entities = &ionoscloud.NicEntities{
				Firewallrules: &ionoscloud.FirewallRules{
					Items: &[]ionoscloud.FirewallRule{
						firewall,
					},
				},
			}
		}
	}
	if primaryNic != nil && primaryNic.Properties != nil && primaryNic.Properties.Ips != nil {
		if len(*primaryNic.Properties.Ips) == 0 {
			*primaryNic.Properties.Ips = nil
		}
	}

	createdServer, apiResponse, err := client.ServersApi.DatacentersServersPost(ctx, d.Get("datacenter_id").(string)).Server(server).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error creating server: %w", err))
		return diags
	}
	d.SetId(*createdServer.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			log.Printf("[DEBUG] failed to create createdServer resource")
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(fmt.Errorf("error waiting for state change for server creation %w", errState))
		return diags
	}

	// get additional data for schema
	createdServer, apiResponse, err = client.ServersApi.DatacentersServersFindById(ctx, d.Get("datacenter_id").(string), *createdServer.Id).Depth(3).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error fetching server: (%w)", err))
		return diags
	}

	firewallRules, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, d.Get("datacenter_id").(string),
		*createdServer.Id, *(*createdServer.Entities.Nics.Items)[0].Id).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching firewall rules: %s", err))
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

	if (*createdServer.Entities.Nics.Items)[0].Id != nil {
		err := d.Set("primary_nic", *(*createdServer.Entities.Nics.Items)[0].Id)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting primary nic %s: %s", d.Id(), err))
			return diags
		}
	}

	if (*createdServer.Entities.Nics.Items)[0].Properties.Ips != nil &&
		len(*(*createdServer.Entities.Nics.Items)[0].Properties.Ips) > 0 &&
		createdServer.Entities.Volumes.Items != nil &&
		len(*createdServer.Entities.Volumes.Items) > 0 &&
		(*createdServer.Entities.Volumes.Items)[0].Properties != nil &&
		(*createdServer.Entities.Volumes.Items)[0].Properties.ImagePassword != nil {

		d.SetConnInfo(map[string]string{
			"type":     "ssh",
			"host":     (*(*createdServer.Entities.Nics.Items)[0].Properties.Ips)[0],
			"password": *(*createdServer.Entities.Volumes.Items)[0].Properties.ImagePassword,
		})
	}
	return resourceServerRead(ctx, d, meta)
}

func resourceCubeServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient
	dcId := d.Get("datacenter_id").(string)

	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a server ID %s %s", d.Id(), err))
		return diags
	}

	if server.Properties.BootVolume != nil {
		apiResponse, err := client.VolumesApi.DatacentersVolumesDelete(ctx, dcId, *server.Properties.BootVolume.Id).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error occured while delete volume %s of server ID %s %w", *server.Properties.BootVolume.Id, d.Id(), err))
			return diags
		}
		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(fmt.Errorf("error getting state change for volumes delete %w", errState))
			return diags
		}
	}

	apiResponse, err = client.ServersApi.DatacentersServersDelete(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a server ID %s %w", d.Id(), err))
		return diags

	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(fmt.Errorf("error getting state change for datacenter delete %w", errState))
		return diags
	}

	d.SetId("")
	return nil
}

// nu are cores, ram si volume.size cube sererul
