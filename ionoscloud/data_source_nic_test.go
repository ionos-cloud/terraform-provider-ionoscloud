//go:build compute || all || nic

package ionoscloud

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/cloud/v2"
)

func testGetNicDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"server_id": {
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"datacenter_id": {
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"lan": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"dhcp": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"dhcpv6": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"ipv6_cidr_block": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"ips": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
			Optional: true,
		},
		"ipv6_ips": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Computed: true,
		},
		"firewall_active": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"firewall_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"mac": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"device_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"pci_slot": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

// unit test
func Test_dataSourceNicRead(t *testing.T) {
	id := "id"
	mac := "testMac"
	testName := "testname"
	dhcp := true
	dhcpv6 := false
	dhcpv6NulBool := ionoscloud.NullableBool{}
	dhcpv6NulBool.Set(&dhcpv6)
	firewallActive := true
	firewallType := "Bidirectional"
	ipv6CidrBlock := "AUTO"
	ipv6CidrBlockNulStr := ionoscloud.NullableString{}
	ipv6CidrBlockNulStr.Set(&ipv6CidrBlock)
	nic := ionoscloud.Nic{
		Id:       &id,
		Type:     nil,
		Href:     nil,
		Metadata: nil,
		Properties: ionoscloud.NicProperties{
			Name:           &testName,
			Mac:            &mac,
			Ips:            nil,
			Ipv6Ips:        nil,
			Dhcp:           &dhcp,
			Dhcpv6:         dhcpv6NulBool,
			Ipv6CidrBlock:  ipv6CidrBlockNulStr,
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
	data := getEmptyTestResourceData(t, testGetNicDataSourceSchema())
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
	if *nic.Properties.Dhcpv6.Get() != data.Get("dhcpv6").(bool) {
		t.Fatalf("expected '%t', got '%s'", *nic.Properties.Dhcpv6.Get(), data.Get("dhcpv6"))
	}
	if *nic.Properties.Ipv6CidrBlock.Get() != data.Get("ipv6_cidr_block").(string) {
		t.Fatalf("expected '%s', got '%s'", *nic.Properties.Ipv6CidrBlock.Get(), data.Get("ipv6CidrBlock"))
	}
	if *nic.Properties.FirewallActive != data.Get("firewall_active").(bool) {
		t.Fatalf("expected '%t', got '%s'", *nic.Properties.FirewallActive, data.Get("firewallActive"))
	}
	if *nic.Properties.FirewallType != data.Get("firewall_type").(string) {
		t.Fatalf("expected '%s', got '%s'", *nic.Properties.FirewallType, data.Get("firewallType"))
	}
}
