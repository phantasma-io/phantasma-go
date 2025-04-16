package util

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func printBytes(message string, bytes []byte) {
	fmt.Printf("%-25s %v [", message, bytes)
	for i, b := range bytes {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf("%#x", b)
	}

	fmt.Printf("]\n")
}

func BigIntTestConversions(t *testing.T, d TestData) {
	bi, _ := big.NewInt(0).SetString(d.number, 10)
	fmt.Printf("\n%s\n", bi)

	// Print original bytes of GO's big.Int
	printBytes("n bytes [go]        ", bi.Bytes())
	// Print bytes in little-endian
	printBytes("n bytes [go, le]    ", ArrayCloneAndReverse(bi.Bytes()))

	// Convert big.Int into Phantasma's array (signed bytes, LE)
	phaBytes := BigIntToPhantasmaByteArray(bi)

	// Print resulting Phantasma's array
	printBytes("n bytes [pha]               ", phaBytes)
	// Printing 2 reference arrays
	printBytes("n bytes [pha/refval]        ", d.phaBytes)
	printBytes("n bytes [c#/refval]         ", d.csharpBytes)

	// Convert Phantasma's array back to GO's format (unsigned bytes, BI)
	phaBytesConvertedBack, sign := BigIntBytesFromCsharpOrPhantasmaByteArray(phaBytes)
	printBytes("n bytes [pha/converted back] ", phaBytesConvertedBack)

	// Convert reference C#'s array to GO's format (unsigned bytes, BI)
	csharpBytesConvertedBack, _ := BigIntBytesFromCsharpOrPhantasmaByteArray(d.csharpBytes)
	printBytes("n bytes [c#/converted back]  ", csharpBytesConvertedBack)

	// Take original unmodified big.Int, convert it to c# representation and compare with reference array
	csharpBytes := BigIntToCsharpByteArray(bi)
	printBytes("n bytes [c#]                 ", csharpBytes)

	// Making big.Int from Phantasma's array converted back to GO's array
	biPha := big.NewInt(0).SetBytes(phaBytesConvertedBack)
	if sign < 0 {
		biPha = big.NewInt(0).Mul(big.NewInt(-1), biPha)
	}
	fmt.Println("n restored [pha]: ", biPha.String())

	// Making big.Int from C#'s array converted back to GO's array
	biCsharp := big.NewInt(0).SetBytes(csharpBytesConvertedBack)
	if sign < 0 {
		biCsharp = big.NewInt(0).Mul(big.NewInt(-1), biCsharp)
	}
	fmt.Println("n restored [c#]: ", biCsharp.String())

	// Ensure we encode here in the same way as Phantasma.Core in C# do
	require.Equal(t, d.phaBytes, phaBytes)

	// Ensure our C# method encodes in the same way as C# itself
	require.Equal(t, d.csharpBytes, csharpBytes)

	// Comparing reference number with numbers recreated from byte arrays
	require.Equal(t, d.number, biPha.String())
	require.Equal(t, d.number, biCsharp.String())
}

func TestBigInteger(t *testing.T) {
	for _, a := range testData {
		BigIntTestConversions(t, a)
	}
}
