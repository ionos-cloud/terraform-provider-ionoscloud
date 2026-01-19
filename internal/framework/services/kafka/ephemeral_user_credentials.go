package kafka

import (
	"context"
	"fmt"

	kafkaSDK "github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/ephemeral/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/utils/validators"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	kafkaService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var _ ephemeral.EphemeralResourceWithConfigure = (*userCredentialsEphemeral)(nil)

type userCredentialsEphemeral struct {
	client *kafkaService.Client
}

type userCredentialsEphemeralModel struct {
	ClusterID            types.String   `tfsdk:"cluster_id"`
	UserID               types.String   `tfsdk:"user_id"`
	Username             types.String   `tfsdk:"username"`
	Location             types.String   `tfsdk:"location"`
	CertificateAuthority types.String   `tfsdk:"certificate_authority"`
	PrivateKey           types.String   `tfsdk:"private_key"`
	Certificate          types.String   `tfsdk:"certificate"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

// NewUserCredentialsEphemeral creates a new user credentials ephemeral resource.
func NewUserCredentialsEphemeral() ephemeral.EphemeralResource {
	return &userCredentialsEphemeral{}
}

// Metadata returns the metadata for the user credentials ephemeral resource.
func (d *userCredentialsEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafka_user_credentials"
}

// Configure configures the user credentials ephemeral resource.
func (d *userCredentialsEphemeral) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

// Schema returns the schema for the user credentials ephemeral resource.
func (d *userCredentialsEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Description: "The ID of the Kafka cluster",
				Required:    true,
				Validators: []validator.String{
					validators.UUIDValidator{},
				},
			},
			"user_id": schema.StringAttribute{
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

// TODO -- Add logic for 'username', similar to the user credentials data source.
// Open creates a new ephemeral resource and populates it by using data fetched from the API
func (d *userCredentialsEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Kafka API client", "Expected configured Kafka client. Please report this issue to the provider developers.")
		return
	}

	var data userCredentialsEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := data.Timeouts.Open(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	clusterID := data.ClusterID.ValueString()
	userID := data.UserID.ValueString()
	location := data.Location.ValueString()

	userCredentials, _, err := d.client.GetUserCredentialsByID(ctx, clusterID, userID, location)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Kafka User Credentials", fmt.Sprintf("Failed to retrieve user credentials for user with ID: %s, cluster ID: %s, error: %s", userID, clusterID, err))
		return
	}
	if hasMissingData(userCredentials) {
		resp.Diagnostics.AddError("Invalid API Response Format Kafka User Credentials", fmt.Sprintf("Expected valid string values in the API response but received 'nil' instead, user ID: %s, cluster ID: %s", userID, clusterID))
		return
	}
	populateUserCredentialsEphemeralModel(&data, userCredentials)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}

// hasMissingData verifies if the API response contains nil values.
func hasMissingData(userCredentials kafkaSDK.UserReadAccess) bool {
	return userCredentials.Metadata.CertificateAuthority == nil || userCredentials.Metadata.PrivateKey == nil || userCredentials.Metadata.Certificate == nil
}

// populateUserCredentialsEphemeralModel populates the user credentials ephemeral model with information retrieved from the API.
func populateUserCredentialsEphemeralModel(data *userCredentialsEphemeralModel, userCredentials kafkaSDK.UserReadAccess) {
	data.CertificateAuthority = types.StringValue(*userCredentials.Metadata.CertificateAuthority)
	data.PrivateKey = types.StringValue(*userCredentials.Metadata.PrivateKey)
	data.Certificate = types.StringValue(*userCredentials.Metadata.Certificate)
}
