package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
)

func dataSourceDbassMongoTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbassMongoTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Description:      "The unique ID of the template.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Optional:         true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the template.",
				Optional:    true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using the name filter.",
				Default:     false,
				Optional:    true,
			},
			"edition": {
				Type:        schema.TypeString,
				Description: "The edition of the template (e.g. enterprise).",
				Computed:    true,
			},
			"cores": {
				Type:        schema.TypeInt,
				Description: "The number of CPU cores.",
				Computed:    true,
			},
			"ram": {
				Type:        schema.TypeInt,
				Description: "The amount of memory in GB.",
				Computed:    true,
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Description: "The amount of storage size in GB.",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDbassMongoTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MongoClient
	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	// Initial checks.
	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("name and ID cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide a template ID or name"))
	}
	retrievedTemplates, _, err := client.GetTemplates(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching dbaas mongo templates: %w", err))
	}

	var templates []mongo.TemplateResponse
	partialMatch := d.Get("partial_match").(bool)
	if retrievedTemplates.Items != nil {
		for _, retrievedTemplate := range *retrievedTemplates.Items {
			// Filter using the template ID or name.
			if (idOk && *retrievedTemplate.Id == id.(string)) ||
				(nameOk && matchesName(retrievedTemplate, name.(string), partialMatch)) {
				templates = append(templates, retrievedTemplate)
			}
		}
	}
	if templates == nil {
		return diag.FromErr(fmt.Errorf("no DBaaS Mongo Template found with the specified criteria"))
	} else if len(templates) > 1 {
		return diag.FromErr(fmt.Errorf("more than one DBaaS Mongo Template found for the specified search criteria"))
	}
	if err := dbaasService.SetMongoDBTemplateData(d, templates[0]); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// matchesName checks if a template has a specific name. allows for partial matching if partialMatch is true
func matchesName(template mongo.TemplateResponse, name string, partialMatch bool) bool {
	if template.Properties == nil || template.Properties.Name == nil {
		log.Printf("[WARN] template %s missing properties, or name", *template.Id)
		return false
	}

	if partialMatch {
		return strings.Contains(*template.Properties.Name, name)
	}

	return *template.Properties.Name == name
}
