package ionoscloud

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourceFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallCreate,
		Read:   resourceFirewallRead,
		Update: resourceFirewallUpdate,
		Delete: resourceFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: resourceFirewallImport,
		},
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_mac": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_range_start": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if v.(int) < 1 && v.(int) > 65534 {
						errors = append(errors, fmt.Errorf("Port start range must be between 1 and 65534"))
					}
					return
				},
			},

			"port_range_end": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if v.(int) < 1 && v.(int) > 65534 {
						errors = append(errors, fmt.Errorf("Port end range must be between 1 and 65534"))
					}
					return
				},
			},
			"icmp_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"icmp_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nic_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient
	fw := &profitbricks.FirewallRule{
		Properties: profitbricks.FirewallruleProperties{
			Protocol: d.Get("protocol").(string),
		},
	}

	if _, ok := d.GetOk("name"); ok {
		fw.Properties.Name = d.Get("name").(string)
	}
	if _, ok := d.GetOk("source_mac"); ok {
		tempSourceMac := d.Get("source_mac").(string)
		fw.Properties.SourceMac = &tempSourceMac
	}
	if _, ok := d.GetOk("source_ip"); ok {
		tempSourceIp := d.Get("source_ip").(string)
		fw.Properties.SourceIP = &tempSourceIp
	}
	if _, ok := d.GetOk("target_ip"); ok {
		tempTargetIp := d.Get("target_ip").(string)
		fw.Properties.TargetIP = &tempTargetIp
	}
	if _, ok := d.GetOk("port_range_start"); ok {
		tempPortRangeStart := d.Get("port_range_start").(int)
		fw.Properties.PortRangeStart = &tempPortRangeStart
	}
	if _, ok := d.GetOk("port_range_end"); ok {
		tempPortRangeEnd := d.Get("port_range_end").(int)
		fw.Properties.PortRangeEnd = &tempPortRangeEnd
	}
	if _, ok := d.GetOk("icmp_type"); ok {
		tempIcmpType, err := strconv.Atoi(d.Get("icmp_type").(string))
		if err != nil {
			return fmt.Errorf("An error occured while creating a firewall rule: %s", err)
		}
		fw.Properties.IcmpType = &tempIcmpType
	}
	if _, ok := d.GetOk("icmp_code"); ok {
		tempIcmpCodee, err := strconv.Atoi(d.Get("icmp_code").(string))
		if err != nil {
			return fmt.Errorf("An error occured while creating a firewall rule: %s", err)
		}
		fw.Properties.IcmpCode = &tempIcmpCodee
	}

	fw, err := client.CreateFirewallRule(d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string), *fw)

	if err != nil {
		return fmt.Errorf("An error occured while creating a firewall rule: %s", err)
	}
	d.SetId(fw.ID)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, fw.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceFirewallRead(d, meta)
}

func resourceFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient
	fw, err := client.GetFirewallRule(d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("An error occured while fetching a firewall rule  dcId: %s server_id: %s  nic_id: %s ID: %s %s", d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string), d.Id(), err)
	}

	d.Set("protocol", fw.Properties.Protocol)
	d.Set("name", fw.Properties.Name)
	d.Set("source_mac", fw.Properties.SourceMac)
	d.Set("source_ip", fw.Properties.SourceIP)
	d.Set("target_ip", fw.Properties.TargetIP)
	d.Set("port_range_start", fw.Properties.PortRangeStart)
	d.Set("port_range_end", fw.Properties.PortRangeEnd)
	d.Set("icmp_type", fw.Properties.IcmpType)
	d.Set("icmp_code", fw.Properties.IcmpCode)
	d.Set("nic_id", d.Get("nic_id").(string))

	return nil
}

func resourceFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient
	properties := profitbricks.FirewallruleProperties{}

	if d.HasChange("name") {
		_, new := d.GetChange("name")
		properties.Name = new.(string)
	}
	if d.HasChange("source_mac") {
		_, new := d.GetChange("source_mac")
		properties.SourceMac = new.(*string)
	}
	if d.HasChange("source_ip") {
		_, new := d.GetChange("source_ip")
		properties.SourceIP = new.(*string)
	}
	if d.HasChange("target_ip") {
		_, new := d.GetChange("target_ip")
		properties.TargetIP = new.(*string)
	}
	if d.HasChange("port_range_start") {
		_, new := d.GetChange("port_range_start")
		properties.PortRangeStart = new.(*int)
	}
	if d.HasChange("port_range_end") {
		_, new := d.GetChange("port_range_end")
		properties.PortRangeEnd = new.(*int)
	}
	if d.HasChange("icmp_type") {
		_, new := d.GetChange("icmp_type")
		tempIcmpType, err := strconv.Atoi(new.(string))
		if err != nil {
			return fmt.Errorf("An error occured while updating a firewall rule: %s", err)
		}
		properties.IcmpType = &tempIcmpType
	}
	if d.HasChange("icmp_code") {
		_, new := d.GetChange("icmp_code")
		tempIcmpCode, err := strconv.Atoi(new.(string))
		if err != nil {
			return fmt.Errorf("An error occured while updating a firewall rule: %s", err)
		}
		properties.IcmpCode = &tempIcmpCode
	}

	fw, err := client.UpdateFirewallRule(d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string), d.Id(), properties)

	if err != nil {
		return fmt.Errorf("An error occured while updating a firewall rule ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, fw.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceFirewallRead(d, meta)
}

func resourceFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient
	resp, err := client.DeleteFirewallRule(d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("An error occured while deleting a firewall rule ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")

	return nil
}
