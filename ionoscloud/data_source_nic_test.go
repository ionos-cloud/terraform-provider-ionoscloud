//go:build compute || all || nic

package ionoscloud

import (
	"context"
	"encoding/json"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"testing"
)

// unit test
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
