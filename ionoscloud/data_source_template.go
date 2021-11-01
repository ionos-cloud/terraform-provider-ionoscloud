package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"strings"
)

func dataSourceTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cores": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"storage_size": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	templates, apiResponse, err := client.TemplatesApi.TemplatesGet(ctx).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching IonosCloud templates %s ", err))
		return diags
	}

	name := d.Get("name").(string)
	cores, coresOk := d.GetOk("cores")
	ram, ramOk := d.GetOk("ram")
	storageSize, storageSizeOk := d.GetOk("storage_size")

	var results []ionoscloud.Template

	if templates.Items != nil {
		for _, tmp := range *templates.Items {
			if strings.Contains(strings.ToLower(*tmp.Properties.Name), strings.ToLower(name)) {
				results = append(results, tmp)
			}
		}
	}

	if coresOk {
		var coresResults []ionoscloud.Template
		for _, tmp := range results {
			cores := float32(cores.(float64))
			if tmp.Properties.Cores != nil && *tmp.Properties.Cores == cores {
				coresResults = append(coresResults, tmp)
			}
		}
		results = coresResults
	}

	if ramOk {
		var ramResults []ionoscloud.Template
		for _, tmp := range results {
			ram := float32(ram.(float64))
			if tmp.Properties.Ram != nil && *tmp.Properties.Ram == ram {
				ramResults = append(ramResults, tmp)
			}
		}
		results = ramResults
	}

	if storageSizeOk {
		var storageSizeResults []ionoscloud.Template
		storageSize := float32(storageSize.(float64))
		for _, tmp := range results {
			if tmp.Properties.StorageSize != nil && *tmp.Properties.StorageSize == storageSize {
				storageSizeResults = append(storageSizeResults, tmp)
			}
		}
		results = storageSizeResults
	}

	if len(results) > 1 {
		diags := diag.FromErr(fmt.Errorf("There is more than one template that match the search criteria "))
		return diags
	}

	if len(results) == 0 {
		diags := diag.FromErr(fmt.Errorf("There are no templates that match the search criteria "))
		return diags
	}

	if results[0].Properties.Name != nil {
		err := d.Set("name", *results[0].Properties.Name)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting name property for image %s: %s", d.Id(), err))
			return diags
		}
	}

	if results[0].Properties.Cores != nil {
		if err := d.Set("cores", *results[0].Properties.Cores); err != nil {
			return diag.FromErr(err)
		}
	}
	if results[0].Properties.Ram != nil {
		if err := d.Set("ram", *results[0].Properties.Ram); err != nil {
			return diag.FromErr(err)
		}
	}
	if results[0].Properties.StorageSize != nil {
		if err := d.Set("storage_size", *results[0].Properties.StorageSize); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(*results[0].Id)

	return nil
}
