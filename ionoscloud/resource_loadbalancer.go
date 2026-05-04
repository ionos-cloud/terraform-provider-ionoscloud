package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func resourceLoadbalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadbalancerCreate,
		ReadContext:   resourceLoadbalancerRead,
		UpdateContext: resourceLoadbalancerUpdate,
		DeleteContext: resourceLoadbalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLoadbalancerImporter,
		},
		Schema: map[string]*schema.Schema{

			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},

			"ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				Optional:    true,
				ForceNew:    true,
			},
			"nic_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace)},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceLoadbalancerCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	rawIds := d.Get("nic_ids").([]any)
	var nicIds []ionoscloud.Nic

	for _, id := range rawIds {
		id := id.(string)
		nicIds = append(nicIds, ionoscloud.Nic{Id: &id})
	}

	name := d.Get("name").(string)
	lb := &ionoscloud.Loadbalancer{
		Properties: &ionoscloud.LoadbalancerProperties{
			Name: &name,
		},
		Entities: &ionoscloud.LoadbalancerEntities{
			Balancednics: &ionoscloud.BalancedNics{
				Items: &nicIds,
			},
		},
	}

	dcid := d.Get("datacenter_id").(string)

	resp, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersPost(ctx, dcid).Loadbalancer(*lb).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("error occurred while creating a loadbalancer %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	d.SetId(*resp.Id)

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			d.SetId("")
		}
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceLoadbalancerRead(ctx, d, meta)
}

func resourceLoadbalancerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	lb, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersFindById(ctx, d.Get("datacenter_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching a lan: %w", err), nil)
	}

	if lb.Properties.Name != nil {
		if err := d.Set("name", *lb.Properties.Name); err != nil {
			return diagutil.ToDiags(d, fmt.Errorf(""), nil)
		}
	}

	if lb.Properties.Ip != nil {
		if err := d.Set("ip", *lb.Properties.Ip); err != nil {
			return diagutil.ToDiags(d, fmt.Errorf(""), nil)
		}
	}

	if lb.Properties.Dhcp != nil {
		if err := d.Set("dhcp", *lb.Properties.Dhcp); err != nil {
			return diagutil.ToDiags(d, fmt.Errorf(""), nil)
		}
	}

	return nil
}

func resourceLoadbalancerUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	properties := &ionoscloud.LoadbalancerProperties{}

	var hasChangeCount = 0
	if d.HasChange("name") {
		_, newVal := d.GetChange("name")
		name := newVal.(string)
		properties.Name = &name
		hasChangeCount++
	}
	if d.HasChange("ip") {
		_, newVal := d.GetChange("ip")
		ip := newVal.(string)
		properties.Ip = &ip
		hasChangeCount++
	}
	if d.HasChange("dhcp") {
		_, newVal := d.GetChange("dhcp")
		dhcp := newVal.(bool)
		properties.Dhcp = &dhcp
		hasChangeCount++
	}

	if hasChangeCount > 0 {
		_, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersPatch(ctx, d.Get("datacenter_id").(string), d.Id()).Loadbalancer(*properties).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			requestLocation, _ := apiResponse.SafeLocation()
			return diagutil.ToDiags(d, fmt.Errorf("error while updating loadbalancer: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
		}
	}

	if d.HasChange("nic_ids") {
		oldNicIds, newNicIds := d.GetChange("nic_ids")

		oldList := oldNicIds.([]any)

		for _, o := range oldList {
			apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsDelete(context.TODO(),
				d.Get("datacenter_id").(string), d.Id(), o.(string)).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				if httpNotFound(apiResponse) {
					/* 404 - nic was not found - in case the nic is removed, VDC removes the nic from load balancers
					that contain it, behind the scenes - therefore our call will yield 404 */
					tflog.Warn(ctx, "nic already removed from load balancer", map[string]interface{}{"nic_id": o.(string), "loadbalancer_id": d.Id()})
				} else {
					return diagutil.ToDiags(d, fmt.Errorf("[load balancer update] an error occurred while deleting a balanced nic: %w", err), nil)
				}
			} else {
				if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
					requestLocation, _ := apiResponse.SafeLocation()
					return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
				}
			}
		}

		newList := newNicIds.([]any)

		for _, o := range newList {
			id := o.(string)
			nic := ionoscloud.Nic{Id: &id}
			_, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsPost(ctx, d.Get("datacenter_id").(string), d.Id()).Nic(nic).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				requestLocation, _ := apiResponse.SafeLocation()
				return diagutil.ToDiags(d, fmt.Errorf("[load balancer update] an error occurred while creating a balanced nic: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
			}
			if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
				requestLocation, _ := apiResponse.SafeLocation()
				return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
			}
		}

	}

	return resourceLoadbalancerRead(ctx, d, meta)
}

func resourceLoadbalancerDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	dcid := d.Get("datacenter_id").(string)
	apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersDelete(ctx, dcid, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("[load balancer delete] an error occurred while deleting a loadbalancer: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	d.SetId("")
	return nil
}

func resourceLoadbalancerImporter(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	importID := d.Id()

	location, parts := splitImportID(importID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf(
			"invalid import identifier: expected one of <location>:<datacenter-id>/<loadbalancer-id> "+
				"or <datacenter-id>/<loadbalancer-id>, got: %s", importID,
		)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("failed validating import identifier %q: %w", importID, err), nil)
	}

	dcId := parts[0]
	lbId := parts[1]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return nil, err
	}

	loadbalancer, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersFindById(ctx, dcId, lbId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("loadbalancer does not exist %q", lbId), nil)
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while trying to fetch the loadbalancer %q, error:%w", lbId, err), nil)
	}

	tflog.Info(ctx, "loadbalancer found", map[string]interface{}{"loadbalancer_id": *loadbalancer.Id, "datacenter_id": dcId})

	d.SetId(*loadbalancer.Id)

	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	if err := d.Set("location", location); err != nil {
		return nil, err
	}

	if loadbalancer.Properties.Name != nil {
		if err := d.Set("name", *loadbalancer.Properties.Name); err != nil {
			return nil, diagutil.ToError(d, err, nil)
		}
	}

	if loadbalancer.Properties.Ip != nil {
		if err := d.Set("ip", *loadbalancer.Properties.Ip); err != nil {
			return nil, diagutil.ToError(d, err, nil)
		}
	}

	if loadbalancer.Properties.Dhcp != nil {
		if err := d.Set("dhcp", *loadbalancer.Properties.Dhcp); err != nil {
			return nil, diagutil.ToError(d, err, nil)
		}
	}

	if loadbalancer.Entities != nil && loadbalancer.Entities.Balancednics != nil &&
		loadbalancer.Entities.Balancednics.Items != nil && len(*loadbalancer.Entities.Balancednics.Items) > 0 {

		var lans []string
		for _, lan := range *loadbalancer.Entities.Balancednics.Items {
			if *lan.Id != "" {
				lans = append(lans, *lan.Id)
			}
		}
		if err := d.Set("nic_ids", lans); err != nil {
			return nil, diagutil.ToError(d, err, nil)
		}
	}

	return []*schema.ResourceData{d}, nil
}
