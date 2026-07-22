package ionoscloud

import (
	"testing"
)

func TestParseInsecureEnv(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    bool
		wantErr bool
	}{
		{"true", "true", true, false},
		{"false", "false", false, false},
		{"TRUE", "TRUE", true, false},
		{"FALSE", "FALSE", false, false},
		{"1", "1", true, false},
		{"0", "0", false, false},
		{"invalid", "notabool", false, true},
		{"empty-string-not-called", "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val == "" {
				// parseInsecureEnv is only called when the env var is non-empty
				return
			}
			got, diags := parseInsecureEnv(tt.val)
			if tt.wantErr {
				if !diags.HasError() {
					t.Errorf("parseInsecureEnv(%q): expected error diagnostic, got none", tt.val)
				}
				return
			}
			if diags.HasError() {
				t.Errorf("parseInsecureEnv(%q): unexpected error: %v", tt.val, diags)
				return
			}
			if got != tt.want {
				t.Errorf("parseInsecureEnv(%q) = %v, want %v", tt.val, got, tt.want)
			}
		})
	}
}

// TestInsecureEnvFalseDoesNotEnableInsecure verifies that IONOS_ALLOW_INSECURE=false
// does not enable insecure mode — the pre-fix bug set insecure=true for any non-empty value.
func TestInsecureEnvFalseDoesNotEnableInsecure(t *testing.T) {
	val, diags := parseInsecureEnv("false")
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if val {
		t.Error("IONOS_ALLOW_INSECURE=false must not enable insecure mode")
	}
}
