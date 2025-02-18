package fileconfiguration

import (
	"crypto/x509"
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// products that do not have a location and will override the endpoint that is used globally
const (
	Autoscaling             = "autoscaling"
	APIGateway              = "apigateway"
	CDN                     = "cdn"
	Cert                    = "cert"
	Cloud                   = "cloud"
	Dataplatform            = "dataplatform"
	DNS                     = "dns"
	Mongo                   = "mongo"
	ObjectStorageManagement = "objectstoragemanagement"
	PSQL                    = "psql"
)

// products that have a location and will override the endpoint that is for each location
const (
	ContainerRegistry = "containerregistry"
	InMemoryDB        = "inmemorydb"
	Kafka             = "kafka"
	Logging           = "logging"
	Mariadb           = "mariadb"
	Monitoring        = "monitoring"
	NFS               = "nfs"
	ObjectStorage     = "objectstorage"
	VPN               = "vpn"
)

// ClientOverrideOptions is a struct that represents the client override options
type ClientOverrideOptions struct {
	// Endpoint is the endpoint that will be overridden
	Endpoint string
	// SkipTLSVerify skips tls verification. Not recommended for production!
	SkipTLSVerify bool
	// Certificate is the certificate that will be used for tls verification
	Certificate string
	// Credentials are the credentials that will be used for authentication
	Credentials Credentials
}

// Credentials are the credentials that will be used for authentication
type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
}

// Endpoint is a struct that represents an endpoint
type Endpoint struct {
	// the location or the region
	// Products that do not have a location and will override the endpoint that is used globally:
	// cloud, objectstoragemanagement, kafka, dns, mongo, psql, dataplatform, creg, autoscaling, apigateway
	// Products that have location-based endpoints: logging, monitoring, containerregistry, vpn, inmemorydb, nfs, objectstorage, mariadb
	Location            string `yaml:"location"`
	Name                string `yaml:"name"`
	SkipTLSVerify       bool   `yaml:"skipTlsVerify"`
	CertificateAuthData string `yaml:"certificateAuthData,omitempty"`
}

// Product is a struct that represents a product
type Product struct {
	// Name is the name of the product
	Name      string     `yaml:"name"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

// Environment is a struct that represents an environment
type Environment struct {
	Name string `yaml:"name"`
	// CertificateAuthData
	CertificateAuthData string `yaml:"certificateAuthData,omitempty"`
	// Products is a list of ionos products for which we will override endpoint, tls verification
	Products []Product `yaml:"products"`
}

// Profiles wrapper to read only the profiles from the config file
type Profiles struct {
	// CurrentProfile active profile for configuration
	CurrentProfile string `yaml:"currentProfile"`
	// Profiles
	Profiles []Profile `yaml:"profiles"`
}

// Profile is a struct that represents a profile and it's Credentials
type Profile struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Credentials Credentials
}

// FileConfig is a struct that represents the loaded configuration
type FileConfig struct {
	// Version of the configuration
	Version float64 `yaml:"version"`
	// CurrentProfile active profile for configuration
	CurrentProfile string `yaml:"currentProfile"`
	// Profiles list of profiles
	Profiles []Profile `yaml:"profiles"`
	// Environments list of environments
	Environments []Environment `yaml:"environments"`
}

// DefaultConfigFileName returns the default file path for the loaded configuration
func DefaultConfigFileName() (string, error) {
	home := ""
	var err error
	if home, err = os.UserHomeDir(); err != nil {
		return home, err
	}
	if home == "" {
		return home, fmt.Errorf("could not determine home directory")
	}
	return filepath.Join(home, ".ionos", "config"), nil

}

func findConfigFile() string {
	filePath := ""
	filePath = os.Getenv(shared.IonosFilePathEnvVar)
	var err error
	if filePath == "" {
		if filePath, err = DefaultConfigFileName(); err != nil {
			return ""
		}
	}
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		return ""
	}
	return filePath
}

// ReadProfilesFromFile reads profiles from yaml file, loads it into a struct and returns it
func ReadProfilesFromFile() *Profiles {
	filePath := findConfigFile()
	if filePath == "" {
		return nil
	}
	var content []byte
	var err error
	if content, err = os.ReadFile(filePath); err != nil {
		return nil
	}
	loadedProfiles := &Profiles{}
	_ = yaml.Unmarshal(content, loadedProfiles)
	return loadedProfiles
}

// ReadConfigFromFile reads yaml file, loads it into a struct and returns it
// IONOS_CONFIG_FILE environment variable can be set to point to the file to be loaded
func ReadConfigFromFile() (*FileConfig, error) {
	filePath := ""
	filePath = os.Getenv(shared.IonosFilePathEnvVar)
	var err error
	if filePath == "" {
		if filePath, err = DefaultConfigFileName(); err != nil {
			return nil, err
		}
	}
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", filePath)
	}
	var content []byte
	if content, err = os.ReadFile(filePath); err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, fmt.Errorf("file %s is empty", filePath)
	}
	loadedConfig := &FileConfig{}
	if err = yaml.Unmarshal(content, loadedConfig); err != nil {
		return nil, fmt.Errorf("while unmarshalling file %s, %w", filePath, err)
	}
	if os.Getenv(shared.IonosCurrentProfileEnvVar) != "" {
		loadedConfig.CurrentProfile = os.Getenv(shared.IonosCurrentProfileEnvVar)
	}
	return loadedConfig, nil
}

// GetCurrentProfile returns the current profile from the loaded configuration
// if the current profile is not set, it returns nil
// if the current profile is set and found in the loaded configuration, it returns the profile
func (fileConfig *FileConfig) GetCurrentProfile() *Profile {
	currentProfile := os.Getenv(shared.IonosCurrentProfileEnvVar)
	if currentProfile == "" {
		currentProfile = fileConfig.CurrentProfile
	}
	if currentProfile == "" {
		shared.SdkLogger.Printf("[WARN] no current profile set")
		return nil
	}
	for _, profile := range fileConfig.Profiles {
		if profile.Name == currentProfile {
			if shared.SdkLogLevel.Satisfies(shared.Debug) {
				shared.SdkLogger.Printf("[DEBUG] using profile %s", profile.Name)
			}
			return &profile
		}
	}
	if shared.SdkLogLevel.Satisfies(shared.Debug) {
		shared.SdkLogger.Printf("[WARN] no profile found for %s", currentProfile)
	}
	return nil
}

func (fileConfig *FileConfig) GetEnvForCurrentProfile() string {
	if fileConfig == nil {
		return ""
	}
	if currentProfile := fileConfig.GetCurrentProfile(); currentProfile != nil {
		return currentProfile.Environment
	}
	return ""
}

// GetProductOverrides returns the overrides for a specific product for the current environment
// if no current environment is found, the first environment is used for the product that matches productName is returned
func (fileConfig *FileConfig) GetProductOverrides(productName string) *Product {
	if fileConfig == nil {
		return nil
	}
	if productName == "" {
		if shared.SdkLogLevel.Satisfies(shared.Debug) {
			shared.SdkLogger.Printf("[DEBUG] cannot get overrides as product name is empty")
		}
		return nil
	}
	currentEnv := fileConfig.GetEnvForCurrentProfile()
	for _, environment := range fileConfig.Environments {
		if currentEnv != "" && environment.Name != currentEnv {
			continue
		}
		for _, product := range environment.Products {
			if product.Name == productName {
				return &product
			}
		}
	}
	if shared.SdkLogLevel.Satisfies(shared.Debug) {
		shared.SdkLogger.Printf("[DEBUG] no environment overrides found for product %s", productName)
	}
	return nil
}

func (fileConfig *FileConfig) GetProductLocationOverrides(productName, location string) *Endpoint {
	if fileConfig == nil {
		return nil
	}
	product := fileConfig.GetProductOverrides(productName)
	if product == nil || len(product.Endpoints) == 0 {
		return nil
	}
	for _, endpoint := range product.Endpoints {
		if endpoint.Location == location {
			return &endpoint
		}
	}
	if shared.SdkLogLevel.Satisfies(shared.Debug) {
		shared.SdkLogger.Printf("[DEBUG] no endpoint overrides found for product %s and location %s", productName, location)
	}
	return nil
}

// AddCertsToClient adds certificates to the http client
func AddCertsToClient(httpClient *http.Client, authorityData string) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if ok := rootCAs.AppendCertsFromPEM([]byte(authorityData)); !ok && shared.SdkLogLevel.Satisfies(shared.Debug) {
		shared.SdkLogger.Printf("No certs appended, using system certs only")
	}
	httpClient.Transport.(*http.Transport).TLSClientConfig.RootCAs = rootCAs
}
