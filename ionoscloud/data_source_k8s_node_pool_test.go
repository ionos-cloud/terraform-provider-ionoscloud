// +build k8s

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceK8sNodePool(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sNodePoolConfigBasic,
			},
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "name", ResourceNameK8sNodePool, "name"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "k8s_version", ResourceNameK8sNodePool, "k8s_version"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "maintenance_window.0.day_of_the_week", ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "maintenance_window.0.time", ResourceNameK8sNodePool, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "auto_scaling.0.min_node_count", ResourceNameK8sNodePool, "auto_scaling.0.min_node_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "auto_scaling.0.max_node_count", ResourceNameK8sNodePool, "auto_scaling.0.max_node_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "cpu_family", ResourceNameK8sNodePool, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "availability_zone", ResourceNameK8sNodePool, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "storage_type", ResourceNameK8sNodePool, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "node_count", ResourceNameK8sNodePool, "node_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "cores_count", ResourceNameK8sNodePool, "cores_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "ram_size", ResourceNameK8sNodePool, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "storage_size", ResourceNameK8sNodePool, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "public_ips.0", ResourceNameK8sNodePool, "public_ips.0"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "public_ips.1", ResourceNameK8sNodePool, "public_ips.1"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "lans.0", ResourceNameK8sNodePool, "lans.0"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "labels.foo", ResourceNameK8sNodePool, "labels.foo"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "labels.color", ResourceNameK8sNodePool, "labels.color"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "annotations.ann1", ResourceNameK8sNodePool, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "annotations.ann2", ResourceNameK8sNodePool, "annotations.ann2"),
				),
			},
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "name", ResourceNameK8sNodePool, "name"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "k8s_version", ResourceNameK8sNodePool, "k8s_version"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "maintenance_window.0.day_of_the_week", ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "maintenance_window.0.time", ResourceNameK8sNodePool, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "auto_scaling.0.min_node_count", ResourceNameK8sNodePool, "auto_scaling.0.min_node_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "auto_scaling.0.max_node_count", ResourceNameK8sNodePool, "auto_scaling.0.max_node_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "cpu_family", ResourceNameK8sNodePool, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "availability_zone", ResourceNameK8sNodePool, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "storage_type", ResourceNameK8sNodePool, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "node_count", ResourceNameK8sNodePool, "node_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "cores_count", ResourceNameK8sNodePool, "cores_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "ram_size", ResourceNameK8sNodePool, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "storage_size", ResourceNameK8sNodePool, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "public_ips.0", ResourceNameK8sNodePool, "public_ips.0"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "public_ips.1", ResourceNameK8sNodePool, "public_ips.1"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "lans.0", ResourceNameK8sNodePool, "lans.0"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "labels.foo", ResourceNameK8sNodePool, "labels.foo"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "labels.color", ResourceNameK8sNodePool, "labels.color"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "annotations.ann1", ResourceNameK8sNodePool, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "annotations.ann2", ResourceNameK8sNodePool, "annotations.ann2"),
				),
			},
		},
	})
}

const testAccDataSourceProfitBricksK8sNodePoolMatchId = testAccCheckK8sNodePoolConfigBasic + `
data ` + K8sNodePoolResource + ` ` + K8sNodePoolDataSourceById + ` {
	k8s_cluster_id  = ` + K8sClusterResource + `.terraform_acctest.id
	id				= ` + K8sNodePoolResource + `.` + K8sNodePoolTestResource + `.id
}
`

const testAccDataSourceProfitBricksK8sNodePoolMatchName = testAccCheckK8sNodePoolConfigBasic + `
data ` + K8sNodePoolResource + ` ` + K8sNodePoolDataSourceByName + ` {
	k8s_cluster_id 	= ` + K8sClusterResource + `.terraform_acctest.id
	name			= ` + K8sNodePoolResource + `.` + K8sNodePoolTestResource + `.name
}
`
