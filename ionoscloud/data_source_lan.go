package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceLan() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLanRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datacenter_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
		return diag.FromErr(errors.New("please provide either the lan id or name"))
	}
	var lan ionoscloud.Lan
	var err error
	var apiResponse *ionoscloud.APIResponse
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	if idOk {
		/* search by ID */
		lan, apiResponse, err = client.LANsApi.DatacentersLansFindById(ctx, datacenterId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching lan with ID %s: %s", id.(string), err))
		}
	} else {
		/* search by name */
		lans, apiResponse, err := client.LANsApi.DatacentersLansGet(ctx, datacenterId.(string)).Depth(1).Filter("name", name.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching lans: %s", err.Error()))
		}

		if lans.Items != nil && len(*lans.Items) > 0 {
			lan = (*lans.Items)[len(*lans.Items)-1]
		} else {
			return diag.FromErr(fmt.Errorf("no lan found with the specified name"))
		}

	}

	if err = setLanData(d, &lan); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
