package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// A block in the state is cleared when the API returns none, e.g. after IPv6 is turned off outside Terraform.
func TestSetLanDataClearsCidrBlocks(t *testing.T) {
	d := resourceLan().Data(&terraform.InstanceState{
		ID: "1",
		Attributes: map[string]string{
			"ipv4_cidr_block": "10.5.0.0/24",
			"ipv6_cidr_block": "2a01:239:400:f100::/64",
		},
	})

	lan := &ionoscloud.Lan{Id: new("1"), Properties: &ionoscloud.LanProperties{Public: new(false)}}
	if err := setLanData(d, lan); err != nil {
		t.Fatalf("setLanData = %v, want nil", err)
	}

	for _, k := range []string{"ipv4_cidr_block", "ipv6_cidr_block"} {
		if v := d.Get(k).(string); v != "" {
			t.Errorf("%s = %q, want empty", k, v)
		}
	}
}
