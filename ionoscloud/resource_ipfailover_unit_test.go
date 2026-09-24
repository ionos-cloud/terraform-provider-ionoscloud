package ionoscloud

import "testing"

const (
	testFailoverIPv4     = "203.0.113.10"
	testFailoverIPv6     = "2001:db8::10"
	testFailoverLocation = "de/fra"
)

func TestParseIPFailoverImportID(t *testing.T) {
	tests := []struct {
		name                                  string
		importID                              string
		wantLocation, wantDC, wantLan, wantIP string
		wantErr                               bool
	}{
		{name: "ipv4", importID: "dc/1/203.0.113.10", wantDC: "dc", wantLan: "1", wantIP: testFailoverIPv4},
		{name: "ipv4 with location", importID: "de/fra:dc/1/203.0.113.10", wantLocation: testFailoverLocation, wantDC: "dc", wantLan: "1", wantIP: testFailoverIPv4},
		{name: "ipv4 with sub-location", importID: "de/fra/2:dc/1/203.0.113.10", wantLocation: "de/fra/2", wantDC: "dc", wantLan: "1", wantIP: testFailoverIPv4},
		{name: "ipv6", importID: "dc/1/2001:db8::10", wantDC: "dc", wantLan: "1", wantIP: testFailoverIPv6},
		{name: "ipv6 with location", importID: "es/vit:dc/1/2001:db8::10", wantLocation: "es/vit", wantDC: "dc", wantLan: "1", wantIP: testFailoverIPv6},
		{name: "ipv6 fully expanded", importID: "dc/2/2001:0db8:0000:0000:0000:0000:0000:0010", wantDC: "dc", wantLan: "2", wantIP: "2001:0db8:0000:0000:0000:0000:0000:0010"},
		{name: "missing ip", importID: "dc/1", wantErr: true},
		{name: "too many parts", importID: "dc/1/203.0.113.10/extra", wantErr: true},
		{name: "empty lan", importID: "dc//203.0.113.10", wantErr: true},
		{name: "location only", importID: "de/fra:", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, dcID, lanID, ip, err := parseIPFailoverImportID(tt.importID)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseIPFailoverImportID(%q) error = nil, want an error", tt.importID)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseIPFailoverImportID(%q) unexpected error: %v", tt.importID, err)
			}
			if location != tt.wantLocation || dcID != tt.wantDC || lanID != tt.wantLan || ip != tt.wantIP {
				t.Errorf("parseIPFailoverImportID(%q) = (%q, %q, %q, %q), want (%q, %q, %q, %q)", tt.importID,
					location, dcID, lanID, ip, tt.wantLocation, tt.wantDC, tt.wantLan, tt.wantIP)
			}
		})
	}
}
