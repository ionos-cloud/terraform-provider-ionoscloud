package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// These unit tests cover the Confidential Computing visibility wiring (enabled_features on the
// server, cpu_architecture.enabled_features on the datacenter) without any API access.

func TestSetServerPropertiesEnabledFeatures(t *testing.T) {
	server := ionoscloud.Server{
		Properties: &ionoscloud.ServerProperties{
			EnabledFeatures: &[]string{"SEV-SNP"},
		},
	}

	m := SetServerProperties(server)

	feats, ok := m["enabled_features"].([]string)
	if !ok {
		t.Fatalf("enabled_features = %T, want []string", m["enabled_features"])
	}
	if len(feats) != 1 || feats[0] != "SEV-SNP" {
		t.Errorf("enabled_features = %v, want [SEV-SNP]", feats)
	}
}

func TestSetDatacenterDataCPUArchEnabledFeatures(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceDatacenter().Schema, nil)

	dc := &ionoscloud.Datacenter{
		Properties: &ionoscloud.DatacenterProperties{
			CpuArchitecture: &[]ionoscloud.CpuArchitectureProperties{
				{
					CpuFamily:       new("AMD_TURIN"),
					EnabledFeatures: &[]string{"SEV-SNP"},
				},
			},
		},
	}

	if err := setDatacenterData(d, dc); err != nil {
		t.Fatalf("setDatacenterData: %v", err)
	}

	arches := d.Get("cpu_architecture").([]any)
	if len(arches) != 1 {
		t.Fatalf("cpu_architecture len = %d, want 1", len(arches))
	}
	entry := arches[0].(map[string]any)
	feats := entry["enabled_features"].([]any)
	if len(feats) != 1 || feats[0].(string) != "SEV-SNP" {
		t.Errorf("enabled_features = %v, want [SEV-SNP]", feats)
	}
}

// serverIsConfidential drives the destroy volume-ordering, derived from the API-reported features
// so it stays correct for imported/drifted servers.
func TestServerIsConfidential(t *testing.T) {
	tests := []struct {
		name   string
		server ionoscloud.Server
		want   bool
	}{
		{name: "nil properties", server: ionoscloud.Server{}, want: false},
		{
			name:   "no features",
			server: ionoscloud.Server{Properties: &ionoscloud.ServerProperties{}},
			want:   false,
		},
		{
			name:   "other feature only",
			server: ionoscloud.Server{Properties: &ionoscloud.ServerProperties{EnabledFeatures: &[]string{"SOMETHING"}}},
			want:   false,
		},
		{
			name:   "sev-snp present",
			server: ionoscloud.Server{Properties: &ionoscloud.ServerProperties{EnabledFeatures: &[]string{"SEV-SNP"}}},
			want:   true,
		},
		{
			name:   "case-insensitive match",
			server: ionoscloud.Server{Properties: &ionoscloud.ServerProperties{EnabledFeatures: &[]string{"sev-snp"}}},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serverIsConfidential(&tt.server); got != tt.want {
				t.Errorf("serverIsConfidential = %v, want %v", got, tt.want)
			}
		})
	}
}
