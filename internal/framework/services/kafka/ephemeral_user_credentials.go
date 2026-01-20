package kafka

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/ephemeral/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/ephemeralvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/utils/validators"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	kafkaService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var _ ephemeral.EphemeralResourceWithConfigure = (*userCredentialsEphemeral)(nil)
var _ ephemeral.EphemeralResourceWithConfigValidators = (*userCredentialsEphemeral)(nil)

type userCredentialsEphemeral struct {
	client *kafkaService.Client
}

type userCredentialsEphemeralModel struct {
	userCredentialsModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
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

func (d *userCredentialsEphemeral) ConfigValidators(ctx context.Context) []ephemeral.ConfigValidator {
	return []ephemeral.ConfigValidator{
		ephemeralvalidator.ExactlyOneOf(
			path.MatchRoot("username"),
			path.MatchRoot("id"),
		),
	}
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
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
