package configlog

import (
	"errors"
	"fmt"
)

// ErrDefaultConfigUnavailable signals that the default config file could not be used,
// either because its path could not be determined (e.g. the home directory is unknown)
// or because the file does not exist. It is non-fatal on its own: callers should
// proceed with other credential sources and only fail when none are available.
var ErrDefaultConfigUnavailable = errors.New("default config file unavailable")

// defaultConfigUnavailableError marks that the default config file could not be used and
// carries a human-readable reason. It unwraps to ErrDefaultConfigUnavailable so callers
// can detect it with errors.Is, while its message holds just the reason for display.
type defaultConfigUnavailableError struct {
	reason string
}

func (e *defaultConfigUnavailableError) Error() string { return e.reason }
func (e *defaultConfigUnavailableError) Unwrap() error { return ErrDefaultConfigUnavailable }

// Summary and detail for the "no credentials" error, shared by both providers so the
// wording stays identical.
const (
	MissingCredentialsSummary = "No credentials found"
	MissingCredentialsDetail  = "Set credentials in one of these ways:\n" +
		"  - token: the `token` provider argument or IONOS_TOKEN env var\n" +
		"  - username/password: the `username`/`password` provider arguments or IONOS_USERNAME/IONOS_PASSWORD env vars\n" +
		"  - file config (check README.md for more details)"
)

// MissingCredentialsHint returns extra context to append to a "no credentials" error
// when the default config file could not be used, so the real cause is not masked. It
// returns an empty string for any other error.
func MissingCredentialsHint(readFileErr error) string {
	if errors.Is(readFileErr, ErrDefaultConfigUnavailable) {
		return fmt.Sprintf("\nThe default config file was not loaded: %s.", readFileErr.Error())
	}
	return ""
}
