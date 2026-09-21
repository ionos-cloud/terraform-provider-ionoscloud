package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ionoscloud_server and ionoscloud_vcpu_server share writers in both directions: an
// ionoscloud_server with type = "VCPU" is read back by the vcpu writer (serverReadForType), and
// both go through the same leaf helpers. A key one writer sets but the other's schema lacks fails
// the apply with "Invalid address to set" - that is the shape of the 6.7.36 enabled_features
// regression. These tests lock the compatibility down so the two schemas cannot drift apart again.
func TestVCPUServerSchemaCompatibleWithServer(t *testing.T) {
	srv := resourceServer().Schema
	vcpu := resourceVCPUServer().Schema

	for name, vcpuAttr := range vcpu {
		srvAttr, ok := srv[name]
		if !ok {
			t.Errorf("ionoscloud_vcpu_server has %q but ionoscloud_server does not; the shared writers would fail to set it", name)
			continue
		}
		if srvAttr.Type != vcpuAttr.Type {
			t.Errorf("attribute %q: ionoscloud_server type %s, ionoscloud_vcpu_server type %s", name, srvAttr.Type, vcpuAttr.Type)
		}
	}

	// The nested volume block is written by the shared SetVolumeProperties helper, so vcpu's keys
	// must be a subset of the enterprise ones too.
	srvVolume := nestedResource(t, srv, "volume")
	vcpuVolume := nestedResource(t, vcpu, "volume")
	for name := range vcpuVolume.Schema {
		if _, ok := srvVolume.Schema[name]; !ok {
			t.Errorf("volume.%s exists on ionoscloud_vcpu_server but not on ionoscloud_server", name)
		}
	}
}

// The 6.7.36 regression was a state-writer setting enabled_features on a resource whose schema did
// not declare it. Every resource whose read goes through setServerConfidentialVisibility must
// declare both keys it writes.
func TestConfidentialVisibilityKeysDeclared(t *testing.T) {
	resources := map[string]map[string]*schema.Schema{
		"ionoscloud_server":      resourceServer().Schema,
		"ionoscloud_vcpu_server": resourceVCPUServer().Schema,
	}

	for name, s := range resources {
		for _, key := range []string{"enabled_features", "confidential"} {
			if _, ok := s[key]; !ok {
				t.Errorf("%s is missing %q, which its state-writer sets", name, key)
			}
		}
	}
}

// Data sources reading through dataSourceServerRead share the same writer, and it sets
// enabled_features whenever the API reports features.
func TestServerDataSourcesDeclareEnabledFeatures(t *testing.T) {
	resources := map[string]map[string]*schema.Schema{
		"ionoscloud_server data source":      dataSourceServer().Schema,
		"ionoscloud_vcpu_server data source": dataSourceVCPUServer().Schema,
	}

	for name, s := range resources {
		if _, ok := s["enabled_features"]; !ok {
			t.Errorf("%s is missing enabled_features, which its shared reader sets", name)
		}
	}
}

func nestedResource(t *testing.T, s map[string]*schema.Schema, key string) *schema.Resource {
	t.Helper()

	attr, ok := s[key]
	if !ok {
		t.Fatalf("attribute %q missing", key)
	}
	nested, ok := attr.Elem.(*schema.Resource)
	if !ok {
		t.Fatalf("attribute %q is not a nested block", key)
	}
	return nested
}
