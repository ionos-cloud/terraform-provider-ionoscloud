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
					CpuFamily:       strPtr("AMD_TURIN"),
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

func strPtr(s string) *string { return &s }
