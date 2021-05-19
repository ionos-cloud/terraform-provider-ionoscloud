package ionoscloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
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
				Type:     schema.TypeString,
				Required: true,
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

func convertIpFailoverList(ips *[]profitbricks.IPFailover) []interface{} {
	if ips == nil {
		return make([]interface{}, 0)
	}

	ret := make([]interface{}, len(*ips), len(*ips))
	for i, ip := range *ips {
		entry := make(map[string]interface{})

		entry["ip"] = ip.IP
		entry["nic_uuid"] = ip.NicUUID

		ret[i] = entry
	}

	return ret
}

func setLanData(d *schema.ResourceData, lan *profitbricks.Lan) error {
	d.SetId(lan.ID)
	if err := d.Set("id", lan.ID); err != nil {
		return err
	}

	if err := d.Set("name", lan.Properties.Name); err != nil {
		return err
	}
	if err := d.Set("ip_failover", convertIpFailoverList(lan.Properties.IPFailover)); err != nil {
		return err
	}
	if err := d.Set("pcc", lan.Properties.PCC); err != nil {
		return err
	}
	if err := d.Set("public", lan.Properties.Public); err != nil {
		return err
	}
	return nil
}

func dataSourceLanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

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
	var lan *profitbricks.Lan
	var err error

	if idOk {
		/* search by ID */
		lan, err = client.GetLan(datacenterId.(string), id.(string))
		if err != nil {
			return fmt.Errorf("an error occurred while fetching lan with ID %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var lans *profitbricks.Lans
		lans, err := client.ListLans(datacenterId.(string))
		if err != nil {
			return fmt.Errorf("an error occurred while fetching lans: %s", err.Error())
		}

		for _, l := range lans.Items {
			if l.Properties.Name == name.(string) {
				/* lan found */
				lan, err = client.GetLan(datacenterId.(string), l.ID)
				if err != nil {
					return fmt.Errorf("an error occurred while fetching lan %s: %s", l.ID, err)
				}
				break
			}
		}
	}

	if lan == nil {
		return errors.New("lan not found")
	}

	if err = setLanData(d, lan); err != nil {
		return err
	}

	return nil
}
