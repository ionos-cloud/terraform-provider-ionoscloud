package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cr "github.com/ionos-cloud/sdk-go-container-registry"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
)

func dataSourceContainerRegistry() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerRegistryRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_subnet_allow_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "The subnet CIDRs that are allowed to connect to the registry. Specify 'a.b.c.d/32' for an individual IP address. __Note__: If this list is empty or not set, there are no restrictions.",
			},
			"garbage_collection_schedule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"days": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_window": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"days": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"storage_usage": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"features": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vulnerability_scanning": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ContainerClient

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	locationValue, locationOk := d.GetOk("location")

	id := idValue.(string)
	name := nameValue.(string)
	location := locationValue.(string)

	if idOk && (nameOk || locationOk) {
		diags := diag.FromErr(errors.New("id and name or location cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !nameOk && !locationOk {
		diags := diag.FromErr(errors.New("please provide the registry id, name or location"))
		return diags
	}

	var registry cr.RegistryResponse
	var err error

	if idOk {
		/* search by ID */
		registry, _, err = client.GetRegistry(ctx, id)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the registry with ID %s: %w", id, err))
			return diags
		}
	} else {
		var results []cr.RegistryResponse

		registries, _, err := client.ListRegistries(ctx)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching container registries: %w", err))
			return diags
		}

		results = *registries.Items
		if nameOk {
			partialMatch := d.Get("partial_match").(bool)

			log.Printf("[INFO] Using data source for container registry by name with partial_match %t and name: %s", partialMatch, name)

			if registries.Items != nil && len(*registries.Items) > 0 {
				var registriesByName []cr.RegistryResponse
				for _, registryItem := range *registries.Items {
					if registryItem.Properties != nil && registryItem.Properties.Name != nil &&
						(partialMatch && strings.Contains(*registryItem.Properties.Name, name) ||
							!partialMatch && strings.EqualFold(*registryItem.Properties.Name, name)) {
						registriesByName = append(registriesByName, registryItem)
					}
				}
				if registriesByName != nil && len(registriesByName) > 0 {
					results = registriesByName
				} else {
					return diag.FromErr(fmt.Errorf("no registry found with the specified criteria: name = %v", name))
				}
			}
		}

		if locationOk {
			var registriesByLocation []cr.RegistryResponse
			for _, registryItem := range results {
				if registryItem.Properties != nil && registryItem.Properties.Name != nil && strings.EqualFold(*registryItem.Properties.Location, location) {
					registriesByLocation = append(registriesByLocation, registryItem)
				}
			}
			if registriesByLocation != nil && len(registriesByLocation) > 0 {
				results = registriesByLocation
			} else {
				return diag.FromErr(fmt.Errorf("no registry found with the specified criteria: location = %v", location))
			}
		}

		switch {
		case len(results) == 0:
			return diag.FromErr(fmt.Errorf("no registry found with the specified criteria: name = %s location = %s", name, location))
		case len(results) > 1:
			return diag.FromErr(fmt.Errorf("more than one registry found with the specified criteria: name = %s location = %s", name, location))
		default:
			registry = results[0]
		}
	}

	if err := crService.SetRegistryData(d, registry); err != nil {
		return diag.FromErr(err)
	}

	return nil

}
