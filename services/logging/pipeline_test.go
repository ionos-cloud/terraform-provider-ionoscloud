package logging

import (
	"os"
	"testing"

	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/stretchr/testify/assert"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
)

func TestClientConfigurationFlowTable(t *testing.T) {
	tests := []struct {
		name           string
		clientOptions  bundle.ClientOptions
		fileConfig     *fileconfiguration.FileConfig
		productName    string
		location       string
		envVar         string
		expectedURL    string
		expectedEnvURL string
	}{
		{
			name: "overrideClientEndpoint",
			clientOptions: bundle.ClientOptions{
				ClientOverrideOptions: fileconfiguration.ClientOverrideOptions{
					Endpoint:      "https://custom.endpoint.com",
					SkipTLSVerify: true,
					Certificate:   "",
					Credentials: fileconfiguration.Credentials{
						Username: "test-user",
						Password: "test-password",
						Token:    "test-token",
					},
				},
				TerraformVersion: "1.0.0",
			},
			fileConfig: &fileconfiguration.FileConfig{
				Version:        1.0,
				CurrentProfile: "default",
				Profiles: []fileconfiguration.Profile{
					{
						Environment: "de/fra",
						Name:        "default",
						Credentials: fileconfiguration.Credentials{
							Username: "user123",
							Password: "pass123",
							Token:    "token123",
						},
					},
				},
				Environments: []fileconfiguration.Environment{
					{
						Name:                "de/fra",
						CertificateAuthData: "cert_data_here",
						Products: []fileconfiguration.Product{
							{
								Name: "logging",
								Endpoints: []fileconfiguration.Endpoint{
									{
										Name:          "https://override.logging.de-fra.ionos.com",
										Location:      "de/fra",
										SkipTLSVerify: false,
									},
								},
							},
						},
					},
				},
			},
			productName:    fileconfiguration.Logging,
			location:       "de/fra",
			envVar:         "https://env.endpoint.com",
			expectedURL:    "https://override.logging.de-fra.ionos.com",
			expectedEnvURL: "https://env.endpoint.com",
		},
		{
			name: "ProductNotDefinedInEnv",
			clientOptions: bundle.ClientOptions{
				ClientOverrideOptions: fileconfiguration.ClientOverrideOptions{
					Endpoint:      "https://custom.endpoint.com",
					SkipTLSVerify: true,
					Certificate:   "",
					Credentials: fileconfiguration.Credentials{
						Username: "test-user",
						Password: "test-password",
						Token:    "test-token",
					},
				},
				TerraformVersion: "1.0.0",
			},
			fileConfig: &fileconfiguration.FileConfig{
				Version:        1.0,
				CurrentProfile: "default",
				Profiles: []fileconfiguration.Profile{
					{
						Environment: "de/fra",
						Name:        "default",
						Credentials: fileconfiguration.Credentials{
							Username: "user123",
							Password: "pass123",
							Token:    "token123",
						},
					},
				},
				Environments: []fileconfiguration.Environment{
					{
						Name:                "de/fra",
						CertificateAuthData: "cert_data_here",
						Products:            []fileconfiguration.Product{},
					},
				},
			},
			productName:    fileconfiguration.Logging,
			location:       "de/fra",
			envVar:         "https://env.endpoint.com",
			expectedURL:    "https://custom.endpoint.com",
			expectedEnvURL: "https://env.endpoint.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.clientOptions, tt.fileConfig)

			// Step 2: Override the client configuration from the loaded config
			loadedconfig.SetClientOptionsFromConfig(client, tt.productName, tt.location)

			// Verify the override worked
			assert.Equal(t, tt.expectedURL, client.GetConfig().Servers[0].URL)

			// Step 3: Change the config URL using ChangeConfigURL
			client.ChangeConfigURL(tt.location)

			// Verify the final URL after change
			assert.Equal(t, tt.expectedURL, client.GetConfig().Servers[0].URL)

			// Test changing the config URL with environment variable
			os.Setenv("IONOS_API_URL_LOGGING", tt.envVar)
			client.ChangeConfigURL("")
			assert.Equal(t, tt.expectedEnvURL, client.GetConfig().Servers[0].URL)
			os.Unsetenv("IONOS_API_URL_LOGGING")
		})
	}
}
