package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/services/monitoring"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/services/objectstorage"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/services/objectstoragemanagement"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	apiGatewayService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/apigateway"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/autoscaling"
	cdnService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	dataplatformService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/inmemorydb"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadb"
	dnsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
	kafkaService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
	loggingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/logging"
	monitoringService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/monitoring"
	nfsService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/nfs"
	objectStorageService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstorage"
	objectStorageManagementService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstoragemanagement"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/vpn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// ClientOptions is the configuration for the provider.
type ClientOptions struct {
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	Token          types.String `tfsdk:"token"`
	Endpoint       types.String `tfsdk:"endpoint"`
	ContractNumber types.String `tfsdk:"contract_number"`
	S3SecretKey    types.String `tfsdk:"s3_secret_key"`
	S3AccessKey    types.String `tfsdk:"s3_access_key"`
	S3Region       types.String `tfsdk:"s3_region"`
	Retries        types.Int64  `tfsdk:"retries"`
	Insecure       types.Bool   `tfsdk:"insecure"`
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
			"s3_secret_key": schema.StringAttribute{
				Optional:    true,
				Description: "Secret key for IONOS Object Storage operations.",
			},
			"s3_access_key": schema.StringAttribute{
				Optional:    true,
				Description: "Access key for IONOS Object Storage operations.",
			},
			"s3_region": schema.StringAttribute{
				Optional:    true,
				Description: "Region for IONOS Object Storage operations.",
			},
			"insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "This field is to be set only for testing purposes. It is not recommended to set this field in production environments.",
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

	if clientOpts.S3SecretKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("secret_key"), "Unknown IONOS secret key", "")
	}

	if clientOpts.S3AccessKey.IsUnknown() {
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
	endpoint := os.Getenv("IONOS_API_URL")
	terraformVersion := req.TerraformVersion
	version := ionoscloud.Version

	insecureStr := os.Getenv("IONOS_ALLOW_INSECURE")
	insecureBool := false
	if insecureStr != "" {
		boolValue, err := strconv.ParseBool(insecureStr)
		if err != nil {
			log.Fatal(err)
		}
		insecureBool = boolValue
	}
	if !clientOpts.Token.IsNull() {
		token = clientOpts.Token.ValueString()
	}

	if !clientOpts.Username.IsNull() {
		username = clientOpts.Username.ValueString()
	}

	if !clientOpts.Password.IsNull() {
		password = clientOpts.Password.ValueString()
	}

	if !clientOpts.S3AccessKey.IsNull() {
		accessKey = clientOpts.S3AccessKey.ValueString()
	}

	if !clientOpts.S3SecretKey.IsNull() {
		secretKey = clientOpts.S3SecretKey.ValueString()
	}

	if !clientOpts.S3Region.IsNull() {
		region = clientOpts.S3Region.ValueString()
	}

	if !clientOpts.Endpoint.IsNull() {
		endpoint = clientOpts.Endpoint.ValueString()
	}
	if !clientOpts.Insecure.IsNull() {
		insecureBool = clientOpts.Insecure.ValueBool()
	}

	if token == "" && (username == "" || password == "") {
		resp.Diagnostics.AddError("missing credentials", "either token or username and password must be set")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	cleanedEndpoint := utils.CleanURL(endpoint)

	if insecureBool == true {
		resp.Diagnostics.AddWarning("insecure mode enabled", "This is not recommended for production environments.")
	}

	client := &services.SdkBundle{
		CDNClient:          cdnService.NewCDNClient(username, password, token, endpoint, version, terraformVersion, insecureBool),
		AutoscalingClient:  autoscalingService.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		CertManagerClient:  cert.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		CloudApiClient:     newCloudapiClient(username, password, token, endpoint, "DEV", terraformVersion, insecureBool),
		ContainerClient:    crService.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		DataplatformClient: dataplatformService.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		DNSClient:          dnsService.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		LoggingClient:      loggingService.NewClient(username, password, token, cleanedEndpoint, terraformVersion, insecureBool),
		MariaDBClient:      mariadb.NewMariaDBClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		MongoClient:        dbaasService.NewMongoClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		NFSClient:          nfsService.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		PsqlClient:         dbaasService.NewPsqlClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		KafkaClient:        kafkaService.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		APIGatewayClient: apiGatewayService.NewClient(
			username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool,
		),
		VPNClient:                     vpn.NewClient(username, password, token, cleanedEndpoint, terraformVersion, insecureBool),
		InMemoryDBClient:              inmemorydb.NewInMemoryDBClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		S3Client:                      objectStorageService.NewClient(accessKey, secretKey, region, endpoint, insecureBool),
		ObjectStorageManagementClient: objectStorageManagementService.NewClient(username, password, token, cleanedEndpoint, version, terraformVersion, insecureBool),
		MonitoringClient:              monitoringService.NewClient(username, password, token, cleanedEndpoint, terraformVersion, insecureBool),
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func newCloudapiClient(username, password, token, endpoint, version, terraformVersion string, insecure bool) *ionoscloud.APIClient {
	newConfig := ionoscloud.NewConfiguration(username, password, token, endpoint)
	newConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, ionoscloud.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)
	if os.Getenv(constant.IonosDebug) != "" {
		newConfig.Debug = true
	}
	newConfig.MaxRetries = constant.MaxRetries
	newConfig.WaitTime = constant.MaxWaitTime
	newConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}
	client := ionoscloud.NewAPIClient(newConfig)
	return client
}

// Resources returns the resources for the provider.
func (p *IonosCloudProvider) Resources(_ context.Context) []func() resource.Resource {
	var finalResult []func() resource.Resource
	resources := [][]func() resource.Resource{
		objectstorage.Resources(),
		objectstoragemanagement.Resources(),
		monitoring.Resources(),
	}

	for _, r := range resources {
		finalResult = append(finalResult, r...)
	}

	return finalResult
}

// DataSources returns the data sources for the provider.
func (p *IonosCloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	var finalResult []func() datasource.DataSource
	dataSources := [][]func() datasource.DataSource{
		objectstorage.DataSources(),
		objectstoragemanagement.DataSources(),
		monitoring.DataSources(),
	}

	for _, r := range dataSources {
		finalResult = append(finalResult, r...)
	}

	return finalResult
}
