//go:build compute || all

package ionoscloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func TestAccNSGBasic(t *testing.T) {
	var nsg ionoscloud.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckNSGDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNSGConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNSGExists(constant.NSGResource+"."+constant.NSGTestResource, &nsg),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "description", "testing-description"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "rule_ids.#", "0"),
				),
			},
			{
				Config: testAccCheckNSGDataSourceMatchId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNSGExists(constant.NSGResource+"."+constant.NSGTestResource, &nsg),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "description", "testing-description"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "rule_ids.#", "0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "name", constant.NSGResource+"."+constant.NSGTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "description", constant.NSGResource+"."+constant.NSGTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rule_ids.#", constant.NSGResource+"."+constant.NSGTestResource, "rule_ids.#"),
				),
			},
			{
				Config: testAccCheckNSGDataSourceMatchName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNSGExists(constant.NSGResource+"."+constant.NSGTestResource, &nsg),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "description", "testing-description"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "rule_ids.#", "0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NSGDataSourceByName, "id", constant.NSGResource+"."+constant.NSGTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NSGDataSourceByName, "description", constant.NSGResource+"."+constant.NSGTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NSGDataSourceByName, "rule_ids.#", constant.NSGResource+"."+constant.NSGTestResource, "rule_ids.#"),
				),
			},
			{
				Config: testAccCheckNSGConfigBasicUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNSGExists(constant.NSGResource+"."+constant.NSGTestResource, &nsg),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "name", "updated-name"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "description", "updated-description"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "rule_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccNSGFirewallRules(t *testing.T) {
	var nsg ionoscloud.SecurityGroup
	var rule1, rule2, rule3 *ionoscloud.FirewallRule
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckNSGRuleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNSGFirewallRulesBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNSGExists(constant.NSGResource+"."+constant.NSGTestResource, &nsg),
					testAccCheckNSGFirewallRuleExists(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", rule1),
					testAccCheckNSGFirewallRuleExists(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", rule2),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "description", "testing-description"),

					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "name", "SG Rule 1"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "target_ip", "22.231.113.66"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "icmp_type", "1"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "icmp_code", "8"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "type", "INGRESS"),

					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "name", "SG Rule 2"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "target_ip", "22.231.113.70"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "port_range_start", "10"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "port_range_end", "270"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckNSGFirewallRulesBasicAddRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNSGExists(constant.NSGResource+"."+constant.NSGTestResource, &nsg),
					testAccCheckNSGFirewallRuleExists(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", rule1),
					testAccCheckNSGFirewallRuleExists(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", rule2),
					testAccCheckNSGFirewallRuleExists(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", rule3),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "description", "testing-description"),

					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "name", "SG Rule 1"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "target_ip", "22.231.113.66"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "icmp_type", "1"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "icmp_code", "8"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "type", "INGRESS"),

					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "name", "SG Rule 2"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "target_ip", "22.231.113.70"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "port_range_start", "10"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "port_range_end", "270"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "type", "EGRESS"),

					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "name", "SG Rule 3"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "source_mac", "00:0a:95:9d:68:15"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "source_ip", "22.231.113.11"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "target_ip", "22.231.113.75"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckNSGFirewallRulesBasicUpdateRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNSGExists(constant.NSGResource+"."+constant.NSGTestResource, &nsg),
					testAccCheckNSGFirewallRuleExists(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", rule1),
					testAccCheckNSGFirewallRuleExists(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", rule2),
					testAccCheckNSGFirewallRuleExists(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", rule3),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "description", "testing-description"),

					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "name", "SG Rule 1"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "target_ip", "22.231.113.66"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "icmp_type", "1"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "icmp_code", "8"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "type", "INGRESS"),

					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "name", "SG Rule 2"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "target_ip", "22.231.113.70"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "port_range_start", "10"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "port_range_end", "270"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_2", "type", "EGRESS"),

					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "name", "SG Rule 3 Updated"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "protocol", "TCP"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "source_mac", "00:0a:95:9d:68:15"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "source_ip", "22.231.113.11"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "target_ip", "22.231.113.75"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "port_range_start", "800"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "port_range_end", "900"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_3", "type", "EGRESS"),
				),
			},
			{
				Config: testAccCheckNSGFirewallRulesBasicDeleteRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNSGExists(constant.NSGResource+"."+constant.NSGTestResource, &nsg),
					testAccCheckNSGFirewallRuleExists(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", rule1),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "description", "testing-description"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "name", "SG Rule 1"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "protocol", "ICMP"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "source_mac", "00:0a:95:9d:68:16"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "source_ip", "22.231.113.64"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "target_ip", "22.231.113.66"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "icmp_type", "1"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "icmp_code", "8"),
					resource.TestCheckResourceAttr(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "type", "INGRESS"),
				),
			},
			{
				Config: testAccCheckNSGFirewallRulesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNSGExists(constant.NSGResource+"."+constant.NSGTestResource, &nsg),
					testAccCheckNSGFirewallRuleExists(constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", rule1),
					resource.TestCheckResourceAttr(constant.NSGResource+"."+constant.NSGTestResource, "name", "testing-name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "name", constant.NSGResource+"."+constant.NSGTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "description", constant.NSGResource+"."+constant.NSGTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rule_ids.#", constant.NSGResource+"."+constant.NSGTestResource, "rule_ids.#"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rules.#", constant.NSGResource+"."+constant.NSGTestResource, "rule_ids.#"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rules.0.id", constant.NSGResource+"."+constant.NSGTestResource, "rule_ids.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rules.0.name", constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rules.0.protocol", constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rules.0.source_mac", constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "source_mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rules.0.source_ip", constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "source_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rules.0.target_ip", constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "target_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rules.0.icmp_type", constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "icmp_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rules.0.icmp_code", constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "icmp_code"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.NSGResource+"."+constant.NGDataSourceByID, "rules.0.type", constant.NSGFirewallRuleResource+"."+constant.NSGFirewallRuleTestResource+"_1", "type"),
				),
			},
		},
	})
}

func testAccCheckNSGDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.NSGResource {
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

func testAccCheckNSGExists(n string, nsg *ionoscloud.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckNSGExists: Not found: %s", n)
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
			return fmt.Errorf("error occured while fetching NSG: %s", rs.Primary.ID)
		}
		if *foundNSG.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		nsg = &foundNSG
		return nil
	}
}

func testAccCheckNSGRuleDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.NSGFirewallRuleResource {
			continue
		}
		_, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesFindById(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["nsg_id"], rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of the network security group rule %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("network security group rule %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckNSGFirewallRuleExists(n string, rule *ionoscloud.FirewallRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

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

const testAccCheckNSGConfigBasic = testAccCheckDatacenterConfigBasic + `
resource ` + constant.NSGResource + ` ` + constant.NSGTestResource + ` {
  name          = "testing-name"
  description   = "testing-description"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}
`

const testAccCheckNSGConfigBasicUpdated = testAccCheckDatacenterConfigBasic + `
resource ` + constant.NSGResource + ` ` + constant.NSGTestResource + ` {
  name          = "updated-name"
  description   = "updated-description"
  datacenter_id	= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}
`

const testAccCheckNSGDataSourceMatchId = testAccCheckNSGConfigBasic + `
data ` + constant.NSGResource + ` ` + constant.NGDataSourceByID + ` {
  id            = ` + constant.NSGResource + `.` + constant.NSGTestResource + `.id
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}
`

const testAccCheckNSGDataSourceMatchName = testAccCheckNSGConfigBasic + `
data ` + constant.NSGResource + ` ` + constant.NSGDataSourceByName + ` {
  name          = "testing-name"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}
`

const testAccCheckNSGFirewallRulesBasic = testAccCheckNSGConfigBasic + firewallRule1 + firewallRule2
const testAccCheckNSGFirewallRulesBasicAddRule = testAccCheckNSGFirewallRulesBasic + firewallRule3
const testAccCheckNSGFirewallRulesBasicUpdateRule = testAccCheckNSGFirewallRulesBasic + firewallRule3Updated
const testAccCheckNSGFirewallRulesBasicDeleteRules = testAccCheckNSGConfigBasic + firewallRule1
const testAccCheckNSGFirewallRulesDataSource = testAccCheckNSGFirewallRulesBasicDeleteRules + `
data ` + constant.NSGResource + ` ` + constant.NGDataSourceByID + ` {
  id            = ` + constant.NSGResource + `.` + constant.NSGTestResource + `.id
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}
`

const firewallRule1 = `
resource ` + constant.NSGFirewallRuleResource + ` ` + constant.NSGFirewallRuleTestResource + `_1` + ` {
  datacenter_id 	= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  nsg_id            = ` + constant.NSGResource + `.` + constant.NSGTestResource + `.id
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
resource ` + constant.NSGFirewallRuleResource + ` ` + constant.NSGFirewallRuleTestResource + `_2` + ` {
  datacenter_id 	= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  nsg_id            = ` + constant.NSGResource + `.` + constant.NSGTestResource + `.id
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
resource ` + constant.NSGFirewallRuleResource + ` ` + constant.NSGFirewallRuleTestResource + `_3` + ` {
  datacenter_id 	= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  nsg_id            = ` + constant.NSGResource + `.` + constant.NSGTestResource + `.id
  protocol          = "TCP"
  name              = "SG Rule 3"
  source_mac        = "00:0a:95:9d:68:15"
  source_ip         = "22.231.113.11"
  target_ip         = "22.231.113.75"
  type              = "EGRESS"
}
`
const firewallRule3Updated = `
resource ` + constant.NSGFirewallRuleResource + ` ` + constant.NSGFirewallRuleTestResource + `_3` + ` {
  datacenter_id 	= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  nsg_id            = ` + constant.NSGResource + `.` + constant.NSGTestResource + `.id	
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
