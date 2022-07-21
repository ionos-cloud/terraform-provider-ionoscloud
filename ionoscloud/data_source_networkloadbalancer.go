package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
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
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
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
				Optional:    true,
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
	client := meta.(SdkBundle).CloudApiClient

	datacenterId := d.Get("datacenter_id").(string)

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	targetLanValue, targetLanOk := d.GetOk("target_lan")

	id := idValue.(string)
	name := nameValue.(string)
	targetLan := targetLanValue.(int)

	if idOk && (nameOk || targetLanOk) {
		return diag.FromErr(errors.New("id and name/target_lan cannot be both specified in the same time, choose between id or a combination of other parameters"))
	}
	if !idOk && !nameOk && !targetLanOk {
		return diag.FromErr(errors.New("please provide either the lan id or other parameter like  target_lan or name"))
	}
	var networkLoadBalancer ionoscloud.NetworkLoadBalancer
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for network loadbalancer by id %s", id)
		networkLoadBalancer, apiResponse, err = client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, datacenterId, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the network loadbalancer %s: %w", id, err))
		}
	} else {
		/* search by name */
		var results []ionoscloud.NetworkLoadBalancer
		if nameOk {
			partialMatch := d.Get("partial_match").(bool)

			log.Printf("[INFO] Using data source for network loadbalancer by name with partial_match %t and name: %s", partialMatch, name)

			if partialMatch {
				networkLoadBalancers, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersGet(ctx, datacenterId).Depth(1).Filter("name", name).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancers: %s", err.Error()))
				}
				if len(*networkLoadBalancers.Items) == 0 {
					return diag.FromErr(fmt.Errorf("no result found with the specified criteria: name with partial match: %s", name))
				}
				results = *networkLoadBalancers.Items
			} else {
				networkLoadBalancers, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersGet(ctx, datacenterId).Depth(1).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancers: %s", err.Error()))
				}

				if networkLoadBalancers.Items != nil && nameOk {
					var resultsByName []ionoscloud.NetworkLoadBalancer
					for _, nlb := range *networkLoadBalancers.Items {
						if nlb.Properties != nil && nlb.Properties.Name != nil && strings.EqualFold(*nlb.Properties.Name, name) {
							tmpNetworkLoadBalancer, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, datacenterId, *nlb.Id).Execute()
							logApiRequestTime(apiResponse)
							if err != nil {
								return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancer with ID %s: %s", *nlb.Id, err.Error()))
							}
							// results = append(results, tmpNetworkLoadBalancer)
							resultsByName = append(resultsByName, tmpNetworkLoadBalancer)
						}
					}
					if resultsByName == nil || len(resultsByName) == 0 {
						return diag.FromErr(fmt.Errorf("no result found with the specified criteria: name = %s", name))
					}
					results = resultsByName
				}
			}
		} else {
			networkLoadBalancers, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersGet(ctx, datacenterId).Depth(1).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancers: %s", err.Error()))
			}
			results = *networkLoadBalancers.Items
		}

		if targetLanOk && targetLan != 0 {
			var targetLanResults []ionoscloud.NetworkLoadBalancer
			if results != nil {
				for _, networkLoadBalancer := range results {
					if networkLoadBalancer.Properties != nil && networkLoadBalancer.Properties.TargetLan != nil && int(*networkLoadBalancer.Properties.TargetLan) == targetLan {
						targetLanResults = append(targetLanResults, networkLoadBalancer)
					}
				}
			}
			if targetLanResults == nil || len(targetLanResults) == 0 {
				return diag.FromErr(fmt.Errorf("no result found with the specified criteria: target_lan = %d", targetLan))
			}
			results = targetLanResults
		}

		if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one network load balancer found with the specified criteria: name = %s", name))
		} else {
			networkLoadBalancer = results[0]
		}
	}

	if err = setNetworkLoadBalancerData(d, &networkLoadBalancer); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
