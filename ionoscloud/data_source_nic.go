package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"strings"
)

func dataSourceNIC() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNicRead,
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
			"lan": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dhcp": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ips": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}
func getNicDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"server_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
		"lan": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"dhcp": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"ips": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
			Optional: true,
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
	}
}
func dataSourceNicRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	t, dIdOk := data.GetOk("datacenter_id")
	st, sIdOk := data.GetOk("server_id")
	if !dIdOk || !sIdOk {
		return diag.FromErr(fmt.Errorf("datacenter id and server id must be set"))
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
	var err error
	var apiResponse *ionoscloud.APIResponse

	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("either id, or name must be set"))
	}
	if idOk {
		nic, apiResponse, err = client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, datacenterId, serverId, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting nic with id %s %w", id.(string), err))
		}
		if nameOk {
			if *nic.Properties.Name != name {
				return diag.FromErr(fmt.Errorf("name of nic (UUID=%s, name=%s) does not match expected name: %s",
					*nic.Id, *nic.Properties.Name, name))
			}
		}
	} else {
		nics, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsGet(ctx, datacenterId, serverId).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching nics: %w ", err))
		}

		var results []ionoscloud.Nic

		if nameOk && nics.Items != nil {
			for _, tempNic := range *nics.Items {
				if tempNic.Properties != nil && tempNic.Properties.Name != nil && strings.EqualFold(*tempNic.Properties.Name, name) {
					results = append(results, tempNic)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no nic found with the specified criteria: name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one nic found with the specified criteria: name = %s", name))
		} else {
			nic = results[0]
		}
	}

	if err := NicSetData(data, &nic); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
