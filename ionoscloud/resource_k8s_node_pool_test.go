// +build k8s

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAcck8sNodepool_Basic(t *testing.T) {
	var k8sNodepool ionoscloud.KubernetesNodePool
	k8sNodepoolName := "terraform_acctest"

	publicIp1 := os.Getenv("TF_ACC_IONOS_PUBLIC_IP_1")
	if publicIp1 == "" {
		t.Errorf("TF_ACC_IONOS_PUBLIC_1 not set; please set it to a valid public IP for the us/las zone")
		t.FailNow()
	}
	publicIp2 := os.Getenv("TF_ACC_IONOS_PUBLIC_IP_2")
	if publicIp2 == "" {
		t.Errorf("TF_ACC_IONOS_PUBLIC_2 not set; please set it to a valid public IP for the us/las zone")
		t.FailNow()
	}
	publicIp3 := os.Getenv("TF_ACC_IONOS_PUBLIC_IP_3")
	if publicIp3 == "" {
		t.Errorf("TF_ACC_IONOS_PUBLIC_3 not set; please set it to a valid public IP for the us/las zone")
		t.FailNow()
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckk8sNodepoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckk8sNodepoolConfigBasic, k8sNodepoolName, publicIp1, publicIp2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sNodepoolExists("ionoscloud_k8s_node_pool.terraform_acctest", &k8sNodepool),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "name", k8sNodepoolName),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.0", publicIp1),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.1", publicIp2),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckk8sNodepoolConfigUpdate, k8sNodepoolName, publicIp1, publicIp2, publicIp3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sNodepoolExists("ionoscloud_k8s_node_pool.terraform_acctest", &k8sNodepool),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "name", k8sNodepoolName),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.0", publicIp1),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.1", publicIp2),
					resource.TestCheckResourceAttr("ionoscloud_k8s_node_pool.terraform_acctest", "public_ips.2", publicIp3),
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
			if apiResponse == nil || apiResponse.StatusCode != 404 {
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

		foundK8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID).Execute()

		if err != nil {
			payload := ""
			if apiResponse != nil {
				payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
			}
			return fmt.Errorf("error occured while fetching k8s node pool: %s %s", rs.Primary.ID, payload)
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
  public_ips        = [ "%s", "%s" ]
}`

const testAccCheckk8sNodepoolConfigUpdate = `
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
  public_ips        = [ "%s", "%s", "%s" ]
}`
