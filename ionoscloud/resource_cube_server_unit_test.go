package ionoscloud

import (
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// Cube and GPU servers share resourceCubeServerImport, which used to index the first NIC before
// checking that the server has one and panicked on a server without NICs.
func TestPrimaryIPFromNics(t *testing.T) {
	nicID := "nic-1"
	ips := []string{"198.51.100.10", "198.51.100.11"}
	withNics := func(items *[]ionoscloud.Nic) *ionoscloud.Server {
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
		{name: "no nics", server: withNics(&[]ionoscloud.Nic{})},
		{name: "nic without properties", server: withNics(&[]ionoscloud.Nic{{Id: &nicID}}), wantNicID: nicID},
		{
			name:      "nic without ips",
			server:    withNics(&[]ionoscloud.Nic{{Id: &nicID, Properties: &ionoscloud.NicProperties{}}}),
			wantNicID: nicID,
		},
		{
			name:      "nic with empty ips",
			server:    withNics(&[]ionoscloud.Nic{{Id: &nicID, Properties: &ionoscloud.NicProperties{Ips: &[]string{}}}}),
			wantNicID: nicID,
		},
		{
			name:      "first ip of the first nic",
			server:    withNics(&[]ionoscloud.Nic{{Id: &nicID, Properties: &ionoscloud.NicProperties{Ips: &ips}}}),
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
	withVolumes := func(items *[]ionoscloud.Volume) *ionoscloud.Server {
		return &ionoscloud.Server{Entities: &ionoscloud.ServerEntities{Volumes: &ionoscloud.AttachedVolumes{Items: items}}}
	}

	tests := []struct {
		name   string
		server *ionoscloud.Server
		want   string
	}{
		{name: "nil entities", server: &ionoscloud.Server{}},
		{name: "nil items", server: withVolumes(nil)},
		{name: "no volumes", server: withVolumes(&[]ionoscloud.Volume{})},
		{name: "volume without properties", server: withVolumes(&[]ionoscloud.Volume{{}})},
		{name: "volume without image", server: withVolumes(&[]ionoscloud.Volume{{Properties: &ionoscloud.VolumeProperties{}}})},
		{name: "volume with image", server: withVolumes(&[]ionoscloud.Volume{{Properties: &ionoscloud.VolumeProperties{Image: &image}}}), want: image},
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
