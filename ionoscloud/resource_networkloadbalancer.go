package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
	cloudapiflowlog "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/flowlog"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
				Type:             schema.TypeString,
				Description:      "A name of that Network Load Balancer",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"flowlog": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     cloudapiflowlog.FlowlogSchemaResource,
				MaxItems: 1,
				Description: `Only 1 flow log can be configured. Only the name field can change as part of an update. Flow logs holistically capture network information such as source and destination 
IP addresses, source and destination ports, number of packets, amount of bytes, 
the start and end time of the recording, and the type of protocol – 
and log the extent to which your instances are being accessed.`,
			},
		},
		Timeouts:      &resourceDefaultTimeouts,
		CustomizeDiff: ForceNewForFlowlogChanges,
	}
}

func resourceNetworkLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

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

	if flowLogs, ok := d.GetOk("flowlog"); ok {
		networkLoadBalancer.Entities = &ionoscloud.NetworkLoadBalancerEntities{
			Flowlogs: &ionoscloud.FlowLogs{
				Items: &[]ionoscloud.FlowLog{},
			},
		}
		if flowLogList, ok := flowLogs.([]any); ok {
			for _, flowLogData := range flowLogList {
				if flowLog, ok := flowLogData.(map[string]interface{}); ok {
					*networkLoadBalancer.Entities.Flowlogs.Items = append(*networkLoadBalancer.Entities.Flowlogs.Items, cloudapiflowlog.GetFlowlogFromMap(flowLog))
				}
			}
		}
	}
	dcId := d.Get("datacenter_id").(string)

	networkLoadBalancerResp, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersPost(ctx, dcId).NetworkLoadBalancer(networkLoadBalancer).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating network loadbalancer: %w, %s", err, responseBody(apiResponse)))
		return diags
	}

	d.SetId(*networkLoadBalancerResp.Id)

	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if cloudapi.IsRequestFailed(errState) {
			d.SetId("")
		}
		return diag.FromErr(errState)
	}

	return resourceNetworkLoadBalancerRead(ctx, d, meta)
}

func resourceNetworkLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)

	networkLoadBalancer, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, dcId, d.Id()).Depth(3).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if httpNotFound(apiResponse) {
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
	client := meta.(services.SdkBundle).CloudApiClient
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
		vInt := int32(v.(int))
		request.Properties.ListenerLan = &vInt
	}

	if d.HasChange("target_lan") {
		_, v := d.GetChange("target_lan")
		vInt := int32(v.(int))
		request.Properties.TargetLan = &vInt
	}

	if d.HasChange("ips") {
		oldIps, newIps := d.GetChange("ips")
		log.Printf("[INFO] network loadbalancer ips changed from %+v to %+v", oldIps, newIps)
		ipsVal := newIps.([]interface{})
		ips := make([]string, 0)
		if ipsVal != nil {
			for _, ip := range ipsVal {
				ips = append(ips, ip.(string))
			}
		}
		if len(ips) > 0 {
			request.Properties.Ips = &ips
		} else {
			diags := diag.FromErr(fmt.Errorf("you can not empty the ips field for networkloadbalancer %s", d.Id()))
			return diags
		}
	}

	if d.HasChange("lb_private_ips") {
		oldLbPrivateIps, newLbPrivateIps := d.GetChange("lb_private_ips")
		log.Printf("[INFO] network loadbalancer lb_private_ips changed from %+v to %+v", oldLbPrivateIps, newLbPrivateIps)
		lbPrivateIpsVal := newLbPrivateIps.([]interface{})
		lbPrivateIps := make([]string, 0)
		if lbPrivateIpsVal != nil {
			for _, privateIp := range lbPrivateIpsVal {
				lbPrivateIps = append(lbPrivateIps, privateIp.(string))
			}
		}
		if len(lbPrivateIps) == 0 {
			diags := diag.FromErr(fmt.Errorf("you can not empty the lbPrivateIps field for networkloadbalancer %s", d.Id()))
			return diags
		}
		request.Properties.LbPrivateIps = &lbPrivateIps
	}

	if d.HasChange("flowlog") {
		old, newV := d.GetChange("flowlog")
		var firstFlowLogId = ""
		if old != nil && len(old.([]any)) > 0 {
			firstFlowLogId = old.([]any)[0].(map[string]any)["id"].(string)
		}

		if newV.([]any) != nil && len(newV.([]any)) > 0 {
			for _, val := range newV.([]any) {
				if flowLogMap, ok := val.(map[string]any); ok {
					flowLog := cloudapiflowlog.GetFlowlogFromMap(flowLogMap)
					fw := cloudapiflowlog.Service{
						D:      d,
						Client: client,
					}
					err := fw.CreateOrPatchForNLB(ctx, dcId, d.Id(), firstFlowLogId, flowLog)
					if err != nil {
						//if we have a create that failed, we do not want to save in state
						// saving in state would mean a diff that would force a re-create
						if firstFlowLogId == "" {
							_ = d.Set("flowlog", nil)
						}
						return diag.FromErr(err)
					}
				}
			}
		}
	}

	_, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersPatch(ctx, dcId, d.Id()).NetworkLoadBalancerProperties(*request.Properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a network loadbalancer ID %s %s \n ApiError: %s", d.Id(), err, responseBody(apiResponse)))
		return diags
	}

	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return diag.FromErr(errState)
	}

	return resourceNetworkLoadBalancerRead(ctx, d, meta)
}

func resourceNetworkLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)

	apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersDelete(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a network loadbalancer %s %w", d.Id(), err))
		return diags
	}

	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		return diag.FromErr(errState)
	}

	d.SetId("")

	return nil
}

func resourceNetworkLoadBalancerImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).CloudApiClient

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
		if httpNotFound(apiResponse) {
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

func setNetworkLoadBalancerData(d *schema.ResourceData, networkLoadBalancer *ionoscloud.NetworkLoadBalancer) error {

	if networkLoadBalancer.Id != nil {
		d.SetId(*networkLoadBalancer.Id)
	}

	if networkLoadBalancer.Properties != nil {
		if networkLoadBalancer.Properties.Name != nil {
			err := d.Set("name", *networkLoadBalancer.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for network load balancer %s: %w", d.Id(), err)
			}
		}

		if networkLoadBalancer.Properties.ListenerLan != nil {
			err := d.Set("listener_lan", *networkLoadBalancer.Properties.ListenerLan)
			if err != nil {
				return fmt.Errorf("error while setting listener_lan property for network load balancer %s: %w", d.Id(), err)
			}
		}

		if networkLoadBalancer.Properties.TargetLan != nil {
			err := d.Set("target_lan", *networkLoadBalancer.Properties.TargetLan)
			if err != nil {
				return fmt.Errorf("error while setting target_lan property for network load balancer %s: %w", d.Id(), err)
			}
		}

		if networkLoadBalancer.Properties.Ips != nil {
			err := d.Set("ips", *networkLoadBalancer.Properties.Ips)
			if err != nil {
				return fmt.Errorf("error while setting ips property for network load balancer %s: %w", d.Id(), err)
			}
		}

		if networkLoadBalancer.Properties.LbPrivateIps != nil {
			err := d.Set("lb_private_ips", *networkLoadBalancer.Properties.LbPrivateIps)
			if err != nil {
				return fmt.Errorf("error while setting lb_private_ips property for network load balancer %s: %w", d.Id(), err)
			}
		}
		if networkLoadBalancer.Entities != nil && networkLoadBalancer.Entities.Flowlogs != nil &&
			networkLoadBalancer.Entities.Flowlogs.Items != nil && len(*networkLoadBalancer.Entities.Flowlogs.Items) > 0 {
			var flowlogs []map[string]any
			for _, flowLog := range *networkLoadBalancer.Entities.Flowlogs.Items {
				result := map[string]any{}
				result, err := utils.DecodeStructToMap(flowLog.Properties)
				if err != nil {
					return err
				}
				result["id"] = *flowLog.Id
				flowlogs = append(flowlogs, result)
			}
			if err := d.Set("flowlog", flowlogs); err != nil {
				return fmt.Errorf("error setting flowlog %w", err)
			}
		}
	}
	return nil
}
