package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"strings"
)

func resourceFirewall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallCreate,
		ReadContext:   resourceFirewallRead,
		UpdateContext: resourceFirewallUpdate,
		DeleteContext: resourceFirewallDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceFirewallImport,
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
				Computed: true,
			},
			"target_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port_range_start": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if v.(int) < 1 && v.(int) > 65534 {
						errors = append(errors, fmt.Errorf("port start range must be between 1 and 65534"))
					}
					return
				},
			},

			"port_range_end": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if v.(int) < 1 && v.(int) > 65534 {
						errors = append(errors, fmt.Errorf("port end range must be between 1 and 65534"))
					}
					return
				},
			},
			"icmp_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"icmp_code": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"datacenter_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"nic_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceFirewallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	firewallProtocol := d.Get("protocol").(string)

	fw := ionoscloud.FirewallRule{
		Properties: &ionoscloud.FirewallruleProperties{
			Protocol: &firewallProtocol,
		},
	}

	if _, ok := d.GetOk("name"); ok {
		firewallName := d.Get("name").(string)
		fw.Properties.Name = &firewallName
	}
	if _, ok := d.GetOk("source_mac"); ok {
		tempSourceMac := d.Get("source_mac").(string)
		fw.Properties.SourceMac = &tempSourceMac
	}
	if _, ok := d.GetOk("source_ip"); ok {
		tempSourceIp := d.Get("source_ip").(string)
		fw.Properties.SourceIp = &tempSourceIp
	}
	if _, ok := d.GetOk("target_ip"); ok {
		tempTargetIp := d.Get("target_ip").(string)
		fw.Properties.TargetIp = &tempTargetIp
	}
	if _, ok := d.GetOk("port_range_start"); ok {
		tempPortRangeStart := int32(d.Get("port_range_start").(int))
		fw.Properties.PortRangeStart = &tempPortRangeStart
	}
	if _, ok := d.GetOk("port_range_end"); ok {
		tempPortRangeEnd := int32(d.Get("port_range_end").(int))
		fw.Properties.PortRangeEnd = &tempPortRangeEnd
	}
	if _, ok := d.GetOk("icmp_type"); ok {
		fwIcmpType := int32(d.Get("icmp_type").(int))
		fw.Properties.IcmpType = &fwIcmpType
	}
	if _, ok := d.GetOk("icmp_code"); ok {
		fwIcmpTypeCode := int32(d.Get("icmp_code").(int))
		fw.Properties.IcmpCode = &fwIcmpTypeCode
	}
	if _, ok := d.GetOk("type"); ok {
		fwType := d.Get("type").(string)
		fw.Properties.Type = &fwType
	}

	fw, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPost(ctx, d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string)).Firewallrule(fw).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a firewall rule: %s", err))
		return diags
	}
	d.SetId(*fw.Id)

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

	return resourceFirewallRead(ctx, d, meta)
}

func resourceFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	fw, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, d.Get("datacenter_id").(string),
		d.Get("server_id").(string), d.Get("nic_id").(string), d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a firewall rule  dcId: %s server_id: %s  nic_id: %s ID: %s %s", d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string), d.Id(), err))
		return diags
	}

	if err := setFirewallData(d, &fw); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	properties := ionoscloud.FirewallruleProperties{}

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		vStr := v.(string)
		properties.Name = &vStr
	}
	if d.HasChange("source_mac") {
		_, v := d.GetChange("source_mac")
		vStr := v.(string)
		properties.SourceMac = &vStr
	}
	if d.HasChange("source_ip") {
		_, v := d.GetChange("source_ip")
		vStr := v.(string)
		properties.SourceIp = &vStr
	}
	if d.HasChange("target_ip") {
		_, v := d.GetChange("target_ip")
		vStr := v.(string)
		properties.TargetIp = &vStr
	}
	if d.HasChange("port_range_start") {
		vInt := int32(d.Get("port_range_start").(int))
		properties.PortRangeStart = &vInt
	}
	if d.HasChange("port_range_end") {
		vInt := int32(d.Get("port_range_end").(int))
		properties.PortRangeEnd = &vInt
	}
	if d.HasChange("icmp_type") {
		vInt := int32(d.Get("icmp_type").(int))
		properties.IcmpType = &vInt
	}
	if d.HasChange("icmp_code") {
		vInt := int32(d.Get("icmp_code").(int))
		properties.IcmpCode = &vInt
	}

	if d.HasChange("type") {
		_, v := d.GetChange("type")
		vStr := v.(string)
		properties.Type = &vStr
	}

	_, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPatch(ctx, d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string), d.Id()).Firewallrule(properties).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a firewall rule ID %s %s", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceFirewallRead(ctx, d, meta)
}

func resourceFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	apiResponse, err := client.FirewallRulesApi.
		DatacentersServersNicsFirewallrulesDelete(
			ctx, d.Get("datacenter_id").(string),
			d.Get("server_id").(string), d.Get("nic_id").(string), d.Id()).
		Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a firewall rule ID %s %s", d.Id(), err))
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

func resourceFirewallImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	client := meta.(*ionoscloud.APIClient)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}/{nic}/{firewall}", d.Id())
	}

	dcId := parts[0]
	serverId := parts[1]
	nicId := parts[2]
	firewallId := parts[3]

	fw, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, dcId,
		serverId, nicId, firewallId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find firewall rule %q", firewallId)
		}
		return nil, fmt.Errorf("an error occured while retrieving firewall rule %q: %q ", firewallId, err)
	}

	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, err
	}
	if err := d.Set("server_id", serverId); err != nil {
		return nil, err
	}
	if err := d.Set("nic_id", nicId); err != nil {
		return nil, err
	}

	if err := setFirewallData(d, &fw); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func setFirewallData(d *schema.ResourceData, firewall *ionoscloud.FirewallRule) error {

	if firewall.Id != nil {
		d.SetId(*firewall.Id)
	}

	if firewall.Properties != nil {

		if firewall.Properties.Protocol != nil {
			err := d.Set("protocol", *firewall.Properties.Protocol)
			if err != nil {
				return fmt.Errorf("error while setting protocol property for firewall %s: %s", d.Id(), err)
			}
		}

		if firewall.Properties.Name != nil {
			err := d.Set("name", *firewall.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for firewall %s: %s", d.Id(), err)
			}
		}

		if firewall.Properties.SourceMac != nil {
			err := d.Set("source_mac", *firewall.Properties.SourceMac)
			if err != nil {
				return fmt.Errorf("error while setting source_mac property for firewall %s: %s", d.Id(), err)
			}
		}

		if firewall.Properties.SourceIp != nil {
			err := d.Set("source_ip", *firewall.Properties.SourceIp)
			if err != nil {
				return fmt.Errorf("error while setting source_ip property for firewall %s: %s", d.Id(), err)
			}
		}

		if firewall.Properties.TargetIp != nil {
			err := d.Set("target_ip", *firewall.Properties.TargetIp)
			if err != nil {
				return fmt.Errorf("error while setting target_ip property for firewall %s: %s", d.Id(), err)
			}
		}

		if firewall.Properties.PortRangeStart != nil {
			err := d.Set("port_range_start", *firewall.Properties.PortRangeStart)
			if err != nil {
				return fmt.Errorf("error while setting port_range_start property for firewall %s: %s", d.Id(), err)
			}
		}

		if firewall.Properties.PortRangeEnd != nil {
			err := d.Set("port_range_end", *firewall.Properties.PortRangeEnd)
			if err != nil {
				return fmt.Errorf("error while setting port_range_end property for firewall %s: %s", d.Id(), err)
			}
		}

		if firewall.Properties.IcmpType != nil {
			err := d.Set("icmp_type", *firewall.Properties.IcmpType)
			if err != nil {
				return fmt.Errorf("error while setting icmp_type property for firewall %s: %s", d.Id(), err)
			}
		}

		if firewall.Properties.IcmpCode != nil {
			err := d.Set("icmp_code", *firewall.Properties.IcmpCode)
			if err != nil {
				return fmt.Errorf("error while setting icmp_code property for firewall %s: %s", d.Id(), err)
			}
		}

		if firewall.Properties.Type != nil {
			err := d.Set("type", *firewall.Properties.Type)
			if err != nil {
				return fmt.Errorf("error while setting type property for firewall %s: %s", d.Id(), err)
			}
		}
	}
	return nil
}
