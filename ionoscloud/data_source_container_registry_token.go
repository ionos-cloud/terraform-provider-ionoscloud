package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cr "github.com/ionos-cloud/sdk-go-container-registry"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/container-registry"
	"log"
	"strings"
)

func dataSourceContainerRegistryToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerRegistryTokenRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"expiry_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registry_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceContainerRegistryTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		diags := diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !nameOk {
		diags := diag.FromErr(errors.New("please provide either the token id or name"))
		return diags
	}

	var token cr.TokenResponse
	var err error

	if idOk {
		/* search by ID */
		token, _, err = client.GetToken(ctx, registryId, id)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the token with ID %s: %s", id, err))
			return diags
		}
	} else {
		var results []cr.TokenResponse

		tokens, _, err := client.ListTokens(ctx, registryId)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching registry tokens: %s", err.Error()))
			return diags
		}

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for container token by name with partial_match %t and name: %s", partialMatch, name)

		if tokens.Items != nil && len(*tokens.Items) > 0 {
			for _, tokenItem := range *tokens.Items {
				if tokenItem.Properties != nil && tokenItem.Properties.Name != nil &&
					(partialMatch && strings.Contains(*tokenItem.Properties.Name, name) ||
						!partialMatch && strings.EqualFold(*tokenItem.Properties.Name, name)) {
					results = append(results, tokenItem)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no token found with the specified name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one token found with the specified criteria: name = %s", name))
		}

		token = results[0]

	}
	if token.Id != nil {
		d.SetId(*token.Id)
	}
	if token.Properties == nil {
		return diag.FromErr(fmt.Errorf("no token properties found with the specified id = %s", *token.Id))
	}

	if err := crService.SetTokenData(d, *token.Properties); err != nil {
		return diag.FromErr(err)
	}

	return nil

}
