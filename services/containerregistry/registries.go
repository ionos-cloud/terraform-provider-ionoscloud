package containerregistry

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cr "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"time"
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

func GetRegistryDataCreate(d *schema.ResourceData) *cr.PostRegistryInput {

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

	return &registry
}

func GetRegistryDataUpdate(d *schema.ResourceData) *cr.PatchRegistryInput {

	registry := cr.PatchRegistryInput{}

	if _, ok := d.GetOk("garbage_collection_schedule"); ok {
		registry.GarbageCollectionSchedule = GetWeeklySchedule(d, "garbage_collection_schedule")
	}

	return &registry
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

func convertToIonosTime(targetTime string) (*cr.IonosTime, error) {
	var ionosTime cr.IonosTime
	layout := "2006-01-02 15:04:05Z"
	convertedTime, err := time.Parse(layout, targetTime)
	if err != nil {
		return nil, fmt.Errorf("an error occured while converting from IonosTime to time.Time: %w", err)
	}
	ionosTime.Time = convertedTime
	return &ionosTime, nil
}
