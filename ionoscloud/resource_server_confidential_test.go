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
			name:       "confidential omits cores and cpu_family",
			config:     map[string]any{"name": "coco", "type": "ENTERPRISE", "confidential": true, "ram": 4096},
			wantCores:  nil,
			wantCPUNil: true,
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
			wantCores: int32Ptr(4),
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

func int32Ptr(v int32) *int32 { return &v }

func int32PtrEqual(a, b *int32) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
