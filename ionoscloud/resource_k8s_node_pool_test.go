package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAcck8sNodepool_Basic(t *testing.T) {
	var k8sNodepool ionoscloud.KubernetesNodePool
	k8sNodepoolName := "terraform_acctest"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckk8sNodepoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckk8sNodepoolConfigBasic, k8sNodepoolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sNodepoolExists("ionoscloud_k8s_node_pool.terraform_acctest", &k8sNodepool),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "name", k8sNodepoolName),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.0", "157.97.108.242"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.1", "217.160.200.54"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckk8sNodepoolConfigUpdate, k8sNodepoolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sNodepoolExists("ionoscloud_k8s_node_pool.terraform_acctest", &k8sNodepool),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "name", k8sNodepoolName),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.0", "157.97.108.242"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.1", "217.160.200.54"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.2", "217.160.200.55"),
					//					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "maintenance_window.0.day_of_the_week", "Tuesday"),
					//					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "maintenance_window.0.time", "11:00:00Z"),
				),
			},
		},
	})
}

func testAccCheckk8sNodepoolDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_k8s_node_pool" {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		_, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID).Execute()

		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode != 404 {
				return fmt.Errorf("K8s node pool still exists %s %s", rs.Primary.ID, string(apiResponse.Payload))
			}
		} else {
			return fmt.Errorf("Unable to fetch k8s node pool %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckk8sNodepoolExists(n string, k8sNodepool *ionoscloud.KubernetesNodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		log.Printf("[INFO] REQ PATH: %+v/%+v", rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID)

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundK8sNodepool, _, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("Error occured while fetching k8s node pool: %s", rs.Primary.ID)
		}
		if *foundK8sNodepool.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
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
  public_ips        = [ "158.222.102.145", "158.222.102.144" ]
}`

const testAccCheckk8sNodepoolConfigUpdate = `
resource "ionoscloud_datacenter" "terraform_acctest" {
  name        = "terraform_acctest"
  location    = "us_las"
  description = "Datacenter created through terraform"
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
  public_ips        = [ "158.222.102.145", "158.222.102.144", "158.222.102.179" ]
}`
