package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceTargetGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTargetGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the target group.",
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Balancing algorithm.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Balancing protocol.",
			},
			"targets": {
				Type:        schema.TypeList,
				Description: "Array of items in the collection",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Description: "The IP of the balanced target VM.",
							Computed:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "The port of the balanced target service; valid range is 1 to 65535.",
							Computed:    true,
						},
						"weight": {
							Type:        schema.TypeInt,
							Description: "Traffic is distributed in proportion to target weight, relative to the combined weight of all targets. A target with higher weight receives a greater share of traffic. Valid range is 0 to 256 and default is 1; targets with weight of 0 do not participate in load balancing but still accept persistent connections. It is best use values in the middle of the range to leave room for later adjustments.",
							Computed:    true,
						},
						"health_check_enabled": {
							Type:        schema.TypeBool,
							Description: "Makes the target available only if it accepts periodic health check TCP connection attempts; when turned off, the target is considered always available. The health check only consists of a connection attempt to the address and port of the target. Default is True.",
							Computed:    true,
						},
						"maintenance_enabled": {
							Type:        schema.TypeBool,
							Description: "Maintenance mode prevents the target from receiving balanced traffic.",
							Computed:    true,
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
						"check_timeout": {
							Type:        schema.TypeInt,
							Description: "The maximum time in milliseconds to wait for a target to respond to a check. For target VMs with 'Check Interval' set, the lesser of the two  values is used once the TCP connection is established.",
							Computed:    true,
						},
						"check_interval": {
							Type:        schema.TypeInt,
							Description: "The interval in milliseconds between consecutive health checks; default is 2000.",
							Computed:    true,
						},
						"retries": {
							Type:        schema.TypeInt,
							Description: "The maximum number of attempts to reconnect to a target after a connection failure. Valid range is 0 to 65535, and default is three reconnection.",
							Computed:    true,
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
							Description: "The path (destination URL) for the HTTP health check request; the default is /.",
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
							Description: "The response returned by the request, depending on the match type.",
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

func dataSourceTargetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the target group id or name"))
	}
	var targetGroup ionoscloud.TargetGroup
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for target group by id %s", id)
		targetGroup, apiResponse, err = client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, id).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the target group %s: %w", id, err))
		}
	} else {
		/* search by name */
		var results []ionoscloud.TargetGroup

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for target group by name with partial_match %t and name: %s", partialMatch, name)

		if partialMatch {
			targetGroups, apiResponse, err := client.TargetGroupsApi.TargetgroupsGet(ctx).Depth(1).Filter("name", name).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching target groups: %w", err))
			}

			results = *targetGroups.Items
		} else {
			targetGroups, apiResponse, err := client.TargetGroupsApi.TargetgroupsGet(ctx).Depth(1).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching target groups: %w", err))
			}

			if targetGroups.Items != nil {
				for _, t := range *targetGroups.Items {
					if t.Properties.Name != nil && strings.ToLower(*t.Properties.Name) == strings.ToLower(name) {
						tmpTargetGroup, apiResponse, err := client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(ctx, *t.Id).Execute()
						logApiRequestTime(apiResponse)
						if err != nil {
							return diag.FromErr(fmt.Errorf("an error occurred while fetching target group with ID %s: %w", *t.Id, err))
						}
						results = append(results, tmpTargetGroup)

					}
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no target group found with the specified criteria: name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one target group found with the specified criteria: name = %s", name))
		}

		targetGroup = results[0]

	}

	if err = setTargetGroupData(d, &targetGroup); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
