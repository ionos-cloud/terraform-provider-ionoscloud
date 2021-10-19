// +build k8s

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAcck8sNodepool_Basic(t *testing.T) {
	var k8sNodepool ionoscloud.KubernetesNodePool
	k8sNodepoolName := "terraform_acctest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckk8sNodepoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckk8sNodepoolConfigBasic, k8sNodepoolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sNodepoolExists("ionoscloud_k8s_node_pool.terraform_acctest", &k8sNodepool),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "name", k8sNodepoolName),
					resource.TestCheckResourceAttrPair("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.0", "ionoscloud_ipblock.terraform_acctest", "ips.0"),
					resource.TestCheckResourceAttrPair("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.1", "ionoscloud_ipblock.terraform_acctest", "ips.1"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "labels.foo", "bar"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "labels.color", "green"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "annotations.ann1", "value1"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "annotations.ann2", "value2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckk8sNodepoolConfigUpdate, k8sNodepoolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sNodepoolExists("ionoscloud_k8s_node_pool.terraform_acctest", &k8sNodepool),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "name", k8sNodepoolName),
					resource.TestCheckResourceAttrPair("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.0", "ionoscloud_ipblock.terraform_acctest", "ips.0"),
					resource.TestCheckResourceAttrPair("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.1", "ionoscloud_ipblock.terraform_acctest", "ips.1"),
					resource.TestCheckResourceAttrPair("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.2", "ionoscloud_ipblock.terraform_acctest", "ips.2"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "labels.foo", "baz"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "labels.color", "red"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "labels.third", "thirdValue"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "annotations.ann1", "value1Changed"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "annotations.ann2", "value2Changed"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "annotations.ann3", "newValue"),
				),
			},
		},
	})
}

func TestAcck8sNodepool_Version(t *testing.T) {
	var k8sNodepool ionoscloud.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckk8sNodepoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckk8sNodepoolConfigVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sNodepoolExists("ionoscloud_k8s_node_pool.terraform_acctest", &k8sNodepool),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "k8s_version", "1.18.5"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckk8sNodepoolConfigIgnoreVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sNodepoolExists("ionoscloud_k8s_node_pool.terraform_acctest", &k8sNodepool),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "k8s_version", "1.18.5"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckk8sNodepoolConfigChangeVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sNodepoolExists("ionoscloud_k8s_node_pool.terraform_acctest", &k8sNodepool),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "k8s_version", "1.19.10"),
				),
			},
		},
	})
}

func testAccCheckk8sNodepoolDestroyCheck(s *terraform.State) error {
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

func testAccCheckk8sNodepoolExists(n string, k8sNodepool *ionoscloud.KubernetesNodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

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

		foundK8sNodepool, _, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID).Execute()

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

const testAccCheckk8sNodepoolConfigBasic = `
resource "ionoscloud_datacenter" "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource "ionoscloud_lan" "terraform_acctest" {
  datacenter_id = "${ionoscloud_datacenter.terraform_acctest.id}"
  public = false
  name = "terraform_acctest"
}

resource "ionoscloud_ipblock" "terraform_acctest" {
  location = ionoscloud_datacenter.terraform_acctest.location
  size = 3
  name = "terraform_acctest"
}

resource "ionoscloud_k8s_cluster" "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "1.18.9"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}

resource "ionoscloud_k8s_node_pool" "terraform_acctest" {
  name        = "%s"
  k8s_version = "${ionoscloud_k8s_cluster.terraform_acctest.k8s_version}"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
  datacenter_id     = "${ionoscloud_datacenter.terraform_acctest.id}"
  k8s_cluster_id    = "${ionoscloud_k8s_cluster.terraform_acctest.id}"
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
  public_ips        = [ ionoscloud_ipblock.terraform_acctest.ips[0], ionoscloud_ipblock.terraform_acctest.ips[1] ]
  lans 	            = [ 1 ]
  labels = {
    foo = "bar"
    color = "green"
  }
  annotations = {
    ann1 = "value1"
    ann2 = "value2"
  }
}`

const testAccCheckk8sNodepoolConfigUpdate = `
resource "ionoscloud_datacenter" "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource "ionoscloud_lan" "terraform_acctest" {
  datacenter_id = "${ionoscloud_datacenter.terraform_acctest.id}"
  public = false
  name = "terraform_acctest"
}

resource "ionoscloud_ipblock" "terraform_acctest" {
  location = ionoscloud_datacenter.terraform_acctest.location
  size = 3
  name = "terraform_acctest"
}

resource "ionoscloud_k8s_cluster" "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "1.18.9"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}

resource "ionoscloud_k8s_node_pool" "terraform_acctest" {
  name        = "%s"
  k8s_version = "${ionoscloud_k8s_cluster.terraform_acctest.k8s_version}"
  auto_scaling {
  	min_node_count = 1
	max_node_count = 2
  }
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
  datacenter_id     = "${ionoscloud_datacenter.terraform_acctest.id}"
  k8s_cluster_id    = "${ionoscloud_k8s_cluster.terraform_acctest.id}"
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
  public_ips        = [ ionoscloud_ipblock.terraform_acctest.ips[0], ionoscloud_ipblock.terraform_acctest.ips[1], ionoscloud_ipblock.terraform_acctest.ips[2] ]
  lans 	            = [ 1 ]
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

const testAccCheckk8sNodepoolConfigVersion = `
resource "ionoscloud_datacenter" "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource "ionoscloud_k8s_cluster" "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "1.18.5"
}

resource "ionoscloud_k8s_node_pool" "terraform_acctest" {
  name        = "test_version"
  k8s_version = "${ionoscloud_k8s_cluster.terraform_acctest.k8s_version}"
  datacenter_id     = "${ionoscloud_datacenter.terraform_acctest.id}"
  k8s_cluster_id    = "${ionoscloud_k8s_cluster.terraform_acctest.id}"
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
}`

const testAccCheckk8sNodepoolConfigIgnoreVersion = `
resource "ionoscloud_datacenter" "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource "ionoscloud_k8s_cluster" "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "1.18.9"
}

resource "ionoscloud_k8s_node_pool" "terraform_acctest" {
  name        = "test_version"
  k8s_version = "${ionoscloud_k8s_cluster.terraform_acctest.k8s_version}"
  datacenter_id     = "${ionoscloud_datacenter.terraform_acctest.id}"
  k8s_cluster_id    = "${ionoscloud_k8s_cluster.terraform_acctest.id}"
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
}`

const testAccCheckk8sNodepoolConfigChangeVersion = `
resource "ionoscloud_datacenter" "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us/las"
  description = "Datacenter created through terraform"
}

resource "ionoscloud_k8s_cluster" "terraform_acctest" {
  name        = "terraform_acctest"
  k8s_version = "1.19.10"
}

resource "ionoscloud_k8s_node_pool" "terraform_acctest" {
  name        = "test_version"
  k8s_version = "${ionoscloud_k8s_cluster.terraform_acctest.k8s_version}"
  datacenter_id     = "${ionoscloud_datacenter.terraform_acctest.id}"
  k8s_cluster_id    = "${ionoscloud_k8s_cluster.terraform_acctest.id}"
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
}`
