package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"

	"github.com/dustinxie/ecc"
	"github.com/phantasma-io/phantasma-go/pkg/util"
)

// Adds extra byte prefix to signify that this is an uncompressed key
func UncompressedPublicKeyTo65Bytes(pubkey []byte) []byte {
	if len(pubkey) == 65 {
		return util.ArrayClone(pubkey)
	}
	pubkey = append([]byte{0x04}, pubkey...)
	return pubkey
}

func CompressPublicKey(uncompressedPublicKey []byte) []byte {
	x := big.NewInt(0).SetBytes(uncompressedPublicKey[:32])
	y := big.NewInt(0).SetBytes(uncompressedPublicKey[32:])

	var prefix byte = 0x02
	_, m := new(big.Int).DivMod(y, big.NewInt(2), new(big.Int))
	if m.Cmp(big.NewInt(0)) != 0 {
		prefix = 0x03
	}

	return append([]byte{prefix}, x.Bytes()...)
}

func PrivateKeyUnmarshal(privKey []byte, curve elliptic.Curve) *ecdsa.PrivateKey {
	pk := new(ecdsa.PrivateKey)
	pk.Curve = curve
	pk.D = new(big.Int).SetBytes(privKey)

	return pk
}

func PublicKeyUnmarshal(pubKey []byte, curve elliptic.Curve) *ecdsa.PublicKey {
	pub := new(ecdsa.PublicKey)
	pub.Curve = curve

	if len(pubKey) > 33 {
		if curve == ecc.P256k1() {
			pub.X, pub.Y = ecc.UnmarshalCompressed(ecc.P256k1(), CompressPublicKey(pubKey))
		} else {
			pub.X, pub.Y = elliptic.UnmarshalCompressed(elliptic.P256(), CompressPublicKey(pubKey))
		}

	} else {
		if curve == ecc.P256k1() {
			pub.X, pub.Y = ecc.UnmarshalCompressed(ecc.P256k1(), pubKey)
		} else {
			pub.X, pub.Y = elliptic.UnmarshalCompressed(elliptic.P256(), pubKey)
		}
	}

	return pub
}
