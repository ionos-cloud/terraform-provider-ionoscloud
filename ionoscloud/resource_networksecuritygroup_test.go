//go:build compute || all || nsg

package ionoscloud

import "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

func TestAccNetworkSecurityGroupBasic(t *testing.T) {
	var nsg ionoscloud.NetworkSecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkSecurityGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkSecurityGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkSecurityGroupExists(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, &nsg),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "name", constant.NetworkSecurityGroupTestResource),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "description", constant.NetworkSecurityGroupTestResource),
				),
			},
			{
				Config: testAccCheckNetworkSecurityGroupConfigBasicUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkSecurityGroupExists(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, &nsg),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "name", constant.NetworkSecurityGroupTestResource+"_updated"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "description", constant.NetworkSecurityGroupTestResource+"_updated"),
				),
			},
		},
	})
}

func TestAccNetworkSecurityGroupFirewallRules(t *testing.T) {
	var nsg ionoscloud.NetworkSecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkSecurityGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkSecurityGroupFirewallRulesBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkSecurityGroupExists(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, &nsg),
				),
			},
			{
				Config: testAccCheckNetworkSecurityGroupFirewallRulesBasicAddRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkSecurityGroupExists(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, &nsg),
				),
			},
			{
				Config: testAccCheckNetworkSecurityGroupFirewallRulesBasicUpdateRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkSecurityGroupExists(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, &nsg),
				),
			},
			{
				Config: testAccCheckNetworkSecurityGroupFirewallRulesBasicDeleteRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkSecurityGroupExists(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, &nsg),
				),
			},
		},
	})
}

func testAccCheckNetworkSecurityGroupDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.NetworkSecurityGroupResource {
			continue
		}

		_, apiResponse, err := client.GetToken(ctx, rs.Primary.Attributes["id"], rs.Primary.ID)

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of the network security group %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("network security group %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckNetworkSecurityGroupExists(n string, nsg *ionoscloud.NetworkSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckNetworkSecurityGroupExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}
		// todo - retrieve the nsg
		// foundServer, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, rs.Primary.Attributes["datacenter_id"],
		// 	rs.Primary.Attributes["server_id"], rs.Primary.Attributes["nic_id"], rs.Primary.ID).Execute()
		// logApiRequestTime(apiResponse)
		//
		// if err != nil {
		// 	return fmt.Errorf("error occured while fetching NetworkSecurityGroup rule: %s", rs.Primary.ID)
		// }
		// if *foundServer.Id != rs.Primary.ID {
		// 	return fmt.Errorf("record not found")
		// }
		//
		// firewall = &foundServer

		return nil
	}
}

const testAccCheckNetworkSecurityGroupConfigBasic = testAccCheckDatacenterConfigBasic + `
resource ` + constant.NetworkSecurityGroupResource + ` ` + constant.NetworkSecurityGroupTestResource + ` {
  name          = ` + constant.NetworkSecurityGroupTestResource + `
  description   = ` + constant.NetworkSecurityGroupTestResource + `
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}
`

const testAccCheckNetworkSecurityGroupConfigBasicUpdated = testAccCheckDatacenterConfigBasic + `
resource ` + constant.NetworkSecurityGroupResource + ` ` + constant.NetworkSecurityGroupTestResource + ` {
  name          = ` + constant.NetworkSecurityGroupTestResource + `_updated` + `
  description   = ` + constant.NetworkSecurityGroupTestResource + `_updated` + `
  datacenter_id	= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}
`

const testAccCheckNetworkSecurityGroupFirewallRulesBasic = testAccCheckNetworkSecurityGroupConfigBasic + firewallRule1 + firewallRule2
const testAccCheckNetworkSecurityGroupFirewallRulesBasicAddRule = testAccCheckNetworkSecurityGroupFirewallRulesBasic + firewallRule3
const testAccCheckNetworkSecurityGroupFirewallRulesBasicUpdateRule = testAccCheckNetworkSecurityGroupFirewallRulesBasic + firewallRule3_updated
const testAccCheckNetworkSecurityGroupFirewallRulesBasicDeleteRules = testAccCheckNetworkSecurityGroupConfigBasic + firewallRule1

const firewallRule1 = `
resource ` + constant.NetworkSecurityGroupFirewallRuleResource + ` ` + constant.NetworkSecurityGroupFirewallRuleTestResource + `_1` + ` {
  datacenter_id 	= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  nsg_id            = ` + constant.NetworkSecurityGroupResource + `.` + constant.NetworkSecurityGroupTestResource + `.id
  protocol          = "ICMP"
  name              = "SG Rule 1"
  source_mac        = "00:0a:95:9d:68:16"
  source_ip         = "22.231.113.64"
  target_ip         = "22.231.113.66"
  icmp_type         = 1
  icmp_code         = 8
  type              = "INGRESS"
}
`
const firewallRule2 = `
resource ` + constant.NetworkSecurityGroupFirewallRuleResource + ` ` + constant.NetworkSecurityGroupFirewallRuleTestResource + `_2` + ` {
  datacenter_id 	= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  nsg_id            = ` + constant.NetworkSecurityGroupResource + `.` + constant.NetworkSecurityGroupTestResource + `.id
  protocol          = "TCP"
  name              = "SG Rule 2"
  source_mac        = "00:0a:95:9d:68:16"
  source_ip         = "22.231.113.64"
  target_ip         = "22.231.113.70"
  port_range_start  = 10
  port_range_end    = 270
  type              = "EGRESS"
}
`
const firewallRule3 = `
resource ` + constant.NetworkSecurityGroupFirewallRuleResource + ` ` + constant.NetworkSecurityGroupFirewallRuleTestResource + `_3` + ` {
  datacenter_id 	= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  nsg_id            = ` + constant.NetworkSecurityGroupResource + `.` + constant.NetworkSecurityGroupTestResource + `.id
  protocol          = "TCP"
  name              = "SG Rule 3"
  source_mac        = "00:0a:95:9d:68:15"
  source_ip         = "22.231.113.11"
  target_ip         = "22.231.113.75"
  type              = "EGRESS"
}
`
const firewallRule3_updated = `
resource ` + constant.NetworkSecurityGroupFirewallRuleResource + ` ` + constant.NetworkSecurityGroupFirewallRuleTestResource + `_3` + ` {
  datacenter_id 	= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  nsg_id            = ` + constant.NetworkSecurityGroupResource + `.` + constant.NetworkSecurityGroupTestResource + `.id	
  protocol          = "TCP"
  name              = "SG Rule 3 Updated"
  source_mac        = "00:0a:95:9d:68:15"
  source_ip         = "22.231.113.11"
  target_ip         = "22.231.113.75"
  type              = "EGRESS"
  port_range_start  = 800
  port_range_end    = 900
}
`
