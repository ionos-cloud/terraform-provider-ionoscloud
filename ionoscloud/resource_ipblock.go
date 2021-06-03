package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func resourceIPBlock() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPBlockCreate,
		ReadContext:   resourceIPBlockRead,
		UpdateContext: resourceIPBlockUpdate,
		DeleteContext: resourceIPBlockDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

	if err != nil {
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while reserving an ip block: %s %s", err, payload))
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

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching an ip block ID %s %s %s", d.Id(), err, payload))
		return diags
	}

	log.Printf("[INFO] IPS: %s", strings.Join(*ipBlock.Properties.Ips, ","))

	if ipBlock.Properties.Ips != nil {
		if err := d.Set("ips", *ipBlock.Properties.Ips); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if ipBlock.Properties.Location != nil {
		if err := d.Set("location", *ipBlock.Properties.Location); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if ipBlock.Properties.Size != nil {
		if err := d.Set("size", *ipBlock.Properties.Size); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if ipBlock.Properties.Name != nil {
		if err := d.Set("name", *ipBlock.Properties.Name); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
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

	if err != nil {
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while updating an ip block ID %s %s %s", d.Id(), err, payload))
		return diags
	}

	return nil
}

func resourceIPBlockDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	_, apiResponse, err := client.IPBlocksApi.IpblocksDelete(ctx, d.Id()).Execute()
	if err != nil {
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while releasing an ipblock ID: %s %s %s", d.Id(), err, payload))
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
