package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	bundleclient "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(constant.FirewallProtocolEnum, false)),
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
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65534)),
			},

			"port_range_end": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65534)),
			},
			"icmp_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"icmp_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"server_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"nic_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceFirewallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	firewall, diags := getFirewallData(d, "", false)
	if diags != nil {
		return diags
	}
	fw, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPost(ctx, d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string)).Firewallrule(firewall).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating a firewall rule: %s", err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}
	d.SetId(*fw.Id)

	// Wait, catching any errors
	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			log.Printf("[DEBUG] firewall resource failed to be created")
			d.SetId("")
		}
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating a firewall rule dcId: %s server_id: %s  "+
			"nic_id: %s %s", d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string), errState), &utils.DiagsOpts{Timeout: schema.TimeoutCreate})
	}

	return resourceFirewallRead(ctx, d, meta)
}

func resourceFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	fw, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, d.Get("datacenter_id").(string),
		d.Get("server_id").(string), d.Get("nic_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			log.Printf("[DEBUG] could not find firewall rule datacenter_id = %s server_id = %s with id = %s", d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Id())
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching a firewall rule dcId: %s server_id: %s  nic_id: %s ID: %s %s",
			d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string), d.Id(), err), nil)
	}

	if err := setFirewallData(d, &fw); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	firewall, diags := getFirewallData(d, "", true)
	if diags != nil {
		return diags
	}
	_, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPatch(ctx, d.Get("datacenter_id").(string), d.Get("server_id").(string), d.Get("nic_id").(string), d.Id()).Firewallrule(*firewall.Properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating a firewall rule: %s", err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return utils.ToDiags(d, fmt.Sprintf("error getting state change for firewall patch %s", errState), &utils.DiagsOpts{Timeout: schema.TimeoutUpdate})
	}

	return resourceFirewallRead(ctx, d, meta)
}

func resourceFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	apiResponse, err := client.FirewallRulesApi.
		DatacentersServersNicsFirewallrulesDelete(
			ctx, d.Get("datacenter_id").(string),
			d.Get("server_id").(string), d.Get("nic_id").(string), d.Id()).
		Execute()

	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while deleting a firewall rule: %s", err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		return utils.ToDiags(d, fmt.Sprintf("error getting state change for firewall delete %s", errState), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}

	d.SetId("")

	return nil
}

func resourceFirewallImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	client := meta.(bundleclient.SdkBundle).CloudApiClient

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
		return nil, utils.ToError(d, "invalid import. Expecting {datacenter}/{server}/{nic}/{firewall}", nil)
	}

	dcId := parts[0]
	serverId := parts[1]
	nicId := parts[2]
	firewallId := parts[3]

	fw, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, dcId,
		serverId, nicId, firewallId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("unable to find firewall rule %q", firewallId), nil)
		}
		return nil, utils.ToError(d, fmt.Sprintf("an error occurred while retrieving firewall rule %q: %s ", firewallId, err), nil)
	}

	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	if err := d.Set("server_id", serverId); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	if err := d.Set("nic_id", nicId); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	if err := setFirewallData(d, &fw); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	return []*schema.ResourceData{d}, nil
}

func getFirewallData(d *schema.ResourceData, path string, update bool) (ionoscloud.FirewallRule, diag.Diagnostics) {

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
		intIcmpType, err := strconv.Atoi(v.(string))
		if err != nil {
			return firewall, utils.ToDiags(d, fmt.Sprintf("could not parse icmpTpye %s: %s", v.(string), err), nil)
		}
		tempIcmpType := int32(intIcmpType)
		firewall.Properties.IcmpType = &tempIcmpType

	}
	if v, ok := d.GetOk(path + "icmp_code"); ok {
		intIcmpCode, err := strconv.Atoi(v.(string))
		if err != nil {
			return firewall, utils.ToDiags(d, fmt.Sprintf("could not parse icmpCode %s: %s", v.(string), err), nil)
		}
		tempIcmpCode := int32(intIcmpCode)
		firewall.Properties.IcmpCode = &tempIcmpCode

	}
	if v, ok := d.GetOk(path + "type"); ok {
		tempType := v.(string)
		firewall.Properties.Type = &tempType

	}
	return firewall, nil
}

func setFirewallData(d *schema.ResourceData, firewall *ionoscloud.FirewallRule) error {

	if firewall.Id != nil {
		d.SetId(*firewall.Id)
	}

	if firewall.Properties != nil {

		if firewall.Properties.Protocol != nil {
			err := d.Set("protocol", *firewall.Properties.Protocol)
			if err != nil {
				return fmt.Errorf("error while setting protocol property for firewall %s: %w", d.Id(), err)
			}
		}

		if firewall.Properties.Name != nil {
			err := d.Set("name", *firewall.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for firewall %s: %w", d.Id(), err)
			}
		}

		if firewall.Properties.SourceMac != nil {
			err := d.Set("source_mac", *firewall.Properties.SourceMac)
			if err != nil {
				return fmt.Errorf("error while setting source_mac property for firewall %s: %w", d.Id(), err)
			}
		}

		if firewall.Properties.SourceIp != nil {
			err := d.Set("source_ip", *firewall.Properties.SourceIp)
			if err != nil {
				return fmt.Errorf("error while setting source_ip property for firewall %s: %w", d.Id(), err)
			}
		}

		if firewall.Properties.TargetIp != nil {
			err := d.Set("target_ip", *firewall.Properties.TargetIp)
			if err != nil {
				return fmt.Errorf("error while setting target_ip property for firewall %s: %w", d.Id(), err)
			}
		}

		if firewall.Properties.PortRangeStart != nil {
			err := d.Set("port_range_start", *firewall.Properties.PortRangeStart)
			if err != nil {
				return fmt.Errorf("error while setting port_range_start property for firewall %s: %w", d.Id(), err)
			}
		}

		if firewall.Properties.PortRangeEnd != nil {
			err := d.Set("port_range_end", *firewall.Properties.PortRangeEnd)
			if err != nil {
				return fmt.Errorf("error while setting port_range_end property for firewall %s: %w", d.Id(), err)
			}
		}

		if firewall.Properties.IcmpType != nil {
			err := d.Set("icmp_type", strconv.Itoa(int(*firewall.Properties.IcmpType)))
			if err != nil {
				return fmt.Errorf("error while setting icmp_type property for firewall %s: %w", d.Id(), err)
			}
		}

		if firewall.Properties.IcmpCode != nil {
			err := d.Set("icmp_code", strconv.Itoa(int(*firewall.Properties.IcmpCode)))
			if err != nil {
				return fmt.Errorf("error while setting icmp_code property for firewall %s: %w", d.Id(), err)
			}
		}

		if firewall.Properties.Type != nil {
			err := d.Set("type", *firewall.Properties.Type)
			if err != nil {
				return fmt.Errorf("error while setting type property for firewall %s: %w", d.Id(), err)
			}
		}
	}
	return nil
}
