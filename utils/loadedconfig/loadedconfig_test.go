package loadedconfig

import (
	"net/http"
	"testing"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
)

func TestSetClientOptionsFromfileConfig(t *testing.T) {
	type args struct {
		clientOptions *clientoptions.TerraformClientOptions
		fileConfig    *fileconfiguration.FileConfig
		productName   string
	}
	tests := []struct {
		name              string
		args              args
		wantClientOptions *clientoptions.TerraformClientOptions
	}{
		{
			name: "NilClientOptions",
			args: args{
				clientOptions: nil,
				fileConfig:    &fileconfiguration.FileConfig{},
				productName:   "testProduct",
			},
			wantClientOptions: nil,
		},
		{
			name: "NilfileConfig",
			args: args{
				clientOptions: &clientoptions.TerraformClientOptions{},
				fileConfig:    nil,
				productName:   "testProduct",
			},
			wantClientOptions: &clientoptions.TerraformClientOptions{},
		},
		{
			name: "MultipleEndpoints",
			args: args{
				clientOptions: &clientoptions.TerraformClientOptions{},
				fileConfig: &fileconfiguration.FileConfig{
					Environments: []fileconfiguration.Environment{
						{
							Name: "testEnv",
							Products: []fileconfiguration.Product{
								{
									Name: "testProduct",
									Endpoints: []fileconfiguration.Endpoint{
										{Name: "endpoint1", SkipTLSVerify: true},
										{Name: "endpoint2", SkipTLSVerify: false},
									},
								},
							},
						},
					},
				},
				productName: "testProduct",
			},
			wantClientOptions: &clientoptions.TerraformClientOptions{
				ClientOptions: shared.ClientOptions{
					Endpoint:      "endpoint1",
					SkipTLSVerify: true,
				},
			},
		},
		{
			name: "SingleEndpoint",
			args: args{
				clientOptions: &clientoptions.TerraformClientOptions{},
				fileConfig: &fileconfiguration.FileConfig{
					Environments: []fileconfiguration.Environment{
						{
							Name: "testEnv",
							Products: []fileconfiguration.Product{
								{
									Name: "testProduct",
									Endpoints: []fileconfiguration.Endpoint{
										{Name: "endpoint1", SkipTLSVerify: true},
									},
								},
							},
						},
					},
				},
				productName: "testProduct",
			},
			wantClientOptions: &clientoptions.TerraformClientOptions{
				ClientOptions: shared.ClientOptions{
					Endpoint:      "endpoint1",
					SkipTLSVerify: true,
				},
			},
		},
		{
			name: "NoEndpoints",
			args: args{
				clientOptions: &clientoptions.TerraformClientOptions{},
				fileConfig: &fileconfiguration.FileConfig{
					Environments: []fileconfiguration.Environment{
						{
							Name: "testEnv",
							Products: []fileconfiguration.Product{
								{
									Name:      "testProduct",
									Endpoints: []fileconfiguration.Endpoint{},
								},
							},
						},
					},
				},
				productName: "testProduct",
			},
			wantClientOptions: &clientoptions.TerraformClientOptions{},
		},
		{
			name: "BadProductName",
			args: args{
				clientOptions: &clientoptions.TerraformClientOptions{},
				fileConfig: &fileconfiguration.FileConfig{
					Environments: []fileconfiguration.Environment{
						{
							Name: "testEnv",
							Products: []fileconfiguration.Product{
								{
									Name:      "testProduct",
									Endpoints: []fileconfiguration.Endpoint{},
								},
							},
						},
					},
				},
				productName: "productDoesNotExist",
			},
			wantClientOptions: &clientoptions.TerraformClientOptions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetGlobalClientOptionsFromFileConfig(tt.args.clientOptions, tt.args.fileConfig, tt.args.productName)
			if tt.args.clientOptions != nil && tt.wantClientOptions != nil {
				if tt.args.clientOptions.Endpoint != tt.wantClientOptions.Endpoint {
					t.Errorf("got %v, want %v", tt.args.clientOptions.Endpoint, tt.wantClientOptions.Endpoint)
				}
				if tt.args.clientOptions.SkipTLSVerify != tt.wantClientOptions.SkipTLSVerify {
					t.Errorf("got %v, want %v", tt.args.clientOptions.SkipTLSVerify, tt.wantClientOptions.SkipTLSVerify)
				}
			}
		})
	}
}

func TestOverrideClientFromfileConfig(t *testing.T) {
	type args struct {
		client      ConfigProviderWithLoaderAndLocation
		productName string
		location    string
	}
	tests := []struct {
		name string
		args args
		want *shared.ServerConfiguration
	}{
		{
			name: "NilfileConfig",
			args: args{
				client: &mockConfigProviderWithLoaderAndLocation{
					fileConfig: nil,
					config:     &shared.Configuration{},
				},
				productName: "testProduct",
				location:    "testLocation",
			},
			want: nil,
		},
		{
			name: "NoProductOverrides",
			args: args{
				client: &mockConfigProviderWithLoaderAndLocation{
					fileConfig: &fileconfiguration.FileConfig{},
					config:     &shared.Configuration{},
				},
				productName: "testProduct",
				location:    "testLocation",
			},
			want: nil,
		},
		{
			name: "SingleEndpoint",
			args: args{
				client: &mockConfigProviderWithLoaderAndLocation{
					fileConfig: &fileconfiguration.FileConfig{
						Environments: []fileconfiguration.Environment{
							{
								Name: "testEnv",
								Products: []fileconfiguration.Product{
									{
										Name: "testProduct",
										Endpoints: []fileconfiguration.Endpoint{
											{Name: "endpoint1", SkipTLSVerify: true, Location: "testLocation"},
										},
									},
								},
							},
						},
					},
					config: &shared.Configuration{
						HTTPClient: &http.Client{
							Transport: &http.Transport{},
						},
					},
				},
				productName: "testProduct",
				location:    "testLocation",
			},
			want: &shared.ServerConfiguration{
				URL:         "endpoint1",
				Description: shared.EndpointOverridden + "testLocation",
			},
		},
		{
			name: "MultipleEndpoints",
			args: args{
				client: &mockConfigProviderWithLoaderAndLocation{
					fileConfig: &fileconfiguration.FileConfig{
						Environments: []fileconfiguration.Environment{
							{
								Name: "testEnv",
								Products: []fileconfiguration.Product{
									{
										Name: "testProduct",
										Endpoints: []fileconfiguration.Endpoint{
											{Name: "endpoint1", SkipTLSVerify: true, Location: "testLocation"},
											{Name: "endpoint2", SkipTLSVerify: false, Location: "testLocation"},
										},
									},
								},
							},
						},
					},
					config: &shared.Configuration{
						HTTPClient: &http.Client{
							Transport: &http.Transport{},
						},
					},
				},
				productName: "testProduct",
				location:    "testLocation",
			},
			want: &shared.ServerConfiguration{
				URL:         "endpoint1",
				Description: shared.EndpointOverridden + "testLocation",
			},
		},
		{
			name: "WrongLocation",
			args: args{
				client: &mockConfigProviderWithLoaderAndLocation{
					fileConfig: &fileconfiguration.FileConfig{
						Environments: []fileconfiguration.Environment{
							{
								Name: "testEnv",
								Products: []fileconfiguration.Product{
									{
										Name: "testProduct",
										Endpoints: []fileconfiguration.Endpoint{
											{Name: "endpoint1", SkipTLSVerify: true, Location: "correctLocation"},
										},
									},
								},
							},
						},
					},
					config: &shared.Configuration{
						HTTPClient: &http.Client{
							Transport: &http.Transport{},
						},
					},
				},
				productName: "testProduct",
				location:    "wrongLocation",
			},
			want: nil,
		},
		{
			name: "EmptyLocation",
			args: args{
				client: &mockConfigProviderWithLoaderAndLocation{
					fileConfig: &fileconfiguration.FileConfig{
						Environments: []fileconfiguration.Environment{
							{
								Name: "testEnv",
								Products: []fileconfiguration.Product{
									{
										Name: "testProduct",
										Endpoints: []fileconfiguration.Endpoint{
											{Name: "endpoint1", SkipTLSVerify: true, Location: "correctLocation"},
										},
									},
								},
							},
						},
					},
					config: &shared.Configuration{
						HTTPClient: &http.Client{
							Transport: &http.Transport{},
						},
					},
				},
				productName: "testProduct",
				location:    "",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetClientOptionsFromConfig(tt.args.client, tt.args.productName, tt.args.location)
			config := tt.args.client.GetConfig()
			if tt.want == nil {
				if len(config.Servers) != 0 {
					t.Errorf("expected no servers, got %v", config.Servers)
				}
			} else {
				if len(config.Servers) != 1 {
					t.Errorf("expected one server, got %v", config.Servers)
				} else {
					if config.Servers[0].URL != tt.want.URL {
						t.Errorf("expected URL %v, got %v", tt.want.URL, config.Servers[0].URL)
					}
					if config.Servers[0].Description != tt.want.Description {
						t.Errorf("expected description %v, got %v", tt.want.Description, config.Servers[0].Description)
					}
				}
			}
		})
	}
}

type mockConfigProviderWithLoaderAndLocation struct {
	fileConfig *fileconfiguration.FileConfig
	config     *shared.Configuration
}

func (m *mockConfigProviderWithLoaderAndLocation) ChangeConfigURL(location string) {
	// Implement the logic to change the config URL based on the location
}

func (m *mockConfigProviderWithLoaderAndLocation) GetFileConfig() *fileconfiguration.FileConfig {
	return m.fileConfig
}

func (m *mockConfigProviderWithLoaderAndLocation) GetConfig() *shared.Configuration {
	return m.config
}

func TestSetClientOptionsFromfileConfigTable(t *testing.T) {
	tests := []struct {
		name              string
		clientOptions     *clientoptions.TerraformClientOptions
		fileConfig        *fileconfiguration.FileConfig
		productName       string
		wantClientOptions *clientoptions.TerraformClientOptions
	}{
		{
			name:              "NilClientOptions",
			clientOptions:     nil,
			fileConfig:        &fileconfiguration.FileConfig{},
			productName:       "testProduct",
			wantClientOptions: nil,
		},
		{
			name:              "NilfileConfig",
			clientOptions:     &clientoptions.TerraformClientOptions{},
			fileConfig:        nil,
			productName:       "testProduct",
			wantClientOptions: &clientoptions.TerraformClientOptions{},
		},
		{
			name:          "MultipleEndpoints",
			clientOptions: &clientoptions.TerraformClientOptions{},
			fileConfig: &fileconfiguration.FileConfig{
				Environments: []fileconfiguration.Environment{
					{
						Name: "testEnv",
						Products: []fileconfiguration.Product{
							{
								Name: "testProduct",
								Endpoints: []fileconfiguration.Endpoint{
									{Name: "endpoint1", SkipTLSVerify: true},
									{Name: "endpoint2", SkipTLSVerify: false},
								},
							},
						},
					},
				},
			},
			productName: "testProduct",
			wantClientOptions: &clientoptions.TerraformClientOptions{
				ClientOptions: shared.ClientOptions{
					Endpoint:      "endpoint1",
					SkipTLSVerify: true,
				},
			},
		},
		{
			name:          "SingleEndpoint",
			clientOptions: &clientoptions.TerraformClientOptions{},
			fileConfig: &fileconfiguration.FileConfig{
				Environments: []fileconfiguration.Environment{
					{
						Name: "testEnv",
						Products: []fileconfiguration.Product{
							{
								Name: "testProduct",
								Endpoints: []fileconfiguration.Endpoint{
									{Name: "endpoint1", SkipTLSVerify: true},
								},
							},
						},
					},
				},
			},
			productName: "testProduct",
			wantClientOptions: &clientoptions.TerraformClientOptions{
				ClientOptions: shared.ClientOptions{
					Endpoint:      "endpoint1",
					SkipTLSVerify: true,
				},
			},
		},
		{
			name:          "NoEndpoints",
			clientOptions: &clientoptions.TerraformClientOptions{},
			fileConfig: &fileconfiguration.FileConfig{
				Environments: []fileconfiguration.Environment{
					{
						Name: "testEnv",
						Products: []fileconfiguration.Product{
							{
								Name:      "testProduct",
								Endpoints: []fileconfiguration.Endpoint{},
							},
						},
					},
				},
			},
			productName:       "testProduct",
			wantClientOptions: &clientoptions.TerraformClientOptions{},
		},
		{
			name:          "BadProductName",
			clientOptions: &clientoptions.TerraformClientOptions{},
			fileConfig: &fileconfiguration.FileConfig{
				Environments: []fileconfiguration.Environment{
					{
						Name: "testEnv",
						Products: []fileconfiguration.Product{
							{
								Name:      "testProduct",
								Endpoints: []fileconfiguration.Endpoint{},
							},
						},
					},
				},
			},
			productName:       "productDoesNotExist",
			wantClientOptions: &clientoptions.TerraformClientOptions{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetGlobalClientOptionsFromFileConfig(tt.clientOptions, tt.fileConfig, tt.productName)
			if tt.clientOptions != nil && tt.wantClientOptions != nil {
				if tt.clientOptions.Endpoint != tt.wantClientOptions.Endpoint {
					t.Errorf("got %v, want %v", tt.clientOptions.Endpoint, tt.wantClientOptions.Endpoint)
				}
				if tt.clientOptions.SkipTLSVerify != tt.wantClientOptions.SkipTLSVerify {
					t.Errorf("got %v, want %v", tt.clientOptions.SkipTLSVerify, tt.wantClientOptions.SkipTLSVerify)
				}
			}
		})
	}
}
