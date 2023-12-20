package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
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
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Application Load Balancer.",
				Optional:    true,
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
	client := meta.(services.SdkBundle).CloudApiClient

	datacenterId := d.Get("datacenter_id").(string)

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the application load balancer id or name"))
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
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the application load balancer while searching by ID %s: %w", id, err))
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
				return diag.FromErr(fmt.Errorf("an error occurred while fetching application load balancers: %w", err))
			}

			results = *applicationLoadBalancers.Items
		} else {
			applicationLoadBalancers, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersGet(ctx, datacenterId).Depth(1).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching application load balancers: %w", err))
			}

			if applicationLoadBalancers.Items != nil {
				for _, alb := range *applicationLoadBalancers.Items {
					if alb.Properties != nil && alb.Properties.Name != nil && strings.ToLower(*alb.Properties.Name) == strings.ToLower(name) {
						tmpAlb, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, datacenterId, *alb.Id).Execute()
						logApiRequestTime(apiResponse)
						if err != nil {
							return diag.FromErr(fmt.Errorf("an error occurred while fetching application load balancer with ID %s: %w", *alb.Id, err))
						}
						results = append(results, tmpAlb)
					}

				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no application load balanacer found with the specified criteria: name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one application load balanacer found with the specified criteria: name = %s", name))
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
			diags := diag.FromErr(fmt.Errorf("error finding flowlog for application loadbalancer: %w, %s", err, responseBody(apiResponse)))
			return diags
		}
	}
	if err = setApplicationLoadBalancerData(d, &applicationLoadBalancer, flowLog); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
