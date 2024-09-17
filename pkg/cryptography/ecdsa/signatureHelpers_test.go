package ecdsa

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureDropRecoveryId(t *testing.T) {
	signatureWithRIDBytes, err := hex.DecodeString(k1SignatureRefWithRID)
	if err != nil {
		panic(err)
	}

	result := SignatureDropRecoveryId(signatureWithRIDBytes)
	result2 := SignatureDropRecoveryId(result)

	// Method should not modify input array
	assert.Equal(t, k1SignatureRefWithRID, hex.EncodeToString(signatureWithRIDBytes))
	assert.NotEqual(t, k1SignatureRef, hex.EncodeToString(signatureWithRIDBytes))

	assert.Equal(t, k1SignatureRef, hex.EncodeToString(result))
	assert.Equal(t, result, result2)
}

func TestSignatureToRSConversions(t *testing.T) {
	{
		signatureBytes, err := hex.DecodeString(k1SignatureRef)
		if err != nil {
			panic(err)
		}

		r, s := SignatureToRS(signatureBytes)
		signatureBytesRecreated := RSToSignatureWithoutRecoveryId(r, s)
		assert.Equal(t, signatureBytes, signatureBytesRecreated)
	}

	{
		signatureBytes, err := hex.DecodeString(r1SignatureRef)
		if err != nil {
			panic(err)
		}

		r, s := SignatureToRS(signatureBytes)
		signatureBytesRecreated := RSToSignatureWithoutRecoveryId(r, s)
		assert.Equal(t, signatureBytes, signatureBytesRecreated)
	}
}
