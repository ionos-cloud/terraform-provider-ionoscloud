package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	cloudapiflowlog "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/flowlog"
)

func dataSourceNetworkLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_lan": {
				Type:        schema.TypeInt,
				Description: "Id of the listening LAN. (inbound)",
				Computed:    true,
			},
			"ips": {
				Type: schema.TypeList,
				Description: "Collection of IP addresses of the Network Load Balancer. (inbound and outbound) IP of the " +
					"listenerLan must be a customer reserved IP for the public load balancer and private IP " +
					"for the private load balancer.",
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target_lan": {
				Type:        schema.TypeInt,
				Description: "Id of the balanced private target LAN. (outbound)",
				Computed:    true,
			},
			"lb_private_ips": {
				Type: schema.TypeList,
				Description: "Collection of private IP addresses with subnet mask of the Network Load Balancer. IPs " +
					"must contain valid subnet mask. If user will not provide any IP then the system will " +
					"generate one IP with /24 subnet.",
				Computed: true,
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
				ForceNew: true,
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

func dataSourceNetworkLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudAPIClient

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return diag.FromErr(errors.New("no datacenter_id was specified"))
	}
	dcID := datacenterId.(string)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the lan id or name"))
	}
	var networkLoadBalancer ionoscloud.NetworkLoadBalancer
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		networkLoadBalancer, apiResponse, err = client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, dcID, id.(string)).Depth(4).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the network loadbalancer %s, dcID %s : %w", id.(string), dcID, err))
		}
	} else {
		/* search by name */
		var networkLoadBalancers ionoscloud.NetworkLoadBalancers

		networkLoadBalancers, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersGet(ctx, dcID).Depth(4).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancers, dcID %s : %w", dcID, err))
		}

		var results []ionoscloud.NetworkLoadBalancer
		if networkLoadBalancers.Items != nil {
			for _, nlb := range *networkLoadBalancers.Items {
				if nlb.Properties != nil && nlb.Properties.Name != nil && strings.EqualFold(*nlb.Properties.Name, name.(string)) {
					tmpNetworkLoadBalancer, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, dcID, *nlb.Id).Depth(4).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diag.FromErr(fmt.Errorf("an error occurred while fetching network loadbalancer with ID %s: %w", *nlb.Id, err))
					}
					results = append(results, tmpNetworkLoadBalancer)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no network load balancer found with the specified criteria: name = %s", name.(string)))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one network load balancer found with the specified criteria: name = %s", name.(string)))
		} else {
			networkLoadBalancer = results[0]
		}
	}

	if err = setNetworkLoadBalancerData(d, &networkLoadBalancer); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
