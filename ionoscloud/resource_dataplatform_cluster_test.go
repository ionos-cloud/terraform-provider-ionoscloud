//go:build all || dataplatform

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"regexp"
	"testing"
)

func TestAccDataplatformClusterBasic(t *testing.T) {
	var DataplatformCluster dataplatform.ClusterResponseData

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDataplatformClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataplatformClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataplatformClusterExists(DataplatformClusterResource+"."+DataplatformClusterTestResource, &DataplatformCluster),
					resource.TestCheckResourceAttr(DataplatformClusterResource+"."+DataplatformClusterTestResource, "name", DataplatformClusterTestResource),
					resource.TestCheckResourceAttr(DataplatformClusterResource+"."+DataplatformClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(DataplatformClusterResource+"."+DataplatformClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(DataplatformClusterResource+"."+DataplatformClusterTestResource, "data_platform_version", DataPlatformVersion),
				),
			},
			{
				Config: testAccDataSourceDataplatformClusterMatchById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceById, "name", DataplatformClusterResource+"."+DataplatformClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceById, "maintenance_window.0.time", DataplatformClusterResource+"."+DataplatformClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceById, "maintenance_window.0.day_of_the_week", DataplatformClusterResource+"."+DataplatformClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceById, "data_platform_version", DataplatformClusterResource+"."+DataplatformClusterTestResource, "data_platform_version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceById, "datacenter_id", DataplatformClusterResource+"."+DataplatformClusterTestResource, "datacenter_id"),
				),
			},
			{
				Config: testAccDataSourceDataplatformClusterMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceByName, "name", DataplatformClusterResource+"."+DataplatformClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceByName, "maintenance_window.0.time", DataplatformClusterResource+"."+DataplatformClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceByName, "maintenance_window.0.day_of_the_week", DataplatformClusterResource+"."+DataplatformClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceByName, "data_platform_version", DataplatformClusterResource+"."+DataplatformClusterTestResource, "data_platform_version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceByName, "datacenter_id", DataplatformClusterResource+"."+DataplatformClusterTestResource, "datacenter_id"),
				),
			},
			{
				Config: testAccDataSourceDataplatformClusterPartialMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceByName, "name", DataplatformClusterResource+"."+DataplatformClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceByName, "maintenance_window.0.time", DataplatformClusterResource+"."+DataplatformClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceByName, "maintenance_window.0.day_of_the_week", DataplatformClusterResource+"."+DataplatformClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceByName, "data_platform_version", DataplatformClusterResource+"."+DataplatformClusterTestResource, "data_platform_version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DataplatformClusterResource+"."+DataplatformClusterTestDataSourceByName, "datacenter_id", DataplatformClusterResource+"."+DataplatformClusterTestResource, "datacenter_id"),
				),
			},
			{
				Config:      testAccDataSourceDataplatformClusterWrongNameError,
				ExpectError: regexp.MustCompile("no Dataplatform Cluster found with the specified name"),
			},
			{
				Config:      testAccDataSourceDataplatformClusterWrongPartialNameError,
				ExpectError: regexp.MustCompile("no Dataplatform Cluster found with the specified name"),
			},
			{
				Config: testAccCheckDataplatformClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataplatformClusterExists(DataplatformClusterResource+"."+DataplatformClusterTestResource, &DataplatformCluster),
					resource.TestCheckResourceAttr(DataplatformClusterResource+"."+DataplatformClusterTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(DataplatformClusterResource+"."+DataplatformClusterTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(DataplatformClusterResource+"."+DataplatformClusterTestResource, "maintenance_window.0.day_of_the_week", "Saturday"),
					resource.TestCheckResourceAttr(DataplatformClusterResource+"."+DataplatformClusterTestResource, "data_platform_version", DataPlatformVersion),
				),
			},
		},
	})
}

func testAccCheckDataplatformClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).DataplatformClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != DataplatformClusterResource {
			continue
		}

		clusterId := rs.Primary.ID

		_, apiResponse, err := client.GetClusterById(ctx, clusterId)

		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of Dataplatform cluster %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("Dataplatform cluster %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckDataplatformClusterExists(n string, cluster *dataplatform.ClusterResponseData) resource.TestCheckFunc {
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

		clusterId := rs.Primary.ID

		foundCluster, _, err := client.GetClusterById(ctx, clusterId)

		if err != nil {
			return fmt.Errorf("an error occured while fetching Dataplatform Cluster %s: %w", rs.Primary.ID, err)
		}
		if *foundCluster.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		cluster = &foundCluster

		return nil
	}
}

const testAccCheckDataplatformClusterConfigBasic = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fra"
  description = "Datacenter for testing Dataplatform Cluster"
}

resource ` + DataplatformClusterResource + ` ` + DataplatformClusterTestResource + ` {
  datacenter_id   		=  ` + DatacenterResource + `.datacenter_example.id
  name 					= "` + DataplatformClusterTestResource + `"
  maintenance_window {
  	day_of_the_week  	= "Sunday"
   	time				= "09:00:00"
  }
  data_platform_version	= ` + DataPlatformVersion + `
}
`

const testAccCheckDataplatformClusterConfigUpdate = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fra"
  description = "Datacenter for testing Dataplatform Cluster"
}

resource ` + DataplatformClusterResource + ` ` + DataplatformClusterTestResource + ` {
  datacenter_id   		=  ` + DatacenterResource + `.datacenter_example.id
  name 					= "` + UpdatedResources + `"
  maintenance_window {
    day_of_the_week  	= "Saturday"
    time				= "10:00:00"
  }
  data_platform_version	= ` + DataPlatformVersion + `
}
`

const testAccDataSourceDataplatformClusterMatchById = testAccCheckDataplatformClusterConfigBasic + `
  data ` + DataplatformClusterResource + ` ` + DataplatformClusterTestDataSourceById + ` {
  id = ` + DataplatformClusterResource + `.` + DataplatformClusterTestResource + `.id
}
`

const testAccDataSourceDataplatformClusterMatchByName = testAccCheckDataplatformClusterConfigBasic + `
  data ` + DataplatformClusterResource + ` ` + DataplatformClusterTestDataSourceByName + ` {
  name = "` + DataplatformClusterTestResource + `"
}
`

const testAccDataSourceDataplatformClusterPartialMatchByName = testAccCheckDataplatformClusterConfigBasic + `
  data ` + DataplatformClusterResource + ` ` + DataplatformClusterTestDataSourceByName + ` {
  name = "test_"
  partial_match = true
}
`

const testAccDataSourceDataplatformClusterWrongNameError = testAccCheckDataplatformClusterConfigBasic + `
 data ` + DataplatformClusterResource + ` ` + DataplatformClusterTestDataSourceByName + ` {
  name = "wrong_name"
}
`

const testAccDataSourceDataplatformClusterWrongPartialNameError = testAccCheckDataplatformClusterConfigBasic + `
  data ` + DataplatformClusterResource + ` ` + DataplatformClusterTestDataSourceByName + ` {
  name = "wrong_name"
  partial_match = true
}
`
