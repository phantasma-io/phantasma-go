package cryptography

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashFrom(t *testing.T) {
	hash, _ := HashFromBytes([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19})
	assert.Equal(t, hash.String(), "4bcb1b15332489764c289b51b1119b1057c9e6dbca85b1ebc6553eae7c69e5e4")

	hash = HashFromString("asjdhweiurhwiuthedkgsdkfjh4otuiheriughdfjkgnsdçfjherslighjsghnoçiljhoçitujgpe8rotu89pearthkjdf.")
	assert.Equal(t, hash.String(), "9b93849b43a088f6d0add08f8ebfd4cd4ba8040515f281926c44954dbf65567d")
}

func TestHashIsNull(t *testing.T) {
	hash := Hash{}
	assert.Equal(t, true, hash.IsNull())
}

func TestHashBytes(t *testing.T) {
	hash, _ := HashFromBytes([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19})
	result := []byte{228, 229, 105, 124, 174, 62, 85, 198, 235, 177, 133, 202, 219, 230, 201, 87, 16, 155, 17, 177, 81, 155, 40, 76, 118, 137, 36, 51, 21, 27, 203, 75}
	assert.Equal(t, result, hash.Bytes())
}

func TestHashParse(t *testing.T) {
	hash, _ := ParseHash("e4e5697cae3e55c6ebb185cadbe6c957109b11b1519b284c76892433151bcb4b")
	assert.Equal(t, "e4e5697cae3e55c6ebb185cadbe6c957109b11b1519b284c76892433151bcb4b", hash.String())
}
