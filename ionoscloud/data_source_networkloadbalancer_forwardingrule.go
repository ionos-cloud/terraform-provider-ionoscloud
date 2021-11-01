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

func dataSourceNetworkLoadBalancerForwardingRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkLoadBalancerForwardingRuleRead,
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

func dataSourceNetworkLoadBalancerForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return errors.New("no datacenter_id was specified")
	}

	networkloadbalancerId, nlbIdOk := d.GetOk("networkloadbalancer_id")
	if !nlbIdOk {
		return errors.New("no networkloadbalancer_id was specified")
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the network loadbalancer forwarding rule id or name")
	}
	var networkLoadBalancerForwardingRule ionoscloud.NetworkLoadBalancerForwardingRule
	var err error
	var apiResponse *ionoscloud.APIResponse

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if idOk {
		/* search by ID */
		networkLoadBalancerForwardingRule, apiResponse, err = client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId.(string), networkloadbalancerId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the network loadbalancer forwarding rule %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var networkLoadBalancerForwardingRules ionoscloud.NetworkLoadBalancerForwardingRules

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		networkLoadBalancerForwardingRules, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesGet(ctx, datacenterId.(string), networkloadbalancerId.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching network loadbalancers forwarding rules: %s", err.Error())
		}

		if networkLoadBalancerForwardingRules.Items != nil {
			for _, c := range *networkLoadBalancerForwardingRules.Items {
				tmpNetworkLoadBalancerForwardingRule, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId.(string), networkloadbalancerId.(string), *c.Id).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return fmt.Errorf("an error occurred while fetching network loadbalancer forwarding rule with ID %s: %s", *c.Id, err.Error())
				}
				if tmpNetworkLoadBalancerForwardingRule.Properties.Name != nil {
					if strings.Contains(*tmpNetworkLoadBalancerForwardingRule.Properties.Name, name.(string)) {
						networkLoadBalancerForwardingRule = tmpNetworkLoadBalancerForwardingRule
						break
					}
				}

			}
		}

	}

	if &networkLoadBalancerForwardingRule == nil {
		return errors.New("network loadbalancer not found")
	}

	if networkLoadBalancerForwardingRule.Id != nil {
		if err := d.Set("id", *networkLoadBalancerForwardingRule.Id); err != nil {
			return err
		}
	}

	if err = setNetworkLoadBalancerForwardingRuleData(d, &networkLoadBalancerForwardingRule); err != nil {
		return err
	}

	return nil
}

func setNetworkLoadBalancerForwardingRuleData(d *schema.ResourceData, networkLoadBalancerForwardingRule *ionoscloud.NetworkLoadBalancerForwardingRule) error {

	if networkLoadBalancerForwardingRule.Id != nil {
		d.SetId(*networkLoadBalancerForwardingRule.Id)
	}

	if networkLoadBalancerForwardingRule.Properties != nil {

		if networkLoadBalancerForwardingRule.Properties.Name != nil {
			err := d.Set("name", *networkLoadBalancerForwardingRule.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for network load balancer forwarding rule %s: %s", d.Id(), err)
			}
		}

		if networkLoadBalancerForwardingRule.Properties.Algorithm != nil {
			err := d.Set("algorithm", *networkLoadBalancerForwardingRule.Properties.Algorithm)
			if err != nil {
				return fmt.Errorf("error while setting algorithm property for network load balancer forwarding rule %s: %s", d.Id(), err)
			}
		}

		if networkLoadBalancerForwardingRule.Properties.Protocol != nil {
			err := d.Set("protocol", *networkLoadBalancerForwardingRule.Properties.Protocol)
			if err != nil {
				return fmt.Errorf("error while setting protocol property for network load balancer forwarding rule %s: %s", d.Id(), err)
			}
		}

		if networkLoadBalancerForwardingRule.Properties.ListenerIp != nil {
			err := d.Set("listener_ip", *networkLoadBalancerForwardingRule.Properties.ListenerIp)
			if err != nil {
				return fmt.Errorf("error while setting listener_ip property for network load balancer forwarding rule %s: %s", d.Id(), err)
			}
		}

		if networkLoadBalancerForwardingRule.Properties.ListenerPort != nil {
			err := d.Set("listener_port", *networkLoadBalancerForwardingRule.Properties.ListenerPort)
			if err != nil {
				return fmt.Errorf("error while setting listener_port property for network load balancer forwarding rule %s: %s", d.Id(), err)
			}
		}

		if networkLoadBalancerForwardingRule.Properties.HealthCheck != nil {
			var healthCheck []interface{}

			healthCheckEntry := make(map[string]interface{})
			if networkLoadBalancerForwardingRule.Properties.HealthCheck.ClientTimeout != nil {
				healthCheckEntry["client_timeout"] = *networkLoadBalancerForwardingRule.Properties.HealthCheck.ClientTimeout
			}

			if networkLoadBalancerForwardingRule.Properties.HealthCheck.ConnectTimeout != nil {
				healthCheckEntry["connect_timeout"] = *networkLoadBalancerForwardingRule.Properties.HealthCheck.ConnectTimeout
			}

			if networkLoadBalancerForwardingRule.Properties.HealthCheck.TargetTimeout != nil {
				healthCheckEntry["target_timeout"] = *networkLoadBalancerForwardingRule.Properties.HealthCheck.TargetTimeout
			}

			if networkLoadBalancerForwardingRule.Properties.HealthCheck.Retries != nil {
				healthCheckEntry["retries"] = *networkLoadBalancerForwardingRule.Properties.HealthCheck.Retries
			}

			healthCheck = append(healthCheck, healthCheckEntry)

			err := d.Set("health_check", healthCheck)
			if err != nil {
				return fmt.Errorf("error while setting health_check property for network load balancer forwarding rule %s: %s", d.Id(), err)
			}

		}

		if networkLoadBalancerForwardingRule.Properties.Targets != nil && len(*networkLoadBalancerForwardingRule.Properties.Targets) > 0 {
			var forwardingRuleTargets []interface{}
			for _, target := range *networkLoadBalancerForwardingRule.Properties.Targets {
				targetEntry := make(map[string]interface{})

				if target.Ip != nil {
					targetEntry["ip"] = *target.Ip
				}

				if target.Port != nil {
					targetEntry["port"] = *target.Port
				}

				if target.Weight != nil {
					targetEntry["weight"] = *target.Weight
				}

				if target.HealthCheck != nil {
					var healthCheck []interface{}

					healthCheckEntry := make(map[string]interface{})

					if target.HealthCheck.Check != nil {
						healthCheckEntry["check"] = *target.HealthCheck.Check
					}

					if target.HealthCheck.CheckInterval != nil {
						healthCheckEntry["check_interval"] = *target.HealthCheck.CheckInterval
					}

					if target.HealthCheck.Maintenance != nil {
						healthCheckEntry["maintenance"] = *target.HealthCheck.Maintenance
					}

					healthCheck = append(healthCheck, healthCheckEntry)
					targetEntry["health_check"] = healthCheck
				}

				forwardingRuleTargets = append(forwardingRuleTargets, targetEntry)
			}

			if len(forwardingRuleTargets) > 0 {
				if err := d.Set("targets", forwardingRuleTargets); err != nil {
					return fmt.Errorf("error while setting targets property for network load balancer forwarding rule  %s: %s", d.Id(), err)
				}
			}
		}

	}
	return nil
}
