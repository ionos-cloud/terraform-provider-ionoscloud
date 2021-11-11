package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
)

func resourceNatGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNatGatewayCreate,
		ReadContext:   resourceNatGatewayRead,
		UpdateContext: resourceNatGatewayUpdate,
		DeleteContext: resourceNatGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNatGatewayImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the NAT gateway",
				Required:    true,
			},
			"public_ips": {
				Type:        schema.TypeList,
				Description: "Collection of public IP addresses of the NAT gateway. Should be customer reserved IP addresses in that location",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lans": {
				Type:        schema.TypeList,
				Description: "A list of Local Area Networks the node pool should be part of",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "Id for the LAN connected to the NAT gateway",
							Required:    true,
						},
						"gateway_ips": {
							Type: schema.TypeList,
							Description: "Collection of gateway IP addresses of the NAT gateway. Will be auto-generated " +
								"if not provided. Should ideally be an IP belonging to the same subnet as the LAN",
							Optional: true,
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

func resourceNatGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	name := d.Get("name").(string)

	natGateway := ionoscloud.NatGateway{
		Properties: &ionoscloud.NatGatewayProperties{
			Name: &name,
		},
	}

	if publicIpsVal, publicIpsOk := d.GetOk("public_ips"); publicIpsOk {
		publicIpsVal := publicIpsVal.([]interface{})
		if publicIpsVal != nil {
			publicIps := make([]string, len(publicIpsVal), len(publicIpsVal))
			for idx := range publicIpsVal {
				publicIps[idx] = fmt.Sprint(publicIpsVal[idx])
			}
			natGateway.Properties.PublicIps = &publicIps
		} else {
			diags := diag.FromErr(fmt.Errorf("you must provide public_ips for nat gateway resource \n"))
			return diags
		}
	}

	if lansVal, lansOK := d.GetOk("lans"); lansOK {
		if lansVal.([]interface{}) != nil {
			updateLans := false
			var lans []ionoscloud.NatGatewayLanProperties

			for lanIndex := range lansVal.([]interface{}) {
				lan := ionoscloud.NatGatewayLanProperties{}
				addLan := false
				if lanID, lanIdOk := d.GetOk(fmt.Sprintf("lans.%d.id", lanIndex)); lanIdOk {
					lanID := int32(lanID.(int))
					lan.Id = &lanID
					addLan = true
				}
				if lanGatewayIps, lanGatewayIpsOk := d.GetOk(fmt.Sprintf("lans.%d.gateway_ips", lanIndex)); lanGatewayIpsOk {
					lanGatewayIps := lanGatewayIps.([]interface{})
					if lanGatewayIps != nil {
						gatewayIps := make([]string, len(lanGatewayIps), len(lanGatewayIps))
						for idx := range lanGatewayIps {
							gatewayIps[idx] = fmt.Sprint(lanGatewayIps[idx])
						}
						lan.GatewayIps = &gatewayIps
					}
				}
				if addLan {
					lans = append(lans, lan)
				}
			}

			if len(lans) > 0 {
				updateLans = true
			}

			if updateLans == true {
				log.Printf("[INFO] NatGateway LANs set to %+v", lans)
				natGateway.Properties.Lans = &lans
			} else {
				diags := diag.FromErr(fmt.Errorf("you must provide lans for the nat gateway resource \n"))
				return diags
			}
		}
	}

	dcId := d.Get("datacenter_id").(string)

	log.Printf("[*****] %+v\n", natGateway)
	natGatewayResp, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysPost(ctx, dcId).NatGateway(natGateway).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating natGateway: %s, %s", err, responseBody(apiResponse)))
		return diags
	}

	d.SetId(*natGatewayResp.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceNatGatewayRead(ctx, d, meta)
}

func resourceNatGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)

	natGateway, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] Successfully retreived nat gateway %s: %+v", d.Id(), natGateway)

	if err := setNatGatewayData(d, &natGateway); err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	return nil
}

func resourceNatGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)
	request := ionoscloud.NatGateway{
		Properties: &ionoscloud.NatGatewayProperties{},
	}

	dcId := d.Get("datacenter_id").(string)

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		vStr := v.(string)
		request.Properties.Name = &vStr
	}

	if d.HasChange("public_ips") {
		oldPublicIps, newPublicIps := d.GetChange("public_ips")
		log.Printf("[INFO] nat gateway public IPs changed from %+v to %+v", oldPublicIps, newPublicIps)
		publicIpsVal := newPublicIps.([]interface{})
		if publicIpsVal != nil {
			publicIps := make([]string, len(publicIpsVal), len(publicIpsVal))
			for idx := range publicIpsVal {
				publicIps[idx] = fmt.Sprint(publicIpsVal[idx])
			}
			request.Properties.PublicIps = &publicIps
		}
	}

	if d.HasChange("lans") {
		oldLANs, newLANs := d.GetChange("lans")
		if newLANs.([]interface{}) != nil {
			updateLans := false
			var lans []ionoscloud.NatGatewayLanProperties

			for lanIndex := range newLANs.([]interface{}) {
				lan := ionoscloud.NatGatewayLanProperties{}
				addLan := false
				if lanID, lanIdOk := d.GetOk(fmt.Sprintf("lans.%d.id", lanIndex)); lanIdOk {
					lanID := int32(lanID.(int))
					lan.Id = &lanID
					addLan = true
				}
				if lanGatewayIps, lanGatewayIpsOk := d.GetOk(fmt.Sprintf("lans.%d.gateway_ips", lanIndex)); lanGatewayIpsOk {
					lanGatewayIps := lanGatewayIps.([]interface{})
					if lanGatewayIps != nil {
						gatewayIps := make([]string, len(lanGatewayIps), len(lanGatewayIps))
						for idx := range lanGatewayIps {
							gatewayIps[idx] = fmt.Sprint(lanGatewayIps[idx])
						}
						lan.GatewayIps = &gatewayIps
					}
				}
				if addLan {
					lans = append(lans, lan)
				}
			}

			if len(lans) > 0 {
				updateLans = true
			}

			if updateLans == true {
				log.Printf("[INFO] nat gateway  LANs changed from %+v to %+v", oldLANs, newLANs)
				request.Properties.Lans = &lans
			}
		}
	}

	_, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysPatch(ctx, dcId, d.Id()).NatGatewayProperties(*request.Properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a nat gateway ID %s %s", d.Id(), err))
		return diags
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceNatGatewayRead(ctx, d, meta)
}

func resourceNatGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)

	apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysDelete(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a nat gateway %s %s", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	d.SetId("")

	return nil
}

func resourceNatGatewayImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{natgateway}", d.Id())
	}

	dcId := parts[0]
	natGatewayId := parts[1]

	natGateway, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, dcId, natGatewayId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find nat gateway  %q", natGatewayId)
		}
		return nil, fmt.Errorf("an error occured while retrieving nat gateway  %q: %q ", natGatewayId, err)
	}

	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, err
	}

	if err := setNatGatewayData(d, &natGateway); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
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
					lanEntry["gateway_ips"] = gatewayIps

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
