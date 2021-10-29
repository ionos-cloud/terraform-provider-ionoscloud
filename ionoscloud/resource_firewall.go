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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
	client := meta.(*ionoscloud.APIClient)

	firewall := getFirewallData(d, "", false)

	fw, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPost(ctx, d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string)).Firewallrule(firewall).Execute()

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
	client := meta.(*ionoscloud.APIClient)

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
	client := meta.(*ionoscloud.APIClient)

	firewall := getFirewallData(d, "", true)

	_, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPatch(ctx, d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string), d.Id()).Firewallrule(*firewall.Properties).Execute()

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
	client := meta.(*ionoscloud.APIClient)

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

func getFirewallData(d *schema.ResourceData, path string, update bool) ionoscloud.FirewallRule {

	firewall := ionoscloud.FirewallRule{
		Properties: &ionoscloud.FirewallruleProperties{},
	}

	if !update {
		if v, ok := d.GetOk(path + "protocol"); ok {
			vStr := v.(string)
			firewall.Properties.Protocol = &vStr
		}
	}

	if v, ok := d.GetOk(path + "name"); ok {
		vStr := v.(string)
		firewall.Properties.Name = &vStr
	}

	if v, ok := d.GetOk(path + "source_mac"); ok {
		val := v.(string)
		firewall.Properties.SourceMac = &val
	}

	if v, ok := d.GetOk(path + "source_ip"); ok {
		val := v.(string)
		firewall.Properties.SourceIp = &val
	}

	if v, ok := d.GetOk(path + "target_ip"); ok {
		val := v.(string)
		firewall.Properties.TargetIp = &val
	}

	if v, ok := d.GetOk(path + "port_range_start"); ok {
		val := int32(v.(int))
		firewall.Properties.PortRangeStart = &val
	}

	if v, ok := d.GetOk(path + "port_range_end"); ok {
		val := int32(v.(int))
		firewall.Properties.PortRangeEnd = &val
	}

	if v, ok := d.GetOk(path + "icmp_type"); ok {
		tempIcmpType := int32(v.(int))
		firewall.Properties.IcmpType = &tempIcmpType

	}
	if v, ok := d.GetOk(path + "icmp_code"); ok {
		tempIcmpCode := int32(v.(int))
		firewall.Properties.IcmpCode = &tempIcmpCode

	}
	if v, ok := d.GetOk(path + "type"); ok {
		tempType := v.(string)
		firewall.Properties.Type = &tempType

	}
	return firewall
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
