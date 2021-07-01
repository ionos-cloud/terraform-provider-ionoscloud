package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nic_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.All(validation.StringIsNotWhiteSpace)},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceLoadbalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

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

	resp, apiResponse, err := client.LoadBalancerApi.DatacentersLoadbalancersPost(ctx, dcid).Loadbalancer(*lb).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while creating a loadbalancer %s", err))
		return diags
	}

	d.SetId(*resp.Id)

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

	return resourceLoadbalancerRead(ctx, d, meta)
}

func resourceLoadbalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	lb, apiResponse, err := client.LoadBalancerApi.DatacentersLoadbalancersFindById(ctx, d.Get("datacenter_id").(string), d.Id()).Execute()
	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a lan ID %s %s", d.Id(), err))
		return diags
	}

	if lb.Properties.Name != nil {
		if err := d.Set("name", *lb.Properties.Name); err != nil {
			diags := diag.FromErr(fmt.Errorf(""))
			return diags
		}
	}

	if lb.Properties.Ip != nil {
		if err := d.Set("ip", *lb.Properties.Ip); err != nil {
			diags := diag.FromErr(fmt.Errorf(""))
			return diags
		}
	}

	if lb.Properties.Dhcp != nil {
		if err := d.Set("dhcp", *lb.Properties.Dhcp); err != nil {
			diags := diag.FromErr(fmt.Errorf(""))
			return diags
		}
	}

	return nil
}

func resourceLoadbalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	properties := &ionoscloud.LoadbalancerProperties{}

	var hasChangeCount = 0
	if d.HasChange("name") {
		_, newValue := d.GetChange("name")
		name := newValue.(string)
		properties.Name = &name
		hasChangeCount++
	}
	if d.HasChange("ip") {
		_, newValue := d.GetChange("ip")
		ip := newValue.(string)
		properties.Ip = &ip
		hasChangeCount++
	}
	if d.HasChange("dhcp") {
		_, newValue := d.GetChange("dhcp")
		dhcp := newValue.(bool)
		properties.Dhcp = &dhcp
		hasChangeCount++
	}

	if hasChangeCount > 0 {
		_, _, err := client.LoadBalancerApi.DatacentersLoadbalancersPatch(ctx, d.Get("datacenter_id").(string), d.Id()).Loadbalancer(*properties).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while updating loadbalancer %s: %s \n", d.Id(), err))
			return diags
		}
	}

	if d.HasChange("nic_ids") {
		oldNicIds, newNicIds := d.GetChange("nic_ids")

		oldList := oldNicIds.([]interface{})

		for _, o := range oldList {
			_, apiResponse, err := client.LoadBalancerApi.DatacentersLoadbalancersBalancednicsDelete(context.TODO(),
				d.Get("datacenter_id").(string), d.Id(), o.(string)).Execute()
			if err != nil {
				if apiResponse != nil && apiResponse.StatusCode == 404 {
					/* 404 - nic was not found - in case the nic is removed, VDC removes the nic from load balancers
					that contain it, behind the scenes - therefore our call will yield 404 */
					fmt.Printf("[WARNING] nic ID %s already removed from load balancer %s\n", o.(string), d.Id())
				} else {
					diags := diag.FromErr(fmt.Errorf("[load balancer update] an error occured while deleting a balanced nic: %s", err))
					return diags
				}
			} else {
				// Wait, catching any errors
				_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
				if errState != nil {
					diags := diag.FromErr(errState)
					return diags
				}
			}
		}

		newList := newNicIds.([]interface{})

		for _, o := range newList {
			id := o.(string)
			nic := ionoscloud.Nic{Id: &id}
			_, apiResponse, err := client.LoadBalancerApi.DatacentersLoadbalancersBalancednicsPost(ctx, d.Get("datacenter_id").(string), d.Id()).Nic(nic).Execute()
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("[load balancer update] an error occured while creating a balanced nic: %s", err))
				return diags
			}
			// Wait, catching any errors
			_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
			if errState != nil {
				diags := diag.FromErr(errState)
				return diags
			}

		}

	}

	return resourceLoadbalancerRead(ctx, d, meta)
}

func resourceLoadbalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcid := d.Get("datacenter_id").(string)
	_, apiResponse, err := client.LoadBalancerApi.DatacentersLoadbalancersDelete(ctx, dcid, d.Id()).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("[load balancer delete] an error occured while deleting a loadbalancer: %s", err))
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
