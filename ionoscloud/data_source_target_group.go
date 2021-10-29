package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"strings"
)

func dataSourceTargetGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTargetGroupRead,
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
				Computed:    true,
				Description: "Algorithm for the balancing.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Protocol of the balancing.",
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
			"health_check": {
				Type:        schema.TypeList,
				Description: "Health check attributes for Application Load Balancer forwarding rule",
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
							Description: "It specifies the maximum time (in milliseconds) to wait for a connection attempt " +
								"to a target VM to succeed. If unset, the default of 5 seconds will be used.",
							Computed: true,
						},
						"target_timeout": {
							Type: schema.TypeInt,
							Description: "argetTimeout specifies the maximum inactivity time (in milliseconds) on the " +
								"target VM side. If unset, the default of 50 seconds will be used.",
							Computed: true,
						},
						"retries": {
							Type: schema.TypeInt,
							Description: "Retries specifies the number of retries to perform on a target VM after a " +
								"connection failure. If unset, the default value of 3 will be used. (valid range: [0, 65535]).",
							Computed: true,
						},
					},
				},
			},
			"http_health_check": {
				Type:        schema.TypeList,
				Description: "Http health check attributes for Target Group",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Description: "The path for the HTTP health check; default: /.",
							Computed:    true,
						},
						"method": {
							Type:        schema.TypeString,
							Description: "The method for the HTTP health check.",
							Computed:    true,
						},
						"match_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"response": {
							Type:        schema.TypeString,
							Description: "The response returned by the request.",
							Computed:    true,
						},
						"regex": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"negate": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceTargetGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the target group id or name")
	}
	var targetGroup ionoscloud.TargetGroup
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if idOk {
		/* search by ID */
		targetGroup, _, err = client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, id.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the target groups %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var targetGroups ionoscloud.TargetGroups

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		targetGroups, _, err := client.TargetGroupsApi.TargetgroupsGet(ctx).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching target groups: %s", err.Error())
		}

		if targetGroups.Items != nil {
			for _, c := range *targetGroups.Items {
				tmpTargetGroup, _, err := client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, *c.Id).Execute()
				if err != nil {
					return fmt.Errorf("an error occurred while fetching target group with ID %s: %s", *c.Id, err.Error())
				}
				if tmpTargetGroup.Properties.Name != nil {
					if strings.Contains(*tmpTargetGroup.Properties.Name, name.(string)) {
						targetGroup = tmpTargetGroup
						break
					}
				}

			}
		}

	}

	if &targetGroup == nil {
		return errors.New("target group not found")
	}

	if err = setTargetGroupData(d, &targetGroup, client); err != nil {
		return err
	}

	return nil
}

func setTargetGroupData(d *schema.ResourceData, targetGroup *ionoscloud.TargetGroup, _ *ionoscloud.APIClient) error {

	if targetGroup.Id != nil {
		d.SetId(*targetGroup.Id)
		if err := d.Set("id", *targetGroup.Id); err != nil {
			return err
		}
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
