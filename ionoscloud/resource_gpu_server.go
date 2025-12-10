package ionoscloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapifirewall"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapinic"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapiserver"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/nsg"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/slice"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func resourceGPUServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGpuServerCreate,
		ReadContext:   resourceGpuServerRead,
		UpdateContext: resourceGpuServerUpdate,
		DeleteContext: resourceGpuServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGpuServerImport,
		},
		CustomizeDiff: checkServerImmutableFields,

		Schema: map[string]*schema.Schema{
			// GPU Server parameters
			"template_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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
			"boot_cdrom": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				Deprecated:       "Please use the 'ionoscloud_server_boot_device_selection' resource for managing the boot device of the server.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
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
				ForceNew: true,
			},
			"ssh_key_path": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"volume.0.ssh_key_path"},
				Optional:      true,
				Computed:      true,
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
							Required:         true,
							ForceNew:         true,
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

								if len(slice.DiffString(slice.AnyToString(sshKeyPath), slice.AnyToString(oldSshKeyPath))) == 0 {
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
							ForceNew: true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2", "ZONE_3"}, true)),
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
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Sets the power state of the gpu server. Possible values: `RUNNING` or `SUSPENDED`.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{constant.VMStateStart, constant.CubeVMStateStop}, true)),
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

	var err error
	var volume *ionoscloud.VolumeProperties
	volume, err = getVolumeData(d, "volume.0.", constant.GpuType)
	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	image, imageAlias, err = getImage(ctx, client, d, *volume)
	if err != nil {
		return diag.FromErr(err)
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
				diags := diag.FromErr(fmt.Errorf("it is mandatory to provide either public image or imageAlias in conjunction with backup unit id property"))
				return diags
			}
			backupUnitID := backupUnitID.(string)
			volume.BackupunitId = &backupUnitID
		}
	}
	if userData, ok := d.GetOk("volume.0.user_data"); ok {
		if image == "" && imageAlias == "" {
			diags := diag.FromErr(fmt.Errorf("it is mandatory to provide either public image or imageAlias that has cloud-init compatibility in conjunction with backup unit id property "))
			return diags
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
	var primaryNic *ionoscloud.Nic
	if server.Entities.Nics != nil && server.Entities.Nics.Items != nil && len(*server.Entities.Nics.Items) > 0 {
		primaryNic = &(*server.Entities.Nics.Items)[0]
	}
	// Nic Arguments
	nic := ionoscloud.Nic{
		Properties: &ionoscloud.NicProperties{},
	}
	if _, ok := d.GetOk("nic"); ok {
		nic, err = cloudapinic.GetNicFromSchemaCreate(d, "nic.0.")
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("gpu error occurred while getting nic from schema: %w", err))
			return diags
		}
	}

	server.Entities.Nics = &ionoscloud.Nics{
		Items: &[]ionoscloud.Nic{
			nic,
		},
	}
	primaryNic = &(*server.Entities.Nics.Items)[0]
	log.Printf("[DEBUG] dhcp nic after %t", *nic.Properties.Dhcp)
	log.Printf("[DEBUG] primaryNic dhcp %t", *primaryNic.Properties.Dhcp)

	firewall := ionoscloud.FirewallRule{
		Properties: &ionoscloud.FirewallruleProperties{},
	}
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

	createdServer, apiResponse, err := client.ServersApi.DatacentersServersPost(ctx, d.Get("datacenter_id").(string)).Server(server).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error creating server: %w", err))
		return diags
	}
	d.SetId(*createdServer.Id)

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			log.Printf("[DEBUG] failed to create createdServer resource")
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("error waiting for state change for server creation %w", errState))
	}

	// get additional data for schema
	createdServer, apiResponse, err = client.ServersApi.DatacentersServersFindById(ctx, d.Get("datacenter_id").(string), *createdServer.Id).Depth(3).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error fetching server: (%w)", err))
		return diags
	}
	if v, ok := d.GetOk("security_groups_ids"); ok {
		raw := v.(*schema.Set).List()
		nsgService := nsg.Service{Client: client, Meta: meta, D: d}
		if diagnostic := nsgService.PutServerNSG(ctx, dcID, *createdServer.Id, raw); diagnostic != nil {
			return diagnostic
		}
	}

	firewallRules, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, d.Get("datacenter_id").(string),
		*createdServer.Id, *(*createdServer.Entities.Nics.Items)[0].Id).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching firewall rules: %w", err))
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
		primaryNicID := *(*createdServer.Entities.Nics.Items)[0].Id
		err := d.Set("primary_nic", primaryNicID)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting primary nic %s: %w", d.Id(), err))
			return diags
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

	// Set inline volumes
	if createdServer.Entities.Volumes != nil && createdServer.Entities.Volumes.Items != nil {
		var inlineVolumeIds []string
		for _, volume := range *createdServer.Entities.Volumes.Items {
			inlineVolumeIds = append(inlineVolumeIds, *volume.Id)
		}

		if err := d.Set("inline_volume_ids", inlineVolumeIds); err != nil {
			return diag.FromErr(utils.GenerateSetError("server", "inline_volume_ids", err))
		}
	}

	if initialState, ok := d.GetOk("vm_state"); ok {
		ss := cloudapiserver.Service{Client: client, Meta: meta, D: d}
		initialState := initialState.(string)

		if strings.EqualFold(initialState, constant.CubeVMStateStop) {
			err := ss.Stop(ctx, dcID, d.Id(), serverType)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceGpuServerRead(ctx, d, meta)
}

func resourceGpuServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	serverId := d.Id()

	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, serverId).Depth(2).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			log.Printf("[DEBUG] cannot find server by id \n")
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occurred while fetching a server ID %s %w", d.Id(), err))
		return diags
	}
	if server.Properties != nil {
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

		if server.Properties.Hostname != nil {
			if err := d.Set("hostname", *server.Properties.Hostname); err != nil {
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

		if server.Properties.VmState != nil {
			if err := d.Set("vm_state", *server.Properties.VmState); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
	}

	// upgrade from version without inline_volume_ids in GPU server
	if _, ok := d.GetOk("inline_volume_ids"); !ok {
		if bootVolume, ok := d.GetOk("boot_volume"); ok {
			bootVolume := bootVolume.(string)
			var inlineVolumeIds []string
			inlineVolumeIds = append(inlineVolumeIds, bootVolume)
			if err := d.Set("inline_volume_ids", inlineVolumeIds); err != nil {
				return diag.FromErr(utils.GenerateSetError("gpu_server", "inline_volume_ids", err))
			}
		}
	}

	if server.Entities != nil && server.Entities.Securitygroups != nil && server.Entities.Securitygroups.Items != nil {
		if err := nsg.SetNSGInResourceData(d, server.Entities.Securitygroups.Items); err != nil {
			return diag.FromErr(err)
		}
	}

	if server.Entities != nil && server.Entities.Volumes != nil && server.Entities.Volumes.Items != nil && len(*server.Entities.Volumes.Items) > 0 &&
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
		ns := cloudapinic.Service{Client: client, Meta: meta, D: d}

		nic, _, err := ns.Get(ctx, dcId, serverId, primarynic.(string), 0)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error occurred while fetching nic %s for server ID %s %s", primarynic.(string), d.Id(), err))
			return diags
		}

		if len(*nic.Properties.Ips) > 0 {
			if err := d.Set("primary_ip", (*nic.Properties.Ips)[0]); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}

		network := cloudapinic.SetNetworkProperties(*nic)

		if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
			network["ips"] = *nic.Properties.Ips
		}

		if firewallId, ok := d.GetOk("firewallrule_id"); ok {
			firewall, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, dcId, serverId, primarynic.(string), firewallId.(string)).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error occurred while fetching firewallrule %s for server ID %s %s", firewallId.(string), serverId, err))
				return diags
			}

			fw := cloudapifirewall.SetProperties(firewall)

			network["firewall"] = []map[string]interface{}{fw}
		}

		networks := []map[string]interface{}{network}
		if err := d.Set("nic", networks); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] unable to save nic to state IonosCloud Server (%s): %w", serverId, err))
		}
	}

	inlineVolumeIds := d.Get("inline_volume_ids")
	if inlineVolumeIds != nil {
		inlineVolumeIds := inlineVolumeIds.([]any)
		var volumes []any

		for i, volumeId := range inlineVolumeIds {
			volume, apiResponse, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, d.Id(), volumeId.(string)).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error retrieving inline volume %w", err))
			}
			volumePath := fmt.Sprintf("volume.%d.", i)
			entry := SetGpuVolumeProperties(volume)
			userData := d.Get(volumePath + "user_data")
			entry["user_data"] = userData
			backupUnit := d.Get(volumePath + "backup_unit_id")
			entry["backup_unit_id"] = backupUnit
			volumes = append(volumes, entry)
		}
		if err := d.Set("volume", volumes); err != nil {
			return diag.FromErr(fmt.Errorf("[DEBUG] Error saving inline volumes to state for GPU server (%s): %w", d.Id(), err))
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

func resourceGpuServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	ss := cloudapiserver.Service{Client: client, Meta: meta, D: d}

	dcId := d.Get("datacenter_id").(string)
	request := ionoscloud.ServerProperties{}

	currentVmState, err := ss.GetVmState(ctx, dcId, d.Id())
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("could not retrieve server vmState: %w", err))
		return diags
	}
	if strings.EqualFold(currentVmState, constant.CubeVMStateStop) && !d.HasChange("vm_state") {
		diags := diag.FromErr(fmt.Errorf("cannot update a suspended GPU Server, must change the state to RUNNING first"))
		return diags
	}

	// Unsuspend a GPU server first, before applying other changes
	if d.HasChange("vm_state") && strings.EqualFold(currentVmState, constant.CubeVMStateStop) {
		err := ss.Start(ctx, dcId, d.Id(), constant.GpuType)
		if err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

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

	if d.HasChange("hostname") {
		_, n := d.GetChange("hostname")
		nStr := n.(string)
		request.Hostname = &nStr
	}

	if d.HasChange("boot_cdrom") {
		_, n := d.GetChange("boot_cdrom")
		bootCdrom := n.(string)

		if utils.IsValidUUID(bootCdrom) {
			ss := cloudapiserver.Service{Client: meta.(bundleclient.SdkBundle).CloudApiClient, Meta: meta, D: d}
			ss.UpdateBootDevice(ctx, dcId, d.Id(), bootCdrom)
		}
	}

	server, apiResponse, err := client.ServersApi.DatacentersServersPatch(ctx, dcId, d.Id()).Server(request).Depth(3).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occurred while updating server ID %s: %w", d.Id(), err))
		return diags
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return diag.FromErr(errState)
	}

	if d.HasChange("security_groups_ids") {
		if v, ok := d.GetOk("security_groups_ids"); ok {
			raw := v.(*schema.Set).List()
			nsgService := nsg.Service{Client: client, Meta: meta, D: d}
			if diagnostic := nsgService.PutServerNSG(ctx, dcId, d.Id(), raw); diagnostic != nil {
				return diagnostic
			}
		}
	}

	// Volume stuff
	if d.HasChange("volume") {
		properties := ionoscloud.VolumeProperties{}
		inlineVolumeIds := d.Get("inline_volume_ids")

		if inlineVolumeIds != nil {
			inlineVolumeIds := inlineVolumeIds.([]interface{})
			for i, volumeId := range inlineVolumeIds {
				volumeIdStr := volumeId.(string)
				volumePath := fmt.Sprintf("volume.%d.", i)
				_, apiResponse, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, dcId, d.Id(), volumeIdStr).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					diags := diag.FromErr(fmt.Errorf("an error occurred while getting a volume dcId: %s server_id: %s ID: %s Response: %s", dcId, d.Id(), volumeId, err))
					return diags
				}
				if v, ok := d.GetOk(volumePath + "name"); ok {
					vStr := v.(string)
					properties.Name = &vStr
				}
				if v, ok := d.GetOk(volumePath + "bus"); ok {
					vStr := v.(string)
					properties.Bus = &vStr
				}

				if changed := d.HasChange(volumePath + "expose_serial"); changed {
					_, newVal := d.GetChange(volumePath + "expose_serial")
					exposeSerial := newVal.(bool)
					properties.ExposeSerial = &exposeSerial
				}

				if d.HasChange(volumePath + "require_legacy_bios") {
					_, newVal := d.GetChange(volumePath + "require_legacy_bios")
					requireLegacyBios := newVal.(bool)
					properties.RequireLegacyBios = &requireLegacyBios
				}

				_, apiResponse, err = client.VolumesApi.DatacentersVolumesPatch(ctx, d.Get("datacenter_id").(string), volumeIdStr).Volume(properties).Execute()
				logApiRequestTime(apiResponse)

				if err != nil {
					diags := diag.FromErr(fmt.Errorf("error patching volume (%s) (%w)", d.Id(), err))
					return diags
				}

				if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
					return diag.FromErr(errState)
				}
			}
		}
	}

	// Nic stuff
	if d.HasChange("nic") {
		nic := &ionoscloud.Nic{}
		nicStr := d.Get("primary_nic").(string)
		for _, n := range *server.Entities.Nics.Items {
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

		if v, ok := d.GetOk("nic.0.ipv6_cidr_block"); ok {
			ipv6Block := v.(string)
			properties.Ipv6CidrBlock = &ipv6Block
		}

		if v, ok := d.GetOk("nic.0.ipv6_ips"); ok {
			raw := v.([]interface{})
			ipv6Ips := make([]string, len(raw))
			if err := utils.DecodeInterfaceToStruct(raw, ipv6Ips); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
			if len(ipv6Ips) > 0 {
				properties.Ipv6Ips = &ipv6Ips
			}
		}

		if d.HasChange("nic.0.dhcpv6") {
			if dhcpv6, ok := d.GetOkExists("nic.0.dhcpv6"); ok {
				dhcpv6 := dhcpv6.(bool)
				properties.Dhcpv6 = &dhcpv6
			} else {
				properties.SetDhcpv6Nil()
			}
		}

		dhcp := d.Get("nic.0.dhcp").(bool)
		fwRule := d.Get("nic.0.firewall_active").(bool)
		properties.Dhcp = &dhcp
		properties.FirewallActive = &fwRule

		if v, ok := d.GetOk("nic.0.firewall_type"); ok {
			vStr := v.(string)
			properties.FirewallType = &vStr
		}

		if d.HasChange("nic.0.firewall") {

			firewallId := d.Get("firewallrule_id").(string)
			update := true
			if firewallId == "" {
				update = false
			}

			firewall, diags := getFirewallData(d, "nic.0.firewall.0.", update)
			if diags != nil {
				return diags
			}

			_, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, dcId, *server.Id, *nic.Id, firewallId).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				if !httpNotFound(apiResponse) {
					diags := diag.FromErr(fmt.Errorf("error occurred at checking existence of firewall %s %w", firewallId, err))
					return diags
				} else if httpNotFound(apiResponse) {
					diags := diag.FromErr(fmt.Errorf("firewall does not exist %s", firewallId))
					return diags
				}
			}
			if update == false {

				firewall, apiResponse, err = client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPost(ctx, dcId, *server.Id, *nic.Id).Firewallrule(firewall).Execute()
			} else {
				firewall, apiResponse, err = client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPatch(ctx, dcId, *server.Id, *nic.Id, firewallId).Firewallrule(*firewall.Properties).Execute()

			}
			logApiRequestTime(apiResponse)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occurred while running firewall rule dcId: %s server_id: %s nic_id %s ID: %s Response: %s", dcId, *server.Id, *nic.Id, firewallId, err))
				return diags
			}

			if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
				return diag.FromErr(fmt.Errorf("on gpu update an error occurred while waiting for state change dcId: %s server_id: %s nic_id %s ID: %s Response: %w", dcId, *server.Id, *nic.Id, firewallId, errState))
			}

			if firewallId == "" && firewall.Id != nil {
				if err := d.Set("firewallrule_id", firewall.Id); err != nil {
					diags := diag.FromErr(err)
					return diags
				}
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
		ns := cloudapinic.Service{Client: client, Meta: meta, D: d}
		_, _, err := ns.Update(ctx, d.Get("datacenter_id").(string), *server.Id, *nic.Id, properties)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error updating nic (%w)", err))
			return diags
		}

		if d.HasChange("nic.0.security_groups_ids") {
			if v, ok := d.GetOk("nic.0.security_groups_ids"); ok {
				raw := v.(*schema.Set).List()
				nsgService := nsg.Service{Client: client, Meta: meta, D: d}
				if diagnostic := nsgService.PutNICNSG(ctx, d.Get("datacenter_id").(string), *server.Id, *nic.Id, raw); diagnostic != nil {
					return diagnostic
				}
			}
		}
	}

	// Suspend a GPU server last, after applying other changes
	if d.HasChange("vm_state") && strings.EqualFold(currentVmState, constant.VMStateStart) {
		_, newVmState := d.GetChange("vm_state")
		if strings.EqualFold(newVmState.(string), constant.CubeVMStateStop) {
			err := ss.Stop(ctx, dcId, d.Id(), constant.GpuType)
			if err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
	}

	return resourceGpuServerRead(ctx, d, meta)
}

func SetGpuVolumeProperties(volume ionoscloud.Volume) map[string]interface{} {

	volumeMap := map[string]interface{}{}
	if volume.Properties != nil {
		utils.SetPropWithNilCheck(volumeMap, "name", volume.Properties.Name)
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
		utils.SetPropWithNilCheck(volumeMap, "expose_serial", volume.Properties.ExposeSerial)
		utils.SetPropWithNilCheck(volumeMap, "require_legacy_bios", volume.Properties.RequireLegacyBios)
	}
	return volumeMap
}

func resourceGpuServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	dcId := d.Get("datacenter_id").(string)

	apiResponse, err := client.ServersApi.DatacentersServersDelete(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while deleting a server ID %s %w", d.Id(), err))
		return diags

	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		return diag.FromErr(fmt.Errorf("error getting state change for gpu server delete %w", errState))
	}

	d.SetId("")
	return nil
}

func resourceGpuServerImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}", d.Id())
	}

	datacenterId := parts[0]
	serverId := parts[1]

	client := meta.(bundleclient.SdkBundle).CloudApiClient

	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, datacenterId, serverId).Depth(3).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find server %q", serverId)
		}
		return nil, fmt.Errorf("error occurred while fetching a server ID %s %w", d.Id(), err)
	}

	d.SetId(*server.Id)

	firstNicItem := (*server.Entities.Nics.Items)[0]
	if server.Entities != nil && server.Entities.Nics != nil && firstNicItem.Properties != nil &&
		firstNicItem.Properties.Ips != nil &&
		len(*firstNicItem.Properties.Ips) > 0 {
		log.Printf("[DEBUG] set primary_ip to %s", (*firstNicItem.Properties.Ips)[0])
		if err := d.Set("primary_ip", (*firstNicItem.Properties.Ips)[0]); err != nil {
			return nil, fmt.Errorf("error while setting primary ip %s: %w", d.Id(), err)
		}
	}

	if server.Entities != nil && server.Entities.Securitygroups != nil && server.Entities.Securitygroups.Items != nil {
		if err := nsg.SetNSGInResourceData(d, server.Entities.Securitygroups.Items); err != nil {
			return nil, err
		}
	}

	if err := d.Set("datacenter_id", datacenterId); err != nil {
		return nil, err
	}
	if server.Properties != nil {
		if server.Properties.Name != nil {
			if err := d.Set("name", *server.Properties.Name); err != nil {
				return nil, fmt.Errorf("error setting name %w", err)
			}
		}
		if server.Properties.Hostname != nil {
			if err := d.Set("hostname", *server.Properties.Hostname); err != nil {
				return nil, fmt.Errorf("error setting hostname %w", err)
			}
		}
		if server.Properties.TemplateUuid != nil {
			if err := d.Set("template_uuid", *server.Properties.TemplateUuid); err != nil {
				return nil, fmt.Errorf("error setting template uuid %w", err)
			}
		}

		if server.Properties.AvailabilityZone != nil {
			if err := d.Set("availability_zone", *server.Properties.AvailabilityZone); err != nil {
				return nil, fmt.Errorf("error setting availability_zone %w", err)
			}
		}
	}

	if server.Entities != nil && server.Entities.Volumes != nil &&
		len(*server.Entities.Volumes.Items) > 0 &&
		(*server.Entities.Volumes.Items)[0].Properties.Image != nil {
		if err := d.Set("boot_image", *(*server.Entities.Volumes.Items)[0].Properties.Image); err != nil {
			return nil, fmt.Errorf("error setting boot_image %w", err)
		}
	}

	if server.Entities != nil && server.Entities.Nics != nil && len(*server.Entities.Nics.Items) > 0 && (*server.Entities.Nics.Items)[0].Id != nil {
		primaryNic := *(*server.Entities.Nics.Items)[0].Id
		if err := d.Set("primary_nic", primaryNic); err != nil {
			return nil, fmt.Errorf("error setting primary_nic %w", err)
		}
		ns := cloudapinic.Service{Client: client, Meta: meta, D: d}

		nic, apiResponse, err := ns.Get(ctx, datacenterId, serverId, primaryNic, 0)
		if err != nil {
			return nil, err
		}

		if len(*nic.Properties.Ips) > 0 {
			if err := d.Set("primary_ip", (*nic.Properties.Ips)[0]); err != nil {
				return nil, fmt.Errorf("error setting primary_ip %w", err)
			}
		}

		network := cloudapinic.SetNetworkProperties(*nic)
		firewallRules, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, primaryNic).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return nil, err
		}

		if firewallRules.Items != nil {
			if len(*firewallRules.Items) > 0 {
				if err := d.Set("firewallrule_id", *(*firewallRules.Items)[0].Id); err != nil {
					return nil, fmt.Errorf("error setting firewallrule_id %w", err)
				}
			}
		}

		if firewallId, ok := d.GetOk("firewallrule_id"); ok {
			firewall, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, primaryNic, firewallId.(string)).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				return nil, err
			}

			fw := cloudapifirewall.SetProperties(firewall)

			network["firewall"] = []map[string]interface{}{fw}
		}

		networks := []map[string]interface{}{network}
		if err := d.Set("nic", networks); err != nil {
			return nil, fmt.Errorf("error setting nic %w", err)
		}
	}

	if server.Properties != nil && server.Properties.BootVolume != nil {
		if server.Properties.BootVolume.Id != nil {
			if err := d.Set("boot_volume", *server.Properties.BootVolume.Id); err != nil {
				return nil, fmt.Errorf("error setting boot_volume %w", err)
			}
		}
		volumeObj, apiResponse, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, datacenterId, serverId, *server.Properties.BootVolume.Id).Execute()
		logApiRequestTime(apiResponse)
		if err == nil {
			volumeItem := map[string]interface{}{}
			if volumeObj.Properties != nil {
				utils.SetPropWithNilCheck(volumeItem, "name", volumeObj.Properties.Name)
				utils.SetPropWithNilCheck(volumeItem, "disk_type", volumeObj.Properties.Type)
				utils.SetPropWithNilCheck(volumeItem, "licence_type", volumeObj.Properties.LicenceType)
				utils.SetPropWithNilCheck(volumeItem, "bus", volumeObj.Properties.Bus)
				utils.SetPropWithNilCheck(volumeItem, "availability_zone", volumeObj.Properties.AvailabilityZone)
			}

			volumesList := []map[string]interface{}{volumeItem}
			if err := d.Set("volume", volumesList); err != nil {
				return nil, fmt.Errorf("error setting volume %w", err)
			}
		}
	}
	if len(parts) > 3 {
		if err := d.Set("firewallrule_id", parts[3]); err != nil {
			return nil, fmt.Errorf("error setting firewallrule_id %w", err)
		}
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
