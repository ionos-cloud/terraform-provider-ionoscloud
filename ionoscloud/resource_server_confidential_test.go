package ionoscloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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

// A Confidential Computing server must be deleted together with its volumes: the API refuses to
// delete one that would leave its confidential boot volume behind (VDC-5-2060), and the volume
// cannot be detached while attached. This drives the request through a stub API so the actual
// deleteVolumes query parameter on the wire is asserted, not just the builder call.
func TestDeleteServerRequestDeleteVolumes(t *testing.T) {
	tests := []struct {
		name         string
		confidential bool
		want         string
	}{
		{name: "confidential deletes its volumes with the server", confidential: true, want: "true"},
		{name: "normal server keeps its volumes", confidential: false, want: "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath, gotDeleteVolumes string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.Method + " " + r.URL.Path
				gotDeleteVolumes = r.URL.Query().Get("deleteVolumes")
				w.WriteHeader(http.StatusAccepted)
			}))
			defer srv.Close()

			cfg := ionoscloud.NewConfiguration("", "", "token", srv.URL)
			client := ionoscloud.NewAPIClient(cfg)

			_, err := deleteServerRequest(context.Background(), client, "dc-id", "server-id", tt.confidential).Execute()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if want := "DELETE /cloudapi/v6/datacenters/dc-id/servers/server-id"; gotPath != want {
				t.Errorf("request = %q, want %q", gotPath, want)
			}
			if gotDeleteVolumes != tt.want {
				t.Errorf("deleteVolumes = %q, want %q", gotDeleteVolumes, tt.want)
			}
		})
	}
}

// deleteVolumes on the server delete takes down every volume still attached, so the confidential
// delete path has to detach the ones Terraform does not own first. Only inline volume blocks are
// owned; anything else belongs to a separate ionoscloud_volume resource and must survive.
func TestDetachableVolumeIDs(t *testing.T) {
	volumes := func(ids ...string) *ionoscloud.Server {
		items := make([]ionoscloud.Volume, 0, len(ids))
		for _, id := range ids {
			items = append(items, ionoscloud.Volume{Id: new(id)})
		}
		return &ionoscloud.Server{Entities: &ionoscloud.ServerEntities{Volumes: &ionoscloud.AttachedVolumes{Items: &items}}}
	}

	tests := []struct {
		name   string
		server *ionoscloud.Server
		inline []any
		want   []string
	}{
		{
			name:   "boot volume only is owned, nothing to detach",
			server: volumes("boot"),
			inline: []any{"boot"},
			want:   nil,
		},
		{
			name:   "separately managed data disk is detached",
			server: volumes("boot", "data"),
			inline: []any{"boot"},
			want:   []string{"data"},
		},
		{
			name:   "several foreign volumes are all detached",
			server: volumes("boot", "data", "backup"),
			inline: []any{"boot"},
			want:   []string{"data", "backup"},
		},
		{
			name:   "no inline volumes means nothing is owned",
			server: volumes("data"),
			inline: nil,
			want:   []string{"data"},
		},
		{
			name:   "nil server",
			server: nil,
			inline: []any{"boot"},
			want:   nil,
		},
		{
			name:   "server without volumes",
			server: &ionoscloud.Server{},
			inline: []any{"boot"},
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detachableVolumeIDs(tt.server, tt.inline)
			if len(got) != len(tt.want) {
				t.Fatalf("detachableVolumeIDs = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("detachableVolumeIDs = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
