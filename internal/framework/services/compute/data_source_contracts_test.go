// //go:build all || objectstorage || objectstoragemanagement
// // +build all objectstorage objectstoragemanagement
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccDataSourceIonosCloudContracts(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIonosCloudContractsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.#"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.contract_number"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.owner"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.status"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.reg_domain"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.cores_per_server"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.ram_per_server"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.ram_per_contract"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.cores_per_contract"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.cores_provisioned"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.das_volume_provisioned"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.hdd_limit_per_contract"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.hdd_limit_per_volume"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.hdd_volume_provisioned"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.k8s_cluster_limit_total"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.k8s_clusters_provisioned"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.nat_gateway_limit_total"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.nat_gateway_provisioned"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.nlb_limit_total"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.nlb_provisioned"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.ram_provisioned"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.reservable_ips"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.reserved_ips_in_use"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.reserved_ips_on_contract"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.ssd_limit_per_contract"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.ssd_limit_per_volume"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.ssd_volume_provisioned"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.security_groups_per_vdc"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.security_groups_per_resource"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_contracts.test", "contracts.0.resource_limits.rules_per_security_group"),
				),
			},
		},
	})
}

func testAccDataSourceIonosCloudContractsConfig() string {
	return `
data "ionoscloud_contracts" "test" {}
`
}
