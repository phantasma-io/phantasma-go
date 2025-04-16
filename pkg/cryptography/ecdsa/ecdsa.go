package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"

	"github.com/dustinxie/ecc"
	hash "github.com/phantasma-io/phantasma-go/pkg/util/hashing"
)

func Sign(message, prikey []byte, curve ECDsaCurve) ([]byte, error) {
	if len(message) == 0 {
		return nil, errors.New("message lenth is 0")
	}
	if len(prikey) == 0 {
		return nil, errors.New("prikey lenth is 0")
	}

	hash := hash.Sha256(message)

	if curve == Secp256k1 {
		pk := PrivateKeyUnmarshal(prikey, ecc.P256k1())

		signature, err := ecc.SignBytes(pk, hash, ecc.LowerS)
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
	if len(message) == 0 {
		return false, errors.New("message lenth is 0")
	}
	if len(signature) == 0 {
		return false, errors.New("signature lenth is 0")
	}
	if len(pubkey) == 0 {
		return false, errors.New("pubkey lenth is 0")
	}

	hash := hash.Sha256(message)
	if curve == Secp256k1 {
		pub := PublicKeyUnmarshal(pubkey, ecc.P256k1())

		return ecc.VerifyBytes(pub, hash, SignatureDropRecoveryId(signature), ecc.Normal), nil
	}
	if curve == Secp256r1 {
		pub := PublicKeyUnmarshal(pubkey, elliptic.P256())

		r, s := SignatureToRS(signature)
		return ecdsa.Verify(pub, hash, r, s), nil
	}

	return false, nil
}
