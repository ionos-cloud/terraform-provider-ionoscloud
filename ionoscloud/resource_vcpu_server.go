package ionoscloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func resourceVCPUServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVCPUServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServerImport,
		},

		Schema: map[string]*schema.Schema{
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
			"cores": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Required: true,
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
				Optional:         true,
				Computed:         true,
				Deprecated:       "Please use the 'ionoscloud_server_boot_device_selection' resource for managing the boot device of the server.",
				Description:      "The associated boot drive, if any. Must be the UUID of a bootable CDROM image that can be retrieved using the ionoscloud_image data source.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"cpu_family": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
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
			"ssh_keys": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: "Public SSH keys are set on the image as authorized keys for appropriate SSH login to the instance using the corresponding private key. This field may only be set in creation requests. When reading, it always returns null. SSH keys are only supported if a public Linux image is used for the volume creation.",
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
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The size of the volume in GB.",
						},
						"disk_type": {
							Type:             schema.TypeString,
							Required:         true,
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
					},
				},
			},
			"vm_state": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Sets the power state of the vcpu server. Possible values: `RUNNING` or `SHUTOFF`.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{constant.VMStateStart, constant.VMStateStop}, true)),
			},
			"nic": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
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
						"dhcpv6": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ipv6_cidr_block": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
						"ipv6_ips": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							Computed: true,
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
						"security_groups_ids": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "The list of Security Group IDs for the NIC",
						},
						"firewall": {
							Description: "Firewall rules created in the server resource. The rules can also be created as separate resources outside the server resource",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
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
										Type:             schema.TypeInt,
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65534)),
									},
									"port_range_end": {
										Type:             schema.TypeInt,
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65534)),
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
			// When deleting the server, we need to delete the volume that was defined INSIDE the
			// server resource. The only way to know which volume was defined INSIDE the server
			// resource is to save the volume ID in this list.
			// NOTE: In the current version, v6.3.6, it is required to define one volume (and only
			// one) inside the server resource, but in the future we consider the possibility of
			// adding more volumes within the server resource, in which case the current code should
			// be revised as changes should also be made for the update function, as the list of
			// inline_volume_ids can be modified.
			"inline_volume_ids": {
				Type:        schema.TypeList,
				Description: "A list that contains the IDs for the volumes defined inside the server resource.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceVCPUServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := d.Set("type", constant.VCPUType); err != nil {
		return diag.FromErr(err)
	}
	return resourceServerCreate(ctx, d, meta)
}
