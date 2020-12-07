package ionoscloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
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
			"nat": {
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	nic := &profitbricks.Nic{
		Properties: &profitbricks.NicProperties{
			Lan: d.Get("lan").(int),
		},
	}
	if _, ok := d.GetOk("name"); ok {
		nic.Properties.Name = d.Get("name").(string)
	}
	if _, ok := d.GetOkExists("dhcp"); ok {
		val := d.Get("dhcp").(bool)
		nic.Properties.Dhcp = &val
	}

	if _, ok := d.GetOk("ip"); ok {
		raw := d.Get("ip").(string)
		ips := strings.Split(raw, ",")
		nic.Properties.Ips = ips
	}
	if _, ok := d.GetOk("firewall_active"); ok {
		raw := d.Get("firewall_active").(bool)
		nic.Properties.FirewallActive = &raw
	}
	if _, ok := d.GetOk("nat"); ok {
		raw := d.Get("nat").(bool)
		nic.Properties.Nat = &raw
	}

	nic, err := client.CreateNic(d.Get("datacenter_id").(string), d.Get("server_id").(string), *nic)
	if err != nil {
		return fmt.Errorf("Error occured while creating a nic: %s", err)
	}
	d.SetId(nic.ID)
	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, nic.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
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
	client := meta.(*profitbricks.Client)
	nic, err := client.GetNic(d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Id())
	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error occured while fetching a nic ID %s %s", d.Id(), err)
	}
	if nic.Properties != nil {
		log.Printf("[INFO] LAN ON NIC: %d", nic.Properties.Lan)
		d.Set("dhcp", nic.Properties.Dhcp)
		d.Set("lan", nic.Properties.Lan)
		d.Set("name", nic.Properties.Name)
		d.Set("ips", nic.Properties.Ips)
		d.Set("firewall_active", nic.Properties.FirewallActive)
	}

	return nil
}

func resourceNicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	properties := profitbricks.NicProperties{}

	if d.HasChange("name") {
		_, n := d.GetChange("name")

		properties.Name = n.(string)
	}
	if d.HasChange("lan") {
		_, n := d.GetChange("lan")
		properties.Lan = n.(int)
	}
	n := d.Get("dhcp").(bool)
	properties.Dhcp = &n

	if d.HasChange("ip") {
		_, raw := d.GetChange("ip")
		ips := strings.Split(raw.(string), ",")
		properties.Ips = ips
	}
	if d.HasChange("nat") {
		_, raw := d.GetChange("nat")
		nat := raw.(bool)
		properties.Nat = &nat
	}

	nic, err := client.UpdateNic(d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Id(), properties)

	if err != nil {
		return fmt.Errorf("Error occured while updating a nic: %s", err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, nic.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceNicRead(d, meta)
}

func resourceNicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	resp, err := client.DeleteNic(d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("An error occured while deleting a nic dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err)
	}
	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
