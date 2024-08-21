package s3

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"

	tfs3 "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/s3"
)

type objectsDataSource struct {
	client *s3.APIClient
}

type objectsDataSourceModel struct {
	Bucket         types.String   `tfsdk:"bucket"`
	Delimiter      types.String   `tfsdk:"delimiter"`
	EncodingType   types.String   `tfsdk:"encoding_type"`
	MaxKeys        types.Int64    `tfsdk:"max_keys"`
	Prefix         types.String   `tfsdk:"prefix"`
	FetchOwner     types.Bool     `tfsdk:"fetch_owner"`
	StartAfter     types.String   `tfsdk:"start_after"`
	CommonPrefixes []types.String `tfsdk:"common_prefixes"`
	Keys           []types.String `tfsdk:"keys"`
	Owners         []types.String `tfsdk:"owners"`
}

// NewObjectsDataSource creates a new data source for fetching objects from an S3 bucket.
func NewObjectsDataSource() datasource.DataSource {
	return &objectsDataSource{}
}

func (d *objectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_objects"
}

func (d *objectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required: true,
			},
			"delimiter": schema.StringAttribute{
				Optional: true,
			},
			"encoding_type": schema.StringAttribute{
				Optional: true,
			},
			"max_keys": schema.Int64Attribute{
				Optional: true,
			},
			"prefix": schema.StringAttribute{
				Optional:   true,
				Validators: []validator.String{stringvalidator.LengthBetween(0, 1024)},
			},
			"fetch_owner": schema.BoolAttribute{
				Optional: true,
			},
			"start_after": schema.StringAttribute{
				Optional: true,
			},
			"common_prefixes": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"keys": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"owners": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *objectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*s3.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *objectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *objectsDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := fetchObjects(ctx, d.client, data); err != nil {
		resp.Diagnostics.AddError("error fetching objects", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func fetchObjects(ctx context.Context, client *s3.APIClient, data *objectsDataSourceModel) error {
	input := &tfs3.ListObjectsV2Input{
		Bucket:       data.Bucket.ValueString(),
		Delimiter:    data.Delimiter.ValueStringPointer(),
		EncodingType: data.EncodingType.ValueStringPointer(),
		MaxKeys:      toInt32(data.MaxKeys.ValueInt64Pointer()),
		Prefix:       data.Prefix.ValueStringPointer(),
		FetchOwner:   data.FetchOwner.ValueBool(),
		StartAfter:   data.StartAfter.ValueStringPointer(),
	}

	var maxKeys, nKeys int64
	if data.MaxKeys.IsNull() {
		maxKeys = 1000
	} else {
		maxKeys = data.MaxKeys.ValueInt64()
	}

	keys := make([]types.String, 0)
	owners := make([]types.String, 0)
	pages := tfs3.NewListObjectsV2Paginator(client, input)
pageLoop:
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error fetching page: %w", err)
		}

		if page.CommonPrefixes != nil {
			data.CommonPrefixes = make([]types.String, len(*page.CommonPrefixes))
			for i, prefix := range *page.CommonPrefixes {
				data.CommonPrefixes[i] = types.StringPointerValue(prefix.Prefix)
			}
		}

		if page.Contents != nil {
			for _, v := range *page.Contents {
				if nKeys >= maxKeys {
					break pageLoop
				}

				keys = append(keys, types.StringPointerValue(v.Key))
				if v.Owner != nil {
					owners = append(owners, types.StringPointerValue(v.Owner.DisplayName))
				}

				nKeys++
			}
		}
	}

	data.Keys = keys
	data.Owners = owners
	return nil
}
