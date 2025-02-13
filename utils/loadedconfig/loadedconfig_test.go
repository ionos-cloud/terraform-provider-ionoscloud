package loadedconfig

import (
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"net/http"
	"testing"
)

func TestSetClientOptionsFromLoadedConfig(t *testing.T) {
	type args struct {
		clientOptions *bundle.ClientOptions
		loadedConfig  *shared.LoadedConfig
		productName   string
	}
	tests := []struct {
		name              string
		args              args
		wantClientOptions *bundle.ClientOptions
	}{
		{
			name: "NilClientOptions",
			args: args{
				clientOptions: nil,
				loadedConfig:  &shared.LoadedConfig{},
				productName:   "testProduct",
			},
			wantClientOptions: nil,
		},
		{
			name: "NilLoadedConfig",
			args: args{
				clientOptions: &bundle.ClientOptions{},
				loadedConfig:  nil,
				productName:   "testProduct",
			},
			wantClientOptions: &bundle.ClientOptions{},
		},
		{
			name: "MultipleEndpoints",
			args: args{
				clientOptions: &bundle.ClientOptions{},
				loadedConfig: &shared.LoadedConfig{
					Environments: []shared.Environment{
						{
							Name: "testEnv",
							Products: []shared.Product{
								{
									Name: "testProduct",
									Endpoints: []shared.Endpoint{
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
			wantClientOptions: &bundle.ClientOptions{
				ClientOverrideOptions: shared.ClientOverrideOptions{
					Endpoint:      "endpoint1",
					SkipTLSVerify: true,
				},
			},
		},
		{
			name: "SingleEndpoint",
			args: args{
				clientOptions: &bundle.ClientOptions{},
				loadedConfig: &shared.LoadedConfig{
					Environments: []shared.Environment{
						{
							Name: "testEnv",
							Products: []shared.Product{
								{
									Name: "testProduct",
									Endpoints: []shared.Endpoint{
										{Name: "endpoint1", SkipTLSVerify: true},
									},
								},
							},
						},
					},
				},
				productName: "testProduct",
			},
			wantClientOptions: &bundle.ClientOptions{
				ClientOverrideOptions: shared.ClientOverrideOptions{
					Endpoint:      "endpoint1",
					SkipTLSVerify: true,
				},
			},
		},
		{
			name: "NoEndpoints",
			args: args{
				clientOptions: &bundle.ClientOptions{},
				loadedConfig: &shared.LoadedConfig{
					Environments: []shared.Environment{
						{
							Name: "testEnv",
							Products: []shared.Product{
								{
									Name:      "testProduct",
									Endpoints: []shared.Endpoint{},
								},
							},
						},
					},
				},
				productName: "testProduct",
			},
			wantClientOptions: &bundle.ClientOptions{},
		},
		{
			name: "BadProductName",
			args: args{
				clientOptions: &bundle.ClientOptions{},
				loadedConfig: &shared.LoadedConfig{
					Environments: []shared.Environment{
						{
							Name: "testEnv",
							Products: []shared.Product{
								{
									Name:      "testProduct",
									Endpoints: []shared.Endpoint{},
								},
							},
						},
					},
				},
				productName: "productDoesNotExist",
			},
			wantClientOptions: &bundle.ClientOptions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetClientOptionsFromLoadedConfig(tt.args.clientOptions, tt.args.loadedConfig, tt.args.productName)
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

func TestOverrideClientFromLoadedConfig(t *testing.T) {
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
			name: "NilLoadedConfig",
			args: args{
				client: &mockConfigProviderWithLoaderAndLocation{
					loadedConfig: nil,
					config:       &shared.Configuration{},
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
					loadedConfig: &shared.LoadedConfig{},
					config:       &shared.Configuration{},
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
					loadedConfig: &shared.LoadedConfig{
						Environments: []shared.Environment{
							{
								Name: "testEnv",
								Products: []shared.Product{
									{
										Name: "testProduct",
										Endpoints: []shared.Endpoint{
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
					loadedConfig: &shared.LoadedConfig{
						Environments: []shared.Environment{
							{
								Name: "testEnv",
								Products: []shared.Product{
									{
										Name: "testProduct",
										Endpoints: []shared.Endpoint{
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
					loadedConfig: &shared.LoadedConfig{
						Environments: []shared.Environment{
							{
								Name: "testEnv",
								Products: []shared.Product{
									{
										Name: "testProduct",
										Endpoints: []shared.Endpoint{
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
					loadedConfig: &shared.LoadedConfig{
						Environments: []shared.Environment{
							{
								Name: "testEnv",
								Products: []shared.Product{
									{
										Name: "testProduct",
										Endpoints: []shared.Endpoint{
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
			OverrideClientEndpoint(tt.args.client, tt.args.productName, tt.args.location)
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
	loadedConfig *shared.LoadedConfig
	config       *shared.Configuration
}

func (m *mockConfigProviderWithLoaderAndLocation) ChangeConfigURL(location string) {
	// Implement the logic to change the config URL based on the location
}

func (m *mockConfigProviderWithLoaderAndLocation) GetLoadedConfig() *shared.LoadedConfig {
	return m.loadedConfig
}

func (m *mockConfigProviderWithLoaderAndLocation) GetConfig() *shared.Configuration {
	return m.config
}

func TestSetClientOptionsFromLoadedConfigTable(t *testing.T) {
	tests := []struct {
		name              string
		clientOptions     *bundle.ClientOptions
		loadedConfig      *shared.LoadedConfig
		productName       string
		wantClientOptions *bundle.ClientOptions
	}{
		{
			name:              "NilClientOptions",
			clientOptions:     nil,
			loadedConfig:      &shared.LoadedConfig{},
			productName:       "testProduct",
			wantClientOptions: nil,
		},
		{
			name:              "NilLoadedConfig",
			clientOptions:     &bundle.ClientOptions{},
			loadedConfig:      nil,
			productName:       "testProduct",
			wantClientOptions: &bundle.ClientOptions{},
		},
		{
			name:          "MultipleEndpoints",
			clientOptions: &bundle.ClientOptions{},
			loadedConfig: &shared.LoadedConfig{
				Environments: []shared.Environment{
					{
						Name: "testEnv",
						Products: []shared.Product{
							{
								Name: "testProduct",
								Endpoints: []shared.Endpoint{
									{Name: "endpoint1", SkipTLSVerify: true},
									{Name: "endpoint2", SkipTLSVerify: false},
								},
							},
						},
					},
				},
			},
			productName: "testProduct",
			wantClientOptions: &bundle.ClientOptions{
				ClientOverrideOptions: shared.ClientOverrideOptions{
					Endpoint:      "endpoint1",
					SkipTLSVerify: true,
				},
			},
		},
		{
			name:          "SingleEndpoint",
			clientOptions: &bundle.ClientOptions{},
			loadedConfig: &shared.LoadedConfig{
				Environments: []shared.Environment{
					{
						Name: "testEnv",
						Products: []shared.Product{
							{
								Name: "testProduct",
								Endpoints: []shared.Endpoint{
									{Name: "endpoint1", SkipTLSVerify: true},
								},
							},
						},
					},
				},
			},
			productName: "testProduct",
			wantClientOptions: &bundle.ClientOptions{
				ClientOverrideOptions: shared.ClientOverrideOptions{
					Endpoint:      "endpoint1",
					SkipTLSVerify: true,
				},
			},
		},
		{
			name:          "NoEndpoints",
			clientOptions: &bundle.ClientOptions{},
			loadedConfig: &shared.LoadedConfig{
				Environments: []shared.Environment{
					{
						Name: "testEnv",
						Products: []shared.Product{
							{
								Name:      "testProduct",
								Endpoints: []shared.Endpoint{},
							},
						},
					},
				},
			},
			productName:       "testProduct",
			wantClientOptions: &bundle.ClientOptions{},
		},
		{
			name:          "BadProductName",
			clientOptions: &bundle.ClientOptions{},
			loadedConfig: &shared.LoadedConfig{
				Environments: []shared.Environment{
					{
						Name: "testEnv",
						Products: []shared.Product{
							{
								Name:      "testProduct",
								Endpoints: []shared.Endpoint{},
							},
						},
					},
				},
			},
			productName:       "productDoesNotExist",
			wantClientOptions: &bundle.ClientOptions{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetClientOptionsFromLoadedConfig(tt.clientOptions, tt.loadedConfig, tt.productName)
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
