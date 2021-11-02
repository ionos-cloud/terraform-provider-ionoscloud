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
	var apiResponse *ionoscloud.APIResponse

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if idOk {
		/* search by ID */
		targetGroup, apiResponse, err = client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)

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

		targetGroups, apiResponse, err := client.TargetGroupsApi.TargetgroupsGet(ctx).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("an error occurred while fetching target groups: %s", err.Error())
		}

		if targetGroups.Items != nil {
			for _, c := range *targetGroups.Items {
				tmpTargetGroup, apiResponse, err := client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, *c.Id).Execute()
				logApiRequestTime(apiResponse)
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

	if err = setTargetGroupData(d, &targetGroup); err != nil {
		return err
	}

	return nil
}
