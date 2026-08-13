package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// datacenterResourceFactory and datacenterDataSetter are populated by the
// ionoscloud package's init() (see RegisterDatacenterBridge), which owns the
// ionoscloud_datacenter SDKv2 resource. This indirection exists to avoid an
// import cycle: internal/framework/provider depends on this package (to
// register the list resource), and ionoscloud's own test files depend on
// internal/framework/provider (to build the muxed acceptance-test provider
// server) — so this package cannot import the ionoscloud package directly.
var (
	datacenterResourceFactory func() *schema.Resource
	datacenterDataSetter      func(*schema.ResourceData, *ionoscloud.Datacenter) error
)

// RegisterDatacenterBridge wires the ionoscloud_datacenter SDKv2 resource
// constructor and field-mapping function into this package, for use by the
// datacenter list resource. Called once from the ionoscloud package's init().
func RegisterDatacenterBridge(resourceFactory func() *schema.Resource, dataSetter func(*schema.ResourceData, *ionoscloud.Datacenter) error) {
	datacenterResourceFactory = resourceFactory
	datacenterDataSetter = dataSetter
}

// datacenterToTfValues converts a Cloud API Datacenter into the identity and
// resource tftypes.Value pair expected by a list.ListResult, by round-tripping
// it through an SDKv2 ResourceData instance (the same field-mapping code the
// ionoscloud_datacenter resource and data source use).
func datacenterToTfValues(dc ionoscloud.Datacenter) (identity *tftypes.Value, resource *tftypes.Value, err error) {
	if datacenterResourceFactory == nil || datacenterDataSetter == nil {
		return nil, nil, fmt.Errorf("datacenter bridge not registered: the ionoscloud package must be imported (for its init() to call RegisterDatacenterBridge)")
	}

	rd := datacenterResourceFactory().Data(nil)

	if err := datacenterDataSetter(rd, &dc); err != nil {
		return nil, nil, fmt.Errorf("failed to populate datacenter resource data: %w", err)
	}

	if dc.Id == nil {
		return nil, nil, fmt.Errorf("datacenter has no id")
	}

	identityData, err := rd.Identity()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get identity data: %w", err)
	}
	if err := identityData.Set("id", *dc.Id); err != nil {
		return nil, nil, fmt.Errorf("failed to set identity id: %w", err)
	}

	identityVal, err := rd.TfTypeIdentityState()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert identity state: %w", err)
	}

	resourceVal, err := rd.TfTypeResourceState()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert resource state: %w", err)
	}

	return identityVal, resourceVal, nil
}
