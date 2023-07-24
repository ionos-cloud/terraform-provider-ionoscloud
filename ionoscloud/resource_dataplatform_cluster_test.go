//go:build all || dataplatform

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
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
					testAccCheckDataplatformClusterExists(constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, &DataplatformCluster),
					resource.TestCheckResourceAttr(constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "name", constant.DataplatformClusterTestResource),
					resource.TestCheckResourceAttr(constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "version", constant.DataPlatformVersion),
				),
			},
			{
				Config: testAccDataSourceDataplatformClusterMatchById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceById, "name", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceById, "maintenance_window.0.time", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceById, "maintenance_window.0.day_of_the_week", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceById, "version", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceById, "datacenter_id", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "datacenter_id"),
				),
			},
			{
				Config: testAccDataSourceDataplatformClusterMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceByName, "name", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceByName, "maintenance_window.0.time", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceByName, "maintenance_window.0.day_of_the_week", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceByName, "version", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceByName, "datacenter_id", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "datacenter_id"),
				),
			},
			{
				Config: testAccDataSourceDataplatformClusterPartialMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceByName, "name", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceByName, "maintenance_window.0.time", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceByName, "maintenance_window.0.day_of_the_week", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceByName, "version", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestDataSourceByName, "datacenter_id", constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "datacenter_id"),
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
					testAccCheckDataplatformClusterExists(constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, &DataplatformCluster),
					resource.TestCheckResourceAttr(constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "maintenance_window.0.day_of_the_week", "Saturday"),
					resource.TestCheckResourceAttr(constant.DataplatformClusterResource+"."+constant.DataplatformClusterTestResource, "version", constant.DataPlatformVersion),
				),
			},
		},
	})
}

func testAccCheckDataplatformClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).DataplatformClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DataplatformClusterResource {
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
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fra"
  description = "Datacenter for testing Dataplatform Cluster"
}

resource ` + constant.DataplatformClusterResource + ` ` + constant.DataplatformClusterTestResource + ` {
  datacenter_id   		=  ` + constant.DatacenterResource + `.datacenter_example.id
  name 					= "` + constant.DataplatformClusterTestResource + `"
  maintenance_window {
  	day_of_the_week  	= "Sunday"
   	time				= "09:00:00"
  }
  version	= ` + constant.DataPlatformVersion + `
}
`

const testAccCheckDataplatformClusterConfigUpdate = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fra"
  description = "Datacenter for testing Dataplatform Cluster"
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
`

const testAccDataSourceDataplatformClusterMatchById = testAccCheckDataplatformClusterConfigBasic + `
  data ` + constant.DataplatformClusterResource + ` ` + constant.DataplatformClusterTestDataSourceById + ` {
  id = ` + constant.DataplatformClusterResource + `.` + constant.DataplatformClusterTestResource + `.id
}
`

const testAccDataSourceDataplatformClusterMatchByName = testAccCheckDataplatformClusterConfigBasic + `
  data ` + constant.DataplatformClusterResource + ` ` + constant.DataplatformClusterTestDataSourceByName + ` {
  name = "` + constant.DataplatformClusterTestResource + `"
}
`

const testAccDataSourceDataplatformClusterPartialMatchByName = testAccCheckDataplatformClusterConfigBasic + `
  data ` + constant.DataplatformClusterResource + ` ` + constant.DataplatformClusterTestDataSourceByName + ` {
  name = "test_"
  partial_match = true
}
`

const testAccDataSourceDataplatformClusterWrongNameError = testAccCheckDataplatformClusterConfigBasic + `
 data ` + constant.DataplatformClusterResource + ` ` + constant.DataplatformClusterTestDataSourceByName + ` {
  name = "wrong_name"
}
`

const testAccDataSourceDataplatformClusterWrongPartialNameError = testAccCheckDataplatformClusterConfigBasic + `
  data ` + constant.DataplatformClusterResource + ` ` + constant.DataplatformClusterTestDataSourceByName + ` {
  name = "wrong_name"
  partial_match = true
}
`
