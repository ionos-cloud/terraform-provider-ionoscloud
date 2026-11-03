package ionoscloud

import (
	"maps"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
)

// hotPlugAttrs are the six hot-plug attribute names hotPlugReconcileProps reads.
var hotPlugAttrs = []string{
	"cpu_hot_plug", "ram_hot_plug", "nic_hot_plug",
	"nic_hot_unplug", "disc_virtio_hot_plug", "disc_virtio_hot_unplug",
}

// hpObject builds a cty object holding the six hot-plug attributes. Any attribute
// not provided in vals defaults to a null bool (i.e. unset in config).
func hpObject(vals map[string]cty.Value) cty.Value {
	attrs := make(map[string]cty.Value, len(hotPlugAttrs))
	for _, a := range hotPlugAttrs {
		attrs[a] = cty.NullVal(cty.Bool)
	}
	maps.Copy(attrs, vals)
	return cty.ObjectVal(attrs)
}

// gotProps maps the resulting VolumeProperties hot-plug pointers by attribute name.
func gotProps(p ionoscloud.VolumeProperties) map[string]*bool {
	return map[string]*bool{
		"cpu_hot_plug":           p.CpuHotPlug,
		"ram_hot_plug":           p.RamHotPlug,
		"nic_hot_plug":           p.NicHotPlug,
		"nic_hot_unplug":         p.NicHotUnplug,
		"disc_virtio_hot_plug":   p.DiscVirtioHotPlug,
		"disc_virtio_hot_unplug": p.DiscVirtioHotUnplug,
	}
}

func TestHotPlugReconcileProps(t *testing.T) {
	tests := []struct {
		name       string
		rawConfig  cty.Value
		path       string
		wantSet    bool
		wantProps  map[string]bool // only non-nil fields; others asserted nil
		wantErrSub string          // non-empty => expect error containing this substring
	}{
		{
			name:      "standalone volume, root path, some attrs set",
			rawConfig: hpObject(map[string]cty.Value{"cpu_hot_plug": cty.True, "ram_hot_plug": cty.True}),
			path:      "",
			wantSet:   true,
			wantProps: map[string]bool{"cpu_hot_plug": true, "ram_hot_plug": true},
		},
		{
			name:      "explicit false is preserved, not treated as unset",
			rawConfig: hpObject(map[string]cty.Value{"cpu_hot_plug": cty.False}),
			path:      "",
			wantSet:   true,
			wantProps: map[string]bool{"cpu_hot_plug": false},
		},
		{
			name:      "all attrs unset -> nothing to reconcile",
			rawConfig: hpObject(nil),
			path:      "",
			wantSet:   false,
			wantProps: nil,
		},
		{
			name:      "unknown attr value is skipped",
			rawConfig: hpObject(map[string]cty.Value{"cpu_hot_plug": cty.UnknownVal(cty.Bool)}),
			path:      "",
			wantSet:   false,
			wantProps: nil,
		},
		{
			name:      "extra non-hot-plug attrs are ignored",
			rawConfig: cty.ObjectVal(map[string]cty.Value{"name": cty.StringVal("vol"), "cpu_hot_plug": cty.True}),
			path:      "",
			wantSet:   true,
			wantProps: map[string]bool{"cpu_hot_plug": true},
		},
		{
			name: "inline volume, list index path",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"volume": cty.TupleVal([]cty.Value{hpObject(map[string]cty.Value{"nic_hot_plug": cty.True})}),
			}),
			path:      "volume.0.",
			wantSet:   true,
			wantProps: map[string]bool{"nic_hot_plug": true},
		},
		{
			name:      "null raw config -> nothing set, no error",
			rawConfig: cty.NullVal(cty.EmptyObject),
			path:      "",
			wantSet:   false,
			wantProps: nil,
		},
		{
			name:      "unknown raw config -> nothing set, no error",
			rawConfig: cty.UnknownVal(cty.EmptyObject),
			path:      "",
			wantSet:   false,
			wantProps: nil,
		},
		{
			name: "index out of range -> error",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"volume": cty.TupleVal([]cty.Value{hpObject(nil)}),
			}),
			path:       "volume.5.",
			wantErrSub: "has 1 elements",
		},
		{
			name: "negative index -> error",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"volume": cty.TupleVal([]cty.Value{hpObject(nil)}),
			}),
			path:       "volume.-1.",
			wantErrSub: "indexes element -1",
		},
		{
			name:       "indexing a non-list -> error",
			rawConfig:  cty.ObjectVal(map[string]cty.Value{"volume": cty.True}),
			path:       "volume.0",
			wantErrSub: "which is not a list",
		},
		{
			name:       "attribute not in schema -> error",
			rawConfig:  hpObject(nil),
			path:       "does_not_exist",
			wantErrSub: "is not in the schema",
		},
		{
			name:       "descending into a non-object -> error",
			rawConfig:  hpObject(map[string]cty.Value{"cpu_hot_plug": cty.True}),
			path:       "cpu_hot_plug.foo",
			wantErrSub: "which is not an object",
		},
		{
			name: "leaf element is null -> nothing set, no error",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"volume": cty.TupleVal([]cty.Value{cty.NullVal(hpObject(nil).Type())}),
			}),
			path:      "volume.0.",
			wantSet:   false,
			wantProps: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			props, set, err := hotPlugReconcileProps(tt.rawConfig, tt.path)

			if tt.wantErrSub != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErrSub)
				}
				if !strings.Contains(err.Error(), tt.wantErrSub) {
					t.Fatalf("error %q does not contain %q", err.Error(), tt.wantErrSub)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if set != tt.wantSet {
				t.Fatalf("set = %v, want %v", set, tt.wantSet)
			}

			got := gotProps(props)
			for _, attr := range hotPlugAttrs {
				want, expected := tt.wantProps[attr]
				switch {
				case expected && got[attr] == nil:
					t.Errorf("%s: got nil, want %v", attr, want)
				case expected && *got[attr] != want:
					t.Errorf("%s: got %v, want %v", attr, *got[attr], want)
				case !expected && got[attr] != nil:
					t.Errorf("%s: got %v, want unset (nil)", attr, *got[attr])
				}
			}
		})
	}
}
