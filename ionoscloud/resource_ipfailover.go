package ionoscloud

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/lanSvc"
	"log"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/uuidgen"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/slice"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func resourceLanIPFailover() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLanIPFailoverCreate,
		ReadContext:   resourceLanIPFailoverRead,
		UpdateContext: resourceLanIPFailoverUpdate,
		DeleteContext: resourceLanIPFailoverDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIpFailoverImporter,
		},
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Failover IP",
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(validation.IsIPAddress)),
			},
			"nicuuid": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Description:      "The UUID of the master NIC",
			},
			"lan_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceLanIPFailoverCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	dcId := d.Get("datacenter_id").(string)
	lanId := d.Get("lan_id").(string)
	ip := d.Get("ip").(string)
	nicUuid := d.Get("nicuuid").(string)

	// First, retrieve the existing IP Failover groups
	lan, apiResponse, err := lanSvc.FindLanById(*client, ctx, dcId, lanId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return diag.FromErr(fmt.Errorf("unable to find the LAN with ID: %s, datacenter ID: %s", lanId, dcId))
		}
		return diag.FromErr(err)
	}
	if lan.Properties == nil || lan.Properties.IpFailover == nil {
		return diag.FromErr(fmt.Errorf("expected a LAN response containing IP failover groups but received 'nil' instead"))
	}

	// Add the new IP failover group to the list
	*lan.Properties.IpFailover = append(*lan.Properties.IpFailover, ionoscloud.IPFailover{
		Ip:      &ip,
		NicUuid: &nicUuid,
	})

	// Modify the LAN using the new list
	lan, apiResponse, err = client.LANsApi.DatacentersLansPatch(ctx, dcId, lanId).Lan(*lan.Properties).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while patching a lans IP failover group, LAN ID: %s, error: %w", lanId, err))
	}

	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	// Use the IP in order to generate the resource ID
	d.SetId(uuidgen.GenerateUuidFromName(ip))
	return nil
}

func resourceLanIPFailoverRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	dcId := d.Get("datacenter_id").(string)
	lanId := d.Get("lan_id").(string)
	ip := d.Get("ip").(string)

	// Retrieve the existing IP Failover groups
	lan, apiResponse, err := lanSvc.FindLanById(*client, ctx, dcId, lanId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("unable to find the LAN with ID: %s, datacenter ID: %s", lanId, dcId))
		}
		return diag.FromErr(err)
	}
	if lan.Properties == nil || lan.Properties.IpFailover == nil {
		return diag.FromErr(fmt.Errorf("expected a LAN response containing IP failover groups but received 'nil' instead"))
	}

	// Iterate through IP Failover groups and select the proper one using the IP (the IP acts like
	// an identifier) and then set the data accordingly.
	ipFailoverGroups := lan.Properties.IpFailover
	ipFailoverGroupFound := false
	if lan.Properties != nil && ipFailoverGroups != nil && len(*ipFailoverGroups) > 0 {
		for _, ipFailoverGroup := range *ipFailoverGroups {
			if *ipFailoverGroup.Ip == ip {
				if err := d.Set("nicuuid", *ipFailoverGroup.NicUuid); err != nil {
					return diag.FromErr(utils.GenerateSetError(constant.ResourceIpFailover, "nicuuid", err))
				}
				ipFailoverGroupFound = true
				break
			}
		}
	}

	// If the IP failover group was not found, set the ID to the empty string in order to remove the
	// resource from the state.
	if !ipFailoverGroupFound {
		log.Printf("[INFO] IP Failover Group with IP: %s was not found in LAN with ID: %s, datacenter ID: %s", ip, lanId, dcId)
		d.SetId("")
	}

	return nil
}

func resourceLanIPFailoverUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	dcId := d.Get("datacenter_id").(string)
	lanId := d.Get("lan_id").(string)
	ip := d.Get("ip").(string)

	if d.HasChange("nicuuid") {
		oldValue, newValue := d.GetChange("nicuuid")
		newNicUuid := newValue.(string)
		oldNicUuid := oldValue.(string)

		// First, retrieve the existing IP Failover groups
		lan, apiResponse, err := lanSvc.FindLanById(*client, ctx, dcId, lanId)
		if err != nil {
			if apiResponse.HttpNotFound() {
				return diag.FromErr(fmt.Errorf("unable to find the LAN with ID: %s, datacenter ID: %s", lanId, dcId))
			}
			return diag.FromErr(err)
		}
		if lan.Properties == nil || lan.Properties.IpFailover == nil {
			return diag.FromErr(fmt.Errorf("expected a LAN response containing IP failover groups but received 'nil' instead"))
		}

		// Add the new IP failover group to the list
		*lan.Properties.IpFailover = append(*lan.Properties.IpFailover, ionoscloud.IPFailover{
			Ip:      &ip,
			NicUuid: &newNicUuid,
		})

		// Remove the old IP failover group from the list
		*lan.Properties.IpFailover = slice.DeleteFrom(*lan.Properties.IpFailover, ionoscloud.IPFailover{
			Ip:      &ip,
			NicUuid: &oldNicUuid,
		})

		_, apiResponse, err = client.LANsApi.DatacentersLansPatch(ctx, dcId, lanId).Lan(*lan.Properties).Execute()
		apiResponse.LogInfo()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while patching the lan with ID: %s, error: %w", lanId, err))
			return diags
		}

		// Wait, catching any errors
		_, errState := cloudapi.GetStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(errState)
			return diags
		}
	}
	return nil
}

func resourceLanIPFailoverDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	dcId := d.Get("datacenter_id").(string)
	lanId := d.Get("lan_id").(string)
	ip := d.Get("ip").(string)
	nicUuid := d.Get("nicuuid").(string)

	// First, retrieve the existing IP Failover groups
	lan, apiResponse, err := lanSvc.FindLanById(*client, ctx, dcId, lanId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return diag.FromErr(fmt.Errorf("unable to find the LAN with ID: %s, datacenter ID: %s", lanId, dcId))
		}
		return diag.FromErr(err)
	}
	if lan.Properties == nil || lan.Properties.IpFailover == nil {
		return diag.FromErr(fmt.Errorf("expected a LAN response containing IP failover groups but received 'nil' instead"))
	}

	// Remove the failover group from the list
	*lan.Properties.IpFailover = slice.DeleteFrom(*lan.Properties.IpFailover, ionoscloud.IPFailover{
		Ip:      &ip,
		NicUuid: &nicUuid,
	})

	_, apiResponse, err = client.LANsApi.DatacentersLansPatch(ctx, dcId, lanId).Lan(*lan.Properties).Execute()
	apiResponse.LogInfo()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while removing an IP failover group with IP: %s for the LAN with ID: %s, datacenter ID: %s, error: %w", ip, lanId, dcId, err))
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

func resourceIpFailoverImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("invalid import ID: %s. Expecting {datacenter}/{lan}/{ip}", d.Id())
	}
	dcId := parts[0]
	lanId := parts[1]
	ip := parts[2]

	client := meta.(services.SdkBundle).CloudApiClient

	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcId, lanId).Execute()
	apiResponse.LogInfo()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("unable to find the LAN with ID: %s, datacenter ID: %s", lanId, dcId)
		}
		return nil, fmt.Errorf("error while fetching LAN with ID: %s, datacenter ID: %s, err: %w", lanId, dcId, err)
	}

	log.Printf("[INFO] lan found: %+v", lan)

	ipFailoverGroups := lan.Properties.IpFailover
	if lan.Properties != nil && ipFailoverGroups != nil && len(*ipFailoverGroups) > 0 {
		for _, ipFailoverGroup := range *ipFailoverGroups {
			// Search for the appropiate IP Failover Group using the provided IP
			if *ipFailoverGroup.Ip == ip {
				// Set all the information only if the IP Failover Group exists
				// Use the IP in order to generate the resource ID
				d.SetId(uuidgen.GenerateUuidFromName(ip))

				if err := d.Set("datacenter_id", dcId); err != nil {
					return nil, utils.GenerateSetError(constant.ResourceIpFailover, "datacenter_id", err)
				}
				if err := d.Set("lan_id", lanId); err != nil {
					return nil, utils.GenerateSetError(constant.ResourceIpFailover, "lan_id", err)
				}
				if err := d.Set("ip", ip); err != nil {
					return nil, utils.GenerateSetError(constant.ResourceIpFailover, "ip", err)
				}
				if err := d.Set("nicuuid", *ipFailoverGroup.NicUuid); err != nil {
					return nil, utils.GenerateSetError(constant.ResourceIpFailover, "nicuuid", err)
				}
				// After we find the IP Failover Group, we can stop searching since the IP is unique
				return []*schema.ResourceData{d}, nil
			}
		}
	}
	return nil, fmt.Errorf("IP Failover Group with IP: %s does not exist in the LAN with ID: %s, datacenter ID: %s", ip, lanId, dcId)
}
