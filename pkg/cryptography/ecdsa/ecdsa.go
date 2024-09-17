package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	hash "github.com/phantasma-io/phantasma-go/pkg/util/hashing"
)

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
	} else if curve == Secp256r1 {
		pk := PrivateKeyUnmarshal(prikey, elliptic.P256())

		r, s, err := ecdsa.Sign(rand.Reader, pk, hash)
		if err != nil {
			return nil, err
		}

		signature := RSToSignatureWithoutRecoveryId(r, s)

		return signature, nil
	}

	return nil, errors.New("unsupported curve")
}

func Verify(message, signature, pubkey []byte, curve ECDsaCurve) (bool, error) {
	hash := hash.Sha256(message)
	if curve == Secp256k1 {

		var uncompressedPubkey []byte
		if len(pubkey) > 33 {
			uncompressedPubkey = UncompressedPublicKeyTo65Bytes(pubkey)
		} else {
			var err error
			uncompressedPubkey, err = DecompressPublicKey(pubkey, Secp256k1)

			if err != nil {
				return false, err
			}
		}

		return secp256k1.VerifySignature(uncompressedPubkey,
			hash,
			SignatureDropRecoveryId(signature)), nil
	}
	if curve == Secp256r1 {
		pub := PublicKeyUnmarshal(pubkey, elliptic.P256())

		r, s := SignatureToRS(signature)
		return ecdsa.Verify(pub, hash, r, s), nil
	}

	return false, nil
}
