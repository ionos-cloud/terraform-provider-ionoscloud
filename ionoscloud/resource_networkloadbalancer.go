package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
)

func resourceNetworkLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkLoadBalancerCreate,
		ReadContext:   resourceNetworkLoadBalancerRead,
		UpdateContext: resourceNetworkLoadBalancerUpdate,
		DeleteContext: resourceNetworkLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNetworkLoadBalancerImport,
		},
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

func resourceNetworkLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	networkLoadBalancer := ionoscloud.NetworkLoadBalancer{
		Properties: &ionoscloud.NetworkLoadBalancerProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		networkLoadBalancer.Properties.Name = &name
	} else {
		diags := diag.FromErr(fmt.Errorf("name must be provided for network loadbalancer"))
		return diags
	}

	if listenerLan, listenerLanOk := d.GetOk("listener_lan"); listenerLanOk {
		listenerLan := int32(listenerLan.(int))
		networkLoadBalancer.Properties.ListenerLan = &listenerLan
	} else {
		diags := diag.FromErr(fmt.Errorf("listener lan must be provided for network loadbalancer"))
		return diags
	}

	if targetLan, targetLanOk := d.GetOk("target_lan"); targetLanOk {
		targetLan := int32(targetLan.(int))
		networkLoadBalancer.Properties.TargetLan = &targetLan
	} else {
		diags := diag.FromErr(fmt.Errorf("target lan must be provided for network loadbalancer"))
		return diags
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

	networkLoadBalancerResp, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersPost(ctx, dcId).NetworkLoadBalancer(networkLoadBalancer).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating network loadbalancer: %s, %s", err, responseBody(apiResponse)))
		return diags
	}

	d.SetId(*networkLoadBalancerResp.Id)

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

	return resourceNetworkLoadBalancerRead(ctx, d, meta)
}

func resourceNetworkLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)

	networkLoadBalancer, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] Successfully retreived network load balancer %s: %+v", d.Id(), networkLoadBalancer)

	if err := setNetworkLoadBalancerData(d, &networkLoadBalancer); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetworkLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)
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
	_, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersPatch(ctx, dcId, d.Id()).NetworkLoadBalancerProperties(*request.Properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a network loadbalancer ID %s %s \n ApiError: %s", d.Id(), err, responseBody(apiResponse)))
		return diags
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceNetworkLoadBalancerRead(ctx, d, meta)
}

func resourceNetworkLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcId := d.Get("datacenter_id").(string)

	apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersDelete(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a network loadbalancer %s %s", d.Id(), err))
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

func resourceNetworkLoadBalancerImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{networkloadbalancer}", d.Id())
	}

	dcId := parts[0]
	networkLoadBalancerId := parts[1]

	networkLoadBalancer, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, dcId, networkLoadBalancerId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find network load balancer %q", networkLoadBalancerId)
		}
		return nil, fmt.Errorf("an error occured while retrieving network load balancer  %q: %q ", networkLoadBalancerId, err)
	}

	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, err
	}

	if err := setNetworkLoadBalancerData(d, &networkLoadBalancer); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
