package proofOfAddress

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/dustinxie/ecc"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/phantasma-io/phantasma-go/pkg/cryptography"
	"github.com/phantasma-io/phantasma-go/pkg/cryptography/ecdsa"
	"github.com/phantasma-io/phantasma-go/pkg/cryptography/neoLegacy"
)

type ProofOfAddressesVerifier struct {
	Message            string
	SignedMessage      string
	SignedMessageBytes []byte

	PhaAddress        string
	PhaPublicKeyBytes []byte

	EthAddress        string
	EthPublicKey      string
	EthPublicKeyBytes []byte

	Neo2Address        string
	Neo2PublicKey      string
	Neo2PublicKeyBytes []byte

	PhaSignature       string
	PhaSignatureBytes  []byte
	EthSignature       string
	EthSignatureBytes  []byte
	Neo2Signature      string
	Neo2SignatureBytes []byte
}

func NewProofOfAddressesVerifier(message string) *ProofOfAddressesVerifier {
	v := &ProofOfAddressesVerifier{}
	v.Message = message

	split := strings.Split(strings.Replace(v.Message, "\r", "", -1), "\n")

	v.SignedMessage = strings.Join(split[:6], "\n")
	fmt.Printf("v.SignedMessage: %s\n", v.SignedMessage)
	v.SignedMessageBytes = []byte(v.SignedMessage)

	v.PhaAddress = split[1][19:]
	fmt.Printf("v.PhaAddress: %s\n", v.PhaAddress)
	phaAddress, err := cryptography.FromString(v.PhaAddress)
	if err != nil {
		panic(err)
	}
	v.PhaPublicKeyBytes = phaAddress.GetPublicKey()

	v.EthAddress = split[2][18:]
	fmt.Printf("v.EthAddress: %s\n", v.EthAddress)
	v.EthPublicKey = split[3][21:]
	fmt.Printf("v.EthPublicKey: %s\n", v.EthPublicKey)
	v.EthPublicKeyBytes, err = hex.DecodeString(v.EthPublicKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("v.EthPublicKey (HEX): 0x%s\n", hex.EncodeToString(v.EthPublicKeyBytes))

	v.Neo2Address = split[4][20:]
	fmt.Printf("v.Neo2Address: %s\n", v.Neo2Address)
	v.Neo2PublicKey = split[5][23:]
	fmt.Printf("v.Neo2PublicKey: %s\n", v.Neo2PublicKey)
	v.Neo2PublicKeyBytes, err = hex.DecodeString(v.Neo2PublicKey)
	if err != nil {
		panic(err)
	}
	v.PhaSignature = split[7][21:]
	fmt.Printf("v.PhaSignature: %s\n", v.PhaSignature)
	v.PhaSignatureBytes, err = hex.DecodeString(v.PhaSignature)
	if err != nil {
		panic(err)
	}
	v.EthSignature = split[8][20:]
	fmt.Printf("v.EthSignature: %s\n", v.EthSignature)
	v.EthSignatureBytes, err = hex.DecodeString(v.EthSignature)
	if err != nil {
		panic(err)
	}
	fmt.Printf("v.EthSignatureBytes len: %d\n", len(v.EthSignatureBytes))
	fmt.Printf("v.EthSignature (HEX): 0x%s\n", hex.EncodeToString(v.EthSignatureBytes))

	v.Neo2Signature = split[9][22:]
	fmt.Printf("v.Neo2Signature: %s\n", v.Neo2Signature)
	v.Neo2SignatureBytes, err = hex.DecodeString(v.Neo2Signature)
	if err != nil {
		panic(err)
	}

	return v
}

func (v *ProofOfAddressesVerifier) VerifyMessage() (bool, string) {
	success := true
	errorMessage := ""

	if !ed25519.Verify(v.PhaPublicKeyBytes, v.SignedMessageBytes, v.PhaSignatureBytes) {
		success = false
		errorMessage += "Phantasma signature is incorrect!\n"
	}

	fmt.Println(len(v.EthSignatureBytes))
	ethRes, err := ecdsa.Verify(v.SignedMessageBytes, v.EthSignatureBytes, v.EthPublicKeyBytes, ecdsa.Secp256k1)
	if err != nil {
		panic(err)
	}
	if !ethRes {
		success = false
		errorMessage += "Ethereum signature is incorrect!\n"
	}

	fmt.Println(len(v.Neo2SignatureBytes))
	neoRes, err := ecdsa.Verify(v.SignedMessageBytes, v.Neo2SignatureBytes, v.Neo2PublicKeyBytes, ecdsa.Secp256r1)
	if err != nil {
		panic(err)
	}
	if !neoRes {
		success = false
		errorMessage += "Neo Legacy signature is incorrect!\n"
	}

	pubEth := ecdsa.PublicKeyUnmarshal(v.EthPublicKeyBytes, ecc.P256k1())
	ethAddressFromPublicKey := crypto.PubkeyToAddress(*pubEth).Hex()

	if v.EthAddress != ethAddressFromPublicKey {
		success = false
		errorMessage += "Ethereum address is incorrect: " + ethAddressFromPublicKey + "\n"
	}

	neo2AddressFromPublicKey := neoLegacy.Address([]byte(v.Neo2PublicKeyBytes))
	if v.Neo2Address != neo2AddressFromPublicKey {
		success = false
		errorMessage += "Neo Legacy address is incorrect: " + neo2AddressFromPublicKey + "\n"
	}

	return success, errorMessage
}
