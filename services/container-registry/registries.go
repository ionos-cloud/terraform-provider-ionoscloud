package container_registry

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cr "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"time"
)

type RegistriesService interface {
	ListRegistries(ctx context.Context) (cr.RegistriesResponse, *cr.APIResponse, error)
	CreateRegistry(ctx context.Context, registryInput cr.PostRegistryInput) (cr.PostRegistryOutput, *cr.APIResponse, error)
	DeleteRegistry(ctx context.Context, registryId string) (*cr.APIResponse, error)
	GetRegistry(ctx context.Context, registryId string) (cr.RegistryResponse, *cr.APIResponse, error)
	PatchRegistry(ctx context.Context, registryId string, registryInput cr.PatchRegistryInput) (cr.RegistryResponse, *cr.APIResponse, error)
	PutRegistry(ctx context.Context, registryId string, registryInput cr.PutRegistryInput) (cr.PutRegistryOutput, *cr.APIResponse, error)
	DeleteRepositories(ctx context.Context, registryId, repositoryId string) (*cr.APIResponse, error)
	ListTokens(ctx context.Context, registryId string) (cr.TokensResponse, *cr.APIResponse, error)
	CreateTokens(ctx context.Context, registryId string, tokenInput cr.PostTokenInput) (cr.PostTokenOutput, *cr.APIResponse, error)
	DeleteToken(ctx context.Context, registryId, tokenId string) (*cr.APIResponse, error)
	GetToken(ctx context.Context, registryId, tokenId string) (cr.TokenResponse, *cr.APIResponse, error)
	PatchToken(ctx context.Context, registryId, tokenId string, tokenInput cr.PatchTokenInput) (cr.TokenResponse, *cr.APIResponse, error)
	PutToken(ctx context.Context, registryId, tokenId string, tokenInput cr.PutTokenInput) (cr.PutTokenOutput, *cr.APIResponse, error)
}

func (c *Client) ListRegistries(ctx context.Context) (cr.RegistriesResponse, *cr.APIResponse, error) {
	registry, apiResponse, err := c.RegistriesApi.RegistriesGet(ctx).Execute()
	if apiResponse != nil {
		return registry, apiResponse, err

	}
	return registry, nil, err
}

func (c *Client) CreateRegistry(ctx context.Context, registryInput cr.PostRegistryInput) (cr.PostRegistryOutput, *cr.APIResponse, error) {
	registry, apiResponse, err := c.RegistriesApi.RegistriesPost(ctx).PostRegistryInput(registryInput).Execute()
	if apiResponse != nil {
		return registry, apiResponse, err

	}
	return registry, nil, err
}

func (c *Client) DeleteRegistry(ctx context.Context, registryId string) (*cr.APIResponse, error) {
	apiResponse, err := c.RegistriesApi.RegistriesDelete(ctx, registryId).Execute()
	if apiResponse != nil {
		return apiResponse, err

	}
	return nil, err
}

func (c *Client) GetRegistry(ctx context.Context, registryId string) (cr.RegistryResponse, *cr.APIResponse, error) {
	registries, apiResponse, err := c.RegistriesApi.RegistriesFindById(ctx, registryId).Execute()
	if apiResponse != nil {
		return registries, apiResponse, err

	}
	return registries, nil, err
}

func (c *Client) PatchRegistry(ctx context.Context, registryId string, registryInput cr.PatchRegistryInput) (cr.RegistryResponse, *cr.APIResponse, error) {
	registries, apiResponse, err := c.RegistriesApi.RegistriesPatch(ctx, registryId).PatchRegistryInput(registryInput).Execute()
	if apiResponse != nil {
		return registries, apiResponse, err

	}
	return registries, nil, err
}

func (c *Client) PutRegistry(ctx context.Context, registryId string, registryInput cr.PutRegistryInput) (cr.PutRegistryOutput, *cr.APIResponse, error) {
	registries, apiResponse, err := c.RegistriesApi.RegistriesPut(ctx, registryId).PutRegistryInput(registryInput).Execute()
	if apiResponse != nil {
		return registries, apiResponse, err

	}
	return registries, nil, err
}

func (c *Client) DeleteRepositories(ctx context.Context, registryId, repositoryId string) (*cr.APIResponse, error) {
	apiResponse, err := c.RepositoriesApi.RegistriesRepositoriesDelete(ctx, registryId, repositoryId).Execute()
	if apiResponse != nil {
		return apiResponse, err

	}
	return nil, err
}

func (c *Client) ListTokens(ctx context.Context, registryId string) (cr.TokensResponse, *cr.APIResponse, error) {
	tokens, apiResponse, err := c.TokensApi.RegistriesTokensGet(ctx, registryId).Execute()
	if apiResponse != nil {
		return tokens, apiResponse, err

	}
	return tokens, nil, err
}

func (c *Client) CreateTokens(ctx context.Context, registryId string, tokenInput cr.PostTokenInput) (cr.PostTokenOutput, *cr.APIResponse, error) {
	token, apiResponse, err := c.TokensApi.RegistriesTokensPost(ctx, registryId).PostTokenInput(tokenInput).Execute()
	if apiResponse != nil {
		return token, apiResponse, err

	}
	return token, nil, err
}

func (c *Client) DeleteToken(ctx context.Context, registryId, tokenId string) (*cr.APIResponse, error) {
	apiResponse, err := c.TokensApi.RegistriesTokensDelete(ctx, registryId, tokenId).Execute()
	if apiResponse != nil {
		return apiResponse, err

	}
	return nil, err
}

func (c *Client) GetToken(ctx context.Context, registryId, tokenId string) (cr.TokenResponse, *cr.APIResponse, error) {
	token, apiResponse, err := c.TokensApi.RegistriesTokensFindById(ctx, registryId, tokenId).Execute()
	if apiResponse != nil {
		return token, apiResponse, err

	}
	return token, nil, err
}

func (c *Client) PatchToken(ctx context.Context, registryId, tokenId string, tokenInput cr.PatchTokenInput) (cr.TokenResponse, *cr.APIResponse, error) {
	token, apiResponse, err := c.TokensApi.RegistriesTokensPatch(ctx, registryId, tokenId).PatchTokenInput(tokenInput).Execute()
	if apiResponse != nil {
		return token, apiResponse, err

	}
	return token, nil, err
}

func (c *Client) PutToken(ctx context.Context, registryId, tokenId string, tokenInput cr.PutTokenInput) (cr.PutTokenOutput, *cr.APIResponse, error) {
	token, apiResponse, err := c.TokensApi.RegistriesTokensPut(ctx, registryId, tokenId).PutTokenInput(tokenInput).Execute()
	if apiResponse != nil {
		return token, apiResponse, err

	}
	return token, nil, err
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

	if _, ok := d.GetOk("maintenance_window"); ok {
		registry.Properties.MaintenanceWindow = GetWeeklySchedule(d, "maintenance_window")
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

	if _, ok := d.GetOk("maintenance_window"); ok {
		registry.MaintenanceWindow = GetWeeklySchedule(d, "maintenance_window")
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

	if time, ok := d.GetOk(fmt.Sprintf("%v.0.time", property)); ok {
		time := time.(string)
		weeklySchedule.Time = &time
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

	if registry.Properties.MaintenanceWindow != nil {
		var schedule []interface{}
		scheduleEntry := SetWeeklySchedule(*registry.Properties.MaintenanceWindow)
		schedule = append(schedule, scheduleEntry)
		if err := d.Set("maintenance_window", schedule); err != nil {
			return utils.GenerateSetError(resourceName, "maintenance_window", err)
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
	utils.SetPropWithNilCheck(storage, "updated_at", storageUsage.UpdatedAt)

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

func SetTokenData(d *schema.ResourceData, token cr.TokenResponse) error {

	resourceName := "registry token"

	if token.Id != nil {
		d.SetId(*token.Id)
	}

	if token.Properties.Credentials != nil {
		var credentials []interface{}
		credentialsEntry := SetCredentials(*token.Properties.Credentials)
		credentials = append(credentials, credentialsEntry)
		if err := d.Set("credentials", credentials); err != nil {
			return utils.GenerateSetError(resourceName, "credentials", err)
		}
	}

	// ToDo: fix diff between expiry_date post and get
	//if token.Properties.ExpiryDate != nil {
	//	timeValue := (*token.Properties.ExpiryDate).Time
	//	if err := d.Set("expiry_date", timeValue.String()); err != nil {
	//		return utils.GenerateSetError(resourceName, "expiry_date", err)
	//	}
	//}

	if token.Properties.Name != nil {
		if err := d.Set("name", *token.Properties.Name); err != nil {
			return utils.GenerateSetError(resourceName, "name", err)
		}
	}

	if token.Properties.Scopes != nil {
		scopes := SetScopes(*token.Properties.Scopes)
		if err := d.Set("scopes", scopes); err != nil {
			return utils.GenerateSetError(resourceName, "scopes", err)
		}
	}

	if token.Properties.Status != nil {
		if err := d.Set("status", *token.Properties.Status); err != nil {
			return utils.GenerateSetError(resourceName, "status", err)
		}
	}

	return nil
}

func SetCredentials(credentials cr.Credentials) map[string]interface{} {

	credentialsEntry := map[string]interface{}{}

	utils.SetPropWithNilCheck(credentialsEntry, "username", credentials.Username)
	utils.SetPropWithNilCheck(credentialsEntry, "password", credentials.Password)

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
	layout := "2006-01-02T15:04:05Z"
	convertedTime, err := time.Parse(layout, targetTime)
	if err != nil {
		return nil, fmt.Errorf("an error occured while converting recovery_target_time to time.Time: %s", err)

	}
	ionosTime.Time = convertedTime
	return &ionosTime, nil
}
