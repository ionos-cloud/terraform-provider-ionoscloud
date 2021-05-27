package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func resourceLoadbalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadbalancerCreate,
		Read:   resourceLoadbalancerRead,
		Update: resourceLoadbalancerUpdate,
		Delete: resourceLoadbalancerDelete,
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

func resourceLoadbalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	raw_ids := d.Get("nic_ids").([]interface{})
	nic_ids := []ionoscloud.Nic{}

	for _, id := range raw_ids {
		id := id.(string)
		nic_ids = append(nic_ids, ionoscloud.Nic{Id: &id})
	}

	name := d.Get("name").(string)
	lb := &ionoscloud.Loadbalancer{
		Properties: &ionoscloud.LoadbalancerProperties{
			Name: &name,
		},
		Entities: &ionoscloud.LoadbalancerEntities{
			Balancednics: &ionoscloud.BalancedNics{
				Items: &nic_ids,
			},
		},
	}

	dcid := d.Get("datacenter_id").(string)
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)
	if cancel != nil {
		defer cancel()
	}

	resp, apiResp, err := client.LoadBalancerApi.DatacentersLoadbalancersPost(ctx, dcid).Loadbalancer(*lb).Execute()
	if err != nil {
		return fmt.Errorf("Error occured while creating a loadbalancer %s", err)
	}

	d.SetId(*resp.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResp.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceLoadbalancerRead(d, meta)
}

func resourceLoadbalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}

	lb, apiResponse, err := client.LoadBalancerApi.DatacentersLoadbalancersFindById(ctx, d.Get("datacenter_id").(string), d.Id()).Execute()
	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("An error occured while fetching a lan ID %s %s", d.Id(), err)
	}

	d.Set("name", *lb.Properties.Name)
	d.Set("ip", *lb.Properties.Ip)
	d.Set("dhcp", *lb.Properties.Dhcp)

	return nil
}

func resourceLoadbalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	properties := &ionoscloud.LoadbalancerProperties{}

	var hasChangeCount int = 0
	if d.HasChange("name") {
		_, new := d.GetChange("name")
		name := new.(string)
		properties.Name = &name
		hasChangeCount++
	}
	if d.HasChange("ip") {
		_, new := d.GetChange("ip")
		ip := new.(string)
		properties.Ip = &ip
		hasChangeCount++
	}
	if d.HasChange("dhcp") {
		_, new := d.GetChange("dhcp")
		dhcp := new.(bool)
		properties.Dhcp = &dhcp
		hasChangeCount++
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)
	if cancel != nil {
		defer cancel()
	}
	if hasChangeCount > 0 {
		_, _, err := client.LoadBalancerApi.DatacentersLoadbalancersPatch(ctx, d.Get("datacenter_id").(string), d.Id()).Loadbalancer(*properties).Execute()
		if err != nil {
			return fmt.Errorf("error while updating loadbalancer %s: %s ", d.Id(), err)
		}
	}

	if d.HasChange("nic_ids") {
		oldNicIds, newNicIds := d.GetChange("nic_ids")

		oldList := oldNicIds.([]interface{})

		for _, o := range oldList {
			_, apiresponse, err := client.LoadBalancerApi.DatacentersLoadbalancersBalancednicsDelete(context.TODO(),
				d.Get("datacenter_id").(string), d.Id(), o.(string)).Execute()
			if err != nil {
				if apiresponse != nil && apiresponse.Response.StatusCode == 404 {
					/* 404 - nic was not found - in case the nic is removed, VDC removes the nic from load balancers
					that contain it, behind the scenes - therefore our call will yield 404 */
					fmt.Printf("[WARNING] nic ID %s already removed from load balancer %s\n", o.(string), d.Id())
				} else {
					return fmt.Errorf("[load balancer update] an error occured while deleting a balanced nic: %s", err)
				}
			} else {
				// Wait, catching any errors
				_, errState := getStateChangeConf(meta, d, apiresponse.Header.Get("location"), schema.TimeoutUpdate).WaitForState()
				if errState != nil {
					return errState
				}
			}
		}

		newList := newNicIds.([]interface{})

		for _, o := range newList {
			id := o.(string)
			nic := ionoscloud.Nic{Id: &id}
			_, apiResp, err := client.LoadBalancerApi.DatacentersLoadbalancersBalancednicsPost(ctx, d.Get("datacenter_id").(string), d.Id()).Nic(nic).Execute()
			if err != nil {
				return fmt.Errorf("[load balancer update] an error occured while creating a balanced nic: %s", err)
			}
			// Wait, catching any errors
			_, errState := getStateChangeConf(meta, d, apiResp.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
			if errState != nil {
				return errState
			}

		}

	}

	return resourceLoadbalancerRead(d, meta)
}

func resourceLoadbalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	dcid := d.Get("datacenter_id").(string)
	_, apiResp, err := client.LoadBalancerApi.DatacentersLoadbalancersDelete(ctx, dcid, d.Id()).Execute()

	if err != nil {
		return fmt.Errorf("[load balancer delete] an error occured while deleting a loadbalancer: %s", err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResp.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
