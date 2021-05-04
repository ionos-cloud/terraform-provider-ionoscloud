package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNic() *schema.Resource {
	return &schema.Resource{
		Create: resourceNicCreate,
		Read:   resourceNicRead,
		Update: resourceNicUpdate,
		Delete: resourceNicDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNicImport,
		},
		Schema: map[string]*schema.Schema{

			"lan": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dhcp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ips": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"firewall_active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	lan := d.Get("lan").(int)
	lanConverted := int32(lan)

	nic := ionoscloud.Nic{
		Properties: &ionoscloud.NicProperties{
			Lan: &lanConverted,
		},
	}
	if _, ok := d.GetOk("name"); ok {
		name := d.Get("name").(string)
		nic.Properties.Name = &name
	}
	if _, ok := d.GetOkExists("dhcp"); ok {
		val := d.Get("dhcp").(bool)
		nic.Properties.Dhcp = &val
	}
	if _, ok := d.GetOk("ip"); ok {
		raw := d.Get("ip").(string)
		ips := strings.Split(raw, ",")
		nic.Properties.Ips = &ips
	}
	if _, ok := d.GetOk("firewall_active"); ok {
		raw := d.Get("firewall_active").(bool)
		nic.Properties.FirewallActive = &raw
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)
	if cancel != nil {
		defer cancel()
	}
	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nic, apiResp, err := client.NetworkInterfacesApi.DatacentersServersNicsPost(ctx, dcid, srvid).Nic(nic).Execute()

	if err != nil {
		return fmt.Errorf("Error occured while creating a nic: %s", err)
	}
	if nic.Id != nil {
		d.SetId(*nic.Id)
	}
	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResp.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}
	return resourceNicRead(d, meta)
}

func resourceNicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()

	rsp, apiresponse, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, dcid, srvid, nicid).Execute()
	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiresponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error occured while fetching a nic ID %s %s", d.Id(), err)
	}

	if rsp.Properties != nil {
		log.Printf("[INFO] LAN ON NIC: %d", rsp.Properties.Lan)
		if rsp.Properties.Dhcp != nil {
			d.Set("dhcp", *rsp.Properties.Dhcp)
		}
		if rsp.Properties.Lan != nil {
			d.Set("lan", *rsp.Properties.Lan)
		}
		if rsp.Properties.Name != nil {
			d.Set("name", *rsp.Properties.Name)
		}
		if rsp.Properties.Ips != nil {
			d.Set("ips", *rsp.Properties.Ips)
		}
		if rsp.Properties.FirewallActive != nil {
			d.Set("firewall_active", *rsp.Properties.FirewallActive)
		}
		if rsp.Properties.Mac != nil {
			d.Set("mac", *rsp.Properties.Mac)
		}
	}

	return nil
}

func resourceNicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client
	properties := ionoscloud.NicProperties{}

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		name := n.(string)
		properties.Name = &name
	}
	if d.HasChange("lan") {
		_, n := d.GetChange("lan")
		lan := n.(int32)
		properties.Lan = &lan
	}

	n := d.Get("dhcp").(bool)
	properties.Dhcp = &n

	if d.HasChange("ip") {
		_, raw := d.GetChange("ip")
		ips := strings.Split(raw.(string), ",")
		properties.Ips = &ips
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)
	if cancel != nil {
		defer cancel()
	}
	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()

	_, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsPatch(ctx, dcid, srvid, nicid).Nic(properties).Execute()

	if err != nil {
		return fmt.Errorf("Error occured while updating a nic: %s", err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceNicRead(d, meta)
}

func resourceNicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}
	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()
	//resp, err := client.DeleteNic(d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Id())
	_, apiresp, err := client.NetworkInterfacesApi.DatacentersServersNicsDelete(ctx, dcid, srvid, nicid).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while deleting a nic dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err)
	}
	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiresp.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
