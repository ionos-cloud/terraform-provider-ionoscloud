package ionoscloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
		CustomizeDiff: checkServerImmutableFields,

		Schema: map[string]*schema.Schema{
			"template_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"cores": {
				Type:     schema.TypeInt,
				Optional: true, // this should be required when the deprecated version will be removed
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Optional: true, // this should be required when the deprecated version will be removed
				Computed: true,
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2"}, true)),
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
				Description:      "server usages: ENTERPRISE or CUBE",
				DiffSuppressFunc: utils.DiffToLower,
				//to do: add in next release
				//ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CUBE", "ENTERPRISE"}, true)),
			},
			"boot_image": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"primary_nic": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the primary network interface",
			},
			"primary_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"firewallrule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"firewallrule_ids": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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
				ConflictsWith: []string{"volume.0.ssh_key_path", "volume.0.ssh_keys", "ssh_keys"},
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Will be renamed to ssh_keys in the future, to allow users to set both the ssh key path or directly the ssh key",
			},
			"ssh_keys": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"volume.0.ssh_key_path", "volume.0.ssh_keys", "ssh_key_path"},
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Description:   "Public SSH keys are set on the image as authorized keys for appropriate SSH login to the instance using the corresponding private key. This field may only be set in creation requests. When reading, it always returns null. SSH keys are only supported if a public Linux image is used for the volume creation.",
			},
			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The size of the volume in GB.",
						},
						"disk_type": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Deprecated:  "Please use ssh_key_path under server level",
							Description: "Public SSH keys are set on the image as authorized keys for appropriate SSH login to the instance using the corresponding private key. This field may only be set in creation requests. When reading, it always returns null. SSH keys are only supported if a public Linux image is used for the volume creation.",
							Computed:    true,
							ForceNew:    true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if k == "volume.0.ssh_key_path.#" {
									if d.Get("ssh_key_path.#") == new {
										return true
									}
								}

								sshKeyPath := d.Get("volume.0.ssh_key_path").([]interface{})
								oldSshKeyPath := d.Get("ssh_key_path").([]interface{})

								difKeypath := utils.DiffSlice(convertSlice(sshKeyPath), convertSlice(oldSshKeyPath))
								if len(difKeypath) == 0 {
									return true
								}

								return false
							},
						},
						"ssh_keys": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Deprecated:  "Please use ssh_keys under server level",
							Computed:    true,
							ForceNew:    true,
							Description: "Public SSH keys are set on the image as authorized keys for appropriate SSH login to the instance using the corresponding private key. This field may only be set in creation requests. When reading, it always returns null. SSH keys are only supported if a public Linux image is used for the volume creation.",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if k == "volume.0.ssh_keys.#" {
									if d.Get("ssh_keys.#") == new {
										return true
									}
								}

								sshKeys := d.Get("volume.0.ssh_keys").([]interface{})
								oldSshKeys := d.Get("ssh_keys").([]interface{})

								if len(utils.DiffSlice(convertSlice(sshKeys), convertSlice(oldSshKeys))) == 0 {
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
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2", "ZONE_3"}, true)),
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								DiffSuppressFunc: utils.DiffEmptyIps,
							},
							Description: "Collection of IP addresses assigned to a nic. Explicitly assigned public IPs need to come from reserved IP blocks, Passing value null or empty array will assign an IP address automatically.",
							Computed:    true,
							Optional:    true,
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"protocol": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: utils.DiffToLower,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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
										ValidateDiagFunc: validation.ToDiagFunc(func(v interface{}, k string) (ws []string, errors []error) {
											if v.(int) < 1 && v.(int) > 65534 {
												errors = append(errors, fmt.Errorf("port start range must be between 1 and 65534"))
											}
											return
										}),
									},
									"port_range_end": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateDiagFunc: validation.ToDiagFunc(func(v interface{}, k string) (ws []string, errors []error) {
											if v.(int) < 1 && v.(int) > 65534 {
												errors = append(errors, fmt.Errorf("port end range must be between 1 and 65534"))
											}
											return
										}),
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
			"label": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     labelResource,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func checkServerImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {

	//we do not want to check in case of resource creation
	if diff.Id() == "" {
		return nil
	}
	if diff.HasChange("image_name") {
		return fmt.Errorf("image_name %s", ImmutableError)
	}

	if diff.HasChange("availability_zone") {
		return fmt.Errorf("availability_zone %s", ImmutableError)
	}
	if diff.HasChange("volume") {
		if diff.HasChange("volume.0.user_data") {
			return fmt.Errorf("volume.0.user_data %s", ImmutableError)
		}

		if diff.HasChange("volume.0.backup_unit_id") {
			return fmt.Errorf("volume.0.backup_unit_id %s", ImmutableError)
		}

		if diff.HasChange("volume.0.disk_type") {
			return fmt.Errorf("volume.0.disk_type %s", ImmutableError)
		}

		if diff.HasChange("volume.0.availability_zone") {
			return fmt.Errorf("volume.0.availability_zone %s", ImmutableError)
		}
	}
	return nil

}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient
	datacenterId := d.Get("datacenter_id").(string)

	serverReq, err := initializeCreateRequests(d)
	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	// create volume object with data to be used for image
	volume, err := getVolumeData(d, "volume.0.")

	if err != nil {
		return diag.FromErr(err)
	}

	// get image and imageAlias
	image, imageAlias, err := getImage(ctx, client, d, *volume)
	if err != nil {
		return diag.FromErr(err)
	}

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
		if utils.IsValidUUID(backupUnitId.(string)) {
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

	// add volume object to serverReq
	serverReq.Entities = &ionoscloud.ServerEntities{
		Volumes: &ionoscloud.AttachedVolumes{
			Items: &[]ionoscloud.Volume{
				{
					Properties: volume,
				},
			},
		},
	}

	// get nic data and add object to serverReq
	if nics, ok := d.GetOk("nic"); ok {
		serverReq.Entities.Nics = &ionoscloud.Nics{
			Items: &[]ionoscloud.Nic{},
		}
		if nics.([]interface{}) != nil {
			for nicIndex := range nics.([]interface{}) {
				nicPath := fmt.Sprintf("nic.%d.", nicIndex)
				nic := getNicData(d, nicPath)
				*serverReq.Entities.Nics.Items = append(*serverReq.Entities.Nics.Items, nic)
				fwRulesPath := nicPath + "firewall"
				if firewallRules, ok := d.GetOk(fwRulesPath); ok {
					fwRules := []ionoscloud.FirewallRule{}
					(*serverReq.Entities.Nics.Items)[nicIndex].Entities = &ionoscloud.NicEntities{
						Firewallrules: &ionoscloud.FirewallRules{
							Items: &[]ionoscloud.FirewallRule{},
						},
					}

					if firewallRules.([]interface{}) != nil && len(firewallRules.([]interface{})) > 0 {
						fwRulesIntf := firewallRules.([]interface{})

						fwRulesProperties := make([]ionoscloud.FirewallruleProperties, len(fwRulesIntf))
						err = utils.DecodeInterfaceToStruct(fwRulesIntf, fwRulesProperties)
						if err != nil {
							return diag.FromErr(fmt.Errorf("could not decode from %s to slice of FirewallRules %w", fwRulesIntf, err))
						}
						for _, prop := range fwRulesProperties {
							FwPropUnsetSetFieldIfNotSetInSchema(&prop, fwRulesPath, d)
							firewall := ionoscloud.FirewallRule{
								Properties: &prop,
							}
							fwRules = append(fwRules, firewall)
						}
					}
					*(*serverReq.Entities.Nics.Items)[nicIndex].Entities.Firewallrules.Items = fwRules
				}
			}
		}
	}

	postServer, apiResponse, err := client.ServersApi.DatacentersServersPost(ctx, datacenterId).Server(serverReq).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error creating server: (%w)", err))
		return diags
	}

	if postServer.Id != nil {
		d.SetId(*postServer.Id)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			log.Printf("[DEBUG] failed to create server resource")
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}

		diags := diag.FromErr(fmt.Errorf("error waiting for state change for server creation %w", errState))
		return diags
	}

	// Logic for labels creation
	ls := LabelsService{ctx: ctx, client: client}
	if err := ls.datacentersServersLabelsCreate(datacenterId, *postServer.Id, d.Get("label")); err != nil {
		return diag.FromErr(err)
	}

	// get additional data for schema
	foundServer, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, datacenterId, *postServer.Id).Depth(4).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error fetching server: %w", err))
		return diags
	}
	if foundServer.Entities.Nics.Items != nil {
		if len(*foundServer.Entities.Nics.Items) > 0 {
			foundFirstNic := (*foundServer.Entities.Nics.Items)[0]
			var rules []string
			if foundFirstNic.Entities != nil && foundFirstNic.Entities.Firewallrules != nil && foundFirstNic.Entities.Firewallrules.Items != nil {
				sentFirstNic := (*serverReq.Entities.Nics.Items)[0]

				if sentFirstNic.Entities != nil && sentFirstNic.Entities.Firewallrules != nil && sentFirstNic.Entities.Firewallrules.Items != nil {
					sentRules := *sentFirstNic.Entities.Firewallrules.Items
					foundRules := *foundFirstNic.Entities.Firewallrules.Items
					if len(sentRules) > 0 {

						if err := d.Set("firewallrule_id", *(foundRules)[0].Id); err != nil {
							diags := diag.FromErr(err)
							return diags
						}
						if len(sentRules) > 0 {
							//keep order of rules
							for _, rule := range sentRules {
								for _, foundRule := range foundRules {
									if reflect.DeepEqual(rule.Properties, foundRule.Properties) {
										rules = append(rules, *foundRule.Id)
									}
								}
							}
						}
					}
				}
			}
			if err := d.Set("firewallrule_ids", rules); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
			if foundFirstNic.Id != nil {
				err := d.Set("primary_nic", *foundFirstNic.Id)
				if err != nil {
					diags := diag.FromErr(fmt.Errorf("error while setting primary nic %s: %w", d.Id(), err))
					return diags
				}
			}
			firstNicProps := foundFirstNic.Properties
			if firstNicProps != nil {
				firstNicIps := firstNicProps.Ips
				if firstNicIps != nil &&
					len(*firstNicIps) > 0 {
					log.Printf("[DEBUG] set primary_ip to %s", (*firstNicIps)[0])
					if err := d.Set("primary_ip", (*firstNicIps)[0]); err != nil {
						diags := diag.FromErr(fmt.Errorf("error while setting primary ip %s: %w", d.Id(), err))
						return diags
					}
				}

				volumeItems := serverReq.Entities.Volumes.Items
				firstVolumeItem := (*volumeItems)[0]
				if firstNicProps.Ips != nil &&
					len(*firstNicIps) > 0 &&
					volumeItems != nil &&
					len(*volumeItems) > 0 &&
					firstVolumeItem.Properties != nil &&
					firstVolumeItem.Properties.ImagePassword != nil {

					d.SetConnInfo(map[string]string{
						"type":     "ssh",
						"host":     (*firstNicIps)[0],
						"password": *firstVolumeItem.Properties.ImagePassword,
					})
				}
			}
		}

	}
	return resourceServerRead(ctx, d, meta)
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	serverId := d.Id()

	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, serverId).Depth(4).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			log.Printf("[DEBUG] cannot find server by id \n")
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a server ID %s %w", d.Id(), err))
		return diags
	}
	if err := setResourceServerData(ctx, client, d, &server); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func SetNetworkProperties(nic ionoscloud.Nic) map[string]interface{} {

	network := map[string]interface{}{}
	if nic.Properties != nil {
		utils.SetPropWithNilCheck(network, "dhcp", nic.Properties.Dhcp)
		utils.SetPropWithNilCheck(network, "firewall_active", nic.Properties.FirewallActive)
		utils.SetPropWithNilCheck(network, "firewall_type", nic.Properties.FirewallType)
		utils.SetPropWithNilCheck(network, "lan", nic.Properties.Lan)
		utils.SetPropWithNilCheck(network, "name", nic.Properties.Name)
		utils.SetPropWithNilCheck(network, "ips", nic.Properties.Ips)
		utils.SetPropWithNilCheck(network, "mac", nic.Properties.Mac)
		if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
			network["ips"] = *nic.Properties.Ips
		}
	}
	return network
}

func SetVolumeProperties(volume ionoscloud.Volume) map[string]interface{} {

	volumeMap := map[string]interface{}{}
	if volume.Properties != nil {
		if volume.Properties.SshKeys != nil && len(*volume.Properties.SshKeys) > 0 {
			var sshKeys []interface{}
			for _, sshKey := range *volume.Properties.SshKeys {
				sshKeys = append(sshKeys, sshKey)
			}
			volumeMap["ssh_keys"] = sshKeys
		}

		utils.SetPropWithNilCheck(volumeMap, "image_password", volume.Properties.ImagePassword)
		utils.SetPropWithNilCheck(volumeMap, "name", volume.Properties.Name)

		utils.SetPropWithNilCheck(volumeMap, "size", volume.Properties.Size)
		utils.SetPropWithNilCheck(volumeMap, "disk_type", volume.Properties.Type)
		utils.SetPropWithNilCheck(volumeMap, "licence_type", volume.Properties.LicenceType)
		utils.SetPropWithNilCheck(volumeMap, "bus", volume.Properties.Bus)
		utils.SetPropWithNilCheck(volumeMap, "availability_zone", volume.Properties.AvailabilityZone)
		utils.SetPropWithNilCheck(volumeMap, "cpu_hot_plug", volume.Properties.CpuHotPlug)
		utils.SetPropWithNilCheck(volumeMap, "ram_hot_plug", volume.Properties.RamHotPlug)
		utils.SetPropWithNilCheck(volumeMap, "nic_hot_plug", volume.Properties.NicHotPlug)
		utils.SetPropWithNilCheck(volumeMap, "nic_hot_unplug", volume.Properties.NicHotUnplug)
		utils.SetPropWithNilCheck(volumeMap, "disc_virtio_hot_plug", volume.Properties.DiscVirtioHotPlug)
		utils.SetPropWithNilCheck(volumeMap, "disc_virtio_hot_unplug", volume.Properties.DiscVirtioHotUnplug)
		utils.SetPropWithNilCheck(volumeMap, "device_number", volume.Properties.DeviceNumber)
		utils.SetPropWithNilCheck(volumeMap, "user_data", volume.Properties.UserData)
		utils.SetPropWithNilCheck(volumeMap, "backup_unit_id", volume.Properties.BackupunitId)
		utils.SetPropWithNilCheck(volumeMap, "boot_server", volume.Properties.BootServer)
	}
	return volumeMap
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

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

	if d.HasChange("cpu_family") {
		_, n := d.GetChange("cpu_family")
		nStr := n.(string)
		request.CpuFamily = &nStr
	}

	if d.HasChange("boot_cdrom") {
		_, n := d.GetChange("boot_cdrom")
		bootCdrom := n.(string)

		if utils.IsValidUUID(bootCdrom) {
			request.BootCdrom = &ionoscloud.ResourceReference{
				Id: &bootCdrom,
			}
		} else {
			diags := diag.FromErr(fmt.Errorf("boot_cdrom has to be a valid UUID, got: %s", bootCdrom))
			return diags
		}
		/* todo: figure out a way of sending a nil bootCdrom to the API (the sdk's omitempty doesn't let us) */
	}

	server, apiResponse, err := client.ServersApi.DatacentersServersPatch(ctx, dcId, d.Id()).Server(request).Depth(3).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while updating server ID %s: %w", d.Id(), err))
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
		_, apiResponse, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, d.Id(), bootVolume).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			volume := ionoscloud.Volume{
				Id: &bootVolume,
			}
			_, apiResponse, err := client.ServersApi.DatacentersServersVolumesPost(ctx, dcId, d.Id()).Volume(volume).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occured while attaching a volume dcId: %s server_id: %s ID: %s Response: %s", dcId, d.Id(), bootVolume, err))
				return diags
			}

			// Wait, catching any errors
			_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
			if errState != nil {
				diags := diag.FromErr(fmt.Errorf("an error occured while waiting for a state change for dcId: %s server_id: %s ID: %s %w", dcId, d.Id(), bootVolume, err))
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

		_, apiResponse, err = client.VolumesApi.DatacentersVolumesPatch(ctx, d.Get("datacenter_id").(string), bootVolume).Volume(properties).Execute()
		logApiRequestTime(apiResponse)

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
		updateNic := true
		primaryNic := d.Get("primary_nic").(string)
		if primaryNic == "" {
			updateNic = false
		}

		nic := &ionoscloud.Nic{}
		for _, n := range *server.Entities.Nics.Items {
			nicStr := d.Get("primary_nic").(string)
			if *n.Id == nicStr {
				nic = &n
				break
			}
		}

		lan := int32(d.Get("nic.0.lan").(int))
		nicProperties := ionoscloud.NicProperties{
			Lan: &lan,
		}

		if v, ok := d.GetOk("nic.0.name"); ok {
			vStr := v.(string)
			nicProperties.Name = &vStr
		}

		if v, ok := d.GetOk("nic.0.ips"); ok {
			raw := v.([]interface{})
			if raw != nil && len(raw) > 0 {
				ips := make([]string, 0)
				for _, rawIp := range raw {
					if rawIp != nil {
						ip := rawIp.(string)
						ips = append(ips, ip)
					}
				}
				if ips != nil && len(ips) > 0 {
					nicProperties.Ips = &ips
				}
			}
		}

		dhcp := d.Get("nic.0.dhcp").(bool)
		fwRule := d.Get("nic.0.firewall_active").(bool)
		nicProperties.Dhcp = &dhcp
		nicProperties.FirewallActive = &fwRule

		if v, ok := d.GetOk("nic.0.firewall_type"); ok {
			vStr := v.(string)
			nicProperties.FirewallType = &vStr
		}
		firstNicFirewallPath := "nic.0.firewall"
		if d.HasChange(firstNicFirewallPath) {
			oldValues, newValues := d.GetChange(firstNicFirewallPath)
			oldFwSlice, newFwSlice, err := GetChangesInFirewallRules(oldValues.([]interface{}), newValues.([]interface{}))
			if err != nil {
				return diag.FromErr(fmt.Errorf("could not get changes for FirewallRules %w", err))
			}

			firewallId := d.Get("firewallrule_id").(string)
			firewallRuleIds := d.Get("firewallrule_ids").([]interface{})
			fwRuleIdsString := convertSlice(firewallRuleIds)

			if nic != nil && nic.Id != nil {
				fs := FirewallService{client: client, meta: meta, schemaData: d}

				firewalls, err := fs.firewallsGet(ctx, dcId, *server.Id, *nic.Id, 1)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching firewall rules: for dcId: %s server_id: %s nic_id %s `%w`", dcId, *server.Id, *nic.Id, err))
				}
				//delete old rules
				for _, ruleProp := range oldFwSlice {
					for _, item := range firewalls {
						if strings.EqualFold(*ruleProp.Name, *item.Properties.Name) {
							if firewallId == *item.Id {
								firewallId = ""
							}
							apiResponse, err = fs.deleteFirewallRule(ctx, dcId, *server.Id, *nic.Id, *item.Id)
							if err != nil {
								return diag.FromErr(fmt.Errorf("an error occured while deleting firewall rule for dcId: %s server_id: %s nic_id %s ID: %s `%w`", dcId, *server.Id, *nic.Id, *item.Id, err))
							}

							if utils.Contains(fwRuleIdsString, *item.Id) {
								fwRuleIdsString = utils.DeleteFromSlice(fwRuleIdsString, *item.Id)
							}
						}
					}
				}
			}

			firewallRules := []ionoscloud.FirewallRule{}
			// create updated rules
			for _, ruleProp := range newFwSlice {
				FwPropUnsetSetFieldIfNotSetInSchema(&ruleProp, firstNicFirewallPath, d)

				fwRule := ionoscloud.FirewallRule{
					Properties: &ruleProp,
				}
				if nic != nil && nic.Id != nil {
					fs := FirewallService{client: client, meta: meta, schemaData: d}
					firewall, _, err := fs.createFirewallRule(ctx, dcId, *server.Id, *nic.Id, fwRule)
					if err != nil {
						return diag.FromErr(err)
					}
					firewallRules = append(firewallRules, *firewall)
					fwRuleIdsString = append(fwRuleIdsString, *firewall.Id)
				} else { //if the nic does not exist, just fw add rule to the list to be created below with the nic
					firewallRules = append(firewallRules, fwRule)
				}
			}
			if firewallId != "" {
				if err := d.Set("firewallrule_id", firewallId); err != nil {
					diags := diag.FromErr(err)
					return diags
				}
			}
			if err := d.Set("firewallrule_ids", utils.ToAnyList(fwRuleIdsString)); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
			if len(firewallRules) > 0 {
				nic.Entities = &ionoscloud.NicEntities{
					Firewallrules: &ionoscloud.FirewallRules{
						Items: &firewallRules,
					},
				}
			}
		}
		mProp, _ := json.Marshal(nicProperties)
		log.Printf("[DEBUG] Updating props: %s", string(mProp))
		var createdNic ionoscloud.Nic
		if !updateNic {
			nic.Properties = &nicProperties
			createdNic, apiResponse, err = client.NetworkInterfacesApi.DatacentersServersNicsPost(ctx, d.Get("datacenter_id").(string), *server.Id).Nic(*nic).Execute()

		} else {
			_, apiResponse, err = client.NetworkInterfacesApi.DatacentersServersNicsPatch(ctx, d.Get("datacenter_id").(string), *server.Id, *nic.Id).Nic(nicProperties).Execute()
		}
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error updating nic (%s)", err))
			return diags
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(fmt.Errorf("error getting state change for nics patch %w", errState))
			return diags
		}

		if createdNic.Id != nil {
			fs := FirewallService{client: client, meta: meta, schemaData: d}
			fwRules, err := fs.firewallsGet(ctx, d.Get("datacenter_id").(string), *server.Id, *createdNic.Id, 1)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occurred while fetching firewall rules: %w", err))
				return diags
			}
			if len(fwRules) > 0 {
				if err := d.Set("firewallrule_id", fwRules[0].Id); err != nil {
					diags := diag.FromErr(err)
					return diags
				}
			}
		}

		_, primaryNicOk := d.GetOk("primary_nic")
		if createdNic.Id != nil && !updateNic && !primaryNicOk {
			if err := d.Set("primary_nic", *createdNic.Id); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
	}

	// Labels logic for update
	if d.HasChanges("label") {
		ls := LabelsService{ctx: ctx, client: client}
		oldLabelsData, newLabelsData := d.GetChange("label")

		// Make some differences to see exactly what labels need to be added/removed.
		labelsToBeCreated := newLabelsData.(*schema.Set).Difference(oldLabelsData.(*schema.Set))
		labelsToBeDeleted := oldLabelsData.(*schema.Set).Difference(newLabelsData.(*schema.Set))

		if err := ls.datacentersServersLabelsDelete(dcId, d.Id(), labelsToBeDeleted); err != nil {
			return diag.FromErr(err)
		}

		if err := ls.datacentersServersLabelsCreate(dcId, d.Id(), labelsToBeCreated); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceServerRead(ctx, d, meta)
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient
	dcId := d.Get("datacenter_id").(string)
	// A bigger depth is required since we need all volumes items.
	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, d.Id()).Depth(2).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a server ID %s %w", d.Id(), err))
		return diags
	}

	if strings.ToLower(*server.Properties.Type) != "cube" {
		diags := deleteAllVolumes(ctx, d, meta, client, &server)
		if diags != nil {
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

func deleteAllVolumes(ctx context.Context, d *schema.ResourceData, meta interface{}, client *ionoscloud.APIClient, server *ionoscloud.Server) diag.Diagnostics {
	dcId := d.Get("datacenter_id").(string)

	if server.Entities == nil || server.Entities.Volumes == nil || server.Entities.Volumes.Items == nil {
		return nil
	}

	volumes := server.Entities.Volumes.Items
	for _, volume := range *volumes {
		apiResponse, err := client.VolumesApi.DatacentersVolumesDelete(ctx, dcId, *volume.Id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error occured while deleting remaining volume %s of server ID %s %w", *volume.Id, d.Id(), err))
			return diags
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(fmt.Errorf("error getting state change for volumes delete %w", errState))
			return diags
		}

	}
	return nil
}

func resourceServerImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter UUID}/{server UUID}", d.Id())
	}

	datacenterId := parts[0]
	serverId := parts[1]

	client := meta.(SdkBundle).CloudApiClient

	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, datacenterId, serverId).Depth(3).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find server %q", serverId)
		}
		return nil, fmt.Errorf("error occured while fetching a server ID %s %w", d.Id(), err)
	}
	var primaryNic ionoscloud.Nic
	d.SetId(*server.Id)
	primaryNicId := ""
	//first we try to get primary nic from parts, then if that fails, we get it from entities.
	if len(parts) > 2 {
		primaryNicId = parts[2]
		if err := d.Set("primary_nic", primaryNicId); err != nil {
			return nil, fmt.Errorf("error setting primary_nic id %w", err)
		}
	} else {
		if server.Entities != nil && server.Entities.Nics != nil && len(*server.Entities.Nics.Items) > 0 {
			primaryNic = (*server.Entities.Nics.Items)[0]
		}
	}
	if primaryNicId != "" {
		if server.Entities != nil && server.Entities.Nics != nil && server.Entities.Nics.Items != nil {
			for _, nic := range *server.Entities.Nics.Items {
				if *nic.Id == primaryNicId {
					primaryNic = nic
					break
				}
			}
		}
	}

	log.Printf("[DEBUG] set primary_ip to %s", (*primaryNic.Properties.Ips)[0])
	if err := d.Set("primary_ip", (*primaryNic.Properties.Ips)[0]); err != nil {
		return nil, fmt.Errorf("error while setting primary ip %s: %w", d.Id(), err)
	}
	if err := d.Set("datacenter_id", datacenterId); err != nil {
		return nil, err
	}

	if len(parts) > 3 {
		if err := d.Set("firewallrule_id", parts[3]); err != nil {
			return nil, fmt.Errorf("error setting firewallrule_id %w", err)
		}
		var rules []string
		rules = append(rules, parts[3])
		if err := d.Set("firewallrule_ids", rules); err != nil {
			return nil, utils.GenerateSetError(ServerResource, "firewallrule_ids", err)
		}
	}

	if err := setResourceServerData(ctx, client, d, &server); err != nil {
		return nil, err
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

// Reads public key from file or directly provided and returns key string if valid
func readPublicKey(pathOrKey string) (string, error) {
	var bytes []byte
	var err error
	if utils.CheckFileExists(pathOrKey) {
		bytes, err = os.ReadFile(pathOrKey)
		if err != nil {

			return "", err
		}
	} else {
		log.Printf("[DEBUG] error opening file, key must have been provided directly %s ", pathOrKey)
		bytes = []byte(pathOrKey)
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	if err != nil {
		return "", fmt.Errorf("error for public key %s, check if path is correct or key is in correct format", pathOrKey)
	}
	return string(ssh.MarshalAuthorizedKey(pubKey)[:]), nil
}

func SetCdromProperties(image ionoscloud.Image) map[string]interface{} {

	cdrom := make(map[string]interface{})
	if image.Properties != nil {
		utils.SetPropWithNilCheck(cdrom, "name", image.Properties.Name)
		utils.SetPropWithNilCheck(cdrom, "description", image.Properties.Description)
		utils.SetPropWithNilCheck(cdrom, "location", image.Properties.Location)
		utils.SetPropWithNilCheck(cdrom, "size", image.Properties.Size)
		utils.SetPropWithNilCheck(cdrom, "cpu_hot_plug", image.Properties.CpuHotPlug)
		utils.SetPropWithNilCheck(cdrom, "cpu_hot_unplug", image.Properties.CpuHotUnplug)
		utils.SetPropWithNilCheck(cdrom, "ram_hot_plug", image.Properties.RamHotPlug)
		utils.SetPropWithNilCheck(cdrom, "ram_hot_unplug", image.Properties.RamHotUnplug)
		utils.SetPropWithNilCheck(cdrom, "nic_hot_plug", image.Properties.NicHotPlug)
		utils.SetPropWithNilCheck(cdrom, "nic_hot_unplug", image.Properties.NicHotUnplug)
		utils.SetPropWithNilCheck(cdrom, "disc_virtio_hot_plug", image.Properties.DiscVirtioHotPlug)
		utils.SetPropWithNilCheck(cdrom, "disc_virtio_hot_unplug", image.Properties.DiscVirtioHotUnplug)
		utils.SetPropWithNilCheck(cdrom, "disc_scsi_hot_plug", image.Properties.DiscScsiHotPlug)
		utils.SetPropWithNilCheck(cdrom, "disc_scsi_hot_unplug", image.Properties.DiscScsiHotUnplug)
		utils.SetPropWithNilCheck(cdrom, "licence_type", image.Properties.LicenceType)
		utils.SetPropWithNilCheck(cdrom, "image_type", image.Properties.ImageType)
		utils.SetPropWithNilCheck(cdrom, "public", image.Properties.Public)
	}

	return cdrom
}

// Initializes server with the required attributes depending on the server type (CUBE or ENTERPRISE)
func initializeCreateRequests(d *schema.ResourceData) (ionoscloud.Server, error) {

	serverType := d.Get("type").(string)

	// create server object and populate with common attributes
	server, err := getServerData(d)
	if err != nil {
		return *server, err
	}

	if serverType != "" {
		server.Properties.Type = &serverType
	}
	switch strings.ToLower(serverType) {
	case "cube":
		if v, ok := d.GetOk("template_uuid"); ok {
			vStr := v.(string)
			server.Properties.TemplateUuid = &vStr
		} else {
			return *server, fmt.Errorf("template_uuid argument is required for %s type of servers\n", serverType)
		}

		if _, ok := d.GetOk("cores"); ok {
			return *server, fmt.Errorf("cores argument can not be set for %s type of servers\n", serverType)
		}

		if _, ok := d.GetOk("ram"); ok {
			return *server, fmt.Errorf("ram argument can not be set for %s type of servers\n", serverType)
		}

		if _, ok := d.GetOk("volume.0.size"); ok {
			return *server, fmt.Errorf("volume.0.size argument can not be set for %s type of servers\n", serverType)
		}
	default: //enterprise
		if _, ok := d.GetOk("template_uuid"); ok {
			return *server, fmt.Errorf("template_uuid argument can not be set only for %s type of servers\n", serverType)
		}

		if v, ok := d.GetOk("cores"); ok {
			vInt := int32(v.(int))
			server.Properties.Cores = &vInt
		} else {
			return *server, fmt.Errorf("cores argument is required for %s type of servers\n", serverType)
		}

		if v, ok := d.GetOk("ram"); ok {
			vInt := int32(v.(int))
			server.Properties.Ram = &vInt
		} else {
			return *server, fmt.Errorf("ram argument is required for %s type of servers\n", serverType)
		}

		if _, ok := d.GetOk("volume.0.size"); !ok {
			return *server, fmt.Errorf("volume.0.size argument is required for %s type of servers\n", serverType)
		}
	}
	return *server, nil
}

func getServerData(d *schema.ResourceData) (*ionoscloud.Server, error) {
	server := ionoscloud.NewServer(*ionoscloud.NewServerPropertiesWithDefaults())

	if v, ok := d.GetOk("availability_zone"); ok {
		vStr := v.(string)
		server.Properties.AvailabilityZone = &vStr
	}

	serverName := d.Get("name").(string)
	server.Properties.Name = &serverName

	if v, ok := d.GetOk("cpu_family"); ok {
		if v.(string) != "" {
			vStr := v.(string)
			server.Properties.CpuFamily = &vStr
		}
	}

	if _, ok := d.GetOk("boot_cdrom"); ok {
		bootCdrom := d.Get("boot_cdrom").(string)
		if utils.IsValidUUID(bootCdrom) {
			server.Properties.BootCdrom = &ionoscloud.ResourceReference{
				Id: &bootCdrom,
			}

		} else {
			return nil, fmt.Errorf("boot_cdrom has to be a valid UUID, got: %s", bootCdrom)
		}
	}

	return server, nil
}

func setResourceServerData(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData, server *ionoscloud.Server) error {
	if server.Id != nil {
		d.SetId(*server.Id)
	}

	datacenterId := d.Get("datacenter_id").(string)
	if server.Properties != nil {
		if server.Properties.Name != nil {
			if err := d.Set("name", *server.Properties.Name); err != nil {
				return fmt.Errorf("error setting name %w", err)
			}
		}

		if server.Properties.Cores != nil {
			if err := d.Set("cores", *server.Properties.Cores); err != nil {
				return fmt.Errorf("error setting cores %w", err)
			}
		}

		if server.Properties.Ram != nil {
			if err := d.Set("ram", *server.Properties.Ram); err != nil {
				return fmt.Errorf("error setting ram %w", err)
			}
		}

		if server.Properties.AvailabilityZone != nil {
			if err := d.Set("availability_zone", *server.Properties.AvailabilityZone); err != nil {
				return fmt.Errorf("error setting availability_zone %w", err)
			}
		}

		if server.Properties.CpuFamily != nil {
			if err := d.Set("cpu_family", *server.Properties.CpuFamily); err != nil {
				return fmt.Errorf("error setting cpu_family %w", err)
			}
		}

		if server.Properties.Type != nil {
			if err := d.Set("type", *server.Properties.Type); err != nil {
				return fmt.Errorf("error setting type %w", err)
			}
		}

		if server.Properties.BootCdrom != nil && server.Properties.BootCdrom.Id != nil {
			if err := d.Set("boot_cdrom", *server.Properties.BootCdrom.Id); err != nil {
				return fmt.Errorf("error setting boot_cdrom %w", err)
			}
		}

		if server.Properties.BootVolume != nil && server.Properties.BootVolume.Id != nil {
			if err := d.Set("boot_volume", *server.Properties.BootVolume.Id); err != nil {
				return fmt.Errorf("error setting bootVolume %w", err)
			}
		}

		if server.Entities != nil && server.Entities.Volumes != nil && server.Entities.Volumes.Items != nil && len(*server.Entities.Volumes.Items) > 0 &&
			(*server.Entities.Volumes.Items)[0].Properties != nil && (*server.Entities.Volumes.Items)[0].Properties.Image != nil {
			if err := d.Set("boot_image", *(*server.Entities.Volumes.Items)[0].Properties.Image); err != nil {
				return fmt.Errorf("error setting boot_image %w", err)
			}
		}
	}

	if server.Entities == nil {
		return nil
	}

	if server.Properties.BootVolume != nil {
		volume, apiResponse, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, datacenterId, d.Id(), *server.Properties.BootVolume.Id).Execute()
		logApiRequestTime(apiResponse)
		if err == nil {
			var volumes []interface{}
			entry := SetVolumeProperties(volume)
			userData := d.Get("volume.0.user_data")
			entry["user_data"] = userData

			backupUnit := d.Get("volume.0.backup_unit_id")
			entry["backup_unit_id"] = backupUnit
			volumes = append(volumes, entry)
			if err := d.Set("volume", volumes); err != nil {
				return fmt.Errorf("error setting volume %w", err)
			}
		}
	}

	_, primaryNicOk := d.GetOk("primary_nic")
	// take nic and firewall from schema if set is used in resource read, else take it from entities
	var nicId string
	if primaryNicOk {
		nicId = d.Get("primary_nic").(string)
	} else if server.Entities.Nics != nil && server.Entities.Nics.Items != nil && len(*server.Entities.Nics.Items) > 0 && (*server.Entities.Nics.Items)[0].Id != nil {
		lastNicIndex := len(*server.Entities.Nics.Items) - 1
		if (*server.Entities.Nics.Items)[lastNicIndex].Id != nil {
			// this might be a terraformer import, so primary_nic might not be set
			// if no nics found until now, get the last one from the list of nics.
			nicId = *(*server.Entities.Nics.Items)[lastNicIndex].Id
		}
	}
	nic, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, datacenterId, d.Id(), nicId).Depth(2).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return err
	}

	if nic.Properties != nil {
		nicEntry := SetNetworkProperties(nic)
		nicEntry["id"] = *nic.Id
		fs := FirewallService{client: client, schemaData: d}
		firewallRuleIds := d.Get("firewallrule_ids").([]interface{})
		var fwRulesEntries []map[string]interface{}

		for _, id := range firewallRuleIds {
			firewall, apiResponse, err := fs.firewallFindById(ctx, datacenterId, d.Id(), nicId, id.(string))
			if err != nil {
				if !apiResponse.HttpNotFound() {
					return fmt.Errorf("error, while searching for firewall rule with id %s (%w)", id.(string), err)
				}
			}
			if firewall == nil {
				continue
			}
			if firewall.Properties != nil && firewall.Properties.Name != nil {
				log.Printf("[DEBUG] found firewall rule with name %s", *firewall.Properties.Name)
			}
			firewallEntry := SetFirewallProperties(*firewall)
			if firewallEntry != nil && len(firewallEntry) != 0 {
				fwRulesEntries = append(fwRulesEntries, firewallEntry)
			}
		}
		if fwRulesEntries != nil && len(fwRulesEntries) > 0 {
			nicEntry["firewall"] = fwRulesEntries
		}
		if len(nicEntry) != 0 {
			nics := []map[string]interface{}{nicEntry}

			if err := d.Set("nic", nics); err != nil {
				return fmt.Errorf("error settings nics %w", err)
			}
		}
	}

	// Labels logic
	ls := LabelsService{ctx: ctx, client: client}
	labels, err := ls.datacentersServersLabelsGet(datacenterId, d.Id(), false)
	if err != nil {
		return err
	}
	if err := d.Set("label", labels); err != nil {
		return err
	}

	//if token != nil {
	//	if err := d.Set("token", *token.Token); err != nil {
	//		return err
	//	}
	//}
	return nil
}
