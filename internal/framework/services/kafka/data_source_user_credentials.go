package kafka

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	kafkaSDK "github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/utils/validators"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	kafkaService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var _ datasource.DataSourceWithConfigure = (*userCredentialsDataSource)(nil)
var _ datasource.DataSourceWithConfigValidators = (*userCredentialsDataSource)(nil)

type userCredentialsDataSource struct {
	client *kafkaService.Client
}

type userCredentialsDataSourceModel struct {
	ClusterID            types.String   `tfsdk:"cluster_id"`
	ID                   types.String   `tfsdk:"id"`
	Username             types.String   `tfsdk:"username"`
	Location             types.String   `tfsdk:"location"`
	CertificateAuthority types.String   `tfsdk:"certificate_authority"`
	PrivateKey           types.String   `tfsdk:"private_key"`
	Certificate          types.String   `tfsdk:"certificate"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

// NewUserCredentialsDataSource creates a new user credentials data source.
func NewUserCredentialsDataSource() datasource.DataSource {
	return &userCredentialsDataSource{}
}

// Metadata returns the metadata for the data source.
func (d *userCredentialsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafka_user_credentials"
}

// Configure configures the data source client.
func (d *userCredentialsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *userCredentialsDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("username"),
			path.MatchRoot("id"),
		),
	}
}

// Schema returns the schema for the data source.
func (d *userCredentialsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Description: "The ID of the Kafka cluster",
				Required:    true,
				Validators: []validator.String{
					validators.UUIDValidator{},
				},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the Kafka user",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDValidator{},
				},
			},
			"username": schema.StringAttribute{
				Description: "The name of the Kafka user",
				Optional:    true,
				Computed:    true,
			},
			"location": schema.StringAttribute{
				Description: "The location of the Kafka user",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(kafkaService.AvailableLocations...),
				},
			},
			"certificate_authority": schema.StringAttribute{
				Description: "PEM for the certificate authority.",
				Computed:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "PEM for the private key.",
				Computed:    true,
			},
			"certificate": schema.StringAttribute{
				Description: "PEM for the certificate.",
				Computed:    true,
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *userCredentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Kafka API client", "Expected configured Kafka client. Please report this issue to the provider developers.")
		return
	}

	var data userCredentialsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := data.Timeouts.Read(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	clusterID := data.ClusterID.ValueString()
	userID := data.ID.ValueString()
	username := data.Username.ValueString()
	location := data.Location.ValueString()

	var userCredentials kafkaSDK.UserReadAccess
	var err error
	// TODO -- Move the GET logic inside another function
	if userID != "" {
		userCredentials, _, err = d.client.GetUserCredentialsByID(ctx, clusterID, userID, location)
		if err != nil {
			resp.Diagnostics.AddError("API Error Reading Kafka User Credentials", fmt.Sprintf("Failed to retrieve user credentials for user with ID: %s, cluster ID: %s, error: %s", userID, clusterID, err))
			return
		}
	} else if username != "" {
		userCredentials, _, err = d.client.GetUserCredentialsByName(ctx, clusterID, username, location)
		if err != nil {
			resp.Diagnostics.AddError("API Error Reading Kafka User Credentials", fmt.Sprintf("Failed to retrieve user credentials for user with name: %s, cluster ID: %s, error: %s", username, clusterID, err))
			return
		}
	}

	if hasMissingData(userCredentials) {
		resp.Diagnostics.AddError("Invalid API Response Format Kafka User Credentials", fmt.Sprintf("Expected valid string values in the API response but received 'nil' instead, user ID: %s, cluster ID: %s", userID, clusterID))
		return
	}
	populateUserCredentialsDataSourceModel(&data, userCredentials)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// populateUserCredentialsDataSourceModel populates the user credentials ephemeral model with information retrieved from the API.
func populateUserCredentialsDataSourceModel(data *userCredentialsDataSourceModel, userCredentials kafkaSDK.UserReadAccess) {
	data.CertificateAuthority = types.StringValue(*userCredentials.Metadata.CertificateAuthority)
	data.PrivateKey = types.StringValue(*userCredentials.Metadata.PrivateKey)
	data.Certificate = types.StringValue(*userCredentials.Metadata.Certificate)
	// TODO -- Check if these two need to be here or in another place
	data.ID = types.StringValue(userCredentials.Id)
	data.Username = types.StringValue(userCredentials.Properties.Name)
}
