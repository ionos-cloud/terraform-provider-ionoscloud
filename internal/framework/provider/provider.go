package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/services/s3"
	s3service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/s3"
)

// ClientOptions is the configuration for the provider.
type ClientOptions struct {
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	Token          types.String `tfsdk:"token"`
	Endpoint       types.String `tfsdk:"endpoint"`
	ContractNumber types.String `tfsdk:"contract_number"`
	SecretKey      types.String `tfsdk:"secret_key"`
	AccessKey      types.String `tfsdk:"access_key"`
	Region         types.String `tfsdk:"region"`
	Retries        types.Int64  `tfsdk:"retries"`
}

// IonosCloudProvider is the provider implementation.
type IonosCloudProvider struct {
}

// New creates a new provider.
func New() provider.Provider {
	return &IonosCloudProvider{}
}

// Metadata returns the metadata for the provider.
func (p *IonosCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ionoscloud"
}

// Schema returns the schema for the provider.
func (p *IonosCloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "IonosCloud username for API operations. If token is provided, token is preferred",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Description: "IonosCloud password for API operations. If token is provided, token is preferred",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Description: "IonosCloud bearer token for API operations.",
			},
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "IonosCloud REST API URL. Usually not necessary to be set, SDKs know internally how to route requests to the API.",
			},
			"retries": schema.Int64Attribute{
				Optional:           true,
				DeprecationMessage: "Timeout is used instead of this functionality",
			},
			"contract_number": schema.StringAttribute{
				Optional:    true,
				Description: "To be set only for reseller accounts. Allows to run terraform on a contract number under a reseller account.",
			},
			"secret_key": schema.StringAttribute{
				Optional:    true,
				Description: "Secret key for IONOS S3 bucket operations.",
			},
			"access_key": schema.StringAttribute{
				Optional:    true,
				Description: "Access key for IONOS S3 bucket operations.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "Region for IONOS S3 bucket operations.",
			},
		},
	}
}

// Configure configures the provider.
func (p *IonosCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var clientOpts ClientOptions
	diags := req.Config.Get(ctx, &clientOpts)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if clientOpts.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("token"), "Unknown IONOS token", "token must be set")
	}

	if clientOpts.SecretKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("secret_key"), "Unknown IONOS secret key", "")
	}

	if clientOpts.AccessKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("access_key"), "Unknown IONOS access key", "")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	token := os.Getenv("IONOS_TOKEN")
	username := os.Getenv("IONOS_USERNAME")
	password := os.Getenv("IONOS_PASSWORD")
	accessKey := os.Getenv("IONOS_S3_ACCESS_KEY")
	secretKey := os.Getenv("IONOS_S3_SECRET_KEY")
	region := os.Getenv("IONOS_S3_REGION")

	if !clientOpts.Token.IsNull() {
		token = clientOpts.Token.ValueString()
	}

	if !clientOpts.Username.IsNull() {
		username = clientOpts.Username.ValueString()
	}

	if !clientOpts.Password.IsNull() {
		password = clientOpts.Password.ValueString()
	}

	if !clientOpts.AccessKey.IsNull() {
		accessKey = clientOpts.AccessKey.ValueString()
	}

	if !clientOpts.SecretKey.IsNull() {
		secretKey = clientOpts.SecretKey.ValueString()
	}

	if !clientOpts.Region.IsNull() {
		region = clientOpts.Region.ValueString()
	}

	if accessKey == "" || secretKey == "" {
		resp.Diagnostics.AddError("s3 keys missing", "access_key and secret_key must be set")
	}

	if token == "" && (username == "" || password == "") {
		resp.Diagnostics.AddError("missing credentials", "either token or username and password must be set")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := s3service.NewClient(accessKey, secretKey, region)
	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources returns the resources for the provider.
func (p *IonosCloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		s3.NewBucketResource,
	}
}

// DataSources returns the data sources for the provider.
func (p *IonosCloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		s3.NewBucketDataSource,
	}
}
