package containerregistry

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cr "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// ListRegistries lists all container registries
func (c *Client) ListRegistries(ctx context.Context) (cr.RegistriesResponse, *shared.APIResponse, error) {
	registry, apiResponse, err := c.sdkClient.RegistriesApi.RegistriesGet(ctx).Execute()
	apiResponse.LogInfo()
	return registry, apiResponse, err
}

// CreateRegistry creates a container registry
func (c *Client) CreateRegistry(ctx context.Context, registryInput cr.PostRegistryInput) (cr.PostRegistryOutput, *shared.APIResponse, error) {
	registry, apiResponse, err := c.sdkClient.RegistriesApi.RegistriesPost(ctx).PostRegistryInput(registryInput).Execute()
	apiResponse.LogInfo()
	return registry, apiResponse, err
}

// DeleteRegistry deletes a container registry
func (c *Client) DeleteRegistry(ctx context.Context, registryID string) (*shared.APIResponse, error) {
	apiResponse, err := c.sdkClient.RegistriesApi.RegistriesDelete(ctx, registryID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// GetRegistry gets a container registry
func (c *Client) GetRegistry(ctx context.Context, registryID string) (cr.RegistryResponse, *shared.APIResponse, error) {
	registries, apiResponse, err := c.sdkClient.RegistriesApi.RegistriesFindById(ctx, registryID).Execute()
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
	log.Printf("[INFO] state of the container registry with ID: %v is: %s ", ID, creg.Metadata.State)
	if utils.IsStateFailed(creg.Metadata.State) {
		return false, fmt.Errorf("container registry  %s is in a failed state", d.Id())
	}
	return strings.EqualFold(creg.Metadata.State, "RUNNING"), nil
}

// PatchRegistry patches a container registry
func (c *Client) PatchRegistry(ctx context.Context, registryID string, registryInput cr.PatchRegistryInput) (cr.RegistryResponse, *shared.APIResponse, error) {
	registries, apiResponse, err := c.sdkClient.RegistriesApi.RegistriesPatch(ctx, registryID).PatchRegistryInput(registryInput).Execute()
	apiResponse.LogInfo()
	return registries, apiResponse, err
}

// PutRegistry puts a container registry
func (c *Client) PutRegistry(ctx context.Context, registryID string, registryInput cr.PutRegistryInput) (cr.PutRegistryOutput, *shared.APIResponse, error) {
	registries, apiResponse, err := c.sdkClient.RegistriesApi.RegistriesPut(ctx, registryID).PutRegistryInput(registryInput).Execute()
	apiResponse.LogInfo()
	return registries, apiResponse, err
}

// DeleteRepositories deletes all repositories in a container registry
func (c *Client) DeleteRepositories(ctx context.Context, registryID, repositoryID string) (*shared.APIResponse, error) {
	apiResponse, err := c.sdkClient.RepositoriesApi.RegistriesRepositoriesDelete(ctx, registryID, repositoryID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// ListTokens lists all tokens in a container registry
func (c *Client) ListTokens(ctx context.Context, registryID string) (cr.TokensResponse, *shared.APIResponse, error) {
	tokens, apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensGet(ctx, registryID).Execute()
	apiResponse.LogInfo()
	return tokens, apiResponse, err
}

// CreateToken creates a token in a container registry
func (c *Client) CreateToken(ctx context.Context, registryID string, tokenInput cr.PostTokenInput) (cr.PostTokenOutput, *shared.APIResponse, error) {
	token, apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensPost(ctx, registryID).PostTokenInput(tokenInput).Execute()
	apiResponse.LogInfo()
	return token, apiResponse, err

}

// DeleteToken deletes a token in a container registry
func (c *Client) DeleteToken(ctx context.Context, registryID, tokenID string) (*shared.APIResponse, error) {
	apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensDelete(ctx, registryID, tokenID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// GetToken gets a token in a container registry
func (c *Client) GetToken(ctx context.Context, registryID, tokenID string) (cr.TokenResponse, *shared.APIResponse, error) {
	token, apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensFindById(ctx, registryID, tokenID).Execute()
	apiResponse.LogInfo()
	return token, apiResponse, err

}

// PatchToken patches a token in a container registry
func (c *Client) PatchToken(ctx context.Context, registryID, tokenID string, tokenInput cr.PatchTokenInput) (cr.TokenResponse, *shared.APIResponse, error) {
	token, apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensPatch(ctx, registryID, tokenID).PatchTokenInput(tokenInput).Execute()
	apiResponse.LogInfo()
	return token, apiResponse, err
}

// PutToken puts a token in a container registry
func (c *Client) PutToken(ctx context.Context, registryID, tokenID string, tokenInput cr.PutTokenInput) (cr.PutTokenOutput, *shared.APIResponse, error) {
	token, apiResponse, err := c.sdkClient.TokensApi.RegistriesTokensPut(ctx, registryID, tokenID).PutTokenInput(tokenInput).Execute()
	apiResponse.LogInfo()
	return token, apiResponse, err
}

// GetRegistryDataCreate get registry data for create
func GetRegistryDataCreate(d *schema.ResourceData) (*cr.PostRegistryInput, error) {

	registry := cr.PostRegistryInput{
		Properties: cr.PostRegistryProperties{},
	}

	if _, ok := d.GetOk("garbage_collection_schedule"); ok {
		registry.Properties.GarbageCollectionSchedule = getWeeklySchedule(d, "garbage_collection_schedule")
	}

	if location, ok := d.GetOk("location"); ok {
		location := location.(string)
		registry.Properties.Location = location
	}

	if name, ok := d.GetOk("name"); ok {
		name := name.(string)
		registry.Properties.Name = name
	}

	if v, ok := d.GetOk("api_subnet_allow_list"); ok {
		raw := v.([]any)
		ips := make([]string, len(raw))
		err := utils.DecodeInterfaceToStruct(raw, ips)
		if err != nil {
			return nil, err
		}
		if len(ips) > 0 {
			registry.Properties.ApiSubnetAllowList = ips
		}
	}

	return &registry, nil
}

// GetRegistryDataUpdate get registry data for update
func GetRegistryDataUpdate(d *schema.ResourceData) (*cr.PatchRegistryInput, error) {

	registry := cr.PatchRegistryInput{}

	if _, ok := d.GetOk("garbage_collection_schedule"); ok {
		registry.GarbageCollectionSchedule = getWeeklySchedule(d, "garbage_collection_schedule")
	}

	// When api_subnet_allow_list = [] in the TF plan:
	// GetOk => _, false
	// GetOkExists => _, true
	// We don't want to ignore the cases in which we modify the api_subnet_allow_list to an empty list, so
	// here we need to use 'GetOkExists'.
	if v, ok := d.GetOkExists("api_subnet_allow_list"); ok { //nolint:staticcheck
		raw := v.([]any)
		ips := make([]string, len(raw))
		err := utils.DecodeInterfaceToStruct(raw, ips)
		if err != nil {
			return nil, err
		}
		registry.ApiSubnetAllowList = ips
	}

	return &registry, nil
}

func getWeeklySchedule(d *schema.ResourceData, property string) *cr.WeeklySchedule {
	var weeklySchedule cr.WeeklySchedule

	if days, ok := d.GetOk(fmt.Sprintf("%v.0.days", property)); ok {
		var daysToAdd []cr.Day
		for _, day := range days.([]any) {
			daysToAdd = append(daysToAdd, cr.Day(day.(string)))
		}
		weeklySchedule.Days = daysToAdd
	}

	if timeField, ok := d.GetOk(fmt.Sprintf("%v.0.time", property)); ok {
		timeStr := timeField.(string)
		weeklySchedule.Time = timeStr
	}

	return &weeklySchedule
}

// SetRegistryData sets the registry data for the container registry
func SetRegistryData(d *schema.ResourceData, registry cr.RegistryResponse) error {

	resourceName := "registry"

	if registry.Id != nil {
		d.SetId(*registry.Id)
	}

	if registry.Properties.GarbageCollectionSchedule != nil {
		var schedule []any
		scheduleEntry := setWeeklySchedule(*registry.Properties.GarbageCollectionSchedule)
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

	if err := d.Set("location", registry.Properties.Location); err != nil {
		return utils.GenerateSetError(resourceName, "location", err)
	}

	if err := d.Set("name", registry.Properties.Name); err != nil {
		return utils.GenerateSetError(resourceName, "name", err)
	}

	if registry.Properties.StorageUsage != nil {
		var storage []any
		storageEntry := setStorageUsage(*registry.Properties.StorageUsage)
		storage = append(storage, storageEntry)
		if err := d.Set("storage_usage", storage); err != nil {
			return utils.GenerateSetError(resourceName, "storage_usage", err)
		}
	}

	if err := d.Set("api_subnet_allow_list", registry.Properties.ApiSubnetAllowList); err != nil {
		return fmt.Errorf("error setting api_subnet_allow_list %w", err)
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

func setWeeklySchedule(weeklySchedule cr.WeeklySchedule) map[string]any {

	schedule := map[string]any{}

	utils.SetPropWithNilCheck(schedule, "time", weeklySchedule.Time)
	utils.SetPropWithNilCheck(schedule, "days", weeklySchedule.Days)

	return schedule
}

func setStorageUsage(storageUsage cr.StorageUsage) map[string]any {

	storage := map[string]any{}

	utils.SetPropWithNilCheck(storage, "bytes", storageUsage.Bytes)
	if storageUsage.UpdatedAt != nil {
		utils.SetPropWithNilCheck(storage, "updated_at", storageUsage.UpdatedAt.String())
	}

	return storage
}

// GetTokenDataCreate returns the token data for the registry token
func GetTokenDataCreate(d *schema.ResourceData) (*cr.PostTokenInput, error) {

	token := cr.PostTokenInput{
		Properties: cr.PostTokenProperties{},
	}

	if expiryDate, ok := d.GetOk("expiry_date"); ok {
		expiryDate, err := convertToTime(expiryDate.(string))
		if err != nil {
			return nil, err
		}
		nullableIonosTime := cr.NullableIonosTime{NullableTime: *(cr.NewNullableTime(expiryDate))}
		token.Properties.ExpiryDate = &nullableIonosTime
	}

	if name, ok := d.GetOk("name"); ok {
		name := name.(string)
		token.Properties.Name = name
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

// GetTokenDataUpdate returns the token data for the registry token
func GetTokenDataUpdate(d *schema.ResourceData) (*cr.PatchTokenInput, error) {

	token := cr.PatchTokenInput{}

	if expiryDate, ok := d.GetOk("expiry_date"); ok {
		expiryDate, err := convertToTime(expiryDate.(string))
		if err != nil {
			return nil, err
		}
		nullableIonosTime := cr.NullableIonosTime{NullableTime: *(cr.NewNullableTime(expiryDate))}
		token.ExpiryDate = &nullableIonosTime
	} else {
		nullableIonosTime := cr.NullableIonosTime{NullableTime: *(cr.NewNullableTime(nil))}
		token.ExpiryDate = &nullableIonosTime
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

// GetScopes returns the scopes for the registry token
func GetScopes(d *schema.ResourceData) []cr.Scope {
	scopes := make([]cr.Scope, 0)

	if scopeValue, ok := d.GetOk("scopes"); ok {
		scopeValue := scopeValue.([]any)
		for _, item := range scopeValue {

			scopeContent := item.(map[string]any)
			connection := cr.Scope{}

			if actions, ok := scopeContent["actions"]; ok {
				actions := actions.([]any)
				var actionsToAdd []string
				if len(actions) > 0 {
					for _, action := range actions {
						actionsToAdd = append(actionsToAdd, action.(string))
					}
				}
				connection.Actions = actionsToAdd
			}

			if name, ok := scopeContent["name"]; ok {
				name := name.(string)
				connection.Name = name
			}

			if scopeType, ok := scopeContent["type"]; ok {
				scopeType := scopeType.(string)
				connection.Type = scopeType
			}

			scopes = append(scopes, connection)
		}

	}

	return scopes

}

// SetTokenData sets the token data for the registry token. Does not set credentials, as they never change once created
func SetTokenData(d *schema.ResourceData, tokenProps cr.TokenProperties) error {

	regToken := "registry token "

	if tokenProps.ExpiryDate != nil {
		timeValue := tokenProps.ExpiryDate.NullableTime
		if err := d.Set("expiry_date", timeValue.Get().String()); err != nil {
			return utils.GenerateSetError(regToken, "expiry_date", err)
		}
	}

	if err := d.Set("name", tokenProps.Name); err != nil {
		return utils.GenerateSetError(regToken, "name", err)
	}

	scopes := setScopes(tokenProps.Scopes)
	if err := d.Set("scopes", scopes); err != nil {
		return utils.GenerateSetError(regToken, "scopes", err)
	}

	if tokenProps.Status != nil {
		if err := d.Set("status", *tokenProps.Status); err != nil {
			return utils.GenerateSetError(regToken, "status", err)
		}
	}

	return nil
}

// SetCredentials sets the credentials for the registry token. username can be set on read, password will be set only on create
func SetCredentials(credentials cr.Credentials) map[string]any {

	credentialsEntry := map[string]any{}
	credentialsEntry["username"] = credentials.Username
	if credentials.Password != "" {
		credentialsEntry["password"] = credentials.Password
	}

	return credentialsEntry
}

func setScopes(scopes []cr.Scope) []any {

	var tokenScopes []any //nolint:prealloc
	for _, scope := range scopes {
		scopeEntry := make(map[string]any)

		if scope.Actions != nil {
			scopeEntry["actions"] = scope.Actions
		}
		scopeEntry["name"] = scope.Name
		scopeEntry["type"] = scope.Type

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
		registryFeatures.VulnerabilityScanning.Enabled = vulnerabilityScanning
		warnings = diag.Diagnostics{crfWarnings.warnCRVScanningOff}

		if vulnerabilityScanning {
			warnings = diag.Diagnostics{crfWarnings.warnCRVScanningOn}
		}
	}

	return registryFeatures, warnings

}

func convertToTime(targetTime string) (*time.Time, error) {
	var convertedTime time.Time
	var err error

	// targetTime might have time zone offset layout (+0000 UTC)
	if convertedTime, err = time.Parse(constant.DatetimeTZOffsetLayout, targetTime); err != nil {
		if convertedTime, err = time.Parse(constant.DatetimeZLayout, targetTime); err != nil {
			return nil, fmt.Errorf("an error occurred while converting from IonosTime string to time.Time: %w", err)
		}
	}
	return &convertedTime, nil
}
