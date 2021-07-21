// +build k8s

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceK8sNodePool_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_k8s_node_pool.test_ds_k8s_node_pool", "name", "test_nodepool"),
					resource.TestCheckResourceAttr("data.ionoscloud_k8s_node_pool.test_ds_k8s_node_pool", "k8s_version", "1.20.8"),
				),
			},
		},
	})
}

func TestAccDataSourceK8sNodePool_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_k8s_node_pool.test_ds_k8s_node_pool", "name", "test_nodepool"),
					resource.TestCheckResourceAttr("data.ionoscloud_k8s_node_pool.test_ds_k8s_node_pool", "k8s_version", "1.20.8"),
					resource.TestCheckResourceAttrSet("data.ionoscloud_k8s_node_pool.test_ds_k8s_node_pool", "id"),
				),
			},
		},
	})

}

const testAccDataSourceProfitBricksK8sNodePoolCreateResources = `
resource "ionoscloud_datacenter" "test_ds_k8s_datacenter" {
	name              = "test_datacenter"
	location          = "us/las"
	description       = "test datacenter"
}

resource "ionoscloud_k8s_cluster" "test_ds_k8s_cluster" {
	name              = "test_cluster"
}

resource "ionoscloud_k8s_node_pool" "test_ds_k8s_node_pool" {
	depends_on 				= [ionoscloud_datacenter.test_ds_k8s_datacenter, ionoscloud_k8s_cluster.test_ds_k8s_cluster]
	name					= "test_nodepool"
	datacenter_id			= ionoscloud_datacenter.test_ds_k8s_datacenter.id
	k8s_cluster_id			= ionoscloud_k8s_cluster.test_ds_k8s_cluster.id
	node_count				= 1
	cpu_family				= "AMD_OPTERON"
	cores_count				= 1
	ram_size				= 2048
	availability_zone 		= "AUTO"
	storage_type			= "HDD"
	storage_size			= 15
	k8s_version				= "1.20.8"
}
`

const testAccDataSourceProfitBricksK8sNodePoolMatchId = `
resource "ionoscloud_datacenter" "test_ds_k8s_datacenter" {
  name              = "test_datacenter"
  location          = "us/las"
  description       = "test datacenter"
}

resource "ionoscloud_k8s_cluster" "test_ds_k8s_cluster" {
  name              = "test_cluster"
}

resource "ionoscloud_k8s_node_pool" "test_ds_k8s_node_pool" {
	depends_on 				= [ionoscloud_datacenter.test_ds_k8s_datacenter, ionoscloud_k8s_cluster.test_ds_k8s_cluster]
  name							= "test_nodepool"
	datacenter_id			= ionoscloud_datacenter.test_ds_k8s_datacenter.id
	k8s_cluster_id		= ionoscloud_k8s_cluster.test_ds_k8s_cluster.id
	node_count				= 1
	cpu_family				= "AMD_OPTERON"
	cores_count				= 1
	ram_size					= 2048
	availability_zone = "AUTO"
	storage_type			= "HDD"
	storage_size			= 15
	k8s_version				= "1.20.8"
  #	public_ips				= [ "158.222.102.239", "158.222.102.241", "158.222.102.242" ]
  # public_ips				= [ ]
  #   public_ips        = [ ]
}

data "ionoscloud_k8s_node_pool" "test_ds_k8s_node_pool" {
	k8s_cluster_id 	= ionoscloud_k8s_cluster.test_ds_k8s_cluster.id
	id				= ionoscloud_k8s_node_pool.test_ds_k8s_node_pool.id
}
`

const testAccDataSourceProfitBricksK8sNodePoolMatchName = `
resource "ionoscloud_datacenter" "test_ds_k8s_datacenter" {
  name              = "test_datacenter"
  location          = "us/las"
  description       = "test datacenter"
}

resource "ionoscloud_k8s_cluster" "test_ds_k8s_cluster" {
  name              = "test_cluster"
}

resource "ionoscloud_k8s_node_pool" "test_ds_k8s_node_pool" {
	depends_on 				= [ionoscloud_datacenter.test_ds_k8s_datacenter, ionoscloud_k8s_cluster.test_ds_k8s_cluster]
  name							= "test_nodepool"
	datacenter_id			= ionoscloud_datacenter.test_ds_k8s_datacenter.id
	k8s_cluster_id		= ionoscloud_k8s_cluster.test_ds_k8s_cluster.id
	node_count				= 1
	cpu_family				= "AMD_OPTERON"
	cores_count				= 1
	ram_size					= 2048
	availability_zone = "AUTO"
	storage_type			= "HDD"
	storage_size			= 15
	k8s_version				= "1.20.8"
  #	public_ips				= [ "158.222.102.239", "158.222.102.241", "158.222.102.242" ]
  # public_ips				= [ ]
  #   public_ips        = [ ]
}

data "ionoscloud_k8s_node_pool" "test_ds_k8s_node_pool" {
	k8s_cluster_id 	= ionoscloud_k8s_cluster.test_ds_k8s_cluster.id
	name			= "test_nodepool"
}
`
