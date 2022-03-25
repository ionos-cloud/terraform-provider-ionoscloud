package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"strings"
)

func dataSourceApplicationLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Application Load Balancer.",
				Optional:    true,
			},
			"listener_lan": {
				Type:        schema.TypeInt,
				Description: "D of the listening (inbound) LAN.",
				Computed:    true,
			},
			"ips": {
				Type:        schema.TypeSet,
				Description: "Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target_lan": {
				Type:        schema.TypeInt,
				Description: "ID of the balanced private target LAN (outbound).",
				Computed:    true,
			},
			"lb_private_ips": {
				Type:        schema.TypeSet,
				Description: "Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceApplicationLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.FromErr(errors.New("please provide either the application load balancer id or name"))
	}

	var applicationLoadBalancer ionoscloud.ApplicationLoadBalancer
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		applicationLoadBalancer, apiResponse, err = client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, datacenterId, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the nat gateway %s: %w", id, err))
		}
	} else {
		/* search by name */
		var applicationLoadBalancers ionoscloud.ApplicationLoadBalancers

		applicationLoadBalancers, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersGet(ctx, datacenterId).Depth(5).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway: %s", err.Error()))
		}

		var results []ionoscloud.ApplicationLoadBalancer

		if applicationLoadBalancers.Items != nil {
			for _, alb := range *applicationLoadBalancers.Items {
				if alb.Properties != nil && alb.Properties.Name != nil && strings.ToLower(*alb.Properties.Name) == strings.ToLower(name) {
					tmpAlb, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, datacenterId, *alb.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway with ID %s: %s", *alb.Id, err.Error()))
					}
					results = append(results, tmpAlb)
				}

			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no application load balanacer found with the specified criteria: name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one application load balanacer found with the specified criteria: name = %s", name))
		} else {
			applicationLoadBalancer = results[0]
		}
	}

	if err = setApplicationLoadBalancerData(d, &applicationLoadBalancer); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
