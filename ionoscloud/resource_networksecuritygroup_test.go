package ionoscloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func TestAccNetworkSecurityGroupBasic(t *testing.T) {
	var nsg ionoscloud.SecurityGroup

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
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "description", "testing-description"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "rule_ids.#", "0"),
				),
			},
			{
				Config: testAccCheckNetworkSecurityGroupConfigBasicUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkSecurityGroupExists(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, &nsg),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "name", "updated-name"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "description", "updated-description"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "rule_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccNetworkSecurityGroupFirewallRules(t *testing.T) {
	var nsg ionoscloud.SecurityGroup
	var rule1, rule2, rule3 *ionoscloud.FirewallRule
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
					testAccCheckNSGFirewallRuleExists(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", rule1),
					testAccCheckNSGFirewallRuleExists(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", rule2),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "description", "testing-description"),

					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "name", "SG Rule 1"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "target_ip", "22.231.113.66"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "icmp_type", "1"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "icmp_code", "8"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "type", "INGRESS"),

					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "name", "SG Rule 2"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "target_ip", "22.231.113.70"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "port_range_start", "10"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "port_range_end", "270"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckNetworkSecurityGroupFirewallRulesBasicAddRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkSecurityGroupExists(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, &nsg),
					testAccCheckNSGFirewallRuleExists(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", rule1),
					testAccCheckNSGFirewallRuleExists(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", rule2),
					testAccCheckNSGFirewallRuleExists(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", rule3),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "description", "testing-description"),

					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "name", "SG Rule 1"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "target_ip", "22.231.113.66"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "icmp_type", "1"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "icmp_code", "8"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "type", "INGRESS"),

					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "name", "SG Rule 2"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "target_ip", "22.231.113.70"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "port_range_start", "10"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "port_range_end", "270"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "type", "EGRESS"),

					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "name", "SG Rule 3"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "source_mac", "00:0a:95:9d:68:15"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "source_ip", "22.231.113.11"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "target_ip", "22.231.113.75"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckNetworkSecurityGroupFirewallRulesBasicUpdateRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkSecurityGroupExists(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, &nsg),
					testAccCheckNSGFirewallRuleExists(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", rule1),
					testAccCheckNSGFirewallRuleExists(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", rule2),
					testAccCheckNSGFirewallRuleExists(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", rule3),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "description", "testing-description"),

					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "name", "SG Rule 1"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "target_ip", "22.231.113.66"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "icmp_type", "1"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "icmp_code", "8"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "type", "INGRESS"),

					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "name", "SG Rule 2"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "target_ip", "22.231.113.70"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "port_range_start", "10"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "port_range_end", "270"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_2", "type", "EGRESS"),

					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "name", "SG Rule 3 Updated"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "source_mac", "00:0a:95:9d:68:15"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "source_ip", "22.231.113.11"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "target_ip", "22.231.113.75"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "port_range_start", "800"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "port_range_end", "900"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_3", "type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckNetworkSecurityGroupFirewallRulesBasicDeleteRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkSecurityGroupExists(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, &nsg),
					testAccCheckNSGFirewallRuleExists(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", rule1),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "description", "testing-description"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "name", "SG Rule 1"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "target_ip", "22.231.113.66"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "icmp_type", "1"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "icmp_code", "8"),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", "type", "INGRESS"),
					// The NSG Resource is not refreshed since the changes are happening in the firewall rules resources so at this point the 'rule_ids' list still has 3 entries
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "rule_ids.#", "3"),
				),
			},
			{
				Config: testAccCheckNetworkSecurityGroupFirewallRulesBasicDeleteRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkSecurityGroupExists(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, &nsg),
					testAccCheckNSGFirewallRuleExists(constant.NetworkSecurityGroupFirewallRuleResource+"."+constant.NetworkSecurityGroupFirewallRuleTestResource+"_1", rule1),
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "name", "testing-name"),
					// New test step refreshes the state so the 'rule_ids' list reflects the correct value
					resource.TestCheckResourceAttr(constant.NetworkSecurityGroupResource+"."+constant.NetworkSecurityGroupTestResource, "rule_ids.#", "1"),
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
		_, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

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

func testAccCheckNetworkSecurityGroupExists(n string, nsg *ionoscloud.SecurityGroup) resource.TestCheckFunc {
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

		foundNSG, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("error occured while fetching NetworkSecurityGroup: %s", rs.Primary.ID)
		}
		if *foundNSG.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		nsg = &foundNSG
		return nil
	}
}

func testAccCheckNSGFirewallRuleExists(n string, rule *ionoscloud.FirewallRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckNSGFirewallRuleExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no NSG firewall rule ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundRule, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["nsg_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("error occured while fetching NSG firewall rule: %s", rs.Primary.ID)
		}
		if *foundRule.Id != rs.Primary.ID {
			return fmt.Errorf("nsg firewall rule not found")
		}
		rule = &foundRule
		return nil
	}
}

const testAccCheckNetworkSecurityGroupConfigBasic = testAccCheckDatacenterConfigBasic + `
resource ` + constant.NetworkSecurityGroupResource + ` ` + constant.NetworkSecurityGroupTestResource + ` {
  name          = "testing-name"
  description   = "testing-description"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}
`

const testAccCheckNetworkSecurityGroupConfigBasicUpdated = testAccCheckDatacenterConfigBasic + `
resource ` + constant.NetworkSecurityGroupResource + ` ` + constant.NetworkSecurityGroupTestResource + ` {
  name          = "updated-name"
  description   = "updated-description"
  datacenter_id	= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}
`

const testAccCheckNetworkSecurityGroupFirewallRulesBasic = testAccCheckNetworkSecurityGroupConfigBasic + firewallRule1 + firewallRule2
const testAccCheckNetworkSecurityGroupFirewallRulesBasicAddRule = testAccCheckNetworkSecurityGroupFirewallRulesBasic + firewallRule3
const testAccCheckNetworkSecurityGroupFirewallRulesBasicUpdateRule = testAccCheckNetworkSecurityGroupFirewallRulesBasic + firewallRule3Updated
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
const firewallRule3Updated = `
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
