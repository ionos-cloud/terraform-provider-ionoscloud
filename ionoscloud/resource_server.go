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
				Optional: true,
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
			"vm_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cdrom": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"cpu_hot_plug": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cpu_hot_unplug": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ram_hot_plug": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ram_hot_unplug": {
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
						"disc_scsi_hot_plug": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disc_scsi_hot_unplug": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"licence_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
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
						"image": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_alias": {
							Type:     schema.TypeString,
							Computed: true,
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
						"image_aliases": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
										Type:     schema.TypeInt,
										Optional: true,
									},
									"icmp_code": {
										Type:     schema.TypeInt,
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

	nic := ionoscloud.Nic{
		Properties: &ionoscloud.NicProperties{},
	}
	firewall := ionoscloud.FirewallRule{
		Properties: &ionoscloud.FirewallruleProperties{},
	}

	// create server object
	serverRequest, err := getServerData(d, false)

	// create volume object with data to be used for image
	volume, err := getVolumeData(d, "volume.0.", false)

	if err != nil {
		return diag.FromErr(err)
	}

	// get image and imageAlias
	image, imageAlias, err := getImage(ctx, client, d, *volume)

	// add remaining properties in volume (dependent in image and imageAlias)
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

	if backupUnitId, ok := d.GetOk("volume.0.backup_unit_id"); ok {
		if IsValidUUID(backupUnitId.(string)) {
			if image == "" && imageAlias == "" {
				diags := diag.FromErr(fmt.Errorf("it is mandatory to provide either public image or imageAlias in conjunction with backup unit id property"))
				return diags
			} else {
				backupUnitId := backupUnitId.(string)
				volume.BackupunitId = &backupUnitId
			}
		}
	}

	if userData, ok := d.GetOk("volume.0.user_data"); ok {
		if image == "" && imageAlias == "" {
			diags := diag.FromErr(fmt.Errorf("it is mandatory to provide either public image or imageAlias that has cloud-init compatibility in conjunction with backup unit id property "))
			return diags
		} else {
			userData := userData.(string)
			volume.UserData = &userData
		}
	}

	// add volume object to server
	serverRequest.Entities = &ionoscloud.ServerEntities{
		Volumes: &ionoscloud.AttachedVolumes{
			Items: &[]ionoscloud.Volume{
				{
					Properties: volume,
				},
			},
		},
	}

	// get nic data and add object to server
	if _, ok := d.GetOk("nic"); ok {
		nic = getNicData(d, "nic.0.")
	}

	serverRequest.Entities.Nics = &ionoscloud.Nics{
		Items: &[]ionoscloud.Nic{
			nic,
		},
	}

	// get firewall data and add object to server
	if _, ok := d.GetOk("nic.0.firewall"); ok {
		firewall = getFirewallData(d, "nic.0.firewall.0.", false)
		(*serverRequest.Entities.Nics.Items)[0].Entities = &ionoscloud.NicEntities{
			Firewallrules: &ionoscloud.FirewallRules{
				Items: &[]ionoscloud.FirewallRule{
					firewall,
				},
			},
		}
	}

	if (*serverRequest.Entities.Nics.Items)[0].Properties.Ips != nil {
		if len(*(*serverRequest.Entities.Nics.Items)[0].Properties.Ips) == 0 {
			*(*serverRequest.Entities.Nics.Items)[0].Properties.Ips = nil
		}
	}

	server, apiResponse, err := client.ServerApi.DatacentersServersPost(ctx, d.Get("datacenter_id").(string)).Server(*serverRequest).Execute()

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

	// get additional data for schema
	server, _, err = client.ServerApi.DatacentersServersFindById(ctx, d.Get("datacenter_id").(string), *server.Id).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error fetching server: %s", err))
		return diags
	}

	firewallRules, _, err := client.NicApi.DatacentersServersNicsFirewallrulesGet(ctx, d.Get("datacenter_id").(string),
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
		serverRequest.Entities.Volumes.Items != nil &&
		len(*serverRequest.Entities.Volumes.Items) > 0 &&
		(*serverRequest.Entities.Volumes.Items)[0].Properties != nil &&
		(*serverRequest.Entities.Volumes.Items)[0].Properties.ImagePassword != nil {

		d.SetConnInfo(map[string]string{
			"type":     "ssh",
			"host":     (*(*server.Entities.Nics.Items)[0].Properties.Ips)[0],
			"password": *(*serverRequest.Entities.Volumes.Items)[0].Properties.ImagePassword,
		})
	}

	return resourceServerRead(ctx, d, meta)
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)
	serverId := d.Id()

	server, apiResponse, err := client.ServerApi.DatacentersServersFindById(ctx, dcId, serverId).Execute()
	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a server ID %s %s", d.Id(), err))
		return diags
	}

	if err := setServerData(ctx, client, d, &server, true); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)

	serverRequest, err := getServerData(d, true)

	if err != nil {
		return diag.FromErr(err)
	}
	server, apiResponse, err := client.ServerApi.DatacentersServersPatch(ctx, dcId, d.Id()).Server(*serverRequest.Properties).Execute()

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
				diags := diag.FromErr(fmt.Errorf("an error occured while attaching a volume dcId: %s server_id: %s ID: %s %s", dcId, d.Id(), bootVolume, err))
				return diags
			}

			// Wait, catching any errors
			_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
			if errState != nil {
				diags := diag.FromErr(errState)
				return diags
			}
		}

		volume, err := getVolumeData(d, "volume.0.", true)

		if err != nil {
			return diag.FromErr(err)
		}
		_, apiResponse, err := client.VolumeApi.DatacentersVolumesPatch(ctx, d.Get("datacenter_id").(string), bootVolume).Volume(*volume).Execute()

		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error patching volume %s %s", d.Id(), err))
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
		var nicId string
		for _, n := range *server.Entities.Nics.Items {
			nicStr := d.Get("primary_nic").(string)
			if *n.Id == nicStr {
				nicId = *n.Id
				break
			}
		}

		nic := getNicData(d, "nic.0.")
		if d.HasChange("nic.0.firewall") {

			firewall := getFirewallData(d, "nic.0.firewall.0.", true)

			firewallId := d.Get("firewallrule_id").(string)

			_, _, err := client.NicApi.DatacentersServersNicsFirewallrulesFindById(ctx, dcId, *server.Id, nicId, firewallId).Execute()

			if err != nil {
				if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode != 404 {
					diags := diag.FromErr(fmt.Errorf("error occured at checking existance of firewall %s %s", firewallId, err))
					return diags
				} else if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
					diags := diag.FromErr(fmt.Errorf("firewall does not exist %s", firewallId))
					return diags
				}
			}

			firewall, apiResponse, err = client.NicApi.DatacentersServersNicsFirewallrulesPatch(ctx, dcId, *server.Id, nicId, firewallId).Firewallrule(*firewall.Properties).Execute()
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
		mProp, _ := json.Marshal(nic.Properties)
		log.Printf("[DEBUG] Updating props: %s", string(mProp))
		_, apiResponse, err := client.NicApi.DatacentersServersNicsPatch(ctx, d.Get("datacenter_id").(string), *server.Id, nicId).Nic(*nic.Properties).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error updating nic %s", err))
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
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a server ID %s %s", d.Id(), err))
		return diags
	}

	if server.Properties.BootVolume != nil {
		_, apiResponse, err := client.VolumeApi.DatacentersVolumesDelete(ctx, dcId, *server.Properties.BootVolume.Id).Execute()

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

	_, apiResponse, err = client.ServerApi.DatacentersServersDelete(ctx, dcId, d.Id()).Execute()
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

	server, apiResponse, err := client.ServerApi.DatacentersServersFindById(ctx, datacenterId, serverId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find server %q", serverId)
		}
		return nil, fmt.Errorf("error occured while fetching a server ID %s %s", d.Id(), err)
	}

	if err := d.Set("datacenter_id", datacenterId); err != nil {
		return nil, err
	}

	if err := setServerData(ctx, client, d, &server, false); err != nil {
		return nil, err
	}
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

func SetCdromProperties(image ionoscloud.Image) map[string]interface{} {

	cdrom := make(map[string]interface{})

	setPropWithNilCheck(cdrom, "name", image.Properties.Name)
	setPropWithNilCheck(cdrom, "description", image.Properties.Description)
	setPropWithNilCheck(cdrom, "location", image.Properties.Location)
	setPropWithNilCheck(cdrom, "size", image.Properties.Size)
	setPropWithNilCheck(cdrom, "cpu_hot_plug", image.Properties.CpuHotPlug)
	setPropWithNilCheck(cdrom, "cpu_hot_unplug", image.Properties.CpuHotUnplug)
	setPropWithNilCheck(cdrom, "ram_hot_plug", image.Properties.RamHotPlug)
	setPropWithNilCheck(cdrom, "ram_hot_unplug", image.Properties.RamHotUnplug)
	setPropWithNilCheck(cdrom, "nic_hot_plug", image.Properties.NicHotPlug)
	setPropWithNilCheck(cdrom, "nic_hot_unplug", image.Properties.NicHotUnplug)
	setPropWithNilCheck(cdrom, "disc_virtio_hot_plug", image.Properties.DiscVirtioHotPlug)
	setPropWithNilCheck(cdrom, "disc_virtio_hot_unplug", image.Properties.DiscVirtioHotUnplug)
	setPropWithNilCheck(cdrom, "disc_scsi_hot_plug", image.Properties.DiscScsiHotPlug)
	setPropWithNilCheck(cdrom, "disc_scsi_hot_unplug", image.Properties.DiscScsiHotUnplug)
	setPropWithNilCheck(cdrom, "licence_type", image.Properties.LicenceType)
	setPropWithNilCheck(cdrom, "image_type", image.Properties.ImageType)
	setPropWithNilCheck(cdrom, "public", image.Properties.Public)

	return cdrom
}

func SetFirewallProperties(firewall ionoscloud.FirewallRule) map[string]interface{} {

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
	return fw
}

func SetVolumeProperties(volume ionoscloud.Volume) map[string]interface{} {

	volumeMap := map[string]interface{}{}

	setPropWithNilCheck(volumeMap, "name", volume.Properties.Name)
	setPropWithNilCheck(volumeMap, "disk_type", volume.Properties.Type)
	setPropWithNilCheck(volumeMap, "size", volume.Properties.Size)
	setPropWithNilCheck(volumeMap, "licence_type", volume.Properties.LicenceType)
	setPropWithNilCheck(volumeMap, "image", volume.Properties.Image)
	setPropWithNilCheck(volumeMap, "image_alias", volume.Properties.ImageAlias)
	setPropWithNilCheck(volumeMap, "bus", volume.Properties.Bus)
	setPropWithNilCheck(volumeMap, "availability_zone", volume.Properties.AvailabilityZone)
	setPropWithNilCheck(volumeMap, "cpu_hot_plug", volume.Properties.CpuHotPlug)
	setPropWithNilCheck(volumeMap, "ram_hot_plug", volume.Properties.RamHotPlug)
	setPropWithNilCheck(volumeMap, "nic_hot_plug", volume.Properties.NicHotPlug)
	setPropWithNilCheck(volumeMap, "nic_hot_unplug", volume.Properties.NicHotUnplug)
	setPropWithNilCheck(volumeMap, "disc_virtio_hot_plug", volume.Properties.DiscVirtioHotPlug)
	setPropWithNilCheck(volumeMap, "disc_virtio_hot_unplug", volume.Properties.DiscVirtioHotUnplug)
	setPropWithNilCheck(volumeMap, "device_number", volume.Properties.DeviceNumber)
	if volume.Properties.SshKeys != nil && len(*volume.Properties.SshKeys) > 0 {
		var sshKeys []interface{}
		for _, sshKey := range *volume.Properties.SshKeys {
			sshKeys = append(sshKeys, sshKey)
		}
		volumeMap["ssh_keys"] = sshKeys
	}

	return volumeMap
}

func SetNetworkProperties(nic ionoscloud.Nic) map[string]interface{} {

	network := map[string]interface{}{}

	setPropWithNilCheck(network, "dhcp", nic.Properties.Dhcp)
	setPropWithNilCheck(network, "nat", nic.Properties.Nat)
	setPropWithNilCheck(network, "firewall_active", nic.Properties.FirewallActive)
	setPropWithNilCheck(network, "lan", nic.Properties.Lan)
	setPropWithNilCheck(network, "name", nic.Properties.Name)
	setPropWithNilCheck(network, "mac", nic.Properties.Mac)

	if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
		network["ips"] = *nic.Properties.Ips
	}

	return network
}

func getServerData(d *schema.ResourceData, update bool) (*ionoscloud.Server, error) {
	server := ionoscloud.Server{
		Properties: &ionoscloud.ServerProperties{},
	}

	if update {
		if d.HasChange("availability_zone") {
			return nil, fmt.Errorf("availability_zone is immutable")
		}
	} else {
		if v, ok := d.GetOk("availability_zone"); ok {
			vStr := v.(string)
			server.Properties.AvailabilityZone = &vStr
		}
	}

	serverName := d.Get("name").(string)
	server.Properties.Name = &serverName

	serverCores := int32(d.Get("cores").(int))
	server.Properties.Cores = &serverCores

	serverRam := int32(d.Get("ram").(int))
	server.Properties.Ram = &serverRam

	if v, ok := d.GetOk("cpu_family"); ok {
		if v.(string) != "" {
			vStr := v.(string)
			server.Properties.CpuFamily = &vStr
		}
	}

	if _, ok := d.GetOk("boot_cdrom"); ok {
		bootCdrom := d.Get("boot_cdrom").(string)
		if IsValidUUID(bootCdrom) {
			server.Properties.BootCdrom = &ionoscloud.ResourceReference{
				Id: &bootCdrom,
			}

		} else {
			return nil, fmt.Errorf("boot_cdrom has to be a valid UUID, got: %s", bootCdrom)
		}
	}

	return &server, nil
}

func setServerData(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData, server *ionoscloud.Server, readFromSchema bool) error {

	if server.Id != nil {
		d.SetId(*server.Id)
	}

	datacenterId := d.Get("datacenter_id").(string)
	if server.Properties != nil {
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

		if server.Properties.AvailabilityZone != nil {
			if err := d.Set("availability_zone", *server.Properties.AvailabilityZone); err != nil {
				return err
			}
		}

		if server.Properties.VmState != nil {
			if err := d.Set("vm_state", *server.Properties.VmState); err != nil {
				return err
			}
		}

		if server.Properties.CpuFamily != nil {
			if err := d.Set("cpu_family", *server.Properties.CpuFamily); err != nil {
				return err
			}
		}
		if server.Properties.BootCdrom != nil && server.Properties.BootCdrom.Id != nil {
			if err := d.Set("boot_cdrom", *server.Properties.BootCdrom.Id); err != nil {
				return err
			}
		}

		if server.Properties.BootVolume != nil && server.Properties.BootVolume.Id != nil {
			if err := d.Set("boot_volume", *server.Properties.BootVolume.Id); err != nil {
				return err
			}
		}

		if server.Entities.Volumes != nil && server.Entities.Volumes.Items != nil && len(*server.Entities.Volumes.Items) > 0 &&
			(*server.Entities.Volumes.Items)[0].Properties.Image != nil {
			if err := d.Set("boot_image", *(*server.Entities.Volumes.Items)[0].Properties.Image); err != nil {
				return err
			}
		}
	}

	if server.Entities == nil {
		return nil
	}

	var cdroms []interface{}
	if server.Entities.Cdroms != nil && server.Entities.Cdroms.Items != nil && len(*server.Entities.Cdroms.Items) > 0 {
		for _, image := range *server.Entities.Cdroms.Items {
			entry := SetCdromProperties(image)
			cdroms = append(cdroms, entry)
		}
		if err := d.Set("cdrom", cdroms); err != nil {
			return err
		}
	}

	if server.Properties.BootVolume != nil {
		volume, _, err := client.ServerApi.DatacentersServersVolumesFindById(ctx, datacenterId, d.Id(), *server.Properties.BootVolume.Id).Execute()
		if err == nil {
			var volumes []interface{}
			entry := SetVolumeProperties(volume)
			userData := d.Get("volume.0.user_data")
			entry["user_data"] = userData

			backupUnit := d.Get("volume.0.backup_unit_id")
			entry["backup_unit_id"] = backupUnit
			volumes = append(volumes, entry)
			if err := d.Set("volume", volumes); err != nil {
				return err
			}
		}
	}

	_, primaryNicOk := d.GetOk("primary_nic")
	_, primaryFirewallOk := d.GetOk("firewallrule_id")

	// take nic and firewall from schema if set is used in resource read, else take it from entities
	if (readFromSchema && primaryNicOk) || (!readFromSchema && server.Entities.Nics != nil && server.Entities.Nics.Items != nil && len(*server.Entities.Nics.Items) > 0 && (*server.Entities.Nics.Items)[0].Id != nil) {
		var nicId string
		if readFromSchema {
			nicId = d.Get("primary_nic").(string)
		} else {
			nicId = *(*server.Entities.Nics.Items)[0].Id
		}

		nic, _, err := client.NicApi.DatacentersServersNicsFindById(ctx, datacenterId, d.Id(), nicId).Execute()
		if err != nil {
			return err
		}
		nicEntry := SetNetworkProperties(nic)

		if (readFromSchema && primaryFirewallOk) || !readFromSchema {
			var firewallId string
			if readFromSchema {
				firewallId = d.Get("firewallrule_id").(string)
			} else {
				firewallId = *(*nic.Entities.Firewallrules.Items)[0].Id
			}

			firewall, _, err := client.NicApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, d.Id(), nicId, firewallId).Execute()
			if err != nil {
				return err
			}

			firewallEntry := SetFirewallProperties(firewall)

			nicEntry["firewall"] = []map[string]interface{}{firewallEntry}

		}

		nics := []map[string]interface{}{nicEntry}

		if err := d.Set("nic", nics); err != nil {
			return err
		}
	}

	return nil
}
