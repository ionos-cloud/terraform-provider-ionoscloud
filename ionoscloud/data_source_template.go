package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strconv"
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

	name, nameOk := d.GetOk("name")
	cores, coresOk := d.GetOk("cores")
	ram, ramOk := d.GetOk("ram")
	storageSize, storageSizeOk := d.GetOk("storage_size")

	var template ionoscloud.Template
	request := client.TemplatesApi.TemplatesGet(ctx).Depth(1)

	if nameOk {
		request = request.Filter("name", name.(string))
	}

	if coresOk {
		request = request.Filter("cores", strconv.Itoa(int(cores.(float64))))
	}

	if ramOk {
		request = request.Filter("ram", strconv.Itoa(int(ram.(float64))))
	}

	if storageSizeOk {
		request = request.Filter("storageSize", strconv.Itoa(int(storageSize.(float64))))
	}

	templates, apiResponse, err := request.Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching IonosCloud templates %s ", err))
		return diags
	}

	if templates.Items != nil && len(*templates.Items) > 0 {
		template = (*templates.Items)[len(*templates.Items)-1]
		log.Printf("[INFO] %v templates found matching the search criteria. Getting the latest template from the list %v", len(*templates.Items), *template.Id)
	} else {
		return diag.FromErr(fmt.Errorf("no template found with the specified criteria: name %s, cores %s, ram %s, storageSize %s", name.(string), strconv.Itoa(int(cores.(float64))), strconv.Itoa(int(ram.(float64))), strconv.Itoa(int(storageSize.(float64)))))
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
