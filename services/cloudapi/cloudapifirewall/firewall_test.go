package cloudapifirewall

import (
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
)

func TestExtractOrderedFirewallIDs(t *testing.T) {
	t.Run("matches when API returns icmpCode/icmpType as explicit null for an unset rule", func(t *testing.T) {
		// This mirrors what DecodeInterfaceToStruct produces for a rule with no icmp_code/icmp_type
		// in config: the fields are left unset (isSet=false), not explicitly nulled.
		sentProps := ionoscloud.FirewallruleProperties{
			Name:           new("SSH"),
			Protocol:       new("TCP"),
			SourceMac:      *ionoscloud.NewNullableString(new("00:0a:95:9d:68:17")),
			SourceIp:       *ionoscloud.NewNullableString(new("87.106.83.92")),
			TargetIp:       *ionoscloud.NewNullableString(new("217.154.64.29")),
			PortRangeStart: new(int32(22)),
			PortRangeEnd:   new(int32(22)),
			Type:           new("EGRESS"),
		}

		foundProps := sentProps
		foundProps.IpVersion = *ionoscloud.NewNullableString(new("IPv4"))
		// API sends these back as an explicit JSON null, which the SDK unmarshals as isSet=true,
		// value=nil - not the same zero value as a field that was never set.
		foundProps.IcmpCode = ionoscloud.NullableInt32{}
		foundProps.IcmpCode.Set(nil)
		foundProps.IcmpType = ionoscloud.NullableInt32{}
		foundProps.IcmpType.Set(nil)

		id := "found-id-1"
		sent := ionoscloud.FirewallRule{Properties: sentProps}
		found := ionoscloud.FirewallRule{Id: &id, Properties: foundProps}

		ids := ExtractOrderedFirewallIDs([]ionoscloud.FirewallRule{found}, []ionoscloud.FirewallRule{sent})
		if len(ids) != 1 || ids[0] != id {
			t.Fatalf("expected [%s], got %v", id, ids)
		}
	})

	t.Run("matches when API returns sourceMac/sourceIp/targetIp as explicit null for an unset rule", func(t *testing.T) {
		sentProps := ionoscloud.FirewallruleProperties{
			Name:           new("SSH"),
			Protocol:       new("TCP"),
			PortRangeStart: new(int32(22)),
			PortRangeEnd:   new(int32(22)),
			Type:           new("EGRESS"),
		}

		foundProps := sentProps
		foundProps.IpVersion = *ionoscloud.NewNullableString(new("IPv4"))
		foundProps.SourceMac = ionoscloud.NullableString{}
		foundProps.SourceMac.Set(nil)
		foundProps.SourceIp = ionoscloud.NullableString{}
		foundProps.SourceIp.Set(nil)
		foundProps.TargetIp = ionoscloud.NullableString{}
		foundProps.TargetIp.Set(nil)

		id := "found-id-2"
		sent := ionoscloud.FirewallRule{Properties: sentProps}
		found := ionoscloud.FirewallRule{Id: &id, Properties: foundProps}

		ids := ExtractOrderedFirewallIDs([]ionoscloud.FirewallRule{found}, []ionoscloud.FirewallRule{sent})
		if len(ids) != 1 || ids[0] != id {
			t.Fatalf("expected [%s], got %v", id, ids)
		}
	})

	t.Run("does not match when a real property differs", func(t *testing.T) {
		sentProps := ionoscloud.FirewallruleProperties{
			Name:     new("SSH"),
			Protocol: new("TCP"),
		}
		foundProps := ionoscloud.FirewallruleProperties{
			Name:     new("other"),
			Protocol: new("TCP"),
		}
		id := "found-id-1"
		sent := ionoscloud.FirewallRule{Properties: sentProps}
		found := ionoscloud.FirewallRule{Id: &id, Properties: foundProps}

		ids := ExtractOrderedFirewallIDs([]ionoscloud.FirewallRule{found}, []ionoscloud.FirewallRule{sent})
		if len(ids) != 0 {
			t.Fatalf("expected no match, got %v", ids)
		}
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		if ids := ExtractOrderedFirewallIDs(nil, nil); len(ids) != 0 {
			t.Fatalf("expected empty, got %v", ids)
		}
	})
}
