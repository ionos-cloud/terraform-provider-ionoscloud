package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const dataSourceK8sNodePoolId = DataSource + "." + K8sNodePoolResource + "." + K8sNodePoolDataSourceById
const dataSourceK8sNodePoolName = DataSource + "." + K8sNodePoolResource + "." + K8sNodePoolDataSourceByName

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
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "name", resourceNameK8sNodePool, "name"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "k8s_version", resourceNameK8sNodePool, "k8s_version"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "maintenance_window.0.day_of_the_week", resourceNameK8sNodePool, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "maintenance_window.0.time", resourceNameK8sNodePool, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "auto_scaling.0.min_node_count", resourceNameK8sNodePool, "auto_scaling.0.min_node_count"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "auto_scaling.0.max_node_count", resourceNameK8sNodePool, "auto_scaling.0.max_node_count"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "cpu_family", resourceNameK8sNodePool, "cpu_family"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "availability_zone", resourceNameK8sNodePool, "availability_zone"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "storage_type", resourceNameK8sNodePool, "storage_type"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "node_count", resourceNameK8sNodePool, "node_count"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "cores_count", resourceNameK8sNodePool, "cores_count"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "ram_size", resourceNameK8sNodePool, "ram_size"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "storage_size", resourceNameK8sNodePool, "storage_size"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "public_ips.0", resourceNameK8sNodePool, "public_ips.0"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "public_ips.1", resourceNameK8sNodePool, "public_ips.1"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "lans.0", resourceNameK8sNodePool, "lans.0"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "labels.foo", resourceNameK8sNodePool, "labels.foo"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "labels.color", resourceNameK8sNodePool, "labels.color"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "annotations.ann1", resourceNameK8sNodePool, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolId, "annotations.ann2", resourceNameK8sNodePool, "annotations.ann2"),
				),
			},
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "name", resourceNameK8sNodePool, "name"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "k8s_version", resourceNameK8sNodePool, "k8s_version"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "maintenance_window.0.day_of_the_week", resourceNameK8sNodePool, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "maintenance_window.0.time", resourceNameK8sNodePool, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "auto_scaling.0.min_node_count", resourceNameK8sNodePool, "auto_scaling.0.min_node_count"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "auto_scaling.0.max_node_count", resourceNameK8sNodePool, "auto_scaling.0.max_node_count"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "cpu_family", resourceNameK8sNodePool, "cpu_family"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "availability_zone", resourceNameK8sNodePool, "availability_zone"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "storage_type", resourceNameK8sNodePool, "storage_type"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "node_count", resourceNameK8sNodePool, "node_count"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "cores_count", resourceNameK8sNodePool, "cores_count"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "ram_size", resourceNameK8sNodePool, "ram_size"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "storage_size", resourceNameK8sNodePool, "storage_size"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "public_ips.0", resourceNameK8sNodePool, "public_ips.0"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "public_ips.1", resourceNameK8sNodePool, "public_ips.1"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "lans.0", resourceNameK8sNodePool, "lans.0"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "labels.foo", resourceNameK8sNodePool, "labels.foo"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "labels.color", resourceNameK8sNodePool, "labels.color"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "annotations.ann1", resourceNameK8sNodePool, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(dataSourceK8sNodePoolName, "annotations.ann2", resourceNameK8sNodePool, "annotations.ann2"),
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
