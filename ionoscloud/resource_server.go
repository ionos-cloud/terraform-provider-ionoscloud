package ionoscloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
				Optional: true,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
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
							Optional: true,
							Computed: true,
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
							Computed:   true,
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

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	server, volume, err := initializeCreateRequests(d)

	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	var sshKeyPath []interface{}
	var publicKeys []string
	var image, imageAlias, imageInput string
	var isSnapshot bool
	var diags diag.Diagnostics
	var password, licenceType string

	dcId := d.Get("datacenter_id").(string)

	serverName := d.Get("name").(string)
	server.Properties.Name = &serverName

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
				diags := diag.FromErr(fmt.Errorf("error fetching sshkey from file (%s) %s", path, err.Error()))
				return diags
			}
			publicKeys = append(publicKeys, publicKey)
		}
		if len(publicKeys) > 0 {
			volume.SshKeys = &publicKeys
		}
	}

	if v, ok := d.GetOk("volume.0.image_name"); ok {
		imageInput = v.(string)
		if err := d.Set("image_name", v.(string)); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	} else if v, ok := d.GetOk("image_name"); ok {
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

		nic.Properties.Dhcp = boolAddr(d.Get("nic.0.dhcp").(bool))
		nic.Properties.FirewallActive = boolAddr(d.Get("nic.0.firewall_active").(bool))

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

		log.Printf("[DEBUG] dhcp nic before%t", *nic.Properties.Dhcp)

		server.Entities.Nics = &ionoscloud.Nics{
			Items: &[]ionoscloud.Nic{
				nic,
			},
		}

		log.Printf("[DEBUG] dhcp nic after %t", *nic.Properties.Dhcp)
		log.Printf("[DEBUG] dhcp %t", *(*server.Entities.Nics.Items)[0].Properties.Dhcp)

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

			(*server.Entities.Nics.Items)[0].Entities = &ionoscloud.NicEntities{
				Firewallrules: &ionoscloud.FirewallRules{
					Items: &[]ionoscloud.FirewallRule{
						firewall,
					},
				},
			}
		}
	}

	if (*server.Entities.Nics.Items)[0].Properties.Ips != nil {
		if len(*(*server.Entities.Nics.Items)[0].Properties.Ips) == 0 {
			*(*server.Entities.Nics.Items)[0].Properties.Ips = nil
		}
	}

	server, apiResponse, err := client.ServersApi.DatacentersServersPost(ctx, d.Get("datacenter_id").(string)).Server(server).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error creating server: %s", err))
		return diags
	}
	d.SetId(*server.Id)

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
	server, _, err = client.ServersApi.DatacentersServersFindById(ctx, d.Get("datacenter_id").(string), *server.Id).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error fetching server: (%s)", err))
		return diags
	}

	firewallRules, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, d.Get("datacenter_id").(string),
		*server.Id, *(*server.Entities.Nics.Items)[0].Id).Execute()

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

	if (*server.Entities.Nics.Items)[0].Id != nil {
		err := d.Set("primary_nic", *(*server.Entities.Nics.Items)[0].Id)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting primary nic %s: %s", d.Id(), err))
			return diags
		}
	}

	if (*server.Entities.Nics.Items)[0].Properties.Ips != nil &&
		len(*(*server.Entities.Nics.Items)[0].Properties.Ips) > 0 &&
		server.Entities.Volumes.Items != nil &&
		len(*server.Entities.Volumes.Items) > 0 &&
		(*server.Entities.Volumes.Items)[0].Properties != nil &&
		(*server.Entities.Volumes.Items)[0].Properties.ImagePassword != nil {

		d.SetConnInfo(map[string]string{
			"type":     "ssh",
			"host":     (*(*server.Entities.Nics.Items)[0].Properties.Ips)[0],
			"password": *(*server.Entities.Volumes.Items)[0].Properties.ImagePassword,
		})
	}

	return resourceServerRead(ctx, d, meta)
}

func GetFirewallResource(d *schema.ResourceData, path string, update bool) ionoscloud.FirewallRule {

	firewall := ionoscloud.FirewallRule{
		Properties: &ionoscloud.FirewallruleProperties{},
	}

	if !update {
		if v, ok := d.GetOk(path + ".protocol"); ok {
			vStr := v.(string)
			firewall.Properties.Protocol = &vStr
		}
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

	if v, ok := d.GetOk(path + ".type"); ok {
		val := v.(string)
		firewall.Properties.Type = &val
	}

	return firewall
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)
	serverId := d.Id()

	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, serverId).Execute()
	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a server ID %s %s", d.Id(), err))
		return diags
	}

	if server.Properties.TemplateUuid != nil {
		if err := d.Set("template_uuid", *server.Properties.TemplateUuid); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
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

	if server.Properties.Type != nil {
		if err := d.Set("type", *server.Properties.Type); err != nil {
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

	if server.Entities.Volumes != nil && server.Entities.Volumes.Items != nil && len(*server.Entities.Volumes.Items) > 0 &&
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

		nic, _, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, dcId, serverId, primarynic.(string)).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error occured while fetching nic %s for server ID %s %s", primarynic.(string), d.Id(), err))
			return diags
		}

		if len(*nic.Properties.Ips) > 0 {
			if err := d.Set("primary_ip", (*nic.Properties.Ips)[0]); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}

		network := map[string]interface{}{}

		setPropWithNilCheck(network, "dhcp", nic.Properties.Dhcp)
		setPropWithNilCheck(network, "firewall_active", nic.Properties.FirewallActive)
		setPropWithNilCheck(network, "firewall_type", nic.Properties.FirewallType)
		setPropWithNilCheck(network, "lan", nic.Properties.Lan)
		setPropWithNilCheck(network, "name", nic.Properties.Name)
		setPropWithNilCheck(network, "ips", nic.Properties.Ips)
		setPropWithNilCheck(network, "mac", nic.Properties.Mac)
		setPropWithNilCheck(network, "device_number", nic.Properties.DeviceNumber)
		setPropWithNilCheck(network, "pci_slot", nic.Properties.PciSlot)

		if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
			network["ips"] = *nic.Properties.Ips
		}

		if firewallId, ok := d.GetOk("firewallrule_id"); ok {
			firewall, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, dcId, serverId, primarynic.(string), firewallId.(string)).Execute()
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error occured while fetching firewallrule %s for server ID %s %s", firewallId.(string), serverId, err))
				return diags
			}

			fw := map[string]interface{}{}
			/*
				"protocol": *firewall.Properties.Protocol,
				"name":     *firewall.Properties.Name,
			*/
			setPropWithNilCheck(fw, "protocol", firewall.Properties.Protocol)
			setPropWithNilCheck(fw, "name", firewall.Properties.Name)
			setPropWithNilCheck(fw, "source_mac", firewall.Properties.SourceMac)
			setPropWithNilCheck(fw, "source_ip", firewall.Properties.SourceIp)
			setPropWithNilCheck(fw, "target_ip", firewall.Properties.TargetIp)
			setPropWithNilCheck(fw, "port_range_start", firewall.Properties.PortRangeStart)
			setPropWithNilCheck(fw, "port_range_end", firewall.Properties.PortRangeEnd)
			setPropWithNilCheck(fw, "icmp_type", firewall.Properties.IcmpType)
			setPropWithNilCheck(fw, "icmp_code", firewall.Properties.IcmpCode)
			setPropWithNilCheck(fw, "type", firewall.Properties.Type)

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
		volumeObj, _, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, serverId, *server.Properties.BootVolume.Id).Execute()

		if err == nil {
			volumeItem := map[string]interface{}{}

			setPropWithNilCheck(volumeItem, "name", volumeObj.Properties.Name)
			setPropWithNilCheck(volumeItem, "disk_type", volumeObj.Properties.Type)
			setPropWithNilCheck(volumeItem, "size", volumeObj.Properties.Size)
			setPropWithNilCheck(volumeItem, "licence_type", volumeObj.Properties.LicenceType)
			setPropWithNilCheck(volumeItem, "bus", volumeObj.Properties.Bus)
			setPropWithNilCheck(volumeItem, "availability_zone", volumeObj.Properties.AvailabilityZone)
			setPropWithNilCheck(volumeItem, "cpu_hot_plug", volumeObj.Properties.CpuHotPlug)
			setPropWithNilCheck(volumeItem, "ram_hot_plug", volumeObj.Properties.CpuHotPlug)
			setPropWithNilCheck(volumeItem, "nic_hot_plug", volumeObj.Properties.CpuHotPlug)
			setPropWithNilCheck(volumeItem, "nic_hot_unplug", volumeObj.Properties.CpuHotPlug)
			setPropWithNilCheck(volumeItem, "disc_virtio_hot_plug", volumeObj.Properties.CpuHotPlug)
			setPropWithNilCheck(volumeItem, "disc_virtio_hot_unplug", volumeObj.Properties.CpuHotPlug)
			setPropWithNilCheck(volumeItem, "device_number", volumeObj.Properties.DeviceNumber)
			setPropWithNilCheck(volumeItem, "pci_slot", volumeObj.Properties.PciSlot)

			userData := d.Get("volume.0.user_data")
			volumeItem["user_data"] = userData

			backupUnit := d.Get("volume.0.backup_unit_id")
			volumeItem["backup_unit_id"] = backupUnit

			volumesList := []map[string]interface{}{volumeItem}
			if err := d.Set("volume", volumesList); err != nil {
				diags := diag.FromErr(fmt.Errorf("[DEBUG] Error saving volume to state for IonosCloud server (%s): %s", d.Id(), err))
				return diags
			}
		}
	}

	bootVolume, ok := d.GetOk("boot_volume")
	if ok && len(bootVolume.(string)) > 0 {
		_, _, err = client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, d.Id(), bootVolume.(string)).Execute()
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
		diags := diag.FromErr(fmt.Errorf("availability_zone is immutable"))
		return diags
	}
	if d.HasChange("cpu_family") {
		_, n := d.GetChange("cpu_family")
		nStr := n.(string)
		request.CpuFamily = &nStr
	}

	if d.HasChange("image_name") {
		diags := diag.FromErr(fmt.Errorf("image_name is immutable"))
		return diags
	}

	if d.HasChange("boot_cdrom") {
		_, n := d.GetChange("boot_cdrom")
		bootCdrom := n.(string)

		if IsValidUUID(bootCdrom) {

			request.BootCdrom = &ionoscloud.ResourceReference{
				Id: &bootCdrom,
			}

		} else {
			diags := diag.FromErr(fmt.Errorf("boot_volume has to be a valid UUID, got: %s", bootCdrom))
			return diags
		}
		/* todo: figure out a way of sending a nil bootCdrom to the API (the sdk's omitempty doesn't let us) */
	}

	server, apiResponse, err := client.ServersApi.DatacentersServersPatch(ctx, dcId, d.Id()).Server(request).Execute()

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
		if d.HasChange("volume.0.user_data") {
			diags := diag.FromErr(fmt.Errorf("volume.0.user_data is immutable and is only allowed to be set on a new volume creation"))
			return diags
		}

		if d.HasChange("volume.0.backup_unit_id") {
			diags := diag.FromErr(fmt.Errorf("volume.0.backup_unit_id\" is immutable and is only allowed to be set on a new volume creation"))
			return diags
		}

		if d.HasChange("volume.0.image_name") {
			diags := diag.FromErr(fmt.Errorf("volume.0.image_name is immutable"))
			return diags
		}

		if d.HasChange("volume.0.disk_type") {
			diags := diag.FromErr(fmt.Errorf("volume.0.disk_type is immutable"))
			return diags
		}

		if d.HasChange("volume.0.availability_zone") {
			diags := diag.FromErr(fmt.Errorf("volume.0.availability_zone is immutable"))
			return diags
		}

		bootVolume := d.Get("boot_volume").(string)
		_, _, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, d.Id(), bootVolume).Execute()

		if err != nil {
			volume := ionoscloud.Volume{
				Id: &bootVolume,
			}
			_, apiResponse, err := client.ServersApi.DatacentersServersVolumesPost(ctx, dcId, d.Id()).Volume(volume).Execute()
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s", dcId, d.Id(), bootVolume, err))
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

		_, apiResponse, err := client.VolumesApi.DatacentersVolumesPatch(ctx, d.Get("datacenter_id").(string), bootVolume).Volume(properties).Execute()

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error patching volume (%s) (%s)", d.Id(), err))
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

		if v, ok := d.GetOk("nic.0.ips"); ok {
			raw := v.([]interface{})
			if raw != nil && len(raw) > 0 {
				ips := make([]string, 0)
				for _, rawIp := range raw {
					ip := rawIp.(string)
					ips = append(ips, ip)
				}
				if ips != nil && len(ips) > 0 {
					properties.Ips = &ips
				}
			}
		}

		properties.Dhcp = boolAddr(d.Get("nic.0.dhcp").(bool))

		properties.FirewallActive = boolAddr(d.Get("nic.0.firewall_active").(bool))

		if v, ok := d.GetOk("nic.0.firewall_type"); ok {
			vStr := v.(string)
			properties.FirewallType = &vStr
		}

		if d.HasChange("nic.0.firewall") {

			firewall := GetFirewallResource(d, "nic.0.firewall.0", true)

			firewallId := d.Get("firewallrule_id").(string)

			_, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, dcId, *server.Id, *nic.Id, firewallId).Execute()

			if err != nil {
				if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode != 404 {
					diags := diag.FromErr(fmt.Errorf("error occured at checking existance of firewall %s %s", firewallId, err))
					return diags
				} else if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
					diags := diag.FromErr(fmt.Errorf("firewall does not exist %s", firewallId))
					return diags
				}
			}

			firewall, apiResponse, err = client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPatch(ctx, dcId, *server.Id, *nic.Id, firewallId).Firewallrule(*firewall.Properties).Execute()
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occured while updating firewall rule dcId: %s server_id: %s nic_id %s ID: %s Response: %s", dcId, *server.Id, *nic.Id, firewallId, err))
				return diags
			}

			// Wait, catching any errors
			_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
			if errState != nil {
				diags := diag.FromErr(errState)
				return diags
			}

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
			diags := diag.FromErr(fmt.Errorf("error updating nic (%s)", err))
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

	server, _, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, d.Id()).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a server ID %s %s", d.Id(), err))
		return diags
	}

	if server.Properties.BootVolume != nil && strings.ToLower(*server.Properties.Type) != "cube" {
		apiResponse, err := client.VolumesApi.DatacentersVolumesDelete(ctx, dcId, *server.Properties.BootVolume.Id).Execute()

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error occured while delete volume %s of server ID %s %s", *server.Properties.BootVolume.Id, d.Id(), err))
			return diags
		}
		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(errState)
			return diags
		}
	}

	apiResponse, err := client.ServersApi.DatacentersServersDelete(ctx, dcId, d.Id()).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a server ID %s %s", d.Id(), err))
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

func resourceServerImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}", d.Id())
	}

	datacenterId := parts[0]
	serverId := parts[1]

	client := meta.(*ionoscloud.APIClient)

	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, datacenterId, serverId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find server %q", serverId)
		}
		return nil, fmt.Errorf("error occured while fetching a server ID %s %s", d.Id(), err)
	}

	d.SetId(*server.Id)

	if err := d.Set("datacenter_id", datacenterId); err != nil {
		return nil, err
	}

	if server.Properties.Name != nil {
		if err := d.Set("name", *server.Properties.Name); err != nil {
			return nil, err
		}
	}

	if server.Properties.Cores != nil {
		if err := d.Set("cores", *server.Properties.Cores); err != nil {
			return nil, err
		}
	}

	if server.Properties.Ram != nil {
		if err := d.Set("ram", *server.Properties.Ram); err != nil {
			return nil, err
		}
	}

	if server.Properties.AvailabilityZone != nil {
		if err := d.Set("availability_zone", *server.Properties.AvailabilityZone); err != nil {
			return nil, err
		}
	}

	if server.Properties.CpuFamily != nil {
		if err := d.Set("cpu_family", *server.Properties.CpuFamily); err != nil {
			return nil, err
		}
	}

	if server.Entities.Volumes != nil &&
		len(*server.Entities.Volumes.Items) > 0 &&
		(*server.Entities.Volumes.Items)[0].Properties.Image != nil {
		if err := d.Set("boot_image", *(*server.Entities.Volumes.Items)[0].Properties.Image); err != nil {
			return nil, err
		}
	}

	if server.Entities.Nics != nil && len(*server.Entities.Nics.Items) > 0 && (*server.Entities.Nics.Items)[0].Id != nil {
		primaryNic := *(*server.Entities.Nics.Items)[0].Id
		if err := d.Set("primary_nic", primaryNic); err != nil {
			return nil, err
		}

		nic, _, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, datacenterId, serverId, primaryNic).Execute()
		if err != nil {
			return nil, err
		}

		if len(*nic.Properties.Ips) > 0 {
			if err := d.Set("primary_ip", (*nic.Properties.Ips)[0]); err != nil {
				return nil, err
			}
		}

		network := map[string]interface{}{}

		setPropWithNilCheck(network, "dhcp", nic.Properties.Dhcp)
		setPropWithNilCheck(network, "firewall_active", nic.Properties.FirewallActive)

		setPropWithNilCheck(network, "lan", nic.Properties.Lan)
		setPropWithNilCheck(network, "name", nic.Properties.Name)
		setPropWithNilCheck(network, "ips", nic.Properties.Ips)
		setPropWithNilCheck(network, "mac", nic.Properties.Mac)

		if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
			network["ips"] = *nic.Properties.Ips
		}

		firewallRules, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, primaryNic).Execute()

		if err != nil {
			return nil, err
		}

		if firewallRules.Items != nil {
			if len(*firewallRules.Items) > 0 {
				if err := d.Set("firewallrule_id", *(*firewallRules.Items)[0].Id); err != nil {
					return nil, err
				}
			}
		}

		if firewallId, ok := d.GetOk("firewallrule_id"); ok {
			firewall, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, primaryNic, firewallId.(string)).Execute()
			if err != nil {
				return nil, err
			}

			fw := map[string]interface{}{}

			setPropWithNilCheck(fw, "protocol", firewall.Properties.Protocol)
			setPropWithNilCheck(fw, "name", firewall.Properties.Name)
			setPropWithNilCheck(fw, "source_mac", firewall.Properties.SourceMac)
			setPropWithNilCheck(fw, "source_ip", firewall.Properties.SourceIp)
			setPropWithNilCheck(fw, "target_ip", firewall.Properties.TargetIp)
			setPropWithNilCheck(fw, "port_range_start", firewall.Properties.PortRangeStart)
			setPropWithNilCheck(fw, "port_range_end", firewall.Properties.PortRangeEnd)
			setPropWithNilCheck(fw, "icmp_type", firewall.Properties.IcmpType)
			setPropWithNilCheck(fw, "icmp_code", firewall.Properties.IcmpCode)

			network["firewall"] = []map[string]interface{}{fw}
		}

		networks := []map[string]interface{}{network}
		if err := d.Set("nic", networks); err != nil {
			return nil, err
		}
	}

	if server.Properties.BootVolume != nil {
		if server.Properties.BootVolume.Id != nil {
			if err := d.Set("boot_volume", *server.Properties.BootVolume.Id); err != nil {
				return nil, err
			}
		}
		volumeObj, _, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, datacenterId, serverId, *server.Properties.BootVolume.Id).Execute()
		if err == nil {
			volumeItem := map[string]interface{}{}

			setPropWithNilCheck(volumeItem, "name", volumeObj.Properties.Name)
			setPropWithNilCheck(volumeItem, "disk_type", volumeObj.Properties.Type)
			setPropWithNilCheck(volumeItem, "size", volumeObj.Properties.Size)
			setPropWithNilCheck(volumeItem, "licence_type", volumeObj.Properties.LicenceType)
			setPropWithNilCheck(volumeItem, "bus", volumeObj.Properties.Bus)
			setPropWithNilCheck(volumeItem, "availability_zone", volumeObj.Properties.AvailabilityZone)

			volumesList := []map[string]interface{}{volumeItem}
			if err := d.Set("volume", volumesList); err != nil {
				return nil, err
			}
		}
	}
	if len(parts) > 3 {
		if err := d.Set("firewallrule_id", parts[3]); err != nil {
			return nil, err
		}
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
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

// Initializes server and volume with the required attributes depending on the server type (CUBE or ENTERPRISE)
func initializeCreateRequests(d *schema.ResourceData) (ionoscloud.Server, ionoscloud.VolumeProperties, error) {
	server := ionoscloud.Server{
		Properties: &ionoscloud.ServerProperties{},
	}
	volume := ionoscloud.VolumeProperties{}

	serverType := d.Get("type").(string)

	if serverType != "" {
		server.Properties.Type = &serverType
	}

	switch strings.ToLower(serverType) {
	case "cube":
		if v, ok := d.GetOk("template_uuid"); ok {
			vStr := v.(string)
			server.Properties.TemplateUuid = &vStr
		} else {
			return server, volume, fmt.Errorf("template_uuid argument is required for %s type of servers\n", serverType)
		}

		if _, ok := d.GetOk("cores"); ok {
			return server, volume, fmt.Errorf("cores argument can not be set for %s type of servers\n", serverType)
		}

		if _, ok := d.GetOk("ram"); ok {
			return server, volume, fmt.Errorf("ram argument can not be set for %s type of servers\n", serverType)
		}

		if _, ok := d.GetOk("volume.0.size"); ok {
			return server, volume, fmt.Errorf("volume.0.size argument can not be set for %s type of servers\n", serverType)
		}
	default:
		if _, ok := d.GetOk("template_uuid"); ok {
			return server, volume, fmt.Errorf("template_uuid argument can not be set only for %s type of servers\n", serverType)
		}

		if v, ok := d.GetOk("cores"); ok {
			vInt := int32(v.(int))
			server.Properties.Cores = &vInt
		} else {
			return server, volume, fmt.Errorf("cores argument is required for %s type of servers\n", serverType)
		}

		if v, ok := d.GetOk("ram"); ok {
			vInt := int32(v.(int))
			server.Properties.Ram = &vInt
		} else {
			return server, volume, fmt.Errorf("ram argument is required for %s type of servers\n", serverType)
		}

		if v, ok := d.GetOk("volume.0.size"); ok {
			vInt := float32(v.(int))
			volume.Size = &vInt
		} else {
			return server, volume, fmt.Errorf("volume.0.size argument is required for %s type of servers\n", serverType)
		}
	}

	return server, volume, nil
}
