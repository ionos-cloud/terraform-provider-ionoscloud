package hash

import (
	"crypto/md5"
)

// MD5 computes the MD5 checksum of the data.
func MD5(data []byte) (string, error) {
	hasher := md5.New() // nolint:gosec
	if _, err := hasher.Write(data); err != nil {
		return "", err
	}

	return string(hasher.Sum(nil)), nil
}
