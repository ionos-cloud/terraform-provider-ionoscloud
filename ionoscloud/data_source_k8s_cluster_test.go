package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceK8sCluster(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sClusterConfigBasic,
			},
			{
				Config: testAccDataSourceK8sClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceById, "name", K8sClusterResource+"."+K8sClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceById, "k8s_version", K8sClusterResource+"."+K8sClusterTestResource, "k8s_version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceById, "maintenance_window.0.day_of_the_week", K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceById, "maintenance_window.0.time", K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceById, "maintenance_window.0.time", K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceById, "api_subnet_allow_list.0", K8sClusterResource+"."+K8sClusterTestResource, "api_subnet_allow_list.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceById, "s3_buckets.0.name", K8sClusterResource+"."+K8sClusterTestResource, "s3_buckets.0.name"),
				),
			},
			{
				Config: testAccDataSourceK8sClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceByName, "name", K8sClusterResource+"."+K8sClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceByName, "k8s_version", K8sClusterResource+"."+K8sClusterTestResource, "k8s_version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceByName, "maintenance_window.0.day_of_the_week", K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceByName, "maintenance_window.0.time", K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceByName, "maintenance_window.0.time", K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceByName, "api_subnet_allow_list.0", K8sClusterResource+"."+K8sClusterTestResource, "api_subnet_allow_list.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+K8sClusterResource+"."+K8sClusterDataSourceByName, "s3_buckets.0.name", K8sClusterResource+"."+K8sClusterTestResource, "s3_buckets.0.name"),
				),
			},
		},
	})
}

const testAccDataSourceK8sClusterMatchId = testAccCheckK8sClusterConfigBasic + `
data ` + K8sClusterResource + ` ` + K8sClusterDataSourceById + `{
  id	= ` + K8sClusterResource + `.` + K8sClusterTestResource + `.id
}
`

const testAccDataSourceK8sClusterMatchName = testAccCheckK8sClusterConfigBasic + `
data ` + K8sClusterResource + ` ` + K8sClusterDataSourceByName + `{
  name	= "` + K8sClusterTestResource + `"
}
`
