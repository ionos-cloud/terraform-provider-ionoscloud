package ionoscloud

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func dataSourcePcc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePccRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						// The id of the cross-connected LAN
						"lan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						// The name of the cross-connected LAN
						"lan_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						// The id of the cross-connected VDC
						"datacenter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						// The name of the cross-connected VDC
						"datacenter_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"connectable_datacenters": {
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
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func convertPccPeers(peers *[]profitbricks.PCCPeer) []interface{} {
	if peers == nil {
		return make([]interface{}, 0)
	}

	ret := make([]interface{}, len(*peers), len(*peers))
	for i, peer := range *peers {
		entry := make(map[string]interface{})

		entry["lan_id"] = peer.LANId
		entry["lan_name"] = peer.LANName
		entry["datacenter_id"] = peer.DataCenterID
		entry["datacenter_name"] = peer.DataCenterName
		entry["location"] = peer.Location

		ret[i] = entry
	}

	return ret
}

func convertConnectableDatacenters(dcs *[]profitbricks.PCCConnectableDataCenter) []interface{} {
	if dcs == nil {
		return make([]interface{}, 0)
	}

	ret := make([]interface{}, len(*dcs), len(*dcs))
	for i, dc := range *dcs {
		entry := make(map[string]interface{})

		entry["id"] = dc.ID
		entry["name"] = dc.Name
		entry["location"] = dc.Location

		ret[i] = entry
	}

	return ret
}

func setPccDataSource(d *schema.ResourceData, pcc *profitbricks.PrivateCrossConnect) error {
	d.SetId(pcc.ID)

	if err := d.Set("id", pcc.ID); err != nil {
		return err
	}

	if err := d.Set("name", pcc.Properties.Name); err != nil {
		return err
	}
	if err := d.Set("description", pcc.Properties.Description); err != nil {
		return err
	}
	if err := d.Set("peers", convertPccPeers(pcc.Properties.Peers)); err != nil {
		return err
	}

	if err := d.Set("connectable_datacenters", convertConnectableDatacenters(pcc.Properties.ConnectableDatacenters)); err != nil {
		return err
	}
	return nil
}

func dataSourcePccRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the pcc id or name")
	}

	var pcc *profitbricks.PrivateCrossConnect
	var err error

	if idOk {
		/* search by ID */
		pcc, err = client.GetPrivateCrossConnect(id.(string))
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the pcc with ID %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var pccs *profitbricks.PrivateCrossConnects
		pccs, err := client.ListPrivateCrossConnects()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching pccs: %s", err.Error())
		}

		for _, p := range pccs.Items {
			if strings.Contains(p.Properties.Name, name.(string)) {
				/* lan found */
				pcc, err = client.GetPrivateCrossConnect(p.ID)
				if err != nil {
					return fmt.Errorf("an error occurred while fetching the pcc with ID %s: %s", p.ID, err)
				}
				break
			}
		}
	}

	if pcc == nil {
		return errors.New("pcc not found")
	}

	if err = setPccDataSource(d, pcc); err != nil {
		return err
	}

	return nil
}
