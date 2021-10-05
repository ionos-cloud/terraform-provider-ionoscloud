package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDatacenter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{

				Config: testAccDatasourceDatacenter,
			},
			{
				Config: testAccDataSourceDatacenterMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_matching", "name", "ionoscloud_datacenter.foobar", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_matching", "location", "ionoscloud_datacenter.foobar", "location"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_matching", "description", "ionoscloud_datacenter.foobar", "description"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_matching", "version", "ionoscloud_datacenter.foobar", "version"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_matching", "features", "ionoscloud_datacenter.foobar", "features"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_matching", "sec_auth_protection", "ionoscloud_datacenter.foobar", "sec_auth_protection"),
				),
			},
			{
				Config: testAccDataSourceDatacenterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_id", "name", "ionoscloud_datacenter.foobar", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_id", "location", "ionoscloud_datacenter.foobar", "location"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_id", "description", "ionoscloud_datacenter.foobar", "description"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_id", "version", "ionoscloud_datacenter.foobar", "version"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_id", "features", "ionoscloud_datacenter.foobar", "features"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_datacenter.test_id", "sec_auth_protection", "ionoscloud_datacenter.foobar", "sec_auth_protection"),
				),
			},
		},
	})

}

const testAccDatasourceDatacenter = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "test_name"
	location = "us/las"
	description = "Test Datacenter Description"
	sec_auth_protection = false
}`

const testAccDataSourceDatacenterMatching = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "test_name"
	location = "us/las"
	description = "Test Datacenter Description"
	sec_auth_protection = false
}

data "ionoscloud_datacenter" "test_matching" {
    name = ionoscloud_datacenter.foobar.name
    location = "us/las"
}`

const testAccDataSourceDatacenterMatchId = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "test_name"
	location = "us/las"
	description = "Test Datacenter Description"
	sec_auth_protection = false
}

data "ionoscloud_datacenter" "test_id" {
    id = ionoscloud_datacenter.foobar.id
}`
