package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
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
			"algorithm": {
				Type:        schema.TypeString,
				Description: "Algorithm for the balancing.",
				Computed:    true,
			},
			"protocol": {
				Type:        schema.TypeString,
				Description: "Protocol of the balancing.",
				Computed:    true,
			},
			"listener_ip": {
				Type:        schema.TypeString,
				Description: "Listening IP. (inbound)",
				Computed:    true,
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

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return diag.FromErr(errors.New("no datacenter_id was specified"))
	}

	networkloadbalancerId, nlbIdOk := d.GetOk("networkloadbalancer_id")
	if !nlbIdOk {
		return diag.FromErr(errors.New("no networkloadbalancer_id was specified"))
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the lan id or name"))
	}
	var networkLoadBalancerForwardingRule ionoscloud.NetworkLoadBalancerForwardingRule
	var err error
	var apiResponse *shared.APIResponse

	if idOk {
		/* search by ID */
		networkLoadBalancerForwardingRule, apiResponse, err = client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId.(string), networkloadbalancerId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the network loadbalancer forwarding rule %s: %w", id.(string), err))
		}
	} else {
		/* search by name */
		var networkLoadBalancerForwardingRules ionoscloud.NetworkLoadBalancerForwardingRules

		networkLoadBalancerForwardingRules, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesGet(ctx, datacenterId.(string), networkloadbalancerId.(string)).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancers forwarding rules: %w", err))
		}

		var results []ionoscloud.NetworkLoadBalancerForwardingRule
		if networkLoadBalancerForwardingRules.Items != nil {
			for _, nlbfr := range networkLoadBalancerForwardingRules.Items {
				if strings.EqualFold(nlbfr.Properties.Name, name.(string)) {
					tmpNetworkLoadBalancerForwardingRule, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId.(string), networkloadbalancerId.(string), *nlbfr.Id).Depth(1).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancer forwarding rule with ID %s: %w", *nlbfr.Id, err))
					}
					results = append(results, tmpNetworkLoadBalancerForwardingRule)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no network load balancer forwarding rule found with the specified criteria: name = %s", name.(string)))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one network load balancer forwarding rule found with the specified criteria: name = %s", name.(string)))
		} else {
			networkLoadBalancerForwardingRule = results[0]
		}
	}

	if err = setNetworkLoadBalancerForwardingRuleData(d, &networkLoadBalancerForwardingRule); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
