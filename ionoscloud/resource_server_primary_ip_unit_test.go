package ionoscloud

import (
	"context"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// The enterprise server importer indexed the primary NIC's IPs unchecked and panicked on a NIC
// without IPs; it now shares setServerPrimaryIPFromNic with the VCPU importer.
func TestSetServerPrimaryIPFromNic(t *testing.T) {
	primaryID, otherID := "nic-primary", "nic-other"
	ips := []string{"198.51.100.10", "198.51.100.11"}
	withNics := func(items *[]ionoscloud.Nic) *ionoscloud.Server {
		return &ionoscloud.Server{Entities: &ionoscloud.ServerEntities{Nics: &ionoscloud.Nics{Items: items}}}
	}

	tests := []struct {
		name   string
		server *ionoscloud.Server
		wantIP string
	}{
		{name: "nil entities", server: &ionoscloud.Server{}},
		{name: "nil nics", server: &ionoscloud.Server{Entities: &ionoscloud.ServerEntities{}}},
		{name: "nil items", server: withNics(nil)},
		{name: "nic with nil id", server: withNics(&[]ionoscloud.Nic{{Properties: &ionoscloud.NicProperties{Ips: &ips}}})},
		{name: "primary nic with nil properties", server: withNics(&[]ionoscloud.Nic{{Id: &primaryID}})},
		{name: "primary nic with nil ips", server: withNics(&[]ionoscloud.Nic{{Id: &primaryID, Properties: &ionoscloud.NicProperties{}}})},
		{name: "primary nic with empty ips", server: withNics(&[]ionoscloud.Nic{{Id: &primaryID, Properties: &ionoscloud.NicProperties{Ips: &[]string{}}}})},
		{name: "only another nic has ips", server: withNics(&[]ionoscloud.Nic{{Id: &otherID, Properties: &ionoscloud.NicProperties{Ips: &ips}}})},
		{
			name: "first ip of the primary nic",
			server: withNics(&[]ionoscloud.Nic{
				{Id: &otherID, Properties: &ionoscloud.NicProperties{Ips: &[]string{"203.0.113.1"}}},
				{Id: &primaryID, Properties: &ionoscloud.NicProperties{Ips: &ips}},
			}),
			wantIP: "198.51.100.10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := resourceServer().Data(nil)
			if err := setServerPrimaryIPFromNic(context.Background(), d, tt.server, primaryID); err != nil {
				t.Fatalf("setServerPrimaryIPFromNic: %v", err)
			}
			if got := d.Get("primary_ip").(string); got != tt.wantIP {
				t.Errorf("primary_ip = %q, want %q", got, tt.wantIP)
			}
		})
	}
}
