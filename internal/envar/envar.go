package envar

import (
	"fmt"
	"os"
	"testing"
)

const (
	// IonosToken is the environment variable name for the Ionos Cloud API token.
	IonosToken = "IONOS_TOKEN"
	// IonosUsername is the environment variable name for the Ionos Cloud API username.
	IonosUsername = "IONOS_USERNAME"
	// IonosPassword is the environment variable name for the Ionos Cloud API password.
	IonosPassword = "IONOS_PASSWORD"
	// IonosS3AccessKey is the environment variable name for the IONOS Object Storage access key.
	IonosS3AccessKey = "IONOS_S3_ACCESS_KEY"
	// IonosS3SecretKey is the environment variable name for the IONOS Object Storage secret key.
	IonosS3SecretKey = "IONOS_S3_SECRET_KEY"
)

// RequireOneOf verifies that at least one environment variable is non-empty or returns an error.
func RequireOneOf(names []string, usageMessage string) (string, string, error) {
	for _, variable := range names {
		value := os.Getenv(variable)

		if value != "" {
			return variable, value, nil
		}
	}

	return "", "", fmt.Errorf("at least one environment variable of %v must be set. Usage: %s", names, usageMessage)
}

// Require verifies that an environment variable is non-empty or returns an error.
func Require(name string, usageMessage string) (string, error) {
	value := os.Getenv(name)

	if value == "" {
		return "", fmt.Errorf("environment variable %s must be set. Usage: %s", name, usageMessage)
	}

	return value, nil
}

// FailIfAllEmpty verifies that at least one environment variable is non-empty or fails the test.
//
// If at least one environment variable is non-empty, returns the first name and value.
func FailIfAllEmpty(t *testing.T, names []string, usageMessage string) (string, string) {
	t.Helper()

	name, value, err := RequireOneOf(names, usageMessage)
	if err != nil {
		t.Fatal(err)
		return "", ""
	}

	return name, value
}

// FailIfEmpty verifies that an environment variable is non-empty or fails the test.
//
// For acceptance tests, this function must be used outside PreCheck functions to set values for configurations.
func FailIfEmpty(t *testing.T, name string, usageMessage string) string {
	t.Helper()

	value := os.Getenv(name)

	if value == "" {
		t.Fatalf("environment variable %s must be set. Usage: %s", name, usageMessage)
	}

	return value
}

// SkipIfEmpty verifies that an environment variable is non-empty or skips the test.
//
// For acceptance tests, this function must be used outside PreCheck functions to set values for configurations.
func SkipIfEmpty(t *testing.T, name string, usageMessage string) string {
	t.Helper()

	value := os.Getenv(name)

	if value == "" {
		t.Skipf("skipping test; environment variable %s must be set. Usage: %s", name, usageMessage)
	}

	return value
}

// SkipIfAllEmpty verifies that at least one environment variable is non-empty or skips the test.
//
// If at least one environment variable is non-empty, returns the first name and value.
func SkipIfAllEmpty(t *testing.T, names []string, usageMessage string) (string, string) {
	t.Helper()

	name, value, err := RequireOneOf(names, usageMessage)
	if err != nil {
		t.Skipf("skipping test because %s.", err)
		return "", ""
	}

	return name, value
}
