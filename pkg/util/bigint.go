package util

import (
	"math/big"
	"slices"
)

func BigIntToCsharpByteArray(n *big.Int) []byte {
	if n.BitLen() == 0 { // Check if big int is zero
		return []byte{0x00}
	}

	var bytes = n.Bytes()

	if n.Sign() == -1 {
		TwosComplementConvertTo(bytes)
	}

	slices.Reverse(bytes) // Converting to little-endian format

	if n.Sign() == -1 {
		if bytes[len(bytes)-1] < 128 {
			bytes = append(bytes, 0xff)
		}
	} else {
		if bytes[len(bytes)-1] >= 128 && bytes[len(bytes)-1] != 0x00 {
			bytes = append(bytes, 0x00)
		}
	}

	return bytes
}

// That's a replication of C# ToSignedByteArray() from BigIntegerExtension.cs
// Adding Phantasma's additional postfixes
func BigIntToPhantasmaByteArray(n *big.Int) []byte {
	var bytes = BigIntToCsharpByteArray(n)

	if n.Sign() == -1 { // Big int is negative
		if len(bytes) == 1 {
			bytes = append(bytes, 0xff, 0xff)
		} else if len(bytes) > 1 && bytes[len(bytes)-1] == 0xff {
			bytes = append(bytes, 0xff)
		}
	} else {
		if bytes[len(bytes)-1] != 0x00 {
			bytes = append(bytes, 0x00)
		}
	}

	return bytes
}

func BigIntBytesFromCsharpOrPhantasmaByteArray(bytes []byte) ([]byte, int) {
	n := make([]byte, len(bytes))
	copy(n, bytes)

	if len(n) == 1 && n[0] == 0x00 {
		return []byte{}, 0 // To keep compatibility with BinWriter/Reader tests, doesn't change the value of big int
	}

	if n[len(n)-1] < 128 {
		// It's a positive number.

		slices.Reverse(n)
		return n, 1
	}

	// It's a negative number.

	slices.Reverse(n)
	TwosComplementConvertFrom(n)

	return n, -1
}

func BigIntFromCsharpOrPhantasmaByteArray(bytes []byte) *big.Int {
	b, sign := BigIntBytesFromCsharpOrPhantasmaByteArray(bytes)

	n := big.NewInt(0)
	if sign < 0 {
		n = n.Mul(n.SetBytes(b), big.NewInt(-1))
	} else {
		n = n.SetBytes(b)
	}

	return n
}
