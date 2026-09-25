package ionoscloud

import (
	"maps"
	"slices"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// inline_volume_ids drives both the volume-block refresh and the delete path, so a state carrying
// an empty list while an inline volume block is still declared has to be repaired: left alone the
// volume block blanks out and every later plan fails with "volume.0.disk_type attribute is
// immutable", which blocks destroy too. An empty list with no volume block is legitimate - every
// disk then belongs to a separate ionoscloud_volume resource - and must not be seeded, or a server
// delete would destroy a disk Terraform does not own.
func TestShouldSeedInlineVolumeIDs(t *testing.T) {
	tests := []struct {
		name  string
		state map[string]any // nil means no prior state at all
		want  bool
	}{
		{
			name:  "no prior state does not seed",
			state: nil,
			want:  false,
		},
		{
			name:  "attribute absent seeds (state predates 6.4.0)",
			state: map[string]any{"volume": []any{map[string]any{"name": "system"}}},
			want:  true,
		},
		{
			name: "empty list with an inline volume block seeds",
			state: map[string]any{
				"inline_volume_ids": []any{},
				"volume":            []any{map[string]any{"name": "system"}},
			},
			want: true,
		},
		{
			name: "empty list with no volume block does not seed",
			state: map[string]any{
				"inline_volume_ids": []any{},
			},
			want: false,
		},
		{
			name: "populated list does not seed",
			state: map[string]any{
				"inline_volume_ids": []any{"vol-1"},
				"volume":            []any{map[string]any{"name": "system"}},
			},
			want: false,
		},
	}

	res := resourceServer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := seedingStateData(t, res, tt.state)
			if got := shouldSeedInlineVolumeIDs(d); got != tt.want {
				t.Errorf("shouldSeedInlineVolumeIDs = %v, want %v", got, tt.want)
			}
		})
	}
}

// seedingStateData builds ResourceData backed by a raw state holding exactly the given attributes,
// which is what lets GetRawState tell "attribute absent" apart from "present but empty" - the
// distinction the seeding decision turns on.
func seedingStateData(t *testing.T, res *schema.Resource, attrs map[string]any) *schema.ResourceData {
	t.Helper()

	if attrs == nil {
		return res.Data(&terraform.InstanceState{ID: "srv-1"})
	}

	jsonMap := map[string]any{"id": "srv-1"}
	maps.Copy(jsonMap, attrs)

	raw, err := schema.JSONMapToStateValue(jsonMap, res.CoreConfigSchema())
	if err != nil {
		t.Fatalf("building raw state: %v", err)
	}
	state, err := res.ShimInstanceStateFromValue(raw)
	if err != nil {
		t.Fatalf("shimming instance state: %v", err)
	}
	state.ID = "srv-1"
	state.RawState = raw

	return res.Data(state)
}

// A server imported without a boot volume must end up with an empty (not null) ownership list:
// a null computed list is planned as "known after apply" forever, which turned the import into an
// update that sent an empty PATCH and never converged.
func TestSeedInlineVolumeIDs(t *testing.T) {
	tests := []struct {
		name  string
		state map[string]any
		want  []string
	}{
		{
			name:  "boot volume is seeded",
			state: map[string]any{"boot_volume": "vol-1"},
			want:  []string{"vol-1"},
		},
		{
			name:  "no boot volume seeds an empty list",
			state: map[string]any{"boot_volume": ""},
			want:  []string{},
		},
		{
			name:  "boot volume absent seeds an empty list",
			state: map[string]any{},
			want:  []string{},
		},
	}

	for resName, res := range map[string]*schema.Resource{"server": resourceServer(), "vcpu_server": resourceVCPUServer()} {
		for _, tt := range tests {
			t.Run(resName+"/"+tt.name, func(t *testing.T) {
				d := seedingStateData(t, res, tt.state)
				if !shouldSeedInlineVolumeIDs(d) {
					t.Fatal("shouldSeedInlineVolumeIDs = false for an imported state, want true")
				}
				got := seedInlineVolumeIDs(d)
				if got == nil || !slices.Equal(got, tt.want) {
					t.Fatalf("seedInlineVolumeIDs = %#v, want %#v", got, tt.want)
				}

				if err := d.Set("inline_volume_ids", got); err != nil {
					t.Fatalf("setting inline_volume_ids: %v", err)
				}
				wantCount := strconv.Itoa(len(tt.want))
				if count, ok := d.State().Attributes["inline_volume_ids.#"]; !ok || count != wantCount {
					t.Errorf("inline_volume_ids.# = %q (present: %v), want %q", count, ok, wantCount)
				}
			})
		}
	}
}

func TestIsEmptyServerPatch(t *testing.T) {
	name := "patched-name"
	cores := int32(2)
	multiQueue := false

	tests := []struct {
		name    string
		request ionoscloud.ServerProperties
		want    bool
	}{
		{name: "no field", request: ionoscloud.ServerProperties{}, want: true},
		{name: "name", request: ionoscloud.ServerProperties{Name: &name}, want: false},
		{name: "cores", request: ionoscloud.ServerProperties{Cores: &cores}, want: false},
		{name: "explicit false", request: ionoscloud.ServerProperties{NicMultiQueue: &multiQueue}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmptyServerPatch(tt.request); got != tt.want {
				t.Errorf("isEmptyServerPatch = %v, want %v", got, tt.want)
			}
		})
	}
}
