package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"strings"
)

func resourceApplicationLoadBalancerForwardingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationLoadBalancerForwardingRuleCreate,
		ReadContext:   resourceApplicationLoadBalancerForwardingRuleRead,
		UpdateContext: resourceApplicationLoadBalancerForwardingRuleUpdate,
		DeleteContext: resourceApplicationLoadBalancerForwardingRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationLoadBalancerForwardingRuleImport,
		},
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Description:  "The name of the Application Load Balancer forwarding rule.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "Balancing protocol.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringInSlice([]string{"HTTP"}, true)),
			},
			"listener_ip": {
				Type:         schema.TypeString,
				Description:  "Listening (inbound) IP.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Description:  "Listening (inbound) port number; valid range is 1 to 65535.",
				Required:     true,
				ValidateFunc: validation.All(validation.IntBetween(1, 65535)),
			},
			"client_timeout": {
				Type:        schema.TypeInt,
				Description: "The maximum time in milliseconds to wait for the client to acknowledge or send data; default is 50,000 (50 seconds).",
				Optional:    true,
				Computed:    true,
			},
			"server_certificates": {
				Type:        schema.TypeSet,
				Description: "Array of items in the collection.",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.All(validation.IsUUID),
				},
			},
			"http_rules": {
				Type:        schema.TypeList,
				Description: "Array of items in that collection",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Description:  "The unique name of the Application Load Balancer HTTP rule.",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"type": {
							Type:         schema.TypeString,
							Description:  "Type of the HTTP rule.",
							Required:     true,
							ValidateFunc: validation.All(validation.StringInSlice([]string{"FORWARD", "STATIC", "REDIRECT"}, true)),
						},
						"target_group": {
							Type:        schema.TypeString,
							Description: "The ID of the target group; mandatory and only valid for FORWARD actions.",
							Optional:    true,
						},
						"drop_query": {
							Type:        schema.TypeBool,
							Description: "Default is false; valid only for REDIRECT actions.",
							Optional:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "The location for redirecting; mandatory and valid only for REDIRECT actions.",
							Optional:    true,
						},
						"status_code": {
							Type:         schema.TypeInt,
							Description:  "Valid only for REDIRECT and STATIC actions. For REDIRECT actions, default is 301 and possible values are 301, 302, 303, 307, and 308. For STATIC actions, default is 503 and valid range is 200 to 599.",
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.All(validation.IntInSlice([]int{301, 302, 303, 307, 308, 200, 503, 599})),
						},
						"response_message": {
							Type:        schema.TypeString,
							Description: "The response message of the request; mandatory for STATIC actions.",
							Optional:    true,
						},
						"content_type": {
							Type:        schema.TypeString,
							Description: "Valid only for STATIC actions.",
							Optional:    true,
							Computed:    true,
						},
						"conditions": {
							Type:        schema.TypeList,
							Description: "An array of items in the collection.The action is only performed if each and every condition is met; if no conditions are set, the rule will always be performed.",
							Optional:    true,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:             schema.TypeString,
										Description:      "Type of the HTTP rule condition.",
										Required:         true,
										ValidateFunc:     validation.All(validation.StringInSlice([]string{"HEADER", "PATH", "QUERY", "METHOD", "HOST", "COOKIE", "SOURCE_IP"}, true)),
										DiffSuppressFunc: utils.DiffToLower,
									},
									"condition": {
										Type:             schema.TypeString,
										Description:      "Matching rule for the HTTP rule condition attribute; mandatory for HEADER, PATH, QUERY, METHOD, HOST, and COOKIE types; must be null when type is SOURCE_IP.",
										Optional:         true,
										ValidateFunc:     validation.All(validation.StringInSlice([]string{"EXISTS", "CONTAINS", "EQUALS", "MATCHES", "STARTS_WITH", "ENDS_WITH"}, true)),
										DiffSuppressFunc: utils.DiffToLower,
									},
									"negate": {
										Type:        schema.TypeBool,
										Description: "Specifies whether the condition is negated or not; the default is False.",
										Optional:    true,
									},
									"key": {
										Type:        schema.TypeString,
										Description: "Must be null when type is PATH, METHOD, HOST, or SOURCE_IP. Key can only be set when type is COOKIES, HEADER, or QUERY.",
										Optional:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "Mandatory for conditions CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH; must be null when condition is EXISTS; should be a valid CIDR if provided and if type is SOURCE_IP.",
										Optional:    true,
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
				ValidateFunc: validation.All(validation.IsUUID),
			},
			"application_loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.IsUUID),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceApplicationLoadBalancerForwardingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	applicationLoadBalancerForwardingRule := ionoscloud.ApplicationLoadBalancerForwardingRule{
		Properties: &ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		applicationLoadBalancerForwardingRule.Properties.Name = &name
	}

	if protocol, protocolOk := d.GetOk("protocol"); protocolOk {
		protocol := protocol.(string)
		applicationLoadBalancerForwardingRule.Properties.Protocol = &protocol
	}

	if listenerIp, listenerIpOk := d.GetOk("listener_ip"); listenerIpOk {
		listenerIp := listenerIp.(string)
		applicationLoadBalancerForwardingRule.Properties.ListenerIp = &listenerIp
	}

	if listenerPort, listenerPortOk := d.GetOk("listener_port"); listenerPortOk {
		listenerPort := int32(listenerPort.(int))
		applicationLoadBalancerForwardingRule.Properties.ListenerPort = &listenerPort
	}

	if clientTimeout, clientTimeoutOk := d.GetOk("client_timeout"); clientTimeoutOk {
		clientTimeout := int32(clientTimeout.(int))
		applicationLoadBalancerForwardingRule.Properties.ClientTimeout = &clientTimeout
	}

	if serverCertificatesVal, serverCertificatesOk := d.GetOk("server_certificates"); serverCertificatesOk {
		serverCertificatesVal := serverCertificatesVal.(*schema.Set).List()
		if serverCertificatesVal != nil {
			serverCertificates := make([]string, 0)
			for _, value := range serverCertificatesVal {
				serverCertificates = append(serverCertificates, value.(string))
			}
			if len(serverCertificates) > 0 {
				applicationLoadBalancerForwardingRule.Properties.ServerCertificates = &serverCertificates
			}
		}
	}

	if _, httpRulesOk := d.GetOk("http_rules"); httpRulesOk {
		if httpRules, err := getAlbHttpRulesData(d); err == nil {
			applicationLoadBalancerForwardingRule.Properties.HttpRules = httpRules
		} else {
			return diag.FromErr(err)
		}
	}

	dcId := d.Get("datacenter_id").(string)
	albId := d.Get("application_loadbalancer_id").(string)

	albForwardingRuleResp, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPost(ctx, dcId, albId).ApplicationLoadBalancerForwardingRule(applicationLoadBalancerForwardingRule).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating application loadbalancer forwarding rule: %w \n ApiError: %s", err, responseBody(apiResponse)))
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

	client := meta.(SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)

	albId := d.Get("application_loadbalancer_id").(string)

	applicationLoadBalancerForwardingRule, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(ctx, dcId, albId, d.Id()).Execute()

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] Successfully retreived application load balancer forwarding rule %s: %+v", d.Id(), applicationLoadBalancerForwardingRule)

	if err := setApplicationLoadBalancerForwardingRuleData(d, &applicationLoadBalancerForwardingRule); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceApplicationLoadBalancerForwardingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

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

	if d.HasChange("client_timeout") {
		_, v := d.GetChange("client_timeout")
		vStr := int32(v.(int))
		request.Properties.ClientTimeout = &vStr
	}

	if d.HasChange("server_certificates") {
		_, v := d.GetChange("server_certificates")
		certificatesValues := v.(*schema.Set).List()
		serverCertificates := make([]string, 0)
		if certificatesValues != nil {
			for _, value := range certificatesValues {
				serverCertificates = append(serverCertificates, value.(string))
			}
		}
		request.Properties.ServerCertificates = &serverCertificates
	}

	if d.HasChange("http_rules") {
		if httpRules, err := getAlbHttpRulesData(d); err == nil {
			request.Properties.HttpRules = httpRules
		} else {
			return diag.FromErr(err)
		}
	}

	_, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPatch(ctx, dcId, albId, d.Id()).ApplicationLoadBalancerForwardingRuleProperties(*request.Properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a application loadbalancer forwarding rule ID %s %w",
			d.Id(), err))
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
	client := meta.(SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	albID := d.Get("application_loadbalancer_id").(string)

	apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesDelete(ctx, dcId, albID, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a application loadbalancer forwarding rule %s %w", d.Id(), err))
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

func resourceApplicationLoadBalancerForwardingRuleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).CloudApiClient

	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{alb}/{alb_forwarding_rule}", d.Id())
	}

	datacenterId := parts[0]
	albId := parts[1]
	ruleId := parts[2]

	albForwardingRule, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId, albId, ruleId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find alb forwarding rule %q", ruleId)
		}
		return nil, fmt.Errorf("an error occured while retrieving the alb forwarding rule %q, %w", ruleId, err)
	}

	if err := d.Set("datacenter_id", datacenterId); err != nil {
		return nil, fmt.Errorf("error while setting datacenter_id property for  alb forwarding rule %q: %w", ruleId, err)
	}
	if err := d.Set("application_loadbalancer_id", albId); err != nil {
		return nil, fmt.Errorf("error while setting application_loadbalancer_id property for  alb forwarding rule %q: %w", ruleId, err)
	}

	if err := setApplicationLoadBalancerForwardingRuleData(d, &albForwardingRule); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func setApplicationLoadBalancerForwardingRuleData(d *schema.ResourceData, applicationLoadBalancerForwardingRule *ionoscloud.ApplicationLoadBalancerForwardingRule) error {

	if applicationLoadBalancerForwardingRule.Id != nil {
		d.SetId(*applicationLoadBalancerForwardingRule.Id)
	}

	if applicationLoadBalancerForwardingRule.Properties != nil {
		if applicationLoadBalancerForwardingRule.Properties.Name != nil {
			err := d.Set("name", *applicationLoadBalancerForwardingRule.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for application load balancer forwarding rule %s: %w", d.Id(), err)
			}
		}

		if applicationLoadBalancerForwardingRule.Properties.Protocol != nil {
			err := d.Set("protocol", *applicationLoadBalancerForwardingRule.Properties.Protocol)
			if err != nil {
				return fmt.Errorf("error while setting protocol property for application load balancer forwarding rule %s: %w", d.Id(), err)
			}
		}

		if applicationLoadBalancerForwardingRule.Properties.ListenerIp != nil {
			err := d.Set("listener_ip", *applicationLoadBalancerForwardingRule.Properties.ListenerIp)
			if err != nil {
				return fmt.Errorf("error while setting listener_ip property for application load balancer forwarding rule %s: %w", d.Id(), err)
			}
		}

		if applicationLoadBalancerForwardingRule.Properties.ListenerPort != nil {
			err := d.Set("listener_port", *applicationLoadBalancerForwardingRule.Properties.ListenerPort)
			if err != nil {
				return fmt.Errorf("error while setting listener_port property for application load balancer forwarding rule %s: %w", d.Id(), err)
			}
		}

		if applicationLoadBalancerForwardingRule.Properties.ClientTimeout != nil {
			err := d.Set("client_timeout", *applicationLoadBalancerForwardingRule.Properties.ClientTimeout)
			if err != nil {
				return fmt.Errorf("error while setting client_timeout property for application load balancer forwarding rule %s: %w", d.Id(), err)
			}
		}

		if applicationLoadBalancerForwardingRule.Properties.ServerCertificates != nil {
			err := d.Set("server_certificates", *applicationLoadBalancerForwardingRule.Properties.ServerCertificates)
			if err != nil {
				return fmt.Errorf("error while setting server_certificates property for application load balancer forwarding rule %s: %w", d.Id(), err)
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
				return fmt.Errorf("error while setting http_rules property for application load balancer forwarding rule  %s: %w", d.Id(), err)
			}
		}
	}
	return nil
}

func getAlbHttpRulesData(d *schema.ResourceData) (*[]ionoscloud.ApplicationLoadBalancerHttpRule, error) {
	var httpRules []ionoscloud.ApplicationLoadBalancerHttpRule

	httpRulesVal := d.Get("http_rules").([]interface{})

	for httpRuleIndex := range httpRulesVal {

		httpRule := ionoscloud.ApplicationLoadBalancerHttpRule{}

		if name, nameOk := d.GetOk(fmt.Sprintf("http_rules.%d.name", httpRuleIndex)); nameOk {
			name := name.(string)
			httpRule.Name = &name
		}

		if typeVal, typeOk := d.GetOk(fmt.Sprintf("http_rules.%d.type", httpRuleIndex)); typeOk {
			typeVal := typeVal.(string)
			httpRule.Type = &typeVal
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

				var conditions []ionoscloud.ApplicationLoadBalancerHttpRuleCondition

				for conditionIndex := range conditionsVal.([]interface{}) {

					condition := ionoscloud.ApplicationLoadBalancerHttpRuleCondition{}

					typeVal := d.Get(fmt.Sprintf("http_rules.%d.conditions.%d.type", httpRuleIndex, conditionIndex)).(string)
					condition.Type = &typeVal

					if conditionVal, conditionOk := d.GetOk(fmt.Sprintf("http_rules.%d.conditions.%d.condition", httpRuleIndex, conditionIndex)); conditionOk {
						conditionVal := conditionVal.(string)
						condition.Condition = &conditionVal
					} else if strings.ToUpper(typeVal) != "SOURCE_IP" {
						return nil, fmt.Errorf("condition must be provided for application loadbalancer forwarding rule http rule condition")
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

					conditions = append(conditions, condition)
				}

				httpRule.Conditions = &conditions
			}

		}

		httpRules = append(httpRules, httpRule)

	}

	return &httpRules, nil
}
