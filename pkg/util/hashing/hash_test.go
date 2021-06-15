package hash

import (
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSha256(t *testing.T) {
	input := []byte("hello")
	data := Sha256(input)

	expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	actual := hex.EncodeToString(data)

	assert.Equal(t, expected, actual)
}

func TestHashDoubleSha256(t *testing.T) {
	input := []byte("hello")
	data := DoubleSha256(input)

	firstSha := Sha256(input)
	doubleSha := Sha256(firstSha)
	expected := hex.EncodeToString(doubleSha)

	actual := hex.EncodeToString(data)
	assert.Equal(t, expected, actual)
}

func TestChecksum(t *testing.T) {
	testCases := []struct {
		data []byte
		sum  uint32
	}{
		{nil, 0xe2e0f65d},
		{[]byte{}, 0xe2e0f65d},
		{[]byte{1, 2, 3, 4}, 0xe272e48d},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.sum, binary.LittleEndian.Uint32(Checksum(tc.data)))
	}
}
