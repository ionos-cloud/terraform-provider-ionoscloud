package loadedconfig

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
)

// TerraformToSDK maps the Terraform location to the SDK location
// todo check if we remove this
var TerraformToSDK = map[string]string{
	"de/fra": "de-fra",
	"de/txl": "de-txl",
	"es/vit": "es-vit",
	"gb/lhr": "gb-lhr",
	"fr/par": "fr-par",
	"us/las": "us-las",
	"us/ewr": "us-ewr",
	"us/mci": "us-mci",
}

type fileConfigProvider interface {
	GetFileConfig() *fileconfiguration.FileConfig
}

// ConfigProviderWithLoader is a shared interface for all services that use the loaded config and also have a shared config
type ConfigProviderWithLoader interface {
	fileConfigProvider
	shared.ConfigProvider
}

// ConfigProviderWithLoaderAndLocation is a shared interface for all services that use the loaded config, have a shared confi
// and also need to change the config URL based on location
type ConfigProviderWithLoaderAndLocation interface {
	ConfigProviderWithLoader
	ChangeConfigURL(ctx context.Context, location string)
}

// SetClientOptionsFromConfig overrides the client configuration with the loaded config
// if the product and location are found in the loaded config
// Any changes here should be reflected in the service overrideClientEndpoint functions for the sdks not using bundle
func SetClientOptionsFromConfig(ctx context.Context, client ConfigProviderWithLoaderAndLocation, productName, location string) {
	// whatever is set, at the end we need to check if the IONOS_API_URL_productname is set and override the endpoint
	defer client.ChangeConfigURL(ctx, location)
	// do not set from config if we use IONOS_API_URL
	if os.Getenv(shared.IonosApiUrlEnvVar) != "" {
		tflog.Debug(ctx, "using endpoint from env", map[string]interface{}{"product": productName, "env": shared.IonosApiUrlEnvVar, "url": os.Getenv(shared.IonosApiUrlEnvVar)})
		return
	}
	fileConfig := client.GetFileConfig()
	if fileConfig == nil {
		return
	}
	config := client.GetConfig()
	if config == nil {
		return
	}
	endpoint := fileConfig.GetProductLocationOverrides(productName, location)
	if endpoint == nil {
		tflog.Warn(ctx, "missing endpoint", map[string]interface{}{"product": productName, "location": location})
		return
	}
	config.Servers = shared.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	config.HTTPClient = &http.Client{}
	config.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)
	if endpoint.SkipTLSVerify || endpoint.CertificateAuthData != "" {
		tflog.Debug(ctx, "endpoint TLS config", map[string]interface{}{"product": productName, "skip_tls_verify": endpoint.SkipTLSVerify, "has_cert_auth_data": endpoint.CertificateAuthData != "", "cert_auth_data_len": len(endpoint.CertificateAuthData)})
	}
}

// SetGlobalClientOptionsFromFileConfig sets the client options from the loaded config if not already set
// mutates clientOptions. Should only be used if the product does not have location overrides
func SetGlobalClientOptionsFromFileConfig(ctx context.Context, clientOptions *clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig, productName string) {
	if clientOptions == nil || fileConfig == nil {
		return
	}
	productOverrides := fileConfig.GetProductOverrides(productName)
	if productOverrides == nil || len(productOverrides.Endpoints) == 0 {
		tflog.Warn(ctx, "missing config for product", map[string]interface{}{"product": productName})
		return
	}
	if len(productOverrides.Endpoints) > 1 {
		tflog.Warn(ctx, "multiple endpoints found for product, using the first one", map[string]interface{}{"product": productOverrides.Name})
	}

	if !clientOptions.SkipTLSVerify {
		clientOptions.SkipTLSVerify = productOverrides.Endpoints[0].SkipTLSVerify
		if clientOptions.SkipTLSVerify {
			tflog.Debug(ctx, "file config TLS", map[string]interface{}{"product": productName, "skip_tls_verify": clientOptions.SkipTLSVerify})
		}
	}
	if clientOptions.Endpoint == "" {
		clientOptions.Endpoint = productOverrides.Endpoints[0].Name
		tflog.Debug(ctx, "file config global endpoint", map[string]interface{}{"product": productName, "endpoint": clientOptions.Endpoint})
	}
	if productOverrides.Endpoints[0].CertificateAuthData != "" {
		tflog.Debug(ctx, "file config certificateAuthData present", map[string]interface{}{"product": productName, "cert_auth_data_len": len(productOverrides.Endpoints[0].CertificateAuthData)})
	}
	clientOptions.Certificate = productOverrides.Endpoints[0].CertificateAuthData
}
