package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/serverutil"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapinic"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapiserver"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/nsg"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func resourceGPUServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGpuServerCreate,
		ReadContext:   resourceCubeServerRead,
		UpdateContext: resourceCubeServerUpdate,
		DeleteContext: serverutil.ResourceCommonServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCubeServerImport,
		},
		CustomizeDiff: checkServerImmutableFields,

		Schema: map[string]*schema.Schema{
			"template_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"hostname": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "The hostname of the resource. Allowed characters are a-z, 0-9 and - (minus). Hostname should not start with minus and should not be longer than 63 characters. If no value provided explicitly, it will be populated with the name of the server",
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(validation.StringIsNotWhiteSpace, validation.StringLenBetween(1, 63))),
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2"}, true)),
			},
			"boot_volume": {
				Type:     schema.TypeString,
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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"image_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ssh_key_path": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"security_groups_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The list of Security Group IDs for the server",
			},
			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
						"licence_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
							ForceNew:         true,
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
							ForceNew:    true,
						},
						"user_data": {
							Type:        schema.TypeString,
							Description: "The cloud-init configuration for the volume as base64 encoded string. The property is immutable and is only allowed to be set on a new volume creation. It is mandatory to provide either 'public image' or 'imageAlias' that has cloud-init compatibility in conjunction with this property.",
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
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
						"expose_serial": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							Description: "If set to `true` will expose the serial id of the disk attached to the server. " +
								"If set to `false` will not expose the serial id. Some operating systems or software solutions require the serial id to be exposed to work properly. " +
								"Exposing the serial can influence licensed software (e.g. Windows) behavior",
						},
						"require_legacy_bios": {
							Type:        schema.TypeBool,
							Description: "Indicates if the image requires the legacy BIOS for compatibility or specific needs.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"vm_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Sets the power state of the gpu server. Possible values: `RUNNING` or `SUSPENDED`.",
				// SUSPENDING a RUNNING GPU server sometimes results in a PAUSED state.
				// The API doesn't really support this yet fully
				// Allow users to set it to PAUSED if they notice their server is in that state after suspending it
				// to prevent Terraform from trying to start it again on the next apply
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{constant.VMStateStart, constant.CubeVMStateStop, "PAUSED"}, true)),
			},
			"nic": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: serverutil.SchemaNicElem,
				},
			},
			"inline_volume_ids": {
				Type:        schema.TypeList,
				Description: "A list that contains the IDs for the volumes defined inside the gpu server resource.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_replace": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When set to true, allows the update of immutable fields by destroying and re-creating the resource.",
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

//nolint:gocyclo
func resourceGpuServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	server := ionoscloud.Server{
		Properties: &ionoscloud.ServerProperties{},
	}

	var image, imageAlias string

	dcID := d.Get("datacenter_id").(string)

	serverName := d.Get("name").(string)
	server.Properties.Name = &serverName

	templateUuid := d.Get("template_uuid").(string)
	server.Properties.TemplateUuid = &templateUuid

	if v, ok := d.GetOk("availability_zone"); ok {
		vStr := v.(string)
		server.Properties.AvailabilityZone = &vStr
	}

	serverType := constant.GpuType
	server.Properties.Type = &serverType

	if v, ok := d.GetOk("hostname"); ok {
		if v.(string) != "" {
			vStr := v.(string)
			server.Properties.Hostname = &vStr
		}
	}

	if _, ok := d.GetOk("boot_volume"); ok {
		return utils.ToDiags(d, "boot_volume argument can be set only in update requests", nil)
	}

	var err error
	var volume *ionoscloud.VolumeProperties
	if _, ok := d.GetOk("volume"); ok {
		volume, err = getVolumeData(d, "volume.0.", constant.GpuType)
		if err != nil {
			return utils.ToDiags(d, err.Error(), nil)
		}
		image, imageAlias, err = getImage(ctx, client, d, *volume)
		if err != nil {
			return utils.ToDiags(d, err.Error(), nil)
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
		if backupUnitID, ok := d.GetOk("volume.0.backup_unit_id"); ok {
			if utils.IsValidUUID(backupUnitID.(string)) {
				if image == "" && imageAlias == "" {
					return utils.ToDiags(d, "it is mandatory to provide either public image or imageAlias in conjunction with backup unit id property", nil)
				}
				backupUnitID := backupUnitID.(string)
				volume.BackupunitId = &backupUnitID
			}
		}
		if userData, ok := d.GetOk("volume.0.user_data"); ok {
			if image == "" && imageAlias == "" {
				return utils.ToDiags(d, "it is mandatory to provide either public image or imageAlias that has cloud-init compatibility in conjunction with backup unit id property ", nil)
			}
			userData := userData.(string)
			volume.UserData = &userData
		}
		server.Entities = &ionoscloud.ServerEntities{
			Volumes: &ionoscloud.AttachedVolumes{
				Items: &[]ionoscloud.Volume{
					{
						Properties: volume,
					},
				},
			},
		}
	}

	var primaryNic *ionoscloud.Nic
	if _, ok := d.GetOk("nic"); ok {
		nic, err := cloudapinic.GetNicFromSchemaCreate(d, "nic.0.")
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("gpu error occurred while getting nic from schema: %s", err), nil)
		}

		server.Entities.Nics = &ionoscloud.Nics{
			Items: &[]ionoscloud.Nic{
				nic,
			},
		}
		primaryNic = &(*server.Entities.Nics.Items)[0]
		log.Printf("[DEBUG] dhcp nic after %t", *nic.Properties.Dhcp)
		log.Printf("[DEBUG] primaryNic dhcp %t", *primaryNic.Properties.Dhcp)

		var firewall ionoscloud.FirewallRule
		if _, ok := d.GetOk("nic.0.firewall"); ok {
			var diags diag.Diagnostics
			firewall, diags = getFirewallData(d, "nic.0.firewall.0.", false)
			if diags != nil {
				return diags
			}
			(*server.Entities.Nics.Items)[0].Entities = &ionoscloud.NicEntities{
				Firewallrules: &ionoscloud.FirewallRules{
					Items: &[]ionoscloud.FirewallRule{
						firewall,
					},
				},
			}
		}

		if primaryNic != nil && primaryNic.Properties != nil && primaryNic.Properties.Ips != nil {
			if len(*primaryNic.Properties.Ips) == 0 {
				*primaryNic.Properties.Ips = nil
			}
		}
	}

	createdServer, apiResponse, err := client.ServersApi.DatacentersServersPost(ctx, d.Get("datacenter_id").(string)).Server(server).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("error creating server: %s", err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}
	d.SetId(*createdServer.Id)

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			log.Printf("[DEBUG] failed to create createdServer resource")
			d.SetId("")
		}
		return utils.ToDiags(d, fmt.Sprintf("error waiting for state change for server creation %s", errState), &utils.DiagsOpts{Timeout: schema.TimeoutCreate})
	}

	// get additional data for schema
	createdServer, apiResponse, err = client.ServersApi.DatacentersServersFindById(ctx, d.Get("datacenter_id").(string), *createdServer.Id).Depth(3).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("error fetching server: (%s)", err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}
	if v, ok := d.GetOk("security_groups_ids"); ok {
		raw := v.(*schema.Set).List()
		nsgService := nsg.Service{Client: client, Meta: meta, D: d}
		if diagnostic := nsgService.PutServerNSG(ctx, dcID, *createdServer.Id, raw); diagnostic != nil {
			return diagnostic
		}
	}

	if createdServer.Entities != nil && createdServer.Entities.Nics != nil && len(*createdServer.Entities.Nics.Items) > 0 {
		firewallRules, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, d.Get("datacenter_id").(string),
			*createdServer.Id, *(*createdServer.Entities.Nics.Items)[0].Id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			requestLocation, _ := apiResponse.Location()
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching firewall rules: %s", err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
		}

		if firewallRules.Items != nil {
			if len(*firewallRules.Items) > 0 {
				if err := d.Set("firewallrule_id", *(*firewallRules.Items)[0].Id); err != nil {
					return utils.ToDiags(d, err.Error(), nil)
				}
			}
		}

		if (*createdServer.Entities.Nics.Items)[0].Id != nil {
			primaryNicID := *(*createdServer.Entities.Nics.Items)[0].Id
			err := d.Set("primary_nic", primaryNicID)
			if err != nil {
				return utils.ToDiags(d, fmt.Sprintf("error while setting primary nic: %s", err), nil)
			}
			if v, ok := d.GetOk("nic.0.security_groups_ids"); ok {
				raw := v.(*schema.Set).List()
				nsgService := nsg.Service{Client: client, Meta: meta, D: d}
				if diagnostic := nsgService.PutNICNSG(ctx, dcID, d.Id(), primaryNicID, raw); diagnostic != nil {
					return diagnostic
				}
			}
		}

		if (*createdServer.Entities.Nics.Items)[0].Properties.Ips != nil &&
			len(*(*createdServer.Entities.Nics.Items)[0].Properties.Ips) > 0 &&
			createdServer.Entities.Volumes != nil &&
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
	}

	// Set inline volumes
	if createdServer.Entities != nil && createdServer.Entities.Volumes != nil && createdServer.Entities.Volumes.Items != nil {
		var inlineVolumeIds []string
		for _, volume := range *createdServer.Entities.Volumes.Items {
			inlineVolumeIds = append(inlineVolumeIds, *volume.Id)
		}

		if err := d.Set("inline_volume_ids", inlineVolumeIds); err != nil {
			return utils.ToDiags(d, utils.GenerateSetError("server", "inline_volume_ids", err).Error(), nil)
		}
	}

	if initialState, ok := d.GetOk("vm_state"); ok {
		ss := cloudapiserver.Service{Client: client, Meta: meta, D: d}
		initialState := initialState.(string)

		if strings.EqualFold(initialState, constant.CubeVMStateStop) ||
			strings.EqualFold(initialState, constant.GpuVMStateStop) {
			if err := ss.Stop(ctx, dcID, d.Id(), serverType); err != nil {
				return utils.ToDiags(d, err.Error(), nil)
			}
		}

	}

	return resourceCubeServerRead(ctx, d, meta)
}
