package cryptography

import (
	"crypto/ed25519"

	"github.com/phantasma.io/phantasma-go/pkg/io"
)

// Ed25519Signature struct
type Ed25519Signature struct {
	bytes []byte
	kind  SignatureKind
}

// NewEd25519Signature instatiates a new signature object
func NewEd25519Signature(bytes []byte) Ed25519Signature {
	return Ed25519Signature{bytes, Ed25519}
}

// Generate a new signature
func Generate(keyPair KeyPair, message []byte) Ed25519Signature {
	return NewEd25519Signature(ed25519.Sign(keyPair.ExpandedPrivateKey(), message))
}

// Kind returns the type of the signature
func (sig Ed25519Signature) Kind() SignatureKind {
	return sig.kind
}

// Verify verifies that the message was signed by one of the addresses passed in
func (sig Ed25519Signature) Verify(message []byte, addresses []Address) bool {

	for _, address := range addresses {
		if !address.IsUser() {
			continue
		}

		pubKey := address.Bytes()[2:]

		if ed25519.Verify(pubKey, message, sig.bytes) {
			return true
		}
	}

	return false
}

// Bytes returns the byte representation of the signature
func (sig Ed25519Signature) Bytes() []byte {
	bw := *io.NewBufBinWriter()
	sig.Serialize(bw.BinWriter)
	return bw.Bytes()
}

// Serialize implements ther Serializable interface
func (sig Ed25519Signature) Serialize(writer *io.BinWriter) {
	writer.WriteVarBytes(sig.bytes)
}

// Deserialize implements ther Serializable interface
func (sig Ed25519Signature) Deserialize(reader *io.BinReader) {
	reader.ReadVarBytes()
}
