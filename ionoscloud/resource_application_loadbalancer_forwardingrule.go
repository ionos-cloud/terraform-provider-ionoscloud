package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
)

func resourceApplicationLoadBalancerForwardingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationLoadBalancerForwardingRuleCreate,
		ReadContext:   resourceApplicationLoadBalancerForwardingRuleRead,
		UpdateContext: resourceApplicationLoadBalancerForwardingRuleUpdate,
		DeleteContext: resourceApplicationLoadBalancerForwardingRuleDelete,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Description: "A name of that Application Load Balancer forwarding rule",
				Required:    true,
			},
			"protocol": {
				Type:        schema.TypeString,
				Description: "rotocol of the balancing.",
				Required:    true,
			},
			"listener_ip": {
				Type:        schema.TypeString,
				Description: "Listening IP. (inbound)",
				Required:    true,
			},
			"listener_port": {
				Type:        schema.TypeInt,
				Description: "Listening port number. (inbound) (range: 1 to 65535)",
				Required:    true,
			},
			"health_check": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Health check attributes for Application Load Balancer forwarding rule",
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
					},
				},
			},
			"server_certificates": {
				Type:        schema.TypeList,
				Description: "Array of items in that collection.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"http_rules": {
				Type:        schema.TypeList,
				Description: "Array of items in that collection",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "A name of that Application Load Balancer http rule",
							Required:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Type of the Http Rule",
							Required:    true,
						},
						"target_group": {
							Type:        schema.TypeString,
							Description: "The UUID of the target group; mandatory for FORWARD action",
							Optional:    true,
						},
						"drop_query": {
							Type:        schema.TypeBool,
							Description: "Default is false; true for REDIRECT action.",
							Optional:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "The location for redirecting; mandatory for REDIRECT action",
							Optional:    true,
						},
						"status_code": {
							Type:        schema.TypeInt,
							Description: "On REDIRECT action it can take the value 301, 302, 303, 307, 308; on STATIC action it is between 200 and 599",
							Optional:    true,
						},
						"response_message": {
							Type:        schema.TypeString,
							Description: "he response message of the request; mandatory for STATIC action",
							Optional:    true,
						},
						"content_type": {
							Type:        schema.TypeString,
							Description: "Will be provided by the PAAS Team; default application/json",
							Optional:    true,
						},
						"conditions": {
							Type:        schema.TypeList,
							Description: "Array of items in that collection",
							Optional:    true,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Description: "Type of the Http Rule condition.",
										Required:    true,
									},
									"condition": {
										Type:        schema.TypeString,
										Description: "Condition of the Http Rule condition.",
										Required:    true,
									},
									"negate": {
										Type:        schema.TypeBool,
										Description: "Specifies whether the condition is negated or not; default: false.",
										Optional:    true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
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
			"application_loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceApplicationLoadBalancerForwardingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	applicationLoadBalancerForwardingRule := ionoscloud.ApplicationLoadBalancerForwardingRule{
		Properties: &ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		applicationLoadBalancerForwardingRule.Properties.Name = &name
	} else {
		diags := diag.FromErr(fmt.Errorf("name must be provided for application loadbalancer forwarding rule"))
		return diags
	}

	if protocol, protocolOk := d.GetOk("protocol"); protocolOk {
		protocol := protocol.(string)
		applicationLoadBalancerForwardingRule.Properties.Protocol = &protocol
	} else {
		diags := diag.FromErr(fmt.Errorf("protocol must be provided for application loadbalancer forwarding rule"))
		return diags
	}

	if listenerIp, listenerIpOk := d.GetOk("listener_ip"); listenerIpOk {
		listenerIp := listenerIp.(string)
		applicationLoadBalancerForwardingRule.Properties.ListenerIp = &listenerIp
	} else {
		diags := diag.FromErr(fmt.Errorf("listner ip must be provided for application loadbalancer forwarding rule"))
		return diags
	}

	if listenerPort, listenerPortOk := d.GetOk("listener_port"); listenerPortOk {
		listenerPort := int32(listenerPort.(int))
		applicationLoadBalancerForwardingRule.Properties.ListenerPort = &listenerPort
	} else {
		diags := diag.FromErr(fmt.Errorf("listner port must be provided for application loadbalancer forwarding rule"))
		return diags
	}

	//if _, healthCheckOk := d.GetOk("health_check.0"); healthCheckOk {
	//	applicationLoadBalancerForwardingRule.Properties.HealthCheck = &ionoscloud.ApplicationLoadBalancerForwardingRuleHealthCheck{}
	//
	//	if clientTimeout, clientTimeoutOk := d.GetOk("health_check.0.client_timeout"); clientTimeoutOk {
	//		clientTimeout := int32(clientTimeout.(int))
	//		applicationLoadBalancerForwardingRule.Properties.HealthCheck.ClientTimeout = &clientTimeout
	//	}
	//
	//}

	if serverCertificatesVal, serverCertificatesOk := d.GetOk("server_certificates"); serverCertificatesOk {
		serverCertificatesVal := serverCertificatesVal.([]interface{})
		if serverCertificatesVal != nil {
			serverCertificates := make([]string, 0)
			for idx, value := range serverCertificatesVal {
				serverCertificates[idx] = value.(string)
			}
			if len(serverCertificates) > 0 {
				applicationLoadBalancerForwardingRule.Properties.ServerCertificates = &serverCertificates
			}
		}
	}

	if httpRulesVal, httpRulesOk := d.GetOk("http_rules"); httpRulesOk {
		if httpRulesVal.([]interface{}) != nil {
			addHttpRules := false

			var httpRules []ionoscloud.ApplicationLoadBalancerHttpRule

			for httpRuleIndex := range httpRulesVal.([]interface{}) {
				httpRule := ionoscloud.ApplicationLoadBalancerHttpRule{}
				addHttpRule := false
				if name, nameOk := d.GetOk(fmt.Sprintf("http_rules.%d.name", httpRuleIndex)); nameOk {
					name := name.(string)
					httpRule.Name = &name
					addHttpRule = true
				} else {
					diags := diag.FromErr(fmt.Errorf("ip must be provided for application loadbalancer forwarding rule http_rule"))
					return diags
				}

				if typeVal, typeOk := d.GetOk(fmt.Sprintf("http_rules.%d.type", httpRuleIndex)); typeOk {
					typeVal := typeVal.(string)
					httpRule.Type = &typeVal
				} else {
					diags := diag.FromErr(fmt.Errorf("type must be provided for application loadbalancer forwarding rule http_rule"))
					return diags
				}

				if targetGroup, targetGroupOk := d.GetOk(fmt.Sprintf("http_rules.%d.target_group", httpRuleIndex)); targetGroupOk {
					targetGroup := targetGroup.(string)
					httpRule.TargetGroup = &targetGroup
				}

				if dropQuery, dropQueryOk := d.GetOk(fmt.Sprintf("http_rules.%d.drop_query", httpRuleIndex)); dropQueryOk {
					dropQuery := dropQuery.(bool)
					httpRule.DropQuery = &dropQuery
				}

				if location, locationOk := d.GetOk(fmt.Sprintf("http_rules.%d.location", httpRuleIndex)); locationOk {
					location := location.(string)
					httpRule.Location = &location
				}

				if statusCode, statusCodeOk := d.GetOk(fmt.Sprintf("http_rules.%d.status_code", httpRuleIndex)); statusCodeOk {
					statusCode := int32(statusCode.(int))
					httpRule.StatusCode = &statusCode
				}

				if responseMessage, responseMessageOk := d.GetOk(fmt.Sprintf("http_rules.%d.response_message", httpRuleIndex)); responseMessageOk {
					responseMessage := responseMessage.(string)
					httpRule.ResponseMessage = &responseMessage
				}

				if contentType, contentTypeOk := d.GetOk(fmt.Sprintf("http_rules.%d.content_type", httpRuleIndex)); contentTypeOk {
					contentType := contentType.(string)
					httpRule.ContentType = &contentType
				}

				if conditionsVal, conditionsOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions", httpRuleIndex)); conditionsOk {
					if conditionsVal.([]interface{}) != nil {

						addConditions := false
						var conditions []ionoscloud.ApplicationLoadBalancerHttpRuleCondition

						for conditionIndex := range conditionsVal.([]interface{}) {

							condition := ionoscloud.ApplicationLoadBalancerHttpRuleCondition{}
							addCondition := false

							if typeVal, typeOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions.%d.type", httpRuleIndex, conditionIndex)); typeOk {
								typeVal := typeVal.(string)
								condition.Type = &typeVal
								addCondition = true
							} else {
								diags := diag.FromErr(fmt.Errorf("type must be provided for application loadbalancer forwarding rule http rule condition"))
								return diags
							}

							if conditionVal, conditionOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions.%d.condition", httpRuleIndex, conditionIndex)); conditionOk {
								conditionVal := conditionVal.(string)
								condition.Condition = &conditionVal
							} else {
								diags := diag.FromErr(fmt.Errorf("condition must be provided for application loadbalancer forwarding rule http rule condition"))
								return diags
							}

							if negate, negateOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions.%d.negate", httpRuleIndex, conditionIndex)); negateOk {
								negate := negate.(bool)
								condition.Negate = &negate
							}

							if key, keyOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions.%d.key", httpRuleIndex, conditionIndex)); keyOk {
								key := key.(string)
								condition.Key = &key
							}

							if value, valueOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions.%d.value", httpRuleIndex, conditionIndex)); valueOk {
								value := value.(string)
								condition.Value = &value
							}

							if addCondition {
								conditions = append(conditions, condition)
							}
						}

						if len(conditions) > 0 {
							addConditions = true
						}

						if addConditions {
							httpRule.Conditions = &conditions
						}

					}

				}

				if addHttpRule {
					httpRules = append(httpRules, httpRule)
				}

			}

			if len(httpRules) > 0 {
				addHttpRules = true
			}

			if addHttpRules == true {
				log.Printf("[INFO] Application load balancer forwarding rule httpRules set to %+v", httpRules)
				applicationLoadBalancerForwardingRule.Properties.HttpRules = &httpRules
			}
		}
	}

	dcId := d.Get("datacenter_id").(string)
	albId := d.Get("application_loadbalancer_id").(string)

	albForwardingRuleResp, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPost(ctx, dcId, albId).ApplicationLoadBalancerForwardingRule(applicationLoadBalancerForwardingRule).Execute()

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating application loadbalancer forwarding rule: %s \n ApiError: %s", err, responseBody(apiResponse)))
		return diags
	}

	d.SetId(*albForwardingRuleResp.Id)

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

	return resourceApplicationLoadBalancerForwardingRuleRead(ctx, d, meta)
}

func resourceApplicationLoadBalancerForwardingRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)

	albId := d.Get("application_loadbalancer_id").(string)

	applicationLoadBalancerForwardingRule, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(ctx, dcId, albId, d.Id()).Execute()

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] Successfully retreived application load balancer forwarding rule %s: %+v", d.Id(), applicationLoadBalancerForwardingRule)

	if applicationLoadBalancerForwardingRule.Properties.Name != nil {
		err := d.Set("name", *applicationLoadBalancerForwardingRule.Properties.Name)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting name property for application load balancer forwarding rule %s: %s", d.Id(), err))
			return diags
		}
	}

	if applicationLoadBalancerForwardingRule.Properties.Protocol != nil {
		err := d.Set("protocol", *applicationLoadBalancerForwardingRule.Properties.Protocol)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting protocol property for application load balancer forwarding rule %s: %s", d.Id(), err))
			return diags
		}
	}

	if applicationLoadBalancerForwardingRule.Properties.ListenerIp != nil {
		err := d.Set("listener_ip", *applicationLoadBalancerForwardingRule.Properties.ListenerIp)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting listener_ip property for application load balancer forwarding rule %s: %s", d.Id(), err))
			return diags
		}
	}

	if applicationLoadBalancerForwardingRule.Properties.ListenerPort != nil {
		err := d.Set("listener_port", *applicationLoadBalancerForwardingRule.Properties.ListenerPort)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting listener_port property for application load balancer forwarding rule %s: %s", d.Id(), err))
			return diags
		}
	}

	//if applicationLoadBalancerForwardingRule.Properties.HealthCheck != nil {
	//	healthCheck := make([]interface{}, 1)
	//
	//	healthCheckEntry := make(map[string]interface{})
	//	if applicationLoadBalancerForwardingRule.Properties.HealthCheck.ClientTimeout != nil {
	//		healthCheckEntry["client_timeout"] = *applicationLoadBalancerForwardingRule.Properties.HealthCheck.ClientTimeout
	//	}
	//	healthCheck[0] = healthCheckEntry
	//	err := d.Set("health_check", healthCheck)
	//	if err != nil {
	//		diags := diag.FromErr(fmt.Errorf("error while setting health_check property for application load balancer forwarding rule %s: %s", d.Id(), err))
	//		return diags
	//	}
	//}

	if applicationLoadBalancerForwardingRule.Properties.ServerCertificates != nil {
		err := d.Set("server_certificates", *applicationLoadBalancerForwardingRule.Properties.ServerCertificates)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting server_certificates property for application load balancer forwarding rule %s: %s", d.Id(), err))
			return diags
		}
	}

	httpRules := make([]interface{}, 0)
	if applicationLoadBalancerForwardingRule.Properties.HttpRules != nil && len(*applicationLoadBalancerForwardingRule.Properties.HttpRules) > 0 {
		httpRules = make([]interface{}, 0)
		for _, rule := range *applicationLoadBalancerForwardingRule.Properties.HttpRules {
			ruleEntry := make(map[string]interface{})

			if rule.Name != nil {
				ruleEntry["name"] = *rule.Name
			}

			if rule.Type != nil {
				ruleEntry["type"] = *rule.Type
			}

			if rule.TargetGroup != nil {
				ruleEntry["target_group"] = *rule.TargetGroup
			}

			if rule.DropQuery != nil {
				ruleEntry["drop_query"] = *rule.DropQuery
			}

			if rule.Location != nil {
				ruleEntry["location"] = *rule.Location
			}

			if rule.StatusCode != nil {
				ruleEntry["status_code"] = *rule.StatusCode
			}

			if rule.ResponseMessage != nil {
				ruleEntry["response_message"] = *rule.ResponseMessage
			}

			if rule.ContentType != nil {
				ruleEntry["content_type"] = *rule.ContentType
			}

			if rule.Conditions != nil {
				conditions := make([]interface{}, 0)
				for _, condition := range *rule.Conditions {
					conditionEntry := make(map[string]interface{})

					if condition.Type != nil {
						conditionEntry["type"] = *condition.Type
					}

					if condition.Condition != nil {
						conditionEntry["condition"] = *condition.Condition
					}

					if condition.Negate != nil {
						conditionEntry["negate"] = *condition.Negate
					}

					if condition.Key != nil {
						conditionEntry["key"] = *condition.Key
					}

					if condition.Value != nil {
						conditionEntry["value"] = *condition.Value
					}

					conditions = append(conditions, conditionEntry)
				}

				ruleEntry["conditions"] = conditions
			}

			httpRules = append(httpRules, ruleEntry)
		}
	}

	if len(httpRules) > 0 {
		if err := d.Set("http_rules", httpRules); err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting http_rules property for application load balancer forwarding rule  %s: %s", d.Id(), err))
			return diags
		}
	}

	return nil
}

func resourceApplicationLoadBalancerForwardingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	request := ionoscloud.ApplicationLoadBalancerForwardingRule{
		Properties: &ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{},
	}

	dcId := d.Get("datacenter_id").(string)
	albId := d.Get("application_loadbalancer_id").(string)

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		vStr := v.(string)
		request.Properties.Name = &vStr
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

	//if d.HasChange("health_check.0") {
	//	_, v := d.GetChange("health_check.0")
	//	if v.(map[string]interface{}) != nil {
	//		updateHealthCheck := false
	//
	//		healthCheck := &ionoscloud.ApplicationLoadBalancerForwardingRuleHealthCheck{}
	//
	//		if d.HasChange("health_check.0.client_timeout") {
	//			_, newValue := d.GetChange("health_check.0.client_timeout")
	//			if newValue != 0 {
	//				updateHealthCheck = true
	//				newValue := int32(newValue.(int))
	//				healthCheck.ClientTimeout = &newValue
	//			}
	//		}
	//
	//		if updateHealthCheck == true {
	//			request.Properties.HealthCheck = healthCheck
	//		}
	//	}
	//}

	if d.HasChange("server_certificates") {
		_, v := d.GetChange("server_certificates")
		if v.([]interface{}) != nil {
			serverCertificates := make([]string, 0)
			for idx, value := range v.([]interface{}) {
				serverCertificates[idx] = value.(string)
			}
			if len(serverCertificates) > 0 {
				request.Properties.ServerCertificates = &serverCertificates
			}
		}
	}

	if d.HasChange("http_rules") {
		_, newHttpRules := d.GetChange("http_rules")
		if newHttpRules.([]interface{}) != nil {
			updateHttpRules := false

			var httpRules []ionoscloud.ApplicationLoadBalancerHttpRule

			for httpRuleIndex := range newHttpRules.([]interface{}) {
				httpRule := ionoscloud.ApplicationLoadBalancerHttpRule{}
				addHttpRule := false
				if name, nameOk := d.GetOk(fmt.Sprintf("http_rules.%d.name", httpRuleIndex)); nameOk {
					name := name.(string)
					httpRule.Name = &name
					addHttpRule = true
				} else {
					diags := diag.FromErr(fmt.Errorf("ip must be provided for application loadbalancer forwarding rule http_rule"))
					return diags
				}

				if typeVal, typeOk := d.GetOk(fmt.Sprintf("http_rules.%d.type", httpRuleIndex)); typeOk {
					typeVal := typeVal.(string)
					httpRule.Type = &typeVal
				} else {
					diags := diag.FromErr(fmt.Errorf("type must be provided for application loadbalancer forwarding rule http_rule"))
					return diags
				}

				if targetGroup, targetGroupOk := d.GetOk(fmt.Sprintf("http_rules.%d.target_group", httpRuleIndex)); targetGroupOk {
					targetGroup := targetGroup.(string)
					httpRule.TargetGroup = &targetGroup
				}

				if dropQuery, dropQueryOk := d.GetOk(fmt.Sprintf("http_rules.%d.drop_query", httpRuleIndex)); dropQueryOk {
					dropQuery := dropQuery.(bool)
					httpRule.DropQuery = &dropQuery
				}

				if location, locationOk := d.GetOk(fmt.Sprintf("http_rules.%d.location", httpRuleIndex)); locationOk {
					location := location.(string)
					httpRule.Location = &location
				}

				if statusCode, statusCodeOk := d.GetOk(fmt.Sprintf("http_rules.%d.status_code", httpRuleIndex)); statusCodeOk {
					statusCode := int32(statusCode.(int))
					httpRule.StatusCode = &statusCode
				}

				if responseMessage, responseMessageOk := d.GetOk(fmt.Sprintf("http_rules.%d.response_message", httpRuleIndex)); responseMessageOk {
					responseMessage := responseMessage.(string)
					httpRule.ResponseMessage = &responseMessage
				}

				if contentType, contentTypeOk := d.GetOk(fmt.Sprintf("http_rules.%d.content_type", httpRuleIndex)); contentTypeOk {
					contentType := contentType.(string)
					httpRule.ContentType = &contentType
				}

				if conditionsVal, conditionsOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions", httpRuleIndex)); conditionsOk {
					if conditionsVal.([]interface{}) != nil {

						addConditions := false
						var conditions []ionoscloud.ApplicationLoadBalancerHttpRuleCondition

						for conditionIndex := range conditionsVal.([]interface{}) {

							condition := ionoscloud.ApplicationLoadBalancerHttpRuleCondition{}
							addCondition := false

							if typeVal, typeOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions.%d.type", httpRuleIndex, conditionIndex)); typeOk {
								typeVal := typeVal.(string)
								condition.Type = &typeVal
								addCondition = true
							} else {
								diags := diag.FromErr(fmt.Errorf("type must be provided for application loadbalancer forwarding rule http rule condition"))
								return diags
							}

							if conditionVal, conditionOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions.%d.condition", httpRuleIndex, conditionIndex)); conditionOk {
								conditionVal := conditionVal.(string)
								condition.Condition = &conditionVal
							} else {
								diags := diag.FromErr(fmt.Errorf("condition must be provided for application loadbalancer forwarding rule http rule condition"))
								return diags
							}

							if negate, negateOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions.%d.negate", httpRuleIndex, conditionIndex)); negateOk {
								negate := negate.(bool)
								condition.Negate = &negate
							}

							if key, keyOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions.%d.key", httpRuleIndex, conditionIndex)); keyOk {
								key := key.(string)
								condition.Key = &key
							}

							if value, valueOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions.%d.value", httpRuleIndex, conditionIndex)); valueOk {
								value := value.(string)
								condition.Value = &value
							}

							if addCondition {
								conditions = append(conditions, condition)
							}
						}

						if len(conditions) > 0 {
							addConditions = true
						}

						if addConditions {
							httpRule.Conditions = &conditions
						}

					}

				}

				if addHttpRule {
					httpRules = append(httpRules, httpRule)
				}

			}

			if len(httpRules) > 0 {
				updateHttpRules = true
			}

			if updateHttpRules == true {
				log.Printf("[INFO] Application load balancer forwarding rule httpRules set to %+v", httpRules)
				request.Properties.HttpRules = &httpRules
			}
		}
	}
	_, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPatch(ctx, dcId, albId, d.Id()).ApplicationLoadBalancerForwardingRuleProperties(*request.Properties).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a application loadbalancer forwarding rule ID %s %s \n ApiError: %s",
			d.Id(), err, responseBody(apiResponse)))
		return diags
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceApplicationLoadBalancerForwardingRuleRead(ctx, d, meta)
}

func resourceApplicationLoadBalancerForwardingRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)
	albID := d.Get("application_loadbalancer_id").(string)

	_, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesDelete(ctx, dcId, albID, d.Id()).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a application loadbalancer forwarding rule %s %s", d.Id(), err))
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
