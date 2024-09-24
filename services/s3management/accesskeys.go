package s3management

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	s3management "github.com/ionos-cloud/sdk-go-s3-management"
)

type AccesskeyResourceModel struct {
	AccessKey       types.String   `tfsdk:"accesskey"`
	SecretKey       types.String   `tfsdk:"secretkey"`
	CanonicalUserId types.String   `tfsdk:"canonical_user_id"`
	ContractUserId  types.String   `tfsdk:"contract_user_id"`
	Description     types.String   `tfsdk:"description"`
	ID              types.String   `tfsdk:"id"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}

// AccessKeyDataSourceModel is used to represent an accesskey for a data source.
type AccessKeyDataSourceModel struct {
	AccessKey       types.String `tfsdk:"accesskey"`
	CanonicalUserId types.String `tfsdk:"canonical_user_id"`
	ContractUserId  types.String `tfsdk:"contract_user_id"`
	Description     types.String `tfsdk:"description"`
	ID              types.String `tfsdk:"id"`
}

func (c *Client) GetAccessKey(ctx context.Context, accessKeyId string) (s3management.AccessKey, *s3management.APIResponse, error) {
	accessKey, apiResponse, err := c.client.AccesskeysApi.AccesskeysFindById(ctx, accessKeyId).Execute()
	apiResponse.LogInfo()
	return accessKey, apiResponse, err
}

func (c *Client) ListAccessKeys(ctx context.Context) (s3management.AccessKeyList, *s3management.APIResponse, error) {
	accessKeys, apiResponse, err := c.client.AccesskeysApi.AccesskeysGet(ctx).Execute()
	apiResponse.LogInfo()
	return accessKeys, apiResponse, err
}

func (c *Client) CreateAccessKey(ctx context.Context, accessKey s3management.AccessKeyCreate, timeout time.Duration) (s3management.AccessKey, *s3management.APIResponse, error) {
	accessKeyResponse, apiResponse, err := c.client.AccesskeysApi.AccesskeysPost(ctx).AccessKeyCreate(accessKey).Execute()
	apiResponse.LogInfo()
	if err == nil && accessKeyResponse.Id != nil {
		err = backoff.Retry(func() error {
			return c.accessKeyAvailableCheck(ctx, *accessKeyResponse.Id)
		}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)))
		if err != nil {
			return accessKeyResponse, apiResponse, fmt.Errorf("failed to wait for accessKey available: %w", err)
		}
	}

	return accessKeyResponse, apiResponse, err
}

func (c *Client) UpdateAccessKey(ctx context.Context, accessKeyId string, accessKey s3management.AccessKeyEnsure, timeout time.Duration) (s3management.AccessKey, *s3management.APIResponse, error) {
	accessKeyResponse, apiResponse, err := c.client.AccesskeysApi.AccesskeysPut(ctx, accessKeyId).AccessKeyEnsure(accessKey).Execute()
	apiResponse.LogInfo()
	if err == nil {
		err = backoff.Retry(func() error {
			return c.accessKeyAvailableCheck(ctx, accessKeyId)
		}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)))
		if err != nil {
			return accessKeyResponse, apiResponse, fmt.Errorf("failed to wait for accessKey available: %w", err)
		}
	}

	return accessKeyResponse, apiResponse, err
}

func (c *Client) DeleteAccessKey(ctx context.Context, accessKeyId string, timeout time.Duration) (*s3management.APIResponse, error) {
	apiResponse, err := c.client.AccesskeysApi.AccesskeysDelete(ctx, accessKeyId).Execute()
	apiResponse.LogInfo()

	if err == nil {
		err = backoff.Retry(func() error {
			return c.accessKeyDeletedCheck(ctx, accessKeyId)
		}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)))
		if err != nil {
			return apiResponse, fmt.Errorf("failed to wait for accessKey available: %w", err)
		}
	}

	return apiResponse, err
}

func SetAccessKeyPropertiesToPlan(plan *AccesskeyResourceModel, accessKey s3management.AccessKey) {

	if accessKey.Properties != nil {
		if accessKey.Properties.AccessKey != nil {
			plan.AccessKey = basetypes.NewStringPointerValue(accessKey.Properties.AccessKey)
		}
		if accessKey.Properties.CanonicalUserId != nil {
			plan.CanonicalUserId = basetypes.NewStringPointerValue(accessKey.Properties.CanonicalUserId)
		}
		if accessKey.Properties.ContractUserId != nil {
			plan.ContractUserId = basetypes.NewStringPointerValue(accessKey.Properties.ContractUserId)
		}
		if accessKey.Properties.Description != nil {
			plan.Description = basetypes.NewStringPointerValue(accessKey.Properties.Description)
		}
		if accessKey.Properties.SecretKey != nil {
			plan.SecretKey = basetypes.NewStringPointerValue(accessKey.Properties.SecretKey)
		}
	}
	if accessKey.Id != nil {
		plan.ID = basetypes.NewStringPointerValue(accessKey.Id)
	}
}

func SetAccessKeyPropertiesToDataSourcePlan(plan *AccessKeyDataSourceModel, accessKey s3management.AccessKey) {

	if accessKey.Properties != nil {
		if accessKey.Properties.AccessKey != nil {
			plan.AccessKey = basetypes.NewStringPointerValue(accessKey.Properties.AccessKey)
		}
		if accessKey.Properties.CanonicalUserId != nil {
			plan.CanonicalUserId = basetypes.NewStringPointerValue(accessKey.Properties.CanonicalUserId)
		}
		if accessKey.Properties.ContractUserId != nil {
			plan.ContractUserId = basetypes.NewStringPointerValue(accessKey.Properties.ContractUserId)
		}
		if accessKey.Properties.Description != nil {
			plan.Description = basetypes.NewStringPointerValue(accessKey.Properties.Description)
		}
	}
	if accessKey.Id != nil {
		plan.ID = basetypes.NewStringPointerValue(accessKey.Id)
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

	if *accessKey.Metadata.Status != "AVAILABLE" {
		return fmt.Errorf("accessKey status is not 'AVAILABLE'")
	}

	return nil
}
