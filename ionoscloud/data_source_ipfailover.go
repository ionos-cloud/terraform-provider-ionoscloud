package ionoscloud

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/uuidgen"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceIpFailover() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIpFailoverRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPAddress),
			},
			"nicuuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lan_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceIpFailoverRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	dcId := d.Get("datacenter_id").(string)
	lanId := d.Get("lan_id").(string)
	ip := d.Get("ip").(string)

	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcId, lanId).Execute()
	apiResponse.LogInfo()
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("unable to find the LAN with ID: %s, datacenter ID: %s", lanId, dcId))
		}
		return diag.FromErr(fmt.Errorf("error while fetching LAN with ID: %s, datacenter ID: %s, err: %w", lanId, dcId, err))
	}
	if lan.Properties == nil || lan.Properties.IpFailover == nil {
		return diag.FromErr(fmt.Errorf("expected a LAN response containing IP failover groups but received 'nil' instead"))
	}

	ipFailoverGroups := lan.Properties.IpFailover
	ipFailoverGroupFound := false
	if lan.Properties != nil && ipFailoverGroups != nil && len(*ipFailoverGroups) > 0 {
		for _, ipFailoverGroup := range *ipFailoverGroups {
			// Search for the appropriate IP Failover Group using the provided IP
			if *ipFailoverGroup.Ip == ip {
				// Set the information only if the IP Failover Group exists
				// Use the IP in order to generate the resource ID
				d.SetId(uuidgen.GenerateUuidFromName(ip))

				if err := d.Set("nicuuid", *ipFailoverGroup.NicUuid); err != nil {
					return diag.FromErr(utils.GenerateSetError(constant.ResourceIpFailover, "nicuuid", err))
				}
				ipFailoverGroupFound = true
				// After we find the IP Failover Group, we can stop searching since the IP is unique
				break
			}
		}
	}

	if !ipFailoverGroupFound {
		return diag.FromErr(fmt.Errorf("IP Failover Group with IP: %s does not exist in the LAN with ID: %s, datacenter ID: %s", ip, lanId, dcId))
	}
	return nil
}
