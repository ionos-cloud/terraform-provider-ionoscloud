package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	cloudapiflowlog "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/flowlog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceApplicationLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Application Load Balancer.",
				Optional:    true,
				Computed:    true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"listener_lan": {
				Type:        schema.TypeInt,
				Description: "ID of the listening (inbound) LAN.",
				Computed:    true,
			},
			"ips": {
				Type:        schema.TypeSet,
				Description: "Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target_lan": {
				Type:        schema.TypeInt,
				Description: "ID of the balanced private target LAN (outbound).",
				Computed:    true,
			},
			"lb_private_ips": {
				Type:        schema.TypeSet,
				Description: "Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"central_logging": {
				Type:        schema.TypeBool,
				Description: "Turn logging on and off for this product. Default value is 'false'.",
				Computed:    true,
			},
			"logging_format": {
				Type:        schema.TypeString,
				Description: "Specifies the format of the logs.",
				Computed:    true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flowlog": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     cloudapiflowlog.FlowlogSchemaDatasource,
				Description: `Flow logs holistically capture network information such as source and destination 
IP addresses, source and destination ports, number of packets, amount of bytes, 
the start and end time of the recording, and the type of protocol – 
and log the extent to which your instances are being accessed.`,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceApplicationLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	datacenterId := d.Get("datacenter_id").(string)

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return utils.ToDiags(d, "id and name cannot be both specified in the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the application load balancer id or name", nil)
	}

	var applicationLoadBalancer ionoscloud.ApplicationLoadBalancer
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for application load balancer by id %s", id)
		applicationLoadBalancer, apiResponse, err = client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, datacenterId, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the application load balancer while searching by ID %s: %s", id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		/* search by name */
		var results []ionoscloud.ApplicationLoadBalancer

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for application load balancer by name with partial_match %t and name: %s", partialMatch, name)

		if partialMatch {
			applicationLoadBalancers, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersGet(ctx, datacenterId).Depth(1).Filter("name", name).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching application load balancers: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
			}

			results = *applicationLoadBalancers.Items
		} else {
			applicationLoadBalancers, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersGet(ctx, datacenterId).Depth(1).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching application load balancers: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
			}

			if applicationLoadBalancers.Items != nil {
				for _, alb := range *applicationLoadBalancers.Items {
					if alb.Properties != nil && alb.Properties.Name != nil && strings.EqualFold(*alb.Properties.Name, name) {
						tmpAlb, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, datacenterId, *alb.Id).Execute()
						logApiRequestTime(apiResponse)
						if err != nil {
							return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching application load balancer with ID %s: %s", *alb.Id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
						}
						results = append(results, tmpAlb)
					}

				}
			}
		}

		if results == nil || len(results) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no application load balanacer found with the specified criteria: name = %s", name), nil)
		} else if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one application load balanacer found with the specified criteria: name = %s", name), nil)
		}

		applicationLoadBalancer = results[0]

	}
	fw := cloudapiflowlog.Service{
		Client: client,
		Meta:   meta,
		D:      d,
	}
	flowLog, apiResponse, err := fw.GetFlowLogForALB(ctx, datacenterId, *applicationLoadBalancer.Id, 2)
	if err != nil {
		if !apiResponse.HttpNotFound() {
			return utils.ToDiags(d, fmt.Sprintf("error finding flowlog for application loadbalancer: %s, %s", err, responseBody(apiResponse)), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	}
	if err = setApplicationLoadBalancerData(d, &applicationLoadBalancer, flowLog); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
