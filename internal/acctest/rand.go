package acctest

import "math/rand"

const (
	// charSetAlphaNum is the alphanumeric character set for use with
	// RandStringFromCharSet.
	charSetAlphaNum = "abcdefghijklmnopqrstuvwxyz012346789"

	// Length of the resource name we wish to generate.
	resourceNameLength = 10
)

// GenerateRandomResourceName builds a unique-ish resource identifier to use in
// tests.
func GenerateRandomResourceName(prefix string) string {
	result := make([]byte, resourceNameLength)
	for i := 0; i < resourceNameLength; i++ {
		result[i] = charSetAlphaNum[randIntRange(0, len(charSetAlphaNum))]
	}
	return prefix + string(result)
}

// randIntRange returns a random integer between min (inclusive) and max
// (exclusive).
func randIntRange(min int, max int) int {
	return rand.Intn(max-min) + min //nolint:gosec
}
