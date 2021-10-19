package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLan(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLanCreateResources,
			},
			{
				Config: testAccDataSourceLanMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_lan."+LanDataSourceById, "name", "ionoscloud_lan."+LanTestResource, "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_lan."+LanDataSourceById, "ip_failover.nic_uuid", "ionoscloud_lan."+LanTestResource, "ip_failover.nic_uuid"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_lan."+LanDataSourceById, "ip_failover.ip", "ionoscloud_lan."+LanTestResource, "ip_failover.ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_lan."+LanDataSourceById, "pcc", "ionoscloud_lan."+LanTestResource, "pcc"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_lan."+LanDataSourceById, "public", "ionoscloud_lan."+LanTestResource, "public"),
				),
			},
			{
				Config: testAccDataSourceLanMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_lan."+LanDataSourceByName, "name", "ionoscloud_lan."+LanTestResource, "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_lan."+LanDataSourceByName, "ip_failover.nic_uuid", "ionoscloud_lan."+LanTestResource, "ip_failover.nic_uuid"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_lan."+LanDataSourceByName, "ip_failover.ip", "ionoscloud_lan."+LanTestResource, "ip_failover.ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_lan."+LanDataSourceByName, "pcc", "ionoscloud_lan."+LanTestResource, "pcc"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_lan."+LanDataSourceByName, "public", "ionoscloud_lan."+LanTestResource, "public"),
				),
			},
		},
	})
}

const testAccDataSourceLanCreateResources = `
resource "ionoscloud_datacenter" "foobar" {
  name              = "test_datasource_lan"
  location          = "de/fra"
  description       = "datacenter for testing the lan terraform data source"
}
resource "ionoscloud_private_crossconnect" "example" {
  name        = "example"
  description = "example description"
}
resource "ionoscloud_lan" ` + LanTestResource + ` {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = false
  name = "` + LanTestResource + `"
  pcc = ionoscloud_private_crossconnect.example.id
}
`

const testAccDataSourceLanMatchId = `
resource "ionoscloud_datacenter" "foobar" {
  name              = "test_datasource_lan"
  location          = "de/fra"
  description       = "datacenter for testing the lan terraform data source"
}
resource "ionoscloud_private_crossconnect" "example" {
  name        = "example"
  description = "example description"
}
resource "ionoscloud_lan" ` + LanTestResource + ` {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = false
  name = "` + LanTestResource + `"
  pcc = ionoscloud_private_crossconnect.example.id
}
data "ionoscloud_lan" ` + LanDataSourceById + ` {
  datacenter_id = ionoscloud_datacenter.foobar.id
  id			= ionoscloud_lan.` + LanTestResource + `.id
}
`

const testAccDataSourceLanMatchName = `
resource "ionoscloud_datacenter" "foobar" {
  name              = "test_datasource_lan"
  location          = "de/fra"
  description       = "datacenter for testing the lan terraform data source"
}
resource "ionoscloud_private_crossconnect" "example" {
  name        = "example"
  description = "example description"
}
resource "ionoscloud_lan" ` + LanTestResource + ` {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = false
  name = "` + LanTestResource + `"
  pcc = ionoscloud_private_crossconnect.example.id
}
data "ionoscloud_lan" ` + LanDataSourceByName + ` {
  datacenter_id = ionoscloud_datacenter.foobar.id
  name			= "` + LanTestResource + `"
}
`
