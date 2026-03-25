package ionoscloud

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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
			"location": {
				Type:        schema.TypeString,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				Optional:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceIpFailoverRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	dcId := d.Get("datacenter_id").(string)
	lanId := d.Get("lan_id").(string)
	ip := d.Get("ip").(string)

	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcId, lanId).Execute()
	apiResponse.LogInfo()
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return diagutil.ToDiags(d, fmt.Errorf("unable to find the LAN with ID: %s, datacenter ID: %s", lanId, dcId), &diagutil.ErrorContext{StatusCode: apiResponse.StatusCode})
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while fetching LAN with ID: %s, datacenter ID: %s, err: %w", lanId, dcId, err), &diagutil.ErrorContext{StatusCode: apiResponse.StatusCode})
	}
	if lan.Properties == nil || lan.Properties.IpFailover == nil {
		return diagutil.ToDiags(d, fmt.Errorf("expected a LAN response containing IP failover groups but received 'nil' instead"), nil)
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
					return diagutil.ToDiags(d, utils.GenerateSetError(constant.ResourceIpFailover, "nicuuid", err), nil)
				}
				ipFailoverGroupFound = true
				// After we find the IP Failover Group, we can stop searching since the IP is unique
				break
			}
		}
	}

	if !ipFailoverGroupFound {
		return diagutil.ToDiags(d, fmt.Errorf("IP Failover Group with IP: %s does not exist in the LAN with ID: %s, datacenter ID: %s", ip, lanId, dcId), nil)
	}
	return nil
}
