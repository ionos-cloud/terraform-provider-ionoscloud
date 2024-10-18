package objectstorage

import (
	"context"
	"encoding/hex"
	"encoding/xml"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	objstorage "github.com/ionos-cloud/sdk-go-s3"

	convptr "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/convptr"
	hash2 "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/hash"
)

const objectLockEnabled = "Enabled"

// ObjectLockConfigurationModel is used to create, update and delete a bucket object lock configuration.
type ObjectLockConfigurationModel struct {
	Bucket            types.String `tfsdk:"bucket"`
	ObjectLockEnabled types.String `tfsdk:"object_lock_enabled"`
	Rule              *rule        `tfsdk:"rule"`
}

type rule struct {
	DefaultRetention *defaultRetention `tfsdk:"default_retention"`
}

type defaultRetention struct {
	Mode  types.String `tfsdk:"mode"`
	Days  types.Int64  `tfsdk:"days"`
	Years types.Int64  `tfsdk:"years"`
}

// CreateObjectLock creates a new bucket object lock configuration.
func (c *Client) CreateObjectLock(ctx context.Context, data *ObjectLockConfigurationModel) error {
	input := buildObjectLockConfigurationFromModel(data)
	bytes, err := xml.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal bucket lifecycle configuration: %w", err)
	}

	md5Sum, err := hash2.MD5(bytes)
	if err != nil {
		return fmt.Errorf("failed to generate MD5 sum: %s", err.Error())
	}
	_, err = c.client.ObjectLockApi.PutObjectLockConfiguration(ctx, data.Bucket.ValueString()).PutObjectLockConfigurationRequest(input).ContentMD5(hex.EncodeToString([]byte(md5Sum))).Execute()
	return err
}

// GetObjectLock gets a bucket object lock configuration.
func (c *Client) GetObjectLock(ctx context.Context, name types.String) (*ObjectLockConfigurationModel, bool, error) {
	output, httpResp, err := c.client.ObjectLockApi.GetObjectLockConfiguration(ctx, name.ValueString()).Execute()
	if httpResp.HttpNotFound() {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	builtModel := buildObjectLockConfigurationModelFromAPIResponse(output, &ObjectLockConfigurationModel{Bucket: name})
	return builtModel, true, nil
}

// GetObjectLockEnabled gets a bucket object lock configuration.
func (c *Client) GetObjectLockEnabled(ctx context.Context, name types.String) (types.Bool, error) {
	output, httpResp, err := c.client.ObjectLockApi.GetObjectLockConfiguration(ctx, name.ValueString()).Execute()
	if httpResp.HttpNotFound() {
		return types.BoolValue(false), nil
	}

	if err != nil {
		return types.BoolNull(), fmt.Errorf("failed to get object lock: %w", err)
	}

	if output != nil && output.ObjectLockEnabled != nil && *output.ObjectLockEnabled == objectLockEnabled {
		return types.BoolValue(true), nil
	}

	return types.BoolValue(false), nil
}

// UpdateObjectLock updates a bucket object lock configuration.
func (c *Client) UpdateObjectLock(ctx context.Context, data *ObjectLockConfigurationModel) error {
	if err := c.CreateObjectLock(ctx, data); err != nil {
		return err
	}

	model, found, err := c.GetObjectLock(ctx, data.Bucket)
	if !found {
		return fmt.Errorf("bucket object lock configuration not found")
	}

	if err != nil {
		return err
	}

	*data = *model
	return nil
}

func buildObjectLockConfigurationModelFromAPIResponse(output *objstorage.GetObjectLockConfigurationOutput, data *ObjectLockConfigurationModel) *ObjectLockConfigurationModel {
	built := &ObjectLockConfigurationModel{
		Bucket:            data.Bucket,
		ObjectLockEnabled: types.StringPointerValue(output.ObjectLockEnabled),
	}
	if output.Rule != nil {
		built.Rule = &rule{
			DefaultRetention: &defaultRetention{
				Mode:  types.StringPointerValue(output.Rule.DefaultRetention.Mode),
				Days:  types.Int64PointerValue(convptr.Int32ToInt64(output.Rule.DefaultRetention.Days)),
				Years: types.Int64PointerValue(convptr.Int32ToInt64(output.Rule.DefaultRetention.Years)),
			},
		}
	}

	return built
}

func buildObjectLockConfigurationFromModel(data *ObjectLockConfigurationModel) objstorage.PutObjectLockConfigurationRequest {
	req := objstorage.PutObjectLockConfigurationRequest{
		ObjectLockEnabled: data.ObjectLockEnabled.ValueStringPointer(),
		Rule: &objstorage.PutObjectLockConfigurationRequestRule{
			DefaultRetention: &objstorage.DefaultRetention{
				Mode:  data.Rule.DefaultRetention.Mode.ValueStringPointer(),
				Days:  convptr.Int64ToInt32(data.Rule.DefaultRetention.Days.ValueInt64Pointer()),
				Years: convptr.Int64ToInt32(data.Rule.DefaultRetention.Years.ValueInt64Pointer()),
			},
		},
	}
	return req
}
