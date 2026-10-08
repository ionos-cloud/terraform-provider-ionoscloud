package bundleclient_test

import (
	"context"
	"strings"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/failover"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
)

// newBundle creates an SdkBundle for unit tests. An optional endpoint sets
// clientOptions.Endpoint, mirroring how providerConfigure propagates IONOS_API_URL.
func newBundle(fileConfig *fileconfiguration.FileConfig, endpoint ...string) *bundleclient.SdkBundle {
	var ep string
	if len(endpoint) > 0 {
		ep = endpoint[0]
	}
	return bundleclient.New(context.Background(), clientoptions.TerraformClientOptions{
		ClientOptions: shared.ClientOptions{Endpoint: ep},
	}, fileConfig)
}

// newCloudFileConfig builds a minimal FileConfig with a Cloud product entry and the given failover options.
func newCloudFileConfig(endpoints []fileconfiguration.Endpoint, fo *failover.Options) *fileconfiguration.FileConfig {
	return &fileconfiguration.FileConfig{
		Environments: []fileconfiguration.Environment{
			{
				Name: "test",
				Products: []fileconfiguration.Product{
					{
						Name:      fileconfiguration.Cloud,
						Endpoints: endpoints,
					},
				},
			},
		},
		Failover: fo,
	}
}

func TestNewCloudAPIClientWithFailover(t *testing.T) {
	globalEp1 := fileconfiguration.Endpoint{Name: "https://ep1.example.com"}
	globalEp2 := fileconfiguration.Endpoint{Name: "https://ep2.example.com"}
	locationEp := fileconfiguration.Endpoint{Name: "https://de-fra.example.com", Location: "de/fra"}

	roundRobinFO := &failover.Options{Strategy: failover.RoundRobin}
	noneFO := &failover.Options{Strategy: failover.None}
	emptyFO := &failover.Options{}
	invalidFO := &failover.Options{Strategy: "random"}

	tests := []struct {
		name            string
		setup           func(t *testing.T)
		fileConfig      *fileconfiguration.FileConfig
		wantErr         bool
		wantErrContains string
		validateClient  func(t *testing.T, client *ionoscloud.APIClient)
	}{
		{
			name:       "nil fileConfig returns default client",
			fileConfig: nil,
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertNotFailoverTransport(t, client)
				assertDefaultServer(t, client)
			},
		},
		{
			name: "IONOS_API_URL set bypasses fileConfig and failover",
			setup: func(t *testing.T) {
				t.Setenv(shared.IonosApiUrlEnvVar, "https://custom.ionos.com")
			},
			fileConfig: newCloudFileConfig([]fileconfiguration.Endpoint{globalEp1}, roundRobinFO),
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertNotFailoverTransport(t, client)
			},
		},
		{
			name: "fileConfig without Cloud product overrides returns default client",
			fileConfig: &fileconfiguration.FileConfig{
				Environments: []fileconfiguration.Environment{
					{Name: "test", Products: []fileconfiguration.Product{{Name: "dns"}}},
				},
			},
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertNotFailoverTransport(t, client)
				assertDefaultServer(t, client)
			},
		},
		{
			name:       "Cloud product with nil failover block behaves like none strategy",
			fileConfig: newCloudFileConfig([]fileconfiguration.Endpoint{globalEp1}, nil),
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertNotFailoverTransport(t, client)
				assertServerURLs(t, client, globalEp1.Name)
			},
		},
		{
			name:       "roundRobin strategy configures failover transport with global endpoints",
			fileConfig: newCloudFileConfig([]fileconfiguration.Endpoint{globalEp1, globalEp2}, roundRobinFO),
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				if _, ok := client.GetConfig().HTTPClient.Transport.(*failover.RoundTripper); !ok {
					t.Errorf("expected *failover.RoundTripper, got %T", client.GetConfig().HTTPClient.Transport)
				}
				assertServerURLs(t, client, globalEp1.Name, globalEp2.Name)
			},
		},
		{
			name:       "none strategy uses first global endpoint transport without failover",
			fileConfig: newCloudFileConfig([]fileconfiguration.Endpoint{globalEp1, globalEp2}, noneFO),
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertNotFailoverTransport(t, client)
				assertServerURLs(t, client, globalEp1.Name)
			},
		},
		{
			name:       "empty strategy behaves like none",
			fileConfig: newCloudFileConfig([]fileconfiguration.Endpoint{globalEp1}, emptyFO),
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertNotFailoverTransport(t, client)
				assertServerURLs(t, client, globalEp1.Name)
			},
		},
		{
			name:            "no global endpoints returns error",
			fileConfig:      newCloudFileConfig([]fileconfiguration.Endpoint{locationEp}, roundRobinFO),
			wantErr:         true,
			wantErrContains: "global resources require global endpoints",
		},
		{
			name:            "invalid strategy returns descriptive error",
			fileConfig:      newCloudFileConfig([]fileconfiguration.Endpoint{globalEp1}, invalidFO),
			wantErr:         true,
			wantErrContains: "invalid failover strategy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}

			bundle := newBundle(tt.fileConfig)
			client, err := bundle.NewCloudAPIClientWithFailover(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.wantErrContains != "" && !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.wantErrContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if client == nil {
				t.Fatal("expected non-nil client, got nil")
			}
			if tt.validateClient != nil {
				tt.validateClient(t, client)
			}
		})
	}
}

func assertNotFailoverTransport(t *testing.T, client *ionoscloud.APIClient) {
	t.Helper()
	if _, ok := client.GetConfig().HTTPClient.Transport.(*failover.RoundTripper); ok {
		t.Error("expected non-failover transport, got *failover.RoundTripper")
	}
}

func assertFailoverTransport(t *testing.T, client *ionoscloud.APIClient) {
	t.Helper()
	if _, ok := client.GetConfig().HTTPClient.Transport.(*failover.RoundTripper); !ok {
		t.Errorf("expected *failover.RoundTripper, got %T", client.GetConfig().HTTPClient.Transport)
	}
}

func assertDefaultServer(t *testing.T, client *ionoscloud.APIClient) {
	t.Helper()
	servers := client.GetConfig().Servers
	if len(servers) == 0 {
		t.Fatal("expected at least one server configured")
	}
	const defaultURL = "https://api.ionos.com/cloudapi/v6"
	if servers[0].URL != defaultURL {
		t.Errorf("expected default server URL %q, got %q", defaultURL, servers[0].URL)
	}
}

func assertServerURLs(t *testing.T, client *ionoscloud.APIClient, wantURLs ...string) {
	t.Helper()
	servers := client.GetConfig().Servers
	if len(servers) != len(wantURLs) {
		t.Fatalf("expected %d servers, got %d: %v", len(wantURLs), len(servers), servers)
	}
	for i, url := range wantURLs {
		if servers[i].URL != url {
			t.Errorf("servers[%d].URL: want %q, got %q", i, url, servers[i].URL)
		}
	}
}

func TestNewCloudAPIClient(t *testing.T) {
	fra1 := fileconfiguration.Endpoint{Name: "https://fra1.example.com", Location: "de/fra"}
	fra2 := fileconfiguration.Endpoint{Name: "https://fra2.example.com", Location: "de/fra"}
	txl := fileconfiguration.Endpoint{Name: "https://txl.example.com", Location: "de/txl"}
	global1 := fileconfiguration.Endpoint{Name: "https://g1.example.com"}
	global2 := fileconfiguration.Endpoint{Name: "https://g2.example.com"}

	roundRobinFO := &failover.Options{Strategy: failover.RoundRobin}
	noneFO := &failover.Options{Strategy: failover.None}
	invalidFO := &failover.Options{Strategy: "random"}

	tests := []struct {
		name            string
		setup           func(t *testing.T)
		fileConfig      *fileconfiguration.FileConfig
		location        string
		endpoint        string
		wantErr         bool
		wantErrContains string
		validateClient  func(t *testing.T, client *ionoscloud.APIClient)
	}{
		{
			name:       "nil fileConfig returns default client",
			fileConfig: nil,
			location:   "de/fra",
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertNotFailoverTransport(t, client)
				assertDefaultServer(t, client)
			},
		},
		{
			name: "IONOS_API_URL set bypasses fileConfig and failover",
			setup: func(t *testing.T) {
				t.Setenv(shared.IonosApiUrlEnvVar, "https://custom.ionos.com")
			},
			fileConfig: newCloudFileConfig([]fileconfiguration.Endpoint{fra1}, roundRobinFO),
			location:   "de/fra",
			endpoint:   "https://custom.ionos.com",
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertNotFailoverTransport(t, client)
				assertServerURLs(t, client, "https://custom.ionos.com/cloudapi/v6")
			},
		},
		{
			name: "fileConfig without Cloud product overrides returns default client",
			fileConfig: &fileconfiguration.FileConfig{
				Environments: []fileconfiguration.Environment{
					{Name: "test", Products: []fileconfiguration.Product{{Name: "dns"}}},
				},
			},
			location: "de/fra",
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertNotFailoverTransport(t, client)
				assertDefaultServer(t, client)
			},
		},
		{
			name:       "location match with roundRobin uses location endpoints with failover",
			fileConfig: newCloudFileConfig([]fileconfiguration.Endpoint{fra1, fra2}, roundRobinFO),
			location:   "de/fra",
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertFailoverTransport(t, client)
				assertServerURLs(t, client, fra1.Name, fra2.Name)
			},
		},
		{
			name:       "location match with none uses first location endpoint",
			fileConfig: newCloudFileConfig([]fileconfiguration.Endpoint{fra1, fra2}, noneFO),
			location:   "de/fra",
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertNotFailoverTransport(t, client)
				assertServerURLs(t, client, fra1.Name)
			},
		},
		{
			name:       "selects only endpoints for the requested location",
			fileConfig: newCloudFileConfig([]fileconfiguration.Endpoint{fra1, txl}, roundRobinFO),
			location:   "de/txl",
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertServerURLs(t, client, txl.Name)
			},
		},
		{
			name:            "location not defined returns error and does not fall back to global",
			fileConfig:      newCloudFileConfig([]fileconfiguration.Endpoint{fra1}, roundRobinFO),
			location:        "us/las",
			wantErr:         true,
			wantErrContains: "no endpoint defined for location",
		},
		{
			name:       "no regional endpoints falls back to global endpoints",
			fileConfig: newCloudFileConfig([]fileconfiguration.Endpoint{global1, global2}, roundRobinFO),
			location:   "de/fra",
			validateClient: func(t *testing.T, client *ionoscloud.APIClient) {
				assertFailoverTransport(t, client)
				assertServerURLs(t, client, global1.Name, global2.Name)
			},
		},
		{
			name:            "invalid strategy returns descriptive error",
			fileConfig:      newCloudFileConfig([]fileconfiguration.Endpoint{fra1}, invalidFO),
			location:        "de/fra",
			wantErr:         true,
			wantErrContains: "invalid failover strategy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}

			bundle := newBundle(tt.fileConfig, tt.endpoint)
			client, err := bundle.NewCloudAPIClient(context.Background(), tt.location)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.wantErrContains != "" && !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.wantErrContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if client == nil {
				t.Fatal("expected non-nil client, got nil")
			}
			if tt.validateClient != nil {
				tt.validateClient(t, client)
			}
		})
	}
}
