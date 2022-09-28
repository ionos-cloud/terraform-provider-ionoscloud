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
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceNatGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	datacenterId := d.Get("datacenter_id").(string)

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
	var natGateway ionoscloud.NatGateway
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for nat gateway by id %s", id)

		natGateway, apiResponse, err = client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, datacenterId, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the nat gateway %s: %w", id, err))
		}
	} else {
		/* search by name */
		var results []ionoscloud.NatGateway

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for nat gateway by name with partial_match %t and name: %s", partialMatch, name)

		if partialMatch {
			natGateways, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysGet(ctx, datacenterId).Depth(1).Filter("name", name).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway: %s", err.Error()))
			}

			results = *natGateways.Items
		} else {
			natGateways, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysGet(ctx, datacenterId).Depth(1).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway: %s", err.Error()))
			}

			if natGateways.Items != nil {
				for _, ng := range *natGateways.Items {
					if ng.Properties != nil && ng.Properties.Name != nil && strings.EqualFold(*ng.Properties.Name, name) {
						tmpNatGateway, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, datacenterId, *ng.Id).Execute()
						logApiRequestTime(apiResponse)
						if err != nil {
							return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway with ID %s: %s", *ng.Id, err.Error()))
						}
						natGateway = tmpNatGateway
						results = append(results, natGateway)
					}

				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no nat gateway found with the specified criteria: name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one nat gateway found with the specified criteria: name = %s", name))
		} else {
			natGateway = results[0]
		}

	}

	if err = setNatGatewayData(d, &natGateway); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
