package objectstoragemanagement

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstoragemanagement"
	objectStorageManagementService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstoragemanagement"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	objectstoragemanagementApi "github.com/ionos-cloud/sdk-go-object-storage-management"
)

var _ datasource.DataSourceWithConfigure = (*accessKeyDataSource)(nil)

// NewAccesskeyDataSource creates a new data source for the accesskey resource.
func NewAccesskeyDataSource() datasource.DataSource {
	return &accessKeyDataSource{}
}

type accessKeyDataSource struct {
	client *objectStorageManagementService.Client
}

// Metadata returns the metadata for the data source.
func (d *accessKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage_accesskey"
}

// Configure configures the data source.
func (d *accessKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = clientBundle.ObjectStorageManagementClient
}

// Schema returns the schema for the data source.
func (d *accessKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID (UUID) of the AccessKey.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the Access key.",
				Computed:    true,
			},
			"accesskey": schema.StringAttribute{
				Optional:    true,
				Description: "Access key metadata is a string of 92 characters.",
				Computed:    true,
			},
			"canonical_user_id": schema.StringAttribute{
				Description: "The canonical user ID which is valid for user-owned buckets.",
				Computed:    true,
			},
			"contract_user_id": schema.StringAttribute{
				Description: "The contract user ID which is valid for contract-owned buckets",
				Computed:    true,
			},
		},
	}
}

// Read reads the data source.
func (d *accessKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("object storage management api client not configured", "The provider client is not configured")
		return
	}

	var data *objectstoragemanagement.AccessKeyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()
	accessKeyID := data.AccessKey.ValueString()
	description := data.Description.ValueString()

	var accessKey objectstoragemanagementApi.AccessKeyRead
	var accessKeys objectstoragemanagementApi.AccessKeyReadList
	var apiResponse *objectstoragemanagementApi.APIResponse
	var err error
	switch {
	case !data.ID.IsNull():
		accessKey, apiResponse, err = d.client.GetAccessKey(ctx, id)

		if apiResponse.HttpNotFound() {
			resp.Diagnostics.AddError("accesskey not found", "The accesskey was not found")
			return
		}
		if err != nil {
			resp.Diagnostics.AddError("an error occurred while fetching the accesskey with", err.Error())
			return
		}
		if !data.AccessKey.IsNull() && *accessKey.Properties.AccessKey != accessKeyID {
			resp.Diagnostics.AddError(
				"accesskeyID does not match",
				fmt.Sprintf(
					"accesskey property of Accesskey (UUID=%s, accesskey=%s) does not match expected accesskey: %s",
					*accessKey.Id, *accessKey.Properties.AccessKey, accessKeyID,
				),
			)
			return
		}
		if !data.Description.IsNull() && *accessKey.Properties.Description != description {
			resp.Diagnostics.AddError(
				"accesskeyID does not match",
				fmt.Sprintf(
					"description of Accesskey (UUID=%s, description=%s) does not match expected description: %s",
					*accessKey.Id, *accessKey.Properties.Description, description,
				),
			)
			return
		}
	case !data.AccessKey.IsNull():
		accessKeys, _, err = d.client.ListAccessKeysFilter(ctx, accessKeyID)
		if err != nil {
			resp.Diagnostics.AddError("an error occurred while fetching the accesskeys", err.Error())
			return
		}
		if len(*accessKeys.Items) != 0 {
			accessKey = (*accessKeys.Items)[0]
		} else {
			resp.Diagnostics.AddError("accesskey not found", "The accesskey was not found")
			return
		}
		if !data.Description.IsNull() && *accessKey.Properties.Description != description {
			resp.Diagnostics.AddError(
				"accesskeyID does not match",
				fmt.Sprintf(
					"description of Accesskey (UUID=%s, description=%s) does not match expected description: %s",
					*accessKey.Id, *accessKey.Properties.Description, description,
				),
			)
			return
		}
	case !data.Description.IsNull():
		accessKeys, _, err = d.client.ListAccessKeys(ctx)
		if err != nil {
			resp.Diagnostics.AddError("an error occurred while fetching the accesskeys", err.Error())
			return
		}
		found := false
		for _, item := range *accessKeys.Items {
			if *item.Properties.Description == description {
				accessKey = item
				found = true
				break
			}
		}
		if !found {
			resp.Diagnostics.AddError("accesskey not found", "The accesskey was not found")
			return
		}
	default:
		resp.Diagnostics.AddError("ID, accesskeyId or description must be set", "ID, accesskeyId or description must be set")
		return
	}

	objectStorageManagementService.SetAccessKeyPropertiesToDataSourcePlan(data, accessKey)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
