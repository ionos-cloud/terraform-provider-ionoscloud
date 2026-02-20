package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cr "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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

//nolint:gocyclo
func dataSourceContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).ContainerClient

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	locationValue, locationOk := d.GetOk("location")

	id := idValue.(string)
	name := nameValue.(string)
	location := locationValue.(string)

	if idOk && (nameOk || locationOk) {
		return utils.ToDiags(d, "id and name or location cannot be both specified in the same time", nil)
	}
	if !idOk && !nameOk && !locationOk {
		return utils.ToDiags(d, "please provide the registry id, name or location", nil)
	}

	var registry cr.RegistryResponse
	var apiResponse *shared.APIResponse
	var err error

	if idOk {
		/* search by ID */
		registry, apiResponse, err = client.GetRegistry(ctx, id)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the registry with ID %s: %s", id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		var results []cr.RegistryResponse

		registries, apiResponse, err := client.ListRegistries(ctx)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching container registries: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}

		results = registries.Items
		if nameOk {
			partialMatch := d.Get("partial_match").(bool)

			log.Printf("[INFO] Using data source for container registry by name with partial_match %t and name: %s", partialMatch, name)

			if registries.Items != nil && len(registries.Items) > 0 {
				var registriesByName []cr.RegistryResponse
				for _, registryItem := range registries.Items {
					if partialMatch && strings.Contains(registryItem.Properties.Name, name) ||
						!partialMatch && strings.EqualFold(registryItem.Properties.Name, name) {
						registriesByName = append(registriesByName, registryItem)
					}
				}
				if len(registriesByName) > 0 {
					results = registriesByName
				} else {
					return utils.ToDiags(d, fmt.Sprintf("no registry found with the specified criteria: name = %v", name), nil)
				}
			}
		}

		if locationOk {
			var registriesByLocation []cr.RegistryResponse
			for _, registryItem := range results {
				if strings.EqualFold(registryItem.Properties.Location, location) {
					registriesByLocation = append(registriesByLocation, registryItem)
				}
			}
			if len(registriesByLocation) > 0 {
				results = registriesByLocation
			} else {
				return utils.ToDiags(d, fmt.Sprintf("no registry found with the specified criteria: location = %v", location), nil)
			}
		}

		switch {
		case len(results) == 0:
			return utils.ToDiags(d, fmt.Sprintf("no registry found with the specified criteria: name = %s location = %s", name, location), nil)
		case len(results) > 1:
			return utils.ToDiags(d, fmt.Sprintf("more than one registry found with the specified criteria: name = %s location = %s", name, location), nil)
		default:
			registry = results[0]
		}
	}

	if err := crService.SetRegistryData(d, registry); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil

}
