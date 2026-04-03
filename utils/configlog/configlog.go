package configlog

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
)

// LoadFileConfigWithLogging wraps fileconfiguration.NewFromEnv() with pre/post logging
// so that users can trace config file loading with TF_LOG=DEBUG.
func LoadFileConfigWithLogging() (*fileconfiguration.FileConfig, error) {
	// Resolve path for logging
	filePath := os.Getenv(shared.IonosFilePathEnvVar)
	source := shared.IonosFilePathEnvVar
	if filePath == "" {
		source = "default"
		defaultPath, err := fileconfiguration.DefaultConfigFileName()
		if err != nil {
			log.Printf("[DEBUG] Could not determine default config file path: %s", err)
		} else {
			filePath = defaultPath
		}
	}

	if filePath != "" {
		var status string
		info, err := os.Stat(filePath) //nolint:gosec // G703 - path from user's own env var
		if err == nil {
			status = "found"
			if f, readErr := os.Open(filePath); readErr != nil { //nolint:gosec // G304 - path from user's own env var
				status = fmt.Sprintf("found but unreadable (permissions: %04o)", info.Mode().Perm())
			} else {
				f.Close()
			}
		} else if os.IsNotExist(err) {
			status = "not found"
		}
		log.Printf("[DEBUG] Config file: %s (source: %s, status: %s)", filePath, source, status)
	}

	fileConfig, err := fileconfiguration.NewFromEnv()
	if err != nil {
		log.Printf("[DEBUG] Config file not loaded: %s", err)
		return nil, err
	}

	log.Printf("[DEBUG] Config file loaded successfully (version: %.1f)", float64(fileConfig.Version))
	logProfileAndEnvironment(fileConfig)
	logFileConfigEndpoints(fileConfig)

	return fileConfig, nil
}

// logProfileAndEnvironment logs profile and environment resolution details.
func logProfileAndEnvironment(fileConfig *fileconfiguration.FileConfig) {
	var parts []string
	parts = append(parts, fmt.Sprintf("%d profile(s), currentProfile: %q", len(fileConfig.Profiles), fileConfig.CurrentProfile))

	envOverride := os.Getenv(shared.IonosCurrentProfileEnvVar)
	if envOverride != "" {
		parts = append(parts, fmt.Sprintf("%s overrides to %q", shared.IonosCurrentProfileEnvVar, envOverride))
	}

	profile := fileConfig.GetCurrentProfile()
	if profile != nil {
		parts = append(parts, fmt.Sprintf("active: %q (environment: %q)", profile.Name, profile.Environment))
	} else {
		currentProfile := envOverride
		if currentProfile == "" {
			currentProfile = fileConfig.CurrentProfile
		}
		if currentProfile != "" {
			parts = append(parts, fmt.Sprintf("no matching profile for %q, available: [%s]", currentProfile, strings.Join(fileConfig.GetProfileNames(), ", ")))
		}
	}

	log.Printf("[DEBUG] Profile resolution: %s", strings.Join(parts, " | "))
}

// logFileConfigEndpoints logs product and endpoint counts from the active environment in the file config.
func logFileConfigEndpoints(fileConfig *fileconfiguration.FileConfig) {
	failoverOpts := fileConfig.GetFailoverOptions()
	if failoverOpts != nil {
		log.Printf("[DEBUG] Failover config: strategy=%q", failoverOpts.Strategy)
	} else {
		log.Printf("[DEBUG] Failover config: not set (default: none)")
	}

	envName := fileConfig.GetEnvForCurrentProfile()
	if envName == "" {
		return
	}

	for _, env := range fileConfig.Environments {
		if env.Name == envName {
			var products []string
			for _, product := range env.Products {
				products = append(products, fmt.Sprintf("%s(%d)", product.Name, len(product.Endpoints)))
			}
			log.Printf("[DEBUG] Environment %q: %d product(s): %s", env.Name, len(env.Products), strings.Join(products, ", "))
			return
		}
	}
}

// LogCredentialResolution logs which credentials were found and from where.
func LogCredentialResolution(creds shared.Credentials, fileConfigUsed bool, profileName string) {
	foundStr := func(name string, found bool) string {
		if found {
			return name + "=found"
		}
		return name + "=not found"
	}

	line := fmt.Sprintf("Credentials: %s, %s, %s",
		foundStr("token", creds.Token != ""),
		foundStr("user/pass", creds.Username != "" && creds.Password != ""),
		foundStr("S3 keys", creds.S3AccessKey != "" && creds.S3SecretKey != ""),
	)

	if fileConfigUsed && profileName != "" {
		var found []string
		if creds.Token != "" {
			found = append(found, "token")
		}
		if creds.Username != "" && creds.Password != "" {
			found = append(found, "user+pass")
		}
		if creds.S3AccessKey != "" && creds.S3SecretKey != "" {
			found = append(found, "S3 keys")
		}
		if len(found) == 0 {
			found = append(found, "none")
		}
		line += fmt.Sprintf(" | file config profile %q: %s", profileName, strings.Join(found, ", "))
	}

	if creds.Token != "" && creds.Username != "" && creds.Password != "" {
		line += " | both token and user/pass provided; token takes precedence"
	}

	if creds.Token != "" {
		line += " | authenticating via token"
	} else if creds.Username != "" && creds.Password != "" {
		line += " | authenticating via user/pass"
	}

	log.Printf("[DEBUG] %s", line)
}

// LogEndpointEnvVars logs which product-specific endpoint env vars are set.
func LogEndpointEnvVars() {
	envVars := map[string]string{
		"IONOS_API_URL":                           "global",
		"IONOS_API_URL_VPN":                       "VPN",
		"IONOS_API_URL_KAFKA":                     "Kafka",
		"IONOS_API_URL_LOGGING":                   "Logging",
		"IONOS_API_URL_MONITORING":                "Monitoring",
		"IONOS_API_URL_MARIADB":                   "MariaDB",
		"IONOS_API_URL_NFS":                       "NFS",
		"IONOS_API_URL_INMEMORYDB":                "InMemoryDB",
		"IONOS_API_URL_OBJECT_STORAGE":            "Object Storage",
		"IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT": "Object Storage Management",
	}

	var set []string
	for envVar, product := range envVars {
		if val := os.Getenv(envVar); val != "" {
			set = append(set, fmt.Sprintf("%s (%s): %s", envVar, product, val))
		}
	}
	sort.Strings(set)
	if len(set) > 0 {
		log.Printf("[DEBUG] Endpoint env vars: %s", strings.Join(set, " | "))
	}
}

// LogTLSConfig logs TLS-related configuration.
func LogTLSConfig(insecureBool bool) {
	var parts []string
	if os.Getenv("IONOS_ALLOW_INSECURE") != "" {
		parts = append(parts, "IONOS_ALLOW_INSECURE is set")
	}
	if insecureBool {
		parts = append(parts, "TLS verification disabled")
	}
	if len(parts) > 0 {
		log.Printf("[DEBUG] TLS: %s", strings.Join(parts, ", "))
	}
}

// LogEndpoint logs the global endpoint configuration.
func LogEndpoint(endpoint string) {
	if endpoint != "" {
		log.Printf("[DEBUG] Global endpoint: %s", endpoint)
	} else {
		log.Printf("[DEBUG] Global endpoint not set, using SDK defaults")
	}
}

// LogS3Region logs the S3 region configuration.
func LogS3Region(region string) {
	if region != "" {
		log.Printf("[DEBUG] S3 region: %s (source: explicit)", region)
	} else {
		log.Printf("[DEBUG] S3 region: eu-central-3 (source: default)")
	}
}

// FormatLocation returns a display string for a location, showing "(no location)" when empty.
func FormatLocation(location string) string {
	if location == "" {
		return "(no location)"
	}
	return location
}
