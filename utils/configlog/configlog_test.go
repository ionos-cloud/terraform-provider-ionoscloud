package configlog

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// captureLog captures log output during fn execution and returns it as a string.
func captureLog(fn func()) string {
	var buf bytes.Buffer
	origOutput := log.Writer()
	origFlags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0) // no timestamps for easier assertion
	defer func() {
		log.SetOutput(origOutput)
		log.SetFlags(origFlags)
	}()
	fn()
	return buf.String()
}

func TestLogProfileAndEnvironment_WithProfiles(t *testing.T) {
	cfg := &fileconfiguration.FileConfig{
		CurrentProfile: "prod",
		Profiles: []fileconfiguration.Profile{
			{Name: "prod", Environment: "production"},
			{Name: "dev", Environment: "development"},
		},
	}

	output := captureLog(func() {
		logProfileAndEnvironment(cfg)
	})

	assertContains(t, output, "Profile resolution:")
	assertContains(t, output, "2 profile(s)")
	assertContains(t, output, `currentProfile: "prod"`)
	assertContains(t, output, `active: "prod" (environment: "production")`)
}

func TestLogProfileAndEnvironment_NoMatchingProfile(t *testing.T) {
	cfg := &fileconfiguration.FileConfig{
		CurrentProfile: "staging",
		Profiles: []fileconfiguration.Profile{
			{Name: "prod"},
			{Name: "dev"},
		},
	}

	output := captureLog(func() {
		logProfileAndEnvironment(cfg)
	})

	assertContains(t, output, `no matching profile for "staging"`)
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

	output := captureLog(func() {
		logProfileAndEnvironment(cfg)
	})

	assertContains(t, output, `overrides to "override-profile"`)
}

func TestLogCredentialResolution_TokenOnly(t *testing.T) {
	output := captureLog(func() {
		LogCredentialResolution(shared.Credentials{Token: "my-token"}, false, "")
	})

	assertContains(t, output, "token=found")
	assertContains(t, output, "user/pass=not found")
	assertContains(t, output, "S3 keys=not found")
	assertContains(t, output, "authenticating via token")
	assertNotContains(t, output, "my-token")
}

func TestLogCredentialResolution_UsernamePassword(t *testing.T) {
	output := captureLog(func() {
		LogCredentialResolution(shared.Credentials{Username: "user", Password: "pass"}, false, "")
	})

	assertContains(t, output, "token=not found")
	assertContains(t, output, "user/pass=found")
	assertContains(t, output, "authenticating via user/pass")
}

func TestLogCredentialResolution_BothTokenAndUserPass(t *testing.T) {
	output := captureLog(func() {
		LogCredentialResolution(shared.Credentials{Token: "tok", Username: "user", Password: "pass"}, false, "")
	})

	assertContains(t, output, "both token and user/pass provided; token takes precedence")
	assertContains(t, output, "authenticating via token")
}

func TestLogCredentialResolution_FileConfigProfile(t *testing.T) {
	output := captureLog(func() {
		LogCredentialResolution(shared.Credentials{Token: "tok", S3AccessKey: "ak", S3SecretKey: "sk"}, true, "myprofile")
	})

	assertContains(t, output, `file config profile "myprofile": token, S3 keys`)
}

func TestLogCredentialResolution_FileConfigNoCredentials(t *testing.T) {
	output := captureLog(func() {
		LogCredentialResolution(shared.Credentials{}, true, "empty")
	})

	assertContains(t, output, `file config profile "empty": none`)
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

	output := captureLog(func() {
		LogEndpointEnvVars()
	})

	if strings.Contains(output, "Endpoint env var") {
		t.Errorf("expected no endpoint env var output, got: %s", output)
	}
}

func TestLogEndpointEnvVars_SomeSet(t *testing.T) {
	t.Setenv("IONOS_API_URL_VPN", "https://vpn.custom.example.com")
	t.Setenv("IONOS_API_URL_KAFKA", "https://kafka.custom.example.com")

	output := captureLog(func() {
		LogEndpointEnvVars()
	})

	assertContains(t, output, "Endpoint env vars:")
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

	output := captureLog(func() {
		logFileConfigEndpoints(cfg)
	})

	assertContains(t, output, `Environment "production": 2 product(s):`)
	assertContains(t, output, "cloud(1)")
	assertContains(t, output, "dns(2)")
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

	output := captureLog(func() {
		logFileConfigEndpoints(cfg)
	})

	assertNotContains(t, output, "product(s)")
}

func TestLogTLSConfig_InsecureSet(t *testing.T) {
	t.Setenv("IONOS_ALLOW_INSECURE", "true")

	output := captureLog(func() {
		LogTLSConfig(true)
	})

	assertContains(t, output, "TLS:")
	assertContains(t, output, "IONOS_ALLOW_INSECURE is set")
	assertContains(t, output, "TLS verification disabled")
}

func TestLogTLSConfig_Secure(t *testing.T) {
	t.Setenv("IONOS_ALLOW_INSECURE", "")

	output := captureLog(func() {
		LogTLSConfig(false)
	})

	assertNotContains(t, output, "TLS")
}

func TestLogEndpoint_Set(t *testing.T) {
	output := captureLog(func() {
		LogEndpoint("https://api.ionos.com")
	})
	assertContains(t, output, "Global endpoint: https://api.ionos.com")
}

func TestLogEndpoint_Empty(t *testing.T) {
	output := captureLog(func() {
		LogEndpoint("")
	})
	assertContains(t, output, "Global endpoint not set, using SDK defaults")
}

func TestLogS3Region_Explicit(t *testing.T) {
	output := captureLog(func() {
		LogS3Region("us-central-1")
	})
	assertContains(t, output, "S3 region: us-central-1")
}

func TestLogS3Region_Default(t *testing.T) {
	output := captureLog(func() {
		LogS3Region("")
	})
	assertContains(t, output, "S3 region: "+constant.DefaultS3Region+" (default)")
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

	output := captureLog(func() {
		cfg, loadErr := LoadFileConfigWithLogging()
		if loadErr == nil {
			t.Error("expected error for unreadable file")
		}
		if cfg != nil {
			t.Error("expected nil config for unreadable file")
		}
	})

	assertContains(t, output, "Config file:")
	assertContains(t, output, "unreadable")
	assertContains(t, output, "0000")
}

func TestLoadFileConfigWithLogging_NoFile(t *testing.T) {
	t.Setenv(shared.IonosFilePathEnvVar, "/tmp/nonexistent-ionos-config-test-file")

	output := captureLog(func() {
		cfg, err := LoadFileConfigWithLogging()
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
		if cfg != nil {
			t.Error("expected nil config for nonexistent file")
		}
	})

	assertContains(t, output, "Config file: /tmp/nonexistent-ionos-config-test-file")
	assertContains(t, output, "source: IONOS_CONFIG_FILE")
	assertContains(t, output, "status: not found")
	assertContains(t, output, "Config file not loaded")
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

	output := captureLog(func() {
		cfg, loadErr := LoadFileConfigWithLogging()
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

	assertContains(t, output, "Config file:")
	assertContains(t, output, "source: IONOS_CONFIG_FILE")
	assertContains(t, output, "status: found")
	assertContains(t, output, "Config file loaded successfully (version: 1.0)")
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
