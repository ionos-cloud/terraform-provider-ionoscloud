package mongov2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	mongov2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mongov2"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var _ datasource.DataSourceWithConfigure = (*templatesDataSource)(nil)

type templateDataSourceModel struct {
	Cores       types.Int32  `tfsdk:"cores"`
	Edition     types.String `tfsdk:"edition"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	RAM         types.Int32  `tfsdk:"ram"`
	StorageSize types.Int32  `tfsdk:"storage_size"`
}

func templateDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cores":        schema.Int32Attribute{Computed: true, Description: "The number of CPU cores."},
		"edition":      schema.StringAttribute{Computed: true, Description: "The edition the template applies to."},
		"id":           schema.StringAttribute{Computed: true, Description: "The ID (UUID) of the instance template."},
		"name":         schema.StringAttribute{Computed: true, Description: "The name of the instance template."},
		"ram":          schema.Int32Attribute{Computed: true, Description: "The amount of memory in gigabytes (GB)."},
		"storage_size": schema.Int32Attribute{Computed: true, Description: "The amount of storage in gigabytes (GB)."},
	}
}

type templatesDataSource struct {
	bundle *bundleclient.SdkBundle
}

type templatesDataSourceModel struct {
	Edition  types.String              `tfsdk:"edition"`
	Items    []templateDataSourceModel `tfsdk:"items"`
	Location types.String              `tfsdk:"location"`
}

// NewTemplatesDataSource creates a new data source for listing MongoDB v2 instance templates.
func NewTemplatesDataSource() datasource.DataSource {
	return &templatesDataSource{}
}

func (d *templatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongo_templates_v2"
}

func (d *templatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *templatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists IONOS CLOUD MongoDB V2 instance templates, with optional filter by edition.",
		Attributes: map[string]schema.Attribute{
			"edition": schema.StringAttribute{
				Optional:    true,
				Description: "Filter templates by edition (ENTERPRISE or BUSINESS).",
			},
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of instance templates.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: templateDataSourceAttributes(),
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The location to query. Available locations: " + mongov2service.AvailableLocationsString() + ".",
			},
		},
	}
}

func (d *templatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data templatesDataSourceModel
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

	list, apiResponse, err := client.ListTemplates(ctx, data.Edition.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to list MongoDB V2 instance templates", diagutil.WrapError(err, &diagutil.ErrorContext{
			StatusCode:     apiResponse.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	items := make([]templateDataSourceModel, 0, len(list.Items))
	for _, t := range list.Items {
		items = append(items, mapTemplateToModel(t))
	}
	data.Items = items

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapTemplateToModel(t mongov3.InstanceTemplateRead) templateDataSourceModel {
	item := templateDataSourceModel{
		ID:          types.StringValue(t.Id),
		Name:        types.StringPointerValue(t.Properties.Name),
		RAM:         types.Int32PointerValue(t.Properties.Ram),
		Cores:       types.Int32PointerValue(t.Properties.Cores),
		StorageSize: types.Int32PointerValue(t.Properties.StorageSize),
	}
	if t.Properties.Edition != nil {
		item.Edition = types.StringValue(string(*t.Properties.Edition))
	}
	return item
}
