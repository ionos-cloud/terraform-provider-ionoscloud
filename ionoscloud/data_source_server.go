package ionoscloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func dataSourceServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerRead,
		Schema: map[string]*schema.Schema{
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
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
						"image_aliases": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"public": {
							Type:     schema.TypeBool,
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
						"image": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_alias": {
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
						"device_number": {
							Type:     schema.TypeInt,
							Computed: true,
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
						"nat": {
							Type:     schema.TypeBool,
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

func intOrDefault(i *int, d int) int {
	if i != nil {
		return *i
	}
	return d
}

func setServerData(d *schema.ResourceData, server *profitbricks.Server) error {
	d.SetId(server.ID)
	if err := d.Set("id", server.ID); err != nil {
		return err
	}

	if err := d.Set("name", server.Properties.Name); err != nil {
		return err
	}
	if err := d.Set("cores", server.Properties.Cores); err != nil {
		return err
	}
	if err := d.Set("ram", server.Properties.RAM); err != nil {
		return err
	}
	if err := d.Set("availability_zone", server.Properties.AvailabilityZone); err != nil {
		return err
	}
	if err := d.Set("vm_state", server.Properties.VMState); err != nil {
		return err
	}
	if err := d.Set("cpu_family", server.Properties.CPUFamily); err != nil {
		return err
	}
	if server.Properties.BootCdrom != nil {
		if err := d.Set("boot_cdrom", server.Properties.BootCdrom.ID); err != nil {
			return err
		}
	}
	if server.Properties.BootVolume != nil {
		if err := d.Set("boot_volume", server.Properties.BootVolume.ID); err != nil {
			return err
		}
	}

	if server.Entities == nil {
		return nil
	}

	var cdroms []interface{}
	if server.Entities.Cdroms != nil {
		cdroms = make([]interface{}, len(server.Entities.Cdroms.Items), len(server.Entities.Cdroms.Items))
		for i, image := range server.Entities.Cdroms.Items {
			entry := make(map[string]interface{})

			entry["id"] = image.ID
			entry["name"] = image.Properties.Name
			entry["description"] = image.Properties.Description
			entry["location"] = image.Properties.Location
			entry["size"] = image.Properties.Size
			entry["cpu_hot_plug"] = image.Properties.CPUHotPlug
			entry["cpu_hot_unplug"] = image.Properties.CPUHotUnplug
			entry["ram_hot_plug"] = image.Properties.RAMHotPlug
			entry["ram_hot_unplug"] = image.Properties.RAMHotUnplug
			entry["nic_hot_plug"] = image.Properties.NicHotPlug
			entry["nic_hot_unplug"] = image.Properties.NicHotUnplug
			entry["disc_virtio_hot_plug"] = image.Properties.DiscVirtioHotPlug
			entry["disc_virtio_hot_unplug"] = image.Properties.DiscVirtioHotUnplug
			entry["disc_scsi_hot_plug"] = image.Properties.DiscScsiHotPlug
			entry["disc_scsi_hot_unplug"] = image.Properties.DiscScsiHotUnplug
			entry["licence_type"] = image.Properties.LicenceType
			entry["image_type"] = image.Properties.ImageType
			entry["image_alias"] = image.Properties.ImageAliases
			entry["public"] = image.Properties.Public

			cdroms[i] = entry
		}
	}

	if err := d.Set("cdroms", cdroms); err != nil {
		return err
	}

	var volumes = make([]interface{}, 0)
	if server.Entities.Volumes != nil {
		volumes = make([]interface{}, len(server.Entities.Volumes.Items), len(server.Entities.Volumes.Items))
		for i, volume := range server.Entities.Volumes.Items {
			entry := make(map[string]interface{})

			entry["id"] = volume.ID
			entry["name"] = volume.Properties.Name
			entry["type"] = volume.Properties.Type
			entry["size"] = volume.Properties.Size
			entry["availability_zone"] = volume.Properties.AvailabilityZone
			entry["image"] = volume.Properties.Image
			entry["image_alias"] = volume.Properties.ImageAlias
			entry["image_password"] = volume.Properties.ImagePassword

			sshKeys := make([]interface{}, len(volume.Properties.SSHKeys), len(volume.Properties.SSHKeys))
			for j, sshKey := range volume.Properties.SSHKeys {
				sshKeys[j] = sshKey
			}
			entry["ssh_keys"] = sshKeys

			entry["bus"] = volume.Properties.Bus
			entry["licence_type"] = volume.Properties.LicenceType
			entry["cpu_hot_plug"] = volume.Properties.CPUHotPlug
			entry["cpu_hot_unplug"] = volume.Properties.CPUHotUnplug
			entry["ram_hot_plug"] = volume.Properties.RAMHotPlug
			entry["ram_hot_unplug"] = volume.Properties.RAMHotUnplug
			entry["nic_hot_plug"] = volume.Properties.NicHotPlug
			entry["nic_hot_unplug"] = volume.Properties.NicHotUnplug
			entry["disc_virtio_hot_plug"] = volume.Properties.DiscVirtioHotPlug
			entry["disc_virtio_hot_unplug"] = volume.Properties.DiscVirtioHotUnplug
			entry["disc_scsi_hot_plug"] = volume.Properties.DiscScsiHotPlug
			entry["disc_scsi_hot_unplug"] = volume.Properties.DiscScsiHotUnplug
			entry["device_number"] = volume.Properties.DeviceNumber

			volumes[i] = entry
		}
	}

	if err := d.Set("volumes", volumes); err != nil {
		return err
	}

	var nics = make([]interface{}, 0)
	if server.Entities.Volumes != nil {
		nics = make([]interface{}, len(server.Entities.Nics.Items), len(server.Entities.Nics.Items))
		for k, nic := range server.Entities.Nics.Items {
			entry := make(map[string]interface{})

			entry["id"] = nic.ID
			entry["name"] = nic.Properties.Name
			entry["mac"] = nic.Properties.Mac

			ips := make([]interface{}, len(nic.Properties.Ips))
			for idx, ip := range nic.Properties.Ips {
				ips[idx] = ip
			}
			entry["ips"] = ips

			entry["dhcp"] = boolOrDefault(nic.Properties.Dhcp, false)
			entry["lan"] = nic.Properties.Lan
			entry["firewall_active"] = boolOrDefault(nic.Properties.FirewallActive, false)
			entry["nat"] = boolOrDefault(nic.Properties.Nat, false)

			firewallRules := make([]interface{}, 0)
			if nic.Entities != nil && nic.Entities.FirewallRules != nil {
				firewallRules = make([]interface{}, len(nic.Entities.FirewallRules.Items))
				for idx, rule := range nic.Entities.FirewallRules.Items {
					ruleEntry := make(map[string]interface{})

					ruleEntry["id"] = rule.ID
					ruleEntry["name"] = rule.Properties.Name
					ruleEntry["protocol"] = rule.Properties.Protocol
					ruleEntry["source_mac"] = stringOrDefault(rule.Properties.SourceMac, "")
					ruleEntry["source_ip"] = stringOrDefault(rule.Properties.SourceIP, "")
					ruleEntry["target_ip"] = stringOrDefault(rule.Properties.TargetIP, "")
					ruleEntry["icmp_code"] = intOrDefault(rule.Properties.IcmpCode, 0)
					ruleEntry["icmp_type"] = intOrDefault(rule.Properties.IcmpType, 0)
					firewallRules[idx] = ruleEntry
				}
			}
			entry["firewall_rules"] = firewallRules

			nics[k] = entry
		}
	}

	if err := d.Set("nics", nics); err != nil {
		return err
	}

	return nil
}

func dataSourceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

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
	var server *profitbricks.Server
	var err error

	if idOk {
		/* search by ID */
		server, err = client.GetServer(datacenterId.(string), id.(string))
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the server with ID %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var servers *profitbricks.Servers
		servers, err := client.ListServers(datacenterId.(string))
		if err != nil {
			return fmt.Errorf("an error occurred while fetching servers: %s", err.Error())
		}

		for _, s := range servers.Items {
			if s.Properties.Name == name.(string) {
				/* server found */
				server, err = client.GetServer(datacenterId.(string), s.ID)
				if err != nil {
					return fmt.Errorf("an error occurred while fetching the server with ID %s: %s", s.ID, err)
				}
				break
			}
		}
	}

	if server == nil {
		return errors.New("server not found")
	}

	if err = setServerData(d, server); err != nil {
		return err
	}

	return nil
}
