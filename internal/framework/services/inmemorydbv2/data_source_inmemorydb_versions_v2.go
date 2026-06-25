package inmemorydbv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	inmemorydbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	inmemorydbv2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/inmemorydbv2"
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
		"can_upgrade_to": schema.ListAttribute{Computed: true, ElementType: types.StringType, Description: "List of versions that a cluster running this version can be upgraded to."},
		"comment":        schema.StringAttribute{Computed: true, Description: "Additional human-readable information about the version lifecycle."},
		"id":             schema.StringAttribute{Computed: true, Description: "The ID (UUID) of the version."},
		"status":         schema.StringAttribute{Computed: true, Description: "The support status of the version."},
		"version":        schema.StringAttribute{Computed: true, Description: "The version for the cluster."},
	}
}

type versionsDataSource struct {
	bundle *bundleclient.SdkBundle
}

type versionsDataSourceModel struct {
	Items    []versionDataSourceModel `tfsdk:"items"`
	Location types.String             `tfsdk:"location"`
}

// NewVersionsDataSource creates a new data source for listing InMemoryDB v2 versions.
func NewVersionsDataSource() datasource.DataSource {
	return &versionsDataSource{}
}

func (d *versionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inmemorydb_versions_v2"
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
		Description: "Lists all supported InMemoryDB v2 versions in a given location.",
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
				Description: "The location to query. Available locations: " + inmemorydbv2service.AvailableLocationsString() + ".",
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
	client, err := d.bundle.NewInMemoryDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create InMemoryDB v2 client", err.Error())
		return
	}

	list, _, err := client.ListVersions(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to list InMemoryDB v2 versions in location %s", location), err.Error())
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

func mapVersionToModel(ctx context.Context, props *inmemorydbv3.SupportedVersion, model *versionDataSourceModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	model.Version = types.StringPointerValue(props.Version)
	model.Status = types.StringPointerValue(props.Status)
	model.Comment = types.StringPointerValue(props.Comment)

	upgrades := make([]types.String, len(props.CanUpgradeTo))
	for i, v := range props.CanUpgradeTo {
		upgrades[i] = types.StringValue(v)
	}
	listVal, diags := types.ListValueFrom(ctx, types.StringType, upgrades)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return diagnostics
	}
	model.CanUpgradeTo = listVal
	return diagnostics
}
