package hash

import (
	"crypto/sha256"
)

// Sha256 hashes the incoming byte slice
// using the sha256 algorithm.
func Sha256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// DoubleSha256 performs sha256 twice on the given data.
func DoubleSha256(data []byte) []byte {

	h1 := Sha256(data)
	hash := Sha256(h1)
	return hash
}

// Checksum returns the checksum for a given piece of data
// using sha256 twice as the hash algorithm.
func Checksum(data []byte) []byte {
	hash := DoubleSha256(data)
	return hash[:4]
}
