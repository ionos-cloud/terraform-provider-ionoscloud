package ionoscloud

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// initializeCreateRequests is the pure request-builder shared by the create path. These unit
// tests cover the Confidential Computing (SEV-SNP) branch without needing an uploaded image or
// any API access: cores/cpu_family must be omitted (derived from the image) and are rejected if
// set, while the normal ENTERPRISE path still requires cores.
func TestInitializeCreateRequestsConfidential(t *testing.T) {
	tests := []struct {
		name       string
		config     map[string]any
		wantErr    string // substring; empty means no error expected
		wantCores  *int32
		wantCPUNil bool
	}{
		{
			name: "confidential omits cores and cpu_family",
			config: map[string]any{
				"name": "coco", "type": "ENTERPRISE", "confidential": true, "ram": 4096,
				"image_name": "some-sev-snp-image",
				"volume":     []any{map[string]any{"disk_type": "HDD"}},
			},
			wantCores:  nil,
			wantCPUNil: true,
		},
		{
			name: "confidential requires a volume block",
			config: map[string]any{
				"name": "coco", "type": "ENTERPRISE", "confidential": true, "ram": 4096,
				"image_name": "some-sev-snp-image",
			},
			wantErr: "confidential requires a volume block",
		},
		{
			name: "confidential requires a boot image",
			config: map[string]any{
				"name": "coco", "type": "ENTERPRISE", "confidential": true, "ram": 4096,
				"volume": []any{map[string]any{"disk_type": "HDD"}},
			},
			wantErr: "confidential requires a boot image",
		},
		{
			name:    "confidential rejects cores",
			config:  map[string]any{"name": "coco", "type": "ENTERPRISE", "confidential": true, "ram": 4096, "cores": 4},
			wantErr: "cores argument must not be set for confidential",
		},
		{
			name:    "confidential rejects cpu_family",
			config:  map[string]any{"name": "coco", "type": "ENTERPRISE", "confidential": true, "ram": 4096, "cpu_family": "AMD_TURIN"},
			wantErr: "cpu_family argument must not be set for confidential",
		},
		{
			name:    "non-confidential enterprise still requires cores",
			config:  map[string]any{"name": "srv", "type": "ENTERPRISE", "ram": 4096},
			wantErr: "cores argument is required",
		},
		{
			name:      "non-confidential enterprise sets cores",
			config:    map[string]any{"name": "srv", "type": "ENTERPRISE", "ram": 4096, "cores": 4},
			wantCores: new(int32(4)),
		},
		{
			name:    "confidential rejected for non-enterprise type",
			config:  map[string]any{"name": "coco", "type": "VCPU", "confidential": true, "ram": 4096},
			wantErr: "only supported for ENTERPRISE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceServer().Schema, tt.config)

			server, err := initializeCreateRequests(d)

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server.Properties.Type == nil || *server.Properties.Type != "ENTERPRISE" {
				t.Errorf("type = %v, want ENTERPRISE", server.Properties.Type)
			}
			if !int32PtrEqual(server.Properties.Cores, tt.wantCores) {
				t.Errorf("cores = %v, want %v", server.Properties.Cores, tt.wantCores)
			}
			if tt.wantCPUNil && server.Properties.CpuFamily != nil {
				t.Errorf("cpu_family = %v, want nil", *server.Properties.CpuFamily)
			}
		})
	}
}

// confidential must be Computed + ForceNew: Computed lets Read derive it from the API so an
// imported/drifted server does not trigger a spurious destroy+recreate; ForceNew makes an actual
// change replace the server.
func TestConfidentialSchemaComputedForceNew(t *testing.T) {
	s := resourceServer().Schema["confidential"]
	if s == nil {
		t.Fatal("confidential attribute missing from schema")
	}
	if !s.Computed {
		t.Error("confidential must be Computed so read-back avoids a spurious ForceNew replace on import/drift")
	}
	if !s.ForceNew {
		t.Error("confidential must be ForceNew")
	}
}

func int32PtrEqual(a, b *int32) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
