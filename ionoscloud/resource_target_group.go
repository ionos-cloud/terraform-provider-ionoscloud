package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
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
				Description:  "The name of the target group.",
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"algorithm": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Balancing algorithm.",
				ValidateFunc: validation.All(validation.StringInSlice([]string{"ROUND_ROBIN", "LEAST_CONNECTION", "RANDOM", "SOURCE_IP"}, true)),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Balancing protocol.",
				ValidateFunc: validation.All(validation.StringInSlice([]string{"HTTP"}, true)),
			},
			"targets": {
				Type:        schema.TypeList,
				Description: "Array of items in the collection.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:         schema.TypeString,
							Description:  "The IP of the balanced target VM.",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "The port of the balanced target service; valid range is 1 to 65535.",
							Required:    true,
						},
						"weight": {
							Type:        schema.TypeInt,
							Description: "Traffic is distributed in proportion to target weight, relative to the combined weight of all targets. A target with higher weight receives a greater share of traffic. Valid range is 0 to 256 and default is 1; targets with weight of 0 do not participate in load balancing but still accept persistent connections. It is best use values in the middle of the range to leave room for later adjustments.",
							Required:    true,
						},
						"health_check_enabled": {
							Type:        schema.TypeBool,
							Description: "Makes the target available only if it accepts periodic health check TCP connection attempts; when turned off, the target is considered always available. The health check only consists of a connection attempt to the address and port of the target. Default is True.",
							Optional:    true,
							Computed:    true,
						},
						"maintenance_enabled": {
							Type:        schema.TypeBool,
							Description: "Maintenance mode prevents the target from receiving balanced traffic.",
							Optional:    true,
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
						"check_timeout": {
							Type:        schema.TypeInt,
							Description: "The maximum time in milliseconds to wait for a target to respond to a check. For target VMs with 'Check Interval' set, the lesser of the two  values is used once the TCP connection is established.",
							Optional:    true,
							Computed:    true,
						},
						"check_interval": {
							Type:        schema.TypeInt,
							Description: "The interval in milliseconds between consecutive health checks; default is 2000.",
							Optional:    true,
							Computed:    true,
						},
						"retries": {
							Type:         schema.TypeInt,
							Description:  "The maximum number of attempts to reconnect to a target after a connection failure. Valid range is 0 to 65535, and default is three reconnection.",
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.All(validation.IntBetween(1, 65535)),
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
							Description: "The path (destination URL) for the HTTP health check request; the default is /.",
							Optional:    true,
							Computed:    true,
						},
						"method": {
							Type:         schema.TypeString,
							Description:  "The method for the HTTP health check.",
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.All(validation.StringInSlice([]string{"HEAD", "PUT", "POST", "GET", "TRACE", "PATCH", "OPTIONS"}, true)),
						},
						"match_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.All(validation.StringInSlice([]string{"STATUS_CODE", "RESPONSE_BODY"}, true)),
						},
						"response": {
							Type:         schema.TypeString,
							Description:  "The response returned by the request, depending on the match type.",
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
	client := meta.(services.SdkBundle).CloudApiClient

	targetGroup := ionoscloud.TargetGroup{
		Properties: &ionoscloud.TargetGroupProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		targetGroup.Properties.Name = &name
	}

	if algorithm, algorithmOk := d.GetOk("algorithm"); algorithmOk {
		algorithm := algorithm.(string)
		targetGroup.Properties.Algorithm = &algorithm
	}

	if protocol, protocolOk := d.GetOk("protocol"); protocolOk {
		protocol := protocol.(string)
		targetGroup.Properties.Protocol = &protocol
	}

	targetGroup.Properties.Targets = getTargetGroupTargetData(d)

	if _, healthCheckOk := d.GetOk("health_check.0"); healthCheckOk {
		targetGroup.Properties.HealthCheck = getTargetGroupHealthCheckData(d)
	}

	if _, httpHealthCheckOk := d.GetOk("http_health_check.0"); httpHealthCheckOk {
		targetGroup.Properties.HttpHealthCheck = getTargetGroupHttpHealthCheckData(d)
	}

	rsp, apiResponse, err := client.TargetGroupsApi.TargetgroupsPost(ctx).TargetGroup(targetGroup).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a target group: %w ", err))
		return diags
	}

	d.SetId(*rsp.Id)
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if cloudapi.IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceTargetGroupRead(ctx, d, meta)
}

func resourceTargetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	rsp, apiResponse, err := client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a target group %s %w", d.Id(), err))
		return diags
	}

	if err := setTargetGroupData(d, &rsp); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTargetGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	targetGroup := ionoscloud.TargetGroupPut{
		Properties: &ionoscloud.TargetGroupProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		targetGroup.Properties.Name = &name
	}

	if algorithm, algorithmOk := d.GetOk("algorithm"); algorithmOk {
		algorithm := algorithm.(string)
		targetGroup.Properties.Algorithm = &algorithm
	}

	if protocol, protocolOk := d.GetOk("protocol"); protocolOk {
		protocol := protocol.(string)
		targetGroup.Properties.Protocol = &protocol
	}

	targetGroup.Properties.Targets = getTargetGroupTargetData(d)

	if _, healthCheckOk := d.GetOk("health_check.0"); healthCheckOk {
		targetGroup.Properties.HealthCheck = getTargetGroupHealthCheckData(d)
	}

	if _, httpHealthCheckOk := d.GetOk("http_health_check.0"); httpHealthCheckOk {
		targetGroup.Properties.HttpHealthCheck = getTargetGroupHttpHealthCheckData(d)
	}

	response, apiResponse, err := client.TargetGroupsApi.TargetgroupsPut(ctx, d.Id()).TargetGroup(targetGroup).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while restoring a targetGroup ID %s %w", d.Id(), err))
		return diags
	}

	d.SetId(*response.Id)

	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceTargetGroupRead(ctx, d, meta)
}

func resourceTargetGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	apiResponse, err := client.TargetGroupsApi.TargetGroupsDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a target group %s %w", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	d.SetId("")

	return nil
}

func resourceTargetGroupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).CloudApiClient

	groupIp := d.Id()

	groupTarget, apiResponse, err := client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, groupIp).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find target group %q", groupIp)
		}
		return nil, fmt.Errorf("an error occured while retrieving the target group %q, %w", groupIp, err)
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
				return fmt.Errorf("error while setting name property for target group %s: %w", d.Id(), err)
			}
		}

		if targetGroup.Properties.Algorithm != nil {
			err := d.Set("algorithm", *targetGroup.Properties.Algorithm)
			if err != nil {
				return fmt.Errorf("error while setting algorithm property for target group %s: %w", d.Id(), err)
			}
		}

		if targetGroup.Properties.Protocol != nil {
			err := d.Set("protocol", *targetGroup.Properties.Protocol)
			if err != nil {
				return fmt.Errorf("error while setting protocol property for target group %s: %w", d.Id(), err)
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

				if target.HealthCheckEnabled != nil {
					targetEntry["health_check_enabled"] = *target.HealthCheckEnabled
				}

				if target.MaintenanceEnabled != nil {
					targetEntry["maintenance_enabled"] = *target.MaintenanceEnabled
				}

				forwardingRuleTargets = append(forwardingRuleTargets, targetEntry)
			}
		}

		if len(forwardingRuleTargets) > 0 {
			if err := d.Set("targets", forwardingRuleTargets); err != nil {
				return fmt.Errorf("error while setting targets property for target group  %s: %w", d.Id(), err)
			}
		}

		if targetGroup.Properties.HealthCheck != nil {
			healthCheck := make([]interface{}, 1)

			healthCheckEntry := make(map[string]interface{})

			if targetGroup.Properties.HealthCheck.CheckTimeout != nil {
				healthCheckEntry["check_timeout"] = *targetGroup.Properties.HealthCheck.CheckTimeout
			}

			if targetGroup.Properties.HealthCheck.CheckInterval != nil {
				healthCheckEntry["check_interval"] = *targetGroup.Properties.HealthCheck.CheckInterval
			}

			if targetGroup.Properties.HealthCheck.Retries != nil {
				healthCheckEntry["retries"] = *targetGroup.Properties.HealthCheck.Retries
			}

			healthCheck[0] = healthCheckEntry
			err := d.Set("health_check", healthCheck)
			if err != nil {
				return fmt.Errorf("error while setting health_check property for target group %s: %w", d.Id(), err)
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
				return fmt.Errorf("error while setting http_health_check property for target group %s: %w", d.Id(), err)
			}
		}

	}
	return nil
}

func getTargetGroupTargetData(d *schema.ResourceData) *[]ionoscloud.TargetGroupTarget {
	targets := make([]ionoscloud.TargetGroupTarget, 0)

	if targetsVal, targetsOk := d.GetOk("targets"); targetsOk {

		if targetsVal.([]interface{}) != nil {

			for targetIndex := range targetsVal.([]interface{}) {
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

				healthCheck := d.Get(fmt.Sprintf("targets.%d.health_check_enabled", targetIndex)).(bool)
				target.HealthCheckEnabled = &healthCheck

				maintenance := d.Get(fmt.Sprintf("targets.%d.maintenance_enabled", targetIndex)).(bool)
				target.MaintenanceEnabled = &maintenance

				targets = append(targets, target)
			}
		}
	}
	return &targets
}

func getTargetGroupHealthCheckData(d *schema.ResourceData) *ionoscloud.TargetGroupHealthCheck {
	healthCheck := ionoscloud.TargetGroupHealthCheck{}

	if checkTimeout, checkTimeoutOk := d.GetOk("health_check.0.check_timeout"); checkTimeoutOk {
		checkTimeout := int32(checkTimeout.(int))
		healthCheck.CheckTimeout = &checkTimeout
	}

	if checkInterval, checkIntervalOk := d.GetOk("health_check.0.check_interval"); checkIntervalOk {
		checkInterval := int32(checkInterval.(int))
		healthCheck.CheckInterval = &checkInterval
	}

	if retries, retriesOk := d.GetOk("health_check.0.retries"); retriesOk {
		retries := int32(retries.(int))
		healthCheck.Retries = &retries
	}

	return &healthCheck
}

func getTargetGroupHttpHealthCheckData(d *schema.ResourceData) *ionoscloud.TargetGroupHttpHealthCheck {
	httpHealthCheck := ionoscloud.TargetGroupHttpHealthCheck{}

	if path, pathOk := d.GetOk("http_health_check.0.path"); pathOk {
		path := path.(string)
		httpHealthCheck.Path = &path
	}

	if method, methodOk := d.GetOk("http_health_check.0.method"); methodOk {
		method := method.(string)
		httpHealthCheck.Method = &method
	}

	if matchType, matchTypeOk := d.GetOk("http_health_check.0.match_type"); matchTypeOk {
		matchType := matchType.(string)
		httpHealthCheck.MatchType = &matchType
	}

	if response, responseOk := d.GetOk("http_health_check.0.response"); responseOk {
		response := response.(string)
		httpHealthCheck.Response = &response
	}

	if regex, regexOk := d.GetOk("http_health_check.0.regex"); regexOk {
		regex := regex.(bool)
		httpHealthCheck.Regex = &regex
	}

	if negate, negateOk := d.GetOk("http_health_check.0.negate"); negateOk {
		negate := negate.(bool)
		httpHealthCheck.Negate = &negate
	}

	return &httpHealthCheck
}
