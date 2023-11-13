//go:build all || autoscaling
// +build all autoscaling

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccAutoscalingGroupImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		ExternalProviders: randomProviderVersion343(),
		CheckDestroy:      testAccCheckAutoscalingGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAutoscalingGroupConfigBasic,
			},
			{
				ResourceName:            constant.AutoscalingGroupResource + "." + constant.AutoscalingGroupTestResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replica_configuration.0.volume.0.ssh_keys", "replica_configuration.0.volume.0.password", "replica_configuration.0.volume.0.image_password"},
			},
		},
	})
}
