package ecdsa

import (
	"math/big"

	"github.com/phantasma-io/phantasma-go/pkg/util"
)

// Removes recovery ID from signature's byte array in 65-byte [R || S || V] format
func SignatureDropRecoveryId(signature []byte) []byte {

	if len(signature) != 65 {
		return util.ArrayClone(signature)
	}

	return signature[:len(signature)-1]
}

func RSToSignatureWithoutRecoveryId(r, s *big.Int) []byte {
	return append(r.Bytes(), s.Bytes()...)
}

// Returns R/S pair
func SignatureToRS(signature []byte) (*big.Int, *big.Int) {
	signature = SignatureDropRecoveryId(signature)

	return big.NewInt(0).SetBytes(signature[:32]), big.NewInt(0).SetBytes(signature[32:])
}
