package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	certSDK "github.com/ionos-cloud/sdk-go-cert-manager"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	certService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func dataSourceCertificateManagerProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProviderRead,
		Schema: map[string]*schema.Schema{
			"location": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The location of the auto-certificate provider",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(constant.Locations, false)),
			},
			"id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Description:      "The ID of the auto-certificate provider",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
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
	client := meta.(services.SdkBundle).CertManagerClient
	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	location := d.Get("location").(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be provided at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the auto-certificate provider ID or name"))
	}

	var provider certSDK.ProviderRead
	var err error

	if idOk {
		id := id.(string)
		provider, _, err = client.GetProvider(ctx, id, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the auto-certificate provider with ID: %v, error: %w", id, err))
		}
	} else {
		providers, _, err := client.ListProviders(ctx, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching auto-certificate providers: %w", err))
		}
		var results []certSDK.ProviderRead
		if providers.Items != nil {
			for _, providerItem := range *providers.Items {
				if providerItem.Properties != nil && providerItem.Properties.Name != nil && strings.EqualFold(*providerItem.Properties.Name, name.(string)) {
					results = append(results, providerItem)
				}
			}
		}

		if len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no auto-certificate provider found with the specified name: %v", name))
		}
		if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one auto-certificate provider found with the specified name: %v", name))
		}
		provider = results[0]
	}

	if err := certService.SetProviderData(d, provider); err != nil {
		return diag.FromErr(err)
	}

	if provider.Properties.ExternalAccountBinding != nil {
		var externalAccountBinding []interface{}
		externalAccountBindingEntry := map[string]interface{}{}
		utils.SetPropWithNilCheck(externalAccountBindingEntry, "key_id", *provider.Properties.ExternalAccountBinding.KeyId)
		externalAccountBinding = append(externalAccountBinding, externalAccountBindingEntry)
		if err := d.Set("external_account_binding", externalAccountBinding); err != nil {
			return diag.FromErr(utils.GenerateSetError("Auto-certificate provider", "external_account_binding", err))
		}
	}
	return nil
}
