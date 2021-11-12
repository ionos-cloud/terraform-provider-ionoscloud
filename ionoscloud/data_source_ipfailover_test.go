package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceIpFailoverMatchFields(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIpFailoverConfigBasic,
			},
			{
				Config: testAccDataSourceIpFailoverConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ipfailoverResourceFullName, "id"),
					resource.TestCheckResourceAttrPair(ipfailoverResourceFullName, "id", DataSource+"."+ResourceIpFailover+"."+ipfailoverName, "id"),
					resource.TestCheckResourceAttrPair(ipfailoverResourceFullName, "nicuuid", DataSource+"."+ResourceIpFailover+"."+ipfailoverName, "nicuuid"),
					resource.TestCheckResourceAttrPair(ipfailoverResourceFullName, "lan_id", DataSource+"."+ResourceIpFailover+"."+ipfailoverName, "lan_id"),
					resource.TestCheckResourceAttrPair(ipfailoverResourceFullName, "datacenter_id", DataSource+"."+ResourceIpFailover+"."+ipfailoverName, "datacenter_id"),
				),
			},
		},
	})
}

var testAccDataSourceIpFailoverConfigBasic = testAccCheckLanIPFailoverConfig + `
data ` + ResourceIpFailover + " " + ipfailoverName + `{
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  id		    = ` + ipfailoverResourceFullName + `.id
}
`
