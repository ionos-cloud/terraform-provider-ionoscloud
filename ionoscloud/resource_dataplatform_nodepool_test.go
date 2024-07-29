//go:build all || dataplatform

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccDataplatformNodePoolBasic(t *testing.T) {
	var DataplatformNodePool dataplatform.NodePoolResponseData

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDataplatformNodePoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataplatformNodePoolConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataplatformNodePoolExists(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, &DataplatformNodePool),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "name", constant.DataplatformNodePoolTestResource),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "node_count", "1"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "cores_count", "1"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "ram_size", "2048"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "storage_size", "10"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "labels.foo", "bar"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "labels.color", "green"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "annotations.ann1", "value1"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "annotations.ann2", "value2"),
				),
			},
			{
				Config: testAccDataSourceDataplatformNodePoolMatchById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "name", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "node_count", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "node_count"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "cpu_family", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "cores_count", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "cores_count"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "ram_size", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "ram_size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "availability_zone", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "storage_type", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "storage_size", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "maintenance_window.0.time", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "maintenance_window.0.day_of_the_week", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "labels.foo", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "labels.foo"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "labels.color", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "labels.color"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "annotations.ann1", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceById, "annotations.ann2", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "annotations.ann2"),
				),
			},
			{
				Config: testAccDataSourceDataplatformNodePoolMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "name", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "node_count", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "node_count"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "cpu_family", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "cores_count", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "cores_count"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "ram_size", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "ram_size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "availability_zone", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "storage_type", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "storage_size", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "maintenance_window.0.time", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "maintenance_window.0.day_of_the_week", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "labels.foo", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "labels.foo"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "labels.color", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "labels.color"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "annotations.ann1", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "annotations.ann2", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "annotations.ann2")),
			},
			{
				Config: testAccDataSourceDataplatformNodePoolPartialMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "name", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "node_count", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "node_count"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "cpu_family", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "cores_count", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "cores_count"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "ram_size", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "ram_size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "availability_zone", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "storage_type", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "storage_size", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "maintenance_window.0.time", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "maintenance_window.0.day_of_the_week", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "labels.foo", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "labels.foo"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "labels.color", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "labels.color"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "annotations.ann1", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestDataSourceByName, "annotations.ann2", constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "annotations.ann2"),
				),
			},
			{
				Config:      testAccDataSourceDataplatformNodePoolWrongNameError,
				ExpectError: regexp.MustCompile("no Dataplatform NodePool found with the specified name"),
			},
			{
				Config:      testAccDataSourceDataplatformNodePoolWrongPartialNameError,
				ExpectError: regexp.MustCompile("no Dataplatform NodePool found with the specified name"),
			},
			{
				Config: testAccDataSourceDataplatformNodePools,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(constant.DataSource+"."+constant.DataplatformNodePoolsDataSource+"."+constant.DataplatformNodePoolsTestDataSource, "node_pools.#"),
				),
			},
			{
				Config: testAccDataSourceDataplatformNodePoolsByName,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(constant.DataSource+"."+constant.DataplatformNodePoolsDataSource+"."+constant.DataplatformNodePoolsTestDataSource, "node_pools.#"),
				),
			},
			{
				Config: testAccDataSourceDataplatformNodePoolsByNamePartialMatch,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(constant.DataSource+"."+constant.DataplatformNodePoolsDataSource+"."+constant.DataplatformNodePoolsTestDataSource, "node_pools.#"),
				),
			},
			{
				Config:      testAccDataSourceDataplatformNodePoolsByNameError,
				ExpectError: regexp.MustCompile("no Dataplatform NodePool found under cluster"),
			},
			{
				Config:      testAccDataSourceDataplatformNodePoolsByNamePartialMatchError,
				ExpectError: regexp.MustCompile("no Dataplatform NodePool found under cluster"),
			},
			{
				Config: testAccCheckDataplatformNodePoolConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataplatformNodePoolExists(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, &DataplatformNodePool),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "name", constant.DataplatformNodePoolTestResource),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "node_count", "2"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "cores_count", "1"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "ram_size", "2048"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "storage_size", "10"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "labels.foo", "bar"),
					resource.TestCheckResourceAttr(constant.DataplatformNodePoolResource+"."+constant.DataplatformNodePoolTestResource, "annotations.ann1", "value1"),
				),
			},
		},
	})
}

func testAccCheckDataplatformNodePoolDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).DataplatformClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DataplatformNodePoolResource {
			continue
		}

		clusterId := rs.Primary.Attributes["cluster_id"]
		nodePoolId := rs.Primary.ID

		_, apiResponse, err := client.GetNodePool(ctx, clusterId, nodePoolId)

		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of Dataplatform Node Pool %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("Dataplatform NodePool %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckDataplatformNodePoolExists(n string, nodePool *dataplatform.NodePoolResponseData) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).DataplatformClient

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

		clusterId := rs.Primary.Attributes["cluster_id"]
		nodePoolId := rs.Primary.ID

		foundNodePool, _, err := client.GetNodePool(ctx, clusterId, nodePoolId)

		if err != nil {
			return fmt.Errorf("an error occurred while fetching Dataplatform Node Pool %s: %w", rs.Primary.ID, err)
		}
		if *foundNodePool.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		nodePool = &foundNodePool

		return nil
	}
}

const testAccCheckDataplatformNodePoolConfigBasic = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fra"
  description = "Datacenter for testing Dataplatform NodePool"
}

resource ` + constant.DataplatformClusterResource + ` ` + constant.DataplatformClusterTestResource + ` {
  datacenter_id   		=  ` + constant.DatacenterResource + `.datacenter_example.id
  name 					= "` + constant.DataplatformNodePoolTestResource + `"
  maintenance_window {
   	day_of_the_week  	= "Sunday"
   	time				= "09:00:00"
  }
  version	= ` + constant.DataPlatformVersion + `
}

resource ` + constant.DataplatformNodePoolResource + ` ` + constant.DataplatformNodePoolTestResource + ` {
  cluster_id    	= ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
  name        		= "` + constant.DataplatformNodePoolTestResource + `"
  node_count        = 1
  cpu_family        = "INTEL_XEON"
  cores_count       = 1
  ram_size          = 2048
  availability_zone = "AUTO"
  storage_type      = "HDD"
  storage_size      = 10
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00"
  }
  labels 			= {
    foo   			= "bar"
    color 			= "green"
  }
  annotations 		= {
    ann1 			= "value1"
    ann2 			= "value2"
  }
}
`

const testAccCheckDataplatformNodePoolConfigUpdate = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fra"
  description = "Datacenter for testing Dataplatform NodePool"
}

resource ` + constant.DataplatformClusterResource + ` ` + constant.DataplatformClusterTestResource + ` {
  datacenter_id   		=  ` + constant.DatacenterResource + `.datacenter_example.id
  name 					= "` + constant.UpdatedResources + `"
  maintenance_window {
  	day_of_the_week  	= "Saturday"
   	time				= "10:00:00"
  }
  version	= ` + constant.DataPlatformVersion + `
}

resource ` + constant.DataplatformNodePoolResource + ` ` + constant.DataplatformNodePoolTestResource + ` {
  cluster_id    = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
  name        = "` + constant.DataplatformNodePoolTestResource + `"
  node_count        = 2
  cpu_family        = "INTEL_XEON"
  cores_count       = 1
  ram_size          = 2048
  availability_zone = "AUTO"
  storage_type      = "HDD"
  storage_size      = 10
  maintenance_window {
    day_of_the_week = "Sunday"
    time            = "10:00:00"
  }
  labels = {
    foo   = "bar"
  }
  annotations = {
    ann1 = "value1"
  }
}
`

const testAccDataSourceDataplatformNodePoolMatchById = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + constant.DataplatformNodePoolResource + ` ` + constant.DataplatformNodePoolTestDataSourceById + ` {
  cluster_id    = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id	
  id = ` + constant.DataplatformNodePoolResource + `.` + constant.DataplatformNodePoolTestResource + `.id
}
`

const testAccDataSourceDataplatformNodePoolMatchByName = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + constant.DataplatformNodePoolResource + ` ` + constant.DataplatformNodePoolTestDataSourceByName + ` {
  cluster_id    = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
  name = "` + constant.DataplatformNodePoolTestResource + `"
}
`

const testAccDataSourceDataplatformNodePoolPartialMatchByName = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + constant.DataplatformNodePoolResource + ` ` + constant.DataplatformNodePoolTestDataSourceByName + ` {
  cluster_id    = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
  name = "test_"
  partial_match = true
}
`

const testAccDataSourceDataplatformNodePoolWrongNameError = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + constant.DataplatformNodePoolResource + ` ` + constant.DataplatformNodePoolTestDataSourceByName + ` {
  cluster_id    = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
  name = "wrong_name"
}
`

const testAccDataSourceDataplatformNodePoolWrongPartialNameError = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + constant.DataplatformNodePoolResource + ` ` + constant.DataplatformNodePoolTestDataSourceByName + ` {
  cluster_id    = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
  name = "wrong_name"
  partial_match = true
}
`
const testAccDataSourceDataplatformNodePools = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + constant.DataplatformNodePoolsDataSource + ` ` + constant.DataplatformNodePoolsTestDataSource + ` {
	cluster_id    = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
}
`

const testAccDataSourceDataplatformNodePoolsByName = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + constant.DataplatformNodePoolsDataSource + ` ` + constant.DataplatformNodePoolsTestDataSource + ` {
  cluster_id    = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
  name = "` + constant.DataplatformNodePoolTestResource + `"
}
`

const testAccDataSourceDataplatformNodePoolsByNamePartialMatch = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + constant.DataplatformNodePoolsDataSource + ` ` + constant.DataplatformNodePoolsTestDataSource + ` {
  cluster_id    = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
  name = "test_"
  partial_match = true
}
`

const testAccDataSourceDataplatformNodePoolsByNameError = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + constant.DataplatformNodePoolsDataSource + ` ` + constant.DataplatformNodePoolsTestDataSource + ` {
  cluster_id    = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
  name = "wrong_name"
}
`

const testAccDataSourceDataplatformNodePoolsByNamePartialMatchError = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + constant.DataplatformNodePoolsDataSource + ` ` + constant.DataplatformNodePoolsTestDataSource + ` {
  cluster_id    = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
  name = "wrong_name"
  partial_match = true
}
`
