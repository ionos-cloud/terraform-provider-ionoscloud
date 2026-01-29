package kafka

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/validators"
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
	userCredentialsModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
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

	userCredentials, diags := getUserCredentials(ctx, *d.client, data.userCredentialsModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if hasMissingData(userCredentials) {
		resp.Diagnostics.AddError("Invalid API Response Format Kafka User Credentials", fmt.Sprintf("Expected valid string values in the API response but received 'nil' instead, user ID: %s, cluster ID: %s", userCredentials.Id, data.ClusterID.ValueString()))
		return
	}

	populateUserCredentialsModel(&data.userCredentialsModel, userCredentials)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
