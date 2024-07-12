package blockchain

import (
	"testing"

	"github.com/phantasma-io/phantasma-go/pkg/cryptography"
	"github.com/phantasma-io/phantasma-go/pkg/io"
	"github.com/stretchr/testify/assert"
)

func TestNewTx(t *testing.T) {
	tx := NewTransaction("mainnet", "main", []byte{0x01, 0x02, 0x03}, 1623519055, []byte("PAYLOAD"))
	assert.Equal(t, "mainnet", tx.NexusName)
	assert.Equal(t, "main", tx.ChainName)
	assert.Equal(t, tx.Hash.String(), "6cb5c800485e64d9e14c2dede12ac3c845c0bbe8a2fd420f46bc28ca905710b6")
}

func TestTxNilPayload(t *testing.T) {
	tx := NewTransaction("mainnet", "main", []byte{0x01, 0x02, 0x03}, 1623519055, nil)
	assert.Equal(t, "mainnet", tx.NexusName)
	assert.Equal(t, "main", tx.ChainName)
	assert.Equal(t, tx.Hash.String(), "b049d3bf5449191d3eb8ea3ea9cdace3712775d509c96ff8743266298e4b077a")
}

func TestTxSign(t *testing.T) {
	tx := NewTransaction("mainnet", "main", []byte{0x01, 0x02, 0x03}, 1623519055, nil)
	assert.Equal(t, "mainnet", tx.NexusName)
	assert.Equal(t, "main", tx.ChainName)
	assert.Equal(t, tx.Hash.String(), "b049d3bf5449191d3eb8ea3ea9cdace3712775d509c96ff8743266298e4b077a")

	kp := cryptography.NewPhantasmaKeys([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x30, 0x31, 0x32})
	tx.Sign(kp)

	assert.True(t, len(tx.Signatures) == 1)
	sigBytes := tx.Signatures[0].Bytes()
	assert.Equal(t, []byte{64, 201, 196, 27, 63, 0, 116, 45, 220, 91, 240, 237, 165, 112, 220, 44, 20, 88, 78, 168, 227, 90, 130, 235, 142, 236, 165, 47, 201, 182, 151, 99, 73, 179, 146, 235, 153, 1, 150, 254, 184, 113, 69, 112, 26, 150, 155, 41, 7, 53, 43, 215, 220, 2, 153, 210, 35, 239, 148, 145, 177, 88, 250, 174, 7}, sigBytes)
}

func TestTxSerialization(t *testing.T) {
	tx := NewTransaction("mainnet", "main", []byte{0x01, 0x02, 0x03}, 1623519055, nil)
	assert.Equal(t, "mainnet", tx.NexusName)
	assert.Equal(t, "main", tx.ChainName)
	assert.Equal(t, tx.Hash.String(), "b049d3bf5449191d3eb8ea3ea9cdace3712775d509c96ff8743266298e4b077a")
	assert.Equal(t, []byte{7, 109, 97, 105, 110, 110, 101, 116, 4, 109, 97, 105, 110, 3, 1, 2, 3, 79, 239, 196, 96, 0}, tx.Bytes(false))

	bw := *io.NewBufBinWriter()
	tx.Serialize(bw.BinWriter, false)
	bytes := bw.Bytes()

	newTx := Transaction{}
	br := *io.NewBinReaderFromBuf(bytes)
	newTx.Deserialize(&br)
	assert.Equal(t, tx.ChainName, newTx.ChainName)
	assert.Equal(t, tx.NexusName, newTx.NexusName)
	assert.Equal(t, tx.Hash, newTx.Hash)
	assert.Equal(t, tx.Script, newTx.Script)
	assert.Equal(t, tx.Signatures, newTx.Signatures)
	assert.Equal(t, tx.Payload, newTx.Payload)

	// just to be sure
	assert.Equal(t, tx, newTx)
}

//TODO
//func TestTxMine(t *testing.T) {}
