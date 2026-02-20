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

func dataSourceContainerRegistryToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerRegistryTokenRead,
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
	client := meta.(bundleclient.SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return utils.ToDiags(d, "id and name cannot be both specified in the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the token id or name", nil)
	}

	var token cr.TokenResponse
	var apiResponse *shared.APIResponse
	var err error

	if idOk {
		/* search by ID */
		token, apiResponse, err = client.GetToken(ctx, registryId, id)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the token with ID %s: %s", id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		var results []cr.TokenResponse

		tokens, apiResponse, err := client.ListTokens(ctx, registryId)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching registry tokens: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for container token by name with partial_match %t and name: %s", partialMatch, name)

		if tokens.Items != nil && len(tokens.Items) > 0 {
			for _, tokenItem := range tokens.Items {
				if partialMatch && strings.Contains(tokenItem.Properties.Name, name) ||
					!partialMatch && strings.EqualFold(tokenItem.Properties.Name, name) {
					results = append(results, tokenItem)
				}
			}
		}

		if len(results) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no token found with the specified name = %s", name), nil)
		} else if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one token found with the specified criteria: name = %s", name), nil)
		}

		token = results[0]

	}
	if token.Id != nil {
		d.SetId(*token.Id)
	}

	if err := crService.SetTokenData(d, token.Properties); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	var credentials []any
	credentialsEntry := crService.SetCredentials(token.Properties.Credentials)
	credentials = append(credentials, credentialsEntry)
	if err := d.Set("credentials", credentials); err != nil {
		return utils.ToDiags(d, utils.GenerateSetError("token", "credentials", err).Error(), nil)
	}
	return nil

}
