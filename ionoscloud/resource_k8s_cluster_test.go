//go:build k8s
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

func TestAccK8sClusterBasic(t *testing.T) {
	var k8sCluster ionoscloud.KubernetesCluster

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckK8sClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sClusterExists(K8sClusterResource+"."+K8sClusterTestResource, &k8sCluster),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "name", K8sClusterTestResource),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "k8s_version", "1.19.10"),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.time", "09:00:00Z"),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "api_subnet_allow_list.0", "1.2.3.4/32"),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "s3_buckets.0.name", "sdktestv66"),
				),
			},
			{
				Config: testAccCheckK8sClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sClusterExists(K8sClusterResource+"."+K8sClusterTestResource, &k8sCluster),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "k8s_version", "1.19.10"),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.day_of_the_week", "Monday"),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.time", "10:30:00Z"),
					resource.TestCheckNoResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "api_subnet_allow_list"),
					resource.TestCheckNoResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "s3_buckets"),
				),
			},
			{
				Config: testAccCheckk8sClusterConfigUpdateVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sClusterExists(K8sClusterResource+"."+K8sClusterTestResource, &k8sCluster),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "k8s_version", "1.20.10"),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.day_of_the_week", "Monday"),
					resource.TestCheckResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "maintenance_window.0.time", "10:30:00Z"),
					resource.TestCheckNoResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "api_subnet_allow_list"),
					resource.TestCheckNoResourceAttr(K8sClusterResource+"."+K8sClusterTestResource, "s3_buckets")),
			},
		},
	})
}

func testAccCheckK8sClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != K8sClusterResource {
			continue
		}

		_, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of k8s cluster %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("k8s cluster %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckK8sClusterExists(n string, k8sCluster *ionoscloud.KubernetesCluster) resource.TestCheckFunc {
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
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("an error occured while fetching k8s Cluster %s: %s", rs.Primary.ID, err)
		}
		if *foundK8sCluster.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		k8sCluster = &foundK8sCluster

		return nil
	}
}

const testAccCheckK8sClusterConfigBasic = `
resource ` + K8sClusterResource + ` ` + K8sClusterTestResource + ` {
  name        = "` + K8sClusterTestResource + `"
  k8s_version = "1.19.10"
  maintenance_window {
    day_of_the_week = "Sunday"
    time            = "09:00:00Z"
  }
  api_subnet_allow_list = ["1.2.3.4/32"]
  s3_buckets { 
     name = "sdktestv66"
  }
}`

const testAccCheckK8sClusterConfigUpdate = `
resource ` + K8sClusterResource + ` ` + K8sClusterTestResource + ` {
  name        = "` + UpdatedResources + `"
  k8s_version = "1.19.14"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "10:30:00Z"
  }
  api_subnet_allow_list = []
  s3_buckets {}
}`

const testAccCheckk8sClusterConfigUpdateVersion = `
resource ` + K8sClusterResource + ` ` + K8sClusterTestResource + ` {
  name        = "` + UpdatedResources + `"
  k8s_version = "1.20.10"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "10:30:00Z"
  }
  api_subnet_allow_list = []
  s3_buckets {}
}`
