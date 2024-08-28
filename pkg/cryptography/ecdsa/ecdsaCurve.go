package ecdsa

type ECDsaCurve uint

const (
	Secp256r1 ECDsaCurve = 0 // We use this for Neo signatures
	Secp256k1 ECDsaCurve = 1 // We use this for Eth/Bsc signatures
)
