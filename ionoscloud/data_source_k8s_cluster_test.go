// +build k8s

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceK8sCluster_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksK8sClusterCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksK8sClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_k8s_cluster.test_ds_k8s_cluster", "name", "test_cluster"),
					resource.TestCheckResourceAttr("data.ionoscloud_k8s_cluster.test_ds_k8s_cluster", "k8s_version", "1.20.8"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_k8s_cluster.test_ds_k8s_cluster", "kube_config"),
				),
			},
		},
	})
}

func TestAccDataSourceK8sCluster_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksK8sClusterCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksK8sClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_k8s_cluster.test_ds_k8s_cluster", "name", "test_cluster"),
					resource.TestCheckResourceAttr("data.ionoscloud_k8s_cluster.test_ds_k8s_cluster", "k8s_version", "1.20.8"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_k8s_cluster.test_ds_k8s_cluster", "kube_config"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_k8s_cluster.test_ds_k8s_cluster", "id"),
				),
			},
		},
	})

}

const testAccDataSourceProfitBricksK8sClusterCreateResources = `
resource "ionoscloud_k8s_cluster" "test_ds_k8s_cluster" {
  name         = "test_cluster"
  k8s_version  = "1.20.8"
}
`

const testAccDataSourceProfitBricksK8sClusterMatchId = `
resource "ionoscloud_k8s_cluster" "test_ds_k8s_cluster" {
  name         = "test_cluster"
  k8s_version  = "1.20.8"
}

data "ionoscloud_k8s_cluster" "test_ds_k8s_cluster" {
  id	= ionoscloud_k8s_cluster.test_ds_k8s_cluster.id
}
`

const testAccDataSourceProfitBricksK8sClusterMatchName = `
resource "ionoscloud_k8s_cluster" "test_ds_k8s_cluster" {
  name         = "test_cluster"
  k8s_version  = "1.20.8"
}

data "ionoscloud_k8s_cluster" "test_ds_k8s_cluster" {
  name	= "test_cluster"
}
`
