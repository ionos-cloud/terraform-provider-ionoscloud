//go:build all || dsaas
// +build all dsaas

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dsaas "github.com/ionos-cloud/sdk-go-autoscaling"
	"regexp"
	"testing"
)

func TestAccDSaaSClusterBasic(t *testing.T) {
	var DSaaSCluster dsaas.ClusterResponseData

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDSaaSClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDSaaSClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDSaaSClusterExists(DSaaSClusterResource+"."+DSaaSClusterTestResource, &DSaaSCluster),
					resource.TestCheckResourceAttr(DSaaSClusterResource+"."+DSaaSClusterTestResource, "name", DSaaSClusterTestResource),
					resource.TestCheckResourceAttr(DSaaSClusterResource+"."+DSaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(DSaaSClusterResource+"."+DSaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(DSaaSClusterResource+"."+DSaaSClusterTestResource, "data_platform_version", "1.1.0"),
				),
			},
			{
				Config: testAccDataSourceDSaaSClusterMatchById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceById, "name", DSaaSClusterResource+"."+DSaaSClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceById, "maintenance_window.0.time", DSaaSClusterResource+"."+DSaaSClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceById, "maintenance_window.0.day_of_the_week", DSaaSClusterResource+"."+DSaaSClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceById, "data_platform_version", DSaaSClusterResource+"."+DSaaSClusterTestResource, "data_platform_version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceById, "datacenter_id", DSaaSClusterResource+"."+DSaaSClusterTestResource, "datacenter_id"),
				),
			},
			{
				Config: testAccDataSourceDSaaSClusterMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceByName, "name", DSaaSClusterResource+"."+DSaaSClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceByName, "maintenance_window.0.time", DSaaSClusterResource+"."+DSaaSClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceByName, "maintenance_window.0.day_of_the_week", DSaaSClusterResource+"."+DSaaSClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceByName, "data_platform_version", DSaaSClusterResource+"."+DSaaSClusterTestResource, "data_platform_version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceByName, "datacenter_id", DSaaSClusterResource+"."+DSaaSClusterTestResource, "datacenter_id"),
				),
			},
			{
				Config: testAccDataSourceDSaaSClusterPartialMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceByName, "name", DSaaSClusterResource+"."+DSaaSClusterTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceByName, "maintenance_window.0.time", DSaaSClusterResource+"."+DSaaSClusterTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceByName, "maintenance_window.0.day_of_the_week", DSaaSClusterResource+"."+DSaaSClusterTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceByName, "data_platform_version", DSaaSClusterResource+"."+DSaaSClusterTestResource, "data_platform_version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DSaaSClusterResource+"."+DSaaSClusterTestDataSourceByName, "datacenter_id", DSaaSClusterResource+"."+DSaaSClusterTestResource, "datacenter_id"),
				),
			},
			{
				Config:      testAccDataSourceDSaaSClusterWrongNameError,
				ExpectError: regexp.MustCompile("no DSaaS Cluster found with the specified name"),
			},
			{
				Config:      testAccDataSourceDSaaSClusterWrongPartialNameError,
				ExpectError: regexp.MustCompile("no DSaaS Cluster found with the specified name"),
			},
			{
				Config: testAccCheckDSaaSClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDSaaSClusterExists(DSaaSClusterResource+"."+DSaaSClusterTestResource, &DSaaSCluster),
					resource.TestCheckResourceAttr(DSaaSClusterResource+"."+DSaaSClusterTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(DSaaSClusterResource+"."+DSaaSClusterTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(DSaaSClusterResource+"."+DSaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Saturday"),
					resource.TestCheckResourceAttr(DSaaSClusterResource+"."+DSaaSClusterTestResource, "data_platform_version", "1.1.0"),
				),
			},
		},
	})
}

func testAccCheckDSaaSClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).DSaaSClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != DSaaSClusterResource {
			continue
		}

		clusterId := rs.Primary.ID

		_, apiResponse, err := client.GetCluster(ctx, clusterId)

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of DSaaS cluster %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("DSaaS cluster %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckDSaaSClusterExists(n string, cluster *dsaas.ClusterResponseData) resource.TestCheckFunc {
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

		clusterId := rs.Primary.ID

		foundCluster, _, err := client.GetCluster(ctx, clusterId)

		if err != nil {
			return fmt.Errorf("an error occured while fetching DSaaS Cluster %s: %s", rs.Primary.ID, err)
		}
		if *foundCluster.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		cluster = &foundCluster

		return nil
	}
}

const testAccCheckDSaaSClusterConfigBasic = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fkb"
  description = "Datacenter for testing DSaaS Cluster"
}

resource ` + DSaaSClusterResource + ` ` + DSaaSClusterTestResource + ` {
	datacenter_id   		=  ` + DatacenterResource + `.datacenter_example.id
  	name 					= "` + DSaaSClusterTestResource + `"
  	maintenance_window {
    	day_of_the_week  	= "Sunday"
    	time				= "09:00:00"
  	}
  	data_platform_version	= "1.1.0"
}
`

const testAccCheckDSaaSClusterConfigUpdate = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/fkb"
  description = "Datacenter for testing DSaaS Cluster"
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
`

const testAccDataSourceDSaaSClusterMatchById = testAccCheckDSaaSClusterConfigBasic + `
data ` + DSaaSClusterResource + ` ` + DSaaSClusterTestDataSourceById + ` {
	id = ` + DSaaSClusterResource + `.` + DSaaSClusterTestResource + `.id
}
`

const testAccDataSourceDSaaSClusterMatchByName = testAccCheckDSaaSClusterConfigBasic + `
data ` + DSaaSClusterResource + ` ` + DSaaSClusterTestDataSourceByName + ` {
	name = "` + DSaaSClusterTestResource + `"
}
`

const testAccDataSourceDSaaSClusterPartialMatchByName = testAccCheckDSaaSClusterConfigBasic + `
data ` + DSaaSClusterResource + ` ` + DSaaSClusterTestDataSourceByName + ` {
	name = "test_"
    partial_match = true
}
`

const testAccDataSourceDSaaSClusterWrongNameError = testAccCheckDSaaSClusterConfigBasic + `
data ` + DSaaSClusterResource + ` ` + DSaaSClusterTestDataSourceByName + ` {
	name = "wrong_name"
}
`

const testAccDataSourceDSaaSClusterWrongPartialNameError = testAccCheckDSaaSClusterConfigBasic + `
data ` + DSaaSClusterResource + ` ` + DSaaSClusterTestDataSourceByName + ` {
	name = "wrong_name"
	partial_match = true
}
`
