package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func resourceIPBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPBlockCreate,
		Read:   resourceIPBlockRead,
		Update: resourceIPBlockUpdate,
		Delete: resourceIPBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceIPBlockCreate(d *schema.ResourceData, meta interface{}) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)
	if cancel != nil {
		defer cancel()
	}

	ipblock, apiResponse, err := client.IPBlocksApi.IpblocksPost(ctx).Ipblock(ipblock).Execute()

	if err != nil {
		return fmt.Errorf("an error occured while reserving an ip block: %s", err)
	}
	d.SetId(*ipblock.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceIPBlockRead(d, meta)
}

func resourceIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	ipBlock, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("an error occured while fetching an ip block ID %s %s", d.Id(), err)
	}

	log.Printf("[INFO] IPS: %s", strings.Join(*ipBlock.Properties.Ips, ","))

	d.Set("ips", *ipBlock.Properties.Ips)
	d.Set("location", *ipBlock.Properties.Location)
	d.Set("size", *ipBlock.Properties.Size)
	d.Set("name", *ipBlock.Properties.Name)

	return nil
}
func resourceIPBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	request := ionoscloud.IpBlockProperties{}

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		name := n.(string)
		request.Name = &name
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)
	if cancel != nil {
		defer cancel()
	}

	_, _, err := client.IPBlocksApi.IpblocksPatch(ctx, d.Id()).Ipblock(request).Execute()

	if err != nil {
		return fmt.Errorf("an error occured while updating an ip block ID %s %s", d.Id(), err)
	}

	return nil

}

func resourceIPBlockDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.IPBlocksApi.IpblocksDelete(ctx, d.Id()).Execute()
	if err != nil {
		return fmt.Errorf("an error occured while releasing an ipblock ID: %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
