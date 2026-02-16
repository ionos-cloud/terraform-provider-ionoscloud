package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceApplicationLoadBalancerForwardingRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationLoadBalancerForwardingRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Application Load Balancer forwarding rule.",
				Optional:    true,
				Computed:    true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"protocol": {
				Type:        schema.TypeString,
				Description: "Balancing protocol",
				Computed:    true,
			},
			"listener_ip": {
				Type:        schema.TypeString,
				Description: "Listening (inbound) IP",
				Computed:    true,
			},
			"listener_port": {
				Type:        schema.TypeInt,
				Description: "Listening (inbound) port number; valid range is 1 to 65535.",
				Computed:    true,
			},
			"client_timeout": {
				Type:        schema.TypeInt,
				Description: "The maximum time in milliseconds to wait for the client to acknowledge or send data; default is 50,000 (50 seconds).",
				Computed:    true,
			},
			"server_certificates": {
				Type:        schema.TypeSet,
				Description: "Array of items in the collection.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"http_rules": {
				Type:        schema.TypeList,
				Description: "Array of items in that collection",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The unique name of the Application Load Balancer HTTP rule.",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Type of the HTTP rule.",
							Computed:    true,
						},
						"target_group": {
							Type:        schema.TypeString,
							Description: "The ID of the target group; mandatory and only valid for FORWARD actions.",
							Computed:    true,
						},
						"drop_query": {
							Type:        schema.TypeBool,
							Description: "Default is false; valid only for REDIRECT actions.",
							Computed:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "The location for redirecting; mandatory and valid only for REDIRECT actions.",
							Computed:    true,
						},
						"status_code": {
							Type:        schema.TypeInt,
							Description: "Valid only for REDIRECT and STATIC actions. For REDIRECT actions, default is 301 and possible values are 301, 302, 303, 307, and 308. For STATIC actions, default is 503 and valid range is 200 to 599.",
							Computed:    true,
						},
						"response_message": {
							Type:        schema.TypeString,
							Description: "The response message of the request; mandatory for STATIC actions.",
							Computed:    true,
						},
						"content_type": {
							Type:        schema.TypeString,
							Description: "Valid only for STATIC actions.",
							Computed:    true,
						},
						"conditions": {
							Type:        schema.TypeList,
							Description: "An array of items in the collection.The action is only performed if each and every condition is met; if no conditions are set, the rule will always be performed.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Description: "Type of the HTTP rule condition.",
										Computed:    true,
									},
									"condition": {
										Type:        schema.TypeString,
										Description: "Matching rule for the HTTP rule condition attribute; mandatory for HEADER, PATH, QUERY, METHOD, HOST, and COOKIE types; must be null when type is SOURCE_IP.",
										Computed:    true,
									},
									"negate": {
										Type:        schema.TypeBool,
										Description: "Specifies whether the condition is negated or not; the default is False.",
										Computed:    true,
									},
									"key": {
										Type:        schema.TypeString,
										Description: "Must be null when type is PATH, METHOD, HOST, or SOURCE_IP. Key can only be set when type is COOKIES, HEADER, or QUERY.",
										Computed:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "Mandatory for conditions CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH; must be null when condition is EXISTS; should be a valid CIDR if provided and if type is SOURCE_IP.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"application_loadbalancer_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceApplicationLoadBalancerForwardingRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	datacenterId := d.Get("datacenter_id").(string)
	albId := d.Get("application_loadbalancer_id").(string)

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return utils.ToDiags(d, "id and name cannot be both specified in the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the application load balancer forwarding rule id or name", nil)
	}

	var applicationLoadBalancerForwardingRule ionoscloud.ApplicationLoadBalancerForwardingRule
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for application load balancer forwarding rule by id %s", id)
		applicationLoadBalancerForwardingRule, apiResponse, err = client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId, albId, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the application load balancer forwarding rule while searching by ID %s: %s", id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		/* search by name */
		var results []ionoscloud.ApplicationLoadBalancerForwardingRule

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for application load balancer forwarding rule by name with partial_match %t and name: %s", partialMatch, name)

		if partialMatch {
			applicationLoadBalancersForwardingRules, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesGet(ctx, datacenterId, albId).Depth(1).Filter("name", name).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching application loadbalancer forwarding rules: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
			}

			results = *applicationLoadBalancersForwardingRules.Items
		} else {
			applicationLoadBalancersForwardingRules, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesGet(ctx, datacenterId, albId).Depth(1).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching application loadbalancer forwarding rules: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
			}

			if applicationLoadBalancersForwardingRules.Items != nil {
				for _, albFr := range *applicationLoadBalancersForwardingRules.Items {
					if albFr.Properties != nil && albFr.Properties.Name != nil && strings.EqualFold(*albFr.Properties.Name, name) {
						tmpAlbFr, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId, albId, *albFr.Id).Execute()
						logApiRequestTime(apiResponse)
						if err != nil {
							return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching application load balancer forwarding rule with ID %s: %s", *albFr.Id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
						}
						results = append(results, tmpAlbFr)
					}
				}
			}
		}

		if results == nil || len(results) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no application load balanacer forwarding rule found with the specified criteria: name = %s", name), nil)
		} else if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one application load balanacer forwarding rule found with the specified criteria: name = %s", name), nil)
		}

		applicationLoadBalancerForwardingRule = results[0]
	}

	if err = setApplicationLoadBalancerForwardingRuleData(d, &applicationLoadBalancerForwardingRule); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
