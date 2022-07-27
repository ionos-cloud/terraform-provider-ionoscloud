//go:build all || k8s
// +build all k8s

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccK8sNodePoolBasic(t *testing.T) {
	var k8sNodepool ionoscloud.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckK8sNodePoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sNodePoolConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "name", K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "k8s_version", "1.20.10"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week", "Monday"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "maintenance_window.0.time", "09:00:00Z"),
					resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "auto_scaling"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "node_count", "1"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cores_count", "2"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_size", "40"),
					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "lans.0.id", LanResource+".terraform_acctest", "id"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.dhcp", "true"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.routes.0.network", "1.2.3.5/24"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.routes.0.gateway_ip", "10.1.5.17"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "labels.foo", "bar"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "labels.color", "green"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "annotations.ann1", "value1"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "annotations.ann2", "value2"),
				),
			},
			{
				Config: testAccDataSourceK8sNodePoolMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "name", ResourceNameK8sNodePool, "name"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "k8s_version", ResourceNameK8sNodePool, "k8s_version"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "maintenance_window.0.day_of_the_week", ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "maintenance_window.0.time", ResourceNameK8sNodePool, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "cpu_family", ResourceNameK8sNodePool, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "availability_zone", ResourceNameK8sNodePool, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "storage_type", ResourceNameK8sNodePool, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "node_count", ResourceNameK8sNodePool, "node_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "cores_count", ResourceNameK8sNodePool, "cores_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "ram_size", ResourceNameK8sNodePool, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "storage_size", ResourceNameK8sNodePool, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "lans.0", ResourceNameK8sNodePool, "lans.0"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "labels.foo", ResourceNameK8sNodePool, "labels.foo"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "labels.color", ResourceNameK8sNodePool, "labels.color"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "annotations.ann1", ResourceNameK8sNodePool, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolId, "annotations.ann2", ResourceNameK8sNodePool, "annotations.ann2"),
				),
			},
			{
				Config: testAccDataSourceK8sNodePoolMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "name", ResourceNameK8sNodePool, "name"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "k8s_version", ResourceNameK8sNodePool, "k8s_version"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "maintenance_window.0.day_of_the_week", ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "maintenance_window.0.time", ResourceNameK8sNodePool, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "cpu_family", ResourceNameK8sNodePool, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "availability_zone", ResourceNameK8sNodePool, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "storage_type", ResourceNameK8sNodePool, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "node_count", ResourceNameK8sNodePool, "node_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "cores_count", ResourceNameK8sNodePool, "cores_count"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "ram_size", ResourceNameK8sNodePool, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "storage_size", ResourceNameK8sNodePool, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "lans.0", ResourceNameK8sNodePool, "lans.0"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "labels.foo", ResourceNameK8sNodePool, "labels.foo"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "labels.color", ResourceNameK8sNodePool, "labels.color"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "annotations.ann1", ResourceNameK8sNodePool, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(DataSourceK8sNodePoolName, "annotations.ann2", ResourceNameK8sNodePool, "annotations.ann2"),
				),
			},
			{
				Config: testAccDataSourceK8sNodePoolPartialMatchName,
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
			{
				Config:      testAccDataSourceK8sNodePoolWrongNameError,
				ExpectError: regexp.MustCompile("no nodepool found with the specified name"),
			},
			{
				Config:      testAccDataSourceK8sNodePoolWrongPartialNameError,
				ExpectError: regexp.MustCompile("no nodepool found with the specified name"),
			},
			{
				Config: testAccCheckK8sNodePoolConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "name", K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "k8s_version", "1.20.10"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week", "Tuesday"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "maintenance_window.0.time", "10:00:00Z"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "auto_scaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "auto_scaling.0.max_node_count", "2"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "node_count", "2"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cores_count", "2"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_size", "40"),
					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "public_ips.0", IpBlockResource+".terraform_acctest", "ips.0"),
					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "public_ips.1", IpBlockResource+".terraform_acctest", "ips.1"),
					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "lans.0.id", LanResource+".terraform_acctest", "id"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.dhcp", "false"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.routes.0.network", "1.2.3.4/24"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.routes.0.gateway_ip", "10.1.5.16"),
					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "lans.1.id", LanResource+".terraform_acctest_updated", "id"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.1.dhcp", "false"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.1.routes.0.network", "1.2.3.5/24"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.1.routes.0.gateway_ip", "10.1.5.17"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.1.routes.1.network", "1.2.3.6/24"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.1.routes.1.gateway_ip", "10.1.5.18"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "labels.foo", "baz"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "labels.color", "red"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "labels.third", "thirdValue"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "annotations.ann1", "value1Changed"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "annotations.ann2", "value2Changed"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "annotations.ann3", "newValue"),
				),
			},
			{
				Config: testAccCheckK8sNodePoolConfigUpdateAgain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "name", K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "k8s_version", "1.21.9"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week", "Tuesday"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "maintenance_window.0.time", "10:00:00Z"),
					resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "auto_scaling"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "node_count", "2"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cores_count", "2"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_size", "40"),
					//resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "public_ips"),
					resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "lans"),
					//resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "labels"),
					//resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "annotations")
				),
			},
		},
	})
}

//func TestAccK8sNodePoolGatewayIP(t *testing.T) {
//	var k8sNodepool ionoscloud.KubernetesNodePool
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//		},
//		ProviderFactories: testAccProviderFactories,
//		CheckDestroy:      testAccCheckK8sNodePoolDestroyCheck,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccCheckK8sNodePoolConfigGatewayIP,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckK8sNodePoolExists(ResourceNameK8sNodePool, &k8sNodepool),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "name", K8sNodePoolTestResource),
//					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "k8s_version", "1.20.10"),
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
//}

func TestAccK8sNodePoolNoOptionalAndNodesDataSource(t *testing.T) {
	var k8sNodepool ionoscloud.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckK8sNodePoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sNodePoolConfigNoOptionalFields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "name", K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "k8s_version", "1.23.6"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "node_count", "2"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cores_count", "1"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_size", "40"),
				),
			},
			{
				Config: testAccCheckK8sNodePoolConfigNoOptionalFieldsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "name", K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "k8s_version", "1.23.6"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "node_count", "1"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cores_count", "1"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_size", "40"),
					resource.TestCheckResourceAttr(DataSource+"."+K8sNodePoolNodesResource+".nodes", "nodes.#", "1")),
			},
		},
	})
}

func testAccCheckK8sNodePoolDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != K8sNodePoolResource {
			continue
		}

		_, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking the destruction of k8s node pool %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("k8s node pool %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckK8sNodePoolExists(n string, k8sNodepool *ionoscloud.KubernetesNodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

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
			return fmt.Errorf("error occured while fetching k8s node pool: %s", rs.Primary.ID)
		}
		if *foundK8sNodepool.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		k8sNodepool = &foundK8sNodepool

		return nil
	}
}

const testAccCheckK8sNodePoolConfigBasic = `
resource ` + DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}
resource ` + LanResource + ` "terraform_acctest" {
  datacenter_id = ` + DatacenterResource + `.terraform_acctest.id
  public = false
  name = "terraform_acctest"
}
resource ` + IpBlockResource + ` "terraform_acctest" {
  location = ` + DatacenterResource + `.terraform_acctest.location
  size = 3
  name = "terraform_acctest"
}
resource ` + K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "1.20.10"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}
resource ` + K8sNodePoolResource + ` ` + K8sNodePoolTestResource + ` {
  datacenter_id     = ` + DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + K8sClusterResource + `.terraform_acctest.id
  name        = "` + K8sNodePoolTestResource + `"
  k8s_version = ` + K8sClusterResource + `.terraform_acctest.k8s_version
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
  lans {
    id   = ` + LanResource + `.terraform_acctest.id
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
resource ` + DatacenterResource + ` "terraform_acctest" {
	name        = "terraform_acctest"
	location    = "us/las"
	description = "Datacenter created through terraform"
}
resource ` + LanResource + ` "terraform_acctest" {
	datacenter_id = ` + DatacenterResource + `.terraform_acctest.id
	public = false
	name = "terraform_acctest"
}
resource ` + LanResource + ` "terraform_acctest_updated" {
	datacenter_id = ` + DatacenterResource + `.terraform_acctest.id
	public = false
	name = "terraform_acctest"
}
resource ` + IpBlockResource + ` "terraform_acctest" {
	location = ` + DatacenterResource + `.terraform_acctest.location
	size = 3
	name = "terraform_acctest"
}
resource ` + K8sClusterResource + ` "terraform_acctest" {
	name        = "terraform_acctest"
	k8s_version = "1.20.14"
	maintenance_window {
		day_of_the_week = "Monday"
		time            = "09:00:00Z"
	}
}
resource ` + K8sNodePoolResource + ` ` + K8sNodePoolTestResource + ` {
  	datacenter_id     = ` + DatacenterResource + `.terraform_acctest.id
  	k8s_cluster_id    = ` + K8sClusterResource + `.terraform_acctest.id
  	name        = "` + K8sNodePoolTestResource + `"
 	 k8s_version = ` + K8sClusterResource + `.terraform_acctest.k8s_version
 	 auto_scaling {
 	 	min_node_count = 1
		max_node_count = 2
  }
  maintenance_window {
    day_of_the_week = "Tuesday"
    time            = "10:00:00Z"
  }
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 2
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
  public_ips        = [ ionoscloud_ipblock.terraform_acctest.ips[0], ionoscloud_ipblock.terraform_acctest.ips[1], ionoscloud_ipblock.terraform_acctest.ips[2]]
  lans {
    id   = ` + LanResource + `.terraform_acctest.id
    dhcp = false
 	routes {
       network   = "1.2.3.4/24"
       gateway_ip = "10.1.5.16"
     }
  }
  lans {
    id   = ` + LanResource + `.terraform_acctest_updated.id
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
resource ` + DatacenterResource + ` "terraform_acctest" {
	name        = "terraform_acctest"
	location    = "us/las"
	description = "Datacenter created through terraform"
}
resource ` + LanResource + ` "terraform_acctest" {
	datacenter_id = ` + DatacenterResource + `.terraform_acctest.id
	public = false
	name = "terraform_acctest"
}
resource ` + LanResource + ` "terraform_acctest_updated" {
	datacenter_id = ` + DatacenterResource + `.terraform_acctest.id
	public = false
	name = "terraform_acctest"
}
resource ` + IpBlockResource + ` "terraform_acctest" {
	location = ` + DatacenterResource + `.terraform_acctest.location
	size = 3
	name = "terraform_acctest"
}
resource ` + K8sClusterResource + ` "terraform_acctest" {
	name        = "terraform_acctest"
    k8s_version = "1.21.9"
	maintenance_window {
		day_of_the_week = "Monday"
		time            = "09:00:00Z"
	}
}
resource ` + K8sNodePoolResource + ` ` + K8sNodePoolTestResource + ` {
  	datacenter_id     = ` + DatacenterResource + `.terraform_acctest.id
  	k8s_cluster_id    = ` + K8sClusterResource + `.terraform_acctest.id
  	name        = "` + K8sNodePoolTestResource + `"
 	 k8s_version = ` + K8sClusterResource + `.terraform_acctest.k8s_version
  maintenance_window {
    day_of_the_week = "Tuesday"
    time            = "10:00:00Z"
  }
  cpu_family        = "INTEL_XEON"
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
resource ` + DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}
resource ` + LanResource + ` "terraform_acctest" {
  datacenter_id = ` + DatacenterResource + `.terraform_acctest.id
  public = false
  name = "terraform_acctest"
}
resource ` + IpBlockResource + ` "terraform_acctest" {
  location = ` + DatacenterResource + `.terraform_acctest.location
  size = 1
  name = "terraform_acctest"
}
resource ` + K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "1.20.10"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
  //public = "false"
}

resource ` + K8sNodePoolResource + ` ` + K8sNodePoolTestResource + ` {
  datacenter_id     = ` + DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + K8sClusterResource + `.terraform_acctest.id
  name        = "` + K8sNodePoolTestResource + `"
  k8s_version = ` + K8sClusterResource + `.terraform_acctest.k8s_version
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  } 
  auto_scaling {
    min_node_count = 1
    max_node_count = 1
  }
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
  //gateway_ip        = ` + IpBlockResource + `.terraform_acctest.ips[0]
  lans {
    id   = ` + LanResource + `.terraform_acctest.id
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

const testAccDataSourceK8sNodePoolMatchId = testAccCheckK8sNodePoolConfigBasic + `
data ` + K8sNodePoolResource + ` ` + K8sNodePoolDataSourceById + ` {
	k8s_cluster_id  = ` + K8sClusterResource + `.terraform_acctest.id
	id				= ` + K8sNodePoolResource + `.` + K8sNodePoolTestResource + `.id
}
`

const testAccDataSourceK8sNodePoolMatchName = testAccCheckK8sNodePoolConfigBasic + `
data ` + K8sNodePoolResource + ` ` + K8sNodePoolDataSourceByName + ` {
	k8s_cluster_id 	= ` + K8sClusterResource + `.terraform_acctest.id
	name			= ` + K8sNodePoolResource + `.` + K8sNodePoolTestResource + `.name
}
`

const testAccDataSourceK8sNodePoolPartialMatchName = testAccCheckK8sNodePoolConfigBasic + `
data ` + K8sNodePoolResource + ` ` + K8sNodePoolDataSourceByName + ` {
	k8s_cluster_id 	= ` + K8sClusterResource + `.terraform_acctest.id
	name			= "` + DataSourcePartial + `"
    partial_match   = true
}
`

const testAccDataSourceK8sNodePoolWrongNameError = testAccCheckK8sNodePoolConfigBasic + `
data ` + K8sNodePoolResource + ` ` + K8sNodePoolDataSourceByName + ` {
	k8s_cluster_id 	= ` + K8sClusterResource + `.terraform_acctest.id
	name			= "wrong_name"
}
`

const testAccDataSourceK8sNodePoolWrongPartialNameError = testAccCheckK8sNodePoolConfigBasic + `
data ` + K8sNodePoolResource + ` ` + K8sNodePoolDataSourceByName + ` {
	k8s_cluster_id 	= ` + K8sClusterResource + `.terraform_acctest.id
	name			= "wrong_name"
    partial_match   = true
}
`

const testAccCheckK8sNodePoolConfigNoOptionalFields = `
resource ` + DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource ` + K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "1.23.6"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}
resource ` + K8sNodePoolResource + ` ` + K8sNodePoolTestResource + ` {
  datacenter_id     = ` + DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + K8sClusterResource + `.terraform_acctest.id
  k8s_version = ` + K8sClusterResource + `.terraform_acctest.k8s_version
  name        = "` + K8sNodePoolTestResource + `"
  auto_scaling {
    min_node_count = 1
    max_node_count = 3
  }
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 2
  cores_count       = 1
  ram_size          = 2048
  storage_size      = 40
}`

const testAccCheckK8sNodePoolConfigNoOptionalFieldsUpdate = `
resource ` + DatacenterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource ` + K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "1.23.6"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}
resource ` + K8sNodePoolResource + ` ` + K8sNodePoolTestResource + ` {
  datacenter_id     = ` + DatacenterResource + `.terraform_acctest.id
  k8s_cluster_id    = ` + K8sClusterResource + `.terraform_acctest.id
  name        = "` + K8sNodePoolTestResource + `"
  k8s_version = ` + K8sClusterResource + `.terraform_acctest.k8s_version
  auto_scaling {
    min_node_count = 1
    max_node_count = 3
  }
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 1
  ram_size          = 2048
  storage_size      = 40
}
data ` + K8sNodePoolNodesResource + ` nodes{
  k8s_cluster_id   = ` + K8sClusterResource + `.terraform_acctest.id
  node_pool_id     = ` + K8sNodePoolResource + `.` + K8sNodePoolTestResource + `.id
}`
