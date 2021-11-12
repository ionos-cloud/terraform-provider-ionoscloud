// +build k8s

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
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
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "k8s_version", "1.19.10"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "maintenance_window.0.day_of_the_week", "Monday"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "maintenance_window.0.time", "09:00:00Z"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "auto_scaling.0.min_node_count", "1"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "auto_scaling.0.max_node_count", "1"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_type", "SSD"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "node_count", "1"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "cores_count", "2"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "ram_size", "2048"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "storage_size", "40"),
					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "public_ips.0", IpBLockResource+".terraform_acctest", "ips.0"),
					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "public_ips.1", IpBLockResource+".terraform_acctest", "ips.1"),
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
				Config: testAccCheckK8sNodePoolConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sNodePoolExists(ResourceNameK8sNodePool, &k8sNodepool),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "name", K8sNodePoolTestResource),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "k8s_version", "1.19.10"),
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
					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "public_ips.0", IpBLockResource+".terraform_acctest", "ips.0"),
					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "public_ips.1", IpBLockResource+".terraform_acctest", "ips.1"),
					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "public_ips.2", IpBLockResource+".terraform_acctest", "ips.2"),
					resource.TestCheckResourceAttrPair(ResourceNameK8sNodePool, "lans.0.id", LanResource+".terraform_acctest_updated", "id"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.dhcp", "false"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.routes.0.network", "1.2.3.4/24"),
					resource.TestCheckResourceAttr(ResourceNameK8sNodePool, "lans.0.routes.0.gateway_ip", "10.1.5.16"),
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
					resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "public_ips"),
					resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "lans"),
					resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "labels"),
					resource.TestCheckNoResourceAttr(ResourceNameK8sNodePool, "annotations")),
			},
		},
	})
}

func testAccCheckK8sNodePoolDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_k8s_node_pool" {
			continue
		}

		_, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
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
resource ` + IpBLockResource + ` "terraform_acctest" {
  location = ` + DatacenterResource + `.terraform_acctest.location
  size = 3
  name = "terraform_acctest"
}
resource ` + K8sClusterResource + ` "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "1.19.10"
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
  public_ips        = [ ` + IpBLockResource + `.terraform_acctest.ips[0], ` + IpBLockResource + `.terraform_acctest.ips[1] ]
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
resource ` + IpBLockResource + ` "terraform_acctest" {
	location = ` + DatacenterResource + `.terraform_acctest.location
	size = 3
	name = "terraform_acctest"
}
resource ` + K8sClusterResource + ` "terraform_acctest" {
	name        = "terraform_acctest"
	k8s_version = "1.19.14"
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
  public_ips        = [ ionoscloud_ipblock.terraform_acctest.ips[0], ionoscloud_ipblock.terraform_acctest.ips[1], ionoscloud_ipblock.terraform_acctest.ips[2] ]
  lans {
    id   = ` + LanResource + `.terraform_acctest_updated.id
    dhcp = false
 	routes {
       network   = "1.2.3.4/24"
       gateway_ip = "10.1.5.16"
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
resource ` + IpBLockResource + ` "terraform_acctest" {
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
  public_ips        = []
  labels = {}
  annotations = {}
}`
