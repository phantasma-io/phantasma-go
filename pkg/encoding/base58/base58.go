package base58

import (
	"bytes"
	"errors"

	"github.com/mr-tron/base58"
	hash "github.com/phantasma-io/phantasma-go/pkg/util/hashing"
)

// CheckDecode implements a base58-encoded string decoding with hash-based
// checksum check.
func CheckDecode(s string) (b []byte, err error) {
	b, err = base58.Decode(s)
	if err != nil {
		return nil, err
	}

	if len(b) < 5 {
		return nil, errors.New("invalid base-58 check string: missing checksum")
	}

	if !bytes.Equal(hash.Checksum(b[:len(b)-4]), b[len(b)-4:]) {
		return nil, errors.New("invalid base-58 check string: invalid checksum")
	}

	// Strip the 4 byte long hash.
	b = b[:len(b)-4]

	return b, nil
}

// CheckEncode encodes given byte slice into a base58 string with hash-based
// checksum appended to it.
func CheckEncode(b []byte) string {
	b = append(b, hash.Checksum(b)...)

	return base58.Encode(b)
}

// Encode encodes given byte slice into a base58 string
func Encode(b []byte) string {
	return base58.Encode(b)
}

// Decode implements a base58-encoded string
func Decode(s string) (b []byte, err error) {
	b, err = base58.Decode(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}
