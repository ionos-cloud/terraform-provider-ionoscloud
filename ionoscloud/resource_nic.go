package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
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
			},
			"ips": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"datacenter_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

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
	if _, ok := d.GetOk("firewall_active"); ok {
		raw := d.Get("firewall_active").(bool)
		nic.Properties.FirewallActive = &raw
	}
	if _, ok := d.GetOk("nat"); ok {
		raw := d.Get("nat").(bool)
		nic.Properties.Nat = &raw
	}

	if v, ok := d.GetOk("ips"); ok {
		raw := v.([]interface{})
		if raw != nil && len(raw) > 0 {
			var ips []string
			for _, rawIp := range raw {
				ip := rawIp.(string)
				ips = append(ips, ip)
			}
			if ips != nil && len(ips) > 0 {
				nic.Properties.Ips = &ips
			}
		}
	}

	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nic, apiResponse, err := client.NicApi.DatacentersServersNicsPost(ctx, dcid, srvid).Nic(nic).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while creating a nic: %s", err))
		return diags
	}
	if nic.Id != nil {
		d.SetId(*nic.Id)
	}
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
	return resourceNicRead(ctx, d, meta)
}

func resourceNicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()

	rsp, apiResponse, err := client.NicApi.DatacentersServersNicsFindById(ctx, dcid, srvid, nicid).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		diags := diag.FromErr(fmt.Errorf("error occured while fetching a nic ID %s %s", d.Id(), err))
		return diags
	}

	if rsp.Properties != nil {
		log.Printf("[INFO] LAN ON NIC: %d", rsp.Properties.Lan)
		if rsp.Properties.Dhcp != nil {
			if err := d.Set("dhcp", *rsp.Properties.Dhcp); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
		if rsp.Properties.Lan != nil {
			if err := d.Set("lan", *rsp.Properties.Lan); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
		if rsp.Properties.Name != nil {
			if err := d.Set("name", *rsp.Properties.Name); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
		if rsp.Properties.Ips != nil && len(*rsp.Properties.Ips) > 0 {
			if err := d.Set("ips", *rsp.Properties.Ips); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
		if rsp.Properties.FirewallActive != nil {
			if err := d.Set("firewall_active", *rsp.Properties.FirewallActive); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
		if rsp.Properties.Mac != nil {
			if err := d.Set("mac", *rsp.Properties.Mac); err != nil {
				diags := diag.FromErr(err)
				return diags
			}
		}
	}

	return nil
}

func resourceNicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

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

	if d.HasChange("ips") {
		_, v := d.GetChange("ips")
		raw := v.([]interface{})
		if raw != nil && len(raw) > 0 {
			var ips []string
			for _, rawIp := range raw {
				ip := rawIp.(string)
				ips = append(ips, ip)
			}
			if ips != nil && len(ips) > 0 {
				properties.Ips = &ips
			}
		}
	}
	if d.HasChange("nat") {
		_, raw := d.GetChange("nat")
		nat := raw.(bool)
		properties.Nat = &nat
	}

	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()

	_, apiResponse, err := client.NicApi.DatacentersServersNicsPatch(ctx, dcid, srvid, nicid).Nic(properties).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occured while updating a nic: %s", err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceNicRead(ctx, d, meta)
}

func resourceNicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()

	_, apiResponse, err := client.NicApi.DatacentersServersNicsDelete(ctx, dcid, srvid, nicid).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a nic dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err))
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
