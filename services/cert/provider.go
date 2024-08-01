package cert

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	certmanager "github.com/ionos-cloud/sdk-go-cert-manager"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

var locationToURL = map[string]string{
	"de/fra": "https://certificate-manager.de-fra.ionos.com",
}

// modifyConfigURL modifies the URL inside the client configuration.
// This function is required in order to make requests to different endpoints based on location.
func (c *Client) modifyConfigURL(location string) {
	clientConfig := c.sdkClient.GetConfig()
	clientConfig.Servers = certmanager.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
}

//nolint:golint
func (c *Client) GetProvider(ctx context.Context, providerID, location string) (certmanager.ProviderRead, *certmanager.APIResponse, error) {
	c.modifyConfigURL(location)
	provider, apiResponse, err := c.sdkClient.ProviderApi.ProvidersFindById(ctx, providerID).Execute()
	apiResponse.LogInfo()
	return provider, apiResponse, err
}

//nolint:golint
func (c *Client) ListProviders(ctx context.Context, location string) (certmanager.ProviderReadList, *certmanager.APIResponse, error) {
	c.modifyConfigURL(location)
	providers, apiResponse, err := c.sdkClient.ProviderApi.ProvidersGet(ctx).Execute()
	apiResponse.LogInfo()
	return providers, apiResponse, err
}

//nolint:golint
func (c *Client) CreateProvider(ctx context.Context, providerPostData certmanager.ProviderCreate, location string) (certmanager.ProviderRead, *certmanager.APIResponse, error) {
	c.modifyConfigURL(location)
	provider, apiResponse, err := c.sdkClient.ProviderApi.ProvidersPost(ctx).ProviderCreate(providerPostData).Execute()
	apiResponse.LogInfo()
	return provider, apiResponse, err
}

//nolint:golint
func (c *Client) UpdateProvider(ctx context.Context, providerID, location string, providerPatchData certmanager.ProviderPatch) (certmanager.ProviderRead, *certmanager.APIResponse, error) {
	c.modifyConfigURL(location)
	provider, apiResponse, err := c.sdkClient.ProviderApi.ProvidersPatch(ctx, providerID).ProviderPatch(providerPatchData).Execute()
	apiResponse.LogInfo()
	return provider, apiResponse, err
}

//nolint:golint
func (c *Client) DeleteProvider(ctx context.Context, providerID, location string) (*certmanager.APIResponse, error) {
	c.modifyConfigURL(location)
	apiResponse, err := c.sdkClient.ProviderApi.ProvidersDelete(ctx, providerID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

//nolint:golint
func (c *Client) IsProviderReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	providerID := d.Id()
	location := d.Get("location").(string)
	provider, _, err := c.GetProvider(ctx, providerID, location)
	if err != nil {
		return false, fmt.Errorf("error checking certificate manager provider status: %w", err)
	}
	if provider.Metadata == nil || provider.Metadata.State == nil {
		return false, fmt.Errorf("metadata or state is empty for certificate manager provider with ID: %v", d.Id())
	}
	if utils.IsStateFailed(*provider.Metadata.State) {
		return false, fmt.Errorf("error while checking if auto-certificate provider is ready, provider ID: %v, state: %v", *provider.Id, *provider.Metadata.State)
	}
	return strings.EqualFold(*provider.Metadata.State, constant.Available), nil
}

//nolint:golint
func (c *Client) IsProviderDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	providerID := d.Id()
	location := d.Get("location").(string)
	provider, apiResponse, err := c.GetProvider(ctx, providerID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("error while checking deletion status for certificate manager provider with ID: %v, error: %w", d.Id(), err)
	}
	if provider.Metadata != nil && provider.Metadata.State != nil && utils.IsStateFailed(*provider.Metadata.State) {
		return false, fmt.Errorf("error while checking if auto-certificate provider is deleted properly, provider ID: %v, state: %v", *provider.Id, *provider.Metadata.State)
	}
	return false, nil
}

//nolint:golint
func GetProviderDataCreate(d *schema.ResourceData) *certmanager.ProviderCreate {
	provider := certmanager.ProviderCreate{
		Properties: &certmanager.Provider{},
	}

	name := d.Get("name").(string)
	provider.Properties.Name = &name
	email := d.Get("email").(string)
	provider.Properties.Email = &email
	server := d.Get("server").(string)
	provider.Properties.Server = &server
	if _, ok := d.GetOk("external_account_binding"); ok {
		keyId := d.Get("external_account_binding.0.key_id").(string)
		keySecret := d.Get("external_account_binding.0.key_secret").(string)
		provider.Properties.ExternalAccountBinding = &certmanager.ProviderExternalAccountBinding{
			KeyId:     &keyId,
			KeySecret: &keySecret,
		}
	}
	return &provider
}

//nolint:golint
func SetProviderData(d *schema.ResourceData, provider certmanager.ProviderRead) error {
	resourceName := "Auto-certificate provider"
	if provider.Id != nil {
		d.SetId(*provider.Id)
	}
	if provider.Properties == nil {
		return fmt.Errorf("response properties should not be empty for auto-certificate provider with ID: %v", *provider.Id)
	}
	if provider.Properties.Name != nil {
		if err := d.Set("name", *provider.Properties.Name); err != nil {
			return utils.GenerateSetError(resourceName, "name", err)
		}
	}
	if provider.Properties.Email != nil {
		if err := d.Set("email", *provider.Properties.Email); err != nil {
			return utils.GenerateSetError(resourceName, "email", err)
		}
	}
	if provider.Properties.Server != nil {
		if err := d.Set("server", *provider.Properties.Server); err != nil {
			return utils.GenerateSetError(resourceName, "server", err)
		}
	}
	return nil
}

//nolint:golint
func GetProviderDataUpdate(d *schema.ResourceData) *certmanager.ProviderPatch {
	provider := certmanager.ProviderPatch{
		Properties: &certmanager.PatchName{},
	}

	if d.HasChange("name") {
		_, newValue := d.GetChange("name")
		newValueStr := newValue.(string)
		provider.Properties.Name = &newValueStr
	}
	return &provider
}
