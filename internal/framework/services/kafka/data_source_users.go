package kafka

import (
	"context"
	"fmt"

	kafkaSDK "github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/utils/validators"
	kafkaService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
)

var _ datasource.DataSourceWithConfigure = (*usersDataSource)(nil)

type usersDataSource struct {
	client *kafkaService.Client
}

type usersDataSourceModel struct {
	ClusterID types.String          `tfsdk:"cluster_id"`
	Location  types.String          `tfsdk:"location"`
	Users     []userDataSourceModel `tfsdk:"users"`
}

type userDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
}

// NewUsersDataSource creates a new users data source.
func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

// Metadata returns the metadata for the users data source.
func (d *usersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafka_users"
}

// Configure configures the data source.
func (d *usersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *bundleclient.Sdkbundle, got: %T. Please report this issue to the provider developers.", req.ProviderData))
	}
	d.client = clientBundle.KafkaClient
}

// Schema returns the schema for the data source.
func (d *usersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Description: "The ID of the Kafka cluster",
				Required:    true,
				Validators: []validator.String{
					validators.UUIDValidator{},
				},
			},
			"location": schema.StringAttribute{
				Description: "The location of the Kafka cluster",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(kafkaService.AvailableLocations...),
				},
			},
			"users": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the Kafka user",
							Computed:    true,
						},
						"username": schema.StringAttribute{
							Description: "The name of the Kafka user",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

// Read retrieves and sets the information about users in state.
func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Kafka API client", "Expected configured Kafka client. Please report this issue to the provider developers.")
	}

	var data usersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := data.ClusterID.ValueString()
	location := data.Location.ValueString()

	users, _, err := d.client.GetUsers(ctx, clusterID, location)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Kafka Users", fmt.Sprintf("Failed to retrieve the list of Kafka users for cluster with ID: %s, error: %s", clusterID, err))
		return
	}
	data.Users = buildUsersFromAPIResp(users)
	// TODO - Timeouts
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildUsersFromAPIResp converts the users info from the API response into a slice of user data source models.
func buildUsersFromAPIResp(users kafkaSDK.UserReadList) []userDataSourceModel {
	result := make([]userDataSourceModel, 0, len(users.Items))
	for _, user := range users.Items {
		result = append(result, buildUserModelFromAPIResp(user))
	}
	return result
}

// buildUserModelFromAPIResp converts the user info from the API response into a user data source model.
func buildUserModelFromAPIResp(user kafkaSDK.UserRead) userDataSourceModel {
	return userDataSourceModel{
		ID:       types.StringValue(user.Id),
		Username: types.StringValue(user.Properties.Name),
	}
}
