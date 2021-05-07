package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"strings"
)

func dataSourceTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cores": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"ram": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"storage_size": {
				Type:     schema.TypeFloat,
				Required: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	templates, _, err := client.TemplatesApi.TemplatesGet(ctx).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while fetching IonosCloud templates %s ", err)
	}

	name := d.Get("name").(string)
	cores, coresOk := d.GetOk("cores")
	ram, ramOk := d.GetOk("ram")
	storageSize, storageSizeOk := d.GetOk("storage_size")

	results := []ionoscloud.Template{}

	if templates.Items != nil {
		for _, tmp := range *templates.Items {
			if strings.Contains(strings.ToLower(*tmp.Properties.Name), strings.ToLower(name)) {
				results = append(results, tmp)
			}
		}
	}

	if coresOk {
		coresResults := []ionoscloud.Template{}
		for _, tmp := range results {
			cores := float32(cores.(float64))
			if tmp.Properties.Cores != nil && *tmp.Properties.Cores == cores {
				coresResults = append(coresResults, tmp)
			}
		}
		results = coresResults
	}

	if ramOk {
		ramResults := []ionoscloud.Template{}
		for _, tmp := range results {
			ram := float32(ram.(float64))
			if tmp.Properties.Ram != nil && *tmp.Properties.Ram == ram {
				ramResults = append(ramResults, tmp)
			}
		}
		results = ramResults
	}

	if storageSizeOk {
		storageSizeResults := []ionoscloud.Template{}
		storageSize := float32(storageSize.(float64))
		for _, tmp := range results {
			if tmp.Properties.StorageSize != nil && *tmp.Properties.StorageSize == storageSize {
				storageSizeResults = append(storageSizeResults, tmp)
			}
		}
		results = storageSizeResults
	}

	if len(results) > 1 {
		return fmt.Errorf("There is more than one template that match the search criteria ")
	}

	if len(results) == 0 {
		return fmt.Errorf("There are no templates that match the search criteria ")
	}

	if results[0].Properties.Name != nil {
		err := d.Set("name", *results[0].Properties.Name)
		if err != nil {
			return fmt.Errorf("Error while setting name property for image %s: %s", d.Id(), err)
		}
	}

	if results[0].Properties.Cores != nil {
		if err := d.Set("cores", *results[0].Properties.Cores); err != nil {
			return err
		}
	}
	if results[0].Properties.Ram != nil {
		if err := d.Set("ram", *results[0].Properties.Ram); err != nil {
			return err
		}
	}
	if results[0].Properties.StorageSize != nil {
		if err := d.Set("storage_size", *results[0].Properties.StorageSize); err != nil {
			return err
		}
	}

	d.SetId(*results[0].Id)

	return nil
}
