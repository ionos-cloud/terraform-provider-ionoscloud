package configlog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// LoadFileConfigWithLogging wraps fileconfiguration.NewFromEnv() with pre/post logging
// so that users can trace config file loading with TF_LOG=DEBUG.
func LoadFileConfigWithLogging(ctx context.Context) (*fileconfiguration.FileConfig, error) {
	// Resolve path for logging
	filePath := os.Getenv(shared.IonosFilePathEnvVar)
	source := shared.IonosFilePathEnvVar
	if filePath == "" {
		source = "default"
		defaultPath, err := fileconfiguration.DefaultConfigFileName()
		if err != nil {
			tflog.Debug(ctx, "could not determine default config file path", map[string]interface{}{"error": err.Error()})
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
				_ = f.Close()
			}
		} else if os.IsNotExist(err) {
			status = "not found"
		}
		tflog.Debug(ctx, "config file", map[string]interface{}{"path": filePath, "source": source, "status": status})
	}

	fileConfig, err := fileconfiguration.NewFromEnv()
	if err != nil {
		tflog.Debug(ctx, "config file not loaded", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	tflog.Debug(ctx, "config file loaded successfully", map[string]interface{}{"version": float64(fileConfig.Version)})
	logProfileAndEnvironment(ctx, fileConfig)
	logFileConfigEndpoints(ctx, fileConfig)

	return fileConfig, nil
}

// logProfileAndEnvironment logs profile and environment resolution details.
func logProfileAndEnvironment(ctx context.Context, fileConfig *fileconfiguration.FileConfig) {
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

	tflog.Debug(ctx, "profile resolution", map[string]interface{}{"details": strings.Join(parts, " | ")})
}

// logFileConfigEndpoints logs product and endpoint counts from the active environment in the file config.
func logFileConfigEndpoints(ctx context.Context, fileConfig *fileconfiguration.FileConfig) {
	failoverOpts := fileConfig.GetFailoverOptions()
	if failoverOpts != nil {
		tflog.Debug(ctx, "failover config", map[string]interface{}{"strategy": failoverOpts.Strategy})
	} else {
		tflog.Debug(ctx, "failover config not set (default: none)")
	}

	envName := fileConfig.GetEnvForCurrentProfile()
	if envName == "" {
		return
	}

	for _, env := range fileConfig.Environments {
		if env.Name == envName {
			type endpointJSON struct {
				URL           string `json:"url"`
				Location      string `json:"location,omitempty"`
				SkipTLSVerify bool   `json:"skipTlsVerify,omitempty"`
				CertAuthBytes int    `json:"certAuthDataBytes,omitempty"`
			}
			type productJSON struct {
				Name      string         `json:"name"`
				Endpoints []endpointJSON `json:"endpoints"`
			}

			products := make([]productJSON, 0, len(env.Products))
			for _, product := range env.Products {
				p := productJSON{Name: product.Name}
				for _, ep := range product.Endpoints {
					e := endpointJSON{URL: ep.Name, Location: ep.Location, SkipTLSVerify: ep.SkipTLSVerify}
					if ep.CertificateAuthData != "" {
						e.CertAuthBytes = len(ep.CertificateAuthData)
					}
					p.Endpoints = append(p.Endpoints, e)
				}
				products = append(products, p)
			}

			jsonBytes, err := json.Marshal(products)
			if err != nil {
				tflog.Debug(ctx, "environment products (failed to marshal)", map[string]interface{}{"environment": env.Name, "product_count": len(env.Products), "error": err.Error()})
			} else {
				tflog.Debug(ctx, "environment products", map[string]interface{}{"environment": env.Name, "product_count": len(env.Products), "products": string(jsonBytes)})
			}
			return
		}
	}
}

// LogCredentialResolution logs which credentials were found and from where.
func LogCredentialResolution(ctx context.Context, creds shared.Credentials, fileConfigUsed bool, profileName string) {
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

	tflog.Debug(ctx, line)
}

// LogEndpointEnvVars logs which product-specific endpoint env vars are set.
func LogEndpointEnvVars(ctx context.Context) {
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
		tflog.Debug(ctx, "endpoint env vars", map[string]interface{}{"vars": strings.Join(set, " | ")})
	}
}

// LogTLSConfig logs TLS-related configuration.
func LogTLSConfig(ctx context.Context, insecureBool bool) {
	var parts []string
	if os.Getenv("IONOS_ALLOW_INSECURE") != "" {
		parts = append(parts, "IONOS_ALLOW_INSECURE is set")
	}
	if insecureBool {
		parts = append(parts, "TLS verification disabled")
	}
	if pinnedCert := os.Getenv(shared.IonosPinnedCertEnvVar); pinnedCert != "" {
		parts = append(parts, fmt.Sprintf("%s is set (%d bytes) — cert pinning active for all products", shared.IonosPinnedCertEnvVar, len(pinnedCert)))
	}
	if len(parts) > 0 {
		tflog.Debug(ctx, "TLS config", map[string]interface{}{"details": strings.Join(parts, ", ")})
	}
}

// LogEndpoint logs the global endpoint configuration.
func LogEndpoint(ctx context.Context, endpoint string) {
	if endpoint != "" {
		tflog.Debug(ctx, "global endpoint", map[string]interface{}{"endpoint": endpoint})
	} else {
		tflog.Debug(ctx, "global endpoint not set, using SDK defaults")
	}
}

// LogS3Region logs the S3 region configuration.
func LogS3Region(ctx context.Context, region string) {
	if region != "" {
		tflog.Debug(ctx, "S3 region", map[string]interface{}{"region": region})
	} else {
		tflog.Debug(ctx, "S3 region (default)", map[string]interface{}{"region": constant.DefaultS3Region})
	}
}

// FormatLocation returns a display string for a location, showing "(no location)" when empty.
func FormatLocation(location string) string {
	if location == "" {
		return "(no location)"
	}
	return location
}
