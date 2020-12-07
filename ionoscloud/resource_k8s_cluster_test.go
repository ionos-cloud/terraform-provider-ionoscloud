package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAcck8sCluster_Basic(t *testing.T) {
	var k8sCluster profitbricks.KubernetesCluster
	k8sClusterName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckk8sClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckk8sClusterConfigBasic, k8sClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sClusterExists("ionoscloud_k8s_cluster.example", &k8sCluster),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "name", k8sClusterName),
				),
			},
			{
				Config: testAccCheckk8sClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sClusterExists("ionoscloud_k8s_cluster.example", &k8sCluster),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "name", "example-renamed"),
				),
			},
		},
	})
}

func testAccCheckk8sClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*profitbricks.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_k8s_cluster" {
			continue
		}

		_, err := client.GetKubernetesCluster(rs.Primary.ID)

		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("K8s cluster still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetch k8s cluster %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckk8sClusterExists(n string, k8sCluster *profitbricks.KubernetesCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*profitbricks.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundK8sCluster, err := client.GetKubernetesCluster(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching k8s Cluster: %s", rs.Primary.ID)
		}
		if foundK8sCluster.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		k8sCluster = foundK8sCluster

		return nil
	}
}

const testAccCheckk8sClusterConfigBasic = `
resource "ionoscloud_k8s_cluster" "example" {
  name        = "%s"
	k8s_version = "1.18.5"
  maintenance_window {
    day_of_the_week = "Sunday"
    time            = "09:00:00Z"
  }
}`

const testAccCheckk8sClusterConfigUpdate = `
resource "ionoscloud_k8s_cluster" "example" {
  name        = "example-renamed"
  k8s_version = "1.18.5"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "10:30:00Z"
  }
}`
