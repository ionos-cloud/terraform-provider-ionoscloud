// +build k8s

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_k8s_cluster" {
			continue
		}

		_, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				payload := ""
				if apiResponse != nil {
					payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
				}
				return fmt.Errorf("k8s cluster still exists %s - an error occurred while checking it %s %s", rs.Primary.ID, err, payload)
			}
		} else {
			return fmt.Errorf("k8s cluster still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckk8sClusterExists(n string, k8sCluster *ionoscloud.KubernetesCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundK8sCluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, rs.Primary.ID).Execute()

		if err != nil {
			payload := ""
			if apiResponse != nil {
				payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
			}
			return fmt.Errorf("error occured while fetching k8s Cluster: %s %s", rs.Primary.ID, payload)
		}
		if *foundK8sCluster.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
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
