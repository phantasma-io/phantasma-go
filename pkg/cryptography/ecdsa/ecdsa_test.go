package ecdsa

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMessage string = "test message"

// Eth address 0xDf738B927DA923fe0A5Fd3aD2192990C68913e6a
var pubKeyCompressed string = "025D3F7F469803C68C12B8F731576C74A9B5308484FD3B425D87C35CAED0A2E398"
var pubKey string = "5d3f7f469803c68c12b8f731576c74a9b5308484fd3b425d87c35caed0a2e398c7ac626d916a1d65e23f673a55e6b16ffc1abd673f3ef6ae8d5e6a0f99784a56"
var pubKey65 string = "045d3f7f469803c68c12b8f731576c74a9b5308484fd3b425d87c35caed0a2e398c7ac626d916a1d65e23f673a55e6b16ffc1abd673f3ef6ae8d5e6a0f99784a56"
var privKey string = "4ed773e5c8edc0487acef0011bc9ae8228287d4843f9d8477ff77c401ac59a49"

var signatureRef string = "55deb9e4d985834192ab8298c3dda18eb7082c2a744ebdf7233d0a93fb00a4a90b8af0b590c04c6d73d796f41c5d41abdbf57ecd795f3f40f3da92420b389376"
var signatureRefWithRID string = "55deb9e4d985834192ab8298c3dda18eb7082c2a744ebdf7233d0a93fb00a4a90b8af0b590c04c6d73d796f41c5d41abdbf57ecd795f3f40f3da92420b38937600"

func TestUncompressedPublicKeyTo65Bytes(t *testing.T) {
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		panic(err)
	}

	result := UncompressedPublicKeyTo65Bytes(pubKeyBytes)
	result2 := UncompressedPublicKeyTo65Bytes(result)

	// Method should not modify input array
	assert.Equal(t, pubKey, hex.EncodeToString(pubKeyBytes))
	assert.NotEqual(t, pubKey65, hex.EncodeToString(pubKeyBytes))

	assert.Equal(t, pubKey65, hex.EncodeToString(result))
	assert.Equal(t, result, result2)
}

func TestCompressedPublicKeyDecompression(t *testing.T) {
	pubKeyBytes, err := hex.DecodeString(pubKeyCompressed)
	if err != nil {
		panic(err)
	}
	pubKeyBytes, err = DecompressPublicKey(pubKeyBytes, Secp256k1)

	assert.Equal(t, pubKey65, hex.EncodeToString(pubKeyBytes))
}

func TestSignatureDropRecoveryId(t *testing.T) {
	signatureWithRIDBytes, err := hex.DecodeString(signatureRefWithRID)
	if err != nil {
		panic(err)
	}

	result := SignatureDropRecoveryId(signatureWithRIDBytes)
	result2 := SignatureDropRecoveryId(result)

	// Method should not modify input array
	assert.Equal(t, signatureRefWithRID, hex.EncodeToString(signatureWithRIDBytes))
	assert.NotEqual(t, signatureRef, hex.EncodeToString(signatureWithRIDBytes))

	assert.Equal(t, signatureRef, hex.EncodeToString(result))
	assert.Equal(t, result, result2)
}

func TestEcdsaSecp256k1Signing(t *testing.T) {
	privKeyBytes, err := hex.DecodeString(privKey)
	if err != nil {
		panic(err)
	}

	m := []byte(testMessage)

	signature, err := Sign(m, privKeyBytes, Secp256k1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Signature: %v\n", signature)
	fmt.Printf("Signature (HEX): %s\n", hex.EncodeToString(signature))

	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		panic(err)
	}
	verificationResult1 := Verify(m, signature, pubKeyBytes, Secp256k1)

	signatureRefBytes, err := hex.DecodeString(signatureRef)
	if err != nil {
		panic(err)
	}
	verificationResult2 := Verify([]byte(testMessage), signatureRefBytes, pubKeyBytes, Secp256k1)

	assert.Equal(t, hex.EncodeToString(signature), signatureRef) // Signatures should be the same, deterministic version is used
	assert.Equal(t, true, verificationResult1)
	assert.Equal(t, true, verificationResult2)
}
