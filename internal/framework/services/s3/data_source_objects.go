package s3

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/s3"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type objectsDataSource struct {
	client *s3.Client
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

	client, ok := req.ProviderData.(*s3.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *s3.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *objectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *s3.ObjectsDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.client.ListObjects(ctx, data); err != nil {
		resp.Diagnostics.AddError("error fetching objects", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
