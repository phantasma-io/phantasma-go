package cryptography

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/phantasma-io/phantasma-go/pkg/encoding/base58"
	"github.com/phantasma-io/phantasma-go/pkg/io"
)

// Length is the length of data
const Length = 34

// AddressKind type
type AddressKind byte

const (
	// Invalid address
	Invalid AddressKind = 0x00
	// User address
	User AddressKind = 0x01
	// System address
	System AddressKind = 0x02
	// Interop address
	Interop AddressKind = 0x03
)

// Address struct
type Address struct {

	// Code to run in PhantasmaVM for this transaction.
	data []byte
	Text string
	kind AddressKind
}

var Null []byte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0}

// NewAddress returns a new Address object
func NewAddress(pubKey []byte) Address {

	if pubKey == nil {
		panic("Passed public key can't be nil!")
	}

	if len(pubKey) != Length {
		panic("publicKey length must be " + strconv.Itoa(Length) + " but length is " + strconv.Itoa(len(pubKey)))
	}

	address := Address{}
	address.data = pubKey

	var prefix string
	switch kind := address.Kind(); kind {
	case User:
		prefix = "P"
	case Interop:
		prefix = "X"
	default:
		prefix = "S"
	}

	address.Text = prefix + base58.Encode(address.data)

	return address
}

// NullAddress returns a new null Address object
func NullAddress() Address {
	return NewAddress(Null)
}

// FromString creates an instance of an Address from a string
func FromString(s string) (Address, error) {

	data, err := base58.Decode(s[1:])
	if err != nil {
		return Address{}, err
	}

	address := NewAddress(data)

	switch prefix := s[:1]; prefix {
	case "P":
		if address.Kind() != User {
			return Address{}, fmt.Errorf("Address has to be of type User")
		}
	case "X":
		if address.Kind() != Interop {
			return Address{}, fmt.Errorf("Address has to be of type Interop")
		}
	case "S":
		if address.Kind() != System {
			return Address{}, fmt.Errorf("Address has to be of type System")
		}
	default:
		return Address{}, fmt.Errorf("Unknown address prefix: " + prefix)
	}

	return address, nil
}

// IsValidAddress verifies if a string is a valid address
func IsValidAddress(text string) bool {
	_, err := FromString(text)
	if err != nil {
		return false
	}

	return true
}

// IsUser verifies if the passed in address is a user address
func (a Address) IsUser() bool {
	return a.Kind() == User
}

// FromKey generates an address from a KeyPair
func FromKey(keyPair KeyPair) Address {
	data := make([]byte, Length)
	data[0] = byte(User)

	if len(keyPair.PublicKey()) == 32 {

		copy(data[2:], keyPair.PublicKey()[0:32])

	} else if len(keyPair.PublicKey()) == 33 {

		copy(data[1:], keyPair.PublicKey()[0:33])

	} else {
		panic("Invalid public key length")
	}

	return NewAddress(data)
}

// IsNull checks if the Address represents a nil Address
func (a Address) IsNull() bool {
	if a.data == nil {
		return true
	}

	empty := make([]byte, Length)
	if !bytes.Equal(a.data, empty) {
		return false
	}

	return true
}

// Kind returns the kind of an address
func (a Address) Kind() AddressKind {
	if a.IsNull() {
		a.kind = System
		return a.kind
	}

	if a.data[0] >= 3 {
		a.kind = Interop
	} else if a.data[0] == 2 {
		a.kind = System
	} else {
		a.kind = User
	}

	return a.kind
}

// String creates the a base58 encoded representation of the address including the address prefix
func (a Address) String() string {
	return a.Text
}

// Bytes returns the data representing an address
func (a Address) Bytes() []byte {
	return a.data
}

// BytesPrefixed returns the data representing an address, including prefix required for binary serialization
func (a Address) BytesPrefixed() []byte {
	return bytes.Join([][]byte{{34}, a.data}, []byte{})
}

// Serialize implements ther Serializable interface
func (a Address) Serialize(writer io.BinWriter) {
	writer.WriteVarBytes(a.data)
}

// Deserialize implements ther Serializable interface
func (a Address) Deserialize(reader io.BinReader) {
	a.data = reader.ReadVarBytes()
}
