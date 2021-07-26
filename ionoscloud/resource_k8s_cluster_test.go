// +build k8s

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckk8sClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckk8sClusterConfigBasic, k8sClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sClusterExists("ionoscloud_k8s_cluster.example", &k8sCluster),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "name", k8sClusterName),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "public", "true"),
				),
			},
			{
				Config: testAccCheckk8sClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sClusterExists("ionoscloud_k8s_cluster.example", &k8sCluster),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "name", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "public", "true"),
				),
			},
		},
	})
}

func TestAcck8sCluster_Version(t *testing.T) {
	var k8sCluster ionoscloud.KubernetesCluster

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckk8sClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckk8sClusterConfigVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sClusterExists("ionoscloud_k8s_cluster.example", &k8sCluster),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "name", "test_version"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "k8s_version", "1.18.5"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckk8sClusterConfigIgnoreVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sClusterExists("ionoscloud_k8s_cluster.example", &k8sCluster),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "name", "test_version_ignore"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "k8s_version", "1.18.5"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckk8sClusterConfigChangeVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sClusterExists("ionoscloud_k8s_cluster.example", &k8sCluster),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "name", "test_version_change"),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "k8s_version", "1.19.10"),
				),
			},
		},
	})
}

func TestAcck8sCluster_S3Subnet(t *testing.T) {
	var k8sCluster ionoscloud.KubernetesCluster
	k8sClusterName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckk8sClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckk8sClusterConfigS3Subnet, k8sClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckk8sClusterExists("ionoscloud_k8s_cluster.example", &k8sCluster),
					resource.TestCheckResourceAttr("ionoscloud_k8s_cluster.example", "name", k8sClusterName),
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
				return fmt.Errorf("an error occurred while checking the destruction of k8s cluster %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("k8s cluster %s still exists", rs.Primary.ID)
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

		foundK8sCluster, _, err := client.KubernetesApi.K8sFindByClusterId(ctx, rs.Primary.ID).Execute()

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

const testAccCheckk8sClusterConfigBasic = `
resource "ionoscloud_k8s_cluster" "example" {
  name        = "%s"
  k8s_version = "1.20.8"
  maintenance_window {
    day_of_the_week = "Sunday"
    time            = "09:00:00Z"
  }
}`

const testAccCheckk8sClusterConfigUpdate = `
resource "ionoscloud_k8s_cluster" "example" {
  name        = "updated"
  k8s_version = "1.20.8"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "10:30:00Z"
  } 
}`

const testAccCheckk8sClusterConfigS3Subnet = `
resource "ionoscloud_k8s_cluster" "example" {
  name        = "%s"
  api_subnet_allow_list = ["1.2.3.4/32",
                           "2002::1234:abcd:ffff:c0a8:101/64", 
                           "1.2.3.4", 
                           "2002::1234:abcd:ffff:c0a8:101" ]
  s3Buckets { 
     name = "sdktestv6"
  }
}`

const testAccCheckk8sClusterConfigVersion = `
resource "ionoscloud_k8s_cluster" "example" {
  name        = "test_version"
  k8s_version = "1.18.5"
}`

const testAccCheckk8sClusterConfigIgnoreVersion = `
resource "ionoscloud_k8s_cluster" "example" {
  name        = "test_version_ignore"
  k8s_version = "1.18.9"
}`

const testAccCheckk8sClusterConfigChangeVersion = `
resource "ionoscloud_k8s_cluster" "example" {
  name        = "test_version_change"
  k8s_version = "1.19.10"
}`
