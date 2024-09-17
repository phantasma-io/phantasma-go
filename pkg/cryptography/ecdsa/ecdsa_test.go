package ecdsa

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMessage string = "test message"

// Eth address 0xDf738B927DA923fe0A5Fd3aD2192990C68913e6a
var k1PubKeyCompressed string = "025D3F7F469803C68C12B8F731576C74A9B5308484FD3B425D87C35CAED0A2E398"
var k1PubKey string = "5d3f7f469803c68c12b8f731576c74a9b5308484fd3b425d87c35caed0a2e398c7ac626d916a1d65e23f673a55e6b16ffc1abd673f3ef6ae8d5e6a0f99784a56"
var k1PubKey65 string = "045d3f7f469803c68c12b8f731576c74a9b5308484fd3b425d87c35caed0a2e398c7ac626d916a1d65e23f673a55e6b16ffc1abd673f3ef6ae8d5e6a0f99784a56"

var r1PubKeyCompressed string = "02183A301779007BF42DD7B5247587585B0524E13989F964C2A8E289A0CDC91F00"
var r1PubKey string = "183A301779007BF42DD7B5247587585B0524E13989F964C2A8E289A0CDC91F001765FCC3B4CEE5ED274C4A8B6D80978BDFED678210458CE264D4A4DAB3923EE6"
var r1PubKey65 string = "04183A301779007BF42DD7B5247587585B0524E13989F964C2A8E289A0CDC91F001765FCC3B4CEE5ED274C4A8B6D80978BDFED678210458CE264D4A4DAB3923EE6"

var privKey string = "4ed773e5c8edc0487acef0011bc9ae8228287d4843f9d8477ff77c401ac59a49"

var k1SignatureRef string = "55deb9e4d985834192ab8298c3dda18eb7082c2a744ebdf7233d0a93fb00a4a90b8af0b590c04c6d73d796f41c5d41abdbf57ecd795f3f40f3da92420b389376"
var k1SignatureRefWithRID string = "55deb9e4d985834192ab8298c3dda18eb7082c2a744ebdf7233d0a93fb00a4a90b8af0b590c04c6d73d796f41c5d41abdbf57ecd795f3f40f3da92420b38937600"

var r1SignatureRef string = "7D375DEEB56530A8E09BB3F4AF9217F922FD3D33EBF02874239A2910E9DEF1BD25119CA641F13C6EBED1BFF4FEB7834F56723F9A9DCFC80B3128F1028B2C3A6B"

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

func TestCompressedPublicKeyDecompression(t *testing.T) {
	pubKeyBytes, err := hex.DecodeString(k1PubKeyCompressed)
	if err != nil {
		panic(err)
	}
	pubKeyBytes, err = DecompressPublicKey(pubKeyBytes, Secp256k1)

	assert.Equal(t, k1PubKey65, hex.EncodeToString(pubKeyBytes))
}

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

	pubKeyBytesCompressed, err := hex.DecodeString(k1PubKeyCompressed)
	if err != nil {
		panic(err)
	}
	pubKeyBytes, err := hex.DecodeString(k1PubKey)
	if err != nil {
		panic(err)
	}
	verificationResult1, err := Verify(m, signature, pubKeyBytes, Secp256k1)
	if err != nil {
		panic(err)
	}

	signatureRefBytes, err := hex.DecodeString(k1SignatureRef)
	if err != nil {
		panic(err)
	}
	verificationResult2, err := Verify([]byte(testMessage), signatureRefBytes, pubKeyBytes, Secp256k1)
	if err != nil {
		panic(err)
	}

	verificationResult3, err := Verify([]byte(testMessage), signatureRefBytes, pubKeyBytesCompressed, Secp256k1)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, hex.EncodeToString(signature), k1SignatureRef) // Signatures should be the same, deterministic version is used
	assert.Equal(t, true, verificationResult1)
	assert.Equal(t, true, verificationResult2)
	assert.Equal(t, true, verificationResult3)
}

func TestEcdsaSecp256r1Signing(t *testing.T) {
	privKeyBytes, err := hex.DecodeString(privKey)
	if err != nil {
		panic(err)
	}

	m := []byte(testMessage)

	signature, err := Sign(m, privKeyBytes, Secp256r1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Signature: %v\n", signature)
	fmt.Printf("Signature (HEX): %s\n", hex.EncodeToString(signature))

	pubKeyBytesCompressed, err := hex.DecodeString(r1PubKeyCompressed)
	if err != nil {
		panic(err)
	}
	pubKeyBytes, err := hex.DecodeString(r1PubKey)
	if err != nil {
		panic(err)
	}
	verificationResult1, err := Verify(m, signature, pubKeyBytes, Secp256r1)
	if err != nil {
		panic(err)
	}

	signatureRefBytes, err := hex.DecodeString(r1SignatureRef)
	if err != nil {
		panic(err)
	}
	verificationResult2, err := Verify([]byte(testMessage), signatureRefBytes, pubKeyBytes, Secp256r1)
	if err != nil {
		panic(err)
	}

	verificationResult3, err := Verify([]byte(testMessage), signatureRefBytes, pubKeyBytesCompressed, Secp256r1)
	if err != nil {
		panic(err)
	}

	assert.NotEqual(t, hex.EncodeToString(signature), r1SignatureRef) // Signatures should not be the same, non-deterministic version is used
	assert.Equal(t, true, verificationResult1)
	assert.Equal(t, true, verificationResult2)
	assert.Equal(t, true, verificationResult3)
}
