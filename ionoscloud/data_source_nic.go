package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapinic"
	cloudapiflowlog "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/flowlog"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func dataSourceNIC() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNicRead,
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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
			"lan": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dhcp": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dhcpv6": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ipv6_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"mac": {
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
			"flowlog": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     cloudapiflowlog.FlowlogSchemaDatasource,
				Description: `Flow logs holistically capture network information such as source and destination 
							IP addresses, source and destination ports, number of packets, amount of bytes, 
							the start and end time of the recording, and the type of protocol – 
							and log the extent to which your instances are being accessed.`,
			},
			"security_groups_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceNicRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	t, dIdOk := data.GetOk("datacenter_id")
	st, sIdOk := data.GetOk("server_id")
	if !dIdOk || !sIdOk {
		return diagutil.ToDiags(data, fmt.Errorf("datacenter id and server id must be set"), nil)
	}
	var datacenterId, serverId string

	if dIdOk {
		datacenterId = t.(string)
	}
	if sIdOk {
		serverId = st.(string)
	}
	var name string
	id, idOk := data.GetOk("id")

	t, nameOk := data.GetOk("name")
	if nameOk {
		name = t.(string)
	}
	var nic ionoscloud.Nic
	ns := cloudapinic.Service{Client: client, Meta: meta, D: data}
	if !idOk && !nameOk {
		return diagutil.ToDiags(data, fmt.Errorf("either id, or name must be set"), nil)
	}
	if idOk {
		foundNic, _, err := ns.Get(ctx, datacenterId, serverId, id.(string), 3)
		if err != nil {
			return diagutil.ToDiags(data, fmt.Errorf("error getting nic with id %s %w", id.(string), err), nil)
		}
		if nameOk {
			if foundNic.Properties != nil && *foundNic.Properties.Name != name {
				return diagutil.ToDiags(data, fmt.Errorf("name of nic (UUID=%s, name=%s) does not match expected name: %s",
					*foundNic.Id, *foundNic.Properties.Name, name), nil)
			}
		}
		nic = *foundNic
	} else {
		nics, err := ns.List(ctx, datacenterId, serverId, 3)
		if err != nil {
			return diagutil.ToDiags(data, fmt.Errorf("an error occurred while fetching nics: %w ", err), nil)
		}

		var results []ionoscloud.Nic

		if nameOk && nics != nil {
			for _, tempNic := range nics {
				if tempNic.Properties != nil && tempNic.Properties.Name != nil && *tempNic.Properties.Name == name {
					results = append(results, tempNic)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diagutil.ToDiags(data, fmt.Errorf("no nic found with the specified criteria: name = %s", name), nil)
		} else if len(results) > 1 {
			return diagutil.ToDiags(data, fmt.Errorf("more than one nic found with the specified criteria: name = %s", name), nil)
		} else {
			nic = results[0]
		}
	}

	if err := cloudapinic.NicSetData(data, &nic); err != nil {
		return diagutil.ToDiags(data, err, nil)
	}

	return nil
}
