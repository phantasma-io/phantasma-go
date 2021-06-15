package cryptography

import (
	"bytes"
	"encoding/hex"
	"strings"

	"github.com/phantasma.io/phantasma-go/pkg/io"
	"github.com/phantasma.io/phantasma-go/pkg/util"
	hashing "github.com/phantasma.io/phantasma-go/pkg/util/hashing"
)

// TODO Hash.Null

// HashLength const
const HashLength = 32

// HexPrefix const
const HexPrefix = "0x"

// Hash struct
type Hash struct {

	// Code to run in PhantasmaVM for this transaction.
	_data []byte
}

// HashFromBytes returns a Hash based on the passed in byte slice
func HashFromBytes(data []byte) (Hash, error) {

	if len(data) != HashLength {
		data = hashing.Sha256(data)
	}

	return Hash{_data: data}, nil
}

// HashFromString creates an instance of Hash from a string
func HashFromString(s string) Hash {
	data := hashing.Sha256([]byte(s))
	return Hash{data}
}

// ParseHash parses a string resulting in a Hash
func ParseHash(s string) (Hash, error) {

	if strings.HasPrefix(s, HexPrefix) {
		s = s[2:]
	}

	if len(s) < 64 || len(s) == 0 {
		panic("string has wrong format")
	}

	bytes, err := hex.DecodeString(s)
	if err != nil {
		return Hash{}, err
	}

	bytes = util.ArrayReverse(bytes)
	return Hash{bytes}, nil
}

// Size returns the length of the underlying byte slice
func (h Hash) Size() int {
	return len(h._data)
}

// IsNull checks if the Hash represents a nil hash
func (h Hash) IsNull() bool {
	if h._data == nil {
		return true
	}

	empty := make([]byte, HashLength)
	if !bytes.Equal(h._data, empty) {
		return false
	}

	return true
}

// String creates the a base16 encoded representation of Hash
func (h Hash) String() string {
	data := util.ArrayReverse(h._data)
	return hex.EncodeToString(data)
}

// Bytes returns the bytes of the hash
func (h Hash) Bytes() []byte {
	return h._data
}

// FromUnpaddedHex creates an instance of Hash from an unpadded hex string
func (h Hash) FromUnpaddedHex(s string) (Hash, error) {

	if strings.HasPrefix(s, HexPrefix) {
		s = s[2:]
	}

	var sb strings.Builder
	sb.WriteString(s)

	for sb.Len() < 64 {
		sb.WriteString("0")
		sb.WriteString("0")
	}

	return ParseHash(sb.String())
}

// GetDifficulty retrieves the current difficulty of the hash
func (h Hash) GetDifficulty() int {
	var result int = 0
	for i := 0; i < len(h._data); i++ {
		var n = h._data[i]

		for j := 0; j < 8; j++ {
			if (n & (1 << j)) != 0 {
				result = 1 + (i << 3) + j
			}
		}
	}

	return 256 - result
}

// Serialize implements ther Serializable interface
func (h Hash) Serialize(writer io.BinWriter) {
	writer.WriteBytes(h._data)
}

// Deserialize implements ther Serializable interface
func (h Hash) Deserialize(reader io.BinReader) {
	reader.ReadBytes(h._data)
}
