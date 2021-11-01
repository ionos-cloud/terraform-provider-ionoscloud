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

func dataSourceApplicationLoadBalancerForwardingRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationLoadBalancerForwardingRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:        schema.TypeString,
				Description: "Protocol of the balancing.",
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
					},
				},
			},
			"server_certificates": {
				Type:        schema.TypeList,
				Description: "Array of items in that collection.",
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
							Description: "A name of that Application Load Balancer http rule",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Type of the Http Rule",
							Computed:    true,
						},
						"target_group": {
							Type:        schema.TypeString,
							Description: "The UUID of the target group; mandatory for FORWARD action",
							Computed:    true,
						},
						"drop_query": {
							Type:        schema.TypeBool,
							Description: "Default is false; true for REDIRECT action.",
							Computed:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "The location for redirecting; mandatory for REDIRECT action",
							Computed:    true,
						},
						"status_code": {
							Type:        schema.TypeInt,
							Description: "On REDIRECT action it can take the value 301, 302, 303, 307, 308; on STATIC action it is between 200 and 599",
							Computed:    true,
						},
						"response_message": {
							Type:        schema.TypeString,
							Description: "he response message of the request; mandatory for STATIC action",
							Computed:    true,
						},
						"content_type": {
							Type:        schema.TypeString,
							Description: "Will be provided by the PAAS Team; default application/json",
							Computed:    true,
						},
						"conditions": {
							Type:        schema.TypeList,
							Description: "Array of items in that collection",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Description: "Type of the Http Rule condition.",
										Computed:    true,
									},
									"condition": {
										Type:        schema.TypeString,
										Description: "Condition of the Http Rule condition.",
										Computed:    true,
									},
									"negate": {
										Type:        schema.TypeBool,
										Description: "Specifies whether the condition is negated or not; default: false.",
										Computed:    true,
									},
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
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
			"application_loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceApplicationLoadBalancerForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return errors.New("no datacenter_id was specified")
	}

	albId, albIdOk := d.GetOk("application_loadbalancer_id")
	if !albIdOk {
		return errors.New("no application_loadbalancer_id was specified")
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the application loadbalancer forwarding rule id or name")
	}
	var applicationLoadBalancerForwardingRule ionoscloud.ApplicationLoadBalancerForwardingRule
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if idOk {
		/* search by ID */
		applicationLoadBalancerForwardingRule, _, err = client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId.(string), albId.(string), id.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the nat gateway %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var applicationLoadBalancersForwardingRules ionoscloud.ApplicationLoadBalancerForwardingRules

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		applicationLoadBalancersForwardingRules, _, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesGet(ctx, datacenterId.(string), albId.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching application loadbalancers: %s", err.Error())
		}

		if applicationLoadBalancersForwardingRules.Items != nil {
			for _, c := range *applicationLoadBalancersForwardingRules.Items {
				tmpAlb, _, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(ctx, datacenterId.(string), albId.(string), *c.Id).Execute()
				if err != nil {
					return fmt.Errorf("an error occurred while fetching nat gateway with ID %s: %s", *c.Id, err.Error())
				}
				if tmpAlb.Properties.Name != nil {
					if strings.Contains(*tmpAlb.Properties.Name, name.(string)) {
						applicationLoadBalancerForwardingRule = tmpAlb
						break
					}
				}

			}
		}

	}

	if &applicationLoadBalancerForwardingRule == nil {
		return errors.New("application loadbalancer forwarding rule not found")
	}

	if err = setApplicationLoadBalancerForwardingRuleData(d, &applicationLoadBalancerForwardingRule); err != nil {
		return err
	}

	return nil
}
