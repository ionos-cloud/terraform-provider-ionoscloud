package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceNatGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNatGatewayRead,
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
			"public_ips": {
				Type:        schema.TypeList,
				Description: "Collection of public IP addresses of the NAT gateway. Should be customer reserved IP addresses in that location",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lans": {
				Type:        schema.TypeList,
				Description: "A list of Local Area Networks the node pool should be part of",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "Id for the LAN connected to the NAT gateway",
							Computed:    true,
						},
						"gateway_ips": {
							Type: schema.TypeList,
							Description: "Collection of gateway IP addresses of the NAT gateway. Will be auto-generated " +
								"if not provided. Should ideally be an IP belonging to the same subnet as the LAN",
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceNatGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return utils.ToDiags(d, "please provide either the nat gateway id or name", nil)
	}
	var natGateway ionoscloud.NatGateway
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		natGateway, apiResponse, err = client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, datacenterId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the nat gateway %s: %s", id.(string), err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		/* search by name */
		var natGateways ionoscloud.NatGateways

		natGateways, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysGet(ctx, datacenterId.(string)).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching nat gateway: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}

		var results []ionoscloud.NatGateway
		if natGateways.Items != nil {
			for _, ng := range *natGateways.Items {
				if ng.Properties != nil && ng.Properties.Name != nil && strings.EqualFold(*ng.Properties.Name, name.(string)) {
					tmpNatGateway, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, datacenterId.(string), *ng.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching nat gateway with ID %s: %s", *ng.Id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
					}
					natGateway = tmpNatGateway
					results = append(results, natGateway)
				}

			}
		}

		if results == nil || len(results) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no nat gateway found with the specified criteria: name = %s", name.(string)), nil)
		} else if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one nat gateway found with the specified criteria: name = %s", name.(string)), nil)
		} else {
			natGateway = results[0]
		}

	}

	if err = setNatGatewayData(d, &natGateway); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
