package ionoscloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
)

// Cube and GPU servers share resourceCubeServerImport, which used to index the first NIC before
// checking that the server has one and panicked on a server without NICs.
func TestPrimaryIPFromNics(t *testing.T) {
	nicID := "nic-1"
	ips := []string{"198.51.100.10", "198.51.100.11"}
	withNics := func(items []ionoscloud.Nic) *ionoscloud.Server {
		return &ionoscloud.Server{Entities: &ionoscloud.ServerEntities{Nics: &ionoscloud.Nics{Items: items}}}
	}

	tests := []struct {
		name      string
		server    *ionoscloud.Server
		wantNicID string
		wantIP    string
	}{
		{name: "nil server", server: nil},
		{name: "nil entities", server: &ionoscloud.Server{}},
		{name: "nil nics", server: &ionoscloud.Server{Entities: &ionoscloud.ServerEntities{}}},
		{name: "nil items", server: withNics(nil)},
		{name: "no nics", server: withNics([]ionoscloud.Nic{})},
		{name: "nic without properties", server: withNics([]ionoscloud.Nic{{Id: &nicID}}), wantNicID: nicID},
		{
			name:      "nic without ips",
			server:    withNics([]ionoscloud.Nic{{Id: &nicID, Properties: ionoscloud.NicProperties{}}}),
			wantNicID: nicID,
		},
		{
			name:      "nic with empty ips",
			server:    withNics([]ionoscloud.Nic{{Id: &nicID, Properties: ionoscloud.NicProperties{Ips: []string{}}}}),
			wantNicID: nicID,
		},
		{
			name:      "first ip of the first nic",
			server:    withNics([]ionoscloud.Nic{{Id: &nicID, Properties: ionoscloud.NicProperties{Ips: ips}}}),
			wantNicID: nicID,
			wantIP:    "198.51.100.10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNicID := ""
			if nic := firstServerNic(tt.server); nic != nil && nic.Id != nil {
				gotNicID = *nic.Id
			}
			if gotNicID != tt.wantNicID {
				t.Errorf("firstServerNic id = %q, want %q", gotNicID, tt.wantNicID)
			}

			gotIP, ok := primaryIPFromNics(tt.server)
			if gotIP != tt.wantIP || ok != (tt.wantIP != "") {
				t.Errorf("primaryIPFromNics = (%q, %v), want %q", gotIP, ok, tt.wantIP)
			}
		})
	}
}

func TestFirstVolumeImage(t *testing.T) {
	image := "img-1"
	withVolumes := func(items []ionoscloud.Volume) *ionoscloud.Server {
		return &ionoscloud.Server{Entities: &ionoscloud.ServerEntities{Volumes: &ionoscloud.AttachedVolumes{Items: items}}}
	}

	tests := []struct {
		name   string
		server *ionoscloud.Server
		want   string
	}{
		{name: "nil entities", server: &ionoscloud.Server{}},
		{name: "nil items", server: withVolumes(nil)},
		{name: "no volumes", server: withVolumes([]ionoscloud.Volume{})},
		{name: "volume without properties", server: withVolumes([]ionoscloud.Volume{{}})},
		{name: "volume without image", server: withVolumes([]ionoscloud.Volume{{Properties: &ionoscloud.VolumeProperties{}}})},
		{name: "volume with image", server: withVolumes([]ionoscloud.Volume{{Properties: &ionoscloud.VolumeProperties{Image: &image}}}), want: image},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := firstVolumeImage(tt.server)
			if got != tt.want || ok != (tt.want != "") {
				t.Errorf("firstVolumeImage = (%q, %v), want %q", got, ok, tt.want)
			}
		})
	}
}

// resourceCubeServerUpdate patched *nic.Id for a primary NIC it could not find (deleted outside
// Terraform, so primary_nic was cleared) and panicked; findServerNic must return nil instead.
func TestFindServerNic(t *testing.T) {
	id1, id2 := "nic-1", "nic-2"
	withNics := func(items []ionoscloud.Nic) *ionoscloud.Server {
		return &ionoscloud.Server{Entities: &ionoscloud.ServerEntities{Nics: &ionoscloud.Nics{Items: items}}}
	}
	nics := []ionoscloud.Nic{{}, {Id: &id1}, {Id: &id2}}

	tests := []struct {
		name   string
		server *ionoscloud.Server
		id     string
		wantID string
	}{
		{name: "nil server", server: nil, id: id1},
		{name: "nil entities", server: &ionoscloud.Server{}, id: id1},
		{name: "nil nics", server: &ionoscloud.Server{Entities: &ionoscloud.ServerEntities{}}, id: id1},
		{name: "nil items", server: withNics(nil), id: id1},
		{name: "empty id", server: withNics(nics), id: ""},
		{name: "not found", server: withNics(nics), id: "nic-3"},
		{name: "found after a nic without id", server: withNics(nics), id: id2, wantID: id2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findServerNic(tt.server, tt.id)
			if tt.wantID == "" {
				if got != nil {
					t.Fatalf("findServerNic = %+v, want nil", got)
				}
				return
			}
			if got == nil || got.Id == nil || *got.Id != tt.wantID {
				t.Fatalf("findServerNic = %+v, want id %q", got, tt.wantID)
			}
		})
	}
}

// ssh_key_path is Optional+Computed and never read back. The importer left it out of the state, so
// the plan right after import showed it as unknown and the apply PATCHed the server.
func TestCubeServerImportSetsSSHKeyPath(t *testing.T) {
	stub := newPatchBodyStub(t, map[string]string{
		"/datacenters/dc-1/servers/srv-1": `{"id":"srv-1","properties":{"name":"cube","type":"CUBE"}}`,
	})
	for name, r := range map[string]func() *schema.Resource{"cube": resourceCubeServer, "gpu": resourceGPUServer} {
		t.Run(name, func(t *testing.T) {
			d := r().Data(&terraform.InstanceState{ID: "dc-1/srv-1"})
			if _, err := resourceCubeServerImport(context.Background(), d, stubMeta(t, stub.url)); err != nil {
				t.Fatalf("Import: %v", err)
			}
			if got, ok := d.State().Attributes["ssh_key_path.#"]; !ok || got != "0" {
				t.Errorf("ssh_key_path.# = %q (set: %v), want \"0\"", got, ok)
			}
		})
	}
}
