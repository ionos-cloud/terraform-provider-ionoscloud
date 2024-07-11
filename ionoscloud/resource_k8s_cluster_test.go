//go:build all || k8s
// +build all k8s

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccK8sClusterBasic(t *testing.T) {
	var k8sCluster ionoscloud.KubernetesCluster

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckK8sClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sClusterExists(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, &k8sCluster),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "name", constant.K8sClusterTestResource),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "maintenance_window.0.time", "09:00:00Z"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "api_subnet_allow_list.0", "1.2.3.4/32"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "s3_buckets.0.name", K8sBucket),
				),
			},
			{
				Config: testAccDataSourceK8sClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByID, "name", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByID, "k8s_version", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "k8s_version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByID, "maintenance_window.0.day_of_the_week", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByID, "maintenance_window.0.time", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByID, "api_subnet_allow_list.0", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "api_subnet_allow_list.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByID, "s3_buckets.0.name", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "s3_buckets.0.name"),
				),
			},
			{
				Config: testAccDataSourceK8sClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByName, "name", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByName, "k8s_version", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "k8s_version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByName, "maintenance_window.0.day_of_the_week", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByName, "maintenance_window.0.time", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByName, "api_subnet_allow_list.0", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "api_subnet_allow_list.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.K8sClusterDataSourceByName, "s3_buckets.0.name", constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "s3_buckets.0.name"),
				),
			},
			{
				Config:      testAccDataSourceK8sClusterWrongNameError,
				ExpectError: regexp.MustCompile("no cluster found with the specified name"),
			},
			{
				Config: testAccCheckK8sClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sClusterExists(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, &k8sCluster),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "public", "true"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "maintenance_window.0.day_of_the_week", "Monday"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "maintenance_window.0.time", "10:30:00Z"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "api_subnet_allow_list.0", "1.2.3.4/32"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "api_subnet_allow_list.1", "1.2.5.6/32"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "s3_buckets.0.name", K8sBucket),
				),
			},
			{
				Config: testAccCheckk8sClusterConfigUpdateVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sClusterExists(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, &k8sCluster),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "public", "true"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "maintenance_window.0.day_of_the_week", "Monday"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "maintenance_window.0.time", "10:30:00Z"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "api_subnet_allow_list.0", "1.2.3.4/32"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, "s3_buckets.0.name", K8sBucket),
				),
			},
		},
	})
}

func TestAccK8sClusterPrivate(t *testing.T) {
	var k8sCluster ionoscloud.KubernetesCluster

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckK8sClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckK8sClusterConfigPrivateCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sClusterExists(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, &k8sCluster),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "name", constant.PrivateK8sClusterTestResource),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "maintenance_window.0.time", "09:00:00Z"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "api_subnet_allow_list.0", "1.2.3.4/32"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "s3_buckets.0.name", K8sBucket),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "public", "false"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "location", "de/fra"),
					resource.TestCheckResourceAttrSet(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "nat_gateway_ip"),
					resource.TestCheckResourceAttr(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "node_subnet", K8sPrivateClusterNodeSubnet),
				),
			},
			{
				Config: testAccDataSourcePrivateK8sClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sClusterExists(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, &k8sCluster),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "name", constant.PrivateK8sClusterTestResource),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "k8s_version", K8sVersion),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "maintenance_window.0.time", "09:00:00Z"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "api_subnet_allow_list.0", "1.2.3.4/32"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "s3_buckets.0.name", K8sBucket),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "public", "false"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "location", "de/fra"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "nat_gateway_ip"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "node_subnet", K8sPrivateClusterNodeSubnet),
				),
			},
		},
	})
}

func TestAccK8sClusters(t *testing.T) {
	var k8sCluster ionoscloud.KubernetesCluster

	filtersErrRegex := fmt.Sprintf("no clusters found")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckK8sClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceK8sClusters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8sClusterExists(constant.K8sClusterResource+"."+constant.K8sClusterTestResource, &k8sCluster),
					testAccCheckK8sClusterExists(constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, &k8sCluster),
					// filter by 'name'
					// public cluster
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterName, "clusters.*", map[string]string{
						"name":                                 "test_k8s_cluster",
						"k8s_version":                          K8sVersion,
						"public":                               "true",
						"maintenance_window.0.day_of_the_week": "Sunday",
						"maintenance_window.0.time":            "09:00:00Z",
						"api_subnet_allow_list.0":              "1.2.3.4/32",
						"s3_buckets.0.name":                    "test_k8s_terraform_v7",
					}),
					// private cluster
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterName, "clusters.*", map[string]string{
						"name":                                 "test_private_k8s_cluster",
						"k8s_version":                          K8sVersion,
						"public":                               "false",
						"maintenance_window.0.day_of_the_week": "Sunday",
						"maintenance_window.0.time":            "09:00:00Z",
						"api_subnet_allow_list.0":              "1.2.3.4/32",
						"s3_buckets.0.name":                    "test_k8s_terraform_v7",
						"location":                             "de/fra",
						"node_subnet":                          "192.168.0.0/16",
					}),
					// filter by 'public'
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterPublic, "entries", constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterPublic, "clusters.#"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterPublic, "clusters.0.k8s_version", constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "k8s_version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterPublic, "clusters.0.public", constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "public"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterPublic, "clusters.0.maintenance_window.0.day_of_the_week", constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterPublic, "clusters.0.maintenance_window.0.time", constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterPublic, "clusters.0.api_subnet_allow_list.0", constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "api_subnet_allow_list.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterPublic, "clusters.0.location", constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterPublic, "clusters.0.nat_gateway_ip", constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "nat_gateway_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterPublic, "clusters.0.node_subnet", constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "node_subnet"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.K8sClustersDataSource+"."+constant.K8sClustersDataSourceFilterPublic, "clusters.0.s3_buckets.0.name", constant.K8sClusterResource+"."+constant.PrivateK8sClusterTestResource, "s3_buckets.0.name"),
				),
			},
			{
				Config:      testAccDataSourceK8sClustersNoMatchError,
				ExpectError: regexp.MustCompile(filtersErrRegex),
			},
		},
	})
}

func testAccCheckK8sClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.K8sClusterResource {
			continue
		}

		_, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking the destruction of k8s cluster %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("k8s cluster %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckK8sClusterExists(n string, k8sCluster *ionoscloud.KubernetesCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

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
			return fmt.Errorf("an error occurred while fetching k8s Cluster %s: %w", rs.Primary.ID, err)
		}
		if *foundK8sCluster.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		k8sCluster = &foundK8sCluster

		return nil
	}
}

const testAccCheckK8sClusterConfigBasic = `
resource ` + constant.K8sClusterResource + ` ` + constant.K8sClusterTestResource + ` {
  name        = "` + constant.K8sClusterTestResource + `"
  k8s_version = "` + K8sVersion + `"
  maintenance_window {
    day_of_the_week = "Sunday"
    time            = "09:00:00Z"
  }
  api_subnet_allow_list = ["1.2.3.4/32"]
  s3_buckets { 
     name = "` + K8sBucket + `"
  }
}`

const testAccCheckK8sClusterConfigUpdate = `
resource ` + constant.K8sClusterResource + ` ` + constant.K8sClusterTestResource + ` {
  name        = "` + constant.UpdatedResources + `"
  k8s_version = "` + K8sVersion + `"
  //public = "true"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "10:30:00Z"
  }
  api_subnet_allow_list = ["1.2.3.4/32", "1.2.5.6/32"]
  s3_buckets {
		name = "` + K8sBucket + `"
	}
}`

const testAccCheckk8sClusterConfigUpdateVersion = `
resource ` + constant.K8sClusterResource + ` ` + constant.K8sClusterTestResource + ` {
  name        = "` + constant.UpdatedResources + `"
  k8s_version = "` + K8sVersion + `"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "10:30:00Z"
  }
  api_subnet_allow_list = ["1.2.3.4/32"]
  s3_buckets {
		name = "` + K8sBucket + `"
	}
}`

const testAccCheckK8sClusterConfigPrivateCluster = `
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + ` {
  location = "de/fra"
  size = 1
  name = "` + constant.IpBlockTestResource + `"
}

resource ` + constant.K8sClusterResource + ` ` + constant.PrivateK8sClusterTestResource + ` {
  name        = "` + constant.PrivateK8sClusterTestResource + `"
  k8s_version = "` + K8sVersion + `"
  maintenance_window {
    day_of_the_week = "Sunday"
    time            = "09:00:00Z"
  }
  api_subnet_allow_list = ["1.2.3.4/32"]
  s3_buckets { 
     name = "` + K8sBucket + `"
  }
  public = false
  location = "de/fra"
  nat_gateway_ip = ` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[0]
  node_subnet = "` + K8sPrivateClusterNodeSubnet + `"
}`

const testAccDataSourceK8sClusterMatchId = testAccCheckK8sClusterConfigBasic + `
data ` + constant.K8sClusterResource + ` ` + constant.K8sClusterDataSourceByID + `{
  id	= ` + constant.K8sClusterResource + `.` + constant.K8sClusterTestResource + `.id
}
`

const testAccDataSourcePrivateK8sClusterMatchId = testAccCheckK8sClusterConfigPrivateCluster + `
data ` + constant.K8sClusterResource + ` ` + constant.PrivateK8sClusterTestResource + `{
  id	= ` + constant.K8sClusterResource + `.` + constant.PrivateK8sClusterTestResource + `.id
}
`

const testAccDataSourceK8sClusterMatchName = testAccCheckK8sClusterConfigBasic + `
data ` + constant.K8sClusterResource + ` ` + constant.K8sClusterDataSourceByName + `{
  name	= "` + constant.K8sClusterTestResource + `"
}
`

const testAccDataSourceK8sClusterWrongNameError = testAccCheckK8sClusterConfigBasic + `
data ` + constant.K8sClusterResource + ` ` + constant.K8sClusterDataSourceByName + `{
  name	= "wrong_name"
}
`

const testAccDataSourceK8sClusters = testAccCheckK8sClusterConfigBasic + "\n" + testAccCheckK8sClusterConfigPrivateCluster + `
data ` + constant.K8sClustersDataSource + ` ` + constant.K8sClustersDataSourceFilterName + `{
  depends_on = [` + constant.K8sClusterResource + `.` + constant.PrivateK8sClusterTestResource + `,` + constant.K8sClusterResource + `.` + constant.K8sClusterTestResource + `]
  filter{
    name = "name"
    value = "test_"
  }
}
data ` + constant.K8sClustersDataSource + ` ` + constant.K8sClustersDataSourceFilterPublic + `{
  depends_on = [` + constant.K8sClusterResource + `.` + constant.PrivateK8sClusterTestResource + `,` + constant.K8sClusterResource + `.` + constant.K8sClusterTestResource + `]
  filter{
    name = "public"
    value = "false"
  }
}
`

const testAccDataSourceK8sClustersNoMatchError = testAccCheckK8sClusterConfigBasic + "\n" + testAccCheckK8sClusterConfigPrivateCluster + `
data ` + constant.K8sClustersDataSource + ` ` + constant.K8sClustersDataSourceFilterName + `{
  depends_on = [` + constant.K8sClusterResource + `.` + constant.PrivateK8sClusterTestResource + `,` + constant.K8sClusterResource + `.` + constant.K8sClusterTestResource + `]
  filter{
    name = "` + filterName + `"
    value = "` + noMatchName + `"
  }
  filter{
	name = "` + filterVersion + `"
    value = "` + noMatchVersion + `"
  }
}
`

const (
	filterName     = "name"
	filterVersion  = "k8s_version"
	noMatchName    = "nomatch"
	noMatchVersion = "8.8"
)
