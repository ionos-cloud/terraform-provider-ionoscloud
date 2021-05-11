package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
)

func resourceNetworkLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkLoadBalancerCreate,
		Read:   resourceNetworkLoadBalancerRead,
		Update: resourceNetworkLoadBalancerUpdate,
		Delete: resourceNetworkLoadBalancerDelete,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Description: "A name of that Network Load Balancer",
				Required:    true,
			},
			"listener_lan": {
				Type:        schema.TypeInt,
				Description: "Id of the listening LAN. (inbound)",
				Required:    true,
			},
			"ips": {
				Type: schema.TypeList,
				Description: "Collection of IP addresses of the Network Load Balancer. (inbound and outbound) IP of the " +
					"listenerLan must be a customer reserved IP for the public load balancer and private IP " +
					"for the private load balancer.",
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
				Description: "Collection of private IP addresses with subnet mask of the Network Load Balancer. IPs " +
					"must contain valid subnet mask. If user will not provide any IP then the system will " +
					"generate one IP with /24 subnet.",
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

func resourceNetworkLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	networkLoadBalancer := ionoscloud.NetworkLoadBalancer{
		Properties: &ionoscloud.NetworkLoadBalancerProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		networkLoadBalancer.Properties.Name = &name
	} else {
		return fmt.Errorf("Name must be provided for network loadbalancer")
	}

	if listenerLan, listenerLanOk := d.GetOk("listener_lan"); listenerLanOk {
		listenerLan := int32(listenerLan.(int))
		networkLoadBalancer.Properties.ListenerLan = &listenerLan
	} else {
		return fmt.Errorf("Listener lan must be provided for network loadbalancer")
	}

	if targetLan, targetLanOk := d.GetOk("target_lan"); targetLanOk {
		targetLan := int32(targetLan.(int))
		networkLoadBalancer.Properties.TargetLan = &targetLan
	} else {
		return fmt.Errorf("Target lan must be provided for network loadbalancer")
	}

	if ipsVal, ipsOk := d.GetOk("ips"); ipsOk {
		ipsVal := ipsVal.([]interface{})
		if ipsVal != nil {
			ips := make([]string, len(ipsVal), len(ipsVal))
			for idx := range ipsVal {
				ips[idx] = fmt.Sprint(ipsVal[idx])
			}
			networkLoadBalancer.Properties.Ips = &ips
		}
	}

	if lbPrivateIpsVal, lbPrivateIpsOk := d.GetOk("lb_private_ips"); lbPrivateIpsOk {
		lbPrivateIpsVal := lbPrivateIpsVal.([]interface{})
		if lbPrivateIpsVal != nil {
			lbPrivateIps := make([]string, len(lbPrivateIpsVal), len(lbPrivateIpsVal))
			for idx := range lbPrivateIpsVal {
				lbPrivateIps[idx] = fmt.Sprint(lbPrivateIpsVal[idx])
			}
			networkLoadBalancer.Properties.LbPrivateIps = &lbPrivateIps
		}
	}

	dcId := d.Get("datacenter_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)
	if cancel != nil {
		defer cancel()
	}

	networkLoadBalancerResp, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersPost(ctx, dcId).NetworkLoadBalancer(networkLoadBalancer).Execute()

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error creating network loadbalancer: %s", err)
	}

	d.SetId(*networkLoadBalancerResp.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceNetworkLoadBalancerRead(d, meta)
}

func resourceNetworkLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	dcId := d.Get("datacenter_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	networkLoadBalancer, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, dcId, d.Id()).Execute()

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
	}

	log.Printf("[INFO] Successfully retreived network load balancer %s: %+v", d.Id(), networkLoadBalancer)

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

	return nil
}

func resourceNetworkLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client
	request := ionoscloud.NetworkLoadBalancer{
		Properties: &ionoscloud.NetworkLoadBalancerProperties{},
	}

	dcId := d.Get("datacenter_id").(string)

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		vStr := v.(string)
		request.Properties.Name = &vStr
	}

	if d.HasChange("listener_lan") {
		_, v := d.GetChange("listener_lan")
		vStr := v.(string)
		request.Properties.Name = &vStr
	}

	if d.HasChange("target_lan") {
		_, v := d.GetChange("target_lan")
		vStr := v.(string)
		request.Properties.Name = &vStr
	}

	if d.HasChange("ips") {
		oldIps, newIps := d.GetChange("ips")
		log.Printf("[INFO] network loadbalancer ips changed from %+v to %+v", oldIps, newIps)
		ipsVal := newIps.([]interface{})
		if ipsVal != nil {
			ips := make([]string, len(ipsVal), len(ipsVal))
			for idx := range ipsVal {
				ips[idx] = fmt.Sprint(ipsVal[idx])
			}
			request.Properties.Ips = &ips
		}
	}

	if d.HasChange("lb_private_ips") {
		oldLbPrivateIps, newLbPrivateIps := d.GetChange("lb_private_ips")
		log.Printf("[INFO] network loadbalancer lb_private_ips changed from %+v to %+v", oldLbPrivateIps, newLbPrivateIps)
		lbPrivateIpsVal := newLbPrivateIps.([]interface{})
		if lbPrivateIpsVal != nil {
			lbPrivateIps := make([]string, len(lbPrivateIpsVal), len(lbPrivateIpsVal))
			for idx := range lbPrivateIpsVal {
				lbPrivateIps[idx] = fmt.Sprint(lbPrivateIpsVal[idx])
			}
			request.Properties.LbPrivateIps = &lbPrivateIps
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

	if cancel != nil {
		defer cancel()
	}
	_, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersPatch(ctx, dcId, d.Id()).NetworkLoadBalancerProperties(*request.Properties).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while updating a network loadbalancer ID %s %s \n ApiError: %s", d.Id(), err, string(apiResponse.Payload))
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceNetworkLoadBalancerRead(d, meta)
}

func resourceNetworkLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	dcId := d.Get("datacenter_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersDelete(ctx, dcId, d.Id()).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while deleting a network loadbalancer %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")

	return nil
}
