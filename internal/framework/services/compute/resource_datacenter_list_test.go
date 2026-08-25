package compute

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func TestDatacenterToTfValues(t *testing.T) {
	id := "11111111-1111-1111-1111-111111111111"
	name := "my-dc"

	dc := ionoscloud.Datacenter{
		Id: &id,
		Properties: &ionoscloud.DatacenterProperties{
			Name: &name,
		},
	}

	identityVal, resourceVal, err := datacenterToTfValues(dc)
	if err != nil {
		t.Fatalf("datacenterToTfValues returned error: %v", err)
	}
	if identityVal == nil {
		t.Fatal("expected non-nil identity value")
	}
	if resourceVal == nil {
		t.Fatal("expected non-nil resource value")
	}

	var identityMap map[string]tftypes.Value
	if err := identityVal.As(&identityMap); err != nil {
		t.Fatalf("failed to decode identity value: %v", err)
	}
	idVal, ok := identityMap["id"]
	if !ok {
		t.Fatal("expected identity value to contain \"id\"")
	}
	var idStr string
	if err := idVal.As(&idStr); err != nil {
		t.Fatalf("failed to decode identity id: %v", err)
	}
	if idStr != id {
		t.Errorf("expected identity id %q, got %q", id, idStr)
	}

	var resourceMap map[string]tftypes.Value
	if err := resourceVal.As(&resourceMap); err != nil {
		t.Fatalf("failed to decode resource value: %v", err)
	}
	resourceIDVal, ok := resourceMap["id"]
	if !ok {
		t.Fatal("expected resource value to contain \"id\"")
	}
	var resourceIDStr string
	if err := resourceIDVal.As(&resourceIDStr); err != nil {
		t.Fatalf("failed to decode resource id: %v", err)
	}
	if resourceIDStr != id {
		t.Errorf("expected resource id %q, got %q", id, resourceIDStr)
	}

	nameVal, ok := resourceMap["name"]
	if !ok {
		t.Fatal("expected resource value to contain \"name\"")
	}
	var nameStr string
	if err := nameVal.As(&nameStr); err != nil {
		t.Fatalf("failed to decode resource name: %v", err)
	}
	if nameStr != name {
		t.Errorf("expected resource name %q, got %q", name, nameStr)
	}
}

func TestDatacenterToTfValuesRequiresID(t *testing.T) {
	dc := ionoscloud.Datacenter{}

	if _, _, err := datacenterToTfValues(dc); err == nil {
		t.Fatal("expected an error when the datacenter has no id")
	}
}
