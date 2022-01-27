package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
)

func dataSourceNatGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNatGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceNatGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var natGateway ionoscloud.NatGateway
	var err error
	var apiResponse *ionoscloud.APIResponse
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if idOk {
		/* search by ID */
		natGateway, apiResponse, err = client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, datacenterId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the nat gateway %s: %w", id.(string), err))
		}
	} else {
		/* search by name */
		natGateways, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysGet(ctx, datacenterId.(string)).Depth(1).Filter("name", name.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateways: %s", err.Error()))
		}

		if natGateways.Items != nil && len(*natGateways.Items) > 0 {
			natGateway = (*natGateways.Items)[len(*natGateways.Items)-1]
			log.Printf("[WARN] %v nat gateways found matching the search criteria. Getting the latest nat gateway from the list %v", len(*natGateways.Items), *natGateway.Id)
		} else {
			return diag.FromErr(fmt.Errorf("no nat gateway found with the specified name %s", name.(string)))
		}

	}

	if err = setNatGatewayData(d, &natGateway); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
