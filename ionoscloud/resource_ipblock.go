package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func resourceIPBlock() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPBlockCreate,
		ReadContext:   resourceIPBlockRead,
		UpdateContext: resourceIPBlockUpdate,
		DeleteContext: resourceIPBlockDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIpBlockImporter,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ips": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ip_consumers": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"k8s_nodepool_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"k8s_cluster_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceIPBlockCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	size := d.Get("size").(int)
	sizeConverted := int32(size)
	location := d.Get("location").(string)
	name := d.Get("name").(string)
	ipblock := ionoscloud.IpBlock{
		Properties: &ionoscloud.IpBlockProperties{
			Size:     &sizeConverted,
			Location: &location,
			Name:     &name,
		},
	}

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	ipblock, apiResponse, err := client.IPBlocksApi.IpblocksPost(ctx).Ipblock(ipblock).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while reserving an ip block: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}
	d.SetId(*ipblock.Id)

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			d.SetId("")
		}
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceIPBlockRead(ctx, d, meta)
}

func resourceIPBlockRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	ipBlock, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching an ip block: %w", err), nil)
	}

	tflog.Info(ctx, "ip block fetched", map[string]interface{}{"ips": strings.Join(*ipBlock.Properties.Ips, ",")})

	if err := IpBlockSetData(d, &ipBlock); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}
func resourceIPBlockUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	request := ionoscloud.IpBlockProperties{}

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		name := n.(string)
		request.Name = &name
	}

	_, apiResponse, err := client.IPBlocksApi.IpblocksPatch(ctx, d.Id()).Ipblock(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating an ip block: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	return nil

}

func resourceIPBlockDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	apiResponse, err := client.IPBlocksApi.IpblocksDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while releasing an ipblock: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	d.SetId("")
	return nil
}

func resourceIpBlockImporter(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	importID := d.Id()

	location, parts := splitImportID(importID, ":")
	if len(parts) != 1 {
		return nil, fmt.Errorf("invalid import identifier: expected one of <location>:<ipblock-id> or <ipblock-id>, got: %s", importID)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, fmt.Errorf("failed validating import identifier %q: %w", importID, err)
	}

	ipBlockId := parts[0]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return nil, err
	}

	ipBlock, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, ipBlockId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("ipBlock does not exist %q", ipBlockId), nil)
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while trying to fetch the ipBlock %q, error:%w", ipBlockId, err), nil)

	}

	tflog.Info(ctx, "ipBlock found", map[string]interface{}{"id": *ipBlock.Id})

	d.SetId(*ipBlock.Id)

	if err := IpBlockSetData(d, &ipBlock); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}

func IpBlockSetData(d *schema.ResourceData, ipBlock *ionoscloud.IpBlock) error {
	if ipBlock == nil {
		return fmt.Errorf("ipblock is empty")
	}

	if ipBlock.Id != nil {
		d.SetId(*ipBlock.Id)
	}

	if ipBlock.Properties.Ips != nil && len(*ipBlock.Properties.Ips) > 0 {
		if err := d.Set("ips", *ipBlock.Properties.Ips); err != nil {
			return err
		}
	}

	if ipBlock.Properties.Location != nil {
		if err := d.Set("location", *ipBlock.Properties.Location); err != nil {
			return err
		}
	}

	if ipBlock.Properties.Size != nil {
		if err := d.Set("size", *ipBlock.Properties.Size); err != nil {
			return err
		}
	}

	if ipBlock.Properties.Name != nil {
		if err := d.Set("name", *ipBlock.Properties.Name); err != nil {
			return err
		}
	}

	if ipBlock.Properties.IpConsumers != nil && len(*ipBlock.Properties.IpConsumers) > 0 {
		var ipConsumers []any
		for _, ipConsumer := range *ipBlock.Properties.IpConsumers {
			ipConsumerEntry := make(map[string]any)
			utils.SetPropWithNilCheck(ipConsumerEntry, "ip", ipConsumer.Ip)
			utils.SetPropWithNilCheck(ipConsumerEntry, "mac", ipConsumer.Mac)
			utils.SetPropWithNilCheck(ipConsumerEntry, "nic_id", ipConsumer.NicId)
			utils.SetPropWithNilCheck(ipConsumerEntry, "server_id", ipConsumer.ServerId)
			utils.SetPropWithNilCheck(ipConsumerEntry, "server_name", ipConsumer.ServerName)
			utils.SetPropWithNilCheck(ipConsumerEntry, "datacenter_id", ipConsumer.DatacenterId)
			utils.SetPropWithNilCheck(ipConsumerEntry, "datacenter_name", ipConsumer.DatacenterName)
			utils.SetPropWithNilCheck(ipConsumerEntry, "k8s_nodepool_uuid", ipConsumer.K8sNodePoolUuid)
			utils.SetPropWithNilCheck(ipConsumerEntry, "k8s_cluster_uuid", ipConsumer.K8sClusterUuid)

			ipConsumers = append(ipConsumers, ipConsumerEntry)
		}

		if len(ipConsumers) > 0 {
			if err := d.Set("ip_consumers", ipConsumers); err != nil {
				return err
			}
		}
	}
	return nil
}
