package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceNic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNicCreate,
		ReadContext:   resourceNicRead,
		UpdateContext: resourceNicUpdate,
		DeleteContext: resourceNicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNicImport,
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
				Default:  true,
			},
			"dhcpv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether this NIC receives an IPv6 address through DHCP.",
			},
			"ipv6_cidr_block": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IPv6 CIDR block assigned to the NIC.",
			},
			"ips": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: utils.DiffEmptyIps,
				},
				Computed:    true,
				Optional:    true,
				Description: "Collection of IP addresses assigned to a nic. Explicitly assigned public IPs need to come from reserved IP blocks, Passing value null or empty array will assign an IP address automatically.",
			},
			"ipv6_ips": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Collection for IPv6 addresses assigned to a nic. Explicitly assigned IPv6 addresses need to come from inside the IPv6 CIDR block assigned to the nic.",
			},
			"firewall_active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"firewall_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pci_slot": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	nic := getNicData(d, "")

	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)

	nic, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsPost(ctx, dcid, srvid).Nic(nic).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while creating a nic: %w", err))
		return diags
	}
	if nic.Id != nil {
		d.SetId(*nic.Id)
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if cloudapi.IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}
	//Sometimes there is an error because the nic is not found after it's created.
	//Probably a read write consistency issue.
	//We're retrying for 5 minutes. 404 - means we keep on trying.
	var foundNic = &ionoscloud.Nic{}
	err = resource.RetryContext(ctx, 5*time.Minute, func() *resource.RetryError {
		var err error
		*foundNic, apiResponse, err = client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, dcid, srvid, *nic.Id).Execute()
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Could not find nic with Id %s , retrying...", *nic.Id)
			return resource.RetryableError(fmt.Errorf("could not find nic, %w", err))
		}
		if err != nil {
			resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	if foundNic == nil || *foundNic.Id == "" {
		return diag.FromErr(fmt.Errorf("could not find nic with id %s after creation ", *nic.Id))
	}

	return diag.FromErr(NicSetData(d, foundNic))
}

func resourceNicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()

	nic, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, dcid, srvid, nicid).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] nic resource with id %s not found", nicid)
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a nic ID %s %w", d.Id(), err))
		return diags
	}

	if err := NicSetData(d, &nic); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	srvId := d.Get("server_id").(string)
	nicId := d.Id()

	nic := getNicData(d, "")

	_, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsPatch(ctx, dcId, srvId, nicId).Nic(*nic.Properties).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while updating a nic: %w", err))
		return diags
	}

	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceNicRead(ctx, d, meta)
}

func resourceNicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()
	apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsDelete(ctx, dcid, srvid, nicid).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a nic dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err))
		return diags
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	d.SetId("")
	return nil
}

func getNicData(d *schema.ResourceData, path string) ionoscloud.Nic {

	nic := ionoscloud.Nic{
		Properties: &ionoscloud.NicProperties{},
	}

	lanInt := int32(d.Get(path + "lan").(int))
	nic.Properties.Lan = &lanInt

	if v, ok := d.GetOk(path + "name"); ok {
		vStr := v.(string)
		nic.Properties.Name = &vStr
	}

	dhcp := d.Get(path + "dhcp").(bool)
	if dhcpv6, ok := d.GetOkExists(path + "dhcpv6"); ok {
		dhcpv6 := dhcpv6.(bool)
		nic.Properties.Dhcpv6 = &dhcpv6
	} else {
		nic.Properties.SetDhcpv6Nil()
	}
	fwActive := d.Get(path + "firewall_active").(bool)
	nic.Properties.Dhcp = &dhcp
	nic.Properties.FirewallActive = &fwActive

	if _, ok := d.GetOk(path + "firewall_type"); ok {
		raw := d.Get(path + "firewall_type").(string)
		nic.Properties.FirewallType = &raw
	}

	if v, ok := d.GetOk(path + "ips"); ok {
		raw := v.([]interface{})
		if raw != nil && len(raw) > 0 {
			ips := make([]string, 0)
			for _, rawIp := range raw {
				if rawIp != nil {
					ip := rawIp.(string)
					ips = append(ips, ip)
				}
			}
			if ips != nil && len(ips) > 0 {
				nic.Properties.Ips = &ips
			}
		}
	}

	if v, ok := d.GetOk(path + "ipv6_ips"); ok {
		raw := v.([]interface{})
		ipv6_ips := make([]string, len(raw))
		utils.DecodeInterfaceToStruct(raw, ipv6_ips)
		if len(ipv6_ips) > 0 {
			nic.Properties.Ipv6Ips = &ipv6_ips
		}
	}

	if v, ok := d.GetOk(path + "ipv6_cidr_block"); ok {
		ipv6_block := v.(string)
		nic.Properties.Ipv6CidrBlock = &ipv6_block
	}

	return nic
}

func NicSetData(d *schema.ResourceData, nic *ionoscloud.Nic) error {
	if nic == nil {
		return fmt.Errorf("nic is empty")
	}

	if nic.Id != nil {
		d.SetId(*nic.Id)
	}

	if nic.Properties != nil {
		log.Printf("[INFO] LAN ON NIC: %d", nic.Properties.Lan)
		if nic.Properties.Dhcp != nil {
			if err := d.Set("dhcp", *nic.Properties.Dhcp); err != nil {
				return fmt.Errorf("error setting dhcp %w", err)
			}
		}

		if nic.Properties.Dhcpv6 != nil {
			if err := d.Set("dhcpv6", *nic.Properties.Dhcpv6); err != nil {
				return fmt.Errorf("error setting dhcpv6 %w", err)
			}
		}
		if nic.Properties.Lan != nil {
			if err := d.Set("lan", *nic.Properties.Lan); err != nil {
				return fmt.Errorf("error setting lan %w", err)
			}
		}
		if nic.Properties.Name != nil {
			if err := d.Set("name", *nic.Properties.Name); err != nil {
				return fmt.Errorf("error setting name %w", err)
			}
		}
		if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
			if err := d.Set("ips", *nic.Properties.Ips); err != nil {
				return fmt.Errorf("error setting ips %w", err)
			}
		}
		//should no be checked for len, we want to set the empty slice anyway as the field is computed and it will not be set by backend
		// if  ipv6_cidr_block is not set on the lan
		if nic.Properties.Ipv6Ips != nil {
			if err := d.Set("ipv6_ips", *nic.Properties.Ipv6Ips); err != nil {
				return fmt.Errorf("error setting ipv6_ips %w", err)
			}
		}
		if nic.Properties.Ipv6CidrBlock != nil {
			if err := d.Set("ipv6_cidr_block", *nic.Properties.Ipv6CidrBlock); err != nil {
				return fmt.Errorf("error setting ipv6_cidr_block %w", err)
			}
		}
		if nic.Properties.FirewallActive != nil {
			if err := d.Set("firewall_active", *nic.Properties.FirewallActive); err != nil {
				return fmt.Errorf("error setting firewall_active %w", err)
			}
		}
		if nic.Properties.FirewallType != nil {
			if err := d.Set("firewall_type", *nic.Properties.FirewallType); err != nil {
				return fmt.Errorf("error setting firewall_type %w", err)
			}
		}
		if nic.Properties.Mac != nil {
			if err := d.Set("mac", *nic.Properties.Mac); err != nil {
				return fmt.Errorf("error setting mac %w", err)
			}
		}
		if nic.Properties.DeviceNumber != nil {
			if err := d.Set("device_number", *nic.Properties.DeviceNumber); err != nil {
				return fmt.Errorf("error setting device_number %w", err)
			}
		}
		if nic.Properties.PciSlot != nil {
			if err := d.Set("pci_slot", *nic.Properties.PciSlot); err != nil {
				return fmt.Errorf("error setting pci_slot %w", err)
			}
		}
	}

	return nil
}

func resourceNicImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}/{nic}", d.Id())
	}
	dcId := parts[0]
	sId := parts[1]
	nicId := parts[2]

	client := meta.(services.SdkBundle).CloudApiClient

	nic, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, dcId, sId, nicId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if !apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the nic %q", nicId)
		}
		return nil, fmt.Errorf("lan does not exist%q", nicId)
	}

	err = d.Set("datacenter_id", dcId)
	if err != nil {
		return nil, err
	}
	err = d.Set("server_id", sId)
	if err != nil {
		return nil, err
	}

	if err := NicSetData(d, &nic); err != nil {
		return nil, err
	}

	log.Printf("[INFO] nic found: %+v", nic)

	return []*schema.ResourceData{d}, nil
}
