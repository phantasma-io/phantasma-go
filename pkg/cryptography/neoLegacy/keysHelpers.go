package neoLegacy

import (
	"github.com/nspcc-dev/neo-go/pkg/crypto/hash"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/phantasma-io/phantasma-go/pkg/encoding/base58"
)

// Code below was taken from https://github.com/nspcc-dev/neo-go v0.78.4

// GetVerificationScript returns NEO VM bytecode with CHECKSIG command for the
// public key.
func getVerificationScript(pubKey []byte) []byte {
	b := pubKey
	b = append([]byte{byte(0x21)}, b...) // PUSHBYTES33
	b = append(b, byte(0xAC))            // CHECKSIG

	return b
}

// GetScriptHash returns a Hash160 of verification script for the key.
func getScriptHash(pubKey []byte) util.Uint160 {
	return hash.Hash160(getVerificationScript(pubKey))
}

// Prefix is the byte used to prepend to addresses when encoding them, it can
// be changed and defaults to 23 (0x17), the standard NEO prefix.
var Prefix = byte(0x17)

// Uint160ToString returns the "NEO address" from the given Uint160.
func uint160ToString(u util.Uint160) string {
	// Dont forget to prepend the Address version 0x17 (23) A
	b := append([]byte{Prefix}, u.BytesBE()...)
	return base58.CheckEncode(b)
}

// Address returns a base58-encoded NEO-specific address based on the key hash.
func Address(pubKey []byte) string {
	return uint160ToString(getScriptHash(pubKey))
}
