package hash

import (
	"crypto/md5"
	"encoding/xml"
)

func MD5(data interface{}) (string, error) {
	// Marshal the struct to JSON
	jsonBytes, err := xml.Marshal(data)
	if err != nil {
		return "", err
	}

	// Create an MD5 hasher
	hasher := md5.New() // nolint:gosec

	// Write the JSON bytes to the hasher
	_, err = hasher.Write(jsonBytes)
	if err != nil {
		return "", err
	}

	// Compute the MD5 checksum
	md5sum := hasher.Sum(nil)
	return string(md5sum), nil
}
