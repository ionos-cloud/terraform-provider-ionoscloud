package fileconfiguration

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// Usage:
//   // Load from default file
//   cfg, err := fileconfiguration.New("")
//   // From explicit path:
//   cfg, err := fileconfiguration.New("/path/to/config")
//   // From env:
//   cfg, err := fileconfiguration.NewFromEnv()
//
//   // Read profiles list:
//   profiles := fileconfiguration.ReadProfilesFromFile()
//   // Switch/get current profile:
//   prof := cfg.GetCurrentProfile()
//   // Get environment names:
//   envs := cfg.GetEnvironmentNames()
//   // Get profile names:
//   profs := cfg.GetProfileNames()
//
//   // Get endpoint for regionless API:
//   ep := cfg.GetProductOverrides("psql")
//   // Get endpoint for region-specific API:
//   ep := cfg.GetProductLocationOverrides("dns", "de/fra")
//   // Get endpoint for region-specific API with fallback to global endpoint: (combine the two funcs above)
//   ep := cfg.GetOverride("dns", "de/fra")
//   ep := cfg.GetOverride("psql", "")

// products that do not have a location and will override the endpoint that is used globally
const (
	Autoscaling             = "autoscaling"
	APIGateway              = "apigateway"
	CDN                     = "cdn"
	Cert                    = "cert"
	Cloud                   = "cloud"
	ContainerRegistry       = "containerregistry"
	Dataplatform            = "dataplatform"
	DNS                     = "dns"
	Mongo                   = "mongo"
	ObjectStorageManagement = "objectstoragemanagement"
	PSQL                    = "psql"
)

// products that have a location and will override the endpoint that is for each location
const (
	InMemoryDB    = "inmemorydb"
	Kafka         = "kafka"
	Logging       = "logging"
	Mariadb       = "mariadb"
	Monitoring    = "monitoring"
	NFS           = "nfs"
	ObjectStorage = "objectstorage"
	VPN           = "vpn"
)

// Endpoint is a struct that represents an endpoint
type Endpoint struct {
	// the location or the region
	// Products that do not have a location and will override the endpoint that is used globally:
	// cloud, objectstoragemanagement, kafka, dns, mongo, psql, dataplatform, creg, autoscaling, apigateway
	// Products that have location-based endpoints: logging, monitoring, containerregistry, vpn, inmemorydb, nfs, objectstorage, mariadb
	Location            string `yaml:"location,omitempty"`
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
	Credentials shared.Credentials
}

// Version wraps float64 so we can control its YAML output.
type Version float64

// MarshalYAML ensures that, e.g., 1.0 is emitted as "1.0" instead of "1".
func (v Version) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!float",
		Value: fmt.Sprintf("%.1f", float64(v)),
	}, nil
}

// FileConfig is a struct that represents the loaded configuration
type FileConfig struct {
	// Version of the configuration
	Version Version `yaml:"version"`
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

// New reads yaml file, loads it into a struct and returns it
// IONOS_CONFIG_FILE environment variable can be set to point to the file to be loaded
func New(filePath string) (*FileConfig, error) {
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

func NewFromEnv() (*FileConfig, error) {
	return New(os.Getenv(shared.IonosFilePathEnvVar))
}

// GetProfileNames returns a list of profile names from the loaded configuration
func (f *FileConfig) GetProfileNames() []string {
	names := make([]string, len(f.Profiles))
	for i, p := range f.Profiles {
		names[i] = p.Name
	}
	return names
}

// GetEnvironmentNames returns a list of environment names from the loaded configuration
func (f *FileConfig) GetEnvironmentNames() []string {
	names := make([]string, len(f.Environments))
	for i, e := range f.Environments {
		names[i] = e.Name
	}
	return names
}

// GetOverride returns the endpoint for a specific product and location
// with fallback to the global endpoint if no location is found.
//
// It is a helper function combining GetProductLocationOverrides and GetProductOverrides
func (f *FileConfig) GetOverride(productName, location string) *Endpoint {
	if locEp := f.GetProductLocationOverrides(productName, location); locEp != nil {
		return locEp
	}

	if prod := f.GetProductOverrides(productName); prod != nil && len(prod.Endpoints) > 0 {
		if prod.Endpoints[0].Location != "" && prod.Endpoints[0].Location != location {
			// Check if we actually got a location-specific endpoint (e.g. the user asked for a wrong location
			// and GetProductOverrides returned the first location-specific endpoint)
			if shared.SdkLogLevel.Satisfies(shared.Debug) {
				shared.SdkLogger.Printf("[DEBUG] Retrieved location-specific (%s) override '%s' for product '%s' "+
					"when a location-less override was expected, discarding...", location, prod.Endpoints[0].Name, productName)
			}
			return nil
		}
		return &prod.Endpoints[0]
	}
	return nil
}

// GetCurrentProfile returns the current profile from the loaded configuration
// if the current profile is not set, it returns nil
// if the current profile is set and found in the loaded configuration, it returns the profile
func (f *FileConfig) GetCurrentProfile() *Profile {
	if f == nil {
		return nil
	}

	currentProfile := os.Getenv(shared.IonosCurrentProfileEnvVar)
	if currentProfile == "" {
		currentProfile = f.CurrentProfile
	}
	if currentProfile == "" {
		shared.SdkLogger.Printf("[WARN] no current profile set")
		return nil
	}
	for _, profile := range f.Profiles {
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

func (f *FileConfig) GetEnvForCurrentProfile() string {
	if f == nil {
		return ""
	}
	if currentProfile := f.GetCurrentProfile(); currentProfile != nil {
		return currentProfile.Environment
	}
	return ""
}

// GetProductOverrides returns the overrides for a specific product for the current environment
// if no current environment is found, the first environment is used for the product that matches productName is returned
func (f *FileConfig) GetProductOverrides(productName string) *Product {
	if f == nil {
		return nil
	}
	if productName == "" {
		if shared.SdkLogLevel.Satisfies(shared.Debug) {
			shared.SdkLogger.Printf("[DEBUG] cannot get overrides as product name is empty")
		}
		return nil
	}
	currentEnv := f.GetEnvForCurrentProfile()
	for _, environment := range f.Environments {
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

// GetProductLocationOverrides returns the overrides for a specific product and location for the current environment
func (f *FileConfig) GetProductLocationOverrides(productName, location string) *Endpoint {
	if f == nil {
		return nil
	}
	product := f.GetProductOverrides(productName)
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
