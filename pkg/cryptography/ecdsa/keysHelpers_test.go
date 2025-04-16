package ecdsa

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUncompressedPublicKeyTo65Bytes(t *testing.T) {
	pubKeyBytes, err := hex.DecodeString(k1PubKey)
	if err != nil {
		panic(err)
	}

	result := UncompressedPublicKeyTo65Bytes(pubKeyBytes)
	result2 := UncompressedPublicKeyTo65Bytes(result)

	// Method should not modify input array
	assert.Equal(t, k1PubKey, hex.EncodeToString(pubKeyBytes))
	assert.NotEqual(t, k1PubKey65, hex.EncodeToString(pubKeyBytes))

	assert.Equal(t, k1PubKey65, hex.EncodeToString(result))
	assert.Equal(t, result, result2)
}
