package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
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
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
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

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the lan id or name"))
	}
	var lan ionoscloud.Lan
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for lan by id %s", id)

		lan, apiResponse, err = client.LANsApi.DatacentersLansFindById(ctx, datacenterId.(string), id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching lan with ID %s: %w", id, err))
		}
	} else {
		/* search by name */
		var results []ionoscloud.Lan

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for lan by name with partial_match %t and name: %s", partialMatch, name)

		if partialMatch {
			lans, apiResponse, err := client.LANsApi.DatacentersLansGet(ctx, datacenterId.(string)).Depth(1).Filter("name", name).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching lans: %s", err.Error()))
			}
			results = *lans.Items
		} else {
			lans, apiResponse, err := client.LANsApi.DatacentersLansGet(ctx, datacenterId.(string)).Depth(1).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching lans: %s", err.Error()))
			}

			if lans.Items != nil {
				for _, l := range *lans.Items {
					if l.Properties != nil && l.Properties.Name != nil && strings.EqualFold(*l.Properties.Name, name) {
						/* lan found */
						lan, apiResponse, err = client.LANsApi.DatacentersLansFindById(ctx, datacenterId.(string), *l.Id).Execute()
						logApiRequestTime(apiResponse)
						if err != nil {
							return diag.FromErr(fmt.Errorf("an error occurred while fetching lan %s: %w", *l.Id, err))
						}
						results = append(results, l)
					}
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no lan found with the specified name: %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one lan found with the specified criteria name: %s", name))
		} else {
			lan = results[0]
		}
	}

	if err = setLanData(d, &lan); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
