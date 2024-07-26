package tftest

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestCheckExpectedAttrs is a custom test function to check the expected attributes of a resource
// use 'dynamic:' prefix to indicate that the value is dynamic and should be replaced with
// the actual value in your Config TestStep terraform plan
// Example:
//
//	   TestCheckExpectedAttrs("ionoscloud_nfs_cluster.example", map[string]string{
//			    "name":                        "example",
//			    "location":                    "de/txl",
//			    "size":                        "2",
//			    "nfs.0.min_version":           "4.2",
//			    "connections.0.datacenter_id": "dynamic:ionoscloud_datacenter.nfs_dc.id",
//			    "connections.0.ip_address":    "dynamic:ionoscloud_server.nfs_server.nic.0.ips.0",
//			    "connections.0.lan":           "dynamic:ionoscloud_lan.nfs_lan.id",
//	   }),
func TestCheckExpectedAttrs(resourceName string, expectedAttrs map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		for attr, expectedValue := range expectedAttrs {
			actualValue, ok := rs.Primary.Attributes[attr]
			if !ok {
				return fmt.Errorf("attribute not found: %s", attr)
			}
			if actualValue != expectedValue {
				return fmt.Errorf("expected %s to be %s, but got %s", attr, expectedValue, actualValue)
			}
		}

		return nil
	}
}
