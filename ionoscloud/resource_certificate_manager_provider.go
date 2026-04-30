package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func resourceCertificateManagerProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: providerCreate,
		ReadContext:   providerRead,
		UpdateContext: providerUpdate,
		DeleteContext: providerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: providerImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the certificate provider",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The email address of the certificate requester",
			},
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The URL of the certificate provider",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The location of the certificate provider",
			},
			"external_account_binding": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The key ID of the external account binding",
						},
						"key_secret": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Sensitive:   true,
							Description: "The secret of the external account binding",
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func providerCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	location := d.Get("location").(string)

	providerCreateData := cert.GetProviderDataCreate(d)
	response, apiResponse, err := client.CreateProvider(ctx, *providerCreateData, location)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating an auto-certificate provider: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	providerID := response.Id
	d.SetId(providerID)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsProviderReady)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while checking the creation status for the auto-certificate provider with ID: %v, error: %w", providerID, err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String()})
	}
	if err := cert.SetProviderData(d, response); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	return nil
}

func providerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	providerID := d.Id()
	location := d.Get("location").(string)
	provider, apiResponse, err := client.GetProvider(ctx, providerID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while fetching auto-certificate provider with ID: %v, error: %w", providerID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	log.Printf("[INFO] Successfully retrieved auto-certificate provider with ID: %v, provider info: %+v", providerID, provider)
	if err := cert.SetProviderData(d, provider); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	return nil
}

func providerUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	providerID := d.Id()
	location := d.Get("location").(string)

	providerUpdateData := cert.GetProviderDataUpdate(d)
	provider, apiResponse, err := client.UpdateProvider(ctx, providerID, location, *providerUpdateData)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating auto-certificate provider with ID: %v, error: %w", providerID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	if err = utils.WaitForResourceToBeReady(ctx, d, client.IsProviderReady); err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while checking the update status for the auto-certificate provider with ID: %v, error: %w", providerID, err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String()})
	}
	if err := cert.SetProviderData(d, provider); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	return nil
}

func providerDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	providerID := d.Id()
	location := d.Get("location").(string)
	apiResponse, err := client.DeleteProvider(ctx, providerID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while deleting auto-certificate provider with ID: %v, error: %w", providerID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsProviderDeleted)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("deletion check failed for auto-certificate provider with ID: %v, error: %w", providerID, err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String()})
	}
	return nil
}

func providerImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, diagutil.ToError(d, fmt.Errorf("invalid import, expected ID in the format: '<location>:<provider_id>'"), nil)
	}
	location := parts[0]
	providerID := parts[1]
	provider, apiResponse, err := client.GetProvider(ctx, providerID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("auto-certificate provider with ID: %v does not exist", providerID), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while trying to import auto-certificate provider with ID: %v, error: %w", providerID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	log.Printf("[INFO] auto-certificate provider found: %+v", provider)
	if err := d.Set("location", location); err != nil {
		return nil, utils.GenerateSetError("Auto-certificate provider", "location", err)
	}
	if err := cert.SetProviderData(d, provider); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}
