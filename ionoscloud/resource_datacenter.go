package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"regexp"
)

func resourceDatacenter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatacenterCreate,
		ReadContext:   resourceDatacenterRead,
		UpdateContext: resourceDatacenterUpdate,
		DeleteContext: resourceDatacenterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			//Datacenter parameters
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A name of that resource",
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"location": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The physical location where the datacenter will be created. This will be where all of your servers live. Property cannot be modified after datacenter creation (disallowed in update requests)",
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description for the datacenter, e.g. staging, production",
				Computed:    true,
			},
			"sec_auth_protection": {
				Type:        schema.TypeBool,
				Description: "Boolean value representing if the data center requires extra protection e.g. two factor protection",
				Optional:    true,
			},
			"version": {
				Type:        schema.TypeInt,
				Description: "The version of that Data Center. Gets incremented with every change",
				Computed:    true,
			},
			"features": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of features supported by the location this data center is part of",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDatacenterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	datacenterName := d.Get("name").(string)
	datacenterLocation := d.Get("location").(string)

	datacenter := ionoscloud.Datacenter{
		Properties: &ionoscloud.DatacenterProperties{
			Name:     &datacenterName,
			Location: &datacenterLocation,
		},
	}

	if attr, ok := d.GetOk("description"); ok {
		attrStr := attr.(string)
		datacenter.Properties.Description = &attrStr
	}

	if attr, ok := d.GetOk("sec_auth_protection"); ok {
		attrStr := attr.(bool)
		datacenter.Properties.SecAuthProtection = &attrStr
	}

	createdDatacenter, apiResponse, err := client.DataCenterApi.DatacentersPost(ctx).Datacenter(datacenter).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf(
			"error creating data center (%s) (%s)", d.Id(), err))
		return diags
	}
	d.SetId(*createdDatacenter.Id)

	log.Printf("[INFO] DataCenter Id: %s", d.Id())

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

	return resourceDatacenterRead(ctx, d, meta)
}

func resourceDatacenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	datacenter, apiResponse, err := client.DataCenterApi.DatacentersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching a data center ID %s %s", d.Id(), err))
		return diags
	}

	if datacenter.Properties.Name != nil {
		err := d.Set("name", *datacenter.Properties.Name)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting name property for datacenter %s: %s", d.Id(), err))
			return diags
		}
	}

	if datacenter.Properties.Location != nil {
		err := d.Set("location", *datacenter.Properties.Location)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting location property for datacenter %s: %s", d.Id(), err))
			return diags
		}
	}

	if datacenter.Properties.Description != nil {
		err := d.Set("description", *datacenter.Properties.Description)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting description property for datacenter %s: %s", d.Id(), err))
			return diags
		}
	}

	if datacenter.Properties.SecAuthProtection != nil {
		err := d.Set("sec_auth_protection", *datacenter.Properties.SecAuthProtection)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting sec_auth_protection property for datacenter %s: %s", d.Id(), err))
			return diags
		}
	}

	if datacenter.Properties.Features != nil && len(*datacenter.Properties.Features) > 0 {
		err := d.Set("features", *datacenter.Properties.Features)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting features property for datacenter %s: %s", d.Id(), err))
			return diags
		}
	}

	return nil
}

func resourceDatacenterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient
	obj := ionoscloud.DatacenterProperties{}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")
		newNameStr := newName.(string)
		obj.Name = &newNameStr
	}

	if d.HasChange("description") {
		_, newDescription := d.GetChange("description")
		newDescriptionStr := newDescription.(string)
		obj.Description = &newDescriptionStr
	}

	if d.HasChange("location") {
		oldLocation, newLocation := d.GetChange("location")
		diags := diag.FromErr(fmt.Errorf("data center is created in %s location. You can not change location of the data center to %s. It requires recreation of the data center", oldLocation, newLocation))
		return diags
	}

	if d.HasChange("sec_auth_protection") {
		_, newSecAuthProtection := d.GetChange("sec_auth_protection")
		newSecAuthProtectionStr := newSecAuthProtection.(bool)
		obj.SecAuthProtection = &newSecAuthProtectionStr
	}

	_, apiResponse, err := client.DataCenterApi.DatacentersPatch(ctx, d.Id()).Datacenter(obj).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while update the data center ID %s %s", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceDatacenterRead(ctx, d, meta)
}

func resourceDatacenterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	_, apiResponse, err := client.DataCenterApi.DatacentersDelete(ctx, d.Id()).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting the data center ID %s %s", d.Id(), err))
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

