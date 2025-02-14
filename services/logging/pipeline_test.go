package logging

import (
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"os"
	"testing"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/stretchr/testify/assert"
)

func TestClientConfigurationFlowTable(t *testing.T) {
	tests := []struct {
		name           string
		clientOptions  bundle.ClientOptions
		loadedConfig   *shared.LoadedConfig
		productName    string
		location       string
		envVar         string
		expectedURL    string
		expectedEnvURL string
	}{
		{
			name: "overrideClientEndpoint",
			clientOptions: bundle.ClientOptions{
				ClientOverrideOptions: shared.ClientOverrideOptions{
					Endpoint:      "https://custom.endpoint.com",
					SkipTLSVerify: true,
					Certificate:   "",
					Credentials: shared.Credentials{
						Username: "test-user",
						Password: "test-password",
						Token:    "test-token",
					},
				},
				TerraformVersion: "1.0.0",
			},
			loadedConfig: &shared.LoadedConfig{
				Version:        1.0,
				CurrentProfile: "default",
				Profiles: []shared.Profile{
					{
						Environment: "de/fra",
						Name:        "default",
						Credentials: shared.Credentials{
							Username: "user123",
							Password: "pass123",
							Token:    "token123",
						},
					},
				},
				Environments: []shared.Environment{
					{
						Name:                "de/fra",
						CertificateAuthData: "cert_data_here",
						Products: []shared.Product{
							{
								Name: "logging",
								Endpoints: []shared.Endpoint{
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
			productName:    shared.Logging,
			location:       "de/fra",
			envVar:         "https://env.endpoint.com",
			expectedURL:    "https://override.logging.de-fra.ionos.com",
			expectedEnvURL: "https://env.endpoint.com",
		},
		{
			name: "ProductNotDefinedInEnv",
			clientOptions: bundle.ClientOptions{
				ClientOverrideOptions: shared.ClientOverrideOptions{
					Endpoint:      "https://custom.endpoint.com",
					SkipTLSVerify: true,
					Certificate:   "",
					Credentials: shared.Credentials{
						Username: "test-user",
						Password: "test-password",
						Token:    "test-token",
					},
				},
				TerraformVersion: "1.0.0",
			},
			loadedConfig: &shared.LoadedConfig{
				Version:        1.0,
				CurrentProfile: "default",
				Profiles: []shared.Profile{
					{
						Environment: "de/fra",
						Name:        "default",
						Credentials: shared.Credentials{
							Username: "user123",
							Password: "pass123",
							Token:    "token123",
						},
					},
				},
				Environments: []shared.Environment{
					{
						Name:                "de/fra",
						CertificateAuthData: "cert_data_here",
						Products:            []shared.Product{},
					},
				},
			},
			productName:    shared.Logging,
			location:       "de/fra",
			envVar:         "https://env.endpoint.com",
			expectedURL:    "https://custom.endpoint.com",
			expectedEnvURL: "https://env.endpoint.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.clientOptions, tt.loadedConfig)

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
