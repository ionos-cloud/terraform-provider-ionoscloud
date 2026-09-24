package ionoscloud

import (
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// volumeTestImages: a child location (pc/txl/1) owns no images — everything usable from
// it lives in the parent location (de/txl), so matching must work off the images.
func volumeTestImages() []ionoscloud.Image {
	img := func(id, name, location, imageType string, aliases ...string) ionoscloud.Image {
		return ionoscloud.Image{
			Id: &id,
			Properties: &ionoscloud.ImageProperties{
				Name: &name, Location: &location, ImageType: &imageType, ImageAliases: &aliases,
			},
		}
	}
	return []ionoscloud.Image{
		img("hdd-txl", "ubuntu-26.04-20260421", "de/txl", "HDD", "ubuntu:latest", "ubuntu:26.04"),
		img("iso-txl", "ubuntu-26.04-iso", "de/txl", "CDROM", "ubuntu:latest_iso"),
		img("hdd-las", "ubuntu-24.04-20240901", "us/las", "HDD", "ubuntu:latest"),
	}
}

func TestFindCompatibleVolumeImage(t *testing.T) {
	tests := []struct {
		name      string
		imageName string
		locations []string
		wantMatch string // expected matched image id, "" for none
	}{
		{
			name: "exact name match", imageName: "ubuntu-26.04-20260421",
			locations: []string{"pc/txl/1", "de/txl"}, wantMatch: "hdd-txl",
		},
		{
			name: "partial name match", imageName: "ubuntu-26.04",
			locations: []string{"pc/txl/1", "de/txl"}, wantMatch: "hdd-txl",
		},
		{
			name: "no match outside the location set", imageName: "ubuntu-26.04",
			locations: []string{"pc/txl/1"}, wantMatch: "",
		},
		{
			name: "wrong type is skipped", imageName: "ubuntu-26.04-iso",
			locations: []string{"de/txl"}, wantMatch: "",
		},
		{
			name: "empty image name matches nothing", imageName: "",
			locations: []string{"de/txl"}, wantMatch: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, _ := findCompatibleVolumeImage(tt.imageName, volumeTestImages(), tt.locations)
			gotMatch := ""
			if match != nil && match.Id != nil {
				gotMatch = *match.Id
			}
			if gotMatch != tt.wantMatch {
				t.Errorf("match = %q, want %q", gotMatch, tt.wantMatch)
			}
		})
	}
}

func TestFindCompatibleVolumeImage_ExactIDMatch(t *testing.T) {
	match, skipped := findCompatibleVolumeImage("HDD-TXL", volumeTestImages(), []string{"pc/txl/1", "de/txl"})

	if match == nil || match.Id == nil || *match.Id != "hdd-txl" {
		t.Fatalf("match id = %v, want hdd-txl", match)
	}

	if skipped != nil {
		t.Fatalf("skipped = %v, want nil", skipped)
	}
}

func TestFindCompatibleVolumeImage_ReturnsSkippedImageWhenFilteredOut(t *testing.T) {
	match, skipped := findCompatibleVolumeImage("ubuntu-26.04-iso", volumeTestImages(), []string{"de/txl"})

	if match != nil {
		t.Fatalf("match = %v, want nil", match)
	}

	if skipped == nil || skipped.Id == nil || *skipped.Id != "iso-txl" {
		t.Fatalf("skipped id = %v, want iso-txl", skipped)
	}
}

func TestFindCompatibleVolumeImage_SkippedOnWrongLocation(t *testing.T) {
	// hdd-las is HDD but in us/las — wrong location; should land in skipped, not match.
	match, skipped := findCompatibleVolumeImage("ubuntu-24.04", volumeTestImages(), []string{"de/txl"})

	if match != nil {
		t.Fatalf("match = %v, want nil", match)
	}
	if skipped == nil || skipped.Id == nil || *skipped.Id != "hdd-las" {
		t.Fatalf("skipped id = %v, want hdd-las", skipped)
	}
}

func TestFindCompatibleVolumeImage_CaseInsensitiveNameMatch(t *testing.T) {
	match, _ := findCompatibleVolumeImage("UBUNTU-26.04-20260421", volumeTestImages(), []string{"de/txl"})

	if match == nil || match.Id == nil || *match.Id != "hdd-txl" {
		t.Fatalf("match id = %v, want hdd-txl", match)
	}
}

func TestParseVolumeImportID(t *testing.T) {
	const dc, srv, vol = "dc-1", "srv-a", "vol-1"

	tests := []struct {
		name                                   string
		importID                               string
		wantLocation, wantDC, wantSrv, wantVol string
		wantErr                                bool
	}{
		{name: "attached", importID: "dc-1/srv-a/vol-1", wantDC: dc, wantSrv: srv, wantVol: vol},
		{name: "attached with location", importID: "de/fra:dc-1/srv-a/vol-1", wantLocation: "de/fra", wantDC: dc, wantSrv: srv, wantVol: vol},
		{name: "unattached", importID: "dc-1/vol-1", wantDC: dc, wantVol: vol},
		{name: "unattached with location", importID: "de/txl:dc-1/vol-1", wantLocation: "de/txl", wantDC: dc, wantVol: vol},
		{name: "volume id only", importID: "vol", wantErr: true},
		{name: "too many parts", importID: "dc-1/srv-a/vol-1/extra", wantErr: true},
		{name: "empty server part", importID: "dc-1//vol-1", wantErr: true},
		{name: "empty volume part", importID: "dc-1/", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, dcID, serverID, volumeID, err := parseVolumeImportID(tt.importID)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseVolumeImportID(%q) error = nil, want an error", tt.importID)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseVolumeImportID(%q) unexpected error: %v", tt.importID, err)
			}
			if location != tt.wantLocation || dcID != tt.wantDC || serverID != tt.wantSrv || volumeID != tt.wantVol {
				t.Errorf("parseVolumeImportID(%q) = (%q, %q, %q, %q), want (%q, %q, %q, %q)", tt.importID,
					location, dcID, serverID, volumeID, tt.wantLocation, tt.wantDC, tt.wantSrv, tt.wantVol)
			}
		})
	}
}

// server_id must stay optional: an unattached volume has no server, and import leaves it null.
func TestResourceVolumeServerIDOptional(t *testing.T) {
	s := resourceVolume().Schema["server_id"]
	if s.Required || !s.Optional {
		t.Fatalf("server_id: Required=%v Optional=%v, want optional", s.Required, s.Optional)
	}
}
