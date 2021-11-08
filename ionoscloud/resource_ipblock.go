package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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

func resourceIPBlockCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

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

	ipblock, apiResponse, err := client.IPBlocksApi.IpblocksPost(ctx).Ipblock(ipblock).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while reserving an ip block: %s", err))
		return diags
	}
	d.SetId(*ipblock.Id)

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

	return resourceIPBlockRead(ctx, d, meta)
}

func resourceIPBlockRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	ipBlock, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching an ip block ID %s %s", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] IPS: %s", strings.Join(*ipBlock.Properties.Ips, ","))

	if err := IpBlockSetData(d, &ipBlock); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
func resourceIPBlockUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	request := ionoscloud.IpBlockProperties{}

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		name := n.(string)
		request.Name = &name
	}

	_, apiResponse, err := client.IPBlocksApi.IpblocksPatch(ctx, d.Id()).Ipblock(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating an ip block ID %s %s", d.Id(), err))
		return diags
	}

	return nil

}

func resourceIPBlockDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	apiResponse, err := client.IPBlocksApi.IpblocksDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while releasing an ipblock ID: %s %s", d.Id(), err))
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

func resourceIpBlockImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	ipBlockId := d.Id()

	ipBlock, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, ipBlockId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the ipBlock %q", ipBlockId)
		}
		return nil, fmt.Errorf("ipBlock does not exist %q", ipBlockId)
	}

	log.Printf("[INFO] ipBlock found: %+v", ipBlock)

	d.SetId(*ipBlock.Id)

	if err := IpBlockSetData(d, &ipBlock); err != nil {
		return nil, err
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
		var ipConsumers []interface{}
		for _, ipConsumer := range *ipBlock.Properties.IpConsumers {
			ipConsumerEntry := make(map[string]interface{})
			setPropWithNilCheck(ipConsumerEntry, "ip", ipConsumer.Ip)
			setPropWithNilCheck(ipConsumerEntry, "mac", ipConsumer.Mac)
			setPropWithNilCheck(ipConsumerEntry, "nic_id", ipConsumer.NicId)
			setPropWithNilCheck(ipConsumerEntry, "server_id", ipConsumer.ServerId)
			setPropWithNilCheck(ipConsumerEntry, "server_name", ipConsumer.ServerName)
			setPropWithNilCheck(ipConsumerEntry, "datacenter_id", ipConsumer.DatacenterId)
			setPropWithNilCheck(ipConsumerEntry, "datacenter_name", ipConsumer.DatacenterName)
			setPropWithNilCheck(ipConsumerEntry, "k8s_nodepool_uuid", ipConsumer.K8sNodePoolUuid)
			setPropWithNilCheck(ipConsumerEntry, "k8s_cluster_uuid", ipConsumer.K8sClusterUuid)

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
