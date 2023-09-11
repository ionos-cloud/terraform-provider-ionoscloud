package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapinic"
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
	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}

	nic, err := cloudapinic.GetNicFromSchema(d, "")
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while getting nic from schema: %w", err))
		return diags
	}

	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	createdNic, apiResponse, err := ns.Create(ctx, dcid, srvid, nic)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while creating a nic: %w", err))
		return diags
	}

	if createdNic.Id != nil {
		d.SetId(*createdNic.Id)
	}

	//Sometimes there is an error because the nic is not found after it's created.
	//Probably a read write consistency issue.
	//We're retrying for 5 minutes. 404 - means we keep on trying.
	var foundNic = &ionoscloud.Nic{}
	err = retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		var err error
		foundNic, apiResponse, err = ns.FindById(ctx, dcid, srvid, *nic.Id, 0)
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Could not find nic with Id %s , retrying...", *nic.Id)
			return retry.RetryableError(fmt.Errorf("could not find nic, %w", err))
		}
		if err != nil {
			retry.NonRetryableError(err)
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
	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}
	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()
	nic, apiResponse, err := ns.FindById(ctx, dcid, srvid, nicid, 0)
	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] nic resource with id %s not found", nicid)
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a nic ID %s %w", d.Id(), err))
		return diags
	}

	if err := NicSetData(d, nic); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}
	dcId := d.Get("datacenter_id").(string)
	srvId := d.Get("server_id").(string)
	nicId := d.Id()

	nic, err := cloudapinic.GetNicFromSchema(d, "")
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("update error occured while getting nic from schema: %w", err))
		return diags
	}

	_, _, err = ns.Update(ctx, dcId, srvId, nicId, *nic.Properties)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while updating a nic: %w", err))
		return diags
	}

	return resourceNicRead(ctx, d, meta)
}

func resourceNicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}
	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()
	_, err := ns.Delete(ctx, dcid, srvid, nicid)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a nic dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err))
		return diags
	}
	d.SetId("")
	return nil
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
		//should not be checked for len, we want to set the empty slice anyway as the field is computed, and it will not be set by backend
		// if ipv6_cidr_block is not set on the lan
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
