package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
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
	client := meta.(SdkBundle).CloudApiClient

	templates, apiResponse, err := client.TemplatesApi.TemplatesGet(ctx).Depth(1).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching IonosCloud templates %w ", err))
		return diags
	}

	name, nameOk := d.GetOk("name")
	cores, coresOk := d.GetOk("cores")
	ram, ramOk := d.GetOk("ram")
	storageSize, storageSizeOk := d.GetOk("storage_size")

	var results []ionoscloud.Template

	if nameOk && templates.Items != nil {
		for _, tmp := range *templates.Items {
			if strings.Contains(strings.ToLower(*tmp.Properties.Name), strings.ToLower(name.(string))) {
				results = append(results, tmp)
			}
		}
	} else if templates.Items != nil {
		results = *templates.Items
	}

	if coresOk {
		cores := float32(cores.(float64))
		if results != nil {
			var coresResults []ionoscloud.Template
			for _, tmp := range results {
				if tmp.Properties.Cores != nil && *tmp.Properties.Cores == cores {
					coresResults = append(coresResults, tmp)
				}
			}
			results = coresResults
		}
	}

	if ramOk {
		ram := float32(ram.(float64))
		if results != nil {
			var ramResults []ionoscloud.Template
			for _, tmp := range results {
				if tmp.Properties.Ram != nil && *tmp.Properties.Ram == ram {
					ramResults = append(ramResults, tmp)
				}
			}
			results = ramResults
		}
	}

	if storageSizeOk {
		storageSize := float32(storageSize.(float64))
		if results != nil {
			var storageSizeResults []ionoscloud.Template
			for _, tmp := range results {
				if tmp.Properties != nil && tmp.Properties.StorageSize != nil && *tmp.Properties.StorageSize == storageSize {
					storageSizeResults = append(storageSizeResults, tmp)
				}
			}
			results = storageSizeResults
		}
	}

	var template ionoscloud.Template

	if results == nil || len(results) == 0 {
		return diag.FromErr(fmt.Errorf("no template found with the specified criteria: name = %s, cores = %v, ram = %v, storage_size = %v", name.(string), cores.(float64), ram.(float64), storageSize.(float64)))
	} else if len(results) > 1 {
		return diag.FromErr(fmt.Errorf("more than one template found with the specified criteria: name = %s, cores = %v, ram = %v, storage_size = %v", name.(string), cores.(float64), ram.(float64), storageSize.(float64)))
	} else {
		template = results[0]
	}

	if err = setTemplateData(d, &template); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setTemplateData(d *schema.ResourceData, template *ionoscloud.Template) error {
	d.SetId(*template.Id)

	if template.Properties != nil {
		if template.Properties.Name != nil {
			err := d.Set("name", *template.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for image %s: %w", d.Id(), err)
			}
		}

		if template.Properties.Cores != nil {
			if err := d.Set("cores", *template.Properties.Cores); err != nil {
				return err
			}
		}
		if template.Properties.Ram != nil {
			if err := d.Set("ram", *template.Properties.Ram); err != nil {
				return err
			}
		}
		if template.Properties.StorageSize != nil {
			if err := d.Set("storage_size", *template.Properties.StorageSize); err != nil {
				return err
			}
		}

	}

	return nil
}
