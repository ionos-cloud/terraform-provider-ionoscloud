package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dsaas "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"regexp"
	"testing"
)

func TestAccDSaaSNodePoolBasic(t *testing.T) {
	var DSaaSNodePool dsaas.NodePoolResponseData

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDSaaSNodePoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDSaaSNodePoolConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDSaaSNodePoolExists(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, &DSaaSNodePool),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "name", DSaaSNodePoolTestResource),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "node_count", "1"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "cores_count", "1"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "ram_size", "2048"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "storage_size", "10"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "labels.foo", "bar"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "labels.color", "green"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "annotations.ann1", "value1"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "annotations.ann2", "value2"),
				),
			},
			{
				Config: testAccDataSourceDSaaSNodePoolMatchById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "name", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "node_count", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "node_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "cpu_family", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "cores_count", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "cores_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "ram_size", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "availability_zone", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "storage_type", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "storage_size", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "maintenance_window.0.time", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "maintenance_window.0.day_of_the_week", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "labels.foo", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "labels.foo"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "labels.color", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "labels.color"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "annotations.ann1", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceById, "annotations.ann2", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "annotations.ann2"),
				),
			},
			{
				Config: testAccDataSourceDSaaSNodePoolMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "name", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "node_count", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "node_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "cpu_family", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "cores_count", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "cores_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "ram_size", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "availability_zone", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "storage_type", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "storage_size", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "maintenance_window.0.time", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "maintenance_window.0.day_of_the_week", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "labels.foo", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "labels.foo"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "labels.color", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "labels.color"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "annotations.ann1", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "annotations.ann2", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "annotations.ann2")),
			},
			{
				Config: testAccDataSourceDSaaSNodePoolPartialMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "name", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "node_count", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "node_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "cpu_family", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "cores_count", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "cores_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "ram_size", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "availability_zone", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "storage_type", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "storage_size", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "maintenance_window.0.time", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "maintenance_window.0.day_of_the_week", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "labels.foo", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "labels.foo"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "labels.color", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "labels.color"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "annotations.ann1", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "annotations.ann1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSNodePoolResource+"."+DSaaSNodePoolTestDataSourceByName, "annotations.ann2", DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "annotations.ann2"),
				),
			},
			{
				Config:      testAccDataSourceDSaaSNodePoolWrongNameError,
				ExpectError: regexp.MustCompile("no DSaaS NodePool found with the specified name"),
			},
			{
				Config:      testAccDataSourceDSaaSNodePoolWrongPartialNameError,
				ExpectError: regexp.MustCompile("no DSaaS NodePool found with the specified name"),
			},
			{
				Config: testAccDataSourceDSaaSNodePools,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(DataSource+"."+DSaaSNodePoolsDataSource+"."+DSaaSNodePoolsTestDataSource, "node_pools.#"),
				),
			},
			{
				Config: testAccDataSourceDSaaSNodePoolsByName,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(DataSource+"."+DSaaSNodePoolsDataSource+"."+DSaaSNodePoolsTestDataSource, "node_pools.#"),
				),
			},
			{
				Config: testAccDataSourceDSaaSNodePoolsByNamePartialMatch,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(DataSource+"."+DSaaSNodePoolsDataSource+"."+DSaaSNodePoolsTestDataSource, "node_pools.#"),
				),
			},
			{
				Config:      testAccDataSourceDSaaSNodePoolsByNameError,
				ExpectError: regexp.MustCompile("no DSaaS NodePool found under cluster"),
			},
			{
				Config:      testAccDataSourceDSaaSNodePoolsByNamePartialMatchError,
				ExpectError: regexp.MustCompile("no DSaaS NodePool found under cluster"),
			},
			{
				Config: testAccCheckDSaaSNodePoolConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDSaaSNodePoolExists(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, &DSaaSNodePool),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "name", DSaaSNodePoolTestResource),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "node_count", "2"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "cpu_family", "INTEL_XEON"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "cores_count", "1"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "ram_size", "2048"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "availability_zone", "AUTO"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "storage_size", "10"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "labels.foo", "bar"),
					resource.TestCheckResourceAttr(DSaaSNodePoolResource+"."+DSaaSNodePoolTestResource, "annotations.ann1", "value1"),
				),
			},
		},
	})
}

func testAccCheckDSaaSNodePoolDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).DSaaSClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != DSaaSNodePoolResource {
			continue
		}

		clusterId := rs.Primary.Attributes["cluster_id"]
		nodePoolId := rs.Primary.ID

		_, apiResponse, err := client.GetNodePool(ctx, clusterId, nodePoolId)

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of DSaaS Node Pool %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("DSaaS NodePool %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckDSaaSNodePoolExists(n string, nodePool *dsaas.NodePoolResponseData) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).DSaaSClient

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
			return fmt.Errorf("an error occured while fetching DSaaS Node Pool %s: %s", rs.Primary.ID, err)
		}
		if *foundNodePool.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		nodePool = &foundNodePool

		return nil
	}
}

const testAccCheckDSaaSNodePoolConfigBasic = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fkb"
  description = "Datacenter for testing DSaaS NodePool"
}

resource ` + DSaaSClusterResource + ` ` + DSaaSClusterTestResource + ` {
	datacenter_id   		=  ` + DatacenterResource + `.datacenter_example.id
  	name 					= "` + DSaaSNodePoolTestResource + `"
  	maintenance_window {
    	day_of_the_week  	= "Sunday"
    	time				= "09:00:00"
  	}
  	data_platform_version	= "1.1.0"
}

resource ` + DSaaSNodePoolResource + ` ` + DSaaSNodePoolTestResource + ` {
  cluster_id    	= ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
  name        		= "` + DSaaSNodePoolTestResource + `"
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

const testAccCheckDSaaSNodePoolConfigUpdate = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fkb"
  description = "Datacenter for testing DSaaS NodePool"
}

resource ` + DSaaSClusterResource + ` ` + DSaaSClusterTestResource + ` {
	datacenter_id   		=  ` + DatacenterResource + `.datacenter_example.id
  	name 					= "` + UpdatedResources + `"
  	maintenance_window {
    	day_of_the_week  	= "Saturday"
    	time				= "10:00:00"
  	}
  	data_platform_version	= "1.1.0"
}

resource ` + DSaaSNodePoolResource + ` ` + DSaaSNodePoolTestResource + ` {
  cluster_id    = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
  name        = "` + DSaaSNodePoolTestResource + `"
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

const testAccDataSourceDSaaSNodePoolMatchById = testAccCheckDSaaSNodePoolConfigBasic + `
data ` + DSaaSNodePoolResource + ` ` + DSaaSNodePoolTestDataSourceById + ` {
  	cluster_id    = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id	
	id = ` + DSaaSNodePoolResource + `.` + DSaaSNodePoolTestResource + `.id
}
`

const testAccDataSourceDSaaSNodePoolMatchByName = testAccCheckDSaaSNodePoolConfigBasic + `
data ` + DSaaSNodePoolResource + ` ` + DSaaSNodePoolTestDataSourceByName + ` {
    cluster_id    = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
	name = "` + DSaaSNodePoolTestResource + `"
}
`

const testAccDataSourceDSaaSNodePoolPartialMatchByName = testAccCheckDSaaSNodePoolConfigBasic + `
data ` + DSaaSNodePoolResource + ` ` + DSaaSNodePoolTestDataSourceByName + ` {
	cluster_id    = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
	name = "test_"
    partial_match = true
}
`

const testAccDataSourceDSaaSNodePoolWrongNameError = testAccCheckDSaaSNodePoolConfigBasic + `
data ` + DSaaSNodePoolResource + ` ` + DSaaSNodePoolTestDataSourceByName + ` {
	cluster_id    = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
	name = "wrong_name"
}
`

const testAccDataSourceDSaaSNodePoolWrongPartialNameError = testAccCheckDSaaSNodePoolConfigBasic + `
data ` + DSaaSNodePoolResource + ` ` + DSaaSNodePoolTestDataSourceByName + ` {
    cluster_id    = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
	name = "wrong_name"
	partial_match = true
}
`
const testAccDataSourceDSaaSNodePools = testAccCheckDSaaSNodePoolConfigBasic + `
data ` + DSaaSNodePoolsDataSource + ` + ` + DSaaSNodePoolsTestDataSource + ` {
	cluster_id    = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
}
`

const testAccDataSourceDSaaSNodePoolsByName = testAccCheckDSaaSNodePoolConfigBasic + `
data ` + DSaaSNodePoolsDataSource + ` + ` + DSaaSNodePoolsTestDataSource + ` {
	cluster_id    = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
	name = "` + DSaaSNodePoolTestResource + `"}
`

const testAccDataSourceDSaaSNodePoolsByNamePartialMatch = testAccCheckDSaaSNodePoolConfigBasic + `
data ` + DSaaSNodePoolsDataSource + ` + ` + DSaaSNodePoolsTestDataSource + ` {
	cluster_id    = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
	name = "test_"
    partial_match = true
}
`

const testAccDataSourceDSaaSNodePoolsByNameError = testAccCheckDSaaSNodePoolConfigBasic + `
data ` + DSaaSNodePoolsDataSource + ` + ` + DSaaSNodePoolsTestDataSource + ` {
	cluster_id    = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
	name = "wrong_name"
}
`

const testAccDataSourceDSaaSNodePoolsByNamePartialMatchError = testAccCheckDSaaSNodePoolConfigBasic + `
data ` + DSaaSNodePoolsDataSource + ` + ` + DSaaSNodePoolsTestDataSource + ` {
	cluster_id    = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
	name = "wrong_name"
	partial_match = true
}
`
