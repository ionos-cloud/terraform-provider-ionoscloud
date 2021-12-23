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
					resource.TestCheckResourceAttrSet(IpfailoverResourceFullName, "id"),
					resource.TestCheckResourceAttrPair(IpfailoverResourceFullName, "id", DataSource+"."+ResourceIpFailover+"."+IpfailoverName, "id"),
					resource.TestCheckResourceAttrPair(IpfailoverResourceFullName, "nicuuid", DataSource+"."+ResourceIpFailover+"."+IpfailoverName, "nicuuid"),
					resource.TestCheckResourceAttrPair(IpfailoverResourceFullName, "lan_id", DataSource+"."+ResourceIpFailover+"."+IpfailoverName, "lan_id"),
					resource.TestCheckResourceAttrPair(IpfailoverResourceFullName, "datacenter_id", DataSource+"."+ResourceIpFailover+"."+IpfailoverName, "datacenter_id"),
				),
			},
		},
	})
}

var testAccDataSourceIpFailoverConfigBasic = testaccchecklanipfailoverconfigBasic + `
data ` + ResourceIpFailover + " " + IpfailoverName + `{
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  id		    = ` + IpfailoverResourceFullName + `.id
}
`
