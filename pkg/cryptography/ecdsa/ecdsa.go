package ecdsa

import (
	"errors"

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

// Removes recovery ID from signature's byte array in 65-byte [R || S || V] format
func SignatureDropRecoveryId(signature []byte) []byte {

	if len(signature) != 65 {
		return util.ArrayClone(signature)
	}

	return signature[:len(signature)-1]
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

	return false
}
