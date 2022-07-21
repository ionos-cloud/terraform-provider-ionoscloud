package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
)

func dataSourceTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching IonosCloud templates %s ", err))
		return diags
	}

	idValue, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	cores, coresOk := d.GetOk("cores")
	ram, ramOk := d.GetOk("ram")
	storageSize, storageSizeOk := d.GetOk("storage_size")

	id := idValue.(string)

	if idOk && (nameOk || coresOk || ramOk || storageSizeOk) {
		return diag.FromErr(fmt.Errorf("id and name/cores/ram/storage_size cannot be both specified in the same time, choose between id or a combination of other parameters"))
	}
	if !idOk && !nameOk && !coresOk && !ramOk && !storageSizeOk {
		return diag.FromErr(fmt.Errorf("please provide either the template id or other parameter like name, cores or ram"))
	}

	var results []ionoscloud.Template
	var template ionoscloud.Template

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for image by id %s", id)
		template, apiResponse, err = client.TemplatesApi.TemplatesFindById(ctx, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the nat gateway rule %s: %s", id, err))
		}
	} else {
		if nameOk && templates.Items != nil {
			for _, tmp := range *templates.Items {
				if tmp.Properties != nil && tmp.Properties.Name != nil && strings.EqualFold(*tmp.Properties.Name, name.(string)) {
					results = append(results, tmp)
				}
			}
			if len(results) == 0 {
				return diag.FromErr(fmt.Errorf("no result found with the specified criteria: name %f", name))
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
				if len(coresResults) == 0 {
					return diag.FromErr(fmt.Errorf("no result found with the specified criteria: cores %f", cores))
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
				if len(ramResults) == 0 {
					return diag.FromErr(fmt.Errorf("no result found with the specified criteria: ram %f", ram))
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
				if len(storageSizeResults) == 0 {
					return diag.FromErr(fmt.Errorf("no result found with the specified criteria: storage sizw %f", storageSize))
				}
				results = storageSizeResults
			}
		}

		if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one template found with the specified criteria: name = %s, cores = %v, ram = %v, storage_size = %v", name.(string), cores.(float64), ram.(float64), storageSize.(float64)))
		} else {
			template = results[0]
		}
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
				return fmt.Errorf("error while setting name property for image %s: %s", d.Id(), err)
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
