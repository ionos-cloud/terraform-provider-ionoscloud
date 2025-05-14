//go:build all || k8s
// +build all k8s

package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccK8sNodePoolBasic(t *testing.T) {
	var k8sNodepool ionoscloud.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckK8sNodePoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sNodePoolConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(constant.ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "name", constant.K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week", "Monday"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "maintenance_window.0.time", "09:00:00Z"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "auto_scaling.#", "0"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "server_type", "DedicatedCore"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "node_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cores_count", "2"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_size", "40"),
					resource.TestCheckResourceAttrPair(constant.ResourceNameK8sNodePool, "lans.0.id", constant.LanResource+".terraform_acctest", "id"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.0.dhcp", "true"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.0.routes.0.network", "1.2.3.5/24"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.0.routes.0.gateway_ip", "10.1.5.17"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "labels.foo", "bar"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "labels.color", "green"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "annotations.ann1", "value1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "annotations.ann2", "value2"),
				),
			},
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolNodesMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(constant.DataSourceK8sNodePoolNodesId, "nodes.0.public_ip"),
					resource.TestCheckResourceAttrSet(constant.DataSourceK8sNodePoolNodesId, "nodes.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolNodesId, "nodes.0.k8s_version", constant.ResourceNameK8sNodePool, "k8s_version"),
				),
			},
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "name", constant.ResourceNameK8sNodePool, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "k8s_version", constant.ResourceNameK8sNodePool, "k8s_version"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "maintenance_window.0.day_of_the_week", constant.ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "maintenance_window.0.time", constant.ResourceNameK8sNodePool, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "server_type", constant.ResourceNameK8sNodePool, "server_type"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "availability_zone", constant.ResourceNameK8sNodePool, "availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "storage_type", constant.ResourceNameK8sNodePool, "storage_type"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "node_count", constant.ResourceNameK8sNodePool, "node_count"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "cores_count", constant.ResourceNameK8sNodePool, "cores_count"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "ram_size", constant.ResourceNameK8sNodePool, "ram_size"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "storage_size", constant.ResourceNameK8sNodePool, "storage_size"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "lans.0", constant.ResourceNameK8sNodePool, "lans.0"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "labels.foo", constant.ResourceNameK8sNodePool, "labels.foo"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "labels.color", constant.ResourceNameK8sNodePool, "labels.color"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "annotations.ann1", constant.ResourceNameK8sNodePool, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolId, "annotations.ann2", constant.ResourceNameK8sNodePool, "annotations.ann2"),
				),
			},
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "name", constant.ResourceNameK8sNodePool, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "k8s_version", constant.ResourceNameK8sNodePool, "k8s_version"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "maintenance_window.0.day_of_the_week", constant.ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "maintenance_window.0.time", constant.ResourceNameK8sNodePool, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "server_type", constant.ResourceNameK8sNodePool, "server_type"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "availability_zone", constant.ResourceNameK8sNodePool, "availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "storage_type", constant.ResourceNameK8sNodePool, "storage_type"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "node_count", constant.ResourceNameK8sNodePool, "node_count"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "cores_count", constant.ResourceNameK8sNodePool, "cores_count"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "ram_size", constant.ResourceNameK8sNodePool, "ram_size"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "storage_size", constant.ResourceNameK8sNodePool, "storage_size"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "lans.0", constant.ResourceNameK8sNodePool, "lans.0"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "labels.foo", constant.ResourceNameK8sNodePool, "labels.foo"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "labels.color", constant.ResourceNameK8sNodePool, "labels.color"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "annotations.ann1", constant.ResourceNameK8sNodePool, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(constant.DataSourceK8sNodePoolName, "annotations.ann2", constant.ResourceNameK8sNodePool, "annotations.ann2"),
				),
			},
			{
				Config:      testAccDataSourceProfitBricksK8sNodePoolWrongNameError,
				ExpectError: regexp.MustCompile("no nodepool found with the specified name"),
			},
			{
				Config: testAccCheckK8sNodePoolConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(constant.ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "name", constant.K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week", "Tuesday"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "maintenance_window.0.time", "10:00:00Z"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "auto_scaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "auto_scaling.0.max_node_count", "2"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "server_type", "VCPU"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "node_count", "2"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cores_count", "2"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_size", "40"),
					resource.TestCheckResourceAttrPair(constant.ResourceNameK8sNodePool, "public_ips.0", constant.IpBlockResource+".terraform_acctest", "ips.0"),
					resource.TestCheckResourceAttrPair(constant.ResourceNameK8sNodePool, "public_ips.1", constant.IpBlockResource+".terraform_acctest", "ips.1"),
					resource.TestCheckResourceAttrPair(constant.ResourceNameK8sNodePool, "lans.0.id", constant.LanResource+".terraform_acctest", "id"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.0.dhcp", "false"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.0.routes.0.network", "1.2.3.4/24"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.0.routes.0.gateway_ip", "10.1.5.16"),
					resource.TestCheckResourceAttrPair(constant.ResourceNameK8sNodePool, "lans.1.id", constant.LanResource+".terraform_acctest_updated", "id"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.1.dhcp", "false"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.1.routes.0.network", "1.2.3.5/24"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.1.routes.0.gateway_ip", "10.1.5.17"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.1.routes.1.network", "1.2.3.6/24"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.1.routes.1.gateway_ip", "10.1.5.18"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "labels.foo", "baz"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "labels.color", "red"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "labels.third", "thirdValue"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "annotations.ann1", "value1Changed"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "annotations.ann2", "value2Changed"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "annotations.ann3", "newValue"),
				),
			},
			{
				Config: testAccCheckK8sNodePoolConfigUpdateAgain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(constant.ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "name", constant.K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "k8s_version", UpgradedK8sVersion),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week", "Tuesday"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "maintenance_window.0.time", "10:00:00Z"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "auto_scaling.#", "0"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "server_type", "DedicatedCore"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "node_count", "2"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cores_count", "2"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_size", "40"),
					// resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "public_ips"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "lans.#", "0"),
					// resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "labels"),
					// resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "annotations")
				),
			},
		},
	})
}

// func TestAccK8sNodePoolGatewayIP(t *testing.T) {
//	var k8sNodepool ionoscloud.KubernetesNodePool
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//		},
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
//		CheckDestroy:      testAccCheckK8sNodePoolDestroyCheck,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccCheckK8sNodePoolConfigGatewayIP,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckK8sNodePoolExists(ResourceNameK8sNodePool, &k8sNodepool),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "name", K8sNodePoolTestResource),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "k8s_version", K8sVersion),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week", "Monday"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "maintenance_window.0.time", "09:00:00Z"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "auto_scaling.0.min_node_count", "1"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "auto_scaling.0.max_node_count", "1"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cpu_family", "INTEL_XEON"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "availability_zone", "AUTO"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_type", "SSD"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "node_count", "1"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cores_count", "2"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "ram_size", "2048"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_size", "40"),
//					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "gateway_ip", IpBlockResource+".terraform_acctest", "ips.0"),
//					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "lans.0.id", LanResource+".terraform_acctest", "id"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.dhcp", "true"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.routes.0.network", "1.2.3.5/24"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.routes.0.gateway_ip", "10.1.5.17"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "labels.foo", "bar"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "labels.color", "green"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "annotations.ann1", "value1"),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "annotations.ann2", "value2"),
//				),
//			},
//		},
//	})
// }

func TestAccK8sNodePoolNoOptionalAndNodesDataSource(t *testing.T) {
	var k8sNodepool ionoscloud.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckK8sNodePoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sNodePoolConfigNoOptionalFields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(constant.ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "name", constant.K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "node_count", "2"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cores_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_size", "40"),
				),
			},
			{
				Config: testAccCheckK8sNodePoolConfigNoOptionalFieldsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(constant.ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "name", constant.K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "node_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cores_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_size", "40")),
			},
		},
	})
}

func TestAccK8sNodePoolCPUFamilyAndServerType(t *testing.T) {
	var k8sNodepool ionoscloud.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckK8sNodePoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sNodePoolConfigCPUFamilyAndServerType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(constant.ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "name", constant.K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "server_type", "DedicatedCore"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "node_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cores_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_size", "40"),
				),
			},
			{
				Config: testAccCheckK8sNodePoolConfigCPUFamilyAndServerTypeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(constant.ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "name", constant.K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "server_type", "VCPU"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "node_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cores_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_size", "40")),
			},
			{
				Config: testAccCheckK8sNodePoolConfigCPUFamilyRemoveServerTypeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(constant.ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "name", constant.K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "server_type", "DedicatedCore"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "node_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cores_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_size", "40")),
			},
		},
	})
}

func TestAccK8sNodePoolCPUFamilyNoServerType(t *testing.T) {
	var k8sNodepool ionoscloud.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckK8sNodePoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sNodePoolConfigCPUFamilyNoServerType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(constant.ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "name", constant.K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "node_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "cores_count", "1"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(constant.ResourceNameK8sNodePool, "storage_size", "40"),
				),
			},
		},
	})
}

func testAccCheckK8sNodePoolDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.K8sNodePoolResource {
			continue
		}

		_, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking the destruction of k8s node pool %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("k8s node pool %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckK8sNodePoolExists(n string, k8sNodepool *ionoscloud.KubernetesNodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		log.Printf("[INFO] REQ PATH: %+v/%+v", rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID)

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundK8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching k8s node pool: %s", rs.Primary.ID)
		}
		if *foundK8sNodepool.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		k8sNodepool = &foundK8sNodepool

		return nil
	}
}

const testAccCheckK8sNodePoolConfigBasic = `
resource ` + constant.DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}
resource ` + constant.LanResource + ` "terraform_acctest" {
  datacenter_id = ` + constant.DatacenterResource + `.terraform_acctest.id
  public = false
  name = "terraform_acctest"
}
resource ` + constant.IpBlockResource + ` "terraform_acctest" {
  location = ` + constant.DatacenterResource + `.terraform_acctest.location
  size = 3
  name = "terraform_acctest"
}
resource ` + constant.K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "` + K8sVersion + `"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}
resource ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolTestResource + ` {
  datacenter_id     = ` + constant.DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + constant.K8sClusterResource + `.terraform_acctest.id
  name        = "` + constant.K8sNodePoolTestResource + `"
  k8s_version = ` + constant.K8sClusterResource + `.terraform_acctest.k8s_version
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
  server_type       = "DedicatedCore"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
  lans {
    id   = ` + constant.LanResource + `.terraform_acctest.id
    dhcp = true
	routes {
       network   = "1.2.3.5/24"
       gateway_ip = "10.1.5.17"
     }
   }  
  labels = {
    foo = "bar"
    color = "green"
  }
  annotations = {
    ann1 = "value1"
    ann2 = "value2"
  }
}`

const testAccCheckK8sNodePoolConfigUpdate = `
resource ` + constant.DatacenterResource + ` "terraform_acctest" {
	name        = "terraform_acctest"
	location    = "us/las"
	description = "Datacenter created through terraform"
}
resource ` + constant.LanResource + ` "terraform_acctest" {
	datacenter_id = ` + constant.DatacenterResource + `.terraform_acctest.id
	public = false
	name = "terraform_acctest"
}
resource ` + constant.LanResource + ` "terraform_acctest_updated" {
	datacenter_id = ` + constant.DatacenterResource + `.terraform_acctest.id
	public = false
	name = "terraform_acctest"
}
resource ` + constant.IpBlockResource + ` "terraform_acctest" {
	location = ` + constant.DatacenterResource + `.terraform_acctest.location
	size = 3
	name = "terraform_acctest"
}
resource ` + constant.K8sClusterResource + ` "terraform_acctest" {
	name        = "terraform_acctest"
	k8s_version = "` + K8sVersion + `"
	maintenance_window {
		day_of_the_week = "Monday"
		time            = "09:00:00Z"
	}
}
resource ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolTestResource + ` {
  	datacenter_id     = ` + constant.DatacenterResource + `.terraform_acctest.id
  	k8s_cluster_id    = ` + constant.K8sClusterResource + `.terraform_acctest.id
  	name        = "` + constant.K8sNodePoolTestResource + `"
 	 k8s_version = ` + constant.K8sClusterResource + `.terraform_acctest.k8s_version
 	 auto_scaling {
 	 	min_node_count = 1
		max_node_count = 2
  }
  maintenance_window {
    day_of_the_week = "Tuesday"
    time            = "10:00:00Z"
  }
  server_type       = "VCPU"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 2
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
  public_ips        = [ ionoscloud_ipblock.terraform_acctest.ips[0], ionoscloud_ipblock.terraform_acctest.ips[1], ionoscloud_ipblock.terraform_acctest.ips[2]]
  lans {
    id   = ` + constant.LanResource + `.terraform_acctest.id
    dhcp = false
 	routes {
       network   = "1.2.3.4/24"
       gateway_ip = "10.1.5.16"
     }
  }
  lans {
    id   = ` + constant.LanResource + `.terraform_acctest_updated.id
    dhcp = false
 	routes {
       network   = "1.2.3.5/24"
       gateway_ip = "10.1.5.17"
     } 	
     routes {
       network   = "1.2.3.6/24"
       gateway_ip = "10.1.5.18"
     }
   }
  labels = {
    foo = "baz"
    color = "red"
    third = "thirdValue"
  }
  annotations = {
    ann1 = "value1Changed"
    ann2 = "value2Changed"
    ann3 = "newValue"
  }
}`

const testAccCheckK8sNodePoolConfigUpdateAgain = `
resource ` + constant.DatacenterResource + ` "terraform_acctest" {
	name        = "terraform_acctest"
	location    = "us/las"
	description = "Datacenter created through terraform"
}
resource ` + constant.LanResource + ` "terraform_acctest" {
	datacenter_id = ` + constant.DatacenterResource + `.terraform_acctest.id
	public = false
	name = "terraform_acctest"
}
resource ` + constant.LanResource + ` "terraform_acctest_updated" {
	datacenter_id = ` + constant.DatacenterResource + `.terraform_acctest.id
	public = false
	name = "terraform_acctest"
}
resource ` + constant.IpBlockResource + ` "terraform_acctest" {
	location = ` + constant.DatacenterResource + `.terraform_acctest.location
	size = 3
	name = "terraform_acctest"
}
resource ` + constant.K8sClusterResource + ` "terraform_acctest" {
	name        = "terraform_acctest"
    k8s_version = "` + UpgradedK8sVersion + `"
	maintenance_window {
		day_of_the_week = "Monday"
		time            = "09:00:00Z"
	}
}
resource ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolTestResource + ` {
  	datacenter_id     = ` + constant.DatacenterResource + `.terraform_acctest.id
  	k8s_cluster_id    = ` + constant.K8sClusterResource + `.terraform_acctest.id
  	name        = "` + constant.K8sNodePoolTestResource + `"
 	 k8s_version = ` + constant.K8sClusterResource + `.terraform_acctest.k8s_version
  maintenance_window {
    day_of_the_week = "Tuesday"
    time            = "10:00:00Z"
  }
  server_type       = "DedicatedCore"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 2
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
  public_ips        = [ ionoscloud_ipblock.terraform_acctest.ips[0], ionoscloud_ipblock.terraform_acctest.ips[1], ionoscloud_ipblock.terraform_acctest.ips[2]]
  labels = {
    foo = "baz"
    color = "red"
    third = "thirdValue"
  }
  annotations = {
    ann1 = "value1Changed"
    ann2 = "value2Changed"
    ann3 = "newValue"
  }
}`
const testAccCheckK8sNodePoolConfigGatewayIP = `
resource ` + constant.DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}
resource ` + constant.LanResource + ` "terraform_acctest" {
  datacenter_id = ` + constant.DatacenterResource + `.terraform_acctest.id
  public = false
  name = "terraform_acctest"
}
resource ` + constant.IpBlockResource + ` "terraform_acctest" {
  location = ` + constant.DatacenterResource + `.terraform_acctest.location
  size = 1
  name = "terraform_acctest"
}
resource ` + constant.K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = ` + K8sVersion + `
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
  //public = "false"
}

resource ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolTestResource + ` {
  datacenter_id     = ` + constant.DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + constant.K8sClusterResource + `.terraform_acctest.id
  name        = "` + constant.K8sNodePoolTestResource + `"
  k8s_version = ` + constant.K8sClusterResource + `.terraform_acctest.k8s_version
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  } 
  auto_scaling {
    min_node_count = 1
    max_node_count = 1
  }
  server_type       = "DedicatedCore"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
  //gateway_ip        = ` + constant.IpBlockResource + `.terraform_acctest.ips[0]
  lans {
    id   = ` + constant.LanResource + `.terraform_acctest.id
    dhcp = true
	routes {
       network   = "1.2.3.5/24"
       gateway_ip = "10.1.5.17"
     }
   }  
  labels = {
    foo = "bar"
    color = "green"
  }
  annotations = {
    ann1 = "value1"
    ann2 = "value2"
  }
}`

const testAccDataSourceProfitBricksK8sNodePoolMatchId = testAccCheckK8sNodePoolConfigBasic + `
data ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolDataSourceById + ` {
	k8s_cluster_id  = ` + constant.K8sClusterResource + `.terraform_acctest.id
	id				= ` + constant.K8sNodePoolResource + `.` + constant.K8sNodePoolTestResource + `.id
}
`

const testAccDataSourceProfitBricksK8sNodePoolNodesMatchId = testAccCheckK8sNodePoolConfigBasic + `
data ` + constant.K8sNodePoolNodesResource + ` ` + constant.K8sNodePoolDataSourceById + ` {
	k8s_cluster_id  = ` + constant.K8sClusterResource + `.terraform_acctest.id
	node_pool_id	= ` + constant.K8sNodePoolResource + `.` + constant.K8sNodePoolTestResource + `.id
}
`

const testAccDataSourceProfitBricksK8sNodePoolMatchName = testAccCheckK8sNodePoolConfigBasic + `
data ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolDataSourceByName + ` {
	k8s_cluster_id 	= ` + constant.K8sClusterResource + `.terraform_acctest.id
	name			= ` + constant.K8sNodePoolResource + `.` + constant.K8sNodePoolTestResource + `.name
}
`

const testAccDataSourceProfitBricksK8sNodePoolWrongNameError = testAccCheckK8sNodePoolConfigBasic + `
data ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolDataSourceByName + ` {
	k8s_cluster_id 	= ` + constant.K8sClusterResource + `.terraform_acctest.id
	name			= "wrong_name"
}
`

const testAccCheckK8sNodePoolConfigNoOptionalFields = `
resource ` + constant.DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource ` + constant.K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "` + K8sVersion + `"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}
resource ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolTestResource + ` {
  datacenter_id     = ` + constant.DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + constant.K8sClusterResource + `.terraform_acctest.id
  k8s_version = ` + constant.K8sClusterResource + `.terraform_acctest.k8s_version
  name        = "` + constant.K8sNodePoolTestResource + `"
  auto_scaling {
    min_node_count = 1
    max_node_count = 3
  }
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 2
  cores_count       = 1
  ram_size          = 2048
  storage_size      = 40
}`

const testAccCheckK8sNodePoolConfigNoOptionalFieldsUpdate = `
resource ` + constant.DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource ` + constant.K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "` + K8sVersion + `"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}
resource ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolTestResource + ` {
  datacenter_id     = ` + constant.DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + constant.K8sClusterResource + `.terraform_acctest.id
  name        = "` + constant.K8sNodePoolTestResource + `"
  k8s_version = ` + constant.K8sClusterResource + `.terraform_acctest.k8s_version
  auto_scaling {
    min_node_count = 1
    max_node_count = 3
  }
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 1
  ram_size          = 2048
  storage_size      = 40
}
data ` + constant.K8sNodePoolNodesResource + ` nodes{
  k8s_cluster_id   = ` + constant.K8sClusterResource + `.terraform_acctest.id
  node_pool_id     = ` + constant.K8sNodePoolResource + `.` + constant.K8sNodePoolTestResource + `.id
}`

const testAccCheckK8sNodePoolConfigCPUFamilyAndServerType = `
resource ` + constant.DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource ` + constant.K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "` + K8sVersion + `"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}
resource ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolTestResource + ` {
  datacenter_id     = ` + constant.DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + constant.K8sClusterResource + `.terraform_acctest.id
  k8s_version = ` + constant.K8sClusterResource + `.terraform_acctest.k8s_version
  name        = "` + constant.K8sNodePoolTestResource + `"
  auto_scaling {
    min_node_count = 1
    max_node_count = 3
  }
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  server_type       = "DedicatedCore"
  cpu_family        = "INTEL_XEON"
  cores_count       = 1
  ram_size          = 2048
  storage_size      = 40
}`

const testAccCheckK8sNodePoolConfigCPUFamilyAndServerTypeUpdate = `
resource ` + constant.DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource ` + constant.K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "` + K8sVersion + `"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}
resource ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolTestResource + ` {
  datacenter_id     = ` + constant.DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + constant.K8sClusterResource + `.terraform_acctest.id
  name        = "` + constant.K8sNodePoolTestResource + `"
  k8s_version = ` + constant.K8sClusterResource + `.terraform_acctest.k8s_version
  auto_scaling {
    min_node_count = 1
    max_node_count = 3
  }
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  server_type       = "VCPU"
  cpu_family        = "INTEL_XEON"
  cores_count       = 1
  ram_size          = 2048
  storage_size      = 40
}
data ` + constant.K8sNodePoolNodesResource + ` nodes{
  k8s_cluster_id   = ` + constant.K8sClusterResource + `.terraform_acctest.id
  node_pool_id     = ` + constant.K8sNodePoolResource + `.` + constant.K8sNodePoolTestResource + `.id
}`

const testAccCheckK8sNodePoolConfigCPUFamilyRemoveServerTypeUpdate = `
resource ` + constant.DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource ` + constant.K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "` + K8sVersion + `"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}
resource ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolTestResource + ` {
  datacenter_id     = ` + constant.DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + constant.K8sClusterResource + `.terraform_acctest.id
  name        = "` + constant.K8sNodePoolTestResource + `"
  k8s_version = ` + constant.K8sClusterResource + `.terraform_acctest.k8s_version
  auto_scaling {
    min_node_count = 1
    max_node_count = 3
  }
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cpu_family        = "INTEL_XEON"
  cores_count       = 1
  ram_size          = 2048
  storage_size      = 40
}
data ` + constant.K8sNodePoolNodesResource + ` nodes{
  k8s_cluster_id   = ` + constant.K8sClusterResource + `.terraform_acctest.id
  node_pool_id     = ` + constant.K8sNodePoolResource + `.` + constant.K8sNodePoolTestResource + `.id
}`

const testAccCheckK8sNodePoolConfigCPUFamilyNoServerType = `
resource ` + constant.DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource ` + constant.K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "` + K8sVersion + `"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}
resource ` + constant.K8sNodePoolResource + ` ` + constant.K8sNodePoolTestResource + ` {
  datacenter_id     = ` + constant.DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + constant.K8sClusterResource + `.terraform_acctest.id
  name        = "` + constant.K8sNodePoolTestResource + `"
  k8s_version = ` + constant.K8sClusterResource + `.terraform_acctest.k8s_version
  auto_scaling {
    min_node_count = 1
    max_node_count = 3
  }
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cpu_family        = "INTEL_XEON"
  cores_count       = 1
  ram_size          = 2048
  storage_size      = 40
}
data ` + constant.K8sNodePoolNodesResource + ` nodes{
  k8s_cluster_id   = ` + constant.K8sClusterResource + `.terraform_acctest.id
  node_pool_id     = ` + constant.K8sNodePoolResource + `.` + constant.K8sNodePoolTestResource + `.id
}`
