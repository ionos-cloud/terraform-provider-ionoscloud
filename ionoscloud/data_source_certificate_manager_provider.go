package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	certSDK "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	certService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceCertificateManagerProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProviderRead,
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location of the auto-certificate provider",
			},
			"id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Description:      "The ID of the auto-certificate provider",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the auto-certificate provider",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of the certificate requester",
			},
			"server": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the certificate provider",
			},
			"external_account_binding": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key ID of the external account binding",
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	location := d.Get("location").(string)

	if idOk && nameOk {
		return utils.ToDiags(d, "ID and name cannot be provided at the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the auto-certificate provider ID or name", nil)
	}

	var provider certSDK.ProviderRead
	var apiResponse *shared.APIResponse
	var err error

	if idOk {
		id := id.(string)
		provider, apiResponse, err = client.GetProvider(ctx, id, location)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the auto-certificate provider with ID: %v, error: %s", id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		providers, apiResponse, err := client.ListProviders(ctx, location)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching auto-certificate providers: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		var results []certSDK.ProviderRead
		if providers.Items != nil {
			for _, providerItem := range providers.Items {
				if strings.EqualFold(providerItem.Properties.Name, name.(string)) {
					results = append(results, providerItem)
				}
			}
		}

		if len(results) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no auto-certificate provider found with the specified name: %v", name), nil)
		}
		if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one auto-certificate provider found with the specified name: %v", name), nil)
		}
		provider = results[0]
	}

	if err := certService.SetProviderData(d, provider); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	if provider.Properties.ExternalAccountBinding != nil {
		var externalAccountBinding []interface{}
		externalAccountBindingEntry := map[string]interface{}{}
		utils.SetPropWithNilCheck(externalAccountBindingEntry, "key_id", *provider.Properties.ExternalAccountBinding.KeyId)
		externalAccountBinding = append(externalAccountBinding, externalAccountBindingEntry)
		if err := d.Set("external_account_binding", externalAccountBinding); err != nil {
			return utils.ToDiags(d, utils.GenerateSetError("Auto-certificate provider", "external_account_binding", err).Error(), nil)
		}
	}
	return nil
}
