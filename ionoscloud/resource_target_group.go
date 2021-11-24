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

func resourceTargetGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTargetGroupCreate,
		ReadContext:   resourceTargetGroupRead,
		UpdateContext: resourceTargetGroupUpdate,
		DeleteContext: resourceTargetGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTargetGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A name of that Target Group",
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"algorithm": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Algorithm for the balancing.",
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Protocol of the balancing.",
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"targets": {
				Type:        schema.TypeList,
				Description: "Array of items in that collection",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:         schema.TypeString,
							Description:  "IP of a balanced target VM",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
			"health_check": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Health check attributes for Application Load Balancer forwarding rule",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connect_timeout": {
							Type: schema.TypeInt,
							Description: "It specifies the maximum time (in milliseconds) to wait for a connection attempt " +
								"to a target VM to succeed. If unset, the default of 5 seconds will be used.",
							Optional: true,
							Computed: true,
						},
						"target_timeout": {
							Type: schema.TypeInt,
							Description: "argetTimeout specifies the maximum inactivity time (in milliseconds) on the " +
								"target VM side. If unset, the default of 50 seconds will be used.",
							Optional: true,
							Computed: true,
						},
						"retries": {
							Type: schema.TypeInt,
							Description: "Retries specifies the number of retries to perform on a target VM after a " +
								"connection failure. If unset, the default value of 3 will be used. (valid range: [0, 65535]).",
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"http_health_check": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Http health check attributes for Target Group",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Description: "The path for the HTTP health check; default: /.",
							Optional:    true,
							Computed:    true,
						},
						"method": {
							Type:        schema.TypeString,
							Description: "The method for the HTTP health check.",
							Optional:    true,
							Computed:    true,
						},
						"match_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"response": {
							Type:         schema.TypeString,
							Description:  "The response returned by the request.",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"regex": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"negate": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceTargetGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	targetGroup := ionoscloud.TargetGroup{
		Properties: &ionoscloud.TargetGroupProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		targetGroup.Properties.Name = &name
	} else {
		diags := diag.FromErr(fmt.Errorf("name must be provided for target group"))
		return diags
	}

	if algorithm, algorithmOk := d.GetOk("algorithm"); algorithmOk {
		algorithm := algorithm.(string)
		targetGroup.Properties.Algorithm = &algorithm
	} else {
		diags := diag.FromErr(fmt.Errorf("algorithm must be provided for target group"))
		return diags
	}

	if protocol, protocolOk := d.GetOk("protocol"); protocolOk {
		protocol := protocol.(string)
		targetGroup.Properties.Protocol = &protocol
	} else {
		diags := diag.FromErr(fmt.Errorf("protocol must be provided for target group"))
		return diags
	}

	if targetsVal, targetsOk := d.GetOk("targets"); targetsOk {
		if targetsVal.([]interface{}) != nil {
			updateTargets := false

			var targets []ionoscloud.TargetGroupTarget

			for targetIndex := range targetsVal.([]interface{}) {
				target := ionoscloud.TargetGroupTarget{}
				addTarget := false
				if ip, ipOk := d.GetOk(fmt.Sprintf("targets.%d.ip", targetIndex)); ipOk {
					ip := ip.(string)
					target.Ip = &ip
					addTarget = true
				} else {
					diags := diag.FromErr(fmt.Errorf("ip must be provided for target group target"))
					return diags
				}

				if port, portOk := d.GetOk(fmt.Sprintf("targets.%d.port", targetIndex)); portOk {
					port := int32(port.(int))
					target.Port = &port
				} else {
					diags := diag.FromErr(fmt.Errorf("port must be provided for target group target"))
					return diags
				}

				if weight, weightOk := d.GetOk(fmt.Sprintf("targets.%d.weight", targetIndex)); weightOk {
					weight := int32(weight.(int))
					target.Weight = &weight
				} else {
					diags := diag.FromErr(fmt.Errorf("weight must be provided for target group target"))
					return diags
				}

				if _, healthCheckOk := d.GetOk(fmt.Sprintf("targets.%d.health_check.0", targetIndex)); healthCheckOk {
					target.HealthCheck = &ionoscloud.TargetGroupTargetHealthCheck{}

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
				log.Printf("[INFO] Target group targets set to %+v", targets)
				targetGroup.Properties.Targets = &targets
			}
		}
	}

	if _, healthCheckOk := d.GetOk("health_check.0"); healthCheckOk {
		targetGroup.Properties.HealthCheck = &ionoscloud.TargetGroupHealthCheck{}

		if connectTimeout, connectTimeoutOk := d.GetOk("health_check.0.connect_timeout"); connectTimeoutOk {
			connectTimeout := int32(connectTimeout.(int))
			targetGroup.Properties.HealthCheck.ConnectTimeout = &connectTimeout
		}

		if targetTimeout, targetTimeoutOk := d.GetOk("health_check.0.target_timeout"); targetTimeoutOk {
			targetTimeout := int32(targetTimeout.(int))
			targetGroup.Properties.HealthCheck.TargetTimeout = &targetTimeout
		}

		if retries, retriesOk := d.GetOk("health_check.0.retries"); retriesOk {
			retries := int32(retries.(int))
			targetGroup.Properties.HealthCheck.Retries = &retries
		}

	}

	if _, httpHealthCheckOk := d.GetOk("http_health_check.0"); httpHealthCheckOk {
		targetGroup.Properties.HttpHealthCheck = &ionoscloud.TargetGroupHttpHealthCheck{}

		if path, pathOk := d.GetOk("http_health_check.0.path"); pathOk {
			path := path.(string)
			targetGroup.Properties.HttpHealthCheck.Path = &path
		}

		if method, methodOk := d.GetOk("http_health_check.0.method"); methodOk {
			method := method.(string)
			targetGroup.Properties.HttpHealthCheck.Method = &method
		}

		if matchType, matchTypeOk := d.GetOk("http_health_check.0.match_type"); matchTypeOk {
			matchType := matchType.(string)
			targetGroup.Properties.HttpHealthCheck.MatchType = &matchType
		}

		if response, responseOk := d.GetOk("http_health_check.0.response"); responseOk {
			response := response.(string)
			targetGroup.Properties.HttpHealthCheck.Response = &response
		}

		if regex, regexOk := d.GetOk("http_health_check.0.regex"); regexOk {
			regex := regex.(bool)
			targetGroup.Properties.HttpHealthCheck.Regex = &regex
		}

		if negate, negateOk := d.GetOk("http_health_check.0.negate"); negateOk {
			negate := negate.(bool)
			targetGroup.Properties.HttpHealthCheck.Negate = &negate
		}

	}

	rsp, apiResponse, err := client.TargetGroupsApi.TargetgroupsPost(ctx).TargetGroup(targetGroup).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a target group: %s ", err))
		return diags
	}

	d.SetId(*rsp.Id)
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

	return resourceTargetGroupRead(ctx, d, meta)
}

func resourceTargetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	rsp, apiResponse, err := client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a target group %s %s", d.Id(), err))
		return diags
	}

	if err := setTargetGroupData(d, &rsp); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTargetGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	input := ionoscloud.TargetGroupProperties{}

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		vStr := v.(string)
		input.Name = &vStr
	}

	if d.HasChange("algorithm") {
		_, v := d.GetChange("algorithm")
		vStr := v.(string)
		input.Algorithm = &vStr
	}

	if d.HasChange("protocol") {
		_, v := d.GetChange("protocol")
		vStr := v.(string)
		input.Protocol = &vStr
	}

	if d.HasChange("targets") {
		oldTargets, newTargets := d.GetChange("targets")
		if newTargets.([]interface{}) != nil {
			updateTargets := false

			var targets []ionoscloud.TargetGroupTarget

			for targetIndex := range newTargets.([]interface{}) {
				target := ionoscloud.TargetGroupTarget{}

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
					target.HealthCheck = &ionoscloud.TargetGroupTargetHealthCheck{}

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
				input.Targets = &targets
			}
		}
	}

	if d.HasChange("health_check.0") {
		_, v := d.GetChange("health_check.0")
		if v.(map[string]interface{}) != nil {
			updateHealthCheck := false

			healthCheck := &ionoscloud.TargetGroupHealthCheck{}

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
				input.HealthCheck = healthCheck
			}
		}
	}

	if d.HasChange("http_health_check.0") {
		_, v := d.GetChange("http_health_check.0")
		if v.(map[string]interface{}) != nil {
			updateHttpHealthCheck := false

			healthCheck := &ionoscloud.TargetGroupHttpHealthCheck{}

			if d.HasChange("http_health_check.0.path") {
				_, newValue := d.GetChange("health_check.0.path")
				if newValue != 0 {
					newValue := newValue.(string)
					healthCheck.Path = &newValue
				}
			}

			if d.HasChange("http_health_check.0.method") {
				_, newValue := d.GetChange("health_check.0.method")
				if newValue != 0 {
					newValue := newValue.(string)
					healthCheck.Method = &newValue
				}
			}

			if d.HasChange("http_health_check.0.match_type") {
				_, newValue := d.GetChange("health_check.0.match_type")
				if newValue != 0 {
					updateHttpHealthCheck = true
					newValue := newValue.(string)
					healthCheck.MatchType = &newValue
				}
			}

			if d.HasChange("http_health_check.0.response") {
				_, newValue := d.GetChange("health_check.0.response")
				if newValue != 0 {
					newValue := newValue.(string)
					healthCheck.Response = &newValue
				}
			}

			if d.HasChange("http_health_check.0.regex") {
				_, newValue := d.GetChange("health_check.0.regex")
				if newValue != 0 {
					newValue := newValue.(bool)
					healthCheck.Regex = &newValue
				}
			}

			if d.HasChange("http_health_check.0.negate") {
				_, newValue := d.GetChange("health_check.0.negate")
				if newValue != 0 {
					newValue := newValue.(bool)
					healthCheck.Negate = &newValue
				}
			}
			if updateHttpHealthCheck == true {
				input.HttpHealthCheck = healthCheck
			}
		}
	}

	_, apiResponse, err := client.TargetGroupsApi.TargetgroupsPatch(ctx, d.Id()).TargetGroupProperties(input).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while restoring a targetGroup ID %s %d", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceTargetGroupRead(ctx, d, meta)
}

func resourceTargetGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	apiResponse, err := client.TargetGroupsApi.TargetGroupsDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a target group %s %s", d.Id(), err))
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

func resourceTargetGroupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	groupIp := d.Id()

	groupTarget, apiResponse, err := client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, groupIp).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find target group %q", groupIp)
		}
		return nil, fmt.Errorf("an error occured while retrieving the target group %q, %q", groupIp, err)
	}

	if err := setTargetGroupData(d, &groupTarget); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func setTargetGroupData(d *schema.ResourceData, targetGroup *ionoscloud.TargetGroup) error {

	if targetGroup.Id != nil {
		d.SetId(*targetGroup.Id)
	}

	if targetGroup.Properties != nil {

		if targetGroup.Properties.Name != nil {
			err := d.Set("name", *targetGroup.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for target group %s: %s", d.Id(), err)
			}
		}

		if targetGroup.Properties.Algorithm != nil {
			err := d.Set("algorithm", *targetGroup.Properties.Algorithm)
			if err != nil {
				return fmt.Errorf("error while setting algorithm property for target group %s: %s", d.Id(), err)
			}
		}

		if targetGroup.Properties.Protocol != nil {
			err := d.Set("protocol", *targetGroup.Properties.Protocol)
			if err != nil {
				return fmt.Errorf("error while setting protocol property for target group %s: %s", d.Id(), err)
			}
		}

		forwardingRuleTargets := make([]interface{}, 0)
		if targetGroup.Properties.Targets != nil && len(*targetGroup.Properties.Targets) > 0 {
			forwardingRuleTargets = make([]interface{}, 0)
			for _, target := range *targetGroup.Properties.Targets {
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

				forwardingRuleTargets = append(forwardingRuleTargets, targetEntry)
			}
		}

		if len(forwardingRuleTargets) > 0 {
			if err := d.Set("targets", forwardingRuleTargets); err != nil {
				return fmt.Errorf("error while setting targets property for target group  %s: %s", d.Id(), err)
			}
		}

		if targetGroup.Properties.HealthCheck != nil {
			healthCheck := make([]interface{}, 1)

			healthCheckEntry := make(map[string]interface{})

			if targetGroup.Properties.HealthCheck.ConnectTimeout != nil {
				healthCheckEntry["connect_timeout"] = *targetGroup.Properties.HealthCheck.ConnectTimeout
			}

			if targetGroup.Properties.HealthCheck.TargetTimeout != nil {
				healthCheckEntry["target_timeout"] = *targetGroup.Properties.HealthCheck.TargetTimeout
			}

			if targetGroup.Properties.HealthCheck.Retries != nil {
				healthCheckEntry["retries"] = *targetGroup.Properties.HealthCheck.Retries
			}

			healthCheck[0] = healthCheckEntry
			err := d.Set("health_check", healthCheck)
			if err != nil {
				return fmt.Errorf("error while setting health_check property for target group %s: %s", d.Id(), err)
			}
		}

		if targetGroup.Properties.HttpHealthCheck != nil {
			httpHealthCheck := make([]interface{}, 1)

			httpHealthCheckEntry := make(map[string]interface{})

			if targetGroup.Properties.HttpHealthCheck.Path != nil {
				httpHealthCheckEntry["path"] = *targetGroup.Properties.HttpHealthCheck.Path
			}

			if targetGroup.Properties.HttpHealthCheck.Method != nil {
				httpHealthCheckEntry["method"] = *targetGroup.Properties.HttpHealthCheck.Method
			}

			if targetGroup.Properties.HttpHealthCheck.MatchType != nil {
				httpHealthCheckEntry["match_type"] = *targetGroup.Properties.HttpHealthCheck.MatchType
			}

			if targetGroup.Properties.HttpHealthCheck.Response != nil {
				httpHealthCheckEntry["response"] = *targetGroup.Properties.HttpHealthCheck.Response
			}

			if targetGroup.Properties.HttpHealthCheck.Regex != nil {
				httpHealthCheckEntry["regex"] = *targetGroup.Properties.HttpHealthCheck.Regex
			}

			if targetGroup.Properties.HttpHealthCheck.Negate != nil {
				httpHealthCheckEntry["negate"] = *targetGroup.Properties.HttpHealthCheck.Negate
			}

			httpHealthCheck[0] = httpHealthCheckEntry
			err := d.Set("http_health_check", httpHealthCheck)
			if err != nil {
				return fmt.Errorf("error while setting http_health_check property for target group %s: %s", d.Id(), err)
			}
		}

	}
	return nil
}
