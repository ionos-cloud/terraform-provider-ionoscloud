package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceCubeServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCubeServerRead,
		Schema: map[string]*schema.Schema{
			"template_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
						"type": {
							Type:     schema.TypeString,
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func setCubeServerData(d *schema.ResourceData, server *ionoscloud.Server, token *ionoscloud.Token) error {

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
		if server.Entities != nil && server.Entities.Volumes != nil && server.Entities.Volumes.Items != nil && len(*server.Entities.Volumes.Items) > 0 &&
			(*server.Entities.Volumes.Items)[0].Properties.Image != nil {
			if err := d.Set("boot_image", *(*server.Entities.Volumes.Items)[0].Properties.Image); err != nil {
				return err
			}
		}
	}

	if server.Entities == nil {
		return nil
	}

	if server.Entities.Cdroms != nil && server.Entities.Cdroms.Items != nil && len(*server.Entities.Cdroms.Items) > 0 {
		cdroms := setServerCDRoms(server.Entities.Cdroms.Items)
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
			entry["boot_server"] = stringOrDefault(volume.Properties.BootServer, "")

			volumes = append(volumes, entry)
		}

		if err := d.Set("volumes", volumes); err != nil {
			return err
		}
	}

	var nicsIntf []interface{}
	if server.Entities != nil {
		if server.Entities.Nics != nil && server.Entities.Nics.Items != nil {
			nicItems := server.Entities.Nics.Items
			if nicItems != nil && len(*nicItems) > 0 {
				var nics []interface{}
				for _, nic := range *server.Entities.Nics.Items {
					nicMap := SetNetworkProperties(nic)
					fw := setFirewallRules(nic)
					nicMap["firewall_rules"] = fw
					utils.SetPropWithNilCheck(nicMap, "id", nic.Id)
					nics = append(nics, nicMap)
				}
				nicsIntf = nics
			}
		}
	}

	if token != nil {
		if err := d.Set("token", *token.Token); err != nil {
			return err
		}
	}

	if nicsIntf == nil || len(nicsIntf) == 0 {
		return fmt.Errorf("no nics found for criteria, please check your filter configuration")
	}
	err := d.Set("nics", nicsIntf)
	if err != nil {
		return fmt.Errorf("error while setting nics: %w", err)
	}

	return nil
}

func dataSourceCubeServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return diag.FromErr(errors.New("no datacenter_id was specified"))
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the server id or name"))
	}
	var server ionoscloud.Server
	var err error
	var apiResponse *shared.APIResponse

	if idOk {
		/* search by ID */
		server, apiResponse, err = client.ServersApi.DatacentersServersFindById(ctx, datacenterId.(string), id.(string)).Depth(5).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the server with ID %s: %w", id.(string), err))
		}
	} else {
		/* search by name */
		servers, apiResponse, err := client.ServersApi.DatacentersServersGet(ctx, datacenterId.(string)).Depth(5).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching servers: %w", err))
		}

		var results []ionoscloud.Server

		if servers.Items != nil {
			for _, s := range *servers.Items {
				if s.Properties != nil && s.Properties.Name != nil && *s.Properties.Name == name.(string) {
					/* server found */
					server, apiResponse, err = client.ServersApi.DatacentersServersFindById(ctx, datacenterId.(string), *s.Id).Depth(4).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diag.FromErr(fmt.Errorf("an error occurred while fetching the server with ID %s: %w", *s.Id, err))
					}
					results = append(results, server)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no server found with the specified criteria: name = %s", name.(string)))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one server found with the specified criteria: name = %s", name.(string)))
		} else {
			server = results[0]
		}

	}

	var token = ionoscloud.Token{}

	if &server != nil && server.Id != nil {
		token, apiResponse, err = client.ServersApi.DatacentersServersTokenGet(ctx, datacenterId.(string), *server.Id).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the server token %s: %w", *server.Id, err))
		}
	}

	if err = setCubeServerData(d, &server, &token); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
