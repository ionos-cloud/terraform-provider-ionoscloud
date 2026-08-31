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

// sharedServerReadWriterKeys are the keys that setResourceServerData (the state-writer behind the
// shared resourceServerRead) calls d.Set on unconditionally, i.e. not gated on the key existing in
// the caller's schema. Every resource wired to resourceServerRead MUST declare all of them, or that
// resource crashes at plan/apply/refresh with "Invalid address to set". Keep this list in sync with
// the unconditional d.Set calls in setResourceServerData (resource_server.go).
var sharedServerReadWriterKeys = []string{"enabled_features", "confidential"}

// sharedServerReadResources are the resources that reuse resourceServerRead / setResourceServerData
// as their read path. ionoscloud_cube_server is deliberately excluded: it has its own reader
// (resourceCubeServerRead) that does not set these keys.
var sharedServerReadResources = map[string]func() *schema.Resource{
	"ionoscloud_server":      resourceServer,
	"ionoscloud_vcpu_server": resourceVCPUServer,
}

// Regression for the 6.7.36 crash (https://github.com/ionos-cloud/terraform-provider-ionoscloud):
// ionoscloud_vcpu_server reads through the shared server state-writer (resourceServerRead ->
// setResourceServerData), which unconditionally calls d.Set for enabled_features and confidential.
// Both keys must exist in every schema wired to that reader or plan/apply/refresh fails with
// "error setting enabled_features Invalid address to set". This guards the whole class, not just
// the one attribute/one resource that broke in 6.7.36.
func TestSharedServerReadWriterKeysPresentInSchemas(t *testing.T) {
	for name, ctor := range sharedServerReadResources {
		s := ctor().Schema
		for _, key := range sharedServerReadWriterKeys {
			if _, ok := s[key]; !ok {
				t.Errorf("%s schema missing %q; shared server state-writer sets it unconditionally and will crash at read", name, key)
			}
		}
	}
}

// Reproduces the exact failing calls from setResourceServerData (resource_server.go:1566, 1576)
// against every shared-reader schema: enabled_features as a []string and confidential as a bool.
// Stronger than the presence check — it also catches a type mismatch (e.g. someone re-declaring
// enabled_features as TypeString), which would still fail d.Set at runtime.
func TestSharedServerReadWriterSetRoundTrip(t *testing.T) {
	for name, ctor := range sharedServerReadResources {
		d := schema.TestResourceDataRaw(t, ctor().Schema, nil)

		if err := d.Set("enabled_features", []string{"SEV-SNP"}); err != nil {
			t.Errorf("%s: d.Set(enabled_features) = %v, want nil", name, err)
		}
		if err := d.Set("confidential", true); err != nil {
			t.Errorf("%s: d.Set(confidential) = %v, want nil", name, err)
		}

		feats := d.Get("enabled_features").([]any)
		if len(feats) != 1 || feats[0].(string) != "SEV-SNP" {
			t.Errorf("%s: enabled_features = %v, want [SEV-SNP]", name, feats)
		}
		if d.Get("confidential").(bool) != true {
			t.Errorf("%s: confidential = false, want true", name)
		}
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
