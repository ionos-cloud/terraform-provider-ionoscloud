package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"strings"
)

func dataSourceApplicationLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNatGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the application loadbalancer",
				Required:    true,
			},
			"listener_lan": {
				Type:        schema.TypeInt,
				Description: "Id of the listening LAN. (inbound)",
				Required:    true,
			},
			"ips": {
				Type: schema.TypeList,
				Description: "Collection of IP addresses of the Application Load Balancer. (inbound and outbound) IP of " +
					"the listenerLan must be a customer reserved IP for the public load balancer and private IP for the private load balancer.",
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target_lan": {
				Type:        schema.TypeInt,
				Description: "Id of the balanced private target LAN. (outbound)",
				Required:    true,
			},
			"lb_private_ips": {
				Type: schema.TypeList,
				Description: "Collection of private IP addresses with subnet mask of the Application Load Balancer. " +
					"IPs must contain valid subnet mask. If user will not provide any IP then the system will generate one IP with /24 subnet.",
				Optional: true,
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

func dataSourceApplicationLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

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
		return errors.New("please provide either the application loadbalancer id or name")
	}
	var applicationLoadBalancer ionoscloud.ApplicationLoadBalancer
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if idOk {
		/* search by ID */
		applicationLoadBalancer, _, err = client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, datacenterId.(string), id.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the nat gateway %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var applicationLoadBalancers ionoscloud.ApplicationLoadBalancers

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		applicationLoadBalancers, _, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersGet(ctx, datacenterId.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching nat gateway: %s", err.Error())
		}

		if applicationLoadBalancers.Items != nil {
			for _, c := range *applicationLoadBalancers.Items {
				tmpAlb, _, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, datacenterId.(string), *c.Id).Execute()
				if err != nil {
					return fmt.Errorf("an error occurred while fetching nat gateway with ID %s: %s", *c.Id, err.Error())
				}
				if applicationLoadBalancer.Properties.Name != nil {
					if strings.Contains(*applicationLoadBalancer.Properties.Name, name.(string)) {
						applicationLoadBalancer = tmpAlb
						break
					}
				}

			}
		}

	}

	if &applicationLoadBalancer == nil {
		return errors.New("application loadbalancer not found")
	}

	if err = setApplicationLoadBalancerData(d, &applicationLoadBalancer); err != nil {
		return err
	}

	return nil
}

func setApplicationLoadBalancerData(d *schema.ResourceData, applicationLoadBalancer *ionoscloud.ApplicationLoadBalancer) error {

	if applicationLoadBalancer.Id != nil {
		d.SetId(*applicationLoadBalancer.Id)
		if err := d.Set("id", *applicationLoadBalancer.Id); err != nil {
			return err
		}
	}

	if applicationLoadBalancer.Properties != nil {
		if applicationLoadBalancer.Properties.Name != nil {
			err := d.Set("name", *applicationLoadBalancer.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for application loadbalancer %s: %s", d.Id(), err)
			}
		}

		if applicationLoadBalancer.Properties.ListenerLan != nil {
			err := d.Set("listener_lan", *applicationLoadBalancer.Properties.ListenerLan)
			if err != nil {
				return fmt.Errorf("error while setting listener_lan property for application loadbalancer %s: %s", d.Id(), err)
			}
		}

		if applicationLoadBalancer.Properties.Ips != nil {
			err := d.Set("ips", *applicationLoadBalancer.Properties.Ips)
			if err != nil {
				return fmt.Errorf("error while setting ips property for application loadbalancer %s: %s", d.Id(), err)
			}
		}

		if applicationLoadBalancer.Properties.TargetLan != nil {
			err := d.Set("target_lan", *applicationLoadBalancer.Properties.TargetLan)
			if err != nil {
				return fmt.Errorf("error while setting target_lan property for application loadbalancer %s: %s", d.Id(), err)
			}
		}

		if applicationLoadBalancer.Properties.LbPrivateIps != nil {
			err := d.Set("lb_private_ips", *applicationLoadBalancer.Properties.LbPrivateIps)
			if err != nil {
				return fmt.Errorf("error while setting lb_private_ips property for application loadbalancer %s: %s", d.Id(), err)
			}
		}
	}
	return nil
}
