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

func dataSourceNetworkLoadBalancerForwardingRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkLoadBalancerForwardingRuleRead,
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
			"algorithm": {
				Type:        schema.TypeString,
				Description: "Algorithm for the balancing.",
				Computed:    true,
			},
			"protocol": {
				Type:        schema.TypeString,
				Description: "Protocol of the balancing.",
				Optional:    true,
			},
			"listener_ip": {
				Type:        schema.TypeString,
				Description: "Listening IP. (inbound)",
				Optional:    true,
			},
			"listener_port": {
				Type:        schema.TypeInt,
				Description: "Listening port number. (inbound) (range: 1 to 65535)",
				Computed:    true,
			},
			"health_check": {
				Type:        schema.TypeList,
				Description: "Health check attributes for Network Load Balancer forwarding rule",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_timeout": {
							Type: schema.TypeInt,
							Description: "ClientTimeout is expressed in milliseconds. This inactivity timeout applies " +
								"when the client is expected to acknowledge or send data. If unset the default of 50 " +
								"seconds will be used.",
							Computed: true,
						},
						"connect_timeout": {
							Type: schema.TypeInt,
							Description: "It specifies the maximum time (in milliseconds) to wait for a connection " +
								"attempt to a target VM to succeed. If unset, the default of 5 seconds will be used.",
							Computed: true,
						},
						"target_timeout": {
							Type: schema.TypeInt,
							Description: "TargetTimeout specifies the maximum inactivity time (in milliseconds) on the " +
								"target VM side. If unset, the default of 50 seconds will be used.",
							Computed: true,
						},
						"retries": {
							Type: schema.TypeInt,
							Description: "Retries specifies the number of retries to perform on a target VM after a " +
								"connection failure. If unset, the default value of 3 will be used.",
							Computed: true,
						},
					},
				},
			},
			"targets": {
				Type:        schema.TypeList,
				Description: "Array of items in that collection",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Description: "IP of a balanced target VM",
							Computed:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "Port of the balanced target service. (range: 1 to 65535)",
							Computed:    true,
						},
						"weight": {
							Type:        schema.TypeInt,
							Description: "Weight parameter is used to adjust the target VM's weight relative to other target VMs",
							Computed:    true,
						},
						"health_check": {
							Type:        schema.TypeList,
							Description: "Health check attributes for Network Load Balancer forwarding rule target",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"check": {
										Type:        schema.TypeBool,
										Description: "Check specifies whether the target VM's health is checked.",
										Computed:    true,
									},
									"check_interval": {
										Type: schema.TypeInt,
										Description: "CheckInterval determines the duration (in milliseconds) between " +
											"consecutive health checks. If unspecified a default of 2000 ms is used.",
										Computed: true,
									},
									"maintenance": {
										Type:        schema.TypeBool,
										Description: "Maintenance specifies if a target VM should be marked as down, even if it is not.",
										Computed:    true,
									},
								},
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
			"networkloadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceNetworkLoadBalancerForwardingRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	datacenterId := d.Get("datacenter_id").(string)
	networkloadbalancerId := d.Get("networkloadbalancer_id").(string)

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	protocolValue, protocolOk := d.GetOk("protocol")
	listenerIpValue, listenerIpOk := d.GetOk("listener_ip")

	id := idValue.(string)
	name := nameValue.(string)
	protocol := protocolValue.(string)
	listenerIp := listenerIpValue.(string)

	if idOk && (nameOk || protocolOk || listenerIpOk) {
		return diag.FromErr(errors.New("id and name/protocol/listener_ip cannot be both specified in the same time, choose between id or a combination of other parameters"))
	}
	if !idOk && !nameOk && !protocolOk && !listenerIpOk {
		return diag.FromErr(errors.New("please provide either the lan id or or other parameter like name or protocol"))
	}
	var networkLoadBalancerForwardingRule ionoscloud.NetworkLoadBalancerForwardingRule
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for network loadbalancers forwarding rule by id %s", id)
		networkLoadBalancerForwardingRule, apiResponse, err = client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId, networkloadbalancerId, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the network loadbalancer forwarding rule %s: %s", id, err))
		}
	} else {
		/* search by name */
		var results []ionoscloud.NetworkLoadBalancerForwardingRule

		if nameOk {
			partialMatch := d.Get("partial_match").(bool)

			log.Printf("[INFO] Using data source for network loadbalancers forwarding rule by name with partial_match %t and name: %s", partialMatch, name)

			//networkLoadBalancerForwardingRules, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesGet(ctx, datacenterId, networkloadbalancerId).Depth(1).Execute()
			//logApiRequestTime(apiResponse)
			//if err != nil {
			//	return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancers forwarding rules: %s", err.Error()))
			//}
			//var results = *networkLoadBalancerForwardingRules.Items

			if partialMatch {
				networkLoadBalancerForwardingRules, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesGet(ctx, datacenterId, networkloadbalancerId).Depth(1).Filter("name", name).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancers forwarding rules: %s", err.Error()))
				}
				results = *networkLoadBalancerForwardingRules.Items
			} else {
				networkLoadBalancerForwardingRules, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesGet(ctx, datacenterId, networkloadbalancerId).Depth(1).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancers forwarding rules: %s", err.Error()))
				}

				if networkLoadBalancerForwardingRules.Items != nil && nameOk { // aici pot schimba networkLoadBalancerForwardingRules.Items cu results
					var nameResults []ionoscloud.NetworkLoadBalancerForwardingRule
					for _, nlbfr := range *networkLoadBalancerForwardingRules.Items { // aici pot schimba networkLoadBalancerForwardingRules.Items cu results
						if nlbfr.Properties != nil && nlbfr.Properties.Name != nil && strings.EqualFold(*nlbfr.Properties.Name, name) {
							tmpNetworkLoadBalancerForwardingRule, apiResponse, err := client.NetworkLoadBalancersApi.
								DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId, networkloadbalancerId, *nlbfr.Id).Depth(1).Execute()
							logApiRequestTime(apiResponse)
							if err != nil {
								return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancer forwarding rule with ID %s: %s", *nlbfr.Id, err.Error()))
							}
							nameResults = append(nameResults, tmpNetworkLoadBalancerForwardingRule)
						}
					}
					results = nameResults
				}
			}
		} else {
			networkLoadBalancerForwardingRules, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesGet(ctx, datacenterId, networkloadbalancerId).Depth(1).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancers forwarding rules: %s", err.Error()))
			}
			results = *networkLoadBalancerForwardingRules.Items
		}

		if protocolOk && protocol != "" {
			var protocolResults []ionoscloud.NetworkLoadBalancerForwardingRule
			if results != nil {
				for _, nlbFwRule := range results {
					if nlbFwRule.Properties != nil && nlbFwRule.Properties.Protocol != nil && strings.EqualFold(*nlbFwRule.Properties.Protocol, protocol) {
						tmpNetworkLoadBalancerForwardingRule, apiResponse, err := client.NetworkLoadBalancersApi.
							DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId, networkloadbalancerId, *nlbFwRule.Id).Depth(1).Execute()
						logApiRequestTime(apiResponse)
						if err != nil {
							return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancer forwarding rule with ID %s: %s", *nlbFwRule.Id, err.Error()))
						}
						protocolResults = append(protocolResults, tmpNetworkLoadBalancerForwardingRule)
					}
				}
			}
			if protocolResults == nil {
				return diag.FromErr(fmt.Errorf("no network load balancer forwarding rule found with the specified criteria: protocol = %s", protocol))
			}
			results = protocolResults
		}

		if listenerIpOk && listenerIp != "" {
			var listenerIpResults []ionoscloud.NetworkLoadBalancerForwardingRule
			if results != nil {
				for _, nlbFwRule := range results {
					if nlbFwRule.Properties != nil && nlbFwRule.Properties.ListenerIp != nil && strings.EqualFold(*nlbFwRule.Properties.ListenerIp, listenerIp) {
						//tmpNetworkLoadBalancerForwardingRule, apiResponse, err := client.NetworkLoadBalancersApi.
						//	DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId, networkloadbalancerId, *nlbFwRule.Id).Depth(1).Execute()
						//logApiRequestTime(apiResponse)
						//if err != nil {
						//	return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancer forwarding rule with ID %s: %s", *nlbFwRule.Id, err.Error()))
						//}
						listenerIpResults = append(listenerIpResults, nlbFwRule)
					}
				}
			}
			if listenerIpResults == nil {
				return diag.FromErr(fmt.Errorf("no network load balancer forwarding rule found with the specified criteria: listener_ip = %s", listenerIp))
			}
			results = listenerIpResults
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no network load balancer forwarding rule found with the specified criteria: name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one network load balancer forwarding rule found with the specified criteria: name = %s", name))
		} else {
			networkLoadBalancerForwardingRule = results[0]
		}
	}

	if err = setNetworkLoadBalancerForwardingRuleData(d, &networkLoadBalancerForwardingRule); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
