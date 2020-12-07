package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceLan_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLanCreateResources,
			},
			{
				Config: testAccDataSourceLanMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_lan.test_lan", "name", "test_ds_lan"),
				),
			},
		},
	})
}

func TestAccDataSourceLan_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLanCreateResources,
			},
			{
				Config: testAccDataSourceLanMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_lan.test_lan", "name", "test_ds_lan"),
				),
			},
		},
	})

}

const testAccDataSourceLanCreateResources = `
resource "ionoscloud_datacenter" "test_ds_lan" {
  name              = "test_datasource_lan"
  location          = "de/fra"
  description       = "datacenter for testing the lan terraform data source"
}
resource "ionoscloud_lan" "test_ds_lan" {
  depends_on        = [ionoscloud_datacenter.test_ds_lan]
  datacenter_id     = ionoscloud_datacenter.test_ds_lan.id
  name              = "test_ds_lan"
  public            = true
}
`

const testAccDataSourceLanMatchId = `
resource "ionoscloud_datacenter" "test_ds_lan" {
  name              = "test_datasource_lan"
  location          = "de/fra"
  description       = "datacenter for testing the lan terraform data source"
}

resource "ionoscloud_lan" "test_ds_lan" {
  depends_on        = [ionoscloud_datacenter.test_ds_lan]
  datacenter_id     = ionoscloud_datacenter.test_ds_lan.id
  name              = "test_ds_lan"
  public            = true
}

data "ionoscloud_lan" "test_lan" {
  datacenter_id = ionoscloud_datacenter.test_ds_lan.id
  id			= ionoscloud_lan.test_ds_lan.id
}
`

const testAccDataSourceLanMatchName = `
resource "ionoscloud_datacenter" "test_ds_lan" {
  name              = "test_datasource_lan"
  location          = "de/fra"
  description       = "datacenter for testing the lan terraform data source"
}

resource "ionoscloud_lan" "test_ds_lan" {
  depends_on        = [ionoscloud_datacenter.test_ds_lan]
  datacenter_id     = ionoscloud_datacenter.test_ds_lan.id
  name              = "test_ds_lan"
  public            = true
}

data "ionoscloud_lan" "test_lan" {
  datacenter_id = ionoscloud_datacenter.test_ds_lan.id
  name			= "test_ds_"
}
`
