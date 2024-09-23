package containerregistry

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cr "github.com/ionos-cloud/sdk-go-container-registry"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func (c *Client) ListRegistries(ctx context.Context) (cr.RegistriesResponse, *cr.APIResponse, error) {
	registry, apiResponse, err := c.sdkClient.RegistriesApi.RegistriesGet(ctx).Execute()
	apiResponse.LogInfo()
	return registry, apiResponse, err
}

func (c *Client) CreateRegistry(ctx context.Context, registryInput cr.PostRegistryInput) (cr.PostRegistryOutput, *cr.APIResponse, error) {
	registry, apiResponse, err := c.sdkClient.RegistriesApi.RegistriesPost(ctx).PostRegistryInput(registryInput).Execute()
	apiResponse.LogInfo()
	return registry, apiResponse, err
}

func (c *Client) DeleteRegistry(ctx context.Context, registryId string) (*cr.APIResponse, error) {
	apiResponse, err := c.sdkClient.RegistriesApi.RegistriesDelete(ctx, registryId).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) GetRegistry(ctx context.Context, registryId string) (cr.RegistryResponse, *cr.APIResponse, error) {
	registries, apiResponse, err := c.sdkClient.RegistriesApi.RegistriesFindById(ctx, registryId).Execute()
	apiResponse.LogInfo()
	return registries, apiResponse, err
}

// IsRegistryDeleted checks whether the container registry is deleted or not
func (c *Client) IsRegistryDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	ID := d.Id()
	_, resp, err := c.GetRegistry(ctx, ID)
	if resp.HttpNotFound() {
		return true, nil
	}
	return false, err
}

// IsRegistryReady checks whether the container registry is in a ready state or not
func (c *Client) IsRegistryReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	ID := d.Id()
	creg, _, err := c.GetRegistry(ctx, ID)
	if err != nil {
		return true, fmt.Errorf("status check failed for container registry creg with ID: %v, error: %w", ID, err)
	}

	if creg.Metadata == nil || creg.Metadata.State == nil {
		return false, fmt.Errorf("metadata or state is empty for container registry with ID: %v", ID)
	}

	log.Printf("[INFO] state of the container registry with ID: %v is: %s ", ID, *creg.Metadata.State)
	return strings.EqualFold(*creg.Metadata.State, "RUNNING"), nil
}

func (c *Client) PatchRegistry(ctx context.Context, registryId string, registryInput cr.PatchRegistryInput) (cr.RegistryResponse, *cr.APIResponse, error) {
	registries, apiResponse, err := c.sdkClient.RegistriesApi.RegistriesPatch(ctx, registryId).PatchRegistryInput(registryInput).Execute()
	apiResponse.LogInfo()
	return registries, apiResponse, err
}

func (c *Client) PutRegistry(ctx context.Context, registryId string, registryInput cr.PutRegistryInput) (cr.PutRegistryOutput, *cr.APIResponse, error) {
	registries, apiResponse, err := c.sdkClient.RegistriesApi.RegistriesPut(ctx, registryId).PutRegistryInput(registryInput).Execute()
	apiResponse.LogInfo()
	return registries, apiResponse, err
}

func (c *Client) DeleteRepositories(ctx context.Context, registryId, repositoryId string) (*cr.APIResponse, error) {
	apiResponse, err := c.sdkClient.RepositoriesApi.RegistriesRepositoriesDelete(ctx, registryId, repositoryId).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) ListTokens(ctx context.Context, registryId string) (cr.TokensResponse, *cr.APIResponse, error) {
	tokens, apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensGet(ctx, registryId).Execute()
	apiResponse.LogInfo()
	return tokens, apiResponse, err

}

func (c *Client) CreateTokens(ctx context.Context, registryId string, tokenInput cr.PostTokenInput) (cr.PostTokenOutput, *cr.APIResponse, error) {
	token, apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensPost(ctx, registryId).PostTokenInput(tokenInput).Execute()
	apiResponse.LogInfo()
	return token, apiResponse, err

}

func (c *Client) DeleteToken(ctx context.Context, registryId, tokenId string) (*cr.APIResponse, error) {
	apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensDelete(ctx, registryId, tokenId).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) GetToken(ctx context.Context, registryId, tokenId string) (cr.TokenResponse, *cr.APIResponse, error) {
	token, apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensFindById(ctx, registryId, tokenId).Execute()
	apiResponse.LogInfo()
	return token, apiResponse, err

}

func (c *Client) PatchToken(ctx context.Context, registryId, tokenId string, tokenInput cr.PatchTokenInput) (cr.TokenResponse, *cr.APIResponse, error) {
	token, apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensPatch(ctx, registryId, tokenId).PatchTokenInput(tokenInput).Execute()
	apiResponse.LogInfo()
	return token, apiResponse, err

}

func (c *Client) PutToken(ctx context.Context, registryId, tokenId string, tokenInput cr.PutTokenInput) (cr.PutTokenOutput, *cr.APIResponse, error) {
	token, apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensPut(ctx, registryId, tokenId).PutTokenInput(tokenInput).Execute()
	apiResponse.LogInfo()
	return token, apiResponse, err

}

// GetRegistryDataCreate get registry data for create
func GetRegistryDataCreate(d *schema.ResourceData) (*cr.PostRegistryInput, error) {

	registry := cr.PostRegistryInput{
		Properties: &cr.PostRegistryProperties{},
	}

	if _, ok := d.GetOk("garbage_collection_schedule"); ok {
		registry.Properties.GarbageCollectionSchedule = GetWeeklySchedule(d, "garbage_collection_schedule")
	}

	if location, ok := d.GetOk("location"); ok {
		location := location.(string)
		registry.Properties.Location = &location
	}

	if name, ok := d.GetOk("name"); ok {
		name := name.(string)
		registry.Properties.Name = &name
	}

	if v, ok := d.GetOk("api_subnet_allow_list"); ok {
		raw := v.([]interface{})
		ips := make([]string, len(raw))
		err := utils.DecodeInterfaceToStruct(raw, ips)
		if err != nil {
			return nil, err
		}
		if len(ips) > 0 {
			registry.Properties.ApiSubnetAllowList = &ips
		}
	}

	return &registry, nil
}

// GetRegistryDataUpdate get registry data for update
func GetRegistryDataUpdate(d *schema.ResourceData) (*cr.PatchRegistryInput, error) {

	registry := cr.PatchRegistryInput{}

	if _, ok := d.GetOk("garbage_collection_schedule"); ok {
		registry.GarbageCollectionSchedule = GetWeeklySchedule(d, "garbage_collection_schedule")
	}

	// When api_subnet_allow_list = [] in the TF plan:
	// GetOk => _, false
	// GetOkExists => _, true
	// We don't want to ignore the cases in which we modify the api_subnet_allow_list to an empty list, so
	// here we need to use 'GetOkExists'.
	if v, ok := d.GetOkExists("api_subnet_allow_list"); ok { //nolint:staticcheck
		raw := v.([]interface{})
		ips := make([]string, len(raw))
		err := utils.DecodeInterfaceToStruct(raw, ips)
		if err != nil {
			return nil, err
		}
		registry.ApiSubnetAllowList = &ips
	}

	return &registry, nil
}

func GetWeeklySchedule(d *schema.ResourceData, property string) *cr.WeeklySchedule {
	var weeklySchedule cr.WeeklySchedule

	if days, ok := d.GetOk(fmt.Sprintf("%v.0.days", property)); ok {
		var daysToAdd []cr.Day
		for _, day := range days.([]interface{}) {
			daysToAdd = append(daysToAdd, cr.Day(day.(string)))
		}
		weeklySchedule.Days = &daysToAdd
	}

	if timeField, ok := d.GetOk(fmt.Sprintf("%v.0.time", property)); ok {
		timeStr := timeField.(string)
		weeklySchedule.Time = &timeStr
	}

	return &weeklySchedule
}

func SetRegistryData(d *schema.ResourceData, registry cr.RegistryResponse) error {

	resourceName := "registry"

	if registry.Id != nil {
		d.SetId(*registry.Id)
	}

	if registry.Properties.GarbageCollectionSchedule != nil {
		var schedule []interface{}
		scheduleEntry := SetWeeklySchedule(*registry.Properties.GarbageCollectionSchedule)
		schedule = append(schedule, scheduleEntry)
		if err := d.Set("garbage_collection_schedule", schedule); err != nil {
			return utils.GenerateSetError(resourceName, "garbage_collection_schedule", err)
		}
	}

	if registry.Properties.Hostname != nil {
		if err := d.Set("hostname", *registry.Properties.Hostname); err != nil {
			return utils.GenerateSetError(resourceName, "hostname", err)
		}
	}

	if registry.Properties.Location != nil {
		if err := d.Set("location", *registry.Properties.Location); err != nil {
			return utils.GenerateSetError(resourceName, "location", err)
		}
	}

	if registry.Properties.Name != nil {
		if err := d.Set("name", *registry.Properties.Name); err != nil {
			return utils.GenerateSetError(resourceName, "name", err)
		}
	}

	if registry.Properties.StorageUsage != nil {
		var storage []interface{}
		storageEntry := SetStorageUsage(*registry.Properties.StorageUsage)
		storage = append(storage, storageEntry)
		if err := d.Set("storage_usage", storage); err != nil {
			return utils.GenerateSetError(resourceName, "storage_usage", err)
		}
	}

	if registry.Properties.ApiSubnetAllowList != nil {
		if err := d.Set("api_subnet_allow_list", *registry.Properties.ApiSubnetAllowList); err != nil {
			return fmt.Errorf("error setting api_subnet_allow_list %w", err)
		}
	}

	if registry.Properties.Features != nil {

		registryFeatures := map[string]any{}
		utils.SetPropWithNilCheck(registryFeatures, "vulnerability_scanning", registry.Properties.Features.VulnerabilityScanning.Enabled)

		if err := d.Set("features", []map[string]any{registryFeatures}); err != nil {
			return utils.GenerateSetError(resourceName, "features", err)
		}
	}

	return nil
}

func SetWeeklySchedule(weeklySchedule cr.WeeklySchedule) map[string]interface{} {

	schedule := map[string]interface{}{}

	utils.SetPropWithNilCheck(schedule, "time", weeklySchedule.Time)
	utils.SetPropWithNilCheck(schedule, "days", weeklySchedule.Days)

	return schedule
}

func SetStorageUsage(storageUsage cr.StorageUsage) map[string]interface{} {

	storage := map[string]interface{}{}

	utils.SetPropWithNilCheck(storage, "bytes", storageUsage.Bytes)
	if storageUsage.UpdatedAt != nil {
		utils.SetPropWithNilCheck(storage, "updated_at", storageUsage.UpdatedAt.String())
	}

	return storage
}

func GetTokenDataCreate(d *schema.ResourceData) (*cr.PostTokenInput, error) {

	token := cr.PostTokenInput{
		Properties: &cr.PostTokenProperties{},
	}

	if expiryDate, ok := d.GetOk("expiry_date"); ok {
		expiryDate, err := convertToIonosTime(expiryDate.(string))
		if err != nil {
			return nil, err
		}
		token.Properties.ExpiryDate = expiryDate
	}

	if name, ok := d.GetOk("name"); ok {
		name := name.(string)
		token.Properties.Name = &name
	}

	if _, ok := d.GetOk("scopes"); ok {
		token.Properties.Scopes = GetScopes(d)
	}

	if status, ok := d.GetOk("status"); ok {
		status := status.(string)
		token.Properties.Status = &status
	}

	return &token, nil
}

func GetTokenDataUpdate(d *schema.ResourceData) (*cr.PatchTokenInput, error) {

	token := cr.PatchTokenInput{}

	if expiryDate, ok := d.GetOk("expiry_date"); ok {
		expiryDate, err := convertToIonosTime(expiryDate.(string))
		if err != nil {
			return nil, err
		}
		token.ExpiryDate = expiryDate
	} else {
		token.ExpiryDate = nil
	}

	if _, ok := d.GetOk("scopes"); ok {
		token.Scopes = GetScopes(d)
	}

	if status, ok := d.GetOk("status"); ok {
		status := status.(string)
		token.Status = &status
	}

	return &token, nil
}

func GetScopes(d *schema.ResourceData) *[]cr.Scope {
	scopes := make([]cr.Scope, 0)

	if scopeValue, ok := d.GetOk("scopes"); ok {
		scopeValue := scopeValue.([]interface{})
		if scopeValue != nil {
			for _, item := range scopeValue {

				scopeContent := item.(map[string]interface{})
				connection := cr.Scope{}

				if actions, ok := scopeContent["actions"]; ok {
					actions := actions.([]interface{})
					var actionsToAdd []string
					if actions != nil && len(actions) > 0 {
						for _, action := range actions {
							actionsToAdd = append(actionsToAdd, action.(string))
						}
					}
					connection.Actions = &actionsToAdd
				}

				if name, ok := scopeContent["name"]; ok {
					name := name.(string)
					connection.Name = &name
				}

				if scopeType, ok := scopeContent["type"]; ok {
					scopeType := scopeType.(string)
					connection.Type = &scopeType
				}

				scopes = append(scopes, connection)
			}
		}

	}

	return &scopes

}

func SetTokenData(d *schema.ResourceData, tokenProps cr.TokenProperties) error {

	regToken := "registry token "

	if tokenProps.Credentials != nil {
		var credentials []interface{}
		credentialsEntry := SetCredentials(*tokenProps.Credentials)
		credentials = append(credentials, credentialsEntry)
		if err := d.Set("credentials", credentials); err != nil {
			return utils.GenerateSetError(regToken, "credentials", err)
		}
	}

	if tokenProps.ExpiryDate != nil {
		timeValue := (*tokenProps.ExpiryDate).Time
		if err := d.Set("expiry_date", timeValue.String()); err != nil {
			return utils.GenerateSetError(regToken, "expiry_date", err)
		}
	}

	if tokenProps.Name != nil {
		if err := d.Set("name", *tokenProps.Name); err != nil {
			return utils.GenerateSetError(regToken, "name", err)
		}
	}

	if tokenProps.Scopes != nil {
		scopes := SetScopes(*tokenProps.Scopes)
		if err := d.Set("scopes", scopes); err != nil {
			return utils.GenerateSetError(regToken, "scopes", err)
		}
	}

	if tokenProps.Status != nil {
		if err := d.Set("status", *tokenProps.Status); err != nil {
			return utils.GenerateSetError(regToken, "status", err)
		}
	}

	return nil
}

func SetCredentials(credentials cr.Credentials) map[string]interface{} {

	credentialsEntry := map[string]interface{}{}

	utils.SetPropWithNilCheck(credentialsEntry, "username", credentials.Username)

	return credentialsEntry
}

func SetScopes(scopes []cr.Scope) []interface{} {

	var tokenScopes []interface{}
	for _, scope := range scopes {
		scopeEntry := make(map[string]interface{})

		if scope.Actions != nil {
			scopeEntry["actions"] = *scope.Actions
		}

		if scope.Name != nil {
			scopeEntry["name"] = *scope.Name
		}

		if scope.Type != nil {
			scopeEntry["type"] = *scope.Type
		}

		tokenScopes = append(tokenScopes, scopeEntry)
	}

	return tokenScopes

}

// GetRegistryFeatures returns the container registry features retrieved from the configuration
// It will also return a list of warnings related to attributes which should be set explicitly
func GetRegistryFeatures(d *schema.ResourceData) (*cr.RegistryFeatures, diag.Diagnostics) {

	crfWarnings := struct {
		warnCRVScanningOmitted diag.Diagnostic
		warnCRVScanningOff     diag.Diagnostic
		warnCRVScanningOn      diag.Diagnostic
	}{
		warnCRVScanningOmitted: diag.Diagnostic{Severity: diag.Warning,
			Summary: "'vulnerability_scanning' is omitted from the config. CR Vulnerability Scanning has been enabled by default.",
			Detail: "Container Registry Vulnerability Scanning is a paid security feature which is enabled by default.\n" +
				"If you do not wish to enable it, ensure 'vulnerability_scanning' is set to false when creating the resource.\n" +
				"Once activated, it cannot be deactivated afterwards for this CR instance.\n" +
				"More details about CR Vulnerability Scanning: https://docs.ionos.com/cloud/managed-services/container-registry/dcd-how-tos/enable-vulnerability-scanning\n" +
				"Price list is available under the 'Data Platform' section for your selected region at: https://docs.ionos.com/support/general-information/price-list",
		},
		warnCRVScanningOff: diag.Diagnostic{Severity: diag.Warning,
			Summary: "'vulnerability_scanning' is disabled for this Container Registry.",
			Detail: "Container Registry Vulnerability Scanning is a paid security feature which we recommend enabling.\n" +
				"More details about CR Vulnerability Scanning: https://docs.ionos.com/cloud/managed-services/container-registry/dcd-how-tos/enable-vulnerability-scanning\n" +
				"Price list is available under the 'Data Platform' section for your selected region at: https://docs.ionos.com/support/general-information/price-list",
		},
		warnCRVScanningOn: diag.Diagnostic{Severity: diag.Warning,
			Summary: "'vulnerability_scanning' has been enabled for this Container Registry.",
			Detail: "Container Registry Vulnerability Scanning is a paid security feature.\n" +
				"More details about CR Vulnerability Scanning: https://docs.ionos.com/cloud/managed-services/container-registry/dcd-how-tos/enable-vulnerability-scanning\n" +
				"Price list is available under the 'Data Platform' section for your selected region at: https://docs.ionos.com/support/general-information/price-list",
		},
	}

	registryFeatures := cr.NewRegistryFeatures()
	var warnings = diag.Diagnostics{crfWarnings.warnCRVScanningOmitted}
	registryFeatures.VulnerabilityScanning = cr.NewFeatureVulnerabilityScanning(true)
	if vulnerabilityScanning, ok := d.GetOkExists("features.0.vulnerability_scanning"); ok { //nolint:staticcheck
		vulnerabilityScanning := vulnerabilityScanning.(bool)
		registryFeatures.VulnerabilityScanning.Enabled = &vulnerabilityScanning
		warnings = diag.Diagnostics{crfWarnings.warnCRVScanningOff}

		if vulnerabilityScanning {
			warnings = diag.Diagnostics{crfWarnings.warnCRVScanningOn}
		}
	}

	return registryFeatures, warnings

}

func convertToIonosTime(targetTime string) (*cr.IonosTime, error) {
	var ionosTime cr.IonosTime
	var convertedTime time.Time
	var err error

	// targetTime might have time zone offset layout (+0000 UTC)
	if convertedTime, err = time.Parse(constant.DatetimeTZOffsetLayout, targetTime); err != nil {
		if convertedTime, err = time.Parse(constant.DatetimeZLayout, targetTime); err != nil {
			return nil, fmt.Errorf("an error occurred while converting from IonosTime string to time.Time: %w", err)
		}
	}
	ionosTime.Time = convertedTime
	return &ionosTime, nil
}
