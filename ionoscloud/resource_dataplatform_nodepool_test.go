package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"regexp"
	"testing"
)

func TestAccDataplatformNodePoolBasic(t *testing.T) {
	var DataplatformNodePool dataplatform.NodePoolResponseData

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDataplatformNodePoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataplatformNodePoolConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataplatformNodePoolExists(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, &DataplatformNodePool),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "name", DataplatformNodePoolTestResource),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "node_count", "1"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "cores_count", "1"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "ram_size", "2048"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "storage_size", "10"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "labels.foo", "bar"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "labels.color", "green"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "annotations.ann1", "value1"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "annotations.ann2", "value2"),
				),
			},
			{
				Config: testAccDataSourceDataplatformNodePoolMatchById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "name", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "node_count", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "node_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "cpu_family", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "cores_count", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "cores_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "ram_size", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "availability_zone", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "storage_type", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "storage_size", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "maintenance_window.0.time", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "maintenance_window.0.day_of_the_week", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "labels.foo", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "labels.foo"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "labels.color", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "labels.color"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "annotations.ann1", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceById, "annotations.ann2", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "annotations.ann2"),
				),
			},
			{
				Config: testAccDataSourceDataplatformNodePoolMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "name", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "node_count", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "node_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "cpu_family", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "cores_count", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "cores_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "ram_size", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "availability_zone", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "storage_type", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "storage_size", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "maintenance_window.0.time", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "maintenance_window.0.day_of_the_week", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "labels.foo", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "labels.foo"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "labels.color", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "labels.color"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "annotations.ann1", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "annotations.ann2", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "annotations.ann2")),
			},
			{
				Config: testAccDataSourceDataplatformNodePoolPartialMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "name", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "node_count", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "node_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "cpu_family", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "cores_count", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "cores_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "ram_size", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "availability_zone", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "storage_type", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "storage_size", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "maintenance_window.0.time", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "maintenance_window.0.day_of_the_week", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "labels.foo", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "labels.foo"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "labels.color", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "labels.color"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "annotations.ann1", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformNodePoolResource+"."+DataplatformNodePoolTestDataSourceByName, "annotations.ann2", DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "annotations.ann2"),
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
					utils.TestNotEmptySlice(DataSource+"."+DataplatformNodePoolsDataSource+"."+DataplatformNodePoolsTestDataSource, "node_pools.#"),
				),
			},
			{
				Config: testAccDataSourceDataplatformNodePoolsByName,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(DataSource+"."+DataplatformNodePoolsDataSource+"."+DataplatformNodePoolsTestDataSource, "node_pools.#"),
				),
			},
			{
				Config: testAccDataSourceDataplatformNodePoolsByNamePartialMatch,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(DataSource+"."+DataplatformNodePoolsDataSource+"."+DataplatformNodePoolsTestDataSource, "node_pools.#"),
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
					testAccCheckDataplatformNodePoolExists(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, &DataplatformNodePool),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "name", DataplatformNodePoolTestResource),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "node_count", "2"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "cores_count", "1"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "ram_size", "2048"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "storage_size", "10"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "labels.foo", "bar"),
					resource.TestCheckResourceAttr(DataplatformNodePoolResource+"."+DataplatformNodePoolTestResource, "annotations.ann1", "value1"),
				),
			},
		},
	})
}

func testAccCheckDataplatformNodePoolDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).DataplatformClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != DataplatformNodePoolResource {
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
		client := testAccProvider.Meta().(SdkBundle).DataplatformClient

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
			return fmt.Errorf("an error occured while fetching Dataplatform Node Pool %s: %w", rs.Primary.ID, err)
		}
		if *foundNodePool.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		nodePool = &foundNodePool

		return nil
	}
}

const testAccCheckDataplatformNodePoolConfigBasic = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fra"
  description = "Datacenter for testing Dataplatform NodePool"
}

resource ` + DataplatformClusterResource + ` ` + DataplatformClusterTestResource + ` {
  datacenter_id   		=  ` + DatacenterResource + `.datacenter_example.id
  name 					= "` + DataplatformNodePoolTestResource + `"
  maintenance_window {
   	day_of_the_week  	= "Sunday"
   	time				= "09:00:00"
  }
  data_platform_version	= "22.09"
}

resource ` + DataplatformNodePoolResource + ` ` + DataplatformNodePoolTestResource + ` {
  cluster_id    	= ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
  name        		= "` + DataplatformNodePoolTestResource + `"
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
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fra"
  description = "Datacenter for testing Dataplatform NodePool"
}

resource ` + DataplatformClusterResource + ` ` + DataplatformClusterTestResource + ` {
  datacenter_id   		=  ` + DatacenterResource + `.datacenter_example.id
  name 					= "` + UpdatedResources + `"
  maintenance_window {
  	day_of_the_week  	= "Saturday"
   	time				= "10:00:00"
  }
  data_platform_version	= "22.09"
}

resource ` + DataplatformNodePoolResource + ` ` + DataplatformNodePoolTestResource + ` {
  cluster_id    = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
  name        = "` + DataplatformNodePoolTestResource + `"
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
data ` + DataplatformNodePoolResource + ` ` + DataplatformNodePoolTestDataSourceById + ` {
  cluster_id    = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id	
  id = ` + DataplatformNodePoolResource + `.` + DataplatformNodePoolTestResource + `.id
}
`

const testAccDataSourceDataplatformNodePoolMatchByName = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + DataplatformNodePoolResource + ` ` + DataplatformNodePoolTestDataSourceByName + ` {
  cluster_id    = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
  name = "` + DataplatformNodePoolTestResource + `"
}
`

const testAccDataSourceDataplatformNodePoolPartialMatchByName = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + DataplatformNodePoolResource + ` ` + DataplatformNodePoolTestDataSourceByName + ` {
  cluster_id    = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
  name = "test_"
  partial_match = true
}
`

const testAccDataSourceDataplatformNodePoolWrongNameError = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + DataplatformNodePoolResource + ` ` + DataplatformNodePoolTestDataSourceByName + ` {
  cluster_id    = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
  name = "wrong_name"
}
`

const testAccDataSourceDataplatformNodePoolWrongPartialNameError = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + DataplatformNodePoolResource + ` ` + DataplatformNodePoolTestDataSourceByName + ` {
  cluster_id    = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
  name = "wrong_name"
  partial_match = true
}
`
const testAccDataSourceDataplatformNodePools = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + DataplatformNodePoolsDataSource + ` + ` + DataplatformNodePoolsTestDataSource + ` {
	cluster_id    = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
}
`

const testAccDataSourceDataplatformNodePoolsByName = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + DataplatformNodePoolsDataSource + ` + ` + DataplatformNodePoolsTestDataSource + ` {
  cluster_id    = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
  name = "` + DataplatformNodePoolTestResource + `"}
`

const testAccDataSourceDataplatformNodePoolsByNamePartialMatch = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + DataplatformNodePoolsDataSource + ` + ` + DataplatformNodePoolsTestDataSource + ` {
  cluster_id    = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
  name = "test_"
  partial_match = true
}
`

const testAccDataSourceDataplatformNodePoolsByNameError = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + DataplatformNodePoolsDataSource + ` + ` + DataplatformNodePoolsTestDataSource + ` {
  cluster_id    = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
  name = "wrong_name"
}
`

const testAccDataSourceDataplatformNodePoolsByNamePartialMatchError = testAccCheckDataplatformNodePoolConfigBasic + `
data ` + DataplatformNodePoolsDataSource + ` + ` + DataplatformNodePoolsTestDataSource + ` {
  cluster_id    = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
  name = "wrong_name"
  partial_match = true
}
`
