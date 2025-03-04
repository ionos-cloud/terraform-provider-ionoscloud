package objectstorage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstorage"
)

var (
	_ datasource.DataSourceWithConfigure = (*objectDataSource)(nil)
)

// NewObjectDataSource creates a new data source for object.
func NewObjectDataSource() datasource.DataSource {
	return &objectDataSource{}
}

type objectDataSource struct {
	client *objectstorage.Client
}

// Metadata returns the metadata for the object data source.
func (d *objectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_object"
}

// Configure configures the data source.
func (d *objectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = clientBundle.S3Client
}

// Schema returns the schema for the object data source.
func (d *objectDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required: true,
			},
			"key": schema.StringAttribute{
				Required: true,
			},
			"body": schema.StringAttribute{
				Computed: true,
			},
			"cache_control": schema.StringAttribute{
				Computed: true,
			},
			"content_length": schema.Int64Attribute{
				Computed: true,
			},
			"content_disposition": schema.StringAttribute{
				Computed: true,
			},
			"content_encoding": schema.StringAttribute{
				Computed: true,
			},
			"content_language": schema.StringAttribute{
				Computed: true,
			},
			"content_type": schema.StringAttribute{
				Computed: true,
			},
			"expires": schema.StringAttribute{
				Computed: true,
			},
			"server_side_encryption": schema.StringAttribute{
				Computed: true,
			},
			"storage_class": schema.StringAttribute{
				Computed: true,
			},
			"website_redirect": schema.StringAttribute{
				Computed: true,
			},
			"server_side_encryption_customer_algorithm": schema.StringAttribute{
				Computed:   true,
				Validators: []validator.String{stringvalidator.OneOf("AES256")},
			},
			"server_side_encryption_customer_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"server_side_encryption_customer_key_md5": schema.StringAttribute{
				Computed: true,
			},
			"server_side_encryption_context": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"request_payer": schema.StringAttribute{
				Computed: true,
			},
			"object_lock_mode": schema.StringAttribute{
				Computed:   true,
				Validators: []validator.String{stringvalidator.OneOf("GOVERNANCE", "COMPLIANCE")},
			},
			"object_lock_retain_until_date": schema.StringAttribute{
				Computed: true,
			},
			"object_lock_legal_hold": schema.StringAttribute{
				Computed:   true,
				Validators: []validator.String{stringvalidator.OneOf("ON", "OFF")},
			},
			"etag": schema.StringAttribute{
				Computed: true,
			},
			"tags": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"metadata": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"version_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"range": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Read the data source
func (d *objectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *objectstorage.ObjectDataSourceModel

	// Read configuration
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, found, err := d.client.GetObjectForDataSource(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	if !found {
		resp.Diagnostics.AddError("Failed to read resource", "Resource not found")
		return
	}

	data = result
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
