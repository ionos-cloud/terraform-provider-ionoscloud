package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func dataSourcePcc() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePccRead,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				Optional:    true,
			},
			"peers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the cross-connected LAN",
						},
						"lan_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the cross-connected LAN",
						},
						"datacenter_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the cross-connected VDC",
						},
						"datacenter_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the cross-connected VDC",
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

func convertPccPeers(peers *[]ionoscloud.Peer) []interface{} {
	if peers == nil {
		return make([]interface{}, 0)
	}

	ret := make([]interface{}, len(*peers), len(*peers))
	for i, peer := range *peers {
		entry := make(map[string]interface{})

		entry["lan_id"] = peer.Id
		entry["lan_name"] = peer.Name
		entry["datacenter_id"] = peer.DatacenterId
		entry["datacenter_name"] = peer.DatacenterName
		entry["location"] = peer.Location

		ret[i] = entry
	}

	return ret
}

func convertConnectableDatacenters(dcs *[]ionoscloud.ConnectableDatacenter) []interface{} {
	if dcs == nil {
		return make([]interface{}, 0)
	}

	ret := make([]interface{}, len(*dcs), len(*dcs))
	for i, dc := range *dcs {
		entry := make(map[string]interface{})

		entry["id"] = dc.Id
		entry["name"] = dc.Name
		entry["location"] = dc.Location

		ret[i] = entry
	}

	return ret
}

func setPccDataSource(d *schema.ResourceData, pcc *ionoscloud.PrivateCrossConnect) error {

	if pcc.Id != nil {
		d.SetId(*pcc.Id)
	}

	if pcc.Properties != nil {
		if pcc.Properties.Name != nil {
			if err := d.Set("name", *pcc.Properties.Name); err != nil {
				return err
			}
		}
		if pcc.Properties.Description != nil {
			if err := d.Set("description", *pcc.Properties.Description); err != nil {
				return err
			}
		}
		if pcc.Properties.Peers != nil {
			if err := d.Set("peers", convertPccPeers(pcc.Properties.Peers)); err != nil {
				return err
			}
		}
		if pcc.Properties.ConnectableDatacenters != nil && len(*pcc.Properties.ConnectableDatacenters) > 0 {
			if err := d.Set("connectable_datacenters", convertConnectableDatacenters(pcc.Properties.ConnectableDatacenters)); err != nil {
				return err
			}
		}
	}
	return nil
}

func dataSourcePccRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("id and name cannot be both specified in the same time"), nil)
	}
	if !idOk && !nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("please provide either the pcc id or name"), nil)
	}

	var pcc ionoscloud.PrivateCrossConnect
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		pcc, apiResponse, err = client.PrivateCrossConnectsApi.PccsFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching the pcc with ID %s: %w", id.(string), err), &diagutil.ErrorContext{StatusCode: apiResponse.StatusCode})
		}
	}
	if nameOk {
		/* search by name */
		var pccs ionoscloud.PrivateCrossConnects
		pccs, apiResponse, err := client.PrivateCrossConnectsApi.PccsGet(ctx).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching pccs: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.StatusCode})
		}

		var results []ionoscloud.PrivateCrossConnect

		if pccs.Items != nil {
			for _, p := range *pccs.Items {
				if p.Properties != nil && p.Properties.Name != nil && *p.Properties.Name == name.(string) {
					pcc, apiResponse, err = client.PrivateCrossConnectsApi.PccsFindById(ctx, *p.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching the pcc with ID %s: %w", *p.Id, err), &diagutil.ErrorContext{StatusCode: apiResponse.StatusCode})
					}
					results = append(results, pcc)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diagutil.ToDiags(d, fmt.Errorf("no pcc found with the specified criteria: name = %s", name), nil)
		} else if len(results) > 1 {
			return diagutil.ToDiags(d, fmt.Errorf("more than one pcc found with the specified criteria: name = %s", name), nil)
		} else {
			pcc = results[0]
		}

	}

	if err = setPccDataSource(d, &pcc); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}
