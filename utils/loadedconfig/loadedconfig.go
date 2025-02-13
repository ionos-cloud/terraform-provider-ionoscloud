package loadedconfig

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"log"
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

type LoadedConfigProvider interface {
	GetLoadedConfig() *shared.LoadedConfig
}

type ConfigProviderWithLoader interface {
	LoadedConfigProvider
	shared.ConfigProvider
}
type ConfigProviderWithLoaderAndLocation interface {
	ConfigProviderWithLoader
	ChangeConfigURL(location string)
}

// OverrideClientEndpoint overrides the client configuration with the loaded config
// if the product and location are found in the loaded config
func OverrideClientEndpoint(client ConfigProviderWithLoaderAndLocation, productName, location string) {
	//whatever is set, at the end we need to check if the IONOS_API_URL_productname is set and use override the endpoint if yes
	defer client.ChangeConfigURL(location)
	//if os.Getenv(ionoscloud.IonosApiUrlEnvVar) != "" {
	//	fmt.Printf("[DEBUG] Using custom endpoint %s\n", os.Getenv(ionoscloud.IonosApiUrlEnvVar))
	//	return
	//}
	loadedConfig := client.GetLoadedConfig()
	if loadedConfig == nil {
		return
	}
	config := client.GetConfig()
	if config == nil {
		return
	}
	endpoint := loadedConfig.GetProductLocationOverrides(productName, location)
	if endpoint == nil {
		log.Printf("[WARN] Missing endpoint for %s in location %s", productName, location)
		return
	}
	config.Servers = shared.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	if endpoint.SkipTLSVerify {
		config.HTTPClient.Transport = utils.CreateTransport(true)
	}
}

//func ChangeConfigURL(client ConfigProviderWithLoaderAndLocation, location string) {
//	config := client.GetConfig()
//	for _, server := range config.Servers {
//		if strings.EqualFold(server.Description, shared.EndpointOverridden+location) ||
//			strings.EqualFold(server.URL, locationToURL[location]) {
//			config.Servers = shared.ServerConfigurations{
//				{
//					URL:         server.URL,
//					Description: shared.EndpointOverridden + location,
//				},
//			}
//			return
//		}
//	}
//}

// SetClientOptionsFromLoadedConfig sets the client options from the loaded config if not already set
// mutates clientOptions. Should only be used if the product does not have location overrides
func SetClientOptionsFromLoadedConfig(clientOptions *bundle.ClientOptions, loadedConfig *shared.LoadedConfig, productName string) {
	if clientOptions == nil || loadedConfig == nil {
		return
	}
	productOverrides := loadedConfig.GetProductOverrides(productName)
	if productOverrides == nil || len(productOverrides.Endpoints) == 0 {
		log.Printf("[WARN] Missing config for %s", productName)
		return
	}
	if len(productOverrides.Endpoints) > 1 {
		log.Printf("[WARN] Multiple endpoints found for product %s, using the first one", productOverrides.Name)
	}

	if !clientOptions.SkipTLSVerify {
		clientOptions.SkipTLSVerify = productOverrides.Endpoints[0].SkipTLSVerify
	}
	if clientOptions.Endpoint == "" {
		clientOptions.Endpoint = productOverrides.Endpoints[0].Name
	}
}
