package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	bundleclient "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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

func resourceLoadbalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	rawIds := d.Get("nic_ids").([]interface{})
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
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("error occurred while creating a loadbalancer %s", err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}

	d.SetId(*resp.Id)

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			d.SetId("")
		}
		return utils.ToDiags(d, errState.Error(), &utils.DiagsOpts{Timeout: schema.TimeoutCreate})
	}

	return resourceLoadbalancerRead(ctx, d, meta)
}

func resourceLoadbalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	lb, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersFindById(ctx, d.Get("datacenter_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching a lan: %s", err), nil)
	}

	if lb.Properties.Name != nil {
		if err := d.Set("name", *lb.Properties.Name); err != nil {
			return utils.ToDiags(d, "", nil)
		}
	}

	if lb.Properties.Ip != nil {
		if err := d.Set("ip", *lb.Properties.Ip); err != nil {
			return utils.ToDiags(d, "", nil)
		}
	}

	if lb.Properties.Dhcp != nil {
		if err := d.Set("dhcp", *lb.Properties.Dhcp); err != nil {
			return utils.ToDiags(d, "", nil)
		}
	}

	return nil
}

func resourceLoadbalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

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
			requestLocation, _ := apiResponse.Location()
			return utils.ToDiags(d, fmt.Sprintf("error while updating loadbalancer: %s", err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
		}
	}

	if d.HasChange("nic_ids") {
		oldNicIds, newNicIds := d.GetChange("nic_ids")

		oldList := oldNicIds.([]interface{})

		for _, o := range oldList {
			apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsDelete(context.TODO(),
				d.Get("datacenter_id").(string), d.Id(), o.(string)).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				if httpNotFound(apiResponse) {
					/* 404 - nic was not found - in case the nic is removed, VDC removes the nic from load balancers
					that contain it, behind the scenes - therefore our call will yield 404 */
					log.Printf("[WARN] nic ID %s already removed from load balancer %s\n", o.(string), d.Id())
				} else {
					return utils.ToDiags(d, fmt.Sprintf("[load balancer update] an error occurred while deleting a balanced nic: %s", err), nil)
				}
			} else {
				if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
					return utils.ToDiags(d, errState.Error(), &utils.DiagsOpts{Timeout: schema.TimeoutUpdate})
				}
			}
		}

		newList := newNicIds.([]interface{})

		for _, o := range newList {
			id := o.(string)
			nic := ionoscloud.Nic{Id: &id}
			_, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsPost(ctx, d.Get("datacenter_id").(string), d.Id()).Nic(nic).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				requestLocation, _ := apiResponse.Location()
				return utils.ToDiags(d, fmt.Sprintf("[load balancer update] an error occurred while creating a balanced nic: %s", err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
			}
			if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
				return utils.ToDiags(d, errState.Error(), &utils.DiagsOpts{Timeout: schema.TimeoutUpdate})
			}
		}

	}

	return resourceLoadbalancerRead(ctx, d, meta)
}

func resourceLoadbalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	dcid := d.Get("datacenter_id").(string)
	apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersDelete(ctx, dcid, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("[load balancer delete] an error occurred while deleting a loadbalancer: %s", err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		return utils.ToDiags(d, errState.Error(), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}

	d.SetId("")
	return nil
}

func resourceLoadbalancerImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, utils.ToError(d, "invalid import. Expecting {datacenter}/{loadbalancer}", nil)
	}

	dcId := parts[0]
	lbId := parts[1]

	client := meta.(bundleclient.SdkBundle).CloudApiClient

	loadbalancer, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersFindById(ctx, dcId, lbId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("loadbalancer does not exist %q", lbId), nil)
		}
		return nil, utils.ToError(d, fmt.Sprintf("an error occurred while trying to fetch the loadbalancer %q, error:%s", lbId, err), nil)
	}

	log.Printf("[INFO] loadbalancer found: %+v", loadbalancer)

	d.SetId(*loadbalancer.Id)

	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	if loadbalancer.Properties.Name != nil {
		if err := d.Set("name", *loadbalancer.Properties.Name); err != nil {
			return nil, utils.ToError(d, err.Error(), nil)
		}
	}

	if loadbalancer.Properties.Ip != nil {
		if err := d.Set("ip", *loadbalancer.Properties.Ip); err != nil {
			return nil, utils.ToError(d, err.Error(), nil)
		}
	}

	if loadbalancer.Properties.Dhcp != nil {
		if err := d.Set("dhcp", *loadbalancer.Properties.Dhcp); err != nil {
			return nil, utils.ToError(d, err.Error(), nil)
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
			return nil, utils.ToError(d, err.Error(), nil)
		}
	}

	return []*schema.ResourceData{d}, nil
}
