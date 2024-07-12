package cryptography

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"strconv"

	"github.com/phantasma-io/phantasma-go/pkg/encoding/base58"
)

// KeyPair a
type KeyPair interface {
	PrivateKey() []byte
	ExpandedPrivateKey() []byte
	PublicKey() []byte
	Address() Address

	Sign(msg []byte) Signature
}

// PrivateKeyLength defines the length of a private key
const PrivateKeyLength = 32

// PhantasmaKeys is the struct that holds the information about keys for the phantasma blockchain
type PhantasmaKeys struct {
	privateKey         []byte
	expandedPrivateKey []byte
	publicKey          []byte
	address            Address
}

// NewPhantasmaKeys instantiates a new PhantasmaKeys object based on the given seed
func NewPhantasmaKeys(seed []byte) PhantasmaKeys {

	if len(seed) != PrivateKeyLength {
		panic("Length of private key has not been met, needs to be " + strconv.Itoa(PrivateKeyLength))
	}

	keys := PhantasmaKeys{}
	pk := ed25519.NewKeyFromSeed(seed)
	keys.publicKey = pk[32:]
	keys.privateKey = pk[:32]
	keys.expandedPrivateKey = pk[:]
	keys.address = FromKey(keys)

	return keys
}

// GeneratePhantasmaKeys creates a new phantasma keypair
func GeneratePhantasmaKeys() PhantasmaKeys {

	seed := make([]byte, PrivateKeyLength)
	rand.Read(seed)
	keys := NewPhantasmaKeys(seed)
	return keys
}

// FromWIF creates a new key pair based on the passed in WIF
func FromWIF(wif string) (PhantasmaKeys, error) {

	if len(wif) == 0 {
		return PhantasmaKeys{}, fmt.Errorf("WIF needs to be set")
	}

	data, err := base58.CheckDecode(wif)
	if err != nil {
		return PhantasmaKeys{}, err
	}

	privateKey := make([]byte, PrivateKeyLength)
	copy(privateKey[0:], data[1:33])

	keys := NewPhantasmaKeys(privateKey)
	return keys, nil
}

func (k PhantasmaKeys) String() string {
	return k.address.String()
}

// Sign generates a signature for the passed in message
func (k PhantasmaKeys) Sign(msg []byte) Signature {
	return Generate(k, msg)
}

// WIF returns the WIF based on the private key
func (k PhantasmaKeys) WIF() string {
	bytes := make([]byte, PrivateKeyLength+2)
	bytes[0] = 0x80

	copy(bytes[1:], k.PrivateKey()[0:32])
	bytes[33] = 0x01
	encoded := base58.CheckEncode(bytes)

	return string(encoded)
}

//ExpandedPrivateKey returns the associated expanded private key
func (k PhantasmaKeys) ExpandedPrivateKey() []byte {
	return k.expandedPrivateKey
}

//PrivateKey returns the associated private key
func (k PhantasmaKeys) PrivateKey() []byte {
	return k.privateKey
}

//PublicKey returns the associated public key
func (k PhantasmaKeys) PublicKey() []byte {
	return k.publicKey
}

//Address returns the associated address
func (k PhantasmaKeys) Address() Address {
	return k.address
}
