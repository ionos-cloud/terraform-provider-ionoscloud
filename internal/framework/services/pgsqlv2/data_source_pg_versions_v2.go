package pgsqlv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	pgsqlv2 "github.com/ionos-cloud/pgsqlv2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

var _ datasource.DataSourceWithConfigure = (*versionsDataSource)(nil)

type versionsDataSource struct {
	bundle *bundleclient.SdkBundle
}

type versionsDataSourceModel struct {
	Location types.String       `tfsdk:"location"`
	Versions []versionModel `tfsdk:"versions"`
}

type versionModel struct {
	ID           types.String `tfsdk:"id"`
	Version      types.String `tfsdk:"version"`
	Status       types.String `tfsdk:"status"`
	Comment      types.String `tfsdk:"comment"`
	CanUpgradeTo types.List   `tfsdk:"can_upgrade_to"`
}

// NewVersionsDataSource creates a new data source for listing PgSQL v2 versions.
func NewVersionsDataSource() datasource.DataSource {
	return &versionsDataSource{}
}

// Metadata returns the metadata for the versions data source.
func (d *versionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pg_versions_v2"
}

// Configure configures the data source.
func (d *versionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.bundle = clientBundle
}

// Schema returns the schema for the versions data source.
func (d *versionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists available PostgreSQL v2 versions.",
		Attributes: map[string]schema.Attribute{
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The region in which to look up available versions.",
			},
			"versions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of available PostgreSQL versions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID (UUID) of the PostgreSQL version.",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "The PostgreSQL version string.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The support status of the version (e.g. BETA, SUPPORTED, RECOMMENDED, DEPRECATED).",
						},
						"comment": schema.StringAttribute{
							Computed:    true,
							Description: "Additional information about the version status.",
						},
						"can_upgrade_to": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "List of versions that this version can be upgraded to.",
						},
					},
				},
			},
		},
	}
}

// Read reads the PgSQL v2 versions data source.
func (d *versionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data versionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location := data.Location.ValueString()

	client, err := d.bundle.NewPgSQLV2Client(location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create PostgreSQL v2 client", err.Error())
		return
	}

	versionList, _, err := client.ListVersions(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to list PostgreSQL v2 versions in location: %s", location), err.Error())
		return
	}

	var versions []versionModel
	for _, v := range versionList.Items {
		var item versionModel
		resp.Diagnostics.Append(mapVersionResponseToModel(&v, &item)...)
		if resp.Diagnostics.HasError() {
			return
		}
		versions = append(versions, item)
	}

	data.Versions = versions

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapVersionResponseToModel maps the API version response to the data source model.
func mapVersionResponseToModel(version *pgsqlv2.PostgresVersionRead, model *versionModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	model.ID = types.StringValue(version.Id)
	props := &version.Properties
	if props.Version != nil {
		model.Version = types.StringValue(*props.Version)
	}
	if props.Status != nil {
		model.Status = types.StringValue(*props.Status)
	}
	if props.Comment != nil {
		model.Comment = types.StringValue(*props.Comment)
	}
	if props.CanUpgradeTo != nil {
		values := make([]attr.Value, len(props.CanUpgradeTo))
		for i, v := range props.CanUpgradeTo {
			values[i] = types.StringValue(v)
		}
		model.CanUpgradeTo, diagnostics = types.ListValue(types.StringType, values)
	}
	return diagnostics
}
