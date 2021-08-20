package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func dataSourceLan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLanRead,
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

func setLanData(d *schema.ResourceData, lan *ionoscloud.Lan) error {
	d.SetId(*lan.Id)
	if err := d.Set("id", *lan.Id); err != nil {
		return err
	}

	if lan.Properties != nil {
		if lan.Properties.Name != nil {
			if err := d.Set("name", *lan.Properties.Name); err != nil {
				return err
			}
		}
		if lan.Properties.IpFailover != nil && len(*lan.Properties.IpFailover) > 0 {
			if err := d.Set("ip_failover", convertIpFailoverList(lan.Properties.IpFailover)); err != nil {
				return err
			}
		}
		if lan.Properties.Pcc != nil {
			if err := d.Set("pcc", *lan.Properties.Pcc); err != nil {
				return err
			}
		}
		if lan.Properties.Public != nil {
			if err := d.Set("public", *lan.Properties.Public); err != nil {
				return err
			}
		}
	}

	return nil
}

func dataSourceLanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).CloudApiClient

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return errors.New("no datacenter_id was specified")
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the lan id or name")
	}
	var lan ionoscloud.Lan
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	if idOk {
		/* search by ID */
		lan, _, err = client.LanApi.DatacentersLansFindById(ctx, datacenterId.(string), id.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching lan with ID %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var lans ionoscloud.Lans

		lans, _, err := client.LanApi.DatacentersLansGet(ctx, datacenterId.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching lans: %s", err.Error())
		}

		found := false
		if lans.Items != nil {
			for _, l := range *lans.Items {
				if l.Properties.Name != nil && *l.Properties.Name == name.(string) {
					/* lan found */
					lan, _, err = client.LanApi.DatacentersLansFindById(ctx, datacenterId.(string), *l.Id).Execute()
					if err != nil {
						return fmt.Errorf("an error occurred while fetching lan %s: %s", *l.Id, err)
					}
					found = true
					break
				}
			}
		}

		if !found {
			return fmt.Errorf("lan not found")
		}
	}

	if err = setLanData(d, &lan); err != nil {
		return err
	}

	return nil
}
