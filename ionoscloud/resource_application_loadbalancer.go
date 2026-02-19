package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	cloudapiflowlog "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/flowlog"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceApplicationLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationLoadBalancerCreate,
		ReadContext:   resourceApplicationLoadBalancerRead,
		UpdateContext: resourceApplicationLoadBalancerUpdate,
		DeleteContext: resourceApplicationLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationLoadBalancerImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The name of the Application Load Balancer.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"listener_lan": {
				Type:        schema.TypeInt,
				Description: "ID of the listening (inbound) LAN.",
				Required:    true,
			},
			"ips": {
				Type:        schema.TypeSet,
				Description: "Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target_lan": {
				Type:        schema.TypeInt,
				Description: "ID of the balanced private target LAN (outbound).",
				Required:    true,
			},
			"lb_private_ips": {
				Type:        schema.TypeSet,
				Description: "Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"central_logging": {
				Type:        schema.TypeBool,
				Description: "Turn logging on and off for this product. Default value is 'false'.",
				Optional:    true,
				Default:     false,
			},
			"logging_format": {
				Type:        schema.TypeString,
				Description: "Specifies the format of the logs.",
				Optional:    true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceApplicationLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

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

	if centralLogging, centralLoggingOk := d.GetOk("central_logging"); centralLoggingOk {
		centralLogging := centralLogging.(bool)
		applicationLoadBalancer.Properties.CentralLogging = &centralLogging
	}

	if loggingFormat, loggingFormatOk := d.GetOk("logging_format"); loggingFormatOk {
		loggingFormat := loggingFormat.(string)
		applicationLoadBalancer.Properties.LoggingFormat = &loggingFormat
	}

	if ipsVal, ipsOk := d.GetOk("ips"); ipsOk {
		ipsVal := ipsVal.(*schema.Set).List()
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
		privateIpsVal := privateIpsVal.(*schema.Set).List()
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
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating application loadbalancer: %w, %s", err, responseBody(apiResponse)))
		return diags
	}

	d.SetId(*applicationLoadbalancer.Id)

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			d.SetId("")
		}
		return diag.FromErr(errState)
	}

	if flowLogs, ok := d.GetOk("flowlog"); ok {
		fw := cloudapiflowlog.Service{
			D:      d,
			Client: client,
		}
		if flowLogList, ok := flowLogs.([]any); ok {
			for _, flowLogData := range flowLogList {
				if flowLogMap, ok := flowLogData.(map[string]interface{}); ok {
					flowLog := cloudapiflowlog.GetFlowlogFromMap(flowLogMap)
					err := fw.CreateOrPatchForALB(ctx, dcId, d.Id(), "", flowLog)
					if err != nil {
						_ = d.Set("", nil)
						// flowlog creation failed, delete the alb
						diags := resourceApplicationLoadBalancerDelete(ctx, d, meta)
						if diags != nil {
							log.Printf("[ERROR] could not delete alb %v", diags)
						}
						diags = diag.FromErr(fmt.Errorf("error creating flowlog for application loadbalancer: %w, %s", err, responseBody(apiResponse)))
						return diags
					}
				}
			}
		}
	}

	return resourceApplicationLoadBalancerRead(ctx, d, meta)
}

func resourceApplicationLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	dcId := d.Get("datacenter_id").(string)

	applicationLoadBalancer, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] Successfully retrieved application loadbalancer %s: %+v", d.Id(), applicationLoadBalancer)
	fw := cloudapiflowlog.Service{
		Client: client,
		Meta:   meta,
		D:      d,
	}
	flowLog, apiResponse, err := fw.GetFlowLogForALB(ctx, dcId, *applicationLoadBalancer.Id, 1)
	if err != nil {
		if !apiResponse.HttpNotFound() {
			diags := diag.FromErr(fmt.Errorf("error finding flowlog for application loadbalancer: %w, %s", err, responseBody(apiResponse)))
			return diags
		}
	}

	return diag.FromErr(setApplicationLoadBalancerData(d, &applicationLoadBalancer, flowLog))
}

func resourceApplicationLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

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
		ipsVal := newIps.(*schema.Set).List()
		ips := make([]string, 0)
		if ipsVal != nil {
			for _, value := range ipsVal {
				ips = append(ips, value.(string))
			}
		}
		request.Properties.Ips = &ips
	}

	if d.HasChange("central_logging") {
		_, v := d.GetChange("central_logging")
		vBool := v.(bool)
		request.Properties.CentralLogging = &vBool
	}

	if d.HasChange("loggingFormat") {
		_, v := d.GetChange("loggingFormat")
		vStr := v.(string)
		request.Properties.LoggingFormat = &vStr
	}

	if d.HasChange("target_lan") {
		_, v := d.GetChange("target_lan")
		vInt := int32(v.(int))
		request.Properties.TargetLan = &vInt
	}

	if d.HasChange("lb_private_ips") {
		_, newPrivateIps := d.GetChange("lb_private_ips")
		privateIpsVal := newPrivateIps.(*schema.Set).List()
		privateIps := make([]string, 0)

		if privateIpsVal != nil {
			for _, value := range privateIpsVal {
				privateIps = append(privateIps, value.(string))
			}
		}
		request.Properties.LbPrivateIps = &privateIps
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
					err := fw.CreateOrPatchForALB(ctx, dcId, d.Id(), firstFlowLogId, flowLog)
					if err != nil {
						// if we have a create that failed, we do not want to save in state
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
	_, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersPatch(ctx, dcId, d.Id()).ApplicationLoadBalancerProperties(*request.Properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while updating application loadbalancer ID %s %w", d.Id(), err))
		return diags
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return diag.FromErr(errState)
	}

	return resourceApplicationLoadBalancerRead(ctx, d, meta)
}

func resourceApplicationLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	dcId := d.Get("datacenter_id").(string)

	apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersDelete(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while deleting an application loadbalancer %s %w", d.Id(), err))
		return diags
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		return diag.FromErr(errState)
	}

	d.SetId("")

	return nil
}

func resourceApplicationLoadBalancerImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()

	location, parts := splitImportID(importID, "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf(
			"invalid import identifier: expected one of <location>:<datacenter-id>/<alb> or "+
				"<datacenter-id>/<alb>, got: %s", importID,
		)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, fmt.Errorf("failed validating import identifier %q: %w", importID, err)
	}

	datacenterId := parts[0]
	albId := parts[1]

	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	alb, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, datacenterId, albId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find alb %q", albId)
		}
		return nil, fmt.Errorf("an error occurred while retrieving the alb %q, %w", albId, err)
	}

	if err := d.Set("datacenter_id", datacenterId); err != nil {
		return nil, fmt.Errorf("error while setting datacenter_id property for alb %q: %w", albId, err)
	}
	fw := cloudapiflowlog.Service{
		Client: client,
		Meta:   meta,
		D:      d,
	}
	flowLog, apiResponse, err := fw.GetFlowLogForALB(ctx, datacenterId, d.Id(), 0)
	if err != nil {
		if !apiResponse.HttpNotFound() {
			return nil, fmt.Errorf("error finding flowlog for application loadbalancer: %w, %s", err, responseBody(apiResponse))
		}
	}
	if err := setApplicationLoadBalancerData(d, &alb, flowLog); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
func setApplicationLoadBalancerData(d *schema.ResourceData, applicationLoadBalancer *ionoscloud.ApplicationLoadBalancer, flowLog *ionoscloud.FlowLog) error {

	if applicationLoadBalancer.Id != nil {
		d.SetId(*applicationLoadBalancer.Id)
	}

	if applicationLoadBalancer.Properties != nil {
		if applicationLoadBalancer.Properties.Name != nil {
			err := d.Set("name", *applicationLoadBalancer.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for application loadbalancer %s: %w", d.Id(), err)
			}
		}

		if applicationLoadBalancer.Properties.ListenerLan != nil {
			err := d.Set("listener_lan", *applicationLoadBalancer.Properties.ListenerLan)
			if err != nil {
				return fmt.Errorf("error while setting listener_lan property for application loadbalancer %s: %w", d.Id(), err)
			}
		}

		if applicationLoadBalancer.Properties.Ips != nil {
			err := d.Set("ips", *applicationLoadBalancer.Properties.Ips)
			if err != nil {
				return fmt.Errorf("error while setting ips property for application loadbalancer %s: %w", d.Id(), err)
			}
		}

		if applicationLoadBalancer.Properties.TargetLan != nil {
			err := d.Set("target_lan", *applicationLoadBalancer.Properties.TargetLan)
			if err != nil {
				return fmt.Errorf("error while setting target_lan property for application loadbalancer %s: %w", d.Id(), err)
			}
		}

		if applicationLoadBalancer.Properties.LbPrivateIps != nil {
			err := d.Set("lb_private_ips", *applicationLoadBalancer.Properties.LbPrivateIps)
			if err != nil {
				return fmt.Errorf("error while setting lb_private_ips property for application loadbalancer %s: %w", d.Id(), err)
			}
		}

		if applicationLoadBalancer.Properties.CentralLogging != nil {
			err := d.Set("central_logging", *applicationLoadBalancer.Properties.CentralLogging)
			if err != nil {
				return fmt.Errorf("error while setting central_logging property for network load balancer %s: %w", d.Id(), err)
			}
		}

		if applicationLoadBalancer.Properties.LoggingFormat != nil {
			err := d.Set("logging_format", *applicationLoadBalancer.Properties.LoggingFormat)
			if err != nil {
				return fmt.Errorf("error while setting logging_format property for network load balancer %s: %w", d.Id(), err)
			}
		}
	}

	if flowLog != nil {
		var flowlogs []map[string]any
		result, err := utils.DecodeStructToMap(flowLog.Properties)
		if err != nil {
			return err
		}
		result["id"] = *flowLog.Id
		flowlogs = append(flowlogs, result)
		if err := d.Set("flowlog", flowlogs); err != nil {
			return fmt.Errorf("error setting flowlog %w", err)
		}
	}
	return nil
}
