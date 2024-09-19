package proofOfAddress

import (
	"encoding/hex"
	"fmt"
	"testing"

	hash "github.com/phantasma-io/phantasma-go/pkg/util/hashing"
	"github.com/stretchr/testify/assert"
)

var testMessage string = `I have signed this message with my Phantasma, Ethereum and Neo Legacy signatures to prove that following addresses belong to me and were derived from private key that belongs to me and to confirm my willingness to swap funds across these addresses upon my request. My public addresses are:
Phantasma address: P2KHhbVZWDv1ZLLoJccN3PUAb9x9BqRnUyH3ZEhu5YwBeJQ
Ethereum address: 0xDf738B927DA923fe0A5Fd3aD2192990C68913e6a
Ethereum public key: 5D3F7F469803C68C12B8F731576C74A9B5308484FD3B425D87C35CAED0A2E398C7AC626D916A1D65E23F673A55E6B16FFC1ABD673F3EF6AE8D5E6A0F99784A56
Neo Legacy address: Ae3aEA6CpvckvypAUShj2CLsy7sfynKUzj
Neo Legacy public key: 183A301779007BF42DD7B5247587585B0524E13989F964C2A8E289A0CDC91F001765FCC3B4CEE5ED274C4A8B6D80978BDFED678210458CE264D4A4DAB3923EE6

Phantasma signature: EEC4FFE0FE71C522E32A99C3C84F111944275D967BA74BA06749CB58EEC48A2E450B0ABCC59D86880AC01DA2AECFEF1FC5C45BFD840761FC420AF8194ACD5B01
Ethereum signature: E3E1FCD85385675F9E3508630570C545DECCD1241C7A8FFF523D2AC500D6F68745E43975DBF871C99504100B8DD6715F036FA51EFF9EB8B79D1E31FD555E78FC
Neo Legacy signature: 8A961AB366DFD9A3EB1ED5FD496E1CE818A67873162A7A21F663BA320324A8B075EAFA432926B7A894FB3DC5138BF0564A10D1AAC1A4B4E8DB8184CD69F829AC`

var testMessagePhaAddressIncorrect string = `I have signed this message with my Phantasma, Ethereum and Neo Legacy signatures to prove that following addresses belong to me and were derived from private key that belongs to me and to confirm my willingness to swap funds across these addresses upon my request. My public addresses are:
Phantasma address: P2KHhbVZWDv1ZLLoJccN3PUAb9x9BqRnUyH3ZEhu5YwBeJW
Ethereum address: 0xDf738B927DA923fe0A5Fd3aD2192990C68913e6a
Ethereum public key: 5D3F7F469803C68C12B8F731576C74A9B5308484FD3B425D87C35CAED0A2E398C7AC626D916A1D65E23F673A55E6B16FFC1ABD673F3EF6AE8D5E6A0F99784A56
Neo Legacy address: Ae3aEA6CpvckvypAUShj2CLsy7sfynKUzj
Neo Legacy public key: 183A301779007BF42DD7B5247587585B0524E13989F964C2A8E289A0CDC91F001765FCC3B4CEE5ED274C4A8B6D80978BDFED678210458CE264D4A4DAB3923EE6

Phantasma signature: EEC4FFE0FE71C522E32A99C3C84F111944275D967BA74BA06749CB58EEC48A2E450B0ABCC59D86880AC01DA2AECFEF1FC5C45BFD840761FC420AF8194ACD5B01
Ethereum signature: E3E1FCD85385675F9E3508630570C545DECCD1241C7A8FFF523D2AC500D6F68745E43975DBF871C99504100B8DD6715F036FA51EFF9EB8B79D1E31FD555E78FC
Neo Legacy signature: 8A961AB366DFD9A3EB1ED5FD496E1CE818A67873162A7A21F663BA320324A8B075EAFA432926B7A894FB3DC5138BF0564A10D1AAC1A4B4E8DB8184CD69F829AC`

var testMessagePhaSignatureIncorrect string = `I have signed this message with my Phantasma, Ethereum and Neo Legacy signatures to prove that following addresses belong to me and were derived from private key that belongs to me and to confirm my willingness to swap funds across these addresses upon my request. My public addresses are:
Phantasma address: P2KHhbVZWDv1ZLLoJccN3PUAb9x9BqRnUyH3ZEhu5YwBeJQ
Ethereum address: 0xDf738B927DA923fe0A5Fd3aD2192990C68913e6a
Ethereum public key: 5D3F7F469803C68C12B8F731576C74A9B5308484FD3B425D87C35CAED0A2E398C7AC626D916A1D65E23F673A55E6B16FFC1ABD673F3EF6AE8D5E6A0F99784A56
Neo Legacy address: Ae3aEA6CpvckvypAUShj2CLsy7sfynKUzj
Neo Legacy public key: 183A301779007BF42DD7B5247587585B0524E13989F964C2A8E289A0CDC91F001765FCC3B4CEE5ED274C4A8B6D80978BDFED678210458CE264D4A4DAB3923EE6

Phantasma signature: EEC5FFE0FE71C522E32A99C3C84F111944275D967BA74BA06749CB58EEC48A2E450B0ABCC59D86880AC01DA2AECFEF1FC5C45BFD840761FC420AF8194ACD5B01
Ethereum signature: E3E1FCD85385675F9E3508630570C545DECCD1241C7A8FFF523D2AC500D6F68745E43975DBF871C99504100B8DD6715F036FA51EFF9EB8B79D1E31FD555E78FC
Neo Legacy signature: 8A961AB366DFD9A3EB1ED5FD496E1CE818A67873162A7A21F663BA320324A8B075EAFA432926B7A894FB3DC5138BF0564A10D1AAC1A4B4E8DB8184CD69F829AC`

func TestProofOfAddressesVerifier(t *testing.T) {
	v := NewProofOfAddressesVerifier(testMessage)

	hash := hash.Sha256([]byte(v.SignedMessage))
	fmt.Printf("Hash: %s\n", hex.EncodeToString(hash))

	success, errorMessage := v.VerifyMessage()
	assert.Equal(t, true, success)
	assert.Equal(t, "", errorMessage)

	v = NewProofOfAddressesVerifier(testMessagePhaAddressIncorrect)
	success, errorMessage = v.VerifyMessage()
	assert.Equal(t, false, success)
	assert.Equal(t, "Phantasma signature is incorrect!\nEthereum signature is incorrect!\nNeo Legacy signature is incorrect!\n", errorMessage)

	v = NewProofOfAddressesVerifier(testMessagePhaSignatureIncorrect)
	success, errorMessage = v.VerifyMessage()
	assert.Equal(t, false, success)
	assert.Equal(t, "Phantasma signature is incorrect!\n", errorMessage)
}
