package ionoscloud

import (
	"maps"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
