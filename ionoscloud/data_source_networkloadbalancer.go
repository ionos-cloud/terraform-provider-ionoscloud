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

func dataSourceNetworkLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_lan": {
				Type:        schema.TypeInt,
				Description: "Id of the listening LAN. (inbound)",
				Computed:    true,
			},
			"ips": {
				Type: schema.TypeList,
				Description: "Collection of IP addresses of the Network Load Balancer. (inbound and outbound) IP of the " +
					"listenerLan must be a customer reserved IP for the public load balancer and private IP " +
					"for the private load balancer.",
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target_lan": {
				Type:        schema.TypeInt,
				Description: "Id of the balanced private target LAN. (outbound)",
				Computed:    true,
			},
			"lb_private_ips": {
				Type: schema.TypeList,
				Description: "Collection of private IP addresses with subnet mask of the Network Load Balancer. IPs " +
					"must contain valid subnet mask. If user will not provide any IP then the system will " +
					"generate one IP with /24 subnet.",
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

func dataSourceNetworkLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

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
	var networkLoadBalancer ionoscloud.NetworkLoadBalancer
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		networkLoadBalancer, apiResponse, err = client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, datacenterId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the network loadbalancer %s: %s", id.(string), err))
		}
	} else {
		/* search by name */
		var networkLoadBalancers ionoscloud.NetworkLoadBalancers

		networkLoadBalancers, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersGet(ctx, datacenterId.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancers: %s", err.Error()))
		}

		if networkLoadBalancers.Items != nil {
			for _, c := range *networkLoadBalancers.Items {
				tmpNetworkLoadBalancer, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, datacenterId.(string), *c.Id).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancer with ID %s: %s", *c.Id, err.Error()))
				}
				if tmpNetworkLoadBalancer.Properties.Name != nil {
					if strings.Contains(*tmpNetworkLoadBalancer.Properties.Name, name.(string)) {
						networkLoadBalancer = tmpNetworkLoadBalancer
						break
					}
				}

			}
		}

	}

	if &networkLoadBalancer == nil {
		return diag.FromErr(errors.New("network loadbalancer not found"))
	}

	if networkLoadBalancer.Id != nil {
		if err := d.Set("id", *networkLoadBalancer.Id); err != nil {
			return diag.FromErr(err)
		}
	}

	if err = setNetworkLoadBalancerData(d, &networkLoadBalancer); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
