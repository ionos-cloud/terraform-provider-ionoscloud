package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"strings"
)

func dataSourceNatGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNatGatewayRead,
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

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
			return fmt.Errorf("an error occurred while fetching the nat gateway %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var natGateways ionoscloud.NatGateways

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		natGateways, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysGet(ctx, datacenterId.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching nat gateway: %s", err.Error())
		}

		if natGateways.Items != nil {
			for _, c := range *natGateways.Items {
				tmpNatGateway, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, datacenterId.(string), *c.Id).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return fmt.Errorf("an error occurred while fetching nat gateway with ID %s: %s", *c.Id, err.Error())
				}
				if tmpNatGateway.Properties.Name != nil {
					if strings.Contains(*tmpNatGateway.Properties.Name, name.(string)) {
						natGateway = tmpNatGateway
						break
					}
				}

			}
		}

	}

	if &natGateway == nil {
		return errors.New("nat gateway not found")
	}

	if natGateway.Id != nil {
		if err := d.Set("id", *natGateway.Id); err != nil {
			return err
		}
	}

	if err = setNatGatewayData(d, &natGateway); err != nil {
		return err
	}

	return nil
}

func setNatGatewayData(d *schema.ResourceData, natGateway *ionoscloud.NatGateway) error {

	if natGateway.Id != nil {
		d.SetId(*natGateway.Id)
	}

	if natGateway.Properties != nil {
		if natGateway.Properties.Name != nil {
			err := d.Set("name", *natGateway.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for nat gateway %s: %s", d.Id(), err)
			}
		}

		if natGateway.Properties.PublicIps != nil {
			err := d.Set("public_ips", *natGateway.Properties.PublicIps)
			if err != nil {
				return fmt.Errorf("error while setting public_ips property for nat gateway %s: %s", d.Id(), err)
			}
		}

		if natGateway.Properties.Lans != nil && len(*natGateway.Properties.Lans) > 0 {
			var natGatewayLans []interface{}
			for _, lan := range *natGateway.Properties.Lans {
				lanEntry := make(map[string]interface{})

				if lan.Id != nil {
					lanEntry["id"] = *lan.Id
				}

				if len(*lan.GatewayIps) > 0 {
					var gatewayIps []interface{}
					for _, gatewayIp := range *lan.GatewayIps {
						gatewayIps = append(gatewayIps, gatewayIp)
					}
					if len(gatewayIps) > 0 {
						lanEntry["gateway_ips"] = gatewayIps
					}
				}

				natGatewayLans = append(natGatewayLans, lanEntry)
			}

			if len(natGatewayLans) > 0 {
				if err := d.Set("lans", natGatewayLans); err != nil {
					return fmt.Errorf("error while setting lans property for nat gateway %s: %s", d.Id(), err)
				}
			}
		}
	}
	return nil
}
