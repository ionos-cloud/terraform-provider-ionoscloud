package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceLan() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLanRead,
		Schema: map[string]*schema.Schema{
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
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"ip_failover": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nic_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"pcc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipv4_cidr_block": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "For public LANs this property is null, for private LANs it contains the private IPv4 CIDR range.",
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func convertIpFailoverList(ips *[]ionoscloud.IPFailover) []interface{} {
	if ips == nil {
		return make([]interface{}, 0)
	}

	ret := make([]interface{}, len(*ips), len(*ips))
	for i, ip := range *ips {
		entry := make(map[string]interface{})

		entry["ip"] = ip.Ip
		entry["nic_uuid"] = ip.NicUuid

		ret[i] = entry
	}

	return ret
}

func dataSourceLanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return utils.ToDiags(d, "no datacenter_id was specified", nil)
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return utils.ToDiags(d, "id and name cannot be both specified in the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the lan id or name", nil)
	}
	var lan ionoscloud.Lan
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		lan, apiResponse, err = client.LANsApi.DatacentersLansFindById(ctx, datacenterId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching lan with ID %s: %s", id.(string), err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		/* search by name */
		var lans ionoscloud.Lans

		lans, apiResponse, err := client.LANsApi.DatacentersLansGet(ctx, datacenterId.(string)).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching lans: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}

		var results []ionoscloud.Lan

		if lans.Items != nil {
			for _, l := range *lans.Items {
				if l.Properties != nil && l.Properties.Name != nil && *l.Properties.Name == name.(string) {
					/* lan found */
					lan, apiResponse, err = client.LANsApi.DatacentersLansFindById(ctx, datacenterId.(string), *l.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching lan %s: %s", *l.Id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
					}
					results = append(results, l)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no lan found with the specified name: %s", name), nil)
		} else if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one lan found with the specified criteria name: %s", name), nil)
		} else {
			lan = results[0]
		}
	}

	if err = setLanData(d, &lan); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
