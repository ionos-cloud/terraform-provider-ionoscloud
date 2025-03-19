package objectstoragemanagement

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	objectstoragemanagement "github.com/ionos-cloud/sdk-go-bundle/products/objectstoragemanagement/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

// AccesskeyResourceModel is used to represent an accesskey
type AccesskeyResourceModel struct {
	AccessKey       types.String   `tfsdk:"accesskey"`
	SecretKey       types.String   `tfsdk:"secretkey"`
	CanonicalUserID types.String   `tfsdk:"canonical_user_id"`
	ContractUserID  types.String   `tfsdk:"contract_user_id"`
	Description     types.String   `tfsdk:"description"`
	ID              types.String   `tfsdk:"id"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}

// AccessKeyDataSourceModel is used to represent an accesskey for a data source.
type AccessKeyDataSourceModel struct {
	AccessKey       types.String `tfsdk:"accesskey"`
	CanonicalUserID types.String `tfsdk:"canonical_user_id"`
	ContractUserID  types.String `tfsdk:"contract_user_id"`
	Description     types.String `tfsdk:"description"`
	ID              types.String `tfsdk:"id"`
}

var ionosAPIURLObjectStorageManagement = "IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT"

// modifyConfigURL modifies the URL inside the client configuration.
// This function is required in order to make requests to different endpoints.
func (c *Client) modifyConfigURL() {
	clientConfig := c.client.GetConfig()
	if os.Getenv(ionosAPIURLObjectStorageManagement) != "" {
		clientConfig.Servers = shared.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLObjectStorageManagement)),
			},
		}
		return
	}
}

// GetAccessKey retrieves an accesskey
func (c *Client) GetAccessKey(ctx context.Context, accessKeyID string) (objectstoragemanagement.AccessKeyRead, *shared.APIResponse, error) {
	c.modifyConfigURL()
	accessKey, apiResponse, err := c.client.AccesskeysApi.AccesskeysFindById(ctx, accessKeyID).Execute()
	if apiResponse.HttpNotFound() {
		return accessKey, apiResponse, nil
	}
	apiResponse.LogInfo()
	return accessKey, apiResponse, err
}

// ListAccessKeys retrieves all accesskeys
func (c *Client) ListAccessKeys(ctx context.Context) (objectstoragemanagement.AccessKeyReadList, *shared.APIResponse, error) {
	c.modifyConfigURL()
	accessKeys, apiResponse, err := c.client.AccesskeysApi.AccesskeysGet(ctx).Execute()
	apiResponse.LogInfo()
	return accessKeys, apiResponse, err
}

// ListAccessKeysFilter retrieves accesskeys using the accessKeyId filter
func (c *Client) ListAccessKeysFilter(ctx context.Context, accessKeyID string) (objectstoragemanagement.AccessKeyReadList, *shared.APIResponse, error) {
	c.modifyConfigURL()
	accessKeys, apiResponse, err := c.client.AccesskeysApi.AccesskeysGet(ctx).FilterAccesskeyId(accessKeyID).Execute()
	apiResponse.LogInfo()
	return accessKeys, apiResponse, err
}

// CreateAccessKey creates an accesskey
func (c *Client) CreateAccessKey(ctx context.Context, accessKey objectstoragemanagement.AccessKeyCreate, timeout time.Duration) (objectstoragemanagement.AccessKeyRead, *shared.APIResponse, error) {
	c.modifyConfigURL()
	accessKeyResponse, apiResponse, err := c.client.AccesskeysApi.AccesskeysPost(ctx).AccessKeyCreate(accessKey).Execute()
	apiResponse.LogInfo()

	if err != nil {
		return accessKeyResponse, apiResponse, err
	}

	err = backoff.Retry(func() error {
		return c.accessKeyAvailableCheck(ctx, accessKeyResponse.Id)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)))
	if err != nil {
		return accessKeyResponse, apiResponse, fmt.Errorf("failed to wait for accessKey available: %w", err)
	}

	return accessKeyResponse, apiResponse, err
}

// UpdateAccessKey updates an accesskey
func (c *Client) UpdateAccessKey(ctx context.Context, accessKeyID string, accessKey objectstoragemanagement.AccessKeyEnsure, timeout time.Duration) (objectstoragemanagement.AccessKeyRead, *shared.APIResponse, error) {
	c.modifyConfigURL()
	accessKeyResponse, apiResponse, err := c.client.AccesskeysApi.AccesskeysPut(ctx, accessKeyID).AccessKeyEnsure(accessKey).Execute()
	apiResponse.LogInfo()

	if err != nil {
		return accessKeyResponse, apiResponse, err
	}

	err = backoff.Retry(func() error {
		return c.accessKeyAvailableCheck(ctx, accessKeyID)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)))
	if err != nil {
		return accessKeyResponse, apiResponse, fmt.Errorf("failed to wait for accessKey available: %w", err)
	}

	return accessKeyResponse, apiResponse, err
}

// DeleteAccessKey deletes an accesskey
func (c *Client) DeleteAccessKey(ctx context.Context, accessKeyID string, timeout time.Duration) (*shared.APIResponse, error) {
	c.modifyConfigURL()
	apiResponse, err := c.client.AccesskeysApi.AccesskeysDelete(ctx, accessKeyID).Execute()
	apiResponse.LogInfo()

	if err != nil {
		return apiResponse, err
	}

	err = backoff.Retry(func() error {
		return c.accessKeyDeletedCheck(ctx, accessKeyID)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)))
	if err != nil {
		return apiResponse, fmt.Errorf("failed to wait for accessKey available: %w", err)
	}

	return apiResponse, err
}

// SetAccessKeyPropertiesToPlan sets accesskey properties from an SDK object to a AccesskeyResourceModel
func SetAccessKeyPropertiesToPlan(plan *AccesskeyResourceModel, accessKey objectstoragemanagement.AccessKeyRead) {
	plan.AccessKey = basetypes.NewStringValue(accessKey.Properties.AccessKey)
	plan.CanonicalUserID = basetypes.NewStringPointerValue(accessKey.Properties.CanonicalUserId)
	plan.ContractUserID = basetypes.NewStringPointerValue(accessKey.Properties.ContractUserId)
	plan.Description = basetypes.NewStringValue(accessKey.Properties.Description)
	// The secret key is present only in the POST response, on subsequent GET calls we don't
	// want to overwrite the secret key with nil, if the value is set just leave it as it is.
	if plan.SecretKey.IsUnknown() {
		plan.SecretKey = basetypes.NewStringValue(accessKey.Properties.SecretKey)
	}
	plan.ID = basetypes.NewStringValue(accessKey.Id)
}

// SetAccessKeyPropertiesToDataSourcePlan sets accesskey properties from an SDK object to a AccessKeyDataSourceModel
func SetAccessKeyPropertiesToDataSourcePlan(plan *AccessKeyDataSourceModel, accessKey objectstoragemanagement.AccessKeyRead) {
	// Here we check the properties because based on the request not all are set and we do not want to overwrite with nil
	if accessKey.Properties.AccessKey != "" {
		plan.AccessKey = basetypes.NewStringValue(accessKey.Properties.AccessKey)
	}
	if accessKey.Properties.CanonicalUserId != nil {
		plan.CanonicalUserID = basetypes.NewStringPointerValue(accessKey.Properties.CanonicalUserId)
	}
	if accessKey.Properties.ContractUserId != nil {
		plan.ContractUserID = basetypes.NewStringPointerValue(accessKey.Properties.ContractUserId)
	}
	if accessKey.Properties.Description != "" {
		plan.Description = basetypes.NewStringValue(accessKey.Properties.Description)
	}
	if accessKey.Id != "" {
		plan.ID = basetypes.NewStringValue(accessKey.Id)
	}
}

func (c *Client) accessKeyDeletedCheck(ctx context.Context, id string) error {
	_, apiResponse, err := c.GetAccessKey(ctx, id)
	if apiResponse.HttpNotFound() {
		return nil
	}

	if err != nil {
		return backoff.Permanent(fmt.Errorf("failed to check if accessKey exists: %w", err))
	}

	return fmt.Errorf("accessKey still exists")
}

func (c *Client) accessKeyAvailableCheck(ctx context.Context, id string) error {
	accessKey, apiResponse, err := c.GetAccessKey(ctx, id)
	if apiResponse.HttpNotFound() {
		return fmt.Errorf("accessKey not found")
	}

	if err != nil {
		return backoff.Permanent(fmt.Errorf("failed to check if accessKey exists: %w", err))
	}

	if accessKey.Metadata.Status != "AVAILABLE" {
		return fmt.Errorf("accessKey status is not 'AVAILABLE'")
	}

	return nil
}
