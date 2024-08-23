package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/iancoleman/strcase"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapinic"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServersRead,
		Schema: map[string]*schema.Schema{
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"filter": dataSourceFiltersSchema(),
			"servers": {
				Type:        schema.TypeList,
				Description: "list of servers",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:        schema.TypeString,
							Description: "The unique ID of the server.",
							Computed:    true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"boot_cdrom": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"boot_volume": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"boot_image": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cdroms": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     cdromsServerDSResource,
						},
						"volumes": {
							Type:     schema.TypeList,
							Computed: true,
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
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"availability_zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_password": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ssh_keys": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"bus": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"licence_type": {
										Type:     schema.TypeString,
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
									"pci_slot": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"backup_unit_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_data": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"boot_server": {
										Type:        schema.TypeString,
										Description: "The UUID of the attached server.",
										Computed:    true,
									},
								},
							},
						},
						"nics": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     nicServerDSResource,
						},
						"labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     labelDataSource,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func dataSourceServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return diag.FromErr(errors.New("no datacenter_id was specified"))
	}
	req := client.ServersApi.DatacentersServersGet(ctx, datacenterId.(string)).Depth(5)
	filters, filtersOk := d.GetOk("filter")
	if filtersOk {
		for _, v := range filters.(*schema.Set).List() {
			filter := v.(map[string]interface{})
			// we want to convert for example cpu_family to cpuFamily
			name := strcase.ToLowerCamel(filter["name"].(string))
			value := filter["value"].(string)
			req = req.Filter(name, value)
			log.Printf("[INFO] Adding filter with name %s and value %s \n", name, value)
		}
	}
	var err error
	var apiResponse *ionoscloud.APIResponse

	/* search by whatever filter is set above */
	servers, apiResponse, err := req.Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching servers: %w", err))
	}
	serverEntry := make(map[string]interface{})
	var serversIntf []interface{}

	if d.Id() == "" {
		d.SetId(datacenterId.(string))
	}
	for _, server := range *servers.Items {
		serverEntry = SetServerProperties(server)
		utils.SetPropWithNilCheck(serverEntry, "id", server.Id)
		// todo: Add token?
		if server.Entities != nil {
			if server.Entities.Nics != nil && server.Entities.Nics.Items != nil {
				nicItems := server.Entities.Nics.Items
				if nicItems != nil && len(*nicItems) > 0 {
					var nics []interface{}
					for _, nic := range *server.Entities.Nics.Items {
						nicMap := cloudapinic.SetNetworkProperties(nic)
						fw := setFirewallRules(nic)
						nicMap["firewall_rules"] = fw
						utils.SetPropWithNilCheck(nicMap, "id", nic.Id)
						nics = append(nics, nicMap)
					}
					if len(nics) > 0 {
						serverEntry["nics"] = nics
					}
				}
			}
			if server.Entities.Volumes != nil && server.Entities.Volumes.Items != nil {
				volumes := setVolumePropertiesToSlice(*server.Entities.Volumes.Items)
				if volumes != nil && len(volumes) > 0 {
					serverEntry["volumes"] = volumes
				}
			}
			if server.Entities.Cdroms != nil {
				if server.Entities.Cdroms.Items != nil && len(*server.Entities.Cdroms.Items) > 0 {
					cdroms := setServerCDRoms(server.Entities.Cdroms.Items)
					if cdroms != nil && len(cdroms) > 0 {
						serverEntry["cdroms"] = cdroms
					}
				}
			}

			if server.Id == nil {
				return diag.FromErr(fmt.Errorf("expected a valid server ID from the API but received nil instead"))
			}
			// Labels logic
			ls := LabelsService{ctx: ctx, client: client}
			labels, err := ls.datacentersServersLabelsGet(datacenterId.(string), *server.Id, true)
			if err != nil {
				return diag.FromErr(err)
			}
			serverEntry["labels"] = labels
		}
		serversIntf = append(serversIntf, serverEntry)
	}
	if serversIntf == nil || len(serversIntf) == 0 {
		return diag.FromErr(fmt.Errorf("no servers found for criteria, please check your filter configuration"))
	}
	err = d.Set("servers", &serversIntf)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error while setting servers: %w", err))
		return diags
	}

	return nil
}

// setVolumePropertiesToSlice returns a slice of volumes
func setVolumePropertiesToSlice(volumesList []ionoscloud.Volume) []interface{} {
	var volumes []interface{}
	if volumesList != nil && len(volumesList) > 0 {
		for _, volume := range volumesList {
			volumeItemMap := SetVolumeProperties(volume)
			utils.SetPropWithNilCheck(volumeItemMap, "id", volume.Id)
			volumes = append(volumes, volumeItemMap)
		}
	}
	return volumes
}
func setFirewallRules(nic ionoscloud.Nic) []interface{} {
	var firewallRules []interface{}
	if nic.Entities != nil && nic.Entities.Firewallrules != nil && nic.Entities.Firewallrules.Items != nil {
		for _, rule := range *nic.Entities.Firewallrules.Items {
			ruleEntry := setFirewallRuleProperties(rule)
			firewallRules = append(firewallRules, ruleEntry)
		}
	}
	return firewallRules
}

func setFirewallRuleProperties(rule ionoscloud.FirewallRule) map[string]interface{} {
	ruleEntry := make(map[string]interface{})
	ruleEntry["id"] = mongo.ToValueDefault(rule.Id)
	if rule.Properties != nil {
		ruleEntry["name"] = mongo.ToValueDefault(rule.Properties.Name)
		ruleEntry["protocol"] = mongo.ToValueDefault(rule.Properties.Protocol)
		ruleEntry["source_mac"] = mongo.ToValueDefault(rule.Properties.SourceMac)
		ruleEntry["source_ip"] = mongo.ToValueDefault(rule.Properties.SourceIp)
		ruleEntry["target_ip"] = mongo.ToValueDefault(rule.Properties.TargetIp)
		ruleEntry["icmp_code"] = int32OrDefault(rule.Properties.IcmpCode, 0)
		ruleEntry["icmp_type"] = int32OrDefault(rule.Properties.IcmpType, 0)
		ruleEntry["port_range_start"] = int32OrDefault(rule.Properties.PortRangeStart, 0)
		ruleEntry["port_range_end"] = int32OrDefault(rule.Properties.PortRangeEnd, 0)
		ruleEntry["type"] = mongo.ToValueDefault(rule.Properties.Type)
	}
	return ruleEntry
}

func SetServerProperties(server ionoscloud.Server) map[string]interface{} {
	serverMap := map[string]interface{}{}
	if server.Properties != nil {
		utils.SetPropWithNilCheck(serverMap, "template_uuid", server.Properties.TemplateUuid)
		utils.SetPropWithNilCheck(serverMap, "name", server.Properties.Name)
		utils.SetPropWithNilCheck(serverMap, "cores", server.Properties.Cores)
		utils.SetPropWithNilCheck(serverMap, "ram", server.Properties.Ram)
		utils.SetPropWithNilCheck(serverMap, "availability_zone", server.Properties.AvailabilityZone)
		utils.SetPropWithNilCheck(serverMap, "cpu_family", server.Properties.CpuFamily)
		utils.SetPropWithNilCheck(serverMap, "type", server.Properties.Type)
		if server.Properties.BootCdrom != nil && server.Properties.BootCdrom.Id != nil {
			utils.SetPropWithNilCheck(serverMap, "boot_cdrom", *server.Properties.BootCdrom.Id)
		}

		if server.Properties.BootVolume != nil && server.Properties.BootVolume.Id != nil {
			utils.SetPropWithNilCheck(serverMap, "boot_volume", *server.Properties.BootVolume.Id)

		}
		if server.Entities != nil && server.Entities.Volumes != nil && server.Entities.Volumes.Items != nil && len(*server.Entities.Volumes.Items) > 0 &&
			(*server.Entities.Volumes.Items)[0].Properties.Image != nil {
			utils.SetPropWithNilCheck(serverMap, "boot_image", (*server.Entities.Volumes.Items)[0].Properties.Image)
		}
	}
	return serverMap
}

func setServerCDRoms(images *[]ionoscloud.Image) []interface{} {
	var cdroms []interface{}
	for _, image := range *images {
		entry := make(map[string]interface{})

		entry["id"] = mongo.ToValueDefault(image.Id)
		entry["name"] = mongo.ToValueDefault(image.Properties.Name)
		entry["description"] = mongo.ToValueDefault(image.Properties.Description)
		entry["location"] = mongo.ToValueDefault(image.Properties.Location)
		entry["size"] = float32OrDefault(image.Properties.Size, 0)
		entry["cpu_hot_plug"] = boolOrDefault(image.Properties.CpuHotPlug, true)
		entry["cpu_hot_unplug"] = boolOrDefault(image.Properties.CpuHotUnplug, true)
		entry["ram_hot_plug"] = boolOrDefault(image.Properties.RamHotPlug, true)
		entry["ram_hot_unplug"] = boolOrDefault(image.Properties.RamHotUnplug, true)
		entry["nic_hot_plug"] = boolOrDefault(image.Properties.NicHotPlug, true)
		entry["nic_hot_unplug"] = boolOrDefault(image.Properties.NicHotUnplug, true)
		entry["disc_virtio_hot_plug"] = boolOrDefault(image.Properties.DiscVirtioHotPlug, true)
		entry["disc_virtio_hot_unplug"] = boolOrDefault(image.Properties.DiscVirtioHotUnplug, true)
		entry["disc_scsi_hot_plug"] = boolOrDefault(image.Properties.DiscScsiHotPlug, true)
		entry["disc_scsi_hot_unplug"] = boolOrDefault(image.Properties.DiscScsiHotUnplug, true)
		entry["licence_type"] = mongo.ToValueDefault(image.Properties.LicenceType)
		entry["image_type"] = mongo.ToValueDefault(image.Properties.ImageType)
		entry["public"] = boolOrDefault(image.Properties.Public, false)

		if image.Properties.ImageAliases != nil {
			var imageAliases []interface{}
			for _, imageAlias := range *image.Properties.ImageAliases {
				imageAliases = append(imageAliases, imageAlias)
			}
			entry["image_aliases"] = imageAliases
		}

		entry["cloud_init"] = mongo.ToValueDefault(image.Properties.CloudInit)

		cdroms = append(cdroms, entry)
	}
	return cdroms
}
