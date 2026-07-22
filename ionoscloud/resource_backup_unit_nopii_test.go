package ionoscloud

import (
	"os"
	"regexp"
	"testing"
)

// TestBackupUnitUpdateEmailNotLoggedAsPII verifies that the backup unit update
// path does not pass raw email addresses as structured fields to tflog.
// This is a source-level guard: if PII logging is re-introduced, this test fails.
func TestBackupUnitUpdateEmailNotLoggedAsPII(t *testing.T) {
	src, err := os.ReadFile("resource_backup_unit.go")
	if err != nil {
		t.Fatalf("could not read resource_backup_unit.go: %v", err)
	}

	// Pattern matches tflog calls that pass oldEmail or newEmail as field values.
	// Any of these indicate PII is being logged.
	piiPatterns := []*regexp.Regexp{
		regexp.MustCompile(`tflog\.[A-Za-z]+\([^)]*"old"\s*:\s*oldEmail`),
		regexp.MustCompile(`tflog\.[A-Za-z]+\([^)]*"new"\s*:\s*newEmail`),
		regexp.MustCompile(`tflog\.[A-Za-z]+\([^)]*oldEmail`),
		regexp.MustCompile(`tflog\.[A-Za-z]+\([^)]*newEmail`),
	}
	for _, pat := range piiPatterns {
		if pat.Match(src) {
			t.Errorf("resource_backup_unit.go passes email PII to tflog (matched %s) — remove email values from log fields", pat)
		}
	}
}
