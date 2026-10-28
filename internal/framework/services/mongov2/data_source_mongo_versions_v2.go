package mongov2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	mongov2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mongov2"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var _ datasource.DataSourceWithConfigure = (*versionsDataSource)(nil)

type versionDataSourceModel struct {
	CanUpgradeTo types.List   `tfsdk:"can_upgrade_to"`
	Comment      types.String `tfsdk:"comment"`
	ID           types.String `tfsdk:"id"`
	Status       types.String `tfsdk:"status"`
	Version      types.String `tfsdk:"version"`
}

func versionDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"can_upgrade_to": schema.ListAttribute{Computed: true, ElementType: types.StringType, Description: "List of versions this version can be upgraded to. Only versions listed here are valid upgrade targets when changing the `version` field. Downgrades are not supported."},
		"comment":        schema.StringAttribute{Computed: true, Description: "Additional information about the version status."},
		"id":             schema.StringAttribute{Computed: true, Description: "The ID (UUID) of the version."},
		"status":         schema.StringAttribute{Computed: true, Description: "The support status of the version (one of `BETA`, `SUPPORTED`, `RECOMMENDED`, `DEPRECATED`)."},
		"version":        schema.StringAttribute{Computed: true, Description: "The MongoDB version string."},
	}
}

type versionsDataSource struct {
	bundle *bundleclient.SdkBundle
}

type versionsDataSourceModel struct {
	Items    []versionDataSourceModel `tfsdk:"items"`
	Location types.String             `tfsdk:"location"`
}

// NewVersionsDataSource creates a new data source for listing MongoDB v2 versions.
func NewVersionsDataSource() datasource.DataSource {
	return &versionsDataSource{}
}

func (d *versionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongo_versions_v2"
}

func (d *versionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *bundleclient.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.bundle = clientBundle
}

func (d *versionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all supported IONOS CLOUD MongoDB V2 versions in a given location.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of supported versions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: versionDataSourceAttributes(),
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The location to query. Available locations: " + mongov2service.AvailableLocationsString() + ".",
			},
		},
	}
}

func (d *versionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data versionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location := data.Location.ValueString()
	client, err := d.bundle.NewMongoDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MongoDB V2 client", err.Error())
		return
	}

	list, apiResponse, err := client.ListVersions(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to list MongoDB V2 versions", diagutil.WrapError(err, &diagutil.ErrorContext{
			StatusCode:     apiResponse.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	items := make([]versionDataSourceModel, 0, len(list.Items))
	for _, v := range list.Items {
		item := versionDataSourceModel{
			ID: types.StringValue(v.Id),
		}
		resp.Diagnostics.Append(mapVersionToModel(ctx, &v.Properties, &item)...)
		if resp.Diagnostics.HasError() {
			return
		}
		items = append(items, item)
	}
	data.Items = items

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapVersionToModel(ctx context.Context, props *mongov3.MongoVersion, model *versionDataSourceModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	model.Version = types.StringPointerValue(props.Version)
	model.Status = types.StringPointerValue(props.Status)
	model.Comment = types.StringPointerValue(props.Comment)

	listVal, diags := types.ListValueFrom(ctx, types.StringType, props.CanUpgradeTo)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return diagnostics
	}
	model.CanUpgradeTo = listVal
	return diagnostics
}
