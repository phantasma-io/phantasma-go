package cryptography

import "github.com/phantasma.io/phantasma-go/pkg/io"

// SignatureKind type
type SignatureKind uint

const (
	// None Signature
	None SignatureKind = iota
	// Ed25519 Signature
	Ed25519
	// ECDSA Signature
	ECDSA
	// Ring Signature
	Ring
)

// Signature a
type Signature interface {
	Kind() SignatureKind
	Verify(message []byte, addresses []Address) bool
	Serialize(*io.BinWriter)
	Deserialize(*io.BinReader)
	Bytes() []byte
}
