package serverutil

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var SchemaNicElem = map[string]*schema.Schema{
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
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Indicates whether this NIC receives an IPv6 address through DHCP.",
	},
	"ipv6_cidr_block": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "IPv6 CIDR block assigned to the NIC.",
	},
	"ips": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
		},
		Description: "Collection of IP addresses assigned to a nic. Explicitly assigned public IPs need to come from reserved IP blocks, Passing value null or empty array will assign an IP address automatically.",
		Computed:    true,
		Optional:    true,
	},
	"ipv6_ips": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv6Address),
		},
		Optional:    true,
		Computed:    true,
		Description: "Collection for IPv6 addresses assigned to a nic. Explicitly assigned IPv6 addresses need to come from inside the IPv6 CIDR block assigned to the nic.",
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
		MaxItems:    1,
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
}

var SchemaTemplatedDatasource = map[string]*schema.Schema{
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
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"hostname": {
		Type:     schema.TypeString,
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
		Elem:     CdromsServerDSResource,
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
				"require_legacy_bios": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates if the image requires the legacy BIOS for compatibility or specific needs.",
				},
			},
		},
	},
	"nics": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     NicServerDSResource,
	},
	"security_groups_ids": {
		Type:     schema.TypeList,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
	"ram": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"cores": {
		Type:     schema.TypeInt,
		Computed: true,
	},
}

// used for the datasource, when the nic is a member of the server object
var NicServerDSResource = &schema.Resource{
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
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"ipv6_ips": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Computed: true,
		},
		"dhcp": {
			Type:     schema.TypeBool,
			Computed: true,
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
			Elem:     FirewallServerDSResource,
		},
		"security_groups_ids": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	},
}

// used for the datasource, when the firewall is a member of the server object
var FirewallServerDSResource = &schema.Resource{
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
}

// used for the datasource, when the cdrom is a member of the server object
var CdromsServerDSResource = &schema.Resource{
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
}
