package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
)

func resourceNetworkLoadBalancerForwardingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkLoadBalancerForwardingRuleCreate,
		ReadContext:   resourceNetworkLoadBalancerForwardingRuleRead,
		UpdateContext: resourceNetworkLoadBalancerForwardingRuleUpdate,
		DeleteContext: resourceNetworkLoadBalancerForwardingRuleDelete,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Description:  "A name of that Network Load Balancer forwarding rule",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"algorithm": {
				Type:         schema.TypeString,
				Description:  "Algorithm for the balancing.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "Protocol of the balancing.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"listener_ip": {
				Type:         schema.TypeString,
				Description:  "Listening IP. (inbound)",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"listener_port": {
				Type:        schema.TypeInt,
				Description: "Listening port number. (inbound) (range: 1 to 65535)",
				Required:    true,
			},
			"health_check": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Health check attributes for Network Load Balancer forwarding rule",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_timeout": {
							Type: schema.TypeInt,
							Description: "ClientTimeout is expressed in milliseconds. This inactivity timeout applies " +
								"when the client is expected to acknowledge or send data. If unset the default of 50 " +
								"seconds will be used.",
							Optional: true,
							Computed: true,
						},
						"connect_timeout": {
							Type: schema.TypeInt,
							Description: "It specifies the maximum time (in milliseconds) to wait for a connection " +
								"attempt to a target VM to succeed. If unset, the default of 5 seconds will be used.",
							Optional: true,
							Computed: true,
						},
						"target_timeout": {
							Type: schema.TypeInt,
							Description: "TargetTimeout specifies the maximum inactivity time (in milliseconds) on the " +
								"target VM side. If unset, the default of 50 seconds will be used.",
							Optional: true,
							Computed: true,
						},
						"retries": {
							Type: schema.TypeInt,
							Description: "Retries specifies the number of retries to perform on a target VM after a " +
								"connection failure. If unset, the default value of 3 will be used.",
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"targets": {
				Type:        schema.TypeList,
				Description: "Array of items in that collection",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Description: "IP of a balanced target VM",
							Required:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "Port of the balanced target service. (range: 1 to 65535)",
							Required:    true,
						},
						"weight": {
							Type:        schema.TypeInt,
							Description: "Weight parameter is used to adjust the target VM's weight relative to other target VMs",
							Required:    true,
						},
						"health_check": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Health check attributes for Network Load Balancer forwarding rule target",
							Optional:    true,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"check": {
										Type:        schema.TypeBool,
										Description: "Check specifies whether the target VM's health is checked.",
										Optional:    true,
										Computed:    true,
									},
									"check_interval": {
										Type: schema.TypeInt,
										Description: "CheckInterval determines the duration (in milliseconds) between " +
											"consecutive health checks. If unspecified a default of 2000 ms is used.",
										Optional: true,
										Computed: true,
									},
									"maintenance": {
										Type:        schema.TypeBool,
										Description: "Maintenance specifies if a target VM should be marked as down, even if it is not.",
										Optional:    true,
										Computed:    true,
									},
								},
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
			"networkloadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNetworkLoadBalancerForwardingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	networkLoadBalancerForwardingRule := ionoscloud.NetworkLoadBalancerForwardingRule{
		Properties: &ionoscloud.NetworkLoadBalancerForwardingRuleProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		networkLoadBalancerForwardingRule.Properties.Name = &name
	} else {
		diags := diag.FromErr(fmt.Errorf("name must be provided for network loadbalancer forwarding rule"))
		return diags
	}

	if algorithm, algorithmOk := d.GetOk("algorithm"); algorithmOk {
		algorithm := algorithm.(string)
		networkLoadBalancerForwardingRule.Properties.Algorithm = &algorithm
	} else {
		diags := diag.FromErr(fmt.Errorf("algorithm must be provided for network loadbalancer forwarding rule"))
		return diags
	}

	if protocol, protocolOk := d.GetOk("protocol"); protocolOk {
		protocol := protocol.(string)
		networkLoadBalancerForwardingRule.Properties.Protocol = &protocol
	} else {
		diags := diag.FromErr(fmt.Errorf("protocol must be provided for network loadbalancer forwarding rule"))
		return diags
	}

	if listenerIp, listenerIpOk := d.GetOk("listener_ip"); listenerIpOk {
		listenerIp := listenerIp.(string)
		networkLoadBalancerForwardingRule.Properties.ListenerIp = &listenerIp
	} else {
		diags := diag.FromErr(fmt.Errorf("listner ip must be provided for network loadbalancer forwarding rule"))
		return diags
	}

	if listenerPort, listenerPortOk := d.GetOk("listener_port"); listenerPortOk {
		listenerPort := int32(listenerPort.(int))
		networkLoadBalancerForwardingRule.Properties.ListenerPort = &listenerPort
	} else {
		diags := diag.FromErr(fmt.Errorf("listner port must be provided for network loadbalancer forwarding rule"))
		return diags
	}

	if _, healthCheckOk := d.GetOk("health_check.0"); healthCheckOk {
		networkLoadBalancerForwardingRule.Properties.HealthCheck = &ionoscloud.NetworkLoadBalancerForwardingRuleHealthCheck{}

		if clientTimeout, clientTimeoutOk := d.GetOk("health_check.0.client_timeout"); clientTimeoutOk {
			clientTimeout := int32(clientTimeout.(int))
			networkLoadBalancerForwardingRule.Properties.HealthCheck.ClientTimeout = &clientTimeout
		}

		if connectTimeout, connectTimeoutOk := d.GetOk("health_check.0.connect_timeout"); connectTimeoutOk {
			connectTimeout := int32(connectTimeout.(int))
			networkLoadBalancerForwardingRule.Properties.HealthCheck.ConnectTimeout = &connectTimeout
		}

		if targetTimeout, targetTimeoutOk := d.GetOk("health_check.0.target_timeout"); targetTimeoutOk {
			targetTimeout := int32(targetTimeout.(int))
			networkLoadBalancerForwardingRule.Properties.HealthCheck.TargetTimeout = &targetTimeout
		}

		if retries, retriesOk := d.GetOk("health_check.0.retries"); retriesOk {
			retries := int32(retries.(int))
			networkLoadBalancerForwardingRule.Properties.HealthCheck.Retries = &retries
		}

	}

	if targetsVal, targetsOk := d.GetOk("targets"); targetsOk {
		if targetsVal.([]interface{}) != nil {
			updateTargets := false

			var targets []ionoscloud.NetworkLoadBalancerForwardingRuleTarget

			for targetIndex := range targetsVal.([]interface{}) {
				target := ionoscloud.NetworkLoadBalancerForwardingRuleTarget{}
				addTarget := false
				if ip, ipOk := d.GetOk(fmt.Sprintf("targets.%d.ip", targetIndex)); ipOk {
					ip := ip.(string)
					target.Ip = &ip
					addTarget = true
				} else {
					diags := diag.FromErr(fmt.Errorf("ip must be provided for network loadbalancer forwarding rule target"))
					return diags
				}

				if port, portOk := d.GetOk(fmt.Sprintf("targets.%d.port", targetIndex)); portOk {
					port := int32(port.(int))
					target.Port = &port
				} else {
					diags := diag.FromErr(fmt.Errorf("port must be provided for network loadbalancer forwarding rule target"))
					return diags
				}

				if weight, weightOk := d.GetOk(fmt.Sprintf("targets.%d.weight", targetIndex)); weightOk {
					weight := int32(weight.(int))
					target.Weight = &weight
				} else {
					diags := diag.FromErr(fmt.Errorf("weight must be provided for network loadbalancer forwarding rule target"))
					return diags
				}

				if _, healthCheckOk := d.GetOk(fmt.Sprintf("targets.%d.health_check.0", targetIndex)); healthCheckOk {
					target.HealthCheck = &ionoscloud.NetworkLoadBalancerForwardingRuleTargetHealthCheck{}

					if check, checkOk := d.GetOk(fmt.Sprintf("targets.%d.health_check.0.check", targetIndex)); checkOk {
						check := check.(bool)
						target.HealthCheck.Check = &check
					}

					if checkInterval, checkIntervalOk := d.GetOk(fmt.Sprintf("targets.%d.health_check.0.check_interval", targetIndex)); checkIntervalOk {
						checkInterval := int32(checkInterval.(int))
						target.HealthCheck.CheckInterval = &checkInterval
					}
					if maintenance, maintenanceOk := d.GetOk(fmt.Sprintf("targets.%d.health_check.0.maintenance", targetIndex)); maintenanceOk {
						maintenance := maintenance.(bool)
						target.HealthCheck.Maintenance = &maintenance
					}

				}

				if addTarget {
					targets = append(targets, target)
				}

			}

			if len(targets) > 0 {
				updateTargets = true
			}

			if updateTargets == true {
				log.Printf("[INFO] Network load balancer forwarding rule targets set to %+v", targets)
				networkLoadBalancerForwardingRule.Properties.Targets = &targets
			}
		}
	}

	dcId := d.Get("datacenter_id").(string)
	nlbId := d.Get("networkloadbalancer_id").(string)

	networkLoadBalancerForwardingRuleResp, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesPost(ctx, dcId, nlbId).NetworkLoadBalancerForwardingRule(networkLoadBalancerForwardingRule).Execute()

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating network loadbalancer: %s \n ApiError: %s", err, responseBody(apiResponse)))
		return diags
	}

	d.SetId(*networkLoadBalancerForwardingRuleResp.Id)

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

	return resourceNetworkLoadBalancerForwardingRuleRead(ctx, d, meta)
}

func resourceNetworkLoadBalancerForwardingRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)

	nlbID := d.Get("networkloadbalancer_id").(string)

	networkLoadBalancerForwardingRule, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleId(ctx, dcId, nlbID, d.Id()).Execute()

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] Successfully retreived network load balancer forwarding rule %s: %+v", d.Id(), networkLoadBalancerForwardingRule)

	if networkLoadBalancerForwardingRule.Properties.Name != nil {
		err := d.Set("name", *networkLoadBalancerForwardingRule.Properties.Name)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting name property for network load balancer forwarding rule %s: %s", d.Id(), err))
			return diags
		}
	}

	if networkLoadBalancerForwardingRule.Properties.Algorithm != nil {
		err := d.Set("algorithm", *networkLoadBalancerForwardingRule.Properties.Algorithm)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting algorithm property for network load balancer forwarding rule %s: %s", d.Id(), err))
			return diags
		}
	}

	if networkLoadBalancerForwardingRule.Properties.Protocol != nil {
		err := d.Set("protocol", *networkLoadBalancerForwardingRule.Properties.Protocol)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting protocol property for network load balancer forwarding rule %s: %s", d.Id(), err))
			return diags
		}
	}

	if networkLoadBalancerForwardingRule.Properties.ListenerIp != nil {
		err := d.Set("listener_ip", *networkLoadBalancerForwardingRule.Properties.ListenerIp)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting listener_ip property for network load balancer forwarding rule %s: %s", d.Id(), err))
			return diags
		}
	}

	if networkLoadBalancerForwardingRule.Properties.ListenerPort != nil {
		err := d.Set("listener_port", *networkLoadBalancerForwardingRule.Properties.ListenerPort)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting listener_port property for network load balancer forwarding rule %s: %s", d.Id(), err))
			return diags
		}
	}

	if networkLoadBalancerForwardingRule.Properties.HealthCheck != nil {
		healthCheck := make([]interface{}, 1)

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

		healthCheck[0] = healthCheckEntry
		err := d.Set("health_check", healthCheck)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting health_check property for network load balancer forwarding rule %s: %s", d.Id(), err))
			return diags
		}

	}

	forwardingRuleTargets := make([]interface{}, 0)
	if networkLoadBalancerForwardingRule.Properties.Targets != nil && len(*networkLoadBalancerForwardingRule.Properties.Targets) > 0 {
		forwardingRuleTargets = make([]interface{}, len(*networkLoadBalancerForwardingRule.Properties.Targets))
		for targetIndex, target := range *networkLoadBalancerForwardingRule.Properties.Targets {
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
				healthCheck := make([]interface{}, 1)

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

				healthCheck[0] = healthCheckEntry
				targetEntry["health_check"] = healthCheck
			}

			forwardingRuleTargets[targetIndex] = targetEntry
		}
	}

	if len(forwardingRuleTargets) > 0 {
		if err := d.Set("targets", forwardingRuleTargets); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting targets property for network load balancer forwarding rule  %s: %s", d.Id(), err))
			return diags
		}
	}

	return nil
}

func resourceNetworkLoadBalancerForwardingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	request := ionoscloud.NetworkLoadBalancerForwardingRule{
		Properties: &ionoscloud.NetworkLoadBalancerForwardingRuleProperties{},
	}

	dcId := d.Get("datacenter_id").(string)
	nlbID := d.Get("networkloadbalancer_id").(string)

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		vStr := v.(string)
		request.Properties.Name = &vStr
	}

	if d.HasChange("algorithm") {
		_, v := d.GetChange("algorithm")
		vStr := v.(string)
		request.Properties.Algorithm = &vStr
	}

	if d.HasChange("protocol") {
		_, v := d.GetChange("protocol")
		vStr := v.(string)
		request.Properties.Protocol = &vStr
	}

	if d.HasChange("listener_ip") {
		_, v := d.GetChange("listener_ip")
		vStr := v.(string)
		request.Properties.ListenerIp = &vStr
	}

	if d.HasChange("listener_port") {
		_, v := d.GetChange("listener_port")
		vStr := int32(v.(int))
		request.Properties.ListenerPort = &vStr
	}

	if d.HasChange("health_check.0") {
		_, v := d.GetChange("health_check.0")
		if v.(map[string]interface{}) != nil {
			updateHealthCheck := false

			healthCheck := &ionoscloud.NetworkLoadBalancerForwardingRuleHealthCheck{}

			if d.HasChange("health_check.0.client_timeout") {
				_, newValue := d.GetChange("health_check.0.client_timeout")
				if newValue != 0 {
					updateHealthCheck = true
					newValue := int32(newValue.(int))
					healthCheck.ClientTimeout = &newValue
				}
			}

			if d.HasChange("health_check.0.connect_timeout") {
				_, newValue := d.GetChange("health_check.0.connect_timeout")
				if newValue != 0 {
					updateHealthCheck = true
					newValue := int32(newValue.(int))
					healthCheck.ConnectTimeout = &newValue
				}
			}

			if d.HasChange("health_check.0.target_timeout") {
				_, newValue := d.GetChange("health_check.0.target_timeout")
				if newValue != 0 {
					updateHealthCheck = true
					newValue := int32(newValue.(int))
					healthCheck.TargetTimeout = &newValue
				}
			}

			if d.HasChange("health_check.0.retries") {
				_, newValue := d.GetChange("health_check.0.retries")
				if newValue != 0 {
					updateHealthCheck = true
					newValue := int32(newValue.(int))
					healthCheck.Retries = &newValue
				}
			}

			if updateHealthCheck == true {
				request.Properties.HealthCheck = healthCheck
			}
		}
	}

	if d.HasChange("targets") {
		oldTargets, newTargets := d.GetChange("targets")
		if newTargets.([]interface{}) != nil {
			updateTargets := false

			var targets []ionoscloud.NetworkLoadBalancerForwardingRuleTarget

			for targetIndex := range newTargets.([]interface{}) {
				target := ionoscloud.NetworkLoadBalancerForwardingRuleTarget{}

				if ip, ipOk := d.GetOk(fmt.Sprintf("targets.%d.ip", targetIndex)); ipOk {
					ip := ip.(string)
					target.Ip = &ip
				}

				if port, portOk := d.GetOk(fmt.Sprintf("targets.%d.port", targetIndex)); portOk {
					port := int32(port.(int))
					target.Port = &port
				}

				if weight, weightOk := d.GetOk(fmt.Sprintf("targets.%d.weight", targetIndex)); weightOk {
					weight := int32(weight.(int))
					target.Weight = &weight
				}

				if _, healthCheckOk := d.GetOk(fmt.Sprintf("targets.%d.health_check.0", targetIndex)); healthCheckOk {
					target.HealthCheck = &ionoscloud.NetworkLoadBalancerForwardingRuleTargetHealthCheck{}

					if check, checkOk := d.GetOk(fmt.Sprintf("targets.%d.health_check.0.check", targetIndex)); checkOk {
						check := check.(bool)
						target.HealthCheck.Check = &check
					}

					if checkInterval, checkIntervalOk := d.GetOk(fmt.Sprintf("targets.%d.health_check.0.check_interval", targetIndex)); checkIntervalOk {
						checkInterval := int32(checkInterval.(int))
						target.HealthCheck.CheckInterval = &checkInterval
					}

					if maintenance, maintenanceOk := d.GetOk(fmt.Sprintf("targets.%d.health_check.0.maintenance", targetIndex)); maintenanceOk {
						maintenance := maintenance.(bool)
						target.HealthCheck.Maintenance = &maintenance
					}
				}

				targets = append(targets, target)
			}

			if len(targets) > 0 {
				updateTargets = true
			}

			if updateTargets == true {
				log.Printf("[INFO] Network load balancer forwarding rule targets changed from %+v to %+v", oldTargets, newTargets)
				request.Properties.Targets = &targets
			}
		}
	}
	_, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesPatch(ctx, dcId, nlbID, d.Id()).NetworkLoadBalancerForwardingRuleProperties(*request.Properties).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a network loadbalancer forwarding rule ID %s %s \n ApiError: %s",
			d.Id(), err, responseBody(apiResponse)))
		return diags
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceNetworkLoadBalancerForwardingRuleRead(ctx, d, meta)
}

func resourceNetworkLoadBalancerForwardingRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)
	nlbID := d.Get("networkloadbalancer_id").(string)

	apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesDelete(ctx, dcId, nlbID, d.Id()).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a network loadbalancer forwarding rule %s %s", d.Id(), err))
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
