package ionoscloud

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"regexp"
	"testing"
)

func TestAccDataSourceNic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCreateDataCenterAndServer,
			},
			{
				Config: testAccDataSourceNicMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "name", fullNicResourceName, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "dhcp", fullNicResourceName, "dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_active", fullNicResourceName, "firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_type", fullNicResourceName, "firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "mac", fullNicResourceName, "mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "pci_slot", fullNicResourceName, "pci_slot"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "lan", fullNicResourceName, "lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "ips", fullNicResourceName, "ips"),
				),
			},
			{
				Config: testAccDataSourceNicMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "name", fullNicResourceName, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "dhcp", fullNicResourceName, "dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_active", fullNicResourceName, "firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_type", fullNicResourceName, "firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "mac", fullNicResourceName, "mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "pci_slot", fullNicResourceName, "pci_slot"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "lan", fullNicResourceName, "lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "ips", fullNicResourceName, "ips"),
				),
			},
			{
				Config:      testAccDataSourceNicMatchNameError,
				ExpectError: regexp.MustCompile(`there are no nics that match the search criteria`),
			},
			{
				Config:      testAccDataSourceNicMatchIdAndNameError,
				ExpectError: regexp.MustCompile(`does not match expected name`),
			},
		},
	})
}

//unit test
func Test_dataSourceNicRead(t *testing.T) {
	id := "id"
	mac := "testMac"
	testName := "testname"
	dhcp := true
	firewallActive := true
	firewallType := "Bidirectional"
	nic := ionoscloud.Nic{
		Id:       &id,
		Type:     nil,
		Href:     nil,
		Metadata: nil,
		Properties: &ionoscloud.NicProperties{
			Name:           &testName,
			Mac:            &mac,
			Ips:            nil,
			Dhcp:           &dhcp,
			Lan:            nil,
			FirewallActive: &firewallActive,
			FirewallType:   &firewallType,
			DeviceNumber:   nil,
			PciSlot:        nil,
		},
		Entities: nil,
	}
	jsonNic, err := json.Marshal(nic)
	if err != nil {
		t.Fatalf("error marshalling nic %+v", nic)
	}
	var ctx = context.TODO()
	data := getEmptyTestResourceData(t, getNicDataSourceSchema())
	meta := getMockedClient(string(jsonNic))

	err = data.Set("datacenter_id", "testValueDatacenter")
	if err != nil {
		log.Println("error setting key on data ", err)
	}
	err = data.Set("server_id", "testValueServer")
	if err != nil {
		log.Println("error setting key on data ", err)
	}
	err = data.Set("id", "id")
	if err != nil {
		log.Println("error setting key on data ", err)
	}

	diag := dataSourceNicRead(ctx, data, meta)
	if diag != nil {
		t.Fatalf("error %+v", diag)
	}

	if *nic.Properties.Name != data.Get("name").(string) {
		t.Fatalf("expected '%s', got '%s'", *nic.Properties.Name, data.Get("name"))
	}
	if *nic.Properties.Mac != data.Get("mac").(string) {
		t.Fatalf("expected '%s', got '%s'", *nic.Properties.Mac, data.Get("mac"))
	}
	if *nic.Properties.Dhcp != data.Get("dhcp").(bool) {
		t.Fatalf("expected '%t', got '%s'", *nic.Properties.Dhcp, data.Get("dhcp"))
	}
	if *nic.Properties.FirewallActive != data.Get("firewall_active").(bool) {
		t.Fatalf("expected '%t', got '%s'", *nic.Properties.FirewallActive, data.Get("firewallActive"))
	}
	if *nic.Properties.FirewallType != data.Get("firewall_type").(string) {
		t.Fatalf("expected '%s', got '%s'", *nic.Properties.FirewallType, data.Get("firewallType"))
	}
}

const dataSourceNicById = nicResource + ".test_nic_data"

const testAccDataSourceNicMatchId = testAccCheckNicConfigBasic + `
data ` + NicResource + ` test_nic_data {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  id = ` + fullNicResourceName + `.id
}
`

const testAccDataSourceNicMatchName = testAccCheckNicConfigBasic +
	`data ` + NicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	name = ` + fullNicResourceName + `.name 
}`

const testAccDataSourceNicMatchNameError = testAccCheckNicConfigBasic +
	`data ` + NicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	name = "DoesNotExist"
}`

const testAccDataSourceNicMatchIdAndNameError = testAccCheckNicConfigBasic +
	`data ` + NicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	id = ` + fullNicResourceName + `.id
	name = "doesNotExist"
}`
