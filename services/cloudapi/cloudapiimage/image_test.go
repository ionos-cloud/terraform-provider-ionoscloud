package cloudapiimage

import (
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// testImages: a child location (pc/txl/1) owns no images — everything usable from it
// lives in the parent location (de/txl), so matching must work off the images themselves.
func testImages() []ionoscloud.Image {
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

func TestGetImageAlias(t *testing.T) {
	tests := []struct {
		name      string
		want      string
		locations []string
		expect    string
	}{
		{
			name: "alias inherited from the parent location", want: "ubuntu:latest",
			locations: []string{"pc/txl/1", "de/txl"}, expect: "ubuntu:latest",
		},
		{
			name: "alias not matched outside the location set", want: "ubuntu:latest",
			locations: []string{"pc/txl/1"}, expect: "",
		},
		{
			name: "alias on a non-HDD image still matches", want: "ubuntu:latest_iso",
			locations: []string{"pc/txl/1", "de/txl"}, expect: "ubuntu:latest_iso",
		},
		{
			name: "case-insensitive with canonical casing", want: "UBUNTU:LATEST",
			locations: []string{"de/txl"}, expect: "ubuntu:latest",
		},
		{
			name: "empty want", want: "",
			locations: []string{"de/txl"}, expect: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetImageAlias(tt.want, testImages(), tt.locations); got != tt.expect {
				t.Errorf("GetImageAlias(%q, %v) = %q, want %q", tt.want, tt.locations, got, tt.expect)
			}
		})
	}
}

func TestMatchImageAlias(t *testing.T) {
	aliases := []string{"ubuntu:latest", "debian:11", "Rocky-Linux:9"}

	tests := []struct {
		name    string
		aliases []string
		want    string
		expect  string
	}{
		{name: "exact match", aliases: aliases, want: "debian:11", expect: "debian:11"},
		{name: "case-insensitive match preserves canonical casing", aliases: aliases, want: "rocky-linux:9", expect: "Rocky-Linux:9"},
		{name: "no match", aliases: aliases, want: "centos:7", expect: ""},
		{name: "empty want", aliases: aliases, want: "", expect: ""},
		{name: "empty alias list", aliases: nil, want: "ubuntu:latest", expect: ""},
		{name: "skips empty aliases", aliases: []string{"", "ubuntu:latest"}, want: "ubuntu:latest", expect: "ubuntu:latest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchImageAlias(tt.aliases, tt.want); got != tt.expect {
				t.Fatalf("MatchImageAlias(%v, %q) = %q, want %q", tt.aliases, tt.want, got, tt.expect)
			}
		})
	}
}
