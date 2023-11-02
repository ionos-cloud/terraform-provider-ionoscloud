//go:build all || autoscaling
// +build all autoscaling

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAutoscalingGroupImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAutoscalingGroupConfigBasic,
			},
			{
				ResourceName:            AutoscalingGroupResource + "." + AutoscalingGroupTestResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replica_configuration.0.volume.0.ssh_key_paths", "replica_configuration.0.volume.0.ssh_key_values", "replica_configuration.0.volume.0.password"},
			},
		},
	})
}
