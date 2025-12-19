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
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
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
