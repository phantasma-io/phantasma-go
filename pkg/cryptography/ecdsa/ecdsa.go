package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/phantasma-io/phantasma-go/pkg/util"
	hash "github.com/phantasma-io/phantasma-go/pkg/util/hashing"
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

// Removes recovery ID from signature's byte array in 65-byte [R || S || V] format
func SignatureDropRecoveryId(signature []byte) []byte {

	if len(signature) != 65 {
		return util.ArrayClone(signature)
	}

	return signature[:len(signature)-1]
}

// Returns R/S pair
func SignatureToRS(signature []byte) (*big.Int, *big.Int) {
	signature = SignatureDropRecoveryId(signature)

	return big.NewInt(0).SetBytes(signature[:32]), big.NewInt(0).SetBytes(signature[32:])
}

func Sign(message, prikey []byte, curve ECDsaCurve) ([]byte, error) {
	// pk, err := crypto.HexToECDSA(privateKeyHex)
	// if err != nil {
	// 	panic(err)
	// }
	// pubKey := append(pk.PublicKey.X.Bytes(), pk.PublicKey.Y.Bytes()...)

	hash := hash.Sha256(message)

	if curve == Secp256k1 {
		signature, err := secp256k1.Sign(hash, prikey)
		if err != nil {
			return nil, err
		}

		return SignatureDropRecoveryId(signature), nil
	}

	return nil, errors.New("unsupported curve")
}

func Verify(message, signature, pubkey []byte, curve ECDsaCurve) bool {
	hash := hash.Sha256(message)
	if curve == Secp256k1 {
		return secp256k1.VerifySignature(UncompressedPublicKeyTo65Bytes(pubkey),
			hash,
			SignatureDropRecoveryId(signature))
	}
	if curve == Secp256r1 {

		pub := new(ecdsa.PublicKey)
		pub.Curve = elliptic.P256()
		pub.X, pub.Y = elliptic.UnmarshalCompressed(elliptic.P256(), CompressPublicKey(pubkey))

		r, s := SignatureToRS(signature)
		return ecdsa.Verify(pub, hash, r, s)
	}

	return false
}
