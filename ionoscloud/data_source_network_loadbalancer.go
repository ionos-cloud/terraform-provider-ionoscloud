package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"strings"
)

func dataSourceNetworkLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkLoadBalancerRead,
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
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceNetworkLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return errors.New("no datacenter_id was specified")
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the lan id or name")
	}
	var networkLoadBalancer ionoscloud.NetworkLoadBalancer
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if idOk {
		/* search by ID */
		networkLoadBalancer, _, err = client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, datacenterId.(string), id.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the network loadbalancer %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var networkLoadBalancers ionoscloud.NetworkLoadBalancers

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		networkLoadBalancers, _, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersGet(ctx, datacenterId.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching network loadbalancers: %s", err.Error())
		}

		if networkLoadBalancers.Items != nil {
			for _, c := range *networkLoadBalancers.Items {
				tmpNetworkLoadBalancer, _, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, datacenterId.(string), *c.Id).Execute()
				if err != nil {
					return fmt.Errorf("an error occurred while fetching network loadbalancer with ID %s: %s", *c.Id, err.Error())
				}
				if tmpNetworkLoadBalancer.Properties.Name != nil {
					if strings.Contains(*tmpNetworkLoadBalancer.Properties.Name, name.(string)) {
						networkLoadBalancer = tmpNetworkLoadBalancer
						break
					}
				}

			}
		}

	}

	if &networkLoadBalancer == nil {
		return errors.New("network loadbalancer not found")
	}

	if err = setNetworkLoadBalancerData(d, &networkLoadBalancer, client); err != nil {
		return err
	}

	return nil
}

func setNetworkLoadBalancerData(d *schema.ResourceData, networkLoadBalancer *ionoscloud.NetworkLoadBalancer, client *ionoscloud.APIClient) error {

	if networkLoadBalancer.Id != nil {
		d.SetId(*networkLoadBalancer.Id)
		if err := d.Set("id", *networkLoadBalancer.Id); err != nil {
			return err
		}
	}

	if networkLoadBalancer.Properties != nil {
		if networkLoadBalancer.Properties.Name != nil {
			err := d.Set("name", *networkLoadBalancer.Properties.Name)
			if err != nil {
				return fmt.Errorf("Error while setting name property for network load balancer %s: %s", d.Id(), err)
			}
		}

		if networkLoadBalancer.Properties.ListenerLan != nil {
			err := d.Set("listener_lan", *networkLoadBalancer.Properties.ListenerLan)
			if err != nil {
				return fmt.Errorf("Error while setting listener_lan property for network load balancer %s: %s", d.Id(), err)
			}
		}

		if networkLoadBalancer.Properties.TargetLan != nil {
			err := d.Set("target_lan", *networkLoadBalancer.Properties.TargetLan)
			if err != nil {
				return fmt.Errorf("Error while setting target_lan property for network load balancer %s: %s", d.Id(), err)
			}
		}

		if networkLoadBalancer.Properties.Ips != nil {
			err := d.Set("ips", *networkLoadBalancer.Properties.Ips)
			if err != nil {
				return fmt.Errorf("Error while setting ips property for network load balancer %s: %s", d.Id(), err)
			}
		}

		if networkLoadBalancer.Properties.LbPrivateIps != nil {
			err := d.Set("lb_private_ips", *networkLoadBalancer.Properties.LbPrivateIps)
			if err != nil {
				return fmt.Errorf("Error while setting lb_private_ips property for network load balancer %s: %s", d.Id(), err)
			}
		}

	}
	return nil
}
