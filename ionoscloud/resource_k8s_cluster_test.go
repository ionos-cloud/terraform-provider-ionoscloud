package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAcck8sCluster_Basic(t *testing.T) {
	var k8sCluster ionoscloud.KubernetesCluster
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
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckk8sClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_k8s_cluster" {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		_, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, rs.Primary.ID).Execute()

		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse == nil || apiResponse.Response.StatusCode != 404 {
				var payload = "<nil>"
				if apiResponse != nil {
					payload = string(apiResponse.Payload)
				}
				return fmt.Errorf("k8s cluster still exists %s %s", rs.Primary.ID, payload)
			}
		} else {
			return fmt.Errorf("Unable to fetch k8s cluster %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckk8sClusterExists(n string, k8sCluster *ionoscloud.KubernetesCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundK8sCluster, _, err := client.KubernetesApi.K8sFindByClusterId(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("Error occured while fetching k8s Cluster: %s", rs.Primary.ID)
		}
		if *foundK8sCluster.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		k8sCluster = &foundK8sCluster

		return nil
	}
}

const testAccCheckk8sClusterConfigBasic = `
resource "ionoscloud_k8s_cluster" "example" {
  name        = "%s"
  k8s_version = "1.20.6"
  maintenance_window {
    day_of_the_week = "Sunday"
    time            = "09:00:00Z"
  }
}`

const testAccCheckk8sClusterConfigUpdate = `
resource "ionoscloud_k8s_cluster" "example" {
  name        = "updated"
  k8s_version = "1.20.6"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "10:30:00Z"
  }
  
}`
