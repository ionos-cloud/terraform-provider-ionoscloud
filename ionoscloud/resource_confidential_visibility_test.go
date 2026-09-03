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

// serverReadWriterKeys are the keys the server state-writers (setResourceServerData for enterprise,
// setResourceVCPUServerData for vcpu) call d.Set on unconditionally, i.e. not gated on the key
// existing in the schema. Every server resource whose writer sets them MUST declare all of them, or
// that resource crashes at plan/apply/refresh with "Invalid address to set". Keep this in sync with
// the unconditional d.Set calls in both writers.
var serverReadWriterKeys = []string{"enabled_features", "confidential"}

// serverConfidentialAwareResources are the resources whose state-writer sets serverReadWriterKeys.
// As of the un-sharing, ionoscloud_server and ionoscloud_vcpu_server have SEPARATE writers, but both
// still set these keys, so both must declare them. ionoscloud_cube_server / ionoscloud_gpu_server are
// excluded: their reader (resourceCubeServerRead) sets neither key.
var serverConfidentialAwareResources = map[string]func() *schema.Resource{
	"ionoscloud_server":      resourceServer,
	"ionoscloud_vcpu_server": resourceVCPUServer,
}

// Regression for the 6.7.36 crash (https://github.com/ionos-cloud/terraform-provider-ionoscloud):
// the vcpu server crashed because its schema lacked enabled_features/confidential while the writer
// set them unconditionally. Whether the writer is shared (pre-fix) or per-type (post un-sharing),
// every server whose writer sets these keys must declare them, or plan/apply/refresh fails with
// "error setting enabled_features Invalid address to set". This guards the whole class, not just
// the one attribute/one resource that broke in 6.7.36.
func TestSharedServerReadWriterKeysPresentInSchemas(t *testing.T) {
	for name, ctor := range serverConfidentialAwareResources {
		s := ctor().Schema
		for _, key := range serverReadWriterKeys {
			if _, ok := s[key]; !ok {
				t.Errorf("%s schema missing %q; its server state-writer sets it unconditionally and will crash at read", name, key)
			}
		}
	}
}

// Exercises the ACTUAL writer helper both server state-writers call (setServerConfidentialVisibility)
// against every consuming schema, with a populated Server carrying SEV-SNP. This is the real drift
// guard after un-sharing: it runs the writer path (not a bare d.Set), so it catches the vcpu writer
// dropping/misspelling a key, a type mismatch, or confidential being derived wrong — the exact class
// that broke in 6.7.36. Both writers delegate here, so testing the helper covers both.
func TestSetServerConfidentialVisibilityAcrossSchemas(t *testing.T) {
	server := &ionoscloud.Server{
		Properties: &ionoscloud.ServerProperties{
			EnabledFeatures: &[]string{"SEV-SNP"},
		},
	}

	for name, ctor := range serverConfidentialAwareResources {
		d := schema.TestResourceDataRaw(t, ctor().Schema, nil)

		if err := setServerConfidentialVisibility(d, server); err != nil {
			t.Errorf("%s: setServerConfidentialVisibility = %v, want nil", name, err)
			continue
		}

		feats := d.Get("enabled_features").([]any)
		if len(feats) != 1 || feats[0].(string) != "SEV-SNP" {
			t.Errorf("%s: enabled_features = %v, want [SEV-SNP]", name, feats)
		}
		if d.Get("confidential").(bool) != true {
			t.Errorf("%s: confidential = false, want true (SEV-SNP present)", name)
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
