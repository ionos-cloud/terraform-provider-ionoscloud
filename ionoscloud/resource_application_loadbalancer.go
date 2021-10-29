package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
)

func resourceApplicationLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationLoadBalancerCreate,
		ReadContext:   resourceApplicationLoadBalancerRead,
		UpdateContext: resourceApplicationLoadBalancerUpdate,
		DeleteContext: resourceApplicationLoadBalancerDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "Name of the application loadbalancer",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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

func resourceApplicationLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	applicationLoadBalancer := ionoscloud.ApplicationLoadBalancer{
		Properties: &ionoscloud.ApplicationLoadBalancerProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		applicationLoadBalancer.Properties.Name = &name
	} else {
		diags := diag.FromErr(fmt.Errorf("name must be provided for application loadbalancer"))
		return diags
	}

	if listenerLan, listenerLanOk := d.GetOk("listener_lan"); listenerLanOk {
		listener := int32(listenerLan.(int))
		applicationLoadBalancer.Properties.ListenerLan = &listener
	} else {
		diags := diag.FromErr(fmt.Errorf("listener_lan must be provided for application loadbalancer"))
		return diags
	}

	if ipsVal, ipsOk := d.GetOk("ips"); ipsOk {
		ipsVal := ipsVal.([]interface{})
		if ipsVal != nil {
			ips := make([]string, 0)
			for _, value := range ipsVal {
				ips = append(ips, value.(string))
			}
			if len(ips) > 0 {
				applicationLoadBalancer.Properties.Ips = &ips
			}
		}
	}

	if targetLan, targetLanOk := d.GetOk("target_lan"); targetLanOk {
		targetLan := int32(targetLan.(int))
		applicationLoadBalancer.Properties.TargetLan = &targetLan
	} else {
		diags := diag.FromErr(fmt.Errorf("target_lan must be provided for application loadbalancer"))
		return diags
	}

	if privateIpsVal, privateIpsOk := d.GetOk("lb_private_ips"); privateIpsOk {
		privateIpsVal := privateIpsVal.([]interface{})
		if privateIpsVal != nil {
			privateIps := make([]string, 0)
			for _, value := range privateIpsVal {
				privateIps = append(privateIps, value.(string))
			}
			if len(privateIps) > 0 {
				applicationLoadBalancer.Properties.LbPrivateIps = &privateIps
			}
		}
	}

	dcId := d.Get("datacenter_id").(string)

	applicationLoadbalancer, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersPost(ctx, dcId).ApplicationLoadBalancer(applicationLoadBalancer).Execute()

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating application loadbalancer: %s, %s", err, responseBody(apiResponse)))
		return diags
	}

	d.SetId(*applicationLoadbalancer.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceApplicationLoadBalancerRead(ctx, d, meta)
}

func resourceApplicationLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)

	applicationLoadBalancer, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, dcId, d.Id()).Execute()

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] Successfully retreived application loadbalancer %s: %+v", d.Id(), applicationLoadBalancer)

	if err := setApplicationLoadBalancerData(d, &applicationLoadBalancer); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceApplicationLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	request := ionoscloud.ApplicationLoadBalancer{
		Properties: &ionoscloud.ApplicationLoadBalancerProperties{},
	}

	dcId := d.Get("datacenter_id").(string)

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		vStr := v.(string)
		request.Properties.Name = &vStr
	}

	if d.HasChange("listener_lan") {
		_, v := d.GetChange("listener_lan")
		vInt := int32(v.(int))
		request.Properties.ListenerLan = &vInt
	}

	if d.HasChange("ips") {
		_, newIps := d.GetChange("ips")
		ipsVal := newIps.([]interface{})
		if ipsVal != nil {
			ips := make([]string, 0)
			for _, value := range ipsVal {
				ips = append(ips, value.(string))
			}
			if len(ips) > 0 {
				request.Properties.Ips = &ips
			}
		}
	}

	if d.HasChange("target_lan") {
		_, v := d.GetChange("target_lan")
		vInt := int32(v.(int))
		request.Properties.TargetLan = &vInt
	}

	if d.HasChange("lb_private_ips") {
		_, newPrivateIps := d.GetChange("lb_private_ips")
		privateIpsVal := newPrivateIps.([]interface{})
		if privateIpsVal != nil {
			privateIps := make([]string, 0)
			for _, value := range privateIpsVal {
				privateIps = append(privateIps, value.(string))
			}
			if len(privateIps) > 0 {
				request.Properties.LbPrivateIps = &privateIps
			}
		}
	}

	_, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersPatch(ctx, dcId, d.Id()).ApplicationLoadBalancerProperties(*request.Properties).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating application loadbalancer ID %s %s", d.Id(), err))
		return diags
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceApplicationLoadBalancerRead(ctx, d, meta)
}

func resourceApplicationLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)

	apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersDelete(ctx, dcId, d.Id()).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting an application loadbalancer %s %s", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	d.SetId("")

	return nil
}

func setApplicationLoadBalancerData(d *schema.ResourceData, applicationLoadBalancer *ionoscloud.ApplicationLoadBalancer) error {

	if applicationLoadBalancer.Id != nil {
		d.SetId(*applicationLoadBalancer.Id)
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
