package ionoscloud

import (
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
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

func TestResolveVolumeImageName(t *testing.T) {
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
			match, _ := resolveVolumeImageName(tt.imageName, volumeTestImages(), tt.locations)
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
