package configlog

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tfsdklog"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// captureLog runs fn with a context that has a tflog provider logger attached
// (writing to a temp file) and returns the captured output.
func captureLog(t *testing.T, fn func(ctx context.Context)) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tflog-*.log")
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	t.Setenv("TF_LOG", "DEBUG")
	t.Setenv("TF_LOG_PATH", f.Name())

	ctx := tfsdklog.RegisterTestSink(context.Background(), t)
	ctx = tfsdklog.NewRootProviderLogger(ctx)

	fn(ctx)

	data, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestLogProfileAndEnvironment_WithProfiles(t *testing.T) {
	cfg := &fileconfiguration.FileConfig{
		CurrentProfile: "prod",
		Profiles: []fileconfiguration.Profile{
			{Name: "prod", Environment: "production"},
			{Name: "dev", Environment: "development"},
		},
	}

	output := captureLog(t, func(ctx context.Context) {
		logProfileAndEnvironment(ctx, cfg)
	})

	assertContains(t, output, "profile resolution")
	assertContains(t, output, "2 profile(s)")
	assertContains(t, output, "currentProfile:")
	assertContains(t, output, "prod")
	assertContains(t, output, "production")
}

func TestLogProfileAndEnvironment_NoMatchingProfile(t *testing.T) {
	cfg := &fileconfiguration.FileConfig{
		CurrentProfile: "staging",
		Profiles: []fileconfiguration.Profile{
			{Name: "prod"},
			{Name: "dev"},
		},
	}

	output := captureLog(t, func(ctx context.Context) {
		logProfileAndEnvironment(ctx, cfg)
	})

	assertContains(t, output, "no matching profile for")
	assertContains(t, output, "staging")
	assertContains(t, output, "prod, dev")
}

func TestLogProfileAndEnvironment_EnvOverride(t *testing.T) {
	t.Setenv(shared.IonosCurrentProfileEnvVar, "override-profile")

	cfg := &fileconfiguration.FileConfig{
		CurrentProfile: "default",
		Profiles: []fileconfiguration.Profile{
			{Name: "default"},
		},
	}

	output := captureLog(t, func(ctx context.Context) {
		logProfileAndEnvironment(ctx, cfg)
	})

	assertContains(t, output, "overrides to")
	assertContains(t, output, "override-profile")
}

func TestLogCredentialResolution_TokenOnly(t *testing.T) {
	output := captureLog(t, func(ctx context.Context) {
		LogCredentialResolution(ctx, shared.Credentials{Token: "my-token"}, false, "")
	})

	assertContains(t, output, "token=found")
	assertContains(t, output, "user/pass=not found")
	assertContains(t, output, "S3 keys=not found")
	assertContains(t, output, "authenticating via token")
	assertNotContains(t, output, "my-token")
}

func TestLogCredentialResolution_UsernamePassword(t *testing.T) {
	output := captureLog(t, func(ctx context.Context) {
		LogCredentialResolution(ctx, shared.Credentials{Username: "user", Password: "pass"}, false, "")
	})

	assertContains(t, output, "token=not found")
	assertContains(t, output, "user/pass=found")
	assertContains(t, output, "authenticating via user/pass")
}

func TestLogCredentialResolution_BothTokenAndUserPass(t *testing.T) {
	output := captureLog(t, func(ctx context.Context) {
		LogCredentialResolution(ctx, shared.Credentials{Token: "tok", Username: "user", Password: "pass"}, false, "")
	})

	assertContains(t, output, "both token and user/pass provided; token takes precedence")
	assertContains(t, output, "authenticating via token")
}

func TestLogCredentialResolution_FileConfigProfile(t *testing.T) {
	output := captureLog(t, func(ctx context.Context) {
		LogCredentialResolution(ctx, shared.Credentials{Token: "tok", S3AccessKey: "ak", S3SecretKey: "sk"}, true, "myprofile")
	})

	assertContains(t, output, "file config profile")
	assertContains(t, output, "myprofile")
	assertContains(t, output, "token, S3 keys")
}

func TestLogCredentialResolution_FileConfigNoCredentials(t *testing.T) {
	output := captureLog(t, func(ctx context.Context) {
		LogCredentialResolution(ctx, shared.Credentials{}, true, "empty")
	})

	assertContains(t, output, "file config profile")
	assertContains(t, output, "empty")
	assertContains(t, output, "none")
}

func TestLogEndpointEnvVars_NoneSet(t *testing.T) {
	for _, env := range []string{
		"IONOS_API_URL", "IONOS_API_URL_VPN", "IONOS_API_URL_KAFKA",
		"IONOS_API_URL_LOGGING", "IONOS_API_URL_MONITORING", "IONOS_API_URL_MARIADB",
		"IONOS_API_URL_NFS", "IONOS_API_URL_INMEMORYDB", "IONOS_API_URL_OBJECT_STORAGE",
		"IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT",
	} {
		t.Setenv(env, "")
	}

	output := captureLog(t, func(ctx context.Context) {
		LogEndpointEnvVars(ctx)
	})

	if strings.Contains(output, "endpoint env vars") {
		t.Errorf("expected no endpoint env var output, got: %s", output)
	}
}

func TestLogEndpointEnvVars_SomeSet(t *testing.T) {
	t.Setenv("IONOS_API_URL_VPN", "https://vpn.custom.example.com")
	t.Setenv("IONOS_API_URL_KAFKA", "https://kafka.custom.example.com")

	output := captureLog(t, func(ctx context.Context) {
		LogEndpointEnvVars(ctx)
	})

	assertContains(t, output, "endpoint env vars")
	assertContains(t, output, "IONOS_API_URL_VPN (VPN): https://vpn.custom.example.com")
	assertContains(t, output, "IONOS_API_URL_KAFKA (Kafka): https://kafka.custom.example.com")
}

func TestLogFileConfigEndpoints_WithEnvironment(t *testing.T) {
	cfg := &fileconfiguration.FileConfig{
		CurrentProfile: "prod",
		Profiles: []fileconfiguration.Profile{
			{Name: "prod", Environment: "production"},
		},
		Environments: []fileconfiguration.Environment{
			{
				Name: "production",
				Products: []fileconfiguration.Product{
					{Name: "cloud", Endpoints: []fileconfiguration.Endpoint{{Name: "https://api.example.com"}}},
					{Name: "dns", Endpoints: []fileconfiguration.Endpoint{{Name: "https://dns1.example.com"}, {Name: "https://dns2.example.com"}}},
				},
			},
		},
	}

	output := captureLog(t, func(ctx context.Context) {
		logFileConfigEndpoints(ctx, cfg)
	})

	assertContains(t, output, "environment products")
	assertContains(t, output, "production")
	assertContains(t, output, "cloud")
	assertContains(t, output, "https://api.example.com")
	assertContains(t, output, "dns")
	assertContains(t, output, "https://dns1.example.com")
	assertContains(t, output, "https://dns2.example.com")
}

func TestLogFileConfigEndpoints_WithTLSAndCertPerEndpoint(t *testing.T) {
	cfg := &fileconfiguration.FileConfig{
		CurrentProfile: "prod",
		Profiles: []fileconfiguration.Profile{
			{Name: "prod", Environment: "staging"},
		},
		Environments: []fileconfiguration.Environment{
			{
				Name: "staging",
				Products: []fileconfiguration.Product{
					{Name: "cloud", Endpoints: []fileconfiguration.Endpoint{
						{Name: "https://api.staging.example.com", SkipTLSVerify: true, CertificateAuthData: "some-cert-data-here"},
					}},
					{Name: "dns", Endpoints: []fileconfiguration.Endpoint{
						{Name: "https://dns.staging.example.com", Location: "de/fra"},
					}},
				},
			},
		},
	}

	output := captureLog(t, func(ctx context.Context) {
		logFileConfigEndpoints(ctx, cfg)
	})

	assertContains(t, output, "skipTlsVerify")
	assertContains(t, output, "certAuthDataBytes")
	assertContains(t, output, "19")
	assertContains(t, output, "de/fra")
	assertContains(t, output, "https://api.staging.example.com")
	assertContains(t, output, "https://dns.staging.example.com")
}

func TestLogFileConfigEndpoints_NoMatchingEnvironment(t *testing.T) {
	cfg := &fileconfiguration.FileConfig{
		CurrentProfile: "prod",
		Profiles: []fileconfiguration.Profile{
			{Name: "prod", Environment: "staging"},
		},
		Environments: []fileconfiguration.Environment{
			{Name: "production"},
		},
	}

	output := captureLog(t, func(ctx context.Context) {
		logFileConfigEndpoints(ctx, cfg)
	})

	assertNotContains(t, output, "product_count")
}

func TestLogTLSConfig_InsecureSet(t *testing.T) {
	t.Setenv("IONOS_ALLOW_INSECURE", "true")

	output := captureLog(t, func(ctx context.Context) {
		LogTLSConfig(ctx, true)
	})

	assertContains(t, output, "TLS config")
	assertContains(t, output, "IONOS_ALLOW_INSECURE is set")
	assertContains(t, output, "TLS verification disabled")
}

func TestLogTLSConfig_Secure(t *testing.T) {
	t.Setenv("IONOS_ALLOW_INSECURE", "")

	output := captureLog(t, func(ctx context.Context) {
		LogTLSConfig(ctx, false)
	})

	assertNotContains(t, output, "TLS config")
}

func TestLogTLSConfig_WithPinnedCert(t *testing.T) {
	t.Setenv("IONOS_ALLOW_INSECURE", "")
	t.Setenv("IONOS_PINNED_CERT", "sha256-fingerprint-here")

	output := captureLog(t, func(ctx context.Context) {
		LogTLSConfig(ctx, false)
	})

	assertContains(t, output, "IONOS_PINNED_CERT is set")
	assertContains(t, output, "23 bytes")
}

func TestLogEndpoint_Set(t *testing.T) {
	output := captureLog(t, func(ctx context.Context) {
		LogEndpoint(ctx, "https://api.ionos.com")
	})
	assertContains(t, output, "global endpoint")
	assertContains(t, output, "https://api.ionos.com")
}

func TestLogEndpoint_Empty(t *testing.T) {
	output := captureLog(t, func(ctx context.Context) {
		LogEndpoint(ctx, "")
	})
	assertContains(t, output, "global endpoint not set, using SDK defaults")
}

func TestLogS3Region_Explicit(t *testing.T) {
	output := captureLog(t, func(ctx context.Context) {
		LogS3Region(ctx, "us-central-1")
	})
	assertContains(t, output, "S3 region")
	assertContains(t, output, "us-central-1")
}

func TestLogS3Region_Default(t *testing.T) {
	output := captureLog(t, func(ctx context.Context) {
		LogS3Region(ctx, "")
	})
	assertContains(t, output, "S3 region")
	assertContains(t, output, constant.DefaultS3Region)
}

func TestFormatLocation(t *testing.T) {
	if FormatLocation("") != "(no location)" {
		t.Error("expected (no location) for empty string")
	}
	if FormatLocation("de/fra") != "de/fra" {
		t.Error("expected de/fra")
	}
}

func TestLoadFileConfigWithLogging_UnreadableFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "ionos-config-noperm-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString("version: 1.0\n"); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	if err := os.Chmod(tmpFile.Name(), 0000); err != nil {
		t.Fatal(err)
	}

	t.Setenv(shared.IonosFilePathEnvVar, tmpFile.Name())

	output := captureLog(t, func(ctx context.Context) {
		cfg, loadErr := LoadFileConfigWithLogging(ctx)
		if loadErr == nil {
			t.Error("expected error for unreadable file")
		}
		if cfg != nil {
			t.Error("expected nil config for unreadable file")
		}
	})

	assertContains(t, output, "config file")
	assertContains(t, output, "unreadable")
	assertContains(t, output, "0000")
}

func TestLoadFileConfigWithLogging_NoFile(t *testing.T) {
	t.Setenv(shared.IonosFilePathEnvVar, "/tmp/nonexistent-ionos-config-test-file")

	output := captureLog(t, func(ctx context.Context) {
		cfg, err := LoadFileConfigWithLogging(ctx)
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
		if cfg != nil {
			t.Error("expected nil config for nonexistent file")
		}
	})

	assertContains(t, output, "/tmp/nonexistent-ionos-config-test-file")
	assertContains(t, output, "IONOS_CONFIG_FILE")
	assertContains(t, output, "not found")
	assertContains(t, output, "config file not loaded")
}

func TestLoadFileConfigWithLogging_ValidFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "ionos-config-test-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := `version: 1.0
currentProfile: test
profiles:
  - name: test
    environment: dev
    credentials:
      token: "test-token"
environments:
  - name: dev
    products: []
`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	t.Setenv(shared.IonosFilePathEnvVar, tmpFile.Name())

	output := captureLog(t, func(ctx context.Context) {
		cfg, loadErr := LoadFileConfigWithLogging(ctx)
		if loadErr != nil {
			t.Errorf("unexpected error: %s", loadErr)
		}
		if cfg == nil {
			t.Fatal("expected non-nil config")
		}
		if cfg.CurrentProfile != "test" {
			t.Errorf("expected currentProfile 'test', got %q", cfg.CurrentProfile)
		}
	})

	assertContains(t, output, "config file")
	assertContains(t, output, "IONOS_CONFIG_FILE")
	assertContains(t, output, "found")
	assertContains(t, output, "config file loaded successfully")
}

func assertContains(t *testing.T, output, expected string) {
	t.Helper()
	if !strings.Contains(output, expected) {
		t.Errorf("expected output to contain %q, got:\n%s", expected, output)
	}
}

func assertNotContains(t *testing.T, output, unexpected string) {
	t.Helper()
	if strings.Contains(output, unexpected) {
		t.Errorf("expected output NOT to contain %q, got:\n%s", unexpected, output)
	}
}
