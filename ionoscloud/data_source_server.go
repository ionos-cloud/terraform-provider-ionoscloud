package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerRead,
		Schema: map[string]*schema.Schema{
			"template_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datacenter_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
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
						"image_aliases": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cloud_init": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
						"type": {
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
					},
				},
			},
			"nics": {
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
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"dhcp": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"lan": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"firewall_active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"firewall_type": {
							Type:     schema.TypeString,
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
						"firewall_rules": {
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
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_mac": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"icmp_code": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"icmp_type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"port_range_start": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"port_range_end": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
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

func boolOrDefault(p *bool, d bool) bool {
	if p != nil {
		return *p
	}
	return d
}

func stringOrDefault(s *string, d string) string {
	if s != nil {
		return *s
	}
	return d
}

func int32OrDefault(i *int32, d int32) int32 {
	if i != nil {
		return *i
	}
	return d
}

func int64OrDefault(i *int64, d int64) int64 {
	if i != nil {
		return *i
	}
	return d
}

func float32OrDefault(f *float32, d float32) float32 {
	if f != nil {
		return *f
	}
	return d
}

func setServerData(d *schema.ResourceData, server *ionoscloud.Server, token *ionoscloud.Token) error {

	if server.Id != nil {
		d.SetId(*server.Id)
		if err := d.Set("id", *server.Id); err != nil {
			return err
		}
	}

	if server.Properties != nil {
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
			entry := make(map[string]interface{})

			entry["id"] = stringOrDefault(image.Id, "")
			entry["name"] = stringOrDefault(image.Properties.Name, "")
			entry["description"] = stringOrDefault(image.Properties.Description, "")
			entry["location"] = stringOrDefault(image.Properties.Location, "")
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
			entry["licence_type"] = stringOrDefault(image.Properties.LicenceType, "")
			entry["image_type"] = stringOrDefault(image.Properties.ImageType, "")
			entry["public"] = boolOrDefault(image.Properties.Public, false)

			if image.Properties.ImageAliases != nil {
				var imageAliases []interface{}
				for _, imageAlias := range *image.Properties.ImageAliases {
					imageAliases = append(imageAliases, imageAlias)
				}
				entry["image_aliases"] = imageAliases
			}

			entry["cloud_init"] = stringOrDefault(image.Properties.CloudInit, "")

			cdroms = append(cdroms, entry)
		}
		if err := d.Set("cdroms", cdroms); err != nil {
			return err
		}
	}

	var volumes []interface{}
	if server.Entities.Volumes != nil && server.Entities.Volumes.Items != nil && len(*server.Entities.Volumes.Items) > 0 {
		for _, volume := range *server.Entities.Volumes.Items {
			entry := make(map[string]interface{})

			entry["id"] = stringOrDefault(volume.Id, "")
			entry["name"] = stringOrDefault(volume.Properties.Name, "")
			entry["type"] = stringOrDefault(volume.Properties.Type, "")
			entry["size"] = float32OrDefault(volume.Properties.Size, 0)
			entry["availability_zone"] = stringOrDefault(volume.Properties.AvailabilityZone, "")
			entry["image_name"] = stringOrDefault(volume.Properties.Image, "")
			entry["image_password"] = stringOrDefault(volume.Properties.ImagePassword, "")

			if volume.Properties.SshKeys != nil && len(*volume.Properties.SshKeys) > 0 {
				var sshKeys []interface{}
				for _, sshKey := range *volume.Properties.SshKeys {
					sshKeys = append(sshKeys, sshKey)
				}
				entry["ssh_keys"] = sshKeys
			}

			entry["bus"] = stringOrDefault(volume.Properties.Bus, "")
			entry["licence_type"] = stringOrDefault(volume.Properties.LicenceType, "")
			entry["cpu_hot_plug"] = boolOrDefault(volume.Properties.CpuHotPlug, true)
			entry["ram_hot_plug"] = boolOrDefault(volume.Properties.RamHotPlug, true)
			entry["nic_hot_plug"] = boolOrDefault(volume.Properties.NicHotPlug, true)
			entry["nic_hot_unplug"] = boolOrDefault(volume.Properties.NicHotUnplug, true)
			entry["disc_virtio_hot_plug"] = boolOrDefault(volume.Properties.DiscVirtioHotPlug, true)
			entry["disc_virtio_hot_unplug"] = boolOrDefault(volume.Properties.DiscVirtioHotUnplug, true)
			entry["device_number"] = int64OrDefault(volume.Properties.DeviceNumber, 0)
			entry["pci_slot"] = int32OrDefault(volume.Properties.PciSlot, 0)
			entry["backup_unit_id"] = stringOrDefault(volume.Properties.BackupunitId, "")
			entry["user_data"] = stringOrDefault(volume.Properties.UserData, "")

			volumes = append(volumes, entry)
		}

		if err := d.Set("volumes", volumes); err != nil {
			return err
		}
	}

	var nics []interface{}
	if server.Entities.Nics != nil && server.Entities.Nics.Items != nil && len(*server.Entities.Nics.Items) > 0 {
		for _, nic := range *server.Entities.Nics.Items {
			entry := make(map[string]interface{})

			entry["id"] = stringOrDefault(nic.Id, "")
			entry["name"] = stringOrDefault(nic.Properties.Name, "")
			entry["mac"] = stringOrDefault(nic.Properties.Mac, "")

			if nic.Properties.Ips != nil {
				var ips []interface{}
				for _, ip := range *nic.Properties.Ips {
					ips = append(ips, ip)
				}
				entry["ips"] = ips
			}

			entry["dhcp"] = boolOrDefault(nic.Properties.Dhcp, false)
			entry["lan"] = int32OrDefault(nic.Properties.Lan, 0)
			entry["firewall_active"] = boolOrDefault(nic.Properties.FirewallActive, false)
			entry["firewall_type"] = stringOrDefault(nic.Properties.FirewallType, "")
			entry["device_number"] = int32OrDefault(nic.Properties.DeviceNumber, 0)
			entry["pci_slot"] = int32OrDefault(nic.Properties.PciSlot, 0)

			if nic.Entities != nil && nic.Entities.Firewallrules != nil && nic.Entities.Firewallrules.Items != nil {
				var firewallRules []interface{}
				for _, rule := range *nic.Entities.Firewallrules.Items {
					ruleEntry := make(map[string]interface{})

					ruleEntry["id"] = stringOrDefault(rule.Id, "")
					ruleEntry["name"] = stringOrDefault(rule.Properties.Name, "")
					ruleEntry["protocol"] = stringOrDefault(rule.Properties.Protocol, "")
					ruleEntry["source_mac"] = stringOrDefault(rule.Properties.SourceMac, "")
					ruleEntry["source_ip"] = stringOrDefault(rule.Properties.SourceIp, "")
					ruleEntry["target_ip"] = stringOrDefault(rule.Properties.TargetIp, "")
					ruleEntry["icmp_code"] = int32OrDefault(rule.Properties.IcmpCode, 0)
					ruleEntry["icmp_type"] = int32OrDefault(rule.Properties.IcmpType, 0)
					ruleEntry["port_range_start"] = int32OrDefault(rule.Properties.PortRangeStart, 0)
					ruleEntry["port_range_end"] = int32OrDefault(rule.Properties.PortRangeEnd, 0)
					ruleEntry["type"] = stringOrDefault(rule.Properties.Type, "")
					firewallRules = append(firewallRules, ruleEntry)
				}
				entry["firewall_rules"] = firewallRules
			}

			nics = append(nics, entry)
		}

		if err := d.Set("nics", nics); err != nil {
			return err
		}
	}

	if token != nil {
		if err := d.Set("token", *token.Token); err != nil {
			return err
		}
	}

	return nil
}

func dataSourceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return errors.New("no datacenter_id was specified")
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the server id or name")
	}
	var server ionoscloud.Server
	var err error
	var apiResponse *ionoscloud.APIResponse

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if idOk {
		/* search by ID */
		server, apiResponse, err = client.ServersApi.DatacentersServersFindById(ctx, datacenterId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the server with ID %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var servers ionoscloud.Servers
		servers, apiResponse, err := client.ServersApi.DatacentersServersGet(ctx, datacenterId.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching servers: %s", err.Error())
		}

		found := false
		if servers.Items != nil {
			for _, s := range *servers.Items {
				if s.Properties.Name != nil && *s.Properties.Name == name.(string) {
					/* server found */
					server, apiResponse, err = client.ServersApi.DatacentersServersFindById(ctx, datacenterId.(string), *s.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return fmt.Errorf("an error occurred while fetching the server with ID %s: %s", *s.Id, err)
					}
					found = true
					break
				}
			}
		}

		if !found {
			return errors.New("server not found")
		}

	}

	var token = ionoscloud.Token{}

	if &server != nil && server.Id != nil {
		token, apiResponse, err = client.ServersApi.DatacentersServersTokenGet(ctx, datacenterId.(string), *server.Id).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("an error occurred while fetching the server token %s: %s", *server.Id, err)
		}
	}

	if err = setServerData(d, &server, &token); err != nil {
		return err
	}

	return nil
}
