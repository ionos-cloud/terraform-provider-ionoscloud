package cloudapinic

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

// injectRawConfig injects the raw config into the resource data mock.
func injectRawConfig(d *schema.ResourceData, config cty.Value) {
	diff := &terraform.InstanceDiff{
		RawConfig: config,
	}
	v := reflect.ValueOf(d).Elem()
	f := v.FieldByName("diff")
	// 'diff' is a private struct field, changing its value requires the use of pointers.
	cf := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	cf.Set(reflect.ValueOf(diff))
}

// toCtyStringValSlice converts a slice of strings into a slice of values
func toCtyStringValSlice(s []string) []cty.Value {
	r := make([]cty.Value, len(s))
	for i, v := range s {
		r[i] = cty.StringVal(v)
	}
	return r
}

// getResourceData returns a resource data mock.
func getResourceData(t *testing.T, testSchema map[string]*schema.Schema, ipAttr string, ips []string) *schema.ResourceData {
	var config map[string]interface{}
	d := schema.TestResourceDataRaw(t, testSchema, config)
	if ipAttr != "ips" && ipAttr != "ipv6_ips" {
		// Error case with an invalid attribute, proper configuration doesn't matter.
		return d
	}
	ipsSet := schema.NewSet(schema.HashString, utils.ToInterfaceSlice(ips))
	if err := d.Set(ipAttr, ipsSet); err != nil {
		t.Fatalf("error setting ips: %s", err)
	}
	// must not call SetVal with empty slice.
	if len(ips) > 0 {
		rawConfigVal := cty.ObjectVal(map[string]cty.Value{
			ipAttr: cty.SetVal(toCtyStringValSlice(ips)),
		})
		injectRawConfig(d, rawConfigVal)
	}
	return d
}

func TestGetNicIPsFromSchema(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"ips": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Computed: true,
		},
		"ipv6_ips": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Computed: true,
		},
	}
	tests := []struct {
		name          string
		ips           []string
		ipAttr        string
		partialErrMsg string
	}{
		{
			name:          "valid IPs",
			ips:           []string{"192.168.8.5", "192.168.8.4", "192.168.8.3"},
			ipAttr:        "ips",
			partialErrMsg: "",
		},
		{
			name:          "valid IPv6 IPs",
			ips:           []string{"2001:db8:85a3:0:0:8a2e:370:7333", "2001:db8:85a3:0:0:8a2e:370:7334"},
			ipAttr:        "ipv6_ips",
			partialErrMsg: "",
		},
		{
			name:          "no IPs",
			ips:           []string{},
			ipAttr:        "ips",
			partialErrMsg: "expected a valid configuration but received null instead",
		},
		{
			name:          "no IPv6 IPs",
			ips:           []string{},
			ipAttr:        "ipv6_ips",
			partialErrMsg: "expected a valid configuration but received null instead",
		},
		{
			name:          "invalid IP attribute",
			ips:           []string{},
			ipAttr:        "invalid_ip_attribute",
			partialErrMsg: "provided attribute is not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := getResourceData(t, testSchema, tt.ipAttr, tt.ips)
			result, err := GetNicIPsFromSchema(d, tt.ipAttr)
			if tt.partialErrMsg == "" && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.partialErrMsg != "" && err == nil {
				t.Errorf("Expected error containing the following message: %v, got no error", tt.partialErrMsg)
			}
			if tt.partialErrMsg == "" && !reflect.DeepEqual(result, tt.ips) {
				t.Errorf("Expected: %v, got: %v", tt.ips, result)
			}
		})
	}
}
