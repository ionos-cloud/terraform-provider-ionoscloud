package logging

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
)

func TestClientConfigurationFlowTable(t *testing.T) {
	tests := []struct {
		name            string
		clientOptions   bundle.ClientOptions
		fileConfig      *fileconfiguration.FileConfig
		productName     string
		location        string
		envVar          string
		expectedURL     string
		expectedEnvURL  string
		expectedCert    string
		expectedSkipTLS bool
	}{
		{
			name: "overrideClientEndpoint",
			clientOptions: bundle.ClientOptions{
				ClientOptions: shared.ClientOptions{
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
			fileConfig: &fileconfiguration.FileConfig{
				Version:        1.0,
				CurrentProfile: "default",
				Profiles: []fileconfiguration.Profile{
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
				Environments: []fileconfiguration.Environment{
					{
						Name: "de/fra",
						Products: []fileconfiguration.Product{
							{
								Name: "logging",
								Endpoints: []fileconfiguration.Endpoint{
									{
										Name:                "https://override.logging.de-fra.ionos.com",
										Location:            "de/fra",
										SkipTLSVerify:       false,
										CertificateAuthData: "cert_data_here",
									},
								},
							},
						},
					},
				},
			},
			productName:     fileconfiguration.Logging,
			location:        "de/fra",
			envVar:          "https://env.endpoint.com",
			expectedURL:     "https://override.logging.de-fra.ionos.com",
			expectedEnvURL:  "https://env.endpoint.com",
			expectedCert:    "",
			expectedSkipTLS: false,
		},
		{
			name: "ProductNotDefinedInEnvUseGlobal",
			clientOptions: bundle.ClientOptions{
				ClientOptions: shared.ClientOptions{
					Endpoint:      "https://custom.endpoint.com",
					SkipTLSVerify: true,
					Certificate:   "cert_data_here",
					Credentials: shared.Credentials{
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
						Credentials: shared.Credentials{
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
			productName:     fileconfiguration.Logging,
			location:        "de/fra",
			envVar:          "https://env.endpoint.com",
			expectedURL:     "https://custom.endpoint.com",
			expectedEnvURL:  "https://env.endpoint.com",
			expectedCert:    "cert_data_here",
			expectedSkipTLS: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.clientOptions, tt.fileConfig)

			// Step 2: Override the client configuration from the loaded config
			loadedconfig.SetClientOptionsFromConfig(client, tt.productName, tt.location)

			// Verify the override worked
			assert.Equal(t, tt.expectedURL, client.GetConfig().Servers[0].URL)
			//assert.Equal(t, tt.expectedCert, client.GetConfig().Certificate)
			assert.Equal(t, tt.expectedSkipTLS, client.GetConfig().HTTPClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify)
			expectedRootCa := shared.AddCertsToClient(tt.expectedCert)
			configRootCA := client.GetConfig().HTTPClient.Transport.(*http.Transport).TLSClientConfig.RootCAs
			assert.True(t, expectedRootCa.Equal(configRootCA))

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
